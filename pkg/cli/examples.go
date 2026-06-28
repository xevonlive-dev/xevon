package cli

import (
	"fmt"

	"github.com/xevonlive-dev/xevon/pkg/terminal"
)

func printFullExamples() {
	printSection("Scanning", []string{
		"xevon scan -t https://example.com",
		"xevon scan -t https://example.com -t https://api.example.com",
		"xevon scan -T targets.txt",
		"xevon scan -t https://example.com --strategy deep",
		"xevon scan -t https://example.com --scanning-profile quick",
		"xevon scan -t https://example.com --scanning-profile full",
		"xevon scan -t https://example.com --only dynamic-assessment",
		"xevon scan -t https://example.com --skip discovery,spidering",
		"xevon scan -t https://example.com -m xss-reflected,sqli-error",
		"xevon scan -t https://example.com --module-tag spring --module-tag injection",
		"xevon scan -t https://example.com --format jsonl -o results.jsonl",
		"xevon scan -t https://example.com --format html -o report.html",
		"xevon scan -t https://example.com --proxy http://127.0.0.1:8080",
		"xevon scan -t https://example.com -c 100 --rate-limit 200",
		"xevon scan -t https://example.com --scanning-max-duration 2h",
		"xevon scan -t https://example.com --ext custom-check.js",
		"xevon scan -t https://example.com --ext-dir ./my-extensions",
		"xevon scan -t https://example.com --only extension --ext custom-check.js",
		"xevon scan -t https://example.com --project-name my-project",
		"xevon scan -t https://example.com --oast-url https://interact.sh/abc123",
		"xevon scan -t https://example.com --known-issue-scan-tags cve,misconfig --known-issue-scan-severities critical,high",
	})

	printSection("Running Single Phases", []string{
		"xevon run discover -t https://example.com",
		"xevon run spidering -t https://example.com",
		"xevon run dynamic-assessment -t https://example.com",
		"xevon run dynamic-assessment -t https://example.com --module-tag spring",
		"xevon run external-harvest -t https://example.com",
		"xevon run known-issue-scan -t https://example.com",
		"xevon run known-issue-scan -t https://example.com --known-issue-scan-tags cve --known-issue-scan-severities critical,high",
		"xevon run extension -t https://example.com --ext custom-check.js",
		"xevon run deparos -t https://example.com",
		"xevon run dast -t https://example.com",
	})

	printSection("Input Modes", []string{
		"xevon scan -I openapi -i openapi.yaml -t https://api.example.com",
		"xevon scan -I burp -i burp-export.xml -t https://example.com",
		"xevon scan -I curl -i requests.txt",
		"xevon scan -I har -i traffic.har",
		"cat urls.txt | xevon scan -i -",
	})

	printSection("Ingestion", []string{
		"xevon ingest -t https://example.com -I openapi -i spec.yaml",
		"xevon ingest -t https://example.com -I burp -i export.xml",
		"cat urls.txt | xevon ingest -i -",
	})

	printSection("Server", []string{
		"xevon server",
		"xevon server --host 0.0.0.0 --service-port 8443",
		"xevon server --no-auth",
		"xevon server -t https://example.com --scan-on-receive",
	})

	printSection("Database & Results", []string{
		"xevon db ls",
		"xevon db ls --table findings",
		"xevon db stats",
		"xevon db clean --scan-uuid my-scan",
		"xevon traffic",
		"xevon traffic login",
		"xevon finding",
		"xevon export --format jsonl -o full-export.jsonl",
		"xevon export --format jsonl --only findings",
		"xevon export --format jsonl --only findings,http",
		"xevon export --format html -o report.html",
	})

	printSection("Strategy & Phases", []string{
		"xevon strategy",
		"xevon phase",
	})

	printSection("Modules", []string{
		"xevon module ls",
		"xevon module enable xss",
		"xevon module disable sqli",
		"xevon scan -M",
	})

	printSection("Extensions", []string{
		"xevon ext ls",
		"xevon ext docs",
		"xevon ext preset",
		`xevon ext eval 'xevon.log("hello")'`,
		"xevon ext eval --ext-file script.js",
	})

	printSection("Scope & Source", []string{
		"xevon scope view",
		"xevon scope set host.include '*.example.com'",
		"xevon source ls",
		"xevon source add --hostname api.example.com --path ./api-source",
		"xevon source scan 1",
	})

	printSection("Agent (AI)", []string{
		"xevon agent query --source ./src --prompt-template security-code-review",
		"xevon agent query --source ./src --prompt-template endpoint-discovery",
		"xevon agent query 'review this code for vulnerabilities'",
		"xevon agent query --agent-label code-review --prompt-file custom-prompt.md",
		"xevon agent --list-templates",
		"xevon agent swarm -t https://example.com --discover",
		"xevon agent swarm -t https://example.com --discover --focus 'API injection'",
		"xevon agent autopilot -t https://example.com",
	})

	printSection("Configuration", []string{
		"xevon config ls",
		"xevon config clean",
		"xevon version",
	})
}

func printSection(title string, examples []string) {
	fmt.Printf("  %s %s\n", terminal.InfoSymbol(), terminal.BoldCyan(title))
	for _, ex := range examples {
		fmt.Printf("    %s %s\n", terminal.ListSymbol(), terminal.Gray(ex))
	}
	fmt.Println()
}
