package cmd

import (
	"github.com/duke-git/lancet/v2/fileutil"
	"lik/lib"
)

func ProjectExist() bool {
	if len(lib.Data.Args) <= 2 {
		return false
	}
	projectName := lib.Data.Args[2]
	if projectName[:1] == "_" {
		lib.Panic("项目名不合法(下划线“_”开始的名称已被禁用)")
	}
	projectDir := lib.Data.Projects + "/" + projectName
	return fileutil.IsDir(projectDir)
}
