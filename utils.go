package main

import (
	"math/rand"
	"time"
)

// Get 获取随机数
func (r *RandomNum) Get() int {
	if r.IsRandom {
		rand.Seed(time.Now().UnixNano())
		return rand.Intn(r.MaxNum-r.MinNum) + r.MinNum
	}
	return r.DefaultNum
}

// Set 设置随机数
func (r *RandomNum) Set(config LagConfig) {
	r.IsRandom = config.RandomLag
	r.DefaultNum = config.DefaultLag
	r.MinNum = config.LagMin
	r.MaxNum = config.LagMax
}
