package web

type GetDataRequest struct {
	Search string `query:"search" validate:"omitempty"`
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
}
