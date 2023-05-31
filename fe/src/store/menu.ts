import { IMenu, existsByRoutePath as existsMenuByRoutePath } from '@/api/system/menu';
import { defineStore } from 'pinia';

export const useMenuStore = defineStore('menu', {
    state: () => {
        return {
            allMenus: <Required<IMenu>[]>[],
            topMenus: <ITopMenu[]>[],
            sidebarMenus: <ISidebarMenu>{},
            activeTopMenuID: <number>0
        }
    },
    getters: {
        currentSidebarMenus: state => {
            if (!state.activeTopMenuID) {
                return [];
            }
            return state.sidebarMenus[state.activeTopMenuID];
        }
    },
    actions: {
        setAllMenus(val: Required<IMenu>[]) {
            this.allMenus = val;
        },
        setTopMenus(val: ITopMenu[]) {
            this.topMenus = val;
        },
        setSidebarMenus(val: ISidebarMenu) {
            this.sidebarMenus = val;
        },
        setActiveTopMenuID(val: number) {
            this.activeTopMenuID = val;
        },
        routeChange(path: string) {
            this.allMenus.forEach(menuItem => {
                if (existsMenuByRoutePath(menuItem.children, path)) {
                    this.setActiveTopMenuID(menuItem.id)
                }
            })
        }
    }
});
