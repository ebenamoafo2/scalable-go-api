package data

import "time"

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitzero"`
	Runtime   Runtime   `json:"runtime,omitzero"` //movie runtime in minutes
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}
