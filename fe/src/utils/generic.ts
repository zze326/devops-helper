import { ElMessageBox, ElLoading } from "element-plus"
import { usePermissStore } from "@/store/permiss"
import { getLoginInfo } from "@/utils/login"
import { App } from "vue"

export const sleepAsync = async function (time: number) {
    return new Promise((resolve) => setTimeout(resolve, time))
}

export const confirm = async (message: string, title = '提示', confirmText = '确定', cancelText = '取消') => {
    return await ElMessageBox.confirm(message, title, {
        confirmButtonText: confirmText,
        cancelButtonText: cancelText,
        type: 'warning'
    }).catch(() => false)
}

export const openLoading = (message = '加载中...') => {
    const loading = ElLoading.service({
        lock: true,
        text: message,
        background: 'rgba(0, 0, 0, 0.7)',
    })
    return loading
}

// 2023-04-04T18:00:19+08:00 => 2023-04-04 18:00:19 
export const formatDate = (date: string) => {
    const inputDate = new Date(date);
    const year = inputDate.getFullYear();
    const month = (inputDate.getMonth() + 1).toString().padStart(2, '0');
    const day = inputDate.getDate().toString().padStart(2, '0');
    const hours = inputDate.getHours().toString().padStart(2, '0');
    const minutes = inputDate.getMinutes().toString().padStart(2, '0');
    const seconds = inputDate.getSeconds().toString().padStart(2, '0');
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
};

export const hasPermission = (code: string): boolean => {
    if (getLoginInfo()?.userinfo.is_super) {
        return true
    }
    const permissStore = usePermissStore()
    return permissStore.currentCodes.includes(code);
}

export const setGlobalFunc = (app: App<Element>) => {
    // 注册全局函数，支持转换 2023-04-04T18:00:19+08:00 => 2023-04-04 18:00:19 
    app.config.globalProperties.$formatDate = formatDate;
    app.config.globalProperties.$hasPermission = hasPermission;
}

export const getWsProtocol = (url: string | null = null): string => {
    // 如果没有指定 URL，则使用当前页面的 URL
    if (!url) {
        url = window.location.href;
    }

    // 根据 URL 的协议返回对应的 WebSocket 协议
    const isSecure = url.startsWith("https");
    return isSecure ? "wss" : "ws";
};

export const getHttpProtocol = (url: string | null = null): string => {
    // 如果没有指定 URL，则使用当前页面的 URL
    if (!url) {
        url = window.location.href;
    }

    // 根据 URL 的协议返回对应的 WebSocket 协议
    const isSecure = url.startsWith("https");
    return isSecure ? "https" : "http";
}