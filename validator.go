package Ant

import (
	"errors"
	"fmt"
	"reflect"
)

type Validator struct {
}

// Struct 结构体检查
func (v *Validator) Struct(s interface{}) error {
	var err []string
	tag := reflect.TypeOf(s)
	value := reflect.ValueOf(s)
	for i := 0; i < tag.NumField(); i++ {
		ruleName := tag.Field(i).Tag.Get("ant")
		fieldName := tag.Field(i).Name
		// 检查是否需要数据检查
		if ruleName != "" {
			// 拿到正则检查对象
			rule, e := rules.getRule(ruleName)
			if e != nil {
				err = append(err, fmt.Sprintf(fieldName+": "+ruleName+" "+e.Error()))
				continue
			}
			// 使用正则对象进行检查
			isMatch := rule.MatchString(fmt.Sprint(value.FieldByName(fieldName).Interface()))
			if !isMatch {
				err = append(err, fmt.Sprintf(fieldName+": 检查错误"))
				continue
			}
		}
	}
	return errors.New(fmt.Sprintln(err))
}
