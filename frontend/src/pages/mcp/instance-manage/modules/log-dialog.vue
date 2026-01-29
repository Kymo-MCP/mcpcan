<template>
  <el-dialog
    v-model="dialogInfo.visible"
    :title="dialogInfo.title"
    width="60vw"
    top="10vh"
    :close-on-click-modal="false"
  >
    <el-scrollbar
      style="height: 70vh; padding-bottom: 20px"
      always
      class="config-info"
      v-loading="dialogInfo.loading"
    >
      <el-button-group class="action-group my-4 px-5" v-if="dialogInfo.logType === 'container'">
        <el-tooltip class="box-item" effect="dark" :content="t('common.download')" placement="top">
          <el-button :icon="Download" @click="handleDownload" />
        </el-tooltip>
        <el-tooltip class="box-item" effect="dark" :content="t('common.refresh')" placement="top">
          <el-button :icon="Refresh" @click="handleRefresh" />
        </el-tooltip>
      </el-button-group>
      <div class="custom-style my-4 px-5">
        <el-segmented v-model="dialogInfo.logType" :options="logTypeOptions" />
      </div>
      <div v-if="dialogInfo.logType === 'container'" class="px-5 py-2 log-info">
        {{ dialogInfo.logContent.logs || '暂无日志内容' }}
      </div>
      <div v-else class="px-5 py-2">
        <template v-if="dialogInfo.accessLogs.length">
          <div
            class="log-line"
            :class="{ expanded: expandedLogIds.includes(item.id) }"
            v-for="item in dialogInfo.accessLogs"
            :key="item.id"
            @click="toggleLogExpand(item.id)"
            :title="
              expandedLogIds.includes(item.id)
                ? t('mcp.instance.log.folded')
                : t('mcp.instance.log.expanded')
            "
          >
            <el-icon class="expand-icon" :class="{ rotated: expandedLogIds.includes(item.id) }">
              <CaretRight />
            </el-icon>
            <span class="log-time">{{ formatTime(item.createdAt) }}</span>
            <el-tag :type="getLevelType(item.level)" size="small" class="log-level-tag">
              {{ getLevelLabel(item.level) }}
            </el-tag>
            <span class="log-event">{{ item.event }}</span>
            <span class="log-separator">|</span>
            <span class="log-detail">{{ formatLogOneLine(item.log) }}</span>
          </div>
        </template>
        <div v-else>暂无访问日志内容</div>
      </div>
    </el-scrollbar>
  </el-dialog>
</template>

<script setup lang="ts">
import { InstanceAPI } from '@/api/mcp/instance'
import { type InstanceResult } from '@/types/instance.ts'
import { Download, Refresh, CaretRight } from '@element-plus/icons-vue'
import { downloadData } from '@/utils/files'
import { ElMessage } from 'element-plus'

const { t } = useI18n()
const logTypeOptions = [
  { label: '容器日志', value: 'container' },
  { label: '访问日志', value: 'access' },
]
const dialogInfo = ref<any>({
  visible: false,
  loading: false,
  title: '日志详情',
  logContent: '',
  accessLogs: [],
  logType: 'container',
  instanceId: '',
})
// logic for expandedLogIds
const expandedLogIds = ref<number[]>([])
/**
 * Handle refresh page
 */
const handleRefresh = () => {
  handleGetContainerLogs()
}

/**
 * Handle download logs
 */
const handleDownload = async () => {
  try {
    const { instanceName, instanceId, logs } = dialogInfo.value.logContent
    await downloadData({
      fileName: `${instanceName}_${instanceId}_logs_${new Date().toISOString().slice(0, 19).replace(/:/g, '-')}`,
      data: logs,
    })
    ElMessage.success(t('action.download'))
  } finally {
  }
}
/**
 * Handle get logs API （container logs）
 */
const handleGetContainerLogs = async () => {
  try {
    dialogInfo.value.loading = true
    const data = await InstanceAPI.logs({
      instanceId: dialogInfo.value.instanceId,
      lines: 100,
    })
    dialogInfo.value.logContent = data
  } finally {
    dialogInfo.value.loading = false
  }
}

// get logs based on filters
const handleGetAccessLogs = async () => {
  try {
    dialogInfo.value.loading = true
    const params: any = {
      instanceId: dialogInfo.value.instanceId,
      pageNum: 1,
      pageSize: 999,
    }
    const { logs: logData } = await InstanceAPI.logsByToken(params)
    dialogInfo.value.accessLogs = logData || []
  } catch (error) {
    dialogInfo.value.accessLogs = []
  } finally {
    dialogInfo.value.loading = false
  }
}
// format log as a single line
const formatLogOneLine = (log: string) => {
  if (!log) return ''
  try {
    const logObj = JSON.parse(log)
    logObj.message = JSON.parse(logObj.message.replace(/\s+/g, ' '))
    return JSON.stringify(logObj)
  } catch {
    return log
  }
}
// getLevelLabel
const getLevelLabel = (level: number): string => {
  const map: Record<number, string> = {
    1: 'Trace',
    2: 'Debug',
    3: 'Info',
    4: 'Warn',
    5: 'Error',
  }
  return map[level] || 'Unknown'
}
// getLevelType
const getLevelType = (
  level: number,
): '' | 'primary' | 'success' | 'warning' | 'danger' | 'info' => {
  const map: Record<number, '' | 'primary' | 'success' | 'warning' | 'danger' | 'info'> = {
    1: 'info',
    2: 'primary',
    3: 'success',
    4: 'warning',
    5: 'danger',
  }
  return map[level] || ''
}
const formatTime = (time: string) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
}
const toggleLogExpand = (logId: number) => {
  const index = expandedLogIds.value.indexOf(logId)
  if (index > -1) {
    expandedLogIds.value.splice(index, 1)
  } else {
    expandedLogIds.value.push(logId)
  }
}

// handle init data
const init = (instanceInfo: InstanceResult) => {
  dialogInfo.value.instanceId = instanceInfo.instanceId
  handleGetContainerLogs()
  handleGetAccessLogs()
  dialogInfo.value.visible = true
}

defineExpose({
  init,
})
</script>
<style lang="scss" scoped>
.action-group {
  position: absolute;
  top: 0;
  right: 0;
}
.custom-style {
  position: absolute;
  bottom: 0;
  right: 0;
  .el-segmented {
    --el-segmented-item-selected-color: var(--el-text-color-primary);
    --el-border-radius-base: 16px;
  }
}
.log-info {
  font-family: 'Monaco, Menlo, "Ubuntu Mono", monospace';
  font-size: 12px;
  line-height: 1.5;
  overflow: auto;
  white-space: pre-wrap;
  word-break: break-all;
  border-radius: 8px;
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
.log-line {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  border-bottom: 1px solid var(--ep-border-color-lighter);
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 12px;
  line-height: 1.8;
  transition: all 0.2s;
  cursor: pointer;

  &:hover {
    background-color: var(--el-color-primary-hover);
  }

  &:last-child {
    border-bottom: none;
  }

  // 展开状态
  &.expanded {
    align-items: flex-start;
    // background-color: var(--ep-bg-color-overlay);

    .log-detail {
      white-space: pre-wrap;
      word-break: break-all;
      overflow: visible;
      text-overflow: unset;
    }

    .expand-icon {
      transform: rotate(90deg);
    }
  }

  .expand-icon {
    flex-shrink: 0;
    color: var(--ep-text-color-secondary);
    font-size: 14px;
    transition: transform 0.2s;

    &:hover {
      color: var(--ep-color-primary);
    }
  }

  .log-time {
    flex-shrink: 0;
    color: var(--ep-text-color-secondary);
    font-size: 12px;
    min-width: 120px;
  }

  .log-level-tag {
    flex-shrink: 0;
    font-weight: 600;
  }

  .log-event {
    flex-shrink: 0;
    color: var(--ep-color-primary);
    font-weight: 500;
    min-width: 60px;
  }

  .log-separator {
    flex-shrink: 0;
    color: var(--ep-border-color);
  }

  .log-detail {
    flex: 1;
    color: var(--ep-text-color-primary);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    transition: all 0.2s;
  }
}
</style>
