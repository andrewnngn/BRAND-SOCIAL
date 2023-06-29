package dto

type PostDto struct {
	TextContent string `json:"text"`
	ImageURL    string `json:"image"`
	VideoURL    string `json:"video"`
}
