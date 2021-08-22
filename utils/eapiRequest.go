package utils

import (
	"crypto/md5"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	netUrl "net/url"
	"strconv"
	"strings"
	"time"
)

func EapiRequest(eapiOption eapiOption, options RequestData) (result string, err error) {
	data := SpliceStr(eapiOption.Path, eapiOption.Json)
	answer, err := CreateNewRequest(Format2Params(data), eapiOption.Url, options)
	if err == nil {
		var decrypted []byte
		decrypted = AesDecryptECB([]byte(answer))
		if log.GetLevel() == log.DebugLevel {
			log.Debugf("[EapiRespBodyJson]: %s", string(decrypted))
		}
		return string(decrypted), nil
	}
	return "", err
}

func SpliceStr(path string, data string) (result string) {
	text := fmt.Sprintf("nobody%suse%smd5forencrypt", path, data)
	MD5 := md5.Sum([]byte(text))
	md5str := fmt.Sprintf("%x", MD5)
	result = fmt.Sprintf("%s-36cd479b6b5-%s-36cd479b6b5-%s", path, data, md5str)
	return result
}

func Format2Params(str string) (data string) {
	data = fmt.Sprintf("params=%X", AesEncryptECB(str))
	return data
}

func ChooseUserAgent() string {
	userAgentList := []string{
		"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1",
		"Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 5.1.1; Nexus 6 Build/LYZ28E) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Mobile Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Mobile/14F89;GameHelper",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_0 like Mac OS X) AppleWebKit/602.1.38 (KHTML, like Gecko) Version/10.0 Mobile/14A300 Safari/602.1",
		"NeteaseMusic/6.5.0.1575377963(164);Dalvik/2.1.0 (Linux; U; Android 9; MIX 2 MIUI/V12.0.1.0.PDECNXM)",
	}
	rand.Seed(time.Now().UnixNano())
	var index int
	index = rand.Intn(len(userAgentList))
	return userAgentList[index]
}

func encodeURIComponent(str string) string {
	r := netUrl.QueryEscape(str)
	r = strings.Replace(r, "+", "%20", -1)
	return r
}

func CreateNewRequest(data string, url string, options RequestData) (answer string, err error) {
	client := &http.Client{}
	reqBody := strings.NewReader(data)
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return "", err
	}

	cookie := map[string]interface{}{}
	for i := 0; i < len(options.Cookies); i++ {
		cookie[options.Cookies[i].Key] = options.Cookies[i].Value
	}

	csrfValue, isok := cookie["__csrf"]
	csrfToken := ""
	if isok {
		csrfToken = fmt.Sprintf("%v", csrfValue)
	}
	header := make(map[string]interface{})
	keys := [...]string{"osver", "deviceId", "mobilename", "channel"}
	for _, val := range keys {
		value, ok := cookie[val]
		if ok {
			header[val] = value
		}
	}
	header["appver"] = func() string {
		val, ok := cookie["appver"]
		if ok {
			return fmt.Sprintf("%v", val)
		}
		return "6.5.0"
	}()
	header["versioncode"] = func() string {
		val, ok := cookie["versioncode"]
		if ok {
			return fmt.Sprintf("%v", val)
		}
		return "164"
	}()
	header["buildver"] = func() string {
		val, ok := cookie["buildver"]
		if ok {
			return fmt.Sprintf("%v", val)
		}
		return strconv.FormatInt(time.Now().Unix(), 10)[0:10]
	}()
	header["resolution"] = func() string {
		val, ok := cookie["resolution"]
		if ok {
			return fmt.Sprintf("%v", val)
		}
		return "1920x1080"
	}()
	header["os"] = func() string {
		val, ok := cookie["os"]
		if ok {
			return fmt.Sprintf("%v", val)
		}
		return "android"
	}()
	header["__csrf"] = csrfToken
	cookieMusicU, ok := cookie["MUSIC_U"]
	if ok {
		header["MUSIC_U"] = cookieMusicU
	}
	cookieMusicA, ok := cookie["MUSIC_A"]
	if ok {
		header["MUSIC_A"] = cookieMusicA
	}

	var cookies string
	for key, val := range header {
		cookies += encodeURIComponent(key) + "=" + encodeURIComponent(fmt.Sprintf("%v", val)) + "; "
	}
	req.Header.Set("Cookie", strings.TrimRight(cookies, "; "))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", ChooseUserAgent())

	if log.GetLevel() == log.DebugLevel {
		log.Debugf("[EapiReq]: %+v", req)
		log.Debugf("[EapiReqBodyJson]: %s", AesDecryptECB([]byte(data)))
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorln(err)
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if log.GetLevel() == log.DebugLevel {
		log.Debugf("[EapiResp]: %+v", resp)
	}

	return string(body), nil
}
