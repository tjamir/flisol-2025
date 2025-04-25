package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/tjamir/flisol-2025/microblog/commons/auth"
	"github.com/tjamir/flisol-2025/microblog/user-service/internal/model"
	"github.com/tjamir/flisol-2025/microblog/user-service/internal/repository"
	"github.com/tjamir/flisol-2025/microblog/user-service/proto"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUserNotFound   = errors.New("usuário não encontrado")
	ErrInvalidPassword = errors.New("senha incorreta")
	ErrEmailExists    = errors.New("email já registrado")

	// Erros de operação
	ErrGettingUser = errors.New("erro ao buscar usuário")
	ErrCreatingUser = errors.New("erro ao criar usuário")
	ErrHashingPassword = errors.New("erro ao gerar hash da senha")
	ErrGeneratingToken = errors.New("erro ao gerar token")
	ErrCheckingEmail = errors.New("erro ao verificar email")
	ErrInvalidToken     = errors.New("token inválido")
)



type UserHandler struct {
	proto.UnimplementedUserServiceServer
	Repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

func (h *UserHandler) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	existingUser, err := h.Repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("%w: %v", ErrCheckingEmail, err).Error())
	}
	if existingUser != nil {
		return nil, status.Error(codes.AlreadyExists, ErrEmailExists.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("%w: %v", ErrHashingPassword, err).Error())
	}

	user := &model.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := h.Repo.CreateUser(user); err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("%w: %v", ErrCreatingUser, err).Error())
	}

	return &proto.RegisterResponse{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}


func (h *UserHandler) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	user, err := h.Repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("%w: %v", ErrGettingUser, err).Error())
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, ErrUserNotFound.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, ErrInvalidPassword.Error())
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("%w: %v", ErrGeneratingToken, err).Error())
	}

	return &proto.LoginResponse{Token: token}, nil
}


func (h *UserHandler) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	user, err := h.Repo.GetUserByID(req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("%w: %v", ErrGettingUser, err).Error())
	}
	
	if user == nil {
		return nil, status.Error(codes.NotFound, ErrUserNotFound.Error())
	}

	return &proto.GetUserResponse{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (h *UserHandler) ValidateToken(ctx context.Context, req *proto.ValidateTokenRequest) (*proto.ValidateTokenResponse, error) {
	userID, err := auth.ValidateToken(req.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, fmt.Errorf("%w: %v", ErrInvalidToken, err).Error())
	}

	return &proto.ValidateTokenResponse{
		UserId: userID,
	}, nil
}

