package lib

import (
	"encoding/json"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/samber/lo"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	encryptForce []string
	encryptRule  []map[string]string
)

func init() {
	encryptForce = []string{
		`FRAMEWORK_VOICE_INIT`,
		`FRAMEWORK_SPEECH`,
		`FRAMEWORK_SLK_N2I`,
		`FRAMEWORK_SLK_I2V`,
		`FRAMEWORK_SLK_FIX`,
		`FRAMEWORK_SETTING`,
		`FRAMEWORK_N2C`,
		`FRAMEWORK_MAP_NAME`,
		`FRAMEWORK_INIT`,
		`FRAMEWORK_ID`,
		`FRAMEWORK_GO_IDS`,
		`FRAMEWORK_FONT`,
		`FRAMEWORK_DESTRUCTABLE`,
		`FRAMEWORK_ALIAS`,
		`DEBUGGING`,
	}
	encryptRule = []map[string]string{
		{`file`: `builtIn/engine`, `repl`: `(J\.\w+) = `, `del`: `J = {}`},
		{`file`: `builtIn/blizzard`, `repl`: `^(\w+) = `, `del`: ``},
		{`file`: `builtIn/i18n`, `repl`: `^(i18n\._\w+) = |^function (i18n\.\w+)\(`, `del`: `i18n = i18n or {}`},
		{`file`: `builtIn/slk`, `repl`: `^function (\w+)\(`, `del`: ``},
		{`file`: `builtIn/union`, `repl`: `^function (\w+)\(`, `del`: ``},
		{`file`: `builtIn/pairx`, `repl`: `^function (pairx)\(`, `del`: ``},
		{`file`: `oop/base/classExpand`, `repl`: `^(\w+) = |^function (\w+)\(`, `del`: ``},
		{`file`: `oop/base/className`, `force`: `^(\w+) = `, `del`: ``},
		{`file`: `oop/base/object`, `repl`: `^function (\w+)\(`, `del`: ``},
		{`file`: `oop/base/oop`, `repl`: `(oop\._\w+) = `, `force`: `function (\w+)\(`, `del`: `oop = oop or {}`},
		{`file`: `oop/base/class`, `force`: `^(\w+) = |^function (\w+)\(`, `del`: ``},
	}
}

func encryptAnalysis() ([]string, []string, []string, []string, []string) {
	var del []string
	repl := []string{"math.sin", "math.cos", "math.tan", "math.asin", "math.acos", "math.atan"}

	forceList := []string{Data.Encrypt + `/force.json`}
	for _, p := range Data.Plugins {
		forceList = append(forceList, Data.Assets+`/war3mapPlugins/`+p+`/encrypt/force.json`)
	}
	for _, f := range forceList {
		if fileutil.IsExist(f) {
			data, _ := ioutil.ReadFile(f)
			var jsonData []string
			err := json.Unmarshal(data, &jsonData)
			if err == nil {
				encryptForce = append(encryptForce, jsonData...)
				sort.Sort(sort.Reverse(sort.StringSlice(encryptForce)))
			}
		}
	}

	ruleList := []string{Data.Encrypt + `/rule.json`}
	for _, p := range Data.Plugins {
		ruleList = append(ruleList, Data.Assets+`/war3mapPlugins/`+p+`/encrypt/rule.json`)
	}
	for _, f := range ruleList {
		if fileutil.IsExist(f) {
			data, _ := ioutil.ReadFile(f)
			var r []map[string]string
			err := json.Unmarshal(data, &r)
			if err == nil {
				encryptRule = append(encryptRule, r...)
			}
		}
	}

	for _, rule := range encryptRule {
		if rule[`del`] != `` {
			del = append(del, rule[`del`])
		}
		paths := []string{
			Data.Library + `/` + rule[`file`] + `.lua`,
			Data.Projects + `/` + Data.ProjectName + `/sublibrary/` + rule[`file`] + `.lua`,
		}
		for _, p := range Data.Plugins {
			paths = append(paths, Data.Assets+`/war3mapPlugins/`+p+`/plulibrary/`+rule[`file`]+`.lua`)
		}
		for _, path := range paths {
			if fileutil.IsExist(path) {
				if rule[`repl`] != `` {
					content, errl := fileutil.ReadFileToString(path)
					if errl != nil {
						Panic(errl)
					}
					rep := strings.Split(rule[`repl`], `|`)
					for _, re := range rep {
						reg, _ := regexp.Compile(`(?m)` + re)
						match := reg.FindAllStringSubmatch(content, -1)
						for _, m := range match {
							repl = append(repl, m[1])
						}
					}
				}
				if rule[`force`] != `` {
					content, errl := fileutil.ReadFileToString(path)
					if errl != nil {
						Panic(errl)
					}
					rep := strings.Split(rule[`force`], `|`)
					for _, re := range rep {
						reg, _ := regexp.Compile(`(?m)` + re)
						match := reg.FindAllStringSubmatch(content, -1)
						for _, m := range match {
							encryptForce = append(encryptForce, m[1])
						}
					}
				}
			}
		}
	}

	sort.Slice(repl, func(i, j int) bool {
		l1 := len(repl[i])
		l2 := len(repl[j])
		if l1 > l2 {
			return true
		}
		if l1 == l2 {
			return repl[i] > repl[j]
		}
		return false
	})

	// facades单独处理
	var facades []string
	paths := []string{
		Data.Library + `/oop/facades/`,
		Data.Projects + `/` + Data.ProjectName + `/sublibrary/oop/facades/`,
	}
	for _, p := range Data.Plugins {
		paths = append(paths, Data.Assets+`/war3mapPlugins/`+p+`/plulibrary/oop/facades/`)
	}
	for _, path := range paths {
		if fileutil.IsDir(path) {
			err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				pLen := len(path)
				if path[pLen-4:pLen] == ".lua" {
					content, cerr := fileutil.ReadFileToString(path)
					if cerr != nil {
						Panic(cerr)
					}
					rep := strings.Split(`(?m)function (\w+)\(`, `|`)
					for _, re := range rep {
						reg, _ := regexp.Compile(re)
						match := reg.FindAllStringSubmatch(content, -1)
						for _, m := range match {
							facades = append(facades, m[1])
						}
					}
				}
				return nil
			})
			if err != nil {
				Panic(err)
			}
		}
	}
	// meta方法单独处理
	var meta []string
	metaBan := []string{
		// class
		`construct`, `destruct`, `restruct`, `remove`,
		// string
		`byte`, `find`, `format`, `gmatch`, `gsub`, `len`, `lower`, `match`,
		`pack`, `packsize`, `rep`, `reverse`, `sub`, `unpack`, `upper`,
		// io
		`close`, `flush`, `lines`, `read`, `seek`, `setvbuf`, `write`,
	}
	metaBanKV := make(map[string]bool)
	for _, b := range metaBan {
		metaBanKV[b] = true
	}
	paths = []string{
		Data.Library + `/oop/class/`,
		Data.Projects + `/` + Data.ProjectName + `/sublibrary/oop/class/`,
	}
	for _, p := range Data.Plugins {
		paths = append(paths, Data.Assets+`/war3mapPlugins/`+p+`/plulibrary/oop/class/`)
	}
	for _, path := range paths {
		if fileutil.IsDir(path) {
			err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				pLen := len(path)
				if path[pLen-4:pLen] == ".lua" {
					content, cerr := fileutil.ReadFileToString(path)
					if cerr != nil {
						Panic(cerr)
					}
					rep := strings.Split(`(?m)class:(\w+)\(`, `|`)
					for _, re := range rep {
						reg, _ := regexp.Compile(re)
						match := reg.FindAllStringSubmatch(content, -1)
						for _, m := range match {
							if metaBanKV[m[1]] {
								continue
							}
							meta = append(meta, m[1])
						}
					}
				}
				return nil
			})
			if err != nil {
				Panic(err)
			}
		}
	}
	meta = lo.Uniq[string](meta)
	return del, repl, encryptForce, facades, meta
}
