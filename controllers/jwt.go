package controllers

import (
	"ecology1/models"
	"errors"
	"fmt"
	"github.com/astaxie/beego/context"
	"github.com/dgrijalva/jwt-go"
	"time"
)

/*func init() {
	//拦截器验证登录
	beego.InsertFilter("/*", beego.BeforeExec, CheckLogin)
}*/

// 一些常量
var (
	TokenExpired     error = errors.New("Token is expired")
	TokenNotValidYet error = errors.New("Token not active yet")
	TokenMalformed   error = errors.New("That's not even a token")
	TokenInvalid     error = errors.New("Couldn't handle this token:")
	SignKey          string
)

// 载荷，可以加一些自己需要的信息
type CustomClaims struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Mobile   string `json:"mobile"`
	Email    string `json:"email"`
	FatherId string `json:"father_id"`
	jwt.StandardClaims
}

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

//验证登录
func CheckLogin(ctx *context.Context) {
	api := ctx.Request.URL.Path
	//不包含在内的验证登录
	if api != "/v1/user/login" && //登录页面
		api != "/v1/user/register" && //注册
		api != "/v1/user/forgetPassword" { //忘记密码
		token := ctx.Request.Header.Get("Authorization")
		if token == "" {
			fmt.Println("拦截：", api)
			ctx.WriteString(`{"code": "500","msg": "未经允许的访问，已拦截！"}`)
			return
		}
		j := NewJWT()
		// parseToken 解析token包含的信息
		tockken, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				ctx.WriteString(`{"code": "401","msg": "未经允许的访问，已拦截！"}`)
				return
			}
			ctx.WriteString(`{"code": "400","msg": "数据格式不正确"}`)
			return
		}
		user := []models.User{}
		models.NewOrm().QueryTable("user").Filter("user_id", tockken.UserID).All(&user)
		if len(user) == 0 {
			user := models.User{
				Name:   tockken.Name,
				UserId: tockken.UserID,
				FatherId:tockken.FatherId,
			}
			erruser := user.Insert()
			if erruser != nil {
				ctx.WriteString(`{"code": "500","msg": "后端服务期错误(db)"}`)
				return
			}
		}
	}
}

// 获取jst里面的value
func GetJwtValues(ctx *context.Context) *CustomClaims {
	token := ctx.Request.Header.Get("Authorization")
	j := NewJWT()
	// parseToken 解析token包含的信息
	tocken, _ := j.ParseToken(token)
	return tocken
}

// 生成令牌
func generateToken(user models.User) (bool, string) {
	j := &JWT{
		[]byte(SignKey),
	}
	claims := CustomClaims{
		user.UserId,
		"",
		"",
		"",
		"",
		jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000),       // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 3600*24*15), // 过期时间 一小时
			Issuer:    "ecology",                             //签名的发行者
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		return false, ""
	} else {
		return true, token
	}
}

// 新建一个jwt实例
func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

// 获取signKey
func GetSignKey() string {
	return SignKey
}

// 这是SignKey
func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

// CreateToken 生成一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析Tokne
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

// 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
