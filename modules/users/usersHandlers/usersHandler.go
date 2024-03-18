package usersHandlers

import (
	"github/Panyakorn4/kwanjai-shop-tutorial/config"
	"github/Panyakorn4/kwanjai-shop-tutorial/modules/entities"
	"github/Panyakorn4/kwanjai-shop-tutorial/modules/users"
	"github/Panyakorn4/kwanjai-shop-tutorial/modules/users/usersUsecases"
	"github/Panyakorn4/kwanjai-shop-tutorial/pkg/kwanjaiauth"

	"github.com/gofiber/fiber/v2"
)

type userHandlerErrCode string

const (
	signUpCustomerErr     userHandlerErrCode = "users-001"
	signInErr             userHandlerErrCode = "users-002"
	refreshPassportErr    userHandlerErrCode = "users-003"
	signOutErr            userHandlerErrCode = "users-004"
	singUpAdminErr        userHandlerErrCode = "users-005"
	generateTokenAdminErr userHandlerErrCode = "users-006"
)

type IUsersHandler interface {
	SignUpCustomer(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	RefreshPassport(c *fiber.Ctx) error
	SignOut(c *fiber.Ctx) error
	SignUpAdmin(c *fiber.Ctx) error
	GenerateAdminToken(c *fiber.Ctx) error
}

type usersHandler struct {
	cfg          config.IConfig
	usersUsecase usersUsecases.IUsersUsecase
}

func UsersHandler(cfg config.IConfig, usersUsecase usersUsecases.IUsersUsecase) IUsersHandler {
	return &usersHandler{
		cfg:          cfg,
		usersUsecase: usersUsecase,
	}
}
func (h *usersHandler) SignUpCustomer(c *fiber.Ctx) error {
	// Request body parser
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomerErr),
			err.Error(),
		).Res()
	}

	// Email validation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomerErr),
			"email pattern is invalid",
		).Res()
	}

	// Insert
	result, err := h.usersUsecase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		}
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}

func (h *usersHandler) SignIn(c *fiber.Ctx) error {
	req := new(users.UserCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code,
			string(signInErr),
			err.Error(),
		).Res()
	}

	passport, err := h.usersUsecase.GetPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code,
			string(signInErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) RefreshPassport(c *fiber.Ctx) error {
	req := new(users.UserRefreshCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(refreshPassportErr),
			err.Error(),
		).Res()
	}

	passport, err := h.usersUsecase.RefreshPassport(req)

	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(refreshPassportErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) SignOut(c *fiber.Ctx) error {
	req := new(users.UserRemoveCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signOutErr),
			err.Error(),
		).Res()
	}

	if err := h.usersUsecase.DeleteOauth(req.OauthId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signOutErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}

func (h *usersHandler) SignUpAdmin(c *fiber.Ctx) error {
	// Request body parser
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomerErr),
			err.Error(),
		).Res()
	}

	// Email validation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomerErr),
			"email pattern is invalid",
		).Res()
	}
	// Insert
	result, err := h.usersUsecase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		}
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}

func (h *usersHandler) GenerateAdminToken(c *fiber.Ctx) error {
	adminToken, err := kwanjaiauth.NewKwanjaiAuth(
		kwanjaiauth.Admin,
		h.cfg.Jwt(),
		nil,
	)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(generateTokenAdminErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK, &struct {
			Token string `json:"token"`
		}{
			Token: adminToken.SignToken(),
		},
	).Res()
}
