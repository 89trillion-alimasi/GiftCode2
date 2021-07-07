package model

import "time"

// GiftCode 表示一个已经被创建的礼品码
type GiftCode struct {
	// Description 礼品描述信息
	Description string `json:"description"`

	// Type 礼品码（指定用户一次性消耗:1，不指定用户限制兑换次数:2，不限用户不限次数兑换:3）
	Type int `json:"type"`

	// ReceivingUser 指定可领取的用户名
	ReceivingUser string `json:"receiving_user"`

	// AvailableTimes 礼品码可领取的次数
	AvailableTimes int `json:"available_times"`

	// ValidPeriod 有效期
	ValidPeriod int64 `json:"valid_period"`

	// CreatTime 礼品码被创建的时间
	CreatTime int64 `json:"create_time"`

	// CreateUser 创建这个礼品码的用户
	CreateUser string `json:"create_user"`

	// GiftPackages 存储礼品包内容
	GiftPackages []GiftPackage `json:"gift_packages"`

	// ReceivedUsers 存储已经领取过该礼品码的用户
	ReceivedUsers map[string]int64 `json:"received_users"`

	// 礼品码已经被领取过的次数
	ReceivedCount int `json:"received_count"`

	// code 存储内部生成的礼品码
	Code string `json:"code"`

	// 礼品码过期时间，内部 config 使用
	Expiration time.Duration `json:"-"`
}

// AddReceivedUser 将领取用户添加到礼品码的领取列表中，同时增加领取次数
func (g *GiftCode) AddReceivedUser(user string) {
	if g.ReceivedUsers == nil {
		g.ReceivedUsers = make(map[string]int64)
	}
	g.ReceivedUsers[user] = time.Now().Unix()
	g.ReceivedCount++
}

// GiftPackage 表示礼品包内容
type GiftPackage struct {
	Name string `json:"name"`
	Num  int    `json:"num"`
}

// VerifyRequest 是验证礼品码的请求实体
type VerifyRequest struct {
	// Code 礼品码
	Code string `json:"code"`
	// User 领取用户
	User string `json:"user"`
}

//用户登陆数据
type UerInfo struct {
	UID     string
	Gold    string
	Diamond string
}
