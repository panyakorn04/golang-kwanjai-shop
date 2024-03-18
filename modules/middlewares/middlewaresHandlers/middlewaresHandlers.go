package middlewaresHandlers

import (
	"github/Panyakorn4/kwanjai-shop-tutorial/config"
	"github/Panyakorn4/kwanjai-shop-tutorial/modules/entities"
	"github/Panyakorn4/kwanjai-shop-tutorial/modules/middlewares/middlewaresUsecases"
	"github/Panyakorn4/kwanjai-shop-tutorial/pkg/kwanjaiauth"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type middlewareHandlersErrCode string

const (
	routerCheckErr middlewareHandlersErrCode = "middlware-001"
	jwtAuthErr     middlewareHandlersErrCode = "middlware-002"
)

type IMiddlewaresHandlers interface {
	Core() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
}

type middlewaresHandlers struct {
	cfg                 config.IConfig
	middlewaresUsecases middlewaresUsecases.IMiddlewaresUsecases
}

func MiddlewaresHandlers(cfg config.IConfig, middlewaresUsecases middlewaresUsecases.IMiddlewaresUsecases) IMiddlewaresHandlers {
	return &middlewaresHandlers{
		cfg:                 cfg,
		middlewaresUsecases: middlewaresUsecases,
	}
}

func (h *middlewaresHandlers) Core() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}
func (h *middlewaresHandlers) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return entities.NewResponse(c).Error(
			fiber.ErrNotFound.Code,
			string(routerCheckErr),
			"router not found",
		).Res()
	}
}

func (h *middlewaresHandlers) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "02/01/2006",
		TimeZone:   "Bangkok/Asia",
	})
}

func (h *middlewaresHandlers) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		result, err := kwanjaiauth.ParseToken(h.cfg.Jwt(), token)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				err.Error(),
			).Res()
		}

		claims := result.Claims
		if !h.middlewaresUsecases.FindAccessToken(claims.Id, token) {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				"no permission to access",
			).Res()
		}

		// Set UserId
		c.Locals("userId", claims.Id)
		c.Locals("userRoleId", claims.RoleId)
		return c.Next()
	}
}
