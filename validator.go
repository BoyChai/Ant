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
		// 检查是否需要数据检查
		if ruleName != "" {
			// 拿到正则检查对象
			rule, e := rules.getRule(ruleName)
			if e != nil {
				err = append(err, fmt.Sprintf(formatFieldName(filePath, fieldName)+": "+ruleName+" "+e.Error()))
				continue
			}
			// 使用正则对象进行检查
			isMatch := rule.MatchString(fmt.Sprint(value.FieldByName(fieldName).Interface()))
			if !isMatch {
				err = append(err, fmt.Sprintf(formatFieldName(filePath, fieldName)+": Check failed"))
				continue
			}
		}
	}
	return formatError(err)
}

// 格式化字段名称
func formatFieldName(parentPath []string, fieldName string) string {
	path := strings.Join(append(parentPath, fieldName), ".")
	return path
}

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
