package controller

const (
	ParameterBindingIsUnsuccessful = 400
	ReceivingUserIsEmpty           = 401
	SpecifyTheNumberOfRedemptions  = 402
	Invalidgiftcodetype            = 404
	GiftCodeDescription            = 405
	PleaseEnterAValidTime          = 406
	PackageContent                 = 408
	CreateGiftCodeFaied            = 409
	CreatedSuccessfully            = 201
	GiftCodeHasExpired             = 410
	GiftCodeErr                    = 411
	FailedToClaim                  = 412
	Successful                     = 200
	InvalidTime                    = 403
	InsertionFailed                = 413
	Userlogin                      = 414
	Registered                     = 202
	InsertionUserFailed            = 507
	DecodeFaild                    = 508
)

var statusText = map[int]string{
	ParameterBindingIsUnsuccessful: "绑定参数未成功",
	ReceivingUserIsEmpty:           "请指定领取用户",
	SpecifyTheNumberOfRedemptions:  "指定可兑换次数",
	Invalidgiftcodetype:            "礼品码类型不合法",
	GiftCodeDescription:            "请输入礼品码描述信息",
	PleaseEnterAValidTime:          "请输入有效期",
	PackageContent:                 "请输入礼包内容",
	CreateGiftCodeFaied:            "创建礼品码失败",
	CreatedSuccessfully:            "成功",
	GiftCodeHasExpired:             "礼品码输入错误或礼品码已过期",
	GiftCodeErr:                    "请输入礼品码",
	FailedToClaim:                  "领取失败",
	Successful:                     "成功",
	InvalidTime:                    "有效时间小于当前本地的时间",
	InsertionFailed:                "插入礼品码失败",
	Userlogin:                      "请输入登陆用户",
	Registered:                     "已为用户注册",
	InsertionUserFailed:            "插入用户信息失败",
	DecodeFaild:                    "解析失败",
}

type Mesg struct {
	Code    int
	Message string
	Data    interface{}
}

// StatusText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.
func StatusText(code int) Mesg {
	return Mesg{
		Code:    code,
		Message: statusText[code],
	}
}

func StatusText1(code int, data interface{}) Mesg {
	return Mesg{
		Code:    code,
		Message: statusText[code],
		Data:    data,
	}
}
