// Package bin embeds the xevon-audit host binary into xevon and
// extracts it to a per-user cache directory on first use.
//
// The binary is built from platform/xevon-audit/ by `make update-audit`,
// which runs `bun run build` and copies the host-platform output to
// _bin/xevon-audit. Cross-compiling xevon requires staging the
// matching xevon-audit-<os>-<arch> blob at _bin/xevon-audit before
// `go build`; the release does this per-target via the goreleaser pre-hook
// build/scripts/stage-audit-blob.sh (builds run with --parallelism 1 because
// the embed path is a single shared file). verifyBlobForHost (verify.go)
// guards against a wrong-platform blob slipping through at runtime.
package bin

import (
	"embed"
)

// binFS holds whatever lives under _bin at compile time. `all:` lets the
// embed directive match an empty _bin (only the tracked .gitkeep) on
// fresh clones — missing-binary surfaces as a runtime extract error
// rather than a build failure.
//
//go:embed all:_bin
var binFS embed.FS

const embeddedName = "_bin/xevon-audit"
