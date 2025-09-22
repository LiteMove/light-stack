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

// 下载文件
export const downloadFile = (id: number) => {
  return request({
    url: `/v1/files/${id}/download`,
    method: 'get',
    responseType: 'blob'
  }).then(response => {
    // 创建下载链接
    const blob = new Blob([response.data])
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url

    // 从响应头获取文件名
    const contentDisposition = response.headers['content-disposition']
    let filename = 'download'
    if (contentDisposition) {
      const filenameMatch = contentDisposition.match(/filename=(.+)/)
      if (filenameMatch) {
        filename = filenameMatch[1].replace(/"/g, '')
      }
    }

    link.download = filename
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  })
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