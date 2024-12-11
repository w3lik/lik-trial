package cmd

import (
	"github.com/pterm/pterm"
	"lik/lib"
	"os/exec"
	"path/filepath"
)

func Kill() {
	bat, err := filepath.Abs(lib.Data.WE + "/bin/kill.bat")
	if err != nil {
		lib.Panic(err)
	}
	cmd := exec.Command(bat)
	err = cmd.Run()
	if err != nil {
		lib.Panic(err)
	}
	pterm.Info.Println("已尝试关闭所有魔兽客户端")
}
