package main

import (
	"fmt"
	"github.com/pterm/pterm"
	"lik/cmd"
	"lik/lib"
)

func main() {
	fmt.Println("Welcome to try the lik.")
	fmt.Println("官方网站 https://www.hunzsig.com")
	fmt.Println("发电作者 https://afdian.com/a/hunzsig")
	if len(lib.Data.Args) < 2 {
		pterm.Error.Println("无效操作！")
		return
	}
	switch lib.Data.Args[1] {
	case "help":
		cmd.Help()
	case "new":
		cmd.New()
	case "we":
		cmd.WE()
	case "model":
		cmd.Model()
	case "clear":
		cmd.Clear()
	case "run":
		cmd.Run()
	case "multi":
		cmd.Multi()
	case "kill":
		cmd.Kill()
	default:
		pterm.Error.Println("命令: `" + lib.Data.Args[1] + "` 不存在!")
	}
}
