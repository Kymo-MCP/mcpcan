import { da } from 'element-plus/es/locales.mjs'
import { useRouter } from 'vue-router'

/**
 * 对象转url字符串参数
 */
export const obj2UrlString = (obj: any) => {
  if (!obj) {
    return ''
  }
  // 使用 Object.entries 获取对象的键值对数组，然后使用 filter 方法过滤掉值为空字符串的项
  const filteredEntries = Object.entries(obj).filter(
    ([, value]) => value !== '' && value !== undefined,
  )
  // 使用 map 方法将每个键值对转换成 'key=value' 的形式，然后使用 join 方法将它们连接成一个字符串
  const queryString = filteredEntries.map(([key, value]) => `${key.trim()}=${value}`).join('&')

  // 如果 queryString 不为空，则在其前面添加 '?'，否则返回空字符串
  return queryString ? `?${queryString}` : ''
}

interface jumpType {
  url: string
  data?: Record<string, any>
  isOpen?: boolean
  isCurrentProject?: boolean
  jumpWithTenantId?: boolean
}
/**
 * Router hooks - must be called inside setup() or functional components
 * @returns
 */
export const useRouterHooks = () => {
  // 在函数调用时（setup内部）才获取 router 实例
  const router = useRouter()

  const jumpToPage = (params: jumpType) => {
    return new Promise((resolve: any) => {
      const { url, data = {}, isOpen = false } = params

      // 如果当前路由存在 layout 参数并且目标没有 layout，则带上当前 layout
      let paramsData = {
        layout: router.currentRoute.value.query?.layout,
        ...data,
      }

      if (isOpen) {
        // 如果是绝对 URL（以 http(s) 或 // 开头），直接打开
        const isAbsolute = /^(https?:)?\/\//i.test(url)
        const queryStr = obj2UrlString(paramsData)
        let openUrl = url

        if (!isAbsolute) {
          try {
            // 使用 Vite 的 base（可能为 '/mcpcan-web/' 或 './'）构造 hash 模式下的完整地址
            const base = (import.meta as any).env?.BASE_URL || '/'
            const origin = window.location.origin
            let nb = String(base)
            if (!nb.startsWith('/')) nb = '/' + nb
            if (!nb.endsWith('/')) nb = nb + '/'

            // 如果传入的 url 包含 '#'，认为调用方已经提供了 hash 部分
            if (url.includes('#')) {
              // 例如 url = '/mcpcan-web/#/path' 或 '#/path'
              // 保持原样，仅添加查询字符串
              // 若 url 是相对路径（以 / 开头），补全 origin
              if (/^\//.test(url)) {
                openUrl = origin + url + queryStr
              } else {
                openUrl = url + queryStr
              }
            } else {
              // 将 path 归一化（去掉前导斜杠）并放入 hash 路由
              const p = url.startsWith('/') ? url.slice(1) : url
              openUrl = `${origin}${nb}#/${p}${queryStr}`
            }
          } catch (err) {
            // 任何异常回退到简单拼接
            openUrl = url + queryStr
          }
        } else {
          openUrl = url + queryStr
        }

        window.open(openUrl)
      } else {
        router
          .push({
            path: url,
            query: {
              ...paramsData,
            },
          })
          .then(() => {
            resolve()
          })
      }
    })
  }

  const jumpBack = () => {
    router.go(-1)
  }

  const reload = () => {
    window.location.reload()
  }

  return {
    jumpToPage,
    jumpBack,
    reload,
  }
}
