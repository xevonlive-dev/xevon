package server

import (
	"context"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/xevonlive-dev/xevon/internal/config"
	"github.com/xevonlive-dev/xevon/internal/runner"
	"github.com/xevonlive-dev/xevon/pkg/agent"
	"github.com/xevonlive-dev/xevon/pkg/core/services"
	"github.com/xevonlive-dev/xevon/pkg/database"
	"github.com/xevonlive-dev/xevon/pkg/http"
	"github.com/xevonlive-dev/xevon/pkg/olium/tool"
	"github.com/xevonlive-dev/xevon/pkg/piolium"
	"github.com/xevonlive-dev/xevon/pkg/queue"
	"go.uber.org/zap"
)

// countCache caches expensive COUNT(*) query results with a TTL.
type countCache struct {
	mu        sync.Mutex
	records   int64
	findings  int64
	updatedAt time.Time
	ttl       time.Duration
}

func newCountCache(ttl time.Duration) *countCache {
	return &countCache{ttl: ttl}
}

// Get returns cached counts, refreshing from the database if expired.
func (cc *countCache) Get(db *database.DB) (records, findings int64) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	if time.Since(cc.updatedAt) < cc.ttl {
		return cc.records, cc.findings
	}

	// Shared cross-request cache refresh — not owned by any single request, so it
	// uses a background context rather than one caller's request context.
	ctx := context.Background()
	if recordCount, err := db.NewSelect().Model((*database.HTTPRecord)(nil)).Count(ctx); err == nil {
		cc.records = int64(recordCount)
	}
	if findingCount, err := db.NewSelect().Model((*database.Finding)(nil)).Count(ctx); err == nil {
		cc.findings = int64(findingCount)
	}
	cc.updatedAt = time.Now()

	return cc.records, cc.findings
}

// scanState tracks a running scan for a specific project.
type scanState struct {
	running bool
	runner  *runner.Runner
	scanID  string
}

// queuedScan represents a scan waiting in a per-project queue.
type queuedScan struct {
	scanID        string
	runner        *runner.Runner
	projectUUID   string
	enqueued      time.Time
	uploadResults bool
}

// Handlers holds the HTTP handlers and their dependencies.
type Handlers struct {
	queue         queue.Queue
	db            *database.DB
	repo          *database.Repository
	recordWriter  *database.RecordWriter
	config        ServerConfig
	settings      *config.Settings
	httpRequester *http.Requester
	services      *services.Services // shared with httpRequester; may be nil
	configWatcher *config.ConfigWatcher
	startTime     time.Time

	// Domain handler groups extracted from this struct. Composed here and
	// wired in NewHandlers; the routes call e.g. h.findings.HandleListFindings.
	findings *findingsHandlers

	// Per-project scan state for API-triggered scans
	scanMu     sync.Mutex
	scanStates map[string]*scanState // keyed by projectUUID

	// Scan queue: when ScanQueueCapacity > 0, scans are queued instead of rejected
	scanQueues map[string]chan *queuedScan

	// Cached scope matcher (lazy-initialized, invalidated on config change)
	scopeMatcherMu sync.RWMutex
	scopeMatcher   *config.ScopeMatcher

	// Prometheus metrics handler
	metricsHandler fiber.Handler

	// Long-lived agent engine (shared across requests for warm session reuse)
	agentEngine *agent.Engine

	// Agent run state for API-triggered agent runs
	agentMu           sync.Mutex
	agentHeavySem     chan struct{} // counting semaphore for heavy runs (autopilot/swarm)
	agentLightSem     chan struct{} // counting semaphore for light runs (query/chat)
	agenticScanStatus map[string]*AgenticScanStatusResponse
	// agenticCancels holds the cancel func for each in-flight agentic run so it
	// can be stopped externally (stop/delete). Guarded by agentMu.
	agenticCancels map[string]context.CancelFunc

	// projectHeavyMu guards projectHeavyActive. The map tallies currently-
	// running heavy agent runs per project so a single tenant can be
	// capped (AgentHeavyPerProject) below the cluster-wide limit
	// (AgentHeavyMax). Entries are deleted when the count drops to 0
	// so an idle project doesn't keep an empty map row.
	projectHeavyMu     sync.Mutex
	projectHeavyActive map[string]int

	// Background cleanup for completed agent run statuses
	agentCleanupStop chan struct{}

	// Cached COUNT query results for server-info endpoint
	counts *countCache

	// Cached piolium availability — `pi` install state doesn't change
	// mid-process, so the per-request audit-harness picker reuses this
	// instead of re-probing PATH and reading ~/.pi/agent/settings.json.
	pioliumAvailableOnce sync.Once
	pioliumAvailable     bool
}

// pioliumAvailableCached wraps piolium.IsAvailable with a sync.Once so the
// PATH lookup + settings.json read happens at most once per server process.
func (h *Handlers) pioliumAvailableCached() bool {
	h.pioliumAvailableOnce.Do(func() {
		h.pioliumAvailable = piolium.IsAvailable()
	})
	return h.pioliumAvailable
}

// NewHandlers creates a new Handlers instance.
// Starts a background goroutine to clean up old agent run records from the database.
func NewHandlers(q queue.Queue, db *database.DB, repo *database.Repository, rw *database.RecordWriter, cfg ServerConfig, settings *config.Settings, httpRequester *http.Requester, svc *services.Services) *Handlers {
	heavyMax := cfg.AgentHeavyMax
	if heavyMax <= 0 {
		heavyMax = 5
	}
	lightMax := cfg.AgentLightMax
	if lightMax <= 0 {
		lightMax = 10
	}

	h := &Handlers{
		queue:              q,
		db:                 db,
		repo:               repo,
		recordWriter:       rw,
		config:             cfg,
		settings:           settings,
		httpRequester:      httpRequester,
		services:           svc,
		startTime:          time.Now(),
		scanStates:         make(map[string]*scanState),
		scanQueues:         make(map[string]chan *queuedScan),
		agentHeavySem:      make(chan struct{}, heavyMax),
		agentLightSem:      make(chan struct{}, lightMax),
		agenticScanStatus:  make(map[string]*AgenticScanStatusResponse),
		agenticCancels:     make(map[string]context.CancelFunc),
		projectHeavyActive: make(map[string]int),
		agentCleanupStop:   make(chan struct{}),
		counts:             newCountCache(10 * time.Second),
	}
	h.findings = &findingsHandlers{db: db, repo: repo}

	// Reconcile zombie runs orphaned by a previous process exit. A scan or agent
	// run left in a non-terminal state has no goroutine behind it after a
	// restart, so mark it failed instead of leaving it "running" forever (which
	// also blocks it from being deleted and shows a NaN duration in the UI).
	if repo != nil {
		rctx := context.Background()
		if n, err := repo.FailRunningScans(rctx, "interrupted by server restart"); err != nil {
			zap.L().Warn("startup: failed to reconcile stale scans", zap.Error(err))
		} else if n > 0 {
			zap.L().Info("startup: marked orphaned scans as failed", zap.Int("count", n))
		}
		if n, err := repo.FailRunningAgenticScans(rctx, "interrupted by server restart"); err != nil {
			zap.L().Warn("startup: failed to reconcile stale agent runs", zap.Error(err))
		} else if n > 0 {
			zap.L().Info("startup: marked orphaned agent runs as failed", zap.Int("count", n))
		}
	}

	// Initialize the agent-browser shell guard from config at startup so a
	// server booted with agent.browser.enable=false refuses agent-browser
	// invocations immediately, before any agent run starts.
	if settings != nil {
		tool.SetBrowserBlocked(!settings.Agent.Browser.IsEnabled())
	}

	if !cfg.NoAgent {
		h.agentEngine = agent.NewEngine(settings, repo)
		go h.agentDBCleanupLoop()
	}
	return h
}

// agentDBCleanupLoop periodically removes old completed/failed agent runs from the
// database and prunes the in-memory map for completed runs.
func (h *Handlers) agentDBCleanupLoop() {
	const ttl = 24 * time.Hour
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-h.agentCleanupStop:
			return
		case <-ticker.C:
			// Prune in-memory map (completed runs older than 1h)
			now := time.Now()
			h.agentMu.Lock()
			for id, status := range h.agenticScanStatus {
				if status.CompletedAt != nil && now.Sub(*status.CompletedAt) > time.Hour {
					delete(h.agenticScanStatus, id)
				}
			}
			h.agentMu.Unlock()

			// Prune DB (completed/failed runs older than 24h)
			if h.repo != nil {
				if n, err := h.repo.DeleteOldAgenticScans(context.Background(), ttl); err == nil && n > 0 {
					zap.L().Debug("Cleaned up old agent runs", zap.Int("count", n))
				}
			}

			// Clean up old agent session directories
			if h.settings != nil {
				sessDir := h.settings.Agent.EffectiveSessionsDir()
				if n, err := agent.CleanupSessionDirs(sessDir, 48*time.Hour); err == nil && n > 0 {
					zap.L().Debug("Cleaned up old session directories", zap.Int("count", n))
				}
			}
		}
	}
}

// HandleHealth handles GET /health.
func (h *Handlers) HandleHealth(c fiber.Ctx) error {
	return c.JSON(HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

// HandleAppInfo handles GET /api/info — returns basic app info.
func (h *Handlers) HandleAppInfo(c fiber.Ctx) error {
	commit := h.config.Commit
	if len(commit) > 7 {
		commit = commit[:7]
	}
	return c.JSON(AppInfoResponse{
		Name:        "xevon",
		Version:     h.config.Version,
		Author:      h.config.Author,
		Docs:        "https://docs.xevon.live",
		LicenseSPDX: "AGPL-3.0-or-later",
		Source:      "https://github.com/xevonlive-dev/xevon",
		BuildTime:   h.config.BuildTime,
		Commit:      commit,
	})
}

// HandleServerInfo handles GET /server-info.
func (h *Handlers) HandleServerInfo(c fiber.Ctx) error {
	metrics := h.queue.Metrics()
	if metrics == nil {
		metrics = &queue.QueueMetrics{}
	}

	commit := h.config.Commit
	if len(commit) > 7 {
		commit = commit[:7]
	}

	resp := ServerInfoResponse{
		Name:        "xevon",
		Version:     h.config.Version,
		Author:      h.config.Author,
		Docs:        "https://docs.xevon.live",
		BuildTime:   h.config.BuildTime,
		Commit:      commit,
		Uptime:      time.Since(h.startTime).Round(time.Second).String(),
		ServiceAddr: h.config.ServiceAddr,
		ProxyAddr:   h.config.IngestProxyAddr,
		QueueDepth:  metrics.Depth,
		License:     h.config.License,
		LicenseSPDX: "AGPL-3.0-or-later",
		Source:      "https://github.com/xevonlive-dev/xevon",
		DemoOnly:    h.config.DemoOnly,
		ViewOnly:    h.config.ViewOnly,
	}

	if h.db != nil {
		resp.TotalRecords, resp.TotalFindings = h.counts.Get(h.db)
	}

	return c.JSON(resp)
}

// getScopeMatcher returns the cached ScopeMatcher, creating it lazily on first call.
func (h *Handlers) getScopeMatcher() *config.ScopeMatcher {
	h.scopeMatcherMu.RLock()
	m := h.scopeMatcher
	h.scopeMatcherMu.RUnlock()
	if m != nil {
		return m
	}

	// Double-check under write lock
	h.scopeMatcherMu.Lock()
	defer h.scopeMatcherMu.Unlock()
	if h.scopeMatcher != nil {
		return h.scopeMatcher
	}
	if h.settings != nil {
		h.scopeMatcher = config.NewScopeMatcher(h.settings.Scope)
	}
	return h.scopeMatcher
}

// resetScopeMatcher invalidates the cached ScopeMatcher so it is rebuilt on next use.
func (h *Handlers) resetScopeMatcher() {
	h.scopeMatcherMu.Lock()
	h.scopeMatcher = nil
	h.scopeMatcherMu.Unlock()
}

// IsScanRunning reports whether any scan is currently running (across all projects).
// Implements metrics.ScanStateProvider.
func (h *Handlers) IsScanRunning() bool {
	h.scanMu.Lock()
	defer h.scanMu.Unlock()
	for _, st := range h.scanStates {
		if st.running {
			return true
		}
	}
	return false
}

// getProjectScanState returns the scan state for a project, creating it if needed.
// Must be called with scanMu held.
func (h *Handlers) getProjectScanState(projectUUID string) *scanState {
	st, ok := h.scanStates[projectUUID]
	if !ok {
		st = &scanState{}
		h.scanStates[projectUUID] = st
	}
	return st
}

// scanQueueWorker processes queued scans for a project sequentially.
func (h *Handlers) scanQueueWorker(projectUUID string, ch chan *queuedScan) {
	for qs := range ch {
		h.scanMu.Lock()
		st := h.getProjectScanState(qs.projectUUID)
		st.running = true
		st.runner = qs.runner
		st.scanID = qs.scanID
		h.scanMu.Unlock()

		h.runBackgroundScan(qs.scanID, qs.runner, qs.projectUUID, qs.uploadResults)
	}
}

// Close releases handler resources including the agent engine pool.
func (h *Handlers) Close() {
	close(h.agentCleanupStop)
	h.scanMu.Lock()
	for _, ch := range h.scanQueues {
		close(ch)
	}
	h.scanQueues = make(map[string]chan *queuedScan)
	h.scanMu.Unlock()
	// Agent engine is in-process (olium); no resources to release.
}

// HandleMetrics handles GET /metrics — serves Prometheus metrics.
func (h *Handlers) HandleMetrics(c fiber.Ctx) error {
	if h.metricsHandler != nil {
		return h.metricsHandler(c)
	}
	return c.Status(fiber.StatusNotFound).SendString("metrics not configured")
}
