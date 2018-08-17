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

	mr := new(MultiReq)
	useTime := time.Duration(0)
	defer func() {
		testTimer := time.Since(start)
		fields := map[string]interface{}{
			"title": "totalReq",
			"测试时间":  testTimer,
		}
		dlog.NewEntry(fields).Info2("请求平均时间", useTime/time.Duration(reqN), "请求次数", reqN, "最大并发数", mr.Max())
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
			mr.Update(1)
			defer mr.Update(-1)
			duration := SendRequest(url, method, mapHeaders, mapParams, body)
			useTime += duration
		})
	}
	reqGW.Wait()
}

type MultiReq struct {
	lock     sync.Mutex
	count    int64
	maxCount int64
}

func (mr *MultiReq) Update(n int64) {
	mr.lock.Lock()
	defer mr.lock.Unlock()
	if n > 0 {
		mr.count += n
	} else {
		// 在减的时候一般是最大值
		if mr.count > mr.maxCount {
			mr.maxCount = mr.count
		}
		mr.count += n
	}
}

func (mr *MultiReq) Max() int64 {
	return mr.maxCount
}
