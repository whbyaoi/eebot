package analysis

import (
	"errors"
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

// 获取在范围内的类别索引
func (jc JJLCategoryKey) GetSpanIndexes(left int, right int) (indexes []int, err error) {
	endpoint := map[string]struct{}{}
	valid := []string{}
	for i := range jc {
		arr := strings.Split(jc[i], "-")
		valid = append(valid, arr...)
		endpoint[arr[0]] = struct{}{}
		endpoint[arr[1]] = struct{}{}
	}
	if left >= right {
		return nil, errors.New("范围左值大于右值")
	}
	if _, ok := endpoint[fmt.Sprintf("%d", left)]; !ok {
		return nil, fmt.Errorf("范围左值 %d 不在可选范围 %v 内", left, valid)
	}
	if _, ok := endpoint[fmt.Sprintf("%d", right)]; !ok {
		return nil, fmt.Errorf("范围左值 %d 不在可选范围 %v 内", right, valid)
	}
	for i := range jc{
		arr := strings.Split(jc[i], "-")
		start, end := arr[0], arr[1]
		s, _ := strconv.ParseInt(start, 10, 64)
		e, _ := strconv.ParseInt(end, 10, 64)
		if s >= int64(left) && e <= int64(right){
			indexes = append(indexes, i)
		}
	}
	return
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

var DefaultJJLCategoryKeys = JJLCategoryKey{"0-1500", "1500-1600", "1600-1700", "1700-1800", "1800-1900", "1900-2000", "2000-2100", "2100-2200", "2200-2300", "2300-2400", "2400-3000"}

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
