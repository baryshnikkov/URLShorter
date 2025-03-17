package link

type createReq struct {
	URL string `json:"url" validate:"required,url"`
}

type createRes struct {
	URL  string `json:"url"`
	Hash string `json:"hash"`
}
