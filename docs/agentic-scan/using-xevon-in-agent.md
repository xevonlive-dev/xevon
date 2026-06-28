# Using the xevon Scanner Skill in Claude Code & Codex

This guide explains how to install and use the `xevon-scanner` skill (`public/skills/xevon-scanner/`) with AI coding agents — Claude Code and OpenAI Codex — to operate the xevon CLI for web vulnerability scanning, security testing, and custom extension authoring.

---

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
  - [JavaScript Extensions](#10-javascript-extensions)
  - [Configuration & Projects](#11-configuration--projects)
- [Natural Language Examples](#natural-language-examples)
- [Tips & Best Practices](#tips--best-practices)

---

## What the Skill Does

The skill teaches the AI agent how to:

1. **Pick the right xevon command** for any security testing task
2. **Construct correct flag combinations** with proper syntax
3. **Follow scanning workflows** end-to-end (ingest → scan → triage → export)
4. **Write custom JavaScript extensions** using the `xevon.*` API
5. **Operate AI agent modes** (query, autopilot, swarm)
6. **Manage data** — browse traffic, filter findings, export reports, clean databases

The skill uses lazy-loaded references: the main `SKILL.md` stays small, and detailed docs are loaded on demand when the agent needs deep flag information or extension authoring guidance.

---

## Skill Structure

```
public/skills/xevon-scanner/
├── SKILL.md                              # Main skill (decision tree, recipes, flags)
└── references/
    ├── scanning-commands.md              # scan, scan-url, scan-request, run
    ├── server-and-ingestion.md           # server, ingest, traffic, traffic replay
    ├── agent-commands.md                 # agent, agent query, autopilot, swarm
    ├── data-and-management.md            # db, module, ext, config, scope, export
    ├── flags-reference.md                # Complete alphabetical flag index
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

**Run only dynamic-assessment (vulnerability scanning):**
```
> Skip discovery, just run the vulnerability modules
```
```bash
xevon run dynamic-assessment -t https://example.com
# or
xevon scan -t https://example.com --only dynamic-assessment
```

**Run only known-issue-scan (Nuclei templates + Kingfisher secrets):**
```
> Run known-issue scan, only critical and high severity
```
```bash
xevon run known-issue-scan -t https://example.com --known-issue-scan-severities critical,high
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

#### Agent (Template-Based)

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

**List available templates:**
```
> What prompt templates are available?
```
```bash
xevon agent --list-templates
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

**With extended timeout:**
```
> Run a comprehensive review with extra time
```
```bash
xevon agent query --max-duration 10m 'comprehensive security review of all handlers'
```

#### Agent Autopilot (Autonomous Scanning)

**Basic autonomous scan:**
```
> Let the AI autonomously scan the target
```
```bash
xevon agent autopilot -t https://example.com
```

**With source code context and focus area:**
```
> Autonomous scan focused on auth bypass, with source code for context
```
```bash
xevon agent autopilot -t https://api.example.com --source ./src --focus "auth bypass"
```

**Custom limits (fewer commands, shorter timeout):**
```
> Limit the agent to 50 commands and 15 minutes
```
```bash
xevon agent autopilot -t https://example.com --max-commands 50 --timeout 15m
```

**Preview the system prompt (dry run):**
```
> Show me what system prompt the autopilot agent would receive
```
```bash
xevon agent autopilot -t https://example.com --dry-run
```

**Custom system prompt:**
```
> Use my own system prompt for autopilot
```
```bash
xevon agent autopilot -t https://example.com --system-prompt my-system-prompt.md
```

**Use a different provider per invocation:**
```
> Run autopilot using the anthropic-api-key provider
```
```bash
xevon agent autopilot -t https://example.com --provider anthropic-api-key --llm-api-key "$ANTHROPIC_API_KEY"
```

**Autopilot tool surface:**
- `bash` (catastrophic patterns hard-rejected: `rm -rf /`, fork bombs, `dd` to block devices, `mkfs` against real devices), `read_file`, `write_file`, `edit_file`, `ls`, `grep`, `glob`, `web_fetch`
- Plus autopilot-only tools: `halt_scan`, `report_finding`, `load_skill` (when skills are loaded), and vigtool (`run_scan`, `list_findings`, etc.) when a DB repository is attached
- Per-tool timeout: 5 minutes (configurable via `agent.olium.call_timeout_sec`)
- Hard cap: 200 findings (soft warning at 50); turn cap from `--intensity`: quick 150, balanced 500 (default), deep 1500

#### Agent Swarm (AI-Guided Multi-Phase Scan)

**Basic swarm scan with discovery (all phases):**
```
> Run the full AI swarm scan
```
```bash
xevon agent swarm --discover -t https://example.com
```

The swarm runs:
1. **Discover** — Native content discovery + spidering (no AI)
2. **Plan** — AI analyzes discovery results, produces an attack plan
3. **Scan** — Native executor with agent-selected modules (no AI)
4. **Triage** — AI reviews findings, confirms/dismisses, suggests follow-ups
5. **Rescan** — Targeted re-scanning from triage recommendations (no AI)
6. **Report** — Structured output from database (no AI)

**Swarm with focus area and source code:**
```
> Swarm scan focused on SQL injection, with source code
```
```bash
xevon agent swarm --discover -t https://example.com --focus "SQL injection" --source ./src
```

**Control rescan iterations:**
```
> Allow up to 3 triage→rescan iterations
```
```bash
xevon agent swarm --discover -t https://example.com --max-iterations 3
```

**Skip discovery and start from planning (use existing DB data):**
```
> I already have traffic in the database, start from planning
```
```bash
xevon agent swarm --discover -t https://example.com --skip discover --start-from plan
```

**Skip triage (just discover → plan → scan):**
```
> Run swarm but skip triage and rescan
```
```bash
xevon agent swarm --discover -t https://example.com --skip triage --skip rescan
```

**Use a scanning profile:**
```
> Run swarm with the deep scanning profile
```
```bash
xevon agent swarm --discover -t https://example.com --profile deep
```

**Preview agent prompts (dry run):**
```
> Show me the prompts without executing
```
```bash
xevon agent swarm --discover -t https://example.com --dry-run
```

**Specific source files for agent context:**
```
> Only include routes.go and handlers.go as context
```
```bash
xevon agent swarm --discover -t https://example.com --source ./src \
  --files "routes.go,handlers.go"
```

**Use a different provider per invocation:**
```
> Run swarm using anthropic-api-key
```
```bash
xevon agent swarm --discover -t https://example.com --provider anthropic-api-key --llm-api-key "$ANTHROPIC_API_KEY"
```

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

**Filter findings by severity:**
```
> Show only high and critical findings
```
```bash
xevon finding --severity high,critical
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

**Slim export (omit raw HTTP request/response bytes, keep metadata):**
```
> Export HTTP records without the raw request/response bodies
```
```bash
xevon export --omit-response --only http -o urls.jsonl
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

### 10. JavaScript Extensions

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

### 11. Configuration & Projects

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
> Enable extensions in dynamic-assessment
```
```bash
xevon config set dynamic-assessment.extensions.enabled true
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
xevon strategy
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
| "Run the full AI pipeline" | `xevon agent swarm --discover -t <url>` |
| "Show me all critical findings" | `xevon finding --severity critical` |
| "Export results as HTML report" | `xevon export --format html -o report.html` |
| "What traffic is in the database?" | `xevon traffic` |
| "Write me an extension that checks for exposed .env files" | Generates a JS extension file |
| "Start the server with auto-scan" | `xevon server -t <url> --scan-on-receive` |
| "Source-aware AI scan" | `xevon agent swarm -t <url> --source ./src` |
| "Multi-phase code audit" | `xevon agent audit --source ./src` |
| "Clean up old scan data" | `xevon db clean --before <date> --force` |

---

## Tips & Best Practices

1. **Start with `scan -t`** — It's the most common command. Add flags incrementally.
2. **Use strategies** — `lite` for quick checks, `balanced` for most cases, `deep` for full coverage. When you have source code, reach for an agent mode (`swarm`, `autopilot`, `audit`, or `query`) instead.
3. **Phase isolation** — Use `--only` or `xevon run <phase>` to iterate on a single phase without re-running the entire pipeline.
4. **Module tags** — Filter modules by technology (`spring`, `nodejs`) or vulnerability class (`xss`, `injection`) to reduce noise.
5. **Watch mode** — Add `--watch 5s` to `traffic`, `finding`, or `db stats` for real-time monitoring during long scans.
6. **Dry-run agents** — Always `--dry-run` first for agent commands to preview prompts before spending AI tokens.
7. **Swarm over autopilot** — Use `agent swarm --discover` for structured scans (lower cost, reproducible). Use `agent autopilot` for exploratory, creative scanning.
8. **Extensions for custom logic** — Write JS extensions instead of modifying core modules. They run alongside built-in modules with `--ext`.
9. **Projects for isolation** — Use `xevon project create` to keep scan data separate across engagements.
10. **Export early** — Run `xevon export --format html -o report.html` to share results as interactive reports.
