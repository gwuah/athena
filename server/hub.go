package server

type Hub struct {
	electrons  map[string]*Electron
	disconnect chan *Electron
	connect    chan *Electron
}

func newHub() *Hub {
	return &Hub{
		disconnect: make(chan *Electron),
		connect:    make(chan *Electron),
		electrons:  make(map[string]*Electron),
	}
}

func (h *Hub) Init() {
	for {
		select {
		case electron := <-h.connect:
			h.electrons[electron.id] = electron
		case electron := <-h.disconnect:
			if _, ok := h.electrons[electron.id]; ok {
				delete(h.electrons, electron.id)
				close(electron.send)
			}
		}
	}
}
