package controller

import (
	"net/http"
	"path/filepath"
	"strconv"

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

// UploadFile 上传文件
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

	// 文件大小限制（50MB）
	const maxFileSize = 50 << 20
	if file.Size > maxFileSize {
		response.Error(c, http.StatusBadRequest, "文件大小不能超过50MB")
		return
	}

	// 检查文件类型
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".txt"}
	ext := filepath.Ext(file.Filename)
	allowed := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			allowed = true
			break
		}
	}
	if !allowed {
		response.Error(c, http.StatusBadRequest, "不支持的文件类型")
		return
	}

	// 上传文件
	uploadedFile, err := fc.fileService.UploadFile(file, userID, tenantID, usageType)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "文件上传失败")
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

// DownloadFile 下载文件
func (fc *FileController) DownloadFile(c *gin.Context) {
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

	// 设置下载头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+file.OriginalName)
	c.Header("Content-Type", "application/octet-stream")

	c.File(file.FilePath)
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

// GetAllFiles 获取所有文件列表（管理员功能）
func (fc *FileController) GetAllFiles(c *gin.Context) {
	// 获取租户信息
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

	// 获取过滤参数
	filters := make(map[string]interface{})
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
