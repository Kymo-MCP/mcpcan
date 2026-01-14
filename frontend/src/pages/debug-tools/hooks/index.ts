import { type InstanceResult, InstanceStatus, ContainerOptions } from '@/types/instance'
import { useMcpStore } from '@/stores/modules/mcp-store'

export const useDebugToolsHooks = () => {
  const { t } = useI18n()
  // const instanceInfo = ref<any>({})
  const toolList = ref<any[]>([])
  const currentTool = ref<any>(null)
  const keyword = ref('')
  const inputJson = ref('{}')
  const outputResult = ref<string>('')
  const history = ref<any[]>([])
  const route = useRoute()
  const loading = ref(false)
  const running = ref(false)
  const instanceId = computed(() => route.query.instanceId as string)
  const { currentInstance: instanceInfo } = toRefs(useMcpStore())

  // service status
  const activeOptions = {
    active: {
      label: t('status.' + InstanceStatus.ACTIVE),
      type: 'success',
      value: InstanceStatus.ACTIVE,
    },
    inactive: {
      label: t('status.' + InstanceStatus.INACTIVE),
      type: 'danger',
      value: InstanceStatus.INACTIVE,
    },
  }
  return {
    activeOptions,
    instanceInfo,
    toolList,
    currentTool,
    keyword,
    inputJson,
    outputResult,
    history,
    route,
    loading,
    running,
    instanceId,
    t,
  }
}
