package collect

import (
	"fmt"
	"reflect"
	"strings"
)

type ObjCollect struct {
	MasterCollect
	objs reflect.Value
	typ  reflect.Type
}

func (arr *ObjCollect) DD() {
	ret := fmt.Sprintf("ObjCollect(%d)(%s):{\n", arr.Count(), arr.typ.String())

	iter := arr.objs.MapRange()

	for iter.Next() {
		ret = ret + fmt.Sprintf("\t%s:\t%+v\n", iter.Key(), iter.Value())
	}
	ret = ret + "}\n"
	fmt.Print(ret)
}

func NewObjCollect(objs interface{}) *ObjCollect {
	vals := reflect.ValueOf(objs)
	typ := reflect.TypeOf(objs).Elem()

	arr := &ObjCollect{
		objs: vals,
		typ:  typ,
	}
	arr.MasterCollect.Parent = arr
	return arr
}

func (arr *ObjCollect) NewEmpty(err ...error) ICollect {
	objs := reflect.MakeSlice(arr.objs.Type(), 0, 0)
	ret := &ObjCollect{
		objs: objs,
		typ:  arr.typ,
	}
	ret.Parent = ret
	if len(err) != 0 {
		ret.SetErr(err[0])
	}
	return ret
}

func (arr *ObjCollect) Count() int {
	return arr.objs.Len()
}

func (arr *ObjCollect) GroupBy(keys ...string) ICollect {
	if arr.Err() != nil {
		return arr
	}

	mapType := reflect.MapOf(reflect.TypeOf(""), arr.objs.Type())
	mapObj := reflect.MakeMap(mapType)

	for i := 0; i < arr.objs.Len(); i++ {

		var ret strings.Builder
		for _, key := range keys {
			if len(ret.String()) == 0 {
				o := arr.objs.Index(i).FieldByName(key).Interface()
				ret.WriteString(fmt.Sprintf("%v", o))
			} else {
				o := arr.objs.Index(i).FieldByName(key).Interface()
				ret.WriteString("|")
				ret.WriteString(fmt.Sprintf("%v", o))
			}
		}
		kKeyType := reflect.ValueOf(ret.String())

		if mapObj.MapIndex(kKeyType) != reflect.ValueOf(nil) {
			sliceArr := reflect.Append(mapObj.MapIndex(kKeyType), arr.objs.Index(i))
			mapObj.SetMapIndex(kKeyType, sliceArr)
		} else {
			sliceArr := reflect.MakeSlice(arr.objs.Type(), 0, 0)
			sliceArr = reflect.Append(sliceArr, arr.objs.Index(i))
			mapObj.SetMapIndex(kKeyType, sliceArr)
		}
	}
	eleTyp := mapObj.Type().Elem()
	newArr := &ObjCollect{
		objs: mapObj,
		typ:  eleTyp,
	}

	newArr.MasterCollect.Parent = newArr
	return newArr
}

// 需要修改成与groupby类似使反射的方式
func (arr *ObjCollect) Sum(k string) (sum int64) {

	for i := 0; i < arr.objs.Len(); i++ {
		sum += arr.objs.Index(i).FieldByName(k).Int()
	}
	return
}
