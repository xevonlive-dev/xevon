---
name: write-jsext
description: Reference for writing custom xevon JavaScript extensions. Use when you need to author a one-off scanner module — passive (reads existing HTTP records) or active (sends new requests) — and run it via the run_extension tool. Covers the module shape, the xevon.* API surface, and the common pitfalls.
license: MIT
allowed-tools:
  - read_file
  - write_file
  - edit_file
  - run_extension
  - run_scan
---

# Writing a Custom xevon JS Extension

You're authoring a single-file JavaScript extension that the xevon scanner
will load and run. Two output paths from here:

1. Hand the file path or inline source to **`run_extension`** to execute it
   immediately and get back findings.
2. Write the file to disk (e.g. into the session dir) for the operator to
   reuse later via `xevon scan --ext path/to/script.js`.

You do **not** need a build step. xevon's embedded JS engine (Sobek) runs
the script directly. ES2017 features work; no `import` / `export` — use
`module.exports`.

## Module Shape

Every extension exports one object via `module.exports`. The shape decides
whether it runs as a passive or active scanner module.

### Passive (read-only, runs against stored HTTP records)

```javascript
module.exports = {
  id: "kebab-case-id",                // unique within the run
  name: "Human-readable name",
  description: "1–2 sentence purpose",
  type: "passive",                     // required
  severity: "low",                     // critical|high|medium|low|info
  confidence: "tentative",             // certain|firm|tentative
  scope: "response",                   // request|response|both — what to read
  tags: ["exposure", "light"],         // free-form classification
  scanTypes: ["per_request"],          // only "per_request" for passive

  scanPerRequest: function(ctx) {
    // ctx.request:  { url, method, headers, body, ... }
    // ctx.response: { status, headers, body, ... }
    // ctx.record:   { uuid, addRemarks(tags) }  (DB-backed records only)
    //
    // Return null when there's nothing to report, OR an array of finding
    // objects (see "Finding shape" below). Throwing is fine — xevon
    // catches per-record errors and continues.
  }
};
```

### Active (sends new requests; can mutate inputs and probe)

```javascript
module.exports = {
  id: "active-extension-id",
  name: "Active Extension",
  description: "What it tests for",
  type: "active",                      // required
  severity: "medium",
  confidence: "tentative",
  tags: ["injection", "custom"],
  scanTypes: ["per_request"],          // most common; per_host / per_insertion_point also exist

  scanPerRequest: function(ctx) {
    // Same ctx as passive, plus you can call xevon.http.* to send
    // new requests, and xevon.parse.* / xevon.utils.* for helpers.
    // Return null or an array of findings.
  }
};
```

### Finding shape (what scanPerRequest returns)

```javascript
return [{
  url: ctx.request ? ctx.request.url : "",  // required
  matched: "the substring or evidence",      // shown verbatim in reports
  name: "Specific bug title",                // overrides module name
  description: "Markdown-friendly evidence + reasoning",
  severity: "high",                          // overrides module severity
  // optional:
  request: ctx.request ? ctx.request.raw : "",
  response: ctx.response ? ctx.response.raw : "",
  tags: ["custom", "graphql"]
}];
```

Or, instead of returning, you can call `xevon.scan.createFinding(obj)`
directly. Returning is preferred — it composes better with the executor's
deduplication and rate-limiting.

## Core API Cheat-Sheet

Only the most common pieces. Full TypeScript defs live at
`pkg/jsext/xevon.d.ts` — read it if you need an obscure helper.

### `xevon.http`

```javascript
xevon.http.get(url, opts)                 // -> HttpResponse
xevon.http.post(url, body, opts)          // -> HttpResponse
xevon.http.request({ url, method, headers, body })
xevon.http.send(rawHttpRequestString)     // raw request, ignores cookies
xevon.http.batch([req1, req2, ...])       // parallel; returns []
xevon.http.session({ headers, cookies })  // persistent jar + defaults
xevon.http.login({ url, body, ... })      // returns authed session
xevon.http.cachedGet(url, opts)           // memoized within the run
```

`HttpResponse` exposes `status`, `headers`, `body`, `raw`, plus helpers like
`.json()`, `.text()`. Treat `headers` as case-insensitive but values as
arrays-or-strings; normalize with `xevon.utils.headerValue(h, "name")`.

### `xevon.parse` and `xevon.utils`

```javascript
xevon.parse.url(url)                      // { scheme, hostname, port, path, query, ... }
xevon.utils.hasDynamicSegment(path)       // true for /users/123, /items/<uuid>
xevon.utils.pathToTemplate(path)          // /users/123 -> /users/*
xevon.utils.regexMatch(haystack, regex)   // boolean
xevon.utils.regexExtract(haystack, regex) // array of matches (capture group 1)
xevon.utils.md5(s) / sha256(s) / base64Encode(s)
```

### `xevon.log`

```javascript
xevon.log.info("...")   // visible in runtime.log + verbose output
xevon.log.warn("...")   // surfaces in operator console
xevon.log.error("...")  // tagged red; doesn't abort the script
```

Don't use `console.log` — it's a no-op in this engine.

### `xevon.db` (only when a repository is wired up — usually true at run time)

```javascript
xevon.db.records.query({ hostname, path, method, limit })
xevon.db.records.annotate(uuid, { risk_score, remarks })
xevon.db.compareResponses(records)        // anomaly grouping helper
xevon.db.findings.query({ scanUUID, severity })
```

### `xevon.scan`

```javascript
xevon.scan.isInScope(host, path)          // honor scope rules in custom checks
xevon.scan.getCurrentScan()               // { uuid: "<scan-uuid>" }
xevon.scan.createFinding(findingObj)      // bypasses return-array path
```

### `xevon.oast` (out-of-band callback testing)

```javascript
var oast = xevon.oast.allocate();         // { url, hostname, ... }
// fire a request that triggers a callback to oast.url
var hits = xevon.oast.poll(oast.id, 30);  // wait up to 30s for callbacks
```

Use this for SSRF, blind XSS, blind SQLi, log4shell-style probes.

### `xevon.agent` (LLM-assisted analysis — optional, may be unavailable)

```javascript
var verdict = xevon.agent.confirmFinding({
  name: "Reflected XSS",
  request: ctx.request.raw,
  response: ctx.response.raw,
  matched: payload
});
if (verdict && verdict.confirmed) { /* ... */ }

xevon.agent.generatePayloads({ type: "xss", context: "html", count: 5 });
```

Always null-check the result — `xevon.agent.*` returns null when no LLM
client is configured.

## Common Pitfalls

1. **No `import` / `export` / `require`** beyond `module.exports`. Stick to
   plain JS plus the `xevon.*` globals.
2. **Don't mutate `ctx.request` / `ctx.response`** — they're shared with
   other modules. If you need a modified request, build a new one with
   `xevon.http.buildRequest(ctx.request.raw, { ... })`.
3. **Regex flags use string syntax**, not literal `/.../i`. The engine
   accepts a string pattern; pass `"(?i)foo"` for case-insensitive.
4. **Skip irrelevant content types.** Most passive checks should bail
   early on CSS / images / fonts — see `internal_url_leak.js` for the
   pattern.
5. **Always set `confidence: "tentative"` for heuristic checks.** Save
   `firm` / `certain` for cases where you've actually verified the bug
   (e.g., OAST callback received, response diff confirmed).
6. **Per-request modules are fan-out.** They run for every record in
   scope — keep them cheap. Use `xevon.http.cachedGet` for repeated
   lookups, and bail on the first signal that the check doesn't apply.

## Minimal Working Example (passive)

```javascript
// detect-debug-headers.js
module.exports = {
  id: "debug-headers",
  name: "Debug headers exposed",
  description: "Flags responses leaking X-Debug-* / X-Powered-By in production-looking apps",
  type: "passive",
  severity: "low",
  confidence: "firm",
  scope: "response",
  tags: ["exposure", "headers", "light"],
  scanTypes: ["per_request"],

  scanPerRequest: function(ctx) {
    if (!ctx.response || !ctx.response.headers) return null;

    var leaked = [];
    var headers = ctx.response.headers;
    for (var name in headers) {
      var lower = name.toLowerCase();
      if (lower.indexOf("x-debug") === 0 ||
          lower === "x-powered-by" ||
          lower === "server") {
        leaked.push(name + ": " + headers[name]);
      }
    }
    if (leaked.length === 0) return null;

    return [{
      url: ctx.request.url,
      matched: leaked.join("; "),
      name: "Debug / framework headers exposed",
      description: "Response includes verbose framework headers:\n" +
        leaked.map(function(l) { return "- `" + l + "`"; }).join("\n"),
      severity: "low"
    }];
  }
};
```

## Iteration Loop

Once you have a draft:

1. **Validate by running it.** Call `run_extension` with `script_source`
   (or `script_path` if you wrote it to disk first). Pass concrete
   targets so you get real findings back, not just a smoke test.
2. **Read the result struct.** `finding_count > 0` is your signal that
   the matcher fires. Zero usually means a regex / scope mistake — add
   `xevon.log.info(...)` lines and re-run.
3. **Tighten before declaring success.** False positives are the
   default — a passing run on one URL doesn't mean the rule is good.
   Run against at least 2–3 targets, including one that should NOT
   match.
4. **Persist when settled.** If the operator wants this to run as part
   of regular scans, write the file under `<sessionDir>/extensions/`
   or a project-level extensions directory.

## When NOT to Write an Extension

If the bug class already has a built-in module (xss, sqli, ssrf, idor,
etc.), prefer `run_scan` with `modules: ["<id>"]` over a hand-written
extension. Extensions are for novel logic that doesn't fit the
generic scanner shape — protocol quirks, app-specific invariants,
correlation across records, custom OAST flows.
