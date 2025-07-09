package server

import (
	"context"
	"database/sql"
	"time"
	pb "user_service/src/internal/interfaces/output/grpc"

	"github.com/google/uuid"
)

type SessionServer struct {
	pb.UnimplementedSessionValidatorServer

	DB *sql.DB
}

func (s *SessionServer) ValidateSession(ctx context.Context, req *pb.ValidateSessionRequest) (*pb.ValidateSessionResponse, error) {
	sessionId := req.GetSessionId()
	uid, err := uuid.Parse(sessionId)
	if err != nil {
		return &pb.ValidateSessionResponse{
			Valid: false,
			Error: "invalid session id",
		}, err
	}
	var userId int
	var expiresAt time.Time
	err = s.DB.QueryRow("SELECT userid, expires_at from sessions where id=$1", uid).Scan(&userId, &expiresAt)
	if err != nil {
		return &pb.ValidateSessionResponse{
			Valid: false,
			Error: "session not found",
		}, err
	}
	if time.Now().After(expiresAt) {
		return &pb.ValidateSessionResponse{
			Valid: false,
			Error: "Session expired",
		}, nil
	}

	return &pb.ValidateSessionResponse{
		Valid:  true,
		UserId: string(rune(userId)),
		Error:  "",
	}, nil
}
