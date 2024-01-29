package v1

import (
	"coze-chat-proxy/common"
	jsoniter "github.com/json-iterator/go"
	"math/rand"
	"time"
)

func GenerateID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	id := "chatcmpl-"
	for i := 0; i < length; i++ {
		id += string(charset[rand.Intn(len(charset))])
	}
	return id
}

func Obj2Bytes(obj interface{}) ([]byte, error) {
	// 创建一个jsonIter的Encoder
	configCompatibleWithStandardLibrary := jsoniter.ConfigCompatibleWithStandardLibrary
	// 将结构体转换为JSON文本并保持顺序
	bytes, err := configCompatibleWithStandardLibrary.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func SetTimer(outTime time.Duration) (*time.Timer, error) {
	if outTime == 0 {
		outTime = common.StreamRequestOutTime
	}
	return time.NewTimer(outTime), nil

}

func TimerReset(timer *time.Timer, outTime time.Duration) error {
	if outTime == 0 {
		outTime = common.StreamRequestOutTime
	}
	timer.Reset(outTime)
	return nil
}
