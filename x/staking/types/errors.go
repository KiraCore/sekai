package types

import "fmt"

var ErrInvalidMonikerLength = fmt.Errorf("invalid moniker length (max 64 bytes)")
var ErrInvalidWebsiteLength = fmt.Errorf("invalid website length (max 64 bytes)")
var ErrInvalidSocialLength = fmt.Errorf("invalid social length (max 64 bytes)")
var ErrInvalidIdentityLength = fmt.Errorf("invalid identity length (max 64 bytes)")
