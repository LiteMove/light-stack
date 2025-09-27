package generator

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FilePackager 文件打包器
type FilePackager struct {
	outputDir string
}

// NewFilePackager 创建文件打包器
func NewFilePackager(outputDir string) *FilePackager {
	return &FilePackager{
		outputDir: outputDir,
	}
}

// PackageToZip 打包文件到ZIP
func (p *FilePackager) PackageToZip(result *GenerateResult, zipFileName string) (string, error) {
	// 确保输出目录存在
	if err := os.MkdirAll(p.outputDir, 0755); err != nil {
		return "", fmt.Errorf("创建输出目录失败: %v", err)
	}

	// 生成完整的ZIP文件路径
	zipPath := filepath.Join(p.outputDir, zipFileName)

	// 创建ZIP文件
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", fmt.Errorf("创建ZIP文件失败: %v", err)
	}
	defer zipFile.Close()

	// 创建ZIP写入器
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 添加文件到ZIP
	for fileName, content := range result.Files {
		if err := p.addFileToZip(zipWriter, fileName, content); err != nil {
			return "", fmt.Errorf("添加文件 %s 到ZIP失败: %v", fileName, err)
		}
	}

	// 添加README文件
	readmeContent := p.generateReadme(result)
	if err := p.addFileToZip(zipWriter, "README.md", readmeContent); err != nil {
		return "", fmt.Errorf("添加README文件失败: %v", err)
	}

	return zipPath, nil
}

// addFileToZip 添加文件到ZIP
func (p *FilePackager) addFileToZip(zipWriter *zip.Writer, fileName, content string) error {
	fileWriter, err := zipWriter.Create(fileName)
	if err != nil {
		return fmt.Errorf("创建ZIP文件条目失败: %v", err)
	}

	_, err = io.WriteString(fileWriter, content)
	if err != nil {
		return fmt.Errorf("写入文件内容失败: %v", err)
	}

	return nil
}

// PackageToFiles 解压文件到目录
func (p *FilePackager) PackageToFiles(result *GenerateResult, targetDir string) error {
	// 确保目标目录存在
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %v", err)
	}

	// 保存每个文件
	for fileName, content := range result.Files {
		fullPath := filepath.Join(targetDir, fileName)

		// 确保文件目录存在
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建文件目录 %s 失败: %v", dir, err)
		}

		// 写入文件
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("写入文件 %s 失败: %v", fullPath, err)
		}
	}

	// 生成README文件
	readmePath := filepath.Join(targetDir, "README.md")
	readmeContent := p.generateReadme(result)
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("写入README文件失败: %v", err)
	}

	return nil
}

// CreateZipFromBytes 从字节创建ZIP
func (p *FilePackager) CreateZipFromBytes(files map[string][]byte) ([]byte, error) {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)
	defer zipWriter.Close()

	for fileName, content := range files {
		fileWriter, err := zipWriter.Create(fileName)
		if err != nil {
			return nil, fmt.Errorf("创建ZIP文件条目失败: %v", err)
		}

		if _, err := fileWriter.Write(content); err != nil {
			return nil, fmt.Errorf("写入文件内容失败: %v", err)
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("关闭ZIP写入器失败: %v", err)
	}

	return buf.Bytes(), nil
}

// ExtractZip 解压ZIP文件
func (p *FilePackager) ExtractZip(zipPath, extractDir string) error {
	// 打开ZIP文件
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("打开ZIP文件失败: %v", err)
	}
	defer reader.Close()

	// 确保解压目录存在
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		return fmt.Errorf("创建解压目录失败: %v", err)
	}

	// 解压每个文件
	for _, file := range reader.File {
		if err := p.extractFile(file, extractDir); err != nil {
			return fmt.Errorf("解压文件 %s 失败: %v", file.Name, err)
		}
	}

	return nil
}

// extractFile 解压单个文件
func (p *FilePackager) extractFile(file *zip.File, extractDir string) error {
	// 构建完整路径
	filePath := filepath.Join(extractDir, file.Name)

	// 检查路径安全性
	if !strings.HasPrefix(filePath, filepath.Clean(extractDir)+string(os.PathSeparator)) {
		return fmt.Errorf("非法文件路径: %s", file.Name)
	}

	// 如果是目录，创建目录
	if file.FileInfo().IsDir() {
		return os.MkdirAll(filePath, file.FileInfo().Mode())
	}

	// 确保文件目录存在
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("创建文件目录失败: %v", err)
	}

	// 打开ZIP文件中的文件
	fileReader, err := file.Open()
	if err != nil {
		return fmt.Errorf("打开ZIP文件中的文件失败: %v", err)
	}
	defer fileReader.Close()

	// 创建目标文件
	targetFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.FileInfo().Mode())
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer targetFile.Close()

	// 复制文件内容
	if _, err := io.Copy(targetFile, fileReader); err != nil {
		return fmt.Errorf("复制文件内容失败: %v", err)
	}

	return nil
}

// GetFileSize 获取文件大小
func (p *FilePackager) GetFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, fmt.Errorf("获取文件信息失败: %v", err)
	}
	return fileInfo.Size(), nil
}

// FileExists 检查文件是否存在
func (p *FilePackager) FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// DeleteFile 删除文件
func (p *FilePackager) DeleteFile(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}
	return nil
}

// CleanupOldFiles 清理旧文件
func (p *FilePackager) CleanupOldFiles(maxAge int64) error {
	// 遍历输出目录
	return filepath.Walk(p.outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查文件年龄
		if info.ModTime().Unix() < maxAge {
			// 删除过期文件
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("删除过期文件 %s 失败: %v", path, err)
			}
		}

		return nil
	})
}

// generateReadme 生成README文件
func (p *FilePackager) generateReadme(result *GenerateResult) string {
	var builder strings.Builder

	builder.WriteString("# 代码生成结果\n\n")
	builder.WriteString(fmt.Sprintf("**表名**: %s\n", result.TableName))
	builder.WriteString(fmt.Sprintf("**业务名称**: %s\n", result.BusinessName))
	builder.WriteString(fmt.Sprintf("**文件数量**: %d\n", result.FileCount))
	builder.WriteString(fmt.Sprintf("**总大小**: %.2f KB\n", float64(result.TotalSize)/1024))
	builder.WriteString(fmt.Sprintf("**生成时间**: %s\n", result.StartTime.Format("2006-01-02 15:04:05")))
	builder.WriteString(fmt.Sprintf("**耗时**: %s\n\n", result.Duration))

	builder.WriteString("## 文件列表\n\n")
	for fileName := range result.Files {
		builder.WriteString(fmt.Sprintf("- %s\n", fileName))
	}

	builder.WriteString("\n## 使用说明\n\n")
	builder.WriteString("1. 将生成的文件复制到对应的项目目录中\n")
	builder.WriteString("2. 根据项目需要调整代码和配置\n")
	builder.WriteString("3. 运行项目进行测试\n\n")

	builder.WriteString("## 注意事项\n\n")
	builder.WriteString("- 生成的代码仅供参考，请根据实际需要进行调整\n")
	builder.WriteString("- 建议在使用前备份现有代码\n")
	builder.WriteString("- 如有问题，请联系开发人员\n")

	return builder.String()
}

// PackageInfo 打包信息
type PackageInfo struct {
	FileName   string `json:"fileName"`   // 文件名
	FilePath   string `json:"filePath"`   // 文件路径
	FileSize   int64  `json:"fileSize"`   // 文件大小
	CreateTime string `json:"createTime"` // 创建时间
	ExpiryTime string `json:"expiryTime"` // 过期时间
}

// GetPackageInfo 获取打包信息
func (p *FilePackager) GetPackageInfo(filePath string) (*PackageInfo, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %v", err)
	}

	return &PackageInfo{
		FileName:   fileInfo.Name(),
		FilePath:   filePath,
		FileSize:   fileInfo.Size(),
		CreateTime: fileInfo.ModTime().Format("2006-01-02 15:04:05"),
		ExpiryTime: fileInfo.ModTime().Add(24 * 7 * time.Hour).Format("2006-01-02 15:04:05"), // 7天后过期
	}, nil
}
