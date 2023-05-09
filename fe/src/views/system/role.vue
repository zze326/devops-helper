<template>
	<div>
		<ElCard>
			<PagerTable :loading="tableDataLoading" :pager-state="pagerState" placeholder="角色/代码" @list-page="listPage"
				@add-click="handleAdd" :add-permiss-code="permiss.add">
				<el-table-column prop="id" label="ID" />
				<el-table-column prop="name" label="名称" />
				<el-table-column prop="code" label="代码" />
				<el-table-column label="操作" v-if="$hasPermission(permiss.edit) || $hasPermission(permiss.delete)">
					<template #default="scope">
						<el-button icon="Edit" type="primary" v-permiss="permiss.edit"
							@click="handleEdit(scope.row)">编辑</el-button>
						<el-button icon="Delete" type="danger" v-permiss="permiss.delete"
							@click="handleDelete(scope.row)">删除</el-button>
					</template>
				</el-table-column>
			</PagerTable>
		</ElCard>

		<ElDialog v-model="editDialog.visible" :title="editDialog.title" width="30%">
			<ElForm ref="editFormRef" :model="editFormData" status-icon :rules="editForm.rules" label-width="120px">
				<ElFormItem label="名称" prop="name">
					<el-input v-model="editFormData.name" type="text" />
				</ElFormItem>
				<ElFormItem label="代码" prop="code">
					<el-input v-model="editFormData.code" />
				</ElFormItem>

				<ElFormItem label="权限" prop="permission_ids">
					<el-tree-select :default-expanded-keys="editFormData.permission_ids" :props="{ label: 'labelWithPath' }"
						@change="handlePermissionTreeSelectChange" style="width: 100%;"
						v-model="editFormData.permission_ids" highlight-current clearable
						:data="editDialog.permissionTreeSelectData" multiple check-strictly placeholder="请选择"
						:filterable="true" :render-after-expand="false">
						<template #default="{ node, data }">
							<span :class="data.checked ? 'checked-item' : ''">{{ data.label }}</span>
						</template>
					</el-tree-select>
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

<script setup lang="ts" name="route">
import { ref, reactive, onMounted, toRefs, nextTick } from 'vue'
import { addApi, editApi, deleteApi, getApi, listPageApi, IRole } from '@/api/system/role'
import { listTreeApi as listTreePermissions } from '@/api/system/permission'
import { PagerState } from '@/utils/pager'
import { confirm } from '@/utils/generic'
import { refershTreeSelectDataWithChecked, genTreeSelectDataWithChecked } from '@/utils/treeSelect'
import _ from 'lodash'
import { ElMessage, FormInstance } from 'element-plus'

const enum permiss {
	add = 'role-add',
	edit = 'role-edit',
	delete = 'role-delete',
}

const editFormRef = ref<FormInstance>()
const state = reactive<{
	editFormData: IRole,
	editDialog: {
		visible: boolean,
		title: string,
		add: boolean,
		permissionTreeSelectData: ITreeSelectData[],
	},
	tableDataLoading: boolean,
}>({
	editFormData: {
		name: '',
		code: '',
	},
	editDialog: {
		visible: false,
		title: '新增角色',
		add: true,
		permissionTreeSelectData: [],
	},
	tableDataLoading: false,
})

const { editFormData, editDialog, tableDataLoading } = toRefs(state)

const editForm = {
	rules: {
		name: [
			{ required: true, message: '请输入角色名称', trigger: 'blur' }
		],
		code: [
			{ required: true, message: '请输入角色代码', trigger: 'blur' }
		],
	},
	submit: () => {
		editFormRef.value?.validate(async (valid: boolean) => {
			if (!valid) {
				return false
			}
			state.editDialog.add ? await addApi(state.editFormData) : await editApi(state.editFormData as Required<IRole>)
			await listPage()
			ElMessage.success(`${state.editDialog.add ? '新增' : '编辑'}成功`);
			state.editDialog.visible = false
		})
	},
}

const pagerState = reactive(new PagerState<IRole>(1, 10))

onMounted(() => {
	listPage()
})
// 分页
const listPage = async () => {
	state.tableDataLoading = true
	let res = await listPageApi(pagerState.getPager(['name', 'code']))
	pagerState.items = res.data.items
	pagerState.total = res.data.total
	state.tableDataLoading = false
}
// 新增
const handleAdd = async () => {
	state.editDialog.add = true
	state.editDialog.title = '新增角色'
	_.assign(state.editFormData, {
		id: undefined,
		name: '',
		code: '',
		permission_ids: [],
	});

	let res = await listTreePermissions()
	state.editDialog.permissionTreeSelectData = genTreeSelectDataWithChecked(res.data, [])
	state.editDialog.visible = true
	nextTick(() => {
		editFormRef.value?.clearValidate()
	})
}

// 编辑
const handleEdit = async (row: Required<IRole>) => {
	state.editDialog.add = false
	state.editDialog.title = '编辑角色'
	let res = await getApi(row.id)
	state.editFormData = _.pick(res.data, ['id', 'name', 'code'])
	state.editFormData.permission_ids = res.data.permissions.map((item) => item.id || 0)
	let res2 = await listTreePermissions()
	state.editDialog.permissionTreeSelectData = genTreeSelectDataWithChecked(res2.data, state.editFormData.permission_ids)
	state.editDialog.visible = true
}

// 删除
const handleDelete = async (row: Required<IRole>) => {
	if (!await confirm("确认删除？")) {
		return
	}
	await deleteApi(row.id)
	await listPage()
}

// 权限树选择
const handlePermissionTreeSelectChange = (value: number[]) => {
	state.editDialog.permissionTreeSelectData = refershTreeSelectDataWithChecked(state.editDialog.permissionTreeSelectData, value)
}

</script>
<style scoped>
.checked-item {
	color: #008000;
}
</style>
