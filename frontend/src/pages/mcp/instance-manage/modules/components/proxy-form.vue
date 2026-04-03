<template>
  <div class="px-5 py-2" v-loading="pageInfo.loading">
    <el-form
      ref="baseInfo"
      :model="pageInfo.formData"
      :rules="pageInfo.rules"
      label-width="auto"
      label-position="left"
    >
      <el-row :gutter="24">
        <el-col :span="18">
          <el-form-item :label="t('mcp.instance.formData.instanceName')" prop="name">
            <el-input
              v-model="pageInfo.formData.name"
              :placeholder="t('mcp.instance.formData.instanceName')"
            />
          </el-form-item>
          <el-form-item :label="t('mcp.template.formData.notes')" prop="notes">
            <el-input
              v-model="pageInfo.formData.notes"
              :rows="2"
              type="textarea"
              :placeholder="t('mcp.template.formData.notes')"
            />
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item prop="iconPath">
            <Upload v-model="pageInfo.formData.iconPath"></Upload>
          </el-form-item>
        </el-col>
      </el-row>
      <el-form-item :label="t('mcp.instance.form.mcpProtocol')" prop="mcpProtocol">
        <el-segmented
          v-model="pageInfo.formData.mcpProtocol"
          :options="protocolOptions"
          @change="handleMcpProtocolChange"
        />
      </el-form-item>
      <el-row>
        <el-col :span="18">
          <el-form-item prop="mcpServers">
            <template #label>{{ t('mcp.instance.formData.mcpServers') }}</template>
            <MonacoEditor
              v-model="pageInfo.formData.mcpServers"
              language="json"
              :height="locale === 'en' ? '600px' : '300px'"
            />
          </el-form-item>
          <!-- Headers 配置：紧跟 mcpServers 下方，创建/编辑均展示 -->
          <InstanceHeaders v-model:headers="pageInfo.formData.headers" />
        </el-col>
        <el-col :span="6">
          <div
            class="pl-3 rounded border border-[var(--ep-border-color-lighter)] text-[var(--ep-text-color-secondary)] text-xs leading-6 tracking-wide"
            v-html="mcpServersTips"
          ></div>
        </el-col>
      </el-row>
    </el-form>
  </div>
  <!-- probe instance dialog model -->
  <ProbeStatus ref="probe"></ProbeStatus>
  <ConfigDialog ref="config"></ConfigDialog>
  <LogDialog ref="log"></LogDialog>
</template>

<script setup lang="ts">
import { useInstanceFormHooks } from '../../hooks/form-instance.ts'
import Upload from '@/components/upload/index.vue'
import { AccessType } from '@/types/instance.ts'
import InstanceHeaders from './instance-headers.vue'
import { ElMessage } from 'element-plus'
import { TemplateAPI } from '@/api/mcp/template'
import { InstanceAPI } from '@/api/mcp/instance'
import MonacoEditor from '@/components/MonacoEditor/index.vue'
import { type InstanceResult } from '@/types/index.ts'
import { cloneDeep } from 'lodash-es'
import ProbeStatus from '../probe-dialog.vue'
import ConfigDialog from '../url-config-dialog.vue'
import LogDialog from '../log-dialog.vue'
import { useMcpStoreHook } from '@/stores'

const { t, locale } = useI18n()
const { query, pageInfo, jumpToPage } = useInstanceFormHooks()
const { currentInstance, envList } = toRefs(useMcpStoreHook())
const { handleGetEnvList } = useMcpStoreHook()
const baseInfo = ref()
const protocolOptions = [
  { label: 'SSE', value: 1 },
  { label: 'STREAMABLE_HTTP', value: 2 },
]

const mcpServersTips = computed(() => {
  return locale.value === 'en'
    ? `If you deploy by uploading a code package, set the working directory in the configuration: <code>"cwd": "/app/codepkg/[codezip_name]"</code>.
When the instance starts, the platform will automatically create a default runtime container for the current MCP configuration. The container includes the <code>mcp-hosting</code> hosting component and uses a protocol adapter to expose this configuration as SSE/Streamable HTTP. It also binds to <code>0.0.0.0:8080</code> to provide the service externally.
The MCPCAN gateway will automatically probe and detect the container's network address and complete the routing setup.
For more startup configuration examples of different types, see <a href="#/template-manage">Deploy Templates</a>.`
    : `若需上传代码包部署，需在配置信息中指定工作目录："cwd": "/app/codepkg/[codezip_name]"。
实例启动时，平台将为当前 MCP 配置自动创建默认运行容器，容器内置 mcp-hosting 托管组件，通过协议适配器将本配置转换为 SSE/Streamable HTTP 协议，同时绑定至0.0.0.0:8080端口对外提供服务；MCPCAN 网关服务将自动探测并识别该容器的网络地址，完成对接。
更多不同类型的启动配置示例，可参考<a href="#/template-manage">部署模板</a>。`
})

/**
 * Handle McpProtocol Changed
 */
const handleMcpProtocolChange = () => {}
const config = ref()
const handleConfig = () => {
  config.value.init(Object.assign(currentInstance.value, pageInfo.value.formData))
}
const probe = ref()
const handleViewStatus = () => {
  probe.value.init(pageInfo.value.formData)
}
const log = ref()
const handleViewLog = () => {
  log.value.init(pageInfo.value.formData)
}
// Handle confirm save
const handleConfirm = async () => {
  baseInfo.value.validate(async (valid: boolean) => {
    if (valid) {
      try {
        pageInfo.value.loading = true
        if (!pageInfo.value.formData.environmentId) {
          await handleGetEnvList()
          pageInfo.value.formData.environmentId = envList.value[0]?.id || ''
        }
        if (!pageInfo.value.formData.environmentId) {
          ElMessage.warning(t('mcp.template.rules.environmentId'))
          return
        }
        if (Array.isArray(pageInfo.value.formData.headers)) {
          pageInfo.value.formData.headers = Object.fromEntries(
            pageInfo.value.formData.headers
              .filter((header: any) => header.key?.trim())
              .map((header: any) => [header.key, header.value]),
          )
        }
        const { instanceId } = await (
          pageInfo.value.formData.instanceId ? InstanceAPI.edit : InstanceAPI.create
        )({
          ...pageInfo.value.formData,
          environmentVariables: pageInfo.value.formData.environmentVariables?.reduce(
            (obj: any, item: any) => ({ ...obj, [item.key]: item.value }),
            {},
          ),
        })
        ElMessage.success(
          pageInfo.value.formData.instanceId ? t('action.edit') : t('action.create'),
        )
        pageInfo.value.formData.instanceId = instanceId
        jumpToPage({
          url: '/new-instance',
          data: {
            instanceId,
            type: query.type,
          },
        })
      } finally {
        pageInfo.value.loading = false
      }
    }
  })
}
/**
 * save as a template
 */
const handleSaveAsTemplate = async () => {
  try {
    pageInfo.value.loading = true
    baseInfo.value.validate(async (valid: boolean) => {
      if (valid) {
        try {
          pageInfo.value.loading = true
          await TemplateAPI.create({
            ...pageInfo.value.formData,
            environmentVariables: pageInfo.value.formData.environmentVariables?.reduce(
              (obj: any, item: any) => ({ ...obj, [item.key]: item.value }),
              {},
            ),
          })
          ElMessage.success(t('action.create'))
        } finally {
          pageInfo.value.loading = false
        }
      }
    })
  } finally {
    pageInfo.value.loading = false
  }
}
const init = async (instance: InstanceResult | null) => {
  await handleGetEnvList()
  if (instance) {
    pageInfo.value.formData = cloneDeep(instance)
  } else {
    pageInfo.value.formData.mcpProtocol = 1
  }
  nextTick(() => {
    pageInfo.value.formData.accessType = AccessType.PROXY
    if (!pageInfo.value.formData.environmentId) {
      pageInfo.value.formData.environmentId = envList.value[0]?.id || ''
    }
  })
}
defineExpose({
  init,
  handleConfig,
  handleViewStatus,
  handleViewLog,
  handleConfirm,
  handleSaveAsTemplate,
})
</script>
<style lang="scss" scoped>
.tip {
  padding: 10px;
  border-radius: 4px;
  font-size: 12px;
  &.tip-warning {
    background-color: #fff1f0;
    border-left: 5px solid var(--el-color-danger);
  }
  &.tip-primary {
    background-color: #409eff1a;
    border-left: 5px solid var(--el-color-primary);
  }
}

/* 让 el-collapse-item 左侧图标与标题内容顶部对齐 */
:deep(.el-collapse-item__header) {
  align-items: flex-start;
}

/* header 内部容器通常也是 flex，保持一致 */
:deep(.el-collapse-item__header .el-collapse-item__title) {
  align-items: flex-start;
}

/* 图标在顶部附近，避免贴边太紧 */
:deep(.el-collapse-item__header .el-collapse-item__arrow) {
  margin-top: 16px;
}
</style>
