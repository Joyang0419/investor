package server

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"protos/micro_auth"
	investor2 "repo/mongodb/investor"
	"repo/mongodb/schema"
)

type Auth struct {
	micro_auth.UnimplementedAuthServiceServer
	query   investor2.IQuery
	command investor2.ICommand
	timeout time.Duration
}

func NewAuth(
	query investor2.IQuery,
	command investor2.ICommand,
	timeout time.Duration,
) micro_auth.AuthServiceServer {
	return &Auth{
		command: command,
		query:   query,
		timeout: timeout,
	}
}

func (s *Auth) Login(ctx context.Context, request *micro_auth.LoginRequest) (response *micro_auth.LoginResponse, err error) {
	if request == nil {
		return nil,
			status.Errorf(
				codes.InvalidArgument,
				"[AuthServer][Login]request is nil",
			)
	}

	if _, err = s.command.Upsert(ctx, s.timeout, schema.Investor{
		ID:            request.ID,
		Email:         request.Email,
		VerifiedEmail: request.VerifiedEmail,
		Name:          request.Name,
		GivenName:     request.GivenName,
		FamilyName:    request.FamilyName,
		Picture:       request.Picture,
		Locale:        request.Locale,
	}); err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"[AuthServer][Login]Upsert err: %v", err,
		)
	}

	// TODO Token
	return &micro_auth.LoginResponse{Token: "TODO JWT TOKEN"}, nil
}
