package types

import (
	"testcase/internal/modules/user/entities"
	"testcase/package/securities"
)

type LoginResponse struct {
	User  entities.User        `json:"user"`
	Token securities.TokenPair `json:"token"`
}
