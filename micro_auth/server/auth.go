package server

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"protos/micro_auth"
	investor2 "repo/mongodb/investor"
	"repo/mongodb/schema"
	"tools/encryption"
)

type TokenInfo struct {
	Email string
}

type Auth struct {
	micro_auth.UnimplementedAuthServiceServer
	query         investor2.IQuery
	command       investor2.ICommand
	timeout       time.Duration
	jwtEncryption encryption.JWTEncryption[TokenInfo]
}

func NewAuth(
	query investor2.IQuery,
	command investor2.ICommand,
	timeout time.Duration,
	jwtEncryption encryption.JWTEncryption[TokenInfo],
) micro_auth.AuthServiceServer {
	return &Auth{
		command:       command,
		query:         query,
		timeout:       timeout,
		jwtEncryption: jwtEncryption,
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
		ID:                 request.ID,
		Email:              request.Email,
		VerifiedEmail:      request.VerifiedEmail,
		Name:               request.Name,
		GivenName:          request.GivenName,
		FamilyName:         request.FamilyName,
		Picture:            request.Picture,
		Locale:             request.Locale,
		LastLoginTimestamp: request.LastLoginTimestamp,
	}); err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"[AuthServer][Login]Upsert err: %v", err,
		)
	}

	token, err := s.jwtEncryption.Encrypt(
		jwt.MapClaims{
			"ID":                 request.ID,
			"Email":              request.Email,
			"LastLoginTimestamp": request.LastLoginTimestamp,
		},
	)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"[AuthServer][Login]jwtEncryption.Encrypt err: %v", err,
		)
	}

	return &micro_auth.LoginResponse{Token: token}, nil
}

func (s *Auth) ValidateToken(ctx context.Context, request *micro_auth.ValidateTokenRequest) (response *micro_auth.ValidateTokenResponse, err error) {
	select {
	case <-ctx.Done():
		return nil,
			status.Errorf(
				codes.Canceled,
				"[AuthServer][ValidateToken]ctx is done",
			)
	default:
		if request == nil {
			return nil,
				status.Errorf(
					codes.InvalidArgument,
					"[AuthServer][ValidateToken]request is nil",
				)
		}

		valid, claims, errDecrypt := s.jwtEncryption.Decrypt(request.Token)
		if errDecrypt != nil {
			return nil, status.Errorf(
				codes.Unauthenticated,
				"[AuthServer][ValidateToken]jwtEncryption.Decrypt err: %v", errDecrypt,
			)
		}

		return &micro_auth.ValidateTokenResponse{
			Valid: valid,
			Email: claims.Email,
		}, nil
	}
}
