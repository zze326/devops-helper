import { createApp } from 'vue';
import { createPinia } from 'pinia';
import * as ElementPlusIconsVue from '@element-plus/icons-vue';
import App from '@/App.vue';
import { initRouter } from '@/router';
import 'element-plus/dist/index.css';
import '@/assets/css/icon.css';
import { setGlobalFunc, hasPermission } from '@/utils/generic'
import { usePermissStore } from '@/store/permiss';

const app = createApp(App);

setGlobalFunc(app)
app.use(createPinia());
initRouter(app)
// 注册elementplus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component);
}

console.log(import.meta.env)

app.directive('permiss', {
    mounted(el, binding) {
        if (!binding.value) {
            return
        }

        if (!hasPermission(binding.value)) {
            el['hidden'] = true;
        }
    },
});

app.mount('#app');

