package openapifile

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/kymo-mcp/mcpcan/pkg/common"
	"github.com/kymo-mcp/mcpcan/pkg/database/model"
	"github.com/kymo-mcp/mcpcan/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// OpenapiFileManager manages OpenAPI files.
type OpenapiFileManager struct {
	config     *common.CodeConfig
	pathPrefix string
}

// NewOpenapiFileManager creates a new OpenapiFileManager instance.
func NewOpenapiFileManager(codeConfig *common.CodeConfig, pathPrefix string) *OpenapiFileManager {
	return &OpenapiFileManager{
		config:     codeConfig,
		pathPrefix: pathPrefix,
	}
}

// OpenapiFileInfo represents information about an OpenAPI file.
type OpenapiFileInfo struct {
	OpenapiFileID   string
	OpenapiFilePath string
	OriginalName    string
	FileSize        int64
	OpenapiFileType model.OpenapiFileType
}

// UploadOpenapiFile uploads an OpenAPI file.
func (m *OpenapiFileManager) UploadOpenapiFile(file multipart.File, header *multipart.FileHeader) (*OpenapiFileInfo, error) {
	// Log upload start information
	logger.Info("Starting OpenAPI file upload",
		zap.String("filename", header.Filename),
		zap.Int64("fileSize", header.Size),
		zap.Int("configMaxSizeMB", m.config.Upload.MaxFileSize))

	// Validate file type
	fileType, err := m.validateFileType(header.Filename)
	if err != nil {
		logger.Error("File type validation failed",
			zap.String("filename", header.Filename),
			zap.Error(err))
		return nil, err
	}

	// Validate file size
	if err := m.validateFileSize(header.Size); err != nil {
		logger.Error("File size validation failed",
			zap.String("filename", header.Filename),
			zap.Int64("fileSize", header.Size),
			zap.Int("maxSizeMB", m.config.Upload.MaxFileSize),
			zap.Error(err))
		return nil, err
	}

	// Generate file ID
	fileID := uuid.New().String()

	// Create file directory structure
	fileDir, err := m.createFileDirectory(fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to create file directory: %v", err)
	}

	// Save original file
	filePath, err := m.saveOriginalFile(file, fileDir, header.Filename)
	if err != nil {
		// Clean up directory
		os.RemoveAll(fileDir)
		return nil, fmt.Errorf("failed to save original file: %v", err)
	}

	// Convert to relative path based on the configured root path
	relFilePath, _ := m.ToRelativePath(filePath)

	fileInfo := &OpenapiFileInfo{
		OpenapiFileID:   fileID,
		OpenapiFilePath: relFilePath,
		OriginalName:    header.Filename,
		FileSize:        header.Size,
		OpenapiFileType: fileType,
	}

	logger.Info("OpenAPI file uploaded successfully",
		zap.String("openapiFileId", fileID),
		zap.String("filePath", relFilePath),
		zap.String("fileType", string(fileType)))

	return fileInfo, nil
}

// validateFileType validates the file type for OpenAPI documents
func (m *OpenapiFileManager) validateFileType(filename string) (model.OpenapiFileType, error) {
	allowedExtensions := []string{".json", ".yaml", ".yml"}
	filename = strings.ToLower(filename)

	for _, ext := range allowedExtensions {
		if strings.HasSuffix(filename, ext) {
			switch ext {
			case ".json":
				return model.OpenapiFileTypeJson, nil
			case ".yaml", ".yml":
				return model.OpenapiFileTypeYaml, nil
			}
		}
	}

	return model.OpenapiFileTypeUnknown, fmt.Errorf("unsupported file type for OpenAPI document, allowed extensions: %v", allowedExtensions)
}

// validateFileSize validates the file size.
func (m *OpenapiFileManager) validateFileSize(size int64) error {
	maxSize := int64(m.config.Upload.MaxFileSize) * 1024 * 1024 // Convert to bytes
	if size > maxSize {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size %d MB", size, m.config.Upload.MaxFileSize)
	}
	return nil
}

// createFileDirectory creates the file directory.
func (m *OpenapiFileManager) createFileDirectory(fileID string) (string, error) {
	// Create directory structure based on configuration: root_path/openapi-{id}
	fileDirName := fmt.Sprintf("openapi-%s", fileID)
	fileDir := filepath.Join(m.pathPrefix, fileDirName)

	if err := os.MkdirAll(fileDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %v", fileDir, err)
	}

	return fileDir, nil
}

// saveOriginalFile saves the original OpenAPI file.
func (m *OpenapiFileManager) saveOriginalFile(file multipart.File, fileDir, filename string) (string, error) {
	// Reset file pointer to the beginning
	file.Seek(0, 0)

	filePath := filepath.Join(fileDir, filename)
	outFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %v", filePath, err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, file); err != nil {
		return "", fmt.Errorf("failed to copy file content: %v", err)
	}

	return filePath, nil
}

// ToRelativePath converts an absolute path to a relative path based on the configured root path.
func (m *OpenapiFileManager) ToRelativePath(absolutePath string) (string, error) {
	// Get the absolute path of the configured root path
	absRootPath, err := filepath.Abs(m.pathPrefix)
	if err != nil {
		return absolutePath, err
	}

	// Get the absolute path of the target path
	absTargetPath, err := filepath.Abs(absolutePath)
	if err != nil {
		return absolutePath, err
	}

	// Calculate the relative path
	relPath, err := filepath.Rel(absRootPath, absTargetPath)
	if err != nil {
		return absolutePath, err
	}

	return relPath, nil
}

// ToAbsolutePath converts a relative path to an absolute path.
func (m *OpenapiFileManager) ToAbsolutePath(relativePath string) (string, error) {
	// If it's already an absolute path, return directly
	if filepath.IsAbs(relativePath) {
		return relativePath, nil
	}

	// Get the absolute path of the configured root path
	absRootPath, err := filepath.Abs(m.pathPrefix)
	if err != nil {
		return "", err
	}

	// Join to get the absolute path
	absolutePath := filepath.Join(absRootPath, relativePath)
	return absolutePath, nil
}

// DeleteOpenapiFile removes an OpenAPI file and its directory
func (m *OpenapiFileManager) DeleteOpenapiFile(filePath string) error {
	// Convert relative path to absolute path
	absFilePath, err := m.ToAbsolutePath(filePath)
	if err != nil {
		return fmt.Errorf("failed to convert to absolute path: %v", err)
	}

	// Get the directory containing the file
	fileDir := filepath.Dir(absFilePath)

	// Check if file directory exists
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		logger.Warn("OpenAPI file directory does not exist", zap.String("path", fileDir))
		return nil // Consider it as already deleted
	}

	// Remove the entire file directory
	if err := os.RemoveAll(fileDir); err != nil {
		logger.Error("Failed to remove OpenAPI file directory",
			zap.String("path", fileDir),
			zap.Error(err))
		return fmt.Errorf("failed to remove OpenAPI file directory: %v", err)
	}

	logger.Info("OpenAPI file deleted successfully", zap.String("path", fileDir))
	return nil
}

// GetOpenapiFileContent reads the content of an OpenAPI file
func (m *OpenapiFileManager) GetOpenapiFileContent(filePath string) (string, error) {
	// Convert relative path to absolute path
	absFilePath, err := m.ToAbsolutePath(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to convert to absolute path: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("OpenAPI file not found: %s", absFilePath)
	}

	// Read file content
	content, err := os.ReadFile(absFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read OpenAPI file: %v", err)
	}

	return string(content), nil
}

// UpdateOpenapiFileContent updates the content of an OpenAPI file
func (m *OpenapiFileManager) UpdateOpenapiFileContent(filePath string, content string) error {
	// Convert relative path to absolute path
	absFilePath, err := m.ToAbsolutePath(filePath)
	if err != nil {
		return fmt.Errorf("failed to convert to absolute path: %v", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(absFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Write file content
	if err := os.WriteFile(absFilePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write OpenAPI file: %v", err)
	}

	logger.Info("OpenAPI file updated successfully", zap.String("path", absFilePath))
	return nil
}

// ValidateOpenapiFile validates if the file is a valid OpenAPI document
func (m *OpenapiFileManager) ValidateOpenapiFile(filePath string) error {
	// Convert relative path to absolute path
	absFilePath, err := m.ToAbsolutePath(filePath)
	if err != nil {
		return fmt.Errorf("failed to convert to absolute path: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(absFilePath); os.IsNotExist(err) {
		return fmt.Errorf("OpenAPI file not found: %s", absFilePath)
	}

	// Read file content for basic validation
	content, err := os.ReadFile(absFilePath)
	if err != nil {
		return fmt.Errorf("failed to read OpenAPI file for validation: %v", err)
	}

	// Basic validation based on file extension
	ext := strings.ToLower(filepath.Ext(absFilePath))
	contentStr := string(content)

	// For JSON files, check if it starts with { or [
	if ext == ".json" {
		trimmed := strings.TrimSpace(contentStr)
		if !strings.HasPrefix(trimmed, "{") && !strings.HasPrefix(trimmed, "[") {
			return fmt.Errorf("invalid JSON format: file does not start with { or [")
		}
	}

	// For YAML files, check for basic YAML structure
	if ext == ".yaml" || ext == ".yml" {
		trimmed := strings.TrimSpace(contentStr)
		if trimmed == "" {
			return fmt.Errorf("invalid YAML format: file is empty")
		}
		// Basic YAML validation - should contain at least one key-value pair or document separator
		if !strings.Contains(trimmed, ":") && !strings.Contains(trimmed, "---") {
			return fmt.Errorf("invalid YAML format: no key-value pairs or document separators found")
		}
	}

	doc, err := parseOpenapiDocument(content, ext)
	if err != nil {
		return fmt.Errorf("failed to parse OpenAPI document: %w", err)
	}
	if err := validateOpenapiSecuritySchemes(doc); err != nil {
		return err
	}

	logger.Info("OpenAPI file validation passed", zap.String("path", absFilePath))
	return nil
}

func parseOpenapiDocument(content []byte, ext string) (map[string]interface{}, error) {
	var doc map[string]interface{}
	switch ext {
	case ".json":
		if err := json.Unmarshal(content, &doc); err != nil {
			return nil, err
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(content, &doc); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported OpenAPI file extension: %s", ext)
	}

	if len(doc) == 0 {
		return nil, fmt.Errorf("OpenAPI document is empty")
	}
	return doc, nil
}

func validateOpenapiSecuritySchemes(doc map[string]interface{}) error {
	components, ok := toStringMap(doc["components"])
	if !ok {
		return nil
	}
	securitySchemes, ok := toStringMap(components["securitySchemes"])
	if !ok {
		return nil
	}

	for schemeName, rawScheme := range securitySchemes {
		scheme, ok := toStringMap(rawScheme)
		if !ok {
			continue
		}
		schemeType, _ := scheme["type"].(string)
		if strings.EqualFold(strings.TrimSpace(schemeType), "http") {
			if _, hasIn := scheme["in"]; hasIn {
				return fmt.Errorf("invalid components: security scheme %q: security scheme of type %q can't have 'in'", schemeName, "http")
			}
		}
	}

	return nil
}

func toStringMap(v interface{}) (map[string]interface{}, bool) {
	switch m := v.(type) {
	case map[string]interface{}:
		return m, true
	case map[interface{}]interface{}:
		converted := make(map[string]interface{}, len(m))
		for k, value := range m {
			key, ok := k.(string)
			if !ok {
				continue
			}
			converted[key] = value
		}
		return converted, true
	default:
		return nil, false
	}
}

// GetFileInfo retrieves information about an OpenAPI file
func (m *OpenapiFileManager) GetFileInfo(filePath string) (*OpenapiFileInfo, error) {
	// Convert relative path to absolute path
	absFilePath, err := m.ToAbsolutePath(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to absolute path: %v", err)
	}

	// Get file info
	fileInfo, err := os.Stat(absFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %v", err)
	}

	// Determine file type based on extension
	fileType := model.OpenapiFileTypeUnknown
	ext := strings.ToLower(filepath.Ext(absFilePath))
	switch ext {
	case ".json":
		fileType = model.OpenapiFileTypeJson
	case ".yaml", ".yml":
		fileType = model.OpenapiFileTypeYaml
	}

	return &OpenapiFileInfo{
		OpenapiFilePath: filePath,
		OriginalName:    filepath.Base(absFilePath),
		FileSize:        fileInfo.Size(),
		OpenapiFileType: fileType,
	}, nil
}
