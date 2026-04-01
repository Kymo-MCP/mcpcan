package biz

import (
	"strings"
	"testing"

	instancepb "github.com/kymo-mcp/mcpcan/api/market/instance"
)

func TestBuildCreateTokens_AutoCreateDefaultToken(t *testing.T) {
	b := &InstanceBiz{}
	tokens := b.buildCreateTokens("instance-1", true, nil)
	if len(tokens) != 1 {
		t.Fatalf("expected 1 token, got %d", len(tokens))
	}
	if tokens[0].InstanceId != "instance-1" {
		t.Fatalf("expected instanceId instance-1, got %s", tokens[0].InstanceId)
	}
	if !tokens[0].Enabled {
		t.Fatalf("expected default token enabled=true")
	}
	if !strings.HasPrefix(tokens[0].Token, "Bearer ") {
		t.Fatalf("expected token to start with Bearer, got %s", tokens[0].Token)
	}
	if len(tokens[0].Usages) != 1 || tokens[0].Usages[0] != "default" {
		t.Fatalf("expected usages=[default], got %#v", tokens[0].Usages)
	}
}

func TestBuildCreateTokens_DisabledAndNoTokens(t *testing.T) {
	b := &InstanceBiz{}
	tokens := b.buildCreateTokens("instance-1", false, nil)
	if len(tokens) != 0 {
		t.Fatalf("expected empty tokens, got %d", len(tokens))
	}
}

func TestBuildCreateTokens_PreserveProvidedTokens(t *testing.T) {
	b := &InstanceBiz{}
	tokens := b.buildCreateTokens("instance-1", true, []*instancepb.McpToken{
		{
			InstanceId: "old-id",
			Token:      "Bearer custom-token",
			Enabled:    true,
			Usages:     []string{"default"},
		},
	})
	if len(tokens) != 1 {
		t.Fatalf("expected 1 token, got %d", len(tokens))
	}
	if tokens[0].InstanceId != "instance-1" {
		t.Fatalf("expected instanceId rewritten to instance-1, got %s", tokens[0].InstanceId)
	}
	if tokens[0].Token != "Bearer custom-token" {
		t.Fatalf("expected original token preserved, got %s", tokens[0].Token)
	}
}

