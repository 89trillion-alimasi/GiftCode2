package service

import (
	"Giftcode-mongo-protobuf/Protobuf"
	"Giftcode-mongo-protobuf/model"
	"google.golang.org/protobuf/proto"
	"log"
	"reflect"
	"testing"
)

//登陆测试
func TestLogin(t *testing.T) {
	type tests struct {
		userid string
		want   model.UerInfo
	}

	test := []tests{
		{userid: "HrKcLApw", want: model.UerInfo{UID: "HrKcLApw", Gold: "10", Diamond: "180"}},
	}

	for _, v := range test {
		got, _, _ := Login(v.userid)

		if reflect.DeepEqual(got, v.want) {
			t.Errorf("actully%v,expected %v", got, v.want)
		}
	}

}

//客户领取礼品测试
func TestClientVerifyGiftCode(t *testing.T) {
	type tests struct {
		userid string
		code   string
		want   Protobuf.GeneralReward
	}
	test := []tests{
		{userid: "HrKcLApw", code: "1ZAZB1gj", want: Protobuf.GeneralReward{
			Code:    int32(1),
			Counter: map[uint32]uint64{1: uint64(30), 2: uint64(200)},
			Changes: map[uint32]uint64{1: uint64(10), 2: uint64(20)},
			Balance: map[uint32]uint64{1: uint64(20), 2: uint64(180)},
			Msg:     "",
			Ext:     "",
		}},
	}
	for _, v := range test {
		got, _, _ := ClientVerifyGiftCode(v.userid, v.code)
		data := &Protobuf.GeneralReward{}
		err := proto.Unmarshal(got, data)
		if err != nil {
			log.Fatal("unmarshaling error: ", err)
		}
		if !reflect.DeepEqual(data, v.want) {
			t.Errorf("actully%v,expected %v", data, v.want)
		}
	}
}
