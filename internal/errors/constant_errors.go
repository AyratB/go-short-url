package customerrors

import "errors"

var ErrDuplicateEntity = errors.New("record with this original url already exists")
