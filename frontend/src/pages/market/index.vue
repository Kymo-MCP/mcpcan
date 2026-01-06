<template>
  <div>
    <!-- <div><el-card></el-card></div> -->
    <div v-loading="loading" class="flex gap-4 mt-4">
      <div class="w-75 type-sticky" v-if="false">
        <el-input
          ref="searchInputRef"
          v-model="keyword"
          :placeholder="t('market.placeholderName')"
          class="mb-4"
          :suffix-icon="Search"
          @keyup.enter="handleQuery"
        />
        <el-card>
          <div class="flex justify-between items-center">
            <span class="my-1">{{ t('market.search') }}</span>
            <el-button v-if="categoryName" type="primary" size="small" @click="clearType">
              {{ t('market.clear') }}
            </el-button>
          </div>
          <div class="mt-3 flex flex-col gap-2">
            <div
              v-for="type in typeMap"
              :key="type.value"
              class="w-full flex items-center gap-3 rounded-md px-3 py-2 text-left type-item"
              :class="{ 'active-type': categoryName === type.value }"
              @click="selectType(type.value)"
            >
              <el-icon>
                <i class="icon iconfont" :class="type.icon"></i>
              </el-icon>
              <span class="text-sm">{{ type.label }}</span>
            </div>
          </div>
        </el-card>
      </div>
      <div class="flex-1">
        <template v-if="marketList.length">
          <el-row :gutter="20">
            <el-col
              v-for="(card, index) in marketList"
              :key="index"
              :xs="24"
              :sm="24"
              :md="12"
              :lg="8"
              :xl="6"
              class="mb-4"
            >
              <McpCard :card="card"></McpCard>
            </el-col>
          </el-row>

          <div class="mt-8 flex justify-end">
            <el-pagination
              background
              :total="pagerConfig.total"
              :current-page="pagerConfig.page"
              :page-size="pagerConfig.pageSize"
              @current-change="handlePageChange"
            />
          </div>
        </template>
        <el-empty v-else class="mt-40"></el-empty>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { Search } from '@element-plus/icons-vue'
import { useMarketListHooks } from './hooks/index.ts'
import McpCard from './modules/mcp-card.vue'
import { MarketAPI } from '@/api/market/index.ts'

const { t, loading, typeMap, categoryName, keyword, searchInputRef, pagerConfig } =
  useMarketListHooks()

const marketList = ref([])

// handle selected type
const selectType = (value: string) => {
  pagerConfig.value.page = 1
  categoryName.value = value
  handleGetMarketList()
}

// handle clear type
const clearType = () => {
  pagerConfig.value.page = 1
  categoryName.value = ''
  keyword.value = ''
  handleGetMarketList()
}

// handle Search
const handleQuery = () => {
  pagerConfig.value.page = 1
  handleGetMarketList()
}

const handlePageChange = (newPage: number) => {
  pagerConfig.value.page = newPage
  handleGetMarketList()
}

// handle get market list data
const handleGetMarketList = async () => {
  try {
    loading.value = true
    const { list, total } = await MarketAPI.list({
      page: pagerConfig.value.page,
      pageSize: pagerConfig.value.pageSize,
      name: keyword.value,
      categoryName: categoryName.value,
    })
    marketList.value = list || []
    pagerConfig.value.total = Number(total || 0)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  handleGetMarketList()
  // searchInputRef.value.$el.getElementsByClassName('el-input__suffix')[0].onclick = handleQuery
})
</script>

<style lang="scss" scoped>
.type-sticky {
  position: sticky;
  top: 0px;
  align-self: flex-start;
  height: fit-content;
  z-index: 1;
}
:deep(.el-input__suffix) {
  cursor: pointer;
}
.type-item {
  transition: all 0.3s ease;
  overflow: hidden;
  background-color: var(--ep-bg-purple-color);
  border: 1px solid transparent;
  &.active-type {
    background-color: var(--ep-bg-purple-color-deep);
  }
  &:hover {
    scale: 1.02;
    cursor: pointer;
    background-color: var(--ep-bg-purple-color-deep);
    border-color: var(--ep-btn-color-top);
  }
}
</style>
