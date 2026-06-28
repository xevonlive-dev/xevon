// Package api holds the shared types for the xevon.* JavaScript API surface:
// the namespace constants, the JSFuncDef declaration + handler factory, the
// APIOptions dependency bundle, and the APIFunction documentation record.
//
// It is a leaf package: the jsext engine and each per-domain function package
// (api/parse, …) import it, but it imports nothing from jsext. This lets API
// domains live in their own subpackages while sharing one source of truth for
// the registry types.
package api

import "github.com/grafana/sobek"

// Namespace constants for the xevon.* JS API.
const (
	NsRoot       = "xevon"
	NsLog        = "xevon.log"
	NsUtils      = "xevon.utils"
	NsParse      = "xevon.parse"
	NsHTTP       = "xevon.http"
	NsScan       = "xevon.scan"
	NsIngest     = "xevon.ingest"
	NsAgent      = "xevon.agent"
	NsDB         = "xevon.db"
	NsDBRecords  = "xevon.db.records"
	NsDBFindings = "xevon.db.findings"
	NsOAST       = "xevon.oast"
	NsRecord     = "xevon.record"
	NsConfig     = "xevon.config"
	NsMCP        = "xevon.mcp"
)

// HandlerFactory creates a JS function handler given runtime dependencies.
// It is called once per VM setup, not per invocation.
type HandlerFactory func(vm *sobek.Runtime, opts APIOptions) func(sobek.FunctionCall) sobek.Value

// JSFuncDef declares a JS API function with metadata and an optional handler factory.
// When MakeHandler is nil, the entry is metadata-only (e.g., dynamic config keys,
// per-request properties like xevon.record.uuid).
type JSFuncDef struct {
	Namespace   string
	Name        string
	Category    string
	Signature   string
	Returns     string
	Description string
	Example     string
	MakeHandler HandlerFactory // nil for metadata-only entries
}

// FullName returns the fully-qualified function name (e.g. "xevon.utils.sha256").
func (d JSFuncDef) FullName() string {
	return d.Namespace + "." + d.Name
}
