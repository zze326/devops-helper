import { defineStore } from 'pinia';

export const useHostStore = defineStore('resource.host', {
    state: () => {
        return {
            activeGroupID: <number | null>null, // 当前激活的分组ID
            reloadGroupCounter: 0, // 重载分组计数器，用于触发分组列表的重新加载
            permissionManage: false // 是否处于权限管理面板
        };
    },
    actions: {
        activeGroup(groupID: number | null) {
            this.activeGroupID = groupID
        },
        reloadGroup() {
            this.reloadGroupCounter++
        }
    }
});
