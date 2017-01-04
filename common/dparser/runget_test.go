package dparser

import (
	"testing"

	"strings"

	check "gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type StringSuite struct {
	Content string
}

func init() {
	check.Suite(&StringSuite{`
        func Get(n int) (list []int) {
            for token := 1; token < 10; token++ {
                list = append(list, rand.Intn(token))
            }
            return
        } 
        `})
}

func (s *StringSuite) TestSplit(c *check.C) {
	lineList := strings.Split(s.Content, "\n")
	c.Assert(lineList, check.HasLen, 8)
}

type GetSuite struct {
	intList string
	strList string
}

func init() {
	check.Suite(&GetSuite{
		intList: `
        func Get(n int) (list []int) {
            for token := 1; token < 10; token++ {
                list = append(list, rand.Intn(token))
            }
			list = append(list, n)
            return
        }
        `,
		strList: `
                func Get(n int) []string {
            return []string{"AA A", "BBB hah", "CCC"}
        }
        `,
	})
}

func (g *GetSuite) TestGetIntList(c *check.C) {
	res := RunGet(g.intList, 1000)
	c.Log(res)
	c.Assert(len(res) > 0, check.Equals, true)
}

func (g *GetSuite) TestGetStrList(c *check.C) {
	res := RunGet(g.strList, 0)
	c.Log(res)
	c.Assert(len(res) > 0, check.Equals, true)
}
