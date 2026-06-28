package api

import (
	"github.com/xevonlive-dev/xevon/internal/config"
	"github.com/xevonlive-dev/xevon/pkg/agent/llm"
	"github.com/xevonlive-dev/xevon/pkg/database"
	"github.com/xevonlive-dev/xevon/pkg/http"
	"github.com/xevonlive-dev/xevon/pkg/oast"
	"github.com/xevonlive-dev/xevon/pkg/output"
)

// APIOptions holds all context needed to set up the xevon.* JS API on a VM.
type APIOptions struct {
	ScriptID   string
	HTTPClient *http.Requester
	ConfigVars map[string]string

	// Scanner context (all nil-safe)
	ScopeMatcher *config.ScopeMatcher
	ScopeConfig  *config.ScopeConfig
	ScanUUID     string
	ProjectUUID  string

	// Security controls (from extensions config)
	AllowExec   bool   // gate for exec() and setEnv(); default false
	SandboxDir  string // base path for file ops; empty = cwd
	ExecTimeout int    // max seconds for exec(); default 30, cap 120

	// Finding emitter for hooks that want to create findings
	FindingEmitter func(*output.ResultEvent)

	// Database repository for ingest API (nil = ingest disabled)
	Repository *database.Repository

	// LLMClient enables xevon.agent.* API (nil = disabled)
	LLMClient llm.Client

	// OASTService enables xevon.oast.* API (nil = disabled)
	OASTService *oast.Service
}
