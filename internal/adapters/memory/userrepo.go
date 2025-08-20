package memory

import (
	"strings"
	"sync"

	"userhub/internal/domain"
)

type UserRepo struct {
	mu     sync.RWMutex
	byID   map[string]*domain.User
	byMail map[string]*domain.User
}

func NewUserRepo() *UserRepo {
	return &UserRepo{byID: map[string]*domain.User{}, byMail: map[string]*domain.User{}}
}

func (r *UserRepo) Create(u *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	k := strings.ToLower(u.Email)
	if _, ok := r.byMail[k]; ok {
		return errDuplicateEmail{}
	}
	cp := *u
	r.byID[u.ID] = &cp
	r.byMail[k] = &cp
	return nil
}

func (r *UserRepo) FindByEmail(email string) (*domain.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.byMail[strings.ToLower(email)]
	if !ok {
		return nil, false
	}
	cp := *u
	return &cp, true
}

func (r *UserRepo) FindByID(id string) (*domain.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.byID[id]
	if !ok {
		return nil, false
	}
	cp := *u
	return &cp, true
}

type errDuplicateEmail struct{}

func (errDuplicateEmail) Error() string { return "email already exists" }
