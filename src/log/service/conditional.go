package log_service

import log_flag "github.com/Rfluid/insta-tools/src/log/flag"

func LogConditionally[T any](
	logFunction func(msg string, args ...[]T),
	msg string,
	args ...[]T,
) {
	if !log_flag.Logs || logFunction == nil {
		return
	}
	logFunction(msg, args...)
}
