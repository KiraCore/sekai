package interx

import (
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/database"
	"github.com/KiraCore/sekai/INTERX/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	legacytx "github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterInterxFaucetRoutes registers faucet services.
func RegisterInterxFaucetRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(config.FaucetRequestURL, FaucetRequest(gwCosmosmux, rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", config.FaucetRequestURL, "This is an API for faucet service.", false)
}

func serveFaucetInfo(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	faucetInfo := types.FaucetAccountInfo{}
	faucetInfo.Address = config.Config.Faucet.Address
	faucetInfo.Balances = common.GetAccountBalances(gwCosmosmux, r.Clone(r.Context()), config.Config.Faucet.Address)

	return faucetInfo, nil, http.StatusOK
}

/**
 * Error Codes
 * 0 : InternalServerError
 * 1 : Fail to send tokens
 * 100: Invalid address
 * 101: Claim time left
 * 102: Invalid token
 * 103: No need to send tokens
 * 104: Can't send tokens, less than minimum amount
 * 105: Not enough tokens
 */
func serveFaucet(r *http.Request, gwCosmosmux *runtime.ServeMux, request types.InterxRequest, rpcAddr string, bech32addr string, token string) (interface{}, interface{}, int) {
	// check address
	faucetAccAddr, err := sdk.AccAddressFromBech32(config.Config.Faucet.Address)
	if err != nil {
		common.GetLogger().Error("[faucet] Invalid bech32addr: ", config.Config.Faucet.Address)
		return common.ServeError(0, "", fmt.Sprintf("internal server error: %s", err), http.StatusInternalServerError)
	}

	// check address
	claimAccAddr, err := sdk.AccAddressFromBech32(bech32addr)
	if err != nil {
		common.GetLogger().Error("[faucet] Invalid bech32addr: ", claimAccAddr)
		return common.ServeError(100, "", fmt.Sprintf("invalid address: %s", err), http.StatusBadRequest)
	}

	// check claim limit
	timeLeft := database.GetClaimTimeLeft(bech32addr)
	if timeLeft > 0 {
		common.GetLogger().Error("[faucet] Claim time left: ", timeLeft)
		return common.ServeError(101, "", fmt.Sprintf("claim limit: %d second(s) left", timeLeft), http.StatusBadRequest)
	}

	availableBalances := common.GetAccountBalances(gwCosmosmux, r.Clone(r.Context()), config.Config.Faucet.Address)
	claimBalances := common.GetAccountBalances(gwCosmosmux, r.Clone(r.Context()), bech32addr)

	availableAmount := new(big.Int)
	availableAmount.SetString("0", 10)
	fmt.Println(availableBalances)
	for _, balance := range availableBalances {
		if balance.Denom == token {
			availableAmount.SetString(balance.Amount, 10)
		}
	}

	claimAmount := new(big.Int)
	claimAmount.SetString("0", 10)
	for _, balance := range claimBalances {
		if balance.Denom == token {
			claimAmount.SetString(balance.Amount, 10)
		}
	}

	faucetAmount := new(big.Int)
	faucetAmountString, ok := config.Config.Faucet.FaucetAmounts[token] // X

	if !ok {
		common.GetLogger().Error("[faucet] Failed to get faucet amount from the configuration")
		return common.ServeError(102, "", "invalid token", http.StatusBadRequest)
	}

	faucetAmount.SetString(faucetAmountString, 10)

	faucetMininumAmount := new(big.Int)
	faucetMininumAmountString, ok := config.Config.Faucet.FaucetMinimumAmounts[token] // M

	if !ok {
		common.GetLogger().Error("[faucet] Failed to get faucet minimum amount from the configuration")
		return common.ServeError(102, "", "invalid token", http.StatusBadRequest)
	}

	faucetMininumAmount.SetString(faucetMininumAmountString, 10)

	coinStr, ok := config.Config.Faucet.FeeAmounts[token]

	if !ok {
		common.GetLogger().Error("[faucet] Failed to get fee amount from the configuration")
		return common.ServeError(102, "", "invalid token", http.StatusBadRequest)
	}

	feeAmount, err := sdk.ParseCoinNormalized(coinStr)

	if err != nil {
		common.GetLogger().Error("[faucet] Failed to parse fee amount from the configuration: ", coinStr)
		return common.ServeError(102, "", "invalid token", http.StatusBadRequest)
	}

	// common.GetLogger().Info("[faucet] Available amount: ", availableAmount)
	// common.GetLogger().Info("[faucet] Claim amount: ", claimAmount)
	// common.GetLogger().Info("[faucet] Faucet amount: ", faucetAmount)
	// common.GetLogger().Info("[faucet] Faucet minimum amount: ", faucetMininumAmount)

	if faucetAmount.Cmp(claimAmount) <= 0 {
		common.GetLogger().Error("[faucet] No need to send tokens: faucetAmount <= claimAmount")
		return common.ServeError(103, "", "no need to send tokens", http.StatusBadRequest)
	}

	claimingAmount := new(big.Int)
	claimingAmount.SetString("0", 10)
	claimingAmount = claimingAmount.Sub(faucetAmount, claimAmount)
	if claimingAmount.Cmp(faucetMininumAmount) <= 0 {
		common.GetLogger().Error("[faucet] Less than minimum amount: faucetAmount-claimAmount <= faucetMininumAmount")
		return common.ServeError(104, "", "can't send tokens, less than minimum amount", http.StatusBadRequest)
	}

	remainingAmount := new(big.Int)
	remainingAmount.SetString("0", 10)
	remainingAmount = remainingAmount.Sub(availableAmount, faucetMininumAmount)
	if claimingAmount.Cmp(remainingAmount) > 0 {
		common.GetLogger().Error("[faucet] Not enough tokens: faucetAmount-claimAmount > availableAmount-faucetMininumAmount")
		return common.ServeError(105, "", "not enough tokens", http.StatusBadRequest)
	}

	// GET AccountNumber and Sequence
	accountNumber, sequence := common.GetAccountNumberSequence(gwCosmosmux, r.Clone(r.Context()), config.Config.Faucet.Address)
	// common.GetLogger().Info("[faucet] accountNumber: ", accountNumber)
	// common.GetLogger().Info("[faucet] sequence: ", sequence)

	msgSend := &bank.MsgSend{
		FromAddress: faucetAccAddr.String(),
		ToAddress:   claimAccAddr.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin(token, sdk.NewIntFromBigInt(claimingAmount))),
	}

	msgs := []sdk.Msg{msgSend}
	fee := legacytx.NewStdFee(200000, sdk.NewCoins(feeAmount)) //Fee handling
	memo := "Faucet Transfer"

	sigs := make([]legacytx.StdSignature, 1)
	signBytes := legacytx.StdSignBytes(common.NodeStatus.Chainid, accountNumber, sequence, 0, fee, msgs, memo)

	sig, err := config.Config.Faucet.PrivKey.Sign(signBytes)
	if err != nil {
		common.GetLogger().Error("[faucet] Failed to sign transaction: ", err)
		panic(err)
	}

	sigs[0] = legacytx.StdSignature{PubKey: config.Config.Faucet.PubKey, Signature: sig}

	stdTx := legacytx.NewStdTx(msgs, fee, sigs, memo)

	txBuilder := config.EncodingCg.TxConfig.NewTxBuilder()
	err = txBuilder.SetMsgs(stdTx.GetMsgs()...)
	if err != nil {
		common.GetLogger().Error("[faucet] Failed to set tx msgs: ", err)
		return common.ServeError(1, "failed to set TX Msgs", err.Error(), http.StatusInternalServerError)
	}

	sigV2, err := stdTx.GetSignaturesV2()
	if err != nil {
		common.GetLogger().Error("[faucet] Failed to get SignatureV2: ", err)
		return common.ServeError(1, "failed to get SignaturesV2", err.Error(), http.StatusInternalServerError)
	}

	sigV2[0].Sequence = sequence

	err = txBuilder.SetSignatures(sigV2...)
	if err != nil {
		common.GetLogger().Error("[faucet] Failed to set SignatureV2: ", err)
		return common.ServeError(1, "failed to set Signatures", err.Error(), http.StatusInternalServerError)
	}

	txBuilder.SetMemo(stdTx.GetMemo())
	txBuilder.SetFeeAmount(stdTx.GetFee())
	txBuilder.SetGasLimit(stdTx.GetGas())

	txBytes, err := config.EncodingCg.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		common.GetLogger().Error("[faucet] Failed to get tx bytes: ", err)
		return common.ServeError(1, "failed to get TX bytes", err.Error(), http.StatusBadRequest)
	}

	// send tokens
	txHash, err := common.BroadcastTransaction(rpcAddr, txBytes)
	if err != nil {
		common.GetLogger().Error("[faucet] Failed to broadcast transaction: ", err)
		return common.ServeError(1, "", err.Error(), http.StatusInternalServerError)
	}

	// add new claim
	database.AddNewClaim(bech32addr, time.Now())

	type FaucetResponse struct {
		Hash string `json:"hash"`
	}
	return FaucetResponse{Hash: txHash}, nil, http.StatusOK
}

// FaucetRequest is a function to handle faucet service.
func FaucetRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		queries := r.URL.Query()
		claims := queries["claim"]
		tokens := queries["token"]

		if len(claims) == 0 && len(tokens) == 0 {
			// common.GetLogger().Info("[faucet] Entering faucet info")
			response.Response, response.Error, statusCode = serveFaucetInfo(r, gwCosmosmux)
		} else if len(claims) == 1 && len(tokens) == 1 {
			// common.GetLogger().Info("[faucet] Entering faucet: claim = ", claims[0], ", token = ", tokens[0])
			response.Response, response.Error, statusCode = serveFaucet(r, gwCosmosmux, request, rpcAddr, claims[0], tokens[0])
		} else {
			common.GetLogger().Error("[faucet] Invalid parameters")
			response.Response, response.Error, statusCode = common.ServeError(0, "", "invalid query parameters", http.StatusBadRequest)
		}

		common.WrapResponse(w, request, *response, statusCode, false)
	}
}
