<template>
  <el-dialog
    v-model="dialogInfo.visible"
    width="1200px"
    align-center
    :show-close="false"
    header-class="header-border"
    footer-class="footer-border"
  >
    <template #title>{{ dialogInfo.title }}</template>
    <el-splitter v-loading="dialogInfo.loading">
      <el-splitter-panel size="50%" class="p-4">
        <div style="height: 75vh">
          <el-form
            ref="baseInfo"
            :model="formData"
            :rules="rules"
            label-width="auto"
            label-position="top"
          >
            <el-form-item :label="t('mcp.instance.formData.instanceName')" prop="name">
              <el-input
                v-model="formData.name"
                :placeholder="t('mcp.instance.formData.instanceName')"
              />
            </el-form-item>
            <el-form-item
              v-if="false"
              :label="t('mcp.instance.formData.environmentId')"
              prop="environmentId"
            >
              <el-select
                v-model="formData.environmentId"
                :placeholder="t('mcp.instance.formData.environmentId')"
              >
                <el-option
                  v-for="(env, index) in envList"
                  :key="index"
                  :label="env.name"
                  :value="env.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item :label="t('mcp.instance.openApi.serviceUrl')" prop="openapiBaseUrl">
              <el-input
                v-model="formData.openapiBaseUrl"
                :placeholder="t('mcp.instance.openApi.serviceUrl')"
              />
            </el-form-item>
            <el-form-item :label="t('mcp.template.formData.notes')" prop="notes">
              <el-input
                v-model="formData.notes"
                :rows="4"
                type="textarea"
                :placeholder="t('mcp.template.formData.notes')"
              />
            </el-form-item>
            <el-form-item :label="t('mcp.template.formData.icon')" prop="iconPath">
              <Upload v-model="formData.iconPath"></Upload>
            </el-form-item>
          </el-form>

          <TokenForm
            v-if="dialogInfo.operation === 'instance'"
            ref="tokenForm"
            :formData="formData.tokens[0]"
          ></TokenForm>
          <div class="mt-8 color-gray text-3 pb-4">
            {{ t('mcp.instance.openApi.tips') }}
          </div>
        </div>
      </el-splitter-panel>
      <el-splitter-panel
        size="50%"
        :min="600"
        class="p-4"
        :class="{
          'cursor-not-allowed': !!formData.templateId && dialogInfo.operation === 'template',
        }"
      >
        <div
          :class="{
            'disabled-click': !!formData.templateId && dialogInfo.operation === 'template',
          }"
        >
          <div class="flex-sub link-hover" v-if="!formData.openapiFileID">
            <div class="flex flex-col" style="height: 75vh">
              <el-card class="mb-3">
                <div class="flex items-center">
                  <mcp-image
                    :src="openapi"
                    width="28"
                    height="28"
                    border-radius="4"
                    class="mr-3"
                  ></mcp-image>
                  <div class="flex-sub">
                    {{ t('mcp.instance.openApi.support') }}
                  </div>
                </div>
              </el-card>
              <div
                class="flex-sub select-api border-rd-1 mt-2 center cursor-pointer"
                @click="handleSelectDocs"
              >
                <div class="text-center">
                  <el-icon class="el-icon--upload" size="67"><Files /></el-icon>
                  <div>{{ t('mcp.instance.openApi.selectDocs') }}</div>
                </div>
              </div>
            </div>
          </div>

          <div v-else class="p-3">
            <el-tree
              style="height: 75vh"
              :data="apiNodeList"
              show-checkbox
              node-key="id"
              :default-expand-all="true"
              :default-checked-keys="defaultCheckedKeys"
              ref="apiTreeRef"
              :render-content="renderContent"
            />
          </div>
        </div>
      </el-splitter-panel>
    </el-splitter>
    <template #footer>
      <div class="center">
        <el-button class="mr-4 w-25" @click="dialogInfo.visible = false">{{
          t('common.cancel')
        }}</el-button>
        <mcp-button class="w-25 mr-4" @click="handleConfirm">
          {{ t('common.save') }}
        </mcp-button>
        <!-- v-if="!formData.instanceId && !formData.templateId" -->
        <mcp-button v-if="dialogInfo.operation === 'instance'" @click="handleSaveAsTemplate"
          >{{ t('mcp.instance.action.asTemplate') }}
        </mcp-button>
      </div>
    </template>
  </el-dialog>
  <Select
    v-model="selectVisible"
    v-model:selected="formData.openapiFileID"
    ref="openAPISelect"
    :title="t('mcp.instance.pageDesc.apiSelectTitle')"
    @confirm="handleGetAPIDetail"
    :options="docsList"
  >
    <template #action>
      <el-upload
        class="mr-8"
        drag
        :action="action"
        :on-success="handleSuccess"
        :before-upload="handleBeforeUpload"
        :headers="headers"
        accept=".yaml, .JSON, application/yaml, application/JSON"
        :auto-upload="true"
        :show-file-list="false"
      >
        <el-icon><UploadFilled /></el-icon>
        <div class="ml-2">
          {{ t('mcp.instance.openApi.localFile') }}
        </div>
      </el-upload>
    </template>
  </Select>
</template>

<script lang="ts" setup>
import baseConfig from '@/config/base_config.ts'
import { Storage } from '@/utils/storage'
import { ElMessage } from 'element-plus'
import { UploadFilled, Files } from '@element-plus/icons-vue'
import McpImage from '@/components/mcp-image/index.vue'
import { openapi } from '@/utils/logo.ts'
import McpButton from '@/components/mcp-button/index.vue'
import Select from '@/components/mcp-select/index.vue'
import Upload from '@/components/upload/index.vue'
import yaml from 'js-yaml'
import { buildApiTree } from '@/utils/json.ts'
import { DocsAPI } from '@/api/api-docs'
import { useMcpStoreHook, useUserStore } from '@/stores'
import { InstanceAPI } from '@/api/mcp/instance'
import { getToken } from '@/utils/system'
import { AccessType, McpProtocol, SourceType, TokenType } from '@/types'
import { ElTooltip } from 'element-plus'
import TokenForm from './components/token-form.vue'
import { TemplateAPI } from '@/api/mcp/template'
import { useRouterHooks } from '@/utils/url'

const { userInfo } = useUserStore()
const { envList } = toRefs(useMcpStoreHook())
const { handleGetEnvList } = useMcpStoreHook()
const { jumpToPage } = useRouterHooks()
const { t } = useI18n()
const emit = defineEmits(['on-refresh'])
const dialogInfo = ref<any>({
  visible: false,
  loading: false,
  title: t('mcp.instance.openApi.importTitle'),
  operation: '',
})
const selectVisible = ref(false)
const docsList = ref<any>([])
const formData = ref({
  instanceId: '',
  templateId: '',
  name: '',
  notes: '',
  iconPath: '',
  environmentId: envList.value[0]?.id,
  openapiBaseUrl: '',
  openapiFileID: '',
  chooseOpenapiFileID: '', // 选择的文档库文件ID
  enabledToken: true,
  sourceType: SourceType.OPENAPI,
  tokens: [
    {
      enabled: true,
      expireAt: '',
      headers: [{ key: 'Authorization', value: '' }],
      publishAt: new Date().getTime(),
      token:
        'Bearer ' +
        getToken(
          JSON.stringify({
            expireAt: Date.now(),
            userId: userInfo.userId,
            username: userInfo.username,
          }),
        ),
      tokenType: TokenType.BEARER,
      usages: ['default'],
    },
  ],
})
const rules = ref({
  name: [
    { required: true, message: t('mcp.instance.rules.name'), trigger: 'blur' },
    // { type: 'string', max: 40, message: t('mcp.instance.rules.nameMax40'), trigger: 'blur' },
  ],
  openapiBaseUrl: [
    { required: true, message: t('mcp.instance.rules.openapiBaseUrl'), trigger: 'blur' },
  ],
  environmentId: [
    { required: true, message: t('mcp.instance.rules.environmentId'), trigger: 'change' },
  ],
})
const baseInfo = ref<any>(null)
const apiNodeList = ref<any[]>([])
const defaultCheckedKeys = ref<any[]>([])
const apiTreeRef = ref(null)
const docObject = ref<any>(null)
const originFileText = ref<any>(null)

const action = ref(
  baseConfig.SERVER_BASE_URL + (window as any).__APP_CONFIG__?.API_BASE + '/market/openapi/upload',
)
const headers = ref({
  Authorization: `Bearer ${Storage.get('token')}`,
})

/**
 * current checked keys
 */
const currentCheckedKeys = computed(() => {
  return (apiTreeRef as any).value.getCheckedKeys().filter((key: string) => key)
})

/**
 * handle Setect a openApi docs
 */
const handleSelectDocs = () => {
  selectVisible.value = true
}

/**
 * Get api list
 */
const handleGetAPIlist = async () => {
  const data = await DocsAPI.list({ page: 1, pageSize: 999 })
  docsList.value = data.list
}

/**
 * api list render
 * @param h
 * @param param1
 */
const renderContent = (h: any, params: any) => {
  const { data } = params as { node: any; data?: any }
  // build tooltip content
  const buildTooltipContent = () => {
    if (!data || !data.method) return ''
    const parts = []
    if (data.id) parts.push(`ID: <span class="font-500">${data.id}</span>`)
    if (data.method)
      parts.push(
        `Method: <span class="${
          {
            GET: 'color-green',
            POST: 'color-orange',
            PUT: 'color-blue',
            DELETE: 'color-red',
          }[data.method as string]
        }">${data.method.toUpperCase()}</span>`,
      )
    if (data.path) parts.push(`Path: <span class="font-500">${data.path}</span>`)
    if (data.summary) parts.push(`Summary: <span class="font-500">${data.summary}</span>`)
    return `<div class="font-bold">${parts.join('\n')}</div>`
  }

  const tooltipContent = buildTooltipContent()

  return h(
    'div',
    {
      class: ' w-full grid grid-cols-10',
    },
    [
      h(
        'div',
        { class: 'ellipsis-one col-span-9' },
        h('div', { class: 'flex' }, [
          h(
            'span',
            {
              class:
                {
                  GET: 'color-green',
                  POST: 'color-orange',
                  PUT: 'color-blue',
                  DELETE: 'color-red',
                }[data.method as string] + ' font-bold',
            },
            data.method,
          ),
          h('div', { class: 'flex-sub ml-2 u-line-1' }, data.label),
        ]),
      ),
      tooltipContent &&
        h(
          'div',
          { class: 'col-span-1 text-right pr-2' },
          h(
            ElTooltip,
            {
              content: tooltipContent,
              placement: 'left',
              rawContent: true,
              effect: 'dark',
              popperClass: 'api-detail-tooltip',
              style: { width: '400px' },
            },
            {
              default: () =>
                h(
                  'span',
                  { class: 'cursor-pointer', style: { color: 'var(--el-color-primary)' } },
                  t('common.more'),
                ),
            },
          ),
        ),
    ].filter(Boolean),
  )
}

const handleValidFile = (rawText: string) => {
  return new Promise((resolve, reject) => {
    try {
      docObject.value = JSON.parse(rawText)
      if (!docObject.value.openapi || docObject.value.openapi < '3.0.0') {
        ElMessage.error(t('mcp.instance.openApi.support3'))
        reject(false)
      }
    } catch {
      try {
        docObject.value = yaml.load(rawText)
        if (!docObject.value.openapi || docObject.value.openapi < '3.0.0') {
          ElMessage.error(t('mcp.instance.openApi.support3'))
          reject(false)
        }
      } catch (yamlErr) {
        console.warn('Error', yamlErr)
        reject(false)
      }
    }
    resolve(true)
  })
}

// handle default checked keys and server url、 version check
const handleDefaultCheckedKeys = (rawText: string) => {
  try {
    // keep original file text
    originFileText.value = rawText
    // try to parse as JSON/YAML (fault tolerance)
    try {
      docObject.value = JSON.parse(rawText)
      formData.value.openapiBaseUrl =
        formData.value.openapiBaseUrl ||
        (docObject.value.servers?.length ? docObject.value.servers[0]?.url : '')
      defaultCheckedKeys.value = []
      const collectIds = (nodes: any[]) => {
        nodes.forEach((n) => {
          if (n.id) defaultCheckedKeys.value.push(n.id)
          if (n.children && n.children.length) collectIds(n.children)
        })
      }
      collectIds(buildApiTree(docObject.value))
    } catch {
      try {
        docObject.value = yaml.load(rawText)
        formData.value.openapiBaseUrl =
          formData.value.openapiBaseUrl ||
          (docObject.value.servers?.length ? docObject.value.servers[0]?.url : '')
        defaultCheckedKeys.value = []
        const collectIds = (nodes: any[]) => {
          nodes.forEach((n) => {
            if (n.id) defaultCheckedKeys.value.push(n.id)
            if (n.children && n.children.length) collectIds(n.children)
          })
        }
        collectIds(buildApiTree(docObject.value))
        console.log(buildApiTree(docObject.value), defaultCheckedKeys.value)
      } catch (yamlErr) {
        console.warn('Error', yamlErr)
      }
    }
  } catch (err) {
    console.error('Error', err)
  }
}

// Handle default node api list
const handleDefaultNodeAPIlist = (rawText: string) => {
  try {
    // keep original text for later use
    originFileText.value = rawText
    // try to parse as JSON/YAML (fault tolerance)
    try {
      docObject.value = JSON.parse(rawText)
      apiNodeList.value = [
        {
          id: 'root',
          label: t('mcp.instance.openApi.interface'),
          children: buildApiTree(docObject.value),
        },
      ]
      const collectIds = (nodes: any[]) => {
        nodes.forEach((n) => {
          if (n.children && n.children.length) collectIds(n.children)
        })
      }
      collectIds(apiNodeList.value)
    } catch {
      // Not JSON, try to parse as YAML
      try {
        docObject.value = yaml.load(rawText)
        apiNodeList.value = [
          {
            id: 'root',
            label: t('mcp.instance.openApi.interface'),
            children: buildApiTree(docObject.value),
          },
        ]
        const collectIds = (nodes: any[]) => {
          nodes.forEach((n) => {
            if (n.children && n.children.length) collectIds(n.children)
          })
        }
        collectIds(apiNodeList.value)
      } catch (yamlErr) {
        console.warn('Error', yamlErr)
      }
    }
  } catch (err) {
    console.error('Error', err)
  }
}

/**
 * get build api tree
 * @param file Handle before upload
 */
const handleBeforeUpload = async (file: File) => {
  const rawText = await file.text()
  // validate file content
  await handleValidFile(rawText)
  // keep original text for later use
  // handleDefaultCheckedKeys(rawText)
  // handleDefaultNodeAPIlist(rawText)
}

/**
 *  Handle update success
 * @param response
 */
const handleSuccess = (response: { code: number; data: { openapiFileId: string } }) => {
  if (response.code !== 0) {
    return
  }
  formData.value.openapiFileID = response.data.openapiFileId
  ElMessage.success(t('action.upload'))
  handleGetAPIlist()
}

// get openapi file detail
const handleGetAPIDetail = async (id: string) => {
  try {
    dialogInfo.value.loading = true
    const { content } = await DocsAPI.fileContent({ openapiFileId: id })
    // validate file content
    await handleValidFile(content)
    handleDefaultCheckedKeys(content)
    handleDefaultNodeAPIlist(content)
  } catch {
    formData.value.openapiFileID = ''
  } finally {
    dialogInfo.value.loading = false
  }
}

/**
 * Handle file upload again
 */
const handleUploadAgain = async () => {
  // Reset form data
  const newDoc = JSON.parse(JSON.stringify(docObject.value))
  for (const path in newDoc.paths) {
    const obj = newDoc.paths[path] || {}
    for (const method in obj) {
      const mKey = method as unknown as string
      if (!currentCheckedKeys.value.includes(obj[mKey].operationId)) {
        delete obj[mKey]
      }
    }
    if (Object.keys(newDoc.paths[path]).length === 0) {
      delete newDoc.paths[path]
    }
  }

  // Convert to YAML and upload
  const yamlText = yaml.dump(newDoc)
  const fd = new FormData()
  const blob = new Blob([yamlText], { type: 'application/x-yaml' })
  fd.append('file', blob, 'openapi.yaml')
  fd.append('baseOpenapiFileID', formData.value.openapiFileID)
  try {
    dialogInfo.value.loading = true
    const res = await fetch(action.value, {
      method: 'POST',
      headers: Object.assign({}, headers.value),
      body: fd,
    })
    const body = await res.json()
    if (body?.code === 0 && body.data?.openapiFileId) {
      formData.value.chooseOpenapiFileID = body.data.openapiFileId
    } else {
      ElMessage.error(t('action.uploadFail'))
    }
  } catch (err) {
    console.error('upload error', err)
    ElMessage.error(t('action.uploadFail'))
  } finally {
    dialogInfo.value.loading = false
  }
}

/**
 * Handle confirm and submit data
 */
const handleConfirm = async () => {
  if (formData.value.templateId && dialogInfo.value.operation === 'template') {
    handleSaveAsTemplate()
    return
  }
  try {
    // handle selected APIs; assemble the selected paths into a new OpenAPI document and upload
    if (originFileText.value) {
      await handleUploadAgain()
    } else {
      ElMessage.error(t('mcp.instance.openApi.validFileFail'))
    }
    baseInfo.value.validate(async (valid: boolean) => {
      if (currentCheckedKeys.value.length === 0) {
        ElMessage.error(t('mcp.instance.openApi.selectAtLeastOne'))
        return
      }
      if (valid) {
        if (!formData.value.instanceId) {
          formData.value.tokens[0].headers = Object.fromEntries(
            formData.value.tokens[0].headers.map((header: any) => [header.key, header.value]),
          )
        }
        dialogInfo.value.loading = true
        // 提交数据
        await (formData.value.instanceId
          ? InstanceAPI.editByOpenAPI(formData.value)
          : InstanceAPI.createByOpenAPI(formData.value))
        dialogInfo.value.visible = false
        ElMessage.success(formData.value.instanceId ? t('action.edit') : t('action.create'))
        emit('on-refresh')
        if (dialogInfo.value.operation === 'create') {
          jumpToPage({
            url: '/instance-manage',
            data: {},
          })
        }
      }
    })
  } finally {
    dialogInfo.value.loading = false
  }
}

/**
 * save as a template
 */
const handleSaveAsTemplate = async () => {
  try {
    dialogInfo.value.loading = true
    baseInfo.value.validate(async (valid: boolean) => {
      if (valid) {
        await (formData.value.templateId ? TemplateAPI.edit : TemplateAPI.create)({
          ...formData.value,
          packageId: formData.value.openapiFileID,
          accessType: AccessType.HOSTING,
          mcpProtocol: McpProtocol.STREAMABLE_HTTP,
          sourceType: SourceType.OPENAPI,
        })
        ElMessage.success(formData.value.templateId ? t('action.edit') : t('action.create'))
        emit('on-refresh')
        dialogInfo.value.visible = false
      }
    })
  } finally {
    dialogInfo.value.loading = false
  }
}

/**
 * Handle instance detail and openAPI data
 * @param id instance id
 */
const handleGetDetail = async (id: string) => {
  try {
    dialogInfo.value.loading = true
    const data = await InstanceAPI.detail({ instanceId: id })
    formData.value = {
      ...data,
      chooseOpenapiFileID: data.packageId,
      openapiFileID: '',
    }
    console.log('instance detail data1', data)

    // get openapi file content
    const { baseOpenapiFileID, content } = await DocsAPI.fileContent({
      openapiFileId: data.packageId,
    })
    // validate file content
    await handleValidFile(content)
    // handle default checked keys
    handleDefaultCheckedKeys(content)
    formData.value.openapiFileID = baseOpenapiFileID
    // get original api document content
    const res = await DocsAPI.fileContent({
      openapiFileId: baseOpenapiFileID,
    })
    handleDefaultNodeAPIlist(res.content)
    console.log('instance detail data2', formData.value)
  } finally {
    dialogInfo.value.loading = false
  }
}

/**
 * handle template detail
 * @param id template id
 */
const handleTemplateDetail = async (id: string) => {
  const data = await TemplateAPI.detail({
    id,
  })
  const tokenValue =
    'Bearer ' +
    getToken(
      JSON.stringify({
        expireAt: Date.now(),
        userId: userInfo.userId,
        username: userInfo.username,
      }),
    )
  formData.value = {
    instanceId: '',
    templateId: data.templateId,
    name: data.name,
    notes: data.notes,
    iconPath: data.iconPath,
    openapiBaseUrl: data.openapiBaseUrl,
    environmentId: envList.value[0]?.id || data.environmentId,
    enabledToken: true,
    openapiFileID: data.packageId,
    chooseOpenapiFileID: '',
    sourceType: SourceType.OPENAPI,
    tokens: [
      {
        enabled: true,
        expireAt: '',
        headers: [{ key: 'Authorization', value: '' }],
        publishAt: new Date().getTime(),
        token: tokenValue,
        tokenType: TokenType.BEARER,
        usages: ['default'],
      },
    ],
  }
  // get openapi file content
  const { baseOpenapiFileID, content } = await DocsAPI.fileContent({
    openapiFileId: data.packageId,
  })
  // validate file content
  await handleValidFile(content)
  // handle default checked keys
  handleDefaultCheckedKeys(content)

  handleDefaultNodeAPIlist(content)
}

/**
 * Init dialog data
 */
const init = async (id?: string, type?: string) => {
  dialogInfo.value.visible = true
  dialogInfo.value.operation = type
  baseInfo.value?.resetFields()
  handleGetAPIlist()
  await handleGetEnvList()
  if (id) {
    // 模板编辑
    if (type === 'template' || type === 'instance') {
      handleTemplateDetail(id)
    } else {
      // get detail
      handleGetDetail(id)
    }
  } else {
    const tokenValue =
      'Bearer ' +
      getToken(
        JSON.stringify({
          expireAt: Date.now(),
          userId: userInfo.userId,
          username: userInfo.username,
        }),
      )
    formData.value = {
      instanceId: '',
      templateId: '',
      name: '',
      notes: '',
      iconPath: '',
      openapiBaseUrl: '',
      environmentId: envList.value[0]?.id,
      enabledToken: true,
      openapiFileID: '',
      chooseOpenapiFileID: '',
      sourceType: SourceType.OPENAPI,
      tokens: [
        {
          enabled: true,
          expireAt: '',
          headers: [{ key: 'Authorization', value: '' }],
          publishAt: new Date().getTime(),
          token: tokenValue,
          tokenType: TokenType.BEARER,
          usages: ['default'],
        },
      ],
    }
  }
}

defineExpose({
  init,
})
</script>

<style lang="scss" scoped>
.api-info {
  border: 1px dashed var(--el-color-primary);
}
.upload-demo {
  width: 100%;
  color: var(--el-color-primary);
  :deep(.el-upload.is-drag) {
    height: 100%;
  }
  .el-icon--upload {
    color: var(--el-color-primary);
  }
  .el-upload__text {
    color: var(--el-color-primary);
  }
}
.select-api {
  border: 1px dashed var(--el-color-primary);
  color: var(--el-color-primary);
  &:hover {
    border-color: var(--el-color-primary-hover);
    color: var(--el-color-primary-hover);
  }
}
:deep(.el-tree-node__content) {
  padding: 0 0 1px 4px;
  margin: 2px 0;
  &:hover {
    background-color: var(--ep-bg-purple-color);
    border-radius: 4px;
  }
}
:deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
  background-color: var(--el-color-primary);
  border-color: var(--el-color-primary);
  border-radius: 4px;
}
:deep(.el-checkbox__input.is-indeterminate .el-checkbox__inner) {
  background-color: var(--el-color-primary);
}
:deep(.el-checkbox__input.is-indeterminate .el-checkbox__inner) {
  border-color: var(--el-color-primary);
}
:deep(.el-checkbox__inner) {
  &:hover {
    border-color: var(--el-color-primary);
  }
}
</style>

<style lang="scss">
.el-dialog__header.header-border {
  background-color: transparent !important;
  border-bottom: 1px solid var(--el-border-color-light) !important;
}
.el-dialog__footer.footer-border {
  background-color: transparent !important;
  border-top: 1px solid var(--el-border-color-light) !important;
}
.api-detail-tooltip {
  max-width: 400px !important;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
