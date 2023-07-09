package utils

import "net/http"

type ParamsContextkey struct{}

func GetField(r *http.Request, index int) string {
	fields := r.Context().Value(ParamsContextkey{}).([]string)
	return fields[index]
}
