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
      <el-row :gutter="24">
        <el-col :span="18">
          <el-form-item prop="mcpServers">
            <template #label>{{ t('mcp.instance.formData.mcpServers') }}</template>
            <MonacoEditor v-model="pageInfo.formData.mcpServers" language="json" height="200px" />
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <div
            class="pl-3 rounded border border-[var(--ep-border-color-lighter)] text-[var(--ep-text-color-secondary)] text-xs leading-6 tracking-wide"
          >
            {{ t('mcp.instance.hostingForm.directTips') }}
          </div>
        </el-col>
      </el-row>
    </el-form>
  </div>
</template>
<script setup lang="ts">
import { useTemplateFormHooks } from '../hooks/form-template.ts'
import Upload from '@/components/upload/index.vue'
import { JsonFormatter } from '@/utils/json'
import { ElLoading, ElMessage, ElMessageBox } from 'element-plus'
import { InstanceAPI } from '@/api/mcp/instance'
import { AccessType } from '@/types/instance.ts'
import { TemplateAPI } from '@/api/mcp/template'
import MonacoEditor from '@/components/MonacoEditor/index.vue'
import { type InstanceResult } from '@/types/index.ts'
import { cloneDeep } from 'lodash-es'

const { t } = useI18n()
const { query, pageInfo, placeholderServer, jumpToPage } = useTemplateFormHooks()
const baseInfo = ref()
const protocolOptions = [
  { label: 'SSE', value: 1 },
  { label: 'STREAMABLE_HTTP', value: 2 },
]
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
        if (!pageInfo.value.formData.instanceId) {
          if (Array.isArray(pageInfo.value.formData.tokens[0].headers)) {
            pageInfo.value.formData.tokens[0].headers = Object.fromEntries(
              pageInfo.value.formData.tokens[0].headers?.map((header: any) => [
                header.key,
                header.value,
              ]),
            )
          }
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
        pageInfo.value.loading = true
        await TemplateAPI.create({
          ...pageInfo.value.formData,
          environmentVariables: pageInfo.value.formData.environmentVariables?.reduce(
            (obj: any, item: any) => ({ ...obj, [item.key]: item.value }),
            {},
          ),
        })
        ElMessage.success(t('action.create'))
        pageInfo.value.loading = false
      }
    })
  } finally {
    pageInfo.value.loading = false
  }
}
const init = (instance: InstanceResult | null) => {
  if (instance) {
    pageInfo.value.formData = cloneDeep(instance)
  } else {
    pageInfo.value.formData.mcpProtocol = 1
  }
  nextTick(() => {
    pageInfo.value.formData.accessType = AccessType.DIRECT
  })
}
defineExpose({
  init,
  handleConfirm,
  handleSaveAsTemplate,
})
</script>
