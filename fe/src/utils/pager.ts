export class Pager implements IPager {
    constructor(page: number, page_size: number, search?: string, search_columns?: string[], select_columns?: string[]) {
        this.page = page;
        this.page_size = page_size;
        this.search = search;
        this.select_columns = select_columns
        if (search && search_columns) {
            this.setSearch(search, search_columns)
        }
    }

    page: number;
    page_size: number;
    search?: string;
    wheres: {
        logic: "or" | "and";
        columns: {
            column: string;
            op: string;
            value: string;
        }[];
    }[] = [];
    select_columns?: string[]


    setSearch(value: string, columns: string[]) {
        this.wheres.push({
            logic: "or",
            columns: columns.map((column) => {
                return {
                    column,
                    op: "~",
                    value
                }
            }
            )
        })
    }

    addWhere = (logic: "or" | "and", column: string, op: string, value: string) => {
        this.wheres.push({
            logic,
            columns: [{
                column,
                op,
                value
            }]
        })
    }
}


export class PagerState<T> implements IPagerState<T> {
    page: number;
    page_size: number;
    search?: string;
    total: number;
    items: T[];

    constructor(page: number, page_size: number, search: string | undefined = undefined) {
        this.page = page;
        this.page_size = page_size;
        this.search = search;
        this.total = 0;
        this.items = [];
    }

    getPager = (search_columns: string[] = ['name'], select_columns: string[] = []): IPager => {
        return new Pager(this.page, this.page_size, this.search, search_columns, select_columns);
    }
}