package response

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        string `json:"expires_in"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type Ernie8kResponse struct {
	ErrorCode        int    `json:"error_code"`
	ErrorMsg         string `json:"error_msg"`
	ID               string `json:"id"`
	Object           string `json:"object"`
	Created          int    `json:"created"`
	SentenceId       int    `json:"sentence_id"`
	IsEnd            bool   `json:"is_end"`
	IsTruncated      bool   `json:"is_truncated"`
	FinishReason     string `json:"finish_reason"`
	SearchInfo       string `json:"search_info"`
	Result           string `json:"result"`
	NeedClearHistory bool   `json:"need_clear_history"`
	Flag             int    `json:"flag"`
	BanRound         int    `json:"ban_round"`
}

type GoldResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		List []*struct {
			Type             string `json:"type"`
			Typename         string `json:"typename"`
			Price            string `json:"price"`
			OpeningPrice     string `json:"openingprice"`
			MaxPrice         string `json:"maxprice"`
			MinPrice         string `json:"minprice"`
			ChangePercent    string `json:"changepercent"`
			LastClosingPrice string `json:"lastclosingprice"`
			TradeAmount      string `json:"tradeamount"`
			UpdateTime       string `json:"updatetime"`
		} `json:"list"`
	} `json:"data"`
}

type NewsResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Channel string `json:"channel"`
		Num     int    `json:"num"`
		List    []struct {
			Title string `json:"title"`
		} `json:"list"`
	} `json:"data"`
}

type TxImgToImgResp struct {
	Response struct {
		ResultImage string `json:"ResultImage"`
		RequestId   string `json:"RequestId"`
	} `json:"response"`
}
