# Installation

xevon ships as a single, statically-linked Go binary with no runtime
dependencies. Pick whichever install method fits your
environment — all of them give you the same `xevon` binary.

## Quick install (recommended)

```bash
curl -fsSL https://xevon.live/install.sh | bash
```

The installer:

- resolves the latest release from the CDN, downloads the matching binary, and
  **verifies its SHA-256 checksum** before installing;
- installs to `~/.local/bin/xevon`;
- adds `~/.local/bin` to your `PATH` in the appropriate shell profile
  (`.zshrc`, `.bashrc`/`.bash_profile`, or `config.fish`).

If `~/.local/bin` was not already on your `PATH`, activate it without
restarting your shell:

```bash
export PATH="$HOME/.local/bin:$PATH"
xevon doctor
```

## npm

The npm package is a thin launcher that pulls the correct prebuilt binary for
your platform as an optional dependency (Node 16+).

```bash
# Global install
npm install -g @xevon/xevon

# Or run once without installing
npx @xevon/xevon scan -h
```

## Build from source

Requires **Go 1.26+**, `git`, and `make`. No C toolchain is needed.

```bash
git clone https://github.com/xevonlive-dev/xevon.git
cd xevon
make build
```

`make build` outputs `bin/xevon` and installs the binary to `$GOPATH/bin`.

> **Do not** run `go build` directly. `make build` injects version metadata and
> guarantees a clean build; an ad-hoc `go build` produces a binary that reports
> the wrong version and may be left stale in the working tree.

To place the binary on your `PATH` explicitly (assuming `$GOPATH/bin` is on it):

```bash
make install
```

Other useful targets:
`make build-all` (cross-compile Linux/macOS/Windows), `make test-unit`
(fast tests), `make lint`. Run `make help` or inspect the `Makefile` for the
full list.

## Docker

```bash
docker pull codiologies/xevon:latest

# Run any command — the entrypoint is the xevon binary
docker run --rm codiologies/xevon:latest scan -h

# Scan a target and write a report to the host
docker run --rm -v "$PWD:/out" codiologies/xevon:latest \
  scan --stateless -t https://example.com --format jsonl -o /out/results
```

## Verify the installation

```bash
xevon version     # prints version, build, and commit
xevon doctor      # validates environment, config, and optional dependencies
```

`xevon doctor` is the fastest way to confirm everything is wired up
correctly — run it after any install or upgrade.

### Auto-install missing dependencies

Most of xevon's core scanning has no external dependencies, but some
optional features need extra tooling (a browser for SPA spidering, nuclei
templates for the known-issue scan, `bun`/`pi` for certain agent drivers). If
`doctor` reports a fixable item, let it install it for you:

```bash
xevon doctor --fix              # auto-install/fix every failing check
xevon doctor --fix --only chrome,nuclei   # fix only specific items
```

`--fix` prints the report, installs what's missing, then re-checks and shows
the updated status. `--only` accepts any of:
`nuclei`, `chrome`, `bun`, `claude`, `agent-browser`, `pi`, `piolium`
(it has no effect without `--fix`).

## Updating

`xevon update` re-runs the official installer to fetch the latest release
and refreshes the local nuclei-templates checkout used by the known-issue scan:

```bash
xevon update                  # update binary + nuclei templates
xevon update --skip-templates # only reinstall the binary
xevon update --skip-binary    # only refresh nuclei templates
xevon update -F               # skip the confirmation prompt
```

The binary update always installs to `~/.local/bin/xevon`. If your running
binary lives elsewhere (e.g. a `make install` build in `$GOPATH/bin` or a
Docker image), `update` prints a warning — ensure `~/.local/bin` precedes the
old location on your `PATH`, or upgrade that copy manually.

> npm and Docker installs are upgraded through their own tooling
> (`npm update -g @xevon/xevon` / `docker pull`), not `xevon update`.

## Where xevon stores data

Everything lives under `~/.xevon/` (override with the `XEVON_HOME`
environment variable):

| Path | Contents |
|------|----------|
| `~/.xevon/xevon-configs.yaml` | Main configuration file |
| `~/.xevon/database-xevon.sqlite` | Default SQLite scan database |
| `~/.xevon/agent-sessions/` | Agentic scan session artifacts |
| `~/.xevon/prompts/` | User prompt templates |

Initialize the workspace and a starter config explicitly with:

```bash
xevon init
```

## Uninstall

```bash
rm -f ~/.local/bin/xevon      # or: $GOPATH/bin/xevon for source installs
rm -rf ~/.xevon               # optional: remove config, database, sessions
```

For npm: `npm uninstall -g @xevon/xevon`.

## Next steps

- [Quick Start](quick-start.md) — run your first scan in under a minute.
- [Native Scan & Stateless Scanning](native-scan.md) — CLI scan recipes.
- [Setting Up the Agent](setup-agent.md) — wire up AI-driven scanning.
- [Configuration Reference](../configuration.md) — every config knob.
- [Troubleshooting](../troubleshooting.md) — fixes for common issues.
