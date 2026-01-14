import request from '@/utils/request'
import baseConfig from '@/config/base_config.ts'
import { type InstanceResult } from '@/types/instance'

export const InstanceAPI = {
  // 实例列表
  list(data: TableData) {
    return request<TableData, List<InstanceResult>>({
      url: `/market/instance/list`,
      method: 'POST',
      data,
    })
  },
  // 创建实例
  create(data: any) {
    return request<any, List>({
      url: `/market/instance/create`,
      method: 'POST',
      data,
    })
  },
  // 创建实例
  createByOpenAPI(data: any) {
    return request<any, List>({
      url: `/market/instance/openapi/create`,
      method: 'POST',
      data,
    })
  },
  // 实例详情
  detail(data: any) {
    return request<any, any>({
      url: `/market/instance/${data.instanceId}`,
      method: 'GET',
      data,
    })
  },
  // 编辑实例
  edit(data: any) {
    return request<any, any>({
      url: `/market/instance/edit`,
      method: 'PUT',
      data,
    })
  },
  // 编辑实例
  editByOpenAPI(data: any) {
    return request<any, any>({
      url: `/market/instance/openapi/edit`,
      method: 'PUT',
      data,
    })
  },
  // 获取实例日志
  logs(data: any) {
    return request<any, any>({
      url: `/market/instance/logs`,
      method: 'POST',
      data,
    })
  },
  // token网关日志
  logsByToken(data: any) {
    return request<any, any>({
      url: `/market/gateway-log/find`,
      method: 'POST',
      data,
    })
  },
  // 停止实例
  stop(data: any) {
    return request<any, any>({
      url: `/market/instance/disabled`,
      method: 'PUT',
      data,
    })
  },
  // 启动实例
  restart(data: any) {
    return request<any, any>({
      url: `/market/instance/restart`,
      method: 'PUT',
      data,
    })
  },
  // 实例状态
  status(instanceId: string) {
    return request<any, any>({
      url: `/market/instance/status/${instanceId}`,
      method: 'GET',
    })
  },
  // 删除实例
  delete(instanceId: string) {
    return request<any, any>({
      url: `/market/instance/${instanceId}`,
      method: 'DELETE',
    })
  },
  // 实例数据统计
  count() {
    return request<any, any>({
      url: `/market/dashboard/statistical`,
      method: 'GET',
    })
  },
  // change token status
  updateTokenStatus(data: { instanceId: string; enabledToken: boolean }) {
    return request<any, any>({
      url: `/market/instance/token/control`,
      method: 'PUT',
      data,
    })
  },
  updateInstanceTokens(data: {
    instanceId: string
    tokens: Array<{
      token: string
      expireAt: number
      publishAt: number
      usages: string[]
    }>
  }) {
    return request<any, any>({
      url: `/market/instance/token/edit`,
      method: 'PUT',
      data,
    })
  },
}

export const TokenAPI = {
  // token list api
  list(data: any) {
    return request<any, any>({
      url: `/market/instance/token/list`,
      method: 'POST',
      data,
    })
  },
  // delete token api
  delete(data: any) {
    return request<any, any>({
      url: `/market/instance/token/delete`,
      method: 'POST',
      data,
    })
  },
  //edit token api
  edit(data: any) {
    return request<any, any>({
      url: `/market/instance/token/edit`,
      method: 'PUT',
      data,
    })
  },
}

export const deBugAPI = {
  // 实例工具列表
  toolList(data: any) {
    return request<any, any>({
      url: `/market/instance/list-tools`,
      method: 'POST',
      data,
    })
  },
  // 实例工具调用
  toolCall(data: any) {
    return request<any, any>({
      url: `/market/instance/call-tools`,
      method: 'POST',
      data,
    })
  },
}

// 列表请求
export interface TableData {
  /** 页码 */
  page: string
  /** 每页显示数量 */
  pageSize: string
  /** 允许传入其他任意类型的参数 */
  [key: string]: any
}
// 列表返回
export interface List<T = any> {
  list: T[]
  page: number
  pageSize: number
  total: number
}
