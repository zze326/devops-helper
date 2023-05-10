<template>
    <el-tree ref="treeSelectRef" @node-click="handleNodeClick" @node-contextmenu="handleContextMenu" default-expand-all
        :data="treeState.data" node-key="id" :render-after-expand="false" :expand-on-click-node="false"
        :highlight-current="true">
        <template #default="{ node, data }">
            <el-input ref="labelInputRef" @blur="handleNewTreeItem" @keyup.enter.prevent="handleLabelInputEnter"
                v-if="data.editable" v-model.trim="data.inputValue" size="small" placeholder="请输入"></el-input>
            <div v-else>
                <span style="position: relative;top:2px">
                    <el-icon v-if="node.expanded && !node.isLeaf">
                        <FolderOpened />
                    </el-icon>
                    <el-icon v-else>
                        <Folder />
                    </el-icon>
                </span>
                <span style="margin-left: 5px;">{{ data.label }} ({{ hostStore.permissionManage ? data.permissionCount :
                    data.hostCount }})</span>
            </div>
        </template>
    </el-tree>
    <el-row v-if="!treeState.data || treeState.data.length === 0">
        <el-col :span="24" style="text-align: center;">
            <el-button size="small" icon="Plus" @click="handleContextMenuItemClick('add-root')">新建</el-button>
        </el-col>
    </el-row>

    <ContextMenu @item-click="handleContextMenuItemClick" :items="contextMenuState.items"
        :visible="contextMenuState.visible" :x="contextMenuState.x" :y="contextMenuState.y"
        @close="contextMenuState.visible = false" />
</template>

<script lang="ts" setup>
import { reactive, getCurrentInstance, nextTick, onMounted, ref, watch } from 'vue'
import { ElMessage, TreeNode, ElTree } from 'element-plus';
import { addApi, listTreeApi, renameApi, deleteApi, IHostGroup, ITreeItem } from '@/api/resource/host-group';
import { confirm, hasPermission } from '@/utils/generic'
import { useHostStore } from '@/store/resource/host'

const enum permiss {
    add = 'host-group-add',
    rename = 'host-group-rename',
    delete = 'host-group-delete',
}

const appInstance = getCurrentInstance()
const treeSelectRef = ref<any>(null)

// IHostGroup[] => ITreeItem[]
const convertTree = (tree: IHostGroup[]): ITreeItem[] => {
    const result: ITreeItem[] = []
    tree.forEach(item => {
        const newItem: ITreeItem = {
            id: item.id ?? 0,
            hostCount: item.hosts?.length ?? 0,
            permissionCount: item.host_group_permissions?.length ?? 0,
            label: item.name,
            children: item.children ? convertTree(item.children) : [],
        }
        result.push(newItem)
    })
    return result
}

// 递归移除 data 中的元素
function removeItems(data: ITreeItem[], callback: (item: ITreeItem) => boolean): void {
    data.forEach((item, index) => {
        // 如果当前元素符合需要移除的条件，则将其从原数组中删除
        if (callback(item)) {
            data.splice(index, 1);
            return;
        }

        // 如果当前元素有子元素，则递归调用该函数来处理子元素
        if (item.children && item.children.length > 0) {
            removeItems(item.children, callback);
        }
    });
}

// 右键菜单状态
const contextMenuState = reactive({
    visible: false,
    x: 0,
    y: 0,
    items: [
        {
            label: '新建根分组',
            icon: 'Folder',
            key: 'add-root',
            permiss: permiss.add
        },
        {
            label: '新建子分组',
            icon: 'FolderAdd',
            key: 'add-sub',
            permiss: permiss.add
        },
        {
            label: '重命名',
            icon: 'Edit',
            key: 'rename',
            permiss: permiss.rename
        },
        {
            label: '删除此分组',
            icon: 'FolderRemove',
            key: 'delete',
            permiss: permiss.delete
        }
    ]
})

// 树状态
const treeState = reactive<{
    contextMenuSelect: {
        data: ITreeItem | null,
        node: TreeNode | null,
        opCommand: string | null,
        newData: ITreeItem | null,
    },
    lastClickKey: number | null,
    data: Array<ITreeItem>
}>({
    contextMenuSelect: {
        data: null,
        node: null,
        opCommand: null,
        newData: null,
    },
    data: [],
    lastClickKey: null,
})

// 右键菜单点击事件
const handleContextMenu = (e: PointerEvent, data: ITreeItem, node: TreeNode) => {
    e.preventDefault();
    if (!contextMenuState.items.some(item => hasPermission(item.permiss))) {
        return
    }
    contextMenuState.x = e.clientX - 10
    contextMenuState.y = e.clientY - 10
    contextMenuState.visible = true;
    node.expanded = true
    treeState.contextMenuSelect.data = data
    treeState.contextMenuSelect.node = node
}

const hostStore = useHostStore()
const handleNodeClick = (data: ITreeItem, nodeAttr: any, treeNode: TreeNode) => {
    if (treeState.lastClickKey === treeSelectRef.value.getCurrentKey()) {
        treeSelectRef.value.setCurrentKey(null)
        treeState.lastClickKey = null
        hostStore.activeGroup(null)
    } else {
        treeState.lastClickKey = treeSelectRef.value.getCurrentKey()
        hostStore.activeGroup(data.id)
    }
}


// 新增输入框
const newGroup = (treeItems: Array<ITreeItem>, newItem: ITreeItem) => {
    treeItems.push(newItem)
    lableInputFocus()
}

// 输入框获取焦点
const lableInputFocus = () => {
    nextTick(() => {
        // 新增的输入框获取焦点
        const inputRef = appInstance?.refs['labelInputRef'] as HTMLInputElement;
        if (inputRef) {
            inputRef.focus();
        }
    })
}

// 处理右键菜单点击事件
const handleContextMenuItemClick = async (command: string) => {
    treeState.contextMenuSelect.opCommand = command
    let newItem = {
        id: 0,
        label: '新建分组',
        editable: true,
        hostCount: 0,
        permissionCount: 0,
    }
    switch (command) {
        case 'add-root':
            treeState.contextMenuSelect.newData = newItem
            newGroup(treeState.data, newItem)
            break;
        case 'add-sub':
            treeState.contextMenuSelect.newData = newItem
            treeState.contextMenuSelect.data!.children = treeState.contextMenuSelect.data!.children || []
            newGroup(treeState.contextMenuSelect.data!.children, newItem)
            break;
        case 'rename':
            treeState.contextMenuSelect.data!.editable = true
            treeState.contextMenuSelect.data!.inputValue = treeState.contextMenuSelect.data!.label
            lableInputFocus()
            break;
        case 'delete':
            if (treeState.contextMenuSelect.data!.children && treeState.contextMenuSelect.data!.children!.length > 0) {
                ElMessage.warning('请移除子分组后再尝试删除')
                return
            }
            if (!await confirm(`确定删除分组 ${treeState.contextMenuSelect.data!.label} 吗？`)) return
            await deleteApi(treeState.contextMenuSelect.data!.id)
            // 如果删除分组下有主机，提示请移除分组下的主机后再尝试删除
            removeItems(treeState.data, item => item.id === treeState.contextMenuSelect.data!.id)
            ElMessage.success('删除成功')
            break
    }
    contextMenuState.visible = false;
}

const handleLabelInputEnter = async (e: any) => {
    treeState.contextMenuSelect.data!.editable = false
    treeState.contextMenuSelect.newData!.editable = false
}

// 处理输入框回车、失去焦点事件
const handleNewTreeItem = async (e: any) => {
    let parentID = 0
    switch (treeState.contextMenuSelect.opCommand) {
        case 'rename':
            if (!e.target.value) {
                treeState.contextMenuSelect.data!.editable = false
                return
            }
            if (treeState.contextMenuSelect.data!.label !== e.target.value) {
                let res = await renameApi({
                    id: treeState.contextMenuSelect.data!.id,
                    name: e.target.value,
                }).catch(() => false)
                if (!res) {
                    return
                }
                treeState.contextMenuSelect.data!.label = e.target.value
                ElMessage.success('重命名成功')
            }
            treeState.contextMenuSelect.data!.editable = false
            break;
        case 'add-sub':
            parentID = treeState.contextMenuSelect.data!.id
        case 'add-root':
            if (!e.target.value) {
                removeItems(treeState.data, item => item.editable === true)
                return
            }
            let reqData: IHostGroup = {
                name: e.target.value,
                parent_id: parentID,
            }

            let res = await addApi(reqData).catch(() => false)
            if (!res) {
                removeItems(treeState.data, item => item.editable === true || item.id === 0)
                return
            }
            treeState.contextMenuSelect.newData!.id = (res as IResp<IHostGroup>).data.id || 0
            treeState.contextMenuSelect.newData!.label = e.target.value
            treeState.contextMenuSelect.newData!.editable = false
            ElMessage.success('新增成功')
            break
    }
}

const reloadTree = async () => {
    let res = await listTreeApi()
    treeState.data = convertTree(res.data)
}

onMounted(async () => {
    await reloadTree()
    watch(() => hostStore.reloadGroupCounter, async () => {
        await reloadTree()
        treeSelectRef.value.setCurrentKey(hostStore.activeGroupID)
    })

    watch(() => hostStore.activeGroupID, () => {
        treeState.lastClickKey = hostStore.activeGroupID ?? 0
        treeSelectRef.value.setCurrentKey(hostStore.activeGroupID)
    })
})
</script>

<style scoped></style>
