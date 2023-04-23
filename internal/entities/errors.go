package entities

import "errors"

var ErrNotFound = errors.New("link not found")
var ErrLinkExists error = errors.New("link exists already")
var ErrLinkNotExists error = errors.New("link does not exist")
var ErrStorage error = errors.New("storage error")
