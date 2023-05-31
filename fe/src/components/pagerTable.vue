<template>
    <div>
        <ElForm v-if="withSearch" ref="queryForm" @submit.prevent :inline="true">
            <ElFormItem label="模糊搜索" prop="search">
                <el-input icon="Plus" ref="searchInput" v-model="pagerState.search" @keyup.enter='listPage'
                    :placeholder="placeholder" clearable />
            </ElFormItem>
            <slot name="queryFormItem"></slot>
            <ElFormItem>
                <ElButton type="primary" icon="Search" @click="listPage">搜索</ElButton>
                <ElButton icon="RefreshLeft" @click="resetClick">重置
                </ElButton>
                <ElButton type="success" icon="Plus" v-if="!disableAdd" v-permiss="addPermissCode" @click="addClick">添加
                </ElButton>
                <slot name="queryFormBtn"></slot>
            </ElFormItem>
        </ElForm>

        <el-table v-loading="loading" :data="pagerState.items" stripe style="width: 100%">
            <slot></slot>
        </el-table>
        <el-pagination style="float: right; margin-top: 10px;margin-bottom: 35px;" v-model:current-page="pagerState.page"
            v-model:page-size="pagerState.page_size" :page-sizes="[10, 20, 30, 50, 100]" :disabled="pagerState.total === 0"
            :background="true" layout="total, sizes, prev, pager, next, jumper" :total="pagerState.total"
            @size-change="listPage" @current-change="listPage" />
    </div>
</template>

<script setup lang="ts">
import { PagerState } from '@/utils/pager';
const props = defineProps({
    "pagerState": {
        type: PagerState,
        required: true
    },
    "withSearch": {
        type: Boolean,
        default: true
    },
    "addPermissCode": {
        type: String,
        default: ""
    },
    "loading": {
        type: Boolean,
        default: false
    },
    "placeholder": {
        type: String,
        default: '名称'
    },
    "disableAdd": {
        type: Boolean,
        default: false
    }
})

const emits = defineEmits(["listPage", "addClick", "resetClick"])

const listPage = () => {
    emits("listPage")
}

const addClick = () => {
    emits("addClick")
}

const resetClick = () => {
    if (emits("resetClick") === undefined) {
        props.pagerState.search = ""
        emits("listPage")
    }
}
</script>
<style scoped></style>
