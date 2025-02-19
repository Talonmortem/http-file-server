package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Talonmortem/http-file-server/internal/database"
	"github.com/gin-gonic/gin"
)

type FileInfo struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Size     string `json:"size"`
	Modified string `json:"modified"`
	IsDir    bool   `json:"is_dir"`
	Notes    string `json:"notes"`
	Owner    string `json:"owner"`
}

func ListFilesHandler(rootDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, _ := c.Get("username")
		if username == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		path := c.Query("path")

		files, err := os.ReadDir(path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot read directory"})
			log.Printf("Cannot read directory %s: %v", path, err)
			return
		}

		var fileList []FileInfo
		for _, file := range files {
			info, _ := file.Info()
			size := " "
			if !file.IsDir() {
				size = formatSize(info.Size())
			}
			//Owner := database.GetOwner(filepath.Join(path, file.Name()))
			//if Owner == username.(string) {
			fileList = append(fileList, FileInfo{
				Name:     file.Name(),
				Path:     filepath.Join(path, file.Name()),
				Size:     size,
				Modified: info.ModTime().Format(time.DateTime),
				IsDir:    file.IsDir(),
				Notes:    database.GetNotes(filepath.Join(path, file.Name())),
				Owner:    database.GetOwner(filepath.Join(path, file.Name())),
			})
			/*} else {
				log.Printf("User %s is not owner of file %s", username, filepath.Join(path, file.Name()))
			}*/
		}

		c.JSON(http.StatusOK, gin.H{"files": fileList})
	}
}

func formatSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}
