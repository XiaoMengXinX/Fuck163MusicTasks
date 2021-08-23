package utils

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/url"
	"strings"
)

// GetLoginStat 查询登录状态
func GetLoginStat(data RequestData, config APIConfig) (result LoginStatData, err error) {
	body, err := MakeRequest(config.NeteaseAPI+"/login/status", data)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &result)
	return result, err
}

// GetTasksData 获取音乐人任务
func GetTasksData(data RequestData, config APIConfig) (result TasksData, err error) {
	body, err := MakeRequest(config.NeteaseAPI+"/musician/tasks", data)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &result)
	return result, err
}

// GetCloudbeanData 获取云豆数量
func GetCloudbeanData(data RequestData, config APIConfig) (result CloudbeanData, err error) {
	body, err := MakeRequest(config.NeteaseAPI+"/musician/cloudbean", data)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &result)
	return result, err
}

// MusicianSign 音乐人签到
func MusicianSign(data RequestData, config APIConfig) (result MusicianSignResult, err error) {
	body, err := MakeRequest(config.NeteaseAPI+"/musician/sign", data)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &result)
	return result, err
}

// UserSign 用户签到 , signType = 0 为安卓端签到 , 1 为 web/PC 签到
func UserSign(data RequestData, signType int, config APIConfig) (result UserSignResult, err error) {
	body, err := MakeRequest(fmt.Sprintf("%s/daily_signin?type=%d", config.NeteaseAPI, signType), data)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &result)
	return result, err
}

// SendMsg 发送私信
func SendMsg(data RequestData, userID []int, msg string, config APIConfig) (result SendMsgResult, err error) {
	var userIDs string
	userIDs = fmt.Sprintf("%d", userID[0])
	for i := 1; i < len(userID); i++ {
		userIDs = fmt.Sprintf("%s,%d", userIDs, userID[i])
	}
	body, err := MakeRequest(fmt.Sprintf("%s/send/text?user_ids=%s&msg=%s", config.NeteaseAPI, userIDs, url.QueryEscape(msg)), data)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &result)
	return result, err
}

// Comment 发送评论
func Comment(data RequestData, commentConfig CommentConfig, config APIConfig) (result CommentResult, err error) {
	var commentIdArg, threadIdArg, idArg string
	if commentConfig.CommentId != 0 {
		commentIdArg = fmt.Sprintf("&commentId=%d", commentConfig.CommentId)
	}
	if commentConfig.ID != 0 {
		idArg = fmt.Sprintf("&id=%d", commentConfig.ID)
	}
	if commentConfig.ThreadId != 0 {
		threadIdArg = fmt.Sprintf("&threadId=%d", commentConfig.ThreadId)
	}
	body, err := MakeRequest(fmt.Sprintf("%s/comment?t=%d&type=%d&content=%s%s%s%s", config.NeteaseAPI, commentConfig.CommentType, commentConfig.ResType, url.QueryEscape(commentConfig.Content), commentIdArg, threadIdArg, idArg), data)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &result)
	return result, err
}

// SendEvent 发布动态
func SendEvent(data RequestData, msg string) (result SendEventResult, err error) {
	var options eapiOption
	options.Path = "/api/share/friends/resource"
	options.Url = "https://music.163.com/eapi/share/friends/resource"
	UUID := uuid.New()
	shareConfig := EventConfig{
		Msg:        msg,
		Type:       "noresource",
		UUID:       strings.Replace(UUID.String(), "-", "", -1),
		Pics:       "[]",
		AddComment: "false",
		Header:     "{}",
		ER:         "true",
	}
	bodyJson, err := json.Marshal(shareConfig)
	if err != nil {
		return SendEventResult{}, err
	}
	options.Json = string(bodyJson)
	body, err := EapiRequest(options, data)
	err = json.Unmarshal([]byte(body), &result)
	return result, err
}

// DelEvent 删除动态
func DelEvent(data RequestData, eventID int) (result PlainResult, err error) {
	var options eapiOption
	options.Path = "/api/event/delete"
	options.Url = "https://music.163.com/eapi/event/delete"
	delEventConfig := DelEventConfig{
		ID:     fmt.Sprintf("%d", eventID),
		Header: "{}",
		ER:     "true",
	}
	bodyJson, err := json.Marshal(delEventConfig)
	if err != nil {
		return PlainResult{}, err
	}
	options.Json = string(bodyJson)
	body, err := EapiRequest(options, data)
	err = json.Unmarshal([]byte(body), &result)
	return result, err
}

// GetUserDetail 获取用户详细
func GetUserDetail(data RequestData, uid int, config APIConfig) (result UserDetailData, err error) {
	body, err := MakeRequest(fmt.Sprintf("%s/user/detail?uid=%d", config.NeteaseAPI, uid), data)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &result)
	return result, err
}

// ObtainCloudBean 领取云豆
func ObtainCloudBean(data RequestData, userMissionId, period int, config APIConfig) (result PlainResult, err error) {
	body, err := MakeRequest(fmt.Sprintf("%s/musician/cloudbean/obtain?id=%d&period=%d", config.NeteaseAPI, userMissionId, period), data)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &result)
	return result, err
}
