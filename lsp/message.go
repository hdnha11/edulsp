package lsp

type Message struct {
	RPC string `json:"jsonrpc"`
}

type Request struct {
	Message
	ID     int    `json:"id"`
	Method string `json:"method"`

	// We will just specify the type of the params in all the Request types
	// Params ...
}

type Response struct {
	Message
	ID *int `json:"id,omitempty"`

	// Result
	// Error
}

type Notification struct {
	Message
	Method string `json:"method"`

	// Params ...
}
