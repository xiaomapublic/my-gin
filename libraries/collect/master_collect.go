package collect

type MasterCollect struct {
	compare  func(interface{}, interface{}) int // 比较函数
	err      error                              // 错误信息
	isCopied bool                               // 是否已经拷贝，如果设置了true，说明已经拷贝，任何操作不影响之前的数组

	ICollect
	Parent ICollect //用于调用子类
}

func (arr *MasterCollect) Err() error {
	return arr.err
}

func (arr *MasterCollect) SetErr(err error) ICollect {
	arr.err = err
	return arr
}

func (arr *MasterCollect) Copy() ICollect {
	if arr.Parent == nil {
		panic("no parent")
	}
	arr.isCopied = true
	return arr.Parent.Copy().SetCompare(arr.compare)
}

func (arr *MasterCollect) SetCompare(compare func(a interface{}, b interface{}) int) ICollect {
	arr.compare = compare
	return arr
}

func (arr *MasterCollect) DD() {
	if arr.Parent == nil {
		panic("DD: not Implement")
	}
	arr.Parent.DD()
}

func (arr *MasterCollect) Append(item interface{}) ICollect {
	if arr.Err() != nil {
		return arr
	}
	return arr.Insert(arr.Count(), item)
}

func (arr *MasterCollect) Insert(index int, obj interface{}) ICollect {
	if arr.Err() != nil {
		return arr
	}
	if arr.Parent == nil {
		panic("no parent")
	}

	if arr.isCopied == false {
		arr.Copy()
		arr.isCopied = true
	}

	return arr.Parent.Insert(index, obj)
}

func (arr *MasterCollect) Count() int {
	if arr.Err() != nil {
		return 0
	}
	if arr.Parent == nil {
		panic("no parent")
	}
	return arr.Parent.Count()
}
