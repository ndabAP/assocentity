package assocentity

import "errors"

const (
	unicodeapostrophe = 39
	uncodedash        = 45
	escapedapos       = '\''
	dashchar          = '-'
)

var (
	// ErrNoEntity occurs if no entity was found.
	ErrNoEntity = errors.New("no entity found")
)
