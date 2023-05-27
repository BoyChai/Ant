package Ant

import (
	"fmt"
	"reflect"
	"strings"
)

type Validator struct {
}

type Error struct {
	Is   bool
	Data []errorData
}
type errorData struct {
	Value string
}

var a int

var b int

// Struct 结构体检查
func (v *Validator) Struct(s interface{}, filePath ...string) Error {
	var err []string
	tag := reflect.TypeOf(s)
	value := reflect.ValueOf(s)
	for i := 0; i < tag.NumField(); i++ {
		ruleName := tag.Field(i).Tag.Get("ant")
		fieldType := tag.Field(i).Type.Kind()
		fieldName := tag.Field(i).Name

		// 嵌套结构体处理
		if fieldType == reflect.Struct {
			e := v.Struct(value.FieldByName(fieldName).Interface(), formatFieldName(filePath, fieldName))
			for i := 0; i < len(e.Data); i++ {
				err = append(err, e.Data[i].Value)
			}
			continue
		}

		// 切片类型处理
		if fieldType == reflect.Slice {
			a++
			fmt.Println("Test:", a)
			sliceValue := value.FieldByName(fieldName)
			// 遍历结构体切片
			for j := 0; j < sliceValue.Len(); j++ {
				b++
				fmt.Println(b)
				sliceFieldName := fieldName + "[" + fmt.Sprint(j) + "]"
				// Test[0]
				elemValue := sliceValue.Index(j)
				if elemValue.Kind() == reflect.Struct {
					e := v.Struct(elemValue.Interface(), formatFieldName(filePath, sliceFieldName))
					for i := 0; i < len(e.Data); i++ {
						err = append(err, e.Data[i].Value)
					}
				}
				// 检查是否需要数据检查
				err = append(err, checkValue(ruleName, elemValue, filePath, fieldName, true, sliceFieldName)...)
			}
			continue
		}

		// 检查是否需要数据检查
		err = append(err, checkValue(ruleName, value, filePath, fieldName, false)...)
	}
	return formatError(err)
}

// 格式化字段名称
func formatFieldName(parentPath []string, fieldName string) string {
	path := strings.Join(append(parentPath, fieldName), ".")
	return path
}

// 格式化错误输出
func formatError(data []string) Error {
	if len(data) <= 0 {
		return Error{}
	}
	err := make([]errorData, len(data))
	for i, value := range data {
		err[i] = errorData{Value: value}
	}
	e := Error{Is: true, Data: err}
	return e
}

func checkValue(ruleName string, value reflect.Value, filePath []string, fieldName string, slice bool, sliceFieldName ...string) []string {
	var err []string
	if ruleName != "" {
		// 拿到正则检查对象
		rule, e := rules.getRule(ruleName)
		if e != nil {
			if len(sliceFieldName) == 0 {
				err = append(err, fmt.Sprintf(formatFieldName(filePath, fieldName)+": "+ruleName+" "+e.Error()))
			} else {
				err = append(err, fmt.Sprintf(formatFieldName(filePath, sliceFieldName[0])+": "+ruleName+" "+e.Error()))
			}
			return err
		}
		// 使用正则对象进行检查
		var isMatch bool
		if !slice {
			isMatch = rule.MatchString(fmt.Sprint(value.FieldByName(fieldName).Interface()))
		} else {
			isMatch = rule.MatchString(fmt.Sprint(value.Interface()))
		}
		if !isMatch {
			if len(sliceFieldName) == 0 {
				err = append(err, fmt.Sprintf(formatFieldName(filePath, fieldName)+": Check failed"))
			} else {
				err = append(err, fmt.Sprintf(formatFieldName(filePath, sliceFieldName[0])+": Check failed"))
			}
			return err
		}
	}
	return err
}
