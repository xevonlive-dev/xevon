# xevon Security Scanner - AI Agent Integration

xevon is a web vulnerability scanner available as a CLI tool. Use these commands for security testing workflows.

## Quick Scan Commands

### Scan a single URL
```bash
# Basic GET scan with JSON output
xevon scan-url https://example.com/api/users --json

# Scan with specific method, body, and headers
xevon scan-url https://example.com/api/login \
  --method POST \
  --body '{"user":"admin","pass":"test"}' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer TOKEN' \
  --json

# Scan with specific modules only
xevon scan-url https://example.com/search?q=test -m xss-reflected,sqli-error --json
```

### Scan a raw HTTP request
```bash
# From stdin
echo -e "GET /api/users HTTP/1.1\r\nHost: example.com\r\nCookie: session=abc\r\n\r\n" | \
  xevon scan-request --json

# From file
xevon scan-request -i request.txt --json

# With target URL override
xevon scan-request -i request.txt --target https://staging.example.com --json
```

## Listing Modules

```bash
# List all scanner modules as JSON
xevon module ls --json

# List only active modules
xevon module ls --json --type active

# Filter modules by keyword
xevon module ls xss --json
```

## Ingesting HTTP Traffic

```bash
# Ingest URLs from a file
xevon ingest -i urls.txt --json

# Ingest from OpenAPI spec
xevon ingest -i openapi.yaml -t https://api.example.com --json

# Ingest from stdin
cat requests.txt | xevon ingest --json
```

## JSON Output

All commands support `--json` / `-j` for structured JSON output on stdout. Human-readable messages go to stderr. This makes it safe to parse output in pipelines:

```bash
result=$(xevon scan-url https://example.com/page --json)
findings=$(echo "$result" | jq '.findings')
count=$(echo "$result" | jq '.findings | length')
```

## Common Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--json` | `-j` | JSON output to stdout |
| `--modules` | `-m` | Comma-separated module IDs |
| `--concurrency` | `-c` | Worker count (default: 50) |
| `--timeout` | | HTTP timeout (default: 15s) |
| `--proxy` | | HTTP/SOCKS5 proxy URL |
| `--db` | | SQLite database path |
