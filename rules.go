package Ant

import (
	"errors"
	"fmt"
	"regexp"
)

var rules = rulesMap{}

type rulesMap map[string]string

// 获取规则对象
func (r *rulesMap) getRule(name string) (*regexp.Regexp, error) {
	if rules[name] != "" {
		regexpObj, err := regexp.Compile(rules[name])
		if err != nil {
			return nil, errors.New(fmt.Sprint("规则编译错误,请联系开发者. err:", err))
		}
		return regexpObj, nil
	}
	return nil, errors.New("未知的规则")
}
