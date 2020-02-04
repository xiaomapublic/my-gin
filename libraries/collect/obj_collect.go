package collect

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type ObjCollect struct {
	MasterCollect
	objs reflect.Value
	typ  reflect.Type
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

func (arr *ObjCollect) DD() {
	ret := fmt.Sprintf("ObjCollect(%d)(%s):{\n", arr.Count(), arr.typ.String())

	iter := arr.objs.MapRange()

	for iter.Next() {
		ret = ret + fmt.Sprintf("\t%s:\t%+v\n", iter.Key(), iter.Value())
	}
	ret = ret + "}\n"
	fmt.Print(ret)
}

func (arr *ObjCollect) GetInterface() interface{} {
	return arr.objs.Interface()
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

// 分组
func (arr *ObjCollect) GroupBy(keys ...string) ICollect {

	if arr.objs.Kind().String() != "slice" {
		arr.err = errors.New("not slice")
	}

	if arr.Err() != nil {
		panic(arr.err)
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

// 求和
func (arr *ObjCollect) Sum(k string) (sum int64) {

	for i := 0; i < arr.objs.Len(); i++ {
		sum += arr.objs.Index(i).FieldByName(k).Int()
	}
	return
}
