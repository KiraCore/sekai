package ante_test

import (
	"encoding/hex"
	"fmt"

	"github.com/KiraCore/sekai/app/ante"
	ethmath "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	apitypes "github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func GenerateEIP712SignBytes() ([]byte, error) {
	types := apitypes.Types{
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
		},
		"delegate": {
			{Name: "param", Type: "string"},
		},
	}

	chainId := ethmath.NewHexOrDecimal256(int64(ante.EthChainID))

	domain := apitypes.TypedDataDomain{
		Name:    "Kira",
		Version: "1",
		ChainId: chainId,
		// VerifyingContract: "0xdef1c0ded9bec7f1a1670819833240f027b25eff",
	}

	typedData := apitypes.TypedData{
		Types:       types,
		PrimaryType: "delegate",
		Domain:      domain,
		Message: apitypes.TypedDataMessage{
			"param": `{"amount":"100000000ukex","to":"kiravaloper13j3w9pdc47e54z2gj4uh37rnnfwxcfcmjh4ful"}`,
		},
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, err
	}

	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, err
	}

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	hashBytes := crypto.Keccak256(rawData)

	return hashBytes, nil
}

func (suite *AnteTestSuite) TestGenerateEIP712SignBytes() {
	suite.SetupTest(false) // reset

	bytes, err := GenerateEIP712SignBytes()
	suite.Require().NoError(err)

	signatureData, err := hex.DecodeString("0a59681b3be1c26a71072989a43cca378ba2726f7183078865117715b9634f7d4c4c31fc008c411bb75739d44afd68a4f561b8176eb6cd38184a762ec83550e31b")
	suite.Require().NoError(err)
	signatureData[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1
	recovered, err := crypto.SigToPub(bytes, signatureData)
	suite.Require().NoError(err)
	recoveredAddr := crypto.PubkeyToAddress(*recovered)
	suite.Require().Equal(hex.EncodeToString(recoveredAddr[:]), "3f48fdb5ee16f729b16f4084ed1577557b6855cf")
}
