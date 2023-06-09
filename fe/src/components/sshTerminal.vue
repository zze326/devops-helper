<template >
    <div style="position: relative;">
        <div ref="terminalRef"></div>
        <input v-if="searchMode" type="text" v-model="searchText" ref="searchInputRef" id="search-input"
            @keyup.enter="search" @keyup.esc="cancelSearch" placeholder="搜索..." />
    </div>
</template>

<script setup lang="ts">
import { reactive, toRefs, onMounted, ref, onUnmounted, nextTick, defineExpose } from 'vue'
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import { WebLinksAddon } from 'xterm-addon-web-links';
import { SearchAddon } from 'xterm-addon-search';
import 'xterm/css/xterm.css'
import { ElMessage } from 'element-plus'
import _ from 'lodash'

const props = defineProps(
    {
        wsUrl: {
            type: String,
            required: true
        },
        inBody: {
            type: Boolean,
            default: false
        },
        paddingBottom: {
            type: Number,
            default: 0
        },
        innerData: {
            type: Boolean,
            default: false
        }
    })

const emits = defineEmits(['close', 'connected', 'reload', 'ctrlU', "message"])

const radiusEm = props.inBody ? 0 : 0.5

const terminalRef = ref<HTMLDivElement>()
const term = new Terminal({
    fontSize: 15,
    cursorBlink: true,
    convertEol: true,
    rightClickSelectsWord: true,
    allowTransparency: true,
    theme: {
        foreground: 'white',
        background: '#2B394C',
    }
})

const state = reactive({
    searchMode: false,
    searchText: '',
    connected: false,
})

const search = () => {
    if (!searchAddon) return
    searchAddon.findNext(searchText.value)
}

const { searchMode, searchText } = toRefs(state)
const searchInputRef = ref<HTMLInputElement>()
const cancelSearch = () => {
    if (!searchInputRef.value) return
    searchInputRef.value.value = "";
    searchInputRef.value!.blur();
    searchMode.value = false;
    term.focus();
}

term.attachCustomKeyEventHandler((event) => {
    console.log("searchMode:", searchMode)
    if (event.key === "f" && event.ctrlKey) {
        searchText.value = '';
        searchMode.value = true;
        nextTick(() => {
            searchInputRef.value!.focus();
        });
        return false;
    } else if (event.key === "Escape" && searchMode.value) {
        cancelSearch()
        return false;
    } else if (event.key === "u" && event.ctrlKey && event.type !== "keyup") {
        emits("ctrlU")
        return false;
    }
    return true;
});

const fitAddon = new FitAddon();
const searchAddon = new SearchAddon();
term.loadAddon(fitAddon);
term.loadAddon(new WebLinksAddon());
term.loadAddon(searchAddon);

const resizeTerminal = () => {
    if (!terminalRef.value) return
    state.connected = true
    terminalRef.value.style.height = props.inBody ? `${document.body.clientHeight - props.paddingBottom}px` : `${(terminalRef.value!.offsetParent as HTMLDivElement)!.offsetParent!.clientHeight - props.paddingBottom}px`
    fitAddon.fit();
    term.resize(term.cols, term.rows)
    if (ws.readyState) {
        var msg = { type: "resize", rows: term.rows, cols: term.cols }
        ws.send(JSON.stringify(msg))
    }
}

const onOpen = () => {
    resizeTerminal()
    nextTick(() => {
        emits("connected")
    })
}

const ws = new WebSocket(props.wsUrl)

const init = () => {
    if (!terminalRef.value) return
    term.open(terminalRef.value);
    term.focus()
    term.writeln('Connecting...');
    ws.onopen = onOpen
    ws.onclose = () => {
        emits("close")
        state.connected = false
        term.write('\n\x1b[31m终端已断开\x1b[0m\n');
        ElMessage.warning({
            message: 'Websocket 连接已关闭',
            type: 'warning',
            center: true,
        });
    }
    ws.onmessage = async e => {
        if (props.innerData) {
            let msgObj = JSON.parse(e.data)
            emits("message", msgObj)
            if (msgObj.clear) {
                term.reset()
            }
            term.write(msgObj.data)
        } else {
            emits("message", e.data)
            term.write(e.data)
        }
    }
    ws.onerror = () => {
        emits("close")
        state.connected = false
        term.write('\x1b[31m终端已断开\x1b[0m\n');
        ElMessage.error({
            message: '建立 Websocket 连接失败',
            type: 'error',
            center: true,
        });
    }

    window.addEventListener("resize", _.debounce(resizeTerminal, 300))

    term.onData(function (input) {
        if (state.connected) {
            var msg = { type: "input", input: input }
            ws.send(JSON.stringify(msg))
        } else {
            if (input === "\r" && !state.connected) {
                emits("reload")
            }
        }
    })
}

const wsSend = (msg: { type: string, input: string }) => {
    if (!state.connected) return
    ws.send(JSON.stringify(msg))
}

const termWrite = (data: string) => {
    term.write(data)
}

const termClear = () => {
    term.clear()
}

defineExpose({
    wsSend,
    termWrite,
    termClear
})

onMounted(() => {
    init()
})
onUnmounted(() => {
    var msg = { type: "close", input: '' }
    ws.send(JSON.stringify(msg))
    ws.close()
})
</script>

<style scoped>
:deep(.xterm) {
    padding: 1em;
    position: unset;
}

:deep(.xterm-viewport) {
    border-radius: v-bind(radiusEm + 'em');
}

#search-input {
    position: absolute;
    top: 10px;
    right: 10px;
    z-index: 1;
    padding: 5px;
    border: none;
    background-color: rgba(255, 255, 255, 0.8);
    border-radius: 3px;
}
</style>
