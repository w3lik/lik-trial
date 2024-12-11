package cmd

import (
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/pterm/pterm"
	"io/fs"
	"lik/lib"
	"os"
)

func New() {
	if lib.Data.ProjectName == "" {
		lib.Panic("不允许存在空名称项目")
	}
	if ProjectExist() {
		lib.Panic("已存在同名项目，你可以输入“run " + lib.Data.ProjectName + "”命令直接测试，或者请使用其他名称")
	}
	var err error
	// 如果没有 projects 目录则生成一个
	check := fileutil.IsDir(lib.Data.Projects)
	if !check {
		err = os.Mkdir(lib.Data.Projects, fs.ModePerm)
		if err != nil {
			lib.Panic(err)
		}
	}
	// 生成项目目录
	projectDir := lib.Data.Projects + "/" + lib.Data.ProjectName
	err = os.Mkdir(projectDir, fs.ModePerm)
	// 复制初始文件
	lib.CopyPath(lib.Data.Vendor+"/lni/map", projectDir+"/w3x/map")
	lib.CopyPath(lib.Data.Vendor+"/lni/table", projectDir+"/w3x/table")
	lib.CopyFile(lib.Data.Vendor+"/lni/war3mapMap.blp", projectDir+"/w3x/war3mapMap.blp")
	// 复制初始脚本
	lib.CopyPathEmbed(lib.Embeds, "embeds/new", projectDir)
	// 生成备份w3x目录
	Backup()
	pterm.Success.Println("项目创建完成！")
	pterm.Info.Println("你可以输入“we " + lib.Data.ProjectName + "”编辑地图信息")
	pterm.Info.Println("或可以输入“run " + lib.Data.ProjectName + "”命令直接调试")
}
