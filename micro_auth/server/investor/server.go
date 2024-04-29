package investor

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

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
	if isDuplicate {
		return nil, fmt.Errorf("[Server][CreateInvestor]login account is duplicate")
	}

	insertIDs, err := s.command.InsertMany(ctx, s.timeout, readyToCreate)
	if errorx.CheckErrorExist(err) {
		return nil, fmt.Errorf("[Server][CreateInvestor]InsertMany err: %w", err)
	}

	//if slicex.CheckLengthFitExpected[](insertIDs.InsertedIDs, 1) {
	//	return nil, fmt.Errorf("[Server][CreateInvestor]InsertMany length err")
	//}

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
