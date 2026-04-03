package biz

import (
	"os"
	"testing"
)

func TestNormalizeGatewayURLForContainer(t *testing.T) {
	t.Setenv("MCP_INTERNAL_BASE_URL", "http://mcp-entry-svc")

	got := normalizeGatewayURLForContainer("http://localhost/mcp-gateway/abc/mcp")
	want := "http://mcp-entry-svc/mcp-gateway/abc/mcp"
	if got != want {
		t.Fatalf("normalize gateway url failed, got=%s want=%s", got, want)
	}
}

func TestNormalizeGatewayURLForContainer_KeepNonGatewayURL(t *testing.T) {
	t.Setenv("MCP_INTERNAL_BASE_URL", "http://mcp-entry-svc")

	raw := "https://api.example.com/v1/mcp"
	got := normalizeGatewayURLForContainer(raw)
	if got != raw {
		t.Fatalf("non-gateway url should not change, got=%s want=%s", got, raw)
	}
}

func TestNormalizeGatewayURLForContainer_NoInternalBase(t *testing.T) {
	_ = os.Unsetenv("MCP_INTERNAL_BASE_URL")

	raw := "http://localhost/mcp-gateway/abc/mcp"
	got := normalizeGatewayURLForContainer(raw)
	if got != raw {
		t.Fatalf("url should keep original when MCP_INTERNAL_BASE_URL empty, got=%s want=%s", got, raw)
	}
}
