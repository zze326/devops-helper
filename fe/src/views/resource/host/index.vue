<template>
	<el-row :gutter="6">
		<el-col :span="4">
			<ElCard>
				<el-text>分组列表</el-text>
				<el-switch v-model="hostStore.permissionManage" v-permiss="permiss.permissionRead" style="float: right; margin-top: -5px;" inline-prompt active-text="终端权限"
					inactive-text="管理主机" />
				<el-divider />
				<Group />
			</ElCard>
		</el-col>
		<el-col :span="20">
			<ElCard>
				<Permission v-if="hostStore.permissionManage"/>
				<Host v-else/>
			</ElCard>
		</el-col>
	</el-row>
</template>

<script setup lang="ts">
import Group from './group.vue'
import Host from './host.vue'
import Permission from './permission.vue'
import { useHostStore } from '@/store/resource/host'

const enum permiss {
	permissionRead = 'host-group-permission-readonly',
}

const hostStore = useHostStore()
</script>

<style scoped>
.el-row,
.el-col,
.el-card {
	height: 100%;
}


</style>