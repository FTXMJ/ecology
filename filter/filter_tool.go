package filter

import (
	db "ecology/db"
	"ecology/models"
	"github.com/gin-gonic/gin"

	"time"
)

// 接口用户身份认证
func ApiClassification(api string) string {
	if api == "/api/v1/quote/ticker" {
		return "other"
	}

	// true=管理员
	if api[16:21] == "admin" {
		return "admin"
	} else if api == "/api/v1/ecology/swagger" || api == "/api/v1/ecology/check" {
		return "nil" // 服务发现 心跳检测 以及 swagger
	} else {
		return "user"
	}
}

// 普通用户的 过滤器处理规则
func UserFilter(ctx *gin.Context, api string) (code int, msg string, c *CustomClaims) {
	token := ctx.Request.Header.Get("Authorization")
	if token == "" {
		return 401, "未经允许的访问，已拦截！", nil
	}
	j := NewJWT()
	// parseToken 解析token包含的信息
	tockken, err := j.ParseToken(token)
	if err != nil {
		if err == TokenExpired {
			return 401, "未经允许的访问，已拦截！", tockken
		}
		return 401, "数据格式不正确", tockken
	}
	o := db.NewEcologyOrm()
	u := models.User{}

	_, nicke_name, err_ping_user := models.PingUserAdmin(token, tockken.UserID)
	if err_ping_user != nil {
		return 500, err_ping_user.Error(), tockken
	}
	err_read := o.Table("user").Where("user_id = ?", tockken.UserID).First(&u)
	if err_read.Error != nil && err_read.Error.Error() == "<QuerySeter> no row found" && tockken.NameSpace == "" {
		o.Begin()
		f, _, err_get_user := models.PingUser(token, tockken.UserID)
		if err_get_user != nil {
			return 500, "用户不存在!", tockken
		}
		user := models.User{}
		if f.(string) != "" {
			user.UserId = tockken.UserID
			user.FatherId = f.(string)
			user.UserName = tockken.Name
		} else {
			user.UserId = tockken.UserID
			user.UserName = tockken.Name
		}

		erruser := o.Create(&user)
		if erruser.Error != nil {
			o.Rollback()
			return 500, "创建用户失败!", tockken
		}
		account_def := models.Account{
			UserId:         tockken.UserID,
			CreateDate:     time.Now().Format("2006-01-02 15:04:05"),
			DynamicRevenue: true,
			StaticReturn:   true,
			PeerState:      true,
		}
		account_def_err := o.Create(&account_def)
		if account_def_err.Error != nil {
			o.Rollback()
			return 500, "创建生态仓库失败!", tockken
		}
		formula := models.Formula{
			EcologyId:      account_def.Id,
			ReturnMultiple: 1,
		}
		err_for := o.Create(&formula)
		if err_for.Error != nil {
			o.Rollback()
			return 500, "创建算力表失败!", tockken
		}
		o.Commit()
	} else if err_read.Error == nil && u.UserName != nicke_name.(string) && tockken.NameSpace == "" {
		o.Begin()
		u.UserName = nicke_name.(string)
		o.Update(&u, "user_name")
		o.Commit()
	} else if err_read.Error != nil && err_read.Error.Error() != "<QuerySeter> no row found" {
		return 500, "后端服务期错误!", tockken
	}
	return 0, "", tockken
}

// 管理员用户的 过滤器处理规则
func AdminFilter(ctx *gin.Context, api string) (code int, msg string, c *CustomClaims) {
	token := ctx.Request.Header.Get("Authorization")

	if token == "" {
		return 401, "未经允许的访问，已拦截！", nil
	}

	j := NewJWT()
	// parseToken 解析token包含的信息
	tockken, err := j.ParseToken(token)
	if err != nil {
		if err == TokenExpired {
			return 401, "未经允许的访问，已拦截！", tockken
		}
		return 401, "数据格式不正确", tockken
	} else if tockken.NameSpace != "admin" {
		return 401, "未经允许的访问，已拦截！", tockken
	}
	_, _, err_ping_user := models.PingUserAdmin(token, tockken.UserID)
	if err_ping_user != nil {
		return 500, err_ping_user.Error(), tockken
	}
	return 0, "", tockken
}
