package filter

import (
	"ecology/models"
	"github.com/gin-gonic/gin"

	"github.com/dgrijalva/jwt-go"

	"errors"
	"time"
)

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
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	Mobile    string `json:"mobile"`
	Email     string `json:"email"`
	FatherId  string `json:"father_id"`
	NameSpace string `json:"name_space"`
	jwt.StandardClaims
}

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

// 获取jst里面的value
func GetJwtValues(ctx *gin.Context) *CustomClaims {
	token := ctx.Request.Header.Get("Authorization")
	j := NewJWT()
	// parseToken 解析token包含的信息
	tocken, _ := j.ParseToken(token)
	return tocken
}

// 生成令牌
func GenerateToken(user models.User) (bool, string) {
	j := &JWT{
		[]byte(SignKey),
	}
	claims := CustomClaims{
		user.UserId,
		"",
		"",
		"",
		"",
		"",
		jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000), // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 60*3), // 过期时间 三分钟
			Issuer:    "ecology",                       //签名的发行者
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
