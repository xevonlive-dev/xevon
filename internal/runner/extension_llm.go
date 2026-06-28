package runner

import (
	"github.com/xevonlive-dev/xevon/internal/config"
	"github.com/xevonlive-dev/xevon/pkg/agent/llm"
	"go.uber.org/zap"
)

// extensionLLMClient builds the olium-backed LLM client that powers the
// xevon.agent.* JS extension API (xevon.agent.complete/ask/chat/...).
// It resolves the provider from agent.olium config.
//
// Returns nil when olium isn't configured (missing credentials, unknown
// provider). That's intentionally non-fatal: extensions that never touch the
// agent API are unaffected, and ones that do surface a clear nil-client error
// instead of aborting scan startup.
func extensionLLMClient(settings *config.Settings) llm.Client {
	if settings == nil {
		return nil
	}
	client, err := llm.NewOliumClient(&settings.Agent.Olium)
	if err != nil {
		zap.L().Debug("xevon.agent JS extension API disabled: olium provider unavailable",
			zap.Error(err))
		return nil
	}
	return client
}
