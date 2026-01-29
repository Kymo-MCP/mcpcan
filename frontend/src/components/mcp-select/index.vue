<template>
  <el-dialog
    v-model:value="currentVisible"
    v-bind="props"
    width="680px"
    top="6vh"
    :show-close="false"
  >
    <template #header>
      <div class="center">{{ props.title }}</div>
    </template>
    <div v-loading="props.loading">
      <el-scrollbar ref="scrollbarRef" max-height="75vh" class="pr-2">
        <div class="mb-4 mt-1 flex align-center w-full pl-2">
          <el-input
            v-model="searchKeyword"
            :suffix-icon="Search"
            :placeholder="t('desc.placeholderKey')"
            class="flex-sub mr-4"
          ></el-input>
          <slot name="action"></slot>
        </div>

        <el-radio-group v-model="currentSelected" class="mr-2 ml-2" v-if="_options.length">
          <el-radio
            v-for="(option, index) in _options"
            :key="index"
            :value="option.id"
            size="large"
            class="w-full radio-item"
          >
            <slot name="options" :option="option">
              <div class="flex justify-between">
                <div class="flex align-center">
                  <el-image :src="zipLogo" style="width: 32px; height: 32px"></el-image>
                  <span class="ml-2"> {{ option.name }}</span>
                </div>
                <div class="flex align-center">
                  {{ t('mcp.template.formData.size') }}：{{ formatFileSize(option.size) }}
                  <span>{{ timestampToDate(option.createdAt) }}</span>
                </div>
              </div>
            </slot>
          </el-radio>
        </el-radio-group>
        <el-empty v-else></el-empty>
      </el-scrollbar>
    </div>
    <template #footer>
      <div class="center">
        <el-button @click="handleCancel" class="mr-4">{{ t('common.cancel') }}</el-button>
        <McpButton @click="handleComfirm">{{ t('common.ok') }}</McpButton>
      </div>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import zipLogo from '@/assets/logo/zip.png'
import { Search } from '@element-plus/icons-vue'
import { formatFileSize, timestampToDate } from '@/utils/system'
import McpButton from '../mcp-button/index.vue'

const { t } = useI18n()

// interface OptionItem {
//   id: string | number
//   name: string
//   [key: string]: any
// }
// const props = withDefaults(
//   defineProps<{
//     modelValue?: boolean
//     modelSelected?: string | number | null
//     title: string
//     options: Array<any>
//   }>(),
//   {
//     modelValue: false,
//     modelSelected: '',
//   },
// )

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  selected: {
    type: String,
    default: '',
  },
  title: {
    type: String,
    default: '',
  },
  options: {
    type: Array<any>,
    default: [],
  },
  loading: {
    type: Boolean,
    default: false,
  },
})

const currentVisible = ref(props.modelValue)
const currentSelected = ref(props.selected)
const searchKeyword = ref('')
const _options = computed(() => {
  const escapedKeyword = searchKeyword.value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  const regex = new RegExp(escapedKeyword, 'i')
  return props.options.filter((item: any) => item.name && regex.test(item.name))
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'update:selected', value: string): void
  (e: 'confirm', value: string): void
}>()

watch(
  () => props.modelValue,
  (newVal) => {
    currentVisible.value = newVal as boolean
  },
)

watch(
  () => currentVisible.value,
  (newVal) => {
    emit('update:modelValue', newVal)
  },
)
watch(
  () => props.selected,
  (newVal) => {
    currentSelected.value = newVal
  },
)

const handleCancel = () => {
  currentVisible.value = false
  currentSelected.value = props.selected
}

const handleComfirm = () => {
  emit('update:selected', currentSelected.value as string)
  emit('confirm', currentSelected.value as string)
  currentVisible.value = false
}
</script>

<style lang="scss" scoped>
.w100 {
  width: 100;
}
.radio-item {
  height: 56px !important;
  background: var(--ep-bg-stripe-dark);

  &:nth-child(2n-1) {
    background: var(--ep-bg-stripe-light);
  }
  &.is-checked {
    background-color: rgba(204, 187, 255, 0.16);
  }
}
:deep(.el-radio-group) {
  width: calc(100% - 20px);
}
:deep(.el-radio) {
  padding: 16px;
  margin-right: 0;
}
:deep(.el-radio__label) {
  width: 100%;
  padding-right: 20px;
  flex: 1;
}
:deep(.el-radio__inner) {
  border: 1px solid var(--el-color-primary);
}
:deep(.el-radio__input.is-checked .el-radio__inner) {
  background: var(--el-color-primary);
  border-color: var(--el-color-primary);
}
:deep(.el-radio__input.is-checked + .el-radio__label) {
  color: var(--el-color-primary);
}
:deep(.el-upload-dragger) {
  height: 100%;
  border: 1px dashed var(--el-color-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 6px;
  &:hover {
    border-color: var(--el-color-primary-hover);
    .el-icon--upload {
      color: var(--el-color-primary-hover);
    }
    .el-upload__text {
      color: var(--el-color-primary-hover);
    }
  }
}
</style>
