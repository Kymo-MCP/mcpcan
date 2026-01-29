<template>
  <div class="agent-page">
    <div class="flex justify-between page-header">
      <div class="header-container">
        {{ t('agent.pageDesc.title') }}
        <span class="desc">{{ t('agent.pageDesc.desc') }}</span>
      </div>
      <el-dropdown trigger="click" class="ml-4" @click.stop :show-arrow="false">
        <mcp-button :icon="Plus">{{ t('agent.action.create') }}</mcp-button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="community" @click="handleNewAgent(AgentType.DIFY)">
              <div class="flex align-center">
                <McpImage :src="dify" fit="contain" width="50" height="20" class="mr-1" />
                {{ t('agent.action.community') }}
              </div>
            </el-dropdown-item>
            <el-dropdown-item command="business" @click="handleNewAgent(AgentType.DIFY_ENTERPRISE)">
              <McpImage :src="dify" fit="contain" width="50" height="20" class="mr-1" />
              {{ t('agent.action.enterprise') }}
            </el-dropdown-item>
            <el-dropdown-item command="business" @click="handleNewAgent(AgentType.COZE)">
              <McpImage
                :src="coze"
                fit="contain"
                width="20"
                height="20"
                borderRadius="4"
                style="background-color: #fff; margin-right: 10px"
              />{{ t('agent.action.enterprise') }}
            </el-dropdown-item>
            <el-dropdown-item command="business" @click="handleNewAgent(AgentType.N8N)">
              <McpImage :src="n8n" fit="contain" width="20" height="20" />
              <span class="ml-2">{{ t('agent.action.n8n') }}</span>
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
    <div v-loading="pageInfo.loading" :element-loading-text="pageInfo.loadingText">
      <TablePlus
        :showOperation="true"
        searchContainer="#agentSearch"
        ref="tablePlus"
        :requestConfig="requestConfig"
        :columns="columns"
        show-view-mode
        view-mode="card"
        :showPage="false"
        v-model:pageConfig="pageConfig"
        :gridConfig="{ xs: 24, sm: 12, md: 12, lg: 8, xl: 6 }"
        :handlerColumnConfig="{
          width: '180px',
          fixed: 'right',
          align: 'center',
        }"
      >
        <template #action>
          <div class="flex justify-between mb-4">
            <div class="center">
              <el-image :src="agentLogo" style="width: 20px; height: 20px"></el-image>
              <span class="desc">{{ t('agent.pageDesc.total') }}：{{ pageConfig.total }}</span>
            </div>
            <div id="agentSearch"></div>
          </div>
        </template>
        <template #accessName="{ row }">
          <div class="flex align-center">
            <el-tooltip effect="dark" placement="top" class="flex-sub" :raw-content="true">
              <div class="flex-sub ml-2 ellipsis-two">{{ (row as any).accessName }}</div>
              <template #content>
                <div style="width: 300px">
                  {{ (row as any).accessName }}
                </div>
              </template>
            </el-tooltip>
          </div>
        </template>
        <template #operation="{ row }">
          <div class="flex align-center justify-center">
            <el-button type="text" size="small" link class="base-btn-link" @click="handleEdit(row)">
              {{ t('common.edit') }}
            </el-button>
            <el-button
              type="text"
              size="small"
              link
              class="base-btn-link"
              v-if="row.accessType !== AgentType.COZE"
              @click="handleConnection(row)"
            >
              {{ t('common.connection') }}
            </el-button>
            <el-button type="danger" size="small" link @click="handleDelete(row)">
              {{ t('common.delete') }}
            </el-button>
          </div>
        </template>
        <template #slotCard="{ row }: { row: any }">
          <el-card>
            <template #header>
              <div class="flex align-center justify-between">
                <div class="flex-sub u-line-1 font-bold text-[16px] cursor-pointer">
                  <el-tooltip :content="row.accessName" placement="top" trigger="click">
                    {{ row.accessName }}
                  </el-tooltip>
                </div>
              </div>
            </template>
            <div class="flex justify-between items-center">
              <div class="center">
                <McpImage
                  :src="logoIcon[row.accessType] || dify"
                  fit="contain"
                  width="50"
                  height="20"
                />
                <div v-if="row.accessType === AgentType.N8N" class="flex-sub ml-2 ellipsis-two">
                  {{ AgentType.N8N }}
                </div>
                <div v-else class="flex-sub ml-2 ellipsis-two">
                  {{
                    row.accessType === AgentType.DIFY
                      ? t('agent.action.community')
                      : t('agent.action.enterprise')
                  }}
                </div>
              </div>
              <div>
                <el-tooltip :content="t('common.edit')" placement="top">
                  <el-icon
                    :size="16"
                    class="mx-2 cursor-pointer link-hover"
                    @click="handleEdit(row)"
                  >
                    <Edit />
                  </el-icon>
                </el-tooltip>
                <el-tooltip
                  v-if="row.accessType !== AgentType.COZE"
                  :content="t('common.connection')"
                  placement="top"
                >
                  <el-icon
                    :size="16"
                    class="mx-2 cursor-pointer link-hover"
                    @click="handleConnection(row)"
                  >
                    <Link />
                  </el-icon>
                </el-tooltip>
                <el-tooltip :content="t('common.delete')" placement="top">
                  <el-icon
                    :size="16"
                    class="mx-2 cursor-pointer link-hover"
                    @click="handleDelete(row)"
                    color="#F56C6C"
                  >
                    <Delete />
                  </el-icon>
                </el-tooltip>
              </div>
            </div>
          </el-card>
        </template>
      </TablePlus>
    </div>

    <FormAgent ref="formAgent" @on-refresh="handleFormSuccess"></FormAgent>
  </div>
</template>

<script setup lang="ts">
import { Plus, More, Edit, Delete, Link } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import TablePlus from '@/components/TablePlus/index.vue'
import McpButton from '@/components/mcp-button/index.vue'
import McpImage from '@/components/mcp-image/index.vue'
import FormAgent from './modules/form-dialog.vue'
import { kymo, dify, coze, n8n } from '@/utils/logo.ts'
import agentLogo from '@/assets/logo/instance.png'
import { useAgentTableHooks } from './index.ts'
import { AgentType } from '@/types/agent'

const { t, tablePlus, columns, pageInfo, requestConfig, pageConfig, AgentAPI, logoIcon } =
  useAgentTableHooks()
// view model：'card' or 'table'
const formAgent = ref()

const handleNewAgent = (accessType: string) => {
  formAgent.value.init(accessType, null)
}

// form data submit success
const handleFormSuccess = () => {
  // refresh table data
  tablePlus.value?.initData()
}

const handleConnection = async (row: any) => {
  try {
    pageInfo.value.loading = true
    if (row.accessType === AgentType.N8N) {
      const { loginStatus } = await AgentAPI.checkN8n({
        accessID: row.accessID,
      })
      if (loginStatus) {
        ElMessage.success(t('agent.action.successConnection'))
      }
    } else {
      await AgentAPI.connectionTest(row)
      ElMessage.success(t('agent.action.successConnection'))
    }
  } catch (error) {
    console.error('Failed to test connection:', error)
  } finally {
    pageInfo.value.loading = false
  }
}

// handle edit formData
const handleEdit = (row: any) => {
  formAgent.value.init(row.accessType, row)
}

const handleDelete = (row: any) => {
  ElMessageBox.confirm(t('agent.pageDesc.deleteDesc'), t('common.warn'), {
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
  }).then(async () => {
    await AgentAPI.delete(row.accessID)
    ElMessage.success(t('action.delete'))
    // refresh table data
    tablePlus.value?.initData()
  })
}

/**
 * Handle init page list data
 */
const init = () => {
  tablePlus.value.initData()
}

onMounted(init)
</script>

<style scoped lang="scss">
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
</style>
