package slicerepository

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/marcokz/movie-final/internal/movie"
)

func InitialMovies() map[int64]movie.Movie {
	return map[int64]movie.Movie{
		1: {
			ID:   1,
			Name: "Map film 1",
			Year: 2016,
		},
		2: {
			ID:   2,
			Name: "Map film 2",
			Year: 2019,
		},
		3: {
			ID:   3,
			Name: "Map film 3",
			Year: 2024,
		},
	}
}

type Storage struct {
	movies      map[int64]movie.Movie
	nextMovieID int64
	mu          sync.RWMutex //когда больше читателей, то лучше RWMutex
}

func New(m map[int64]movie.Movie) *Storage {
	return &Storage{movies: m, nextMovieID: 4}
}

func (s *Storage) CreateMovie(ctx context.Context, m movie.Movie) error {
	s.mu.Lock() // race condition будет ждать пока все не сделают Unlock
	defer s.mu.Unlock()

	for _, oneMap := range s.movies { //долгий поиск по всем фильмам
		if oneMap.Name == m.Name && oneMap.Year == m.Year {
			return errors.New("this movie is already exists")
		}
	}

	m.ID = s.nextMovieID
	atomic.AddInt64(&s.nextMovieID, 1) //инкрементация счетчика для ID
	s.movies[m.ID] = m

	return nil
}

func (s *Storage) GetMovies(ctx context.Context) []movie.Movie { //получение данных
	s.mu.RLock()
	defer s.mu.RUnlock()

	movies := make([]movie.Movie, 0, len(s.movies))
	for _, movie := range s.movies {
		movies = append(movies, movie)
	}
	return movies
}

func (s *Storage) GetMoviesByID(ctx context.Context, id int64) (movie.Movie, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if m, ok := s.movies[id]; ok {
		return m, nil
	}

	return movie.Movie{}, movie.ErrNotFound
}

func (s *Storage) UpdateMovie(ctx context.Context, m movie.Movie) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.movies[m.ID]; ok {
		s.movies[m.ID] = m
		return nil
	}

	return movie.ErrNotFound
	/*  // is it good?
	if _, exists := s.movies[m.ID]; !exists {
		return movie.ErrNotFound
	}
	s.movies[m.ID] = m
	return nil
	*/
}

func (s *Storage) DeleteMovieByID(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.movies[id]; ok {
		delete(s.movies, id)
		return nil
	}

	return movie.ErrNotFound
}
