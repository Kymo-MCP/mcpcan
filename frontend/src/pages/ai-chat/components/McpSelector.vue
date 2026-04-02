<template>
  <el-dialog
    :model-value="visible"
    :title="title || t('aiChat.mcpInstances')"
    width="60%"
    center
    append-to-body
    @update:model-value="emit('update:visible', $event)"
    @open="initMcpList"
  >
    <div class="h-[500px] flex gap-4 mt-2">
      <!-- Left Column: MCP List (Multi-select) -->
      <div
        class="w-1/2 flex flex-col border rounded-lg overflow-hidden border-[var(--ep-border-color)]"
      >
        <div
          class="p-3 bg-[var(--ep-bg-color-page)] border-b border-[var(--ep-border-color)] flex justify-between items-center gap-2"
        >
          <span
            class="text-xs font-bold text-[var(--ep-text-color-secondary)] uppercase tracking-wider shrink-0"
            >{{ t('aiChat.instances') }}</span
          >
          <div class="flex-1" />
          <span class="text-xs text-[var(--ep-text-color-secondary)]">
            {{ selectedMcpIds.length }} {{ t('aiChat.selected') }}
          </span>
        </div>
        <div
          class="flex-1 overflow-y-auto p-2"
          v-infinite-scroll="loadMoreMcp"
          :infinite-scroll-disabled="mcpLoading || !mcpHasMore"
          :infinite-scroll-immediate="false"
        >
          <div
            v-if="!mcpInstances || mcpInstances.length === 0"
            class="flex flex-col items-center justify-center h-full text-[var(--ep-text-color-secondary)]"
          >
            {{ t('aiChat.noInstancesAvailable') }}
          </div>
          <div
            v-else-if="filteredMcpInstances.length === 0"
            class="flex flex-col items-center justify-center h-full text-[var(--ep-text-color-secondary)]"
          >
            {{ t('aiChat.noMatchesFound') }}
          </div>
          <div
            v-for="instance in filteredMcpInstances"
            :key="instance.instanceId"
            class="flex items-center p-3 rounded mb-1 cursor-pointer transition-colors border border-transparent"
            :class="
              selectedMcpIds.includes(instance.instanceId)
                ? 'bg-[var(--el-color-primary-light-9)] border-[var(--el-color-primary-light-5)]'
                : 'hover:bg-[var(--ep-fill-color)]'
            "
            @click="toggleMcpSelection(instance)"
          >
            <el-checkbox
              :model-value="selectedMcpIds.includes(instance.instanceId)"
              @change="toggleMcpSelection(instance)"
              @click.stop
              class="mr-3"
            />
            <div class="flex flex-col overflow-hidden m-l-2 flex-1">
              <span
                class="font-semibold text-sm truncate"
                :class="
                  selectedMcpIds.includes(instance.instanceId)
                    ? 'text-[var(--el-color-primary)]'
                    : ''
                "
                >{{ instance.name || instance.instanceName }}</span
              >
              <span class="text-xs text-[var(--ep-text-color-secondary)] truncate">
                {{
                  instance.mcpProtocol === 1 ? 'SSE' : instance.mcpProtocol === 2 ? 'HTTP' : 'STDIO'
                }}
                •
                <el-text
                  :type="activeOptions[instance.status as keyof typeof activeOptions]?.type"
                  link
                >
                  {{ activeOptions[instance.status as keyof typeof activeOptions]?.label }}
                </el-text>
              </span>

              <!-- Token Selection -->
              <div
                v-if="selectedMcpIds.includes(instance.instanceId) && instance.enabledToken"
                class="m-y-2 m-x-2"
                @click.stop
              >
                <el-select
                  v-model="selectedTokens[instance.instanceId]"
                  :placeholder="t('aiChat.selectToken')"
                  size="small"
                  class="w-full"
                  @change="handleTokenChange"
                  :loading="!instanceTokens[instance.instanceId] && instance.enabledToken"
                >
                  <el-option
                    v-for="token in instanceTokens[instance.instanceId] || []"
                    :key="token.token"
                    :value="token.token"
                    :disabled="!token.enabled"
                  >
                    <div class="max-w-100 flex flex-wrap gap-1 py-1">
                      <div v-if="token.usages && token.usages.length" class="flex flex-wrap gap-1">
                        <el-tag
                          v-for="(tag, num) in token.usages"
                          :key="num"
                          type="info"
                          size="small"
                          effect="light"
                        >
                          <div class="ellipsis-one max-w-25 text-xs">
                            {{ tag }}
                          </div>
                        </el-tag>
                      </div>
                      <el-tag
                        :type="token.enabled ? 'success' : 'danger'"
                        size="small"
                        effect="plain"
                        class="shrink-0"
                      >
                        {{ token.enabled ? t('aiChat.tokenActive') : t('aiChat.tokenInactive') }}
                      </el-tag>
                    </div>
                  </el-option>
                </el-select>
              </div>
            </div>
          </div>

          <!-- Load More Button -->
          <div v-if="mcpHasMore" class="mt-2">
            <el-button v-loading="mcpLoading" class="w-full" size="small" @click="loadMoreMcp">
              {{ t('aiChat.loadMore') }}
            </el-button>
          </div>
        </div>
      </div>

      <!-- Right Column: Config Preview -->
      <div
        class="w-1/2 flex flex-col border rounded-lg overflow-hidden border-[var(--ep-border-color)]"
      >
        <div class="p-3 bg-[var(--ep-bg-color-page)] border-b border-[var(--ep-border-color)]">
          <span
            class="text-xs font-bold text-[var(--ep-text-color-secondary)] uppercase tracking-wider"
            >{{ t('aiChat.mergedConfiguration') }}</span
          >
        </div>
        <div class="flex-1 overflow-hidden bg-[var(--ep-bg-color-deep)] relative">
          <MonacoEditor
            v-model="mergedMcpConfig"
            :options="{
              minimap: { enabled: false },
              scrollBeyondLastLine: false,
              automaticLayout: true,
              readOnly: false,
            }"
            language="json"
            :height="'100%'"
          />
          <el-tooltip :content="t('aiChat.copyConfig')" placement="top">
            <el-icon
              class="absolute top-2 right-4 z-10 cursor-pointer hover:text-[var(--el-color-primary)] bg-[var(--ep-bg-color)] rounded p-1"
              size="20"
              @click="copyConfig"
            >
              <CopyDocument />
            </el-icon>
          </el-tooltip>
        </div>
      </div>
    </div>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="emit('update:visible', false)">{{ t('aiChat.cancel') }}</el-button>
        <el-button type="primary" @click="confirmMcpSelection">{{ t('aiChat.confirm') }}</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { CopyDocument } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import MonacoEditor from '@/components/MonacoEditor/index.vue'
import { useI18n } from 'vue-i18n'
import type { InstanceResult } from '@/types/instance'
import { InstanceStatus, McpProtocol } from '@/types/instance'
import { TokenAPI } from '@/api/mcp/instance'

const { t } = useI18n()

const props = defineProps<{
  visible: boolean
  title?: string
  mcpInstances?: InstanceResult[]
  mcpConfig?: string
  mcpLoading?: boolean
  mcpHasMore?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:visible', visible: boolean): void
  (e: 'update:mcpConfig', config: string): void
  (e: 'save-mcp-config', config: string): void
  (e: 'load-mcp', page: number, append: boolean, query: string): void
}>()

const selectedMcpIds = ref<string[]>([])
const selectedTokens = ref<Record<string, string>>({})
const instanceTokens = ref<Record<string, any[]>>({})
const mcpPage = ref(1)

const manualConfig = ref('{}')
const isManualEdit = ref(false)
const isSyncingFromEditor = ref(false)

const activeOptions = {
  active: {
    label: t('status.' + InstanceStatus.ACTIVE),
    type: 'success',
    value: InstanceStatus.ACTIVE,
  },
  inactive: {
    label: t('status.' + InstanceStatus.INACTIVE),
    type: 'danger',
    value: InstanceStatus.INACTIVE,
  },
}

const filteredMcpInstances = computed(() => {
  if (!props.mcpInstances) return []
  return props.mcpInstances
})

const mergedMcpConfig = computed({
  get: () => manualConfig.value,
  set: (val) => {
    manualConfig.value = val
    emit('update:mcpConfig', val)
    // Sync left-side selection from editor content
    syncSelectionFromConfig(val)
  },
})

/**
 * Parse the JSON config from the editor and sync the left-side checkbox selection.
 * Reads `instanceId` from each server entry in `mcpServers` to determine which instances are selected.
 */
const syncSelectionFromConfig = (configStr: string) => {
  if (isSyncingFromEditor.value) return
  isSyncingFromEditor.value = true
  try {
    const conf = JSON.parse(configStr)
    if (conf.mcpServers) {
      const ids: string[] = []
      const configKeys = Object.keys(conf.mcpServers)

      for (const key of configKeys) {
        const server = conf.mcpServers[key]
        if (server.instanceId) {
          ids.push(server.instanceId)
          // Also sync token from config if present
          if (server.headers?.Authorization) {
            selectedTokens.value[server.instanceId] = server.headers.Authorization
          }
        }
      }

      // Fallback: for entries without instanceId, try matching accessType===1 instances by sourceConfig keys
      const unmatchedKeys = configKeys.filter((k) => !conf.mcpServers[k].instanceId)
      if (unmatchedKeys.length > 0 && props.mcpInstances) {
        for (const inst of props.mcpInstances) {
          if (ids.includes(inst.instanceId)) continue
          if (inst.accessType === 1) {
            try {
              const source = JSON.parse((inst as any).sourceConfig || inst.mcpConfig || '{}')
              if (source.mcpServers) {
                const sourceKeys = Object.keys(source.mcpServers)
                if (sourceKeys.some((k) => unmatchedKeys.includes(k))) {
                  ids.push(inst.instanceId)
                }
              }
            } catch (e) {}
          } else {
            // Try matching by serverKey pattern: mcp-{id前8位}
            const generatedKey = `mcp-${inst.instanceId.slice(0, 8)}`
            if (unmatchedKeys.includes(generatedKey)) {
              ids.push(inst.instanceId)
            }
          }
        }
      }

      // Only update if actually different to avoid unnecessary triggers
      const sorted1 = [...ids].sort().join(',')
      const sorted2 = [...selectedMcpIds.value].sort().join(',')
      if (sorted1 !== sorted2) {
        selectedMcpIds.value = ids
      }
    } else {
      if (selectedMcpIds.value.length > 0) {
        selectedMcpIds.value = []
      }
    }
  } catch (e) {
    // Invalid JSON, don't change selection
  } finally {
    isSyncingFromEditor.value = false
  }
}

watch(
  selectedMcpIds,
  () => {
    if (isSyncingFromEditor.value) return
    const val = generateMergedConfigFromSelection()
    manualConfig.value = val
    emit('update:mcpConfig', val)
  },
  { deep: true },
)

watch(
  () => props.mcpConfig,
  (val) => {
    if (val && val !== manualConfig.value) {
      manualConfig.value = val
    }
  },
  { immediate: true },
)

const generateInstanceConfig = (instance: InstanceResult) => {
  if (instance.accessType === 1) {
    try {
      const source = (instance as any).sourceConfig || instance.mcpConfig
      const parsed = JSON.parse(source || '{}')
      // Inject instanceId and instanceName into each server entry
      if (parsed.mcpServers) {
        for (const key of Object.keys(parsed.mcpServers)) {
          parsed.mcpServers[key].instanceId = instance.instanceId
          parsed.mcpServers[key].instanceName = instance.name || instance.instanceName
        }
      }
      return parsed
    } catch (e) {
      return { mcpServers: {} }
    }
  }

  const publicPath = (window as any).__APP_CONFIG__?.PUBLIC_PATH || ''
  const instancePath = (instance as any).publicProxyPath || ''
  const cleanPublicPath = publicPath

  let fullUrl = `${window.location.origin}${cleanPublicPath}${instancePath}`.replace(
    /([^:]\/)\/+/g,
    '$1',
  )

  const serverKey = `mcp-${instance.instanceId.slice(0, 8)}`

  const resolveServerType = () => {
    const protocol = (instance as any).proxyProtocol || instance.mcpProtocol

    if (protocol === McpProtocol.SSE || protocol === 1 || protocol === 'sse') {
      return 'sse'
    }

    if (
      protocol === McpProtocol.STREAMABLE_HTTP ||
      protocol === 2 ||
      protocol === 'streamable-http' ||
      protocol === 'streamable_http'
    ) {
      return 'streamable_http'
    }

    const lowerUrl = fullUrl.toLowerCase()
    if (lowerUrl.endsWith('/sse')) {
      return 'sse'
    }
    if (lowerUrl.endsWith('/mcp')) {
      return 'streamable_http'
    }
    return 'sse'
  }

  const serverConfig: any = {
    url: fullUrl,
    type: resolveServerType(),
    instanceId: instance.instanceId,
    instanceName: instance.name || instance.instanceName,
  }

  if (instance.enabledToken) {
    // Priority 1: Manually selected token
    if (selectedTokens.value[instance.instanceId]) {
      serverConfig.headers = {
        Authorization: selectedTokens.value[instance.instanceId],
      }
    } else {
      // Priority 2: Auto-select valid token from loaded tokens
      const tokens = instanceTokens.value[instance.instanceId] || instance.tokens
      if (tokens && tokens.length > 0) {
        const now = Date.now()
        const validToken = tokens.find((t: any) => t.enabled && (!t.expireAt || t.expireAt > now))

        if (validToken) {
          // Auto-select in UI too
          selectedTokens.value[instance.instanceId] = validToken.token
          serverConfig.headers = {
            Authorization: validToken.token,
          }
        }
      }
    }
  }

  const result = {
    mcpServers: {
      [serverKey]: serverConfig,
    },
  }

  return result
}

const generateMergedConfigFromSelection = () => {
  if (selectedMcpIds.value.length === 0) return '{}'
  const merged = { mcpServers: {} as any }
  selectedMcpIds.value.forEach((id) => {
    const instance = props.mcpInstances?.find((i) => i.instanceId === id)
    if (instance) {
      const conf = generateInstanceConfig(instance)
      if (conf && conf.mcpServers) {
        // Need to merge carefully, or overwrite?
        // Assuming unique keys by serverKey logic (mcp-id-prefix)
        Object.assign(merged.mcpServers, conf.mcpServers)
      }
    }
  })
  return JSON.stringify(merged, null, 4)
}

const toggleMcpSelection = async (instance: InstanceResult) => {
  const id = instance.instanceId
  if (selectedMcpIds.value.includes(id)) {
    selectedMcpIds.value = selectedMcpIds.value.filter((i) => i !== id)
    // Clear token selection when deselected? Maybe or keep it for re-selection
    delete selectedTokens.value[id]
  } else {
    selectedMcpIds.value.push(id)
    // Auto-select token if enabled
    if (instance.enabledToken) {
      // Load tokens if not already loaded or if using instance.tokens which might be empty
      if (!instanceTokens.value[id] || instanceTokens.value[id].length === 0) {
        try {
          // Fix: The response structure from TokenAPI.list is { tokens: [...] } not just []
          const response = await TokenAPI.list({
            instanceId: id,
          })
          // @ts-ignore
          const tokens = response.tokens || response
          // @ts-ignore
          instanceTokens.value[id] = Array.isArray(tokens) ? tokens : []
        } catch (e) {
          console.error('Failed to load tokens for instance', id, e)
        }
      }

      const tokens = instanceTokens.value[id] || instance.tokens
      if (tokens && tokens.length > 0) {
        const now = Date.now()
        // @ts-ignore
        const validToken = tokens.find((t: any) => t.enabled && (!t.expireAt || t.expireAt > now))
        if (validToken) {
          selectedTokens.value[id] = validToken.token
        }
      }
    }
  }
}

const handleTokenChange = () => {
  // Force re-generation of config
  const val = generateMergedConfigFromSelection()
  manualConfig.value = val
  emit('update:mcpConfig', val)
}

watch(
  selectedTokens,
  () => {
    const val = generateMergedConfigFromSelection()
    manualConfig.value = val
    emit('update:mcpConfig', val)
  },
  { deep: true },
)

const initMcpList = async () => {
  mcpPage.value = 1
  emit('load-mcp', 1, false, '')

  if (!props.mcpConfig || props.mcpConfig === '{}') {
    selectedMcpIds.value = []
  } else if (props.mcpInstances) {
    try {
      const conf = JSON.parse(props.mcpConfig)
      if (conf.mcpServers) {
        const existingServers = Object.keys(conf.mcpServers)
        const matchedIds: string[] = []

        // Build a set of instanceIds found in config for quick lookup
        const configInstanceIds = new Set<string>()
        for (const key of existingServers) {
          const server = conf.mcpServers[key]
          if (server.instanceId) {
            configInstanceIds.add(server.instanceId)
          }
        }

        for (const inst of props.mcpInstances) {
          const generatedKey = `mcp-${inst.instanceId.slice(0, 8)}`
          // Match by instanceId field in config OR by serverKey pattern
          const matchedByInstanceId = configInstanceIds.has(inst.instanceId)
          const matchedByKey = existingServers.includes(generatedKey)

          if (matchedByInstanceId || matchedByKey) {
            matchedIds.push(inst.instanceId)

            // Pre-load tokens and set selected token if enabled
            if (inst.enabledToken) {
              // Try to get token from config
              const matchedKey = matchedByInstanceId
                ? existingServers.find((k) => conf.mcpServers[k].instanceId === inst.instanceId)
                : generatedKey
              const currentConfig = matchedKey ? conf.mcpServers[matchedKey] : null
              if (currentConfig && currentConfig.headers && currentConfig.headers.Authorization) {
                selectedTokens.value[inst.instanceId] = currentConfig.headers.Authorization
              }

              // Check if we need to load tokens
              if (!instanceTokens.value[inst.instanceId]) {
                try {
                  const response = await TokenAPI.list({ instanceId: inst.instanceId })
                  // @ts-ignore
                  const tokens = response.tokens || response
                  // @ts-ignore
                  instanceTokens.value[inst.instanceId] = Array.isArray(tokens) ? tokens : []
                } catch (e) {}
              }

              // If not set from config, try auto-select
              if (!selectedTokens.value[inst.instanceId]) {
                const tokens = instanceTokens.value[inst.instanceId] || inst.tokens
                if (tokens && tokens.length > 0) {
                  const now = Date.now()
                  // @ts-ignore
                  const validToken = tokens.find(
                    (t: any) => t.enabled && (!t.expireAt || t.expireAt > now),
                  )
                  if (validToken) selectedTokens.value[inst.instanceId] = validToken.token
                }
              }
            }
          } else if (inst.accessType === 1) {
            try {
              const source = JSON.parse((inst as any).sourceConfig || inst.mcpConfig || '{}')
              if (source.mcpServers) {
                const sourceKeys = Object.keys(source.mcpServers)
                if (sourceKeys.some((k) => existingServers.includes(k))) {
                  matchedIds.push(inst.instanceId)
                }
              }
            } catch (e) {}
          }
        }
        if (matchedIds.length > 0) {
          selectedMcpIds.value = matchedIds
        }
      }
    } catch (e) {
      selectedMcpIds.value = []
    }
  }
}

const loadMoreMcp = () => {
  if (props.mcpLoading || !props.mcpHasMore) return
  mcpPage.value++
  emit('load-mcp', mcpPage.value, true, '')
}

const confirmMcpSelection = () => {
  emit('update:mcpConfig', mergedMcpConfig.value)
  emit('save-mcp-config', mergedMcpConfig.value)
  emit('update:visible', false)
}

const copyConfig = async () => {
  try {
    await navigator.clipboard.writeText(mergedMcpConfig.value)
    ElMessage.success(t('common.copySuccess') || 'Copied')
  } catch (err) {
    ElMessage.error(t('common.copyFailed') || 'Copy failed')
  }
}
</script>

<style scoped>
/* Add any necessary styles here */
</style>
