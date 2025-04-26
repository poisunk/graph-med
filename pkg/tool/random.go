package tool

import (
	"math/rand"
	"time"
)

const (
	RAND_STRING_KIND_NUM   = 0 // 纯数字
	RAND_STRING_KIND_LOWER = 1 // 小写字母
	RAND_STRING_KIND_UPPER = 2 // 大写字母
	RAND_STRING_KIND_ALL   = 3 // 数字、大小写字母
)

// RandomString 随机字符串
func RandomString(size int, kind int) string {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		if is_all {
			ikind = randGen.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + randGen.Intn(scope))
	}
	return string(result)
}
