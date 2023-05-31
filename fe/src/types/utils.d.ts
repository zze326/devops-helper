import { ComponentCustomProperties } from "@vue/runtime-core";
import { IDataDictItem } from "@/api/system/data-dict-item";

declare module '@vue/runtime-core' {
    interface ComponentCustomProperties {
        $formatDate: (date: string) => string;
        $hasPermission: (code: string) => boolean;
        $listDataDictItems: (typeCode: string) => Promise<IDataDictItem[]>;
    }
}
export default ComponentCustomProperties;