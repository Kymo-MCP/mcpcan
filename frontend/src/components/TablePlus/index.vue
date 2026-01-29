<template>
  <!-- defer延迟渲染 -->
  <Teleport defer :to="props.searchContainer">
    <div class="flex">
      <el-input
        ref="searchInputRef"
        v-model="formData[formConfig[0].key]"
        v-bind="{ ...formConfig[0].props }"
        @keyup.enter="handleQuery"
        :suffix-icon="Search"
      ></el-input>
      <el-button style="width: 32px" class="ml-3" @click="initData">
        <el-icon><Refresh /></el-icon>
      </el-button>
      <el-popover
        v-model:visible="showMoreSearch"
        placement="bottom"
        :width="464"
        trigger="click"
        :show-arrow="false"
        append-to-body
        style="padding: 24px"
      >
        <template #reference>
          <el-button style="width: 32px">
            <el-icon><i class="icon iconfont MCP-shaixuan1"></i></el-icon>
          </el-button>
        </template>

        <FormPlus
          ref="searchFromRef"
          :form-config="formConfig"
          :form-data="formData"
          name="formPlus"
          label-width="108"
          label-position="left"
        >
          <template #handler>
            <el-form-item class="search-buttons flex-sub text-right">
              <div class="flex-sub flex justify-end">
                <el-button @click="resetFields" class="mr-2">{{ t('common.reseat') }}</el-button>
                <GlareHover
                  width="auto"
                  height="auto"
                  background="transparent"
                  border-color="#222222"
                  border-radius="4px"
                  glare-color="#ffffff"
                  :glare-opacity="0.3"
                  :glare-size="300"
                  :transition-duration="800"
                  :play-once="false"
                >
                  <el-button type="primary" @click="handleQuery" class="base-btn">{{
                    t('common.ok')
                  }}</el-button>
                </GlareHover>
              </div>
            </el-form-item>
          </template>
        </FormPlus>
      </el-popover>
      <el-button
        v-if="props.showViewMode"
        style="width: 32px"
        class="ml-3"
        @click="changeViewMode"
        :title="viewMode === 'table' ? '卡片视图' : '表格视图'"
      >
        <el-icon v-if="viewMode === 'table'"><Grid /></el-icon>
        <el-icon v-if="viewMode === 'card'"><List /></el-icon>
      </el-button>
    </div>
  </Teleport>

  <slot name="action"></slot>
  <el-table
    v-bind="$attrs"
    v-if="viewMode === 'table'"
    ref="dataTableRef"
    v-loading="loading"
    :data="list"
    :header-cell-style="{
      'background-color': 'var(--ep-bg-grey)',
    }"
    class="data-table__content"
    :row-key="rowKey"
    @selection-change="handleSelectionChange"
  >
    <el-table-column
      type="selection"
      width="55"
      v-if="props.multiple"
      :selectable="isRowSelectable"
      reserve-selection
    ></el-table-column>
    <el-table-column
      v-for="(column, index) in props.columns"
      :key="index || column.dataIndex"
      :label="column.label"
      :prop="column.dataIndex"
      align="left"
      min-width="100"
      v-bind="{ ...column.props }"
    >
      <template #header="scope">
        <slot :name="column.headSlot" :row="scope.row"></slot>
      </template>
      <template #default="scope">
        <!-- 优先使用插槽 -->
        <slot :name="column.dataIndex" :row="scope.row" :index="column.dataIndex">
          <!-- 其次使用 customRender -->
          <template v-if="typeof column.customRender === 'function'">
            <!-- 用 component 动态渲染 h 函数返回的 VNode -->
            <component
              :is="column.customRender({ row: scope.row, index: scope.$index })"
              v-if="isVNode(column.customRender({ row: scope.row, index: scope.$index }))"
            />
            <!-- 普通值（字符串/数字等）直接显示 -->
            <template v-else>
              {{ column.customRender({ row: scope.row, index: scope.$index }) }}
            </template>
          </template>

          <!-- 最后使用默认值 -->
          <span v-else>
            {{ scope.row[column.dataIndex] || '--' }}
          </span>
        </slot>
      </template>
    </el-table-column>
    <el-table-column
      v-if="props.showOperation"
      v-bind="props.handlerColumnConfig"
      :fixed="props.handlerColumnConfig?.fixed || 'right'"
      :width="props.handlerColumnConfig?.width || '240px'"
      :label="t('common.operation')"
    >
      <template #default="scope">
        <slot name="operation" :row="scope.row" :index="scope.$index"></slot>
      </template>
    </el-table-column>
    <template #empty>
      <el-empty :image-size="200" :description="t('status.noData')" />
    </template>
  </el-table>
  <div v-else-if="props.showViewMode && viewMode === 'card'" v-loading="loading">
    <el-row :gutter="20">
      <el-col v-for="(row, index) in list" :key="index" v-bind="props.gridConfig" class="mb-4">
        <slot name="slotCard" :row="row" :index="index">
          <el-card shadow="hover">
            <div class="text-center">{{ '数据为空' }}</div>
          </el-card>
        </slot>
      </el-col>
    </el-row>
    <div v-if="!list.length" class="mt-8">
      <el-empty :image-size="200" :description="t('status.noData')" />
    </div>
  </div>
  <div class="mt-8 flex justify-end" v-if="showPage">
    <el-pagination
      background
      :total="_pagerConfig.total"
      :current-page="_pagerConfig.page"
      :page-size="_pagerConfig.pageSize"
      @current-change="handlePageChange"
    />
  </div>
</template>

<script setup lang="ts">
import { isVNode } from 'vue'
import { cloneDeep } from 'lodash-es'
import FormPlus from '../FormPlus/index.vue'
import { Search, Refresh, List, Grid } from '@element-plus/icons-vue'
import GlareHover from '../Animation/GlareHover.vue'
import { AccessType } from '@/types/instance'

defineOptions({
  inheritAttrs: false,
})

const { t } = useI18n()
const searchInputRef = ref()
const viewMode = ref<'table' | 'card'>('table')
interface RequestConfig {
  // eslint-disable-next-line @typescript-eslint/no-unsafe-function-type
  api: Function
  searchQuery: {
    model: Record<string, any>
  }
}

interface TableColumn {
  label: string
  dataIndex: string
  customRender?: (params: { row: Record<string, any>; index: number }) => any
  [key: string]: any
}

interface PageConfig {
  page: number
  pageSize: number
  total: number
}
interface HandlerColumnConfig {
  width: string | null
  fixed: string | null
  align?: string | null
}
interface GridConfig {
  xs?: number
  sm?: number
  md?: number
  lg?: number
  xl?: number
}

const props = withDefaults(
  defineProps<{
    requestConfig: RequestConfig
    columns: TableColumn[]
    pageConfig: PageConfig
    showOperation: boolean
    searchContainer: string
    handlerColumnConfig: HandlerColumnConfig | null | undefined
    // eslint-disable-next-line @typescript-eslint/no-unsafe-function-type
    queryFormatter?: Function
    showPage?: boolean
    showViewMode?: boolean
    viewMode?: 'card' | 'table' | string
    multiple?: boolean
    rowKey?: string
    gridConfig?: GridConfig
  }>(),
  {
    showPage: () => true,
    showViewMode: () => false,
    multiple: () => false,
    rowKey: () => 'id',
    gridConfig: () => ({ xs: 24, sm: 12, md: 8, lg: 6, xl: 6 }),
  },
)

const showMoreSearch = ref(false)
const loading = ref(false)
const list = ref<unknown[]>([])
// 跨页选择：全局已选数据
const selectedRows = ref<any[]>([])
const dataTableRef =
  ref<InstanceType<(typeof import('element-plus/lib/components/table/src/table.vue'))['default']>>()
const _pagerConfig = ref(Object.assign({}, props.pageConfig))

const emit = defineEmits<{
  (e: 'update:pageConfig', value: PageConfig): void
  (e: 'resetFields', value: any): void
  (e: 'update:viewMode', value: 'card' | 'table'): void
  (e: 'on-selection-change', value: any[]): void
}>()
/**
 * Handle page change event
 */
const handlePageChange = (newPage: number) => {
  _pagerConfig.value.page = newPage
  initData()
}

/**
 * init search data config
 */
const formConfig = computed(() => {
  const config = cloneDeep(props.columns)
    .filter((f: any) => f.searchConfig)
    .map((o: any) => {
      return {
        key: o.dataIndex,
        ...o.searchConfig,
        span: o.searchConfig.span || 4,
      }
    })
  config.push({
    span: 5,
    component: 'slot',
    slotName: 'handler',
  })
  return config
})

/**
 * init form data
 */
const formData = ref<any>({})
const initFormData = () => {
  cloneDeep(props.columns)
    .filter((f: any) => f.searchConfig)
    .forEach((o: any) => {
      if (!formData.value.hasOwnProperty.call(o.dataIndex)) formData.value[o.dataIndex] = undefined
    })
}

/**
 * Handle search event
 */
const handleQuery = () => {
  showMoreSearch.value = false
  initData()
}

/**
 * reset form data
 */
const searchFromRef = ref()
const resetFields = () => {
  showMoreSearch.value = false
  initFormData()
  searchFromRef.value?.resetFields()
  emit('resetFields', null)
  initData()
}

/**
 * Handle change view mode
 */
const changeViewMode = () => {
  viewMode.value = viewMode.value === 'table' ? 'card' : 'table'
  emit('update:viewMode', viewMode.value)
  initData()
}

const handleSelectionChange = (selection: any[]) => {
  selectedRows.value = selection
  emit('on-selection-change', selectedRows.value)
}

// Element Plus selection guard: return false to make the row's checkbox disabled
const isRowSelectable = (row: any) => {
  return row?.accessType !== AccessType.DIRECT
}

//Init search data
const initSearchQuery = () => {
  if (props.queryFormatter) {
    props.queryFormatter(formData.value)
  }
}

const customize = (searchData: object) => {
  formData.value = Object.assign({ ...formData.value }, searchData)
}

/**
 * Handle init table data
 */
const initData = async () => {
  await initSearchQuery()
  try {
    loading.value = true
    const data = await props.requestConfig.api({
      page: props.pageConfig.page,
      pageSize: props.pageConfig.pageSize,
      ...formData.value,
      ...props.requestConfig.searchQuery.model,
    })

    list.value = data.list || []
    _pagerConfig.value.total = Number(data.total) || data.list?.length || 0
    emit('update:pageConfig', _pagerConfig.value)
  } catch (error) {
    console.error('数据加载失败', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  initFormData()
  // 搜索图标按钮注册搜索事件
  searchInputRef.value.$el.getElementsByClassName('el-input__suffix')[0].onclick = handleQuery
  viewMode.value = (props.viewMode as 'table' | 'card') || 'table'
})

// 提供方法：获取/设置已选项
defineExpose({
  initData,
  resetFields,
  customize,
  getSelectedRows: () => selectedRows.value,
  setSelectedRows: (rows: any[]) => {
    selectedRows.value = rows || []
  },
})
</script>

<style lang="scss">
.el-popover.el-popper {
  padding: 24px 24px 12px !important;
}
.el-table__body tr.current-row > td.el-table__cell {
  // background-color: var(--ep-bg-purple-color);
}
.el-table--striped .el-table__body tr.el-table__row--striped.current-row td.el-table__cell {
  // background-color: var(--ep-bg-purple-color);
}
.el-input__suffix {
  cursor: pointer;
}
.el-checkbox__input.is-checked .el-checkbox__inner {
  background-color: var(--el-color-primary);
  border-color: var(--ep-pager-border);
}
.el-checkbox__input.is-indeterminate .el-checkbox__inner {
  background-color: var(--el-color-primary);
  border-color: var(--ep-pager-border);
}
.el-checkbox__inner:hover {
  border-color: var(--ep-pager-border);
}
</style>
