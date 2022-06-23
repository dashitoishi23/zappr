package commonmodels

type PagedResponse[T any] struct {
	Items []T   `json:"items"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
	Pages int   `json:"pages"`
	Err   error `json:"-"`
}

type PagedRequest[T any] struct {
	Entity T   `json:"entity"`
	Page   int `json:"page"`
	Size   int `json:"size"`
}

func (p PagedResponse[T]) Failed() error { return p.Err }
