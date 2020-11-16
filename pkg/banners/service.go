package banners

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"
	"sync"
)

// Service type
type Service struct {
	mu    sync.RWMutex
	items []*Banner
	index int64
}

// NewService some
func NewService() *Service {
	return &Service{items: make([]*Banner, 0)}
}

// Banner some
type Banner struct {
	ID      int64
	Title   string
	Content string
	Button  string
	Link    string
	Image   string
}

// All returns everything
func (s *Service) All(ctx context.Context) ([]*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.items != nil {
		return s.items, nil
	}

	return nil, errors.New("No banners")
}

// ByID returns by id
func (s *Service) ByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, banner := range s.items {
		if banner.ID == id {
			return banner, nil
		}
	}

	return nil, errors.New("item not found")
}

// Save saves
func (s *Service) Save(ctx context.Context, item *Banner) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if item.ID == 0 {
		s.index++
		img := ""
		if item.Image != "" {
			im := strings.Split(item.Image, ".")
			img = strconv.FormatInt(s.index, 10) + "." + im[1]
		}
		newBanner := &Banner{
			ID:      s.index,
			Title:   item.Title,
			Content: item.Content,
			Button:  item.Button,
			Link:    item.Link,
			Image:   img,
		}
		s.items = append(s.items, newBanner)
		return newBanner, nil
	}
	sBanner, err := s.ByID(ctx, item.ID)
	if err != nil {
		log.Print(err)
		//http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return nil, errors.New("item not found")
	}
	sBanner.Button = item.Button
	sBanner.Title = item.Title
	sBanner.Content = item.Content
	sBanner.Link = item.Link
	if item.Image != "" {
		im := strings.Split(item.Image, ".")
		img := strconv.FormatInt(item.ID, 10) + "." + im[1]
		sBanner.Image = img
	}
	return sBanner, nil

}

// RemoveByID removes
func (s *Service) RemoveByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sBanner, err := s.ByID(ctx, id)
	if err != nil {
		log.Print(err)
		return nil, errors.New("item not found")
	}
	for i, banner := range s.items {
		if banner.ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			break
		}
	}

	return sBanner, nil
}
