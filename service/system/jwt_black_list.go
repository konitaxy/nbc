package system

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/system"
)

type JwtService struct{}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: JsonInBlacklist
//@description: 拉黑jwt
//@param: jwtList model.JwtBlacklist
//@return: err error

func (jwtService *JwtService) JsonInBlacklist(jwtList system.JwtBlacklist) (err error) {
	err = global.GVA_DB.Create(&jwtList).Error
	if err != nil {
		return
	}
	global.BlackCache.SetDefault(jwtList.Jwt, struct{}{})
	return
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: IsBlacklist
//@description: 判断JWT是否在黑名单内部
//@param: jwt string
//@return: bool

func (jwtService *JwtService) IsBlacklist(jwt string) bool {
	_, ok := global.BlackCache.Get(jwt)
	return ok
	// err := global.GVA_DB.Where("jwt = ?", jwt).First(&system.JwtBlacklist{}).Error
	// isNotFound := errors.Is(err, gorm.ErrRecordNotFound)
	// return !isNotFound
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetRedisJWT
//@description: 从redis取jwt
//@param: userName string
//@return: err error, redisJWT string

func (jwtService *JwtService) GetRedisJWT(userName string) (err error, redisJWT string) {
	redisJWT, err = global.GVA_REDIS.Get(context.Background(), userName).Result()
	return err, redisJWT
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetRedisJWT
//@description: jwt存入redis并设置过期时间
//@param: jwt string, userName string
//@return: err error

func (jwtService *JwtService) SetRedisJWT(jwt string, userName string) (err error) {
	// 此处过期时间等于jwt过期时间
	timer := time.Duration(global.GVA_CONFIG.JWT.ExpiresTime) * time.Second
	err = global.GVA_REDIS.Set(context.Background(), userName, jwt, timer).Err()
	return err
}

func (jwtService *JwtService) SetRedisSnapJWT(jwt string, userName string) (err error) {
	// 此处过期时间等于jwt过期时间
	timer := time.Duration(global.GVA_CONFIG.JWT.SnapExpiresTime) * time.Second
	err = global.GVA_REDIS.Set(context.Background(), userName, jwt, timer).Err()
	return err
}

// IAM 用户专用的 Redis key 前缀
const IAMJwtPrefix = "iam_jwt_"

// GetRedisIAMJWT 从redis取IAM用户jwt
func (jwtService *JwtService) GetRedisIAMJWT(email string) (err error, redisJWT string) {
	redisJWT, err = global.GVA_REDIS.Get(context.Background(), IAMJwtPrefix+email).Result()
	return err, redisJWT
}

// SetRedisIAMJWT IAM用户jwt存入redis并设置过期时间
func (jwtService *JwtService) SetRedisIAMJWT(jwt string, email string) (err error) {
	timer := time.Duration(global.GVA_CONFIG.JWT.ExpiresTime) * time.Second
	err = global.GVA_REDIS.Set(context.Background(), IAMJwtPrefix+email, jwt, timer).Err()
	return err
}
func (jwtService *JwtService) RedisSetUserStatus(userName string, status uint) (err error) {
	// 此处过期时间等于jwt过期时间
	timer := time.Duration(global.GVA_CONFIG.JWT.SnapExpiresTime) * time.Second

	err = global.GVA_REDIS.Set(context.Background(), fmt.Sprintf("client_status_%s", userName), status, timer).Err()
	return err
}

func (jwtService *JwtService) RedisIsFreeze(userName string) bool {
	// 此处过期时间等于jwt过期时间

	status := global.GVA_REDIS.Get(context.Background(), fmt.Sprintf("client_status_%s", userName)).Val()
	return status == "3"
}

func LoadAll() {
	var data []string
	err := global.GVA_DB.Model(&system.JwtBlacklist{}).Select("jwt").Find(&data).Error
	if err != nil {
		global.GVA_LOG.Error("加载数据库jwt黑名单失败!", zap.Error(err))
		return
	}
	for i := 0; i < len(data); i++ {
		global.BlackCache.SetDefault(data[i], struct{}{})
	} // jwt黑名单 加入 BlackCache 中
}
