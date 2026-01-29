<template>
  <div>
    <el-tabs v-model="activeName" class="custom-tabs">
      <el-tab-pane name="first">
        <template #label>
          <span class="mx-2">
            <el-icon>
              <i class="icon iconfont MCP-daimaguanli"></i>
            </el-icon>
            <span class="ml-2">{{ t('sideMenu.package') }}</span>
          </span>
        </template>
        <Code></Code>
      </el-tab-pane>
      <el-tab-pane name="second">
        <template #label>
          <span class="mx-2">
            <el-icon>
              <i class="icon iconfont MCP-wenjian"></i>
            </el-icon>
            <span class="ml-2">open API</span>
          </span>
        </template>
        <Docs></Docs>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import Code from '../script-code/code-list/index.vue'
import Docs from '../api-docs/docs-list/index.vue'

const { t } = useI18n()
const activeName = ref('first')
</script>

<style lang="scss" scoped>
.custom-tabs {
  :deep(.el-tabs__header) {
    // margin-bottom: 24px;
    border-bottom: none;
  }

  :deep(.el-tabs__nav-wrap) {
    &::after {
      display: none;
    }
  }

  :deep(.el-tabs__nav) {
    background-color: var(--ep-home-glass);
    padding: 4px;
    border-radius: 8px;
    display: inline-flex;
  }

  :deep(.el-tabs__active-bar) {
    display: none;
  }

  :deep(.el-tabs__item) {
    padding: 0;
    height: 40px;
    line-height: 40px;
    border-radius: 6px;
    transition: all 0.3s;
    margin-bottom: 0;
    border: none;
    font-size: 14px;
    font-weight: 500;
    &.is-active {
      background-color: var(--ep-bg-purple-color-deep);
      box-shadow:
        0 1px 3px 0 rgba(0, 0, 0, 0.1),
        0 1px 2px -1px rgba(0, 0, 0, 0.1);
    }

    &:hover:not(.is-active) {
      color: var(--el-color-primary);
      background-color: rgba(255, 255, 255, 0.6);
    }
  }
}
</style>
