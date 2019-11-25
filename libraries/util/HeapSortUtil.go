/**
 * 堆排序
 */
package util

/**
 *	替换元素
 *
 */
func replace(data []int, i int, j int) {
	if i != j {
		data[i], data[j] = data[j], data[i]
	}
}

/**
 *	创建小根堆
 *
 */
func createMinHeap(data []int, lastIndex int) {
	// 最后一个节点的父节点下标
	fatherIndex := (lastIndex - 1) / 2
	for i := fatherIndex; 0 <= i; i-- {
		// 记录当前节点
		k := i
		/**
		 * 将整个树变为最小根堆，包括每个子树的根节点为最小
		 * 判断子节点是否存在，如果存在则创建最小根堆
		 */
		for 2*k+1 <= lastIndex {
			// 假设最小值为左节点
			smallIndex := 2*k + 1
			// 判断右节点是否存在，如果右节点小于等于最后一个节点，则存在
			if smallIndex+1 <= lastIndex {
				// 判断左右节点，哪个值更小
				if data[smallIndex] > data[smallIndex+1] {
					// 如果左节点大于右节点，则最小值为右节点
					smallIndex++
				}
			}
			// 如果父节点大于最小子节点，则交换位置
			if data[k] > data[smallIndex] {
				replace(data, k, smallIndex)
				// 交换位置后，再判断换位置后的节点作为根节点的子堆是否需要调整
				k = smallIndex
			} else {
				break
			}
		}
	}
}

/**
 * 获取指定数量的最大值数组
 *
 */
func GetMaxNumber(count int, data []int) []int {
	if len(data) < count {
		count = len(data)
	}
	// 创建存储排序数据的切片
	maxNumberArr := make([]int, count)
	for i := 0; i < len(data); i++ {
		// 创建指定大小的最小根堆
		if data[i] > maxNumberArr[0] {
			maxNumberArr[0] = data[i]
			createMinHeap(maxNumberArr, count-1)
		}
	}
	return maxNumberArr
}

/**
 * 小根堆，倒序排序
 *
 */
func SmallHeapDesc(data []int) []int {
	// 将浅拷贝转换为深拷贝
	result := make([]int, len(data), len(data))
	copy(result, data)

	// 排序：从下标0开始，切片结尾指针逐步前移进行小根堆排序，并将最小值移动到数组末尾，然后用剩下的数组继续组成最小根堆
	for i := 0; i < len(result)-1; i++ {
		createMinHeap(result, len(result)-1-i)
		replace(result, 0, len(result)-1-i)
	}
	return result
}

/**
 *	小根堆，正序排序
 *
 */
func SmallHeapAsc(data []int) (result []int) {
	// 将浅拷贝转换为深拷贝
	result = make([]int, len(data), len(data))
	copy(result, data)
	// 排序：从下标0开始，切片指针逐步后移进行小根堆排序，并将最小值移动到数组开头，然后用剩下的数组继续组成最小根堆
	for i := 0; i < len(result)-1; i++ {
		// 切片截取，并未开辟新内存地址，只是指针的移动
		createMinHeap(result[i:], len(result[i:])-1)
	}
	return
}

/**
 *	创建大根堆
 *
 */
func createMaxHeap(data []int, lastIndex int) {
	// 获取最大下标的父节点
	for i := (lastIndex - 1) / 2; i >= 0; i-- {
		k := i
		// 如果父节点的左子节点存在
		for 2*k+1 <= lastIndex {
			// 假设最大值为左子节点
			largeIndex := 2*k + 1
			// 如果父节点的右子节点存在，则将左右节点作比较
			if largeIndex+1 <= lastIndex {
				// 如果右子节点值大于左子节点的值，则将最大值下标标记为右子节点
				if data[largeIndex] < data[largeIndex+1] {
					largeIndex++
				}
			}

			// 将父节点值与最大值子节点值作比较
			if data[k] < data[largeIndex] {
				// 如果子节点值大于父节点值，则将数据位置交换
				replace(data, k, largeIndex)
				// 交换数据后，将父节点下标标记为交换数据后的子节点下标，以子节点下标作为父节点，来构建最大根堆
				k = largeIndex
			} else {
				break
			}
		}
	}
}

/**
 *	获取指定数量的最小值数组
 *
 */
func GetMinNumber(count int, data []int) []int {
	// 如果待排序数据总数小于需要排序的数据总数，则按照待排序数据总数来建堆
	if len(data) < count {
		count = len(data)
	}

	// 将切片浅拷贝转换为深拷贝
	minNumberArr := make([]int, count, count)
	copy(minNumberArr, data[:count])

	// 创建原始最大根堆
	createMaxHeap(minNumberArr, count-1)
	// 将原数据与最大根堆的根节点作比较，如果根节点数据大于当前数据，则替换根节点数据，然后重新创建最大根堆
	for i := count; i < len(data); i++ {
		if minNumberArr[0] > data[i] {
			minNumberArr[0] = data[i]
			createMaxHeap(minNumberArr, count-1)
		}
	}
	return minNumberArr
}

/**
 *	大根堆，正序排序
 *
 */
func LargeHeapAsc(data []int) []int {
	// 将切片浅拷贝转换为深拷贝
	result := make([]int, len(data), len(data))
	copy(result, data)
	// 排序：从下标0开始，切片结尾指针逐步前移进行最大根堆排序，并将最大值移到数组末尾，然后用剩下的数组继续组成最大根堆
	for i := 0; i < len(result)-1; i++ {
		createMaxHeap(result, len(result)-1-i)
		replace(result, 0, len(result)-1-i)
	}
	return result
}

/**
 *	大根堆，倒序排序
 */
func LargeHeapDesc(data []int) (result []int) {
	// 将切片浅拷贝转换为深拷贝
	result = make([]int, len(data), len(data))
	copy(result, data)

	// 排序：从下标0开始，逐步后移进行大根堆排序，并将最大值移动到数组开头
	for i := 0; i < len(result)-1; i++ {
		createMaxHeap(result[i:], len(result[i:])-1)
	}
	return
}
