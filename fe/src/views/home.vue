<template>
	<Header />
	<Sidebar />
	<div class="content-box" :class="{ 'content-collapse': sidebar.collapse }">
		<Tags></Tags>
		<router-view class="content" v-slot="{ Component }">
			<transition name="move" mode="out-in">
				<keep-alive :include="tags.nameList">
					<component :is="Component"></component>
				</keep-alive>
			</transition>
		</router-view>
	</div>
</template>
<script setup lang="ts">
import { onBeforeMount } from 'vue';
import { useSidebarStore } from '@/store/sidebar';
import { useTagsStore } from '@/store/tags';
import { useMenuStore } from '@/store/menu';
import { listTreeMenusForCurrentUserApi } from '@/api/app'
import { IMenu } from '@/api/system/menu';



const sidebar = useSidebarStore();
const tags = useTagsStore();
const menu = useMenuStore();


// 将菜单数据转换为树形结构
function toMenuItems(menus: IMenu[]): IMenuItem[] {
	return menus.map(menu => {
		const menuItem: IMenuItem = {
			icon: menu.icon,
			index: menu.route?.route_path || menu.id?.toString() || '',
			name: menu.name,
			permiss: menu.id,
		}
		if (menu.children) {
			menuItem.subs = toMenuItems(menu.children);
		}
		return menuItem;
	})
}

onBeforeMount(async () => {
	const res = await listTreeMenusForCurrentUserApi()
	// 设置顶部菜单
	menu.setTopMenus(res.data.map(menuItem => {
		return <ITopMenu>{
			id: menuItem.id,
			icon: menuItem.icon,
			name: menuItem.name,
		}
	}))
	// 设置侧边栏菜单
	const sidebarMenus: ISidebarMenu = {};
	res.data.forEach(topMenu => {
		const children = toMenuItems(topMenu.children ?? []);
		sidebarMenus[topMenu.id] = children;
	});
	menu.setSidebarMenus(sidebarMenus)
	menu.setActiveTopMenuID(menu.topMenus[0].id);
})
</script>
