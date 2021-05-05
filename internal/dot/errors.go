package dot

import "errors"

var (
	ErrInvalidEdge   = errors.New("invalid edge")
	ErrAttrNameEmpty = errors.New("attr name empty")
	ErrNodeNameEmpty = errors.New("node name empty")
	ErrNoDotSource   = errors.New("no dot source")
	ErrMissingState  = errors.New("missing state")
)
