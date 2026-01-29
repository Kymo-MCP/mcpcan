import { AgentAPI } from '@/api/agent/index'
import { timestampToDate } from '@/utils/system'
import { kymo, dify, coze, n8n } from '@/utils/logo.ts'
import { AgentType } from '@/types/agent'

export const useAgentTableHooks = () => {
  const { t } = useI18n()
  const tablePlus = ref()
  const pageInfo = ref({
    loading: false,
    loadingText: t('agent.pageDesc.connectionDesc'),
  })
  const logoIcon = ref<any>({
    DifyEnterprise: dify,
    Dify: dify,
    COZE: coze,
    N8N: n8n,
  })
  const columns = ref<any>([
    {
      label: t('agent.columns.accessName'),
      dataIndex: 'accessName',
      searchConfig: {
        component: 'el-input',
        label: t('agent.columns.accessName'),
        props: {
          placeholder: t('agent.columns.accessName'),
        },
      },
    },
    {
      label: t('agent.columns.accessType'),
      dataIndex: 'accessType',
      searchConfig: {
        component: 'el-select',
        label: t('agent.columns.accessType'),
        props: {
          placeholder: t('agent.columns.accessType'),
          options: [
            { label: t('agent.action.community'), value: AgentType.DIFY },
            { label: t('agent.action.enterprise'), value: AgentType.DIFY_ENTERPRISE },
          ],
        },
      },
      customRender: ({ row }: { row: any }) => {
        return row.accessType === AgentType.DIFY
          ? t('agent.action.community')
          : t('agent.action.enterprise')
      },
    },
    {
      label: t('agent.columns.dbHost'),
      dataIndex: 'dbHost',
    },
    {
      dataIndex: 'createTime',
      label: t('api.columns.createdAt'),
      customRender: ({ row }: { row: any }) => {
        return timestampToDate(row.createTime)
      },
    },
    {
      dataIndex: 'updateTime',
      label: t('api.columns.updatedAt'),
      customRender: ({ row }: { row: any }) => {
        return timestampToDate(row.updateTime)
      },
    },
  ])
  const requestConfig = {
    api: AgentAPI.list,
    searchQuery: {
      model: {},
    },
  }
  const pageConfig = ref({
    total: 0,
    page: 1,
    pageSize: 10,
  })

  return { t, columns, requestConfig, tablePlus, pageConfig, pageInfo, AgentAPI, logoIcon }
}
