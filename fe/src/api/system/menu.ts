import request from '@/utils/request';
import { IFrontendRoute } from '@/api/system/frontend-route'

export interface IMenu {
    id?: number;
    name: string;
    icon?: string;
    route_id: number;
    parent_id: number;
    enabled: boolean;
    sort: number;
    children?: IMenu[];
    route?: IFrontendRoute;
}

export const listTreeApi = (): IRespData<Required<IMenu>[]> => request.get("/Menu/ListTree")
export const createApi = (data: IMenu): IRespData<null> => request.post("/Menu/Add", data)
export const deleteApi = (id: number): IRespData<null> => request.delete("/Menu/Delete", { params: { id } })
export const editApi = (data: IMenu): IRespData<null> => request.post("/Menu/Edit", data)
export const updateStatus = (data: { id: number, enabled: boolean }): IRespData<null> => request.post("/Menu/UpdateStatus", data)
export const getApi = (id: number): IRespData<Required<IMenu>> => request.get("/Menu/Get", { params: { id } })


// 递归查询指定菜单列表中是否存在路由路径等于指定路由路径的菜单
export function existsByRoutePath(menus: IMenu[] | undefined, path: string): boolean {
    if (menus === undefined) return false;
	for (let menuItem of menus) {
		if (menuItem.route?.route_path === path) {
			return true
		}
		if (menuItem.children && existsByRoutePath(menuItem.children, path)) {
			return true
		}
	}
	return false
}