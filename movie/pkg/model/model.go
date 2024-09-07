package model

import model "github.com/Sahas001/movieapp/metadata/pkg"

type MovieDetails struct {
	Rating   *float64       `json:"rating,omitempty"`
	Metadata model.Metadata `json:"metadata"`
}
