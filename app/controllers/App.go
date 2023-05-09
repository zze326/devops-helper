package controllers

import (
	"fmt"
	"github.com/revel/revel"
	o_system "github.com/zze326/devops-helper/app/models/orm/system"
	v_system "github.com/zze326/devops-helper/app/models/view/system"
	gormc "github.com/zze326/devops-helper/app/modules/gormc/app/controllers"
	"github.com/zze326/devops-helper/app/results"
	"github.com/zze326/devops-helper/app/startups"
	"github.com/zze326/devops-helper/app/utils"
	"gorm.io/gorm"
	"sort"
)

type App struct {
	gormc.Controller
}

func (c App) Ping() revel.Result {
	return results.JsonOk()
}

// Login 登录
func (c App) Login(loginData v_system.AppLogin) revel.Result {
	userModel := new(o_system.User)
	if err := c.DB.First(userModel, "username = ?", loginData.Username).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return results.JsonError(fmt.Errorf("用户不存在"))
		}
		return results.JsonError(err)
	}

	if userModel.Password != utils.EncodeMD5(loginData.Password) {
		return results.JsonError(fmt.Errorf("密码错误"))
	}

	return createLoginResult(userModel)
}

// RefreshLogin 刷新 Token
func (c App) RefreshLogin(_requestUserID int) revel.Result {
	userModel := new(o_system.User)
	if err := c.DB.First(userModel, _requestUserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return results.JsonError(fmt.Errorf("用户不存在"))
		}
		return results.JsonError(err)
	}
	return createLoginResult(userModel)
}

// createLoginResult 创建登录结果
func createLoginResult(userModel *o_system.User) revel.Result {
	token, accessExpire, refreshAfter, err := utils.GenJwtToken(revel.Config.StringDefault("jwt.secret", ""), int64(revel.Config.IntDefault("jwt.expire", 86400)), uint(userModel.ID), userModel.Username, userModel.RealName)
	if err != nil {
		return results.JsonError(err)
	}

	type userinfo struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		RealName string `json:"real_name"`
	}

	type jwtResp struct {
		Token        string   `json:"token"`
		Userinfo     userinfo `json:"userinfo"`
		AccessExpire int64    `json:"access_expire"`
		RefreshAfter int64    `json:"refresh_after"`
	}

	return results.JsonOkData(jwtResp{
		Token:        token,
		AccessExpire: accessExpire,
		RefreshAfter: refreshAfter,
		Userinfo: userinfo{
			ID:       uint(userModel.ID),
			Username: userModel.Username,
			RealName: userModel.RealName,
		},
	})
}

// ListRoutesAndPermissionCodesForCurrentUser 获取当前用户的前端路由和权限
func (c App) ListRoutesAndPermissionCodesForCurrentUser(_requestUserID int) revel.Result {
	requestUser := new(o_system.User)
	if err := c.DB.Preload("Roles.Permissions").First(requestUser, _requestUserID).Error; err != nil {
		return results.JsonError(err)
	}

	_, permissionMap, err := o_system.Permission{}.ListTree(c.DB, false, true, false)
	if err != nil {
		return results.JsonError(err)
	}

	frontendRouteMap := make(map[int]*o_system.FrontendRoute, len(requestUser.Roles)*5)
	permissionCodeMap := make(map[string]bool, len(requestUser.Roles)*5)
	permissionCodes := make([]string, 0)
	for _, role := range requestUser.Roles {
		for _, permission := range role.Permissions {
			loopAddSubPermissionFrontendRoute(permissionMap[permission.ID], frontendRouteMap, permissionCodeMap)
		}
	}

	for permissionCode, _ := range permissionCodeMap {
		permissionCodes = append(permissionCodes, permissionCode)
	}

	var frontendRouteModels []*o_system.FrontendRoute
	for _, route := range frontendRouteMap {
		frontendRouteModels = append(frontendRouteModels, route)
	}

	return results.JsonOkData(struct {
		FrontendRoutes  []*o_system.FrontendRoute `json:"routes"`
		PermissionCodes []string                  `json:"permission_codes"`
	}{
		FrontendRoutes:  frontendRouteModels,
		PermissionCodes: permissionCodes,
	})
}

func loopAddSubPermissionFrontendRoute(permission *o_system.Permission, frontendRouteMap map[int]*o_system.FrontendRoute, permissionCodeMap map[string]bool) {
	for _, route := range permission.FrontendRoutes {
		if _, exists := frontendRouteMap[route.ID]; !exists {
			frontendRouteMap[route.ID] = route
		}
	}

	if !permissionCodeMap[permission.Code] {
		permissionCodeMap[permission.Code] = true
	}

	if len(permission.Children) > 0 {
		for _, child := range permission.Children {
			loopAddSubPermissionFrontendRoute(child, frontendRouteMap, permissionCodeMap)
		}
	}
}

// ListTreeMenusForCurrentUser 获取当前用户的菜单树
func (c App) ListTreeMenusForCurrentUser(_requestUserID int) revel.Result {
	requestUser := new(o_system.User)
	if err := c.DB.Preload("Roles.Permissions").First(requestUser, _requestUserID).Error; err != nil {
		return results.JsonError(err)
	}

	var menus []*o_system.Menu
	// 从数据库中获取所有菜单记录
	if err := c.DB.Preload("Route").Find(&menus).Error; err != nil {
		return results.JsonError(err)
	}

	// 构建 ID 到权限映射的 map
	menuMap := make(map[int]*o_system.Menu)
	for _, menu := range menus {
		menuMap[menu.ID] = menu
	}

	_, permissionMap, err := o_system.Permission{}.ListTree(c.DB, true, false, false)
	if err != nil {
		return results.JsonError(err)
	}

	uniqueMenuIDs := make(map[int]bool, len(requestUser.Roles)*5)
	for _, role := range requestUser.Roles {
		for _, permission := range role.Permissions {
			// 添加关联权限及子权限关联的菜单
			loopAddSubPermissionMenu(uniqueMenuIDs, permissionMap[permission.ID])
		}
	}

	var relatedMenus = make(map[int]*o_system.Menu, len(uniqueMenuIDs))
	for k, _ := range uniqueMenuIDs {
		// 添加关联菜单的父菜单
		loopAddParentMenu(menuMap[k], relatedMenus, menuMap)
	}

	// 找出根节点并返回
	var roots []*o_system.Menu
	// 对每个权限，将其作为子节点添加到对应父节点的 Children 列表中
	for _, menu := range relatedMenus {
		parent, ok := menuMap[menu.ParentID]
		if !ok {
			// 如果该权限没有父节点，则认为它是根节点
			roots = append(roots, menu)
			continue
		}
		parent.Children = append(parent.Children, menu)
	}

	// 排序
	for _, menu := range menus {
		if menu.Children != nil {
			sort.Slice(menu.Children, func(i, j int) bool {
				return menu.Children[i].Sort < menu.Children[j].Sort
			})
		}
	}

	sort.Slice(roots, func(i, j int) bool {
		return roots[i].Sort < roots[j].Sort
	})

	return results.JsonOkData(roots)
}

// RefreshPermission 刷新权限
func (c App) RefreshPermission() revel.Result {
	// 刷新后端权限 - casbin
	if err := startups.RefreshCasbin(); err != nil {
		return results.JsonError(err)
	}
	return results.JsonOk()
}

func loopAddParentMenu(menu *o_system.Menu, menuMap map[int]*o_system.Menu, allMenuMap map[int]*o_system.Menu) {
	if _, exists := menuMap[menu.ID]; !exists {
		menuMap[menu.ID] = menu
	}

	if menu.ParentID != 0 {
		loopAddParentMenu(allMenuMap[menu.ParentID], menuMap, allMenuMap)
	}
}

func loopAddSubPermissionMenu(uniqueMenuIDs map[int]bool, permission *o_system.Permission) {
	for _, menu := range permission.Menus {
		if !uniqueMenuIDs[menu.ID] {
			uniqueMenuIDs[menu.ID] = true
		}
	}

	if len(permission.Children) > 0 {
		for _, child := range permission.Children {
			loopAddSubPermissionMenu(uniqueMenuIDs, child)
		}
	}
}
