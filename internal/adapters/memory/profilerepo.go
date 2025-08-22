package memory

import (
	"sync"
	"time"

	"userhub/internal/domain"
)

type ProfileRepo struct {
	mu   sync.RWMutex
	data map[string]*domain.Profile
}

func NewProfileRepo() *ProfileRepo {
	return &ProfileRepo{data: map[string]*domain.Profile{}}
}

func (r *ProfileRepo) Upsert(p *domain.Profile) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *p
	cp.UpdatedAt = time.Now().UTC()
	r.data[p.UserID] = &cp
	return nil
}

func (r *ProfileRepo) FindByUserID(userID string) (*domain.Profile, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.data[userID]
	if !ok {
		return nil, false
	}
	cp := *p
	return &cp, true
}
