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
            <MonacoEditor v-model="pageInfo.formData.mcpServers" language="json" height="200px" />
          </el-form-item>
          <!-- Headers 配置：紧跟 mcpServers 下方 -->
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
</template>

<script setup lang="ts">
import { useTemplateFormHooks } from '../hooks/form-template.ts'
import Upload from '@/components/upload/index.vue'
import { AccessType } from '@/types/instance.ts'
import InstanceHeaders from '../../instance-manage/modules/components/instance-headers.vue'
import { ElMessage } from 'element-plus'
import { TemplateAPI } from '@/api/mcp/template'
import { InstanceAPI } from '@/api/mcp/instance'
import MonacoEditor from '@/components/MonacoEditor/index.vue'
import { type InstanceResult } from '@/types/index.ts'
import { cloneDeep } from 'lodash-es'
import { useMcpStoreHook } from '@/stores'

const { t, locale } = useI18n()
const { query, pageInfo, jumpToPage } = useTemplateFormHooks()
const { envList } = toRefs(useMcpStoreHook())
const { handleGetEnvList } = useMcpStoreHook()
const baseInfo = ref()
const protocolOptions = [
  { label: 'SSE', value: 1 },
  { label: 'STREAMABLE_HTTP', value: 2 },
]

const mcpServersTips = computed(() => {
  return locale.value === 'en'
    ? `MCP service SSE/STREAMABLE_HTTP protocol configuration is currently in proxy mode, and the traffic will be forwarded to the MCP configuration provided through the platform gateway.
      After saving, the gateway access configuration will be displayed on the list page. You can also view <a href="#/template-manage">Template List</a> which provides multiple startup examples. `
    : `MCP服务SSE/STEAMABLE_HTTP协议配置当前为代理模式，流量会通过此平台网关转发到此配置提供的MCP 配置中，保存后会在列表页显示网关访问配置。也可以查看<a href="#/template-manage">部署模板</a> 提供了多个启动示例。`
})

/**
 * Handle McpProtocol Changed
 */
const handleMcpProtocolChange = () => {}

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
        if (!pageInfo.value.formData.instanceId) {
          if (Array.isArray(pageInfo.value.formData.tokens?.[0]?.headers)) {
            pageInfo.value.formData.tokens[0].headers = Object.fromEntries(
              pageInfo.value.formData.tokens[0].headers?.map((header: any) => [
                header.key,
                header.value,
              ]),
            )
          }
        }
        let parsedHeaders = pageInfo.value.formData.headers
        if (Array.isArray(parsedHeaders)) {
          parsedHeaders = Object.fromEntries(
            parsedHeaders.filter((header: any) => header.key?.trim()).map((header: any) => [header.key, header.value])
          )
        }
        const { instanceId } = await (
          pageInfo.value.formData.instanceId ? InstanceAPI.edit : InstanceAPI.create
        )({
          ...pageInfo.value.formData,
          headers: parsedHeaders,
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
          if (!pageInfo.value.formData.environmentId) {
            await handleGetEnvList()
            pageInfo.value.formData.environmentId = envList.value[0]?.id || ''
          }
          let parsedHeaders = pageInfo.value.formData.headers
          if (Array.isArray(parsedHeaders)) {
            parsedHeaders = Object.fromEntries(
              parsedHeaders.filter((header: any) => header.key?.trim()).map((header: any) => [header.key, header.value])
            )
          }
          await TemplateAPI.create({
            ...pageInfo.value.formData,
            headers: parsedHeaders,
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
</style>
