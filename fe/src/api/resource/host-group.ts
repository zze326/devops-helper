import request from '@/utils/request';
import { IHost } from './host';
import { IHostGroupPermission } from './host_group_permission';

export interface IHostGroup {
    id?: number
    name: string
    parent_id: number
    hosts?: IHost[]
    host_group_permissions?: IHostGroupPermission[]
    children?: IHostGroup[]
}

// 创建主机组
export const addApi = (data: IHostGroup): IRespData<Required<IHostGroup>> => request.post("/HostGroup/Add", data)
// 获取主机组树
export const listTreeApi = (): IRespData<Required<IHostGroup>[]> => request.get("/HostGroup/ListTree")
// 获取主机组服务器树
export const listTreeWithHostsApi = (): IRespData<Required<IHostGroup>[]> => request.get("/HostGroup/ListTreeWithHosts")
// 重命名主机组
export const renameApi = (data: Partial<IHostGroup>): IRespData<null> => request.put("/HostGroup/Rename", data)
// 删除主机组
export const deleteApi = (id: number): IRespData<null> => request.delete("/HostGroup/Delete", { params: { id } })