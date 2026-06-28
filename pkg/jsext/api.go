package jsext

import (
	"github.com/grafana/sobek"
)

// SetupAPI installs the xevon.* global namespace on a sobek VM.
func SetupAPI(vm *sobek.Runtime, opts APIOptions) {
	// Create top-level xevon object
	xevon := vm.NewObject()
	_ = vm.Set("xevon", xevon)

	// Set up config variables
	configObj := vm.NewObject()
	for k, v := range opts.ConfigVars {
		_ = configObj.Set(k, v)
	}
	_ = xevon.Set("config", configObj)

	// Register all functions via the declarative registry
	registerFuncs(vm, opts, allFuncDefs())
}
