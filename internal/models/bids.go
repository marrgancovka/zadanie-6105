package models

import (
	"github.com/satori/uuid"
	"time"
)

type TypeAuthor string

type TypeDecision string

const (
	Organization TypeAuthor = "Organization"
	User         TypeAuthor = "User"
)

type BidRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	TenderId    uuid.UUID  `json:"tenderId"`
	AuthorType  TypeAuthor `json:"authorType"`
	AuthorId    uuid.UUID  `json:"authorId"`
}

type BidResponse struct {
	Id          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      TypeStatus `json:"status"`
	TenderId    uuid.UUID  `json:"tenderId"`
	AuthorType  TypeAuthor `json:"authorType"`
	AuthorId    uuid.UUID  `json:"authorId"`
	Version     int        `json:"version"`
	CreatedAt   time.Time  `json:"createdAt"`
}

type BidEditRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
