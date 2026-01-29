<template>
  <el-dialog
    v-model="dialogInfo.visible"
    :title="dialogInfo.title"
    :show-close="false"
    :close-on-click-modal="false"
    width="680px"
    top="20vh"
  >
    <template #header>
      <div class="flex items-center w-full">
        <span>{{ dialogInfo.title }}</span
        >:
        <McpImage
          :src="currentPlatform?.icon"
          fit="contain"
          width="80"
          height="20"
          :key="formModel.accessType"
        />
        <span class="text-purple">{{ currentPlatform?.name }}</span>
      </div>
    </template>
    <div v-loading="dialogInfo.loading" :element-loading-text="dialogInfo.loadingText">
      <el-form
        ref="formRef"
        :model="formModel"
        :rules="rules"
        :label-width="locale === 'en' ? '160px' : '110px'"
      >
        <el-form-item
          v-if="formModel.accessType === AgentType.COZE"
          :label="t('agent.formData.subType')"
          prop="subType"
        >
          <el-radio-group v-model="formModel.subType" fill="var(--el-color-primary)">
            <el-radio-button :label="t('agent.columns.person')" :value="CozeType.PERSON" />
            <el-radio-button :label="t('agent.columns.team')" :value="CozeType.TEAM" />
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="t('agent.formData.accessName')" prop="accessName">
          <el-input
            v-model="formModel.accessName"
            :placeholder="t('agent.placeholder.accessName')"
            maxlength="64"
            show-word-limit
          />
        </el-form-item>
        <el-form-item
          v-if="formModel.accessType === AgentType.COZE && formModel.subType === CozeType.TEAM"
          :label="t('agent.formData.enterpriseId')"
          prop="enterpriseID"
        >
          <el-input
            v-model="formModel.enterpriseID"
            type="text"
            :placeholder="t('agent.placeholder.enterpriseId')"
          />
        </el-form-item>
        <template
          v-if="
            formModel.accessType === AgentType.DIFY_ENTERPRISE ||
            formModel.accessType === AgentType.DIFY
          "
        >
          <el-form-item :label="t('agent.formData.dbHost')" prop="dbHost">
            <el-input v-model="formModel.dbHost" :placeholder="t('agent.placeholder.dbHost')" />
          </el-form-item>
          <el-form-item :label="t('agent.formData.dbPort')" prop="dbPort">
            <el-input-number
              v-model="formModel.dbPort"
              :min="1"
              :max="65535"
              :controls="false"
              :placeholder="t('agent.placeholder.dbPort')"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item :label="t('agent.formData.dbUser')" prop="dbUser">
            <el-input
              v-model="formModel.dbUser"
              :placeholder="t('agent.placeholder.dbUser')"
              autocomplete="off"
            />
          </el-form-item>
          <el-form-item :label="t('agent.formData.dbPassword')" prop="dbPassword">
            <el-input
              v-model="formModel.dbPassword"
              type="password"
              show-password
              :placeholder="t('agent.placeholder.dbPassword')"
              autocomplete="new-password"
            />
          </el-form-item>

          <el-form-item :label="t('agent.formData.dbName')" prop="dbName">
            <el-input v-model="formModel.dbName" :placeholder="t('agent.placeholder.dbName')" />
          </el-form-item>
        </template>
        <el-form-item
          label="UID"
          prop="cozeUserID"
          v-if="formModel.accessType === AgentType.COZE && formModel.subType === CozeType.PERSON"
        >
          <el-input
            v-model="formModel.cozeUserID"
            :placeholder="t('agent.placeholder.cozeUserID')"
            maxlength="32"
            show-word-limit
          />
        </el-form-item>
        <!-- N8N -->
        <template v-if="formModel.accessType === 'N8N'">
          <el-form-item :label="t('agent.formData.baseUrl')" prop="baseUrl">
            <el-input v-model="formModel.baseUrl" :placeholder="t('agent.placeholder.baseUrl')" />
          </el-form-item>
          <el-form-item :label="t('agent.formData.username')" prop="username">
            <el-input
              v-model="formModel.username"
              :placeholder="t('agent.placeholder.username')"
              autocomplete="off"
            />
          </el-form-item>
          <el-form-item :label="t('agent.formData.password')" prop="password">
            <el-input
              v-model="formModel.password"
              type="password"
              show-password
              :placeholder="t('agent.placeholder.password')"
              autocomplete="new-password"
            />
          </el-form-item>
        </template>
      </el-form>
    </div>

    <template #footer>
      <div class="center">
        <el-button @click="handleCancel" class="mr-2">{{ t('common.cancel') }}</el-button>
        <mcp-button @click="handleSubmit" :loading="dialogInfo.loading">{{
          t('common.save')
        }}</mcp-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import type { FormInstance, FormRules } from 'element-plus'
import McpButton from '@/components/mcp-button/index.vue'
import { AgentAPI } from '@/api/agent/index'
import { ElMessage, ElMessageBox } from 'element-plus'
import McpImage from '@/components/mcp-image/index.vue'
import { dify, coze, n8n } from '@/utils/logo.ts'
import { AgentType, CozeType } from '@/types/agent'

const { t, locale } = useI18n()
const dialogInfo = ref<{
  title: string
  visible: boolean
  loading: boolean
  loadingText: string
}>({
  title: t('agent.pageDesc.formTitle'),
  visible: false,
  loading: false,
  loadingText: t('agent.formData.loadingText'),
})

// formData model
const formModel = ref<{
  accessID: string
  accessName: string
  accessType: string
  dbHost: string
  dbPort: number
  dbUser: string
  dbPassword: string
  dbName: string
  enterpriseID?: string
  subType?: string
  baseUrl?: string
  username?: string
  password?: string
  cozeUserID?: string
}>({
  accessID: '',
  accessName: '',
  accessType: '',
  dbHost: '',
  dbPort: 5432,
  dbUser: '',
  dbPassword: '',
  dbName: '',
  enterpriseID: '',
  subType: CozeType.PERSON,
  baseUrl: '',
  username: '',
  password: '',
  cozeUserID: '',
})
const platformList = ref([
  { type: AgentType.DIFY, icon: dify, name: t('agent.action.community') },
  { type: AgentType.COZE, icon: coze, name: t('agent.action.enterprise') },
  { type: AgentType.DIFY_ENTERPRISE, icon: dify, name: t('agent.action.enterprise') },
  { type: AgentType.N8N, icon: n8n, name: t('agent.action.n8n') },
])
const currentPlatform = computed(() => {
  return platformList.value.find((item) => item.type === formModel.value.accessType)
})
const formRef = ref<FormInstance>()
// validation rules
const rules: FormRules<typeof formModel.value> = {
  accessName: [
    { required: true, message: t('agent.placeholder.accessName'), trigger: 'blur' },
    { min: 2, max: 64, message: t('agent.placeholder.accessNameLength'), trigger: 'blur' },
  ],
  subType: [{ required: true, message: t('agent.placeholder.subType'), trigger: 'change' }],
  accessType: [{ required: true, message: t('agent.placeholder.accessType'), trigger: 'change' }],
  enterpriseID: [{ required: true, message: t('agent.placeholder.enterpriseId'), trigger: 'blur' }],
  dbHost: [{ required: true, message: t('agent.placeholder.dbHost'), trigger: 'blur' }],
  dbPort: [
    {
      required: true,
      validator: (_rule, value, callback) => {
        if (!value || value < 1 || value > 65535)
          return callback(new Error(t('agent.placeholder.dbPortLength')))
        callback()
      },
      trigger: 'change',
    },
  ],
  dbUser: [{ required: true, message: t('agent.placeholder.dbUser'), trigger: 'blur' }],
  dbPassword: [{ required: true, message: t('agent.placeholder.dbPassword'), trigger: 'blur' }],
  dbName: [{ required: true, message: t('agent.placeholder.dbName'), trigger: 'blur' }],
  baseUrl: [{ required: true, message: t('agent.placeholder.baseUrl'), trigger: 'blur' }],
  username: [{ required: true, message: t('agent.placeholder.username'), trigger: 'blur' }],
  password: [{ required: true, message: t('agent.placeholder.password'), trigger: 'blur' }],
  cozeUserID: [
    { required: true, message: t('agent.placeholder.cozeUserID'), trigger: 'blur' },
    { max: 32, message: t('agent.placeholder.accessNameLength'), trigger: 'blur' },
  ],
}

// cancel dialog
const handleCancel = () => {
  dialogInfo.value.visible = false
}

// submit
const emit = defineEmits<{
  (e: 'submit', payload: typeof formModel.value): void
  (e: 'on-refresh'): void
}>()

// handle form submit
const handleSubmit = () => {
  if (formModel.value.accessType === AgentType.COZE) {
    handleSaveCoze()
  } else if (formModel.value.accessType === AgentType.N8N) {
    handleSaveN8n()
  } else {
    handleSaveDify()
  }
}

// handle save dify
const handleSaveDify = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      dialogInfo.value.loading = true
      const params = {
        accessName: formModel.value.accessName,
        accessType: formModel.value.accessType,
        dbHost: formModel.value.dbHost,
        dbPort: formModel.value.dbPort,
        dbUser: formModel.value.dbUser,
        dbPassword: formModel.value.dbPassword,
        dbName: formModel.value.dbName,
      }
      await AgentAPI.connectionTest(params)
        .then(async () => {
          await (formModel.value.accessID
            ? AgentAPI.update({ accessID: formModel.value.accessID, ...params })
            : AgentAPI.create(params))
          ElMessage.success(formModel.value.accessID ? t('action.edit') : t('action.create'))
          emit('on-refresh')
          dialogInfo.value.visible = false
          formRef.value?.resetFields()
        })
        .catch(async () => {
          const result = await ElMessageBox.confirm(
            t('agent.action.failConnectionTips'),
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
          )
          if (result === 'confirm') {
            await (formModel.value.accessID
              ? AgentAPI.update({ accessID: formModel.value.accessID, ...params })
              : AgentAPI.create(params))
            ElMessage.success(formModel.value.accessID ? t('action.edit') : t('action.create'))
            emit('on-refresh')
            dialogInfo.value.visible = false
            formRef.value?.resetFields()
          }
        })
    } catch {
    } finally {
      dialogInfo.value.loading = false
    }
  })
}

// handle save coze
const handleSaveCoze = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      dialogInfo.value.loading = true
      const params = {
        accessName: formModel.value.accessName,
        accessType: formModel.value.accessType,
        enterpriseID: formModel.value.subType === CozeType.TEAM ? formModel.value.enterpriseID : '',
        subType: formModel.value.subType,
        cozeUserID: formModel.value.cozeUserID || '',
      }
      await (formModel.value.accessID
        ? AgentAPI.update({
            accessID: formModel.value.accessID,
            ...params,
          })
        : AgentAPI.create({
            ...params,
          }))
      ElMessage.success(formModel.value.accessID ? t('action.edit') : t('action.create'))
      emit('on-refresh')
      dialogInfo.value.visible = false
      formRef.value?.resetFields()
    } finally {
      dialogInfo.value.loading = false
    }
  })
}

// handle save n8n
const handleSaveN8n = async () => {
  if (!formRef.value) return
  const validate = await formRef.value.validate()
  if (!validate) return
  try {
    dialogInfo.value.loading = true
    const { loginStatus, message, pluginStatus } = await AgentAPI.checkN8n({
      baseUrl: formModel.value.baseUrl!,
      username: formModel.value.username!,
      password: formModel.value.password!,
    })
    // don't install plugin
    if (loginStatus && !pluginStatus) {
      // confirm install plugin
      ElMessageBox.confirm(t('agent.pageDesc.boxtips'), t('common.warn'), {
        confirmButtonText: t('agent.action.authInstall'),
        cancelButtonText: t('agent.action.manualInstall'),
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
          dialogInfo.value.loading = true
          dialogInfo.value.loadingText = t('agent.formData.installingPlugin')
          const { success } = await AgentAPI.installPlugin({
            baseUrl: formModel.value.baseUrl!,
            username: formModel.value.username!,
            password: formModel.value.password!,
          })
          if (success) {
            handleSaveN8n()
          }
        })
        .catch(() => {
          handleSaveN8n()
        })
    } else if (loginStatus && pluginStatus) {
      const params = {
        accessName: formModel.value.accessName,
        accessType: formModel.value.accessType,
        baseUrl: formModel.value.baseUrl!,
        username: formModel.value.username!,
        password: formModel.value.password!,
      }
      // login and plugin success
      await (formModel.value.accessID
        ? AgentAPI.update({
            accessID: formModel.value.accessID,
            ...params,
          })
        : AgentAPI.create(params))
      ElMessage.success(formModel.value.accessID ? t('action.edit') : t('action.create'))
      emit('on-refresh')
      dialogInfo.value.visible = false
      formRef.value?.resetFields()
    }
  } finally {
    dialogInfo.value.loading = false
  }
}

/**
 * Init and open dialog
 */
const init = (accessType: string, row: any) => {
  dialogInfo.value.visible = true
  if (row) {
    formModel.value = row
  } else {
    formRef.value?.resetFields()
    formModel.value = {
      accessID: '',
      accessName: '',
      accessType,
      dbHost: '',
      dbPort: 5432,
      dbUser: '',
      dbPassword: '',
      dbName: '',
      enterpriseID: '',
      subType: CozeType.PERSON,
      cozeUserID: '',
    }
  }
}
defineExpose({
  init,
})
</script>

<style scoped>
.dialog-footer {
  display: inline-flex;
  gap: 8px;
}
</style>
