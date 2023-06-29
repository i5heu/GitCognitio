package types

type Message struct {
	ID         string `json:"id"`
	mustBeAuth bool
	Type       string `json:"type"`
	Path       string `json:"path"`
	MetaData   string `json:"metaData"`
	Data       string `json:"data"`
}
