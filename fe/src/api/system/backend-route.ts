import request from '@/utils/request';

export interface IBackendRoute {
    id?: number
    name: string;
    path: string;
}

// 创建后端路由
export const addApi = (data: IBackendRoute): IRespData<null> => request.post("/BackendRoute/Add", data)
// 编辑后端路由
export const editApi = (data: Required<IBackendRoute>): IRespData<null> => request.post("/BackendRoute/Edit", data)
// 删除后端路由
export const deleteApi = (id: number): IRespData<null> => request.post("/BackendRoute/Delete", { id })
// 根据 ID 获取后端路由
export const getApi = (id: number): IRespData<Required<IBackendRoute>> => request.get("/BackendRoute/Get", { params: { id } })
// 分页查询后端路由
export const listPageApi = (data: IPager): IRespData<IPageData<Required<IBackendRoute>>> => request.post("/BackendRoute/ListPage", data)
// 查询所有后端路由
export const listAllApi = (): IRespData<Required<IBackendRoute>[]> => request.get("/BackendRoute/ListAll")