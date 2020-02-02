package collect

type ICollect interface {

	// ICollect错误信息，链式调用的时候需要检查下这个error是否存在，每次调用之后都检查一下
	Err() error
	// 设置ICollect的错误信息
	SetErr(error) ICollect
	// 获取数组长度，对所有Collection生效
	Count() int
	// 增加一个元素。
	Insert(index int, item interface{}) ICollect
	// 复制当前数组
	Copy() ICollect
	// 设置比较函数，理论上所有Collection都能设置比较函数，但是强烈不建议基础Collection设置
	SetCompare(func(a interface{}, b interface{}) int) ICollect
	// 获取数组长度，对所有Collection生效
	GroupBy(...string) ICollect
	// 打印出当前map结构
	DD()
	// 获取sum值
	Sum(string) int64
}
