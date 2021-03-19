package keeper

import (
	"github.com/KiraCore/sekai/x/tokens/types"
)

func addTokens(origin, addings []string) []string {
	for _, adding := range addings {
		index := types.FindTokenIndex(origin, adding)
		if index >= 0 {
			continue
		}
		origin = append(origin, adding) // add into the array
	}
	return origin
}

func removeTokens(origin, removings []string) []string {
	for _, removing := range removings {
		index := types.FindTokenIndex(origin, removing)
		if index < 0 {
			continue
		}
		// fast remove from array
		origin[index] = origin[len(origin)-1] // set last element to index
		origin = origin[:len(origin)-1]       // remove last element
	}
	return origin
}
