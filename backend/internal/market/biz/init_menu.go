package biz

import (
	"context"
	"fmt"
	"strings"

	"github.com/kymo-mcp/mcpcan/internal/market/config"
	"github.com/kymo-mcp/mcpcan/pkg/database/model"
	"github.com/kymo-mcp/mcpcan/pkg/database/repository/mysql"
	"github.com/kymo-mcp/mcpcan/pkg/menu"
)

func (a *App) createAdminRoleMenus(adminUser *model.SysUser) error {
	menuSyncPolicy := strings.ToLower(strings.TrimSpace(a.config.Init.MenuSyncPolicy))
	if menuSyncPolicy == config.InitMenuSyncPolicyOff {
		fmt.Println("Menu sync is disabled by init policy, skipping admin menu sync.")
		return nil
	}

	existingRole, err := mysql.SysRoleRepo.FindByName(context.Background(), adminUser.GetUsername())
	if err != nil {
		return err
	}
	if existingRole == nil {
		return fmt.Errorf("admin role not found")
	}

	menus := menu.GetMenus(a.config.CodeMode)

	menuList, _, err := mysql.SysMenuRepo.FindWithPagination(context.Background(), 1, 9999, "")
	if err != nil {
		return fmt.Errorf("failed to find menus: %w", err)
	}

	var mcpcanMenus []*model.SysMenu
	for _, dbMenu := range menuList {
		if !strings.HasPrefix(dbMenu.GetPermission(), "mcpcan") {
			continue
		}
		mcpcanMenus = append(mcpcanMenus, dbMenu)
	}

	// 根据 Permission 过滤出 db 需要新增和删除的菜单
	dbMenusToAdd, err := findAddMenu(context.Background(), menus, mcpcanMenus)
	if err != nil {
		return fmt.Errorf("failed to find add menu: %w", err)
	}

	dbMenusToDelete := make([]*model.SysMenu, 0)
	if menuSyncPolicy == config.InitMenuSyncPolicyFullSync {
		dbMenusToDelete = findDeleteMenu(context.Background(), menus, mcpcanMenus)
	}

	var menuIds []int64
	for _, dbMenuDelete := range dbMenusToDelete {
		menuIds = append(menuIds, dbMenuDelete.MenuID)
	}

	if len(menuIds) > 0 {
		// 删除数据库中的菜单, 和关联该菜单的角色
		if err := mysql.SysMenuRepo.BatchDelete(context.Background(), menuIds); err != nil {
			return fmt.Errorf("failed to delete menus: %w", err)
		}
		if err := mysql.SysRolesMenusRepo.BatchDeleteByMenuID(context.Background(), menuIds); err != nil {
			return fmt.Errorf("failed to delete roles menus: %w", err)
		}
	}

	for _, addMenu := range dbMenusToAdd {
		// 递归创建
		if err := createMenu(context.Background(), adminUser, existingRole, addMenu, nil); err != nil {
			return fmt.Errorf("failed to create menu: %w", err)
		}
	}

	if menuSyncPolicy == config.InitMenuSyncPolicyFullSync {
		for _, m := range menuChildToList(menus) {
			dbMenu, err := mysql.SysMenuRepo.FindByPermission(context.Background(), m.Permission)
			if err != nil {
				continue
			}
			if dbMenu.GetPath() != m.Path || dbMenu.GetEngTitle() != m.EngTitle ||
				dbMenu.GetTitle() != m.Title || dbMenu.GetMenuSort() != int64(m.Sort) {
				dbMenu.Title = &m.Title
				dbMenu.EngTitle = &m.EngTitle
				dbMenu.MenuSort = &m.Sort
				dbMenu.Path = &m.Path
				err := mysql.SysMenuRepo.Update(context.Background(), dbMenu)
				if err != nil {
					return fmt.Errorf("failed to update menu: %w", err)
				}
			}
		}
	}

	return nil
}

func menuChildToList(menus []*menu.Menu) []*menu.Menu {
	var leafMenus []*menu.Menu
	for _, m := range menus {
		if m.Children != nil {
			leafMenus = append(leafMenus, menuChildToList(m.Children)...)
			leafMenus = append(leafMenus, m)
		} else {
			leafMenus = append(leafMenus, m)
		}
	}
	return leafMenus
}

func findDeleteMenu(ctx context.Context, menus []*menu.Menu, mcpcanMenus []*model.SysMenu) []*model.SysMenu {
	var dbMenusToDelete []*model.SysMenu
	for _, dbMenu := range mcpcanMenus {
		found := false
		for _, m := range menuChildToList(menus) {
			if dbMenu.GetPermission() == m.Permission {
				found = true
				break
			}
		}
		if !found {
			dbMenusToDelete = append(dbMenusToDelete, dbMenu)
		}
	}
	return dbMenusToDelete
}

func findAddMenu(ctx context.Context, menus []*menu.Menu, mcpcanMenus []*model.SysMenu) ([]*menu.Menu, error) {
	var dbMenusToAdd []*menu.Menu = []*menu.Menu{}
	for _, m := range menus {
		permission := m.Permission
		found := false
		for _, dbMenu := range mcpcanMenus {
			if dbMenu.GetPermission() == permission {
				found = true
				break
			}
		}

		if m.Children != nil {
			child, err := findAddMenu(ctx, m.Children, mcpcanMenus)
			if err != nil {
				return nil, fmt.Errorf("failed to find add child menu: %w", err)
			}
			// 如果父类不需要创建，则塞入子类中
			if !found {
				m.Children = child
			} else {
				dbMenusToAdd = append(dbMenusToAdd, child...)
			}
		}

		if !found {
			dbMenusToAdd = append(dbMenusToAdd, m)
		}
	}
	return dbMenusToAdd, nil
}

func createMenu(ctx context.Context, adminUser *model.SysUser, existingRole *model.SysRole, menu *menu.Menu, parentMenu *model.SysMenu) error {
	create := &model.SysMenu{
		Permission: &menu.Permission,
		Title:      &menu.Title,
		Type:       &menu.Type,
		MenuSort:   &menu.Sort,
		Path:       &menu.Path,
		EngTitle:   &menu.EngTitle,
		SubCount:   len(menu.Children),
		CreateBy:   adminUser.Username,
		UpdateBy:   adminUser.Username,
	}
	if parentMenu != nil {
		create.PID = &parentMenu.MenuID
	}
	// 创建菜单
	if err := mysql.SysMenuRepo.Create(ctx, create); err != nil {
		return fmt.Errorf("failed to create menu: %w", err)
	}
	// 管理员角色关联菜单
	if err := mysql.SysRolesMenusRepo.BatchCreate(ctx, []*model.SysRolesMenus{
		{
			RoleID: int64(existingRole.RoleID),
			MenuID: create.MenuID,
		},
	}); err != nil {
		return fmt.Errorf("failed to create roles menus: %w", err)
	}

	// 递归遍历子菜单并创建
	for _, child := range menu.Children {
		if err := createMenu(ctx, adminUser, existingRole, child, create); err != nil {
			return fmt.Errorf("failed to create child menu: %w", err)
		}
	}
	return nil
}
