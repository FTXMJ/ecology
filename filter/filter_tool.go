package filter

import (
	"ecology/models"

	"github.com/astaxie/beego/context"

	"fmt"
	"time"
)

// 接口用户身份认证
func ApiClassification(api string) string {
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
func UserFilter(ctx *context.Context, api string) {
	token := ctx.Request.Header.Get("Authorization")
	if token == "" {
		fmt.Println("拦截：", api)
		ctx.WriteString(`{"code": "401","msg": "未经允许的访问，已拦截！"}`)
		fmt.Println(GenerateToken(models.User{
			UserId: "77e3732c1e4541bebf3782b43631b8b1",
		}))
		return
	}
	j := NewJWT()
	// parseToken 解析token包含的信息
	tockken, err := j.ParseToken(token)
	if err != nil {
		if err == TokenExpired {
			ctx.WriteString(`{"code": "401","msg": "未经允许的访问，已拦截R！"}`)
			return
		}
		ctx.WriteString(`{"code": "401","msg": "数据格式不正确"}`)
		return
	}
	o := models.NewOrm()
	u := models.User{
		UserId: tockken.UserID,
	}
	_, nicke_name, err_ping_user := models.PingUserAdmin(token, tockken.UserID)
	if err_ping_user != nil {
		ctx.WriteString(`{"code": "500","msg": "` + err_ping_user.Error() + `"}`)
		return
	}
	err_read := o.Read(&u, "user_id")
	if err_read != nil && err_read.Error() == "<QuerySeter> no row found" && tockken.NameSpace == "" {
		o.Begin()
		f, _, err_get_user := models.PingUser(token, tockken.UserID)
		if err_get_user != nil {
			ctx.WriteString(`{"code": "500","msg": "用户不存在!"}`)
			return
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

		_, erruser := o.Insert(&user)
		if erruser != nil {
			o.Rollback()
			ctx.WriteString(`{"code": "500","msg": "创建用户失败!"}`)
			return
		}
		account_def := models.Account{
			UserId:         tockken.UserID,
			CreateDate:     time.Now().Format("2006-01-02 15:04:05"),
			DynamicRevenue: true,
			StaticReturn:   true,
			PeerState:      true,
		}
		_, account_def_err := o.Insert(&account_def)
		if account_def_err != nil {
			o.Rollback()
			ctx.WriteString(`{"code": "500","msg": "创建生态仓库失败!"}`)
			return
		}
		formula := models.Formula{
			EcologyId:      account_def.Id,
			ReturnMultiple: 1,
		}
		_, err_for := o.Insert(&formula)
		if err_for != nil {
			o.Rollback()
			ctx.WriteString(`{"code": "500","msg": "创建算力表失败"}`)
			return
		}
		o.Commit()
	} else if err_read == nil && u.UserName != nicke_name.(string) && tockken.NameSpace == "" {
		o.Begin()
		u.UserName = nicke_name.(string)
		o.Update(&u, "user_name")
		o.Commit()
	} else if err_read != nil && err_read.Error() != "<QuerySeter> no row found" {
		ctx.WriteString(`{"code": "500","msg": "后端服务期错误"}`)
		return
	}
}

// 管理员用户的 过滤器处理规则
func AdminFilter(ctx *context.Context, api string) {
	token := ctx.Request.Header.Get("Authorization")

	if token == "" {
		ctx.WriteString(`{"code": "401","msg": "未经允许的访问，已拦截！"}`)
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
		ctx.WriteString(`{"code": "401","msg": "数据格式不正确"}`)
		return
	} else if tockken.NameSpace != "admin" {
		ctx.WriteString(`{"code": "401","msg": "未经允许的访问，已拦截！"}`)
		return
	}
	_, _, err_ping_user := models.PingUserAdmin(token, tockken.UserID)
	if err_ping_user != nil {
		ctx.WriteString(`{"code": "500","msg": "` + err_ping_user.Error() + `"}`)
		return
	}
}
