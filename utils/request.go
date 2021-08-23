package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// MakeRequest 创建 http 请求
func MakeRequest(url string, data RequestData) ([]byte, error) {
	client := &http.Client{}
	var req *http.Request

	urlArr := strings.Split(url, "?")
	if len(urlArr) == 2 {
		urlArr[1] = fmt.Sprintf("%s&timestamp=%d", urlArr[1], time.Now().Unix()) // 给请求加入时间戳，避免缓存机制
		// url = urlArr[0] + "?" + netUrl.PathEscape(urlArr[1])
		url = urlArr[0] + "?" + urlArr[1]
	} else {
		url = fmt.Sprintf("%s?timestamp=%d", url, time.Now().Unix())
	}

	if data.Body == "" {
		req, _ = http.NewRequest("POST", url, nil)
	} else {
		req, _ = http.NewRequest("POST", url, strings.NewReader(data.Body))
	}

	for i := 0; i < len(data.Cookies); i++ {
		req.AddCookie(&http.Cookie{Name: data.Cookies[i].Key, Value: data.Cookies[i].Value, HttpOnly: true})
	}

	for i := 0; i < len(data.Headers); i++ {
		req.Header.Add(data.Headers[i].Key, data.Headers[i].Value)
	}

	if log.GetLevel() == log.DebugLevel {
		log.Debugf("[HttpReq]: %+v", req)
		log.Debugf("[HttpRespBody]: %s", data.Body)
	}

	resp, err := client.Do(req)
	if err != nil {
		return []byte(""), err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorln(err)
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), nil
	}

	if log.GetLevel() == log.DebugLevel {
		log.Debugf("[HttpResp]: %+v", resp)
		log.Debugf("[HttpRespBody]: %s", string(body))
	}

	return body, nil
}
