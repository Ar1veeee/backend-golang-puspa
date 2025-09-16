package services

import (
	"backend-golang/internal/constants"
	"backend-golang/internal/helpers"
)

type TokenService interface {
	GenerateAccessToken(userId, role string) (string, error)
}

type tokenService struct{}

func NewTokenService() TokenService {
	return &tokenService{}
}

func (s *tokenService) GenerateAccessToken(userId, role string) (string, error) {
	return helpers.GenerateToken(userId, constants.Role(role))
}
