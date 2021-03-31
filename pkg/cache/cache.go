package cache

import (
	"errors"
	"time"

	"github.com/new-adventure-areolite/grpc-app-server/pd/fight"
	"github.com/patrickmn/go-cache"
)

// Store ...
type Store struct {
	cache *cache.Cache
}

// ErrNotFound ...
var ErrNotFound = errors.New("not found")

// Add ...
func (s *Store) Add(hero *fight.Hero) error {
	return s.cache.Add(hero.Name, hero, cache.NoExpiration)
}

// Update updates the cache, if not exist, add it.
func (s *Store) Update(hero *fight.Hero) error {
	_, err := s.Get(hero.Name)
	switch {
	case err == ErrNotFound:
		return s.Add(hero)

	case err != nil:
		return err

	default:
		break
	}

	return s.cache.Replace(hero.Name, hero, cache.NoExpiration)
}

// Get ...
func (s *Store) Get(name string) (*fight.Hero, error) {
	value, ok := s.cache.Get(name)
	if !ok {
		return nil, ErrNotFound
	}

	hero, ok := value.(*fight.Hero)
	if !ok {
		return nil, errors.New("invalid type")
	}

	return hero, nil
}

// List ...
func (s *Store) List() []*fight.Hero {
	items := s.cache.Items()
	var heros []*fight.Hero
	for name := range items {
		value := items[name]
		hero, ok := value.Object.(*fight.Hero)
		if !ok {
			continue
		}
		heros = append(heros, hero)
	}
	return heros
}

// HeroStore ...
var HeroStore = Store{
	cache: cache.New(5*time.Minute, 10*time.Minute),
}

// // Store ...
// var Store = cache.New(5*time.Minute, 10*time.Minute)

// // HeroList ...
// const HeroList = "hero_list"
