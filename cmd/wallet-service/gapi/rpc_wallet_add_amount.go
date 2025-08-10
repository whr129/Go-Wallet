package gapi

import (
	"context"

	db "github.com/whr129/go-wallet/cmd/wallet-service/db/sqlc"
	"github.com/whr129/go-wallet/cmd/wallet-service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertAccountToProto(account db.Account) *pb.Account {
	return &pb.Account{
		Id:        account.ID,
		UserId:    account.UserID,
		Balance:   account.Balance,
		Currency:  account.Currency,
		CreatedAt: timestamppb.New(account.CreatedAt),
		IsDeleted: account.IsDeleted,
	}
}

func (server Server) AddAcountbalance(ctx context.Context, req *pb.AddAccountBalanceRequest) (*pb.AddAccountBalanceResponse, error) {
	accountID := req.GetAccountId()
	amount := req.GetAmount()

	if accountID <= 0 || amount <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "account ID and amount must be provided and amount must be greater than zero")
	}

	// Call the store method to add balance
	arg := db.AddAccountBalanceParams{
		ID:     accountID,
		Amount: amount,
	}
	account, err := server.store.AddAccountBalance(ctx, arg)

	if err != nil {
		if db.ErrorCode(err) == db.ErrRecordNotFound.Error() {
			return nil, status.Errorf(codes.NotFound, "account with ID %d not found ", accountID)
		}
		return nil, status.Errorf(codes.Internal, "failed to add account balance: %v", err)
	}

	return &pb.AddAccountBalanceResponse{
		Account: convertAccountToProto(account),
	}, nil
}
