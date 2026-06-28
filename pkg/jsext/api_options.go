package jsext

import (
	"github.com/xevonlive-dev/xevon/internal/config"
	"github.com/xevonlive-dev/xevon/pkg/agent/llm"
	"github.com/xevonlive-dev/xevon/pkg/database"
	"github.com/xevonlive-dev/xevon/pkg/oast"
)

// APIOptions lives in pkg/jsext/api; it is re-exported here via an alias in
// api_aliases.go.

// EngineOptions provides scanner context to the JS engine.
// These come from the runner at engine creation time.
type EngineOptions struct {
	ScopeMatcher *config.ScopeMatcher
	ScopeConfig  *config.ScopeConfig
	ScanUUID     string
	Repository   *database.Repository

	// LLMClient enables xevon.agent.* API (nil = disabled)
	LLMClient llm.Client

	// OASTService enables xevon.oast.* API (nil = disabled)
	OASTService *oast.Service
}
