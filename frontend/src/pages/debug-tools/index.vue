<template>
  <div class="tool-page h-full flex flex-col" v-loading="loading">
    <!-- Instance Detail Info Bar -->
    <el-card class="mb-4 text-sm" shadow="hover">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-4">
          <mcp-image
            :src="instanceInfo.iconPath"
            width="32"
            height="32"
            :key="instanceInfo.instanceId"
          ></mcp-image>
          <div class="flex items-center gap-2">
            <span class="font-bold text-lg">{{
              instanceInfo.instanceName || 'Unknown Instance'
            }}</span>
            <el-tag
              :type="activeOptions[instanceInfo.status as keyof typeof activeOptions]?.type"
              effect="dark"
              size="small"
              round
            >
              {{
                activeOptions[instanceInfo.status as keyof typeof activeOptions]?.label ||
                'Unknown status'
              }}
            </el-tag>
          </div>
          <el-divider direction="vertical" />
          <div class="text-gray-500 u-line-1">ID: {{ instanceInfo.instanceId }}</div>
          <el-divider direction="vertical" />
          <div class="text-gray-500 u-line-1">Base URL: {{ configUrl }}</div>
        </div>
        <div>
          <el-button v-if="layout" @click="handleBack" class="link-hover">
            <el-icon class="mr-2">
              <i class="icon iconfont MCP-fanhui"></i>
            </el-icon>
            {{ t('common.back') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Main Splitter Area -->
    <div class="flex-1 overflow-hidden">
      <el-splitter class="h-full border rounded-lg bg-page overflow-hidden">
        <!-- Panel 1: Tool List -->
        <el-splitter-panel size="350px" class="flex flex-col" :resizable="false">
          <div class="p-3 border-b bg-header font-medium flex justify-between">
            <span>{{ t('mcp.debugTool.list') }}</span>
            <el-tag size="small" type="info">{{ toolList.length }}</el-tag>
          </div>
          <div class="p-2 border-b">
            <el-input
              v-model="keyword"
              :placeholder="t('mcp.debugTool.keyPlaceholder')"
              :prefix-icon="Search"
              clearable
            />
          </div>
          <el-scrollbar class="flex-1">
            <div
              v-for="tool in filteredTools"
              :key="tool.name"
              class="p-3 cursor-pointer tool-item transition-colors border-l-4 border-transparent"
              :class="{ 'active-tool': currentTool?.name === tool.name }"
              @click="handleSelectTool(tool)"
            >
              <div class="font-bold">{{ tool.name }}</div>
              <div class="text-xs mt-1 line-clamp-2">{{ tool.description }}</div>
            </div>
            <el-empty v-if="filteredTools.length === 0" />
          </el-scrollbar>
        </el-splitter-panel>

        <!-- Panel 2: Input Parameters -->
        <el-splitter-panel class="flex flex-col relative">
          <div class="p-3 border-b bg-header font-medium flex justify-between items-center">
            <span>{{ t('mcp.debugTool.input') }}</span>
            <el-switch
              v-if="currentTool"
              v-model="jsonMode"
              active-text="JSON"
              inactive-text="Form"
              size="small"
            />
          </div>
          <el-scrollbar class="flex-1 p-4 px-2">
            <template v-if="currentTool">
              <div class="mb-4 mx-2">
                <div class="text-sm text-secondary mb-2">Description</div>
                <div class="bg-block p-2 rounded text-sm text-regular">
                  {{ currentTool.description || 'No description provided.' }}
                </div>
              </div>

              <!-- Form Mode -->
              <div v-if="!jsonMode && currentTool.inputSchema?.properties" class="mx-2">
                <SchemaForm
                  v-model="paramsForm"
                  :schema="currentTool.inputSchema"
                  @update:modelValue="handleFormChange"
                />
              </div>

              <!-- JSON Mode or No Schema -->
              <div v-else class="h-full flex flex-col mx-2">
                <div class="mb-2 text-sm font-bold text-regular">Arguments (JSON)</div>
                <div class="border rounded flex-1 relative min-h-[300px]">
                  <el-input
                    v-model="inputJson"
                    type="textarea"
                    :rows="15"
                    placeholder="{}"
                    class="font-mono text-sm h-full"
                    resize="none"
                    @input="handleJsonChange"
                  />
                  <div
                    v-if="jsonError"
                    class="text-error text-xs mt-1 absolute bottom-2 left-2 right-2 bg-card p-1 opacity-90 border border-red-200"
                  >
                    {{ jsonError }}
                  </div>
                </div>
              </div>
            </template>
            <el-empty v-else />
          </el-scrollbar>
        </el-splitter-panel>

        <!-- Panel 3: Operation Bar -->
        <el-splitter-panel
          size="200px"
          class="flex flex-col bg-block border-l border-r"
          :resizable="false"
        >
          <div class="p-3 border-b bg-header font-medium text-center">
            {{ t('mcp.debugTool.action') }}
          </div>
          <div class="flex-1 flex flex-col items-center p-4 gap-4">
            <el-button
              type="primary"
              size="large"
              class="w-full base-btn"
              :icon="VideoPlay"
              :loading="running"
              :disabled="!currentTool || !!jsonError"
              @click="handleRunTool"
            >
              {{ t('mcp.debugTool.run') }}
            </el-button>
            <el-button
              type="primary"
              size="large"
              class="w-full base-btn !m-l-0"
              :icon="DocumentCopy"
              :disabled="!inputJson"
              @click="handleCopyInput"
            >
              {{ t('mcp.debugTool.copyInput') }}
            </el-button>
            <el-button
              type="primary"
              size="large"
              class="w-full base-btn !m-l-0"
              :icon="DocumentCopy"
              :disabled="!outputResult"
              @click="handleCopyOutput"
            >
              {{ t('mcp.debugTool.copyOutput') }}
            </el-button>

            <!-- <el-divider content-position="center" class="!my-2">History</el-divider>

            <el-scrollbar class="w-full flex-1">
              <div
                v-for="(item, idx) in history"
                :key="idx"
                class="history-item text-xs p-2 mb-2 rounded bg-card shadow-sm border cursor-pointer hover:bg-hover"
                @click="restoreHistory(item)"
              >
                <div class="font-bold truncate text-regular">{{ item.tool }}</div>
                <div class="text-secondary">{{ formatTime(item.timestamp) }}</div>
                <div class="mt-1">
                  <el-tag size="small" :type="item.status === 'success' ? 'success' : 'danger'">{{
                    item.status
                  }}</el-tag>
                </div>
              </div>
            </el-scrollbar> -->
          </div>
        </el-splitter-panel>

        <!-- Panel 4: Output Parameters -->
        <el-splitter-panel class="flex flex-col">
          <div class="p-3 border-b bg-header font-medium flex justify-between items-center">
            <span>{{ t('mcp.debugTool.output') }}</span>
            <el-button v-if="outputResult" size="small" link @click="clearOutput">Clear</el-button>
          </div>
          <el-scrollbar class="flex-1 p-0">
            <div v-if="parsedOutput" class="p-4">
              <!-- Result Status -->
              <el-alert
                v-if="parsedOutput.result?.isError"
                :title="t('mcp.debugTool.runFail')"
                type="error"
                show-icon
                :closable="false"
                class="mb-4"
              />
              <el-alert
                v-else
                :title="t('mcp.debugTool.runSuccess')"
                type="success"
                show-icon
                :closable="false"
                class="mb-4"
              />

              <!-- Content List -->
              <div
                v-if="parsedOutput.result?.content && parsedOutput.result.content.length"
                class="flex flex-col gap-4"
              >
                <div
                  v-for="(item, idx) in parsedOutput.result.content"
                  :key="idx"
                  class="rounded border bg-block overflow-hidden"
                >
                  <!-- Text Content -->
                  <div
                    v-if="item.type === 'text'"
                    class="p-3 text-sm font-mono whitespace-pre-wrap break-words text-regular"
                  >
                    {{ item.text }}
                  </div>

                  <!-- Image Content -->
                  <div v-else-if="item.type === 'image'" class="p-3 flex justify-center bg-card">
                    <img
                      :src="`data:${item.mimeType};base64,${item.data}`"
                      class="max-w-full max-h-[400px] object-contain border rounded"
                    />
                  </div>

                  <!-- Fallback -->
                  <div v-else class="p-3 text-sm text-secondary italic">
                    Unknown content type: {{ item.type }}
                  </div>
                </div>
              </div>
              <div v-else class="text-center text-secondary py-8">No content returned</div>

              <!-- Raw Details Toggle -->
              <div class="mt-6 border-t pt-4">
                <div
                  class="flex items-center gap-1 text-xs text-secondary cursor-pointer select-none font-bold"
                  @click="showRawOutput = !showRawOutput"
                >
                  <el-icon><ArrowRight v-if="!showRawOutput" /><ArrowDown v-else /></el-icon>
                  <span>Raw JSON Output</span>
                </div>
                <div v-if="showRawOutput" class="mt-2 text-xs">
                  <div class="bg-block p-2 rounded border font-mono whitespace-pre-wrap break-all">
                    {{ outputResult }}
                  </div>
                </div>
              </div>
            </div>
            <el-empty v-else />
          </el-scrollbar>
        </el-splitter-panel>
      </el-splitter>
    </div>
  </div>
</template>

<script setup lang="ts">
import {
  Search,
  Refresh,
  VideoPlay,
  ArrowRight,
  ArrowDown,
  DocumentCopy,
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useDebugToolsHooks } from './hooks/index.ts'
import McpImage from '@/components/mcp-image/index.vue'
import SchemaForm from './components/SchemaForm.vue'
import {
  generateDefaultValue,
  isPropertyRequired,
  resolveRef,
  normalizeUnionType,
} from '@/utils/schemaUtils'
import { deBugAPI } from '@/api/mcp/instance.ts'
import { AccessType } from '@/types'
import { setClipboardData } from '@/utils/system'
import { useRouterHooks } from '@/utils/url'

const {
  activeOptions,
  instanceInfo,
  toolList,
  currentTool,
  keyword,
  inputJson,
  outputResult,
  history,
  route,
  loading,
  running,
  instanceId,
  t,
} = useDebugToolsHooks()
const layout = useLayout()
const { jumpBack } = useRouterHooks()
const configUrl = computed(() => {
  if (instanceInfo.value.accessType === AccessType.DIRECT) {
    const mcpServers = JSON.parse(instanceInfo.value.sourceConfig).mcpServers
    return mcpServers[Object.keys(mcpServers)[0]].url
  }
  return `${window.location.origin}${(window as any).__APP_CONFIG__?.PUBLIC_PATH}${instanceInfo.value.publicProxyPath}`
})
// Computed
const filteredTools = computed(() => {
  if (!keyword.value) return toolList.value
  const lower = keyword.value.toLowerCase()
  return toolList.value.filter(
    (t) =>
      t.name.toLowerCase().includes(lower) ||
      (t.description && t.description.toLowerCase().includes(lower)),
  )
})

const showRawOutput = ref(false)
const parsedOutput = computed(() => {
  if (!outputResult.value) return null
  try {
    return JSON.parse(outputResult.value)
  } catch (e) {
    return null
  }
})

const jsonError = computed(() => {
  try {
    JSON.parse(inputJson.value)
    return ''
  } catch (e: any) {
    return e.message
  }
})

// Methods
const formatTime = (ts: number) => {
  return new Date(ts).toLocaleTimeString()
}

// handle get tool list
const getTools = async () => {
  try {
    loading.value = true
    const list = await deBugAPI.toolList({
      instanceId: instanceId.value || '',
      domain: configUrl.value,
    })
    toolList.value = list || []
  } finally {
    loading.value = false
  }
}

const paramsForm = ref<Record<string, any>>({})
const jsonMode = ref(false)

const handleSelectTool = (tool: any) => {
  currentTool.value = tool
  jsonMode.value = false
  // Generate default values from schema
  const defaults: Record<string, any> = {}
  if (tool.inputSchema?.properties) {
    Object.entries(tool.inputSchema.properties).forEach(([key, schema]: [string, any]) => {
      const resolved = resolveRef(schema, tool.inputSchema)
      const val = generateDefaultValue(resolved, key, tool.inputSchema)
      if (val !== undefined) {
        defaults[key] = val
      }
    })
  }

  paramsForm.value = defaults
  inputJson.value = JSON.stringify(defaults, null, 2)
  outputResult.value = ''
}

const handleFormChange = () => {
  inputJson.value = JSON.stringify(paramsForm.value, null, 2)
}

const handleJsonChange = (val: string) => {
  try {
    paramsForm.value = JSON.parse(val)
  } catch (e) {
    // Ignore parse errors while typing
  }
}

const handleRunTool = async () => {
  if (jsonError.value) return
  try {
    running.value = true
    const params = JSON.parse(inputJson.value)
    const data = await deBugAPI.toolCall({
      instanceId: instanceId.value || '',
      toolName: currentTool.value.name,
      domain: configUrl.value,
      arguments: JSON.stringify(params),
    })

    const displayOutput = {
      tool: currentTool.value.name,
      arguments: params,
      result: normalizeUnionType(data),
    }
    outputResult.value = JSON.stringify(displayOutput, null, 2)
    // Add to history
    history.value.unshift({
      tool: currentTool.value.name,
      arguments: inputJson.value,
      timestamp: Date.now(),
      status: data.isError ? 'error' : 'success',
    })
  } catch (err) {
    outputResult.value = JSON.stringify({ error: 'Execution Failed', details: err }, null, 2)
  } finally {
    running.value = false
  }
}

const restoreHistory = (item: any) => {
  const tool = toolList.value.find((t) => t.name === item.tool)
  if (tool) {
    currentTool.value = tool
    inputJson.value = item.params
  }
}

const refreshInstance = () => {
  getTools()
}

const clearOutput = () => {
  outputResult.value = ''
}

const handleCopyInput = () => {
  setClipboardData(inputJson.value)
  ElMessage.success(t('action.copy'))
}

const handleCopyOutput = () => {
  setClipboardData(outputResult.value)
  ElMessage.success(t('action.copy'))
}

// back last class page
const handleBack = () => {
  jumpBack()
}

onMounted(() => {
  getTools()
})
</script>

<style lang="scss" scoped>
.tool-page {
  /* Fix height to fill remaining viewport if needed, assuming header is ~60px */
  height: calc(100vh - 60px);
  margin: -20px;
  padding: 16px;
  background-color: var(--ep-bg-color-page);
}

.bg-page {
  border: 1px solid var(--ep-border-color);
}
.bg-header {
  background-color: var(--ep-fill-color);
  border-bottom: 1px solid var(--ep-border-color-light);
}
.bg-block {
  background-color: var(--ep-fill-color);
  border: 1px solid var(--ep-border-color-lighter);
}
.bg-card {
  background-color: var(--ep-bg-color-overlay);
}
.text-primary {
  color: var(--ep-text-color-primary);
}
.text-regular {
  color: var(--ep-text-color-regular);
}
.text-secondary {
  color: var(--ep-text-color-secondary);
}
.text-placeholder {
  color: var(--ep-text-color-placeholder);
}

.hover\:bg-hover:hover {
  background-color: var(--ep-fill-color);
}
.tool-item {
  border-bottom: 1px solid var(--ep-border-color-lighter);
  transition: all 0.3s ease;
}
.tool-item:hover {
  background-color: var(--el-color-primary-hover);
  transform: scale(1.02);
}

.active-tool {
  background-color: var(--el-color-primary);
}
:deep(.el-splitter) {
  --el-splitter-border-color: var(--ep-border-color);
}
</style>
