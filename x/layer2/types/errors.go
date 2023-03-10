package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// layer2 module errors
var (
	ErrInvalidDappBondDenom          = errors.Register(ModuleName, 1, "invalid dapp bond denom")
	ErrLowAmountToCreateDappProposal = errors.Register(ModuleName, 2, "low amount to create dapp proposal")
	ErrDappDoesNotExist              = errors.Register(ModuleName, 3, "dapp not found")
	ErrDappAlreadyExists             = errors.Register(ModuleName, 4, "dapp already exists")
	ErrMaxDappBondReached            = errors.Register(ModuleName, 5, "max dapp bond reached")
	ErrNotEnoughUserDappBond         = errors.Register(ModuleName, 6, "not enough user dapp bond")
	ErrUserDappBondDoesNotExist      = errors.Register(ModuleName, 7, "user dapp bond does not exist")
)
