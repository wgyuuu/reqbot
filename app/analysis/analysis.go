package analysis

import (
	"strings"

	"github.com/wgyuuu/reqbot/common/dlog"
)

const (
	IndexKey int = iota
	IndexValue
)

type AnalysisCNF struct {
	currKey   string
	mapConfig map[string]*ParamsData
}

func New() *AnalysisCNF {
	return &AnalysisCNF{
		mapConfig: make(map[string]*ParamsData),
	}
}

func (a *AnalysisCNF) Accept(data string) {
	data = PreTreat(data)
	if len(data) > 2 && data[:1] == "[" && data[len(data)-1:] == "]" {
		a.putData(strings.ToUpper(data[1:len(data)-1]), newParamsData())
		return
	}

	params := a.getCurrParams()
	if params == nil {
		return
	}

	params.Put(data)
}

func (a *AnalysisCNF) GetDefaultValue(keyType string, reqN int) string {
	params := a.getParams(keyType)
	if params == nil {
		return ""
	}
	return params.GetDefaultValue(reqN)
}

func (a *AnalysisCNF) GetValues(keyType string, reqN int) map[string]string {
	params := a.getParams(keyType)
	if params == nil {
		return make(map[string]string)
	}
	return params.GetValues(reqN)
}

func (a *AnalysisCNF) putData(key string, params *ParamsData) {
	ok := KeyMap[key]
	if !ok {
		dlog.Errorf("key(%s) error.\n", key)
		return
	}
	a.currKey = key
	a.mapConfig[key] = params
}

func (a *AnalysisCNF) getParams(key string) *ParamsData {
	ok := KeyMap[key]
	if !ok {
		return nil
	}

	params := a.mapConfig[key]
	return params
}

func (a *AnalysisCNF) getCurrParams() (params *ParamsData) {
	if a.currKey == "" {
		return nil
	}
	params = a.mapConfig[a.currKey]
	return
}

// 去空格& //注释
func PreTreat(data string) string {
	lineList := strings.Split(data, "\n")

	for i, line := range lineList {
		list := strings.Split(line, " ")

		for k, str := range list {
			idxHttp := strings.Index(str, "http")
			idxIgnore := strings.Index(str, "//")

			if idxHttp == -1 && idxIgnore > -1 {
				list[k] = str[:idxIgnore]
				list = list[:k]
				break
			}
		}
		line = strings.Join(list, " ")
		lineList[i] = strings.TrimSpace(line)
	}
	return strings.Join(lineList, "\n")
}
