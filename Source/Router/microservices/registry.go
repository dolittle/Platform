package microservices

import "sync"

type Registry struct {
	microservices map[Identity]Microservice
	lock          sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		microservices: make(map[Identity]Microservice),
		lock:          sync.RWMutex{},
	}
}

func (r *Registry) Upsert(info Microservice) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.microservices[info.Identity] = info
}

func (r *Registry) Delete(info Microservice) {
	r.lock.Lock()
	defer r.lock.Unlock()

	delete(r.microservices, info.Identity)
}

func (r *Registry) Clear() {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.microservices = make(map[Identity]Microservice)
}

func (r *Registry) Get(id Identity) (Microservice, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	info, ok := r.microservices[id]
	return info, ok
}

func (r *Registry) All() []Microservice {
	r.lock.RLock()
	defer r.lock.RUnlock()

	infos := make([]Microservice, 0, len(r.microservices))
	for _, info := range r.microservices {
		infos = append(infos, info)
	}
	return infos
}
