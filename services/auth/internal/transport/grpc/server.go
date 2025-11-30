package grpc

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/sigame/auth/internal/domain"
	"github.com/sigame/auth/internal/service"
	pb "github.com/sigame/auth/proto"
)

// Server implements gRPC server for auth service
type Server struct {
	pb.UnimplementedAuthServiceServer
	authService *service.AuthService
}

// NewServer creates a new gRPC server instance
func NewServer(authService *service.AuthService) *Server {
	return &Server{
		authService: authService,
	}
}

// ValidateToken validates a JWT access token
func (s *Server) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	if req.Token == "" {
		return &pb.ValidateTokenResponse{
			Valid: false,
			Error: "token is required",
		}, nil
	}

	claims, err := s.authService.ValidateAccessToken(ctx, req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid: false,
			Error: err.Error(),
		}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:     true,
		UserId:    claims.UserID.String(),
		Username:  claims.Username,
		AvatarUrl: "", // No avatar support yet
	}, nil
}

// GetUserInfo retrieves user information by ID
func (s *Server) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	if req.UserId == "" {
		return &pb.GetUserInfoResponse{
			Found: false,
			Error: "user_id is required",
		}, nil
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return &pb.GetUserInfoResponse{
			Found: false,
			Error: "invalid user ID format",
		}, nil
	}

	user, err := s.authService.GetUserInfo(ctx, userID)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return &pb.GetUserInfoResponse{
				Found: false,
				Error: "user not found",
			}, nil
		}
		return &pb.GetUserInfoResponse{
			Found: false,
			Error: err.Error(),
		}, nil
	}

	return &pb.GetUserInfoResponse{
		Found:     true,
		UserId:    user.ID.String(),
		Username:  user.Username,
		AvatarUrl: "", // No avatar support yet
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}
