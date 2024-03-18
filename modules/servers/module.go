package servers

import (
	"github/Panyakorn4/kwanjai-shop-tutorial/modules/middlewares/middlewaresHandlers"
	"github/Panyakorn4/kwanjai-shop-tutorial/modules/middlewares/middlewaresRepositories"
	"github/Panyakorn4/kwanjai-shop-tutorial/modules/middlewares/middlewaresUsecases"
	monitorHandlers "github/Panyakorn4/kwanjai-shop-tutorial/modules/monitor/handlersHandlers"
	"github/Panyakorn4/kwanjai-shop-tutorial/modules/users/usersHandlers"
	"github/Panyakorn4/kwanjai-shop-tutorial/modules/users/usersRepositories"
	"github/Panyakorn4/kwanjai-shop-tutorial/modules/users/usersUsecases"

	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
}

type moduleFactory struct {
	r   fiber.Router
	s   *server
	mid middlewaresHandlers.IMiddlewaresHandlers
}

func InitModule(r fiber.Router, s *server, mid middlewaresHandlers.IMiddlewaresHandlers) IModuleFactory {
	return &moduleFactory{
		r:   r,
		s:   s,
		mid: mid,
	}
}

func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandlers {
	repository := middlewaresRepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecases(repository)
	return middlewaresHandlers.MiddlewaresHandlers(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.s.cfg)
	m.r.Get("/", handler.HealthCheck)

}

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UsersRepository(m.s.db)
	usecase := usersUsecases.UsersUsecase(m.s.cfg, repository)
	handler := usersHandlers.UsersHandler(m.s.cfg, usecase)

	router := m.r.Group("/users")
	router.Post("/signup", handler.SignUpCustomer)
	router.Post("/signin", handler.SignIn)
	router.Post("/refresh", handler.RefreshPassport)
	router.Post("/signout", handler.SignOut)
	router.Post("/signup-admin", handler.SignOut)

	router.Get("/secret", m.mid.JwtAuth(), handler.GenerateAdminToken)
}
