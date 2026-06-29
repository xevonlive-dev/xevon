# xevon-scanner

Claude Code skill for operating the [xevon](https://xevon.live/) web vulnerability scanner CLI. See the [documentation](https://docs.xevon.live/) for full details.

## What is xevon?

[xevon](https://xevon.live/) represents Next-Generation Vulnerability Discovery Powered by Agentic AI and Built for Scale. It combines traditional DAST scanning with AI-powered analysis to find vulnerabilities in web applications. Full documentation is available at [docs.xevon.live](https://docs.xevon.live/). Key capabilities:

- **Multi-phase scanning** — discovery, spidering, SPA analysis, audit, and SAST
- **Flexible input** — scan URLs directly, or import from OpenAPI specs, Burp exports, HAR files, cURL commands, and more
- **AI agent modes** — autonomous scanning (autopilot, swarm), foreground archon audits, and AI-assisted code review
- **Extensible** — write custom scanner modules in JavaScript
- **Source-aware** — whitebox scanning that combines static analysis with dynamic testing

---

# Using the xevon Scanner Skill in Claude Code, Codex or any agents

This guide explains how to install and use the `xevon-scanner` skill (`skills/xevon-scanner/`) with AI coding agents — Claude Code and OpenAI Codex — to operate the xevon CLI for web vulnerability scanning, security testing, and custom extension authoring.

## Table of Contents

- [What the Skill Does](#what-the-skill-does)
- [Skill Structure](#skill-structure)
- [Installation](#installation)
- [Usage Examples by Category](#usage-examples-by-category)
  - [Scanning](#1-scanning)
  - [Input Formats](#2-input-formats)
  - [Phase Control](#3-phase-control)
  - [Module Filtering](#4-module-filtering)
  - [Server & Ingestion](#5-server--ingestion)
  - [AI Agent Modes](#6-ai-agent-modes)
  - [Traffic & Results Browsing](#7-traffic--results-browsing)
  - [Data Management](#8-data-management)
  - [Export & Reports](#9-export--reports)
  - [Whitebox / Source-Aware Scanning](#10-whitebox--source-aware-scanning)
  - [JavaScript Extensions](#11-javascript-extensions)
  - [Configuration & Projects](#12-configuration--projects)
- [Natural Language Examples](#natural-language-examples)
- [Tips & Best Practices](#tips--best-practices)

---

## What the Skill Does

The skill teaches the AI agent how to:

1. **Pick the right xevon command** for any security testing task
2. **Construct correct flag combinations** with proper syntax
3. **Follow scanning workflows** end-to-end (ingest → scan → triage → export)
4. **Write custom JavaScript extensions** using the `xevon.*` API
5. **Execute ad-hoc JavaScript** with `xevon js` for scripting and automation
6. **Operate AI agent modes** (query, autopilot, swarm, archon)
7. **Manage data** — browse traffic, filter findings, export reports, clean databases

The skill uses lazy-loaded references: the main `SKILL.md` stays small, and detailed docs are loaded on demand when the agent needs deep flag information or extension authoring guidance.

---

## Skill Structure

```
skills/xevon-scanner/
├── SKILL.md                              # Main skill (decision tree, recipes, flags)
└── references/
    ├── scanning-commands.md              # scan, scan-url, scan-request, run
    ├── server-and-ingestion.md           # server, ingest, traffic, traffic replay
    ├── agent-commands.md                 # agent, agent query, autopilot, swarm, archon, session
    ├── data-and-management.md            # db, module, ext, js, config, scope, source, export
    ├── flags-reference.md                # Complete alphabetical flag index
    ├── session-auth-config.md            # Auth config YAML format, extract rules
    └── writing-extensions.md             # JS extension API and examples
```

---

## Installation

**Option A: Install via npx / bunx (recommended)**

```bash
bunx skills add xevon/skills --skill xevon-scanner --agent <agent-name> --yes
```

or with `npx`:

```bash
npx skills add xevon/skills --skill xevon-scanner --agent <agent-name> --yes
```

Replace `<agent-name>` with your agent (e.g., `claude-code`, `codex`). This fetches the skill from the [xevon/skills](https://github.com/xevonlive-dev/skills) repository and registers it automatically.

**Option B: Clone and copy manually**

```bash
git clone https://github.com/xevonlive-dev/skills.git
cd skills
```

Then copy the skill folder to your agent's configuration directory:

```bash
# For Claude Code
cp -R xevon-scanner ~/.claude

# For other agents
cp -R xevon-scanner ~/.agents
```

Once installed, the skill **auto-triggers** when you mention keywords like `scan`, `xevon`, `agent autopilot`, `vulnerability scanner`, `openapi scan`, etc. In Claude Code, you can also invoke it explicitly with `/xevon-scanner`.

---

## Usage Examples by Category

### 1. Scanning

**Basic scan against a single target:**
```
> Scan https://example.com for vulnerabilities
```
```bash
xevon scan -t https://example.com
```

**Multiple targets:**
```
> Scan both https://example.com and https://api.example.com
```
```bash
xevon scan -t https://example.com -t https://api.example.com
```

**Targets from a file:**
```
> I have a list of URLs in targets.txt, scan all of them
```
```bash
xevon scan -T targets.txt
```

**Scan with a specific strategy:**
```
> Do a deep scan of https://example.com with discovery and spidering
```
```bash
xevon scan -t https://example.com --strategy deep
```

**Scan a single URL with custom method, headers, and body:**
```
> Test the login endpoint for vulnerabilities: POST to https://api.example.com/login with JSON credentials
```
```bash
xevon scan-url https://api.example.com/login \
  --method POST \
  --body '{"username":"admin","password":"test123"}' \
  -H "Content-Type: application/json"
```

**Scan a raw HTTP request from a file:**
```
> I captured a raw request in request.txt, scan it
```
```bash
xevon scan-request -i request.txt
```

**Scan a raw request from stdin:**
```
> Scan this raw request from the terminal
```
```bash
echo -e "GET /api/users?id=1 HTTP/1.1\r\nHost: example.com\r\nAuthorization: Bearer tok123\r\n" | xevon scan-request
```

**Scan with a proxy (e.g., Burp Suite):**
```
> Scan https://example.com and route traffic through Burp
```
```bash
xevon scan -t https://example.com --proxy http://127.0.0.1:8080
```

**High-speed scan with tuned concurrency:**
```
> Scan fast — 100 workers, 200 req/s, max 5 per host
```
```bash
xevon scan -t https://example.com -c 100 --rate-limit 200 --max-per-host 5
```

**Scan and output results as JSONL:**
```
> Scan and save results as JSON lines
```
```bash
xevon scan -t https://example.com --format jsonl -o results.jsonl
```

**Scan and generate an HTML report:**
```
> Scan and produce an interactive HTML report
```
```bash
xevon scan -t https://example.com --format html -o report.html
```

**Scan with CI-friendly output (for CI/CD pipelines):**
```
> Scan and output only JSONL findings, no color or banners
```
```bash
xevon scan -t https://example.com --ci-output-format
```

**Scan with custom scanning profile:**
```
> Use the aggressive scanning profile
```
```bash
xevon scan -t https://example.com --scanning-profile aggressive
```

**Scan with strict origin scope:**
```
> Only scan URLs on the exact same origin
```
```bash
xevon scan -t https://example.com --scope-origin strict
```

---

### 2. Input Formats

**OpenAPI 3.x spec with explicit base URL:**
```
> Scan my API using the OpenAPI spec
```
```bash
xevon scan -I openapi -i api-spec.yaml -t https://api.example.com
```

**OpenAPI spec using servers from the spec:**
```
> Use the server URLs defined in the spec itself
```
```bash
xevon scan -I openapi -i api-spec.yaml --spec-url
```

**OpenAPI with auth header and parameter values:**
```
> Scan the spec with bearer auth and set the user_id parameter to 42
```
```bash
xevon scan -I openapi -i spec.yaml -t https://api.example.com \
  --spec-header "Authorization: Bearer eyJ..." \
  --spec-var "user_id=42"
```

**Swagger 2.0 spec:**
```
> Import a Swagger 2.0 spec and scan
```
```bash
xevon scan -I swagger -i swagger.json -t https://api.example.com
```

**Burp Suite XML export:**
```
> I exported traffic from Burp, scan it
```
```bash
xevon scan -I burp -i burp-export.xml -t https://example.com
```

**HAR (HTTP Archive) file:**
```
> Scan my browser-recorded HAR file
```
```bash
xevon scan -I har -i traffic.har
```

**cURL commands file:**
```
> I have a file of curl commands, scan them all
```
```bash
xevon scan -I curl -i curl-commands.txt
```

**Postman collection:**
```
> Import and scan my Postman collection
```
```bash
xevon scan -I postman -i collection.json -t https://api.example.com
```

**Nuclei templates:**
```
> Run these Nuclei templates against the target
```
```bash
xevon scan -I nuclei -i templates/ -t https://example.com
```

**Piped URLs from stdin:**
```
> Pipe a list of URLs into the scanner
```
```bash
cat urls.txt | xevon scan -i -
```

---

### 3. Phase Control

**Run only discovery (content enumeration):**
```
> Just run content discovery against the target
```
```bash
xevon run discover -t https://example.com
# or
xevon scan -t https://example.com --only discovery
```

**Run only spidering (headless browser crawling):**
```
> Spider the target with a headless browser
```
```bash
xevon run spidering -t https://example.com
```

**Run only audit (vulnerability scanning):**
```
> Skip discovery, just run the vulnerability modules
```
```bash
xevon run dynamic-assessment -t https://example.com
# or
xevon scan -t https://example.com --only dynamic-assessment
```

**Run only SPA (security posture assessment via Nuclei):**
```
> Run Nuclei-based security posture assessment, only critical and high severity
```
```bash
xevon run spa -t https://example.com --spa-severities critical,high
```

**Run only SAST (static analysis on source code):**
```
> Run static analysis on my Go app, filter for gin-related rules
```
```bash
xevon run sast --source /path/to/app --rule gin
```

**Run only external harvest (Wayback, Common Crawl, OTX):**
```
> Gather URLs from external intelligence sources
```
```bash
xevon run external-harvest -t https://example.com
```

**Skip specific phases:**
```
> Scan but skip discovery and spidering
```
```bash
xevon scan -t https://example.com --skip discovery,spidering
```

**Run only JavaScript extensions:**
```
> Run only my custom extension, skip built-in modules
```
```bash
xevon scan -t https://example.com --only extension --ext ./custom-check.js
# or
xevon run ext -t https://example.com --ext ./custom-check.js
```

**Phase aliases reference:**

| Alias | Resolves To |
|-------|-------------|
| `deparos`, `discover` | `discovery` |
| `spitolas` | `spidering` |
| `dynamic-assessment` | `audit` |
| `ext` | `extension` |

---

### 4. Module Filtering

**List all available scanner modules:**
```
> Show me all scanner modules
```
```bash
xevon module ls
# or
xevon scan -M
```

**Filter modules by keyword:**
```
> Show me all XSS-related modules
```
```bash
xevon module ls xss
```

**List only active modules with verbose details:**
```
> Show active modules with descriptions
```
```bash
xevon module ls --type active -v
```

**Scan with specific modules only:**
```
> Only run reflected XSS and error-based SQL injection modules
```
```bash
xevon scan -t https://example.com -m xss-reflected,sqli-error
```

**Filter modules by tags (OR logic):**
```
> Scan with modules tagged 'spring' or 'injection'
```
```bash
xevon scan -t https://example.com --module-tag spring --module-tag injection
```

**Combine module IDs and tags:**
```
> Run sqli-error plus all XSS-tagged modules
```
```bash
xevon scan -t https://example.com -m sqli-error --module-tag xss
```

**Enable/disable modules persistently:**
```
> Disable all SQL injection modules, enable all XSS modules
```
```bash
xevon module disable sqli
xevon module enable xss
```

**Enable by exact module ID:**
```
> Enable only the reflected XSS module
```
```bash
xevon module enable active-xss-reflected --id
```

---

### 5. Server & Ingestion

**Start the API server (default port 9002):**
```
> Start the xevon server
```
```bash
xevon server
```

**Start server on custom port without auth:**
```
> Start the server on port 8443 with no authentication
```
```bash
xevon server --service-port 8443 --no-auth
```

**Start server with scan-on-receive (auto-scan ingested traffic):**
```
> Start the server and auto-scan every request that comes in
```
```bash
xevon server -t https://example.com --scan-on-receive
```

**Start server with transparent proxy for recording:**
```
> Start server with a recording proxy on port 8080
```
```bash
xevon server --ingest-proxy-port 8080
```

**High-concurrency server:**
```
> Start a high-throughput server with 200 workers
```
```bash
xevon server -c 200 --mem-buffer 50000
```

**Ingest an OpenAPI spec locally:**
```
> Import the spec into the database without scanning
```
```bash
xevon ingest -t https://api.example.com -I openapi -i spec.yaml
```

**Ingest and auto-scan:**
```
> Import the spec and scan immediately
```
```bash
xevon ingest -t https://api.example.com -I openapi -i spec.yaml -S
```

**Ingest Burp export:**
```
> Import Burp traffic into the database
```
```bash
xevon ingest -t https://example.com -I burp -i export.xml
```

**Remote ingest to a running server:**
```
> Send traffic to the xevon server running on localhost
```
```bash
xevon ingest -s http://localhost:9002 -I openapi -i spec.yaml
```

**Ingest without fetching responses:**
```
> Store request-only records, don't make network requests
```
```bash
xevon ingest -t https://example.com -I burp -i export.xml --disable-fetch-response
```

---

### 6. AI Agent Modes

#### Agent Query (Template-Based Code Review)

Template-based code review + endpoint discovery runs through the `query` subcommand. The parent `xevon agent` only supports `--list-templates` and `--list-agents`.

**Security code review:**
```
> Review my source code for security vulnerabilities
```
```bash
xevon agent query --prompt-template security-code-review --source ./src
```

**Endpoint discovery from source:**
```
> Find all API endpoints in my source code
```
```bash
xevon agent query --prompt-template endpoint-discovery --source ./src
```

**Review specific files only:**
```
> Review only auth.go and middleware.go for security issues
```
```bash
xevon agent query --prompt-template security-code-review --source ./src \
  --files "src/auth.go,src/middleware.go"
```

**Append extra instructions to a template:**
```
> Code review, but focus on authentication and authorization
```
```bash
xevon agent query --prompt-template security-code-review --source ./src \
  --append "Focus specifically on authentication and authorization vulnerabilities"
```

**Use a custom prompt file:**
```
> Run the agent with my own prompt template
```
```bash
xevon agent query --prompt-file custom-prompt.md --source ./src
```

**Select a specific agent backend:**
```
> Use the Claude backend for code review
```
```bash
xevon agent query --agent claude --prompt-template security-code-review --source ./src
```

**Dry-run to preview the rendered prompt:**
```
> Show me what prompt would be sent to the agent
```
```bash
xevon agent query --prompt-template security-code-review --source ./src --dry-run
```

**Save agent output to a file:**
```
> Save the review results to a JSON file
```
```bash
xevon agent query --prompt-template security-code-review --source ./src \
  --output review-results.json
```

**List available templates and backends:**
```
> What prompt templates and agent backends are available?
```
```bash
xevon agent --list-templates
xevon agent --list-agents
```

**List or inspect agent sessions:**
```
> Show all agent run sessions
```
```bash
xevon agent session
xevon agent session --id <session-id>
```

**Built-in templates include:**
- `security-code-review` — Comprehensive security review
- `injection-sinks` — Find injection sinks
- `auth-bypass` — Auth bypass vectors
- `secret-detection` — Hardcoded secrets
- `endpoint-discovery` — API endpoints from source
- `api-input-gen` — Generate test inputs
- `curl-command-gen` — Generate cURL commands
- `attack-surface-mapper` — Map attack surface
- `nextjs-security-audit` — Next.js security review
- `react-xss-audit` — React XSS audit
- `cors-csrf-review` — CORS/CSRF config audit

#### Agent Query (Freeform Prompt)

**Inline prompt:**
```
> Ask the agent to review code for vulnerabilities
```
```bash
xevon agent query 'review this code for SQL injection vulnerabilities'
```

**Named prompt flag:**
```
> Analyze the authentication flow
```
```bash
xevon agent query --prompt 'analyze the authentication flow for bypass vectors'
```

**Pipe prompt from stdin:**
```
> Pipe a prompt to the agent
```
```bash
echo "check for SSRF in the URL-fetching handler" | xevon agent query --stdin
```

**Custom prompt file with a specific backend:**
```
> Run a custom prompt file through Claude
```
```bash
xevon agent query --agent claude --prompt-file custom-prompt.md
```

**With extended timeout:**
```
> Run a comprehensive review with extra time
```
```bash
xevon agent query --max-duration 10m 'comprehensive security review of all handlers'
```

#### Agent Autopilot (Autonomous Scanning)

Autopilot runs a single autonomous operator session that drives the xevon CLI. When `--source` is set, archon-audit runs first to build whitebox context, then the operator takes over. Use `--intensity quick|balanced|deep` to bundle limits, timeout, archon mode, and browser into one flag. Default timeout is `6h`.

**Basic autonomous scan:**
```
> Let the AI autonomously scan the target
```
```bash
xevon agent autopilot -t https://example.com
```

**Natural-language prompt (target/source/focus auto-extracted):**
```
> Scan VAmPI source at ~/src/VAmPI running on localhost:3005
```
```bash
xevon agent autopilot "scan VAmPI source at ~/src/VAmPI on localhost:3005"
```

**With source code context and focus area:**
```
> Autonomous scan focused on auth bypass, with source code for context
```
```bash
xevon agent autopilot -t https://api.example.com --source ./src --focus "auth bypass"
```

**Intensity presets (quick/balanced/deep):**
```
> Fast autopilot for CI; deep autopilot for a pentest
```
```bash
xevon agent autopilot -t https://example.com --source ./src --intensity quick   # 30 cmds, 1h
xevon agent autopilot -t https://example.com --intensity deep                    # 300 cmds, 12h, browser
```

**Scan only a PR diff or recent commits:**
```
> Scan the last 3 commits / a PR branch
```
```bash
xevon agent autopilot -t https://example.com --source ./src --last-commits 3
xevon agent autopilot -t https://example.com --source ./src --diff main...feature-branch
```

**Browser-based auth preflight:**
```
> Log in via the browser before autopilot starts
```
```bash
xevon agent autopilot -t https://example.com --browser --credentials "admin/admin123"
xevon agent autopilot -t https://example.com --browser --auth-required \
  --browser-start-url https://example.com/login
```

**Disable or tune archon-audit (runs by default when --source is set):**
```bash
xevon agent autopilot -t https://example.com --source ./src --archon=off
xevon agent autopilot -t https://example.com --source ./src --archon deep
```

**Custom limits (explicit override of presets):**
```
> Limit the agent to 50 commands and 15 minutes
```
```bash
xevon agent autopilot -t https://example.com --max-commands 50 --timeout 15m
```

**Upload results to cloud storage:**
```bash
xevon agent autopilot -t https://example.com --source ./src --upload-results
```

**Preview the system prompt (dry run):**
```
> Show me what system prompt the autopilot agent would receive
```
```bash
xevon agent autopilot -t https://example.com --dry-run
```

**Use a different agent backend:**
```
> Run autopilot with Codex
```
```bash
xevon agent autopilot -t https://example.com --agent codex
```

**Resume a previous session or override the olium provider:**
```bash
xevon agent autopilot --resume ~/.xevon/agent-sessions/<uuid>
xevon agent autopilot -t https://example.com --provider anthropic-api-key
```

#### Agent Swarm (Full-Scope and Targeted)

Swarm is AI-guided scanning — targeted by default, full-scope with `--discover`. Intensity presets (`--intensity quick|balanced|deep`) bundle discovery, triage, code-audit, browser, duration, and iteration defaults. Explicit flags override. Default swarm duration is `12h`; set with `--max-duration`.

**Deep analysis of a single endpoint:**
```
> Deep scan the users API endpoint for vulnerabilities
```
```bash
xevon agent swarm -t https://example.com/api/users
```

**Natural-language prompt:**
```
> Scan source at ~/src/app running on localhost:3005
```
```bash
xevon agent swarm "scan source at ~/src/app on localhost:3005"
```

**From a curl command:**
```
> Analyze this curl command for vulnerabilities
```
```bash
xevon agent swarm --input "curl -X POST https://example.com/api/login -d '{\"user\":\"admin\"}'"
```

**Pipe raw HTTP from stdin (auto-detected):**
```
> Scan this raw HTTP request
```
```bash
echo -e "POST /api/search HTTP/1.1\r\nHost: example.com\r\n\r\nq=test" | xevon agent swarm
```

**Intensity presets:**
```
> Quick swarm for CI; deep swarm for a full pentest
```
```bash
xevon agent swarm -t https://example.com/api/users?id=1 --intensity quick
xevon agent swarm -t https://example.com --source ./src --intensity deep
```

**Focus on a specific vulnerability:**
```
> Focus on SQL injection in the users endpoint
```
```bash
xevon agent swarm -t https://example.com/api/users --vuln-type sqli
```

**Source-aware swarm (discovers routes from source):**
```
> Scan my app with source code context
```
```bash
xevon agent swarm -t http://localhost:3000 --source ~/projects/my-app
```

**Source-aware with specific files:**
```
> Analyze only the API routes and user model
```
```bash
xevon agent swarm -t http://localhost:8080 --source ./backend \
  --files src/routes/api.js,src/models/user.js
```

**Source analysis only (extract routes, no scan):**
```
> Just extract routes from my source code
```
```bash
xevon agent swarm -t http://localhost:3000 --source ./src --source-analysis-only
```

**Background archon-audit (runs in parallel, requires --source):**
```
> Scan and audit the codebase at the same time
```
```bash
xevon agent swarm -t http://localhost:3000 --source ./src --archon         # lite
xevon agent swarm -t http://localhost:3000 --source ./src --archon deep    # 10-phase
```

**Scan only changed code:**
```
> Focus the swarm on the last N commits or a PR
```
```bash
xevon agent swarm -t https://example.com --source ./src --last-commits 3
xevon agent swarm -t https://example.com --source ./src --diff main...feature-branch
```

**Browser-based auth capture:**
```
> Swarm needs to log in through the UI first
```
```bash
xevon agent swarm -t https://example.com --browser --browser-auth \
  --credentials "username=admin,password=secret"
```

**Upload results:**
```bash
xevon agent swarm -t https://example.com --source ./src --upload-results
```

**Custom instructions:**
```
> Focus the agent on GraphQL parsing vulnerabilities
```
```bash
xevon agent swarm -t https://example.com/graphql --instruction "Focus on GraphQL parsing"
```

**Triage + rescan loop:**
```bash
xevon agent swarm -t https://example.com/api/users --triage --max-iterations 5
```

**Preview prompts:**
```
> Show me what the swarm agent would do
```
```bash
xevon agent swarm -t https://example.com/api/users --dry-run
```

#### Agent Archon (Foreground Whitebox Audit)

Agent archon runs the embedded `archon-audit` harness against a source tree as a foreground audit. Findings are imported into the xevon database; raw output is streamed to the console and persisted to `{session}/runtime.log`.

**Quick 3-phase audit (secrets + SAST triage):**
```
> Run a fast archon audit on this repo
```
```bash
xevon agent archon --mode lite --source .
```

**Deep 10-phase audit:**
```
> Full archon audit of my codebase
```
```bash
xevon agent archon --mode deep --source .
```

**Audit a remote repo (auto-clones via source-aware storage):**
```bash
xevon agent archon --mode lite --source https://github.com/org/repo
```

**Balanced 6-phase audit (archon always uses Claude; route OpenAI workloads through `--driver=piolium` on `xevon agent audit`):**
```bash
xevon agent archon --mode balanced --source ~/code/myapp
```

**Construct PoCs for confirmed findings in a prior audit:**
```bash
xevon agent archon --mode confirm --source ./audit-with-findings
```

**Revisit an existing audit tree with fresh context:**
```bash
xevon agent archon --mode revisit --source ./prior-audit-tree
```

**Read-only progress check on an in-progress run:**
```bash
xevon agent archon --mode status --source ./in-progress-audit
```

Valid modes: `lite`, `balanced` (alias `scan`), `deep`, `revisit`, `confirm`, `merge`, `diff`, `status`, `mock`. Valid agents: `claude` (default), `codex`. Add `--no-stream` to suppress console output (log is still written).

---

### 7. Traffic & Results Browsing

**Browse all stored HTTP traffic:**
```
> Show me the HTTP traffic in the database
```
```bash
xevon traffic
```

**Fuzzy search traffic:**
```
> Show traffic related to login
```
```bash
xevon traffic login
```

**Tree view (hierarchical URL structure):**
```
> Show traffic as a directory tree
```
```bash
xevon traffic --tree
```

**Burp-style colored output:**
```
> Show traffic in Burp Suite style
```
```bash
xevon traffic --burp
```

**Filter by host, method, status:**
```
> Show POST and PUT requests to api.example.com that returned 200
```
```bash
xevon traffic --host api.example.com --method POST,PUT --status 200
```

**Filter by date range:**
```
> Show traffic from January 2024
```
```bash
xevon traffic --from 2024-01-01 --to 2024-01-31
```

**Search in request/response body:**
```
> Find traffic containing "password" in the body
```
```bash
xevon traffic --body password
```

**Search in headers:**
```
> Find traffic with JWT tokens in headers
```
```bash
xevon traffic --header "Bearer"
```

**Custom columns:**
```
> Show host, method, path, status, and auth columns
```
```bash
xevon traffic --columns HOST,METHOD,PATH,STATUS,AUTH
```

**Watch mode (auto-refresh):**
```
> Monitor traffic in real-time, refresh every 5 seconds
```
```bash
xevon traffic --watch 5s
```

**View raw HTTP request/response:**
```
> Show raw traffic for the last 5 records
```
```bash
xevon traffic --raw --limit 5
```

**Browse findings:**
```
> Show all vulnerability findings
```
```bash
xevon finding
```

**Load findings from a file or stdin:**
```
> Import findings from a JSONL file
```
```bash
xevon finding load -i findings.jsonl
# or from stdin
cat findings.jsonl | xevon finding load -i -
```

**Filter findings by severity:**
```
> Show only high and critical findings
```
```bash
xevon finding --severity high,critical
```

**Filter findings by module type or source:**
```
> Show only active module findings from audit
```
```bash
xevon finding --module-type active --finding-source audit
```

**View a specific finding in Burp-style format:**
```
> Show finding #42 with full HTTP details
```
```bash
xevon finding --id 42 --burp
```

**Custom finding columns:**
```
> Show findings with tags and confidence
```
```bash
xevon finding --columns ID,SEVERITY,MODULE,MATCHED_AT,TAGS,CONFIDENCE
```

**Search findings:**
```
> Find SQL injection findings
```
```bash
xevon finding --search "sql injection"
```

**Watch findings in real-time:**
```
> Monitor findings as they come in
```
```bash
xevon finding --watch 5s
```

**Replay stored traffic (re-send requests):**
```
> Replay login-related requests and compare responses
```
```bash
xevon traffic replay login
```

**Replay and replace stored responses:**
```
> Replay requests to api.example.com and update stored responses
```
```bash
xevon traffic replay --host api.example.com --in-replace
```

---

### 8. Data Management

**Database statistics:**
```
> Show me database stats
```
```bash
xevon db stats
```

**Detailed stats with host breakdown:**
```
> Show detailed stats broken down by host
```
```bash
xevon db stats --detailed
```

**Stats for a specific host:**
```
> Stats for example.com only
```
```bash
xevon db stats --host example.com
```

**Live-updating stats:**
```
> Watch database stats, refresh every 10 seconds
```
```bash
xevon db stats --watch 10s
```

**List database records with filters:**
```
> Show findings table, critical and high severity
```
```bash
xevon db ls --table findings --severity critical,high
```

**List available tables and columns:**
```
> What tables are in the database? What columns does findings have?
```
```bash
xevon db ls --list-tables
xevon db ls --list-columns --table findings
```

**Clean records by hostname:**
```
> Delete all records for old-target.com
```
```bash
xevon db clean --host old-target.com --force
```

**Clean old records with dry-run preview:**
```
> Preview what would be deleted before January 2024
```
```bash
xevon db clean --before 2024-01-01 --dry-run
```

**Clean only findings (keep HTTP records):**
```
> Delete info-severity findings but keep the HTTP records
```
```bash
xevon db clean --findings-only --severity info --force
```

**Clean orphaned findings:**
```
> Remove findings without associated HTTP records
```
```bash
xevon db clean --orphans
```

**Reset entire database:**
```
> Wipe the entire database and start fresh
```
```bash
xevon db clean --force
```

**Reclaim disk space after deletion:**
```
> Vacuum the database to reclaim space
```
```bash
xevon db clean --vacuum
```

**Seed database with sample data (for development/testing):**
```
> Populate the database with sample data
```
```bash
xevon db seed
```

**View the raw runtime log for a past scan or agent session:**
```
> Show me the log for run 550e8400-...
```
```bash
xevon log 550e8400-e29b-41d4-a716-446655440000
xevon log <uuid> -f              # follow like `tail -f`
xevon log <uuid> --tail 500      # last 500 lines
xevon log <uuid> --full          # whole log
xevon log <uuid> --strip-ansi    # plain text, safe for grep/pipe
```

**List all scan + agent sessions:**
```
> Show every scan and agent session, with log availability
```
```bash
xevon log ls
xevon log --tui                   # interactive picker
```

**Import an archon audit folder or JSONL export:**
```
> Import these scan results into the database
```
```bash
xevon import /path/to/archon-output-harbor/   # archon audit folder
xevon import scan-results.jsonl                # JSONL from `xevon export --format jsonl`
```

---

### 9. Export & Reports

**Full JSONL export:**
```
> Export everything from the database as JSONL
```
```bash
xevon export --format jsonl -o full-export.jsonl
```

**Export only findings:**
```
> Export just the findings
```
```bash
xevon export --format jsonl --only findings -o findings.jsonl
```

**Export findings and HTTP records:**
```
> Export findings and associated HTTP traffic
```
```bash
xevon export --format jsonl --only findings,http -o results.jsonl
```

**HTML report:**
```
> Generate an interactive HTML report
```
```bash
xevon export --format html -o report.html
```

**Lite export (omit raw HTTP data):**
```
> Export URLs only, without raw request/response data
```
```bash
xevon export --lite --only http -o urls.jsonl
```

**Export with search filter:**
```
> Export only records matching example.com
```
```bash
xevon export --search "example.com" -o filtered.jsonl
```

**Database-level export as CSV:**
```
> Export HTTP records as CSV
```
```bash
xevon db export -f csv -o records.csv
```

**Export as Markdown:**
```
> Export records as a Markdown report
```
```bash
xevon db export -f markdown -o report.md
```

**Export raw requests only:**
```
> Export just the raw HTTP requests
```
```bash
xevon db export -f raw --request-only -o requests.txt
```

**Export filtered by host and date:**
```
> Export records for example.com from 2024 onwards
```
```bash
xevon db export -f csv -o records.csv --host example.com --from 2024-01-01
```

**Export a single record by UUID:**
```
> Export record abc12345
```
```bash
xevon db export --uuid abc12345
```

**Export module registry:**
```
> Export all available scanner modules
```
```bash
xevon export --only modules
```

---

### 10. Whitebox / Source-Aware Scanning

**Scan with local source code:**
```
> Whitebox scan with source code in ./src
```
```bash
xevon scan -t https://example.com --source ./src --strategy whitebox
```

**Scan with source cloned from Git:**
```
> Clone the repo and run a whitebox scan
```
```bash
xevon scan -t https://example.com \
  --source-url https://github.com/org/repo --strategy whitebox
```

**Link source code to a hostname first, then scan:**
```
> Link the source repo to example.com, then whitebox scan
```
```bash
xevon source add --hostname example.com --path ./src
xevon scan -t https://example.com --strategy whitebox
```

**Link source with metadata:**
```
> Link source with language and framework info
```
```bash
xevon source add --hostname api.example.com --path ./src -l go -f gin
```

**Link source from Git URL:**
```
> Clone and link a GitHub repo
```
```bash
xevon source add --hostname example.com --git https://github.com/org/repo
```

**List linked source repos:**
```
> Show all linked source repositories
```
```bash
xevon source ls
```

**Run SAST only:**
```
> Run static analysis on the source code
```
```bash
xevon run sast --source /path/to/app
```

**SAST with rule filtering:**
```
> Run SAST, only gin-related rules
```
```bash
xevon run sast --source /path/to/app --rule gin
```

**Ad-hoc SAST on a local path or Git URL:**
```
> Run static analysis on an ad-hoc path without linking source
```
```bash
xevon scan -t https://example.com --sast-adhoc /path/to/app
xevon scan -t https://example.com --sast-adhoc https://github.com/org/repo
```

---

### 11. JavaScript Extensions

**Install preset examples:**
```
> Install the example extension scripts
```
```bash
xevon ext preset
```

**View the extension API reference:**
```
> Show me the extension API docs
```
```bash
xevon ext docs
xevon ext docs --example          # with code examples
xevon ext docs http               # filter by namespace
```

**List loaded extensions:**
```
> Show currently loaded extensions
```
```bash
xevon ext ls
xevon ext ls --type active        # active extensions only
```

**Lint extensions for errors:**
```
> Validate my extension files for syntax errors and unknown API calls
```
```bash
xevon ext lint custom-check.js
xevon ext lint ./my-extensions/
```

**Quick-test JS code inline:**
```
> Test a JS expression
```
```bash
xevon ext eval 'xevon.log.info("hello from extension")'
xevon ext eval 'xevon.utils.md5("password")'
```

**Evaluate a JS file:**
```
> Run a JS script file
```
```bash
xevon ext eval --ext-file script.js
```

**Run a custom extension against a target:**
```
> Run my custom scanner extension
```
```bash
xevon run extension -t https://example.com --ext custom-check.js
# or
xevon run ext -t https://example.com --ext custom-check.js
```

**Run extension alongside built-in modules:**
```
> Run built-in modules plus my custom extension
```
```bash
xevon scan -t https://example.com --ext custom-check.js
```

**Run only extensions (skip built-in modules):**
```
> Run only my custom extensions
```
```bash
xevon scan -t https://example.com --only extension --ext custom-check.js
```

**Load multiple extensions:**
```
> Run three extensions together
```
```bash
xevon scan -t https://example.com --ext check1.js --ext check2.js --ext check3.js
```

**Load all extensions from a directory:**
```
> Run all extensions in my extensions folder
```
```bash
xevon scan -t https://example.com --ext-dir ./my-extensions/
```

**Ask the agent to write an extension:**
```
> Write me a passive extension that checks for missing security headers
```

The agent will generate a JS file like:

```javascript
module.exports = {
  id: "missing-security-headers",
  name: "Missing Security Headers",
  type: "passive",
  severity: "low",
  confidence: "certain",
  scope: "response",
  tags: ["headers", "misconfiguration", "light"],
  scanTypes: ["per_request"],

  scanPerRequest: function(ctx) {
    if (!ctx.response) return null;
    var headers = ctx.response.headers;
    var missing = [];

    if (!headers["strict-transport-security"]) missing.push("HSTS");
    if (!headers["x-content-type-options"]) missing.push("X-Content-Type-Options");
    if (!headers["x-frame-options"] && !headers["content-security-policy"]) {
      missing.push("X-Frame-Options/CSP");
    }

    if (missing.length === 0) return null;

    return {
      url: ctx.request.url,
      name: "Missing Security Headers: " + missing.join(", "),
      severity: "low",
      description: "Response is missing: " + missing.join(", ")
    };
  }
};
```

**Ask the agent to write an AI-augmented extension:**
```
> Write an active extension that uses AI to generate XSS payloads
```

The agent will generate a JS file using `xevon.agent.generatePayloads()` and `xevon.agent.analyzeResponse()`.

**Execute ad-hoc JavaScript with xevon js:**
```
> Run a JS script that queries the database and flags high-risk records
```
```bash
xevon js --code 'var records = xevon.db.records.query({ hostname: "example.com" }); records.forEach(function(r) { if (xevon.utils.hasDynamicSegment(xevon.parse.url(r.url).path)) xevon.db.records.annotate(r.uuid, { risk_score: 50 }); })'
```

**Execute a JS file with target context:**
```
> Run my scanner script against example.com
```
```bash
xevon js --target https://example.com --code-file my-scanner.js
```

**Quick hash or encode from the terminal:**
```
> MD5 hash the string "password"
```
```bash
xevon js --format text --code 'xevon.utils.md5("password")'
```

**YAML extension (simple pattern matching):**
```
> Write a YAML extension that detects stack traces and SQL errors
```

```yaml
id: error-pattern-detector
name: Verbose Error Pattern Detector
type: passive
severity: suspect
confidence: tentative
scope: response
tags: [error, information-disclosure, light]
scanTypes: [per_request]
patterns:
  - name: "Stack Trace Detected"
    regex: "(?:at\\s+[\\w.$]+\\(|Traceback \\(most recent|Exception in thread)"
    severity: suspect
  - name: "SQL Error Message"
    regex: "(?:mysql_|pg_|sqlite_|ORA-\\d{5}|SQLSTATE\\[)"
    severity: medium
```

---

### 12. Configuration & Projects

**First-time setup (create ~/.xevon with defaults):**
```
> Initialize xevon
```
```bash
xevon init
xevon init --force      # regenerate API key + re-extract preset data
```

**Reset the installation back to a clean state:**
```
> Wipe my xevon installation and start over
```
```bash
xevon config clean      # prompts for confirmation
xevon config clean -F   # skip prompt
```

**Run a health check on the installation:**
```
> Is my xevon install healthy?
```
```bash
xevon doctor
```

**View all configuration:**
```
> Show the current xevon config
```
```bash
xevon config ls
```

**View a specific config section:**
```
> Show scope configuration
```
```bash
xevon config ls scope
xevon config ls scanning_pace
xevon config ls server
```

**Set configuration values:**
```
> Set the default strategy to deep
```
```bash
xevon config set scanning_strategy.default_strategy deep
```

**Set scope mode:**
```
> Set origin scope to strict
```
```bash
xevon config set scope.origin.mode strict
```

**Enable extensions globally:**
```
> Enable extensions in audit
```
```bash
xevon config set audit.extensions.enabled true
```

**View scope rules:**
```
> Show current scope rules
```
```bash
xevon scope view
xevon scope view host
```

**View scanning strategies:**
```
> Show available strategies and their phases
```
```bash
xevon strategy ls
```

**Create and manage projects:**
```
> Create a project, then switch to it
```
```bash
xevon project create my-project
xevon project list
xevon project use my-project
```

**Scope CLI operations to a project:**
```
> Scan within a specific project
```
```bash
xevon scan -t https://example.com --project-name my-project
```

**Project-scoped database access:**
```
> Show stats for my-project
```
```bash
XEVON_PROJECT=my-project xevon db stats
```

**Authentication management utilities:**
```
> List sessions, lint an auth file, load auth, or generate TOTP
```
```bash
xevon auth list
xevon auth lint <auth-file>
xevon auth load <auth-file>
xevon auth totp --secret <base32-secret>
```

**Authenticated scanning with session config:**
```
> Scan with authentication config
```
```bash
# Two flags: --auth-file for files/bare names, --auth for inline sessions:
xevon scan -t https://example.com --auth-file auth.yaml              # bundle file
xevon scan -t https://example.com --auth-file admin-session.yaml     # single-session file
xevon scan -t https://example.com --auth-file <session-name>         # bare name in session_dir
xevon scan -t https://example.com --auth "admin:Cookie:sid=abc"      # inline
```

---

## Natural Language Examples

These are examples of natural language prompts you can give to Claude Code or Codex with the skill installed. The agent will translate them into the correct xevon commands.

| You Say | Agent Runs |
|---------|------------|
| "Scan example.com" | `xevon scan -t https://example.com` |
| "Deep scan with spidering" | `xevon scan -t <url> --strategy deep` |
| "Import my Burp export and scan it" | `xevon scan -I burp -i export.xml` |
| "Scan my OpenAPI spec with auth" | `xevon scan -I openapi -i spec.yaml -t <url> --spec-header "Authorization: Bearer ..."` |
| "Only run XSS modules" | `xevon scan -t <url> --module-tag xss` |
| "Review my code for security issues" | `xevon agent query --prompt-template security-code-review --source ./src` |
| "Autonomous scan focused on injection" | `xevon agent autopilot -t <url> --focus "injection"` |
| "Run the full AI scan" | `xevon agent swarm --discover -t <url>` |
| "Deep scan this endpoint for SQLi" | `xevon agent swarm -t <url> --vuln-type sqli` |
| "Scan with source code context" | `xevon agent swarm -t <url> --source ./src` |
| "Run this JS script against the API" | `xevon js --code-file script.js --target <url>` |
| "MD5 hash this string" | `xevon js --format text --code 'xevon.utils.md5("...")'` |
| "Show me all critical findings" | `xevon finding --severity critical` |
| "Import findings from a file" | `xevon finding load -i findings.jsonl` |
| "Show active module findings in Burp format" | `xevon finding --module-type active --burp` |
| "Run scan for CI/CD pipeline" | `xevon scan -t <url> --ci-output-format` |
| "Export results as HTML report" | `xevon export --format html -o report.html` |
| "What traffic is in the database?" | `xevon traffic` |
| "Write me an extension that checks for exposed .env files" | Generates a JS extension file |
| "Lint my extension for errors" | `xevon ext lint custom-check.js` |
| "Start the server with auto-scan" | `xevon server -t <url> --scan-on-receive` |
| "Whitebox scan with my source code" | `xevon scan -t <url> --source ./src --strategy whitebox` |
| "Ad-hoc SAST on a local path" | `xevon scan -t <url> --sast-adhoc /path/to/app` |
| "Clean up old scan data" | `xevon db clean --before <date> --force` |
| "Seed the database with sample data" | `xevon db seed` |
| "List agent sessions" | `xevon agent session` |
| "Show auth utilities" | `xevon auth list` |
| "Scan with auth file" | `xevon scan -t <url> --auth-file auth.yaml` |
| "Scan VAmPI source on localhost:3005" (natural-language) | `xevon agent autopilot "scan VAmPI source at ~/src/VAmPI on localhost:3005"` |
| "Fast archon audit of this repo" | `xevon agent archon --mode lite --source .` |
| "Deep archon audit of a GitHub repo" | `xevon agent archon --mode deep --source https://github.com/org/repo` |
| "Quick swarm for CI" | `xevon agent swarm -t <url> --intensity quick` |
| "Deep autopilot pentest" | `xevon agent autopilot -t <url> --intensity deep` |
| "Scan the last 3 commits" | `xevon agent autopilot -t <url> --source ./src --last-commits 3` |
| "Scan this PR for regressions" | `xevon agent autopilot -t <url> --source ./src --diff main...feature-branch` |
| "Tail the live log for that scan" | `xevon log <uuid> -f` |
| "List all scan + agent sessions" | `xevon log ls` |
| "Import archon audit results" | `xevon import /path/to/archon-output/` |
| "Import a JSONL scan export" | `xevon import results.jsonl` |
| "Initialize xevon" | `xevon init` |
| "Reset xevon to a clean install" | `xevon config clean` |
| "Health-check my install" | `xevon doctor` |

---

## Tips & Best Practices

1. **Start with `scan -t`** — It's the most common command. Add flags incrementally.
2. **Use strategies** — `lite` for quick checks, `balanced` for most cases, `deep` for full coverage, `whitebox` when you have source code.
3. **Phase isolation** — Use `--only` or `xevon run <phase>` to iterate on a single phase without re-running the entire pipeline.
4. **Module tags** — Filter modules by technology (`spring`, `nodejs`) or vulnerability class (`xss`, `injection`) to reduce noise.
5. **Watch mode** — Add `--watch 5s` to `traffic`, `finding`, or `db stats` for real-time monitoring during long scans.
6. **Dry-run agents** — Always `--dry-run` first for agent commands to preview prompts before spending AI tokens.
7. **Swarm over autopilot** — Use `agent swarm --discover` for structured scans (lower cost, reproducible). Use `agent autopilot` for exploratory, creative scanning.
8. **Extensions for custom logic** — Write JS extensions instead of modifying core modules. They run alongside built-in modules with `--ext`.
9. **Projects for isolation** — Use `xevon project create` to keep scan data separate across engagements.
10. **Export early** — Run `xevon export --format html -o report.html` to share results as interactive reports.
11. **Intensity presets** — Use `--intensity quick` for CI/PR pipelines, `--intensity balanced` (default) for day-to-day, `--intensity deep` for full pentests. Explicit flags override the preset.
12. **Archon for whitebox** — Reach for `agent archon --mode lite/balanced/deep` when you have source code and want a structured multi-phase audit with PoC construction (`--mode confirm`) and incremental runs (`--mode diff`).
13. **Follow sessions live** — Long-running agent runs leave a `runtime.log`; tail it with `xevon log <uuid> -f` from another terminal instead of attaching to the process.
14. **Start fresh between engagements** — `xevon config clean` wipes `~/.xevon/` and `xevon init` reinitializes it with a new API key and preset data.

## Resources

- **Website**: [xevon.live](https://xevon.live/)
- **Documentation**: [docs.xevon.live](https://docs.xevon.live/)
- **GitHub**: [github.com/xevonlive-dev/xevon](https://github.com/xevonlive-dev/xevon)
- **Skills Repository**: [github.com/xevonlive-dev/skills](https://github.com/xevonlive-dev/skills)
