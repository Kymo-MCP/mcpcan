package biz

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kymo-mcp/mcpcan/internal/authz/biz"
	"github.com/kymo-mcp/mcpcan/internal/market/config"
	"github.com/kymo-mcp/mcpcan/pkg/database/model"
	"github.com/kymo-mcp/mcpcan/pkg/database/repository/mysql"
)

func (a *App) createAdminUser() (*model.SysUser, error) {
	ctx := context.Background()
	userBiz := biz.NewUserBiz()
	initConfig := config.GetConfig()

	// Prepare admin user parameters
	now := time.Now()
	username := initConfig.Init.AdminUsername
	password := initConfig.Init.AdminPassword
	nickname := initConfig.Init.AdminNickname
	roleName := initConfig.Init.AdminRoleName
	level := initConfig.Init.AdminRoleLevel
	dataScope := initConfig.Init.AdminDataScope
	deptName := initConfig.Init.AdminDeptName
	if roleName == "" {
		roleName = username
	}
	if level == 0 {
		level = 1
	}
	if dataScope == "" {
		dataScope = string(model.DataScopeAll)
	}
	if deptName == "" {
		deptName = "总部"
	}
	isAdmin := true
	enabled := true

	adminUser := &model.SysUser{
		Username:   &username,
		NickName:   &nickname,
		IsAdmin:    isAdmin,
		Enabled:    &enabled,
		Source:     stringPtr("PLATFORM"),
		CreateTime: &now,
		UpdateTime: &now,
	}

	adminDept := &model.SysDept{
		Name:    deptName,
		Enabled: 1,
		Source:  model.DeptSourcePlatform,
	}

	adminRole := &model.SysRole{
		Name:        roleName,
		Description: &roleName,
		Level:       &level,
		DataScope:   &dataScope,
		CreateTime:  &now,
		UpdateTime:  &now,
	}

	// Check if admin user already exists
	existingUser, err := mysql.SysUserRepo.FindByUsername(ctx, *adminUser.Username)
	if err == nil && existingUser != nil {
		passwordPolicy := strings.ToLower(strings.TrimSpace(initConfig.Init.AdminPasswordPolicy))
		if passwordPolicy == config.InitAdminPasswordPolicyForceReset {
			fmt.Println("Admin user already exists, resetting password by init policy...")
			err = userBiz.SetUserPassword(ctx, existingUser, password)
			if err != nil {
				return nil, fmt.Errorf("failed to update admin password: %v", err)
			}
		} else {
			fmt.Println("Admin user already exists, skip password reset by init policy.")
		}

		// Ensure admin dept is created and linked
		adminDept, _ = createAdminDept(ctx, adminDept)
		if adminDept != nil {
			if existingUser.DeptID == nil || *existingUser.DeptID == 0 {
				existingUser.DeptID = &adminDept.DeptID
				_ = mysql.SysUserRepo.Update(ctx, existingUser)
			}
		}

		a.AdminUser = existingUser
		return existingUser, nil
	}

	// Create admin role
	adminRole, err = createAdminRole(ctx, adminRole)
	if err != nil {
		return nil, fmt.Errorf("failed to create admin role: %v", err)
	}

	// Create admin dept
	adminDept, err = createAdminDept(ctx, adminDept)
	if err != nil {
		return nil, fmt.Errorf("failed to create admin dept: %v", err)
	}
	adminUser.DeptID = &adminDept.DeptID

	// Create user
	err = userBiz.CreateUser(ctx, adminUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create admin user: %v", err)
	}

	// Re-fetch user to ensure UserID is set correctly
	adminUser, err = mysql.SysUserRepo.FindByUsername(ctx, *adminUser.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created admin user: %v", err)
	}

	// Set user password
	err = userBiz.SetUserPassword(ctx, adminUser, password)
	if err != nil {
		return nil, fmt.Errorf("failed to set admin password: %v", err)
	}

	// Assign admin role to user
	err = userBiz.AssignRolesToUser(ctx, adminUser.UserID, []uint{adminRole.RoleID})
	if err != nil {
		return nil, fmt.Errorf("failed to assign admin role to user: %v", err)
	}

	fmt.Printf("Admin user created successfully with ID: %d\n", adminUser.UserID)
	a.AdminUser = adminUser
	return adminUser, nil
}

func createAdminRole(ctx context.Context, adminRole *model.SysRole) (*model.SysRole, error) {
	// Check if admin role already exists
	existingRole, err := mysql.SysRoleRepo.FindByName(ctx, adminRole.Name)
	if err == nil && existingRole != nil {
		return existingRole, nil
	}

	// Create role
	err = mysql.SysRoleRepo.Create(ctx, adminRole)
	if err != nil {
		return nil, fmt.Errorf("failed to create admin role: %v", err)
	}

	fmt.Printf("Admin role created successfully with ID: %d\n", adminRole.RoleID)
	return adminRole, nil
}

func createAdminDept(ctx context.Context, adminDept *model.SysDept) (*model.SysDept, error) {
	// Check if admin dept already exists
	existingDept, err := mysql.SysDeptRepo.FindByName(ctx, adminDept.Name)
	if err == nil && existingDept != nil {
		return existingDept, nil
	}

	// Create dept
	err = mysql.SysDeptRepo.Create(ctx, adminDept)
	if err != nil {
		return nil, fmt.Errorf("failed to create admin dept: %v", err)
	}

	fmt.Printf("Admin dept created successfully with ID: %d\n", adminDept.DeptID)
	return adminDept, nil
}

func stringPtr(s string) *string {
	return &s
}
