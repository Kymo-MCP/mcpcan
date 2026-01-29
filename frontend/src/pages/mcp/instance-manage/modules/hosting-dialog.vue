<template>
  <div>
    <el-dialog
      v-model="dialogInfo.visible"
      :title="'托管模式'"
      :show-close="false"
      :close-on-click-modal="false"
      class="access-type-dialog"
      width="850px"
      top="10vh"
      destroy-on-close
      append-to-body
    >
      <el-scrollbar height="70vh">
        <HostForm ref="hostForm"></HostForm>
      </el-scrollbar>
      <template #footer>
        <div
          :class="
            dialogInfo.instanceInfo.formData?.instanceId
              ? 'flex justify-between items-center'
              : 'text-center'
          "
        >
          <div v-if="dialogInfo.instanceInfo.formData?.instanceId" class="flex">
            <el-button link type="primary" @click="handleConfig"> 访问配置 </el-button>
            <el-divider direction="vertical" class="!h-4 !my-auto" />
            <el-button link type="warning" @click="handleViewStatus"> 状态探测 </el-button>
            <el-divider direction="vertical" class="!h-4 !my-auto" />
            <el-button link type="success" @click="handleViewLog"> 查看日志 </el-button>
          </div>
          <div class="flex justify-center">
            <mcp-button @click="handleConfirm" class="mr-4"> 保存并运行 </mcp-button>
            <mcp-button plain @click="handleSaveAsTemplate" class="mr-4"> 另存为模板 </mcp-button>
            <el-button @click="handleClose">退出</el-button>
          </div>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import HostForm from './components/host-form.vue'
import McpButton from '@/components/mcp-button/index.vue'

const hostForm = ref()
const dialogInfo = ref<any>({
  visible: false,
  instanceInfo: {},
})

const handleConfig = () => {
  hostForm.value.handleConfig()
}
const handleClose = () => {
  dialogInfo.value.visible = false
}
const handleViewStatus = () => {
  hostForm.value.handleViewStatus()
}
const handleViewLog = () => {
  hostForm.value.handleViewLog()
}
const handleConfirm = () => {
  hostForm.value.handleConfirm()
}
const handleSaveAsTemplate = () => {
  hostForm.value.handleSaveAsTemplate()
}
const init = (instance: any) => {
  dialogInfo.value.instanceInfo = instance || {}
  dialogInfo.value.visible = true
  nextTick(() => {
    hostForm.value.init(instance)
  })
}
defineExpose({
  init,
})
</script>
