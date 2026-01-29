<template>
  <div>
    <el-dialog
      v-model="dialogInfo.visible"
      :title="t('mcp.instance.accessType.title')"
      append-to-body
      width="1200px"
      class="access-type-dialog"
    >
      <div class="flex justify-center gap-8 py-8">
        <div
          v-for="item in accessOptions"
          :key="item.value"
          class="access-card flex flex-col items-center p-6 cursor-pointer transition-all duration-300 rounded-xl border bg-card"
          @click="handleSelect(item)"
          @mouseenter="currentHover = item.value"
          @mouseleave="currentHover = null"
        >
          <GradientText
            v-if="currentHover === item.value"
            text=""
            :colors="['#ffaa40', '#9c40ff', '#ffaa40']"
            :animation-speed="10"
            :show-border="false"
          >
            <div class="w-full h-full flex flex-col items-center">
              <el-icon :size="48" class="mb-4 !text-[48px]" :color="item.color">
                <i class="icon iconfont !text-[48px]" :class="item.icon"></i>
              </el-icon>
              <div class="text-lg font-bold" :class="item.subLabel ? '' : 'mb-2'">
                {{ item.label }}
              </div>
              <div class="text-size-sm mb-1" v-if="item.subLabel">({{ item?.subLabel }})</div>
              <div class="tip tip-primary desc">
                {{ item.description }}
              </div>
              <div class="mt-4 text-size-sm w-full" v-if="item.supportTypes">
                <div class="font-bold text-left">
                  {{ t('mcp.instance.accessType.protocolType') }}
                </div>
                <div v-if="item.supportTypes.includes(McpProtocol.SSE)" class="item-control">
                  {{ t('mcp.instance.accessType.SSE') }}
                </div>
                <div
                  v-if="item.supportTypes.includes(McpProtocol.STREAMABLE_HTTP)"
                  class="item-control"
                >
                  {{ t('mcp.instance.accessType.STREAMABLE_HTTP') }}
                </div>
                <div v-if="item.supportTypes.includes(McpProtocol.STDIO)" class="item-control">
                  {{ t('mcp.instance.accessType.STDIO') }}
                </div>
              </div>
            </div>
          </GradientText>
          <div v-else class="w-full h-full flex flex-col items-center">
            <el-icon :size="48" class="mb-4 !text-[48px]" :color="item.color">
              <i class="icon iconfont !text-[48px]" :class="item.icon"></i>
            </el-icon>
            <div class="text-lg font-bold" :class="item.subLabel ? '' : 'mb-2'">
              {{ item.label }}
            </div>
            <div class="text-size-sm mb-1" v-if="item.subLabel">({{ item?.subLabel }})</div>
            <div class="tip tip-primary desc">
              {{ item.description }}
            </div>
            <div class="mt-4 text-size-sm w-full" v-if="item.supportTypes">
              <div class="font-bold text-left">{{ t('mcp.instance.accessType.protocolType') }}</div>
              <div v-if="item.supportTypes.includes(McpProtocol.SSE)" class="item-control">
                {{ t('mcp.instance.accessType.SSE') }}
              </div>
              <div
                v-if="item.supportTypes.includes(McpProtocol.STREAMABLE_HTTP)"
                class="item-control"
              >
                {{ t('mcp.instance.accessType.STREAMABLE_HTTP') }}
              </div>
              <div v-if="item.supportTypes.includes(McpProtocol.STDIO)" class="item-control">
                {{ t('mcp.instance.accessType.STDIO') }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>
    <component
      ref="formComponent"
      :is="currentModal?.formComponent"
      @on-refresh="emit('on-refresh')"
    ></component>
  </div>
</template>

<script setup lang="ts">
import GradientText from '@/components/Animation/GradientText.vue'
import { useInstanceFormHooks } from '../hooks/form-instance.ts'
import { AccessType, McpProtocol } from '@/types/instance.ts'
import HostingDialog from './hosting-dialog.vue'
import ProxyDialog from './proxy-dialog.vue'
import DirectDialog from './direct-dialog.vue'
import OpenApiDialog from './open-api-dialog.vue'
import { type InstanceResult } from '@/types/instance.ts'

const emit = defineEmits(['select', 'on-refresh'])
const { t } = useI18n()
const { pageInfo, jumpToPage, path } = useInstanceFormHooks()
const dialogInfo = ref({
  visible: false,
})
const currentHover = ref<AccessType | null>()
const formComponent = ref()
const accessOptions = [
  {
    label: t('mcp.instance.accessType.hosting'),
    value: AccessType.HOSTING,
    icon: 'MCP-anquan',
    formComponent: HostingDialog,
    color: '#67C23A',
    supportTypes: [McpProtocol.SSE, McpProtocol.STREAMABLE_HTTP, McpProtocol.STDIO],
    description: t('mcp.instance.accessType.hostingDesc'),
  },
  {
    label: t('mcp.instance.accessType.proxy'),
    value: AccessType.PROXY,
    icon: 'MCP-daili',
    formComponent: ProxyDialog,
    color: '#E6A23C',
    supportTypes: [McpProtocol.SSE, McpProtocol.STREAMABLE_HTTP],
    description: t('mcp.instance.accessType.proxyDesc'),
  },
  {
    label: t('mcp.instance.accessType.direct'),
    value: AccessType.DIRECT,
    icon: 'MCP-zhilian',
    formComponent: DirectDialog,
    color: '#409EFF',
    supportTypes: [McpProtocol.SSE, McpProtocol.STREAMABLE_HTTP],
    description: t('mcp.instance.accessType.directDesc'),
  },
  {
    label: t('mcp.instance.accessType.api'),
    subLabel: t('mcp.instance.accessType.supportApi'),
    value: 4,
    icon: 'MCP-wenjian1',
    formComponent: OpenApiDialog,
    color: '#ff8eb9',
    // supportTypes: [],
    description: t('mcp.instance.accessType.openAPIDesc'),
  },
]
const currentModal = computed(() => {
  return accessOptions.find((option) => option.value === pageInfo.value.accessType)
})

const handleSelect = async (item: any) => {
  pageInfo.value.accessType = item.value
  dialogInfo.value.visible = false
  // dialog view
  if (item.value === 4) {
    nextTick(() => {
      formComponent.value?.init(null, path === '/template-manage' ? 'template' : 'instance')
    })
    return
  }
  if (path === '/template-manage') {
    jumpToPage({
      url: '/new-template',
      data: {
        type: item.value,
      },
    })
    return
  }
  // page view with create
  jumpToPage({
    url: '/new-instance',
    data: {
      type: item.value,
    },
  })
}

const init = (instance: InstanceResult | null) => {
  if (instance) {
    pageInfo.value.accessType = instance.accessType
    if (instance.instanceId) {
      nextTick(() => {
        formComponent.value?.init(instance)
      })
      return
    }
  }
  dialogInfo.value.visible = true
}
defineExpose({
  init,
})
</script>

<style lang="scss" scoped>
.access-card {
  width: 240px;
  // height: 280px;
  border-color: var(--ep-border-color-lighter);
  color: var(--ep-text-color-secondary);
  font-family: 'Monaco, Menlo, "Ubuntu Mono", monospace';
  font-size: 12px;
  word-break: normal;
  &:hover {
    transform: scale(1.1);
    box-shadow: 0 12px 32px rgba(0, 0, 0, 0.1);
    border: 1px solid var(--el-color-primary);
    z-index: 10;
  }
}
.bg-card {
  background-color: var(--ep-home-glass);
}
.tip {
  padding: 10px;
  border-radius: 4px;
  font-size: 12px !important;
  &.tip-warning {
    background-color: #fff1f0;
    border-left: 5px solid var(--el-color-danger);
  }
  &.tip-primary {
    background-color: #409eff1a;
    border-left: 5px solid var(--el-color-primary);
  }
}
.desc {
  text-align: justify;
  text-justify: inter-word;
  word-break: normal;
  overflow-wrap: break-word;
}
.item-control {
  font-size: 12px;
  padding-left: 12px;
  &::before {
    content: '• ';
    color: var(--el-color-primary);
  }
}
</style>
