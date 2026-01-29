<template>
  <div>
    <!-- head model -->
    <div class="flex justify-between page-header">
      <div class="header-container">
        {{ t('mcp.instance.pageDesc.list') }}
        <span class="desc">{{ t('mcp.instance.pageDesc.desc') }}</span>
      </div>
      <el-dropdown
        trigger="click"
        class="ml-4"
        @click.stop
        :show-arrow="false"
        @command="(cmd: string) => handleCommand(cmd, {} as InstanceResult)"
      >
        <mcp-button :icon="Plus">{{ t('mcp.instance.action.add') }}</mcp-button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="handleAddInstance">
              <el-icon><i class="icon iconfont MCP-a-1"></i></el-icon>
              {{ t('mcp.instance.action.customize') }}
            </el-dropdown-item>
            <el-dropdown-item command="handleAddByTemplate">
              <el-icon><i class="icon iconfont MCP-a-1"></i></el-icon>
              {{ t('mcp.instance.action.byTemplate') }}
            </el-dropdown-item>
            <!-- <el-dropdown-item command="handleAddByDocs">
              <el-icon><i class="icon iconfont MCP-a-1"></i></el-icon>
              {{ t('mcp.instance.action.byDocs') }}
            </el-dropdown-item> -->
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>

    <!-- count data model -->
    <el-row v-if="!layout" justify="space-around">
      <el-col :span="6" class="center" v-for="(dataCount, index) in dataCountList" :key="index">
        <GlareHover
          width="338px"
          height="124px"
          background="transparent"
          border-color="#222222"
          border-radius="6px"
          glare-color="#ffffff"
          :glare-opacity="0.3"
          :glare-size="300"
          :transition-duration="800"
          :play-once="false"
          class="data-card"
        >
          <div
            class="data-card flex justify-around align-center"
            @click="handleSearchByCount(dataCount.type)"
          >
            <div class="center">
              <el-image
                :src="dataCount.icon"
                style="width: 20px; height: 20px"
                class="mr-2"
              ></el-image>
              {{ dataCount.title }}
            </div>
            <span class="count" v-count-to="dataCount.count"></span>
          </div>
        </GlareHover>
      </el-col>
    </el-row>

    <!-- table model -->
    <div v-loading="load.status" :element-loading-text="load.text">
      <TablePlus
        ref="tablePlus"
        searchContainer="#instanceSearch"
        :showOperation="true"
        :requestConfig="requestConfig"
        :columns="columns"
        show-view-mode
        v-model:view-mode="viewMode"
        :multiple="selection.showSelect"
        :rowKey="selection.rowKey"
        :row-class-name="tableRowClassName"
        :cell-class-name="tableRowClassName"
        v-model:pageConfig="pageConfig"
        :gridConfig="{ xs: 24, sm: 12, md: 12, lg: 8, xl: 6 }"
        :handlerColumnConfig="{
          fixed: 'right',
          width: '380px',
          align: 'center',
        }"
        @on-selection-change="handleTableSelect"
      >
        <template #action>
          <div class="flex justify-between mb-4">
            <div class="center">
              <el-image :src="instanceLogo" style="width: 20px; height: 20px"></el-image>
              <span class="desc"
                >{{ t('mcp.instance.pageDesc.total') }}：{{ pageConfig.total }}</span
              >
              <span class="ml-4 cursor-pointer base-btn-link font-bold center" @click="handleSync">
                <el-icon class="mr-1"><Share /></el-icon>
                {{ t('agent.pageDesc.platFormSync') }}

                <span v-if="selection.selectList.length" class="ml-1">
                  {{ t('agent.sync.selected') }} {{ selection.selectList.length }}
                  {{ t('agent.sync.unit2') }}
                </span>
              </span>
            </div>
            <div id="instanceSearch"></div>
          </div>
        </template>
        <template #enabledTokenHeader>
          <div class="flex align-center">
            <span class="mr-2">{{ t('mcp.instance.enabledToken') }}</span>
            <el-popover placement="top" width="250">
              <div>{{ t('mcp.instance.tableHeadDesc.token') }}</div>
              <template #reference>
                <el-icon class="cursor-pointer"><Warning /></el-icon>
              </template>
            </el-popover>
          </div>
        </template>
        <template #instanceName="{ row }">
          <div class="flex align-center">
            <mcp-image :src="row.iconPath" width="32" height="32" :key="row.instanceId"></mcp-image>
            <el-tooltip effect="dark" placement="top" class="ml-6" :raw-content="true" width="300">
              <div class="flex-sub ml-2 ellipsis-two">{{ row.instanceName }}</div>
              <template #content>
                <div class="title-instance">
                  <div>{{ row.instanceName }}</div>
                  <div class="text-primary text-bold">ID:{{ row.instanceId }}</div>
                  <div class="text-success text-bold">
                    {{ t('mcp.instance.containerName') }}：{{ row.containerName }}
                  </div>
                </div>
              </template>
            </el-tooltip>
          </div>
        </template>
        <template #enabledToken="{ row }">
          <el-switch
            v-if="row.accessType !== AccessType.DIRECT"
            v-model="row.enabledToken"
            style="--el-switch-on-color: #13ce66"
            inline-prompt
            :active-text="t('common.on')"
            :inactive-text="t('common.off')"
            :loading="row.loading"
            @change="handleEabledToken(row)"
          ></el-switch>
          <span v-else class="color-gray">{{ t('mcp.instance.pageDesc.noToken') }}</span>
        </template>
        <template #status="{ row }">
          <el-text :type="activeOptions[row.status as keyof typeof activeOptions]?.type" link>
            {{ activeOptions[row.status as keyof typeof activeOptions]?.label }}
          </el-text>
        </template>
        <template #containerStatus="{ row }">
          <div v-if="row.accessType === AccessType.HOSTING">
            <el-text
              :type="containerOptions[row.containerStatus as keyof typeof containerOptions]?.type"
              link
            >
              <span class="mr-2">
                {{ containerOptions[row.containerStatus as keyof typeof containerOptions]?.label }}
              </span>
            </el-text>
            <el-popover v-if="!row.containerIsReady" placement="right" width="250">
              <div>{{ t('mcp.instance.pageDesc.packStatusTips') }}</div>
              <template #reference>
                <el-icon class="cursor-pointer"><Warning /></el-icon>
              </template>
            </el-popover>
          </div>

          <span v-else>--</span>
        </template>
        <template #createdAt="{ row }">
          <span>{{ timestampToDate(row.createdAt) }}</span>
        </template>
        <template #publicProxyConfig="{ row }">
          <el-link
            type="success"
            class="base-btn-link"
            :underline="false"
            @click="handleViewConfig(row)"
          >
            {{ t('mcp.instance.viewConfig') }}
          </el-link>
        </template>
        <template #operation="{ row }">
          <div class="text-center">
            <el-button
              type="primary"
              size="small"
              link
              class="base-btn-link"
              @click="handleEditInstance(row)"
            >
              {{ t('env.run.action.edit') }}
            </el-button>
            <el-button
              type="primary"
              size="small"
              link
              class="base-btn-link"
              @click="handleViewAllLog(row)"
              v-if="row.accessType !== AccessType.DIRECT"
            >
              {{ t('mcp.instance.action.logs') }}
            </el-button>
            <el-button
              type="primary"
              size="small"
              link
              class="base-btn-link"
              @click="handleDebugTools(row)"
            >
              {{ t('mcp.instance.action.debugTool') }}
            </el-button>
            <el-button
              v-if="row.accessType === AccessType.HOSTING"
              type="primary"
              size="small"
              link
              class="base-btn-link"
              @click="
                row.status === InstanceStatus.INACTIVE
                  ? handleRestartInstance(row)
                  : handleStopInstance(row)
              "
            >
              {{
                row.status === InstanceStatus.INACTIVE
                  ? t('mcp.instance.action.start')
                  : t('mcp.instance.action.stop')
              }}
            </el-button>
            <el-button
              v-if="row.accessType !== AccessType.DIRECT && row.status === InstanceStatus.ACTIVE"
              type="primary"
              size="small"
              link
              class="base-btn-link"
              @click="handleRestartInstance(row)"
            >
              {{ t('mcp.instance.action.reStart') }}
            </el-button>
            <el-button
              v-if="row.status === InstanceStatus.ACTIVE"
              type="primary"
              size="small"
              link
              class="base-btn-link"
              @click="handleViewStatus(row)"
            >
              {{ t('mcp.instance.action.probe') }}
            </el-button>
            <el-button
              size="small"
              type="danger"
              link
              @click="handleDeleteInstance(row.instanceId)"
            >
              {{ t('mcp.instance.action.delete') }}
            </el-button>
          </div>
        </template>
        <template #slotCard="{ row }: { row: any }">
          <SpotlightCard
            :class-name="row.checked ? 'hover-scale-card selected-row' : 'hover-scale-card'"
            spotlight-color="var(--ep-bg-purple-color)"
          >
            <div class="h-[130px] flex flex-col">
              <div class="flex-grow-1 flex h-0">
                <div class="mr-2">
                  <mcp-image
                    :src="row.iconPath"
                    width="32"
                    height="32"
                    :key="row.instanceId"
                  ></mcp-image>
                </div>
                <div class="flex-grow-1 flex flex-col">
                  <div class="flex justify-between">
                    <div class="flex items-center flex-grow-1">
                      <div class="max-w-[110px] u-line-1 font-bold text-[16px] cursor-pointer">
                        <el-tooltip :content="row.instanceName" placement="top" trigger="click">
                          {{ row.instanceName }}
                        </el-tooltip>
                      </div>
                      <div class="ml-1 text-sm">
                        <div
                          class="flex color-[#67C23A] items-center"
                          v-if="row.accessType === AccessType.HOSTING"
                        >
                          <el-icon :size="16" class="mr-1" color="#67C23A">
                            <i class="icon iconfont MCP-anquan"></i>
                          </el-icon>
                          {{ t('mcp.type.hosting') }}
                        </div>
                        <div
                          class="flex color-[#E6A23C] items-center"
                          v-if="row.accessType === AccessType.PROXY"
                        >
                          <el-icon :size="16" class="mr-1" color="#E6A23C">
                            <i class="icon iconfont MCP-daili"></i>
                          </el-icon>
                          {{ t('mcp.type.proxy') }}
                        </div>
                        <div
                          class="flex color-[#409EFF] items-center"
                          v-if="row.accessType === AccessType.DIRECT"
                        >
                          <el-icon :size="16" class="mr-1" color="#409EFF">
                            <i class="icon iconfont MCP-zhilian"></i>
                          </el-icon>
                          {{ t('mcp.type.direct') }}
                        </div>
                      </div>
                    </div>
                    <div class="flex items-center h-[32px]">
                      <el-tooltip v-if="row.accessType === AccessType.HOSTING" placement="top">
                        <el-icon
                          :size="16"
                          class="mx-2"
                          :color="row.containerStatus === 'running' ? '#67C23A' : '#F56C6C'"
                        >
                          <i class="icon iconfont MCP-MCPshili"></i>
                        </el-icon>
                        <template #content>
                          <span>
                            {{
                              containerOptions[row.containerStatus as keyof typeof containerOptions]
                                ?.label
                            }}
                          </span>
                        </template>
                      </el-tooltip>
                      <el-tag>{{
                        mcpProtocolOptions.find((item) => item.value === row.mcpProtocol)?.label
                      }}</el-tag>
                    </div>
                  </div>
                  <div
                    class="mt-1 line-height-[20px] flex-grow-1 h-0 text-justify break-words pr-1 text-xs leading-normal ellipsis-three"
                  >
                    <el-tooltip placement="top" trigger="click">
                      {{ row.notes }}
                      <template #content>
                        <div style="width: 300px; text-align: justify; word-break: break-word">
                          {{ row.notes }}
                        </div>
                      </template>
                    </el-tooltip>
                  </div>
                </div>
              </div>
              <div class="flex justify-between mt-2">
                <div class="flex items-center">
                  <div class="ml-2">
                    <el-tooltip :content="t('mcp.instance.action.logs')" placement="top">
                      <el-icon
                        :size="16"
                        class="mx-2 cursor-pointer link-hover"
                        @click="handleViewAllLog(row)"
                      >
                        <Document />
                      </el-icon>
                    </el-tooltip>
                    <el-tooltip :content="t('mcp.instance.action.probe')" placement="top">
                      <el-icon
                        :size="16"
                        class="mx-2 cursor-pointer link-hover"
                        @click="handleViewStatus(row)"
                      >
                        <i class="icon iconfont MCP-tancerenwu"></i>
                      </el-icon>
                    </el-tooltip>
                    <el-tooltip :content="t('mcp.instance.action.debugTool')" placement="top">
                      <el-icon
                        :size="16"
                        class="mx-2 cursor-pointer link-hover"
                        @click="handleDebugTools(row)"
                      >
                        <i class="icon iconfont MCP-tool"></i>
                      </el-icon>
                    </el-tooltip>
                    <el-tooltip :content="t('env.run.action.edit')" placement="top">
                      <el-icon
                        :size="16"
                        class="mx-2 cursor-pointer link-hover"
                        @click="handleEditInstance(row)"
                      >
                        <Edit />
                      </el-icon>
                    </el-tooltip>
                    <el-tooltip :content="t('mcp.instance.action.delete')" placement="top">
                      <el-icon
                        :size="16"
                        class="mx-2 cursor-pointer link-hover"
                        @click="handleDeleteInstance(row.instanceId)"
                        color="#F56C6C"
                      >
                        <Delete />
                      </el-icon>
                    </el-tooltip>
                  </div>
                </div>
                <div class="flex">
                  <el-tooltip
                    :content="t('mcp.instance.card.containerControl')"
                    placement="top"
                    trigger="click"
                  >
                    <el-switch
                      v-if="row.accessType !== AccessType.DIRECT"
                      v-model="row.status"
                      style="--el-switch-on-color: #13ce66"
                      inline-prompt
                      size="small"
                      :active-text="t('mcp.instance.card.start')"
                      :inactive-text="t('mcp.instance.card.stop')"
                      :active-value="InstanceStatus.ACTIVE"
                      :inactive-value="InstanceStatus.INACTIVE"
                      :loading="row.loading"
                      @click="handleSwitchInstance(row)"
                    ></el-switch>
                  </el-tooltip>
                  <mcp-button size="small" class="ml-2" @click="handleViewConfig(row)">{{
                    t('mcp.instance.card.configUrl')
                  }}</mcp-button>
                </div>
              </div>
            </div>
            <el-checkbox
              v-model="row.checked"
              class="check-box"
              @change="handleSelectedWithCard(row)"
              :disabled="row.accessType === AccessType.DIRECT"
            ></el-checkbox>
          </SpotlightCard>
        </template>
      </TablePlus>
    </div>

    <!-- view detail model -->
    <InstanceDetail ref="instanceDetail"></InstanceDetail>
    <AccessTypeDialog ref="accessTypeDialog" @on-refresh="init"></AccessTypeDialog>
    <!-- view config model -->
    <ViewConfig ref="viewConfig" @on-refresh="init"></ViewConfig>
    <!-- probe instance dialog model -->
    <ProbeStatus ref="probe"></ProbeStatus>
    <!-- select template -->
    <Select
      v-model="selectVisible"
      :loading="templateLoading"
      ref="packageSelect"
      :title="t('mcp.instance.action.selectTempalte')"
      :options="templateList"
      @confirm="handleConfirmSelect"
    >
      <template #options="{ option }">
        <div class="flex justify-between">
          <div class="flex align-center w-full">
            <mcp-image :src="option.iconPath" width="32" height="32"></mcp-image>
            <el-tooltip effect="dark" placement="top" class="ml-2" :raw-content="true">
              <div class="flex-sub ml-2 ellipsis-one">{{ option.name }}</div>
              <template #content>
                <div style="width: 300px">
                  {{ option.name }}
                </div>
              </template>
            </el-tooltip>
          </div>
        </div>
      </template>
    </Select>
    <!-- Create a intance by openAPI docs -->
    <OpenAPIDialog ref="openAPIDialog" @on-refresh="init"></OpenAPIDialog>
    <!-- select agent with nameSpace  -->
    <AgentSyncDialog ref="agentSyncDialog"></AgentSyncDialog>
    <TaskList ref="taskList"></TaskList>
    <LogDialog ref="logDialog"></LogDialog>
  </div>
</template>

<script setup lang="ts">
import {
  Plus,
  More,
  Warning,
  Share,
  Document,
  Operation,
  Edit,
  Delete,
} from '@element-plus/icons-vue'
import { timestampToDate } from '@/utils/system'
import TablePlus from '@/components/TablePlus/index.vue'
import { useInstanceTableHooks } from './hooks/index.ts'
import { ElMessage, ElMessageBox } from 'element-plus'
import InstanceDetail from './modules/instance-detail.vue'
import OpenAPIDialog from './modules/open-api-dialog.vue'
import instanceLogo from '@/assets/logo/instance.png'
import GlareHover from '@/components/Animation/GlareHover.vue'
import McpButton from '@/components/mcp-button/index.vue'
import ViewConfig from './modules/view-config.vue'
import ProbeStatus from './modules/probe-dialog.vue'
import Select from '@/components/mcp-select/index.vue'
import AccessTypeDialog from './modules/access-type.vue'
import { TemplateAPI } from '@/api/mcp/template'
import McpImage from '@/components/mcp-image/index.vue'
import { AccessType, InstanceStatus, SourceType } from '@/types/instance'
import { type InstanceResult } from '@/types/instance.ts'
import AgentSyncDialog from './modules/agent-sync-dialog.vue'
import TaskList from './modules/task-list.vue'
import SpotlightCard from '@/components/Animation/SpotlightCard.vue'
import LogDialog from './modules/log-dialog.vue'

const { t } = useI18n()
const layout = useLayout()
const {
  load,
  query,
  columns,
  jumpToPage,
  tablePlus,
  requestConfig,
  pageConfig,
  activeOptions,
  containerOptions,
  InstanceAPI,
  instanceCount,
  instanceDetail,
  viewConfig,
  dataCountList,
  probe,
  openAPIDialog,
  selectVisible,
  templateLoading,
  templateList,
  timer,
  selection,
  agentSyncDialog,
  currentInstance,
  meta,
  mcpProtocolOptions,
} = useInstanceTableHooks()
const viewMode = ref('card')
const accessTypeDialog = ref()
const handleAddInstance = () => {
  accessTypeDialog.value.init()
}
watch(
  () => viewMode.value,
  () => {
    selection.value.selectList = []
  },
  { deep: true },
)
/**
 * Handle create a instance by template list
 */

const handleAddByTemplate = async () => {
  try {
    selectVisible.value = true
    templateLoading.value = true
    const data = await TemplateAPI.list({ page: 1, pageSize: 999 })
    templateList.value = data.list.map((template: any) => ({
      id: template.templateId,
      name: template.name,
      ...template,
    }))
  } finally {
    templateLoading.value = false
  }
}

/**
 * Handle add a instance by openAPI server
 */
const handleAddByDocs = () => {
  openAPIDialog.value.init()
}

// handle enabled token switch
const handleEabledToken = async (row: InstanceResult) => {
  try {
    row.loading = true
    await InstanceAPI.updateTokenStatus({
      instanceId: row.instanceId,
      enabledToken: row.enabledToken,
    })
  } catch (error) {
    row.enabledToken = !row.enabledToken
  } finally {
    row.loading = false
  }
}

/**
 * handle confirm selected template
 * @param templateId - selected of templateId
 */
const handleConfirmSelect = (templateId: string) => {
  const template = templateList.value.find((item: any) => item.templateId === templateId)
  if (
    templateList.value.find((item: any) => item.templateId === templateId).sourceType ===
    SourceType.OPENAPI
  ) {
    openAPIDialog.value.init(templateId, 'create')
    return
  }
  jumpToPage({
    url: '/new-instance',
    data: {
      templateId,
      type: template.accessType,
    },
  })
}

/**
 * Handle search instance list by count data
 * @param type - dimension of count
 */
const handleSearchByCount = (type: string) => {
  switch (type) {
    case 'total':
      tablePlus.value.resetFields()
      return
    case 'active':
      tablePlus.value.customize({
        status: InstanceStatus.ACTIVE,
        accessType: null,
        mcpProtocol: null,
      })
      break
    case 'inactive':
      tablePlus.value.customize({
        status: InstanceStatus.INACTIVE,
        accessType: null,
        mcpProtocol: null,
      })
      break
    case 'hosting':
      tablePlus.value.customize({
        status: null,
        accessType: AccessType.HOSTING,
      })
      break
    default:
      tablePlus.value.resetFields()
  }
  tablePlus.value.initData()
}

/**
 * Handle view instance detail info
 * @param row - item of instance
 */
const handleViewDetail = (row: InstanceResult) => {
  instanceDetail.value.init(row)
}

/**
 * handle view instance server log
 * @param row - item of instance
 */
const handleViewLog = (row: InstanceResult) => {
  jumpToPage({
    url: '/instance-log',
    data: {
      instanceId: row.instanceId,
    },
    // isOpen: !meta.hideLayout,
  })
}

/**
 * handle view instance access logs
 * @param row - item of instance
 */
const handleViewAccessLog = (row: InstanceResult) => {
  jumpToPage({
    url: '/token-log',
    data: {
      instanceId: row.instanceId,
    },
  })
}
const logDialog = ref()
const handleViewAllLog = (row: InstanceResult) => {
  logDialog.value.init(row)
}

/**
 * Handle eidt the instance form
 * @param row - instance form data
 */
const handleEditInstance = (row: InstanceResult) => {
  if (!row.sourceType) {
    ElMessage.error('未知类型，无法编辑')
    return
  }
  currentInstance.value = row
  if (row.sourceType === SourceType.OPENAPI) {
    openAPIDialog.value.init(row.instanceId, 'edit')
    return
  }
  // accessTypeDialog.value.init(row)
  jumpToPage({
    url: '/new-instance',
    data: {
      instanceId: row.instanceId,
      type: row.accessType,
    },
  })
}

const handleSwitchInstance = async (row: InstanceResult) => {
  try {
    row.loading = true
    if (row.status === InstanceStatus.ACTIVE) {
      await InstanceAPI.restart({
        instanceId: row.instanceId,
      })
      ElMessage.success(t('mcp.instance.action.restart'))
    } else {
      await InstanceAPI.stop({
        instanceId: row.instanceId,
      })
      ElMessage.success(t('mcp.instance.action.stopInstance'))
    }
  } catch {
    row.status =
      row.status === InstanceStatus.ACTIVE ? InstanceStatus.INACTIVE : InstanceStatus.ACTIVE
  } finally {
    row.loading = false
  }
}

/**
 * handle stop instance server
 * @param instanceId - instance id
 */
const handleStopInstance = async (instanceId: string) => {
  try {
    load.value.status = true
    load.value.text = t('mcp.instance.loading.stop')
    await InstanceAPI.stop({
      instanceId,
    })
    ElMessage.success(t('mcp.instance.action.stopInstance'))
    init()
  } finally {
    load.value.status = false
  }
}

/**
 * Handle start instance
 * @param instanceId - instance id
 */
const handleRestartInstance = async (instanceId: string) => {
  try {
    load.value.status = true
    load.value.text = t('mcp.instance.loading.restart')
    await InstanceAPI.restart({
      instanceId,
    })
    ElMessage.success(t('mcp.instance.action.restart'))
    init()
  } finally {
    load.value.status = false
  }
}

/**
 * Handle view public proxy config
 * @param instanceInfo - instance form data
 */
const handleViewConfig = async (instanceInfo: InstanceResult) => {
  viewConfig.value.init(instanceInfo)
}

/**
 * Handle probe insatnce status
 * @param instanceInfo - instance form data
 */
const handleViewStatus = async (instanceInfo: InstanceResult) => {
  probe.value.init(instanceInfo)
}

/**
 * Handle debug tools
 * @param row - item of instance data
 */
const handleDebugTools = (instanceInfo: InstanceResult) => {
  currentInstance.value = instanceInfo
  jumpToPage({
    url: '/debug-tools',
    data: {
      instanceId: instanceInfo.instanceId,
      layout: false,
    },
  })
}

/**
 * Handle delete instance
 * @param instanceId - instance id
 */
const handleDeleteInstance = (instanceId: string) => {
  ElMessageBox.confirm(t('mcp.instance.action.deleteBox'), t('common.warn'), {
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
  })
    .then(async () => {
      await InstanceAPI.delete(instanceId)
      ElMessage({
        type: 'success',
        message: t('action.delete'),
      })
      init()
    })
    .catch(() => {})
}

/**
 * Collect more opearation
 * @param callback - function name
 * @param row - item of instance data
 */
const handleCommand = (callback: string, row: InstanceResult) => {
  switch (callback) {
    case 'handleViewLog':
      handleViewLog(row)
      break
    case 'handleEditInstance':
      handleEditInstance(row)
      break
    case 'handleStopInstance':
      handleStopInstance(row.instanceId)
      break
    case 'handleRestartInstance':
      handleRestartInstance(row.instanceId)
      break
    case 'handleViewStatus':
      handleViewStatus(row)
      break
    case 'handleDebugTools':
      handleDebugTools(row)
      break
    case 'handleDeleteInstance':
      handleDeleteInstance(row.instanceId)
      break
    case 'handleAddInstance':
      handleAddInstance()
      break
    case 'handleAddByTemplate':
      handleAddByTemplate()
      break
    case 'handleAddByDocs':
      handleAddByDocs()
      break
    case 'handleViewAccessLog':
      handleViewAccessLog(row)
      break
    default:
      ElMessage.warning(`未找到 "${callback}" 对应的操作`)
  }
}
/**
 * Handle table select change
 * @param selectionList - selected instance list
 */
const handleSync = () => {
  if (!selection.value.selectList.length) {
    ElMessage.warning(t('agent.pageDesc.mustSelectMCP'))
    return
  }
  agentSyncDialog.value.init(selection.value.selectList)
}
const handleTableSelect = (selectionList: InstanceResult[]) => {
  selection.value.selectList = selectionList
  // 有选择项暂停定时任务；否则启动定时任务
  if (!selection.value.selectList.length) {
    init()
    if (timer.value) {
      return
    } else {
      timer.value = setInterval(init, 30000)
    }
  } else {
    clearInterval(timer.value)
    timer.value = 0
  }
}

const tableRowClassName = ({ row }: { row: any }) => {
  if (selection.value.selectList.find((item) => item.instanceId === row.instanceId)) {
    return 'selected-row'
  }
  return ''
}

// selected with card mode
const handleSelectedWithCard = (row: InstanceResult) => {
  const index = selection.value.selectList.findIndex((item) => item.instanceId === row.instanceId)
  if (index === -1) {
    selection.value.selectList.push(row)
  } else {
    selection.value.selectList.splice(index, 1)
  }
  // 有选择项暂停定时任务；否则启动定时任务
  if (!selection.value.selectList.length) {
    init()
    if (timer.value) {
      return
    } else {
      timer.value = setInterval(init, 30000)
    }
  } else {
    clearInterval(timer.value)
    timer.value = 0
  }
}

/**
 * Handle get count data
 */
const handleGetCount = async () => {
  const data = await InstanceAPI.count()
  instanceCount.value = data
}

/**
 * Init data for page
 */
const init = () => {
  if (query.type) {
    handleSearchByCount(query.type as string)
  } else {
    tablePlus.value.initData()
  }
  handleGetCount()
}

onBeforeUnmount(() => {
  if (timer.value) {
    clearInterval(timer.value)
  }
})

onActivated(() => {
  init()
  if (!timer.value) {
    timer.value = setInterval(init, 30000)
  }
})

onDeactivated(() => {
  if (timer.value) {
    clearInterval(timer.value)
    timer.value = 0
  }
})

onMounted(() => {
  init()
  if (timer.value) {
    return
  } else {
    timer.value = setInterval(init, 30000)
  }
})
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
.data-card {
  width: 338px;
  height: 124px;
  background: url('@/assets/images/data-card.png') no-repeat center center;
  background-size: 100% 100%;
  margin-bottom: 24px;
  cursor: pointer;
  color: white;
  .count {
    font-size: 40px;
  }
}

.danger-btn {
  width: 100px;
  border-radius: 4px;
  &.el-button {
    background-color: rgba(12, 25, 207, 0.08) !important;
  }
}
.option-item {
  width: 580px;
}
.title-instance {
  width: 300px;
}
:deep(.el-table) .selected-row {
  --el-table-tr-bg-color: var(--ep-bg-purple-color-deep);
}
.hover-scale-card {
  transition: transform 0.3s;
  border: 1px solid var(--el-color-primary);
  border-radius: 8px;
  position: relative;
  padding: 16px;
  &:hover {
    transform: scale(1.02);
  }
  .check-box {
    position: absolute;
    top: 50%;
    left: 20px;
    transform: translateY(-50%);
  }
}
.selected-row {
  background-color: var(--ep-bg-purple-color-deep);
}
</style>
