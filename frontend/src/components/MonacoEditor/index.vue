<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import * as monaco from 'monaco-editor'

// Vite 环境下 Worker 配置
import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'
import jsonWorker from 'monaco-editor/esm/vs/language/json/json.worker?worker'
import { useSystemStoreHook } from '@/stores/modules/system-store'

const systemStore = useSystemStoreHook()
// 设置 Worker 环境
self.MonacoEnvironment = {
  getWorker(_, label) {
    if (label === 'json') {
      return new jsonWorker()
    }
    return new editorWorker()
  },
}

const props = withDefaults(
  defineProps<{
    modelValue: string
    language?: string
    height?: string
    readOnly?: boolean
  }>(),
  {
    language: 'json',
    height: '400px',
    readOnly: false,
  },
)

const emit = defineEmits(['update:modelValue', 'change'])

const editorRef = ref<HTMLElement | null>(null)
let editorInstance: monaco.editor.IStandaloneCodeEditor | null = null
let isSettingValue = false // 防止死循环

// 初始化编辑器
const initEditor = () => {
  if (!editorRef.value) return
  editorInstance = monaco.editor.create(editorRef.value, {
    value: props.modelValue,
    language: props.language,
    theme: systemStore.themeType === 'dark' ? 'vs-dark' : 'vs',
    automaticLayout: true, // 自动布局
    readOnly: props.readOnly,
    minimap: { enabled: false }, // 关闭小地图
    scrollBeyondLastLine: false,
    fontSize: 14,
    fontFamily: "'Fira Code', 'Consolas', monospace",
    formatOnPaste: true,
    tabSize: 2,
    lineNumbers: 'on', // 行数显示
    lineNumbersMinChars: 3, // 缩短行号宽度
    folding: true, // 代码折叠
    foldingStrategy: 'indentation', // 折叠策略
    padding: { top: 10 },
  })

  // 监听内容变化
  editorInstance.onDidChangeModelContent(() => {
    if (editorInstance) {
      const value = editorInstance.getValue()
      isSettingValue = true
      emit('update:modelValue', value)
      emit('change', value)
      isSettingValue = false
    }
  })

  // 初始格式化 (如果是 JSON)
  if (props.language === 'json') {
    setTimeout(() => {
      editorInstance?.getAction('editor.action.formatDocument')?.run()
    }, 300)
  }
}

// 监听值变化
watch(
  () => props.modelValue,
  (newValue) => {
    if (editorInstance && !isSettingValue && newValue !== editorInstance.getValue()) {
      editorInstance.setValue(newValue)
    }
  },
)

// 监听主题变化
watch(
  () => systemStore.themeType,
  (newTheme) => {
    if (editorInstance) {
      monaco.editor.setTheme(systemStore.themeType === 'dark' ? 'vs-dark' : 'vs')
    }
  },
)

// 监听语言变化
watch(
  () => props.language,
  (newLang) => {
    if (editorInstance) {
      const model = editorInstance.getModel()
      if (model) {
        monaco.editor.setModelLanguage(model, newLang)
      }
    }
  },
)

onMounted(() => {
  nextTick(() => {
    initEditor()
  })
})

onBeforeUnmount(() => {
  if (editorInstance) {
    editorInstance.dispose()
  }
})
</script>

<template>
  <div
    ref="editorRef"
    class="monaco-editor-container"
    :style="{
      height: height,
      width: '100%',
      border: '1px solid var(--el-border-color)',
      borderRadius: '4px',
    }"
  ></div>
</template>

<style scoped>
.monaco-editor-container {
  overflow: hidden;
}
</style>
