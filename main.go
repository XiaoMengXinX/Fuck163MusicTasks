package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/XiaoMengXinX/Fuck163MusicTasks/utils"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

type LogFormatter struct{}

func (s *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006/01/02 15:04:05")
	var msg string
	msg = fmt.Sprintf("%s [%s] %s (%s:%d)\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Message, path.Base(entry.Caller.File), entry.Caller.Line)
	return []byte(msg), nil
}

var Config utils.Config
var APIConfig utils.APIConfig
var CommentLeg utils.RandomNum
var EventLeg utils.RandomNum
var MsgLeg utils.RandomNum
var ProcessingUser int
var configFileName = flag.String("c", "config.json", "Config filename") // 从 cli 参数读取配置文件名
var printVersion = flag.Bool("v", false, "Print version")

var (
	RUNTIME_VERSION = fmt.Sprintf(runtime.Version())                     // 编译环境
	VERSION         = ""                                                 // 程序版本
	COMMIT_SHA      = ""                                                 // 编译哈希
	BUILD_TIME      = ""                                                 // 编译日期
	BUILD_OS        = ""                                                 // 编译系统
	BUILD_ARCH      = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH) // 运行环境
)

func init() {
	output := io.MultiWriter(os.Stdout)
	log.SetOutput(output)
	log.SetFormatter(&log.TextFormatter{
		DisableColors:          false,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})
	log.SetFormatter(new(LogFormatter))
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true)
}

func main() {
	defer func() {
		err := recover()
		if err != nil {
			log.Errorln(err)
		}
	}()

	flag.Parse() // 解析命令行参数

	if *printVersion {
		fmt.Printf(`Fuck163MusicTasks %s (%s)
Build Hash: %s
Build Date: %s
Build OS: %s
Build ARCH: %s
`, VERSION, RUNTIME_VERSION, COMMIT_SHA, BUILD_TIME, BUILD_OS, BUILD_ARCH)
		os.Exit(0)
	}

	func() { // 读取配置文件
		configFile, err := os.Open(*configFileName)
		if err != nil {
			log.Fatal(err)
		}
		defer func(configFile *os.File) {
			err := configFile.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(configFile)
		configFileData, err := ioutil.ReadAll(configFile)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(configFileData, &Config)
		if err != nil {
			log.Fatal(err)
		}
	}()

	if Config.DEBUG { // 检查是否开启 DEBUG 模式
		log.SetLevel(log.DebugLevel)
	}

	APIConfig.NeteaseAPI = Config.NeteaseAPI // 设置自定义网易云 api

	CommentLeg.Set(Config.CommentReplyConfig.LagConfig) // 设置延迟
	EventLeg.Set(Config.EventSendConfig.LagConfig)
	MsgLeg.Set(Config.SendMsgConfig.LagConfig)

	for ProcessingUser = 0; ProcessingUser < len(Config.Users); ProcessingUser++ { // 开始执行自动任务
		data := utils.RequestData{
			Cookies: Config.Users[ProcessingUser].Cookies,
		}

		userData, err := utils.GetLoginStat(data, APIConfig)
		if err != nil {
			log.Errorln(err)
		}

		err = AutoTasks(userData, data)
		if err != nil {
			log.Errorln(err)
		}
	}
}

func AutoTasks(userData utils.LoginStatData, data utils.RequestData) error {
	defer func() {
		err := recover()
		if err != nil {
			log.Errorln(err)
		}
	}()
	err := UserSignTask(userData, data)
	if err != nil {
		log.Errorln(err)
	}
	userDetail, err := utils.GetUserDetail(data, userData.Data.Account.Id, APIConfig)
	if err != nil {
		return err
	}
	if strings.Contains(userDetail.CurrentExpert.RoleName, "音乐人") {
		autoTasks, err := CheckCloudBean(userData, data)
		if err != nil {
			return err
		}
		if len(autoTasks) != 0 {
			log.Printf("[%s] 正在运行自动任务中", userData.Data.Profile.Nickname)
			for i := 0; i < len(autoTasks); i++ {
				func() {
					defer func() {
						err := recover()
						if err != nil {
							log.Errorln(err)
						}
					}()
					switch autoTasks[i] {
					case 399000:
						log.Printf("[%s] 执行音乐人签到任务中", userData.Data.Profile.Nickname)
						result, err := utils.MusicianSign(data, APIConfig)
						if err != nil {
							log.Println(err)
						}
						if result.Code == 200 {
							log.Printf("[%s] 音乐人签到成功", userData.Data.Profile.Nickname)
						} else {
							log.Printf("[%s] 音乐人签到失败: %s", userData.Data.Profile.Nickname, result.Message)
						}
					case 398000:
						log.Printf("[%s] 执行发送动态任务中", userData.Data.Profile.Nickname)
						err := SendEventTask(userData, data)
						if err != nil {
							log.Println(err)
						}
						log.Printf("[%s] 发送动态任务执行完成", userData.Data.Profile.Nickname)
					case 393001:
						log.Printf("[%s] 执行回复评论任务中", userData.Data.Profile.Nickname)
						commentConfig := utils.CommentConfig{
							CommentType: 2,
							ResType:     0,
							ID:          Config.CommentReplyConfig.RepliedComment[ProcessingUser].ID,
							CommentId:   Config.CommentReplyConfig.RepliedComment[ProcessingUser].CommentId,
						}
						err := ReplyCommentTask(userData, commentConfig, data)
						if err != nil {
							log.Println(err)
						}
						log.Printf("[%s] 发送回复评论执行完成", userData.Data.Profile.Nickname)
					case 395002:
						log.Printf("[%s] 执行发送私信任务中", userData.Data.Profile.Nickname)
						err := SendMsgTask(userData, Config.SendMsgConfig.UserID[ProcessingUser], data)
						if err != nil {
							log.Println(err)
						}
						log.Printf("[%s] 发送私信任务执行完成", userData.Data.Profile.Nickname)
					}
				}()
			}
			log.Printf("[%s] 所有任务执行完成，正在重新检查并领取云豆", userData.Data.Profile.Nickname)
			autoTasks, err = CheckCloudBean(userData, data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func UserSignTask(userData utils.LoginStatData, data utils.RequestData) error {
	result, err := utils.UserSign(data, 0, APIConfig)
	if err != nil {
		return err
	}
	if result.Code != 200 {
		log.Printf("[%s] %s (%s)", userData.Data.Profile.Nickname, result.Msg, "Android")
	} else {
		log.Printf("[%s] 签到成功，获得 %d 经验 (%s)", userData.Data.Profile.Nickname, result.Point, "Android")
	}

	result, err = utils.UserSign(data, 1, APIConfig)
	if err != nil {
		return err
	}
	if result.Code != 200 {
		log.Printf("[%s] %s (%s)", userData.Data.Profile.Nickname, result.Msg, "web/PC")
	} else {
		log.Printf("[%s] 签到成功，获得 %d 经验 (%s)", userData.Data.Profile.Nickname, result.Point, "Android")
	}
	return nil
}

func SendEventTask(userData utils.LoginStatData, data utils.RequestData) error {
	failedTimes := 0
	for i := 0; i < 3; {
		if failedTimes >= 5 {
			return fmt.Errorf("[%s] 发送动态累计 %d 次失败, 已自动退出", userData.Data.Profile.Nickname, failedTimes)
		}
		msg := randomText(Config.Content)
		sendResult, err := utils.SendEvent(data, msg)
		if err != nil {
			return err
		}
		if sendResult.Code == 200 {
			log.Printf("[%s] 发送动态成功, 内容: \"%s\"", userData.Data.Profile.Nickname, msg)
			i++
			randomLag := EventLeg.Get()
			log.Printf("[%s] 延时 %d 秒", userData.Data.Profile.Nickname, randomLag)
			time.Sleep(time.Duration(randomLag) * time.Second)
			delResult, err := utils.DelEvent(data, sendResult.Event.Id)
			if err != nil {
				return err
			}
			if delResult.Code != 200 {
				log.Errorf("[%s] 删除动态失败, 动态ID: %d, 代码: %d, 原因: \"%s\"", userData.Data.Profile.Nickname, sendResult.Event.Id, delResult.Code, delResult.Message)
			} else {
				log.Printf("[%s] 删除动态成功, 动态ID: %d", userData.Data.Profile.Nickname, sendResult.Event.Id)
			}
		} else {
			log.Errorf("[%s] 发送动态失败, 内容: \"%s\", 代码: %d, 原因: \"%s\"", userData.Data.Profile.Nickname, msg, sendResult.Code, sendResult.Message)
			failedTimes++
		}
		randomLag := EventLeg.Get()
		log.Printf("[%s] 延时 %d 秒", userData.Data.Profile.Nickname, randomLag)
		time.Sleep(time.Duration(randomLag) * time.Second)
	}
	return nil
}

func ReplyCommentTask(userData utils.LoginStatData, commentConfig utils.CommentConfig, data utils.RequestData) error {
	replyToID := commentConfig.CommentId
	failedTimes := 0
	for i := 0; i < 5; {
		if failedTimes >= 5 {
			return fmt.Errorf("[%s] 回复评论累计 %d 次失败, 已自动退出", userData.Data.Profile.Nickname, failedTimes)
		}
		msg := randomText(Config.Content)
		commentConfig.CommentId = replyToID
		commentConfig.CommentType = 2
		commentConfig.Content = msg
		replyResult, err := utils.Comment(data, commentConfig, APIConfig)
		if err != nil {
			return err
		}
		if replyResult.Code == 200 {
			log.Printf("[%s] 回复评论成功, 歌曲ID: %d, 评论ID: %d, 内容: \"%s\"", userData.Data.Profile.Nickname, commentConfig.ID, commentConfig.CommentId, msg)
			i++
			randomLag := CommentLeg.Get()
			log.Printf("[%s] 延时 %d 秒", userData.Data.Profile.Nickname, randomLag)
			time.Sleep(time.Duration(randomLag) * time.Second)
			commentConfig.CommentId = replyResult.Comment.CommentId
			commentConfig.CommentType = 0
			commentConfig.Content = ""
			delResult, err := utils.Comment(data, commentConfig, APIConfig)
			if err != nil {
				return err
			}
			if delResult.Code != 200 {
				log.Errorf("[%s] 删除评论失败, 歌曲ID: %d, 评论ID: %d, 代码: %d, 原因: \"%s\"", userData.Data.Profile.Nickname, commentConfig.ID, commentConfig.CommentId, delResult.Code, delResult.Message)
			} else {
				log.Printf("[%s] 删除评论成功, 歌曲ID: %d, 评论ID: %d", userData.Data.Profile.Nickname, commentConfig.ID, commentConfig.CommentId)
			}
		} else {
			log.Errorf("[%s] 回复评论失败, 歌曲ID: %d, 评论ID: %d, 内容: \"%s\", 代码: %d, 原因: \"%s\"", userData.Data.Profile.Nickname, commentConfig.ID, commentConfig.CommentId, msg, replyResult.Code, replyResult.Message)
			failedTimes++
		}
		randomLag := CommentLeg.Get()
		log.Printf("[%s] 延时 %d 秒", userData.Data.Profile.Nickname, randomLag)
		time.Sleep(time.Duration(randomLag) * time.Second)
	}
	return nil
}

func SendMsgTask(userData utils.LoginStatData, userIDs []int, data utils.RequestData) error {
	failedTimes := 0
	for i := 0; i < 5; {
		if failedTimes >= 5 {
			return fmt.Errorf("[%s] 发送私信累计 %d 次失败, 已自动退出, 是不是工具人把你拉黑了(", userData.Data.Profile.Nickname, failedTimes)
		}
		var userID int
		if len(userIDs) == 1 {
			userID = userIDs[0]
		} else {
			rand.Seed(time.Now().UnixNano())
			userID = userIDs[rand.Intn(len(userIDs)-1)]
		}
		msg := randomText(Config.Content)
		sendResult, err := utils.SendMsg(data, []int{userID}, msg, APIConfig)
		if err != nil {
			return err
		}
		if sendResult.Code == 200 {
			log.Printf("[%s] 发送私信成功, 用户ID: %d, 内容: \"%s\"", userData.Data.Profile.Nickname, userID, msg)
			i++
		} else {
			if len(sendResult.Blacklist) != 0 {
				log.Errorf("[%s] 发送私信失败, 用户ID: %d, 内容: \"%s\", 代码: %d, 您已被目标用户拉黑", userData.Data.Profile.Nickname, userID, msg, sendResult.Code)
			} else {
				log.Errorf("[%s] 发送私信失败, 用户ID: %d, 内容: \"%s\", 代码: %d", userData.Data.Profile.Nickname, userID, msg, sendResult.Code)
			}
			failedTimes++
		}
		randomLag := MsgLeg.Get()
		log.Printf("[%s] 延时 %d 秒", userData.Data.Profile.Nickname, randomLag)
		time.Sleep(time.Duration(randomLag) * time.Second)
	}
	return nil
}

func CheckCloudBean(userData utils.LoginStatData, data utils.RequestData) ([]int, error) {
	cloudBeanData, err := utils.GetCloudbeanData(data, APIConfig)
	if err != nil {
		return []int{}, err
	}
	log.Printf("[%s] 账号当前云豆数: %d", userData.Data.Profile.Nickname, cloudBeanData.Data.CloudBean)
	log.Printf("[%s] 获取音乐人任务中...", userData.Data.Profile.Nickname)
	tasksData, err := utils.GetTasksData(data, APIConfig)
	if err != nil {
		return []int{}, err
	}
	var autoTasks []int
	for i := 0; i < len(tasksData.Data.List); i++ {
		if tasksData.Data.List[i].Status == 20 {
			log.Printf("[%s] 「%s」任务已完成，正在领取云豆", userData.Data.Profile.Nickname, tasksData.Data.List[i].Description)
			result, err := utils.ObtainCloudBean(data, tasksData.Data.List[i].UserMissionId, tasksData.Data.List[i].Period, APIConfig)
			if err != nil {
				log.Errorln(err)
			}
			if result.Code == 200 {
				log.Printf("[%s] 领取「%s」任务云豆成功", userData.Data.Profile.Nickname, tasksData.Data.List[i].Description)
			} else {
				log.Printf("[%s] 领取「%s」任务云豆失败: %s", userData.Data.Profile.Nickname, tasksData.Data.List[i].Description, result.Message)
			}
		}
		if autoTaskAvail(tasksData.Data.List[i].MissionId) && tasksData.Data.List[i].Status != 100 && tasksData.Data.List[i].Status != 20 {
			log.Printf("[%s] 任务「%s」任务未完成，已添加到任务列表", userData.Data.Profile.Nickname, tasksData.Data.List[i].Description)
			autoTasks = append(autoTasks, tasksData.Data.List[i].MissionId)
		}
	}
	if len(autoTasks) == 0 {
		log.Printf("[%s] 后面的任务，明天再来探索吧！", userData.Data.Profile.Nickname)
	}
	return autoTasks, err
}

func autoTaskAvail(val int) bool {
	availAutoTasks := []int{399000, 398000, 393001, 395002}
	for i := 0; i < len(availAutoTasks); i++ {
		if val == availAutoTasks[i] {
			return true
		}
	}
	return false
}

func randomText(textSlice []string) string {
	rand.Seed(time.Now().UnixNano())
	return textSlice[rand.Intn(len(textSlice)-1)]
}
