package hubs

type Hub struct {
	Addr   string
	Serial string
}

func (h *Hub) endpoint() string {
	return "http://" + h.Addr
}
