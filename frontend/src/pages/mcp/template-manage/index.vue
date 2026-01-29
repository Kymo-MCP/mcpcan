<template>
  <div>
    <!-- head model -->
    <div class="flex justify-between page-header">
      <div class="header-container">
        {{ t('mcp.template.pageDesc.list') }}
        <span class="desc">{{ t('mcp.template.pageDesc.desc') }}</span>
      </div>

      <mcp-button :icon="Plus" @click="handleAddTemplate">{{
        t('mcp.template.action.add')
      }}</mcp-button>
    </div>

    <!-- table model -->
    <TablePlus
      :showOperation="true"
      searchContainer="#templateSearch"
      ref="tablePlus"
      :requestConfig="requestConfig"
      :columns="columns"
      v-model:pageConfig="pageConfig"
      :handlerColumnConfig="{
        width: '220px',
        fixed: 'right',
        align: 'center',
      }"
    >
      <template #action>
        <div class="flex justify-between mb-4">
          <div class="center">
            <el-image :src="instanceLogo" style="width: 20px; height: 20px"></el-image>
            <span class="desc">{{ t('mcp.template.pageDesc.total') }}：{{ pageConfig.total }}</span>
          </div>
          <div id="templateSearch"></div>
        </div>
      </template>
      <template #name="{ row }">
        <div class="flex align-center">
          <mcp-image
            :src="baseUrl + row.iconPath"
            width="32"
            height="32"
            :key="row.templateId"
          ></mcp-image>
          <el-tooltip effect="dark" placement="top" class="ml-6" :raw-content="true">
            <div class="flex-sub ml-2 ellipsis-two">{{ row.name }}</div>
            <template #content>
              <div style="width: 300px">
                {{ row.name }}
              </div>
            </template>
          </el-tooltip>
        </div>
      </template>
      <template #operation="{ row }">
        <el-button size="small" link @click="handleEditTemplate(row)" class="base-btn-link">
          {{ t('env.run.action.edit') }}
        </el-button>
        <el-button size="small" link @click="handleCreatInstance(row)" class="base-btn-link">
          {{ t('mcp.template.action.createInstance') }}
        </el-button>
        <el-button type="danger" link @click="handleDeleteTemplate(row)">
          {{ t('mcp.instance.action.delete') }}
        </el-button>
      </template>
    </TablePlus>
    <AccessTypeDialog ref="accessTypeDialog"></AccessTypeDialog>
    <!-- Create a intance by openAPI docs -->
    <OpenAPIDialog ref="openAPIDialog" @on-refresh="init"></OpenAPIDialog>
  </div>
</template>

<script setup lang="ts">
import { Plus, More } from '@element-plus/icons-vue'
import TablePlus from '@/components/TablePlus/index.vue'
import { useTemplateTableHooks } from './hooks/index.ts'
import { ElMessage, ElMessageBox } from 'element-plus'
import { TemplateAPI } from '@/api/mcp/template'
import { useRouterHooks } from '@/utils/url'
import instanceLogo from '@/assets/logo/instance.png'
import McpButton from '@/components/mcp-button/index.vue'
import McpImage from '@/components/mcp-image/index.vue'
import { SourceType, type TemplateResult } from '@/types/index.ts'
import OpenAPIDialog from '../instance-manage/modules/open-api-dialog.vue'
import AccessTypeDialog from '../instance-manage/modules/access-type.vue'

const { t } = useI18n()
const { tablePlus, columns, requestConfig, pageConfig } = useTemplateTableHooks()
const { jumpToPage } = useRouterHooks()
const openAPIDialog = ref()
const baseUrl = (window as any).__APP_CONFIG__?.PUBLIC_PATH || ''
const accessTypeDialog = ref()

/**
 * Handle create a tamplate
 */
const handleAddTemplate = () => {
  accessTypeDialog.value.init()
  // jumpToPage({
  //   url: '/new-template',
  //   data: {},
  // })
}

/**
 * Handle edit the template
 * @param row - item of template
 */
const handleEditTemplate = (row: TemplateResult) => {
  if (row.sourceType === SourceType.OPENAPI) {
    openAPIDialog.value.init(row.templateId, 'template')
    return
  }
  jumpToPage({
    url: '/new-template',
    data: {
      templateId: row.templateId,
      type: row.accessType,
    },
  })
}

/**
 * Handle delete a template
 * @param templateId - template key
 */
const handleDeleteTemplate = (row: TemplateResult) => {
  ElMessageBox.confirm(
    `${t('common.confirm') + t('status.delete')}“${row.name}”?`,
    t('common.warn'),
    {
      confirmButtonText: t('common.ok'),
      cancelButtonText: t('common.cancel'),
      type: 'warning',
      customClass: 'tips-box',
      center: true,
      showClose: false,
      confirmButtonClass: 'is-plain el-button--danger danger-btn',
      customStyle: {
        width: '517px',
        height: '247px',
      },
    },
  ).then(async () => {
    await TemplateAPI.delete(row.templateId)
    ElMessage({
      type: 'success',
      message: t('status.delete') + t('status.success'),
    })
    init()
  })
}

/**
 * Handle create a instance by template
 */
const handleCreatInstance = (row: TemplateResult) => {
  if (row.sourceType === SourceType.OPENAPI) {
    openAPIDialog.value.init(row.templateId, 'instance')
    return
  }
  jumpToPage({
    url: '/new-instance',
    data: { templateId: row.templateId, type: row.accessType },
  })
}

/**
 * Handle more operation function
 * @param callback - function name
 * @param row - item of template
 */
const handleCommand = (callback: string, row: TemplateResult) => {
  switch (callback) {
    case 'handleCreatInstance':
      handleCreatInstance(row)
      break
    case 'handleDeleteTemplate':
      handleDeleteTemplate(row)
      break
    default:
      ElMessage.warning(`未找到 "${callback}" 对应的操作`)
  }
}

/**
 * Handle init template list
 */
const init = () => {
  tablePlus.value.initData()
}

onMounted(init)
</script>
<style lang="scss" scoped>
.page-header {
  margin-bottom: 24px;
  .header-container {
    font-size: 20px;
  }
}
.desc {
  font-size: 16px;
  color: #999999;
  margin-left: 16px;
}
:deep(.danger-btn) {
  width: 100px;
  border-radius: 4px;
  &.el-button {
    background-color: rgba(12, 25, 207, 0.08) !important;
  }
}
</style>
