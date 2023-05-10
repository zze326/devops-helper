import request from '@/utils/request';
import { IHostGroup } from './host-group';

export interface IHost {
    id?: number
    name: string
    host: string
    port: number
    username: string
    password?: string
    private_key?: string
    use_key: boolean
    desc?: string
    host_groups?: IHostGroup[]
    host_group_ids?: number[]
    save_session?: boolean
    update_password_and_private_key?: boolean
}

// 创建主机
export const addApi = (data: IHost): IRespData<Required<IHost>> => request.post("/Host/Add", data)
// 编辑主机
export const editApi = (data: IHost): IRespData<IHost> => request.put("/Host/Edit", data)
// 分页获取主机
export const listPageApi = (data: IPager, hostGroupID: number): IRespData<IPageData<Required<IHost>>> => request.post("/Host/ListPage", data, { params: { hostGroupID } })
// 主机
export const deleteApi = (id: number): IRespData<null> => request.delete("/Host/Delete", { params: { id } })
// 获取主机
export const getApi = (id: number): IRespData<Required<IHost>> => request.get("/Host/Get", { params: { id } })
// 获取密码和私钥
export const getPasswordAndPrivateKeyApi = (id: number): IRespData<Partial<IHost>> => request.get("/Host/GetPasswordAndPrivateKey", { params: { id } })
// 测试终端是否可以 ssh 连接
export const testSSHApi = (id: number): IRespData<null> => request.get("/Host/TestSSH", { params: { id } })
// 下载文件
export const downloadFileApi = (id: number, path: string): Promise<Blob> => request.get('/Host/DownloadFile', {
    params: { id, path },
    responseType: 'blob' // 将响应数据类型设置为 Blob 对象。
})