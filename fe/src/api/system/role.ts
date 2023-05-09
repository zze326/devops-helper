import request from '@/utils/request';
import { IPermission } from '@/api/system/permission';

export interface IRole {
    id?: number
    name: string;
    code: string;
    permissions?: IPermission[]
    permission_ids?: number[]
}

// 创建角色
export const addApi = (data: IRole): IRespData<null> => request.post("/Role/Add", data)
// 编辑角色
export const editApi = (data: Required<IRole>): IRespData<null> => request.post("/Role/Edit", data)
// 删除角色
export const deleteApi = (id: number): IRespData<null> => request.post("/Role/Delete", { id })
// 根据 ID 获取角色
export const getApi = (id: number): IRespData<Required<IRole>> => request.get("/Role/Get", { params: { id } })
// 分页查询角色
export const listPageApi = (data: IPager): IRespData<IPageData<Required<IRole>>> => request.post("/Role/ListPage", data)
// 查询所有角色
export const listAllApi = (): IRespData<Required<IRole>[]> => request.get("/Role/ListAll")