<template>
	<el-row style="height: 100%;" :gutter="12">
		<el-col :span="3">
			<el-card>
				<div class="logo">DevOps Helper</div>
				<el-tree ref="treeSelectRef" @node-click="handleNodeClick" default-expand-all :data="treeState.data"
					node-key="id" :render-after-expand="false" :expand-on-click-node="false" :highlight-current="true">
					<template #default="{ node, data }">
						<span style="position: relative;top:2px">
							<el-icon v-if="data.isHost">
								<Monitor />
							</el-icon>
							<el-icon v-else-if="node.expanded && !node.isLeaf">
								<FolderOpened />
							</el-icon>
							<el-icon v-else>
								<Folder />
							</el-icon>
						</span>
						<span style="margin-left: 5px;">{{ data.label }}</span><span style="margin-left: 5px;"
							v-if="!data.isHost">({{
								data.hostCount
							}})</span>
					</template>
				</el-tree>
			</el-card>
		</el-col>
		<el-col :span="21">
			<el-card>
				<el-tabs @tab-click="handleTabClick" @tab-remove="handleTabRemove" v-model="tabState.activeTab" closable
					type="border-card">
					<!-- 占位标签页模板 -->
					<el-tab-pane label="DevOps Helper" :name="-1" v-if="!tabState.data.length">
						<div class="placeholder-content">
							<p>
								<el-icon style="position: relative;top:14px">
									<InfoFilled />
								</el-icon>
								当前没有打开任何终端连接
							</p>
						</div>
					</el-tab-pane>
					<el-tab-pane :lazy="true" :key="item.key" :name="idx" :label="item.label"
						v-for="item, idx in tabState.data">
						<template #label>
							<div :class="{ 'connected': item.connected, 'disconnected': !item.connected }">
								{{ item.label }}
								<span style="position: relative; top: 2px">
									<el-icon v-if="item.connected">
										<SuccessFilled />
									</el-icon>
									<el-icon v-else>
										<CircleCloseFilled />
									</el-icon>
								</span>
							</div>
						</template>
						<SshTerminal @reload="item.key = Date.now()" @connected="item.connected = true"
							@close="item.connected = false" :ws-url="item.terminalUrl" :in-body="false"
							:padding-bottom="125" @ctrl-u="handleTerminalCtrlU">
						</SshTerminal>
						<el-drawer v-if="$hasPermission(permiss.fileManagerReadonly)" v-model="fileManagerState.visible" title="文件管理器" direction="rtl"
							:before-close="handleFileManagerDrawerClose" size="50%">
							<FileManager @download-file="handleDownloadFile" :visible="fileManagerState.visible"
								:ws-url="fileManagerState.wsUrl" @close="handleFileManagerClose" />
						</el-drawer>
					</el-tab-pane>
				</el-tabs>
			</el-card>
		</el-col>
	</el-row>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, reactive, ref } from 'vue';
import { TabPaneName, TabsPaneContext, TreeNode } from 'element-plus';
import { getWsProtocol, getHttpProtocol } from '@/utils/generic'
import { getLoginInfo } from '@/utils/login'
import { ITreeItem } from './group.vue'
import { IHostGroup } from '@/api/resource/host-group';
import { listTreeWithHostsApi } from '@/api/resource/host-group';
import FileManager from './file-manager.vue'
import { confirm } from '@/utils/generic'

const enum permiss {
	fileManagerReadonly = 'host-terminal-file-manager-readonly',
}

interface TabItem {
	name: number
	terminalUrl: string
	label: string
	connected: boolean
	key: number
}

const treeSelectRef = ref<any>(null)

// 树状态
const treeState = reactive<{
	data: Array<ITreeItem>
}>({
	data: [],
})


// tab 状态
const tabState = reactive<{ data: TabItem[], activeTab: number }>({
	data: [],
	activeTab: 0,
})

const fileManagerState = reactive<{
	visible: boolean;
	showHiddenFiles: boolean;
	wsUrl: string;
	path: string;
}>({
	visible: false,
	showHiddenFiles: false,
	wsUrl: "",
	path: "/",
})

const handleFileManagerClose = () => {
	fileManagerState.visible = false
}

const handleTerminalCtrlU = () => {
	const wsUrl = `${getWsProtocol()}://${import.meta.env.VITE_BASE_URL}/Host/SFTPFileManage?id=${tabState.data[tabState.activeTab].name}&token=${getLoginInfo()?.token}`
	fileManagerState.wsUrl = wsUrl
	fileManagerState.visible = true
}

const handleDownloadFile = async (path: string) => {
	window.open(`${getHttpProtocol()}://${import.meta.env.VITE_BASE_URL}/Host/DownloadFile?id=${tabState.data[tabState.activeTab].name}&path=${path}&token=${getLoginInfo()?.token}`)
	// 可以优化，服务端生成临时认证链接，前端直接下载
	// let loading = openLoading("正在下载文件，请稍后...")
	// let res = await downloadFileApi(tabState.data[tabState.activeTab].name, path).catch(() => false)
	// loading.close()
	// if (res) saveAs(res as Blob, basename(path))
}

const convertTree = (tree: IHostGroup[]): ITreeItem[] => {
	const result: ITreeItem[] = []
	tree.forEach(item => {
		const newItem: ITreeItem = {
			id: item.id ?? 0,
			hostCount: item.hosts?.length ?? 0,
			permissionCount: item.host_group_permissions?.length ?? 0,
			label: item.name,
			children: item.children ? convertTree(item.children) : [],
			isHost: false
		}

		if (item.hosts) {
			item.hosts.forEach(host => {
				newItem.children?.push({
					id: host.id ?? 0,
					hostCount: 0,
					permissionCount: 0,
					label: host.name,
					isHost: true,
				})
			})
		}

		result.push(newItem)
	})
	return result
}

const handleTabRemove = (idx: TabPaneName) => {
	tabState.data = tabState.data.filter((_, index) => index !== idx)
	if (tabState.activeTab >= (idx as number)) {
		tabState.activeTab--
	}
	if (tabState.activeTab < 0) {
		tabState.activeTab = 0
	}
}

const handleTabClick = (idx: TabsPaneContext) => {
	if ((idx.paneName as number) < 0) {
		return
	}
	let currentItem = tabState.data[idx.paneName as number]
	treeSelectRef.value.setCurrentKey(currentItem.name)
	// currentItem.key = Date.now()
}

const handleNodeClick = (data: ITreeItem, nodeAttr: any, treeNode: TreeNode) => {
	if (data.isHost) {
		const wsUrl = `${getWsProtocol()}://${import.meta.env.VITE_BASE_URL}/Host/Terminal?id=${data.id}&token=${getLoginInfo()?.token}`
		tabState.data.push({
			name: data.id,
			terminalUrl: wsUrl,
			label: data.label,
			connected: false,
			key: Date.now(),
		})
		tabState.activeTab = tabState.data.length - 1
	}
}

// 处理文件管理器关闭
const handleFileManagerDrawerClose = async () => {
	if (!await confirm("确认关闭？")) {
		return
	}
	fileManagerState.visible = false
}

onMounted(async () => {
	document.title = `终端管理`
	let res = await listTreeWithHostsApi()
	treeState.data = convertTree(res.data)
	tabState.activeTab = -1
})

onUnmounted(() => {
	fileManagerState.visible = false
})
</script>

<style scoped>
.logo {
	font-size: 24px;
	font-weight: bold;
	color: #409EFF;
	margin-bottom: 20px;
	text-align: center;
}

.el-tree :deep(.el-tree-node__content) {
	padding: 2px 10px;
	transition: background-color 0.3s;
}

.el-tree :deep(.el-tree-node__content:hover) {
	background-color: #E6F7FF;
}

:deep(.el-card) {
	height: 100%;
}

:deep(.el-card__body) {
	height: 100%;
}

:deep(.el-tabs) {
	height: 100%;
}

:deep(.el-tabs__content) {
	height: 100%;
}

.placeholder-content {
	color: #999;
	display: flex;
	flex-direction: column;
	height: 100%;
	justify-content: center;
	text-align: center;
}

.placeholder-content i {
	font-size: 40px;
	margin-bottom: 20px;
}

.connected {
	color: #67C23A;
}

.disconnected {
	color: #F56C6C;
}
</style>
