<template>
  <PagerTable :loading="tableDataLoading" :pager-state="pagerState" :add-permiss-code="permiss.add" @list-page="listPage"
    @add-click="handleAdd" placeholder="名称/主机名/IP">
    <template #queryFormBtn>
      <el-button type="warning" icon="CaretRight" v-permiss="permiss.terminalTree" @click="toTerminalTree()">终端
      </el-button>
      <el-button type="info" icon="Document" @click="terminalSessionView">终端会话记录</el-button>
    </template>
    <el-table-column prop="id" label="ID" />
    <el-table-column prop="name" label="名称" />
    <el-table-column prop="host" label="主机名/IP" />
    <el-table-column label="操作"
      v-if="$hasPermission(permiss.edit) || $hasPermission(permiss.delete) || $hasPermission(permiss.terminalTree)">
      <template #default="scope">
        <el-button type="primary" icon="Edit" v-permiss="permiss.edit" @click="handleEdit(scope.row)">编辑</el-button>
        <el-button type="danger" icon="Delete" v-permiss="permiss.delete" @click="handleDelete(scope.row)">删除
        </el-button>
        <el-button type="warning" :loading="scope.row.loading" icon="CaretRight" v-permiss="permiss.terminalSingle"
          @click="toTerminalSingle(scope.row)">终端
        </el-button>
      </template>
    </el-table-column>
  </PagerTable>

  <ElDialog v-model="editDialog.visible" :title="editDialog.title" width="35%">
    <ElForm ref="editFormRef" :model="editFormData" status-icon :rules="editForm.rules" label-width="80px">
      <ElFormItem label="分组" prop="host_group_ids">
        <el-tree-select style="width: 100%;" v-model="editFormData.host_group_ids" :props="{ label: 'labelWithPath' }"
          highlight-current clearable :data="editDialog.hostGroupTreeSelectData" multiple check-strictly placeholder="请选择"
          :filterable="true" :render-after-expand="false">
          <template #default="{ node, data }">
            <span>{{ data.label }}</span>
          </template>
        </el-tree-select>
      </ElFormItem>
      <ElFormItem label="名称" prop="name">
        <el-input v-model="editFormData.name" type="text" />
      </ElFormItem>
      <ElFormItem label="连接地址" prop="ssh_uri">
        <el-input style="display: none;" v-model="editFormData.ssh_uri"></el-input>
        <div class="inputs">
          <ElFormItem prop="username">
            <el-input class="input1" v-model="editFormData.username" placeholder="用户名">
              <template #prepend>ssh</template>
            </el-input>
          </ElFormItem>
          <ElFormItem prop="host">
            <el-input class="input2" v-model="editFormData.host" placeholder="主机名/IP">
              <template #prepend>@</template>
            </el-input>
          </ElFormItem>
          <ElFormItem prop="port">
            <el-input class="input3" v-model.number="editFormData.port" placeholder="端口">
              <template #prepend>-p</template>
            </el-input>
          </ElFormItem>
        </div>
      </ElFormItem>
      <el-row :gutter="24">
        <el-col :span="18" v-if="!editFormData.use_key">
          <ElFormItem label="密码" prop="password">
            <el-input :disabled="!editDialog.showPassword" class="el-password"
              :type="editDialog.showPassword ? 'text' : 'password'" v-model="editFormData.password" placeholder="请输入密码">
              <template #append>
                <el-icon @click="changePasswordView">
                  <component :is="!editDialog.showPassword ? 'View' : 'Hide'"></component>
                </el-icon>
              </template>
            </el-input>
          </ElFormItem>
        </el-col>
        <el-col :span="6">
          <ElFormItem label="秘钥认证" prop="use_key">
            <el-switch v-model="editFormData.use_key" />
          </ElFormItem>
        </el-col>
      </el-row>

      <ElFormItem v-if="editFormData.use_key" label="私钥" prop="private_key">
        <el-input v-if="editDialog.showPassword" v-model="editFormData.private_key" :autosize="{ minRows: 4, maxRows: 8 }"
          type="textarea" placeholder="私钥，通常位于 $HOME/.ssh/id_rsa">
        </el-input>
        <el-input v-else :disabled="true" class="el-password" :type="editDialog.showPassword ? 'text' : 'password'"
          v-model="emptyText" placeholder="点击右侧小眼睛查看/编辑秘钥">
          <template #append>
            <el-icon @click="changePasswordView">
              <component :is="!editDialog.showPassword ? 'View' : 'Hide'"></component>
            </el-icon>
          </template>
        </el-input>
      </ElFormItem>

      <ElFormItem label="描述" prop="desc">
        <el-input v-model="editFormData.desc" :autosize="{ minRows: 2, maxRows: 4 }" type="textarea"
          placeholder="描述内容..." />
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

  <el-drawer destroy-on-close v-model="sessionDrawerVisible" title="终端会话记录" direction="rtl" size="50%">
    <TerminalSession />
  </el-drawer>
</template>

<script setup lang="ts">
import { nextTick, onMounted, reactive, ref, toRefs, watch } from 'vue'
import {
  addApi,
  deleteApi,
  editApi,
  getApi,
  getPasswordAndPrivateKeyApi,
  IHost,
  listPageApi,
  testSSHApi
} from '@/api/resource/host'
import { listTreeApi as listTreeHostGroups } from '@/api/resource/host-group'
import { genTreeSelectData } from '@/utils/treeSelect'
import { PagerState } from '@/utils/pager';
import _ from 'lodash'
import { ElMessage, FormInstance } from 'element-plus';
import { useHostStore } from '@/store/resource/host'
import { confirm } from '@/utils/generic'
import { useRouter } from 'vue-router';
import TerminalSession from './terminal-session.vue';

const enum permiss {
  add = 'host-add',
  edit = 'host-edit',
  delete = 'host-delete',
  terminalSingle = 'host-terminal-single',
  terminalTree = 'host-terminal-tree',
}

const state = reactive<
  IPageState<IHost &
  { ssh_uri?: string }> &
  {
    editDialog: { showPassword: boolean, hostGroupTreeSelectData: ITreeSelectData[] }
    sessionDrawerVisible: boolean
  }>({
    editFormData: {
      name: '',
      host: '',
      port: 22,
      username: '',
      password: '',
      private_key: '',
      use_key: false,
      desc: '',
      host_group_ids: [],
    },
    editDialog: {
      visible: false,
      title: '新增主机',
      add: true,
      showPassword: false,
      hostGroupTreeSelectData: [],
    },
    tableDataLoading: false,
    sessionDrawerVisible: false,
  })

const hostStore = useHostStore()
const editFormRef = ref<FormInstance>()
const { editFormData, editDialog, tableDataLoading, sessionDrawerVisible } = toRefs(state)
const pagerState = reactive(new PagerState<IHost>(1, 10))
const emptyText = ref('')
const router = useRouter()

const editForm = {
  rules: {
    host_group_ids: [
      { required: true, message: '请选择分组', trigger: 'blur' }
    ],
    name: [
      { required: true, message: '请输入名称', trigger: 'blur' }
    ],
    host: [
      { required: true, message: '请输入主机名或IP', trigger: 'blur' },
      {
        pattern: /^(?!:\/\/)(?=.{1,253}$)(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.){0,3}([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9]|[\d]{1,3})$/,
        message: '请输入合法的主机名或IP地址',
        trigger: 'blur'
      }
    ],
    ssh_uri: [
      { required: true, message: '请输入', trigger: 'blur' }
    ],
    private_key: [
      { required: true, message: '请输入私钥', trigger: 'blur' },
      {
        pattern: /^(-----BEGIN OPENSSH PRIVATE KEY-----)\n([A-Za-z0-9+/=\s]+)\n(-----END OPENSSH PRIVATE KEY-----)$/,
        message: '请输入合法的私钥文本',
        trigger: 'blur'
      }
    ],
    port: [
      { required: true, message: '请输入端口', trigger: 'blur' }
    ],
    username: [
      { required: true, message: '请输入用户名', trigger: 'blur' }
    ],
    password: [
      { required: true, message: '请输入密码', trigger: 'blur' }
    ],
  },
  submit: () => {
    editFormRef.value?.validate(async (valid: boolean) => {
      if (!valid) {
        return false
      }

      let editReqData = _.omit(state.editFormData, ['ssh_uri']) as IHost
      if (!state.editFormData.update_password_and_private_key) {
        editReqData = _.omit(editReqData, ['password', 'private_key'])
      }
      state.editDialog.add ? await addApi(state.editFormData) : await editApi(editReqData)
      await listPage()
      ElMessage.success(`${state.editDialog.add ? '新增' : '编辑'}成功`);
      hostStore.reloadGroup();
      state.editDialog.visible = false
    })
  },
}

watch(
  [() => state.editFormData.username, () => state.editFormData.host, () => state.editFormData.port],
  ([username, host, port]) => {
    state.editFormData.ssh_uri = `ssh://${username}@${host}:${port}`
  }
)

onMounted(() => {
  watch(() => hostStore.activeGroupID, () => {
    listPage()
  })
  listPage()
})

// 测试终端是否可以连接 & 打开终端标签页
const toTerminalSingle = async (row: Required<IHost & { loading: boolean }>) => {
  row.loading = true
  let res = await testSSHApi(row.id).catch(() => false)
  row.loading = false
  if (!res) {
    return
  }
  if (!await confirm("将在新标签页打开终端，是否确认？")) {
    return
  }
  const { href } = router.resolve({
    name: 'host-terminal-single',
    query: {
      id: row.id,
    }
  });
  window.open(href, "_blank");
}

const toTerminalTree = async () => {
  if (!await confirm("将在新标签页打开终端，是否确认？")) {
    return
  }
  const { href } = router.resolve({
    name: 'host-terminal-tree',
  });
  window.open(href, "_blank");
}

// 分页
const listPage = async () => {
  tableDataLoading.value = true
  let res = await listPageApi(pagerState.getPager(['name', 'host'], ['id', 'name', 'host']), hostStore.activeGroupID || 0)
  pagerState.items = res.data.items
  pagerState.total = res.data.total
  tableDataLoading.value = false
}

// 新增
const handleAdd = async () => {
  let res = await listTreeHostGroups()
  state.editDialog.hostGroupTreeSelectData = genTreeSelectData(res.data)
  state.editDialog.add = true
  state.editDialog.title = '新增主机'
  _.assign(state.editFormData, {
    id: null,
    name: '',
    host: '',
    port: 22,
    username: '',
    password: null,
    private_key: null,
    use_key: false,
    desc: '',
  });
  state.editFormData.host_group_ids = hostStore.activeGroupID ? [hostStore.activeGroupID] : []
  state.editDialog.visible = true
  nextTick(() => {
    editFormRef.value?.clearValidate()
  })
}

// 编辑
const handleEdit = async (row: Required<IHost>) => {
  let [groupTreeRes, hostRes] = await Promise.all([listTreeHostGroups(), getApi(row.id)])
  state.editDialog.hostGroupTreeSelectData = genTreeSelectData(groupTreeRes.data)
  state.editDialog.add = false
  state.editDialog.showPassword = false
  state.editFormData.update_password_and_private_key = false
  state.editDialog.title = '编辑主机'
  state.editFormData = _.pick(hostRes.data, ['id', 'name', 'host', 'port', 'username', 'password', 'use_key', 'desc'])
  state.editFormData.private_key = 'tmp'
  state.editFormData.ssh_uri = `ssh://${state.editFormData.username}@${state.editFormData.host}:${state.editFormData.port}`
  state.editFormData.host_group_ids = hostRes.data.host_groups.map(item => item.id ?? 0)
  state.editDialog.visible = true
  nextTick(() => {
    editFormRef.value?.clearValidate()
  })
}

// 删除
const handleDelete = async (row: Required<IHost>) => {
  if (!await confirm("确认删除？")) {
    return
  }
  await deleteApi(row.id)
  await listPage()
  hostStore.reloadGroup();
}

// 查看密码或秘钥
const changePasswordView = async () => {
  if (!state.editDialog.add && !state.editDialog.showPassword && !state.editFormData.update_password_and_private_key) {
    let res = await getPasswordAndPrivateKeyApi(state.editFormData.id ?? 0)
    state.editFormData.password = res.data.password
    state.editFormData.private_key = res.data.private_key
    state.editFormData.update_password_and_private_key = true
  }
  state.editDialog.showPassword = !state.editDialog.showPassword
}

const terminalSessionView = async () => {
  state.sessionDrawerVisible = true
}

</script>

<style scoped>
.inputs {
  display: flex;
  align-items: center;
}

:deep(.input1 .el-input__wrapper) {
  border-radius: 0;
  width: 60px;
}

:deep(.input2 .el-input__wrapper) {
  border-radius: 0;
  width: 127px;
}

:deep(.input2 .el-input-group__prepend) {
  border-bottom-left-radius: 0;
  border-top-left-radius: 0;
}

:deep(.input3 .el-input-group__prepend) {
  border-bottom-left-radius: 0;
  border-top-left-radius: 0;
}

:deep(.el-password .el-input-group__prepend) {
  border-bottom-left-radius: 0;
  border-top-left-radius: 0;
}

:deep(.el-password .el-input-group__append) {
  cursor: pointer;
}

.input1 {
  position: relative;
  left: 1px;
}

.input3 {
  position: relative;
  left: -1px;
}
</style>
