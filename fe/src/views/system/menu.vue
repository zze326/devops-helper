<template>
	<div>
		<ElCard>
			<ElForm ref="queryForm" :inline="true">
				<ElFormItem>
					<ElButton type="success" icon="Plus" @click="handleAdd()" v-permiss="permiss.add">添加</ElButton>
				</ElFormItem>
			</ElForm>
			<el-table v-loading="tableDataLoading" :data="tableData" style="width: 100%; margin-bottom: 20px" row-key="id"
				border :default-expand-all="false">
				<el-table-column prop="name" label="名称" />
				<el-table-column prop="icon" label="图标">
					<template #default="scope">
						<div v-if="scope.row.icon">
							<el-icon style="position: relative; top:2px">
								<component :is="scope.row.icon"></component>
							</el-icon>
							{{ scope.row.icon }}
						</div>
						<el-tag v-else type="info">未设置</el-tag>
					</template>
				</el-table-column>
				<el-table-column prop="sort" label="排序" sortable />
				<el-table-column label="创建时间">
					<template #default="scope">
						{{ $formatDate(scope.row.created_at) }}
					</template>
				</el-table-column>
				<el-table-column label="启用">
					<template #default="scope">
						<el-switch :disabled="!$hasPermission(permiss.enable ?? '')" v-model="scope.row.enabled"
							@change="handleChangeStatus(scope.row)"
							style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949" />
					</template>
				</el-table-column>
				<el-table-column label="操作" min-width="120px"
					v-if="$hasPermission(permiss.edit) || $hasPermission(permiss.delete)">
					<template #default="scope">
						<ElButton icon="Plus" type="success" v-permiss="permiss.add" @click="handleAdd(scope.row.id)">添加
						</ElButton>
						<el-button icon="Edit" type="primary" v-permiss="permiss.edit"
							@click="handleEdit(scope.row)">编辑</el-button>
						<el-button icon="Delete" type="danger" v-permiss="permiss.delete"
							@click="handleDelete(scope.row)">删除</el-button>
					</template>
				</el-table-column>
			</el-table>
		</ElCard>

		<ElDialog v-model="editDialog.visible" :title="editDialog.title" width="30%">
			<ElForm ref="editFormRef" :model="editFormData" status-icon :rules="editForm.rules" label-width="80px">
				<ElFormItem label="上级菜单" prop="parent_id">
					<el-tree-select style="width: 100%;" v-model.number="editFormData.parent_id"
						:data="editDialog.treeSelectData" check-strictly placeholder="请选择上级菜单，不选则为顶级菜单" :filterable="true"
						:render-after-expand="false" />
				</ElFormItem>
				<ElFormItem label="名称" prop="name">
					<el-input v-model.trim="editFormData.name" type="text" />
				</ElFormItem>
				<el-row>
					<el-col :span="15">
						<ElFormItem label="图标" prop="icon">
							<el-select style="width: 100%;" clearable v-model="editFormData.icon" filterable
								placeholder="请选择">
								<template #prefix v-if="editFormData.icon">
									<el-icon>
										<component :is="editFormData.icon"></component>
									</el-icon>
								</template>
								<el-option v-for="[key] of Object.entries(ElementPlusIconsVue)" :key="key" :value="key">
									<el-icon>
										<component :is="key"></component>
									</el-icon> {{ key }}
								</el-option>
							</el-select>
						</ElFormItem>
					</el-col>
					<el-col :span="9">
						<ElFormItem label="排序" prop="sort">
							<el-input-number v-model.number="editFormData.sort" :min="0" :max="999" />
						</ElFormItem>
					</el-col>
				</el-row>
				<el-row>
					<el-col :span="15">
						<ElFormItem label="路由" prop="route_id">
							<el-select style="width: 100%;" clearable v-model="editFormData.route_id" filterable
								placeholder="请选择">
								<el-option v-for="item in editDialog.routeSelectData" :key="item.id" :value="item.id"
									:label="`${item.title} - ${item.route_path}`">
								</el-option>
							</el-select>
						</ElFormItem>
					</el-col>
					<el-col :span="9">
						<ElFormItem label="启用" prop="enabled">
							<el-switch v-model="editFormData.enabled"
								style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949" />
						</ElFormItem>
					</el-col>
				</el-row>
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

<script setup lang="ts" name="sys-menu">
import { reactive, toRefs, ref, onMounted } from 'vue'
import { ElMessage, FormInstance } from 'element-plus'
import { confirm } from '@/utils/generic'
import { listTreeApi, createApi, getApi, editApi, updateStatus, IMenu, deleteApi } from '@/api/system/menu';
import { listAllApi as listAllRoutes, IFrontendRoute } from '@/api/system/frontend-route';
import { genTreeSelectData } from '@/utils/treeSelect'
import * as ElementPlusIconsVue from '@element-plus/icons-vue';
import _ from 'lodash';

const enum permiss {
	add = 'menu-add',
	edit = 'menu-edit',
	delete = 'menu-delete',
	enable = "menu-enable-or-disable",
}

const state = reactive<Required<IPageState<Partial<IMenu>>> & {
	editDialog: {
		routeSelectData: Required<IFrontendRoute>[],
		treeSelectData: ITreeSelectData[],
	},
}
>({
	editDialog: {
		visible: false,
		add: true,
		title: '',
		routeSelectData: [],
		treeSelectData: [],
	},
	editFormData: {
		name: '',
		icon: '',
		route_id: undefined,
		parent_id: undefined,
		enabled: true,
		sort: 0,
	},
	tableData: [],
	tableDataLoading: false,
})

const handleAdd = async (to = undefined) => {
	state.editDialog.add = true
	state.editDialog.title = '新增菜单'
	state.editFormData = {
		name: '',
		icon: '',
		route_id: undefined,
		parent_id: to,
		enabled: true,
		sort: 0,
	}
	let [menuRes, routeRes] = await Promise.all([listTreeApi(), listAllRoutes()])
	state.editDialog.treeSelectData = genTreeSelectData(menuRes.data)
	state.editDialog.routeSelectData = routeRes.data
	state.editDialog.visible = true
}

const handleEdit = async (row: Required<IMenu>) => {
	state.editDialog.add = false
	state.editDialog.title = '编辑菜单'
	let res = await getApi(row.id)
	state.editFormData = _.pick(res.data, ['id', 'name', 'route_id', 'parent_id', 'icon', 'enabled', 'sort'])
	state.editFormData.route_id = state.editFormData.route_id === 0 ? undefined : state.editFormData.route_id
	state.editFormData.parent_id = state.editFormData.parent_id === 0 ? undefined : state.editFormData.parent_id
	let [menuRes, routeRes] = await Promise.all([listTreeApi(), listAllRoutes()])
	state.editDialog.treeSelectData = genTreeSelectData(menuRes.data)
	state.editDialog.routeSelectData = routeRes.data
	state.editDialog.visible = true
	console.log(state.editFormData)
}

const handleChangeStatus = async (row: Required<IMenu>) => {
	await updateStatus({ id: row.id, enabled: row.enabled })
	ElMessage.success(`${row.enabled ? '启用' : '禁用'}成功`)
}

const handleDelete = async (row: Required<IMenu>) => {
	if (!await confirm("确认删除？")) {
		return
	}
	console.log(row)
	await deleteApi(row.id)
	await listTableData()
	ElMessage.success('删除成功')
}

const editFormRef = ref<FormInstance>()
const editForm = {
	rules: {
		name: [
			{ required: true, message: '请输入名称', trigger: 'blur' }
		],
		enabled: [
			{ required: true, message: '请选择是否启用', trigger: 'blur' }
		],
	},
	submit: () => {
		editFormRef.value?.validate(async (valid: boolean) => {
			if (!valid) {
				return false
			}
			state.editDialog.add ? await createApi(state.editFormData as IMenu) : await editApi(state.editFormData as IMenu)
			state.editDialog.visible = false
			ElMessage.success(`${state.editDialog.add ? '新增' : '编辑'}成功`);
			await listTableData()
		})
	},
}

const { editDialog, editFormData, tableData, tableDataLoading } = toRefs(state)

const listTableData = async () => {
	state.tableDataLoading = true
	let res = await listTreeApi()
	tableData.value = res.data
	state.tableDataLoading = false
}

onMounted(async () => {
	await listTableData()
})

</script>
<style scoped></style>
