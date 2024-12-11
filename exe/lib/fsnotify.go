package lib

import (
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/fsnotify/fsnotify"
	"github.com/pterm/pterm"
	"io/fs"
	"strings"
	"sync"
)

import (
	"os"
	"path/filepath"
)

type NotifyFile struct {
	watch *fsnotify.Watcher
}

var (
	modTimes sync.Map
	notifyC  []string
)

func NewNotifyFile() *NotifyFile {
	w := new(NotifyFile)
	w.watch, _ = fsnotify.NewWatcher()
	return w
}

// WatchDir 监控目录
func (this *NotifyFile) WatchDir(dir string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if strings.Index(path, ".git") > 0 {
				return nil
			}
			paa, err2 := filepath.Abs(path)
			if err2 != nil {
				return err2
			}
			err2 = this.watch.Add(paa)
			if err2 != nil {
				return err2
			}
		}
		return nil
	})
	go this.WatchEvent()
}

func getName(name string) string {
	if name[len(name)-4:] != ".lua" {
		return ""
	}
	name = strings.Replace(name, Data.Pwd+"\\", "", -1)
	name = strings.Replace(name, ".lua", "", -1)
	name = strings.Replace(name, "\\", "/", -1)
	if strings.Index(name, "builtIn") > 0 {
		name = ""
	} else if strings.Index(name, "library") == 0 {
		name = strings.Replace(name, "library/", "", -1)
		name = strings.Replace(name, "/", ".", -1)
	} else if strings.Index(name, "projects") == 0 {
		name = strings.Replace(name, "projects/"+Data.ProjectName+"/", "", -1)
		name = strings.Replace(name, "/", ".", -1)
	} else if strings.Index(name, "assets/war3mapPlugins") == 0 {
		name = strings.Replace(name, "assets/war3mapPlugins/", "", -1)
		name = strings.Replace(name, "/", ".", -1)
	} else if strings.Index(name, "assets/war3mapUI") == 0 {
		name = strings.Replace(name, "assets/war3mapUI/", "", -1)
		name = strings.Replace(name, "/", ".", -1)
	} else {
		name = ""
	}
	return name
}

func create(name string) {
	exes := []string{"war3.exe"}
	if ExeRunningQty(exes) < 1 {
		os.Exit(0)
		return
	}
	name = getName(name)
	if name == "" {
		return
	}
	if !InArray(name, notifyC) {
		f := Data.War3 + "/fwhc.txt"
		notifyC = append(notifyC, name)
		err := FilePutContents(f, strings.Join(notifyC, "|"), fs.ModePerm)
		if err != nil {
			pterm.Error.Println(err.Error())
		}
	}
}

func modify(name string) {
	exes := []string{"war3.exe"}
	if ExeRunningQty(exes) < 1 {
		os.Exit(0)
		return
	}
	mt, _ := modTimes.Load(name)
	if mt != nil && mt.(int64) > 0 && mt.(int64) != GetModTime(name) {
		name = getName(name)
		if name == "" {
			return
		}
		f := Data.War3 + "/fwht.txt"
		var cons []string
		if fileutil.IsExist(f) {
			con, err := fileutil.ReadFileToString(f)
			if err != nil {
				pterm.Error.Println(err.Error())
				return
			}
			if len(con) > 0 {
				cons = strings.Split(con, "|")
			}
		}
		if !InArray(name, cons) {
			cons = append(cons, name)
			err := FilePutContents(f, strings.Join(cons, "|"), fs.ModePerm)
			if err != nil {
				pterm.Error.Println(err.Error())
				return
			}
		}
	}
	modTimes.Store(name, GetModTime(name))
}

func remove(name string) {
	exes := []string{"war3.exe"}
	if ExeRunningQty(exes) < 1 {
		os.Exit(0)
		return
	}
	name = getName(name)
	if name == "" {
		return
	}
	for k, v := range notifyC {
		if v == name {
			notifyC = append(notifyC[:k], notifyC[(k+1):]...)
			break
		}
	}
	f := Data.War3 + `/` + "fwhc.txt"
	if fileutil.IsExist(f) {
		err := FilePutContents(f, strings.Join(notifyC, "|"), fs.ModePerm)
		if err != nil {
			pterm.Error.Println(err.Error())
		}
	}
}

func (this *NotifyFile) WatchEvent() {
	defer this.watch.Close()
	for {
		select {
		case ev := <-this.watch.Events:
			{
				if ev.Op&fsnotify.Create == fsnotify.Create {
					file, err := os.Stat(ev.Name)
					if err == nil {
						if file.IsDir() {
							this.watch.Add(ev.Name)
						} else {
							create(ev.Name)
						}
					}
				}

				if ev.Op&fsnotify.Write == fsnotify.Write {
					modify(ev.Name)
				}

				if ev.Op&fsnotify.Remove == fsnotify.Remove || ev.Op&fsnotify.Rename == fsnotify.Rename {
					fi, err := os.Stat(ev.Name)
					if err == nil && fi.IsDir() {
						this.watch.Remove(ev.Name)
					} else {
						remove(ev.Name)
					}
				}

				//if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
				//	pterm.Info.Println("修改权限 : ", ev.Name)
				//}
			}
		case err := <-this.watch.Errors:
			{
				pterm.Error.Println("error : ", err)
				return
			}
		}
	}
}

func Hot() {
	if Data.War3 == "" {
		return
	}
	pterm.Success.Println("全局热更新生效中...")
	watch := NewNotifyFile()
	watch.WatchDir(Data.Pwd + "/library")
	watch.WatchDir(Data.Pwd + "/assets/war3mapPlugins")
	watch.WatchDir(Data.Pwd + "/assets/war3mapUI")
	watch.WatchDir(Data.Pwd + "/projects/" + Data.ProjectName)
	select {}
}
