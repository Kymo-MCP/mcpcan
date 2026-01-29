<template>
  <div>
    <!-- Header Info -->
    <div class="flex justify-between items-start">
      <div class="flex gap-6">
        <!-- Logo -->
        <!-- <div
          class="w-24 h-24 rounded-2xl border border-gray-200 flex items-center justify-center bg-white shadow-sm overflow-hidden"
        > -->
        <mcp-image :src="iconUrl(currentMCP)" width="96" height="96"></mcp-image>
        <!-- <img :src="iconUrl(currentMCP)" alt="Logo" class="w-full h-full object-contain" /> -->
        <!-- </div> -->

        <!-- Content -->
        <div class="flex flex-col justify-between py-1">
          <h1 class="text-3xl font-bold leading-tight">
            {{ locale === 'zh-cn' ? currentMCP.name : currentMCP.nameEn }}
          </h1>

          <div class="flex items-center gap-3 text-gray-500 text-sm mt-1">
            <span class="font-bold text-base">{{ currentMCP.githubOwner }}</span>
            <span class="w-1 h-1 rounded-full bg-gray-400"></span>
            <span>{{ t('market.updateTime') }}: {{ timestampToDate(currentMCP.updateTime) }}</span>
          </div>
          <div class="flex items-center gap-4 mt-3">
            <el-tag
              type="primary"
              effect="dark"
              class="!border-none px-2 py-1.5 text-sm"
              v-for="(type, index) in currentMCP.categoryIds"
              :key="index"
            >
              {{ translationTag(type.code) }}
            </el-tag>
            <div class="flex items-center gap-2 text-gray-600">
              <el-icon class="cursor-pointer">
                <i class="icon iconfont MCP-fork"></i>
              </el-icon>
              <span class="text-base font-medium">{{
                githubNumber(currentMCP.githubForksCount)
              }}</span>
            </div>
            <div class="flex items-center gap-2 text-gray-600">
              <el-icon class="cursor-pointer">
                <i class="icon iconfont MCP-GitHub"></i>
              </el-icon>
              <span class="text-base font-medium">{{
                githubNumber(currentMCP.githubStargazersCount)
              }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Github Button -->
      <div>
        <div v-if="layout" class="mb-4 text-right">
          <el-button @click="handleBack" class="link-hover">
            <el-icon>
              <i class="icon iconfont MCP-fanhui"></i>
            </el-icon>
            {{ t('common.back') }}
          </el-button>
        </div>
        <div>
          <el-button round size="large" @click="goGithub">
            <template #icon>
              <el-icon class="cursor-pointer link-hover">
                <i class="icon iconfont MCP-GitHub"></i>
              </el-icon>
            </template>
            GitHub
          </el-button>
        </div>
      </div>
    </div>

    <el-divider />

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <!-- Left Column -->
      <div class="lg:col-span-2 flex flex-col gap-6">
        <!-- Description Card -->
        <el-card>
          <h2 class="text-lg font-bold mb-4">{{ t('market.descWhatIsMCP') }}</h2>
          <p class="leading-relaxed text-sm">
            {{ locale === 'zh-cn' ? currentMCP.description : currentMCP.descriptionEn }}
          </p>
        </el-card>

        <!-- Readme Card -->
        <el-card>
          <!-- <h2 class="text-3xl font-bold mb-8">MCP-Auth （readme.md）</h2> -->
          <div class="prose max-w-none dark:prose-invert" v-html="renderedReadme"></div>
        </el-card>
      </div>

      <!-- Right Column -->
      <div class="lg:col-span-1 flex flex-col gap-6 config-sticky" v-if="currentMCP.status === 1">
        <!-- Config Card -->
        <el-card>
          <h2 class="text-lg font-bold mb-4">{{ t('market.config') }}</h2>
          <div
            class="bg-gray-100 dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700 overflow-hidden"
          >
            <pre
              class="text-xs font-mono text-gray-600 dark:text-gray-300 overflow-x-auto whitespace-pre-wrap break-all"
              >{{ JsonFormatter.format(currentMCP.configTemplate) }}
            </pre>
          </div>
        </el-card>

        <!-- Action Button -->
        <mcp-button @click="handleCreate" size="large" class="w-full">{{
          t('market.quickCreate')
        }}</mcp-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useMarketDetailHooks } from './hooks/market-detail.ts'
import { timestampToDate, githubNumber } from '@/utils/system'
import { JsonFormatter } from '@/utils/json'
import MarkdownIt from 'markdown-it'
import baseConfig from '@/config/base_config.ts'
import McpButton from '@/components/mcp-button/index.vue'
import McpImage from '@/components/mcp-image/index.vue'
import { AccessType } from '@/types/instance.ts'

const layout = useLayout()
const { t, locale, jumpBack, jumpToPage, currentMCP } = useMarketDetailHooks()

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
})
const renderedReadme = computed(() => {
  return md.render(currentMCP.githubReadme || '')
})
const translationTag = (code: string) => {
  return t('market.type.' + code)
}
const iconUrl = computed(() => {
  return (card: any) => {
    if (card.iconUrl) {
      return card.iconUrl
    } else if (card.githubOwnerAvatarUrl) {
      return card.githubOwnerAvatarUrl
    } else {
      return '@/assets/logo.png'
    }
  }
})
// back last class page
const handleBack = () => {
  jumpBack()
}

const goGithub = () => {
  window.open(currentMCP.githubRepoUrl, '_blank')
}

const handleCreate = () => {
  jumpToPage({
    url: '/new-instance',
    data: {
      from: 'market',
      mcpId: currentMCP.id,
      type: AccessType.HOSTING,
    },
  })
}
</script>

<style lang="scss" scoped>
:deep(.el-button) {
  --el-button-hover-bg-color: #fff;
  --el-button-hover-text-color: var(--el-color-primary);
  --el-button-hover-border-color: var(--el-color-primary-light-7);
}
.config-sticky {
  position: sticky;
  top: 0px;
  align-self: flex-start;
  height: fit-content;
  z-index: 1;
}
</style>
