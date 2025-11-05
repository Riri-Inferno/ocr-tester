package utils

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// FileType
type FileType int

const (
    FileTypeUnknown FileType = iota
    FileTypePDF
    FileTypeJPEG
    FileTypePNG
    FileTypeGIF
    FileTypeWebP
)

// MaxFileSize
const MaxFileSize = 20 * 1024 * 1024

// MaxPDFPages
const MaxPDFPages = 5

// ファイルタイプを検出
func DetectFileType(data []byte) (FileType, string) {
    if len(data) < 512 {
        return FileTypeUnknown, ""
    }

    // PDFチェック
    if bytes.HasPrefix(data, []byte("%PDF")) {
        return FileTypePDF, "application/pdf"
    }

    // JPEGチェック
    if bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF}) {
        return FileTypeJPEG, "image/jpeg"
    }

    // PNGチェック
    if bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47}) {
        return FileTypePNG, "image/png"
    }

    // GIFチェック
    if bytes.HasPrefix(data, []byte("GIF87a")) || bytes.HasPrefix(data, []byte("GIF89a")) {
        return FileTypeGIF, "image/gif"
    }

    // WebPチェック
    if len(data) > 12 && bytes.HasPrefix(data[8:12], []byte("WEBP")) {
        return FileTypeWebP, "image/webp"
    }

    return FileTypeUnknown, ""
}

// ファイルのバリデーション
func ValidateFile(data []byte, filename string) error {
    // サイズチェック
    if len(data) > MaxFileSize {
        return fmt.Errorf("file size exceeds maximum allowed size of %d MB", MaxFileSize/(1024*1024))
    }

    // ファイルタイプチェック
    fileType, _ := DetectFileType(data)
    if fileType == FileTypeUnknown {
        ext := ""
        if idx := strings.LastIndex(filename, "."); idx >= 0 {
            ext = strings.ToLower(filename[idx+1:])
        }
        return fmt.Errorf("unsupported file type: %s", ext)
    }

    // PDFの場合はページ数チェック
    if fileType == FileTypePDF {
        pageCount, err := GetPDFPageCount(data)
        if err != nil {
            return fmt.Errorf("failed to read PDF: %w", err)
        }
        if pageCount > MaxPDFPages {
            return fmt.Errorf("PDF page count (%d) exceeds maximum allowed pages (%d)", pageCount, MaxPDFPages)
        }
    }

    return nil
}

// PDFのページ数を取得
func GetPDFPageCount(data []byte) (int, error) {
    reader := bytes.NewReader(data)
    // nilを渡すとデフォルト設定が使用される
    ctx, err := api.ReadContext(reader, nil)
    if err != nil {
        return 0, fmt.Errorf("failed to read PDF context: %w", err)
    }
    return ctx.PageCount, nil
}

// PDF全体を画像配列に変換
func ConvertPDFToImages(pdfData []byte) ([][]byte, error) {
    pageCount, err := GetPDFPageCount(pdfData)
    if err != nil {
        return nil, err
    }

    if pageCount > MaxPDFPages {
        pageCount = MaxPDFPages
    }

    // PDFを直接Base64エンコードして送信
    // Gemini APIがPDFを直接サポートしているため画像変換は不要
    images := [][]byte{pdfData}
    
    return images, nil
}

// 必要に応じて画像をリサイズ
func ResizeImageIfNeeded(imageData []byte) ([]byte, error) {
    // 画像サイズが大きすぎる場合のリサイズ処理
    // 現時点では、Gemini APIの制限内であればそのまま返す
    return imageData, nil
}

// OCR処理用にファイルを準備
func ProcessFileForOCR(data []byte) ([]byte, string, error) {
    fileType, mimeType := DetectFileType(data)
    
    switch fileType {
    case FileTypePDF:
        // PDFはそのままBase64エンコードして送る
        return data, "application/pdf", nil
    case FileTypeJPEG, FileTypePNG, FileTypeGIF, FileTypeWebP:
        // 画像もそのまま送る
        return data, mimeType, nil
    default:
        return nil, "", fmt.Errorf("unsupported file type")
    }
}
