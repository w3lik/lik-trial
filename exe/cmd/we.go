package cmd

import (
	"github.com/pterm/pterm"
	"lik/lib"
	"os"
	"os/exec"
)

func WE() {
	if len(os.Args) <= 2 {
		cmd := exec.Command(lib.Data.WE + "/WE.exe")
		_, err := cmd.Output()
		if err != nil {
			lib.Panic(err)
		}
		pterm.Info.Println("WE正在打开")
		return
	}
	if !ProjectExist() {
		lib.Panic("项目" + lib.Data.ProjectName + "不存在")
	}
	// 判断是否已有we在修改项目，提示一下
	exes := []string{"worldedit.exe", "worldeditydwe.exe", "worldeditkkwe.exe"}
	if lib.ExeRunningQty(exes) > 0 {
		pterm.Warning.Println("提示：检测到已有WE开启中，如果你是重复调用了we命令，请保留一个进行修改!")
	}
	w3xDir := lib.Data.Temp + "/" + lib.Data.ProjectName
	w3xFire := lib.Data.Temp + "/" + lib.Data.ProjectName + ".w3x"
	// 检查上一次we的修改数据是否未保存
	buoyFire := lib.Data.Temp + "/" + lib.Data.ProjectName + "/.we"
	mtW := lib.GetModTime(w3xFire)
	mtB := lib.GetModTime(buoyFire)
	if mtW > mtB {
		// 如果地图文件比we打开时新（说明有额外保存过）把保存后的文件拆包并同步
		cmd := exec.Command(lib.Data.W3x2lni+"/w2l.exe", "lni", w3xFire)
		_, err := cmd.Output()
		if err != nil {
			lib.Panic(err)
		}
		Backup() // 以编辑器为主版本
		pterm.Info.Println("同步完毕[检测到之前有过使用we命令进行地图保存行为，正在进行同步备份]")
	}
	_ = os.Remove(buoyFire)
	//
	Pickup()
	cmd := exec.Command(lib.Data.W3x2lni+"/w2l.exe", "obj", w3xDir, w3xFire)
	_, err := cmd.Output()
	if err != nil {
		lib.Panic(err)
	}
	lib.CopyFile(lib.Data.Vendor+"/lni/.we", buoyFire)
	cmd = exec.Command(lib.Data.WE+"/WE.exe", "-loadfile", w3xFire)
	_, err = cmd.Output()
	if err != nil {
		lib.Panic(err)
	}
	pterm.Info.Println("WE正在配图并打开")
}
