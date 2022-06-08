/*************************************************************
     FileName: GoLand->src->finger.go
         Date: 2021/10/26 14:58
       Author: 苦咖啡
        Email: voilet@qq.com
         blog: http://blog.kukafei520.net
      Version: 0.0.1
      History:
**************************************************************/

package sub

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

func FileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func getMacAddrs() (macAddrs []string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return macAddrs
	}

	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}

		macAddrs = append(macAddrs, macAddr)
	}
	return macAddrs
}

func GetFinger() string {
	devopsDeiverInfo := "/etc/.devops_agent"
	devops_uid := "/usr/local/devops/.info"
	uid := FileExists(devopsDeiverInfo)
	dev_id := FileExists(devops_uid)
	if uid || dev_id {
		file, err := os.Open(devopsDeiverInfo)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		content, err := ioutil.ReadAll(file)
		return string(content)
	} else {
		str := getMacAddrs()[0]
		w := md5.New()
		_, _ = io.WriteString(w, str)
		md5str2 := fmt.Sprintf("%x", w.Sum(nil))
		err := ioutil.WriteFile(devopsDeiverInfo, []byte(md5str2), 0666)
		if err != nil {
			fmt.Println(err)
		}
		err2 := ioutil.WriteFile(devops_uid, []byte(md5str2), 0666)
		if err2 != nil {
			fmt.Println(err2)
		}
		return md5str2
	}
}

type RespInfo struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Pass string `json:"pass"`
	Port int    `json:"port"`
}

func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func GetToken(f, r string) string {
	return md5V(f + r)
}

func GetPort() int {
	randString := Generate(5)
	fin := GetFinger()
	tk := GetToken(fin, randString)
	url := "http://" + viper.GetString("ws.ip") + "/api/safe/info?fingerprint=" + fin + "&token=" + tk + "&rand=" + randString
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return 0
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	var resp RespInfo
	_ = json.Unmarshal(body, &resp)
	return resp.Port
}

func GetSshPort() int {
	port := 22
	fi, err := os.Open("/etc/ssh/sshd_config")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return port
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if CheckPort(string(a)) {
			st := strings.Split(string(a), " ")
			p, err := strconv.Atoi(st[1])
			if err != nil {
				port = 22
			} else {
				port = p
			}
		}
	}
	return port
}
func CheckPort(s string) bool {
	matched, err := regexp.MatchString("^Port", s)
	if err != nil {
		return false
	}
	return matched
}
