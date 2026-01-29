<template>
  <div>
    <el-dialog
      v-model="dialogInfo.visible"
      :title="'直连模式'"
      :show-close="false"
      :close-on-click-modal="false"
      class="access-type-dialog"
      width="600px"
      top="10vh"
      append-to-body
    >
      <el-scrollbar height="70vh">
        <DirectForm ref="directForm"></DirectForm>
      </el-scrollbar>
      <template #footer>
        <div class="flex justify-center">
          <mcp-button @click="handleConfirm" class="mr-4"> 保存并运行 </mcp-button>
          <mcp-button plain @click="handleSaveAsTemplate" class="mr-4"> 另存为模板 </mcp-button>
          <el-button @click="handleClose">退出</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import McpButton from '@/components/mcp-button/index.vue'
import DirectForm from './components/direct-form.vue'

const directForm = ref()
const dialogInfo = ref({
  visible: false,
})

const handleConfirm = () => {
  directForm.value?.handleConfirm()
}

const handleSaveAsTemplate = () => {
  directForm.value?.handleSaveAsTemplate()
}

const handleClose = () => {
  dialogInfo.value.visible = false
}

const init = () => {
  dialogInfo.value.visible = true
  nextTick(() => {
    directForm.value?.init()
  })
}
defineExpose({
  init,
})
</script>
