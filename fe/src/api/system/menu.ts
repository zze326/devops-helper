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
