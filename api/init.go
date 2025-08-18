package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net"
	"net/http"
	"picp/config"
	"picp/logger"
	"picp/utils"
	"picp/web"
	"sync"
	"time"
)

var engine *gin.Engine
var loginToken map[string]int64
var loginTokenLock sync.RWMutex

const cookieTokenKey = "token"

func Run(server net.Listener) error {
	engine = gin.New()
	loginToken = make(map[string]int64)
	engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Output: logger.GetWitter(),
	}), gin.Recovery())
	go checkTokenExpire()
	initApi(engine.Group("/api", LoginMiddleware))
	engine.POST("/api/login", doLogin)
	web.Init(engine, isLogin, logout)
	return engine.RunListener(server)
}

func checkTokenExpire() {
	ticker := time.NewTicker(time.Hour)
	for {
		select {
		case <-ticker.C:
			loginTokenLock.Lock()
			for key, value := range loginToken {
				if time.Now().Unix()-value > int64(config.Common.GetCookieMaxAge()) {
					delete(loginToken, key)
				}
			}
			loginTokenLock.Unlock()
		}
	}
}

func isLogin(ctx *gin.Context) bool {
	if config.Common.User == "" || config.Common.Password == "" {
		return true
	}
	token, _ := ctx.Cookie(cookieTokenKey)
	if token != "" {
		loginTokenLock.RLock()
		if value, ok := loginToken[token]; ok {
			loginTokenLock.RUnlock()
			if time.Now().Unix() < value {
				return true
			} else {
				loginTokenLock.Lock()
				delete(loginToken, token)
				loginTokenLock.Unlock()
			}
		} else {
			loginTokenLock.RUnlock()
		}
	}
	return false
}

func logout(ctx *gin.Context) {
	token, _ := ctx.Cookie(cookieTokenKey)
	if token != "" {
		ctx.SetCookie(cookieTokenKey, "", -1, "/", "", false, true)
		loginTokenLock.Lock()
		defer loginTokenLock.Unlock()
		delete(loginToken, token)
	}
}

func LoginMiddleware(ctx *gin.Context) {
	if isLogin(ctx) {
		ctx.Next()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 401,
		"msg":  "请先登录",
	})
	ctx.Abort()
}

type LoginQuery struct {
	User     string `json:"user" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func generateToken() (string, error) {
	for {
		id, err := uuid.NewUUID()
		if err != nil {
			return "", err
		}
		token := id.String()
		if _, ok := loginToken[token]; !ok {
			loginToken[token] = time.Now().Unix() + int64(config.Common.GetCookieMaxAge())
			return token, nil
		}
	}
}

func doLogin(ctx *gin.Context) {
	if isLogin(ctx) {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "已登陆",
		})
		return
	}
	var query LoginQuery
	err := ctx.BindJSON(&query)
	if err != nil {
		replayError(ctx, err)
		return
	}
	loginTokenLock.Lock()
	defer loginTokenLock.Unlock()
	if query.User != config.Common.User || utils.Sha1Sum(query.Password) != config.Common.Password {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "用户名或密码错误",
		})
		return
	}
	token, err := generateToken()
	if err != nil {
		replayError(ctx, err)
		return
	}
	ctx.SetCookie(cookieTokenKey, token, config.Common.GetCookieMaxAge(), "/", "", false, true)
	replaySuccess(ctx, nil)
}

type LoginSetting struct {
	Password string `json:"password"`
	User     string `json:"user"`
	MaxAge   int    `json:"max_age"`
}

func getLoginSetting(ctx *gin.Context) {
	loginTokenLock.RLock()
	defer loginTokenLock.RUnlock()
	replaySuccess(ctx, LoginSetting{
		User:   config.Common.User,
		MaxAge: config.Common.CookieMaxAge,
	})
}
func setLoginSetting(ctx *gin.Context) {
	var query LoginSetting
	err := ctx.BindJSON(&query)
	if err != nil {
		replayError(ctx, err)
		return
	}
	if query.User == "" {
		query.Password = ""
	} else {
		query.Password = utils.Sha1Sum(query.Password)
	}
	loginTokenLock.Lock()
	defer loginTokenLock.Unlock()
	err = config.SetLoginSetting(query.User, query.Password, query.MaxAge)
	if err != nil {
		replayError(ctx, err)
		return
	}
	token, _ := ctx.Cookie(cookieTokenKey)
	for key := range loginToken {
		if key != token {
			delete(loginToken, key)
		}
	}
	replaySuccess(ctx, nil)
}
