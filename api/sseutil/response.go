package sseutil

import (
	"go-rest-api-boilerplate/types"
)

type SseRes struct {
	PkgCode types.SseStatus `json:"pkg_code"`
	PkgData interface{}     `json:"pkg_data"`
}

func NewSseRes(pkgCode types.SseStatus, pkgData interface{}) *SseRes {
	return &SseRes{
		PkgCode: pkgCode,
		PkgData: pkgData,
	}
}
