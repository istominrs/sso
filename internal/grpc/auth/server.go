package auth

import (
	"context"
	ssov1 "github.com/GolangLessons/protos/gen/go/sso"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emailField    = "email"
	passwordField = "password"
	appIDField    = "app_id"
	tokenField    = "token"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error)
	Logout(ctx context.Context, token string) (success bool, err error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
	// TODO: remake
	validator *validator.Validate
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{
		validator: validator.New(),
		auth:      auth,
	})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	if err := validateLogin(s.validator, req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		// TODO: ...
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(s.validator, req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		// TODO: ...
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	if err := validateIsAdmin(s.validator, req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		// TODO: ...
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func (s *serverAPI) Logout(ctx context.Context,
	req *ssov1.LogoutRequest,
) (*ssov1.LogoutResponse, error) {
	if err := validateLogout(s.validator, req); err != nil {
		return nil, err
	}

	success, err := s.auth.Logout(ctx, req.GetToken())
	if err != nil {
		// TODO: ...
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LogoutResponse{
		Success: success,
	}, nil
}

func validateLogin(validator *validator.Validate, req *ssov1.LoginRequest) error {
	if err := validator.Var(req.GetEmail(), "required,email"); err != nil {
		return validateError(err, emailField)
	}

	if err := validator.Var(req.GetPassword(), "required,min=8"); err != nil {
		return validateError(err, passwordField)
	}

	if err := validator.Var(req.GetAppId(), "required"); err != nil {
		return validateError(err, appIDField)
	}

	return nil
}

func validateRegister(validator *validator.Validate, req *ssov1.RegisterRequest) error {
	if err := validator.Var(req.GetEmail(), "required,email"); err != nil {
		return validateError(err, emailField)
	}

	if err := validator.Var(req.GetPassword(), "required,min=8"); err != nil {
		return validateError(err, passwordField)
	}

	return nil
}

func validateIsAdmin(validator *validator.Validate, req *ssov1.IsAdminRequest) error {
	if err := validator.Var(req.GetUserId(), "required"); err != nil {
		return validateError(err, appIDField)
	}

	return nil
}

func validateLogout(validator *validator.Validate, req *ssov1.LogoutRequest) error {
	if err := validator.Var(req.GetToken(), "required"); err != nil {
		return validateError(err, tokenField)
	}

	return nil
}

func validateError(errs error, field string) error {
	for _, err := range errs.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			return status.Errorf(codes.InvalidArgument, "%s is required", field)
		case "email":
			return status.Error(codes.InvalidArgument, "invalid email format")
		case "min":
			return status.Errorf(codes.InvalidArgument, "%s is too short", field)
		default:
			return status.Errorf(codes.InvalidArgument, "invalid %s", field)
		}
	}

	return nil
}
