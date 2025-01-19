package domain

type OrderStatus string

const (
	Created     OrderStatus = "created"
	Approved    OrderStatus = "approved"
	Disapproved OrderStatus = "disapproved"
)

func (o OrderStatus) String() string {
	return string(o)
}
