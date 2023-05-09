<template>
	<div>
		<ElCard>
			<ElForm ref="queryForm" :inline="true">
				<ElFormItem v-permiss='permiss.add'>
					<ElButton type="success" icon="Plus" @click="handleAdd()">添加</ElButton>
				</ElFormItem>
				<ElFormItem v-permiss='permiss.refresh'>
					<ElButton type="warning" icon="Refresh" @click="handleRefresh()">刷新系统权限
					</ElButton>
				</ElFormItem>
			</ElForm>

			<el-table v-loading="tableDataLoading" :data="tableData" style="width: 100%; margin-bottom: 20px" row-key="id"
				border :default-expand-all="false">
				<el-table-column prop="name" label="名称" />
				<el-table-column prop="code" label="代码" />
				<el-table-column label="创建时间">
					<template #default="scope">
						{{ $formatDate(scope.row.created_at) }}
					</template>
				</el-table-column>
				<el-table-column label="操作"
					v-if="$hasPermission(permiss.add) || $hasPermission(permiss.edit) || $hasPermission(permiss.delete)">
					<template #default="scope">
						<ElButton icon="Plus" type="success" v-permiss='permiss.add' @click="handleAdd(scope.row.id)">添加
						</ElButton>
						<el-button icon="Edit" type="primary" v-permiss='permiss.edit'
							@click="handleEdit(scope.row)">编辑</el-button>
						<el-button icon="Delete" type="danger" v-permiss='permiss.delete'
							@click="handleDelete(scope.row)">删除</el-button>
					</template>
				</el-table-column>
			</el-table>
		</ElCard>
		<ElDialog v-model="editDialog.visible" :title="editDialog.title" width="32%" @closed="resetEditFormData">
			<ElForm ref="editFormRef" :model="editFormData" status-icon :rules="editForm.rules" label-width="120px">
				<ElFormItem label="权限名称" prop="name">
					<el-input v-model="editFormData.name" type="text" />
				</ElFormItem>
				<ElFormItem label="权限代码" prop="code">
					<el-input v-model.trim="editFormData.code" type="text" />
				</ElFormItem>
				<ElFormItem label="父权限" prop="parent_id">
					<el-tree-select style="width: 100%;" v-model="editFormData.parent_id" clearable
						:data="editDialog.permissionTreeSelectData" check-strictly placeholder="请选择" :filterable="true"
						:render-after-expand="false" />
				</ElFormItem>
				<ElFormItem label="子权限" prop="children_ids">
					<el-select style="width: 100%;" multiple clearable v-model="editFormData.children_ids" filterable
						placeholder="请选择">
						<el-option v-for="item in editDialog.permissionSelectData" :key="item.id" :value="item.id"
							:label="`${item.name} - ${item.code}`">
						</el-option>
					</el-select>
				</ElFormItem>
				<ElFormItem label="菜单集合" prop="menu_ids">
					<el-tree-select @check-change="handleMenuSelectChange" style="width: 100%;"
						v-model="editFormData.menu_ids" clearable multiple show-checkbox
						:data="editDialog.menuTreeSelectData" check-strictly placeholder="请选择" :filterable="true"
						:render-after-expand="false" />
				</ElFormItem>
				<ElFormItem label="前端路由" prop="frontend_route_ids">
					<el-select style="width: 100%;" multiple clearable v-model="editFormData.frontend_route_ids" filterable
						placeholder="请选择">
						<el-option v-for="item in editDialog.frontendRouteSelectData" :key="item.id" :value="item.id"
							:label="`${item.title} - ${item.route_path}`">
						</el-option>
					</el-select>
				</ElFormItem>
				<ElFormItem label="后端路由" prop="backend_route_ids">
					<el-select style="width: 100%;" multiple clearable v-model="editFormData.backend_route_ids" filterable
						placeholder="请选择">
						<el-option v-for="item in editDialog.backendRouteSelectData" :key="item.id" :value="item.id"
							:label="`${item.name} - ${item.path}`">
						</el-option>
					</el-select>
				</ElFormItem>
			</ElForm>
			<template #footer>
				<span class="dialog-footer">
					<ElButton @click="editDialog.visible = false">取消</ElButton>
					<ElButton type="primary" @click="editForm.submit">
						确认
					</ElButton>
				</span>
			</template>
		</ElDialog>
	</div>
</template>

<script setup lang="ts" name="system-user">
import { ref, reactive, onMounted, toRefs, nextTick } from 'vue'
import { confirm } from '@/utils/generic'
import _ from 'lodash'
import { ElMessage, FormInstance } from 'element-plus'

import { addApi, listAllTopApi, listTreeApi, getApi, editApi, deleteApi, IPermission } from '@/api/system/permission'
import { listTreeApi as listTreeMenus, getApi as getMenu } from '@/api/system/menu';
import { listAllApi as listAllFrontendRoutes, IFrontendRoute } from '@/api/system/frontend-route';
import { listAllApi as listAllBackendRoutes, IBackendRoute } from '@/api/system/backend-route'
import { genTreeSelectData } from '@/utils/treeSelect'
import { refreshPermissionApi } from '@/api/app'


const enum permiss {
	add = 'permission-add',
	edit = 'permission-edit',
	delete = 'permission-delete',
	refresh = 'permission-refresh'
}

const editFormRef = ref<FormInstance>()
const state = reactive<{
	editFormData: IPermission,
	editDialog: {
		visible: boolean,
		title: string,
		add: boolean,
		permissionSelectData: Required<IPermission>[],
		frontendRouteSelectData: Required<IFrontendRoute>[],
		backendRouteSelectData: Required<IBackendRoute>[],
		menuTreeSelectData: ITreeSelectData[],
		permissionTreeSelectData: ITreeSelectData[],
	},
	tableData: Required<IPermission>[],
	tableDataLoading: boolean,
}>({
	editFormData: {
		id: 0,
		name: '',
		code: '',
		parent_id: undefined,
		children_ids: [],
		menu_ids: [],
		frontend_route_ids: [],
		backend_route_ids: []
	},
	editDialog: {
		visible: false,
		title: '新增权限',
		add: true,
		permissionSelectData: [],
		frontendRouteSelectData: [],
		backendRouteSelectData: [],
		menuTreeSelectData: [],
		permissionTreeSelectData: [],
	},
	tableData: [],
	tableDataLoading: false,
})

const { editFormData, editDialog, tableData, tableDataLoading } = toRefs(state)

const editForm = {
	rules: {
		name: [
			{ required: true, message: '请输入权限名称', trigger: 'blur' }
		],
		code: [
			{ required: true, message: '请输入权限标题', trigger: 'blur' }
		],
	},
	submit: () => {
		editFormRef.value?.validate(async (valid: boolean) => {
			if (!valid) {
				return false
			}
			state.editDialog.add ? await addApi(state.editFormData) : await editApi(state.editFormData as Required<IPermission>)
			await listTableData()
			ElMessage.success(`${state.editDialog.add ? '新增' : '编辑'}成功`);
			state.editDialog.visible = false
		})
	},
}

const resetEditFormData = () => {
	_.assign(state.editFormData, {
		id: 0,
		name: '',
		code: '',
		parent_id: undefined,
		children_ids: [],
		menu_ids: [],
		frontend_route_ids: [],
		backend_route_ids: []
	})
}

const handleAdd = async (to?: number) => {
	state.editDialog.add = true;
	state.editDialog.title = '新增权限';
	state.editFormData.parent_id = to
	let [permissionRes, menuRes, frontendRouteRes, backendRouteRes] = await Promise.all([listAllTopApi(), listTreeMenus(), listAllFrontendRoutes(), listAllBackendRoutes()])
	state.editDialog.permissionSelectData = permissionRes.data
	state.editDialog.menuTreeSelectData = genTreeSelectData(menuRes.data)
	state.editDialog.frontendRouteSelectData = frontendRouteRes.data
	state.editDialog.backendRouteSelectData = backendRouteRes.data
	state.editDialog.permissionTreeSelectData = genTreeSelectData(state.tableData)
	state.editDialog.visible = true;

	nextTick(() => {
		editFormRef.value?.clearValidate()
	})
}

const handleEdit = async (row: Required<IPermission>) => {
	state.editDialog.add = false;
	state.editDialog.title = '编辑权限';
	let res = await getApi(row.id)
	state.editFormData = _.pick(res.data, ['id', 'name', 'code', 'parent_id', 'children_ids', 'menu_ids', 'frontend_route_ids', 'backend_route_ids'])
	state.editFormData.parent_id = state.editFormData.parent_id || undefined
	let [permissionRes, menuRes, frontendRouteRes, backendRouteRes] = await Promise.all([listAllTopApi(), listTreeMenus(), listAllFrontendRoutes(), listAllBackendRoutes()])
	state.editDialog.permissionSelectData = _.cloneDeep(permissionRes.data)
	state.editDialog.menuTreeSelectData = genTreeSelectData(menuRes.data)
	state.editDialog.frontendRouteSelectData = frontendRouteRes.data
	state.editDialog.backendRouteSelectData = backendRouteRes.data
	state.editDialog.permissionTreeSelectData = genTreeSelectData(state.tableData)
	state.editDialog.permissionSelectData = _.concat(state.editDialog.permissionSelectData, res.data.children as Required<IPermission>[])
	state.editDialog.permissionSelectData = _.filter(state.editDialog.permissionSelectData, (item) => {
		return item.id !== state.editFormData.id && item.id !== state.editFormData.parent_id
	})

	state.editDialog.visible = true;
}

const handleDelete = async (row: Required<IPermission>) => {
	let res = await confirm(`确定删除权限【${row.name}】吗？`)
	if (!res) {
		return
	}
	await deleteApi(row.id)
	await listTableData()
	ElMessage.success('删除成功');
}

const handleRefresh = async () => {
	await refreshPermissionApi()
	ElMessage.success('刷新成功');
}

const handleMenuSelectChange = async (val: Required<ITreeSelectData>, isSelect: boolean) => {
	let res = await getMenu(val.value)
	if (res.data.route_id === 0) return;
	if (isSelect) {
		if (state.editFormData.frontend_route_ids?.includes(res.data.route_id)) {
			return
		}
		state.editFormData.frontend_route_ids = _.concat(state.editFormData.frontend_route_ids || [], res.data.route_id)
	}
	else {
		state.editFormData.frontend_route_ids = _.filter(state.editFormData.frontend_route_ids, (item) => {
			return item !== res.data.route_id
		})
	}
}

onMounted(async () => {
	await listTableData()
})

const listTableData = async () => {
	tableDataLoading.value = true
	let res = await listTreeApi()
	state.tableData = res.data
	tableDataLoading.value = false
}


</script>

<style scoped></style>
