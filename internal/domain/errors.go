package domain

import "errors"

var ErrNotFound = errors.New("not found")
var ErrDuplicateCode = errors.New("duplicate code")
var ErrCheckPass = errors.New("check isn't pass")
