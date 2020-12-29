package teststaking

import (
	tmcrypto "github.com/tendermint/tendermint/crypto"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/KiraCore/sekai/x/staking/types"
)

// GetTmConsPubKey gets the validator's public key as a tmcrypto.PubKey.
func GetTmConsPubKey(v types.Validator) tmcrypto.PubKey {
	return v.GetConsPubKey()
}

// ToTmValidator casts an SDK validator to a tendermint type Validator.
func ToTmValidator(v types.Validator) (*tmtypes.Validator, error) {
	tmPk := GetTmConsPubKey(v)
	return tmtypes.NewValidator(tmPk, v.ConsensusPower()), nil
}

// ToTmValidators casts all validators to the corresponding tendermint type.
func ToTmValidators(v []types.Validator) ([]*tmtypes.Validator, error) {
	validators := make([]*tmtypes.Validator, len(v))
	var err error
	for i, val := range v {
		validators[i], err = ToTmValidator(val)
		if err != nil {
			return nil, err
		}
	}

	return validators, nil
}
