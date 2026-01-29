<template>
  <div v-loading="pageInfo.loading" :element-loading-text="pageInfo.loadingText">
    <!-- 头部区域 -->
    <!-- <div class="flex justify-between page-header">
      <div class="header-container">
        {{ t('api.pageDesc.title') }} <span class="desc">{{ t('api.pageDesc.desc') }}</span>
      </div>
      <mcp-button :icon="UploadFilled" @click="handleUpdatePackage">{{
        t('api.action.upload')
      }}</mcp-button>
    </div> -->
    <div class="flex-sub center link-hover mb-2">
      <el-upload
        class="upload-demo"
        drag
        :action="action"
        multiple
        :on-success="handleSuccess"
        :headers="headers"
        accept=".json, .yaml, application/json, application/yaml"
      >
        <div class="flex justify-between align-center upload-dragger-content">
          <div class="title">
            <div>{{ t('api.action.upload') }}</div>
            <div class="desc">
              <div class="my-1">JSON {{ t('api.pageDesc.openAI') }} (.json)</div>
              <div>YAML {{ t('api.pageDesc.openAI') }} (.yaml)</div>
            </div>
          </div>
          <div>
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              {{ t('code.desc.suport') }}
            </div>
          </div>
          <div class="footer">
            {{ t('code.desc.describe') }}
            <div class="desc">
              <div class="ml-8 mt-2">{{ t('api.pageDesc.text1') }}</div>
              <div class="ml-8 mt-2">{{ t('api.pageDesc.text2') }}</div>
              <div class="ml-8 mt-2">{{ t('api.pageDesc.text3') }}</div>
              <div class="ml-8 mt-2">{{ t('api.pageDesc.text4') }}</div>
              <div class="ml-8 mt-2">{{ t('api.pageDesc.text5') }}</div>
            </div>
          </div>
        </div>
      </el-upload>
    </div>

    <TablePlus
      :showOperation="true"
      searchContainer="#apiDocsSearch"
      ref="tablePlus"
      :requestConfig="requestConfig"
      :columns="columns"
      v-model:pageConfig="pageConfig"
      :handlerColumnConfig="{
        width: '180px',
        fixed: 'right',
        align: 'center',
      }"
    >
      <template #action>
        <div class="flex justify-between mb-4">
          <div class="center">
            <el-image :src="codeLogo" style="width: 20px; height: 20px"></el-image>
            <span class="desc">{{ t('api.pageDesc.total') }}：{{ pageConfig.total }}</span>
          </div>
          <div id="apiDocsSearch"></div>
        </div>
      </template>
      <template #name="{ row }">
        <div class="flex align-center">
          <el-image
            :src="zipLogo"
            style="width: 32px; height: 32px"
            fit="cover"
            class="mr-2"
          ></el-image>
          <el-tooltip effect="dark" placement="top" class="flex-sub" :raw-content="true">
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
        <el-button type="text" size="small" link @click="handleViewCode(row)" class="base-btn-link">
          {{ t('common.view') }}
        </el-button>
        <el-button type="text" size="small" link @click="handleDownload(row)" class="base-btn-link">
          {{ t('common.download') }}
        </el-button>
        <el-button type="danger" size="small" link @click="handleDelete(row)">{{
          t('common.delete')
        }}</el-button>
      </template>
    </TablePlus>
  </div>
</template>

<script setup lang="ts">
import { UploadFilled, More } from '@element-plus/icons-vue'
import { DocsAPI } from '@/api/api-docs/index'
import TablePlus from '@/components/TablePlus/index.vue'
import { useDocsTableHooks } from './hooks.ts'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouterHooks } from '@/utils/url'
import codeLogo from '@/assets/logo/code.png'
import zipLogo from '@/assets/logo/zip.png'
import McpButton from '@/components/mcp-button/index.vue'

const { t, tablePlus, columns, requestConfig, pageConfig, pageInfo, action, headers } =
  useDocsTableHooks()
const { jumpToPage } = useRouterHooks()

/**
 * Handle jump to the page of update code package
 */
const handleUpdatePackage = () => {
  jumpToPage({
    url: '/update-api-docs',
  })
}
const handleSuccess = (response: { code: number; data: { path: string } }) => {
  if (response.code !== 0) {
    return
  }
  init()
  ElMessage.success(t('action.upload'))
}
/**
 * Handle view code package info
 */
const handleViewCode = (code: any) => {
  jumpToPage({
    url: '/view-api-docs',
    data: {
      id: code.id,
      name: code.name,
    },
  })
}

/**
 * Handle more operation events
 * @param callback - function name
 * @param row - item of code package
 */
const handleCommand = (callback: string, row: any) => {
  switch (callback) {
    case 'handleDownload':
      handleDownload(row)
      break
    case 'handleDelete':
      handleDelete(row)
      break
    default:
      ElMessage.warning(`未找到 "${callback}" 对应的操作`)
  }
}

/**
 * Handle download code package
 */
const handleDownload = async (code: any) => {
  try {
    pageInfo.value.loading = true
    const response = await DocsAPI.download(code)
    const blobUrl = URL.createObjectURL(
      new Blob([response.data], { type: response.headers['content-type'] }),
    )
    const link = document.createElement('a')
    link.href = blobUrl
    link.download =
      response.headers['content-disposition']
        ?.split('filename=')[1]
        ?.match(/filename=("?)(.*?)\1/) || code.name
    document.body.appendChild(link)
    link.click()
  } finally {
    pageInfo.value.loading = false
  }
}

/**
 * Handle delete code package
 */
const handleDelete = (code: any) => {
  ElMessageBox.confirm(`${t('api.action.deleteBox')}“${code.name}”?`, t('common.warn'), {
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
    await DocsAPI.delete(code.id)
    ElMessage({
      type: 'success',
      message: t('action.delete'),
    })
    init()
  })
}

/**
 * Handle init page list
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
.upload-demo {
  width: 100%;
  color: var(--el-color-primary);
  .title {
    font-size: 20px;
    font-weight: 600;
    text-align: left;
    .desc {
      font-size: 14px;
      color: #999999;
      font-weight: 400;
      margin-left: 0;
    }
  }
  .footer {
    font-family:
      PingFangSC,
      PingFang SC;
    font-weight: 600;
    font-size: 20px;
    // color: #cccccc;
    line-height: 28px;
    .desc {
      font-family:
        PingFangSC,
        PingFang SC;
      font-weight: 400;
      font-size: 14px;
      color: #999999;
      line-height: 20px;
      text-align: left;
      font-style: normal;
    }
  }
  :deep(.el-upload-dragger) {
    border: 1px dashed var(--el-color-primary);
    &:hover {
      border-color: var(--el-color-primary-hover);
      .el-icon--upload {
        color: var(--el-color-primary-hover);
      }
      .el-upload__text {
        color: var(--el-color-primary-hover);
      }
      .title {
        color: var(--el-color-primary-hover);
      }
      .footer {
        color: var(--el-color-primary-hover);
      }
    }
  }
  .el-icon--upload {
    color: var(--el-color-primary);
  }
  .el-upload__text {
    color: var(--el-color-primary);
  }
}
</style>
