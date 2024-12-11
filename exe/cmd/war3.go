package cmd

import (
	"bufio"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/samber/lo"
	"io"
	"io/fs"
	"lik/lib"
	"os"
	"regexp"
	"strings"
)

// 处理 war3map.j 文件
func war3mapJ(name string) {
	if lib.Data.BuildModeName != "_release" {
		war3mapJass := lib.Data.BuildDstPath + "/map/war3map.j"
		if fileutil.IsExist(war3mapJass) {
			var war3mapContentBuilder strings.Builder
			jassFile, err := os.OpenFile(war3mapJass, os.O_RDONLY, 0666)
			defer jassFile.Close()
			if err != nil {
				lib.Panic(err)
			}
			srcReader := bufio.NewReader(jassFile)
			for {
				str, err2 := srcReader.ReadString('\n')
				if err2 != nil {
					if err2 == io.EOF {
						break
					} else {
						lib.Panic(err2)
					}
				}
				if strings.Trim(str, " ")[0:2] == "//" {
					continue
				}
				if strings.HasSuffix(str, " \r\n") {
					continue
				}
				if strings.HasPrefix(str, "\r\n") {
					continue
				}
				war3mapContentBuilder.WriteString(str)
			}
			war3mapContent := war3mapContentBuilder.String()
			// 处理J
			execGlobals := []string{"unit prevReadToken = null"}
			var execLua []string
			var execCheatSha1 []string

			nLen := 996
			lua := "exec-lua:"
			luas := strings.Split(lua, "")
			for _, v := range luas {
				nn := lib.NanoOL(nLen)
				execGlobals = append(execGlobals, "string "+nn+"=\""+v+"\"")
				execLua = append(execLua, nn)
			}
			sha1s := strings.Split(name, "")
			for _, v := range sha1s {
				nn := lib.NanoOL(nLen)
				execGlobals = append(execGlobals, "string "+nn+"=\""+v+"\"")
				execCheatSha1 = append(execCheatSha1, nn)
			}
			gi := 10
			for {
				if gi <= 0 {
					break
				}
				nn := lib.NanoOL(nLen)
				execGlobals = append(execGlobals, "string "+nn+"=\""+lib.CharRand()+"\"")
				gi -= 1
			}

			execGlobals = lo.Shuffle[string](execGlobals)

			reContent := strings.Join(execGlobals, "\r\n") + "\r\nendglobals"
			reg, _ := regexp.Compile("endglobals")
			war3mapContent = reg.ReplaceAllString(war3mapContent, reContent)

			reLua := "call Cheat(" + strings.Join(execLua, "+") + "+\"\\\"\"+" + strings.Join(execCheatSha1, "+") + "+\"\\\"\")"
			reContent = "function InitGlobals takes nothing returns nothing\r\n    " + "set prevReadToken = CreateUnit(Player(15),'hfoo',0,0,0)\r\n    " + reLua
			reg, _ = regexp.Compile("function InitGlobals takes nothing returns nothing")
			war3mapContent = reg.ReplaceAllString(war3mapContent, reContent)

			// merge
			err = lib.FilePutContents(war3mapJass, war3mapContent, fs.ModePerm)
			if err != nil {
				lib.Panic(err)
			}
		}
	}
}

func War3() {
	name := Lua()
	war3mapJ(name)
}
