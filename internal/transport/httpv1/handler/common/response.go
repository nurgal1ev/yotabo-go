package common

type Body[T any] struct {
	Data     T                   `json:"data"`
	Errors   map[string][]string `json:"errors"`
	Messages []string            `json:"messages"`
}

type HumaAPIResponse[T any] struct {
	Body Body[T]
}

func NewHumaResponse[T any](data T, messages ...string) *HumaAPIResponse[T] {
	return &HumaAPIResponse[T]{
		Body: Body[T]{
			Data:     data,
			Messages: messages,
			Errors:   nil,
		},
	}
}
