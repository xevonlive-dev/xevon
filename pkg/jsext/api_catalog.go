package jsext

import "sync"

// Category constants for grouping API functions in display output.
const (
	CatLogging    = "Logging"
	CatEncoding   = "Encoding & Decoding"
	CatHashing    = "Hashing"
	CatStrings    = "Strings & Generation"
	CatSystem     = "System & Control"
	CatFileIO     = "File I/O"
	CatExtract    = "Data Extraction"
	CatDetection  = "Detection"
	CatHTTP       = "HTTP Requests"
	CatScan       = "Scan Control"
	CatIngest     = "Ingestion"
	CatSource     = "Source Code"
	CatConfig     = "Configuration"
	CatDBRecords  = "DB Records"
	CatDBFindings = "DB Findings"
	CatDBAnalysis = "DB Analysis"
	CatAgent      = "Agent"
	CatRecord     = "Current Record"
)

// APICatalogEntry is an APIFunction with an assigned display category.
type APICatalogEntry struct {
	Category string
	APIFunction
}

// ─── xevon.log ────────────────────────────────────────────────────────────
const exLogInfo = `xevon.log.info("scanning: " + ctx.request.url)`
const exLogWarn = `xevon.log.warn("unexpected status: " + resp.status)`
const exLogError = `xevon.log.error("request failed: " + ctx.request.url)`
const exLogDebug = `xevon.log.debug("payload: " + payload)`

// ─── xevon.utils / Encoding ───────────────────────────────────────────────
const exBase64Encode = `var encoded = xevon.utils.base64Encode("admin:password")`
const exBase64Decode = `var decoded = xevon.utils.base64Decode("YWRtaW4=")`
const exURLEncode = `var q = xevon.utils.urlEncode("key=val&foo=bar")`
const exURLDecode = `var raw = xevon.utils.urlDecode("key%3Dval")`
const exHTMLEncode = `var safe = xevon.utils.htmlEncode("<script>")`
const exHTMLDecode = `var raw = xevon.utils.htmlDecode("&lt;script&gt;")`

// ─── xevon.utils / Hashing ────────────────────────────────────────────────
const exSHA1 = `var hash = xevon.utils.sha1("hello")`
const exSHA256 = `var hash = xevon.utils.sha256("hello")`
const exMD5 = `var hash = xevon.utils.md5("hello")`

// ─── xevon.utils / Strings ────────────────────────────────────────────────
const exRandomString = `var canary = "VGNM" + xevon.utils.randomString(8)`

// ─── xevon.utils / System ─────────────────────────────────────────────────
const exSleep = `xevon.utils.sleep(500) // wait 500ms`
const exExec = `var result = xevon.utils.exec("curl -s http://example.com")`
const exGetEnv = `var token = xevon.utils.getEnv("API_TOKEN")`
const exSetEnv = `xevon.utils.setEnv("MY_VAR", "value")`

// ─── xevon.utils / File I/O ───────────────────────────────────────────────
const exGlob = `var files = xevon.utils.glob("/tmp/wordlists/*.txt")`
const exReadFile = `var data = xevon.utils.readFile("/tmp/tokens.txt")`
const exReadLines = `var lines = xevon.utils.readLines("/tmp/wordlist.txt")`
const exWriteFile = `xevon.utils.writeFile("/tmp/out.txt", "result data")`
const exMkdir = `xevon.utils.mkdir("/tmp/results")`

// ─── xevon.utils / Data Extraction ────────────────────────────────────────
const exJSONExtract = `var name = xevon.utils.jsonExtract('{"a":{"b":"val"}}', "a.b")`
const exRegexMatch = `if (xevon.utils.regexMatch(body, "error|exception")) { ... }`
const exRegexExtract = `var ver = xevon.utils.regexExtract(header, "Apache/([0-9.]+)")`
const exRegexFindAll = `var emails = xevon.utils.regexFindAll(body, "[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}")`
const exParseURL = `var host = xevon.utils.parse_url("https://sub.example.com/path", "%d")`
const exParseURLFile = `xevon.utils.parse_url_file("urls.txt", "%s://%d", "hosts.txt")`

// ─── xevon.utils / Detection ──────────────────────────────────────────────
const exDetectAnomaly = `var ranked = xevon.utils.detectAnomaly(responses)`

// ─── xevon.utils / TOTP ──────────────────────────────────────────────────
const exTOTPCode = `var otp = xevon.utils.totpCode("JBSWY3DPEHPK3PXP"); // {code: "123456", expires_in: 18}`

// ─── xevon.http ───────────────────────────────────────────────────────────
const exHTTPGet = `var resp = xevon.http.get("https://example.com/api")`
const exHTTPPost = `var resp = xevon.http.post("https://example.com/api", "data=test")`
const exHTTPRequest = `var resp = xevon.http.request({method: "PUT", url: "https://example.com", body: "{}"})`
const exHTTPSend = `var built = insertion.buildRequest(payload); var resp = xevon.http.send(built)`

// ─── xevon.scan ───────────────────────────────────────────────────────────
const exScanListModules = `var mods = xevon.scan.listModules()`
const exScanIsInScope = `if (xevon.scan.isInScope("example.com", "/admin")) { ... }`
const exScanGetScope = `var scope = xevon.scan.getScope()`
const exScanSetScope = `xevon.scan.setScope({host: {include: ["*.example.com"]}})`
const exScanCreateFinding = `xevon.scan.createFinding({url: ctx.request.url, matched: "admin", name: "Admin panel", severity: "medium"})`
const exScanGetCurrentScan = `var scan = xevon.scan.getCurrentScan()`
const exScanStartNewScan = `var r = xevon.scan.startNewScan({targets: ["https://example.com/api"], modules: ["xss", "sqli"]})`

// ─── xevon.ingest ─────────────────────────────────────────────────────────
const exIngestURL = `var r = xevon.ingest.url("https://example.com/api/users")`
const exIngestURLs = "var r = xevon.ingest.urls(\"https://example.com/a\\nhttps://example.com/b\")"
const exIngestCurl = `var r = xevon.ingest.curl("curl -X POST https://example.com/api -d 'data=test'")`
const exIngestRaw = "var r = xevon.ingest.raw(\"GET /api HTTP/1.1\\r\\nHost: example.com\\r\\n\\r\\n\")"
const exIngestOpenAPI = `var r = xevon.ingest.openapi(specJSON, {base_url: "https://api.example.com"})`
const exIngestPostman = `var r = xevon.ingest.postman(collectionJSON)`

// ─── xevon.db ─────────────────────────────────────────────────────────────
const exDBRecordsQuery = `var records = xevon.db.records.query({hostname: "example.com", limit: 10})`
const exDBRecordsGet = `var record = xevon.db.records.get("uuid-string")`
const exDBRecordsGetRelated = `var related = xevon.db.records.getRelated("uuid-string", {limit: 5})`
const exDBRecordsAnnotate = `xevon.db.records.annotate("uuid-string", {risk_score: 80, remarks: ["suspicious"]})`
const exDBFindingsQuery = `var findings = xevon.db.findings.query({severity: ["high", "critical"], limit: 20})`
const exDBFindingsGet = `var finding = xevon.db.findings.get(42)`
const exDBFindingsGetByRecord = `var findings = xevon.db.findings.getByRecord("uuid-string")`
const exDBFindingsCreate = `xevon.db.findings.create({module_id: "my-ext", module_name: "My Extension", severity: "high", description: "Found issue"})`
const exDBCompareResponses = `var result = xevon.db.compareResponses(records)`

// ─── xevon.record ─────────────────────────────────────────────────────────
const exRecordUUID = `var uuid = xevon.record.uuid`
const exRecordAnnotate = `xevon.record.annotate({risk_score: 80, remarks: ["suspicious"]})`
const exRecordAddRiskScore = `xevon.record.addRiskScore(10)`
const exRecordAddRemarks = `xevon.record.addRemarks(["admin-path", "needs-review"])`

// ─── xevon.config ─────────────────────────────────────────────────────────
const exConfigKey = `var token = xevon.config.auth_token`

var (
	apiCatalogOnce  sync.Once
	apiCatalogCache []APICatalogEntry
)

// APICatalog returns all API functions in display order, each with a category.
// The catalog is derived from allFuncDefs(), ensuring it never drifts from the
// actual registered handlers.
func APICatalog() []APICatalogEntry {
	apiCatalogOnce.Do(func() {
		defs := allFuncDefs()
		apiCatalogCache = make([]APICatalogEntry, 0, len(defs))
		for _, d := range defs {
			apiCatalogCache = append(apiCatalogCache, APICatalogEntry{
				Category: d.Category,
				APIFunction: APIFunction{
					Namespace:   d.Namespace,
					Name:        d.Name,
					Signature:   d.Signature,
					Returns:     d.Returns,
					Description: d.Description,
					Example:     d.Example,
				},
			})
		}
	})
	return apiCatalogCache
}
