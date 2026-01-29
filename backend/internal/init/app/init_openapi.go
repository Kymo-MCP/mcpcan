package app

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/kymo-mcp/mcpcan/pkg/common"
	"github.com/kymo-mcp/mcpcan/pkg/database/model"
	"github.com/kymo-mcp/mcpcan/pkg/database/repository/mysql"
	"github.com/kymo-mcp/mcpcan/pkg/logger"
	"github.com/kymo-mcp/mcpcan/pkg/openapifile"
)

// initOpenapi 初始化OpenAPI文档
func (a *App) initOpenapi(ctx context.Context) error {
	logger.Info("Starting OpenAPI document initialization")

	apiFile := "./init-data/openapi-file/mcpcan-openapi.json"
	if _, err := os.Stat(apiFile); os.IsNotExist(err) {
		logger.Info("OpenAPI document source file does not exist, skipping initialization",
			zap.String("apiFile", apiFile))
		return nil
	}

	file, err := os.Open(apiFile)
	if err != nil {
		return fmt.Errorf("failed to open openapi file: %w", err)
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat openapi file: %w", err)
	}

	fileName := filepath.Base(apiFile)
	header := &multipart.FileHeader{
		Filename: fileName,
		Size:     fileStat.Size(),
		Header:   make(map[string][]string),
	}
	fileReader := &fileReaderWrapper{file: file}

	codeConfig := &common.CodeConfig{
		Upload: common.UploadConfig{
			MaxFileSize: 100,
		},
	}
	fileManager := openapifile.NewOpenapiFileManager(codeConfig, a.config.Storage.OpenapiFilePath)

	fileInfo, err := fileManager.UploadOpenapiFile(fileReader, header)
	if err != nil {
		return fmt.Errorf("failed to upload openapi file: %w", err)
	}

	if err := fileManager.ValidateOpenapiFile(fileInfo.OpenapiFilePath); err != nil {
		fileManager.DeleteOpenapiFile(fileInfo.OpenapiFilePath)
		return fmt.Errorf("failed to validate openapi file: %w", err)
	}

	existingPackages, err := mysql.McpOpenapiPackageRepo.FindAll(ctx)
	if err != nil {
		fileManager.DeleteOpenapiFile(fileInfo.OpenapiFilePath)
		return fmt.Errorf("failed to query openapi packages: %w", err)
	}

	var existingBase *model.McpOpenapiPackage
	for _, pkg := range existingPackages {
		if pkg.OriginalName == fileName && pkg.BaseOpenapiFileID == "" {
			existingBase = pkg
			break
		}
	}

	if existingBase != nil {
		if _, err := os.Stat(existingBase.OpenapiFilePath); err == nil {
			if err := fileManager.DeleteOpenapiFile(fileInfo.OpenapiFilePath); err != nil {
				logger.Warn("Failed to delete duplicated openapi file",
					zap.String("filePath", fileInfo.OpenapiFilePath),
					zap.Error(err))
			}
			logger.Info("OpenAPI base package already exists, skipping initialization",
				zap.String("openapiFileId", existingBase.OpenapiFileID),
				zap.String("filePath", existingBase.OpenapiFilePath))
			return nil
		}
	}

	for _, pkg := range existingPackages {
		if pkg.OriginalName == fileName && pkg.BaseOpenapiFileID == "" {
			if err := fileManager.DeleteOpenapiFile(pkg.OpenapiFilePath); err != nil {
				logger.Warn("Failed to delete old openapi file",
					zap.String("openapiFileId", pkg.OpenapiFileID),
					zap.String("filePath", pkg.OpenapiFilePath),
					zap.Error(err))
			}
			if err := mysql.McpOpenapiPackageRepo.DeleteByOpenapiFileID(ctx, pkg.OpenapiFileID); err != nil {
				logger.Warn("Failed to delete old openapi record",
					zap.String("openapiFileId", pkg.OpenapiFileID),
					zap.Error(err))
			}
		}
	}

	openapiPackage := &model.McpOpenapiPackage{
		OpenapiFileID:     fileInfo.OpenapiFileID,
		OpenapiFileType:   fileInfo.OpenapiFileType,
		OpenapiFilePath:   fileInfo.OpenapiFilePath,
		OriginalName:      fileInfo.OriginalName,
		FileSize:          fileInfo.FileSize,
		BaseOpenapiFileID: "",
	}

	if err := mysql.McpOpenapiPackageRepo.Create(ctx, openapiPackage); err != nil {
		fileManager.DeleteOpenapiFile(fileInfo.OpenapiFilePath)
		return fmt.Errorf("failed to save openapi package: %w", err)
	}

	return nil
}
