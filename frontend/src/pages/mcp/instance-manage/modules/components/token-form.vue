<template>
  <el-form
    ref="formRef"
    :model="localFormData"
    label-width="auto"
    label-position="top"
    class="mx-2"
  >
    <el-form-item prop="tokenType" v-if="false">
      <!-- :label="t('mcp.token.authentication')" -->
      <div class="w-full u-line-1" style="white-space: nowrap">
        Authorization：{{ localFormData.token }}
      </div>
    </el-form-item>
    <el-form-item prop="enabledTransport" class="enabledTransport">
      <template #label>
        <div class="w-full flex justify-between items-center">
          <span class="mr-2"> {{ t('mcp.token.passthrough') }} {{ 'Headers' }} </span>
          <div class="center">
            <div
              class="cursor-pointer border border-style-solid border-rd-md border-white ml-2 p-1 center bg-gray-600 color-white hover-scale-110"
              @click="handleAddHeader"
            >
              <el-icon>
                <Plus />
              </el-icon>
            </div>
          </div>
        </div>
      </template>
      <div
        v-for="(item, index) in localFormData.headers"
        :key="index"
        class="flex items-center my-2 pr-3"
      >
        <el-row :gutter="12" class="flex-sub align-center">
          <el-col :span="7">
            <div class="flex h-full items-center justify-end">
              <el-dropdown
                v-if="tokenTypeOptions.some((tokenType) => item.key === tokenType.label)"
                trigger="click"
                class="h-full w-full flex justify-end"
                :show-arrow="false"
              >
                <div class="center cursor-pointer">
                  {{ item.key }}
                  <el-icon class="text-purple ml-1"><Sort /></el-icon>
                </div>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="handleTokenTypeChange(1, index)">
                      Authorization(Bearer)
                    </el-dropdown-item>
                    <el-dropdown-item @click="handleTokenTypeChange(2, index)">
                      Api-Key
                    </el-dropdown-item>
                    <el-dropdown-item @click="handleTokenTypeChange(3, index)">
                      X-API-key
                    </el-dropdown-item>
                    <el-dropdown-item @click="handleTokenTypeChange(4, index)">
                      Authorization(Basic)
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
              <el-input
                v-else
                v-model="item.key"
                :placeholder="t('mcp.instance.token.headersKey')"
                class="flex-sub"
              >
              </el-input>
              <span class="ml-2">:</span>
            </div>
          </el-col>
          <el-col :span="15" class="flex">
            <div class="flex">
              <el-input
                v-model="item.value"
                :placeholder="t('mcp.instance.token.headersValue')"
                class="flex-sub"
              ></el-input>
              <div
                v-if="tokenTypeOptions.some((tokenType) => item.key === tokenType.label)"
                class="text-purple cursor-pointer ml-2"
                @click="handleChangeBasic(index)"
              >
                {{ Number(item.tokenType) === 4 ? t('mcp.token.account') : '  ' }}
              </div>
            </div>
          </el-col>
          <el-col :span="2">
            <div
              class="cursor-pointer border border-style-solid delete-header border-white px-1 ml-2 center bg-red-100/50 color-white hover-bg-red-400/90 hover-scale-105"
              @click="handleDeleteHeader(index)"
            >
              <el-icon><Minus /></el-icon>
            </div>
          </el-col>
        </el-row>
      </div>
    </el-form-item>
    <el-form-item prop="usages">
      <!-- :label="t('mcp.instance.token.tag')" -->
      <el-tag
        v-for="(tag, num) in localFormData.usages"
        :key="num"
        :closable="showClosAble(tag)"
        class="mx-2 my-1"
        :disable-transitions="false"
        @close="handleCloseTag(num)"
        color="var(--ep-bg-purple-color)"
      >
        {{ tag }}
      </el-tag>
      <el-input
        v-if="showTagInput"
        ref="InputRef"
        v-model="inputValue"
        class="w-5 mx-2"
        style="width: 100px"
        size="small"
        @keyup.enter="handleTagConfirm"
        @blur="handleTagConfirm"
      />
      <el-button v-else class="mx-2" size="small" @click="showTagInput = true">
        + {{ t('mcp.instance.token.newTag') }}
      </el-button>
    </el-form-item>
  </el-form>
  <el-dialog v-model="userDataKey.visible" width="400px" top="30vh" :show-close="false">
    <el-form :model="userDataKey" class="p-4" label-width="80px">
      <el-form-item :label="t('login.username')" prop="username">
        <el-input
          v-model="userDataKey.username"
          :placeholder="t('login.message.username.required')"
        />
      </el-form-item>
      <el-form-item :label="t('login.password')" prop="password">
        <el-input
          v-model="userDataKey.password"
          type="password"
          :placeholder="t('login.message.password.required')"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="center">
        <mcp-button @click="handleConfirmAccount" class="w-25">{{ t('common.ok') }}</mcp-button>
      </div>
    </template>
  </el-dialog>
</template>
<script setup lang="ts">
import { Plus, Sort, Minus } from '@element-plus/icons-vue'
import { getToken } from '@/utils/system'
import McpButton from '@/components/mcp-button/index.vue'
import { useUserStore } from '@/stores'

const props = defineProps<{ formData: any }>()
const emit = defineEmits(['update:formData'])
const { t } = useI18n()
const { userInfo } = useUserStore()
const showTagInput = ref(false)
const inputValue = ref('')
const localFormData = computed({
  get: () => {
    return props.formData
  },
  set: (val) => {
    const emitData = { ...val }
    emit('update:formData', emitData)
  },
})
const tokenTypeOptions = [
  { label: 'Authorization', value: 1 },
  { label: 'Api-Key', value: 2 },
  { label: 'X-API-key', value: 3 },
  { label: 'Authorization', value: 4 },
]
// account password dialog data
const userDataKey = ref({
  visible: false,
  username: '',
  password: '',
  index: 0,
})
const showClosAble = computed(() => {
  return (tag: string) => {
    return [
      'user_id',
      'user_name',
      'space_id',
      'space_name',
      'intelligent_access_id',
      'intelligent_access_name',
      'intelligent_access_type',
      'default',
    ].some((keyword) => tag.includes(keyword))
      ? false
      : true
  }
})
// handle add header
const handleAddHeader = () => {
  localFormData.value.headers.push({ key: '', value: '' })
  emit('update:formData', { ...localFormData.value })
}
const handleDeleteHeader = (index: number) => {
  localFormData.value.headers.splice(index, 1)
  emit('update:formData', { ...localFormData.value })
}
// handle tag delete
const handleCloseTag = (num: number) => {
  localFormData.value.usages.splice(num, 1)
}
const handleTagConfirm = () => {
  const value = inputValue.value.trim()
  if (value && !localFormData.value.usages.includes(value)) {
    localFormData.value.usages.push(value)
  }
  inputValue.value = ''
  showTagInput.value = false
}
// handle token type change and clear token value
const handleTokenTypeChange = (tokenType: number, index: number) => {
  localFormData.value.tokenType = tokenType
  localFormData.value.headers[index].tokenType = tokenType
  localFormData.value.headers[index].key = tokenTypeOptions[tokenType - 1].label
  emit('update:formData', { ...localFormData.value })
  handleGetTokenValue(index)
}

// handle get token value
const handleGetTokenValue = (index: number) => {
  if (Number(localFormData.value.tokenType) === 1) {
    localFormData.value.headers[index].value =
      'Bearer ' +
      getToken(
        JSON.stringify({
          expireAt: localFormData.value.expireAt,
          userId: userInfo.userId,
          username: userInfo.username,
        }),
      )
  } else if (Number(localFormData.value.tokenType) === 2) {
    localFormData.value.headers[index].value = getToken(
      JSON.stringify({
        expireAt: localFormData.value.expireAt,
        userId: userInfo.userId,
        username: userInfo.username,
      }),
    )
  } else if (Number(localFormData.value.tokenType) === 3) {
    localFormData.value.headers[index].value = getToken(
      JSON.stringify({
        expireAt: localFormData.value.expireAt,
        userId: userInfo.userId,
        username: userInfo.username,
      }),
    )
  } else if (Number(localFormData.value.tokenType) === 4) {
    localFormData.value.headers[index].value =
      'Basic ' + btoa(`${userDataKey.value.username}:${userDataKey.value.password}`) // Base64 编码
  }
  emit('update:formData', { ...localFormData.value })
}

// handle show account dialog
const handleChangeBasic = (index: number) => {
  userDataKey.value.visible = !userDataKey.value.visible
  userDataKey.value.index = index
}

// handle confirm account interchange token
const handleConfirmAccount = () => {
  handleGetTokenValue(userDataKey.value.index)
  userDataKey.value.visible = false
}
// init header[0]
const init = () => {
  if (!localFormData.value.headers[0]) {
    localFormData.value.headers[0] = { key: 'Authorization', value: localFormData.value.token }
  }
}
onMounted(init)
</script>
<style lang="scss" scoped>
:deep(.el-form-item__label) {
  width: 100% !important;
}
:deep(.el-form-item__content) {
  display: block;
  width: 100%;
}
:deep(.el-tag) {
  color: var(--el-color-primary);
  border-color: var(--ep-border-color);
  .el-tag__close {
    color: var(--el-color-primary);
    &:hover {
      background-color: var(--ep-bg-purple-color-deep);
    }
  }
}
.delete-header {
  width: 24px;
  height: 24px;
  border-radius: 4px;
}
</style>
