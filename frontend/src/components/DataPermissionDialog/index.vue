<template>
  <el-dialog
    v-model="dialogVisible"
    :title="t('dataPermission.title')"
    width="780px"
    align-center
    :show-close="true"
    destroy-on-close
    class="dp-dialog"
    @closed="handleClosed"
  >
    <div v-loading="loading" class="dp-body">
      <!-- 一级 Tab：访问授权 / 访问黑名单 -->
      <div class="dp-tab-bar">
        <div
          :class="['dp-tab-item', activeTab === 'access' && 'active']"
          @click="activeTab = 'access'"
        >
          {{ t('dataPermission.tab.access') }}
        </div>
        <div
          :class="['dp-tab-item', activeTab === 'blacklist' && 'active']"
          @click="activeTab = 'blacklist'"
        >
          {{ t('dataPermission.tab.blacklist') }}
        </div>
      </div>

      <!-- ========== 访问授权 ========== -->
      <template v-if="activeTab === 'access'">
        <!-- 单选：全部人员 / 部分人员 -->
        <el-radio-group v-model="isAllPersonnel" class="dp-radio-row">
          <el-radio :value="true">{{ t('dataPermission.personnelScope.all') }}</el-radio>
          <el-radio :value="false">{{ t('dataPermission.personnelScope.partial') }}</el-radio>
        </el-radio-group>

        <!-- 全部人员提示 -->
        <div v-if="isAllPersonnel" class="dp-all-tip">
          {{ t('dataPermission.allPersonnelTip') }}
        </div>

        <!-- 部分人员 -->
        <template v-if="!isAllPersonnel">
          <!-- 二级维度 Tab -->
          <div class="dp-dim-tabs">
            <span
              v-for="dim in dimensions"
              :key="dim.key"
              :class="['dp-dim-tab', activeDimension === dim.key && 'active']"
              @click="activeDimension = dim.key"
            >
              {{ dim.label }}
            </span>
          </div>

          <!-- 左右分栏 -->
          <div class="dp-split-panel">
            <!-- 左侧：选择面板 -->
            <div class="dp-left">
              <div class="dp-panel-header h-10">
                <el-checkbox
                  v-model="isLeftAllChecked"
                  :indeterminate="isLeftIndeterminate"
                  @change="handleLeftAllChange"
                >
                  {{ t('dataPermission.all') }}（{{ leftTotal }}）
                </el-checkbox>
              </div>
              <div class="dp-panel-search">
                <el-input
                  v-model="leftKeyword"
                  :placeholder="leftSearchPlaceholder"
                  clearable
                  size="default"
                />
              </div>
              <el-scrollbar class="dp-panel-list">
                <!-- 组织架构：树 -->
                <template v-if="activeDimension === 'dept'">
                  <el-tree
                    ref="deptTreeRef"
                    :data="deptTreeData"
                    show-checkbox
                    node-key="id"
                    :default-checked-keys="selectedDeptIds"
                    :props="{ label: 'name', children: 'children' }"
                    :filter-node-method="filterDeptNode"
                    @check="handleDeptCheck"
                  />
                </template>
                <!-- 角色：checkbox 列表 -->
                <template v-if="activeDimension === 'role'">
                  <el-checkbox-group v-model="selectedRoleIds">
                    <div v-for="role in filteredRoleList" :key="role.id" class="dp-check-row">
                      <el-checkbox :value="role.id">{{ role.name }}</el-checkbox>
                    </div>
                  </el-checkbox-group>
                  <el-empty
                    v-if="!filteredRoleList.length"
                    :description="t('dataPermission.noData')"
                    :image-size="48"
                  />
                </template>
                <!-- 人员：表格样式 -->
                <template v-if="activeDimension === 'user'">
                  <table class="dp-user-table">
                    <thead>
                      <tr>
                        <th style="width: 32px">
                          <el-checkbox
                            v-model="isUserAllChecked"
                            :indeterminate="isUserIndeterminate"
                            @change="handleUserAllChange"
                          />
                        </th>
                        <th>{{ t('dataPermission.nickname') }}</th>
                        <th>{{ t('dataPermission.username') }}</th>
                        <th>{{ t('dataPermission.deptName') }}</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="user in filteredUserList" :key="user.id">
                        <td>
                          <el-checkbox
                            :model-value="selectedUserIds.includes(user.id)"
                            @change="(val: boolean) => toggleUser(user.id, val)"
                          />
                        </td>
                        <td>{{ user.username || '-' }}</td>
                        <td>{{ user.nickName || '-' }}</td>
                        <td>{{ user.deptName || '-' }}</td>
                      </tr>
                    </tbody>
                  </table>
                  <el-empty
                    v-if="!filteredUserList.length"
                    :description="t('dataPermission.noData')"
                    :image-size="48"
                  />
                </template>
              </el-scrollbar>
            </div>

            <!-- 右侧：已选面板 -->
            <div class="dp-right">
              <div class="dp-panel-header h-10 dp-right-header">
                <span class="dp-selected-summary">
                  {{ t('dataPermission.selectedDept') }}{{ selectedDeptIds.length
                  }}{{ t('dataPermission.unit') }}，{{ t('dataPermission.selectedRole')
                  }}{{ selectedRoleIds.length }}{{ t('dataPermission.unit') }}，{{
                    t('dataPermission.selectedMember')
                  }}{{ selectedUserIds.length }}{{ t('dataPermission.unit') }}
                </span>
                <span class="dp-clear-btn" @click="handleClearAccess">{{
                  t('dataPermission.clear')
                }}</span>
              </div>
              <div class="dp-panel-search">
                <el-input
                  v-model="rightKeyword"
                  :placeholder="t('dataPermission.placeholder.searchSelectedUser')"
                  clearable
                  size="default"
                />
              </div>
              <el-scrollbar class="dp-panel-list">
                <!-- 已选部门 -->
                <div
                  v-for="item in filteredSelectedDeptDisplayList"
                  :key="'dept-' + item.id"
                  class="dp-selected-item"
                >
                  <el-avatar :size="32" class="dp-avatar dp-avatar-dept">
                    <el-icon :size="16"><OfficeBuilding /></el-icon>
                  </el-avatar>
                  <span class="dp-selected-name">{{ item.name }}</span>
                  <span v-if="!isProtectedAdminDept(item.id)" class="dp-remove-btn" @click="removeSelectedDept(item.id)">移除</span>
                </div>
                <!-- 已选角色 -->
                <div
                  v-for="item in filteredSelectedRoleDisplayList"
                  :key="'role-' + item.id"
                  class="dp-selected-item"
                >
                  <el-avatar :size="32" class="dp-avatar dp-avatar-role">
                    <el-icon :size="16"><UserFilled /></el-icon>
                  </el-avatar>
                  <span class="dp-selected-name">{{ item.name }}</span>
                  <span v-if="!isProtectedAdminRole(item)" class="dp-remove-btn" @click="removeSelectedRole(item.id)">移除</span>
                </div>
                <!-- 已选人员 -->
                <div
                  v-for="item in filteredSelectedUserDisplayList"
                  :key="'user-' + item.id"
                  class="dp-selected-item"
                >
                  <el-avatar :size="32" class="dp-avatar dp-avatar-user">
                    {{ (item.username || item.nickName || '').charAt(0) }}
                  </el-avatar>
                  <span class="dp-selected-name dp-selected-name-user">
                    <span>{{ item.username || item.nickName }}</span>
                    <span v-if="item.deptName" class="dp-selected-sub"
                      >{{ item.deptName }}{{ item.nickName ? `，${item.nickName}` : '' }}</span
                    >
                  </span>
                  <span v-if="!isProtectedAdminUser(item.id)" class="dp-remove-btn" @click="removeSelectedUser(item.id)">移除</span>
                </div>
                <el-empty
                  v-if="
                    !filteredSelectedDeptDisplayList.length &&
                    !filteredSelectedRoleDisplayList.length &&
                    !filteredSelectedUserDisplayList.length
                  "
                  :description="t('dataPermission.noData')"
                  :image-size="48"
                />
              </el-scrollbar>
            </div>
          </div>
        </template>
      </template>

      <!-- ========== 访问黑名单 ========== -->
      <template v-if="activeTab === 'blacklist'">
        <div class="dp-split-panel dp-blacklist-panel">
          <!-- 左侧：选人 -->
          <div class="dp-left">
            <div class="dp-panel-header">
              <el-checkbox
                v-model="isBlacklistAllChecked"
                :indeterminate="isBlacklistIndeterminate"
                @change="handleBlacklistAllChange"
              >
                {{ t('dataPermission.all') }}（{{ userList.length }}）
              </el-checkbox>
            </div>
            <div class="dp-panel-search">
              <el-input
                v-model="blacklistKeyword"
                :placeholder="t('dataPermission.placeholder.searchUser')"
                clearable
                size="default"
              />
            </div>
            <el-scrollbar class="dp-panel-list">
              <table class="dp-user-table">
                <thead>
                  <tr>
                    <th style="width: 32px">
                      <el-checkbox
                        v-model="isBlacklistAllChecked"
                        :indeterminate="isBlacklistIndeterminate"
                        @change="handleBlacklistAllChange"
                      />
                    </th>
                    <th>{{ t('dataPermission.nickname') }}</th>
                    <th>{{ t('dataPermission.username') }}</th>
                    <th>{{ t('dataPermission.deptName') }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="user in filteredBlacklistUserList" :key="user.id">
                    <td>
                      <el-checkbox
                        :model-value="selectedBlacklistUserIds.includes(user.id)"
                        @change="(val: boolean) => toggleBlacklistUser(user.id, val)"
                      />
                    </td>
                    <td>{{ user.username || '-' }}</td>
                    <td>{{ user.nickName || '-' }}</td>
                    <td>{{ user.deptName || '-' }}</td>
                  </tr>
                </tbody>
              </table>
              <el-empty
                v-if="!filteredBlacklistUserList.length"
                :description="t('dataPermission.noData')"
                :image-size="48"
              />
            </el-scrollbar>
          </div>

          <!-- 右侧：已选黑名单 -->
          <div class="dp-right">
            <div class="dp-panel-header dp-right-header">
              <span class="dp-selected-summary">
                {{ t('dataPermission.selectedBlacklist') }}（{{
                  selectedBlacklistUserIds.length
                }}）{{ t('dataPermission.unit') }}
              </span>
              <span class="dp-clear-btn" @click="selectedBlacklistUserIds = []">{{
                t('dataPermission.clear')
              }}</span>
            </div>
            <div class="dp-panel-search">
              <el-input
                v-model="rightBlacklistKeyword"
                :placeholder="t('dataPermission.placeholder.searchSelectedUser')"
                clearable
                size="default"
              />
            </div>
            <el-scrollbar class="dp-panel-list">
              <div
                v-for="item in filteredSelectedBlacklistDisplayList"
                :key="item.id"
                class="dp-selected-item"
              >
                <el-avatar :size="32" class="dp-avatar dp-avatar-user">
                  {{ (item.username || item.nickName || '').charAt(0) }}
                </el-avatar>
                <span class="dp-selected-name">{{ item.username || item.nickName }}</span>
                <span class="dp-remove-btn" @click="removeBlacklistUser(item.id)">移除</span>
              </div>
              <el-empty
                v-if="!filteredSelectedBlacklistDisplayList.length"
                :description="t('dataPermission.noData')"
                :image-size="48"
              />
            </el-scrollbar>
          </div>
        </div>
      </template>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">
          {{ t('common.ok') }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ElMessage } from 'element-plus'
import { OfficeBuilding, UserFilled } from '@element-plus/icons-vue'
import { DeptAPI, RoleAPI, UserAPI } from '@/api/system/index'
import { DataPermissionAPI } from '@/api/data-permission/index'
import type { DataPermissionResult } from '@/api/data-permission/index'
import { useUserStore } from '@/stores/modules/user-store.ts'

const { t } = useI18n()

const { userInfo } = useUserStore()
const isAdmin = computed(() => userInfo.isAdmin === true)

// 判定管理员账户不可被脱保的辅助函数
const isProtectedAdminDept = (deptId: number) => {
  return isAdmin.value && String(deptId) === String(userInfo.deptId)
}
const isProtectedAdminRole = (roleItem: any) => {
  return isAdmin.value && (roleItem.name === 'admin' || (userInfo.roleIds || []).some((id: any) => String(id) === String(roleItem.id)))
}
const isProtectedAdminUser = (userId: number) => {
  return isAdmin.value && String(userId) === String(userInfo.userId)
}

// ===================== 基础状态 =====================
const dialogVisible = ref(false)
const loading = ref(false)
const submitLoading = ref(false)
const resourceName = ref('')
const resourceId = ref('')
const resourceType = ref('')

// ===================== Tab 状态 =====================
const activeTab = ref<'access' | 'blacklist'>('access')
const activeDimension = ref<'dept' | 'role' | 'user'>('dept')
const isAllPersonnel = ref(true)

// ===================== 数据源 =====================
const deptTreeData = ref<any[]>([])
const roleList = ref<any[]>([])
const userList = ref<any[]>([])

// ===================== 已选 ID =====================
const selectedDeptIds = ref<number[]>([])
const selectedRoleIds = ref<number[]>([])
const selectedUserIds = ref<number[]>([])
const selectedBlacklistUserIds = ref<number[]>([])

// ===================== 搜索关键字 =====================
const leftKeyword = ref('')
const rightKeyword = ref('')
const blacklistKeyword = ref('')
const rightBlacklistKeyword = ref('')

// ===================== Tree ref =====================
const deptTreeRef = ref<any>(null)

// ===================== 维度配置 =====================
const dimensions = computed(() => [
  { key: 'dept' as const, label: t('dataPermission.dimension.dept') },
  { key: 'role' as const, label: t('dataPermission.dimension.role') },
  { key: 'user' as const, label: t('dataPermission.dimension.user') },
])

// ===================== 资源类型映射 =====================
const RESOURCE_TYPE_MAP: Record<string, string> = {
  instance: 'dataPermission.type.instance',
  openapi_package: 'dataPermission.type.openapiPackage',
  code_package: 'dataPermission.type.codePackage',
  instance_tokens: 'dataPermission.type.instanceTokens',
  intelligent_access: 'dataPermission.type.intelligentAccess',
}

const resourceTypeLabel = computed(() => {
  const key = RESOURCE_TYPE_MAP[resourceType.value]
  return key ? t(key) : resourceType.value
})

// ===================== 左侧搜索 placeholder =====================
const leftSearchPlaceholder = computed(() => {
  const map: Record<string, string> = {
    dept: t('dataPermission.placeholder.searchDept'),
    role: t('dataPermission.placeholder.searchRole'),
    user: t('dataPermission.placeholder.searchUser'),
  }
  return map[activeDimension.value] || ''
})

// ===================== 左侧总数 =====================
const leftTotal = computed(() => {
  if (activeDimension.value === 'dept') return flattenDeptCount(deptTreeData.value)
  if (activeDimension.value === 'role') return roleList.value.length
  if (activeDimension.value === 'user') return userList.value.length
  return 0
})

/** 递归计算树节点总数 */
const flattenDeptCount = (nodes: any[]): number => {
  let count = 0
  for (const n of nodes) {
    count++
    if (n.children?.length) count += flattenDeptCount(n.children)
  }
  return count
}

// ===================== 搜索过滤 =====================
const filterDeptNode = (value: string, data: any) => {
  if (!value) return true
  return data.name?.toLowerCase().includes(value.toLowerCase())
}

// 左侧搜索联动维度切换时清空
watch(activeDimension, () => {
  leftKeyword.value = ''
})

watch(leftKeyword, (val) => {
  if (activeDimension.value === 'dept') {
    deptTreeRef.value?.filter(val)
  }
})

const filteredRoleList = computed(() => {
  if (!leftKeyword.value) return roleList.value
  const kw = leftKeyword.value.toLowerCase()
  return roleList.value.filter((r: any) => r.name?.toLowerCase().includes(kw))
})

const filteredUserList = computed(() => {
  if (!leftKeyword.value) return userList.value
  const kw = leftKeyword.value.toLowerCase()
  return userList.value.filter(
    (u: any) => u.username?.toLowerCase().includes(kw) || u.nickName?.toLowerCase().includes(kw),
  )
})

const filteredBlacklistUserList = computed(() => {
  if (!blacklistKeyword.value) return userList.value
  const kw = blacklistKeyword.value.toLowerCase()
  return userList.value.filter(
    (u: any) => u.username?.toLowerCase().includes(kw) || u.nickName?.toLowerCase().includes(kw),
  )
})

// ===================== 左侧全选 / 半选（按维度） =====================
const isLeftAllChecked = computed({
  get: () => {
    if (activeDimension.value === 'dept') {
      const total = flattenDeptCount(deptTreeData.value)
      return total > 0 && selectedDeptIds.value.length === total
    }
    if (activeDimension.value === 'role') {
      return roleList.value.length > 0 && selectedRoleIds.value.length === roleList.value.length
    }
    if (activeDimension.value === 'user') {
      return userList.value.length > 0 && selectedUserIds.value.length === userList.value.length
    }
    return false
  },
  set: () => {},
})

const isLeftIndeterminate = computed(() => {
  if (activeDimension.value === 'dept') {
    const total = flattenDeptCount(deptTreeData.value)
    return selectedDeptIds.value.length > 0 && selectedDeptIds.value.length < total
  }
  if (activeDimension.value === 'role') {
    return selectedRoleIds.value.length > 0 && selectedRoleIds.value.length < roleList.value.length
  }
  if (activeDimension.value === 'user') {
    return selectedUserIds.value.length > 0 && selectedUserIds.value.length < userList.value.length
  }
  return false
})

const handleLeftAllChange = (val: boolean) => {
  if (activeDimension.value === 'dept') {
    if (val) {
      const allIds = flattenDeptIds(deptTreeData.value)
      selectedDeptIds.value = allIds
      nextTick(() => deptTreeRef.value?.setCheckedKeys(allIds))
    } else {
      selectedDeptIds.value = []
      nextTick(() => deptTreeRef.value?.setCheckedKeys([]))
    }
  } else if (activeDimension.value === 'role') {
    selectedRoleIds.value = val ? roleList.value.map((r: any) => r.id) : []
  } else if (activeDimension.value === 'user') {
    selectedUserIds.value = val ? userList.value.map((u: any) => u.id) : []
  }
}

/** 递归获取所有部门 id */
const flattenDeptIds = (nodes: any[]): number[] => {
  const ids: number[] = []
  for (const n of nodes) {
    ids.push(n.id)
    if (n.children?.length) ids.push(...flattenDeptIds(n.children))
  }
  return ids
}

// ===================== 人员全选（user 维度表格表头） =====================
const isUserAllChecked = computed({
  get: () => userList.value.length > 0 && selectedUserIds.value.length === userList.value.length,
  set: () => {},
})
const isUserIndeterminate = computed(() => {
  return selectedUserIds.value.length > 0 && selectedUserIds.value.length < userList.value.length
})
const handleUserAllChange = (val: boolean) => {
  selectedUserIds.value = val ? userList.value.map((u: any) => u.id) : []
}
const toggleUser = (id: number, val: boolean) => {
  if (val) {
    if (!selectedUserIds.value.includes(id)) selectedUserIds.value.push(id)
  } else {
    selectedUserIds.value = selectedUserIds.value.filter((v) => v !== id)
  }
}

// ===================== 黑名单全选 =====================
const isBlacklistAllChecked = computed({
  get: () =>
    userList.value.length > 0 && selectedBlacklistUserIds.value.length === userList.value.length,
  set: () => {},
})
const isBlacklistIndeterminate = computed(() => {
  return (
    selectedBlacklistUserIds.value.length > 0 &&
    selectedBlacklistUserIds.value.length < userList.value.length
  )
})
const handleBlacklistAllChange = (val: boolean) => {
  selectedBlacklistUserIds.value = val ? userList.value.map((u: any) => u.id) : []
}
const toggleBlacklistUser = (id: number, val: boolean) => {
  if (val) {
    if (!selectedBlacklistUserIds.value.includes(id)) selectedBlacklistUserIds.value.push(id)
  } else {
    selectedBlacklistUserIds.value = selectedBlacklistUserIds.value.filter((v) => v !== id)
  }
}

// ===================== Tree check 事件 =====================
const handleDeptCheck = () => {
  if (deptTreeRef.value) {
    selectedDeptIds.value = deptTreeRef.value.getCheckedKeys(false) as number[]
  }
}

// ===================== 右侧已选人员展示（访问授权） =====================
const selectedUserDisplayList = computed(() => {
  return userList.value.filter((u: any) => selectedUserIds.value.includes(u.id))
})

const filteredSelectedUserDisplayList = computed(() => {
  if (!rightKeyword.value) return selectedUserDisplayList.value
  const kw = rightKeyword.value.toLowerCase()
  return selectedUserDisplayList.value.filter(
    (u: any) => u.username?.toLowerCase().includes(kw) || u.nickName?.toLowerCase().includes(kw),
  )
})

// ===================== 右侧已选部门展示 =====================
/** 递归获取所有部门节点的扁平列表 */
const flattenDeptNodes = (nodes: any[]): any[] => {
  const result: any[] = []
  for (const n of nodes) {
    result.push(n)
    if (n.children?.length) result.push(...flattenDeptNodes(n.children))
  }
  return result
}

const selectedDeptDisplayList = computed(() => {
  const allNodes = flattenDeptNodes(deptTreeData.value)
  return allNodes.filter((n: any) => selectedDeptIds.value.includes(n.id))
})

const filteredSelectedDeptDisplayList = computed(() => {
  if (!rightKeyword.value) return selectedDeptDisplayList.value
  const kw = rightKeyword.value.toLowerCase()
  return selectedDeptDisplayList.value.filter((n: any) => n.name?.toLowerCase().includes(kw))
})

// ===================== 右侧已选角色展示 =====================
const selectedRoleDisplayList = computed(() => {
  return roleList.value.filter((r: any) => selectedRoleIds.value.includes(r.id))
})

const filteredSelectedRoleDisplayList = computed(() => {
  if (!rightKeyword.value) return selectedRoleDisplayList.value
  const kw = rightKeyword.value.toLowerCase()
  return selectedRoleDisplayList.value.filter((r: any) => r.name?.toLowerCase().includes(kw))
})

// ===================== 右侧移除操作（访问授权） =====================
const removeSelectedDept = (id: number) => {
  selectedDeptIds.value = selectedDeptIds.value.filter((v) => v !== id)
  nextTick(() => deptTreeRef.value?.setCheckedKeys(selectedDeptIds.value))
}

const removeSelectedRole = (id: number) => {
  selectedRoleIds.value = selectedRoleIds.value.filter((v) => v !== id)
}

const removeSelectedUser = (id: number) => {
  selectedUserIds.value = selectedUserIds.value.filter((v) => v !== id)
}

// ===================== 右侧已选人员展示（黑名单） =====================
const selectedBlacklistDisplayList = computed(() => {
  return userList.value.filter((u: any) => selectedBlacklistUserIds.value.includes(u.id))
})

const filteredSelectedBlacklistDisplayList = computed(() => {
  if (!rightBlacklistKeyword.value) return selectedBlacklistDisplayList.value
  const kw = rightBlacklistKeyword.value.toLowerCase()
  return selectedBlacklistDisplayList.value.filter(
    (u: any) => u.username?.toLowerCase().includes(kw) || u.nickName?.toLowerCase().includes(kw),
  )
})

// ===================== 右侧移除操作（黑名单） =====================
const removeBlacklistUser = (id: number) => {
  selectedBlacklistUserIds.value = selectedBlacklistUserIds.value.filter((v) => v !== id)
}

// ===================== 清除访问授权已选 =====================
const handleClearAccess = () => {
  if (isAdmin.value) {
    selectedDeptIds.value = selectedDeptIds.value.filter(id => isProtectedAdminDept(id))
    selectedRoleIds.value = selectedRoleIds.value.filter(id => {
      const role = roleList.value.find(r => r.id === id)
      return role ? isProtectedAdminRole(role) : false
    })
    selectedUserIds.value = selectedUserIds.value.filter(id => isProtectedAdminUser(id))
    nextTick(() => {
      deptTreeRef.value?.setCheckedKeys(selectedDeptIds.value)
    })
  } else {
    selectedDeptIds.value = []
    selectedRoleIds.value = []
    selectedUserIds.value = []
    nextTick(() => deptTreeRef.value?.setCheckedKeys([]))
  }
}

// ===================== 事件 =====================
const emit = defineEmits<{
  (e: 'on-refresh'): void
}>()

export interface DataPermissionInitParams {
  id: string
  name: string
  type: string
}

// ===================== 加载数据源 =====================
const loadSourceData = async () => {
  const [deptRes, roleRes, userRes] = await Promise.all([
    DeptAPI.deptTree(),
    RoleAPI.allList(),
    UserAPI.list({ page: 1, pageSize: 9999 }),
  ])
  deptTreeData.value = deptRes || []
  roleList.value = roleRes?.list || []
  userList.value = userRes?.list || []
}

// ===================== 加载已有权限数据 =====================
const loadPermissionData = async () => {
  try {
    const data: DataPermissionResult = await DataPermissionAPI.get({
      dataType: resourceType.value,
      dataId: resourceId.value,
    })
    isAllPersonnel.value = data.isAllPersonnel ?? true
    selectedDeptIds.value = data.deptIds || []
    selectedRoleIds.value = data.roleIds || []
    selectedUserIds.value = data.userIds || []
    selectedBlacklistUserIds.value = data.blacklistUserIds || []

    nextTick(() => {
      if (deptTreeRef.value) {
        deptTreeRef.value.setCheckedKeys(selectedDeptIds.value)
      }
    })
  } catch {
    // 无已有数据
  }
}

const init = async (params: DataPermissionInitParams) => {
  resetForm()
  resourceId.value = params.id
  resourceName.value = params.name
  resourceType.value = params.type
  dialogVisible.value = true

  try {
    loading.value = true
    await loadSourceData()
    await loadPermissionData()
  } catch (e) {
    console.error(t('dataPermission.loadFailed'), e)
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  try {
    submitLoading.value = true
    await DataPermissionAPI.save({
      dataType: resourceType.value,
      dataId: resourceId.value,
      isAllPersonnel: isAllPersonnel.value,
      deptIds: isAllPersonnel.value ? [] : selectedDeptIds.value,
      roleIds: isAllPersonnel.value ? [] : selectedRoleIds.value,
      userIds: isAllPersonnel.value ? [] : selectedUserIds.value,
      blacklistUserIds: selectedBlacklistUserIds.value,
    })
    ElMessage.success(t('dataPermission.saveSuccess'))
    dialogVisible.value = false
    emit('on-refresh')
  } finally {
    submitLoading.value = false
  }
}

const resetForm = () => {
  activeTab.value = 'access'
  activeDimension.value = 'dept'
  isAllPersonnel.value = true
  selectedDeptIds.value = []
  selectedRoleIds.value = []
  selectedUserIds.value = []
  selectedBlacklistUserIds.value = []
  leftKeyword.value = ''
  rightKeyword.value = ''
  blacklistKeyword.value = ''
  rightBlacklistKeyword.value = ''
}

const handleClosed = () => {
  resourceName.value = ''
  resourceId.value = ''
  resourceType.value = ''
  resetForm()
}

defineExpose({ init })
</script>

<style scoped lang="scss">
/* ===== 弹窗固定高度 ===== */
.dp-dialog {
  :deep(.el-dialog__body) {
    padding: 16px 20px;
  }
}

.dp-body {
  height: 520px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* ===== 一级 Tab 按钮样式 ===== */
.dp-tab-bar {
  display: flex;
  gap: 0;
  margin-bottom: 16px;
  flex-shrink: 0;
}

.dp-tab-item {
  padding: 6px 20px;
  font-size: 14px;
  cursor: pointer;
  border: 1px solid var(--el-color-primary);
  color: var(--el-color-primary);
  transition: all 0.2s;
  user-select: none;

  &:first-child {
    border-radius: 4px 0 0 4px;
  }

  &:last-child {
    border-radius: 0 4px 4px 0;
    border-left: none;
  }

  &.active {
    background: var(--el-color-primary);
    color: #fff;
  }
}

/* ===== 单选行 ===== */
.dp-radio-row {
  margin-bottom: 12px;
  flex-shrink: 0;
}

/* ===== 二级维度 Tab ===== */
.dp-dim-tabs {
  display: flex;
  gap: 24px;
  margin-bottom: 12px;
  border-bottom: 1px solid var(--el-border-color-light);
  flex-shrink: 0;
}

.dp-dim-tab {
  padding: 6px 0;
  font-size: 14px;
  cursor: pointer;
  color: var(--el-text-color-regular);
  border-bottom: 2px solid transparent;
  transition: all 0.2s;

  &.active {
    color: var(--el-color-primary);
    border-bottom-color: var(--el-color-primary);
    font-weight: 500;
  }
}

/* ===== 全部人员居中提示 ===== */
.dp-all-tip {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 500;
}

/* ===== 左右分栏 ===== */
.dp-split-panel {
  display: flex;
  flex: 1;
  min-height: 0;
  border: 1px solid var(--el-border-color-light);
  border-radius: 4px;
  overflow: hidden;
}

.dp-left,
.dp-right {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-width: 0;
  overflow: hidden;
}

.dp-left {
  border-right: 1px solid var(--el-border-color-light);
}

.dp-panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  border-bottom: 1px solid var(--el-border-color-lighter);
  flex-shrink: 0;
  font-size: 13px;
}

.dp-right-header {
  .dp-selected-summary {
    color: var(--el-text-color-regular);
    font-size: 12px;
  }

  .dp-clear-btn {
    color: var(--el-color-danger);
    cursor: pointer;
    font-size: 12px;

    &:hover {
      opacity: 0.8;
    }
  }
}

.dp-panel-search {
  padding: 8px 12px;
  flex-shrink: 0;
}

.dp-panel-list {
  flex: 1;
  min-height: 0;
  padding: 0 12px;
}

/* ===== checkbox 行 ===== */
.dp-check-row {
  padding: 6px 0;
  border-bottom: 1px solid var(--el-border-color-extra-light);

  &:last-child {
    border-bottom: none;
  }
}

/* ===== 人员表格 ===== */
.dp-user-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;

  th,
  td {
    padding: 6px 8px;
    text-align: left;
  }

  th {
    font-weight: 500;
    color: var(--el-text-color-secondary);
    background: var(--el-fill-color-lighter);
    border-bottom: 1px solid var(--el-border-color-light);
  }

  td {
    border-bottom: 1px solid var(--el-border-color-extra-light);
  }

  tbody tr:hover {
    background: var(--el-fill-color-light);
  }
}

/* ===== 右侧已选列表项 ===== */
.dp-selected-item {
  display: flex;
  align-items: center;
  padding: 8px 4px;
  gap: 10px;
  border-bottom: 1px solid var(--el-border-color-extra-light);
  position: relative;

  &:last-child {
    border-bottom: none;
  }

  .dp-remove-btn {
    display: none;
    position: absolute;
    right: 4px;
    color: var(--el-color-danger);
    cursor: pointer;
    font-size: 12px;
    white-space: nowrap;

    &:hover {
      opacity: 0.8;
    }
  }

  &:hover {
    background-color: var(--el-fill-color-light);
    border-radius: 4px;
  }

  &:hover .dp-remove-btn {
    display: inline;
  }
}

.dp-avatar {
  flex-shrink: 0;
  font-size: 12px;
}

.dp-avatar-dept {
  background: #9c27b0;
  color: #fff;
}

.dp-avatar-role {
  background: #9c27b0;
  color: #fff;
}

.dp-avatar-user {
  background: var(--el-color-primary-light-5);
  color: #fff;
}

.dp-selected-name {
  font-size: 13px;
  color: var(--el-text-color-primary);
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dp-selected-name-user {
  display: flex;
  flex-direction: column;
  line-height: 1.3;
}

.dp-selected-sub {
  font-size: 11px;
  color: var(--el-text-color-secondary);
}

/* ===== 黑名单面板高度与访问授权一致 ===== */
.dp-blacklist-panel {
  flex: 1;
  min-height: 0;
}
</style>
