import request from '@/utils/request';
import { IRole } from '@/api/system/role';

export interface IUser {
    id?: number
    username: string;
    phone: string;
    email: string;
    real_name: string;
    roles?: IRole[];
}

// 创建用户
export const addApi = (data: IUser): IRespData<null> => request.post("/User/Add", data)
// 编辑用户
export const editApi = (data: Required<IUser>): IRespData<null> => request.post("/User/Edit", data)
// 删除用户
export const deleteApi = (id: number): IRespData<null> => request.post("/User/Delete", { id })
// 获取用户
export const getApi = (id: number): IRespData<Required<IUser>> => request.get("/User/Get", { params: { id } })
// 分页查询用户
export const listPageApi = (data: IPager): IRespData<IPageData<Required<IUser>>> => request.post("/User/ListPage", data)
// 重置用户密码
export const resetPwd = (data: { id: number, new_password: string }): IRespData<null> => request.put("/User/ResetPwd", data)
// 查询所有用户
export const listAllApi = (): IRespData<Required<IUser>[]> => request.get("/User/ListAll")