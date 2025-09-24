package controller

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/LiteMove/light-stack/internal/middleware"
	"github.com/LiteMove/light-stack/internal/service"
	"github.com/LiteMove/light-stack/pkg/response"
	"github.com/gin-gonic/gin"
)

// FileController 文件控制器
type FileController struct {
	fileService *service.FileService
}

// NewFileController 创建文件控制器实例
func NewFileController(fileService *service.FileService) *FileController {
	return &FileController{
		fileService: fileService,
	}
}

// UploadFile 上传文件（支持新的存储架构）
func (fc *FileController) UploadFile(c *gin.Context) {
	// 获取当前用户和租户信息
	userID := middleware.GetUserIDFromContext(c)
	tenantID, _ := middleware.GetTenantIDFromContext(c)

	// 获取文件
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "文件上传失败")
		return
	}

	// 获取使用类型（可选）
	usageType := c.PostForm("usageType")

	// 获取是否公开（可选，默认为false）
	isPublic := c.PostForm("isPublic") == "true"

	// 上传文件（现在由FileService根据租户配置处理所有验证）
	uploadedFile, err := fc.fileService.UploadFile(file, userID, tenantID, usageType, isPublic)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, uploadedFile.ToProfile())
}

// GetFile 获取文件信息
func (fc *FileController) GetFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文件ID")
		return
	}

	file, err := fc.fileService.GetFileByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "文件不存在")
		return
	}

	response.Success(c, file.ToProfile())
}

// DownloadFile 下载文件（支持新的存储架构）
func (fc *FileController) DownloadFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文件ID")
		return
	}

	// 获取当前用户ID
	userID := middleware.GetUserIDFromContext(c)

	// 获取文件下载URL
	downloadURL, err := fc.fileService.GetDownloadURL(id, userID)
	if err != nil {
		if err.Error() == "access denied" {
			response.Error(c, http.StatusForbidden, "访问被拒绝")
			return
		}
		if err.Error() == "file not found: record not found" {
			response.Error(c, http.StatusNotFound, "文件不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "获取下载链接失败")
		return
	}

	// 如果是外部URL（OSS等），重定向到该URL
	if strings.HasPrefix(downloadURL, "http://") || strings.HasPrefix(downloadURL, "https://") {
		c.Redirect(http.StatusFound, downloadURL)
		return
	}

	// 对于本地存储的私有文件，需要通过认证后下载
	// 这里downloadURL应该是类似 /static/private/tenant_1/2024/01/01/file.jpg 的格式
	// 我们需要将其转换为实际的文件路径进行下载
	file, err := fc.fileService.GetFileByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "文件不存在")
		return
	}

	// 如果是本地存储，直接提供文件下载
	if file.StorageType == "local" {
		// 设置下载头
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename*=UTF-8''"+url.QueryEscape(file.OriginalName))
		c.Header("Content-Type", file.MimeType)
		c.Header("Content-Length", strconv.FormatInt(file.FileSize, 10))

		// 构建实际文件路径
		// file.FilePath 格式类似：private/tenant_1/2024/01/01/filename.ext
		actualPath := filepath.Join("uploads", file.FilePath)

		// 检查文件是否存在
		if _, err := os.Stat(actualPath); os.IsNotExist(err) {
			response.Error(c, http.StatusNotFound, "文件已被删除")
			return
		}

		c.File(actualPath)
		return
	}

	// OSS等外部存储，返回下载URL
	response.Success(c, gin.H{"downloadUrl": downloadURL})
}

// GetDownloadURL 获取文件下载URL（新增API）
func (fc *FileController) GetDownloadURL(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文件ID")
		return
	}

	// 获取当前用户ID
	userID := middleware.GetUserIDFromContext(c)

	// 获取文件下载URL
	downloadURL, err := fc.fileService.GetDownloadURL(id, userID)
	if err != nil {
		if err.Error() == "access denied" {
			response.Error(c, http.StatusForbidden, "访问被拒绝")
			return
		}
		if err.Error() == "file not found: record not found" {
			response.Error(c, http.StatusNotFound, "文件不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "获取下载链接失败")
		return
	}

	response.Success(c, gin.H{"downloadUrl": downloadURL})
}

// DeleteFile 删除文件
func (fc *FileController) DeleteFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文件ID")
		return
	}

	err = fc.fileService.DeleteFile(id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "删除文件失败")
		return
	}

	response.Success(c, gin.H{"message": "文件删除成功"})
}

// GetUserFiles 获取用户文件列表
func (fc *FileController) GetUserFiles(c *gin.Context) {
	// 获取当前用户和租户信息
	userID := middleware.GetUserIDFromContext(c)
	tenantID, _ := middleware.GetTenantIDFromContext(c)

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	// 获取文件列表
	files, total, err := fc.fileService.GetFilesByUser(userID, tenantID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取文件列表失败")
		return
	}

	// 转换为Profile格式
	profiles := make([]interface{}, len(files))
	for i, file := range files {
		profiles[i] = file.ToProfile()
	}

	response.Success(c, gin.H{
		"files": profiles,
		"pagination": gin.H{
			"page":     page,
			"pageSize": pageSize,
			"total":    total,
		},
	})
}

// GetAllFiles 获取所有文件列表（按租户）
func (fc *FileController) GetAllFiles(c *gin.Context) {
	// 获取租户信息
	tenantID, _ := middleware.GetTenantIDFromContext(c)

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 20
	}

	// 获取过滤参数
	filters := make(map[string]interface{})
	if filename := c.Query("filename"); filename != "" {
		filters["filename"] = filename
	}
	if fileType := c.Query("fileType"); fileType != "" {
		filters["file_type"] = fileType
	}
	if usageType := c.Query("usageType"); usageType != "" {
		filters["usage_type"] = usageType
	}
	if uploadUserID := c.Query("uploadUserId"); uploadUserID != "" {
		filters["upload_user_id"] = uploadUserID
	}

	// 获取文件列表
	files, total, err := fc.fileService.GetAllFiles(tenantID, page, pageSize, filters)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取文件列表失败")
		return
	}

	// 转换为Profile格式
	profiles := make([]interface{}, len(files))
	for i, file := range files {
		profile := file.ToProfile()
		// 添加上传用户信息
		if file.UploadUser != nil {
			profiles[i] = gin.H{
				"id":           profile.ID,
				"tenantId":     profile.TenantID,
				"originalName": profile.OriginalName,
				"fileName":     profile.FileName,
				"filePath":     profile.FilePath,
				"fileSize":     profile.FileSize,
				"fileType":     profile.FileType,
				"mimeType":     profile.MimeType,
				"md5":          profile.MD5,
				"uploadUserId": profile.UploadUserID,
				"usageType":    profile.UsageType,
				"createdAt":    profile.CreatedAt,
				"updatedAt":    profile.UpdatedAt,
				"uploadUser": gin.H{
					"id":       file.UploadUser.ID,
					"username": file.UploadUser.Username,
					"nickname": file.UploadUser.Nickname,
					"avatar":   file.UploadUser.Avatar,
				},
			}
		} else {
			profiles[i] = profile
		}
	}

	response.Success(c, gin.H{
		"files": profiles,
		"pagination": gin.H{
			"page":     page,
			"pageSize": pageSize,
			"total":    total,
		},
	})
}
