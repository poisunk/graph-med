package utils

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

const JWT_SECRET = "GRAPH-MED-98492"

var uuidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func IsValidUUID(s string) bool {
	return uuidPattern.MatchString(s)
}

// HashPassword 密码加密
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePasswords 比较输入的密码与哈希值是否匹配
func ComparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// GenerateUUID 生成UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// InList 判断元素是否在列表中
func InList[T comparable](element T, list []T) bool {
	for _, item := range list {
		if item == element {
			return true
		}
	}
	return false
}
