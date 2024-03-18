package utils

import (
	"time"
)

// GetCurrentTime retorna a hora atual
func GetCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
