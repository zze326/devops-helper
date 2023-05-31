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

interface IWhereColumn {
    column: string
    op: string
    value: number | string | boolean
}

interface IPager {
    page: number;
    page_size: number;
    wheres: {
        logic: "or" | "and";
        columns: IWhereColumn[]
    }[]
    addWhere: (logic: "or" | "and", ...columns: IWhereColumn[]) => void;
    setSearch: (value: string, columns: string[]) => void;
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
        is_super: boolean;
    }
    access_expire: number;
    refresh_after: number;
}

interface ITreeSelectData {
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

interface IDataDictItem {
    id?: number
    label: string
    value: number
    sort: number
    data_dict_id: number
}
