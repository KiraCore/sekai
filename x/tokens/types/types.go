package types

import (
	fmt "fmt"
	"math"
	"strconv"
)

// constants
var (
	RateDecimal            int64   = 9
	RateDecimalDenominator float64 = math.Pow10(9)
	RateMaximum            float64 = math.Pow10(10)
)

// NewTokenAlias generates a new token alias struct.
func NewTokenAlias(
	expiration uint32,
	enactment uint32,
	allowedVoteTypes []VoteType,
	symbol string,
	name string,
	icon string,
	decimals uint32,
	denoms []string,
	status ProposalStatus,
) *TokenAlias {
	return &TokenAlias{
		Expiration:       expiration,
		Enactment:        enactment,
		AllowedVoteTypes: allowedVoteTypes,
		Symbol:           symbol,
		Name:             name,
		Icon:             icon,
		Decimals:         decimals,
		Denoms:           denoms,
		Status:           status,
	}
}

// NewTokenRate generates a new token rate struct.
func NewTokenRate(
	denom string,
	rate uint64,
	feePayments bool,
) *TokenRate {
	return &TokenRate{
		Denom:       denom,
		Rate:        rate,
		FeePayments: feePayments,
	}
}

// ToHumanReadable returns human readable struct from raw TokenRate struct
func (m *TokenRate) ToHumanReadable() *TokenRateHumanReadable {
	return &TokenRateHumanReadable{
		Denom:       m.Denom,
		Rate:        fmt.Sprintf("%f", float64(m.Rate)/RateDecimalDenominator),
		FeePayments: m.FeePayments,
	}
}

// ToRaw returns raw TokenRate struct from human readable one
func (m *TokenRateHumanReadable) ToRaw() *TokenRate {
	rate, err := strconv.ParseFloat(m.Rate, 64)
	if err != nil {
		panic("invalid human readable token rate")
	}
	return &TokenRate{
		Denom:       m.Denom,
		Rate:        uint64(rate * RateDecimalDenominator),
		FeePayments: m.FeePayments,
	}
}
