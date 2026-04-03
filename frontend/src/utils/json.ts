/**
 * JSON 格式化工具类
 * 支持格式化（带缩进）、压缩（无空格）、验证JSON有效性
 */
export const JsonFormatter = {
  /**
   * 格式化JSON字符串（带缩进和换行）
   * @param jsonString 原始JSON字符串
   * @param indent 缩进空格数，默认2
   * @returns 格式化后的JSON字符串或错误信息
   */
  format: (jsonString: string, indent: number = 2): string => {
    try {
      const processedStr = jsonString.replace(/\\"/g, '"').replace(/'/g, '"')
      // 先解析验证JSON有效性
      const parsed = JSON.parse(processedStr)
      // 格式化并返回
      return JSON.stringify(parsed, null, indent)
    } catch (error) {
      return jsonString
    }
  },

  /**
   * 压缩JSON字符串（移除所有空格和换行）
   * @param jsonString 原始JSON字符串
   * @returns 压缩后的JSON字符串或错误信息
   */
  minify: (jsonString: string): string => {
    try {
      const parsed = JSON.parse(jsonString)
      return JSON.stringify(parsed)
    } catch (error) {
      return `压缩失败：${(error as Error).message}`
    }
  },

  /**
   * 验证JSON字符串是否有效
   * @param jsonString 待验证的JSON字符串
   * @returns 验证结果对象 { valid: boolean, error?: string }
   */
  validate: (jsonString: string): { valid: boolean; error?: string } => {
    try {
      JSON.parse(jsonString)
      return { valid: true }
    } catch (error) {
      return { valid: false, error: (error as Error).message }
    }
  },

  /**
   * 将JavaScript对象转换为格式化的JSON字符串
   * @param obj 任意JavaScript对象
   * @param indent 缩进空格数，默认2
   * @returns 格式化后的JSON字符串
   */
  stringify: (obj: unknown, indent: number = 2): string => {
    try {
      return JSON.stringify(obj, null, indent)
    } catch (error) {
      return `序列化失败：${(error as Error).message}`
    }
  },
}

// 简化版核心逻辑（模拟 APIfox 解析过程）
const HTTP_METHODS = new Set([
  'get',
  'post',
  'put',
  'delete',
  'patch',
  'options',
  'head',
  'trace',
])

export function getOperationNodeId(path: string, method: string, opDetail: any): string {
  const customId = opDetail?.operationId || opDetail?.id
  if (typeof customId === 'string' && customId.trim()) {
    return customId
  }
  return `${method.toUpperCase()} ${path}`
}

export function buildApiTree(openapiJson: { paths: any }) {
  const tagMap = new Map() // 存储 { tagName: { label: tagName, children: 接口列表 } }

  // 遍历所有接口路径
  Object.entries(openapiJson.paths || {}).forEach(([path, methods]) => {
    // 遍历每个路径下的请求方法（get/post 等）
    Object.entries(methods || {}).forEach(([method, opDetail]) => {
      if (!HTTP_METHODS.has(method.toLowerCase())) {
        return
      }
      const validTags = Array.isArray(opDetail?.tags)
        ? opDetail.tags
            .filter((tag: any) => typeof tag === 'string' && tag.trim())
            .map((tag: string) => tag.trim())
        : []
      const tags = validTags.length > 0 ? validTags : ['未分组'] // 无 tag 或空 tag 时归为“未分组”
      tags.forEach((tag: any) => {
        // 初始化 tag 节点
        if (!tagMap.has(tag)) {
          tagMap.set(tag, {
            label: tag,
            children: [],
            id: '',
          })
        }
        // 将接口加入 tag 分组
        tagMap.get(tag).children.push({
          label: `${path}`,
          summary: opDetail.summary || '无描述',
          path,
          method: method.toUpperCase(),
          id: getOperationNodeId(path, method, opDetail),
        })
      })
    })
  })

  // 转换为数组，便于渲染树形组件
  return Array.from(tagMap.values())
}
