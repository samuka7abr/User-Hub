package app

import (
	"errors"
	"strings"
	"time"

	"userhub/internal/domain"
)

type UserRepo interface {
	Create(*domain.User) error
	FindByEmail(string) (*domain.User, bool)
	FindByID(string) (*domain.User, bool)
}

type ProfileRepo interface {
	Upsert(*domain.Profile) error
	FindByUserID(string) (*domain.Profile, bool)
}

type Hasher interface {
	Hash(string) (string, error)
	Verify(string, string) (bool, error)
}

type TokenMaker interface {
	Make(string, time.Duration) (string, error)
}

type Service struct {
	Repo     UserRepo
	Profiles ProfileRepo
	Hash     Hasher
	Token    TokenMaker
}

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrEmailExists = errors.New("email already exists")
var ErrInvalidInput = errors.New("invalid input")
var ErrNotFound = errors.New("not found")

func (s *Service) Signup(email, password string) (*domain.User, error) {
	if !validEmail(email) || len(password) < 8 {
		return nil, ErrInvalidInput
	}
	if u, ok := s.Repo.FindByEmail(email); ok && u != nil {
		return nil, ErrEmailExists
	}
	hash, err := s.Hash.Hash(password)
	if err != nil {
		return nil, err
	}
	u := &domain.User{ID: domain.NewID(), Email: strings.ToLower(email), PassHash: hash, CreatedAt: time.Now().UTC()}
	if err := s.Repo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) Login(email, password string) (string, string, error) {
	u, ok := s.Repo.FindByEmail(email)
	if !ok {
		return "", "", ErrInvalidCredentials
	}
	ok2, err := s.Hash.Verify(u.PassHash, password)
	if err != nil || !ok2 {
		return "", "", ErrInvalidCredentials
	}
	at, err := s.Token.Make(u.ID, 15*time.Minute)
	if err != nil {
		return "", "", err
	}
	rt, err := s.Token.Make(u.ID, 24*time.Hour)
	if err != nil {
		return "", "", err
	}
	return at, rt, nil
}

func (s *Service) GetUser(id string) (*domain.User, bool) {
	return s.Repo.FindByID(id)
}

func (s *Service) GetProfile(userID string) (*domain.Profile, error) {
	if u, ok := s.Repo.FindByID(userID); !ok || u == nil {
		return nil, ErrNotFound
	}
	if p, ok := s.Profiles.FindByUserID(userID); ok {
		return p, nil
	}
	return &domain.Profile{UserID: userID, Name: "", Bio: "", AvatarURL: "", UpdatedAt: time.Time{}}, nil
}

func (s *Service) UpdateProfile(userID, name, bio, avatar string) (*domain.Profile, error) {
	if u, ok := s.Repo.FindByID(userID); !ok || u == nil {
		return nil, ErrNotFound
	}
	p := &domain.Profile{UserID: userID, Name: name, Bio: bio, AvatarURL: avatar, UpdatedAt: time.Now().UTC()}
	if err := s.Profiles.Upsert(p); err != nil {
		return nil, err
	}
	out, _ := s.Profiles.FindByUserID(userID)
	return out, nil
}

func validEmail(e string) bool {
	return strings.Contains(e, "@") && len(e) >= 6 && len(e) <= 254
}
