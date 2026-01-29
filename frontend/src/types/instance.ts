import { t } from '@/utils/i18n'
// serverce status
export enum InstanceStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
}
// Type of access
export enum AccessType {
  UNKONWN = 0, // Unkonwn
  DIRECT = 1, // direct
  PROXY = 2, // proxy
  HOSTING = 3, // hosting
}

// MCP Protocol
export enum McpProtocol {
  UNKONWN = 0, // Unkonwn
  SSE = 1, // SSE
  STREAMABLE_HTTP = 2, // STREAMABLE_HTTP
  STDIO = 3, // STDIO
}

// List of Container States
export enum ContainerOptions {
  PENDING = 'pending',
  RUNNING = 'running',
  INIT_TIMEOUT_STOP = 'init-timeout-stop',
  RUN_TIMEOUT_STOP = 'run-timeout-stop',
  EXCEPTION_FORCE_STOP = 'exception-force-stop',
  MANUAL_STOP = 'manual-stop',
  CREATE_FAILED = 'create-failed',
  RUNNING_UNREADY = 'running-unready',
}

// Source of Instance
export enum SourceType {
  UNKONWN = 0, // Unkonwn
  MARKET = 1, // market
  TEMPLATE = 2, // template
  CUSTOM = 3, // custom
  OPENAPI = 4, // openapi
}

// Token Type of Instance
export enum TokenType {
  UNKONWN = 0, // Unkonwn
  BEARER = 1, // Bearer
  API_KEY = 2, // API Key
  X_API_KEY = 3, // X-API-KEY
  BASIC = 4, // Basic
}

// the default of instance
export const InstanceData = computed(() => {
  return {
    // IMGADDRESS = 'ccr.ccs.tencentyun.com/itqm-private/mcp-hosting:v2', // imgAddress
    IMGADDRESS: 'ccr.ccs.tencentyun.com/itqm-private/mcp-hosting:v2.1', // imgAddress
    PORT: 8080, // port
    INITSCRIPT: `#!/bin/bash
# ${t('desc.initScript')}
echo 'Initialization completed.'
`, // InitScript
    TIP_IMGADDRESS: `
  1.${t('desc.imgAddressTip1')}<br />
  2. ${t('desc.imgAddressTip2')}<br />
    &nbsp &nbsp ${t('desc.baseCommand')}：tar、wget、zip、unzip<br />
    &nbsp &nbsp Python ${t('desc.env')}：Python 3.12.11，${t('desc.tools')} uv 0.7.12、uvx 0.7.12<br />
    &nbsp &nbsp Node.js ${t('desc.env')}：Node.js v18.20.1，${t('desc.tools')} npm 9.6.6、npx 9.6.6<br />
  `,
    TIP_IMGADDRESS_DEFAULT: `3.${t('desc.imgAddressTip3')} `,
    TIP_MCP_SERVER: `
    {
      "mcpServers": {
        "everything": {
          "args": [
            "-y",
            "@modelcontextprotocol/server-everything"
          ],
          "command": "npx"
        }
      }
    }
  `, // placeholderServer
    PACKAGE_PATH: '/app/codepkg/', // the default of package path
    COMMAND_TIP: 'mcp-hosting --port=%d --mcp-servers-config /app/mcp-servers.json', // default start command
  }
})

export enum NodeVisible {
  RWO = 'ReadWriteOnce',
  ROM = 'ReadOnlyMany',
  RWM = 'ReadWriteMany',
}

export interface VolumeMountsItme {
  type: string
  nodeName: string
  hostPath: string
  mountPath: string
  pvcName: string
  readOnly: boolean
}

// base-form
export interface InstanceForm {
  instanceName: string
  accessType: AccessType
  mcpProtocol: McpProtocol
  notes: string
  mcpConfig: string
  environmentId: string
  iconPath: string
  servicePath: string
  enabledToken: boolean
}

// creat-instance-form
export interface InstanceCreate extends InstanceForm {
  sourceType: SourceType
  name: string
  notes: string
  mcpServers: string
  packageId: string
  environmentId: string
  port: number
  environmentVariables: { key: string; value: string }[]
  volumeMounts: { [key: string]: any }[]
  initScript: string
  command: string
}

// list-result-form
export interface InstanceResult extends InstanceForm {
  instanceId: string
  containerName: string
  containerStatus: ContainerOptions
  status: InstanceStatus
  publicProxyConfig: string
  createdAt: string
  environmentName: string
  containerIsReady: boolean
  tokens: Array<{
    enabled: boolean
    token: string
    expireAt: number
    publishAt: number
    headers: { [key: string]: string }
    usages: string[]
    tokenType: TokenType | null
  }>
  [key: string]: any
}

// template-form-by-instance
export interface TemplateForm extends InstanceForm {
  name: string
}

// template-result-by-instance
export interface TemplateResult extends InstanceForm {
  templateId: string
  name: string
  environmentName: string
  createdAt: string
  [key: string]: any
}
