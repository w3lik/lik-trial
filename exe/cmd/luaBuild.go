package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/pterm/pterm"
	lua "github.com/yuin/gopher-lua"
	"io/fs"
	"lik/lib"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func luaBuild() string {
	if lib.Data.BuildModeName != "_build" {
		lib.Panic("luaBuild")
	}
	var err error
	pterm.Debug.Println("请稍候...")
	//
	var scriptIn []string
	var implIn []string
	var builtIn []map[string]string
	var asyncRandIn [][]string
	var settingMap map[string]string
	var assetsCheck []string
	// 字串数字化
	chao, _ := lib.CharChao()
	sLua := "____={"
	for k, v := range chao {
		i := 3
		nPrev := lib.Nano(lib.Rand(20, 36))
		for {
			if i == 0 {
				if k == `\` {
					k = `\\`
				}
				sLua += nPrev + "='" + k + "',"
				break
			}
			if i == 3 {
				sLua += v + "='" + nPrev + "',"
			} else {
				nCur := lib.Nano(lib.Rand(20, 36))
				sLua += nPrev + "='" + nCur + "',"
				nPrev = nCur
			}
			i = i - 1
		}
	}
	sLua += "} FRAMEWORK_STRING()"
	scriptIn = append(scriptIn, sLua)
	// Inner-Engine
	builtInSort1 := []string{"engine", "console/test", "blizzard", "setting", "i18n"}
	for _, v := range builtInSort1 {
		src := lib.Data.Library + "/builtIn/" + v + ".lua"
		name := strings.Replace(v, "/", "_", -1)
		code, errc := fileutil.ReadFileToString(src)
		if errc != nil {
			lib.Panic(errc)
		}
		dst := lib.Data.BuildDstPath + "/map/" + name + ".lua"
		m := map[string]string{
			"name": name,
			"dst":  dst,
			"code": code,
		}
		builtIn = append(builtIn, m)
		if v == "setting" {
			settingMap = m
		}
	}
	// library
	librarySort := []string{
		"variable", "common", "foundation", "japi", "ability",
		"oop/base", "oop/class", "oop/facades", "oop/modifier", "oop/preposition", "oop/normalizer", "oop/superposition", "oop/limiter",
	}
	for _, n := range librarySort {
		lsf, _ := filepath.Abs(lib.Data.Library + `/` + n)
		err = filepath.Walk(lsf, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			pLen := len(path)
			if path[pLen-4:pLen] == ".lua" {
				requireRelation := strings.Replace(path, lsf, "", 1)
				requireRelation = requireRelation[1:]
				requireRelation = strings.Replace(requireRelation, ".lua", "", -1)
				requireRelation = strings.Replace(requireRelation, "\\", ".", -1)
				requireRelation = strings.Replace(requireRelation, "/", ".", -1)
				requireFile := strings.Replace(requireRelation, ".", "/", -1)
				dst := lib.Data.BuildDstPath + "/map/library/" + n + "/" + requireFile + ".lua"
				name := "library." + strings.Replace(n, `/`, `.`, -1) + "." + requireRelation
				code, errc := fileutil.ReadFileToString(path)
				if errc != nil {
					lib.Panic(errc)
				}
				builtIn = append(builtIn, map[string]string{
					"name": name,
					"dst":  dst,
					"code": code,
				})
			}
			return nil
		})
		if err != nil {
			lib.Panic(err)
		}
	}
	builtInSort2 := []string{"slk", "union", "system", "pairx", "init"}
	for _, v := range builtInSort2 {
		src := lib.Data.Library + "/builtIn/" + v + ".lua"
		name := "builtIn." + v
		dst := lib.Data.BuildDstPath + "/map/" + strings.Replace(name, `.`, `/`, -1) + ".lua"
		code, errc := fileutil.ReadFileToString(src)
		if errc != nil {
			lib.Panic(errc)
		}
		builtIn = append(builtIn, map[string]string{
			"name": name,
			"dst":  dst,
			"code": code,
		})
	}

	//------------------------------

	L := lua.NewState()
	defer L.Close()

	// 外置 lua引擎
	require(L, lib.Data.Library+"/builtIn/console/go.lua")
	require(L, lib.Data.Library+"/builtIn/pairx.lua")
	requireEmbed(L, lib.Embeds, "embeds/builtIn/driver.lua")
	requireEmbed(L, lib.Embeds, "embeds/builtIn/json.lua")
	requireEmbed(L, lib.Embeds, "embeds/builtIn/go.lua")
	requireEmbed(L, lib.Embeds, "embeds/builtIn/slk.lua")
	requireEmbed(L, lib.Embeds, "embeds/builtIn/union.lua")

	// 加载 slk
	requireEmbed(L, lib.Embeds, "embeds/builtIn/system.lua")

	// 加载 项目 assets
	project := lib.Data.Projects + "/" + lib.Data.ProjectName
	assetsDir, _ := filepath.Abs(project + "/assets")
	hasAssets := fileutil.IsDir(assetsDir)
	if hasAssets {
		err = filepath.Walk(assetsDir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			pLen := len(path)
			if path[pLen-4:pLen] == ".lua" {
				require(L, path)
				requireRelation := strings.Replace(path, assetsDir, "", 1)
				requireRelation = requireRelation[1:]
				requireRelation = strings.Replace(requireRelation, ".lua", "", -1)
				requireRelation = strings.Replace(requireRelation, "\\", ".", -1)
				requireRelation = strings.Replace(requireRelation, "/", ".", -1)
				requireFile := strings.Replace(requireRelation, ".", "/", -1)
				name := "assets." + requireRelation
				dst := lib.Data.BuildDstPath + "/map/assets/" + requireFile + ".lua"
				code, errc := fileutil.ReadFileToString(path)
				if errc != nil {
					lib.Panic(errc)
				}
				builtIn = append(builtIn, map[string]string{
					"name": name,
					"dst":  dst,
					"code": code,
				})
			}
			return nil
		})
		if err != nil {
			lib.Panic(err)
		}
	}

	// UI
	lib.CopyPath(lib.Data.Vendor+"/lni/assets/UI", lib.Data.BuildDstPath+"/map/UI")

	// SELECTION
	fn := L.GetGlobal("GO_RESULT_SELECTION")
	if err = L.CallByParam(lua.P{Fn: fn, NRet: 1, Protect: true}); err != nil {
		lib.Panic(err)
	}
	selection := L.ToString(-1)
	selectionDir := lib.Data.Assets + "/war3mapSelection/" + selection
	selectionDirIs := fileutil.IsDir(selectionDir)
	if false == selectionDirIs {
		selectionDir = lib.Data.Vendor + "/lni/assets/Selection"
		selection = "Lni"
	}
	lib.CopyPath(selectionDir, lib.Data.BuildDstPath+"/resource/ReplaceableTextures/Selection")
	pterm.Info.Println("【选择圈】引入：" + selection)

	// FONT
	fn = L.GetGlobal("GO_RESULT_FONT")
	if err = L.CallByParam(lua.P{Fn: fn, NRet: 1, Protect: true}); err != nil {
		lib.Panic(err)
	}
	font := L.ToString(-1)
	fontFile := lib.Data.Assets + "/war3mapFont/" + font + ".ttf"
	if false == fileutil.IsExist(fontFile) {
		font = "default"
	}
	if font == "default" {
		fontFile = lib.Data.Vendor + "/lni/assets/fonts.ttf"
	}
	lib.CopyFile(fontFile, lib.Data.BuildDstPath+"/map/fonts.ttf")
	pterm.Info.Println("【字体】引入：" + font)

	// LOADING
	fn = L.GetGlobal("GO_RESULT_LOADING")
	if err = L.CallByParam(lua.P{Fn: fn, NRet: 1, Protect: true}); err != nil {
		lib.Panic(err)
	}
	_ = os.Remove(lib.Data.BuildDstPath + "/resource/Framework/LoadingScreen.mdx")
	loading := L.ToString(-1)
	loadingPath := lib.Data.Assets + "/war3MapLoading/" + loading
	loadingFile := lib.Data.Assets + "/war3MapLoading/" + loading + ".tga"
	loaded := false
	if fileutil.IsDir(loadingPath) {
		lib.CopyFile(lib.Data.Vendor+"/lni/assets/LoadingScreenDir.mdx", lib.Data.BuildDstPath+"/resource/Framework/LoadingScreen.mdx")
		loadingSites := []string{"pic", "bc", "bg"}
		for _, s := range loadingSites {
			loadingFile = loadingPath + "/" + s + ".tga"
			if fileutil.IsExist(loadingFile) {
				lib.CopyFile(loadingFile, lib.Data.BuildDstPath+"/resource/Framework/LoadingScreen"+s+".tga")
			} else {
				pterm.Error.Println("【载入图碎片】引入：" + s + " 不存在")
			}
		}
		loaded = true
	} else if fileutil.IsExist(loadingFile) {
		lib.CopyFile(lib.Data.Vendor+"/lni/assets/LoadingScreenFile.mdx", lib.Data.BuildDstPath+"/resource/Framework/LoadingScreen.mdx")
		lib.CopyFile(loadingFile, lib.Data.BuildDstPath+"/resource/Framework/LoadingScreen.tga")
		loaded = true
	}
	if loaded {
		w3i, _ := fileutil.ReadFileToString(lib.Data.BuildDstPath + "/table/w3i.ini")
		w3i = strings.Replace(w3i, "\r\n", "\n", -1)
		w3i = strings.Replace(w3i, "\r", "\n", -1)
		w3ia := strings.Split(w3i, "\n")
		canReplace := false
		for k, v := range w3ia {
			if strings.Index(v, "[载入图]") != -1 {
				canReplace = true
			}
			if canReplace && strings.Index(v, "路径") != -1 {
				w3ia[k] = `路径 = "Framework\\LoadingScreen.mdx"`
				break
			}
		}
		err = lib.FilePutContents(lib.Data.BuildDstPath+"/table/w3i.ini", strings.Join(w3ia, "\r\n"), fs.ModePerm)
		if err != nil {
			lib.Panic(err)
		}
		pterm.Info.Println("【载入图】引入：" + loading)
	} else {
		pterm.Error.Println("【载入图】：" + loading + " 不存在")
		time.Sleep(time.Second)
	}

	// PREVIEW
	fn = L.GetGlobal("GO_RESULT_PREVIEW")
	if err = L.CallByParam(lua.P{Fn: fn, NRet: 1, Protect: true}); err != nil {
		lib.Panic(err)
	}
	preview := L.ToString(-1)
	previewFile := lib.Data.Assets + "/war3mapPreview/" + preview + ".tga"
	if fileutil.IsExist(previewFile) {
		lib.CopyFile(previewFile, lib.Data.BuildDstPath+"/resource/war3mapPreview.tga")
		pterm.Info.Println("【预览图】引入：" + preview)
	}

	// ICONS
	fn = L.GetGlobal("GO_RESULT_ICONS")
	if err = L.CallByParam(lua.P{Fn: fn, NRet: 1, Protect: true}); err != nil {
		lib.Panic(err)
	}
	icons := L.ToString(-1)
	var iData [][]string
	err = json.Unmarshal([]byte(icons), &iData)
	if err != nil {
		lib.Panic(err)
	}
	count := 0
	for _, i := range iData {
		assetsCheck = append(assetsCheck, i[1])
		iconFile := lib.Data.Assets + "/war3mapIcon/" + i[0] + ".tga"
		dstFile := lib.Data.BuildDstPath + "/resource/war3mapIcon/" + i[1] + ".tga"
		if fileutil.IsExist(dstFile) {
			pterm.Warning.Println("【重复】：" + i[0] + `->` + i[1])
			time.Sleep(time.Second)
			continue
		}
		if fileutil.IsExist(iconFile) {
			lib.CopyFile(iconFile, dstFile)
			count += 1
		} else {
			pterm.Error.Println("【图标】：文件不存在 " + i[0])
			time.Sleep(time.Second)
		}
	}
	pterm.Info.Println("【图标】引入：" + strconv.Itoa(count) + "个")

	// MODEL
	fn = L.GetGlobal("GO_RESULT_MODEL")
	if err = L.CallByParam(lua.P{Fn: fn, NRet: 1, Protect: true}); err != nil {
		lib.Panic(err)
	}
	models := L.ToString(-1)
	var mData [][]string
	err = json.Unmarshal([]byte(models), &mData)
	if err != nil {
		lib.Panic(err)
	}
	count = 0
	count2 := 0
	for _, i := range mData {
		assetsCheck = append(assetsCheck, i[1])
		modelFile := lib.Data.Assets + "/war3mapModel/" + i[0] + ".mdx"
		if fileutil.IsExist(modelFile) {
			dstFile := lib.Data.BuildDstPath + "/resource/war3mapModel/" + i[1] + ".mdx"
			if fileutil.IsExist(dstFile) {
				pterm.Warning.Println("【重复】：" + i[0] + `->` + i[1])
				time.Sleep(time.Second)
				continue
			}
			lib.CopyFile(modelFile, dstFile)
			count = count + 1
			modelPortraitFile := lib.Data.Assets + "/war3mapModel/" + i[0] + "_Portrait.mdx"
			if fileutil.IsExist(modelPortraitFile) {
				lib.CopyFile(modelPortraitFile, lib.Data.BuildDstPath+"/resource/war3mapModel/"+i[1]+"_Portrait.mdx")
			}
			modelStr, _ := fileutil.ReadFileToString(modelFile)
			modelStr = strings.Replace(modelStr, "\r", "", -1)
			modelStr = strings.Replace(modelStr, "\r", "", -1)
			reg, _ := regexp.Compile("(?i)war3mapTextures(.*?)(.blp)")
			textures := reg.FindAllString(modelStr, -1)
			if len(textures) > 0 {
				ti := 0
				for _, t := range textures {
					t = strings.Replace(t, "war3mapTextures", "", -1)
					textureFile := lib.Data.Assets + "/war3mapTextures" + t
					if fileutil.IsExist(textureFile) {
						lib.CopyFile(textureFile, lib.Data.BuildDstPath+"/resource/war3mapTextures"+t)
						ti += 1
					} else {
						pterm.Error.Println("【贴图】：文件不存在 " + t)
						time.Sleep(time.Second)
					}
				}
				count2 = count2 + ti
			}
		} else {
			pterm.Error.Println("【模型】：文件不存在 " + i[0])
			time.Sleep(time.Second)
		}
	}
	pterm.Info.Println("【模型】引入：" + strconv.Itoa(count) + "个，附带贴图" + strconv.Itoa(count2) + "张")

	// SOUND
	fn = L.GetGlobal("GO_RESULT_SOUND")
	if err = L.CallByParam(lua.P{Fn: fn, NRet: 1, Protect: true}); err != nil {
		lib.Panic(err)
	}
	sounds := L.ToString(-1)
	var sData [][]string
	err = json.Unmarshal([]byte(sounds), &sData)
	if err != nil {
		lib.Panic(err)
	}
	voiceData := make(map[string]map[string]string)
	voiceData["vcm"] = make(map[string]string)
	voiceData["v3d"] = make(map[string]string)
	voiceData["vwp"] = make(map[string]string)
	count = 0
	count2 = 0
	count3 := 0
	count4 := 0
	for _, i := range sData {
		assetsCheck = append(assetsCheck, i[1])
		if i[2] == "vcm" || i[2] == "v3d" || i[2] == "bgm" {
			soundFile := lib.Data.Assets + "/war3mapSound/" + i[0] + ".mp3"
			if fileutil.IsExist(soundFile) {
				lib.CopyFile(soundFile, lib.Data.BuildDstPath+"/resource/war3mapSound/"+i[2]+"/"+i[1]+".mp3")
				if i[2] == "vcm" {
					voiceData["vcm"][i[1]] = strconv.Itoa(lib.VoiceDuration(soundFile))
					count = count + 1
				} else if i[2] == "v3d" {
					voiceData["v3d"][i[1]] = strconv.Itoa(lib.VoiceDuration(soundFile))
					count2 = count2 + 1
				} else if i[2] == "bgm" {
					count3 = count3 + 1
				} else {
					pterm.Error.Println("【声乐】：无效类型 " + i[2])
					time.Sleep(time.Second)
				}
			} else {
				pterm.Error.Println("【声乐】：文件不存在 " + i[0])
				time.Sleep(time.Second)
			}
		} else if i[2] == "vwp" {
			wPath, _ := filepath.Abs(lib.Data.Assets + "/war3mapSound/weapon/" + i[0])
			isWPathDir := fileutil.IsDir(wPath)
			if isWPathDir {
				err = filepath.Walk(wPath, func(path string, info fs.FileInfo, err error) error {
					if err != nil {
						return err
					}
					pLen := len(path)
					if path[pLen-4:pLen] == ".mp3" {
						mp3Relation := strings.Replace(path, wPath, "", 1)
						mp3Relation = strings.Replace(mp3Relation, "\\", "/", -1)
						soundFile := wPath + mp3Relation
						if fileutil.IsExist(soundFile) {
							mp3Relation = strings.Replace(mp3Relation, "/", "_", -1)
							lib.CopyFile(soundFile, lib.Data.BuildDstPath+"/resource/war3mapSound/"+i[2]+"/"+i[1]+mp3Relation)
							mp3Relation = strings.Replace(mp3Relation, ".mp3", "", -1)
							voiceData["vwp"][i[1]+mp3Relation] = strconv.Itoa(lib.VoiceDuration(soundFile))
						}
					}
					return nil
				})
				if err != nil {
					lib.Panic(err)
				}
				count4 = count4 + 1
			} else {
				pterm.Error.Println("【打声】：套件不存在 " + i[0])
				time.Sleep(time.Second)
			}
		}
	}
	pterm.Info.Println("【音效】Vcm 引入：" + strconv.Itoa(count) + "个")
	pterm.Info.Println("【音效】V3d 引入：" + strconv.Itoa(count2) + "个")
	pterm.Info.Println("【打声】Vwp 引入：" + strconv.Itoa(count4) + "个")
	pterm.Info.Println("【乐曲】Bgm 引入：" + strconv.Itoa(count3) + "个")

	//
	fn = L.GetGlobal("GO_RESULT_SLK")
	if err = L.CallByParam(lua.P{Fn: fn, NRet: 1, Protect: true}); err != nil {
		lib.Panic(err)
	}
	// get lua function results
	slkData := L.ToString(-1)
	var slData []map[string]interface{}
	err = json.Unmarshal([]byte(slkData), &slData)
	if err != nil {
		lib.Panic(err)
	}
	// 拼接魔码
	iniKeys := []string{"ability", "unit", "destructable"}
	iniF6 := make(map[string]string)
	reg, _ := regexp.Compile("\\[[A-Za-z][A-Za-z\\d]{3}]")
	var idIni []string
	for _, k := range iniKeys {
		iniF6[k] = ini(lib.Data.BuildDstPath + "/table/" + k + ".ini")
		matches := reg.FindAllString(iniF6[k], -1)
		for _, v := range matches {
			v = strings.Replace(v, "[", "", 1)
			v = strings.Replace(v, "]", "", 1)
			idIni = append(idIni, v)
		}
	}
	idIniByte, _ := json.Marshal(idIni)
	fn = L.GetGlobal("SLK_GO_INI")
	if err = L.CallByParam(lua.P{
		Fn:      fn,
		NRet:    0,
		Protect: true,
	}, lua.LString(idIniByte)); err != nil {
		lib.Panic(err)
	}
	slkIniBuilder := make(map[string]*strings.Builder)
	for _, k := range iniKeys {
		slkIniBuilder[k] = &strings.Builder{}
	}
	var idCli []string
	for _, sda := range slData {
		_slk := make(map[string]string)
		_hash := make(map[string]interface{})
		_id := ""
		_parent := ""
		for key, val := range sda {
			if key[:1] == "_" {
				_hash[key] = val
				if key == "_id" {
					switch v := val.(type) {
					case string:
						_id = v
					}
				} else if key == "_parent" {
					switch v := val.(type) {
					case string:
						_parent = "\"" + strings.Replace(v, "\\", "\\\\", -1) + "\""
					}
				}
			} else {
				var newVal string
				v := reflect.ValueOf(val)
				valType := reflect.TypeOf(val).Kind()
				switch valType {
				case reflect.String:
					newVal = "\"" + strings.Replace(v.String(), "\\", "\\\\", -1) + "\""
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
					reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					newVal = strconv.FormatInt(v.Int(), 10)
				case reflect.Float32, reflect.Float64:
					d := fmt.Sprintf("%v", v)
					if -1 == strings.Index(d, ".") {
						newVal = strconv.FormatInt(int64(v.Float()), 10)
					} else {
						newVal = strconv.FormatFloat(v.Float(), 'f', 2, 64)
					}
				case reflect.Slice, reflect.Array:
					newVal = "{"
					for i := 0; i < v.Len(); i++ {
						var n string
						vv := v.Index(i)
						ve := vv.Elem()
						d := fmt.Sprintf("%v", vv)
						switch ve.Kind() {
						case reflect.String:
							n = "\"" + d + "\""
						case reflect.Float64:
							if -1 == strings.Index(d, ".") {
								n = strconv.FormatInt(int64(ve.Float()), 10)
							} else {
								n = strconv.FormatFloat(ve.Float(), 'f', 2, 64)
							}
						}
						if newVal == "{" {
							newVal += n
						} else {
							newVal += "," + n
						}
					}
					newVal += "}"
				}
				_slk[key] = newVal
			}
		}
		if _id != "" && _parent != "" && len(_slk) > 0 {
			idCli = append(idCli, _id)
			_class := _hash["_class"].(string)
			sbr := slkIniBuilder[_class]
			sbr.WriteString("[" + _id + "]")
			sbr.WriteString("\n_parent=" + _parent)
			for k, v := range _slk {
				sbr.WriteString("\n" + k + "=" + v)
			}
			sbr.WriteString("\n\n")
		}
	}
	// 合并ini
	csTableDir := lib.Data.BuildDstPath + "/table"
	for k, v := range slkIniBuilder {
		if iniF6[k] == "" {
			err = lib.FilePutContents(csTableDir+"/"+k+".ini", v.String(), fs.ModePerm)

		} else {
			err = lib.FilePutContents(csTableDir+"/"+k+".ini", iniF6[k]+"\n\n"+v.String(), fs.ModePerm)
		}
		if err != nil {
			lib.Panic(err)
		}
	}
	// 为 setting 补充：物编ID
	settingContent, _ := settingMap["code"]
	if len(idCli) > 0 {
		for k, v := range idCli {
			idCli[k] = "'" + v + "'"
		}
		idCliStrs := strings.Join(idCli, ",")
		settingContent = strings.Replace(settingContent, "FRAMEWORK_GO_IDS = {}", "FRAMEWORK_GO_IDS = {"+idCliStrs+"}", 1)
	}
	// 补充：字体设定
	settingContent = strings.Replace(settingContent, "FRAMEWORK_FONT = ''", `FRAMEWORK_FONT = "`+font+`"`, 1)
	// 补充：音频设定
	var voiceCreator []string
	for k, v := range voiceData {
		for alias, dur := range v {
			switch k {
			case "vcm":
				voiceCreator = append(voiceCreator, `VcmCreator("`+alias+`",`+dur+`)`)
			case "v3d":
				voiceCreator = append(voiceCreator, `V3dCreator("`+alias+`",`+dur+`)`)
			case "vwp":
				voiceCreator = append(voiceCreator, `VwpCreator("`+alias+`",`+dur+`)`)
			}
		}
	}
	settingContent = strings.Replace(settingContent, "__FRAMEWORK_VOICE_INIT__ = 117532", strings.Join(voiceCreator, "\n"), 1)
	// map name
	wj, _ := fileutil.ReadFileToString(lib.Data.BuildDstPath + "/map/war3map.j")
	reg, _ = regexp.Compile("SetMapName\\(\"(.*)\"\\)")
	sm := reg.FindAllStringSubmatch(wj, 1)
	if len(sm) > 0 {
		mapName := sm[0][1]
		settingContent = strings.Replace(settingContent, "FRAMEWORK_MAP_NAME = ''", `FRAMEWORK_MAP_NAME = "`+mapName+`"`, 1)
	}
	// setting
	settingMap["code"] = settingContent

	// project sublibrary
	slbDir, _ := filepath.Abs(project + "/sublibrary")
	if fileutil.IsDir(slbDir) {
		err = filepath.Walk(slbDir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			pLen := len(path)
			if path[pLen-4:pLen] == ".lua" {
				lc, errc := fileutil.ReadFileToString(path)
				if errc != nil {
					lib.Panic(errc)
				}
				requireRelation := strings.Replace(path, slbDir, "", 1)
				requireRelation = requireRelation[1:]
				requireRelation = strings.Replace(requireRelation, ".lua", "", -1)
				requireRelation = strings.Replace(requireRelation, "\\", ".", -1)
				requireRelation = strings.Replace(requireRelation, "/", ".", -1)
				requireFile := strings.Replace(requireRelation, ".", "/", -1)
				name := "sublibrary." + requireRelation
				dst := lib.Data.BuildDstPath + "/map/sublibrary/" + requireFile + ".lua"
				code := lc
				builtIn = append(builtIn, map[string]string{
					"name": name,
					"dst":  dst,
					"code": code,
				})
				asyncRandIn = asyncRandReg(asyncRandIn, lc)
			}
			return nil
		})
		if err != nil {
			lib.Panic(err)
		}
	}

	// Plugins
	fn = L.GetGlobal("GO_RESULT_PLUGINS")
	if err = L.CallByParam(lua.P{Fn: fn, NRet: 1, Protect: true}); err != nil {
		lib.Panic(err)
	}
	plugins := L.ToString(-1)
	var pluginsData []string
	err = json.Unmarshal([]byte(plugins), &pluginsData)
	if err != nil {
		lib.Panic(err)
	}
	for _, i := range pluginsData {
		pluPath, _ := filepath.Abs(lib.Data.Assets + "/war3mapPlugins/" + i + "/plulibrary")
		isP := fileutil.IsDir(pluPath)
		if isP {
			pterm.Info.Println("【插件】引入：" + i)
			lib.Data.Plugins = append(lib.Data.Plugins, i)
			err = filepath.Walk(pluPath, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				pLen := len(path)
				if path[pLen-4:pLen] == ".lua" {
					lc, errc := fileutil.ReadFileToString(path)
					if errc != nil {
						lib.Panic(errc)
					}
					requireRelation := strings.Replace(path, pluPath, "", 1)
					requireRelation = requireRelation[1:]
					requireRelation = strings.Replace(requireRelation, ".lua", "", -1)
					requireRelation = strings.Replace(requireRelation, "\\", ".", -1)
					requireRelation = strings.Replace(requireRelation, "/", ".", -1)
					requireFile := strings.Replace(requireRelation, ".", "/", -1)
					name := "war3mapPlugins." + i + ".plulibrary." + requireRelation
					dst := lib.Data.BuildDstPath + "/map/war3mapPlugins/" + i + "/plulibrary/" + requireFile + ".lua"
					code := lc
					builtIn = append(builtIn, map[string]string{
						"name": name,
						"dst":  dst,
						"code": code,
					})
					asyncRandIn = asyncRandReg(asyncRandIn, lc)
				}
				return nil
			})
			if err != nil {
				lib.Panic(err)
			}
		} else {
			pterm.Error.Println("【插件】：" + i + " 不存在")
			time.Sleep(time.Second)
		}
	}

	// project scripts
	sDir, _ := filepath.Abs(project + "/scripts")
	err = filepath.Walk(sDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		pLen := len(path)
		if path[pLen-4:pLen] == ".lua" {
			lc, errc := fileutil.ReadFileToString(path)
			if errc != nil {
				lib.Panic(errc)
			}
			requireRelation := strings.Replace(path, sDir, "", 1)
			requireRelation = requireRelation[1:]
			requireRelation = strings.Replace(requireRelation, ".lua", "", -1)
			requireRelation = strings.Replace(requireRelation, "\\", ".", -1)
			requireRelation = strings.Replace(requireRelation, "/", ".", -1)
			requireFile := strings.Replace(requireRelation, ".", "/", -1)
			name := "scripts." + requireRelation
			dst := lib.Data.BuildDstPath + "/map/scripts/" + requireFile + ".lua"
			code := lc
			builtIn = append(builtIn, map[string]string{
				"name": name,
				"dst":  dst,
				"code": code,
			})
			asyncRandIn = asyncRandReg(asyncRandIn, lc)
		}
		return nil
	})
	if err != nil {
		lib.Panic(err)
	}

	// UI Kits
	fn = L.GetGlobal("GO_RESULT_UI")
	if err = L.CallByParam(lua.P{Fn: fn, NRet: 1, Protect: true}); err != nil {
		lib.Panic(err)
	}
	uis := L.ToString(-1)
	var uiData []string
	err = json.Unmarshal([]byte(uis), &uiData)
	if err != nil {
		lib.Panic(err)
	}
	var fdfs []string
	for _, i := range uiData {
		uiPath := lib.Data.Assets + "/war3mapUI/" + i
		isP := fileutil.IsDir(uiPath)
		if isP {
			uiTips := "【套件】引入：" + i
			mainLua := uiPath + "/main.lua"
			if !fileutil.IsExist(mainLua) {
				pterm.Error.Println("【套件】：" + i + " 核心main.lua未定义")
				time.Sleep(time.Second)
			} else {
				fdf := uiPath + "/main.fdf"
				uiTips += `，确认main`
				if fileutil.IsExist(fdf) {
					fdfs = append(fdfs, "UI\\"+i+".fdf")
					lib.CopyFile(fdf, lib.Data.BuildDstPath+"/map/UI/"+i+".fdf")
					uiTips += `，已引入fdf`
				}
				uiAssets := uiPath + "/assets"
				isA := fileutil.IsDir(uiAssets)
				if isA {
					lib.CopyPath(uiAssets, lib.Data.BuildDstPath+"/resource/war3mapUI/"+i+"/assets")
					uiTips += `，已引入assets`
				}
				// scripts
				if fileutil.IsDir(uiPath + "/scripts") {
					uiScripts, _ := filepath.Abs(uiPath + "/scripts")
					err = filepath.Walk(uiScripts, func(path string, info fs.FileInfo, err error) error {
						if err != nil {
							return err
						}
						pLen := len(path)
						if path[pLen-4:pLen] == ".lua" {
							lc, errc := fileutil.ReadFileToString(path)
							if errc != nil {
								lib.Panic(errc)
							}
							requireRelation := strings.Replace(path, uiScripts, "", 1)
							requireRelation = requireRelation[1:]
							requireRelation = strings.Replace(requireRelation, ".lua", "", -1)
							requireRelation = strings.Replace(requireRelation, "\\", ".", -1)
							requireRelation = strings.Replace(requireRelation, "/", ".", -1)
							requireFile := strings.Replace(requireRelation, ".", "/", -1)
							name := "war3mapUI." + i + ".scripts." + requireRelation
							dst := lib.Data.BuildDstPath + "/map/war3mapUI/" + i + "/scripts/" + requireFile + ".lua"
							code := lc
							builtIn = append(builtIn, map[string]string{
								"name": name,
								"dst":  dst,
								"code": code,
							})
							asyncRandIn = asyncRandReg(asyncRandIn, lc)
						}
						return nil
					})
					if err != nil {
						lib.Panic(err)
					}
					uiTips += `，已引入scripts`
				}
				// main
				name := "war3mapUI." + i + ".main"
				dst := lib.Data.BuildDstPath + "/map/war3mapUI/" + i + "/main.lua"
				code, errc := fileutil.ReadFileToString(mainLua)
				if errc != nil {
					lib.Panic(errc)
				}
				builtIn = append(builtIn, map[string]string{
					"name": name,
					"dst":  dst,
					"code": code,
				})
			}
			pterm.Info.Println(uiTips)
		} else {
			pterm.Error.Println("【套件】：" + i + " 不存在")
			time.Sleep(time.Second)
		}
	}
	// toc
	tocFile := lib.Data.BuildDstPath + "/map/UI/framework.toc"
	if !fileutil.IsExist(tocFile) {
		lib.Panic("tocFile not exist")
	}
	toc, err2 := fileutil.ReadFileToString(tocFile)
	if err2 != nil {
		lib.Panic(err2)
	}
	if len(fdfs) > 0 {
		toc += "\r\n" + strings.Join(fdfs, "\r\n")
	}
	err = lib.FilePutContents(tocFile, toc+"\r\n", fs.ModePerm)
	if err != nil {
		lib.Panic(err)
	}

	// japi
	japiCodes := ""
	if lib.Data.War3 != "" {
		s, _ := lib.Embeds.ReadFile("embeds/impl/japi.lua")
		japiCodes = string(s)
		if japiCodes != "" {
			implIn = append(implIn, japiCodes)
		}
	}
	// impl file
	implName := "impl"
	implDst := lib.Data.BuildDstPath + "/map/" + implName + ".lua"
	implCode := strings.Join(implIn, "\n")
	builtIn = append(builtIn, map[string]string{
		"name": implName,
		"dst":  implDst,
		"code": implCode,
	})

	// 处理代码
	for _, v := range builtIn {
		if v["name"] == "" {
			pterm.Debug.Println(v)
			lib.Panic("builtIn")
		}
		scriptIn = append(scriptIn, lib.RequireLua(v["name"]))
	}

	// core
	coreName := "core"
	coreDst := lib.Data.BuildDstPath + "/map/" + coreName + ".lua"
	coreByte, _ := lib.Embeds.ReadFile("embeds/builtIn/core.lua")
	coreCode := string(coreByte)
	coreCode = strings.Replace(coreCode, "--FRAMEWORK_SCRIPT_IN", strings.Join(scriptIn, "\n"), -1)
	coreMap := map[string]string{
		"name": coreName,
		"dst":  coreDst,
		"code": coreCode,
	}
	builtIn = append(builtIn, coreMap)
	builtIn = lib.EncryptLua(builtIn, assetsCheck)
	for _, v := range builtIn {
		if v["dst"] != "" {
			lib.DirCheck(v["dst"])
			err = lib.FilePutContents(v["dst"], v["code"], os.ModePerm)
			if err != nil {
				lib.Panic(err)
			}
		}
	}
	return coreMap["name"]
}
