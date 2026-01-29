<template>
  <div class="px-5 py-2" v-loading="pageInfo.loading">
    <el-form
      ref="baseInfo"
      :model="pageInfo.formData"
      :rules="pageInfo.rules"
      label-width="auto"
      label-position="left"
    >
      <el-row :gutter="24">
        <el-col :span="18">
          <el-form-item :label="t('mcp.instance.formData.instanceName')" prop="name">
            <el-input
              v-model="pageInfo.formData.name"
              :placeholder="t('mcp.instance.formData.instanceName')"
            />
          </el-form-item>
          <el-form-item :label="t('mcp.template.formData.notes')" prop="notes">
            <el-input
              v-model="pageInfo.formData.notes"
              :rows="2"
              type="textarea"
              :placeholder="t('mcp.template.formData.notes')"
            />
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item prop="iconPath">
            <Upload v-model="pageInfo.formData.iconPath"></Upload>
          </el-form-item>
        </el-col>
      </el-row>
      <el-form-item :label="t('mcp.instance.form.mcpProtocol')" prop="mcpProtocol">
        <el-segmented
          v-model="pageInfo.formData.mcpProtocol"
          :options="protocolOptions"
          @change="handleMcpProtocolChange"
        />
      </el-form-item>
      <el-row v-if="pageInfo.formData.mcpProtocol === 3">
        <el-col :span="18">
          <el-form-item prop="mcpServers">
            <template #label>{{ t('mcp.instance.formData.mcpServers') }}</template>
            <MonacoEditor v-model="pageInfo.formData.mcpServers" language="json" height="200px" />
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <div
            class="pl-3 rounded border border-[var(--ep-border-color-lighter)] text-[var(--ep-text-color-secondary)] text-xs leading-6 tracking-wide"
            v-html="mcpServersTips"
          ></div>
        </el-col>
      </el-row>
      <el-form-item prop="packageId">
        <template #label>
          <span class="mr-2">{{ t('mcp.instance.formData.packageId') }}</span>
        </template>
        <div class="flex items-center">
          <el-tooltip
            :disabled="!isPackageNameOverflow"
            :content="currentPackage.name"
            placement="top"
          >
            <div
              v-if="pageInfo.formData.packageId"
              class="mr-2 max-w-60 flex items-center"
              ref="packageNameRef"
            >
              <el-image
                :src="zipLogo"
                style="width: 16px; height: 16px"
                class="mr-1 shrink-0"
              ></el-image>
              <span class="u-line-1" @mouseenter="checkPackageNameOverflow">
                {{ currentPackage.name }}
              </span>
            </div>
          </el-tooltip>
          <mcp-button size="small" class="base-btn" @click="handleSelectPackage">
            {{ t('common.select') }}
          </mcp-button>
          <el-tooltip class="box-item" effect="light" placement="top-start">
            <el-button class="base-btn-link" link @click="downloadDialogVisible = true">
              {{ t('mcp.instance.hostingForm.codeExample') }}
            </el-button>
            <template #content>
              <div class="w-40" v-html="codeTips"></div>
            </template>
          </el-tooltip>
        </div>
      </el-form-item>

      <!-- 命令 -->
      <el-row :gutter="24">
        <el-col :span="16">
          <el-form-item :label="t('mcp.instance.formData.initScript')" prop="initScript">
            <el-input
              v-model="pageInfo.formData.initScript"
              :rows="6"
              type="textarea"
              :placeholder="t('mcp.instance.formData.initScript')"
            />
          </el-form-item>
          <el-form-item :label="t('mcp.instance.formData.command')" prop="command">
            <el-input
              v-model="pageInfo.formData.command"
              :rows="6"
              type="textarea"
              :placeholder="
                showCommand ? InstanceData.COMMAND_TIP : t('mcp.instance.formData.command')
              "
              :disabled="showCommand"
            />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <div class="template-list-container">
            <div class="text-sm font-bold mb-2">
              {{ t('mcp.instance.hostingForm.commandExample') }}
            </div>
            <el-popover
              :visible="!!selectedCommand.commands"
              placement="right"
              :title="selectedCommand?.name"
              :show-arrow="false"
              popper-style="width: 340px;"
            >
              <template #reference>
                <div>
                  <el-scrollbar height="175px" always>
                    <div class="flex flex-col gap-2 pr-2 position-relative">
                      <div
                        v-for="(item, index) in deployTemplateData.configurations"
                        :key="index"
                        class="flex items-center gap-3 rounded-md mx-1 p-1 text-left type-item"
                        :class="{ 'active-type': selectedCommand?.name === item.name }"
                        @click="handleSelectCommand(item)"
                      >
                        <div class="font-bold truncate" :title="item.name">
                          {{ item.name }}
                        </div>
                      </div>
                    </div>
                  </el-scrollbar>
                  <div class="tip tip-primary mt-2" v-html="commandTips"></div>
                </div>
              </template>
              <div v-if="selectedCommand" class="p-4 w-[300px]">
                <div class="mb-4">
                  <div class="flex justify-between items-center mb-1">
                    <span class="font-bold text-sm">{{
                      t('mcp.instance.hostingForm.dependenceExample')
                    }}</span>
                    <el-button
                      link
                      type="primary"
                      @click="handleUseCommand(selectedCommand?.commands?.install, 'initScript')"
                    >
                      {{ t('common.use') }}
                    </el-button>
                  </div>
                  <div class="p-2 rounded text-xs text-gray-500 break-all">
                    {{ selectedCommand?.commands?.install }}
                  </div>
                </div>
                <div v-if="!showCommand" class="mb-2">
                  <div class="flex justify-between items-center mb-1">
                    <span class="font-bold text-sm">{{
                      t('mcp.instance.hostingForm.startExample')
                    }}</span>
                    <el-button
                      link
                      type="primary"
                      @click="handleUseCommand(selectedCommand.commands.start, 'command')"
                    >
                      {{ t('common.use') }}
                    </el-button>
                  </div>
                  <div class="p-2 rounded text-xs text-gray-500 break-all">
                    {{ selectedCommand?.commands?.start }}
                  </div>
                </div>
                <div class="tip tip-primary mb-2">
                  {{
                    locale === 'en' ? selectedCommand?.description_en : selectedCommand?.description
                  }}
                </div>
                <div class="text-xs text-orange-500">
                  {{ t('mcp.instance.hostingForm.commandTips') }}
                </div>
              </div>
            </el-popover>
          </div>
        </el-col>
      </el-row>
      <div class="tip tip-primary mt-2" v-html="commandDesc"></div>
      <div class="font-size-3 mt-4">
        <div class="tip tip-primary mt-2" v-html="buildNotes"></div>
      </div>

      <el-form-item :label="t('mcp.instance.hostingForm.accessUrl')" class="mt-6">
        <el-radio-group>
          <el-radio-button value="ip">
            <el-popover
              :title="t('mcp.instance.hostingForm.listenUrl')"
              width="260"
              placement="bottom-start"
            >
              <template #reference>
                <div class="w-24">0.0.0.0</div>
              </template>
              <div>
                {{ t('mcp.instance.hostingForm.listenUrlTips1') }}
              </div>
            </el-popover>
          </el-radio-button>
          <el-radio-button
            value="port"
            :class="pageInfo.formData.mcpProtocol !== 3 ? 'deep-form' : ''"
          >
            <el-popover
              :title="t('mcp.instance.hostingForm.listenPort')"
              width="260"
              placement="bottom-start"
            >
              <template #reference>
                <div v-if="pageInfo.formData.mcpProtocol === 3" class="w-24 flex items-center">
                  <div class="w-24">8080</div>
                </div>
                <div v-else class="w-30 flex items-center">
                  <el-input
                    v-model.number="pageInfo.formData.port"
                    type="number"
                    :placeholder="t('mcp.instance.formData.port')"
                    class="no-border-input w-24"
                  />
                </div>
              </template>
              <div>
                {{
                  pageInfo.formData.mcpProtocol === 3
                    ? t('mcp.instance.hostingForm.listenPortTips1')
                    : t('mcp.instance.hostingForm.listenPortTips2')
                }}
              </div>
            </el-popover>
          </el-radio-button>
          <el-radio-button
            value="path"
            :class="pageInfo.formData.mcpProtocol !== 3 ? 'deep-form' : ''"
          >
            <el-popover
              :title="t('mcp.instance.hostingForm.listenPath')"
              width="260"
              placement="bottom-start"
            >
              <template #reference>
                <!-- STDIO 协议 -->
                <div v-if="pageInfo.formData.mcpProtocol === 3" class="w-24 flex items-center">
                  <div class="flex-1">/sse</div>
                </div>
                <div v-else class="w-30 flex items-center">
                  <el-input
                    v-model="pageInfo.formData.servicePath"
                    :placeholder="t('mcp.template.formData.servicePath')"
                    class="no-border-input w-24"
                  />
                </div>
              </template>
              <div>
                {{
                  pageInfo.formData.mcpProtocol === 3
                    ? t('mcp.instance.hostingForm.listenPathTips1')
                    : t('mcp.instance.hostingForm.listenPathTips2')
                }}
              </div>
            </el-popover>
          </el-radio-button>
        </el-radio-group>
        <el-radio-group v-if="pageInfo.formData.mcpProtocol === 3" class="mt-4">
          <el-radio-button value="ip">
            <el-popover
              :title="t('mcp.instance.hostingForm.listenUrl')"
              width="260"
              placement="bottom-start"
            >
              <template #reference>
                <div class="w-24">0.0.0.0</div>
              </template>
              <div>
                {{ t('mcp.instance.hostingForm.listenUrlTips1') }}
              </div>
            </el-popover>
          </el-radio-button>
          <el-radio-button value="port">
            <el-popover
              :title="t('mcp.instance.hostingForm.listenPort')"
              width="260"
              placement="bottom-start"
            >
              <template #reference>
                <div class="w-24 flex items-center">
                  <div class="w-24">8080</div>
                </div>
              </template>
              <div>
                {{ t('mcp.instance.hostingForm.listenPortTips1') }}
              </div>
            </el-popover>
          </el-radio-button>
          <el-radio-button value="path">
            <el-popover
              :title="t('mcp.instance.hostingForm.listenPath')"
              width="260"
              placement="bottom-start"
            >
              <template #reference>
                <div class="w-24 flex items-center" @click="handleChangePath">
                  <div class="flex-1">/mcp</div>
                </div>
              </template>
              <div>
                {{ t('mcp.instance.hostingForm.listenPathTips3') }}
              </div>
            </el-popover>
          </el-radio-button>
        </el-radio-group>
      </el-form-item>
      <!-- 环境默认 -->
      <el-form-item
        v-show="false"
        :label="t('mcp.instance.formData.environmentId')"
        prop="environmentId"
      >
        <el-select
          v-model="pageInfo.formData.environmentId"
          :placeholder="t('mcp.instance.formData.environmentId')"
          @change="handleChangeEnvironmentId"
          disabled
        >
          <el-option
            v-for="(env, index) in envList"
            :key="index"
            :label="env.name"
            :value="env.id"
          />
        </el-select>
      </el-form-item>
      <el-collapse v-model="currentOpenCollapse" :expand-icon-position="'left'" accordion>
        <el-collapse-item name="1">
          <template #title>
            <div>
              <span class="mr-1 font-bold">{{
                t('mcp.instance.formData.environmentVariables')
              }}</span>
              <span
                class="rounded border border-[var(--ep-border-color-lighter)] text-[var(--ep-text-color-secondary)] text-xs leading-6 tracking-wide"
              >
                {{ t('mcp.instance.hostingForm.envTips') }}
              </span>
            </div>
          </template>
          <el-form-item prop="environmentVariables">
            <!-- :label="t('mcp.instance.formData.environmentVariables')" -->
            <div class="w-full mt-1 ml-1">
              <div
                class="flex align-center mb-2"
                v-for="(env, index) in pageInfo.formData.environmentVariables"
                :key="index"
              >
                <el-input
                  v-model="env.key"
                  :placeholder="t('mcp.instance.formData.key')"
                  class="mr-2"
                />
                <el-input
                  v-model="env.value"
                  :placeholder="t('mcp.instance.formData.value')"
                  class="mr-2"
                />
                <el-icon
                  class="cursor-pointer"
                  color="#F56C6C"
                  @click="handleDeleteEnvVariable(index)"
                >
                  <Remove />
                </el-icon>
              </div>
            </div>
            <el-button class="add-env" :icon="Plus" plain @click="handleAddEnvVariable">
              {{ t('mcp.instance.formData.add-env') }}
            </el-button>
          </el-form-item>
          <div class="tip tip-primary">
            {{ t('mcp.instance.hostingForm.envNotes') }}
          </div>
        </el-collapse-item>
        <el-collapse-item v-if="pageInfo.formData.environmentId" name="2">
          <template #title>
            <div>
              <span class="mr-1 font-bold">{{ t('mcp.instance.hostingForm.volume') }}</span>
              <span
                class="rounded border border-[var(--ep-border-color-lighter)] text-[var(--ep-text-color-secondary)] text-xs leading-6 tracking-wide"
              >
                {{ t('mcp.instance.hostingForm.volumeTips') }}
              </span>
            </div>
          </template>
          <el-form-item prop="volumeConfig">
            <el-card
              class="w-full mb-2"
              v-for="(volume, index) in pageInfo.formData.volumeMounts"
              :key="index"
            >
              <div>
                <div class="flex align-center justify-between">
                  <span>
                    <span class="text-red">*</span>{{ t('mcp.instance.formData.volumeType') }}
                  </span>
                  <el-icon
                    color="#F56C6C"
                    class="cursor-pointer"
                    @click="handleDeleteVolume(index)"
                  >
                    <CircleClose />
                  </el-icon>
                </div>
                <el-select
                  v-model="volume.type"
                  :placeholder="t('mcp.instance.formData.volumeType')"
                  style="width: 200px"
                >
                  <el-option label="HostPath" value="hostPath" />
                  <el-option
                    v-if="currentEnvironment.environment === 'kubernetes'"
                    label="PVC"
                    value="pvc"
                  />
                  <el-option
                    v-if="currentEnvironment.environment === 'docker'"
                    label="Volume"
                    value="volume"
                  />
                </el-select>
              </div>
              <template v-if="currentEnvironment.environment === 'kubernetes'">
                <template v-if="volume.type === 'hostPath'">
                  <div class="flex align-end">
                    <div>
                      <div>
                        <span class="text-red">*</span>{{ t('mcp.instance.formData.nodeName') }}
                      </div>
                      <el-select
                        v-model="volume.nodeName"
                        :placeholder="t('mcp.instance.formData.nodeName')"
                        style="width: 200px"
                      >
                        <el-option
                          v-for="(node, index) in nodeList"
                          :key="index"
                          :label="node.name"
                          :value="node.name"
                        />
                      </el-select>
                    </div>
                    <el-button class="mr-2 ml-2" :icon="RefreshRight" />
                    <div>
                      <div>
                        <span class="text-red">*</span>{{ t('mcp.instance.formData.hostPath') }}
                      </div>
                      <el-input
                        v-model="volume.hostPath"
                        :placeholder="t('mcp.instance.formData.hostPath')"
                      />
                    </div>
                  </div>
                  <div class="flex align-end">
                    <div class="mr-2">
                      <div>
                        <span class="text-red">*</span>{{ t('mcp.instance.formData.mountPath') }}
                      </div>
                      <el-input
                        v-model="volume.mountPath"
                        :placeholder="t('mcp.instance.formData.mountPath')"
                        style="width: 200px"
                      />
                    </div>
                    <div>
                      <div>{{ t('mcp.instance.formData.readOnly') }}</div>
                      <el-switch v-model="volume.readOnly" />
                    </div>
                  </div>
                </template>

                <template v-if="volume.type === 'pvc'">
                  <div class="flex align-end">
                    <div>
                      <div>
                        <span class="text-red">*</span>{{ t('mcp.instance.formData.pvcName') }}
                      </div>
                      <el-select
                        v-model="volume.pvcName"
                        placeholder="PVC"
                        style="width: 200px"
                        @change="handlePvcChange($event, volume)"
                      >
                        <el-option
                          v-for="(pvc, index) in pvcList"
                          :key="index"
                          :value="pvc.name"
                          :disabled="disabledPvcNode(pvc)"
                        >
                          <span>{{ pvc.name }}</span>
                          <el-tag v-if="disabledPvcNode(pvc)" color="orange">
                            {{ t('mcp.template.formData.isBind') }}
                          </el-tag>
                        </el-option>
                      </el-select>
                    </div>
                    <el-button class="mr-2 ml-2" :icon="RefreshRight" />
                    <div class="mr-2">
                      <div>
                        <span class="text-red">*</span>{{ t('mcp.instance.formData.mountPath') }}
                      </div>
                      <el-input
                        v-model="volume.mountPath"
                        :placeholder="t('mcp.instance.formData.mountPath')"
                      />
                    </div>
                    <div>
                      <div>{{ t('mcp.instance.formData.readOnly') }}</div>
                      <el-switch
                        v-model="volume.readOnly"
                        :disabled="disabledReadOnly(volume.pvcName)"
                      />
                    </div>
                  </div>
                  <el-space v-if="!selectedPvc(volume.pvcName)">
                    <el-text type="secondary">
                      {{ t('mcp.template.formData.accessModes') }}:
                      {{ selectedPvc(volume.pvcName).accessModes?.join(', ') }}
                    </el-text>
                    <template v-if="selectedPvc(volume.pvcName).pods?.length > 0">
                      <el-text type="secondary">
                        {{ t('mcp.template.formData.pods') }}:
                        {{ selectedPvc(volume.pvcName).pods.join(', ') }}
                      </el-text>
                    </template>
                  </el-space>
                </template>
              </template>
              <template v-if="currentEnvironment.environment === 'docker'">
                <template v-if="volume.type === 'hostPath'">
                  <div class="flex align-end">
                    <div class="mr-2">
                      <div>
                        <span class="text-red">*</span>{{ t('mcp.instance.formData.hostPath') }}
                      </div>
                      <el-input
                        v-model="volume.hostPath"
                        :placeholder="t('mcp.instance.formData.hostPath')"
                      />
                    </div>
                    <div class="mr-2">
                      <div>
                        <span class="text-red">*</span>{{ t('mcp.instance.formData.mountPath') }}
                      </div>
                      <el-input
                        v-model="volume.mountPath"
                        :placeholder="t('mcp.instance.formData.mountPath')"
                        style="width: 200px"
                      />
                    </div>
                    <div>
                      <div>{{ t('mcp.instance.formData.readOnly') }}</div>
                      <el-switch v-model="volume.readOnly" />
                    </div>
                  </div>
                </template>

                <template v-if="volume.type === 'volume'">
                  <div class="flex align-end">
                    <div>
                      <div><span class="text-red">*</span>{{ 'volumeName' }}</div>
                      <el-select
                        v-model="volume.volumeName"
                        placeholder="volumeName"
                        style="width: 200px"
                      >
                        <el-option
                          v-for="(pvc, index) in volumeList"
                          :key="index"
                          :value="pvc.name"
                        >
                          <span>{{ pvc.name }}</span>
                        </el-option>
                      </el-select>
                    </div>
                    <el-button class="mr-2 ml-2" :icon="RefreshRight" />
                    <div class="mr-2">
                      <div>
                        <span class="text-red">*</span>{{ t('mcp.instance.formData.mountPath') }}
                      </div>
                      <el-input
                        v-model="volume.mountPath"
                        :placeholder="t('mcp.instance.formData.mountPath')"
                      />
                    </div>
                    <div>
                      <div>{{ t('mcp.instance.formData.readOnly') }}</div>
                      <el-switch v-model="volume.readOnly" />
                    </div>
                  </div>
                </template>
              </template>
            </el-card>
            <el-button class="add-env" :icon="Plus" plain @click="handleAddVolume">
              {{ t('mcp.instance.formData.addVolume') }}
            </el-button>
          </el-form-item>
        </el-collapse-item>
        <el-collapse-item v-if="!pageInfo.formData.instanceId" name="3">
          <template #title>
            <div>
              <div class="mr-1 font-bold">{{ t('mcp.instance.hostingForm.headerTitle') }}</div>
              <Transition name="tip-fade">
                <div
                  v-show="currentOpenCollapse !== '3'"
                  class="flex-1 tip tip-primary my-2 line-height-[18px]"
                  style="margin-left: -16px"
                  v-html="headerTitleTips"
                ></div>
              </Transition>
            </div>
          </template>
          <div>
            <TokenForm ref="tokenForm" :formData="pageInfo.formData.tokens[0]"></TokenForm>
          </div>
          <Transition name="tip-fade">
            <div
              v-show="currentOpenCollapse === '3'"
              class="flex-1 tip tip-primary my-2 line-height-[18px]"
              v-html="headerContentTips"
            ></div>
          </Transition>
        </el-collapse-item>
      </el-collapse>
    </el-form>
  </div>
  <!-- 代码包选择 -->
  <Select
    v-model="selectVisible"
    v-model:selected="pageInfo.formData.packageId"
    red="packageSelect"
    :title="t('mcp.instance.formData.packageId')"
    :options="packageList"
  >
    <template #options="{ option }">
      <div class="flex justify-between">
        <div class="flex align-center">
          <el-image :src="zipLogo" style="width: 32px; height: 32px"></el-image>
          <span class="ml-2"> {{ option.name }}</span>
        </div>
        <div class="flex align-center">
          {{ t('mcp.template.formData.size') }}：{{ formatFileSize(option.size) }}
          <span class="ml-2">{{ timestampToDate(option.createdAt) }}</span>
          <el-button
            type="primary"
            size="small"
            link
            class="base-btn-link ml-2"
            @click="handleViewCode(option)"
          >
            {{ t('mcp.instance.action.view') }}
          </el-button>
        </div>
      </div>
    </template>
    <template #action>
      <el-upload
        class="mr-8"
        drag
        :action="action"
        :on-success="handleSuccess"
        :headers="headers"
        accept=".zip, .tar, .tar.gz, application/zip, application/x-tar, application/gzip"
        :auto-upload="true"
        :show-file-list="false"
      >
        <el-icon><UploadFilled /></el-icon>
        <div class="ml-2">
          {{ t('mcp.instance.openApi.localFile') }}
        </div>
      </el-upload>
    </template>
  </Select>

  <!-- 下载示例代码包弹窗 -->
  <el-dialog
    v-model="downloadDialogVisible"
    :title="t('mcp.instance.pageDesc.downloadCode')"
    width="520px"
    append-to-body
    :close-on-click-modal="false"
  >
    <div class="flex flex-col gap-4">
      <div
        v-for="(item, index) in exampleList"
        :key="index"
        class="flex items-center justify-between p-4 border rounded-lg hover:shadow-md transition-shadow"
      >
        <div class="flex items-center gap-3">
          <el-image :src="zipLogo" style="width: 32px; height: 32px" />
          <div class="flex flex-col">
            <span class="font-bold text-base">{{ item.name }}</span>
            <span class="text-xs text-gray-500">{{ item.description }}</span>
          </div>
        </div>
        <el-button class="base-btn-link" link @click="handleDownloadExample(item)">
          <el-icon class="mr-1"><Download /></el-icon>
          {{ t('common.download') }}
        </el-button>
      </div>
    </div>
  </el-dialog>
  <!-- probe instance dialog model -->
  <ProbeStatus ref="probe"></ProbeStatus>
  <ConfigDialog ref="config"></ConfigDialog>
  <LogDialog ref="log"></LogDialog>
</template>

<script setup lang="ts">
import {
  Plus,
  Remove,
  CircleClose,
  RefreshRight,
  Download,
  UploadFilled,
} from '@element-plus/icons-vue'
import { useInstanceFormHooks } from '../../hooks/form-instance.ts'
import Upload from '@/components/upload/index.vue'
import McpButton from '@/components/mcp-button/index.vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import TokenForm from './token-form.vue'
import zipLogo from '@/assets/logo/zip.png'
import { AccessType, InstanceData, McpProtocol, NodeVisible } from '@/types/instance'
import {
  type VolumeMountsItme,
  type PvcForm,
  type Code,
  type InstanceResult,
} from '@/types/index.ts'
import { useMcpStoreHook } from '@/stores'
import Select from '@/components/mcp-select/index.vue'
import { formatFileSize, timestampToDate } from '@/utils/system'
import { InstanceAPI } from '@/api/mcp/instance'
import { TemplateAPI } from '@/api/mcp/template'
import { cloneDeep } from 'lodash-es'
import baseConfig from '@/config/base_config.ts'
import { Storage } from '@/utils/storage'
import ProbeStatus from '../probe-dialog.vue'
import ConfigDialog from '../url-config-dialog.vue'
import LogDialog from '../log-dialog.vue'
import deployTemplateData from '@/config/deploy-temlate-data.json'
import MonacoEditor from '@/components/MonacoEditor/index.vue'
import { CodeAPI } from '@/api/code/index'

const { t, locale } = useI18n()
const {
  pageInfo,
  query,
  jumpToPage,
  showCommand,
  disabledPvcNode,
  disabledReadOnly,
  selectedPvc,
  selectVisible,
  exampleList,
} = useInstanceFormHooks()
const { packageList, envList, nodeList, pvcList, volumeList, currentInstance } =
  toRefs(useMcpStoreHook())
const {
  handleGetPackageList,
  handleGetEnvList,
  handleGetNodeList,
  handleGetPvcList,
  handleGetVolumeList,
} = useMcpStoreHook()
const baseInfo = ref()
const currentOpenCollapse = ref()
const protocolOptions = [
  { label: 'STDIO', value: 3 },
  { label: 'SSE', value: 1 },
  { label: 'STREAMABLE_HTTP', value: 2 },
]
const action = ref(
  baseConfig.SERVER_BASE_URL + (window as any).__APP_CONFIG__?.API_BASE + '/market/code/upload',
)
const headers = ref({
  Authorization: `Bearer ${Storage.get('token')}`,
})
const handleSuccess = (response: { code: number; data: { path: string } }) => {
  if (response.code !== 0) {
    return
  }
  handleGetPackageList()
  ElMessage.success(t('action.upload'))
}

const mcpServersTips = computed(() => {
  return locale.value === 'en'
    ? `MCP service SSE/STREAMABLE_HTTP protocol configuration is currently in proxy mode, and the traffic will be forwarded to the MCP configuration provided through the platform gateway.
            After saving, the gateway access configuration will be displayed on the list page. You can also view <a href="#/template-manage">Template List</a> which provides multiple startup examples.`
    : `MCP服务SSE/STEAMABLE_HTTP协议配置当前为代理模式，流量会通过此平台网关转发到此配置提供的
            MCP 配置中，保存后会在列表页显示网关访问配置。也可以查看
            <a href="#/template-manage">部署模板</a> 提供了多个启动示例。`
})
const codeTips = computed(() => {
  return locale.value === 'en'
    ? ` Example MCP service code package, download the code files to understand how to start on this platform in different programming languages. You can also check the
            <a href="#/template-manage"> Template List</a> which provides multiple code package startup examples.`
    : ` 示例MCP服务代码包，下载代码文件可用于了解不同编程语言在此平台的启动方式。也可以查看<a href="#/template-manage"> 部署模板</a>提供了多个代码包启动示例。`
})
const commandTips = computed(() => {
  return locale.value === 'en'
    ? `Note: The platform provides a default container for MCP service startup. You can view the startup example commands by selecting the above runtime environment examples. You can also check the
            <a href="#/template-manage"> Template List</a> which provides multiple startup examples.`
    : ` 说明：平台针对MCP服务启动提供了默认容器，选择以上运行环境示例可以查看启动示例命令。也可以到<a href="#/template-manage">部署模板</a>中查看启动示例。`
})

const commandDesc = computed(() => {
  return locale.value === 'en'
    ? `
      <div class="font-bold mb-2">Execution order of service startup commands inside the MCP container</div>
        <div class="space-y-2">
          <div>
            <span class="font-bold">Code package extraction</span>
            <div class="text-xs text-[var(--ep-text-color-secondary)] mt-1">
              The system will automatically download the code package and extract it to the /app/codepkg/ directory.
            </div>
            <div class="text-xs text-orange-500 mt-1">
              Important: If the archive is named code.zip and contains a top-level folder named code, the final path after extraction will be /app/codepkg/code. When running dependency or start commands afterwards, be sure to run <code>cd /app/codepkg</code> first to enter the correct directory.
            </div>
          </div>
          <div>
            <span class="font-bold">Dependency commands</span>
            <div class="text-xs text-[var(--ep-text-color-secondary)] mt-1">
              Avoid using blocking instructions when running dependency commands. If a command blocks, all subsequent actions will fail to execute.
            </div>
          </div>
          <div>
            <span class="font-bold">Start commands</span>
            <div class="text-xs text-[var(--ep-text-color-secondary)] mt-1">
              Start commands are executed from the system root path by default. If a previous dependency command changed directory into the project (using <code>cd</code>) and did not exit that directory, the start command will inherit that working directory and no additional directory change is required.
            </div>
          </div>
    </div>
  `
    : `
      <div class="font-bold mb-2">MCP 容器内服务启动命令执行顺序</div>
        <div class="space-y-2">
          <div>
            <span class="font-bold">代码包解压</span>
            <div class="text-xs text-[var(--ep-text-color-secondary)] mt-1">
              系统会自动下载代码文件，并解压至 /app/codepkg/ 目录。
            </div>
            <div class="text-xs text-orange-500 mt-1">
              特别注意：若压缩包为 code.zip 且内置顶层文件夹 code，解压后最终路径为
              /app/codepkg/code。后续执行依赖命令与启动命令时，需先执行 cd /app/codepkg
              进入对应目录。
            </div>
          </div>
          <div>
            <span class="font-bold">依赖命令</span>
            <div class="text-xs text-[var(--ep-text-color-secondary)] mt-1">
              执行依赖命令时，需避免出现阻塞型指令。一旦发生命令阻塞，将直接导致后续所有动作无法正常执行。
            </div>
          </div>
          <div>
            <span class="font-bold">启动命令</span>
            <div class="text-xs text-[var(--ep-text-color-secondary)] mt-1">
              启动命令默认在系统根路径下执行。若上一步依赖命令中已执行 cd
              操作进入项目目录，且执行后未退出该目录，启动命令会直接沿用此路径，无需再次执行目录切换操作。
            </div>
          </div>
        </div>
    `
})

const buildNotes = computed(() => {
  return locale.value === 'en'
    ? `
      <span class="font-bold light:text-black dark:text-white">MCP hosting container build notes:</span>
      <div class="mb-2">
        Built from the debian:bookworm-slim image to balance image size and compatibility (full glibc support), and to support common native extensions for Python and Node.js.
      </div>
      <div class="font-bold mb-1">Preinstalled components and versions:</div>
      <ul class="list-disc pl-5 space-y-1">
        <li>
          System utilities: curl, git, wget, tar, zip, unzip and basic build tools (build-essential).
        </li>
        <li>
          Python environment:
          <ul class="list-disc pl-5 mt-1 space-y-1">
            <li>Managed by pyenv with multiple versions preinstalled: Python 3.12.8 (default), 3.11.9, 3.10.14.</li>
            <li>Includes modern package managers: poetry and uv (including uvx), configured with Alibaba Cloud PyPI mirrors.</li>
          </ul>
        </li>
        <li>
          Node.js environment:
          <ul class="list-disc pl-5 mt-1 space-y-1">
            <li>Managed by nvm with Node.js v22 (default) and v20 preinstalled.</li>
            <li>Includes npm, yarn and pnpm, configured with domestic registry mirrors.</li>
          </ul>
        </li>
        ${
          pageInfo.value.formData.mcpProtocol === McpProtocol.STDIO
            ? '<li v-if="pageInfo.formData.mcpProtocol === McpProtocol.STDIO">Default start command: <code>mcp-hosting</code>. After starting, it acts as a gateway service that converts the standard MCP STDIO protocol to HTTP (SSE), enabling remote invocation.</li>'
            : ''
        }

      </ul>
  `
    : `
      <span class="font-bold light:text-black dark:text-white">MCP托管容器构建说明：</span>
      <div class="mb-2">
        基于 debian:bookworm-slim 镜像构建，兼顾镜像体积与兼容性（完整 glibc 支持），适配主流
        Python/Node.js 原生扩展。
      </div>
      <div class="font-bold mb-1">预装组件及版本说明：</div>
      <ul class="list-disc pl-5 space-y-1">
        <li>
          系统基础命令：curl、git、wget、tar、zip、unzip 及基础编译依赖（build-essential）。
        </li>
        <li>
          Python 环境：
          <ul class="list-disc pl-5 mt-1 space-y-1">
            <li>采用 pyenv 进行多版本管理，预装 Python 3.12.8（默认）、3.11.9、3.10.14。</li>
            <li>集成现代包管理工具：poetry、uv（含 uvx），均已配置阿里云 PyPI 加速源。</li>
          </ul>
        </li>
        <li>
          Node.js 环境：
          <ul class="list-disc pl-5 mt-1 space-y-1">
            <li>采用 nvm 进行多版本管理，预装 Node.js v22（默认）、v20。</li>
            <li>集成全套包管理工具：npm、yarn、pnpm，均已配置国内镜像源。</li>
          </ul>
        </li>
        ${
          pageInfo.value.formData.mcpProtocol === McpProtocol.STDIO
            ? '<li>默认启动命令：mcp-hosting。启动后作为网关服务，将标准 MCP STDIO 协议转换为HTTP（SSE）协议，提供远程调用能力。</li>'
            : ''
        }
      </ul>
  `
})

const headerTitleTips = computed(() => {
  return locale.value === 'en'
    ? `If this field is left empty, headers from the client will be passed through to the MCP service by default. If the MCPServers configuration defines header parameters, precedence is: MCPServers configuration headers > client request headers.`
    : `此项不配置时：默认透传来自客户端的Headers到MCP服务中。如果MCPServers配置中存在Header参数，优先级为：MCPServers配置中Header > 客户端的Headers`
})

const headerContentTips = computed(() => {
  return locale.value === 'en'
    ? `
      Notes:
      <br/>
      1. Custom headers do not need to be actively submitted by the client; the gateway will automatically include them in requests forwarded to the MCP service.
      <br/>
      2. If there are duplicate header names among client request headers, MCP configuration headers, and custom pass-through headers, the precedence order is: custom pass-through headers > MCP configuration headers > client request headers. That is, when retrieving values for duplicate header names, the custom pass-through configuration is used first, followed by the MCP configuration, and finally the value provided by the client request.`
    : `注意事项：
      <br/>
      1.自定义 Header 无需客户端主动提交，网关在转发流量时，会自动将其携带至 MCP 服务的请求中。
      <br/>
      2.若客户端请求 Header、MCP 配置 Header、自定义透传 Header 三者存在同名项，优先级顺序为：自定义透传 Header > MCP 配置 Header > 客户端请求 Header。即同名 Header 取值时，优先采用自定义透传配置，其次为 MCP 配置，最后为客户端请求传入的值。`
})

/**
 * Current environment
 */
const currentEnvironment = computed(() => {
  return envList.value?.find((item: any) => item.id === pageInfo.value.formData.environmentId) || {}
})
/**
 * The current selected code package
 */
const currentPackage = computed(() => {
  return (
    packageList.value?.find(
      (item: { id: string }) => item.id === pageInfo.value.formData.packageId,
    ) || { name: '', size: '', createdAt: '' }
  )
})
const handleSelectPackage = () => {
  handleGetPackageList()
  selectVisible.value = true
}
/**
 * Handle view code package
 */
const handleViewCode = (code: Code) => {
  jumpToPage({
    url: '/view-code-package',
    data: {
      id: code.id,
      name: code.name,
    },
    isOpen: true,
  })
}
/**
 * Handle add an environment variables
 */
const handleAddEnvVariable = () => {
  pageInfo.value.formData.environmentVariables.push({ key: '', value: '' })
}
/**
 * Handle delete environment variables
 * @param index - Index of environment variables
 */
const handleDeleteEnvVariable = (index: number | string) => {
  pageInfo.value.formData.environmentVariables.splice(Number(index), 1)
}
/**
 * Handle McpProtocol Changed
 */
const handleMcpProtocolChange = (servicePath: number) => {
  pageInfo.value.formData.servicePath = ['', '/sse', '/mcp', '/sse'][servicePath]
}

/**
 * Handle delete Volume mounting
 * @param index - Index of need Delete Volume
 */
const handleDeleteVolume = (index: number | string) => {
  pageInfo.value.formData.volumeMounts.splice(Number(index), 1)
}
/**
 * Update readOnly status
 * @param key - $event
 * @param volume - item of pvc
 */
const handlePvcChange = (_key: any, volume: VolumeMountsItme) => {
  const pvc = pvcList.value.find((pvc: PvcForm) => pvc.name === volume.pvcName)
  if (pvc) {
    const accessModes = pvc.accessModes || []
    let readOnlyValue = false
    if (accessModes.includes(NodeVisible.ROM)) {
      readOnlyValue = true
    } else if (accessModes.includes(NodeVisible.RWM)) {
      readOnlyValue = false
    } else if (accessModes.includes(NodeVisible.RWO)) {
      readOnlyValue = false
    }
    // Update readOnly status
    volume.readOnly = readOnlyValue
  }
}
/**
 * Handle add a volume
 */
const handleAddVolume = () => {
  pageInfo.value.formData.volumeMounts.push({
    type: '',
    nodeName: '', // HostPath name
    hostPath: '',
    mountPath: '',
    pvcName: '',
    volumeName: '',
    readOnly: false,
  })
}
/**
 * Handle change environment event
 */
const handleChangeEnvironmentId = async (e: number | undefined) => {
  pageInfo.value.formData.environmentId = e
  if (envList.value.find((item) => item.id === e)?.environment === 'docker') {
    handleGetVolumeList(pageInfo.value.formData.environmentId)
  } else {
    handleGetNodeList(pageInfo.value.formData.environmentId)
    handleGetPvcList(pageInfo.value.formData.environmentId)
  }
}

/**
 * Handle change path
 */
const handleChangePath = () => {
  pageInfo.value.formData.servicePath =
    pageInfo.value.formData.servicePath === '/sse' ? '/mcp' : '/sse'
}

const config = ref()
const handleConfig = () => {
  config.value.init(Object.assign(currentInstance.value, pageInfo.value.formData))
}
const probe = ref()
const handleViewStatus = () => {
  probe.value.init(pageInfo.value.formData)
}
const log = ref()
const handleViewLog = () => {
  log.value.init(pageInfo.value.formData)
}
// Handle confirm save
const handleConfirm = async () => {
  if (pageInfo.value.formData.instanceId) {
    const result = await ElMessageBox.confirm(
      t('mcp.instance.pageDesc.confirmEdit'),
      t('common.warn'),
      {
        confirmButtonText: t('common.ok'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
        customClass: 'tips-box',
        center: true,
        showClose: false,
        confirmButtonClass: 'is-plain el-button--danger danger-btn',
        customStyle: {
          width: '517px',
          height: '247px',
        },
      },
    )
    if (result !== 'confirm') return
  }
  baseInfo.value.validate(async (valid: boolean) => {
    if (valid) {
      try {
        pageInfo.value.loading = true
        if (!pageInfo.value.formData.instanceId) {
          if (Array.isArray(pageInfo.value.formData.tokens[0].headers)) {
            pageInfo.value.formData.tokens[0].headers = Object.fromEntries(
              pageInfo.value.formData.tokens[0].headers?.map((header: any) => [
                header.key,
                header.value,
              ]),
            )
          }
        }
        const { instanceId } = await (
          pageInfo.value.formData.instanceId ? InstanceAPI.edit : InstanceAPI.create
        )({
          ...pageInfo.value.formData,
          environmentVariables: pageInfo.value.formData.environmentVariables?.reduce(
            (obj: any, item: any) => ({ ...obj, [item.key]: item.value }),
            {},
          ),
        })
        ElMessage.success(
          pageInfo.value.formData.instanceId ? t('action.edit') : t('action.create'),
        )
        pageInfo.value.formData.instanceId = instanceId
        // 访问路径上添加instanceId
        jumpToPage({
          url: '/new-instance',
          data: {
            instanceId,
            type: query.type,
          },
        })
      } finally {
        pageInfo.value.loading = false
      }
    } else {
      ElMessage.warning(t('mcp.template.rules.validForm'))
    }
  })
}
/**
 * save as a template
 */
const handleSaveAsTemplate = () => {
  try {
    pageInfo.value.loading = true
    baseInfo.value.validate(async (valid: boolean) => {
      if (valid) {
        try {
          pageInfo.value.loading = true
          await TemplateAPI.create({
            ...pageInfo.value.formData,
            environmentVariables: pageInfo.value.formData.environmentVariables?.reduce(
              (obj: any, item: any) => ({ ...obj, [item.key]: item.value }),
              {},
            ),
          })
          ElMessage.success(t('action.create'))
        } finally {
          pageInfo.value.loading = false
        }
      }
    })
  } finally {
    pageInfo.value.loading = false
  }
}

const downloadDialogVisible = ref(false)

const handleDownloadExample = async (item: any) => {
  const { list } = await CodeAPI.list({
    page: 1,
    pageSize: 10,
    name: item.name,
  })
  if (list[0]) {
    await handleDownload(list[0])
  }
}
/**
 * Handle download code package
 */
const handleDownload = async (code: any) => {
  try {
    pageInfo.value.loading = true
    const response = await CodeAPI.download(code)
    const blobUrl = URL.createObjectURL(
      new Blob([response.data], { type: response.headers['content-type'] }),
    )
    const link = document.createElement('a')
    link.href = blobUrl
    link.download =
      response.headers['content-disposition']
        ?.split('filename=')[1]
        ?.match(/filename=("?)(.*?)\1/) || code.name
    document.body.appendChild(link)
    link.click()
  } finally {
    pageInfo.value.loading = false
  }
}
const selectedCommand = ref({}) as any

const handleSelectCommand = (item: any) => {
  if (selectedCommand.value.name === item.name) {
    selectedCommand.value = {}
  } else {
    selectedCommand.value = item
  }
}

const handleUseCommand = (command: string, type: string) => {
  if (type === 'initScript') {
    pageInfo.value.formData.initScript = command
  } else if (type === 'command') {
    pageInfo.value.formData.command = command
  }
}
const packageNameRef = ref()
const isPackageNameOverflow = ref(false)
const checkPackageNameOverflow = () => {
  if (packageNameRef.value) {
    const el = packageNameRef.value
    isPackageNameOverflow.value = el.scrollWidth > el.clientWidth
  }
}

const init = async (instance: InstanceResult | null) => {
  if (instance) {
    pageInfo.value.formData = cloneDeep(instance)
  } else {
    pageInfo.value.formData.mcpProtocol = 3
    pageInfo.value.formData.servicePath = '/sse'
  }
  await handleGetEnvList() // 获取环境变量列表
  await handleGetPackageList() // 获取包列表
  nextTick(() => {
    pageInfo.value.formData.accessType = AccessType.HOSTING
    // 默认选中第一个环境变量;所以上诉均没有return
    handleChangeEnvironmentId(envList.value[0]?.id)
  })
}

defineExpose({
  init,
  handleConfig,
  handleViewStatus,
  handleViewLog,
  handleConfirm,
  handleSaveAsTemplate,
})
</script>
<style lang="scss" scoped>
.add-env {
  width: 100%;
  border: 1px dashed var(--el-border-color);
}
.tip {
  padding: 10px;
  border-radius: 4px;
  font-size: 12px;
  &.tip-warning {
    background-color: #fff1f0;
    border-left: 5px solid var(--el-color-danger);
  }
  &.tip-primary {
    background-color: #409eff1a;
    border-left: 5px solid var(--el-color-primary);
  }
}
.upload-demo {
  width: 100%;
  color: var(--el-color-primary);
  :deep(.el-upload.is-drag) {
    height: 100%;
  }
  .el-icon--upload {
    color: var(--el-color-primary);
  }
  .el-upload__text {
    color: var(--el-color-primary);
  }
}
:deep(.el-input__wrapper) {
  background: var(--ep-bg-form);
}
:deep(.el-select__wrapper) {
  background: var(--ep-bg-form);
}
:deep(.el-textarea__inner) {
  background: var(--ep-bg-form);
}
:deep(.el-button) {
  background-color: transparent;
}
:deep(.deep-form .el-radio-button__inner) {
  padding: 0;
}
:deep(.no-border-input .el-input__wrapper) {
  box-shadow: none !important;
  // background: transparent !important;
  padding: 0 11px !important;
}
.type-item {
  font-size: 14px;
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
/* 让 el-collapse-item 左侧图标与标题内容顶部对齐 */
:deep(.el-collapse-item__header) {
  align-items: flex-start;
}

/* header 内部容器通常也是 flex，保持一致 */
:deep(.el-collapse-item__header .el-collapse-item__title) {
  align-items: flex-start;
}

/* 图标在顶部附近，避免贴边太紧 */
:deep(.el-collapse-item__header .el-collapse-item__arrow) {
  margin-top: 16px;
}
</style>
