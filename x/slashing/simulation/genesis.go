package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/KiraCore/sekai/x/slashing/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/simulation"
)

// Simulation parameter constants
const (
	SignedBlocksWindow       = "signed_blocks_window"
	MinSignedPerWindow       = "min_signed_per_window"
	DowntimeInactiveDuration = "downtime_inactive_duration"
)

// GenSignedBlocksWindow randomized SignedBlocksWindow
func GenSignedBlocksWindow(r *rand.Rand) int64 {
	return int64(simulation.RandIntBetween(r, 10, 1000))
}

// GenMinSignedPerWindow randomized MinSignedPerWindow
func GenMinSignedPerWindow(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(int64(r.Intn(10)), 1)
}

// GenDowntimeInactiveDuration randomized DowntimeInactiveDuration
func GenDowntimeInactiveDuration(r *rand.Rand) time.Duration {
	return time.Duration(simulation.RandIntBetween(r, 60, 60*60*24)) * time.Second
}

// RandomizedGenState generates a random GenesisState for slashing
func RandomizedGenState(simState *module.SimulationState) {
	var signedBlocksWindow int64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, SignedBlocksWindow, &signedBlocksWindow, simState.Rand,
		func(r *rand.Rand) { signedBlocksWindow = GenSignedBlocksWindow(r) },
	)

	var minSignedPerWindow sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MinSignedPerWindow, &minSignedPerWindow, simState.Rand,
		func(r *rand.Rand) { minSignedPerWindow = GenMinSignedPerWindow(r) },
	)

	var DowntimeInactiveDuration time.Duration
	simState.AppParams.GetOrGenerate(
		simState.Cdc, DowntimeInactiveDuration, &DowntimeInactiveDuration, simState.Rand,
		func(r *rand.Rand) { DowntimeInactiveDuration = GenDowntimeInactiveDuration(r) },
	)

	params := types.NewParams(
		signedBlocksWindow, minSignedPerWindow, DowntimeInactiveDuration,
	)

	slashingGenesis := types.NewGenesisState(params, []types.SigningInfo{}, []types.ValidatorMissedBlocks{})

	bz, err := json.MarshalIndent(&slashingGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated slashing parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(slashingGenesis)
}
