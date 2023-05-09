import axios, { AxiosInstance, AxiosError, AxiosResponse, InternalAxiosRequestConfig } from 'axios';
import { setLoginInfo, getLoginInfo, logout } from '@/utils/login'
import { sleepAsync } from '@/utils/generic'
import { ElMessage } from 'element-plus';
import { refreshLoginApi } from '@/api/app'

const service: AxiosInstance = axios.create({
    baseURL: `${window.location.protocol}//${import.meta.env.VITE_BASE_URL}`,
    timeout: 120000
});

service.interceptors.request.use(
    async (config: InternalAxiosRequestConfig<any>) => {
        config.headers = config.headers || {};
        let loginInfo = getLoginInfo()
        if (loginInfo) {
            config.headers["Authorization"] = "Bearer " + loginInfo.token;
            if (loginInfo.refresh_after && loginInfo.refresh_after < (Date.now() / 1000) && loginInfo.access_expire > (Date.now() / 1000) && config.url?.toLowerCase() !== "/app/refreshlogin") {
                let res = await refreshLoginApi()
                if (res.code === 1) {
                    setLoginInfo(res.data)
                    console.log("refresh token success!")
                }
            }
        }
        config.headers["Content-Type"] = "application/json";
        return config;
    },
    (error: AxiosError) => {
        console.log(error);
        return Promise.reject();
    }
);

service.interceptors.response.use(
    (response: AxiosResponse) => {
        console.log("response: ", response.data)
        if (response.status === 200) {
            if (response.data.code === 0) {
                ElMessage.error(response.data.msg);
                Promise.reject();
                throw new Error('request failed');
            }
            return response.data;
        } else {
            Promise.reject();
        }
    },
    async (error: AxiosError) => {
        if (error.response?.status === 401) {
            ElMessage.error("登录信息已过期，请重新登录");
            await sleepAsync(1500)
            logout();
            return Promise.reject(error);
        } else if (error.response?.status === 403) {
            ElMessage.error("没有权限");
            return Promise.reject(error);
        }
        ElMessage.error(error.message);
        return Promise.reject(error);
    }
);

export default service;
