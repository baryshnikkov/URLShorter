package link

type CreateReq struct {
	URL string `json:"url" validate:"required,url"`
}

type CreateRes struct {
	URL  string `json:"url"`
	Hash string `json:"hash"`
}
