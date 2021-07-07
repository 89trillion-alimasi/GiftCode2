package service

import (
	Protobuf "Giftcode-mongo-protobuf/Protobuf"
	"Giftcode-mongo-protobuf/config"
	"Giftcode-mongo-protobuf/model"
	"Giftcode-mongo-protobuf/model/Dboperation"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"strconv"
	"time"
)

// SaveGiftCode 将礼品码存储到config 中
func SaveGiftCode(gift *model.GiftCode) error {
	gift.CreatTime = time.Now().Unix()
	bs, err := jsoniter.Marshal(gift)
	if err != nil {
		return err
	}
	ex := time.Unix(gift.ValidPeriod, 0).Sub(time.Now())
	return config.RDB.Set(gift.Code, string(bs), ex).Err()
}

// QueryGiftCode 查询礼品码信息
func QueryGiftCode(code string) *model.GiftCode {
	cmd := config.RDB.Get(code)
	if cmd.Err() != nil {
		logrus.Warnf("礼品码未找到: %s", cmd.Err())
		return nil
	}
	var giftCode model.GiftCode
	err := jsoniter.Unmarshal([]byte(cmd.Val()), &giftCode)
	if cmd.Err() != nil {
		logrus.Warnf("礼品码反序列化失败: %s", err)
		return nil
	}
	return &giftCode
}

// VerifyGiftCode 验证礼品码
func VerifyGiftCode(req model.VerifyRequest) (*model.GiftCode, error) {
	// 从 config 中获取礼品码
	gifCode := QueryGiftCode(req.Code)
	if gifCode == nil {
		return nil, errors.New("礼品码不存在或已过期")
	}

	// 检查礼品码类型
	switch gifCode.Type {
	case 1: // 1 - 指定用户一次性消耗
		if req.User != gifCode.ReceivingUser {
			return nil, errors.New("当前礼品码已经指定用户，您输入的用户无权领取")
		}
		if _, err := One_time(gifCode); err != nil {
			return nil, err
		}

		// TODO: 当前没有真实用户体系，所以这里模拟添加奖励
		logrus.Infof("用户 %s 添加奖励完成", req.User)
		return gifCode, nil

	case 2: // 2 - 不指定用户限制兑换次数

		// 如果礼品码正好还剩一次可以领取，领取后需要删除
		if gifCode.AvailableTimes == 1 {
			if _, err := One_time(gifCode); err != nil {
				return nil, err
			}
			// TODO: 当前没有真实用户体系，所以这里模拟添加奖励
			logrus.Infof("用户 %s 添加奖励完成", req.User)
			return gifCode, nil
		} else {
			// 如果礼品码剩余领取次数大于 1，领取后对可用次数 -1
			gifCode.AvailableTimes--
			// 将领取用户加入到已获取礼品码 领取列表中
			gifCode.AddReceivedUser(req.User)

			// 重新将礼品码存入 redis
			bs, err := jsoniter.Marshal(gifCode)
			if err != nil {
				return nil, err
			}

			logrus.Infof("用户 %s 添加奖励完成", req.User)
			return gifCode, config.RDB.Set(gifCode.Code, string(bs), gifCode.Expiration).Err()
		}
	case 3: // 3 - 不限用户不限次数兑换 无用户限制 无兑换次数限制这里不做处理
		return gifCode, nil
	default: // 非法填写的礼品码类型
		logrus.Infof("用户 %s 添加奖励完成", req.User)
		return gifCode, nil
	}
}

//用户登录和注册
func Login(user_id string) (model.UerInfo, string) {
	user, _ := Dboperation.Searchmongo("uid", user_id)

	var err string
	if user.UID == "" {
		fmt.Println("UID  is  empty")
		err = "用户不存在，已为用户主动注册"

		newUID := Userinfo(8)
		var flag bool
		for flag = true; flag; {
			//是否存在用户
			if Dboperation.ExistUserId(newUID, "user_info") {
				newUID = Userinfo(8)
			} else {
				flag = false
			}
		}
		if Dboperation.Insearchmongo(model.UerInfo{
			UID:     newUID,
			Gold:    "0",
			Diamond: "0",
		}) {
			user = model.UerInfo{
				UID:     newUID,
				Gold:    "0",
				Diamond: "0",
			}
		} else {
			err = "插入不成功"
			return model.UerInfo{}, err
		}

	} else if user.UID != "" {
		err = "用户存在成功返回"
	}
	return user, err
}

//客户请求获取礼品码并返回protobuf数据
func ClientVerifyGiftCode(user_id, gifcode string) ([]byte, string) {
	var user model.VerifyRequest
	//先判断用户id是否存在
	uuser, err1 := Dboperation.Searchmongo("uid", user_id)
	if err1 != nil {
		return nil, "用户不存在，请先注册"
	}
	//把用户值和礼品码值赋予给user
	user.User = user_id
	user.Code = gifcode
	//在验证礼品码有效性
	gift, err2 := VerifyGiftCode(user)

	if err2 != nil {
		return nil, "礼品码已失效或者礼包已领取"
	}

	giftpackage := gift.GiftPackages

	var (
		gold    int
		diamond int
	)

	//礼品包的金币和钻石的值拿出来
	for _, v := range giftpackage {
		if v.Name == "金币" {
			gold = v.Num
		} else {
			diamond = v.Num
		}
	}

	changes := make(map[uint32]uint64) //道具增加量
	balance := make(map[uint32]uint64) //当前余额
	counter := make(map[uint32]uint64) //总共：增加量+余额

	changes[1] = uint64(gold)
	changes[2] = uint64(diamond)

	//金币数更新
	gold_mongo, _ := strconv.Atoi(uuser.Gold)
	gold_new := strconv.Itoa(gold_mongo + gold)
	uuser.Gold = gold_new

	//钻石数更新
	diamond_mongo, _ := strconv.Atoi(uuser.Diamond)
	dianmond_new := strconv.Itoa(diamond_mongo + diamond)
	uuser.Diamond = dianmond_new

	//用户余额
	balance[1] = uint64(gold_mongo)
	balance[2] = uint64(diamond_mongo)

	//用户礼包总和
	counter[1] = uint64(gold_mongo + gold)
	counter[2] = uint64(diamond_mongo + diamond)

	//更新数据库
	Dboperation.Updatamongo("uid", uuser.UID, "gold", uuser.Gold)
	Dboperation.Updatamongo("uid", uuser.UID, "diamond", uuser.Diamond)

	protobufinfo := Protobuf.GeneralReward{
		Code:    int32(1),
		Msg:     "",
		Changes: changes,
		Balance: balance,
		Counter: counter,
		Ext:     "",
	}
	out, err := proto.Marshal(&protobufinfo)

	if err != nil {
		return nil, "序列化protobuf失败"
	}

	return out, "success"
}

//只能取一次礼品
func One_time(code *model.GiftCode) (*model.GiftCode, error) {
	cmd := config.RDB.Del(code.Code)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	if cmd.Val() != 1 {
		return nil, errors.New("礼品码领取失败: 礼品码删除失败")
	}
	return code, nil
}
