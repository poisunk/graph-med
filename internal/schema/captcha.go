package schema

type CaptchaResp struct {
	CaptchaKey  string `json:"captcha_key"`
	ImageBase64 string `json:"image_base64"`
	ThumbBase64 string `json:"thumb_base64"`
	ParentSize  int    `json:"parent_size"`
	ChildSize   int    `json:"child_size"`
}
