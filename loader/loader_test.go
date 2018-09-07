package loader

import (
	"testing"

	"github.com/radu-matei/azure-functions-golang-worker/rpc"
)

func TestParseEntrypoint(t *testing.T) {
	testMetadata := rpc.RpcFunctionMetadata{Name: "Run", Directory: "", ScriptFile: "", EntryPoint: "Run", Bindings: nil}

	// Should fail
	_, err := parseEntrypoint(&testMetadata)
	if err == nil {
		t.Fatalf("parseEntrypoint() should have failed but did not")
	}

	testMetadata.ScriptFile = "../testData/parseTestData.go"

	azFuncArgs, err := parseEntrypoint(&testMetadata)
	if err != nil {
		t.Fatalf("parseEntrypoint() failed: %s", err)
	}

	if len(azFuncArgs) != 2 {
		t.Fatalf("function length is %d, should be 2", len(azFuncArgs))
	}

	if azFuncArgs[0].Name != "req" {
		t.Fatalf("first function name is %s, should be run", azFuncArgs[0].Name)
	}

	if azFuncArgs[1].Name != "ctx" {
		t.Fatalf("first function name is %s, should be ctx", azFuncArgs[1].Name)
	}
}
