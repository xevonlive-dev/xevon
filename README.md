<p align="center">
  <img alt="xevon" src="static/xevon-logo-minimal.png" height="140" />
  <br />
  <strong>xevon — High-fidelity vulnerability scanner fusing agentic AI with native speed, modularity, and precision</strong>

  <p align="center">
  <a href="https://www.xevon.live"><img src="https://img.shields.io/badge/Internet-0078D4?style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyNCIgaGVpZ2h0PSIyNCIgdmlld0JveD0iMCAwIDI0IDI0IiBmaWxsPSJub25lIiBzdHJva2U9IndoaXRlIiBzdHJva2Utd2lkdGg9IjIiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIgc3Ryb2tlLWxpbmVqb2luPSJyb3VuZCI+PGNpcmNsZSBjeD0iMTIiIGN5PSIxMiIgcj0iMTAiLz48cGF0aCBkPSJNMTIgMmExNC41IDE0LjUgMCAwIDAgMCAyMCAxNC41IDE0LjUgMCAwIDAgMC0yMCIvPjxwYXRoIGQ9Ik0yIDEyaDIwIi8+PC9zdmc+&labelColor=black&color=black"></a>
  <a href="https://docs.xevon.live/"><img src="https://img.shields.io/badge/Documentation-0078D4?style=flat&logo=GitBook&logoColor=8be9fd&labelColor=black&color=black"></a>
  <a href="https://twitter.com/xevon"><img src="https://img.shields.io/badge/xevon-0078D4?style=flat&logo=X&logoColor=f8f8f2&labelColor=black&color=black"></a>
  <a href="https://www.linkedin.com/company/xevonlive"><img src="https://img.shields.io/badge/LinkedIn-0078D4?style=flat&logo=linkedin&logoColor=white&labelColor=black&color=black"></a>
  </p>
</p>

---

xevon provides two complementary scanning modes:

- **Native Scan** (`xevon scan`): Fast, powerful, and flexible. Deterministic multi-phase scanning with 250+ modules across content discovery, browser/SPA spidering, and active/passive audit — covering injection, access control, file/path, API/protocol, framework-specific, cloud/infra, and out-of-band (OAST) vulnerability classes.

- **Agentic Scan** (`xevon agent`): Thoroughly audits your codebase. AI-driven scanning that autonomously plans attacks, selects modules, generates custom extensions, and triages results — combining deep source-code audit with autonomous and targeted vulnerability scanning.

| UI Dashboard | Traffic Dashboard |
|:---:|:---:|
| ![Dashboard 1](https://github.com/xevonlive-dev/docs/blob/main/images/xevon-main-workbench.png?raw=true) | ![Dashboard 2](https://github.com/xevonlive-dev/docs/blob/main/images/xevon-ui-dashboard-2.png?raw=true) |

| Native scan | Agentic Scan |
|:---:|:---:|
| ![Native scan](https://github.com/xevonlive-dev/docs/blob/main/images/xevon-cli-native-scan.png?raw=true) | ![Agentic Scan](https://github.com/xevonlive-dev/docs/blob/main/images/xevon-cli-agent-audit-1.png?raw=true) |

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
  - [Option A — Quick Install Script](#option-a--quick-install-script-recommended)
  - [Option B — npm](#option-b--npm)
  - [Option C — Docker](#option-c--docker)
  - [Option D — Build from Source](#option-d--build-from-source)
- [First-Time Setup](#first-time-setup)
- [Running the Dashboard (Workbench UI)](#running-the-dashboard-workbench-ui)
- [Using xevon Without the Dashboard (CLI Only)](#using-xevon-without-the-dashboard-cli-only)
  - [Native Scan (CLI)](#native-scan-cli)
  - [Agentic Scan (CLI)](#agentic-scan-cli)
  - [Source-Code Audit (CLI)](#source-code-audit-cli)
  - [Traffic & Database (CLI)](#traffic--database-cli)
- [Server Mode](#server-mode)
- [Authenticated Scanning](#authenticated-scanning)
- [JavaScript Extensions](#javascript-extensions)
- [Full CLI Reference](#full-cli-reference)
- [Configuration](#configuration)
- [Development / Build from Source](#development--build-from-source)
- [Security](#security)
- [License](#license)

---

## Prerequisites

| Requirement | Version | Notes |
|---|---|---|
| **Go** | 1.22+ | Only needed for build-from-source |
| **Git** | any | Only needed for build-from-source |
| **Make** | any | Only needed for build-from-source |
| **Node.js / npm** | 16+ | Only for npm install method |
| **Docker** | any | Only for Docker install method |
| **OS** | Linux, macOS, Windows (WSL) | Windows native via WSL recommended |

> All install methods produce a **single self-contained binary** — no runtime dependencies needed for scanning.

---

## Installation

### Option A — Quick Install Script (Recommended)

Works on Linux, macOS, and WSL:

```bash
curl -fsSL https://xevon.live/install.sh | bash
```

The installer:
- Downloads the correct prebuilt binary for your platform
- Verifies the SHA-256 checksum
- Installs to `~/.local/bin/xevon`
- Adds `~/.local/bin` to your `PATH` in `.zshrc` / `.bashrc` / `.bash_profile`

Activate immediately without restarting your shell:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Verify:

```bash
xevon version
xevon doctor
```

---

### Option B — npm

```bash
# Install globally
npm install -g @xevon/xevon

# Or run once without installing
npx @xevon/xevon scan -h
```

---

### Option C — Docker

```bash
# Pull the image
docker pull codiologies/xevon:latest

# Run any command (xevon is the entrypoint)
docker run --rm codiologies/xevon:latest scan -h

# Scan a target and save results to your current directory
docker run --rm -v "$PWD:/out" codiologies/xevon:latest \
  scan --stateless -t https://example.com --format jsonl -o /out/results.jsonl
```

---

### Option D — Build from Source

**Requires:** Go 1.22+, Git, Make.

```bash
# 1. Clone the repository
git clone https://github.com/xevonlive-dev/xevon.git
cd xevon

# 2. Build the binary (outputs to bin/xevon AND installs to $GOPATH/bin)
make build

# 3. Verify
xevon version
xevon doctor
```

> ⚠️ Always use `make build`, not `go build` directly. `make build` injects version metadata and produces a correct binary.

Other build targets:

```bash
make build-all    # cross-compile Linux + macOS + Windows
make install      # install to $GOPATH/bin explicitly
make test-unit    # run fast unit tests
make test         # run all tests
make lint         # run linter
make fmt          # format code
```

---

## First-Time Setup

After installation, initialize the workspace:

```bash
# Initialize ~/.xevon/ with default config
xevon init

# Run health check — shows what's installed and what's missing
xevon doctor

# Auto-fix missing optional dependencies (Chrome, Nuclei, Bun, etc.)
xevon doctor --fix
```

xevon stores everything under `~/.xevon/`:

| Path | Contents |
|---|---|
| `~/.xevon/xevon-configs.yaml` | Main config file |
| `~/.xevon/database-xevon.sqlite` | SQLite scan database |
| `~/.xevon/agent-sessions/` | Agentic scan session artifacts |
| `~/.xevon/prompts/` | Custom prompt templates |

Override the home directory:
```bash
XEVON_HOME=/custom/path xevon server
```

---

## Running the Dashboard (Workbench UI)

The dashboard is a web UI served by the xevon API server. It requires the server to be running.

### Step 1 — Start the server

```bash
# Start with no authentication (local use only)
xevon server

# Start with an API key (recommended for any network-accessible setup)
xevon server -k my-secret-key

# Start with a custom port
xevon server --service-port 9002 -k my-key

# Start with transparent proxy for recording browser traffic
xevon server -k my-key --ingest-proxy-port 9003
```

### Step 2 — Open the Dashboard

Open your browser and navigate to:

```
http://localhost:9002
```

The dashboard (xevon Workbench) will load automatically. It includes:

| Tab | What it does |
|---|---|
| **Dashboard** | Overview — stats, recent scans, findings summary |
| **Findings** | Browse, filter, and export vulnerability findings |
| **HTTP Records** | Inspect recorded traffic |
| **OAST** | Out-of-band interaction callbacks |
| **Modules** | List and configure scanner modules |
| **Extensions** | Manage JavaScript extension scripts |
| **Ingest** | Upload specs (OpenAPI, Burp, HAR, curl) |
| **Native Scan** | Configure and launch native scans |
| **Agentic Scan** | Configure and launch AI-driven scans |
| **Database** | Browse the raw SQLite database |
| **Settings** | Configure server settings |

### Step 3 — Create a Project (optional)

Click **[PROJECT: ▼]** in the top bar → **+ New Project** to organize scans.

### Connecting with an API key

If you started the server with `-k my-secret-key`, log in at `http://localhost:9002/login` using that key.

### Workbench — Agentic Scan via Dashboard

1. Go to **Agentic Scan** tab
2. Enter your target URL or paste a raw HTTP request/curl command
3. Select mode: **Autopilot**, **Swarm**, or **Audit**
4. Select intensity: **Quick**, **Balanced**, or **Deep**
5. *(Optional)* Upload source code (zip/folder) for whitebox analysis
6. Click **START SCAN**

> 📌 **Note on Source Upload:** The "Upload source code" widget on the Agentic Scan page uploads the file/folder to the server's local disk under `~/.xevon/repos/`. If you're running the server locally, you can also just pass `--source /path/to/code` directly via CLI (see below), which is simpler.

---

## Using xevon Without the Dashboard (CLI Only)

You do **not** need the dashboard or a running server for most operations. Everything the dashboard does can be done from the CLI.

### Native Scan (CLI)

```bash
# Basic scan of a target URL
xevon scan -t https://example.com

# Deep scan (all phases + extended discovery)
xevon scan -t https://example.com --strategy deep

# Scan specific modules only
xevon scan -t https://example.com -m xss-reflected,sqli-error

# Scan from an OpenAPI/Swagger spec
xevon scan -I openapi -i api-spec.yaml -t https://api.example.com

# Scan using the spec's own servers list
xevon scan -I openapi -i api-spec.yaml --spec-url

# Scan from a Burp Suite export
xevon scan -I burp -i burp-export.xml -t https://example.com

# Scan from a HAR file
xevon scan -I har -i traffic.har

# Scan from a cURL command file
xevon scan -I curl -i requests.txt

# Scan multiple URLs from a file
xevon scan -T targets.txt

# Pipe URLs from stdin
cat urls.txt | xevon scan

# Scan a single raw HTTP request (from file or stdin)
xevon scan-request -i raw-request.txt

# Scan a single URL with custom method/headers
xevon scan-url https://api.example.com/login \
  --method POST \
  --body '{"user":"admin","pass":"test"}' \
  -H "Content-Type: application/json"

# Run a single phase only
xevon run discovery -t https://example.com
xevon run spidering -t https://example.com
xevon run audit -t https://example.com

# Filter modules by tag
xevon scan -t https://example.com --module-tag spring --module-tag injection

# Generate an HTML report
xevon scan -t https://example.com --format html -o report.html

# Generate JSONL output (machine-readable)
xevon scan -t https://example.com --format jsonl -o results.jsonl

# Stateless scan (no database writes — pure stdout)
xevon scan --stateless -t https://example.com

# Scan with proxy (Burp Suite integration)
xevon scan -t https://example.com --proxy http://127.0.0.1:8080

# Tune performance
xevon scan -t https://example.com -c 20 --rate-limit 30 --timeout 30s
```

**Strategies:**

| Strategy | Discovery | Spidering | Known Issues | Audit |
|---|:---:|:---:|:---:|:---:|
| `lite` | ✗ | ✗ | ✗ | ✓ |
| `balanced` *(default)* | ✓ | ✓ | ✓ | ✓ |
| `deep` | ✓ | ✓ | ✓ | ✓ |
| `whitebox` | ✓ | ✗ | ✓ | ✓ |

```bash
xevon scan -t https://example.com --strategy lite      # fast — audit only
xevon scan -t https://example.com --strategy balanced  # default
xevon scan -t https://example.com --strategy deep      # thorough
xevon scan -t https://example.com --strategy whitebox --source ./src
```

---

### Agentic Scan (CLI)

AI-driven scanning. Set up an AI provider first (see [Configuration](#configuration)).

#### Autopilot — Autonomous AI Scanning

```bash
# Basic autonomous scan
xevon agent autopilot -t https://example.com

# With source code context (AI reads code + runs scan)
xevon agent autopilot -t https://example.com --source ./src

# Natural-language prompt (target auto-extracted)
xevon agent autopilot "scan VAmPI at localhost:3005 with source at ~/src/VAmPI"

# Intensity presets (bundle scan settings + audit mode)
xevon agent autopilot -t https://example.com --intensity quick     # fast / CI
xevon agent autopilot -t https://example.com --intensity balanced  # default
xevon agent autopilot -t https://example.com --intensity deep      # full pentest

# Focus on PR diff only (change-aware scanning)
xevon agent autopilot -t https://example.com --source ./src --diff main...feature-branch
xevon agent autopilot -t https://example.com --source ./src --last-commits 3

# Custom wall-clock time budget
xevon agent autopilot -t https://example.com --max-duration 30m

# With browser-based auth preflight
xevon agent autopilot -t https://example.com --browser --credentials "admin/admin123"

# Run AI triage on findings after scan
xevon agent autopilot -t https://example.com --triage

# Preview the rendered system prompt without running
xevon agent autopilot -t https://example.com --dry-run
```

#### Swarm — AI-Guided Targeted Scanning

```bash
# Target a specific endpoint
xevon agent swarm -t https://example.com/api/users

# Full-scope scan (discovery + planning + scan + triage)
xevon agent swarm -t https://example.com --discover

# Source-aware full-scope scan
xevon agent swarm -t https://example.com --source ./src --discover

# From a curl command
xevon agent swarm --input "curl -X POST https://example.com/api/login -d '{\"user\":\"admin\"}'"

# Focus on a vulnerability type
xevon agent swarm -t https://example.com/api/search --vuln-type sqli

# Natural-language prompt
xevon agent swarm "scan source at ~/src/app on localhost:3005"

# Intensity presets
xevon agent swarm -t https://example.com --intensity quick
xevon agent swarm -t https://example.com --intensity deep

# Source analysis only (extract routes, no scan)
xevon agent swarm -t https://example.com --source ./src --source-analysis-only

# Scan only changed files
xevon agent swarm -t https://example.com --source ./src --diff main...feature-branch

# Custom instructions to guide the agent
xevon agent swarm -t https://example.com --instruction "Focus on GraphQL parsing"
```

#### Query — Single-Shot Code Review / Prompt

```bash
# Security code review of a local codebase
xevon agent query --prompt-template security-code-review --source ./src

# Endpoint discovery from source
xevon agent query --prompt-template endpoint-discovery --source ./src

# Inline prompt
xevon agent query "review this code for SQL injection vulnerabilities"

# List available templates
xevon agent --list-templates

# Dry-run (preview rendered prompt without sending to AI)
xevon agent query --prompt-template security-code-review --source ./src --dry-run

# Save output to a file
xevon agent query --prompt-template security-code-review --source ./src \
  --output review-results.json
```

---

### Source-Code Audit (CLI)

Deep AI-driven static analysis of your source code. No target URL required.

```bash
# Quick audit (default: auto-picks claude or codex CLI)
xevon agent audit --source ./src

# Choose audit depth
xevon agent audit --source ./src --mode lite        # fast
xevon agent audit --source ./src --mode balanced    # default
xevon agent audit --source ./src --mode deep        # thorough

# Force a specific driver
xevon agent audit --driver=audit --source ./src      # claude/codex only
xevon agent audit --driver=piolium --source ./src    # Pi-native (requires pi + piolium)
xevon agent audit --driver=both --source ./src       # run both back-to-back

# Audit a remote GitHub repo (auto-clones)
xevon agent audit --source https://github.com/org/repo --mode lite

# Audit from a cloud archive
xevon agent audit --source gs://my-bucket/app.tar.gz

# Interactive mode (you drive the audit manually)
xevon agent audit -i --source ./src

# Check progress on an in-progress audit
xevon agent audit --mode status --source ./in-progress-audit-dir

# Second pass / revisit
xevon agent audit --driver=audit --mode revisit --source ./prior-audit-tree

# Construct PoCs for confirmed findings
xevon agent audit --driver=audit --mode confirm --source ./audit-with-findings
```

**Audit intensity presets:**

| Preset | Mode chain | Clone depth |
|---|---|---|
| `quick` | `lite` | shallow |
| `balanced` *(default)* | `balanced` | shallow |
| `deep` | `deep,confirm` | full history |

```bash
xevon agent audit --intensity deep --source ./src
```

---

### Traffic & Database (CLI)

```bash
# Browse recorded HTTP traffic (TUI)
xevon traffic

# Filter traffic
xevon traffic --host api.example.com --method POST
xevon traffic --tree        # hierarchical host view
xevon traffic login         # fuzzy search

# JSONL output for scripting
xevon traffic -j --host api.example.com

# Browse findings (TUI)
xevon finding
xevon finding --severity high,critical
xevon finding --module-type active
xevon finding -j --severity critical   # JSONL output

# Database stats
xevon db stats
xevon db stats --detailed

# Export results
xevon export --format jsonl -o results.jsonl
xevon export --format html -o report.html
xevon export --only findings -o findings.jsonl

# View scan session logs
xevon log <scan-uuid>              # view runtime log
xevon log <scan-uuid> -f           # follow (tail -f style)
xevon log ls                       # list all sessions

# Replay a stored request with mutations (confirm a finding)
xevon replay --record-uuid <uuid> -m 'name=id,payload=1 OR 1=1'
xevon replay --finding-id 42 -m 'name=q,payload=<svg/onload=alert(1)>'
```

---

## Server Mode

The server exposes a REST API + the web dashboard on port 9002 (default). You can also use the REST API directly without the dashboard.

```bash
# Start server (no auth — local only)
xevon server

# Start with API key authentication
xevon server -k my-secret-key

# Start with specific host + port
xevon server --host 0.0.0.0 --service-port 8443 -k my-key

# Enable transparent proxy (records all browser/app traffic)
xevon server -k my-key --ingest-proxy-port 9003

# Auto-scan new traffic as it arrives
xevon server -k my-key -t https://example.com --scan-on-receive

# Start in view-only mode (read-only dashboard — no scan triggers)
xevon server -k my-key --view-only
```

### Ingest traffic to a running server

```bash
# Ingest URLs
cat urls.txt | xevon ingest -s http://localhost:9002

# Ingest an OpenAPI spec
xevon ingest -s http://localhost:9002 -i api.yaml -I openapi

# Ingest and auto-scan
xevon ingest -s http://localhost:9002 -i api.yaml -I openapi -S
```

### REST API (no dashboard needed)

```bash
# Trigger a scan via API
curl -X POST http://localhost:9002/api/agent/run/autopilot \
  -H "Authorization: Bearer my-secret-key" \
  -H "Content-Type: application/json" \
  -d '{"target":"https://example.com","intensity":"balanced"}'

# List findings
curl http://localhost:9002/api/findings \
  -H "Authorization: Bearer my-secret-key"

# Get server info
curl http://localhost:9002/server-info
```

Full REST API reference: [docs.xevon.live/api-overview](https://docs.xevon.live/api-overview)

---

## Authenticated Scanning

```bash
# Inline auth (name:Header:value — repeatable for multi-session IDOR testing)
xevon scan -t https://example.com \
  --auth "admin:Cookie:session_id=abc123" \
  --auth "user:Cookie:session_id=xyz789"

# Load session from YAML/JSON file
xevon scan -t https://example.com --auth-file ./admin-session.yaml

# Custom HTTP headers
xevon scan -t https://example.com -H "Authorization: Bearer token123"

# Browser-based auth (xevon drives Chromium to log in and capture cookies)
xevon agent autopilot -t https://example.com \
  --browser --credentials "username=admin,password=secret"
```

Auth file format (`session.yaml`):
```yaml
name: admin
headers:
  Cookie: "session=abc123; csrf=xyz"
  Authorization: "Bearer eyJ..."
```

---

## JavaScript Extensions

```bash
# Execute inline JavaScript
xevon js --code 'let r = xevon.http.get(TARGET); console.log(r.status)' \
  -t https://example.com

# Run a JS file
xevon js --code-file ./my-script.js -t https://example.com --timeout 60s

# List loaded extensions
xevon ext ls

# Browse the API with examples
xevon ext docs --example

# Install starter extension scripts
xevon ext preset

# Lint an extension file
xevon ext lint --ext custom-check.js

# Run a custom extension against a target
xevon run extension -t https://example.com --ext ./custom-check.js
```

---

## Full CLI Reference

```
Scanning:
  xevon scan                Run a native multi-phase scan
  xevon run <phase>         Run a single scan phase (alias: scan --only <phase>)
  xevon scan-url <url>      Quick scan of a single URL
  xevon scan-request        Scan from a raw HTTP request (file or stdin)

Agentic scan (AI-driven):
  xevon agent autopilot     Autonomous AI-driven vulnerability scanning
  xevon agent swarm         AI-guided targeted or full-scope vulnerability scanning
  xevon agent query         Single-shot prompt (code review, endpoint discovery)
  xevon agent audit         Source-code audit (xevon-audit + piolium drivers)
  xevon agent session       Browse / inspect agent session artifacts
  xevon olium | xevon ol   Interactive olium TUI (or one-shot via -p)

Server & ingestion:
  xevon server              Start API server + web dashboard
  xevon ingest              Ingest traffic to a running server
  xevon storage             Manage cloud object storage (ls/upload/download/rm)

Data & projects:
  xevon db                  Database operations (stats, export, clean, seed, ls)
  xevon finding             Browse and manage vulnerability findings
  xevon traffic             Browse and replay HTTP traffic records
  xevon replay              Mutate a stored request and diff baseline vs replay
  xevon project             Manage projects (create, list, use)
  xevon scope               Manage scope / allow-list rules
  xevon import              Import findings or data from external sources
  xevon export              Export scan results (JSONL, HTML)
  xevon log                 View / follow runtime logs for scan sessions

Extensions & auth:
  xevon js                  Execute JavaScript / TypeScript inline or from file
  xevon ext                 Manage JavaScript extensions (ls, docs, preset, lint, eval)
  xevon auth                Manage auth sessions (list, load, lint, totp)

Setup & introspection:
  xevon init                Initialize ~/.xevon/ workspace and default config
  xevon config              Manage configuration (ls, set, path, clean)
  xevon strategy            View scanning strategies and phases
  xevon module              Inspect / enable / disable scanner modules
  xevon doctor              Diagnose environment and dependencies
  xevon version             Show version, build, and commit info
  xevon update              Update binary and nuclei templates
```

### Key Flags

```
Native Scan:
  -t, --target         Target URL
  -T, --target-file    File with target URLs (one per line)
  -i, --input          Input file (- for stdin)
  -I, --input-mode     Input format: urls, openapi, swagger, burp, curl, har, postman, nuclei
  -m, --modules        Modules to run (comma-separated IDs, or 'all')
      --strategy       Strategy preset: lite, balanced, deep, whitebox
      --only           Run single phase only
      --skip           Skip a phase
  -c, --concurrency    Concurrent workers (default: 50)
  -r, --rate-limit     Max requests/sec (0 = unlimited)
      --timeout        Per-request timeout (default: 15s)
      --proxy          HTTP/SOCKS5 proxy URL

Authentication:
      --auth           Inline session: name:Header:value (repeatable)
      --auth-file      Session YAML file (repeatable)
  -H, --header         Custom HTTP header (repeatable)

Output:
  -j, --json           JSON output
      --format         Output format: console, jsonl, html
  -o, --output         Output file path
      --silent         Suppress all output except findings
  -v, --verbose        Verbose logging

Agentic Scan:
      --source         Path to source code for whitebox / AI analysis
      --provider       AI provider: openai-codex-oauth, anthropic-api-key,
                       anthropic-oauth, openai-api-key, anthropic-cli,
                       anthropic-vertex, google-vertex
      --model          Model ID override
      --llm-api-key    API key (for anthropic-api-key / openai-api-key)
      --oauth-token    OAuth bearer token (anthropic-oauth)
      --intensity      Intensity preset: quick, balanced, deep
      --diff           Diff range for change-focused scans (e.g. main...feature)
      --last-commits   Shorthand for --diff HEAD~N
      --discover       Run discovery+spidering before planning (swarm)
      --triage         Run AI triage pass on findings
      --audit          Background audit mode: lite, balanced, deep, off
      --driver         Audit driver: auto, both, audit, piolium (agent audit)
      --mode           Audit mode: lite, balanced, deep, revisit, confirm...
      --max-duration   Wall-clock time cap for the agent (e.g. 2h, 30m)
      --dry-run        Preview rendered prompt without executing
```

---

## Configuration

Configuration lives at `~/.xevon/xevon-configs.yaml`. View and edit:

```bash
# List current config
xevon config ls

# Set a value
xevon config set agent.olium.provider anthropic-api-key
xevon config set agent.olium.llm_api_key sk-...

# View config file path
xevon config path

# Reset to defaults
xevon config clean
```

### Setting up an AI Provider (for Agentic Scans)

```bash
# Option 1: Anthropic API key
xevon config set agent.olium.provider anthropic-api-key
xevon config set agent.olium.llm_api_key sk-ant-...

# Option 2: OpenAI API key
xevon config set agent.olium.provider openai-api-key
xevon config set agent.olium.llm_api_key sk-...

# Option 3: Local Ollama (default — no API key needed)
# Make sure Ollama is running with a model pulled:
ollama pull gemma3
# Then xevon will use it automatically (default provider: openai-compatible → localhost:11434)

# Option 4: Claude CLI (requires Claude Max subscription)
claude setup-token   # one-time setup
xevon config set agent.olium.provider anthropic-cli

# Option 5: Anthropic OAuth (Claude Code)
xevon config set agent.olium.provider anthropic-oauth
```

### Minimal `xevon-configs.yaml` example

```yaml
agent:
  olium:
    provider: anthropic-api-key
    llm_api_key: sk-ant-YOUR_KEY_HERE

server:
  api_keys:
    - my-secret-key
```

---

## Development / Build from Source

```bash
# Full build (binary + install)
make build

# Cross-compile for all platforms
make build-all

# Run all tests
make test

# Fast unit tests only
make test-unit

# End-to-end tests (requires Docker)
make test-e2e

# Lint
make lint

# Format code
make fmt

# See all targets
make help
```

See [HACKING.md](HACKING.md) for the full build guide, codebase map, and module development guide.

### Building the Dashboard (optional)

The compiled dashboard is already embedded in the binary. If you want to rebuild it from source:

```bash
# Requires: bun 1.3.11+ (https://bun.sh)
cd platform/xevon-workbench
bun install
bun run build
# Output goes to public/ui/ — picked up automatically by make build
```

---

## Security

xevon is an offensive security tool. Two parts are intentionally permissive:

- **Agent mode runs with no sandbox** — the LLM has full shell, file, and network access on the host. Run agent mode in a disposable container/VM scoped to the engagement.
- **Extensions can run arbitrary commands** — treat untrusted extensions like untrusted code.

Report vulnerabilities in xevon itself privately to [contact@xevon.live](mailto:contact@xevon.live). See [SECURITY.md](SECURITY.md) for the full disclosure policy.

---

## License

xevon is released under the [GNU Affero General Public License v3.0](LICENSE). Derivative works must remain open under the same terms.

Developed by [@codiologies](https://github.com/codiologies).
