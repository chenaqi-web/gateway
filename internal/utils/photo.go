package utils

import (
	"errors"
	"fmt"
	"gateway/internal/config"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

// 定义错误变量，方便调用方用 errors.Is 判断
var (
	ErrNoFileProvided   = errors.New("no file provided")
	ErrFileTooLarge     = errors.New("file too large")
	ErrInvalidImageType = errors.New("invalid image type")
)

// 允许的图片扩展名列表（可配置化）
var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	//".gif":  true,
	//".webp": true,
	//".avif": true,
}

// ValidateImage 验证上传的图片文件是否合法
// 检查项包括：文件是否存在、大小是否超限、扩展名是否允许、MIME类型是否为图片
func ValidateImage(cfg *config.Config, file *multipart.FileHeader) error {
	// 1. 检查文件是否存在
	if file == nil {
		return ErrNoFileProvided
	}

	// 2. 检查文件大小是否在允许范围内
	if file.Size <= 0 || file.Size > cfg.Storage.MaxSize {
		return fmt.Errorf("%w: max %d bytes, got %d bytes",
			ErrFileTooLarge, cfg.Storage.MaxSize, file.Size)
	}

	// 3. 检查文件扩展名是否在允许列表中
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return fmt.Errorf("%w: extension %s not allowed",
			ErrInvalidImageType, ext)
	}

	// 4. 打开文件，读取前512字节用于检测MIME类型
	source, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer source.Close()

	// 读取文件头（前512字节），http.DetectContentType 通过文件头魔数判断真实类型
	header := make([]byte, 512)
	count, err := source.Read(header)
	if err != nil {
		return fmt.Errorf("failed to read file header: %w", err)
	}
	if count < 512 {
		// 文件太小，无法可靠检测类型
		return fmt.Errorf("%w: file too small for type detection", ErrInvalidImageType)
	}

	// 5. 检测真实MIME类型
	contentType := http.DetectContentType(header[:count])
	if !strings.HasPrefix(contentType, "image/") {
		return fmt.Errorf("%w: detected MIME type %s does not match image/*",
			ErrInvalidImageType, contentType)
	}

	return nil
}

// IsValidImageExtension 检查扩展名是否允许
func IsValidImageExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return allowedExtensions[ext]
}
