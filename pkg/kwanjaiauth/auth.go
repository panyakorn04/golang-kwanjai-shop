package kwanjaiauth

import (
	"errors"
	"fmt"
	"github/Panyakorn4/kwanjai-shop-tutorial/config"
	"github/Panyakorn4/kwanjai-shop-tutorial/modules/users"
	"math"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	Admin   TokenType = "admin"
	ApiKey  TokenType = "apikey"
)

type kwanjaiAuth struct {
	mapClaims *kwanjaiMapClaims
	cfg       config.IJwtConfig
}

type kwanjaiAdmin struct {
	*kwanjaiAuth
}

type kwanjaiMapClaims struct {
	Claims *users.UserClaims `json:"claims"`
	jwt.RegisteredClaims
}
type IKwanjaiAuth interface {
	SignToken() string
}

type IKwanjaiAdmin interface {
	SignToken() string
}

func ParseToken(cfg config.IJwtConfig, tokenString string) (*kwanjaiMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &kwanjaiMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.SecretKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*kwanjaiMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}
func RepeatToken(cfg config.IJwtConfig, claims *users.UserClaims, exp int64) string {
	obj := &kwanjaiAuth{
		cfg: cfg,
		mapClaims: &kwanjaiMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "kwanjai-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeRepeatAdapter(exp),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
	return obj.SignToken()
}

func jwtTimeDurationCal(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func (a *kwanjaiAuth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.SecretKey())
	return ss
}
func (a *kwanjaiAdmin) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.AdminKey())
	return ss
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func NewKwanjaiAuth(tokenType TokenType, cfg config.IJwtConfig, claims *users.UserClaims) (IKwanjaiAuth, error) {
	switch tokenType {
	case Access:
		return newAccessToken(cfg, claims), nil
	case Refresh:
		return newRefreshToken(cfg, claims), nil
	case Admin:
		return newAdminToken(cfg), nil
	default:
		return nil, fmt.Errorf("unknown token type")
	}
}

func newAccessToken(cfg config.IJwtConfig, claims *users.UserClaims) IKwanjaiAuth {
	return &kwanjaiAuth{
		cfg: cfg,
		mapClaims: &kwanjaiMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "kwanjai-api",
				Subject:   "access-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.AccessExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newRefreshToken(cfg config.IJwtConfig, claims *users.UserClaims) IKwanjaiAuth {
	return &kwanjaiAuth{
		cfg: cfg,
		mapClaims: &kwanjaiMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "kwanjai-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.RefreshExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func ParseAdminToken(cfg config.IJwtConfig, tokenString string) (*kwanjaiMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &kwanjaiMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.AdminKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*kwanjaiMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

func newAdminToken(cfg config.IJwtConfig) IKwanjaiAuth {
	return &kwanjaiAdmin{
		kwanjaiAuth: &kwanjaiAuth{
			cfg: cfg,
			mapClaims: &kwanjaiMapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "kwanjai-api",
					Subject:   "admin-token",
					Audience:  []string{"admin"},
					ExpiresAt: jwtTimeDurationCal(300),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}
