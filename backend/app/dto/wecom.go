package dto

type SignatureOptions struct {
	Signature string `form:"msg_signature"`
	TimeStamp string `form:"timestamp"`
	Nonce     string `form:"nonce"`
	EchoStr   string `form:"echostr"`
}

type WecomConfig struct {
	Type           string `json:"type"`
	CorpID         string `json:"corp_id"`
	Token          string `json:"token"`
	EncodingAESKey string `json:"encoding_aes_key"`
	Secret         string `json:"secret"`
}
