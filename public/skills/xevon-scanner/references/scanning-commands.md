# Scanning Commands Reference

Complete flag reference for `scan`, `scan-url`, `scan-request`, and `run` commands.

## Table of Contents

- [scan](#scan)
- [scan-url](#scan-url)
- [scan-request](#scan-request)
- [run](#run)
- [Strategy and Phase Interaction](#strategy-and-phase-interaction)

---

## scan

**Usage:** `xevon scan [flags]`

Run a full vulnerability scan pipeline. Supports multiple targets, input formats, phase control, and strategy presets.

### Output flags (scan & run)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | — | Write findings to specified output file |
| `--stats` | — | bool | `false` | Show live progress stats during scanning |
| `--include-response` | — | bool | `false` | Include full HTTP response body in output |
| `--stateless` | — | bool | `false` | Use a temporary database, export results to `--output`, then discard |
| `--upload-results` | — | bool | `false` | Upload scan results to cloud storage after completion (requires storage config) |

Stateless mode is great for ephemeral CI/CD runs — it creates a temp SQLite file, runs the full scan against it, writes the export/report to `--output`, then deletes the DB (including WAL/SHM sidecars). Requires `--output`; mutually exclusive with `--db`. Combine with `--format jsonl` or `--format html` for shareable artifacts.

### Request flags (scan & run)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--header` | `-H` | []string | — | Add custom HTTP header (repeatable, e.g. -H 'Auth: Bearer token') |
| `--advanced-options` | `-a` | map | — | Module-specific options as key=value (e.g. -a xss.dom=true) |
| `--retries` | — | int | `1` | Retry attempts for failed requests |
| `--stream` | — | bool | `false` | Process targets as a stream without buffering or deduplication |

### Input Format flags (scan & run)

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--required-only` | bool | `false` | Parse only required fields from input format (ignore optional) |
| `--skip-format-validation` | bool | `false` | Skip validation of input file format |

### Other flags (scan & run)

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--auth-file` | []string | — | Path to auth file (YAML/JSON, single session or `sessions:` bundle), or bare name resolved against session_dir. Repeatable. |
| `--auth` | []string | — | Inline session in `name:Header:value` format. Repeatable. |
| `--oast-url` | string | — | Fixed out-of-band callback URL (overrides auto-generated interactsh URL) |
| `--pilot` | bool | `false` | Enable AI pilot-driven crawling |

### Content Discovery flags (scan & run)

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--discover` | bool | `false` | Enable content discovery phase before scanning |
| `--discover-max-time` | duration | `1h` | Max time for content discovery per target |
| `--fuzz-wordlist` | string | — | Custom fuzz wordlist path (enables fuzzing during discovery) |
| `--no-prefix-breaker` | bool | `false` | Disable per-prefix circuit breaker that stops trap-directory recursion |

### Browser Spidering flags (scan & run)

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--spider` | — | bool | `false` | Enable browser-based spidering phase before scanning |
| `--spider-max-time` | — | duration | `30m` | Max time for spidering per target |
| `--browser-engine` | `-E` | string | `chromium` | Browser engine: chromium, ungoogled, fingerprint |
| `--browsers` | `-b` | int | `1` | Number of parallel browser instances for spidering |
| `--headless` | — | bool | `true` | Run browser in headless mode |
| `--no-cdp` | — | bool | `false` | Disable Chrome DevTools Protocol event listener detection |
| `--no-forms` | — | bool | `false` | Disable automatic form detection and filling during spidering |

### External Harvest flags (scan & run)

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--external-harvest` | bool | `false` | Enable external intelligence gathering phase (Wayback, CT logs, etc.) |

### KnownIssueScan flags (scan & run)

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--known-issue-scan-tags` | []string | — | Nuclei template tags to include |
| `--known-issue-scan-exclude-tags` | []string | — | Nuclei template tags to exclude |
| `--known-issue-scan-severities` | []string | — | Filter Nuclei templates by severity (critical,high,medium,low,info) |
| `--known-issue-scan-templates-dir` | string | — | Custom Nuclei templates directory |

### SAST flags (scan & run)

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--rule` | string | — | Filter SAST rules by fuzzy name match |
| `--sast-adhoc` | string | — | Ad-hoc SAST scan: local path or git URL (auto-detected, results not saved to database) |

### Examples

```bash
# Basic scan
xevon scan -t https://example.com

# Multiple targets
xevon scan -t https://example.com -t https://api.example.com

# Targets from file
xevon scan -T targets.txt

# Deep strategy with discovery
xevon scan -t https://example.com --strategy deep

# Phase isolation
xevon scan -t https://example.com --only dynamic-assessment
xevon scan -t https://example.com --only ext --ext ./custom-check.js
xevon scan -t https://example.com --skip discovery,spidering

# Specific modules
xevon scan -t https://example.com -m xss-reflected,sqli-error

# Custom scanning profile
xevon scan -t https://example.com --scanning-profile aggressive

# JSONL output
xevon scan -t https://example.com --format jsonl -o results.jsonl

# HTML report
xevon scan -t https://example.com --format html -o report.html

# With proxy
xevon scan -t https://example.com --proxy http://127.0.0.1:8080

# Speed tuning
xevon scan -t https://example.com -c 100 --rate-limit 200

# Whitebox scanning
xevon scan -t https://example.com --source ./src --strategy whitebox

# Whitebox via git clone
xevon scan -t https://example.com --source https://github.com/org/repo --strategy whitebox

# OpenAPI scan
xevon scan -I openapi -i openapi.yaml -t https://api.example.com

# Burp import scan
xevon scan -I burp -i burp-export.xml -t https://example.com

# Pipe from stdin
cat urls.txt | xevon scan -i -

# Filter modules by tag
xevon scan -t https://example.com --module-tag spring --module-tag injection

# Run extension during scan
xevon scan -t https://example.com --ext custom-check.js

# Extensions-only scan
xevon scan -t https://example.com --only extension --ext custom-check.js
```

---

## scan-url

**Usage:** `xevon scan-url <url> [flags]`

Scan a single URL for vulnerabilities. Designed for quick, targeted scans and AI agent integration. Returns JSON output with findings.

### scan-url specific flags

**Spidering:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--spider` | bool | `false` | Run browser-based spidering before scanning |

**Discovery:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--discover` | bool | `false` | Run content discovery before scanning |

**Harvest:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--external-harvest` | bool | `false` | Run external intelligence harvesting before scanning |

**Request:**

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--header` | `-H` | []string | — | Custom header (repeatable) |

**Other:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--method` | string | `GET` | HTTP method |
| `--body` | string | — | Request body |
| `--known-issue-scan` | bool | `false` | Run known issue scan (Nuclei/Kingfisher) |
| `--no-passive` | bool | `false` | Skip passive modules |
| `--no-insertion-points` | bool | `false` | Skip insertion point testing |

### Examples

```bash
# Simple GET scan
xevon scan-url https://example.com/api/users

# POST with body
xevon scan-url https://example.com/login \
  --method POST --body '{"user":"admin","pass":"test"}' \
  -H "Content-Type: application/json"

# With discovery phase
xevon scan-url https://example.com --discover

# Specific modules, no passive
xevon scan-url https://example.com/api -m xss-reflected --no-passive
```

---

## scan-request

**Usage:** `xevon scan-request [flags]`

Read a raw HTTP request from file or stdin and run scanner modules against it. Designed for pipeline integration and AI agent workflows.

### scan-request specific flags

**Spidering:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--spider` | bool | `false` | Run browser-based spidering before scanning |

**Discovery:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--discover` | bool | `false` | Run content discovery before scanning |

**Harvest:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--external-harvest` | bool | `false` | Run external intelligence harvesting before scanning |

**Other:**

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--input` | `-i` | string | `-` (stdin) | Input file or stdin |
| `--target` | — | string | — | Override target URL (scheme://host) |
| `--known-issue-scan` | — | bool | `false` | Run known issue scan |
| `--no-passive` | — | bool | `false` | Skip passive modules |
| `--no-insertion-points` | — | bool | `false` | Skip insertion point testing |

### Examples

```bash
# From file
xevon scan-request -i raw-request.txt

# From stdin
echo -e "GET /api/users HTTP/1.1\r\nHost: example.com\r\n" | xevon scan-request

# With target override
xevon scan-request -i request.txt --target https://staging.example.com

# With discovery
xevon scan-request -i request.txt --discover
```

---

## run

**Usage:** `xevon run <phase> [flags]`

**Aliases:** `r`

Run a single scan phase directly. Equivalent to `xevon scan --only <phase>`.

### Valid phases

| Phase | Aliases |
|-------|---------|
| `ingestion` | — |
| `discovery` | `deparos`, `discover` |
| `external-harvest` | — |
| `known-issue-scan` | — |
| `spidering` | `spitolas` |
| `sast` | — |
| `dynamic-assessment` | `audit`, `dast`, `assessment` |
| `extension` | `ext` |

The `run` command accepts the same flag groups as `scan`: Spidering, Discovery, Harvest, KnownIssueScan, SAST, Input Format, Request, Output, and Other (--oast-url, --pilot).

### Examples

```bash
xevon run discover -t https://example.com
xevon run spidering -t https://example.com
xevon run audit -t https://example.com
xevon run audit -t https://example.com --module-tag spring
xevon run external-harvest -t https://example.com
xevon run known-issue-scan -t https://example.com
xevon run known-issue-scan -t https://example.com --known-issue-scan-tags cve --known-issue-scan-severities critical,high
xevon run sast --sast-adhoc /path/to/app
xevon run sast --sast-adhoc /path/to/app --rule gin
xevon run extension -t https://example.com --ext custom-check.js
xevon run ext -t https://example.com --ext ./my-scanner.js
xevon run deparos -t https://example.com
xevon run audit -t https://example.com
```

---

## Strategy and Phase Interaction

### Precedence

1. `--only <phase>` overrides everything — only that phase runs, heuristics disabled
2. `--skip <phase>` disables specific phases while keeping all others
3. `--strategy <name>` sets baseline phase configuration
4. Individual phase flags (`--discover`, `--spider`, etc.) override strategy settings
5. Config file `scanning_strategy.default_strategy` provides the lowest-precedence default

### Heuristics

- Default: `--heuristics-check basic`
- Levels: `none`, `basic`, `advanced`
- `basic` probes target root pages to detect content type (HTML / JSON / blank) and skips spidering for non-HTML targets
- `advanced` adds deep HTML analysis to detect SPA frameworks and optimize phase selection
- `none` runs all enabled phases unconditionally
- `--skip-heuristics` is shorthand for `--heuristics-check=none`
- `--only` automatically disables heuristics
- Precedence: `--skip-heuristics` > `--heuristics-check` > config > `basic`

### Intensity Presets

`--intensity quick|balanced|deep` is a cross-cutting preset that maps to a scanning profile + strategy. It is also honored by `agent autopilot` and `agent swarm` with backend-specific defaults. Explicit flags always override the preset — e.g. `--intensity deep --scanning-profile foo` applies `deep`'s strategy but your custom profile.

### Scanning Pace

Speed settings have a layered precedence:

1. CLI flags (`-c`, `--rate-limit`, `--max-per-host`) — highest
2. `--scanning-max-duration` — overrides `scanning_pace.max_duration`
3. Config `scanning_pace` section — per-phase max_duration and duration_factor
4. Built-in defaults — lowest

### CI Output

- `--ci-output-format` enables CI-friendly output: JSONL findings only, no color, no banners
- Equivalent to combining `--format jsonl --silent`
- Useful for CI/CD pipelines that parse JSON output

### Valid `--only` Phases

The following phases can be used with `--only` and `--skip`:

`ingestion`, `discovery`, `external-harvest`, `known-issue-scan`, `spidering`, `sast`, `audit`, `extension`

### HTML Format Constraints

- `--format html` requires `-o/--output`
- In `scan` mode with `--only`, HTML is only supported for `discovery` and `spidering` phases
- The `export` command supports HTML for all data

### SAST Constraints

- `--sast-adhoc` accepts either a local path or a git URL (auto-detected)
- Git URLs are cloned to a temp directory automatically
