import request from '@/utils/request';
import { IMenu } from '@/api/system/menu';
import { IFrontendRoute } from '@/api/system/frontend-route'

export interface ILoginReq {
    username: string;
    password: string;
}

export interface ILoginPermission {
    routes: Required<IFrontendRoute>[];
    permission_codes: string[];
}

// 登录
export const loginApi = (data: ILoginReq): IRespData<ILoginInfo> => request.post("/App/Login", data)
// 刷新 token
export const refreshLoginApi = (): IRespData<ILoginInfo> => request.get("/App/RefreshLogin")
// 获取当前用户菜单树
export const listTreeMenusForCurrentUserApi = (): IRespData<Required<IMenu>[]> => request.get("/App/ListTreeMenusForCurrentUser")
// 查询当前用户可用的前端路由和权限
export const listRoutesAndPermissionCodesForCurrentUserApi = (): IRespData<Required<ILoginPermission>> => request.get("/App/ListRoutesAndPermissionCodesForCurrentUser")
// 刷新权限
export const refreshPermissionApi = (): IRespData<null> => request.get("/App/RefreshPermission")