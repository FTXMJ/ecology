package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

//生成所用的 32 位 hash
func Shengchengstr(task_sum, user_id string, operate_type string) string {
	curr_time := time.Now().Unix()
	q := strconv.Itoa(int(curr_time))
	str := task_sum + user_id + operate_type + q
	return GetStringSha256(str)
}

func GetStringSha256(in string) string {
	hash := sha256.New()
	hash.Write([]byte(in))
	bytes := hash.Sum(nil)
	byteCode := hex.EncodeToString(bytes)[:32]
	return byteCode
}


