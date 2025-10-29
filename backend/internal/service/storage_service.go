package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vayura/backend/config"
)

type StorageService interface {
	SaveAvatar(ctx context.Context, userID uint, file *multipart.FileHeader) (string, error)
}

type storageService struct {
	cfg *config.Config
}

func NewStorageService(cfg *config.Config) StorageService {
	return &storageService{cfg: cfg}
}

func (s *storageService) SaveAvatar(ctx context.Context, userID uint, file *multipart.FileHeader) (string, error) {
	uploadPath := s.cfg.Storage.UploadDir
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Validasi ukuran file
	if file.Size > 2*1024*1024 {
		return "", fmt.Errorf("file too large (max 2MB)")
	}

	// Validasi ekstensi file
	ext := strings.ToLower(filepath.Ext(file.Filename))
	validExt := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	if !validExt[ext] {
		return "", fmt.Errorf("invalid file type (only jpg, jpeg, png allowed)")
	}

	// Generate nama file unik
	newFileName := fmt.Sprintf("%d_%d%s", userID, time.Now().Unix(), ext)
	fullPath := filepath.Join(uploadPath, newFileName)

	// Simpan file
	if err := saveUploadedFile(file, fullPath); err != nil {
		return "", err
	}

	return "/" + fullPath, nil
}

func saveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.ReadFrom(src)
	return err
}
