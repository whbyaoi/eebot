package analysis

import (
	"fmt"
	"strconv"
	"strings"
)

var GangUpCategoryKeys = []string{"单排", "双排", "三黑", "四黑"}

type JJLCategoryKey []string

// 获取位置在i的范围，如果值不规范会panic
func (jc JJLCategoryKey) GetRange(i int) [2]float64 {
	arr := strings.Split(jc[i], "-")
	start, end := arr[0], arr[1]
	s, _ := strconv.ParseFloat(start, 64)
	e, _ := strconv.ParseFloat(end, 64)
	return [2]float64{s, e}
}

// 获取jjl所在的位置
func (jc JJLCategoryKey) Index(jjl float64) int {
	for i := range jc {
		rg := jc.GetRange(i)
		if rg[0] <= jjl && jjl < rg[1] {
			return i
		}
	}
	return -1
}

var DefaultJJLCategoryKeys = JJLCategoryKey{"0-1500", "1500-1700", "1700-1800", "1800-1900", "1900-2000", "2000-2100", "2100-2200", "2200-2500"}

func JJLCategoriesByInterval(interval int) (categories JJLCategoryKey) {
	for i := 1000; i <= 2500; i += interval {
		categories = append(categories, fmt.Sprintf("%d-%d", i, i+interval))
	}
	return
}

type AppraiseCategoryKey []string

var DefaultAppraiseCategoryKeys = AppraiseCategoryKey{"sss", "ss", "s", "a", "b", "c", "d"}

func (ac AppraiseCategoryKey) Appraise(level float64) string {
	if level < 25 {
		return "d"
	} else if level >= 25 && level < 50 {
		return "c"
	} else if level >= 50 && level < 75 {
		return "b"
	} else if level >= 75 && level < 92 {
		return "a"
	} else if level >= 92 && level < 97 {
		return "s"
	} else if level >= 97 && level < 99 {
		return "ss"
	} else {
		return "sss"
	}
}
