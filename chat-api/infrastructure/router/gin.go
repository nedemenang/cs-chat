package router

import (
	"chat-api/adapter/api/action"
	"chat-api/adapter/presenter"
	"chat-api/adapter/services"
	"chat-api/adapter/validator"
	"chat-api/infrastructure/common"
	"chat-api/infrastructure/config"
	"chat-api/usecase"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"chat-api/adapter/logger"
	"chat-api/adapter/repository"
)

type ginEngine struct {
	router     *gin.Engine
	log        logger.Logger
	db         repository.NoSQL
	validator  validator.Validator
	port       Port
	ctxTimeout time.Duration
}

func newGinServer(
	log logger.Logger,
	db repository.NoSQL,
	validator validator.Validator,
	port Port,
	t time.Duration,
) *ginEngine {
	return &ginEngine{
		router:     gin.New(),
		log:        log,
		db:         db,
		validator:  validator,
		port:       port,
		ctxTimeout: t,
	}
}

func (g ginEngine) Listen() {
	gin.SetMode(gin.ReleaseMode)
	gin.Recovery()

	g.setAppHandlers(g.router)

	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		Addr:         fmt.Sprintf(":%d", g.port),
		Handler:      g.router,
	}

	g.log.WithFields(logger.Fields{"port": g.port}).Infof("Starting HTTP Server")
	if err := server.ListenAndServe(); err != nil {
		g.log.WithError(err).Fatalln("Error starting HTTP server")
	}
}

func (g ginEngine) setAppHandlers(router *gin.Engine) {
	router.Use(g.CORSMiddleware())

	router.GET("/health", g.healthcheck())

	v1 := router.Group("/v1")

	v1.POST("/channel", g.AuthenticationMiddleware(), g.buildCreateChannelAction())
	v1.POST("/message", g.buildCreateMessageAction())
	v1.GET("/channel/:id", g.AuthenticationMiddleware(), g.buildGetChannelByIdAction())
	v1.PUT("/channel/:id", g.AuthenticationMiddleware(), g.buildUpdateChannelStatusAction())
	v1.GET("/channel", g.AuthenticationMiddleware(), g.buildGetChannelsByQueryAction())

	v1.POST("/user", g.buildCreateUserAction())

	v1.GET("/user/:email", g.AuthenticationMiddleware(), g.buildGetUserByEmailAction())
	v1.POST("/user/login", g.buildLoginUserAction())

}

func (g ginEngine) healthcheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		action.HealthCheck(c.Writer, c.Request)
	}
}

// logging middleware
func (g ginEngine) LoggingMiddleWare() gin.HandlerFunc {

	return func(c *gin.Context) {
		start := time.Now()
		header := c.Request.Header
		_, ok := header["Authorization"]
		if ok {
			delete(header, "Authorization")
		}

		c.Next()
		g.log.WithFields(logger.Fields{
			"name":       "chat-api",
			"status":     c.Writer.Status(),
			"path":       c.Request.URL.Path,
			"method":     c.Request.Method,
			"header":     header,
			"ip":         c.ClientIP(),
			"latency":    time.Since(start).Milliseconds(),
			"user-agent": c.Request.UserAgent(),
		}).Infof("request handled")
	}
}

// CORSMiddleware handler for CORS
func (g ginEngine) CORSMiddleware() gin.HandlerFunc {
	allowedOrigins := common.GetEnv("ALLOWED_ORIGINS", "http://localhost:3000")
	return func(c *gin.Context) {
		origin := c.GetHeader("origin")
		// g.log.Debugf("Allowed Origins: %s, origin: %s", allowedOrigins, origin)
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == http.MethodOptions {
			if _, exists := common.Find(strings.Split(allowedOrigins, ","), origin); !exists {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Request.Header.Del("Origin")
		c.Next()
	}
}

func (g ginEngine) AuthenticationMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		cfg := config.GetConfig()
		if c.Request.Header["Authorization"] != nil {
			tokenSlice := strings.Split(c.Request.Header["Authorization"][0], " ")
			if len(tokenSlice) != 2 {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			token, err := jwt.Parse(strings.TrimSpace(tokenSlice[1]), func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return []byte(cfg.AccessSecret), nil
			})

			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			if token.Valid {
				c.Next()
			}
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

	}
}

func (g ginEngine) buildCreateMessageAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewCreateMessageInteractor(
				repository.NewChannelNoSQL(g.db),
				presenter.NewCreateMessagePresenter(),
				g.ctxTimeout,
			)

			act = action.NewCreateMessageAction(uc, g.log, g.validator)
		)
		act.Execute(c.Writer, c.Request)
	}
}

func (g ginEngine) buildCreateChannelAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewCreateChannelInteractor(
				repository.NewChannelNoSQL(g.db),
				presenter.NewCreateChannelPresenter(),
				g.ctxTimeout,
			)

			act = action.NewCreateChannelAction(uc, g.log, g.validator)
		)

		act.Execute(c.Writer, c.Request)
	}
}

func (g ginEngine) buildCreateUserAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewCreateUserInteractor(
				repository.NewUserNoSQL(g.db),
				services.NewAuthenticationUtility(g.log),
				presenter.NewCreateUserPresenter(),
				g.ctxTimeout,
			)

			act = action.NewCreateUserAction(uc, g.log, g.validator)
		)

		act.Execute(c.Writer, c.Request)
	}
}

func (g ginEngine) buildGetChannelByIdAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewGetChannelByIdInteractor(
				repository.NewChannelNoSQL(g.db),
				presenter.NewGetChannelByIdPresenter(),
				g.ctxTimeout,
			)
			act = action.NewGetChannelByIdAction(uc, g.log, g.validator)
		)

		q := c.Request.URL.Query()
		q.Add("id", c.Param("id"))
		c.Request.URL.RawQuery = q.Encode()

		act.Execute(c.Writer, c.Request)
	}
}

func (g ginEngine) buildGetChannelsByQueryAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewGetChannelByQueryInteractor(
				repository.NewChannelNoSQL(g.db),
				repository.NewUserNoSQL(g.db),
				presenter.NewGetChannelsByQueryPresenter(),
				g.ctxTimeout,
			)
			act = action.NewGetChannelsByQueryAction(uc, g.log)
		)

		q := c.Request.URL.Query()
		q.Add("sort", c.Query("sort"))
		q.Add("limit", c.Query("limit"))
		q.Add("page", c.Query("page"))
		q.Add("repEmail", c.Query("repEmail"))
		q.Add("userEmail", c.Query("userEmail"))
		q.Add("currentStatus", c.Query("currentStatus"))
		c.Request.URL.RawQuery = q.Encode()

		act.Execute(c.Writer, c.Request)
	}
}

func (g ginEngine) buildGetUserByEmailAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewUserByEmailInteractor(
				repository.NewUserNoSQL(g.db),
				presenter.NewGetUserByEmailPresenter(),
				g.ctxTimeout,
			)
			act = action.NewGetUserByEmailAction(uc, g.log, g.validator)
		)

		q := c.Request.URL.Query()
		q.Add("email", c.Param("email"))
		c.Request.URL.RawQuery = q.Encode()

		act.Execute(c.Writer, c.Request)
	}
}

func (g ginEngine) buildLoginUserAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewLoginUserInteractor(
				repository.NewUserNoSQL(g.db),
				services.NewAuthenticationUtility(g.log),
				presenter.NewLoginPresenter(),
				g.ctxTimeout,
			)

			act = action.NewLoginUserAction(uc, g.log, g.validator)
		)

		act.Execute(c.Writer, c.Request)
	}
}

func (g ginEngine) buildUpdateChannelStatusAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewUpdateChannelStatusInteractor(
				repository.NewChannelNoSQL(g.db),
				presenter.NewUpdateChannelStatusPresenter(),
				g.ctxTimeout,
			)
			act = action.NewUpdateChannelAction(uc, g.log, g.validator)
		)

		q := c.Request.URL.Query()
		q.Add("id", c.Param("id"))
		c.Request.URL.RawQuery = q.Encode()

		act.Execute(c.Writer, c.Request)
	}
}
