package dparser

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	tmpFilePath string = "/tmp/reqbot_buffer.go"
)

var mustPKGMap = map[string]bool{
	"json":    true,
	"fmt":     true,
	"reflect": true,
	"flag":    true,
}

func RunGet(src string, n int) string {
	tmpFile, err := os.Create(tmpFilePath)
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpFile.Name())

	w := bufio.NewWriter(tmpFile)
	w.WriteString("package main\n\n")
	for _, pkg := range chkPackages(src) {
		w.WriteString(fmt.Sprintf("import \"%s\"\n", pkg))
	}
	w.WriteString(fmt.Sprintf(`
	var n *int = flag.Int("n", 0, "request count")
	
	func init() {
		flag.Parse()
	}
    func main() {
        res := Get(*n)
        switch reflect.TypeOf(res).Kind() {
        case reflect.Slice, reflect.Struct, reflect.Map, reflect.Array:
            data, _ := json.Marshal(res)
            fmt.Print(string(data))
        default:
            fmt.Printf("%%v", res)
        }
    }
    %s`, src))
	w.Flush()

	cmd := exec.Command("go", "run", tmpFile.Name(), "-n", strconv.Itoa(n))
	res, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	return string(res)
}

func chkPackages(src string) (pkgList []string) {
	var pkgMap = make(map[string]bool)
	for k, v := range mustPKGMap {
		pkgMap[k] = v
	}
	for _, line := range strings.Split(src, "\n") {
		for _, item := range strings.SplitN(strings.Trim(line, " "), " ", -1) {
			index := strings.Index(item, ".")
			if index == -1 {
				continue
			}
			pkgMap[item[:index]] = true
		}
	}

	for k := range pkgMap {
		switch k {
		case "json":
			pkgList = append(pkgList, "encoding/json")
		case "rand":
			pkgList = append(pkgList, "math/rand")
		case "time":
			pkgList = append(pkgList, "time")
		case "fmt":
			pkgList = append(pkgList, "fmt")
		case "reflect":
			pkgList = append(pkgList, "reflect")
		case "flag":
			pkgList = append(pkgList, "flag")
		}
	}
	return
}
