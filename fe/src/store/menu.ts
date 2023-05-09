import { defineStore } from 'pinia';

export const useMenuStore = defineStore('menu', {
    state: () => {
        return {
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
        setTopMenus(val: ITopMenu[]) {
            this.topMenus = val;
        },
        setSidebarMenus(val: ISidebarMenu) {
            this.sidebarMenus = val;
        },
        setActiveTopMenuID(val: number) {
            this.activeTopMenuID = val;
        }
    }
});
