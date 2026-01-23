import { useStorage } from '@vueuse/core'

export const useMarketListHooks = () => {
  const { t } = useI18n()
  const loading = ref(false)
  const searchInputRef = ref()
  const keyword = useStorage('market-list-keyword', '')
  const categoryName = useStorage('market-list-category', '')
  const typeCount = ref<any>([])
  const typeMap = ref([
    { label: t('market.type.test'), value: 'test', icon: 'MCP-ceshi', count: 0 },
    {
      label: t('market.type.browser-automation'),
      value: 'browser-automation',
      icon: 'MCP-liulanqi',
      count: 0,
    },
    { label: t('market.type.search'), value: 'search', icon: 'MCP-xinxijiansuo', count: 0 },
    {
      label: t('market.type.communication'),
      value: 'communication',
      icon: 'MCP-jiaoliu',
      count: 0,
    },
    {
      label: t('market.type.developer-tools'),
      value: 'developer-tools',
      icon: 'MCP-kaifagongju',
      count: 0,
    },
    {
      label: t('market.type.entertainment-and-media'),
      value: 'entertainment-and-media',
      icon: 'MCP-meiti',
      count: 0,
    },
    { label: t('market.type.file-systems'), value: 'file-systems', icon: 'MCP-wenjian1', count: 0 },
    {
      label: t('market.type.finance'),
      value: 'finance',
      icon: 'MCP-financialStatements',
      count: 0,
    },
    {
      label: t('market.type.knowledge-and-memory'),
      value: 'knowledge-and-memory',
      icon: 'MCP-zhishi',
      count: 0,
    },
    {
      label: t('market.type.location-services'),
      value: 'location-services',
      icon: 'MCP-position',
      count: 0,
    },
    {
      label: t('market.type.art-and-culture'),
      value: 'art-and-culture',
      icon: 'MCP-yishu',
      count: 0,
    },
    {
      label: t('market.type.research-and-data'),
      value: 'research-and-data',
      icon: 'MCP-kexuejiaoyu',
      count: 0,
    },
    {
      label: t('market.type.calendar-management'),
      value: 'calendar-management',
      icon: 'MCP-richeng',
      count: 0,
    },
    {
      label: t('market.type.efficiency-tools'),
      value: 'efficiency-tools',
      icon: 'MCP-xiaoshuaigongju',
      count: 0,
    },
    { label: t('market.type.other'), value: 'other', icon: 'MCP-qita', count: 0 },
  ])
  const pagerConfig = useStorage('market-list-pager', {
    total: 0,
    page: 1,
    pageSize: 12,
  })

  return {
    t,
    loading,
    typeMap,
    typeCount,
    categoryName,
    keyword,
    searchInputRef,
    pagerConfig,
  }
}
