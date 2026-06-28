# Data & Management Commands Reference

Complete reference for `init`, `import`, `log`, `doctor`, `db`, `finding`, `module`, `extensions` (alias `ext`), `js`, `config`, `scope`, `source`, `auth`, `project`, `storage`, `strategy`, `export`, and `version` commands.

## Table of Contents

- [init](#init)
- [import](#import)
- [log](#log)
- [doctor](#doctor)
- [db](#db)
- [db list / ls](#db-list)
- [db stats](#db-stats)
- [db export](#db-export)
- [db clean](#db-clean)
- [db seed](#db-seed)
- [finding](#finding)
- [finding load](#finding-load)
- [export (top-level)](#export)
- [module](#module)
- [extensions](#extensions)
- [js](#js)
- [config](#config)
- [config clean](#config-clean)
- [scope](#scope)
- [source](#source)
- [auth](#auth)
- [project](#project)
- [storage](#storage)
- [storage ls](#storage-ls)
- [storage upload](#storage-upload)
- [storage download](#storage-download)
- [storage results](#storage-results)
- [storage presign](#storage-presign)
- [storage rm](#storage-rm)
- [strategy](#strategy)
- [version](#version)

---

## init

**Usage:** `xevon init [--force]`

Create `~/.xevon/` with a default config file (including a freshly generated API key), database schema, default scanning profiles, prompt templates, JavaScript extensions, and SAST rules. Idempotent: skips components that already exist unless `--force` is passed.

With `--force`, the existing config is rewritten with a new API key and the preset directories (`profiles/`, `extensions/`, `prompts/`) are re-extracted from the embedded assets.

```bash
xevon init
xevon init --force   # regenerate API key + re-extract preset data
```

Output (on success):

```
✓ xevon initialized successfully!
  ℹ Config:   ~/.xevon/xevon-configs.yaml
  ℹ Database: ~/.xevon/database-xevon.sqlite
  ℹ Docs:     https://docs.xevon.live
```

---

## import

**Usage:** `xevon import <path>`

Import scan data into the current project's database. Two input formats are supported (auto-detected by path type):

- **Audit output folder** (directory): produced by `xevon agent audit` (xevon-audit or piolium leg) — contains `audit-state.json` and `findings-draft/`. Creates a new `agentic_scan` row plus all findings; severity breakdown is printed on completion.
- **JSONL file** (regular file): each line is a JSON object wrapped in an envelope like `{"type": "http_record", "data": {...}}` or `{"type": "finding", "data": {...}}`. This matches the output of `xevon export --format jsonl`. Records are saved via `SaveRecordsBatch` (batch size 500); findings are deduplicated on save.

```bash
xevon import /path/to/xevon-results/         # audit output folder
xevon import scan-results.jsonl                 # JSONL export
xevon import /tmp/demo/juice-shop.jsonl

# After an interactive audit (`agent audit -i`), turn the on-disk results into a report:
xevon import ./src/xevon-results --format html -o audit-report.html
```

Notes:
- Imported findings inherit the current project's UUID and default `finding_source = "import"` when the field is empty.
- Unknown envelope types are counted and reported at the end (e.g. for forward-compatibility).
- Use `--project-id` / `--project-name` (or `XEVON_PROJECT`) to target a specific project.

---

## log

**Usage:** `xevon log [uuid] [flags]`

View raw `runtime.log` for a native scan or agentic scan session. When called without a UUID, behaves like `log ls` and lists all sessions.

### Log resolution order

1. Agentic session file: `~/.xevon/agent-sessions/<uuid>/runtime.log`
2. Native session file: `~/.xevon/native-sessions/<uuid>/runtime.log`
3. DB fallback: `scan_logs` table (used when `scanning_strategy.scan_logs.persist_logs` is disabled)

The legacy `run.log` filename is also resolved for older sessions. Agent audit child runs whose UUID does not match a session directory fall back to their parent's `SessionDir` column.

### log flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--tail` | `-n` | int | `200` | Show the last N lines (0 = none, -1 = all) |
| `--full` | — | bool | `false` | Show the full log (shortcut for `--tail -1`) |
| `--follow` | `-f` | bool | `false` | Follow log output as it is written (tail -f). Auto-enabled when the session is still running, unless `--follow=false` is set explicitly |
| `--strip-ansi` | — | bool | `false` | Strip ANSI color codes from output |
| `--tui` / `--no-tui` | — | bool | — | Enable / force-disable interactive picker (affects `log ls` behaviour) |

### log ls

`xevon log ls` prints a merged table of native + agentic sessions (kind, UUID, status, target, log availability, size, creation time). Status is color-coded. A tip at the bottom shows how to tail a specific UUID.

### Examples

```bash
# Table of sessions
xevon log ls
xevon log                               # same as `log ls`

# Interactive picker
xevon log --tui

# Stream a session's log, auto-following if it's still running
xevon log 550e8400-e29b-41d4-a716-446655440000

# Tail the last 500 lines
xevon log <uuid> --tail 500

# Full log, no follow, strip ANSI for grep
xevon log <uuid> --full --strip-ansi | grep -i sqli

# Force follow
xevon log <uuid> -f
```

---

## doctor

**Usage:** `xevon doctor`

Run a health check on the installation: verifies the config and database paths, external binaries (`claude`, `codex`, `nuclei`, `kingfisher`), directory permissions, and extracted preset data. Use this after `xevon init` or when an `agent` backend refuses to launch.

```bash
xevon doctor
```

---

## db

**Usage:** `xevon db <subcommand> [flags]`

Manage database records. Parent command for `clean`, `export`, `list` (`ls`), `seed`, and `stats`.

### Shared db flags (persistent)

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--table` | string | — | Database table: http_records, findings, scans |
| `--search` | string | — | Quick search across record fields |

---

## db list

**Usage:** `xevon db list [flags]` (aliases: `ls`)

List database records with filtering, sorting, and display options.

### Display flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--tree` | bool | `false` | Hierarchical tree format |
| `--raw` | bool | `false` | Full raw HTTP request and response |
| `--list-tables` | bool | `false` | List all database table names |
| `--list-columns` | bool | `false` | List column names for the current table |

### Pagination flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--limit` | `-n` | int | `100` | Max records to display |
| `--offset` | `-o` | int | `0` | Records to skip |

### Column selection flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--columns` | []string | — | Columns to include |
| `--exclude-columns` | []string | — | Columns to exclude |

### Filter flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--host` | string | — | Filter by hostname pattern (wildcard supported) |
| `--method` | []string | — | Filter by HTTP method |
| `--status` | []int | — | Filter by HTTP status code |
| `--path` | string | — | Filter by URL path pattern |
| `--scan-id` | string | — | Filter by scan session ID |
| `--severity` | string | — | Filter findings by severity |
| `--min-risk` | int | `0` | Show only records with risk score at or above this value |
| `--remark` | string | — | Filter records containing this text in remarks |
| `--module-type` | string | — | Filter findings by module type (active, passive, nuclei, secret-scan, agent, source-tools, oast, extension) |
| `--finding-source` | string | — | Filter findings by source (audit, spa, agent, oast, source-tools, extension) |
| `--from` | string | — | Records after date (YYYY-MM-DD or RFC3339) |
| `--to` | string | — | Records before date (YYYY-MM-DD or RFC3339) |
| `--header` | string | — | Search within HTTP header names and values |
| `--body` | string | — | Search in request/response body |

### Sorting flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--sort` | string | `created_at` | Sort field: uuid, created_at, sent_at, method, status_code, response_time |
| `--asc` | bool | `false` | Sort ascending |

### Examples

```bash
xevon db ls
xevon db ls --table findings
xevon db ls --table scans
xevon db ls --table findings --severity critical,high
xevon db ls --host example.com --method POST --status 200
xevon db ls --list-tables
xevon db ls --list-columns --table findings
xevon db ls --tree
xevon db ls --raw --limit 5
```

---

## db stats

**Usage:** `xevon db stats [flags]`

Show database statistics including record counts, finding breakdowns, and host summaries.

### stats-specific flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--detailed` | bool | `false` | Show per-host and per-module breakdown |
| `--scan-id` | string | — | Stats for a specific scan session |
| `--host` | string | — | Stats for a specific hostname |

### Examples

```bash
xevon db stats
xevon db stats --detailed
xevon db stats --host example.com
xevon db stats --watch 10s
```

---

## db export

**Usage:** `xevon db export [flags]`

Export database records in various formats.

### export-specific flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--format` | `-f` | string | `jsonl` | Export format: jsonl, json, raw, csv, markdown, markdown-table |
| `--output` | `-o` | string | stdout | Output file path |
| `--host` | — | string | — | Filter by hostname pattern |
| `--method` | — | []string | — | Filter by HTTP method |
| `--status` | — | []int | — | Filter by status code |
| `--path` | — | string | — | Filter by URL path pattern |
| `--scan-id` | — | string | — | Filter by scan session ID |
| `--severity` | — | string | — | Filter by severity level |
| `--from` | — | string | — | Export records created after this date (YYYY-MM-DD) |
| `--to` | — | string | — | Export records created before this date (YYYY-MM-DD) |
| `--limit` | — | int | `0` (unlimited) | Max records to export |
| `--offset` | — | int | `0` | Records to skip |
| `--uuid` | — | string | — | Export single record by UUID |
| `--request-only` | — | bool | `false` | Export only HTTP requests, omitting responses (raw format only) |

### Examples

```bash
xevon db export -f jsonl -o records.jsonl
xevon db export -f csv -o records.csv --host example.com
xevon db export -f markdown -o report.md
xevon db export -f raw --request-only -o requests.txt
xevon db export --uuid abc12345
```

---

## db clean

**Usage:** `xevon db clean [flags]`

Delete database records with filtering. Destructive operations require `--force`.

### clean-specific flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | `false` | Delete all records (requires `--force`) |
| `--host` | string | — | Delete records matching hostname |
| `--scan-id` | string | — | Delete records by scan session |
| `--before` | string | — | Delete records before date (YYYY-MM-DD) |
| `--status` | []int | — | Delete by HTTP status code |
| `--severity` | string | — | Delete findings by severity |
| `--dry-run` | bool | `false` | Show what would be deleted without deleting |
| `--vacuum` | bool | `false` | Reclaim disk space after deletion (SQLite) |
| `--orphans` | bool | `false` | Delete findings with no matching HTTP record |
| `--findings-only` | bool | `false` | Delete findings only, keep HTTP records |

### Special behavior

- `--force` with no filter flags: resets the entire SQLite database (deletes file + recreates)
- `--all` without `--force`: error
- Without `--force`: interactive confirmation prompt

### Examples

```bash
xevon db clean --scan-id my-scan
xevon db clean --host old-target.com --force
xevon db clean --before 2024-01-01 --dry-run
xevon db clean --all --force
xevon db clean --orphans
xevon db clean --findings-only --severity info
xevon db clean --vacuum
xevon db clean --force  # reset entire database
```

---

## db seed

**Usage:** `xevon db seed [flags]`

Populate database with sample data for development and testing.

### Examples

```bash
xevon db seed
```

---

## finding

**Usage:** `xevon finding [search-term] [flags]` (aliases: `findings`)

Browse vulnerability findings with fuzzy search, filtering, raw display, and column selection.

### Finding-specific filter flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--severity` | — | string | — | Filter by severity (comma-separated: critical,high,medium,low,info) |
| `--scan-id` | — | string | — | Filter by scan session ID |
| `--module-type` | — | string | — | Filter by module type (active, passive, nuclei, secret-scan, agent, source-tools, oast, extension) |
| `--finding-source` | — | string | — | Filter by finding source (audit, spa, agent, oast, source-tools, extension) |
| `--id` | — | int | `0` | Filter by finding ID |

### Display flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--raw` | bool | `false` | Show full raw HTTP request and response for each finding |
| `--burp` | bool | `false` | Display in Burp Suite-style format (colored request/response) |
| `--columns` | []string | — | Columns to show (comma-separated, e.g. ID,SEVERITY,MODULE) |
| `--exclude-columns` | []string | — | Columns to hide (comma-separated) |

### Pagination and sorting flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--limit` | `-n` | int | `100` | Maximum findings to display |
| `--offset` | `-o` | int | `0` | Number of findings to skip (for pagination) |
| `--sort` | — | string | `found_at` | Sort by: found_at, created_at, severity, module, confidence |
| `--asc` | — | bool | `false` | Sort in ascending order |

### Additional filter flags

Also accepts: `--host`, `--method`, `--status`, `--path`, `--from`, `--to`, `--search`, `--header`, `--body`, `--source`.

### Available columns

ID, SEVERITY, CONFIDENCE, MODULE, MODULE_ID, SHORT_DESC, DESCRIPTION, TYPE, SOURCE, MATCHED_AT, FOUND_AT, SCAN_UUID, TAGS

Default columns: ID, SEVERITY, MODULE, SHORT_DESC, TYPE, SOURCE, MATCHED_AT

### Examples

```bash
xevon finding
xevon finding --severity high,critical
xevon finding --search "sql injection"
xevon finding --module-type active
xevon finding --finding-source audit
xevon finding --id 42
xevon finding --burp
xevon finding --raw
xevon finding --columns ID,SEVERITY,MODULE,MATCHED_AT,TAGS
xevon finding --sort severity --asc
xevon finding --watch 5s
```

---

## finding load

**Usage:** `xevon finding load [file] [flags]`

Import findings from a file or stdin.

### Examples

```bash
xevon finding load findings.jsonl
cat findings.jsonl | xevon finding load
```

---

## export

**Usage:** `xevon export [flags]`

Top-level export command. Exports database tables and module registry as JSONL or HTML.

### export flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--format` | — | string | `jsonl` | Export format: html, jsonl |
| `--output` | `-o` | string | — | Output file (required for html) |
| `--only` | — | []string | all | Export only these tables (repeatable: http, findings, scans, modules, oast, source-repos, scopes) |
| `--omit-response` | — | bool | `false` | Omit raw HTTP request/response bytes from output (keeps metadata, smaller files) |
| `--search` | — | string | — | Fuzzy search filter across URLs, paths, hostnames, methods, content types, and sources |
| `--limit` | — | int | `0` (unlimited) | Max records per table |

### Examples

```bash
xevon export --format jsonl -o full-export.jsonl
xevon export --format jsonl --only findings
xevon export --format jsonl --only findings,http
xevon export --format html -o report.html
xevon export --only modules
xevon export --omit-response --only http -o urls.jsonl
xevon export --search "example.com" -o filtered.jsonl
```

---

## module

**Usage:** `xevon module [flags]` (aliases: `mo`)

Manage scanner modules. Lists active and passive modules with their scan scope, severity, and enabled status.

### Subcommands

| Command | Aliases | Description |
|---------|---------|-------------|
| `module ls [filter]` | `list` | List available modules (optional fuzzy filter) |
| `module enable <search>` | `e` | Enable modules matching search |
| `module disable <search>` | `d` | Disable modules matching search |

### module ls flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--type` | string | `all` | Filter: all, active, passive |
| `--list-enabled` | bool | `false` | Show only enabled modules |
| `--tags` | bool | `false` | Show only unique module tags |
| `--verbose` / `-v` | bool | `false` | Show long description and confirmation criteria |

### module enable/disable flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--id` | bool | `false` | Match exact module ID instead of fuzzy |

### Examples

```bash
xevon module ls
xevon module ls xss                 # fuzzy filter
xevon module ls --type active
xevon module ls --list-enabled
xevon module ls -v                  # verbose with descriptions
xevon module enable xss             # enable all xss modules
xevon module disable sqli           # disable all sqli modules
xevon module enable active-xss-reflected --id  # exact ID
xevon scan -M                       # shortcut to list modules
```

---

## extensions

**Usage:** `xevon extensions [filter]` (aliases: `ext`)

Manage JavaScript extensions for custom scanning logic.

### Subcommands

| Command | Aliases | Description |
|---------|---------|-------------|
| `ext docs [function]` | `doc`, `api` | Show API reference |
| `ext eval [code]` | `run`, `exec` | Evaluate JavaScript code with xevon.* APIs available |
| `ext lint [file]` | — | Validate extension files for syntax errors and unknown API calls |
| `ext ls [filter]` | `list` | List loaded extensions |
| `ext preset [name]` | `presets`, `init` | Install example presets |

### ext ls flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--type` | string | `all` | Filter: all, active, passive, pre_hook, post_hook |

### ext docs flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--example` | bool | `false` | Show usage examples for each function |

### ext lint flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--ext-file` | string | — | Path to extension file to validate |

### ext eval flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--stdin` | bool | `false` | Read JS code from stdin |
| `--ext-file` | string | — | Path to JS file to evaluate |

### Examples

```bash
xevon ext ls
xevon ext ls --type active
xevon ext docs
xevon ext docs http             # filter API docs by namespace
xevon ext docs --example        # with code examples
xevon ext preset                # install all presets
xevon ext preset my-scanner     # install specific preset
xevon ext eval 'xevon.log.info("hello")'
xevon ext eval --ext-file script.js
echo 'xevon.utils.md5("hello")' | xevon ext eval --stdin
```

---

## js

**Usage:** `xevon js [flags]`

Execute JavaScript code with access to the full `xevon.*` API surface. Reads from stdin by default, or use `--code` / `--code-file` for inline or file input. TypeScript files (`.ts`) are auto-transpiled.

### Input methods (mutually exclusive, in order of precedence)

1. `--code` — Inline JavaScript code
2. `--code-file` — Path to JavaScript/TypeScript file
3. stdin (default) — Read JS code from piped input

### js flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--code` | string | — | Inline JavaScript code to execute |
| `--code-file` | string | — | Path to JavaScript/TypeScript file (auto-transpiles `.ts`) |
| `--target` | string | — | Set TARGET variable in JS scope (URL string) |
| `--timeout` | duration | `30s` | Execution timeout (e.g., `60s`, `2m`) |
| `--format` | string | `json` | Output format: `json` or `text` |

### Available API

The JS VM provides access to all `xevon.*` namespaces:

| Namespace | Description |
|-----------|-------------|
| `xevon.http` | HTTP requests, sessions, batch, replay, sequence, auth testing, GraphQL, caching |
| `xevon.utils` | Encoding, hashing, diff, similarity, JWT, CSS selectors, multipart, file I/O |
| `xevon.parse` | URL, HTTP request/response, HTML, headers, cookies, query, JSON, form parsing |
| `xevon.scan` | Module listing, scope, finding creation, scan control |
| `xevon.db` | HTTP record and finding queries, annotations, comparison |
| `xevon.ingest` | URL, curl, raw HTTP, OpenAPI, Postman ingestion |
| `xevon.source` | Source code file listing, reading, searching |
| `xevon.agent` | AI-augmented analysis (ask, chat, complete, generatePayloads, analyzeResponse, confirmFinding) |
| `xevon.oast` | Out-of-band testing (enabled, payload, poll) |
| `xevon.log` | Logging (info, warn, error, debug) |
| `xevon.config` | Read-only config variables |
| `xevon.payloads(type)` | Built-in payload wordlists (xss, sqli, ssti, ssrf, lfi, etc.) |

### Return value

- Returns `undefined`/`null` → no output
- Otherwise → JSON-stringified return value on stdout
- With `--format text` → JSON strings are unquoted

### Examples

```bash
# Inline code
xevon js --code 'xevon.http.get("https://example.com/api/health")'

# From a file
xevon js --code-file scanner-script.js

# TypeScript auto-transpilation
xevon js --code-file scanner.ts

# From stdin (ideal for agent/pipe workflows)
echo 'xevon.utils.md5("password123")' | xevon js

# With target context (accessible as TARGET variable)
xevon js --target https://example.com --code 'xevon.http.get(TARGET + "/api/users")'

# Custom timeout and text output
xevon js --timeout 60s --format text --code 'xevon.utils.sha256("hello")'

# Query database records
xevon js --code 'JSON.stringify(xevon.db.records.query({ hostname: "example.com", limit: 5 }))'

# Ingest and scan
xevon js --code 'xevon.ingest.url("https://example.com/api/users"); xevon.scan.startNewScan({ targets: ["https://example.com"] })'

# Use AI to generate payloads
xevon js --code 'JSON.stringify(xevon.agent.generatePayloads({ type: "xss", context: "HTML attribute", count: 5 }))'
```

### Differences from `xevon ext eval`

| Feature | `xevon js` | `xevon ext eval` |
|---------|---------------|---------------------|
| Input methods | `--code`, `--code-file`, stdin | positional arg, `--ext-file`, `--stdin` |
| Target context | `--target` sets `TARGET` variable | Not available |
| Timeout | Configurable via `--timeout` | Not configurable |
| Output format | `--format json\|text` | Direct output |
| Use case | General scripting, automation | Quick extension testing |

---

## config

**Usage:** `xevon config <subcommand>`

Manage configuration settings.

### Subcommands

| Command | Aliases | Description |
|---------|---------|-------------|
| `config ls [filter]` | `list`, `view` | Display current configuration |
| `config set <key> <value>` | — | Set a configuration value |
| `config clean` | — | Clean/reset configuration |

### Examples

```bash
xevon config ls
xevon config ls scope           # filter by section
xevon config ls scanning_pace
xevon config ls server          # view server config
xevon config ls --force         # show sensitive values (unredacted)

xevon config set scanning_strategy.default_strategy deep
xevon config set scope.origin.mode strict
xevon config set audit.extensions.enabled true
xevon config set notify.enabled true
```

Config file location: `~/.xevon/xevon-configs.yaml`

---

## config clean

**Usage:** `xevon config clean [-F/--force]`

Reset xevon to a clean state: remove the entire `~/.xevon/` directory (config, database, extensions, prompts, SAST rules, session directories) and regenerate fresh defaults by running the same bootstrap as `xevon init`.

Prompts for `yes` confirmation unless `-F/--force` is passed.

```bash
xevon config clean
xevon config clean -F            # skip confirmation
```

Warning shown:

```
✗ Warn: This will remove ~/.xevon (config, database, and all local data)
Proceed? (type 'yes' to confirm):
```

Use this at the start of a new engagement or when an installation has drifted out of sync with the binary. After cleaning, the default API key is regenerated and all preset data (profiles, prompts, extensions, SAST rules, xevon-audit harness) is re-extracted.

---

## scope

**Usage:** `xevon scope [flags]` (aliases: `sc`)

Manage scan scope rules for filtering traffic.

### Subcommands

| Command | Aliases | Description |
|---------|---------|-------------|
| `scope view [component]` | `ls`, `list` | Display current scope configuration |
| `scope set <key> <value>` | — | Set a scope configuration value |

### Scope Components

host, path, status_code, request_content_type, response_content_type, request_string, response_string

### Examples

```bash
xevon scope view
xevon scope view host           # view host scope only
xevon scope set origin.mode strict
```

---

## source

**Usage:** `xevon source [flags]` (aliases: `src`)

Manage application source code links for whitebox scanning and SAST.

### Subcommands

| Command | Aliases | Description |
|---------|---------|-------------|
| `source ls` | `list` | List linked source repos |
| `source add` | — | Link source code to a hostname |
| `source rm <id>` | — | Remove a source repo link |
| `source scan <id>` | — | Run third-party security tools |

### source add flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--hostname` | `-H` | string | — | Target hostname (**required**) |
| `--path` | `-p` | string | — | Filesystem path to source root |
| `--git` | `-g` | string | — | Git URL to clone |
| `--name` | `-n` | string | dir basename | Display name |
| `--language` | `-l` | string | — | Primary language |
| `--framework` | `-f` | string | — | Framework (express, django, spring, etc.) |
| `--repo-type` | — | string | auto-detected | Type: git, folder, archive |
| `--scan-uuid` | — | string | — | Link to specific scan UUID |
| `--tag` | — | []string | — | Tags (repeatable) |

Note: `--path` and `--git` are mutually exclusive; one is required.

### Examples

```bash
xevon source ls
xevon source add --hostname api.example.com --path ./api-source
xevon source add --hostname example.com --git https://github.com/org/repo
xevon source add --hostname api.example.com --path ./src -l go -f gin
xevon source scan 1
xevon source rm 2
```

---

## auth

**Usage:** `xevon auth <subcommand> [flags]`

Manage authentication configurations and utilities.

### Subcommands

| Command | Aliases | Description |
|---------|---------|-------------|
| `auth lint` | — | Validate auth config files for errors and warnings |
| `auth list` | `ls` | List authentication configs |
| `auth load` | — | Load auth configs from a file or stdin into the database |
| `auth totp` | — | Generate a TOTP code from a base32 secret |

### Examples

```bash
xevon auth list
xevon auth lint auth-config.yaml
xevon auth load auth-config.yaml
xevon auth totp --secret JBSWY3DPEHPK3PXP
```

---

## project

**Usage:** `xevon project <subcommand> [flags]`

Manage projects for multi-tenancy scan data scoping.

### Subcommands

| Command | Aliases | Description |
|---------|---------|-------------|
| `project config` | — | View or update project configuration |
| `project create` | — | Create a new project |
| `project list` | `ls` | List all projects |
| `project use` | — | Switch to a project |

### Examples

```bash
xevon project list
xevon project create --name my-project
xevon project use my-project
xevon project config
```

---

## storage

**Usage:** `xevon storage <subcommand>`

Manage cloud-storage objects scoped to the **active project** (selected via `--project-id`, `--project-name`, or `XEVON_PROJECT`). Mirrors the REST endpoints under `/api/storage/*`.

**Requires** `storage.enabled: true` in `xevon-configs.yaml` (or `XEVON_STORAGE_ENABLED=true`) plus `storage.driver`, `storage.bucket`, `storage.access_key`, and `storage.secret_key`. When storage is disabled, every subcommand prints a tip showing how to enable it and exits cleanly (no error).

### Subcommands

| Command | Aliases | Description |
|---------|---------|-------------|
| `storage ls` | `list` | List objects under the project's prefix |
| `storage upload` | — | Upload a single local file |
| `storage download` | `get` | Download an object by key |
| `storage results` | — | Download a scan's `results.tar.gz` bundle by scan UUID |
| `storage presign` | — | Generate a presigned GET or PUT URL |
| `storage rm` | `delete` | Delete one or more objects |

### Source-archive integration

`agent audit` accepts `--source gs://<project>/<key>` for source archives. The archive is downloaded, extracted (`.zip / .tar.gz / .tar.bz2 / .tar.xz`), and cleaned up automatically.

### Result-bundle integration

Pass `--upload-results` to `scan`, `agent autopilot`, `agent swarm`, `agent audit`, or `agent query` to bundle the session/output and push it to storage at the end of the run. Native scan bundles land at `native-scans/<scan-uuid>/results.tar.gz`; agentic bundles land at `agentic-scans/<uuid>/results.tar.gz`. `xevon storage results <uuid>` checks both prefixes.

---

## storage ls

**Usage:** `xevon storage ls [flags]` (alias `list`)

List objects under the active project's prefix in the configured bucket.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--prefix` | string | — | Limit results to keys under this prefix (e.g. `ugc/`, `native-scans/`) |
| `--tree` | bool | `false` | Render objects as a directory tree (directories first, then files) |
| `--json` | bool | `false` | Output as JSON |

```bash
xevon storage ls
xevon storage ls --prefix ugc/
xevon storage ls --tree
xevon storage ls --json
```

---

## storage upload

**Usage:** `xevon storage upload <file> [flags]`

Upload a single local file to the active project's storage. Without `--key`, the file is stored under `ugc/<basename>` (matching `POST /api/storage/upload-source`). Pass `--key` to choose an explicit object key. Directories are rejected.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--key` | string | `ugc/<basename>` | Object key |
| `--content-type` | string | — | `Content-Type` to set on the object |

```bash
# Default key (ugc/report.pdf)
xevon storage upload ./report.pdf

# Explicit key
xevon storage upload ./report.pdf --key reports/q4.pdf

# Pin a content type
xevon storage upload ./report.pdf --content-type application/pdf
```

---

## storage download

**Usage:** `xevon storage download <key> [flags]` (alias `get`)

Download an object from the active project's storage by full key. Streams to stdout by default; use `-o` to write to a file.

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | — | Write to this file instead of stdout |

```bash
# Stream to a file
xevon storage download ugc/report.pdf -o report.pdf

# Stream to stdout (for pipelines)
xevon storage download ugc/notes.txt | grep TODO
```

---

## storage results

**Usage:** `xevon storage results <scan-uuid> [flags]`

Download the `results.tar.gz` bundle for a native or agentic scan run. Tries `native-scans/<uuid>/results.tar.gz` first, then `agentic-scans/<uuid>/results.tar.gz`. Default output filename is `results-<uuid>.tar.gz` in the current directory.

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | `results-<uuid>.tar.gz` | Write to this file |

```bash
xevon storage results 550e8400-e29b-41d4-a716-446655440000
xevon storage results 550e8400-e29b-41d4-a716-446655440000 -o results.tgz
```

---

## storage presign

**Usage:** `xevon storage presign --key <key> [flags]`

Generate a presigned URL for direct GET (download) or PUT (upload) against the active project's storage. Mirrors `POST /api/storage/presign`.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--key` | string | — | Object key (**required**) |
| `--method` | string | `GET` | HTTP method: `GET` or `PUT` |
| `--expiry` | duration | `1h` | URL validity duration (e.g. `30m`, `1h`, `24h`) |
| `--json` | bool | `false` | Output as JSON `{url, key, method, expiry_seconds}` |

```bash
# 1h GET URL
xevon storage presign --key ugc/foo.tar.gz --method GET --expiry 1h

# 30m PUT URL (machine-readable)
xevon storage presign --key uploads/new.tar.gz --method PUT --expiry 30m --json
```

---

## storage rm

**Usage:** `xevon storage rm <key> [<key>...] [flags]` (alias `delete`)

Permanently delete one or more objects from the active project's storage. Prompts for confirmation (typed `yes`) unless `-F` / `--force` is set.

```bash
# Single delete (prompts for confirmation)
xevon storage rm ugc/foo.tar.gz

# Bulk delete, skip confirmation
xevon storage rm ugc/a.pdf ugc/b.pdf -F
```

---

## strategy

**Usage:** `xevon strategy [flags]` (aliases: `st`, `phase`)

Display scanning strategies and their phase configurations.

### Subcommands

| Command | Aliases | Description |
|---------|---------|-------------|
| `strategy ls` | `list` | List available strategies |

### Examples

```bash
xevon strategy
xevon strategy ls
xevon phase              # alias for strategy
```

---

## version

**Usage:** `xevon version`

Show version, build time, commit, and author information. Supports `--json` for machine-readable output.
