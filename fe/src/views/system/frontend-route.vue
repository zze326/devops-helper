<template>
	<div>
		<ElCard>
			<PagerTable :loading="tableDataLoading" :pager-state="pagerState" @list-page="listPage" @add-click="handleAdd"
				:add-permiss-code="permiss.add" placeholder="名称/标题">
				<el-table-column prop="id" label="ID" />
				<el-table-column prop="name" label="名称" />
				<el-table-column prop="title" label="标题" />
				<el-table-column prop="route_path" label="路由路径" />
				<el-table-column prop="component_path" label="组件路径" />
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
				<ElFormItem label="标题" prop="title">
					<el-input v-model="editFormData.title" type="text" />
				</ElFormItem>
				<ElFormItem label="路由路径" prop="route_path">
					<el-input v-model.number="editFormData.route_path" />
				</ElFormItem>
				<ElFormItem label="组件路径" prop="component_path">
					<el-input v-model.number="editFormData.component_path" />
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
import { addApi, editApi, deleteApi, getApi, listPageApi, IFrontendRoute } from '@/api/system/frontend-route'
import { PagerState } from '@/utils/pager'
import { confirm } from '@/utils/generic'
import _ from 'lodash'
import { ElMessage, FormInstance } from 'element-plus'

const enum permiss {
	add = 'frontend-route-add',
	edit = 'frontend-route-edit',
	delete = 'frontend-route-delete',
}

const editFormRef = ref<FormInstance>()
const state = reactive<IPageState<IFrontendRoute>>({
	editFormData: {
		name: '',
		title: '',
		route_path: '',
		component_path: ''
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
		title: [
			{ required: true, message: '请输入标题', trigger: 'blur' }
		],
		route_path: [
			{ required: true, message: '请输入路由路径', trigger: 'blur' }
		],
		component_path: [
			{ required: true, message: '请输入组件路径', trigger: 'blur' }
		]
	},
	submit: () => {
		editFormRef.value?.validate(async (valid: boolean) => {
			if (!valid) {
				return false
			}
			state.editDialog.add ? await addApi(state.editFormData) : await editApi(state.editFormData as Required<IFrontendRoute>)
			await listPage()
			ElMessage.success(`${state.editDialog.add ? '新增' : '编辑'}成功`);
			state.editDialog.visible = false
		})
	},
}

const pagerState = reactive(new PagerState<IFrontendRoute>(1, 10))

onMounted(() => {
	listPage()
})
// 分页
const listPage = async () => {
	tableDataLoading.value = true
	let res = await listPageApi(pagerState.getPager(['name', 'title']))
	pagerState.items = res.data.items
	pagerState.total = res.data.total
	tableDataLoading.value = false
}
// 新增
const handleAdd = () => {
	state.editDialog.add = true
	state.editDialog.title = '新增路由'
	_.assign(state.editFormData, {
		name: '',
		title: '',
		route_path: '',
		component_path: ''
	});
	state.editDialog.visible = true
	nextTick(() => {
		editFormRef.value?.clearValidate()
	})
}

// 编辑
const handleEdit = async (row: Required<IFrontendRoute>) => {
	state.editDialog.add = false
	state.editDialog.title = '编辑路由'
	let res = await getApi(row.id)
	state.editFormData = _.pick(res.data, ['id', 'name', 'title', 'route_path', 'component_path'])
	state.editDialog.visible = true
}

// 删除
const handleDelete = async (row: Required<IFrontendRoute>) => {
	if (!await confirm("确认删除？")) {
		return
	}
	await deleteApi(row.id)
	await listPage()
}

</script>
<style scoped></style>
