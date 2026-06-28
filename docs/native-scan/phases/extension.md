# Extension Scanning

Extension scanning runs custom JavaScript or YAML extension modules against your targets. The `extension` phase skips all built-in Go scanner modules and runs only your extensions, giving you full control over the scanning logic.

## Quick Start

```bash
# Run a single extension script
xevon run extension -t https://example.com --ext my-scanner.js

# Run all extensions from a directory
xevon run extension -t https://example.com --ext-dir ./my-extensions/

# Run extensions alongside built-in modules (normal scan + extensions)
xevon scan -t https://example.com --ext my-scanner.js
```

## The Extension Phase

When you use `--only extension` (or its alias `ext`), xevon:

1. Skips all discovery, spidering, SPA, and ingestion phases
2. Disables all built-in Go scanner modules
3. Enables extensions automatically
4. Runs only your JS/YAML extension modules during dynamic assessment

```bash
# These are equivalent:
xevon run extension -t https://example.com
xevon scan -t https://example.com --only extension
xevon scan -t https://example.com --only ext
```

## Loading Extensions

### CLI Flags

| Flag | Description |
|------|-------------|
| `--ext <path>` | Load a specific extension script (repeatable) |
| `--ext-dir <dir>` | Override the extensions directory |

Both flags automatically enable extensions. They can be combined with any scan command.

```bash
# Load multiple extension scripts
xevon run extension -t https://example.com \
  --ext ./idor_detector.js \
  --ext ./header_leak.js

# Use a custom extensions directory
xevon run extension -t https://example.com --ext-dir ~/my-extensions/

# Mix both: directory + additional scripts
xevon scan -t https://example.com --ext-dir ~/extensions/ --ext ./extra-check.js
```

### Config File

Enable and configure extensions in `xevon-configs.yaml`:

```yaml
dynamic-assessment:
  extensions:
    enabled: true
    extension_dir: "~/.xevon/extensions/"   # Default extensions directory
    custom_dir:                                 # Additional script paths
      - ~/extra-extensions/custom-scanner.js
    variables:                                  # Custom variables accessible in scripts
      api_key: "your-key"
    allow_exec: false                           # Allow shell command execution
    sandbox_dir: "~/.xevon/sandbox/"         # Sandbox for exec operations
```

## Extension Types

Extensions plug into four points of the scanner pipeline:

| Type | Runs When | Use For |
|------|-----------|---------|
| `active` | During dynamic assessment | Send payloads, detect vulnerabilities |
| `passive` | Analyzing captured traffic | Inspect request/response without new traffic |
| `pre_hook` | Before each request is sent | Modify requests, skip assets, inject headers |
| `post_hook` | After a finding is emitted | Escalate severity, drop false positives |

Both JavaScript and YAML formats support all four types.

## Managing Extensions

```bash
# List loaded extensions
xevon extensions ls
xevon ext ls

# Filter by type
xevon ext ls --type active
xevon ext ls --type passive

# Show detailed info
xevon ext ls -v

# Show API documentation
xevon ext docs
xevon ext docs --example    # Include usage examples

# Install preset extensions
xevon ext preset            # Install all presets
xevon ext preset idor       # Install a specific preset

# Evaluate ad-hoc JavaScript
xevon ext eval 'xevon.log.info("hello")'
xevon ext eval --ext-file script.js
echo 'xevon.utils.md5("test")' | xevon ext eval --stdin
```

## Preset Extensions

xevon ships with starter extension presets. Install them with `xevon ext preset`:

| Preset | Type | Description |
|--------|------|-------------|
| `reflected_param_scanner` | active | Detect reflected parameters in responses |
| `idor_detector` | active | Detect Insecure Direct Object References |
| `ai_xss_scanner` | active | AI-augmented XSS scanning |
| `sensitive_header_leak` | passive | Detect sensitive information in response headers |
| `error_pattern_detector` | passive | Detect error patterns and stack traces |
| `ai_false_positive_filter` | post_hook | AI-powered false positive filtering |
| `ai_response_analyzer` | passive (YAML) | AI-augmented response analysis |
| `add_auth_header` | pre_hook | Inject authorization headers |
| `skip_static_assets` | pre_hook | Skip scanning of static asset URLs |
| `tag_critical_domains` | post_hook | Tag findings from critical domains |

Presets are installed to `~/.xevon/extensions/`.

## Extensions vs Built-in Modules

| | `--only extension` | `--ext` with normal scan |
|-|-------------------|--------------------------|
| Built-in Go modules | Disabled | Enabled |
| Extension modules | Enabled | Enabled |
| Discovery/Spidering | Disabled | Per strategy |
| Use case | Test extensions in isolation | Augment built-in scanning |

Use `--only extension` when developing or testing extensions. Use `--ext` with a normal scan to add extensions on top of built-in modules.

## Common Scenarios

```bash
# Develop and test a new extension
xevon run extension -t https://example.com --ext ./my-new-scanner.js

# Run only YAML-based extensions
xevon run extension -t https://example.com --ext-dir ./yaml-extensions/

# Extensions + full balanced scan
xevon scan -t https://example.com --ext ./custom-check.js

# Extensions with an OpenAPI spec
xevon scan -i api-spec.yaml -I openapi --only extension --ext ./api-fuzzer.js

# Extensions with custom variables
# (set variables in xevon-configs.yaml, access via xevon.config.get("key"))
xevon run extension -t https://example.com --ext ./scanner-with-config.js

# Run extension against a single URL
xevon scan-url https://example.com/api/users --ext ./idor_detector.js
```

## Further Reading

- [Writing Extensions](../customization/writing-extensions.md) — Full guide to writing JS and YAML extensions
- [Extension API Reference](../api-references/extensions.md) — Complete `xevon.*` API surface
- `xevon ext docs` — Built-in API documentation with examples
