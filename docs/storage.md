# Cloud Storage

xevon can read and write to S3-compatible object storage so scan inputs, exports, and result bundles can flow between machines, CI runs, and the server. This guide covers configuration, the `gs://` URL scheme, the `xevon storage` subcommands, and the storage-aware flags on `import`, `export`, `scan`, and the agent commands.

## What it does

- **Project-scoped object access.** Every key is prefixed with the active project's UUID, so two projects sharing the same bucket cannot read each other's data.
- **S3-compatible.** A single minio-go-based client talks to GCS (HMAC), S3, or self-hosted MinIO. Choose the driver in config; the rest is the same.
- **`gs://` URLs everywhere.** `xevon import` and `xevon export -o` accept `gs://<project-uuid>/<key>` directly — downloads to a temp file before importing, uploads from a temp file after exporting.
- **Result auto-archival.** `xevon scan --upload-results` and the agent commands' `--upload-results` flag tar+gzip the result tree and push it to a conventional key under `native-scans/` or `agentic-scans/`.

## Enabling storage

Storage is **disabled by default**. Enable it in `~/.xevon/xevon-configs.yaml`:

```yaml
storage:
  enabled: true
  driver: gcs            # gcs | s3 | minio  (default: gcs)
  bucket: my-xevon-bucket
  region: us-central1    # bucket region
  access_key: ${STORAGE_ACCESS_KEY}
  secret_key: ${STORAGE_SECRET_KEY}
  use_ssl: true          # default: true
  path_style: false      # default: false (set true for MinIO)
  endpoint: ""           # auto-detected for gcs/s3; required for minio
```

Endpoint defaults:
- `gcs` → `storage.googleapis.com` (uses HMAC keys, not service-account JSON)
- `s3` → `s3.amazonaws.com`
- `minio` → no default; you must set `endpoint` (e.g. `minio.example.com:9000`)

Required fields when enabled: `bucket`, `access_key`, `secret_key`. For MinIO, `endpoint` is also required.

### Toggling at runtime

The `XEVON_STORAGE_ENABLED` env var overrides the YAML `enabled` setting:

```bash
XEVON_STORAGE_ENABLED=true  xevon storage ls    # force-enable for one command
XEVON_STORAGE_ENABLED=false xevon scan --upload-results https://target  # force-disable
```

Truthy values: `1`, `true`, `yes`, `on`. Falsy: `0`, `false`, `no`, `off`. Unset means "use the YAML config".

## The `gs://` URL scheme

```
gs://<project-uuid>/<key>
```

- `<project-uuid>` is the xevon project UUID, **not** a GCS bucket name. The bucket comes from `storage.bucket` in config; the project UUID is the *prefix inside* that bucket.
- `<key>` is the object key within the project. May contain `/` for nested paths. Path traversal (`..`, leading `/`, backslashes) is rejected.

Example: with `storage.bucket: my-xevon-bucket`, the URL `gs://abc-123/imports/foo.tar.gz` resolves to the S3 object `s3://my-xevon-bucket/abc-123/imports/foo.tar.gz`.

When you pass a `gs://` URL whose project UUID differs from the active project, xevon logs an info-level notice but still proceeds — useful for cross-project copies, but worth noticing.

## Output placeholders

In any `-o` flag that writes to storage (or locally), two placeholders are expanded:

- `{ts}` → UTC timestamp `YYYY-MM-DDTHH-MM-SSZ` (filename-safe)
- `{project-uuid}` → active project UUID

```bash
xevon export --format html -o gs://{project-uuid}/exports/scan-{ts}.html
# → gs://<active-project>/exports/scan-2026-04-26T12-34-56Z.html
```

## Conventional key prefixes

xevon uses a small set of conventional prefixes inside each project. You can write to other prefixes freely; these are the ones the tooling auto-targets:

| Prefix | Used by | Contents |
|---|---|---|
| `ugc/<filename>` | `xevon storage upload` (default key) | Free-form user uploads |
| `imports/<basename>-<ts>.<ext>` | `xevon import --upload` (default key) | Bundles of import sources |
| `native-scans/<scan-uuid>/results.tar.gz` | `xevon scan --upload-results` | Native scan output bundle |
| `agentic-scans/<run-uuid>/results.tar.gz` | `xevon agent {autopilot,swarm,audit} --upload-results` | Agent session bundle |

## The `xevon storage` command

All subcommands operate on the active project. Every subcommand fails with a clear error if storage is disabled.

```bash
xevon storage ls                              # list all objects in the project
xevon storage ls --prefix ugc/                # filter by prefix
xevon storage ls --tree                       # render as a directory tree
xevon storage ls --json

xevon storage upload ./bundle.tar.gz                         # → ugc/bundle.tar.gz
xevon storage upload ./bundle.tar.gz --key imports/manual.tar.gz
xevon storage upload ./bundle.tar.gz --content-type application/gzip

xevon storage download <key>                  # writes to stdout by default
xevon storage download ugc/bundle.tar.gz -o ./local.tar.gz
xevon storage get ...                         # alias for download

xevon storage results <scan-uuid>             # downloads the scan's result bundle
xevon storage results <run-uuid>  -o ./run.tar.gz

xevon storage presign --key ugc/foo.tar.gz                    # default: GET, 1h
xevon storage presign --key uploads/incoming.tar.gz --method PUT --expiry 30m
xevon storage presign --key ugc/foo.tar.gz --json

xevon storage rm imports/old.tar.gz                            # prompts for confirmation
xevon storage rm imports/a.tar.gz imports/b.tar.gz -F          # batch + force
xevon storage delete ...                       # alias for rm
```

`ls` and `rm` accept their `list` / `delete` aliases. `download` accepts `get`.

## Storage-aware import / export

### Import from storage

```bash
xevon import gs://<project-uuid>/imports/litellm-audit.tar.gz
```

`xevon import` accepts `gs://` URLs as the input argument. The flow:

1. Download the object to a temp file.
2. If it's a `.tar.gz` / `.tgz` / `.zip`, extract it to a temp dir; otherwise treat as JSONL.
3. Detect audit folders (`audit-state.json`) or JSONL envelopes and run the matching importer.
4. For audit imports, copy the resolved folder verbatim into `~/.xevon/agent-sessions/<run-uuid>/audit/` and stamp `session_dir` and `storage_url` (the original `gs://` URL) on the agentic_scan row.
5. Clean up the temp dir.

### Push the import source up after importing

```bash
xevon import ./local-audit-folder --upload                # default key: imports/<base>-<ts>.tar.gz
xevon import ./local-audit-folder --upload-key imports/manual.zip   # zip instead
```

`--upload` bundles the local folder to `.tar.gz` (or `.zip` if `--upload-key` ends in `.zip`) and uploads it. Single files are uploaded as-is. `--upload-key` implies `--upload`. The flag is silently ignored if the input was already a `gs://` URL.

### Export to storage

Any export format can target a `gs://` URL via `-o`:

```bash
xevon export --format jsonl    -o gs://{project-uuid}/exports/data-{ts}.jsonl
xevon export --format html     -o gs://{project-uuid}/exports/report-{ts}.html
xevon export --format pdf      -o gs://{project-uuid}/exports/report-{ts}.pdf
xevon export --format markdown -o gs://{project-uuid}/exports/report-{ts}.md
xevon export --format bundle   -o gs://{project-uuid}/exports/bundle-{ts}.tar.gz \
                                    --scan-uuid <run-uuid>
```

xevon writes to a temp file locally, then uploads on success. If the upload fails, the export is reported as failed (the temp file is cleaned up regardless).

### Bundle export

The `bundle` format (alias `gz`) emits a single `.tar.gz` containing JSONL data, an HTML report, a manifest, and any agent session directories named with `--scan-uuid`:

```
<basename>/
  manifest.json
  export.jsonl
  report.html
  sessions/<uuid>/...    # verbatim copy of ~/.xevon/agent-sessions/<uuid>/
```

`--scan-uuid <uuid>` is repeatable. Each value is an `agentic_scans` row UUID; the bundle pulls in the matching `~/.xevon/agent-sessions/<uuid>/` directory. Missing or unknown UUIDs are warned-and-skipped, so the bundle still ships even if a session dir was pruned. When exactly one `--scan-uuid` is given and resolves to an `agentic_scans` row, the HTML report's target and duration are auto-filled from that row. CLI `--report-target` / `--report-duration` flags still take precedence.

`-o` is required and must end in `.tar.gz` or `.tgz`. The top-level directory name inside the tarball is the basename of `-o` minus the archive extension.

## Auto-uploading scan results

Native scan:

```bash
xevon scan https://target.com --upload-results
# → uploads native-scans/<scan-uuid>/results.tar.gz and stamps storage_url on the scan row
```

Agent runs:

```bash
xevon agent autopilot --input req.txt           --upload-results
xevon agent swarm     --input req.txt --discover --upload-results
xevon agent audit    --source ./repo --upload-results
# → uploads agentic-scans/<run-uuid>/results.tar.gz and stamps storage_url on the agentic_scan row
```

Both flags require storage to be enabled; they emit a warning and skip the upload otherwise (the scan or run still completes locally). The bundle includes the session directory plus, for native scans, the configured output formats and `runtime.log` when `persist_logs` is on.

## Round-tripping a bundle

Because `xevon import` understands tar.gz / zip archives and detects nested audit folders or JSONL files within them, you can ship a `bundle` export to another machine and re-import it:

```bash
# machine A
xevon export --format bundle -o /tmp/snapshot.tar.gz --scan-uuid <run-uuid>
xevon storage upload /tmp/snapshot.tar.gz --key snapshots/snapshot.tar.gz

# machine B
xevon import gs://<project-uuid>/snapshots/snapshot.tar.gz
```

Caveat: an imported bundle currently re-creates a new `agentic_scans` row from the embedded audit folder while the embedded `export.jsonl` re-imports the findings. The findings are de-duplicated, but the agentic_scan row is fresh — the new row will reference the same source data but won't share its UUID with the original.

## Security notes

- **Project isolation.** All keys are validated and prefixed with the project UUID server-side; clients cannot read or write objects outside their project.
- **Path traversal.** Keys containing `..`, `\`, or absolute paths are rejected by `storage.ValidateKey` before any backend call.
- **Presigned URLs** are bounded by `--expiry` (default 1h). They inherit the project prefix; they cannot be used to reach outside the project.
- **Credentials live in config.** Use `${ENV_VAR}` interpolation in `xevon-configs.yaml` rather than committing keys; the loader expands env vars at config-load time.

## Troubleshooting

| Symptom | Likely cause | Fix |
|---|---|---|
| `cloud storage is not enabled` | `storage.enabled` is false (or `XEVON_STORAGE_ENABLED=false`) | Set `storage.enabled: true` or `XEVON_STORAGE_ENABLED=true` |
| `storage.bucket must not be empty` | Bucket missing in config | Set `storage.bucket` in YAML |
| `storage.endpoint is required when driver is minio` | MinIO needs an explicit endpoint | Set `storage.endpoint: minio.example.com:9000` |
| `path contains traversal or invalid characters` | Key contains `..`, `\`, or escapes its project | Use a clean relative key like `exports/foo.html` |
| `Storage URL project (X) differs from active project (Y)` | `gs://` URL refers to a different project | Either `xevon project use <X>` or accept the cross-project copy |
| `failed to upload to gs://...` after a successful export | Network or credential issue at upload time | Re-run; check bucket permissions, IAM/HMAC key validity |
| `--upload-results specified but storage is not enabled` (warning, scan still completes) | Upload skipped at runtime | Enable storage, or drop the flag |

## Related

- [Configuration reference](configuration.md) — full YAML layout including the `storage:` block
- [Projects and multi-tenancy](projects.md) — how project UUIDs scope storage access
- [Output and reporting](output-and-reporting.md) — `--format` options and report metadata
- [Server mode](server-mode/) — REST API endpoints under `/api/storage/`
