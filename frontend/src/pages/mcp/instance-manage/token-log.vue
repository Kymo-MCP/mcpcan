<template>
  <div class="token-log-page">
    <div class="mb-4">
      <el-link v-if="layout" link @click="handleBack" class="link-hover" underline="never">
        <el-icon class="mr-2">
          <i class="icon iconfont MCP-fanhui"></i>
        </el-icon>
        {{ t('common.back') }}
      </el-link>
    </div>
    <!-- search bar -->
    <div class="filter-bar mb-4 p-3 border-rd-2">
      <!-- level -->
      <div class="grid grid-cols-12">
        <div class="col-span-6 flex items-center gap-3 mb-3">
          <span class="font-bold filter-label">{{ t('mcp.instance.log.level') }}：</span>
          <div class="level-tags flex items-center gap-2">
            <el-tag
              v-for="item in levelOptions"
              :key="item.value"
              :type="item.type"
              :effect="level === item.value ? 'dark' : 'light'"
              size="small"
              class="level-tag"
              @click="handleLevelChange(item.value)"
            >
              {{ item.label }}
            </el-tag>
          </div>
        </div>
        <div class="col-span-6 flex items-center">
          <span class="color-gray font-size-3">
            {{ t('mcp.instance.log.timeTips') }}. {{ last24hText }}
          </span>
          <!-- <el-popover placement="top" width="300">
              <div>{{ t('mcp.instance.log.timeTips') }}</div>
              <template #reference>
                <el-icon class="cursor-pointer"><Warning /></el-icon>
              </template>
            </el-popover> -->
        </div>
      </div>

      <div class="grid grid-cols-12">
        <!-- instance list -->
        <div class="col-span-6 flex items-center gap-3 mb-3">
          <span class="font-bold filter-label">{{ t('mcp.instance.log.instance') }}：</span>
          <el-select
            v-model="selectedInstanceId"
            :placeholder="t('mcp.instance.log.instance')"
            size="small"
            clearable
            filterable
            @change="handleInstanceChange"
            style="width: 400px"
          >
            <el-option
              v-for="instance in instanceList"
              :key="instance.instanceId"
              :label="instance.name"
              :value="instance.instanceId"
            />
          </el-select>
        </div>

        <!-- time range -->
        <div class="col-span-6 flex items-center justify-between gap-3 mb-3">
          <div class="flex items-center gap-3">
            <span class="font-bold filter-label">{{ t('mcp.instance.log.timeRange') }}： </span>
            <div>
              <el-date-picker
                v-model="dateRange"
                type="datetimerange"
                :range-separator="t('mcp.instance.log.to')"
                :start-placeholder="t('mcp.instance.log.startTime')"
                :end-placeholder="t('mcp.instance.log.endTime')"
                size="small"
                class="w-full"
                @change="handleDateRangeChange"
              />
              <!-- :disabled-date="disabledDate" -->
            </div>
          </div>
          <el-button size="small" @click="handleGetLogs" :loading="loading">
            <el-icon><Refresh /></el-icon>
            {{ t('common.refresh') }}
          </el-button>
        </div>
      </div>

      <!-- Trace ID + refresh button -->
      <div class="grid grid-cols-12">
        <div class="col-span-6 flex items-center gap-3 mb-3">
          <span class="font-bold filter-label">Trace ID：</span>
          <el-input
            v-model="traceId"
            placeholder="请输入 Trace ID"
            size="small"
            clearable
            @clear="handleGetLogs"
            @keyup.enter="handleGetLogs"
            style="width: 400px"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>
        <!-- Token list -->
        <div class="col-span-6 flex items-center gap-3 mb-3">
          <span class="font-bold filter-label">Token：</span>
          <el-select
            v-model="selectedToken"
            :placeholder="t('mcp.instance.log.token')"
            size="small"
            clearable
            filterable
            @change="handleGetLogs"
            style="width: 400px"
            :disabled="!selectedInstanceId"
          >
            <el-option
              v-for="(tokenItem, index) in tokenList"
              :key="index"
              :label="`${tokenItem.token.substring(0, 30)}... (${tokenItem.usages.join(', ')})`"
              :value="tokenItem.token"
            />
          </el-select>
        </div>
      </div>
    </div>

    <!-- logs content -->
    <div class="logs-container" v-loading="loading">
      <el-scrollbar class="logs-box" v-if="logList.length > 0">
        <div
          class="log-line"
          :class="{ expanded: expandedLogIds.includes(item.id) }"
          v-for="item in logList"
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
      </el-scrollbar>
      <el-empty v-else :description="t('mcp.instance.log.empty')" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { InstanceAPI } from '@/api/mcp/instance'
import { Refresh, Search, CaretRight, Warning } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useRouterHooks } from '@/utils/url'

interface LogItem {
  id: number
  instanceId: string
  tokenHeader: string
  token: string
  usages: string[]
  level: number
  event: string
  createdAt: string
  log: string
  traceId: string
}

interface InstanceItem {
  instanceId: string
  name: string
}

interface TokenItem {
  token: string
  usages: string[]
}

const route = useRoute()
const layout = useLayout()
const logList = ref<LogItem[]>([])
const loading = ref(false)
const dateRange = ref<[Date, Date] | null>(null)
const level = ref<string | number>('')
const traceId = ref('')
const { t } = useI18n()
const { jumpBack } = useRouterHooks()

const instanceList = ref<InstanceItem[]>([])
const tokenList = ref<TokenItem[]>([])
const selectedInstanceId = ref<string>('')
const selectedToken = ref<string>('')

// logic for expandedLogIds
const expandedLogIds = ref<number[]>([])

// level options
const levelOptions = [
  { label: t('mcp.instance.log.all'), value: '', type: 'primary' as const },
  { label: 'Trace', value: 1, type: 'primary' as const },
  { label: 'Debug', value: 2, type: 'info' as const },
  { label: 'Info', value: 3, type: 'primary' as const },
  { label: 'Warn', value: 4, type: 'warning' as const },
  { label: 'Error', value: 5, type: 'danger' as const },
]

// get instanceId and token by route params
const instanceId = computed(
  () => selectedInstanceId.value || (route.query.instanceId as string) || '',
)
const token = computed(() => selectedToken.value || (route.query.token as string) || '')

// 最近24小时时间段展示
const last24hText = computed(() => {
  const end = new Date()
  const start = new Date(end.getTime() - 24 * 60 * 60 * 1000)
  const fmt = (d: Date) =>
    d.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    })
  return `${fmt(start)} ~ ${fmt(end)}`
})

// level change
const handleLevelChange = (val: string | number) => {
  level.value = val
  handleGetLogs()
}

// handle log expand
const toggleLogExpand = (logId: number) => {
  const index = expandedLogIds.value.indexOf(logId)
  if (index > -1) {
    expandedLogIds.value.splice(index, 1)
  } else {
    expandedLogIds.value.push(logId)
  }
}

// const disabledDate = (date: Date) => {
//   return date.getTime() > Date.now() || date.getTime() < Date.now() - 24 * 60 * 60 * 1000
// }

const handleDateRangeChange = (value: [Date, Date] | null) => {
  if (!value || value.length !== 2) {
    handleGetLogs()
    return
  }

  const [startDate, endDate] = value
  const timeDiff = endDate.getTime() - startDate.getTime()
  const maxDiff = 24 * 60 * 60 * 1000 // 24小时的毫秒数

  if (timeDiff > maxDiff) {
    ElMessage.warning(t('mcp.instance.log.timeRangeError'))
    // 自动调整结束时间为开始时间 + 24小时
    const newEndDate = new Date(startDate.getTime() + maxDiff)
    dateRange.value = [startDate, newEndDate]
  }

  handleGetLogs()
}

// get instance list
const getInstanceList = async () => {
  try {
    const { list } = await InstanceAPI.list({ page: '1', pageSize: '999' })
    instanceList.value = (list || []).map((item: any) => ({
      instanceId: item.instanceId,
      name: item.instanceName || item.instanceId,
    }))
  } catch (error) {
    console.error(t('mcp.instance.log.instanceError'), error)
    ElMessage.error(t('mcp.instance.log.instanceError'))
  }
}

// handle get token list
const getTokenList = async (instanceId: string) => {
  if (!instanceId) {
    tokenList.value = []
    return
  }
  try {
    const data = await InstanceAPI.detail({ instanceId })
    tokenList.value = (data.tokens || []).map((item: any) => ({
      token: item.token,
      usages: item.usages || [],
    }))
  } catch (error) {
    console.error(t('mcp.instance.log.tokenFail'), error)
    ElMessage.error(t('mcp.instance.log.tokenFail'))
    tokenList.value = []
  }
}

// change instance
const handleInstanceChange = async (val: string) => {
  selectedToken.value = ''
  tokenList.value = []
  logList.value = []
  if (val) {
    handleGetLogs()
    await getTokenList(val)
  }
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

// get logs based on filters
const handleGetLogs = async () => {
  if (!instanceId.value) {
    ElMessage.warning(t('mcp.instance.log.getInstance'))
    return
  }

  loading.value = true
  try {
    const params: any = {
      instanceId: instanceId.value,
      token: token.value,
      startTime: dateRange.value ? new Date(dateRange.value[0]).getTime().toString() : '',
      endTime: dateRange.value ? new Date(dateRange.value[1]).getTime().toString() : '',
      level: level.value,
      traceId: traceId.value,
      pageNum: 1,
      pageSize: 999,
    }

    const { logs: logData } = await InstanceAPI.logsByToken(params)
    logList.value = logData || []
  } catch (error) {
    console.error(t('mcp.instance.log.logFail'), error)
    ElMessage.error(t('mcp.instance.log.logFail'))
    logList.value = []
  } finally {
    loading.value = false
  }
}

// back last class page
const handleBack = () => {
  jumpBack()
}

const init = async () => {
  // get instance list
  await getInstanceList()
  // initialize selected instanceId and token if present in route
  selectedInstanceId.value = route.query.instanceId as string
  selectedToken.value = route.query.token as string
  if (route.query.instanceId) {
    await getTokenList(selectedInstanceId.value)
  }
  handleGetLogs()
}

onMounted(init)
</script>

<style lang="scss" scoped>
.token-log-page {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.filter-bar {
  // background: var(--ep-bg-color);
  box-shadow: 0 0px 4px rgba(0, 0, 0, 0.08);

  .filter-label {
    display: inline-block;
    min-width: 80px;
    text-align: right;
    color: var(--ep-text-color-regular);
  }

  .color-gray {
    color: var(--ep-text-color-secondary);
  }
}

.level-tags {
  .level-tag {
    cursor: pointer;
    transition: all 0.3s;
    user-select: none;
    font-weight: 500;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
    }
  }
}

.logs-container {
  flex: 1;
  height: 0;
  background: var(--ep-bg-color-page);
  border-radius: 8px;
  overflow: hidden;
}

.logs-box {
  height: 100%;
  padding: 12px;
  box-shadow: 0 0px 4px rgba(0, 0, 0, 0.8);
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
    min-width: 160px;
  }

  .log-level-tag {
    flex-shrink: 0;
    font-weight: 600;
  }

  .log-event {
    flex-shrink: 0;
    color: var(--ep-color-primary);
    font-weight: 500;
    min-width: 200px;
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
