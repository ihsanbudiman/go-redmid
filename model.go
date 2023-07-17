package redmid

type Model struct {
	Status int    `json:"status"`
	Data   []byte `json:"data"`
}
