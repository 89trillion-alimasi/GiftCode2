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
		c.JSON(ParameterBindingIsUnsuccessful, gin.H{MESSAGE: StatusText(ParameterBindingIsUnsuccessful)})
		return
	}

	// 检查礼品码类型
	switch gift.Type {
	case 1: // 1 - 指定用户一次性消耗
		// 一次性消耗强制设置领取次数为 1
		gift.AvailableTimes = 1
		if gift.ReceivingUser == "" {
			c.JSON(ReceivingUserIsEmpty, gin.H{MESSAGE: StatusText(ReceivingUserIsEmpty)})
			return
		}
	case 2: // 2 - 不指定用户限制兑换次数
		if gift.AvailableTimes <= 0 {
			c.JSON(SpecifyTheNumberOfRedemptions, gin.H{MESSAGE: StatusText(SpecifyTheNumberOfRedemptions)})
			return
		}
	case 3: // 3 - 不限用户不限次数兑换 无用户限制 无兑换次数限制这里不做处理
	default: // 非法填写的礼品码类型
		c.JSON(Invalidgiftcodetype, gin.H{MESSAGE: StatusText(Invalidgiftcodetype)})
		return
	}

	// 检查礼品码描述
	if gift.Description == "" {
		c.JSON(GiftCodeDescription, gin.H{MESSAGE: StatusText(GiftCodeDescription)})
		return
	}

	// 检查礼品码有效期是否为空
	if gift.ValidPeriod == 0 {
		c.JSON(PleaseEnterAValidTime, gin.H{MESSAGE: StatusText(PleaseEnterAValidTime)})
		return
	}

	// 检查礼品码有效期是否正确
	//vp, err := time.ParseInLocation("2006-01-02 15:04:05", gift.ValidPeriod, time.Local)
	//if err != nil {
	//	c.JSON(IncorrectTimeFormat, gin.H{MESSAGE: StatusText(IncorrectTimeFormat)})
	//	return
	//}
	//exp := vp.Sub(time.Now())
	if gift.ValidPeriod <= time.Now().Unix() {
		c.JSON(InvalidTime, gin.H{MESSAGE: StatusText(InvalidTime)})
		return
	}

	// 检查礼品码包含的礼品包是否为空
	if gift.GiftPackages == nil || len(gift.GiftPackages) == 0 {
		c.JSON(PackageContent, gin.H{MESSAGE: StatusText(PackageContent)})
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
				c.String(CreateGiftCodeFaied, StatusText(CreateGiftCodeFaied))
				return
			}
		}

	}
	logrus.Infof("生成新的礼品码: %s", gift.Code)

	if err := service.SaveGiftCode(&gift); err != nil {
		c.JSON(InsertionFailed, gin.H{MESSAGE: StatusText(InsertionFailed)})
		return
	}

	// 返回礼品码
	c.JSON(CreatedSuccessfully, gin.H{"code": StatusText(CreatedSuccessfully)})
}

//查询礼品码
func QueryGiftCode(c *gin.Context) {
	// 检测是否输入了礼品码
	code := c.Query("code")
	if code == "" {
		c.JSON(GiftCodeErr, gin.H{MESSAGE: StatusText(GiftCodeErr)})
		return
	}

	// 尝试从 redis 中查找该礼品码
	giftCode := service.QueryGiftCode(code)
	if giftCode == nil {
		c.JSON(GiftCodeHasExpired, gin.H{MESSAGE: StatusText(GiftCodeHasExpired)})
		return
	}

	// 返回礼品码
	c.JSON(200, giftCode)
}

//登陆
func Login(c *gin.Context) {
	user_id := c.Query("userid")
	if len(user_id) == 0 {
		c.JSON(Userlogin, gin.H{MESSAGE: StatusText(Userlogin)})
		return
	}

	user, err := service.Login(user_id)

	if err == "用户存在成功返回" {
		c.JSON(Successful, gin.H{"Status": err, "User": user})
		return
	} else if err == "用户不存在，已为用户主动注册" {
		c.JSON(Registered, gin.H{"Status": err, "NewUser": user})
		return
	} else if err == "插入不成功" {
		c.JSON(InsertionUserFailed, gin.H{MESSAGE: StatusText(InsertionUserFailed)})
	}

}

//验证礼品码
func VerifyGiftCode(c *gin.Context) {
	var req model.VerifyRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(ParameterBindingIsUnsuccessful, gin.H{MESSAGE: StatusText(PleaseEnterAValidTime)})
		return
	}

	// 检查礼品码
	if req.Code == "" {
		c.JSON(GiftCodeErr, gin.H{MESSAGE: StatusText(GiftCodeErr)})
		return
	}

	// 检查领取用户
	if req.User == "" {
		c.JSON(ReceivingUserIsEmpty, gin.H{MESSAGE: StatusText(ReceivingUserIsEmpty)})
		return
	}

	gifCode, err := service.VerifyGiftCode(req)
	if err != nil {
		c.JSON(FailedToClaim, gin.H{MESSAGE: StatusText(FailedToClaim)})
		return
	}
	c.JSON(Successful, gin.H{"GiftPackages": gifCode.GiftPackages})
}

func ClientVerifyGiftCode(c *gin.Context) {
	var req model.VerifyRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(ParameterBindingIsUnsuccessful, gin.H{MESSAGE: StatusText(ParameterBindingIsUnsuccessful)})
		return
	}

	// 检查礼品码
	if req.Code == "" {
		c.JSON(GiftCodeErr, gin.H{MESSAGE: StatusText(GiftCodeErr)})
		return
	}

	// 检查领取用户
	if req.User == "" {
		c.JSON(ReceivingUserIsEmpty, gin.H{MESSAGE: StatusText(ReceivingUserIsEmpty)})
		return
	}

	result, err := service.ClientVerifyGiftCode(req.User, req.Code)

	if result == nil {
		switch err {
		case "用户不存在，请先注册", "礼品码已失效或者礼包已领取":
			c.JSON(FailedToClaim, gin.H{MESSAGE: "用户不存在或者礼品码已失效请检查重试"})
		case "序列化protobuf失败":
			c.JSON(FailedToClaim, gin.H{MESSAGE: "序列化protobuf失败"})
		}
		return
	} else {

	}
	c.Writer.Write(result)
	c.JSON(Successful, gin.H{MESSAGE: "获取成功"})
}
