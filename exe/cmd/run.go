package cmd

import (
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/pterm/pterm"
	"lik/lib"
	"os"
	"os/exec"
	"time"
)

func runTest(mode string, w3xFire string, times int) {
	cmd := exec.Command(lib.Data.WE+"/bin/WEconfig.exe", "-launchwar3", "-loadfile", w3xFire)
	_, _ = cmd.Output()
	pterm.Info.Println("尝试启动中")
	ticker := time.NewTicker(time.Second)
	pterm.Info.Println(<-ticker.C)
	exes := []string{"war3.exe"}
	if lib.ExeRunningQty(exes) > 0 {
		pterm.Success.Println("War3已启动 " + lib.Data.War3)
		if mode == "-h" {
			// 热更新
			if lib.Data.War3 != "" {
				lib.Hot()
			}
		}
	} else {
		if times >= 3 {
			pterm.Error.Println("War3启动失败，请检查环境")
			return
		}
		pterm.Warning.Println("War3启动失败,1秒后再次重试")
		time.Sleep(time.Second)
		runTest(mode, w3xFire, times+1)
	}
}

func Run() {
	if !ProjectExist() {
		if lib.Data.ProjectName != "" {
			lib.Panic("项目不存在：" + lib.Data.ProjectName)
		}
	}
	// 模式
	cache := false
	mode := "-h"
	lib.Data.BuildModeName = ""
	modeLni := "slk"
	if len(os.Args) >= 4 {
		if len(os.Args[3]) == 3 {
			mode = os.Args[3][0:2]
			if os.Args[3][2:3] == "~" {
				cache = true
			}
		} else if len(os.Args[3]) == 2 {
			mode = os.Args[3]
		}
	}
	if mode == "-t" {
		lib.Data.BuildModeName = "_test"
		modeLni = "obj"
		cache = false
	} else if mode == "-h" {
		lib.Data.BuildModeName = "_hot"
		modeLni = "obj"
	} else if mode == "-b" {
		lib.Data.BuildModeName = "_build"
	} else if mode == "-d" {
		lib.Data.BuildModeName = "_dist"
	} else if mode == "-r" {
		lib.Data.BuildModeName = "_release"
	} else {
		lib.Panic("模式错误：" + mode)
	}
	lib.Data.BuildDstPath = lib.Data.Temp + "/" + lib.Data.BuildModeName + "/" + lib.Data.ProjectName
	dstW3xFire := lib.Data.BuildDstPath + "/pack.w3x"
	if cache == false {
		if lib.Data.BuildModeName == "_release" {
			pterm.Debug.Println("准备发布打包")
		} else {
			spinner, _ := pterm.DefaultSpinner.Start("构建地图中")
			temProjectDir := lib.Data.Temp + "/" + lib.Data.ProjectName
			temProjectW3xFire := lib.Data.Temp + "/" + lib.Data.ProjectName + ".w3x"
			buoyFire := lib.Data.Temp + "/" + lib.Data.ProjectName + "/.we"
			mtW := lib.GetModTime(temProjectW3xFire)
			mtB := lib.GetModTime(buoyFire)
			if mtW > mtB {
				// 如果地图文件比we打开时新（说明有额外保存过）把保存后的文件拆包并同步
				cmd := exec.Command(lib.Data.W3x2lni+"/w2l.exe", "lni", temProjectW3xFire)
				_, err := cmd.Output()
				if err != nil {
					lib.Panic(err)
				}
				_ = os.Remove(buoyFire)
				lib.CopyFile(lib.Data.Vendor+"/lni/.we", buoyFire)
				Backup() // 以编辑器为主版本
				pterm.Info.Println("同步完毕[检测到有新的地图保存行为，以‘WE’为主版本]")
			} else {
				Pickup() // 以project为主版本
				pterm.Info.Println("同步完毕[检测到没有新的地图保存行为，以‘project’为主版本]")
			}
			// 复制一份最新的临时文件
			_ = os.RemoveAll(lib.Data.BuildDstPath)
			lib.CopyPath(temProjectDir, lib.Data.BuildDstPath)
			spinner.Success("构建完毕：" + lib.Data.BuildModeName)
		}
		// 调整代码，以支持war3
		t := time.Now()
		War3()
		pterm.Success.Println("资源及代码处理完成，耗时：" + time.Since(t).String())
		// 生成地图
		spinner, _ := pterm.DefaultSpinner.Start("打包地图中")
		t = time.Now()
		cmd := exec.Command(lib.Data.W3x2lni+"/w2l.exe", modeLni, lib.Data.BuildDstPath, dstW3xFire)
		_, err := cmd.Output()
		if err != nil {
			lib.Panic(err)
		}
		// 检查标志
		spinner.Success(`打包地图完成[` + modeLni + `]，耗时` + time.Since(t).String())
		if mode == "-t" {
			pterm.Info.Println(`>>> 临时地图已生成，位置:` + dstW3xFire + ` <<<`)
			return
		}
	} else {
		// 检查数据
		if !fileutil.IsDir(lib.Data.BuildDstPath) {
			pterm.Error.Println("未找到最近一次的地图资源缓存数据")
			return
		}
		// 若存在删除pack.w3x，以免递归打包增加结果文件大小
		if fileutil.IsExist(dstW3xFire) {
			_ = os.Remove(dstW3xFire)
		}
		// 生成地图
		spinner, _ := pterm.DefaultSpinner.Start("打包地图中")
		t := time.Now()
		cmd := exec.Command(lib.Data.W3x2lni+"/w2l.exe", modeLni, lib.Data.BuildDstPath, dstW3xFire)
		_, err := cmd.Output()
		if err != nil {
			lib.Panic(err)
		}
		pterm.Info.Println(`使用最近一次[` + lib.Data.BuildModeName + `]地图资源缓存数据`)
		spinner.Success(`打包地图完成[` + modeLni + `]，耗时` + time.Since(t).String())
	}
	if lib.Data.War3 != "" {
		_ = os.Remove(lib.Data.War3 + "/fwht.txt")
		_ = os.Remove(lib.Data.War3 + "/fwhc.txt")
		_ = os.Remove(lib.Data.War3 + "/dz_w3_plugin.dll")
		_ = os.Remove(lib.Data.War3 + "/version.dll")
		wtPath := lib.Data.War3 + "/Maps/Test"
		if !fileutil.IsDir(wtPath) {
			_ = os.Mkdir(wtPath, os.ModePerm)
		}
	}
	pterm.Info.Println("即将准备地图测试")
	exes := []string{"war3.exe"}
	if lib.ExeRunningQty(exes) > 0 {
		pterm.Warning.Println(">>> 请先关闭当前war3!!! <<<")
		return
	}
	// 跑
	runTest(mode, dstW3xFire, 0)
}
