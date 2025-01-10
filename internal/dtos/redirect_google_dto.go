package dtos

type RequestRedirectGoogleAuth struct {
	Code  string `json:"code"`
	State string `json:"state"`
	Scope string `json:"scope"`
}
