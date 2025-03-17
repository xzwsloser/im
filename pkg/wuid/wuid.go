package wuid

import (
	"fmt"
	"sort"
	"strconv"
)

// @Description: 根据第一个 Id 和 第二个 Id 组合成第三个 Id
func CombineId(aid, bid string) string {
	ids := []string{aid, bid}

	sort.Slice(ids, func(i, j int) bool {
		a, _ := strconv.ParseUint(ids[i], 0, 64)
		b, _ := strconv.ParseUint(ids[j], 0, 64)
		return a < b
	})

	return fmt.Sprintf("%s_%s", ids[0], ids[1])
}
