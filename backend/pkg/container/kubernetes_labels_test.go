package container

import "testing"

func TestSanitizeK8sLabels_DropsTraefikAndInvalidLabels(t *testing.T) {
	input := map[string]string{
		"app":                                     "mcp-ext-123",
		"managed-by":                              "mcpcan",
		"traefik.enable":                          "true",
		"traefik.http.routers.r.rule":             "HostRegexp(`{host:.+}`) && PathPrefix(`/mcp-gateway/abc`)",
		"invalid value key":                       "ok",
		"valid-key":                               "bad/value",
		"mcp.instance.type":                       "proxy-translator",
		"traefik.http.services.long.server.port":  "61180",
	}

	got := sanitizeK8sLabels(input)

	if _, ok := got["traefik.enable"]; ok {
		t.Fatalf("expected traefik label to be dropped")
	}
	if _, ok := got["traefik.http.routers.r.rule"]; ok {
		t.Fatalf("expected traefik router label to be dropped")
	}
	if _, ok := got["invalid value key"]; ok {
		t.Fatalf("expected invalid key label to be dropped")
	}
	if _, ok := got["valid-key"]; ok {
		t.Fatalf("expected invalid value label to be dropped")
	}

	if got["app"] != "mcp-ext-123" {
		t.Fatalf("expected app label to be preserved")
	}
	if got["managed-by"] != "mcpcan" {
		t.Fatalf("expected managed-by label to be preserved")
	}
	if got["mcp.instance.type"] != "proxy-translator" {
		t.Fatalf("expected mcp.instance.type label to be preserved")
	}
}

func TestSanitizeK8sLabels_AllTraefik_ShouldBecomeEmpty(t *testing.T) {
	input := map[string]string{
		"traefik.enable":                      "true",
		"traefik.http.routers.demo.rule":      "PathPrefix(`/mcp-gateway/abc`)",
		"traefik.http.services.demo.port":     "61180",
	}

	got := sanitizeK8sLabels(input)
	if len(got) != 0 {
		t.Fatalf("expected all-traefik labels to be dropped, got=%v", got)
	}
}
