import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router';
import { getLoginInfo } from '@/utils/login';
import {setPermissionCodes} from '@/utils/permiss';
import { usePermissStore } from '@/store/permiss';
import { App } from 'vue';
import { listRoutesAndPermissionCodesForCurrentUserApi } from '@/api/app';
import Home from '@/views/home.vue';
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

let modules = import.meta.glob("@/views/**/*.vue");

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        redirect: '/dashboard',
    },
    {
        path: '/login',
        name: 'Login',
        meta: {
            title: '登录',
        },
        component: modules['/src/views/login.vue'],
    },
    {
        path: '/403',
        name: '403',
        meta: {
            title: '没有权限',
        },
        component: modules['/src/views/403.vue'],
    },
    {
        path: '/404',
        name: '404',
        meta: {
            title: '找不到资源',
        },
        component: modules['/src/views/404.vue'],
    },
    {
        path: '/resource/host-terminal-single',
        name: 'host-terminal-single',
        meta: {
            title: '主机终端 - 单主机',
        },
        component: modules['/src/views/resource/host/terminal-single.vue'],
    },
    {
        path: '/resource/host-terminal-tree',
        name: 'host-terminal-tree',
        meta: {
            title: '主机终端 - 树状目录',
        },
        component: modules['/src/views/resource/host/terminal-tree.vue'],
    },
];

const router = createRouter({
    history: createWebHashHistory(),
    routes,
});

router.beforeEach(async (to, from, next) => {
    NProgress.start()
    document.title = `${to.meta.title} | devops helper`;
    let loginInfo = getLoginInfo();
    if (!loginInfo && to.path !== '/login') {
        next('/login');
    } else {
        next();
    }
    // else if (to.meta.permiss && !permiss.key.includes(to.meta.permiss)) {
    //     // 如果没有权限，则进入403
    //     next('/403');
    // } 
});

router.afterEach(() => {
    NProgress.done()
})

export const setHomeRoutes = async () => {
    if (router.hasRoute('Home')) {
        return;
    }
    let { data } = await listRoutesAndPermissionCodesForCurrentUserApi();
    let respRoutes = data.routes;
    let permissionCodes = data.permission_codes;
    router.addRoute({
        path: '/',
        name: 'Home',
        component: Home,
        children: respRoutes.map(item => {
            return {
                path: item.route_path,
                name: item.name,
                meta: {
                    title: item.title,
                    permiss: item.id,
                },
                component: modules[`/src/views${item.component_path}`]
            };
        })
    })
    try {
        usePermissStore().setCodes(permissionCodes);
    } catch (error) {
        setPermissionCodes(permissionCodes);
    }
}

if (getLoginInfo()) {
    await setHomeRoutes();
}


export const resetHomeRoutes = async () => {
    if (router.hasRoute('Home')) {
        router.removeRoute('Home');
    }
    await setHomeRoutes();
}

export const initRouter = (app: App<Element>) => {
    app.use(router);
}
