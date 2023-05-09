import { defineStore } from 'pinia';
import { getPermissionCodes, setPermissionCodes } from '@/utils/permiss';

export const usePermissStore = defineStore('permiss', {
    state: () => {
        return {
            perrmissionCodes: <string[]>[],
        }
    },
    getters: {
        currentCodes: state => {
            if (state.perrmissionCodes.length === 0) {
                state.perrmissionCodes = getPermissionCodes() ?? []
            }
            return state.perrmissionCodes;
        }
    },
    actions: {
        setCodes(val: string[]) {
            this.perrmissionCodes = val;
            setPermissionCodes(val);
        },
    }
});
