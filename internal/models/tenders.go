package models

import (
	"github.com/satori/uuid"
	"time"
)

type TypeService string

const (
	ServiceTypeConstruction TypeService = "Construction"
	ServiceTypeDelivery     TypeService = "Delivery"
	ServiceTypeManufacture  TypeService = "Manufacture"
)

type TypeStatus string

const (
	StatusCreated   TypeStatus = "Created"
	StatusPublished TypeStatus = "Published"
	StatusClosed    TypeStatus = "Closed"
)

type TendersRequest struct {
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	ServiceType     TypeService `json:"serviceType"`
	OrganizationId  uuid.UUID   `json:"organizationId"`
	CreatorUsername string      `json:"creatorUsername"`
}

type TendersResponse struct {
	Id             uuid.UUID   `json:"id"`
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	Status         TypeStatus  `json:"status"`
	ServiceType    TypeService `json:"serviceType"`
	Version        int         `json:"version"`
	CreatedAt      time.Time   `json:"createdAt"`
	OrganizationId uuid.UUID   `json:"organizationId"`
}

type TenderEditRequest struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	ServiceType TypeService `json:"serviceType"`
}
