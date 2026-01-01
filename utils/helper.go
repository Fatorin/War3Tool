package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateRandomPassword 生成隨機密碼
func GenerateRandomPassword(length int) string {
	rand.Seed(time.Now().UnixNano())

	password := ""
	for i := 0; i < length; i++ {
		password += fmt.Sprintf("%d", rand.Intn(10))
	}

	return password
}
