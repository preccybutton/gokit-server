package Service

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"golang.org/x/time/rate"
	"strconv"
	"xuexi/util"
)

type UserRequest struct {
	Uid int `json:"uid"`
	Method string `json:"method"`
	Token string `json:"token"`
}

type UserResponse struct {
	Result string `json:"result"`
}

//验证中间件
func CheckTokenMiddleware() endpoint.Middleware{
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r := request.(UserRequest)
			uc := UserClaim{}
			getToken, err := jwt.ParseWithClaims(r.Token, &uc, func(token *jwt.Token) (i interface{}, e error) {
				return []byte(secKey), nil
			})
			if getToken != nil && getToken.Valid{	//判断是否合法
				newCtx := context.WithValue(ctx, "LoginUser", getToken.Claims.(*UserClaim).Uname)
				return next(newCtx, request)
			}else if ve, ok := err.(*jwt.ValidationError); ok{	//验证不通过
				if ve.Errors&jwt.ValidationErrorMalformed != 0{
					fmt.Println("错误的token")
				}else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0{
					fmt.Println("token过期或未启用")
				}else{
					fmt.Println("无法处理这个token", err)
				}
			}
			return nil, util.NewMyError(403, "error token")
		}
	}
}

//日志中间件
func UserServiceLogMiddleware(logger log.Logger) endpoint.Middleware{
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r := request.(UserRequest)
			logger.Log("method", r.Method, "event","get user","userid", r.Uid)
			return next(ctx, request)
		}
	}
}

//限流功能的中间件
func RateLimit(limit *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !limit.Allow(){ // 如果令牌桶不通过
				return nil, util.NewMyError(429, "too many")	//使用自己的实现error的struct，设置code为429，message为too many
			}
			return next(ctx, request)	//通过则继续执行原来的endpoint
		}
	}
}

func GenUserEndPoint(userService IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest)
		fmt.Println("当前登录用户名是",ctx.Value("LoginUser"))
		result := "nothing"
		if r.Method == "GET"{
			result = userService.GetName(r.Uid) + strconv.Itoa(util.ServicePort)
		}else if r.Method == "DELETE"{
			err := userService.DelUser(r.Uid)
			if err != nil{
				result = err.Error()
			}else{
				result = fmt.Sprintf("%d 用户删除成功", r.Uid)
			}
		}
		return UserResponse{ result}, nil
	}
}