package main

import (
	"net/http"
)

// Config 配置文件结构
type Config struct {
	DEBUG bool `json:"DEBUG"`
	Users []struct {
		Cookies []*http.Cookie `json:"Cookies"`
	} `json:"Users"`
	EventSendConfig struct {
		LagConfig LagConfig `json:"LagConfig"`
	} `json:"EventSendConfig"`
	CommentConfig struct {
		RepliedComment []struct {
			MusicID   int `json:"MusicID"`
			CommentID int `json:"CommentID"`
		} `json:"RepliedComment"`
		LagConfig LagConfig `json:"LagConfig"`
	} `json:"CommentConfig"`
	SendMsgConfig struct {
		UserID    [][]int   `json:"UserID"`
		LagConfig LagConfig `json:"LagConfig"`
	}
	SendMlogConfig struct {
		PicFolder string    `json:"PicFolder"`
		MusicIDs  []int     `json:"MusicIDs"`
		LagConfig LagConfig `json:"LagConfig"`
	} `json:"SendMlogConfig"`
	AutoGetVipGrowthpoint bool     `json:"AutoGetVipGrowthpoint"`
	Content               []string `json:"Content"`
	Cron                  struct {
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
