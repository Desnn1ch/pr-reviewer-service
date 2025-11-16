package response

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Error struct {
	Error ErrorBody `json:"error"`
}
