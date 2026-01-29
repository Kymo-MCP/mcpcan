import { useMcpStoreHook, useUserStore } from '@/stores'
import { useRouterHooks } from '@/utils/url'
import {
  type PvcResult,
  AccessType,
  McpProtocol,
  SourceType,
  InstanceData,
  NodeVisible,
} from '@/types/index'
import { getToken } from '@/utils/system'

export const useInstanceFormHooks = () => {
  const { t } = useI18n()
  const { jumpBack, jumpToPage } = useRouterHooks()
  const router = useRouter()
  const { query, meta, path } = useRoute()
  const selectVisible = ref(false)
  const originForm = ref<any>()
  const { userInfo } = useUserStore()
  const { currentMCP } = useMcpStoreHook()

  const tokenValue =
    'Bearer ' +
    getToken(
      JSON.stringify({
        expireAt: Date.now(),
        userId: userInfo.userId,
        username: userInfo.username,
      }),
    )
  const pageInfo = ref<any>({
    visible: false,
    loading: false,
    title: t('mcp.instance.formData.title'),
    formData: {
      sourceType: SourceType.CUSTOM,
      name: '',
      accessType: '',
      mcpProtocol: '',
      imgAddress: InstanceData.value.IMGADDRESS,
      notes: '',
      mcpServers: '',
      iconPath: '',
      packageId: '',
      environmentId: '',
      port: InstanceData.value.PORT,
      environmentVariables: [],
      volumeMounts: [],
      initScript: InstanceData.value.INITSCRIPT,
      command: '',
      enabledToken: true,
      tokens: [
        {
          enabled: true,
          expireAt: '',
          publishAt: new Date().getTime(),
          headers: [{ key: 'Authorization', value: '' }],
          token: tokenValue,
          usages: ['default'],
        },
      ],
    },
    rules: {
      name: [
        { required: true, message: t('mcp.instance.rules.name'), trigger: 'blur' },
        // { type: 'string', max: 40, message: t('mcp.instance.rules.nameMax40'), trigger: 'blur' },
      ],
      accessType: [
        { required: true, message: t('mcp.template.rules.deployType'), trigger: 'change' },
      ],
      mcpProtocol: [
        { required: true, message: t('mcp.template.rules.deployType'), trigger: 'change' },
      ],
      imgAddress: [
        { required: true, message: t('mcp.template.rules.imgAddress'), trigger: 'blur' },
      ],
      mcpServers: [
        {
          required: true,
          validator: (rule: any, value: string, callback: (error?: string | Error) => void) => {
            let parsed
            if (!value) return callback(new Error(t('mcp.template.rules.mcpServers.must')))
            try {
              parsed = JSON.parse(value)
            } catch {
              // Capture JSON parsing errors and return custom prompts
              return callback(new Error(t('mcp.template.rules.mcpServers.format')))
            }
            const regex = /^[A-Za-z_-][A-Za-z0-9_-]*$/
            // Get server name
            const serverName = Object.keys(parsed.mcpServers)[0]
            // const formatted = JSON.stringify(parsed, null, 2)
            if (!parsed.mcpServers)
              return callback(new Error(t('mcp.template.rules.mcpServers.name')))
            if (!serverName)
              return callback(new Error(t('mcp.template.rules.mcpServers.serverName')))
            if (!regex.test(serverName)) {
              return callback(new Error(t('mcp.template.rules.mcpServers.regexServerName')))
            }

            // 1.Verification when the current deployment mode is SSE or steamableHttp
            if (
              [AccessType.DIRECT, AccessType.PROXY].includes(pageInfo.value.formData.accessType) &&
              [McpProtocol.SSE, McpProtocol.STREAMABLE_HTTP].includes(
                pageInfo.value.formData.mcpProtocol,
              )
            ) {
              if (!parsed.mcpServers[serverName].url)
                return callback(new Error(t('mcp.template.rules.mcpServers.url')))
              if (parsed.mcpServers[serverName].type) {
                if (!['sse', 'streamable-http'].includes(parsed.mcpServers[serverName].type)) {
                  return callback(new Error(t('mcp.template.rules.mcpServers.type')))
                }
              }
              if (parsed.mcpServers[serverName].transport) {
                if (!['sse', 'streamable-http'].includes(parsed.mcpServers[serverName].transport)) {
                  return callback(new Error(t('mcp.template.rules.mcpServers.transport')))
                }
              }
            }
            // 2.The current protocol is STDIO
            if (showCommand.value) {
              if (!parsed.mcpServers[serverName].command) {
                return callback(new Error(t('mcp.template.rules.mcpServers.command')))
              }
            }

            callback()
          },
          trigger: 'blur',
        },
      ],
      environmentId: [
        { required: true, message: t('mcp.template.rules.environmentId'), trigger: 'change' },
      ],
    },
    tooltip: {
      imgAddress: InstanceData.value.TIP_IMGADDRESS + InstanceData.value.TIP_IMGADDRESS_DEFAULT,
    },
  })
  const { pvcList } = toRefs(useMcpStoreHook())
  const exampleList = [
    {
      name: 'python-mcp-sys-monitor',
      language: 'Python',
      description: 'Python 版本的 MCP Server 示例代码',
    },
    {
      name: 'nodejs-mcp-sys-monitor',
      language: 'Node.js',
      description: 'Node.js 版本的 MCP Server 示例代码',
    },
    {
      name: 'binary-mcp-sys-monitor',
      language: '二进制',
      description: '二进制版本的 MCP Server 示例代码',
    },
  ]
  /**
   * mcpServers placeholder
   */
  const placeholderServer = computed(() => {
    return t('mcp.instance.formData.mcpServersPlaceholder') + InstanceData.value.TIP_MCP_SERVER
  })

  /**
   * condition of show imgAddress
   */
  const showImgAddress = computed(() => {
    return Number(pageInfo.value.formData.accessType) === AccessType.HOSTING
  })

  /**
   * condition of show mcpServers
   */
  const showMcpServers = computed(() => {
    return !(
      pageInfo.value.formData.accessType === AccessType.HOSTING &&
      (pageInfo.value.formData.mcpProtocol === McpProtocol.SSE ||
        pageInfo.value.formData.mcpProtocol === McpProtocol.STREAMABLE_HTTP)
    )
  })

  /**
   * condition of show command
   */
  const showCommand = computed(() => {
    return (
      pageInfo.value.formData.accessType === AccessType.HOSTING &&
      pageInfo.value.formData.mcpProtocol === McpProtocol.STDIO
    )
  })

  /**
   * condition of show Server Path
   */
  const showServicePath = computed(() => {
    return (
      pageInfo.value.formData.accessType === AccessType.HOSTING &&
      (pageInfo.value.formData.mcpProtocol === McpProtocol.SSE ||
        pageInfo.value.formData.mcpProtocol === McpProtocol.STREAMABLE_HTTP)
    )
  })

  /**
   * PVC node inaccessible condition judgment
   */
  const disabledPvcNode = computed(() => {
    return (pvc: PvcResult) =>
      pvc.accessModes?.includes('ReadWriteOnce') && pvc.pods && pvc.pods.length > 0
  })

  /**
   * selected of pvc
   */
  const selectedPvc = computed(() => {
    return (pvcName: string) => pvcList.value?.find((pvc: PvcResult) => pvc.name === pvcName) || {}
  })

  /**
   * The read-only attribute of PVC mode cannot be modified
   */
  const disabledReadOnly = computed(() => {
    return (pvcName: string) =>
      pvcList.value
        ?.find((pvc: PvcResult) => pvc.name === pvcName)
        ?.accessModes?.includes(NodeVisible.ROM)
  })
  return {
    query,
    router,
    meta,
    path,
    jumpBack,
    jumpToPage,
    userInfo,
    pageInfo,
    originForm,
    placeholderServer,
    showImgAddress,
    showMcpServers,
    showCommand,
    showServicePath,
    disabledPvcNode,
    selectedPvc,
    disabledReadOnly,
    selectVisible,
    currentMCP,
    exampleList,
  }
}
