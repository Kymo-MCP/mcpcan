<template>
  <el-dialog v-model="dialogInfo.visible" width="480px" align-center :show-close="false">
    <template #header>
      <div class="center">{{ dialogInfo.title }}</div>
    </template>
    <el-scrollbar ref="scrollbarRef" always class="config-info">
      <div class="py-5 px-5">{{ config }}</div>
      <el-tooltip
        class="box-item"
        effect="dark"
        :content="t('mcp.instance.token.copyUrl')"
        placement="top"
      >
        <el-icon class="base-btn-link copy-icon-url" size="18" @click="handleCopy('url')">
          <Link />
        </el-icon>
      </el-tooltip>
      <el-tooltip
        class="box-item"
        effect="dark"
        :content="t('mcp.instance.token.copyToken')"
        placement="top"
      >
        <el-icon class="base-btn-link copy-icon-token" size="18" @click="handleCopy('token')">
          <Key />
        </el-icon>
      </el-tooltip>
      <el-tooltip
        class="box-item"
        effect="dark"
        :content="t('mcp.instance.token.copyAll')"
        placement="top"
      >
        <el-icon class="base-btn-link copy-icon" size="18" @click="handleCopy('config')">
          <CopyDocument />
        </el-icon>
      </el-tooltip>
      <div class="custom-style my-4 px-5">
        <el-segmented v-model="dialogInfo.pathType" :options="pathTypeOptions" />
      </div>
    </el-scrollbar>
    <template #footer>
      <div class="center">
        <el-button class="w-25 mr-2" @click="dialogInfo.visible = false">{{
          t('common.close')
        }}</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { Operation, CopyDocument, Link, Key, Edit } from '@element-plus/icons-vue'
import { setClipboardData, timestampToDate } from '@/utils/system'
import { ElMessage, ElMessageBox } from 'element-plus'
import { AccessType, TokenType, type InstanceResult } from '@/types'
import { JsonFormatter } from '@/utils/json'
import { InstanceAPI, TokenAPI } from '@/api/mcp/instance'

const { t } = useI18n()
const dialogInfo = ref<any>({
  visible: false,
  loading: false,
  title: '访问配置',
  instanceInfo: {
    instanceId: '',
  },
  pathType: 'steamable_http',
})
const pathTypeOptions = [
  { label: 'SSE', value: 'sse' },
  { label: 'STREAMABLE_HTTP', value: 'steamable_http' },
]
// token list
const tokenList = ref<Array<any>>([])
const configUrl = computed(() => {
  if (dialogInfo.value.instanceInfo.accessType === AccessType.DIRECT) {
    const mcpServers = JSON.parse(dialogInfo.value.instanceInfo.sourceConfig).mcpServers
    return mcpServers[Object.keys(mcpServers)[0]].url
  }
  let publicPath = dialogInfo.value.instanceInfo.publicProxyPath
  if (dialogInfo.value.pathType) {
    const lastSlashIndex = publicPath.lastIndexOf('/')
    if (lastSlashIndex !== -1) {
      publicPath =
        publicPath.substring(0, lastSlashIndex + 1) +
        (dialogInfo.value.pathType === 'sse' ? 'sse' : 'mcp')
    }
  }
  return `${window.location.origin}${(window as any).__APP_CONFIG__?.PUBLIC_PATH}${publicPath}`
})
const configToken = computed(() => {
  if (dialogInfo.value.instanceInfo.accessType === AccessType.DIRECT) {
    const mcpServers = JSON.parse(dialogInfo.value.instanceInfo.sourceConfig).mcpServers
    return mcpServers[Object.keys(mcpServers)[0]].token || 'No Data'
  }
  return `${tokenList.value[0].token}`
})
// config Info
const config = computed(() => {
  if (dialogInfo.value.instanceInfo.accessType === AccessType.DIRECT) {
    return JsonFormatter.format(dialogInfo.value.instanceInfo.sourceConfig, 4)
  }
  if (dialogInfo.value.instanceInfo.enabledToken) {
    if (!tokenList.value) return JsonFormatter.format(`{}`, 4)
    if (tokenList.value.length) {
      return JsonFormatter.format(
        `{
          "mcpServers": {
                "mcp-${dialogInfo.value.instanceInfo.instanceId.slice(0, 8)}": {
                      "url": "${configUrl.value}",
                      "headers": {
                            "Authorization": "${tokenList.value[0]?.token}"
                      }
                }
          }
      }`,
        4,
      )
    }
  }
  return JsonFormatter.format(
    `{
      "mcpServers": {
          "mcp-${dialogInfo.value.instanceInfo.instanceId.slice(0, 8)}": {
              "url": "${configUrl.value}"
          }
      }
  }`,
    4,
  )
})

/**
 * Handle copy config info
 */
const handleCopy = async (type: string) => {
  if (type === 'url') {
    await setClipboardData(configUrl.value)
  } else if (type === 'token') {
    await setClipboardData(configToken.value)
  } else {
    await setClipboardData(config.value)
  }
  ElMessage.success(t('action.copy'))
}

const handleTokenList = async () => {
  dialogInfo.value.instanceInfo.loading = true
  try {
    const { tokens } = await TokenAPI.list({
      instanceId: dialogInfo.value.instanceInfo.instanceId,
    })
    // reverse the token list to show the latest created token on top
    tokenList.value = (tokens || [])
      .map((token: any) => ({
        ...token,
        expire: token.expireAt !== 0 && token.expireAt < Date.now(),
      }))
      .reverse()
  } finally {
    dialogInfo.value.instanceInfo.loading = false
  }
}

const init = (instanceInfo: any) => {
  dialogInfo.value.instanceInfo = instanceInfo
  handleTokenList()
  dialogInfo.value.visible = true
}

defineExpose({
  init,
})
</script>

<style lang="scss" scoped>
.custom-style {
  position: absolute;
  bottom: 0;
  right: 0;
  .el-segmented {
    --el-segmented-item-selected-color: var(--el-text-color-primary);
    --el-border-radius-base: 16px;
  }
}
.config-info {
  font-family: 'Monaco, Menlo, "Ubuntu Mono", monospace';
  font-size: 12px;
  line-height: 1.8;
  white-space: pre;
  word-break: normal;
  border-radius: 8px;
  background: var(--ep-bg-color-deep);
  border-radius: 8px;
  box-sizing: border-box;
}

.copy-icon-url {
  position: absolute;
  top: 12px;
  right: 72px;
  cursor: pointer;
}
.copy-icon-token {
  position: absolute;
  top: 12px;
  right: 42px;
  cursor: pointer;
}
.copy-icon {
  position: absolute;
  top: 12px;
  right: 12px;
  cursor: pointer;
}
</style>
