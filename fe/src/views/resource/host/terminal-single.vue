<template>
	<SshTerminal :ws-url="wsUrl" @reload="handleTerminalReload" :in-body="true" @ctrl-u="handleTerminalCtrlU"></SshTerminal>
	<FileManager @download-file="handleDownloadFile" :visible="fileManagerState.visible" :ws-url="fileManagerState.wsUrl"
		@close="handleFileManagerClose" />
</template>

<script setup lang="ts">
import { onMounted, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { getWsProtocol, getHttpProtocol } from '@/utils/generic'
import { getLoginInfo } from '@/utils/login'
import { getApi } from '@/api/resource/host'
import FileManager from './file-manager.vue'

const router = useRouter();
const { id } = router.currentRoute.value.query;
const wsUrl = `${getWsProtocol()}://${import.meta.env.VITE_BASE_URL}/Host/Terminal?id=${id}&token=${getLoginInfo()?.token}`
// const wsUrl = `${getWsProtocol()}://${import.meta.env.VITE_BASE_URL}/Host/TerminalSessionLog?id=4&token=${getLoginInfo()?.token}`

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


const handleTerminalCtrlU = () => {
	const wsUrl = `${getWsProtocol()}://${import.meta.env.VITE_BASE_URL}/Host/SFTPFileManage?id=${id}&token=${getLoginInfo()?.token}`
	fileManagerState.wsUrl = wsUrl
	fileManagerState.visible = true
}

const handleFileManagerClose = () => {
	fileManagerState.visible = false
}

const handleDownloadFile = async (path: string) => {
	window.open(`${getHttpProtocol()}://${import.meta.env.VITE_BASE_URL}/Host/DownloadFile?id=${id}&path=${path}&token=${getLoginInfo()?.token}`)
}

const handleTerminalReload = () => {
	router.go(0)
}


onMounted(async () => {
	if (!id) {
		router.push('/404')
		return
	}
	let res = await getApi(parseInt(id as string))
	document.title = `终端-${res.data.name}-${res.data.host}`
})

</script>

<style scoped></style>
