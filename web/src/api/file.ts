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

// 获取文件信息
export const getFile = (id: number) => {
  return request({
    url: `/v1/files/${id}`,
    method: 'get'
  })
}

// 获取文件预览（用于图片预览）
export const getFilePreview = (id: number) => {
  return request({
    url: `/v1/files/${id}/download`,
    method: 'get',
    responseType: 'blob'
  }).then(response => {
    // 检查响应是否有效
    if (!response.data) {
      throw new Error('获取文件预览失败：无响应数据')
    }

    // 检查响应状态
    if (response.status !== 200) {
      throw new Error(`获取文件预览失败：HTTP ${response.status}`)
    }

    // 获取Content-Type
    const contentType = response.headers['content-type'] || response.headers['Content-Type'] || 'application/octet-stream'

    // 创建Blob URL用于预览
    const blob = new Blob([response.data], { type: contentType })
    const url = window.URL.createObjectURL(blob)

    return url
  }).catch(error => {
    console.error('获取文件预览失败:', error)
    throw error
  })
}

// 下载文件
export const downloadFile = (id: number, fileInfo?: { originalName: string; mimeType: string }) => {
  return request({
    url: `/v1/files/${id}/download`,
    method: 'get',
    responseType: 'blob'
  }).then(response => {
    // 检查响应是否有效
    if (!response.data) {
      throw new Error('下载文件失败：无响应数据')
    }

    // 检查响应状态
    if (response.status !== 200) {
      throw new Error(`下载文件失败：HTTP ${response.status}`)
    }

    // 获取响应中的Content-Type，如果没有则使用文件信息中的mimeType
    let contentType = response.headers['content-type'] || response.headers['Content-Type'] || 'application/octet-stream'

    // 如果响应的Content-Type是通用类型，且我们有文件的真实MIME类型，使用真实的
    if (fileInfo && (contentType === 'application/octet-stream' || contentType === 'application/binary')) {
      contentType = fileInfo.mimeType
      console.log('Using file mimeType instead of response Content-Type:', contentType)
    }

    console.log('Final Content-Type:', contentType)

    // 创建下载链接，使用正确的MIME类型
    const blob = new Blob([response.data], { type: contentType })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url

    // 从响应头获取文件名
    const contentDisposition = response.headers['content-disposition'] || response.headers['Content-Disposition']
    let filename = fileInfo?.originalName || `file_${id}`

    if (contentDisposition) {
      console.log('Content-Disposition:', contentDisposition)
      // 支持多种文件名格式
      const patterns = [
        /filename\*=UTF-8''([^;]+)/i,  // RFC 5987 格式
        /filename="([^"]+)"/i,         // 标准格式
        /filename=([^;]+)/i            // 简单格式
      ]

      for (const pattern of patterns) {
        const match = contentDisposition.match(pattern)
        if (match && match[1]) {
          filename = match[1]
          // 如果是 URL 编码的，解码它
          if (filename.includes('%')) {
            try {
              filename = decodeURIComponent(filename)
            } catch (e) {
              console.warn('Failed to decode filename:', filename)
            }
          }
          break
        }
      }
    }

    // 如果文件名没有扩展名，尝试从Content-Type推断
    if (!filename.includes('.') && contentType !== 'application/octet-stream') {
      const extension = getExtensionFromMimeType(contentType)
      if (extension) {
        filename += extension
      }
    }

    console.log('Download filename:', filename)
    console.log('Blob type:', blob.type)

    // 设置下载属性并触发下载
    link.download = filename
    link.style.display = 'none'
    document.body.appendChild(link)

    // 使用 setTimeout 确保 DOM 操作完成
    setTimeout(() => {
      link.click()
      console.log('Download triggered for:', filename)
    }, 10)

    // 清理资源
    setTimeout(() => {
      if (document.body.contains(link)) {
        document.body.removeChild(link)
      }
      window.URL.revokeObjectURL(url)
    }, 1000)

    return filename
  }).catch(error => {
    console.error('文件下载失败:', error)

    // 检查是否是网络错误或权限错误
    if (error.response) {
      const status = error.response.status
      if (status === 401) {
        throw new Error('未授权：请重新登录')
      } else if (status === 403) {
        throw new Error('权限不足：无法下载该文件')
      } else if (status === 404) {
        throw new Error('文件不存在')
      } else {
        throw new Error(`下载失败：HTTP ${status}`)
      }
    } else {
      throw new Error('网络错误：请检查网络连接')
    }
  })
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