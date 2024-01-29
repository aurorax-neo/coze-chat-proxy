package v1

import (
	"coze-chat-proxy/common"
	jsoniter "github.com/json-iterator/go"
	"math/rand"
	"strconv"
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

func SetTimer(isStream bool, defaultTimeout time.Duration) (*time.Timer, error) {
	var outTimeStr string
	if isStream {
		outTimeStr = common.StreamRequestOutTime
	} else {
		outTimeStr = common.RequestOutTime
	}
	if outTimeStr != "" {
		outTime, err := strconv.ParseInt(outTimeStr, 10, 64)
		if err != nil {

			return nil, err
		}
		return time.NewTimer(time.Duration(outTime) * time.Second), nil
	}
	return time.NewTimer(defaultTimeout), nil
}

func TimerReset(isStream bool, timer *time.Timer, defaultTimeout time.Duration) error {
	var outTimeStr string
	if isStream {
		outTimeStr = common.StreamRequestOutTime
	} else {
		outTimeStr = common.RequestOutTime
	}
	if outTimeStr != "" {
		outTime, err := strconv.ParseInt(outTimeStr, 10, 64)
		if err != nil {
			return err
		}
		timer.Reset(time.Duration(outTime) * time.Second)
		return nil
	}
	timer.Reset(defaultTimeout)
	return nil
}
