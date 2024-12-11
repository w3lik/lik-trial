package lib

import (
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/mitchellh/go-ps"
	"github.com/pterm/pterm"
	"github.com/tcolgate/mp3"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var charChao map[string]string
var charChaoAnti map[string]string

func stack() string {
	var buf [2 << 10]byte
	res := string(buf[:runtime.Stack(buf[:], true)])
	res = strings.Replace(res, "Z:/Workspace/war3/lik/lik", "", -1)
	return res
}

func Panic(what interface{}) {
	t := reflect.TypeOf(what)
	switch t.Kind() {
	case reflect.String:
		pterm.Error.Println(what)
	case reflect.Ptr:
		pterm.Error.Println(what)
		pterm.Debug.Println(stack())
	default:
		pterm.Error.Println(t.Kind())
	}
	os.Exit(0)
}

// FilePutContents file_put_contents()
func FilePutContents(filename string, data string, mode os.FileMode) error {
	return os.WriteFile(filename, []byte(data), mode)
}

// GetModTime 获取文件(架)修改时间 返回unix时间戳
func GetModTime(path string) int64 {
	modTime := int64(0)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		u := info.ModTime().Unix()
		if u > modTime {
			modTime = u
		}
		return nil
	})
	if err != nil {
		return 0
	}
	return modTime
}

// Rand rand()
// Range: [0, 2147483647]
func Rand(min, max int) int {
	if min > max {
		panic("min: min cannot be greater than max")
	}
	// PHP: getrandmax()
	if int31 := 1<<31 - 1; max > int31 {
		panic("max: max can not be greater than " + strconv.Itoa(int31))
	}
	if min == max {
		return min
	}
	return rand.Intn(max+1-min) + min
}

// InArray in_array()
// haystack supported types: slice, array or map
func InArray(needle interface{}, haystack interface{}) bool {
	val := reflect.ValueOf(haystack)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(needle, val.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(needle, val.MapIndex(k).Interface()) {
				return true
			}
		}
	default:
		panic("haystack: haystack type muset be slice, array or map")
	}

	return false
}

func ExeRunningQty(names []string) int {
	qty := 0
	pa, _ := ps.Processes()
	for _, p := range pa {
		for _, n := range names {
			if strings.ToLower(p.Executable()) == strings.ToLower(n) {
				qty += 1
			}
		}
	}
	return qty
}

func Nano(n int) string {
	if n < 1 {
		return ""
	}
	var err error
	s1 := ""
	s2 := ""
	if n > 5 {
		s1, err = gonanoid.Generate("_abcDEFghiJKLmnopqrSTUvwxYZ", 5)
		if err != nil {
			Panic(err)
		}
		s2, err = gonanoid.Generate("_0123456789abcDEFghiJKLmnopqrSTUvwxYZ", n-5)
		if err != nil {
			Panic(err)
		}
	} else {
		s1, err = gonanoid.Generate("_abcDEFghiJKLmnopqrSTUvwxYZ", n)
		if err != nil {
			Panic(err)
		}
	}
	return s1 + s2
}

func NanoOL(n int) string {
	s, err := gonanoid.Generate("abcDEFghiJKLmnopqrSTUvwxYZ", n)
	if err != nil {
		Panic(err)
	}
	return s
}

func Zebra(n int) string {
	letterStr := "abcdefghijklmnopqrstuvwxyxzABCDEFGHIJKLMNOPQRSTUVWXYXZ"
	numberStr := "0123456789"
	if n < 4 {
		n = 4
	}
	if n%2 != 0 {
		n = n + 1
	}
	r := ""
	for j := 0; j < n; j++ {
		if j%2 == 0 {
			k := rand.Intn(len(letterStr) - 1)
			r += letterStr[k:(k + 1)]
		} else {
			k := rand.Intn(len(numberStr) - 1)
			r += numberStr[k:(k + 1)]
		}
	}
	return r
}

func RequireLua(s string) string {
	return "require(\"" + s + "\")"
}

func VoiceDuration(soundFile string) int {
	skipped := 0
	r, err := os.Open(soundFile)
	defer r.Close()
	if err != nil {
		Panic(err)
		return 0
	}
	d := mp3.NewDecoder(r)
	var f mp3.Frame
	dur := 0
	for {
		if err := d.Decode(&f, &skipped); err != nil {
			if err.Error() == "EOF" {
				break
			} else {
				Panic(err)
				return 0
			}
		}
		dur += int(float64(time.Millisecond) * (1000 / float64(f.Header().SampleRate())) * float64(f.Samples()))
	}
	dur = dur / 1000000
	return dur
}

func CharRand() string {
	allStr := "01234abcdefghijklmnopqrstuvwxyxzABCDEFGHIJKLMNOPQRSTUVWXYXZ56789"
	i := rand.Intn(len(allStr) - 1)
	return allStr[i:(i + 1)]
}

func CharChao() (map[string]string, map[string]string) {
	if charChao == nil {
		charChao = make(map[string]string)
		charChaoAnti = make(map[string]string)
		allStr := `01234abcdefghijklmnopqrstuvwxyxz._:!,/\+-ABCDEFGHIJKLMNOPQRSTUVWXYXZ56789`
		for i := 0; i < len(allStr); i++ {
			r := Zebra(10)
			charChao[allStr[i:(i+1)]] = r
			charChaoAnti[r] = allStr[i:(i + 1)]
		}
	}
	return charChao, charChaoAnti
}

func MBSplit(str string, size int) []string {
	var sp []string
	lenInByte := len(str)
	if lenInByte <= 0 {
		return sp
	}
	count := 0
	i0 := 0
	i := 0
	for {
		if i >= lenInByte {
			break
		}
		curByte := str[i]
		byteLen := 1
		if curByte > 0 && curByte <= 127 {
			byteLen = 1 // 1字节字符
		} else if curByte >= 192 && curByte <= 223 {
			byteLen = 2 // 双字节字符
		} else if curByte >= 224 && curByte <= 239 {
			byteLen = 3 // 汉字
		} else if curByte >= 240 && curByte <= 247 {
			byteLen = 4 // 4字节字符
		}
		count = count + 1 // 字符的个数（长度）
		i = i + byteLen   // 重置下一字节的索引
		if count >= size {
			sp = append(sp, str[i0:i])
			i0 = i
			count = 0
		} else if i > lenInByte {
			sp = append(sp, str[i0:lenInByte])
		}
	}
	return sp
}
