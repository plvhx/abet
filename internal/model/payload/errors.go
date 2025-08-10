package payload

type ResponseError[T any] struct {
    Error T `json:"error"`
}
