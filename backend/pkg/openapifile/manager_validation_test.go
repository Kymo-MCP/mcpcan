package openapifile

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/kymo-mcp/mcpcan/pkg/common"
	"github.com/kymo-mcp/mcpcan/pkg/logger"
)

var initLoggerOnce sync.Once

func ensureLogger(t *testing.T) {
	t.Helper()
	initLoggerOnce.Do(func() {
		if err := logger.Init("error", "console"); err != nil {
			t.Fatalf("failed to init logger: %v", err)
		}
	})
}

func TestValidateOpenapiFile_HTTPAuthWithIn_ShouldFail(t *testing.T) {
	t.Parallel()
	ensureLogger(t)
	tmpDir := t.TempDir()
	manager := NewOpenapiFileManager(&common.CodeConfig{
		Upload: common.UploadConfig{MaxFileSize: 10},
	}, tmpDir)

	filePath := filepath.Join(tmpDir, "invalid-openapi.yaml")
	content := `openapi: 3.0.0
info:
  title: Invalid Auth API
  version: 1.0.0
paths:
  /ping:
    get:
      operationId: ping
      responses:
        '200':
          description: ok
components:
  securitySchemes:
    jwtAuth:
      type: http
      scheme: bearer
      in: header
`
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	err := manager.ValidateOpenapiFile(filePath)
	if err == nil {
		t.Fatalf("expected validation to fail, got nil")
	}
	if !strings.Contains(err.Error(), "can't have 'in'") {
		t.Fatalf("expected security scheme validation error, got: %v", err)
	}
}

func TestValidateOpenapiFile_HTTPAuthWithoutIn_ShouldPass(t *testing.T) {
	t.Parallel()
	ensureLogger(t)
	tmpDir := t.TempDir()
	manager := NewOpenapiFileManager(&common.CodeConfig{
		Upload: common.UploadConfig{MaxFileSize: 10},
	}, tmpDir)

	filePath := filepath.Join(tmpDir, "valid-openapi.yaml")
	content := `openapi: 3.0.0
info:
  title: Valid Auth API
  version: 1.0.0
paths:
  /ping:
    get:
      operationId: ping
      responses:
        '200':
          description: ok
components:
  securitySchemes:
    jwtAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
`
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	if err := manager.ValidateOpenapiFile(filePath); err != nil {
		t.Fatalf("expected validation success, got: %v", err)
	}
}
