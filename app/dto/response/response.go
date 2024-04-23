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

type ZhiHuTopicResponse struct {
	Data []struct {
		Target struct {
			ID    int    `json:"id"`
			Title string `json:"title"`
			URL   string
		} `json:"target"`
		DetailText string `json:"detail_text"`
	} `json:"data"`
}

type TxImgToImgResp struct {
	Response struct {
		ResultImage string `json:"ResultImage"`
		RequestId   string `json:"RequestId"`
	} `json:"response"`
}
