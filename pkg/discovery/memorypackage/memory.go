package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Sahas001/movieapp/pkg/discovery"
)

type (
	serviceName string
	instanceID  string
)

type Registry struct {
	sync.RWMutex
	serviceAddrs map[serviceName]map[instanceID]*serviceInstance
}

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

func NewRegistry() *Registry {
	return &Registry{serviceAddrs: map[serviceName]map[instanceID]*serviceInstance{}}
}

func (r *Registry) Register(ctx context.Context, ID string, Name string, hostPort string) error {
	r.Lock()

	defer r.Unlock()

	// could not use the server name and instance id as it is due to conflict so did some modification

	if _, ok := r.serviceAddrs[serviceName(Name)]; !ok {
		r.serviceAddrs[serviceName(Name)] = map[instanceID]*serviceInstance{}
	}
	r.serviceAddrs[serviceName(Name)][instanceID(ID)] = &serviceInstance{hostPort: hostPort, lastActive: time.Now()}
	return nil
}

func (r *Registry) Deregister(ctx context.Context, ID string, Name string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName(Name)]; !ok {
		return nil
	}
	delete(r.serviceAddrs[serviceName(Name)], instanceID(ID))
	return nil
}

func (r *Registry) ReportHealthyState(ctx context.Context, ID string, Name string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName(Name)]; !ok {
		return errors.New("service is not registered yet")
	}

	if _, ok := r.serviceAddrs[serviceName(Name)][instanceID(ID)]; !ok {
		return errors.New("service instance is not registered yet")
	}
	r.serviceAddrs[serviceName(Name)][instanceID(ID)].lastActive = time.Now()
	return nil
}

func (r *Registry) ServiceAddresses(ctx context.Context, Name string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()

	if len(r.serviceAddrs[serviceName(Name)]) == 0 {
		return nil, discovery.ErrNotFound
	}

	var res []string
	for _, i := range r.serviceAddrs[serviceName(Name)] {
		if i.lastActive.Before(time.Now().Add(-time.Second * 5)) {
			continue
		}
		res = append(res, i.hostPort)
	}
	return res, nil
}
