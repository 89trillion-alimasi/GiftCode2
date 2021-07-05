package controller

import (
	"Giftcode-mongo-protobuf/config"
	"Giftcode-mongo-protobuf/model"
	"Giftcode-mongo-protobuf/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	MESSAGE = "message"
)

//创建礼品码
func CreateGiftCode(c *gin.Context) {
	var gift model.GiftCode
	// 绑定post参数
	if err := c.Bind(&gift); err != nil {
		c.JSON(400, gin.H{MESSAGE: err.Error()})
		return
	}

	// 检查礼品码类型
	switch gift.Type {
	case 1: // 1 - 指定用户一次性消耗
		// 一次性消耗强制设置领取次数为 1
		gift.AvailableTimes = 1
		if gift.ReceivingUser == "" {
			c.JSON(400, gin.H{MESSAGE: "请指定领取用户"})
			return
		}
	case 2: // 2 - 不指定用户限制兑换次数
		if gift.AvailableTimes <= 0 {
			c.JSON(400, gin.H{MESSAGE: "请指定可兑换次数"})
			return
		}
	case 3: // 3 - 不限用户不限次数兑换 无用户限制 无兑换次数限制这里不做处理
	default: // 非法填写的礼品码类型
		c.JSON(400, gin.H{MESSAGE: "礼品码类型不合法"})
		return
	}

	// 检查礼品码描述
	if gift.Description == "" {
		c.JSON(400, gin.H{MESSAGE: "请输入礼品码描述信息"})
		return
	}

	// 检查礼品码有效期是否为空
	if gift.ValidPeriod == "" {
		c.JSON(400, gin.H{MESSAGE: "请输入有效期"})
		return
	}

	// 检查礼品码有效期是否正确
	vp, err := time.ParseInLocation("2006-01-02 15:04:05", gift.ValidPeriod, time.Local)
	if err != nil {
		c.JSON(400, gin.H{MESSAGE: "输入的礼品码有效期格式不正确"})
		return
	}
	exp := vp.Sub(time.Now())
	if exp <= 0 {
		c.JSON(400, gin.H{MESSAGE: "输入的礼品码有效期小于当前时间"})
		return
	}
	gift.Expiration = exp

	// 检查礼品码包含的礼品包是否为空
	if gift.GiftPackages == nil || len(gift.GiftPackages) == 0 {
		c.JSON(400, gin.H{MESSAGE: "请输入礼包内容"})
		return
	}
	var n = 5
	for {
		// 生成 8 位礼品码
		gift.Code = service.RandStringBytesMask(8)
		cmd := config.RDB.Get(gift.Code)

		if cmd.Err() != nil {
			break
		} else {
			if n > 0 {
				continue
			} else {
				c.JSON(400, gin.H{MESSAGE: "5次创建礼品码失败，请重新创建"})
				return
			}
		}

	}
	logrus.Infof("生成新的礼品码: %s", gift.Code)

	// 将礼品码存储到 config
	if err := service.SaveGiftCode(&gift); err != nil {
		c.JSON(400, gin.H{MESSAGE: err.Error()})
		return
	}

	// 返回礼品码
	c.JSON(200, gin.H{"code": gift.Code})
}

//查询礼品码
func QueryGiftCode(c *gin.Context) {
	// 检测是否输入了礼品码
	code := c.Query("code")
	if code == "" {
		c.JSON(400, gin.H{MESSAGE: "请输入礼品码"})
		return
	}

	// 尝试从 config 中查找该礼品码
	giftCode := service.QueryGiftCode(code)
	if giftCode == nil {
		c.JSON(400, gin.H{MESSAGE: "礼品码输入错误或礼品码已过期"})
		return
	}

	// 返回礼品码
	c.JSON(200, giftCode)
}

//登陆
func Login(c *gin.Context) {
	user_id := c.Query("userid")
	if len(user_id) == 0 {
		c.String(400, "请输入用户id")
		return
	}

	user, err := service.Login(user_id)

	if err == "用户存在成功返回" {
		c.JSON(200, gin.H{"Status": err, "User": user})
		return
	} else if err == "用户不存在，已为用户主动注册" {
		c.JSON(201, gin.H{"Status": err, "NewUser": user})
		return
	} else if err == "插入不成功" {
		c.String(507, err)
	}

}

//验证礼品码
func VerifyGiftCode(c *gin.Context) {
	var req model.VerifyRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(400, gin.H{MESSAGE: err.Error()})
		return
	}

	// 检查礼品码
	if req.Code == "" {
		c.JSON(400, gin.H{MESSAGE: "请输入礼品码"})
		return
	}

	// 检查领取用户
	if req.User == "" {
		c.JSON(400, gin.H{MESSAGE: "请输入领取用户"})
		return
	}

	gifCode, err := service.VerifyGiftCode(req)
	if err != nil {
		c.JSON(400, gin.H{MESSAGE: err.Error()})
		return
	}
	c.JSON(200, gin.H{"GiftPackages": gifCode.GiftPackages})
}

func ClientVerifyGiftCode(c *gin.Context) {
	var req model.VerifyRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(400, gin.H{MESSAGE: err.Error()})
		return
	}

	// 检查礼品码
	if req.Code == "" {
		c.JSON(401, gin.H{MESSAGE: "请输入礼品码"})
		return
	}

	// 检查领取用户
	if req.User == "" {
		c.JSON(402, gin.H{MESSAGE: "请输入领取用户"})
		return
	}

	result, err := service.ClientVerifyGiftCode(req.User, req.Code)

	if result == nil {
		switch err {
		case "用户不存在，请先注册", "礼品码已失效或者礼包已领取":
			c.JSON(403, gin.H{MESSAGE: "用户不存在或者礼品码已失效请检查重试"})
		case "序列化protobuf失败":
			c.JSON(404, gin.H{MESSAGE: "序列化protobuf失败"})
		}
		return
	} else {

	}
	c.Writer.Write(result)
	c.JSON(200, gin.H{MESSAGE: "获取成功"})
}
