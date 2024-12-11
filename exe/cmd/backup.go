package cmd

import (
	"github.com/pterm/pterm"
	"lik/lib"
	"os"
)

func Backup() {
	// 生成备份w3x目录
	tempDir := lib.Data.Temp + "/" + lib.Data.ProjectName
	w3xDir := lib.Data.Projects + "/" + lib.Data.ProjectName + "/w3x/map"
	if lib.GetModTime(tempDir+"/map") > lib.GetModTime(w3xDir) {
		_ = os.RemoveAll(w3xDir)
		lib.CopyPath(tempDir+"/map", w3xDir)
		pterm.Info.Println("备份完成[temp(地图备份)->w3x/map]")
	}
	war3mapMap := lib.Data.Projects + "/" + lib.Data.ProjectName + "/w3x/war3mapMap.blp"
	war3mapMapLni := lib.Data.Vendor + "/lni/war3mapMap.blp"
	if lib.GetModTime(tempDir+"/resource/war3mapMap.blp") > lib.GetModTime(war3mapMap) {
		_ = os.Remove(war3mapMap)
		lib.CopyFile(tempDir+"/resource/war3mapMap.blp", war3mapMap)
		pterm.Info.Println("更新同步[temp(war3mapMap)->w3x/war3mapMap]")
	} else if lib.GetModTime(war3mapMapLni) > lib.GetModTime(war3mapMap) {
		_ = os.Remove(war3mapMap)
		lib.CopyFile(war3mapMapLni, war3mapMap)
		pterm.Info.Println("更新同步[lni(war3mapMap)->w3x/war3mapMap]")
	}
	tableDir := lib.Data.Projects + "/" + lib.Data.ProjectName + "/w3x/table"
	if lib.GetModTime(tempDir+"/table") > lib.GetModTime(tableDir) {
		_ = os.RemoveAll(tableDir)
		lib.CopyPath(tempDir+"/table", tableDir)
		pterm.Success.Println("同步完成[temp(原生物编)->w3x/table]")
	}
}
