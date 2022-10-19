package microservices

import (
	"k8s.io/apimachinery/pkg/types"
	"sync"
)

type entry struct {
	info        Microservice
	newestPodID types.UID
}

type Registry struct {
	microservices map[Identity]entry
	lock          sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		microservices: make(map[Identity]entry),
		lock:          sync.RWMutex{},
	}
}

func (r *Registry) Upsert(info Microservice, podID types.UID) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.microservices[info.Identity] = entry{
		info:        info,
		newestPodID: podID,
	}
}

func (r *Registry) Delete(info Microservice, podID types.UID) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if entry := r.microservices[info.Identity]; entry.newestPodID == podID {
		delete(r.microservices, info.Identity)
	}
}

func (r *Registry) Clear() {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.microservices = make(map[Identity]entry)
}

func (r *Registry) Get(id Identity) (Microservice, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	entry, ok := r.microservices[id]
	return entry.info, ok
}

func (r *Registry) All() []Microservice {
	r.lock.RLock()
	defer r.lock.RUnlock()

	infos := make([]Microservice, 0, len(r.microservices))
	for _, entry := range r.microservices {
		infos = append(infos, entry.info)
	}
	return infos
}
