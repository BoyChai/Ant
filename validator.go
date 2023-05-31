package Ant

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Validator struct {
	Parity validatorParity
}

type validatorParity string

const Ant validatorParity = "ant"
const Custom validatorParity = "custom"

type Error struct {
	Is   bool
	Data []errorData
}
type errorData struct {
	Value string
}

// String 字符串检查
func (v *Validator) String(value string, rule string) Error {
	var err []string
	e := v.checkValue(rule, value)
	if e != nil {
		err = append(err, "String: "+e.Error())
		return formatError(err)
	}
	return formatError(err)
}

func (v *Validator) Type(t interface{}, value interface{}, rule string) {

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

		// 切片类型处理
		if fieldType == reflect.Slice || fieldType == reflect.Array {
			sliceValue := value.FieldByName(fieldName)
			// 遍历结构体切片
			for j := 0; j < sliceValue.Len(); j++ {
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
				e := v.checkValue(ruleName, fmt.Sprint(elemValue.Interface()))
				if e != nil {
					err = append(err, fmt.Sprintf(formatFieldName(filePath, sliceFieldName)+e.Error()))
				}
				//err = append(err, checkValue(ruleName, elemValue, filePath, fieldName, true, sliceFieldName)...)
			}
			continue
		}
		//// 检查是否需要数据检查
		e := v.checkValue(ruleName, fmt.Sprint(value.FieldByName(fieldName).Interface()))
		if e != nil {
			err = append(err, fmt.Sprintf(formatFieldName(filePath, fieldName)+e.Error()))
		}
		//err = append(err, checkValue(ruleName, value, filePath, fieldName, false)...)

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

// 数据检查
func (v Validator) checkValue(ruleName, value string) error {
	if ruleName != "" {
		switch v.Parity {
		case Custom:
			rule := custom.getRule(ruleName)
			if rule == nil {
				return errors.New(": " + ruleName + " unknown rule")
			}
			err := rule(value)
			if err != nil {
				return errors.New(": " + err.Error())
			}

			return nil
		default:
			rule, err := rules.getRule(ruleName)
			if err != nil {
				return errors.New(": " + ruleName + " " + err.Error())
			}
			isMatch := rule.MatchString(value)
			if !isMatch {
				return errors.New(": check failed")
			}
			return nil
		}
	}
	return nil
}

// 数据检查
//func checkValue(ruleName string, value reflect.Value, filePath []string, fieldName string, slice bool, sliceFieldName ...string) []string {
//	var err []string
//	if ruleName != "" {
//		// 拿到正则检查对象
//		rule, e := rules.getRule(ruleName)
//		if e != nil {
//			if len(sliceFieldName) == 0 {
//				err = append(err, fmt.Sprintf(formatFieldName(filePath, fieldName)+": "+ruleName+" "+e.Error()))
//			} else {
//				err = append(err, fmt.Sprintf(formatFieldName(filePath, sliceFieldName[0])+": "+ruleName+" "+e.Error()))
//			}
//			return err
//		}
//		// 使用正则对象进行检查
//		var isMatch bool
//		if !slice {
//			isMatch = rule.MatchString(fmt.Sprint(value.FieldByName(fieldName).Interface()))
//		} else {
//			isMatch = rule.MatchString(fmt.Sprint(value.Interface()))
//		}
//		if !isMatch {
//			if len(sliceFieldName) == 0 {
//				err = append(err, fmt.Sprintf(formatFieldName(filePath, fieldName)+": Check failed"))
//			} else {
//				err = append(err, fmt.Sprintf(formatFieldName(filePath, sliceFieldName[0])+": Check failed"))
//			}
//			return err
//		}
//	}
//	return err
//}
