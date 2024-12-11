package lib

import (
	"embed"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/pterm/pterm"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type DataStruct struct {
	Args          []string
	Pwd           string
	Log           string
	War3          string
	Temp          string
	Assets        string
	Plugins       []string
	Projects      string
	ProjectName   string
	BuildModeName string
	BuildDstPath  string
	Vendor        string
	Library       string
	Encrypt       string
	W3x2lni       string
	WE            string
}

var (
	//go:embed embeds
	Embeds embed.FS
	Data   DataStruct
)

func init() {
	rand.Seed(time.Now().UnixNano())
	pterm.EnableDebugMessages()
	pterm.Debug.Prefix.Text = "调整"
	pterm.Info.Prefix.Text = "提示"
	pterm.Success.Prefix.Text = "完成"
	pterm.Warning.Prefix.Text = "警告"
	pterm.Error.Prefix.Text = "错误"
	if Data.Pwd == "" {
		if !fileutil.IsExist("./conf") {
			CopyEmbed(Embeds, "embeds/conf.example", "./conf")
			pterm.Info.Println("请配置<conf>文件中的魔兽1.27版本客户端路径(WAR3)")
			os.Exit(0)
		}
		c, err := os.ReadFile("./conf")
		if err != nil {
			Panic(err)
		}
		content := string(c)
		reg, _ := regexp.Compile("#(.*)")
		content = reg.ReplaceAllString(content, "")
		content = strings.Replace(content, "\r\n", "\n", -1)
		content = strings.Replace(content, "\r", "\n", -1)
		split := strings.Split(content, "\n")
		conf := make(map[string]string)
		for _, iniItem := range split {
			if len(iniItem) > 0 {
				itemSplit := strings.Split(iniItem, "=")
				itemKey := strings.Trim(itemSplit[0], " ")
				itemKey = strings.ToLower(strings.Trim(itemSplit[0], " "))
				conf[itemKey] = strings.Trim(itemSplit[1], " ")
			}
		}
		if conf["root"] != "" {
			Data.Pwd = conf["root"]
		}
		if conf["war3"] != "" {
			Data.War3 = conf["war3"]
			if !fileutil.IsExist(Data.War3 + "/War3.exe") {
				pterm.Warning.Println("请正确配置<conf>文件中的魔兽客户端路径")
				pterm.Error.Println("当前魔兽路径无效:" + Data.War3)
				os.Exit(0)
			}
		} else {
			pterm.Info.Println("请配置<conf>文件中的魔兽客户端路径")
			pterm.Info.Println("可参考<conf.example>复制一个<conf>进行配置")
			os.Exit(0)
		}
		if Data.Pwd == "" {
			Data.Pwd, _ = os.Getwd()
		}
		Data.Pwd, _ = filepath.Abs(Data.Pwd)
		Data.Projects = Data.Pwd + "/projects"
		Data.Log = Data.Pwd + "/log"
		Data.Temp = Data.Pwd + "/temp"
		Data.Assets = Data.Pwd + "/assets"
		Data.Vendor = Data.Pwd + "/vendor"
		Data.Library = Data.Pwd + "/library"
		Data.Encrypt = Data.Pwd + "/encrypt"
		Data.W3x2lni = Data.Vendor + "/w3x2lni"
		Data.WE = Data.Vendor + "/WE"
		Data.Args = os.Args
		if !fileutil.IsDir(Data.Assets) || !fileutil.IsDir(Data.Vendor) || !fileutil.IsDir(Data.Library) || !fileutil.IsDir(Data.W3x2lni) || !fileutil.IsDir(Data.WE) {
			pterm.Error.Println("项目结构发现错误，已中止运行")
			os.Exit(0)
		}
		ck, e := fileutil.ReadFileToString(Data.Library + "/builtIn/setting.lua")
		if e != nil || !strings.Contains(ck, `I2023Y5#DTOX`) {
			pterm.Error.Println("非法核心，已中止运行")
			os.Exit(0)
		}
	}
	if len(Data.Args) >= 3 {
		Data.ProjectName = Data.Args[2]
	} else {
		Data.ProjectName = ""
	}
}
