package analysis

import (
	"fmt"
	"strings"

	"github.com/wgyuuu/reqbot/common/dparser"
)

type ParamsData struct {
	params []Value
}

type Value struct {
	ing   bool // 处理没结束
	vType Type
	info  []string // 1:key 0:value
}

func newParamsData() *ParamsData {
	return new(ParamsData)
}

func (p *ParamsData) Put(data string) {
	var value Value
	if l := len(p.params); l > 0 && p.params[l-1].ing {
		value = p.params[l-1]
	}

	if len(value.info) == 0 {
		list := strings.Split(data, SplitSymbol)

		value.info = make([]string, len(list))
		if len(list) == 2 {
			value.info[1] = strings.TrimSpace(list[0])
		}
		p.params = append(p.params, value)

		strValue := strings.TrimSpace(list[len(list)-1])
		if len(strValue) == 0 {
			return
		}

		if strValue[0] == '`' {
			if length := len(strValue); length > 1 && strValue[length-1] == '`' {
				value.info[0] = strValue[1 : len(strValue)-1]
			} else {
				value.ing = true
				value.info[0] = strValue[1:]
			}
		} else if length := len(strValue); length >= 2 && strValue[0] == '"' && strValue[length-1] == '"' {
			value.info[0] = strValue[1 : length-1]
		} else {
			value.info[0] = strValue
		}
	} else {
		if len(data) > 0 && data[len(data)-1] == '`' {
			value.ing = false
			data = data[:len(data)-1]
		}
		value.info[0] = fmt.Sprintf("%s\n%s", value.info[0], data)
	}

	if !value.ing && strings.Index(value.info[0], "func Get(") != -1 {
		value.vType = Func
	}
	p.params[len(p.params)-1] = value
}

func (p *ParamsData) GetDefaultValue(reqN int) string {
	if len(p.params) != 1 && len(p.params[0].info) != 1 {
		return ""
	}
	if value := p.params[0]; value.vType == Func {
		return dparser.RunGet(value.info[0], reqN)
	}
	return p.params[0].info[0]
}

func (p *ParamsData) GetValues(reqN int) map[string]string {
	mapValues := make(map[string]string)
	for _, value := range p.params {
		if len(value.info) == 1 {
			value.info = append(value.info, "")
		}

		switch value.vType {
		case Func:
			mapValues[value.info[1]] = dparser.RunGet(value.info[0], reqN)
		default:
			mapValues[value.info[1]] = value.info[0]
		}
	}
	return mapValues
}
