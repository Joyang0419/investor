package investor

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"protos/micro_auth"
	investor2 "repo/mongodb/investor"
	"repo/mongodb/schema/investor"
	"tools/errorx"
	"tools/slicex"
)

type Server struct {
	micro_auth.UnimplementedInvestorServiceServer
	query   investor2.IQuery
	command investor2.ICommand
	timeout time.Duration
}

func NewServer(
	query investor2.IQuery,
	command investor2.ICommand,
	timeout time.Duration,
) micro_auth.InvestorServiceServer {
	return &Server{
		command: command,
		query:   query,
		timeout: timeout,
	}
}

func (s *Server) GetInvestors(ctx context.Context, params *micro_auth.QueryInvestorsParams) (*micro_auth.InvestorsResponse, error) {
	investors, err := s.query.GetInvestors(
		ctx,
		s.timeout,
		investor2.GetInvestorsOptFilter{
			InvestorIDs:   params.Id,
			LoginAccounts: params.Username,
			Page:          params.Page,
			PageSize:      params.PageSize,
		},
	)
	if errorx.CheckErrorExist(err) {
		return nil, fmt.Errorf("[Server][GetInvestors]GetInvestors err: %w", err)
	}

	var investorsResponse micro_auth.InvestorsResponse
	for idx := range investors {
		investorsResponse.Investors = append(investorsResponse.Investors, &micro_auth.Investor{
			Id:       investors[idx].InvestorID,
			Username: investors[idx].LoginAccount,
			Password: investors[idx].Password,
		})
	}

	return &investorsResponse, nil
}

// TODO 調整, 要串接Google Oauth
func (s *Server) CreateInvestor(ctx context.Context, input *micro_auth.CreateInvestorInput) (*micro_auth.Investor, error) {
	readyToCreate := []investor.Schema{
		{
			LoginAccount: input.Username,
			Password:     input.Password,
		},
	}
	isDuplicate, err := s.checkLoginAccountDuplicate(ctx, input.Username)
	if errorx.CheckErrorExist(err) {
		return nil, fmt.Errorf("[Server][CreateInvestor]CheckLoginAccountDuplicate err: %w", err)
	}

	// TODO 記得參考這裡, 這是GRPC 操作error 的基操
	if isDuplicate {
		return nil, status.Errorf(
			codes.AlreadyExists,
			"[Server][CreateInvestor] Investor with username '%s' already exists",
			input.Username,
		)

		// TODO delete 我故意留著的， MicroErrOperationConflict 這也沒必要了
		// TODO 下面全部都要修status.Errorf(
		//return nil, errorx.New(
		//	"[Server][CreateInvestor]",
		//	error2.MicroErrOperationConflict,
		//	"",
		//)
	}

	insertIDs, err := s.command.InsertMany(ctx, s.timeout, readyToCreate)
	if errorx.CheckErrorExist(err) {
		return nil, fmt.Errorf("[Server][CreateInvestor]InsertMany err: %w", err)
	}

	typeAsserted, typeAssertOk := insertIDs.InsertedIDs[0].(primitive.ObjectID)
	if !typeAssertOk {
		return nil, fmt.Errorf("[Server][CreateInvestor]InsertMany type assert err")
	}

	return &micro_auth.Investor{
		Id:       typeAsserted.Hex(),
		Username: input.Username,
		Password: input.Password,
	}, nil
}

func (s *Server) checkLoginAccountDuplicate(ctx context.Context, loginAccount string) (bool, error) {
	investors, err := s.query.GetInvestors(ctx, s.timeout,
		investor2.GetInvestorsOptFilter{
			LoginAccounts: []string{loginAccount},
		},
	)
	if errorx.CheckErrorExist(err) {
		return false, fmt.Errorf("[Server][checkLoginAccountDuplicate]GetInvestors err: %w", err)
	}

	if slicex.IsNotEmpty(investors) {
		return true, nil
	}

	return false, nil
}
