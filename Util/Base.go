package Util

import (
	"PanIndex/config"
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"io/fs"
	"math/rand"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ShortDur(d time.Duration) string {
	v, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", d.Seconds()), 64)
	return fmt.Sprint(v) + "s"
}
func GetCurBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		return ""
	}
	n = n + len(start)
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		return ""
	}
	str = string([]byte(str)[:m])
	return str
}

func GetRandomString(n int) string {
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

func GetNextOrPrevious(slice []fs.FileInfo, fs fs.FileInfo, flag int) fs.FileInfo {
	index := 0
	for p, v := range slice {
		if v.Name() == fs.Name() {
			index = p
		}
	}
	if flag == 0 {
		return fs
	} else if flag == -1 {
		if index > 0 {
			index = index + flag
		} else {
			return nil
		}
	} else if flag == 1 {
		if index < len(slice)-1 {
			index = index + flag
		} else {
			return nil
		}
	}
	return slice[index]
}

func FilterFiles(slice []fs.FileInfo, fullPath, sColumn, sOrder string) []fs.FileInfo {
	sort.Slice(slice, func(i, j int) bool {
		d1 := 0
		if slice[i].IsDir() {
			d1 = 1
		}
		d2 := 0
		if slice[j].IsDir() {
			d2 = 1
		}
		if d1 > d2 {
			return true
		} else if d1 == d2 {
			if sColumn == "file_name" {
				c := strings.Compare(slice[i].Name(), slice[j].Name())
				if sOrder == "desc" {
					return c >= 0
				} else {
					return c <= 0
				}
			} else if sColumn == "file_size" {
				if sOrder == "desc" {
					return slice[i].Size() >= slice[j].Size()
				} else {
					return slice[i].Size() <= slice[j].Size()
				}
			} else if sColumn == "last_op_time" {
				if sOrder == "desc" {
					return slice[i].ModTime().After(slice[j].ModTime())
				} else {
					return slice[i].ModTime().Before(slice[j].ModTime())
				}
			} else {
				return slice[i].ModTime().After(slice[j].ModTime())
			}
		} else {
			return false
		}
	})
	arr := []fs.FileInfo{}
	for _, v := range slice {
		fileId := filepath.Join(fullPath, v.Name())
		if config.GloablConfig.HideFileId != "" {
			listSTring := strings.Split(config.GloablConfig.HideFileId, ",")
			sort.Strings(listSTring)
			i := sort.SearchStrings(listSTring, fileId)
			if i < len(listSTring) && listSTring[i] == fileId {
				continue
			}
		}
		if !v.IsDir() {
			arr = append(arr, v)
		}
	}
	return arr
}

func DetermineEncoding(r *bufio.Reader) encoding.Encoding {
	_, err := r.Peek(1024)
	if err != nil {
		log.Error("get code error")
		return unicode.UTF8
	}
	return simplifiedchinese.GBK
}

func CheckPwd(PwdDirIds, path, pwd string) (bool, bool) {
	hasPath := false
	pwdOk := false
	s := strings.Split(PwdDirIds, ",")
	if PwdDirIds == "" {
		return hasPath, pwdOk
	}
	for _, v := range s {
		if strings.Split(v, ":")[0] == path {
			hasPath = true
		}
		if v == path+":"+pwd {
			pwdOk = true
		}
	}
	return hasPath, pwdOk
}
func GetPwdFromCookie(pwd, pathName string) string {
	s := strings.Split(pwd, ",")
	if len(s) > 0 {
		for _, v := range s {
			if strings.Split(v, ":")[0] == pathName {
				return strings.Split(v, ":")[1]
			}
		}
	}
	return ""
}

const (
	VAL   = 0x3FFFFFFF
	INDEX = 0x0000003D
)

var (
	alphabet = []byte("abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

/** implementation of short url algorithm **/
func Transform(longURL string) ([4]string, error) {
	md5Str := getMd5Str(longURL)
	//var hexVal int64
	var tempVal int64
	var result [4]string
	var tempUri []byte
	for i := 0; i < 4; i++ {
		tempSubStr := md5Str[i*8 : (i+1)*8]
		hexVal, err := strconv.ParseInt(tempSubStr, 16, 64)
		if err != nil {
			return result, nil
		}
		tempVal = int64(VAL) & hexVal
		var index int64
		tempUri = []byte{}
		for i := 0; i < 6; i++ {
			index = INDEX & tempVal
			tempUri = append(tempUri, alphabet[index])
			tempVal = tempVal >> 5
		}
		result[i] = string(tempUri)
	}
	return result, nil
}

/** generate md5 checksum of URL in hex format **/
func getMd5Str(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	c := m.Sum(nil)
	return hex.EncodeToString(c)
}
