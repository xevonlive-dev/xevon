package httpmsg

import "testing"

// reset restores the process-global UA state so each subtest is independent.
func reset() {
	uaMu.Lock()
	uaOverride = ""
	buildVersion = ""
	uaMu.Unlock()
}

func TestDefaultUserAgent_BuiltinWhenUnset(t *testing.T) {
	reset()
	if got := DefaultUserAgent(); got != BuiltinUserAgent {
		t.Fatalf("unset override: got %q, want builtin %q", got, BuiltinUserAgent)
	}
}

func TestSetDefaultUserAgent_EmptyIsNoOp(t *testing.T) {
	reset()
	SetDefaultUserAgent("   ") // blank must not clobber the builtin default
	if got := DefaultUserAgent(); got != BuiltinUserAgent {
		t.Fatalf("blank override should be ignored: got %q, want %q", got, BuiltinUserAgent)
	}
}

func TestSetDefaultUserAgent_OverrideWins(t *testing.T) {
	reset()
	const ua = "Mozilla/5.0 (compatible; xevon; +https://github.com/xevonlive-dev/xevon)"
	SetDefaultUserAgent("  " + ua + "  ") // surrounding whitespace is trimmed
	if got := DefaultUserAgent(); got != ua {
		t.Fatalf("override: got %q, want %q", got, ua)
	}
}

func TestDefaultUserAgent_VersionPlaceholderExpansion(t *testing.T) {
	reset()
	SetBuildVersion("v9.9.9")
	SetDefaultUserAgent("Mozilla/5.0 (compatible; xevon/{version}; +https://github.com/xevonlive-dev/xevon)")
	want := "Mozilla/5.0 (compatible; xevon/v9.9.9; +https://github.com/xevonlive-dev/xevon)"
	if got := DefaultUserAgent(); got != want {
		t.Fatalf("version expansion: got %q, want %q", got, want)
	}
}

func TestDefaultUserAgent_VersionPlaceholderFallsBackToDev(t *testing.T) {
	reset()
	SetDefaultUserAgent("xevon/{version}")
	if got := DefaultUserAgent(); got != "xevon/dev" {
		t.Fatalf("empty build version: got %q, want %q", got, "xevon/dev")
	}
}
