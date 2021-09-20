package main

import util "github.com/XiaoMengXinX/Music163Api-Go/utils"

// Config 配置文件结构
type Config struct {
	NeteaseAPI string `json:"NeteaseAPI"`
	DEBUG      bool   `json:"DEBUG"`
	Users      []struct {
		Cookies util.Cookies `json:"Cookies"`
	} `json:"Users"`
	EventSendConfig struct {
		LagConfig LagConfig `json:"LagConfig"`
	} `json:"EventSendConfig"`
	CommentReplyConfig struct {
		RepliedComment []struct {
			ID        int `json:"ID"`
			CommentId int `json:"CommentId"`
		} `json:"RepliedComment"`
		LagConfig LagConfig `json:"LagConfig"`
	} `json:"CommentReplyConfig"`
	SendMsgConfig struct {
		UserID    [][]int   `json:"UserID"`
		LagConfig LagConfig `json:"LagConfig"`
	}
	Content []string `json:"Content"`
	Cron    struct {
		Enabled    bool      `json:"Enabled"`
		Expression string    `json:"Expression"`
		EnableLag  bool      `json:"EnableLag"`
		LagConfig  LagConfig `json:"LagConfig"`
	} `json:"Cron"`
}

// LagConfig 延迟设置
type LagConfig struct {
	LagBetweenSendAndDelete bool `json:"LagBetweenSendAndDelete"`
	RandomLag               bool `json:"RandomLag"`
	DefaultLag              int  `json:"DefaultLag"`
	LagMin                  int  `json:"LagMin"`
	LagMax                  int  `json:"LagMax"`
}

// RandomNum 随机数设置
type RandomNum struct {
	IsRandom   bool
	DefaultNum int
	MinNum     int
	MaxNum     int
}
