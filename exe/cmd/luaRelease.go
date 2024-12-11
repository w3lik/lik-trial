package cmd

import (
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/pterm/pterm"
	"lik/lib"
	"os"
	"strings"
)

func luaRelease() string {
	if lib.Data.BuildModeName != "_release" {
		lib.Panic("luaRelease")
	}
	// check dist
	spinner, _ := pterm.DefaultSpinner.Start("dist阶段测试结果检测中...")
	hotPath := lib.Data.Temp + "/_hot/" + lib.Data.ProjectName
	buildPath := lib.Data.Temp + "/_build/" + lib.Data.ProjectName
	distPath := lib.Data.Temp + "/_dist/" + lib.Data.ProjectName
	hotMod := lib.GetModTime(hotPath)
	buildMod := lib.GetModTime(buildPath)
	distMod := lib.GetModTime(distPath)
	if buildMod > distMod {
		pterm.Warning.Println("分析【析构】：检测到dist数据或比Build数据旧，为了避免打包到旧版本，可重新测试dist阶段")
	}
	if hotMod > distMod {
		pterm.Warning.Println("分析【析构】：检测到dist数据或比Hot数据旧，为了避免打包到旧版本，可重新测试dist阶段")
	}
	if !fileutil.IsDir(distPath) {
		spinner.Fail("分析【析构】：请先完成dist阶段测试")
		os.Exit(0)
	}
	// clone dist
	_ = os.RemoveAll(lib.Data.BuildDstPath)
	lib.CopyPath(distPath, lib.Data.BuildDstPath)
	// check connect、release
	connectFile := lib.Data.BuildDstPath + "/map/.connect"
	if !fileutil.IsExist(connectFile) {
		pterm.Error.Println("分析【析构】：dist阶段测试结果不符合过渡预期，请重新测试")
		os.Exit(0)
	}
	spinner.Success("分析【析构】：dist阶段测试结果检测已通过")
	// replace
	connectStr, err := fileutil.ReadFileToString(connectFile)
	if err != nil {
		lib.Panic(err)
	}
	_ = os.Remove(connectFile)
	connect := strings.Split(connectStr, "|")
	distFile := lib.Data.BuildDstPath + "/map/" + connect[1] + ".lua"
	releaseFile := lib.Data.BuildDstPath + "/map/" + connect[2] + ".lua"
	if !fileutil.IsExist(distFile) || !fileutil.IsExist(releaseFile) {
		pterm.Error.Println("分析【析构】：dist阶段测试结果不符合调试预期，请重新测试")
		os.Exit(0)
	}
	releaseCode, err2 := fileutil.ReadFileToString(releaseFile)
	if err2 != nil {
		lib.Panic(err2)
	}
	_ = os.Remove(releaseFile)
	err = lib.FilePutContents(distFile, releaseCode, os.ModePerm)
	if err != nil {
		lib.Panic(err)
	}
	_ = os.Remove(lib.Data.BuildDstPath + "/pack.w3x")
	pterm.Info.Println("分析【析构】：release核心已覆盖dist核心")
	return connect[0]
}
