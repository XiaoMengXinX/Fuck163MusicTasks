package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/XiaoMengXinX/Music163Api-Go/api"
	"github.com/XiaoMengXinX/Music163Api-Go/types"
	"github.com/XiaoMengXinX/Music163Api-Go/utils"
	"github.com/robfig/cron/v3"
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

var config Config
var commentLag RandomNum
var eventLag RandomNum
var msgLag RandomNum
var mlogLag RandomNum
var processingUser int
var circleID string
var configFileName = flag.String("c", "config.json", "Config filename") // 从 cli 参数读取配置文件名
var printVersion = flag.Bool("v", false, "Print version")
var isDEBUG = flag.Bool("d", false, "DEBUG mode")

var (
	runtimeVersion = fmt.Sprintf(runtime.Version())                     // 编译环境
	version        = ""                                                 // 程序版本
	commitSHA      = ""                                                 // 编译哈希
	buildTime      = ""                                                 // 编译日期
	buildARCH      = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH) // 运行环境
)

func init() {
	checkPathExists("./log")
	timeStamp := time.Now().Local().Format("2006-01-02")
	logFile := fmt.Sprintf("./log/%v.log", timeStamp)
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Error(err)
	}
	output := io.MultiWriter(file, os.Stdout)
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
Build ARCH: %s
`, version, runtimeVersion, commitSHA, buildTime, buildARCH)
		os.Exit(0)
	}

	func() { // 读取配置文件
		configFile, err := os.Open(*configFileName)
		if err != nil {
			log.Errorln(err)
			log.Fatal("读取配置文件失败")
		}
		defer func(configFile *os.File) {
			err := configFile.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(configFile)
		configFileData, err := ioutil.ReadAll(configFile)
		if err != nil {
			log.Errorln(err)
			log.Fatal("读取配置文件失败")
		}
		err = json.Unmarshal(configFileData, &config)
		if err != nil {
			log.Errorln(err)
			log.Fatal("读取配置文件失败, 请检查你的 JSON 格式是否正确")
		}
	}()

	if config.DEBUG { // 检查是否开启 DEBUG 模式
		log.SetLevel(log.DebugLevel)
	}

	commentLag.Set(config.CommentConfig.LagConfig) // 设置延迟
	eventLag.Set(config.EventSendConfig.LagConfig)
	msgLag.Set(config.SendMsgConfig.LagConfig)
	mlogLag.Set(config.SendMlogConfig.LagConfig)

	startTasks()
	startCron()
}

func startCron() {
	if config.Cron.Enabled {
		location, err := time.LoadLocation("Asia/Hong_Kong")
		if err != nil {
			log.Fatal(err)
		}
		parser := cron.NewParser(
			cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
		)
		c := cron.New(cron.WithLocation(location), cron.WithParser(parser))
		var entryID cron.EntryID
		entryID, err = c.AddFunc(fmt.Sprintf("%s", config.Cron.Expression), func() {
			entry := c.Entry(entryID)
			log.Printf("[Cron] 任务已运行, 下次运行时间 %s", entry.Next)
			if config.Cron.EnableLag {
				lag := RandomNum{}
				config.Cron.LagConfig.RandomLag = true
				lag.Set(config.Cron.LagConfig)
				randomLag := lag.Get()
				if randomLag != 0 {
					log.Printf("[Cron] 随机延时 %d 秒", randomLag)
					time.Sleep(time.Duration(randomLag) * time.Second)
				}
			}
			startTasks()
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("[Cron] 任务创建成功, 表达式: %s", config.Cron.Expression)
		c.Start()
		entry := c.Entry(entryID)
		log.Printf("[Cron] 任务已启动, 下次运行时间 %s", entry.Next)
		select {}
	}
}

func startTasks() {
	for processingUser = 0; processingUser < len(config.Users); processingUser++ { // 开始执行自动任务
		data := utils.RequestData{
			Cookies: config.Users[processingUser].Cookies,
		}
		userData, err := api.GetLoginStatus(data)
		if err != nil {
			log.Errorln(err)
		}
		if userData.Profile.UserId == 0 {
			log.Errorf("获取 User[%d] 登录状态失败, 请检查 MUSIC_U 是否失效", processingUser)
		} else {
			err = autoTasks(userData, data)
			if err != nil {
				log.Errorln(err)
			}
		}
	}
}

func autoTasks(userData types.LoginStatusData, data utils.RequestData) error {
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
	userDetail, err := api.GetUserDetail(data, userData.Account.Id)
	if err != nil {
		return err
	}
	if strings.Contains(userDetail.CurrentExpert.RoleName, "网易音乐人") {
		artistDetail, err := api.GetArtistDetail(data, int64(userDetail.Profile.ArtistId))
		parseCircleID(artistDetail)
		autoTasks, err := checkCloudBean(userData, data)
		if err != nil {
			return err
		}
		if len(autoTasks) != 0 {
			log.Printf("[%s] 正在完成音乐人任务中", userData.Profile.Nickname)
			for i := 0; i < len(autoTasks); i++ {
				musicianTasks(userData, data, autoTasks, i)
			}
			log.Printf("[%s] 音乐人任务执行完成, 正在重新检查并领取云豆", userData.Profile.Nickname)
			time.Sleep(time.Duration(10) * time.Second)
			_, err = checkCloudBean(userData, data)
			if err != nil {
				return err
			}
		}
	}
	if config.AutoGetVipGrowthpoint {
		err := vipGrowthpointTask(userData, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func musicianTasks(userData types.LoginStatusData, data utils.RequestData, autoTasks []int, i int) {
	defer func() {
		err := recover()
		if err != nil {
			log.Errorln(err)
		}
	}()
	switch autoTasks[i] {
	case 749006:
		log.Printf("[%s] 执行音乐人签到任务中", userData.Profile.Nickname)
		result, err := api.MusicianSign(data)
		if err != nil {
			log.Println(err)
		}
		if result.Code == 200 {
			log.Printf("[%s] 音乐人签到成功", userData.Profile.Nickname)
		} else {
			log.Printf("[%s] 音乐人签到失败: %s", userData.Profile.Nickname, result.Message)
		}
	case 740004:
		log.Printf("[%s] 执行发送动态任务中", userData.Profile.Nickname)
		err := sendEventTask(userData, data)
		if err != nil {
			log.Println(err)
		}
		log.Printf("[%s] 发送动态任务执行完成", userData.Profile.Nickname)
	case 732004:
		log.Printf("[%s] 执行回复评论任务中", userData.Profile.Nickname)
		commentConfig := api.CommentConfig{
			ResType:      api.ResTypeMusic,
			ResID:        config.CommentConfig.RepliedComment[processingUser].MusicID,
			CommentID:    config.CommentConfig.RepliedComment[processingUser].CommentID,
			ForwardEvent: false,
		}
		err := replyCommentTask(userData, commentConfig, data)
		if err != nil {
			log.Println(err)
		}
		log.Printf("[%s] 发送回复评论执行完成", userData.Profile.Nickname)
	case 755001:
		log.Printf("[%s] 执行发送私信任务中", userData.Profile.Nickname)
		err := sendMsgTask(userData, config.SendMsgConfig.UserID[processingUser], data)
		if err != nil {
			log.Println(err)
		}
		log.Printf("[%s] 发送私信任务执行完成", userData.Profile.Nickname)
	case 744005:
		log.Printf("[%s] 执行发送 Mlog 任务中", userData.Profile.Nickname)
		err := sendMlogTask(userData, data)
		if err != nil {
			log.Println(err)
		}
		log.Printf("[%s] 发送 Mlog 任务执行完成", userData.Profile.Nickname)
	case 755000:
		log.Printf("[%s] 执行发送主创说任务中", userData.Profile.Nickname)
		commentConfig := api.CommentConfig{
			ResType:      api.ResTypeMusic,
			ResID:        config.CommentConfig.RepliedComment[processingUser].MusicID,
			ForwardEvent: false,
		}
		err := musicianSaidTask(userData, commentConfig, data)
		if err != nil {
			log.Println(err)
		}
		log.Printf("[%s] 发送主创说任务执行完成", userData.Profile.Nickname)
	case 740005:
		log.Printf("[%s] 执行访问云圈任务中", userData.Profile.Nickname)
		err := getCircleTask(data)
		if err != nil {
			log.Println(err)
		}
		log.Printf("[%s] 访问云圈任务执行完成", userData.Profile.Nickname)
	}
}

func getCircleTask(data utils.RequestData) error {
	if circleID != "" {
		result, err := api.GetCircle(data, circleID)
		if err != nil {
			return err
		}
		if result.Code != 200 {
			return fmt.Errorf("%s", result.Message)
		}
	}
	return nil
}

func userSignTask(userData types.LoginStatusData, data utils.RequestData) error {
	result, err := api.UserSign(data, 0)
	if err != nil {
		return err
	}
	if result.Code != 200 {
		log.Printf("[%s] %s (%s)", userData.Profile.Nickname, result.Msg, "Android")
	} else {
		log.Printf("[%s] 签到成功 (%s)", userData.Profile.Nickname, "Android")
	}

	result, err = api.UserSign(data, 1)
	if err != nil {
		return err
	}
	if result.Code != 200 {
		log.Printf("[%s] %s (%s)", userData.Profile.Nickname, result.Msg, "web/PC")
	} else {
		log.Printf("[%s] 签到成功 (%s)", userData.Profile.Nickname, "Android")
	}
	return nil
}

func sendEventTask(userData types.LoginStatusData, data utils.RequestData) error {
	failedTimes := 0
	for i := 0; i < 1; {
		if failedTimes >= 5 {
			return fmt.Errorf("[%s] 发送动态累计 %d 次失败, 已自动退出", userData.Profile.Nickname, failedTimes)
		}
		msg := randomText(config.Content)
		sendResult, err := api.SendEvent(data, msg, []string{})
		if err != nil {
			return err
		}
		if sendResult.Code == 200 {
			log.Printf("[%s] 发送动态成功, 动态ID: %d, 内容: \"%s\"", userData.Profile.Nickname, sendResult.Event.Id, msg)
			i++
			if config.EventSendConfig.LagConfig.LagBetweenSendAndDelete {
				randomLag := eventLag.Get()
				if randomLag != 0 {
					log.Printf("[%s] 延时 %d 秒", userData.Profile.Nickname, randomLag)
					time.Sleep(time.Duration(randomLag) * time.Second)
				}
			}
			delResult, err := api.DelEvent(data, sendResult.Event.Id)
			if err != nil {
				return err
			}
			if delResult.Code != 200 {
				log.Errorf("[%s] 删除动态失败, 动态ID: %d, 代码: %d, 原因: \"%s\"", userData.Profile.Nickname, sendResult.Event.Id, delResult.Code, delResult.Message)
			} else {
				log.Printf("[%s] 删除动态成功, 动态ID: %d", userData.Profile.Nickname, sendResult.Event.Id)
			}
		} else {
			log.Errorf("[%s] 发送动态失败, 内容: \"%s\", 代码: %d, 原因: \"%s\"", userData.Profile.Nickname, msg, sendResult.Code, sendResult.Message)
			failedTimes++
		}
		randomLag := eventLag.Get()
		if randomLag != 0 {
			log.Printf("[%s] 延时 %d 秒", userData.Profile.Nickname, randomLag)
			time.Sleep(time.Duration(randomLag) * time.Second)
		}
	}
	return nil
}

func replyCommentTask(userData types.LoginStatusData, commentConfig api.CommentConfig, data utils.RequestData) error {
	replyToID := commentConfig.CommentID
	failedTimes := 0
	for i := 0; i < 2; {
		if failedTimes >= 5 {
			return fmt.Errorf("[%s] 回复评论累计 %d 次失败, 已自动退出", userData.Profile.Nickname, failedTimes)
		}
		msg := randomText(config.Content)
		commentConfig.CommentID = replyToID
		commentConfig.Content = msg
		replyResult, err := api.ReplyComment(data, commentConfig)
		if err != nil {
			return err
		}
		if replyResult.Code == 200 {
			log.Printf("[%s] 回复评论成功, 歌曲ID: %d, 评论ID: %d, 内容: \"%s\"", userData.Profile.Nickname, commentConfig.ResID, commentConfig.CommentID, msg)
			i++
			if config.CommentConfig.LagConfig.LagBetweenSendAndDelete {
				randomLag := commentLag.Get()
				if randomLag != 0 {
					log.Printf("[%s] 延时 %d 秒", userData.Profile.Nickname, randomLag)
					time.Sleep(time.Duration(randomLag) * time.Second)
				}
			}
			commentConfig.CommentID = replyResult.Comment.CommentId
			commentConfig.ResType = api.ResTypeMusic
			commentConfig.Content = ""
			delResult, err := api.DelComment(data, commentConfig)
			if err != nil {
				return err
			}
			if delResult.Code != 200 {
				log.Errorf("[%s] 删除评论失败, 歌曲ID: %d, 评论ID: %d, 代码: %d", userData.Profile.Nickname, commentConfig.ResID, commentConfig.CommentID, delResult.Code)
			} else {
				log.Printf("[%s] 删除评论成功, 歌曲ID: %d, 评论ID: %d", userData.Profile.Nickname, commentConfig.ResID, commentConfig.CommentID)
			}
		} else {
			log.Errorf("[%s] 回复评论失败, 歌曲ID: %d, 评论ID: %d, 内容: \"%s\", 代码: %d", userData.Profile.Nickname, commentConfig.ResID, commentConfig.CommentID, msg, replyResult.Code)
			failedTimes++
		}
		randomLag := commentLag.Get()
		if randomLag != 0 {
			log.Printf("[%s] 延时 %d 秒", userData.Profile.Nickname, randomLag)
			time.Sleep(time.Duration(randomLag) * time.Second)
		}
	}
	return nil
}

func sendMsgTask(userData types.LoginStatusData, userIDs []int, data utils.RequestData) error {
	failedTimes := 0
	for i := 0; i < 2; {
		if failedTimes >= 5 {
			return fmt.Errorf("[%s] 发送私信累计 %d 次失败, 已自动退出, 是不是工具人把你拉黑了(", userData.Profile.Nickname, failedTimes)
		}
		var userID int
		if len(userIDs) == 1 {
			userID = userIDs[0]
		} else {
			rand.Seed(time.Now().UnixNano())
			userID = userIDs[rand.Intn(len(userIDs)-1)]
		}
		msg := randomText(config.Content)
		sendResult, err := api.SendTextMsg(data, []int{userID}, msg)
		if err != nil {
			return err
		}
		if sendResult.Code == 200 {
			log.Printf("[%s] 发送私信成功, 用户ID: %d, 内容: \"%s\"", userData.Profile.Nickname, userID, msg)
			i++
		} else {
			if len(sendResult.Blacklist) != 0 {
				log.Errorf("[%s] 发送私信失败, 用户ID: %d, 内容: \"%s\", 代码: %d, 您已被目标用户拉黑", userData.Profile.Nickname, userID, msg, sendResult.Code)
			} else {
				log.Errorf("[%s] 发送私信失败, 用户ID: %d, 内容: \"%s\", 代码: %d", userData.Profile.Nickname, userID, msg, sendResult.Code)
			}
			failedTimes++
		}
		randomLag := msgLag.Get()
		if randomLag != 0 {
			log.Printf("[%s] 延时 %d 秒", userData.Profile.Nickname, randomLag)
			time.Sleep(time.Duration(randomLag) * time.Second)
		}
	}
	return nil
}

func sendMlogTask(userData types.LoginStatusData, data utils.RequestData) error {
	if !checkPathExists(config.SendMlogConfig.PicFolder) {
		return fmt.Errorf("[%s] \"%s\" 图片文件夹不存在, 无法发送 Mlog", userData.Profile.Nickname, config.SendMlogConfig.PicFolder)
	}
	files, err := os.ReadDir(config.SendMlogConfig.PicFolder)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("[%s] \"%s\" 图片文件夹为空, 无法发送 Mlog", userData.Profile.Nickname, config.SendMlogConfig.PicFolder)
	}
	rand.Seed(time.Now().UnixNano())
	fileName := files[rand.Intn(len(files))].Name()
	musicID := config.SendMlogConfig.MusicIDs[rand.Intn(len(config.SendMlogConfig.MusicIDs))]
	text := randomText(config.Content)
	mlogData, err := api.SendPicMlog(data, text, musicID, []string{fmt.Sprintf("%s/%s", config.SendMlogConfig.PicFolder, fileName)})
	if err != nil {
		return err
	}
	if mlogData.Code != 200 {
		log.Errorf("[%s] 发送 Mlog 失败, 代码: %d, 原因: \"%s\"", userData.Profile.Nickname, mlogData.Code, mlogData.Message)
	} else {
		log.Printf("[%s] 发送 Mlog 成功, 动态ID: %d, 内容: \"%s\", 图片: \"%s\"", userData.Profile.Nickname, mlogData.Data.Event.Id, text, fmt.Sprintf("%s/%s", config.SendMlogConfig.PicFolder, fileName))
	}
	randomLag := mlogLag.Get()
	if randomLag != 0 {
		log.Printf("[%s] 延时 %d 秒", userData.Profile.Nickname, randomLag)
		time.Sleep(time.Duration(randomLag) * time.Second)
	}
	result, err := api.DelEvent(data, mlogData.Data.Event.Id)
	if err != nil {
		return err
	}
	if result.Code != 200 {
		log.Errorf("[%s] 删除 Mlog 失败, 动态ID: %d, 代码: %d, 原因: \"%s\"", userData.Profile.Nickname, mlogData.Data.Event.Id, result.Code, result.Message)
	} else {
		log.Printf("[%s] 删除 Mlog 成功, 动态ID: %d", userData.Profile.Nickname, mlogData.Data.Event.Id)
	}
	return nil
}

func musicianSaidTask(userData types.LoginStatusData, commentConfig api.CommentConfig, data utils.RequestData) error {
	msg := randomText(config.Content)
	commentConfig.Content = msg
	replyResult, err := api.AddComment(data, commentConfig)
	if err != nil {
		return err
	}
	if replyResult.Code == 200 {
		log.Printf("[%s] 发送评论成功, 歌曲ID: %d, 评论ID: %d, 内容: \"%s\"", userData.Profile.Nickname, commentConfig.ResID, commentConfig.CommentID, msg)
		if config.CommentConfig.LagConfig.LagBetweenSendAndDelete {
			randomLag := commentLag.Get()
			if randomLag != 0 {
				log.Printf("[%s] 延时 %d 秒", userData.Profile.Nickname, randomLag)
				time.Sleep(time.Duration(randomLag) * time.Second)
			}
		}
		commentConfig.CommentID = replyResult.Comment.CommentId
		commentConfig.ResType = api.ResTypeMusic
		commentConfig.Content = ""
		delResult, err := api.DelComment(data, commentConfig)
		if err != nil {
			return err
		}
		if delResult.Code != 200 {
			log.Errorf("[%s] 删除评论失败, 歌曲ID: %d, 评论ID: %d, 代码: %d", userData.Profile.Nickname, commentConfig.ResID, commentConfig.CommentID, delResult.Code)
		} else {
			log.Printf("[%s] 删除评论成功, 歌曲ID: %d, 评论ID: %d", userData.Profile.Nickname, commentConfig.ResID, commentConfig.CommentID)
		}
	} else {
		log.Errorf("[%s] 发送评论失败, 歌曲ID: %d, 评论ID: %d, 内容: \"%s\", 代码: %d", userData.Profile.Nickname, commentConfig.ResID, commentConfig.CommentID, msg, replyResult.Code)
	}
	return nil
}

func vipGrowthpointTask(userData types.LoginStatusData, data utils.RequestData) error {
	log.Printf("[%s] 正在检查会员状态", userData.Profile.Nickname)
	vipStat, err := api.GetVipInfo(data)
	if err != nil {
		return err
	}
	if vipStat.Data.RedVipLevel == 0 {
		log.Printf("[%s] 无会员权限，跳过领取成长值任务", userData.Profile.Nickname)
		return nil
	}
	log.Printf("[%s] 检查成功，正在领取会员任务成长值", userData.Profile.Nickname)
	_, err = api.VipTaskRewardAll(data)
	return err
}

func checkCloudBean(userData types.LoginStatusData, data utils.RequestData) ([]int, error) {
	cloudBeanData, err := api.GetCloudbeanNum(data)
	if err != nil {
		return []int{}, err
	}
	log.Printf("[%s] 账号当前云豆数: %d", userData.Profile.Nickname, cloudBeanData.Data.CloudBean)
	log.Printf("[%s] 获取音乐人任务中...", userData.Profile.Nickname)
	dailyTasks, err := api.GetMusicianDailyTasks(data)
	if err != nil {
		return []int{}, err
	}
	weeklyTasks, err := api.GetMusicianWeeklyTasks(data)
	if err != nil {
		return []int{}, err
	}
	var isObtainCloudBean bool
	var autoTasks []int
	for _, task := range dailyTasks.Data.List {
		if task.Status == 20 {
			log.Printf("[%s] 「%s」任务已完成, 正在领取云豆", userData.Profile.Nickname, task.Description)
			isObtainCloudBean = true
			result, err := api.ObtainCloudbean(data, task.UserMissionId, task.Period)
			if err != nil {
				log.Errorln(err)
			}
			if result.Code == 200 {
				log.Printf("[%s] 领取「%s」任务云豆成功, 云豆+%s", userData.Profile.Nickname, task.Description, task.RewardWorth)
			} else {
				log.Errorf("[%s] 领取「%s」任务云豆失败: %s", userData.Profile.Nickname, task.Description, result.Message)
			}
		} else if autoTaskAvail(task.MissionId) && task.Status != 100 {
			log.Printf("[%s] 任务「%s」任务未完成或进行中", userData.Profile.Nickname, task.Description)
			autoTasks = append(autoTasks, task.MissionId)
		}
	}
	for _, task := range weeklyTasks.Data.List {
		for _, s := range task.UserStageTargetList {
			if s.Status == 20 {
				log.Printf("[%s] 「%s」任务已完成, 正在领取云豆", userData.Profile.Nickname, task.Description)
				isObtainCloudBean = true
				if s.UserMissionId != 0 {
					result, err := api.ObtainCloudbean(data, int(s.UserMissionId), task.Period)
					if err != nil {
						log.Errorln(err)
					}
					if result.Code == 200 {
						log.Printf("[%s] 领取「%s」任务云豆成功, 云豆+%d", userData.Profile.Nickname, task.Description, s.Worth)
					} else {
						log.Errorf("[%s] 领取「%s」任务云豆失败: %s", userData.Profile.Nickname, task.Description, result.Message)
					}
				}
			} else if autoTaskAvail(task.MissionId) && task.Status != 100 {
				log.Printf("[%s] 任务「%s」任务未完成或进行中", userData.Profile.Nickname, task.Description)
				autoTasks = append(autoTasks, task.MissionId)
			}
		}
	}
	if isObtainCloudBean {
		time.Sleep(time.Duration(10) * time.Second)
		cloudBeanData, err = api.GetCloudbeanNum(data)
		if err != nil {
			return []int{}, err
		}
		log.Printf("[%s] 账号当前云豆数: %d", userData.Profile.Nickname, cloudBeanData.Data.CloudBean)
	}
	if len(autoTasks) == 0 {
		log.Printf("[%s] 后面的任务, 明天再来探索吧！", userData.Profile.Nickname)
	}
	return autoTasks, err
}

func randomText(textSlice []string) string {
	rand.Seed(time.Now().UnixNano())
	return textSlice[rand.Intn(len(textSlice)-1)]
}

func checkPathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Errorln(err)
		}
		return false
	}
	log.Errorln(err)
	return false
}

func parseCircleID(artistDetail types.ArtistDetailData) {
	for _, d := range artistDetail.Data.Blocks {
		if d.Code == "PERSONAL_MY_CIRCLE" {
			for _, creative := range d.Creatives {
				for _, r := range creative.Resources {
					if r.ResourceType != nil && r.ResourceId != nil {
						if *r.ResourceType == "CIRCLE" && *r.ResourceId != "" {
							circleID = *r.ResourceId
							return
						}
					}
				}
			}
		}
	}
}

func autoTaskAvail(val int) bool { //
	availAutoTasks := []int{740004, 744005, 732004, 755001, 749006, 755000, 740005}
	for i := 0; i < len(availAutoTasks); i++ {
		if val == availAutoTasks[i] {
			return true
		}
	}
	return false
}
