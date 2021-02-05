package token

import (
	"errors"
	"gorm.io/gorm"
	"strconv"
	"ziyun/util/log"
	"ziyun/util/mysql"
)

type ClientDetails struct {
	gorm.Model
	// client 的标识
	//ClientId string
	// client 的密钥
	ClientSecret string
	// 访问令牌有效时间，秒
	AccessTokenValiditySeconds int
	// 刷新令牌有效时间，秒
	RefreshTokenValiditySeconds int
	// 重定向地址，授权码类型中使用
	RegisteredRedirectUri string
	// 可以使用的授权类型
	//AuthorizedGrantTypes []string
}

type UserDetails struct {
	gorm.Model
	// 用户标识
	//UserId string
	// 用户名 唯一
	Username string
	// 用户密码
	Password string
	// 用户具有的权限
	//Authorities []string // 具备的权限
	// client 的标识
	ClientId string
}

type OAuth2Details struct {
	Client *ClientDetails
	User   *UserDetails
}

////////////////////////以下为mysql相关操作//////////////////////////////////
//此为示例程序，没把client信息和user信息预写入到redis中，都是直接访问数据库获得信息，实际应用中需要修改.
//端上需要通过用户名+密码+clientID换取userID，申请并生成token的后续操作都是通过userID和clientID完成。
//userid通过数据库获取，不同company的user均在一个库中， 不需要companyid_userid方式来区分用户

// 访问数据库，返回用户所属client的信息
func GetClientDetailsMysql(companyid string) (*ClientDetails, error) {

	client := new(ClientDetails)
	mysql.Db.AutoMigrate(client)

	//查找数据
	err := mysql.Db.Where("ID = ?", companyid).Find(&client).Error
	if err != nil {
		log.Logrus.Errorln(companyid, "GetClientDetailsMysql: ", err.Error())
		return nil, err
	} else if client.ID == 0 {
		log.Logrus.Errorln(companyid, "GetClientDetailsMysql failure: not found")
		return nil, errors.New("not found")
	} else {
		log.Logrus.Debugln(companyid, "GetUserDetailsMysql ok")
		return client, nil
	}
}

// 访问数据库，并返回用户的信息
func GetUserDetailsMysql(userid string) (*UserDetails, error) {
	user := new(UserDetails)
	mysql.Db.AutoMigrate(user)

	//查找数据
	err := mysql.Db.Where("ID = ?", userid).Find(&user).Error
	if err != nil {
		log.Logrus.Errorln(userid, "GetUserDetailsMysql: ", err.Error())
		return nil, err
	} else if user.ID == 0 {
		log.Logrus.Errorln(userid, "GetUserDetailsMysql failure: not found")
		return nil, errors.New("not found")
	} else {
		log.Logrus.Debugln(userid, "GetUserDetailsMysql ok")
		return user, nil
	}
}

// 访问数据库，验证用户名密码是否正确, 并返回用户ID
func CheckUserMysql(username, password, companyid string) (string, error) {
	user := new(UserDetails)
	mysql.Db.AutoMigrate(user)

	//查找数据
	err := mysql.Db.Where("Username = ? AND Password = ? AND Client_Id = ? ", username, password, companyid).Find(&user).Error
	if err != nil {
		log.Logrus.Errorln(username, "CheckUserMysql: ", err.Error())
		return "", err
	} else if user.ID == 0 {
		log.Logrus.Debugln(username, "CheckUserMysql failure: not found")
		return "", errors.New("not found")
	} else {
		log.Logrus.Debugln(username, "CheckUserMysql ok ", strconv.Itoa(int(user.ID)))
		return strconv.Itoa(int(user.ID)), nil
	}
}

/*
////////////////////////以下为client在redis的操作//////////////////////////////////
//把Client相关信息写入redis, hash结构存储client属性
func AddClientRedis(cli *ClientDetails) error{

	IDstr := strconv.Itoa(int(cli.ID))

	pipe := redis.RedisClient.Pipeline()
	//pipe.HSet("client_"+IDstr+"_info", "ID", IDstr)
	pipe.HSet("client_"+IDstr+"_info", "ClientSecret", cli.ClientSecret)
	pipe.HSet("client_"+IDstr+"_info", "AccessTokenValiditySeconds", strconv.Itoa(int(cli.AccessTokenValiditySeconds)))
	pipe.HSet("client_"+IDstr+"_info", "RefreshTokenValiditySeconds", strconv.Itoa(int(cli.RefreshTokenValiditySeconds)))
	pipe.HSet("client_"+IDstr+"_info", "RegisteredRedirectUri", cli.RegisteredRedirectUri)

	if _, err := pipe.Exec(); err != nil {
		log.Logrus.Errorln(IDstr, "AddClientRedis: ", err.Error())
		return err
	}else{
		return nil
	}
}

//从redis中删除Client相关信息
func DelClientRedis(clientId string) error{

	pipe := redis.RedisClient.Pipeline()
	pipe.Del("client_"+clientId+"_info")

	if _, err := pipe.Exec(); err != nil {
		log.Logrus.Errorln(clientId, "AddClientRedis: ", err.Error())
		return err
	}else{
		return nil
	}
}

//向set中增加clientid。 用set结构存储cilentid， set的key为 client_total_info
func AddClientMgrRedis(clientId string) error{

	pipe := redis.RedisClient.Pipeline()
	pipe.SAdd("client_total_info", clientId)

	if _, err := pipe.Exec(); err != nil {
		log.Logrus.Errorln(clientId, "AddClientMgrRedis: ", err.Error())
		return err
	}else{
		return nil
	}
}

//从set中删除clientid,,
func DelClientMgrRedis(clientId string) error{

	pipe := redis.RedisClient.Pipeline()
	pipe.SRem("client_total_info", clientId)

	if _, err := pipe.Exec(); err != nil {
		log.Logrus.Errorln(clientId, "DelClientMgrRedis: ", err.Error())
		return err
	}else{
		return nil
	}
}

////////////////////////以下为user在redis的操作//////////////////////////////////
//把user相关信息写入redis, hash结构存储user属性
func AddUserRedis(user *UserDetails) error{
	IDstr := strconv.Itoa(int(user.ID))

	pipe := redis.RedisClient.Pipeline()
	//pipe.HSet("user_"+IDstr+"_info", "ID", IDstr)
	pipe.HSet("user_"+IDstr+"_info", "Username", user.Username)
	pipe.HSet("user_"+IDstr+"_info", "Password", user.Password)

	if _, err := pipe.Exec(); err != nil {
		log.Logrus.Errorln(IDstr, "AddUserRedis: ", err.Error())
		return err
	}else{
		return nil
	}
}

//从redis中删除user相关信息
func DelUserRedis(userId string) error{
	IDstr := userId

	pipe := redis.RedisClient.Pipeline()
	pipe.Del("user_"+IDstr+"_info")

	if _, err := pipe.Exec(); err != nil {
		log.Logrus.Errorln(userId, "AddClientRedis: ", err.Error())
		return err
	}else{
		return nil
	}
}

//根据id校验用户密码
func CheckUserRedis(userId, password string) error{
	IDstr := userId

	pipe := redis.RedisClient.Pipeline()
	pw := pipe.HGet("user_"+IDstr+"_info", "Password")
	if _, err := pipe.Exec(); err != nil {
		log.Logrus.Errorln(IDstr, "CheckUserRedis: ", err.Error())
		return err
	}

	//密码正确与否都错误都返回nil
	if pw.Val() == password {
		return nil
	}else {
		return ErrInvalidUsernameAndPassword
	}
}

//增加userid管理到redis, 用set结构存储userid， set的key为 user_total_info
func AddUserMgrRedis(userId string) error{
	pipe := redis.RedisClient.Pipeline()
	pipe.SAdd("user_total_info", userId)

	if _, err := pipe.Exec(); err != nil {
		log.Logrus.Errorln(userId, "AddUserMgrRedis: ", err.Error())
		return err
	}else{
		return nil
	}
}

func DelUserMgrRedis(userId string) error{
	pipe := redis.RedisClient.Pipeline()
	pipe.SRem("user_total_info", userId)

	if _, err := pipe.Exec(); err != nil {
		log.Logrus.Errorln(userId, "DelClientMgrRedis: ", err.Error())
		return err
	}else{
		return nil
	}
}
*/
