interface IResp<T> {
    code: number;
    msg: string;
    data: T;
}

type IRespData<T> = Promise<IResp<T>>

interface IPageData<T> {
    total: number;
    items: T[];
}

interface IPager {
    page: number;
    page_size: number;
    wheres: {
        logic: "or" | "and";
        columns: {
            column: string;
            op: string;
            value: string;
        }[]
    }[]
}

interface IPagerState<T> {
    page: number;
    page_size: number;
    search?: string;
    total: number;
    items: T[];
    getPager(columns?: string[]): IPager;
}

interface ILoginInfo {
    token: string;
    userinfo: {
        id: number;
        username: string;
        real_name: string;
    }
    access_expire: number;
    refresh_after: number;
}

interface ITreeSelectData{
    label: string,
    value?: number,
    children?: ITreeSelectData[],
    checked?: boolean,
    nameWithPath?: string,
}

interface IPageState<T> {
    editFormData: T,
    editDialog: { visible: boolean, title: string, add: boolean },
    tableDataLoading: boolean,
    tableData?: T[],
}