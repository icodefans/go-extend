package function

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

// 生成随机标识
func Randid(prefix uint8, len uint8) (uint64, error) {
	var (
		number string
	)
	if prefix == 0 {
		number = fmt.Sprintf(fmt.Sprint("%0", len, "d"),
			rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(int64(math.Pow(10, float64(len)))),
		)
	} else {
		number = fmt.Sprintf(fmt.Sprint(`%d%0`, len, "d"), prefix,
			rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(int64(math.Pow(10, float64(len)))),
		)
	}
	return strconv.ParseUint(number, 10, 64)
}
