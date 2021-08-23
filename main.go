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

// LogFormatter 自定义 log 格式
type LogFormatter struct{}

// Format 自定义 log 格式
func (s *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006/01/02 15:04:05")
	var msg string
	msg = fmt.Sprintf("%s [%s] %s (%s:%d)\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Message, path.Base(entry.Caller.File), entry.Caller.Line)
	return []byte(msg), nil
}

var config utils.Config
var apiConfig utils.APIConfig
var commentLeg utils.RandomNum
var eventLeg utils.RandomNum
var msgLeg utils.RandomNum
var processingUser int
var configFileName = flag.String("c", "config.json", "Config filename") // 从 cli 参数读取配置文件名
var printVersion = flag.Bool("v", false, "Print version")
var isDEBUG = flag.Bool("d", false, "DEBUG mode")

var (
	runtimeVersion = fmt.Sprintf(runtime.Version())                     // 编译环境
	version        = ""                                                 // 程序版本
	commitSHA      = ""                                                 // 编译哈希
	buildTime      = ""                                                 // 编译日期
	buildOS        = ""                                                 // 编译系统
	buildARCH      = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH) // 运行环境
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
	log.SetReportCaller(true)
	flag.Parse() // 解析命令行参数
	if *isDEBUG {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	defer func() {
		err := recover()
		if err != nil {
			log.Errorln(err)
		}
	}()

	if *printVersion {
		fmt.Printf(`Fuck163MusicTasks %s (%s)
Build Hash: %s
Build Date: %s
Build OS: %s
Build ARCH: %s
`, version, runtimeVersion, commitSHA, buildTime, buildOS, buildARCH)
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
		err = json.Unmarshal(configFileData, &config)
		if err != nil {
			log.Fatal(err)
		}
	}()

	if config.DEBUG { // 检查是否开启 DEBUG 模式
		log.SetLevel(log.DebugLevel)
	}

	apiConfig.NeteaseAPI = config.NeteaseAPI // 设置自定义网易云 api

	commentLeg.Set(config.CommentReplyConfig.LagConfig) // 设置延迟
	eventLeg.Set(config.EventSendConfig.LagConfig)
	msgLeg.Set(config.SendMsgConfig.LagConfig)

	for processingUser = 0; processingUser < len(config.Users); processingUser++ { // 开始执行自动任务
		data := utils.RequestData{
			Cookies: config.Users[processingUser].Cookies,
		}

		userData, err := utils.GetLoginStat(data, apiConfig)
		if err != nil {
			log.Errorln(err)
		}

		err = autoTasks(userData, data)
		if err != nil {
			log.Errorln(err)
		}
	}
}

func autoTasks(userData utils.LoginStatData, data utils.RequestData) error {
	defer func() {
		err := recover()
		if err != nil {
			log.Errorln(err)
		}
	}()
	err := userSignTask(userData, data)
	if err != nil {
		log.Errorln(err)
	}
	userDetail, err := utils.GetUserDetail(data, userData.Data.Account.Id, apiConfig)
	if err != nil {
		return err
	}
	if strings.Contains(userDetail.CurrentExpert.RoleName, "音乐人") {
		autoTasks, err := checkCloudBean(userData, data)
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
						result, err := utils.MusicianSign(data, apiConfig)
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
						err := sendEventTask(userData, data)
						if err != nil {
							log.Println(err)
						}
						log.Printf("[%s] 发送动态任务执行完成", userData.Data.Profile.Nickname)
					case 393001:
						log.Printf("[%s] 执行回复评论任务中", userData.Data.Profile.Nickname)
						commentConfig := utils.CommentConfig{
							CommentType: 2,
							ResType:     0,
							ID:          config.CommentReplyConfig.RepliedComment[processingUser].ID,
							CommentId:   config.CommentReplyConfig.RepliedComment[processingUser].CommentId,
						}
						err := replyCommentTask(userData, commentConfig, data)
						if err != nil {
							log.Println(err)
						}
						log.Printf("[%s] 发送回复评论执行完成", userData.Data.Profile.Nickname)
					case 395002:
						log.Printf("[%s] 执行发送私信任务中", userData.Data.Profile.Nickname)
						err := sendMsgTask(userData, config.SendMsgConfig.UserID[processingUser], data)
						if err != nil {
							log.Println(err)
						}
						log.Printf("[%s] 发送私信任务执行完成", userData.Data.Profile.Nickname)
					}
				}()
			}
			log.Printf("[%s] 所有任务执行完成，正在重新检查并领取云豆", userData.Data.Profile.Nickname)
			autoTasks, err = checkCloudBean(userData, data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func userSignTask(userData utils.LoginStatData, data utils.RequestData) error {
	result, err := utils.UserSign(data, 0, apiConfig)
	if err != nil {
		return err
	}
	if result.Code != 200 {
		log.Printf("[%s] %s (%s)", userData.Data.Profile.Nickname, result.Msg, "Android")
	} else {
		log.Printf("[%s] 签到成功，获得 %d 经验 (%s)", userData.Data.Profile.Nickname, result.Point, "Android")
	}

	result, err = utils.UserSign(data, 1, apiConfig)
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

func sendEventTask(userData utils.LoginStatData, data utils.RequestData) error {
	failedTimes := 0
	for i := 0; i < 3; {
		if failedTimes >= 5 {
			return fmt.Errorf("[%s] 发送动态累计 %d 次失败, 已自动退出", userData.Data.Profile.Nickname, failedTimes)
		}
		msg := randomText(config.Content)
		sendResult, err := utils.SendEvent(data, msg)
		if err != nil {
			return err
		}
		if sendResult.Code == 200 {
			log.Printf("[%s] 发送动态成功, 内容: \"%s\"", userData.Data.Profile.Nickname, msg)
			i++
			if config.EventSendConfig.LagConfig.LagBetweenSendAndDelete {
				randomLag := eventLeg.Get()
				log.Printf("[%s] 延时 %d 秒", userData.Data.Profile.Nickname, randomLag)
				time.Sleep(time.Duration(randomLag) * time.Second)
			}
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
		randomLag := eventLeg.Get()
		log.Printf("[%s] 延时 %d 秒", userData.Data.Profile.Nickname, randomLag)
		time.Sleep(time.Duration(randomLag) * time.Second)
	}
	return nil
}

func replyCommentTask(userData utils.LoginStatData, commentConfig utils.CommentConfig, data utils.RequestData) error {
	replyToID := commentConfig.CommentId
	failedTimes := 0
	for i := 0; i < 5; {
		if failedTimes >= 5 {
			return fmt.Errorf("[%s] 回复评论累计 %d 次失败, 已自动退出", userData.Data.Profile.Nickname, failedTimes)
		}
		msg := randomText(config.Content)
		commentConfig.CommentId = replyToID
		commentConfig.CommentType = 2
		commentConfig.Content = msg
		replyResult, err := utils.Comment(data, commentConfig, apiConfig)
		if err != nil {
			return err
		}
		if replyResult.Code == 200 {
			log.Printf("[%s] 回复评论成功, 歌曲ID: %d, 评论ID: %d, 内容: \"%s\"", userData.Data.Profile.Nickname, commentConfig.ID, commentConfig.CommentId, msg)
			i++
			if config.CommentReplyConfig.LagConfig.LagBetweenSendAndDelete {
				randomLag := commentLeg.Get()
				log.Printf("[%s] 延时 %d 秒", userData.Data.Profile.Nickname, randomLag)
				time.Sleep(time.Duration(randomLag) * time.Second)
			}
			commentConfig.CommentId = replyResult.Comment.CommentId
			commentConfig.CommentType = 0
			commentConfig.Content = ""
			delResult, err := utils.Comment(data, commentConfig, apiConfig)
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
		randomLag := commentLeg.Get()
		log.Printf("[%s] 延时 %d 秒", userData.Data.Profile.Nickname, randomLag)
		time.Sleep(time.Duration(randomLag) * time.Second)
	}
	return nil
}

func sendMsgTask(userData utils.LoginStatData, userIDs []int, data utils.RequestData) error {
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
		msg := randomText(config.Content)
		sendResult, err := utils.SendMsg(data, []int{userID}, msg, apiConfig)
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
		randomLag := msgLeg.Get()
		log.Printf("[%s] 延时 %d 秒", userData.Data.Profile.Nickname, randomLag)
		time.Sleep(time.Duration(randomLag) * time.Second)
	}
	return nil
}

func checkCloudBean(userData utils.LoginStatData, data utils.RequestData) ([]int, error) {
	cloudBeanData, err := utils.GetCloudbeanData(data, apiConfig)
	if err != nil {
		return []int{}, err
	}
	log.Printf("[%s] 账号当前云豆数: %d", userData.Data.Profile.Nickname, cloudBeanData.Data.CloudBean)
	log.Printf("[%s] 获取音乐人任务中...", userData.Data.Profile.Nickname)
	tasksData, err := utils.GetTasksData(data, apiConfig)
	if err != nil {
		return []int{}, err
	}
	var autoTasks []int
	for i := 0; i < len(tasksData.Data.List); i++ {
		if tasksData.Data.List[i].Status == 20 {
			log.Printf("[%s] 「%s」任务已完成，正在领取云豆", userData.Data.Profile.Nickname, tasksData.Data.List[i].Description)
			result, err := utils.ObtainCloudBean(data, tasksData.Data.List[i].UserMissionId, tasksData.Data.List[i].Period, apiConfig)
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
