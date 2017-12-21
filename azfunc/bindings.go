package azfunc

import (
	"reflect"

	"github.com/radu-matei/azure-functions-golang-worker/rpc"
)

// TODO - add other binding and trigger types
// TODO - in the end, every trigger is a binding - does it make sense to have separate types for trigger / binding?
const (
	// HTTPTriggerType represents a HTTP trigger in function load request from the host
	HTTPTriggerType = "httpTrigger"

	// BlobTriggerType represents a blob trigger in function load request from host
	BlobTriggerType = "blobTrigger"

	// HTTPBindingType represents a HTTP binding in function load request from the host
	HTTPBindingType = "http"

	// BlobBindingType represents a blob binding in function load request from the host
	BlobBindingType = "blob"
)

// StringToType - Because we don't have go/types information, we need to map the type info from the AST (which is string) to the actual types - see loader.go:83
// investiage automatically adding here all types from package azfunc
var StringToType = map[string]reflect.Type{
	"*azfunc.HTTPRequest": reflect.TypeOf((*HTTPRequest)(nil)),
	"*azfunc.Context":     reflect.TypeOf((*Context)(nil)),
	"*azfunc.BlobInput":   reflect.TypeOf((*BlobInput)(nil)),
}

// Func contains a function symbol with in and out param types
type Func struct {
	Func             reflect.Value
	Bindings         map[string]*rpc.BindingInfo
	In               []reflect.Type
	NamedInArgs      map[string]reflect.Type
	Out              []reflect.Type
	NamedOutBindings map[string]reflect.Value
}

// Context contains the runtime context of the function
type Context struct {
	FunctionID   string
	InvocationID string
}
