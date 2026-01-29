import request from '@/utils/request'
import baseConfig from '@/config/base_config.ts'

export const TemplateAPI = {
  // 模板列表
  list(data: TableData) {
    return request<any, List>({
      url: `/market/template/list`,
      method: 'POST',
      data,
    })
  },
  // 创建模板
  create(data: any) {
    return request<any, List>({
      url: `/market/template/create`,
      method: 'POST',
      data,
    })
  },
  // 删除模板
  delete(templateId: string) {
    return request<any, List>({
      url: `/market/template/${templateId}`,
      method: 'DELETE',
    })
  },
  // 模板详情
  detail(data: any) {
    return request<any, any>({
      url: `/market/template/${data.id}`,
      method: 'GET',
      data,
    })
  },
  // 编辑模板
  edit(data: any) {
    return request<any, any>({
      url: `/market/template/edit`,
      method: 'PUT',
      data,
    })
  },
  // 可用案例
  cases() {
    return request<any, any>({
      url: `/market/dashboard/available-cases`,
      method: 'GET',
      data: {},
    })
  },
}

// 列表请求
export interface TableData {
  /** 页码 */
  page: number
  /** 每页显示数量 */
  pageSize: number
  /** 允许传入其他任意类型的参数 */
  [key: string]: any
}
// 列表返回
export interface List {
  list: any
  data: object
}
