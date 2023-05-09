import request from '@/utils/request';

export interface IPermission {
    id?: number;
    name: string;
    code: string;
    parent_id?: number;
    children_ids?: number[];
    children?: IPermission[];
    menu_ids?: number[];
    frontend_route_ids?: number[];
    backend_route_ids?: number[];
}

// 获取权限
export const getApi = (id: number): IRespData<Required<IPermission>> => request.get("/Permission/Get", { params: { id } })
// 创建权限
export const addApi = (data: IPermission): IRespData<null> => request.post("/Permission/Add", data)
// 编辑权限
export const editApi = (data: IPermission): IRespData<null> => request.post("/Permission/Edit", data)
// 查询所有顶级权限
export const listAllTopApi = (): IRespData<Required<IPermission>[]> => request.get("/Permission/ListAllTop")
// 查询权限树
export const listTreeApi = (): IRespData<Required<IPermission>[]> => request.get("/Permission/ListTree")
// 删除权限
export const deleteApi = (id: number): IRespData<null> => request.delete("/Permission/Delete", { params: { id } })