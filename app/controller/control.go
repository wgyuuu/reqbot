package controller

import (
	"bufio"
	"os"
	"strconv"

	"github.com/wgyuuu/reqbot/app/analysis"
	"github.com/wgyuuu/reqbot/common/dlog"
	"github.com/wgyuuu/reqbot/common/util"

	"io"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func ProcCNF(filePath string) error {
	f, err := os.OpenFile(filePath, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(f)

	var buffer string
	var any *analysis.AnalysisCNF
	for {
		line, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		buffer += string(line)
		if isPrefix {
			continue
		}
		if strings.Count(buffer, "`")%2 == 1 {
			buffer += "\n"
			continue
		}

		any = procData(any, buffer)
		buffer = ""
	}

	wg.Wait()
	return nil
}

func procData(any *analysis.AnalysisCNF, data string) *analysis.AnalysisCNF {
	if any == nil {
		if strings.TrimSpace(data) == "{" {
			any = analysis.New()
		}
		return any
	}

	if strings.TrimSpace(data) == "}" {
		wg.Add(1)
		go util.SafeFunc(func() {
			defer wg.Done()
			startReq(any)
		})

		return nil
	}

	any.Accept(data)
	return any
}

func startReq(any *analysis.AnalysisCNF) {
	start := time.Now()
	reqN, _ := strconv.Atoi(any.GetDefaultValue(analysis.Count, 0))

	defer func() {
		timer := time.Since(start)
		fields := map[string]interface{}{
			"title": "total_req",
			"timer": timer,
		}
		dlog.NewEntry(fields).Info2("avg_timer", timer/time.Duration(reqN), "req_count", reqN)
	}()

	var reqGW sync.WaitGroup
	for n := 0; n < reqN; n++ {
		url := any.GetDefaultValue(analysis.Url, n)
		method := any.GetDefaultValue(analysis.Method, n)
		mapHeaders := any.GetValues(analysis.Header, n)
		mapParams := any.GetValues(analysis.Params, n)
		body := any.GetDefaultValue(analysis.Body, n)

		reqGW.Add(1)
		go util.SafeFunc(func() {
			defer reqGW.Done()
			SendRequest(url, method, mapHeaders, mapParams, body)
		})
	}

	reqGW.Wait()
}
