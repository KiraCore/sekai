package functions

import (
	middleware "github.com/KiraCore/sekai/middleware"
	sekaitypes "github.com/KiraCore/sekai/types"
)

// GetAllFunctions is a function to get all functions registered
func GetAllFunctions() sekaitypes.FunctionList {
	return middleware.GetFunctionList()
}
