package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

//GetConfig 獲取配置內容
func GetConfig(config, key string) (string, error) {
	flag1 := Strfind(config, key, 0)
	if flag1 == -1 {
		return "", errors.New("no key in this config. ")
	}
	flag1 = Strfind(config, ":", flag1)
	flag2 := Strfind(config, ",", flag1)
	if flag2 == -1 {
		flag2 = len(config)
	}
	return config[flag1+1 : flag2], nil
}

//Strfind 尋找字串
func Strfind(sors string, find string, ind int) int {
	findlen := len(find)
	for i := ind; i < len(sors)-findlen+1; i++ {
		if sors[i:i+findlen] == find {
			return i
		}
	}

	return -1
}

//Getdata 取得資料
func Getdata(inputdata []string, getinputname string) (string, error) {
	for i := 0; i < len(inputdata); i++ {
		//fmt.Println(inputdata[i])
		data := strings.Split(inputdata[i], "=")
		if data[0] == getinputname {
			dataa := strings.Split(data[1], ";")
			return dataa[0], nil
		}
	}
	//fmt.Println(getinputname)
	return "", errors.New("no id")
}

//GetNowtime 取得現在時間
func GetNowtime() string {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02 03:04:05")
}
func randInt(min int, max int) byte {
	rand.Seed(time.Now().UnixNano())
	return byte(min + rand.Intn(max-min))
}

//MakeToken 產生token
func MakeToken(l int) string {
	var result bytes.Buffer
	var temp byte
	for i := 0; i < l; {
		if randInt(65, 91) != temp {
			temp = randInt(65, 91)
			result.WriteByte(temp)
			i++
		}
	}
	return result.String()
}

//Checktime 比較兩個時間
func Checktime(time1, time2 string) bool {
	t1, _ := time.Parse("2006-01-02 15:04:05", time1)
	t2, err := time.Parse("2006-01-02 15:04:05", time2)
	if err == nil && t1.Before(t2) {
		return true
	}
	return false
}

//WebOutputJSON 用來回傳普通的執行結果json
func WebOutputJSON(w http.ResponseWriter, ret int, reason string, i interface{}) {
	out := &WebResult{ret, reason, i}
	b, err := json.Marshal(out)
	if err != nil {
		return
	}
	w.Write(b)
}

//WebOutputAPIJson 用來回傳api的json
func WebOutputAPIJson(w http.ResponseWriter, i interface{}) {
	b, err := json.Marshal(i)
	if err != nil {
		return
	}

	w.Write(b)
}

//RemoveNoreSpace 去掉多餘的空白
func RemoveNoreSpace(text string) string {
	text = strings.TrimSpace(text)
	flag := 0
	RetText := ""
	for i := 0; i < len(text); i++ {
		if text[i] == ' ' {
			flag = flag + 1
			if flag > 1 {
				continue
			}

		} else {
			flag = 0
		}
		RetText = RetText + string(text[i])

	}
	return RetText
}

//RunLocalCommand 執行命令返回資料
func RunLocalCommand(command string) (string, error) {
	RunCommand := exec.Command("/bin/sh", "-c", command)
	var out, er bytes.Buffer
	RunCommand.Stdout = &out
	RunCommand.Stderr = &er
	RunCommand.Start()
	RunCommand.Wait()
	RunCommandRet := out.String()
	RunCommandErr := er.String()
	if RunCommandErr != "" {

		return "", errors.New(RunCommandErr)
	}
	return RunCommandRet, nil

}

//RemovePoint 移除前後單引號
func RemovePoint(str string) string {
	fg := Strfind(str, "'", 0)
	if fg == -1 {
		return str
	}
	fg1 := Strfind(str, "'", fg+1)
	return str[fg+1 : fg1]
}

//RemoveDoublePoint 移除前後雙引號
func RemoveDoublePoint(str string) string {
	fg := Strfind(str, "\"", 0)
	if fg == -1 {
		return str
	}
	fg1 := Strfind(str, "\"", fg+1)
	return str[fg+1 : fg1]
}

//RemoveSliceLastSpace 移除slice最後一個空白
func RemoveSliceLastSpace(data []string) []string {
	if data[len(data)-1] == "" {
		data = append(data[:len(data)-1], data[len(data)-1+1:]...)
	}
	return data
}
