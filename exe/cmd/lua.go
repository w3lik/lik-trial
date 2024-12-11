package cmd

import (
	"embed"
	"github.com/duke-git/lancet/v2/fileutil"
	lua "github.com/yuin/gopher-lua"
	"lik/lib"
	"regexp"
	"strings"
)

func require(L *lua.LState, file string) {
	err := L.DoFile(file)
	if err != nil {
		lib.Panic(err)
	}
}

func requireEmbed(L *lua.LState, e embed.FS, file string) {
	s, _ := e.ReadFile(file)
	err := L.DoString(string(s))
	if err != nil {
		lib.Panic(err)
	}
}

func asyncRandReg(asyncRandIn [][]string, lc string) [][]string {
	reg, _ := regexp.Compile(`async.rand\((.+?),(.+?)\)`)
	ar := reg.FindAllStringSubmatch(lib.ZipLua(lc), -1)
	if len(ar) > 0 {
		for _, a := range ar {
			a[1] = strings.Trim(a[1], " ")
			a[2] = strings.Trim(a[2], " ")
			if a[1] != a[2] {
				asyncRandIn = append(asyncRandIn, []string{a[1], a[2]})
			}
		}
	}
	return asyncRandIn
}

func ini(file string) string {
	content, errIni := fileutil.ReadFileToString(file)
	if errIni != nil {
		return ""
	}
	return content
}

func Lua() string {
	coreName := ""
	if lib.Data.BuildModeName == "_release" {
		coreName = luaRelease()
	} else if lib.Data.BuildModeName == "_dist" {
		coreName = luaDist()
	} else if lib.Data.BuildModeName == "_build" {
		coreName = luaBuild()
	} else {
		coreName = luaTest()
	}
	return coreName
}
