package zconf

import "errors"


var(
	ErrSetDefaultValue  = errors.New("set default value error")
	ErrOpenConfigFile   = errors.New("file does not exist or is corrupted")
	ErrInvalidConfigFile   = errors.New("invalid config content format")
	ErrUnmarshalConfig  = errors.New("unmarshal config error")
)