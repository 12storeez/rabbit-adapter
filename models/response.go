package models

type OkResponse struct {
	Ok bool `json:"ok"`
}

type ErrResponse struct {
	Err string `json:"error"`
}
