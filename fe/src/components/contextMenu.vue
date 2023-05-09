<template>
    <teleport to="body">
        <div v-if="visible" class="context-menu" :style="{ top: `${y}px`, left: `${x}px` }" @mouseleave="handleMouseLeave">
            <ul>
                <li v-permiss="item.permiss" :style="{ color: item.key.includes('delete') ? 'red' : '' }"
                    v-for="(item, index) in items" :key="index" @click="handleItemClick(item)">
                    <el-icon>
                        <component :is="item.icon"></component>
                    </el-icon>
                    <span>{{ item.label }}</span>
                </li>
            </ul>
        </div>
    </teleport>
</template>
  
<script lang="ts" setup>
import { Teleport } from 'vue';

interface MenuItem {
    label: string;
    icon?: string;
    key: string;
    permiss?: string;
}

const emits = defineEmits(['close', 'itemClick']);

defineProps({
    x: {
        type: Number,
        default: 0
    },
    y: {
        type: Number,
        default: 0
    },
    visible: {
        type: Boolean,
        default: false
    },
    items: {
        type: Array<MenuItem>,
        default: [],
        required: true
    }
});

const handleMouseLeave = () => {
    emits('close');
};

const handleItemClick = (item: MenuItem) => {
    emits('itemClick', item.key);
};
</script>
  
<style scoped>
.context-menu {
    position: absolute;
    background-color: #fff;
    border-radius: 4px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
    z-index: 9999;
}

.context-menu ul {
    list-style: none;
    padding: 0;
    margin: 0;
}

.context-menu li {
    display: flex;
    align-items: center;
    padding: 8px 16px;
    cursor: pointer;
}

.context-menu li:hover {
    background-color: #f6f7f8;
}

.context-menu li+li {
    border-top: 1px solid #e4e7ed;
}

.context-menu li span {
    margin-left: 8px;
    font-size: 13px;
}
</style>