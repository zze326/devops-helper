import request from '@/utils/request';
import { IHostGroup } from './host-group';

export interface IHostGroupPermission {
    id?: number
    type: number
    ref_id: number
    host_groups?: IHostGroup[]
}

export interface IHostGroupPermissionListItem {
    id: number
    type: number
    ref_id: number
    host_group_names_str: string
    ref_name: string
}

// 创建主机组权限
export const addApi = (data: IHostGroupPermission): IRespData<Required<IHostGroupPermission>> => request.post("/HostGroupPermission/Add", data)
// 获取主机组权限
export const listPageApi = (data: IPager, hostGroupID: number): IRespData<IPageData<Required<IHostGroupPermissionListItem>>> => request.post("/HostGroupPermission/ListPage", data, { params: { hostGroupID } })
// 编辑主机组权限
export const editApi = (data: Partial<IHostGroupPermission>): IRespData<null> => request.put("/HostGroupPermission/Edit", data)
// 删除主机组权限
export const deleteApi = (id: number): IRespData<null> => request.delete("/HostGroupPermission/Delete", { params: { id } })
// 获取主机组权限
export const getApi = (id: number): IRespData<Required<IHostGroupPermission>> => request.get("/HostGroupPermission/Get", { params: { id } })