package authgrpc

import (
	"context"
	"errors"
	ssov1 "github.com/wnikx/contracts/gen/go/sso"
	"github.com/wnikx/sso/internal/services/auth"
	"github.com/wnikx/sso/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const emptyValue = 0

type Auth interface {
	Login(ctx context.Context, email string, password string, appId int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userId int64, err error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type serverAPI struct {
	auth Auth
	ssov1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {

	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.AppId))

	if err != nil {

		// TODO
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func validateLogin(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "email or password is empty")
	}

	if req.AppId == emptyValue {
		return status.Error(codes.InvalidArgument, "appId is empty")
	}
	return nil
}

func (s *serverAPI) RegisterNewUser(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userId, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &ssov1.RegisterResponse{
		UserId: userId,
	}, nil
}

func validateRegister(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "email or password is empty")
	}
	return nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}
	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func validateIsAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "userId is empty")
	}
	return nil
}
