package util

//取并集：合并数组(重复元素只保留一份)
func MergeUint64(a, b []uint64) []uint64 {
	la := len(a)
	lb := len(b)
	if la == 0 {
		return b
	} else if lb == 0 {
		return a
	}
	num := 0
	for _, s1 := range b {
		isExist := false
		for _, s2 := range a {
			if s1 == s2 {
				isExist = true
				break
			}
		}
		if isExist {
			num++
		}
	}

	if num == la {
		return b
	} else if num == lb {
		return a
	}

	mlen := la + lb - num
	c := make([]uint64, mlen)
	copy(c, a)
	if num == 0 {
		copy(c[la:], b)
		return c
	}
	index := la
	for _, s1 := range b {
		isExist := false
		for _, s2 := range a {
			if s1 == s2 {
				isExist = true
				break
			}
		}
		if !isExist {
			c[index] = s1
			index++
		}
	}

	return c
}

//取并集：合并数组(重复元素只保留一份)
func MergeUint32(a, b []uint32) []uint32 {
	la := len(a)
	lb := len(b)
	if la == 0 {
		return b
	} else if lb == 0 {
		return a
	}
	num := 0
	for _, s1 := range b {
		isExist := false
		for _, s2 := range a {
			if s1 == s2 {
				isExist = true
				break
			}
		}
		if isExist {
			num++
		}
	}

	if num == la {
		return b
	} else if num == lb {
		return a
	}

	mlen := la + lb - num
	c := make([]uint32, mlen)
	copy(c, a)
	if num == 0 {
		copy(c[la:], b)
		return c
	}
	index := la
	for _, s1 := range b {
		isExist := false
		for _, s2 := range a {
			if s1 == s2 {
				isExist = true
				break
			}
		}
		if !isExist {
			c[index] = s1
			index++
		}
	}

	return c
}

//取并集：合并数组(重复元素只保留一份)
func Merge(a, b []string) []string {
	la := len(a)
	lb := len(b)
	if la == 0 {
		return b
	} else if lb == 0 {
		return a
	}
	num := 0
	for _, s1 := range b {
		isExist := false
		for _, s2 := range a {
			if s1 == s2 {
				isExist = true
				break
			}
		}
		if isExist {
			num++
		}
	}

	if num == la {
		return b
	} else if num == lb {
		return a
	}

	mlen := la + lb - num
	c := make([]string, mlen)
	copy(c, a)
	if num == 0 {
		copy(c[la:], b)
		return c
	}
	index := la
	for _, s1 := range b {
		isExist := false
		for _, s2 := range a {
			if s1 == s2 {
				isExist = true
				break
			}
		}
		if !isExist {
			c[index] = s1
			index++
		}
	}

	return c
}

//取交集：两个字符串数组
func GetContainStringArray(a, b []string) []string {
	if len(a) == 0 || len(b) == 0 {
		return []string{}
	}
	res := []string{}
	for _, c := range a {
		isExist := false
		for _, m := range b {
			if c == m {
				isExist = true
				break
			}
		}
		if isExist {
			res = append(res, c)
		}
	}
	return res
}

//取交集：两个uint64数组
func GetContainuUint64Array(a, b []uint64) []uint64 {
	if len(a) == 0 || len(b) == 0 {
		return []uint64{}
	}

	res := []uint64{}
	var c map[uint64]bool = make(map[uint64]bool, len(b))

	for _, k := range b {
		c[k] = true
	}

	for _, v := range a {
		if _, ok := c[v]; ok {
			res = append(res, v)
		}
	}
	return res
}

//v是否在arr中
func ContainUint32(arr []uint32, v uint32) bool {
	if len(arr) == 0 {
		return false
	}
	for _, num := range arr {
		if num == v {
			return true
		}
	}
	return false
}

func ContainUint64(arr []uint64, v uint64) bool {
	if len(arr) == 0 {
		return false
	}
	for _, num := range arr {
		if num == v {
			return true
		}
	}
	return false
}

func Contain(arr []string, v string) bool {
	if len(arr) == 0 {
		return false
	}
	for _, num := range arr {
		if num == v {
			return true
		}
	}
	return false
}
