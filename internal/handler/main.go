package handler

import (
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo-contrib/pprof"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/do"

	"codebase/pkg/httpx-echo"

	storyW "codebase/internal/models"
)

type Config struct {
	Container *do.Injector
	Mode      string
	Origins   []string
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func New(cfg *Config) (http.Handler, error) {
	r := echo.New()
	r.Validator = &CustomValidator{validator: validator.New()}
	r.Pre(middleware.RemoveTrailingSlash())
	if cfg.Mode == "debug" {
		r.Debug = true
		pprof.Register(r)
	}

	r.JSONSerializer = httpx.SegmentJSONSerializer{}
	// r.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	// 	Format: "${time_rfc3339}\t${method}\t${uri}\t${status}\t${latency_human}\n",
	// }))
	r.Use(middleware.Recover())

	routesAPIv1 := r.Group("/api/v1")
	{
		cors := middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     cfg.Origins,
			AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
			AllowCredentials: true,
			MaxAge:           60 * 60,
		})
		routesAPIv1.Use(cors)

		groupStock := &GroupStock{cfg}
		{
			routesAPIv1.GET("/stocks/group/:group", groupStock.GetStocksByGroupHandler)
			routesAPIv1.GET("/stocks/exchange/:exchange", groupStock.GetStocksByExchangeHandler)
			routesAPIv1.GET("/stocks/info/:symbol", groupStock.GetStockInfoHandler)
		}

		routesAPIv1Protected := routesAPIv1.Group("/protected")
		{
			routesAPIv1Protected.Use(echojwt.WithConfig(echojwt.Config{
				SigningKey:  []byte(os.Getenv("JWT_SECRET")),
				TokenLookup: "header:Authorization:Bearer ",
				ContextKey:  storyW.ContextJWTKey,
			}))

			routesAPIv1Protected.Use(WithAutheticatedJWTTokenData())
		}
	}

	r.GET("", func(c echo.Context) error {
		return c.String(http.StatusOK, "👻️")
	})

	return r, nil
}
