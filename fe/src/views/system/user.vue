<template>
	<div>
		<ElCard>
			<PagerTable :loading="tableDataLoading" :pager-state="pagerState" @list-page="listPage" @add-click="handleAdd"
				:add-permiss-code="permiss.add" placeholder="用户名/真实姓名">
				<el-table-column prop="id" label="ID" />
				<el-table-column prop="username" label="用户名" />
				<el-table-column prop="real_name" label="真实姓名" />
				<el-table-column prop="phone" label="电话" />
				<el-table-column prop="email" label="邮箱" />
				<el-table-column min-width="120" label="操作"
					v-if="$hasPermission(permiss.edit) || $hasPermission(permiss.delete) || $hasPermission(permiss.resetPwd)">
					<template #default="scope">
						<el-button icon="Edit" type="primary" v-permiss="permiss.edit"
							@click="handleEdit(scope.row)">编辑</el-button>
						<el-button icon="Delete" type="danger" v-permiss="permiss.delete"
							@click="handleDelete(scope.row)">删除</el-button>
						<el-button icon="RefreshLeft" type="warning" v-permiss="permiss.resetPwd"
							@click="handleResetPwd(scope.row)">重置密码</el-button>
					</template>
				</el-table-column>
			</PagerTable>
		</ElCard>

		<ElDialog v-model="editDialog.visible" :title="editDialog.title" width="20%">
			<ElForm ref="editFormRef" :model="editFormData" status-icon :rules="editForm.rules" label-width="80px">
				<ElFormItem label="用户名" prop="username">
					<el-input v-model="editFormData.username" type="text" />
				</ElFormItem>
				<ElFormItem label="真实姓名" prop="real_name">
					<el-input v-model="editFormData.real_name" type="text" />
				</ElFormItem>
				<ElFormItem label="电话" prop="phone">
					<el-input v-model="editFormData.phone" />
				</ElFormItem>
				<ElFormItem label="邮箱" prop="email">
					<el-input v-model="editFormData.email" />
				</ElFormItem>
				<ElFormItem label="角色" prop="role_ids">
					<el-select v-model="editFormData.role_ids" style="width: 100%;" multiple placeholder="请选择">
						<el-option v-for="item in editDialog.roleSelectData" :key="item.id" :label="item.name"
							:value="item.id ?? 0" />
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

		<ElDialog v-model="resetPwdDialog.visible" :title="resetPwdDialog.title" width="20%">
			<ElForm ref="resetPwdFormRef" :model="resetPwdFormData" status-icon :rules="resetPwdForm.rules"
				label-width="100px">
				<ElFormItem label="新密码" prop="new_password1">
					<el-input type="password" v-model.trim="resetPwdFormData.new_password1" />
				</ElFormItem>
				<ElFormItem label="确认新密码" prop="new_password2">
					<el-input type="password" v-model.trim="resetPwdFormData.new_password2" />
				</ElFormItem>
			</ElForm>
			<template #footer>
				<span class="dialog-footer">
					<ElButton @click="resetPwdDialog.visible = false">取消</ElButton>
					<ElButton type="primary" @click="resetPwdForm.submit">
						确认
					</ElButton>
				</span>
			</template>
		</ElDialog>
	</div>
</template>

<script setup lang="ts" name="system-user">
import { ref, reactive, onMounted, toRefs, nextTick } from 'vue'
import { addApi, editApi, deleteApi, getApi, listPageApi, resetPwd, IUser } from '@/api/system/user'
import { IRole, listAllApi as listAllRoles } from '@/api/system/role'
import { PagerState } from '@/utils/pager'
import { confirm } from '@/utils/generic'
import _ from 'lodash'
import { ElMessage, FormInstance } from 'element-plus'

const enum permiss {
	add = 'user-add',
	edit = 'user-edit',
	delete = 'user-delete',
	resetPwd = 'user-reset-pwd',
}

const pagerState = reactive(new PagerState<IUser>(1, 10))

const editFormRef = ref<FormInstance>()
const resetPwdFormRef = ref<FormInstance>()

const state = reactive<IPageState<IUser> & {
	editFormData: { role_ids?: number[] },
	editDialog: { roleSelectData: IRole[] },
	resetPwdFormData: { id: number, new_password1: string, new_password2: string },
	resetPwdDialog: { visible: boolean, title: string },
	tableDataLoading: boolean,
}>({
	editFormData: {
		id: 0,
		username: '',
		real_name: '',
		phone: '',
		email: '',
		role_ids: [],
	},
	editDialog: {
		visible: false,
		title: '新增用户',
		add: true,
		roleSelectData: [],
	},
	resetPwdFormData: {
		id: 0,
		new_password1: '',
		new_password2: ''
	},
	resetPwdDialog: {
		visible: false,
		title: '重置密码',
	},
	tableDataLoading: false,
})

const { editFormData, editDialog, resetPwdFormData, resetPwdDialog, tableDataLoading } = toRefs(state)

const editForm = {
	rules: {
		username: [
			{ required: true, message: '请输入用户名', trigger: 'blur' }
		],
		real_name: [
			{ required: true, message: '请输入真实姓名', trigger: 'blur' },
			{ pattern: /^[\u4e00-\u9fa5]{0,}$/, message: '请输入汉字', trigger: 'blur' }
		],
		phone: [
			{ pattern: /^1[3456789]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
		],
		email: [
			{
				pattern: /^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+/,
				message: '请输入正确的邮箱', trigger: 'blur'
			}
		]
	},
	submit: () => {
		editFormRef.value?.validate(async (valid: boolean) => {
			if (!valid) {
				return false
			}
			let res = state.editDialog.add ? await addApi(state.editFormData) : await editApi(state.editFormData as Required<IUser>)
			await listPage()
			ElMessage.success(`${state.editDialog.add ? res.msg : '编辑成功'}`);
			state.editDialog.visible = false
		})
	},
}


const resetPwdForm = {
	rules: {
		new_password1: [
			{ required: true, message: '请输入新密码', trigger: 'blur' },
			{ min: 6, max: 20, message: '新密码长度应为 6 到 24 个字符', trigger: 'blur' },
			{ pattern: /^[a-zA-Z0-9_]+$/, message: '新密码只能包含字母、数字和下划线', trigger: 'blur' },
		],
		new_password2: [
			{ required: true, message: '请再次输入新密码', trigger: 'blur' },
			{ min: 6, max: 20, message: '新密码长度应为 6 到 24 个字符', trigger: 'blur' },
			{ pattern: /^[a-zA-Z0-9_]+$/, message: '新密码只能包含字母、数字和下划线', trigger: 'blur' },
			{ validator: validateNewPassword, trigger: 'blur' }
		]
	},
	submit: () => {
		resetPwdFormRef.value?.validate(async (valid: boolean) => {
			if (!valid) {
				return false
			}

			let reqData = {
				id: state.resetPwdFormData.id,
				new_password: state.resetPwdFormData.new_password2
			}

			let res = await resetPwd(reqData)
			ElMessage.success(res.msg);
			state.resetPwdDialog.visible = false
		})
	}
}

function validateNewPassword(rule: any, value: string, callback: Function) {
	if (value !== state.resetPwdFormData.new_password1) {
		callback(new Error('两次输入的密码不相同'));
	} else {
		callback();
	}
}


// 新增
const handleAdd = async () => {
	state.editDialog.add = true
	state.editDialog.title = '新增用户'
	_.assign(state.editFormData, {
		id: 0,
		username: '',
		real_name: '',
		phone: '',
		email: '',
		role_ids: [],
	});

	let res = await listAllRoles()
	state.editDialog.roleSelectData = res.data
	state.editDialog.visible = true
	nextTick(() => {
		editFormRef.value?.clearValidate()
	})
}

// 编辑
const handleEdit = async (row: Required<IUser>) => {
	state.editDialog.add = false
	state.editDialog.title = '编辑用户'

	let [userRes, allRolesRes] = await Promise.all([getApi(row.id), listAllRoles()])
	state.editFormData = _.pick(userRes.data, ['id', 'username', 'phone', 'email', 'real_name'])
	state.editFormData.role_ids = userRes.data.roles.map((item: IRole) => item.id ?? 0)
	state.editDialog.roleSelectData = allRolesRes.data
	state.editDialog.visible = true
}

// 删除
const handleDelete = async (row: Required<IUser>) => {
	if (!await confirm("确认删除？")) {
		return
	}
	await deleteApi(row.id)
	await listPage()
}

// 重置密码
const handleResetPwd = async (row: Required<IUser>) => {
	_.assign(state.resetPwdFormData, {
		id: row.id,
		new_password1: '',
		new_password2: '',
	});
	state.resetPwdDialog.visible = true
	nextTick(() => {
		resetPwdFormRef.value?.clearValidate()
	})
}

onMounted(() => {
	listPage()
})
// 分页
const listPage = async () => {
	state.tableDataLoading = true
	let res = await listPageApi(pagerState.getPager(['username', 'real_name']))
	pagerState.items = res.data.items
	pagerState.total = res.data.total
	state.tableDataLoading = false
}

</script>

<style scoped></style>
