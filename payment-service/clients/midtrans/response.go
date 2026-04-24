package midtrans

type MidtransResponse struct {
	Code   int          `json:"code"`
	Status string       `json:"status"`
	Data   MidtransData `json:"data"`
}

type MidtransData struct {
	Token       string `json:"token"`
	RedirectUrl string `json:"redirect_url"`
}
