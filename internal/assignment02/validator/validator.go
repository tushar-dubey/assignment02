package validator

import (
	"assignment02/rpc"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateGetRequest(req *rpc.GetVideosRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Skip,
			validation.Skip,
		),
		validation.Field(&req.Count,
			validation.Required,
			validation.Max(100),
		),
	)
}

func ValidateSearchRequest(req *rpc.SearchVideosRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Query,
			validation.Required,
		),
		validation.Field(&req.Skip,
			validation.Skip),
		validation.Field(&req.Count,
			validation.Required,
			validation.Max(100),
		),
	)
}
