package quake

import "errors"

var ErrNoMatchingType error = errors.New("No matching type")
var ErrCorruptedMessage error = errors.New("Corrupted message")
