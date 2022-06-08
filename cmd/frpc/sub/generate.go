/*************************************************************
     FileName: GoLand->src->generate.go
         Date: 2021/10/18 18:20
       Author: 苦咖啡
        Email: voilet@qq.com
         blog: http://blog.kukafei520.net
      Version: 0.0.1
      History:
**************************************************************/

package sub

import (
	"crypto/rand"
	"math"
	"math/big"
	"strconv"
)

func Generate(length int) string {
	const base = 36
	size := big.NewInt(base)
	n := make([]byte, length)
	for i, _ := range n {
		c, _ := rand.Int(rand.Reader, size)
		n[i] = strconv.FormatInt(c.Int64(), base)[0]
	}
	return string(n)
}

// RangeRand 生成区间[-m, n]的安全随机数
func RangeRand(min, max int64) int64 {
	if min > max {
		panic("the min is greater than max!")
	}

	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))

		return result.Int64() - i64Min
	} else {
		result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
		return min + result.Int64()
	}
}
