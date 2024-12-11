package cmd

import (
	"github.com/pterm/pterm"
	"lik/lib"
	"os"
)

func Clear() {
	if lib.Data.War3 != "" {
		_ = os.Remove(lib.Data.War3 + "/dz_w3_plugin.ini")
		pterm.Success.Println(`清理魔兽存档完毕`)
	}
	_ = os.RemoveAll(lib.Data.Temp)
	pterm.Success.Println(`清理临时区完毕`)
}
