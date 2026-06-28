<p align="center">
  <a href="https://github.com/xevon"><img alt="xevon" src="https://avatars.githubusercontent.com/u/266502139?s=200&v=4" height="140" /></a>
  <br />
  <strong>xevon - High-fidelity vulnerability scanner with native scan precision and agentic scan intelligence.</strong>
  <br />

  <p align="center"><a href="https://xevon.live">xevon.live</a> - <a href="https://docs.xevon.live"> docs.xevon.live</a></p>
</p>

# xevon-Audit

`xevon-audit` is an autonomous agent within [xevon](https://xevon.live/) that performs comprehensive security audits on your repository, focusing on uncovering exploitable vulnerabilities with high accuracy.

> [!WARNING]
> A full audit run can take a few hours. Go enjoy your coffee ☕ and take a walk. Don't worry, it's worth the wait.

## Why?

Static analysis tools bury you in false positives. Manual audits are thorough but slow and expensive. `xevon-audit` runs security audits as a multi-agent pipeline: each phase builds on the last — from gathering advisories, to flagging candidates, to proposing attack paths, to debating exploitability, to a final verification pass that kills false positives. The workflow is resumable and incremental — re-run after a code change and only affected phases re-execute.

The goal is simple: spend machine time instead of human time, and only surface findings that are real.

## Install as a standalone CLI (recommended)

```bash
npm install -g @xevon/xevon-audit
```

Or via curl:

```bash
curl -fsSL https://cdn.xevon.live/xevon-audit/install.sh | bash
```

### Install as a Skill

Already using the [`skills`](https://www.npmjs.com/package/skills) CLI? Add `xevon-audit` as a skill so your agent learns the modes, flags, and example commands automatically:

```bash
npx skills@latest add https://github.com/xevonlive-dev/xevon-audit --skill xevon-audit --agent claude-code --yes
```

The skill teaches the agent how to install the binary, pick a mode, resume interrupted runs, and produce machine-readable output. The CLI binary itself still needs to be installed via one of the methods above.

| xevon Audit | xevon Audit in Codex |
|:---:|:---:|
| ![xevon Audit](https://github.com/xevonlive-dev/docs/blob/main/images/audit/xevon-audit-standalone.png?raw=true) | ![xevon Audit in Codex](https://github.com/xevonlive-dev/docs/blob/main/images/audit/xevon-audit-in-codex.png?raw=true) |

## Requirements

[`claude`](https://www.npmjs.com/package/@anthropic-ai/claude-code) and/or [`codex`](https://www.npmjs.com/package/@openai/codex) on PATH, plus the matching API key or ambient subscription auth.

## Quickstart

```bash
# just start the audit
xevon-audit run --mode balanced --target /path/to/repo/

# Run a deep audit interactively (auto-installs the harness for the session
# and removes it on exit — leave-no-trace)
xevon-audit run --mode deep --agent claude -i

# Headless deep audit, abort if cost exceeds $20
xevon-audit run --mode deep --agent codex --max-cost 20

# Preflight: binary, auth, content, real message round-trip
xevon-audit verify claude
```

### Auth overrides

Three flags on `xevon-audit run` swap auth in for the lifetime of one run, then
restore the original state on exit:

```bash
# Pass an API key via flag (claude → ANTHROPIC_API_KEY, codex → OPENAI_API_KEY)
xevon-audit run --mode deep --agent claude --api-key sk-ant-...

# Set CLAUDE_CODE_OAUTH_TOKEN for the subprocess / SDK
xevon-audit run --mode deep --agent claude --oauth-token sk-ant-oat01-...

# Temporarily replace ~/.codex/auth.json with a custom file
# (~/.claude/.credentials.json for --agent claude). The original file is moved
# to <target>.xevon-audit-backup before the run and restored on exit.
xevon-audit run --mode deep --agent codex --oauth-cred-file ./codex-auth.json
```

Secrets passed via flag are redacted in the `[auth] applied: …` log line. The
cred-file backup is also restored on Ctrl-C / SIGTERM.

### Machine-readable output

Every command supports `--json` for tooling. Logs stay on stderr; structured
JSON goes to stdout (single object for verify/uninstall, NDJSON event
stream for `run`).

```bash
xevon-audit verify claude --json | jq .ok
xevon-audit run --mode lite --agent claude --json | jq -c 'select(.kind == "phaseEnd")'
```

## Audit Modes

xevon-Audit ships a handful of audit modes — each is a different phase graph
selecting how thorough vs. fast the run is, and what it focuses on. The
short version:

| Mode | Use it when |
|------|-------------|
| `lite` | You want a fast surface scan (secrets + SAST + PoC) on a plain folder. |
| `balanced` | You want a real audit but not the full deep pipeline — middle ground. |
| `deep` | You want the full multi-agent pipeline: highest signal, longest run. |
| `revisit` | You already have a complete `deep` result and want a second anti-anchored pass. |
| `reinvest` | You want cross-agent re-verification of existing CRIT/HIGH findings (claude ↔ codex). |
| `confirm` | You want findings exercised against a live or booted target. |
| `diff` | A small change landed; re-run only deep phases the diff affects. |
| `merge` | You ran xevon-audit multiple times and want one normalized findings tree. |
| `longshot` | Architecture-anchored audits feel exhausted — bottom-up file-by-file hail-mary. |
| `refresh` | You don't want to pick: router resolves to `revisit` or fresh `deep`. |

Run `xevon-audit list` for the live view — descriptions, phase counts, and
observed median runtime from your prior runs. Canonical phase definitions
live in [`src/content/command-defs/`](./src/content/command-defs/);
overriding or extending them (per-user or per-project) is documented in
[`CUSTOMIZATION.md`](./CUSTOMIZATION.md).

### `deep` mode phases

`deep` is a 12-phase pipeline where each phase feeds the next — intel first,
then static analysis, then adversarial review, then PoC and reporting. `D1`/`D2`
run in parallel; `D2` and `D3` are skipped on a no-git target.

| Phase | What it does |
|-------|--------------|
| `D1` Intelligence Pass (CVE) | Collect known CVEs/advisories for the project's stack and dependencies. |
| `D2` Intelligence Pass (History) | Mine git history for security-relevant commits, regressions, and risky changes. *(git only)* |
| `D3` Patch Audit | Inspect prior security fixes for incomplete patches or bypasses. *(git only)* |
| `D4` Threat Model | Build the attack surface and threat model from architecture and entry points. |
| `D5` Code Scan | Static (SAST-style) scan for vulnerable code patterns; enumerates cross-service edges. |
| `D6` Deep Probe | Targeted, manual-style investigation of the most suspicious areas. |
| `D7` Access Audit | Authn/authz and access-control review — broken access, IDOR, privilege escalation. |
| `D8` Review Panel | Adversarial review chamber: debates exploitability and kills false positives (also folds in taint reasoning + variant expansion). |
| `D9` Intent Reconciliation | Reconcile survivors against intended behavior to drop by-design "findings". |
| `D10` PoC Authoring | Write concrete proof-of-concept exploits for the surviving findings. |
| `D11` Finding Finalize | Normalize and finalize findings into the canonical `xevon-results/findings/` tree. |
| `D12` Report Compose | Assemble the final audit report. |

### Output cleanup

Completed `deep`/`confirm` runs, including successful resumes of those modes, automatically prune
raw workspaces after success so the final `xevon-results/` tree contains only durable
deliverables: state JSON, `file-state.json`, `attack-surface/`, finalized
`findings/` + `findings-theoretical/`, mode reports, and `confirm-workspace/`
for confirmation runs. Failed or aborted runs keep raw directories for resume
and debugging. Use `--strip-raw` or `xevon-audit strip <path>` for modes that
do not auto-prune.

### Resuming an interrupted audit

If a run is killed mid-way (quota limit, SIGINT, `--max-cost` cap, crash),
the audit stays non-complete in `xevon-results/audit-state.json`. Pick it up
where it left off — completed phases are skipped, stale `in_progress`
phases are quarantined and retried:

```bash
xevon-audit resume ./repo                       # auto-detect mode + audit
xevon-audit run --mode deep --resume            # explicit form
```

## Project Structure

```
xevon-audit/
├── src/
│   ├── cli/                # run / setup / verify / uninstall entry points
│   ├── engine/             # orchestrator, phase parser, state, harness, modes
│   ├── adapters/           # claude/codex CLI + SDK adapters, platform detect
│   ├── content/            # vendored audit methodology
│   │   ├── agent-defs/         # 31 specialist agent prompts (.md)
│   │   ├── command-defs/       # 9 mode workflows (lite/balanced/deep/…)
│   │   ├── skills/             # 20 standalone workflow skills
│   │   ├── harnesses/          # platform-specific frontmatter (claude, codex)
│   │   ├── sdk-variants/       # generated SDK-safe variants (gitignored)
│   │   └── skills-lock.json    # skill version locks
│   ├── content-bundle.json # build-time inlined content for the compiled binary
│   ├── content-loader.ts   # resolves vendored content + per-user overrides
│   └── index.ts            # CLI entry point
├── build/                  # release packaging (build.ts, install.sh)
├── scripts/                # transform-content.ts, sync helpers
└── tests/
    └── fixtures/           # sample runtime logs (`xevon-audit-runtime.log`, `xevon-audit-json-output.jsonl`) captured from real runs
```

Audit output lands in `xevon-results/` inside the target repository (e.g. `xevon-results/audit-state.json`, `xevon-results/findings/`, `xevon-results/final-audit-report.md`).

Per-user customization (overrides for agents, commands, skills) lives in `~/.config/xevon-audit/`. See [`CUSTOMIZATION.md`](./CUSTOMIZATION.md).

## Development

```bash
bun install
bun run dev -- run --mode lite --agent claude --target ./fixtures/tiny-vuln
bun test
bun run build           # current platform binary → build/dist/ + ~/.local/bin/xevon-audit
bun run build:all       # all 4 targets (host platform also installed)
bun run npm-publish     # publish @xevon/xevon-audit (single bundled pkg) to npm
```

`bun run build` also copies the just-built binary to `~/.local/bin/xevon-audit`
for fast local testing. Override the destination with `XEVON_AUDIT_BIN_DIR=…` or
skip the copy with `XEVON_AUDIT_BUILD_NO_INSTALL=1`.

## License

xevon-audit is made with ♥ by [codiologies](https://github.com/codiologies) and it is released under the MIT license.
