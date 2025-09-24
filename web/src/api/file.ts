import request from '@/utils/request'

export interface FileProfile {
  id: number
  tenantId: number
  originalName: string
  fileName: string
  filePath: string
  fileSize: number
  fileType: string
  mimeType: string
  md5: string
  uploadUserId: number
  usageType: string
  storageType: string
  isPublic: boolean
  accessUrl: string
  createdAt: string
  updatedAt: string
}

export interface FileListResponse {
  files: FileProfile[]
  pagination: {
    page: number
    pageSize: number
    total: number
  }
}

// 上传文件
export const uploadFile = (formData: FormData) => {
  return request({
    url: '/v1/files/upload',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 获取文件信息（包含access_url用于下载、预览、复制链接）
export const getFile = (id: number) => {
  return request({
    url: `/v1/files/${id}`,
    method: 'get'
  })
}

// 获取私有文件内容（用于预览和下载）
export const getPrivateFileContent = (id: number) => {
  return request({
    url: `/v1/files/${id}/private`,
    method: 'get',
    responseType: 'blob'  // 重要：设置响应类型为blob以获取文件内容
  })
}

// 下载私有文件
export const downloadPrivateFile = async (id: number, originalName: string) => {
  try {
    const res = await getPrivateFileContent(id)
    const blob = new Blob([res.data], { type: 'application/octet-stream' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = originalName
    link.style.display = 'none'
    document.body.appendChild(link)

    // 触发下载
    setTimeout(() => {
      link.click()
      console.log('Private file download triggered for:', originalName)
    }, 10)

    // 清理DOM和URL对象
    setTimeout(() => {
      if (document.body.contains(link)) {
        document.body.removeChild(link)
      }
      window.URL.revokeObjectURL(url)
    }, 1000)

    return true
  } catch (err) {
    console.error('私有文件下载失败:', err)
    return false
  }
}
export const downloadFileByUrl = (accessUrl: string, originalName: string) => {
  const link = document.createElement('a')
  link.href = accessUrl
  link.download = originalName
  link.style.display = 'none'
  document.body.appendChild(link)
  
  // 触发下载
  setTimeout(() => {
    link.click()
    console.log('Download triggered for:', originalName)
  }, 10)
  
  // 清理DOM
  setTimeout(() => {
    if (document.body.contains(link)) {
      document.body.removeChild(link)
    }
  }, 1000)
}

// 复制文件链接到剪贴板
export const copyFileUrl = async (accessUrl: string) => {
  try {
    await navigator.clipboard.writeText(accessUrl)
    return true
  } catch (err) {
    console.error('复制链接失败:', err)
    // 降级方案：使用传统方法
    const textArea = document.createElement('textarea')
    textArea.value = accessUrl
    textArea.style.position = 'fixed'
    textArea.style.left = '-999999px'
    textArea.style.top = '-999999px'
    document.body.appendChild(textArea)
    textArea.focus()
    textArea.select()
    
    try {
      document.execCommand('copy')
      document.body.removeChild(textArea)
      return true
    } catch (err) {
      document.body.removeChild(textArea)
      return false
    }
  }
}

// 预览图片文件（直接使用access_url）
export const previewImage = (accessUrl: string) => {
  return accessUrl // 直接返回URL，用于img标签的src
}

// 根据MIME类型获取文件扩展名
const getExtensionFromMimeType = (mimeType: string): string => {
  const mimeToExt: Record<string, string> = {
    // 图片类型
    'image/jpeg': '.jpg',
    'image/jpg': '.jpg',
    'image/png': '.png',
    'image/gif': '.gif',
    'image/webp': '.webp',
    'image/svg+xml': '.svg',
    'image/bmp': '.bmp',
    'image/tiff': '.tiff',

    // 文档类型
    'application/pdf': '.pdf',
    'application/msword': '.doc',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document': '.docx',
    'application/vnd.ms-excel': '.xls',
    'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet': '.xlsx',
    'application/vnd.ms-powerpoint': '.ppt',
    'application/vnd.openxmlformats-officedocument.presentationml.presentation': '.pptx',

    // 文本类型
    'text/plain': '.txt',
    'text/html': '.html',
    'text/css': '.css',
    'text/javascript': '.js',
    'application/json': '.json',
    'text/xml': '.xml',
    'application/xml': '.xml',

    // 压缩文件
    'application/zip': '.zip',
    'application/x-rar-compressed': '.rar',
    'application/x-7z-compressed': '.7z',
    'application/x-tar': '.tar',
    'application/gzip': '.gz',

    // 音频
    'audio/mpeg': '.mp3',
    'audio/wav': '.wav',
    'audio/ogg': '.ogg',
    'audio/mp4': '.m4a',

    // 视频
    'video/mp4': '.mp4',
    'video/avi': '.avi',
    'video/mov': '.mov',
    'video/wmv': '.wmv',
    'video/webm': '.webm'
  }

  return mimeToExt[mimeType] || ''
}

// 删除文件
export const deleteFile = (id: number) => {
  return request({
    url: `/v1/files/${id}`,
    method: 'delete'
  })
}

// 获取用户文件列表
export const getUserFiles = (params: {
  page?: number
  pageSize?: number
  filename?: string
  fileType?: string
  usageType?: string
}) => {
  return request({
    url: '/v1/files',
    method: 'get',
    params
  })
}

// 获取所有文件列表（管理员）
export const getAllFiles = (params: {
  page?: number
  pageSize?: number
  filename?: string
  fileType?: string
  usageType?: string
  uploadUserId?: number
}) => {
  return request({
    url: '/v1/files',
    method: 'get',
    params
  })
}