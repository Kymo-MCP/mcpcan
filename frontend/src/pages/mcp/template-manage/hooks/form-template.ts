import { useRouterHooks } from '@/utils/url'
import { useMcpStoreHook } from '@/stores'
import { type PvcResult, AccessType, McpProtocol, InstanceData, NodeVisible } from '@/types/index'

export const useTemplateFormHooks = () => {
  const { t } = useI18n()
  const { jumpBack, jumpToPage } = useRouterHooks()
  const router = useRouter()
  const { query, meta } = useRoute()
  const selectVisible = ref(false)
  const originForm = ref<any>()
  const pageInfo = ref<any>({
    visible: false,
    loading: false,
    title: t('mcp.template.formData.title'),
    formData: {
      templateId: '',
      name: '',
      accessType: '',
      mcpProtocol: '',
      imgAddress: InstanceData.value.IMGADDRESS,
      mcpServers: '',
      notes: '',
      iconPath: '',
      packageId: '',
      environmentId: '',
      port: InstanceData.value.PORT,
      servicePath: '',
      environmentVariables: [],
      volumeMounts: [],
      initScript: InstanceData.value.INITSCRIPT,
      command: '',
    },
    rules: {
      name: [{ required: true, message: t('mcp.template.rules.name'), trigger: 'blur' }],
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
            } catch (error) {
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
    return (pvcName: string) => pvcList.value?.find((pvc: PvcResult) => pvc.name === pvcName) || []
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
    jumpBack,
    jumpToPage,
    exampleList,
    router,
    query,
    meta,
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
  }
}
