package main

import (
	"flag"
	"fmt"
	"github.com/XiaoMengXinX/Music163Api-Go/api"
	"github.com/XiaoMengXinX/Music163Api-Go/utils"
	log "github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

// LogFormatter 自定义 log 格式
type LogFormatter struct{}

// Format 自定义 log 格式
func (s *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006/01/02 15:04:05")
	var msg string
	msg = fmt.Sprintf("%s [%s] %s (%s:%d)\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Message, path.Base(entry.Caller.File), entry.Caller.Line)
	return []byte(msg), nil
}

func init() {
	flag.Parse() // 解析命令行参数
	output := io.MultiWriter(os.Stdout)
	log.SetOutput(output)
	log.SetFormatter(&log.TextFormatter{
		DisableColors:          false,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})
	log.SetFormatter(new(LogFormatter))
	log.SetReportCaller(true)
	if *isDEBUG {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

var (
	isDEBUG = flag.Bool("d", false, "DEBUG mode")
)
var (
	rmPre = regexp.MustCompile(`(.*)MUSIC_U=`)
	rmSuf = regexp.MustCompile(`;(.*)`)
)

func main() {
	qrKey, err := api.GetQrUnikey(utils.RequestData{})
	if err != nil {
		log.Fatal(err)
	}
	qr, err := qrcode.New(fmt.Sprintf("https://music.163.com/login?codekey=%s", qrKey.Unikey), qrcode.High)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(qr.ToSmallString(false))
	fmt.Println("请使用网易云手机客户端扫描二维码")
	for {
		loginData, header, err := api.CheckQrLogin(utils.RequestData{}, qrKey.Unikey)
		if err != nil {
			log.Fatal(err)
		}
		if loginData.Code == 802 {
			fmt.Println(loginData.Message)
		}
		if loginData.Code == 803 {
			fmt.Println(loginData.Message)
			MUSIC_U := rmSuf.ReplaceAllString(rmPre.ReplaceAllString(header, ""), "")
			if MUSIC_U != "" {
				fmt.Printf("[MUSIC_U] %s\n", MUSIC_U)
			} else {
				log.Errorln("解析 MUSIC_U 失败，请重新登陆")
			}
			break
		}
		time.Sleep(time.Duration(1) * time.Second)
	}
}
