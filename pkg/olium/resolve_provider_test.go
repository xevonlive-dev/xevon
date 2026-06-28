package olium

import (
	"strings"
	"testing"

	"github.com/xevonlive-dev/xevon/internal/config"
)

// TestResolveProviderOpenAICompatibleModel guards the openai-compatible model
// resolution: an empty agent.olium.model must fall back to
// custom_provider.model_id, an explicit model must win over it, and both empty
// must error with an actionable message. This is the regression test for the
// default-collision bug where DefaultOliumConfig shipped Model="gemma4:latest",
// which shadowed custom_provider.model_id (model was neither "" nor the
// DefaultModel sentinel, so the fallback never fired).
func TestResolveProviderOpenAICompatibleModel(t *testing.T) {
	const baseURL = "http://localhost:11434/v1"

	cases := []struct {
		name        string
		model       string // opts.Model (agent.olium.model / --model)
		customModel string // opts.CustomModelID (custom_provider.model_id)
		wantModel   string
		wantErr     string // substring; empty = expect nil
	}{
		{
			name:        "empty model falls back to custom_provider.model_id",
			model:       "",
			customModel: "qwen3.6:latest",
			wantModel:   "qwen3.6:latest",
		},
		{
			name:        "DefaultModel sentinel falls back to custom_provider.model_id",
			model:       DefaultModel,
			customModel: "qwen3.6:latest",
			wantModel:   "qwen3.6:latest",
		},
		{
			name:        "explicit model wins over custom_provider.model_id",
			model:       "llama3.3:70b",
			customModel: "qwen3.6:latest",
			wantModel:   "llama3.3:70b",
		},
		{
			name:        "no model anywhere errors",
			model:       "",
			customModel: "",
			wantErr:     "model is required",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, providerName, gotModel, err := resolveProvider(Options{
				Provider:      "openai-compatible",
				Model:         c.model,
				CustomBaseURL: baseURL,
				CustomModelID: c.customModel,
			})
			if c.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil (model=%q)", c.wantErr, gotModel)
				}
				if !strings.Contains(err.Error(), c.wantErr) {
					t.Fatalf("error %q does not contain %q", err, c.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if providerName != "openai-compatible" {
				t.Errorf("provider name = %q, want openai-compatible", providerName)
			}
			if gotModel != c.wantModel {
				t.Errorf("resolved model = %q, want %q", gotModel, c.wantModel)
			}
		})
	}
}

// TestDefaultConfigDoesNotShadowCustomModelID is the end-to-end guard for the
// original bug report: with the shipped default config, a user who changes ONLY
// custom_provider.model_id (and never touches agent.olium.model) must have that
// model honored. It starts from config.DefaultOliumConfig() — the real shipped
// default — and wires it into olium.Options exactly as pkg/cli/agent_olium.go
// does when no CLI flags are passed.
//
// This catches the root cause that TestResolveProviderOpenAICompatibleModel
// cannot: if someone re-introduces a non-empty Model default in
// DefaultOliumConfig() (e.g. Model: "gemma4:latest"), that value would shadow
// the distinct model_id here and the resolved model would come back wrong.
func TestDefaultConfigDoesNotShadowCustomModelID(t *testing.T) {
	// Shipped default config, then the user's scenario: change only model_id.
	cfg := config.DefaultOliumConfig()
	cfg.CustomProvider.ModelID = "qwen3.6:latest"

	// Sanity-check the precondition that makes the test meaningful: the default
	// model must NOT equal the model_id we set, otherwise the assertion below
	// would pass even if Model were shadowing model_id.
	if cfg.Model == cfg.CustomProvider.ModelID {
		t.Fatalf("test precondition broken: default Model %q must differ from the model_id under test", cfg.Model)
	}

	// Mirror pkg/cli/agent_olium.go Options wiring with no CLI flags set, i.e.
	// firstNonEmptyString("", cfg.X) collapses to cfg.X.
	opts := Options{
		Provider:      cfg.Provider,
		Model:         cfg.Model, // no --model flag
		CustomBaseURL: cfg.CustomProvider.BaseURL,
		CustomModelID: cfg.CustomProvider.ModelID,
		CustomAPIKey:  cfg.CustomProvider.APIKey,
	}

	_, providerName, gotModel, err := resolveProvider(opts)
	if err != nil {
		t.Fatalf("unexpected error resolving default openai-compatible config: %v", err)
	}
	if providerName != "openai-compatible" {
		t.Fatalf("default provider = %q, want openai-compatible", providerName)
	}
	if gotModel != "qwen3.6:latest" {
		t.Errorf("resolved model = %q, want %q — agent.olium.model is shadowing custom_provider.model_id (check DefaultOliumConfig().Model is empty)", gotModel, "qwen3.6:latest")
	}
}
