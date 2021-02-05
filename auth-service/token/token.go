package token

import (
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"
	"ziyun/util/log"
	"ziyun/util/redis"
)

//access_token用jwt编码
//refresh-token用base64编码

var (
	ErrNotSupportOperation        = errors.New("no support operation")
	ErrRefreshTokenExpired        = errors.New("refresh_token expired")
	ErrInvalidRefreshToken        = errors.New("refresh_token Invalid")
	ErrInvalidClientIdAndUserId   = errors.New("invalid ClientID or UserID")
	ErrInvalidUsernameAndPassword = errors.New("invalid username, password")
)

var iJwt TokenEnhancer

//初始化
func InitJwtToken(secretKey string) {
	iJwt = NewJwtTokenEnhancer(secretKey)
}

//构建access token 和refresh token
func MakeToken(userID, clientID string) (string, string, error) {

	if userID == "" || clientID == "" {
		return "", "", ErrInvalidClientIdAndUserId
	}

	//后续从redis中读取， delay
	user_detail, err := GetUserDetailsMysql(userID)
	if err != nil {
		return "", "", ErrInvalidUsernameAndPassword
	}

	//后续从redis中读取， delay
	client_detail, err := GetClientDetailsMysql(clientID)
	if err != nil {
		return "", "", ErrInvalidClientIdAndUserId
	}

	detail := &OAuth2Details{
		client_detail,
		user_detail,
	}

	//构建Rrfresh Token, 采用base64方式
	rtoken := base64.StdEncoding.EncodeToString([]byte(clientID + "_" + userID + "_" + strconv.FormatInt(time.Now().Unix(), 10)))
	log.Logrus.Debugln("MakeToken: make rtoken: ", rtoken)
	//构建Access Token, 采用jwt方式
	atoken, err := iJwt.Enhance(detail)
	if err != nil {
		return "", "", ErrInvalidClientIdAndUserId
	}
	log.Logrus.Debugln("MakeToken: make atoken: ", atoken)

	// 保存新生成令牌
	SetRTokenToRedis(strconv.Itoa(int(detail.User.ID)), rtoken)

	return atoken, rtoken, nil
}

//校验atoken
func CheckAToken(atoken string) (string, error) {
	_, err := iJwt.Extract(atoken)

	if err == nil {
		return "true", nil
	} else {
		return "false", err
	}
}

//申请aToken和rToken
func Token(username, password, companyid string) (string, string, error) {

	// 校验公司id，用户id，密码是否为空
	if username == "" || password == "" || companyid == "" {
		return "", "", ErrInvalidUsernameAndPassword
	}

	// 访问数据库，验证用户名密码是否正确, 并返回userid
	// 登录成功后需要把用户信息写入到redis, 以后都由userid来驱动.
	userid, err := CheckUserMysql(username, password, companyid)
	if err != nil {
		return "", "", ErrInvalidUsernameAndPassword
	}

	return MakeToken(userid, companyid)
}

//带着rtoken来请求新的一对新token
func RefreshToken(rtoken string) (string, string, error) {

	//判空
	if rtoken == "" {
		return "", "", ErrInvalidRefreshToken
	}

	//解码refresh_token, 获取companyid， userid
	decodeBytes, err := base64.StdEncoding.DecodeString(rtoken)
	if err != nil {
		return "", "", err
	}
	st := strings.Split(string(decodeBytes), "_")
	clientId := st[0]
	userId := st[1]
	log.Logrus.Debugln("RefreshToken: decode rtoken: ", clientId, userId)

	//判断上传的refresh_token和redis中保存的是否一致，
	if rtoken_redis, err := GetRTokenFromRedis(userId); err != nil {
		log.Logrus.Errorln("RefreshToken: access redis failure: ")
		return "", "", err
	} else if rtoken_redis == "" {
		//redis中无refresh_token
		log.Logrus.Errorln("RefreshToken: get rtoken from redis not found: ", clientId, userId)
		return "", "", ErrRefreshTokenExpired
	} else if rtoken != rtoken_redis {
		//不一致
		log.Logrus.Errorln("RefreshToken: get rtoken from redis dismatch: ", rtoken, rtoken_redis)
		return "", "", ErrInvalidRefreshToken
	}

	//redis中有refresh token， 并且和用户上传的一致, 重新生成一组
	return MakeToken(userId, clientId)
}

//redis中只保存refresh token, 采用string存储, key是 userid_token
func SetRTokenToRedis(userid, rtoken string) error {
	key := userid + "_rtoken"

	pipe := redis.RedisClient.Pipeline()
	pipe.Set(key, rtoken, 7*24*time.Hour)

	if _, err := pipe.Exec(); err != nil {
		log.Logrus.Errorln(userid, "SetRTokenToRedis failure: ", rtoken, err.Error())
		return err
	} else {
		log.Logrus.Debugln(userid, "SetRTokenToRedis ok: ", key, rtoken)
		return nil
	}
}

//redis中只保存refresh token, 采用string存储, key是 userid_token
func GetRTokenFromRedis(userid string) (string, error) {
	key := userid + "_rtoken"

	pipe := redis.RedisClient.Pipeline()
	rtoken := pipe.Get(key)

	if _, err := pipe.Exec(); err != nil {
		log.Logrus.Errorln(userid, "GetRTokenFromRedis failure: ", rtoken, err.Error())
		return "", err
	} else {
		log.Logrus.Debugln(userid, "GetRTokenFromRedis ok: ", rtoken)
		return rtoken.Val(), nil
	}
}
