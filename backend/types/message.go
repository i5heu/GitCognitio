package types

type Message struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Path string `json:"path"`
	Data string `json:"data"`
}
