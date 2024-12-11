package lib

import (
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/pterm/pterm"
	"github.com/samber/lo"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func ZipLua(src string) string {
	content, err := fileutil.ReadFileToString(src)
	if err != nil {
		content = src
	}
	content = strings.Replace(content, "\r\n", "\n", -1)
	content = strings.Replace(content, "\r", "\n", -1)
	reg, _ := regexp.Compile("\\s*--.*\\[\\[[\\s\\S]*?\\]\\]")
	content = reg.ReplaceAllString(content, "")
	reg, _ = regexp.Compile("\\s*--.*")
	content = reg.ReplaceAllString(content, "")
	cta := strings.Split(content, "\n")
	var ctn []string
	for _, c := range cta {
		c = strings.Trim(c, " ")
		if len(c) > 0 {
			ctn = append(ctn, c)
		}
	}
	return strings.Join(ctn, " ")
}

// EncryptLua 混淆Lua
func EncryptLua(data []map[string]string, assstsIn []string) []map[string]string {
	// 智能分析混淆
	var csvData [][]string
	var _codesBuilder strings.Builder
	spiltSym := ` 0X0X0X0X0X0X0 `
	spinner, _ := pterm.DefaultSpinner.Start("智能分析启动中...")
	// 特殊优化index
	blizzardIdx := -1
	for i, v := range data {
		if v["name"] == "blizzard" {
			blizzardIdx = i
		}
		_codesBuilder.WriteString(v["code"] + ` `)
	}
	// 优化blizzard
	if blizzardIdx >= 0 {
		mixed := _codesBuilder.String()
		reg, _ := regexp.Compile(`(?m)^(\w+) = `)
		blizzardKeys := reg.FindAllStringSubmatch(data[blizzardIdx]["code"], -1)
		for _, b := range blizzardKeys {
			reg, _ = regexp.Compile(`\b` + b[1] + `\b`)
			c := strings.Count(mixed, b[1])
			if c < 2 {
				data[blizzardIdx]["code"] = strings.Replace(data[blizzardIdx]["code"], b[1]+` =`, `--`+b[1]+` =`, -1)
			}
		}
	}
	_codesBuilder.Reset()
	//
	for _, v := range data {
		_codesBuilder.WriteString(ZipLua(v["code"]) + spiltSym)
	}
	content := _codesBuilder.String()
	_codesBuilder.Reset()
	//
	spinner.UpdateText("中文汉字等效分析中...")
	reg, _ := regexp.Compile(`[一-龥]{1,}`)
	mcn := reg.FindAllString(content, -1)
	wcn := make(map[string]int)
	wci := rand.Intn(33) + 3
	parseCheck := make(map[string]int)
	cnParseKV := make(map[string]string)
	var cnParse [][]string
	if len(mcn) > 0 {
		for _, s := range mcn {
			sl := MBSplit(s, 1)
			var code []string
			for _, c := range sl {
				if wcn[c] == 0 {
					wcn[c] = wci
					wci = wci + 1 + rand.Intn(11)
					if cnParseKV[c] == "" {
						cnParseKV[c] = `[` + strconv.Itoa(wcn[c]) + `]=` + `'` + c + `'`
					}
				}
				code = append(code, strconv.Itoa(wcn[c]))
			}
			if parseCheck[s] == 0 {
				parseCheck[s] = 1
				cnParse = append(cnParse, []string{s, "J.N2C(" + strings.Join(code, ",") + ")"})
			} else {
				parseCheck[s] += 1
			}
		}
	}
	sort.Slice(cnParse, func(i, j int) bool {
		return len(cnParse[i][0]) > len(cnParse[j][0])
	})
	pterm.Info.Println("分析【混淆】：中文汉字分析出 " + strconv.Itoa(len(cnParse)) + " 个预备数据")
	//
	parseCheck = make(map[string]int)
	spinner.UpdateText("字串等效分析中...")
	var strParse []string
	reg, _ = regexp.Compile(`"[0-9A-Za-z._:!,/+\\-]{2,}"`)
	m := reg.FindAllString(content, -1)
	if len(m) > 0 {
		for _, s := range m {
			strParse = append(strParse, s)
			s2 := s[1 : len(s)-1]
			parseCheck[s2] += 1
		}
	}
	strParse = lo.Uniq[string](strParse)
	pterm.Info.Println("分析【混淆】：字串等效分析出 " + strconv.Itoa(len(strParse)) + " 个预备数据")

	// 资源检查
	for _, a := range assstsIn {
		if parseCheck[a] <= 1 {
			pterm.Warning.Println("分析【资源】：" + a + " 可能未曾使用")
		}
	}

	// 中文汉字混淆
	spinner.UpdateText("中文汉字混淆中...")
	for _, p := range cnParse {
		k, v := p[0], p[1]
		csvData = append(csvData, []string{"汉字混淆", k, v})
		content = strings.Replace(content, `"`+k+`"`, v, -1)
		content = strings.Replace(content, `<`+k+`>'`, `<"..`+v+`..">`, -1)
		content = strings.Replace(content, `[`+k+`]`, `["..`+v+`.."]`, -1)
	}
	var parseSetArr []string
	for _, v := range cnParseKV {
		parseSetArr = append(parseSetArr, v)
	}
	content = strings.Replace(content, `FRAMEWORK_N2C = {}`, `FRAMEWORK_N2C = {`+strings.Join(parseSetArr, ",")+`}`, 1)
	pterm.Success.Println("分析【混淆】：中文汉字已混淆 " + strconv.Itoa(len(cnParse)) + " 个数据")

	// 字串混淆
	spinner.UpdateText("双引号字串混淆中...")
	strParseKV := make(map[string]string)
	chao, _ := CharChao()
	for k, v := range chao {
		csvData = append(csvData, []string{"Chao", k, v})
	}
	for _, str := range strParse {
		eky := str[1 : len(str)-1]
		ess := strings.Split(strings.Replace(eky, `\\`, `\`, -1), "")
		for k, e := range ess {
			chae := strings.Split(chao[e], "")
			for x, y := 0, len(chae)-1; x < y; x, y = x+1, y-1 {
				chae[x], chae[y] = chae[y], chae[x]
			}
			ess[k] = strings.Join(chae, "")
		}
		for i, j := 0, len(ess)-1; i < j; i, j = i+1, j-1 {
			ess[i], ess[j] = ess[j], ess[i]
		}
		mysterious := Zebra(8)
		esd := "___('" + mysterious + strings.Join(ess, mysterious) + "')"
		strParseKV[eky] = esd
	}
	strParseKV2 := make(map[string]string)
	var strParse2 []string
	for _, v := range strParseKV {
		strParseKV2[v] = Nano(13)
		strParse2 = append(strParse2, strParseKV2[v]+` = `+v)
	}
	coreStr := strings.Join(strParse2, ` `)
	content = strings.Replace(content, `FRAMEWORK_STRING()`, coreStr+` `, 1)
	for k, v := range strParseKV {
		v2 := strParseKV2[v]
		csvData = append(csvData, []string{"字串混淆", k, v2})
		content = strings.Replace(content, `"`+k+`"`, v2, -1)
	}
	pterm.Success.Println("分析【混淆】：字串混淆已混淆 " + strconv.Itoa(len(strParse)) + " 个数据")

	spinner.UpdateText("处理词根中...")
	// 删除特定内容
	del, repl, force, facades, meta := encryptAnalysis()
	for _, w := range del {
		csvData = append(csvData, []string{"销毁词根", w})
		content = strings.Replace(content, w, "", -1)
	}
	pterm.Success.Println("分析【混淆】：销毁词根已处理 " + strconv.Itoa(len(del)) + " 个数据")
	wl := 0
	for _, w := range repl {
		wn := NanoOL(Rand(7, 11))
		csvData = append(csvData, []string{"替换词根", w, wn})
		content = strings.Replace(content, w, wn, -1)
		wl += 1
	}
	pterm.Success.Println("分析【混淆】：替换词根已处理 " + strconv.Itoa(wl) + " 个数据")
	for _, w := range facades {
		wn := NanoOL(Rand(7, 11))
		csvData = append(csvData, []string{"门面词根", w, wn})
		reg, _ = regexp.Compile(`\b` + w + `\(`)
		content = reg.ReplaceAllString(content, wn+`(`)
	}
	pterm.Success.Println("分析【混淆】：门面词根已处理 " + strconv.Itoa(len(facades)) + " 个数据")
	for _, w := range meta {
		wn := NanoOL(Rand(7, 11))
		csvData = append(csvData, []string{"函数词根", w, wn})
		content = strings.Replace(content, ` :`+w+`(`, ` :`+wn+`(`, -1)
		content = strings.Replace(content, `:`+w+`(`, `:`+wn+`(`, -1)
		content = strings.Replace(content, `):`+w+`(`, `):`+wn+`(`, -1)
		content = strings.Replace(content, `(class).`+w+`(self`, `(class).`+wn+`(self`, -1)
	}
	pterm.Success.Println("分析【混淆】：函数词根已处理 " + strconv.Itoa(len(meta)) + " 个数据")
	for _, w := range force {
		wn := NanoOL(Rand(7, 11))
		csvData = append(csvData, []string{"强制词根", w, wn})
		reg, _ = regexp.Compile(`\b` + w + `\b`)
		content = reg.ReplaceAllString(content, wn)
	}
	pterm.Success.Println("分析【混淆】：强制词根已处理 " + strconv.Itoa(len(force)) + " 个数据")
	// csv
	csvDir := Data.Temp + `/_encrypt/`
	if !fileutil.IsDir(csvDir) {
		_ = os.Mkdir(csvDir, os.ModePerm)
	}
	csvFile := csvDir + Data.ProjectName + Data.BuildModeName + `.csv`
	_ = os.Remove(csvFile)
	var csvStrs []string
	for _, v := range csvData {
		v[len(v)-1] = strings.Replace(v[len(v)-1], `,`, `，`, -1)
		csvStrs = append(csvStrs, strings.Join(v, `,`))
	}
	err := FilePutContents(csvFile, "\xEF\xBB\xBF"+strings.Join(csvStrs, "\n"), os.ModePerm)
	if err != nil {
		Panic(err)
	}
	contents := strings.Split(content, spiltSym)
	for i, v := range data {
		v["code"] = contents[i]
	}
	spinner.Success("分析【混淆】：合计" + strconv.Itoa(len(data)) + " 段代码已经处理完毕")
	return data
}
