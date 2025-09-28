package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/LiteMove/light-stack/internal/modules/files/service"
	"github.com/LiteMove/light-stack/internal/shared/middleware"
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

	// 获取是否公开（如果未指定，则使用租户配置的默认值）
	isPublicStr := c.PostForm("isPublic")
	var isPublic bool
	if isPublicStr != "" {
		isPublic = isPublicStr == "true"
	} else {
		// 使用租户配置的默认公开设置，由FileService处理
		isPublic = false // 这里先设为false，在FileService中会应用租户默认设置
	}

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

// GetPrivateFile 获取私有文件内容（需要权限验证）
func (fc *FileController) GetPrivateFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的文件ID")
		return
	}

	// 获取当前用户和租户信息
	userID := middleware.GetUserIDFromContext(c)
	tenantID, _ := middleware.GetTenantIDFromContext(c)

	// 获取文件信息并验证权限
	file, fileContent, err := fc.fileService.GetPrivateFileContent(id, userID, tenantID)
	if err != nil {
		if err.Error() == "file not found" {
			response.Error(c, http.StatusNotFound, "文件不存在")
		} else if err.Error() == "access denied" {
			response.Error(c, http.StatusForbidden, "无权访问此文件")
		} else {
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// 设置响应头
	c.Header("Content-Type", file.MimeType)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", file.OriginalName))
	c.Header("Content-Length", fmt.Sprintf("%d", len(fileContent)))

	// 返回文件内容
	c.Data(http.StatusOK, file.MimeType, fileContent)
}

// 下载、预览、复制链接功能已移除，统一使用文件的 access_url 字段
// 客户端可以直接使用 GetFile 接口获取文件信息，然后使用 access_url 进行操作

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
				"storageType":  profile.StorageType,
				"isPublic":     profile.IsPublic,
				"accessUrl":    profile.AccessURL,
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
