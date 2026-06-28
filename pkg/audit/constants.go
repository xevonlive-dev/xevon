package audit

// HarnessName is the on-disk identifier for the xevon-audit harness.
// Drives the env-var prefix, finding-id prefix, and AgenticScan.Mode
// value. The on-disk output dir (`<source>/xevon-results/`) and the
// session subdir (`<session>/xevon-results/`) are set independently
// on HarnessSpec to match the upstream binary's output directory name.
const HarnessName = "audit"
