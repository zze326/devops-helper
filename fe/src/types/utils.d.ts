import { ComponentCustomProperties } from "@vue/runtime-core";

declare module '@vue/runtime-core' {
    interface ComponentCustomProperties {
        $formatDate: (date: string) => string;
        $hasPermission: (code: string) => boolean;
    }
}
export default ComponentCustomProperties;