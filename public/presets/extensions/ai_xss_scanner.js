// ai_xss_scanner.js
// Active module: AI-augmented XSS scanner.
//
// Uses xevon.agent.generatePayloads() to obtain context-aware XSS payloads
// for each insertion point, then confirms any reflections with
// xevon.agent.analyzeResponse() before reporting — reducing false positives
// and adapting payloads to the detected parameter context.
//
// Requires: LLM client configured in the `agent.llm` section of
// xevon-configs.yaml (ANTHROPIC_API_KEY env var by default).
// Falls back gracefully (no-op) when the agent API is unavailable.

module.exports = {
  id: "ai-xss-scanner",
  name: "AI-Augmented XSS Scanner",
  description: "Generates context-aware XSS payloads via LLM and confirms findings before reporting",
  type: "active",
  severity: "high",
  confidence: "firm",
  tags: ["xss", "ai", "agent", "injection"],
  scanTypes: ["per_insertion_point"],

  scanPerInsertionPoint: function(ctx, insertion) {
    // Gracefully skip if agent API is not configured
    if (typeof xevon.agent === "undefined") return null;

    // Generate context-aware payloads for this specific parameter
    var payloads;
    try {
      payloads = xevon.agent.generatePayloads({
        type: "xss",
        parameter: insertion.name,
        context: "HTTP parameter — value is reflected in HTML response",
        count: 5
      });
    } catch (e) {
      xevon.log.warn("ai-xss-scanner: generatePayloads failed: " + e);
      return null;
    }
    if (!payloads || payloads.length === 0) return null;

    // Capture a baseline response for comparison (clean canary, not XSS)
    var baselineReq = insertion.buildRequest("VGNM_BASELINE_" + xevon.utils.randomString(6));
    var baselineResp = xevon.http.send(baselineReq);

    var findings = [];

    for (var i = 0; i < payloads.length; i++) {
      var payload = payloads[i];
      var attackReq = insertion.buildRequest(payload);
      var attackResp = xevon.http.send(attackReq);

      if (!attackResp || !attackResp.body) continue;

      // Quick heuristic: skip if no XSS-related content in response
      var body = attackResp.body;
      var looksReflected = body.indexOf(payload) !== -1 ||
        body.indexOf("<script>") !== -1 ||
        body.indexOf("onerror=") !== -1 ||
        body.indexOf("onload=") !== -1;
      if (!looksReflected) continue;

      // Confirm with LLM to reduce false positives
      var analysis;
      try {
        analysis = xevon.agent.analyzeResponse({
          request: attackReq,
          response: attackResp.raw,
          vulnerability_type: "xss",
          payload: payload,
          baseline_response: baselineResp ? baselineResp.raw : ""
        });
      } catch (e) {
        xevon.log.warn("ai-xss-scanner: analyzeResponse failed: " + e);
        // Fall back to heuristic-only result at low confidence
        analysis = {
          vulnerable: true,
          confidence: "low",
          evidence: "payload reflected (no LLM confirmation)",
          details: ""
        };
      }

      if (!analysis.vulnerable) continue;
      if (analysis.confidence === "low") continue; // require at least medium confidence

      findings.push({
        matched: payload,
        url: ctx.request.url,
        name: "AI-Confirmed XSS: " + insertion.name,
        description: "LLM-confirmed XSS in parameter '" + insertion.name +
          "'. Evidence: " + analysis.evidence +
          (analysis.details ? "\n\nDetails: " + analysis.details : ""),
        severity: analysis.confidence === "high" ? "high" : "medium",
        request: attackReq,
        response: attackResp.raw
      });
    }

    return findings.length > 0 ? findings : null;
  }
};
