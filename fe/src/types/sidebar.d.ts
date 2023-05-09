interface IMenuItem {
    icon?: string;
    index: string;
    name: string;
    permiss?: number;
    subs?: IMenuItem[];
}


interface ITopMenu {
    id: number;
    icon?: string;
    name: string;
}

interface ISidebarMenu {
    [key: number]: IMenuItem[]
}