package link

type Response struct {
	Message string      `json:"message,omitempty"`
	Content interface{} `json:"content,omitempty"`
}
