import request from '@/utils/request';

export interface IFrontendRoute {
    id?: number
    name: string;
    title: string;
    component_path: string;
    route_path: string;
}

// 创建前端路由
export const addApi = (data: IFrontendRoute): IRespData<null> => request.post("/FrontendRoute/Add", data)
// 编辑前端路由
export const editApi = (data: Required<IFrontendRoute>): IRespData<null> => request.post("/FrontendRoute/Edit", data)
// 删除前端路由
export const deleteApi = (id: number): IRespData<null> => request.post("/FrontendRoute/Delete", { id })
// 根据 ID 获取前端路由
export const getApi = (id: number): IRespData<Required<IFrontendRoute>> => request.get("/FrontendRoute/Get", { params: { id } })
// 分页查询前端路由
export const listPageApi = (data: IPager): IRespData<IPageData<Required<IFrontendRoute>>> => request.post("/FrontendRoute/ListPage", data)
// 查询所有前端路由
export const listAllApi = (): IRespData<Required<IFrontendRoute>[]> => request.get("/FrontendRoute/ListAll")