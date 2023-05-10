<template>
	<div>
		<ElCard>
			<PagerTable :loading="tableDataLoading" :pager-state="pagerState" @list-page="listPage"
				:add-permiss-code="permiss.add" @add-click="handleAdd" placeholder="名称/路由路径">
				<el-table-column prop="id" label="ID" width="100" />
				<el-table-column prop="name" label="名称" />
				<el-table-column prop="path" label="路由路径" min-width="150" />
				<el-table-column label="操作" v-if="$hasPermission(permiss.edit) || $hasPermission(permiss.delete)">
					<template #default="scope">
						<el-button type="primary" icon="Edit" v-permiss="permiss.edit"
							@click="handleEdit(scope.row)">编辑</el-button>
						<el-button type="danger" icon="Delete" v-permiss="permiss.delete"
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
				<ElFormItem label="路由路径" prop="path">
					<el-input v-model.number="editFormData.path" />
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
import { addApi, editApi, deleteApi, getApi, listPageApi, IBackendRoute } from '@/api/system/backend-route'
import { PagerState } from '@/utils/pager'
import { confirm } from '@/utils/generic'
import _ from 'lodash'
import { ElMessage, FormInstance } from 'element-plus'

const enum permiss {
	add = 'backend-route-add',
	edit = 'backend-route-edit',
	delete = 'backend-route-delete',
}

const editFormRef = ref<FormInstance>()
const state = reactive<IPageState<IBackendRoute>>({
	editFormData: {
		name: '',
		path: '',
	},
	editDialog: {
		visible: false,
		title: '新增路由',
		add: true,
	},
	tableDataLoading: false,
})

const { editFormData, editDialog, tableDataLoading } = toRefs(state)

const editForm = {
	rules: {
		name: [
			{ required: true, message: '请输入名称', trigger: 'blur' }
		],
		path: [
			{ required: true, message: '请输入路由路径', trigger: 'blur' }
		],
	},
	submit: () => {
		editFormRef.value?.validate(async (valid: boolean) => {
			if (!valid) {
				return false
			}
			state.editDialog.add ? await addApi(state.editFormData) : await editApi(state.editFormData as Required<IBackendRoute>)
			await listPage()
			ElMessage.success(`${state.editDialog.add ? '新增' : '编辑'}成功`);
			state.editDialog.visible = false
		})
	},
}

const pagerState = reactive(new PagerState<IBackendRoute>(1, 10))

onMounted(() => {
	listPage()
})
// 分页
const listPage = async () => {
	state.tableDataLoading = true
	let res = await listPageApi(pagerState.getPager(['name', 'path']))
	pagerState.items = res.data.items
	pagerState.total = res.data.total
	state.tableDataLoading = false
}
// 新增
const handleAdd = () => {
	state.editDialog.add = true
	state.editDialog.title = '新增路由'
	_.assign(state.editFormData, {
		name: '',
		path: '',
	});
	state.editDialog.visible = true
	nextTick(() => {
		editFormRef.value?.clearValidate()
	})
}

// 编辑
const handleEdit = async (row: Required<IBackendRoute>) => {
	state.editDialog.add = false
	state.editDialog.title = '编辑路由'
	let res = await getApi(row.id)
	state.editFormData = _.pick(res.data, ['id', 'name', 'path'])
	state.editDialog.visible = true
}

// 删除
const handleDelete = async (row: Required<IBackendRoute>) => {
	if (!await confirm("确认删除？")) {
		return
	}
	await deleteApi(row.id)
	await listPage()
}

</script>
<style scoped></style>
