<template>
  <div class="sidebar">
    <el-menu class="sidebar-el-menu" :default-active="onRoutes" :collapse="sidebar.collapse" background-color="#324157"
      text-color="#bfcbd9" active-text-color="#20a0ff" unique-opened router>
      <template v-for="item in menu.currentSidebarMenus">
        <template v-if="item.subs && item.subs.length > 0">
          <el-sub-menu :index="item.index" :key="item.index">
            <template #title>
              <el-icon>
                <component :is="item.icon"></component>
              </el-icon>
              <span>{{ item.name }}</span>
            </template>
            <template v-for="subItem in item.subs">
              <el-sub-menu v-if="subItem.subs && subItem.subs.length > 0" :index="subItem.index" :key="subItem.index">
                <template #title>{{ subItem.name }}</template>
                <el-menu-item v-for="(threeItem, i) in subItem.subs" :key="i" :index="threeItem.index">
                  {{ threeItem.name }}
                </el-menu-item>
              </el-sub-menu>
              <el-menu-item v-else :index="subItem.index">
                {{ subItem.name }}
              </el-menu-item>
            </template>
          </el-sub-menu>
        </template>
        <template v-else>
          <el-menu-item :index="item.index" :key="item.index">
            <el-icon>
              <component :is="item.icon"></component>
            </el-icon>
            <template #title>{{ item.name }}</template>
          </el-menu-item>
        </template>
      </template>
    </el-menu>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useSidebarStore } from '@/store/sidebar';
import { useMenuStore } from '@/store/menu';
import { useRoute } from 'vue-router';

const menu = useMenuStore()
const route = useRoute();
const onRoutes = computed(() => {
  return route.path;
});

const sidebar = useSidebarStore();
</script>

<style scoped>
.sidebar {
  display: block;
  position: absolute;
  left: 0;
  top: 70px;
  bottom: 0;
  overflow-y: scroll;
}

.sidebar::-webkit-scrollbar {
  width: 0;
}

.sidebar-el-menu:not(.el-menu--collapse) {
  width: 250px;
}

.el-menu--collapse {
  width: 65px;
}

.sidebar>ul {
  height: 100%;
}
</style>
