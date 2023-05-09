<template>
    <PagerTable :loading="tableDataLoading" :pager-state="pagerState" @list-page="listPage" :add-permiss-code="permiss.add"
        @add-click="handleAdd" placeholder="用户/角色名称">
        <el-table-column prop="id" label="ID" />
        <el-table-column prop="type" label="类型">
            <template #default="scope">
                <el-tag v-if="scope.row.type === 1">用户</el-tag>
                <el-tag type="info" v-else-if="scope.row.type === 2">角色</el-tag>
            </template>
        </el-table-column>
        <el-table-column prop="ref_name" label="名称" />
        <el-table-column label="关联分组" :min-width="120" v-if="!hostStore.activeGroupID" :show-overflow-tooltip="true">
            <template #default="{ row }">
                <div class="tag-container">
                    <span class="tag-click" :key="item.split('_')[0]"
                        @click="handleGroupTagClick(parseInt(item.split('_')[0]))"
                        v-for="item in (row as IHostGroupPermissionListItem).host_group_names_str.split(',')">
                        {{ item.split("_")[1] }}
                    </span>
                </div>
            </template>
        </el-table-column>

        <el-table-column label="操作" v-if="$hasPermission(permiss.edit) || $hasPermission(permiss.delete)">
            <template #default="scope">
                <el-button type="primary" icon="Edit" v-permiss="permiss.edit" @click="handleEdit(scope.row)">编辑</el-button>
                <el-button type="danger" icon="Delete" v-permiss="permiss.delete"
                    @click="handleDelete(scope.row)">删除</el-button>
            </template>
        </el-table-column>
    </PagerTable>

    <ElDialog v-model="editDialog.visible" :title="editDialog.title" width="30%">
        <ElForm ref="editFormRef" :model="editFormData" status-icon :rules="editForm.rules" label-width="80px">
            <ElFormItem label="类型" prop="type">
                <el-radio-group v-model="editFormData.type">
                    <el-radio :label="1">用户</el-radio>
                    <el-radio :label="2">角色</el-radio>
                </el-radio-group>
            </ElFormItem>
            <ElFormItem v-if="editFormData.type === 1" label="选择用户" prop="user_id">
                <el-select v-model="editFormData.user_id" placeholder="选择用户" style="width: 100%;">
                    <el-option v-for="item in editDialog.allUsers" :key="item.id" :label="item.real_name"
                        :value="item.id ?? 0" />
                </el-select>
            </ElFormItem>
            <ElFormItem v-if="editFormData.type === 2" label="选择角色" prop="role_id">
                <el-select v-model="editFormData.role_id" placeholder="选择角色" style="width: 100%;">
                    <el-option v-for="item in editDialog.allRoles" :key="item.id" :label="item.name"
                        :value="item.id ?? 0" />
                </el-select>
            </ElFormItem>
            <ElFormItem label="分组" prop="host_group_ids">
                <el-tree-select style="width: 100%;" v-model="editFormData.host_group_ids"
                    :props="{ label: 'labelWithPath' }" highlight-current clearable
                    :data="editDialog.hostGroupTreeSelectData" multiple check-strictly placeholder="请选择" :filterable="true"
                    :render-after-expand="false">
                    <template #default="{ node, data }">
                        <span>{{ data.label }}</span>
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
</template>

<script setup lang="ts">
import { reactive, toRefs, ref, nextTick, onMounted, watch } from 'vue'
import { IHostGroupPermission, IHostGroupPermissionListItem, addApi, editApi, getApi, deleteApi, listPageApi } from '@/api/resource/host-group-permission'
import { PagerState } from '@/utils/pager'
import { ElMessage, FormInstance } from 'element-plus'
import { listAllApi as listAllRoles, IRole } from '@/api/system/role'
import { listAllApi as listAllUsers, IUser } from '@/api/system/user'
import { listTreeApi as listTreeHostGroups } from '@/api/resource/host-group'
import { genTreeSelectData } from '@/utils/treeSelect'
import { useHostStore } from '@/store/resource/host'
import _ from 'lodash'


const enum permiss {
    add = 'host-group-permission-add',
    edit = 'host-group-permission-edit',
    delete = 'host-group-permission-delete',
}

const state = reactive<IPageState<Partial<IHostGroupPermission> & { host_group_ids?: number[], user_id?: number, role_id?: number }> & { editDialog: { allUsers: IUser[], allRoles: IRole[], hostGroupTreeSelectData: ITreeSelectData[] } }>({
    editFormData: {
        id: undefined,
        type: 1,
        ref_id: undefined,
        host_groups: [],
        host_group_ids: [],
        user_id: undefined,
        role_id: undefined,
    },
    editDialog: {
        visible: false,
        title: '新增权限',
        add: true,
        allUsers: [],
        allRoles: [],
        hostGroupTreeSelectData: []
    },
    tableDataLoading: false,
})

const editFormRef = ref<FormInstance>()
const editForm = {
    rules: {
        type: [
            { required: true, message: '请选择权限类型', trigger: 'blur' }
        ],
        user_id: [
            { required: true, message: "请选择用户", trigger: 'blur' }
        ],
        role_id: [
            { required: true, message: "请选择角色", trigger: 'blur' }
        ],
        host_group_ids: [
            { required: true, message: '请选择主机分组', trigger: 'blur' }
        ],
    },
    submit: () => {
        editFormRef.value?.validate(async (valid: boolean) => {
            if (!valid) {
                return false
            }
            let reqData = _.pick(state.editFormData, ['id', 'type', 'host_group_ids', 'ref_id'])
            reqData.ref_id = state.editFormData.type === 1 ? state.editFormData.user_id : state.editFormData.role_id
            if (state.editDialog.add) {
                reqData = _.omit(reqData, ['id'])
            }
            state.editDialog.add ? await addApi(reqData as IHostGroupPermission) : await editApi(reqData)
            await listPage()
            ElMessage.success(`${state.editDialog.add ? '新增' : '编辑'}成功`);
            hostStore.reloadGroup()
            state.editDialog.visible = false
        })
    },
}


const { editFormData, editDialog, tableDataLoading } = toRefs(state)
const pagerState = reactive(new PagerState<IHostGroupPermissionListItem>(1, 10))
const hostStore = useHostStore()

const listPage = async () => {
    tableDataLoading.value = true
    let res = await listPageApi(pagerState.getPager([], []), hostStore.activeGroupID || 0)
    pagerState.items = res.data.items
    pagerState.total = res.data.total
    tableDataLoading.value = false
}

const handleGroupTagClick = (groupID: number) => {
    hostStore.activeGroup(groupID)
}

const handleAdd = async () => {
    state.editDialog.add = true
    state.editDialog.title = '新增权限'
    let [allUsersRes, allRolesRes, hostGroupTreeDataRes] = await Promise.all([listAllUsers(), listAllRoles(), listTreeHostGroups()])
    state.editDialog.allUsers = allUsersRes.data.filter(item => item.username !== 'admin')
    state.editDialog.allRoles = allRolesRes.data.filter(item => item.code !== 'admin')
    state.editDialog.hostGroupTreeSelectData = genTreeSelectData(hostGroupTreeDataRes.data)
    _.assign(state.editFormData, {
        id: undefined,
        type: 1,
        ref_id: undefined,
        host_group_ids: [],
        host_groups: [],
    });
    state.editFormData.host_group_ids = hostStore.activeGroupID ? [hostStore.activeGroupID] : []
    state.editDialog.visible = true
    nextTick(() => {
        editFormRef.value?.clearValidate()
    })
}

const handleEdit = async (row: Required<IHostGroupPermission>) => {
    state.editDialog.add = false
    state.editDialog.title = '编辑权限'
    let [allUsersRes, allRolesRes, hostGroupTreeDataRes, hostGroupPermissionRes] = await Promise.all([listAllUsers(), listAllRoles(), listTreeHostGroups(), getApi(row.id)])
    state.editDialog.allUsers = allUsersRes.data.filter(item => item.username !== 'admin')
    state.editDialog.allRoles = allRolesRes.data.filter(item => item.code !== 'admin')
    state.editDialog.hostGroupTreeSelectData = genTreeSelectData(hostGroupTreeDataRes.data)

    state.editFormData = _.pick(hostGroupPermissionRes.data, ['id', 'type', 'ref_id'])
    if (state.editFormData.type === 1) {
        state.editFormData.user_id = state.editFormData.ref_id
    } else if (state.editFormData.type === 2) {
        state.editFormData.role_id = state.editFormData.ref_id
    }

    state.editFormData.host_group_ids = hostGroupPermissionRes.data.host_groups.map(item => item.id ?? 0)
    state.editDialog.visible = true
}
const handleDelete = async (row: Required<IHostGroupPermission>) => {
    if (!await confirm("确认删除？")) {
        return
    }
    await deleteApi(row.id)
    await listPage()
    hostStore.reloadGroup();
}

onMounted(async () => {
    watch(() => hostStore.activeGroupID, async () => {
        await listPage()
    })
    await listPage()
})

</script>

<style scoped></style>
