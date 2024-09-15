package repo

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/satori/uuid"
	"zadanie-6105/internal/models"
	"zadanie-6105/internal/myErrors"
)

type TenderRepoPostgres struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *TenderRepoPostgres {
	return &TenderRepoPostgres{
		db: db,
	}
}

func (r *TenderRepoPostgres) SelectTendersList(limit, offset int32, serviceType []string) ([]*models.TendersResponse, error) {
	var args []interface{}
	args = append(args, models.StatusPublished, limit, offset)

	query := `
        SELECT id, name, description, status, service_type, created_at, version, organization_id
        FROM tender
        WHERE status = $1`

	if len(serviceType) > 0 && len(serviceType) < 3 {
		query += " AND service_type = ANY($4)"
		args = append(args, pq.Array(serviceType))
	}

	query += `
        ORDER BY name ASC
        LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, myErrors.ErrBadRequest
	}
	defer rows.Close()

	var tenders []*models.TendersResponse
	for rows.Next() {
		var tender models.TendersResponse
		if err = rows.Scan(&tender.Id, &tender.Name, &tender.Description, &tender.Status, &tender.ServiceType, &tender.CreatedAt, &tender.Version, &tender.OrganizationId); err != nil {
			return nil, err
		}
		tenders = append(tenders, &tender)
	}

	return tenders, nil
}
func (r *TenderRepoPostgres) CheckUsernameOrganization(creatorUsername string, organizationId uuid.UUID) (bool, error) {
	query := `
        SELECT COUNT(*)
        FROM organization_responsible
        JOIN employee ON organization_responsible.user_id = employee.id
        WHERE employee.username = $1 AND organization_responsible.organization_id = $2`

	var count int
	err := r.db.QueryRow(query, creatorUsername, organizationId).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, myErrors.ErrForbidden
		}
		return false, myErrors.ErrBadRequest
	}
	return count > 0, nil
}
func (r *TenderRepoPostgres) CheckUsernameTender(username string, tenderId uuid.UUID) (bool, error) {
	query := `
        SELECT COUNT(*)
        FROM tender
        WHERE creator_username = $1 AND id = $2`

	var count int
	err := r.db.QueryRow(query, username, tenderId).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, myErrors.ErrForbidden
		}
		return false, myErrors.ErrBadRequest
	}
	return count > 0, nil
}
func (r *TenderRepoPostgres) CreateTender(tenderData *models.TendersRequest) (*models.TendersResponse, error) {
	query := `
        INSERT INTO tender (name, description, service_type, organization_id, creator_username, status)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, name, description, status, service_type, organization_id, created_at, version`

	var tender models.TendersResponse
	err := r.db.QueryRow(query, tenderData.Name, tenderData.Description, tenderData.ServiceType, tenderData.OrganizationId, tenderData.CreatorUsername, models.StatusCreated).Scan(
		&tender.Id, &tender.Name, &tender.Description, &tender.Status, &tender.ServiceType, &tender.OrganizationId, &tender.CreatedAt, &tender.Version)
	if err != nil {
		return nil, myErrors.ErrBadRequest
	}

	return &tender, nil
}
func (r *TenderRepoPostgres) SelectUserTenders(limit, offset int32, username string) ([]*models.TendersResponse, error) {
	query := `
        SELECT id, name, description, status, service_type, organization_id, created_at, version
        FROM tender
        WHERE creator_username = $1
        ORDER BY name ASC
        LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, username, limit, offset)
	if err != nil {
		return nil, myErrors.ErrBadRequest
	}
	defer rows.Close()

	var tenders []*models.TendersResponse
	for rows.Next() {
		var tender models.TendersResponse
		if err = rows.Scan(&tender.Id, &tender.Name, &tender.Description, &tender.Status, &tender.ServiceType, &tender.OrganizationId, &tender.CreatedAt, &tender.Version); err != nil {
			return nil, err
		}
		tenders = append(tenders, &tender)
	}

	return tenders, nil

}
func (r *TenderRepoPostgres) SelectTenderStatus(tenderId uuid.UUID) (string, error) {
	query := `SELECT status FROM tender WHERE id = $1`

	var status string
	err := r.db.QueryRow(query, tenderId).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", myErrors.ErrTenderNotFound
		}
		return "", myErrors.ErrBadRequest
	}
	return status, nil
}
func (r *TenderRepoPostgres) EditStatusTender(tenderId uuid.UUID, status string) (*models.TendersResponse, error) {
	query := `
        UPDATE tender
        SET status = $1, updated_at = CURRENT_TIMESTAMP
        WHERE id = $2
        RETURNING id, name, description, status, service_type, organization_id, created_at, version`

	var tender models.TendersResponse
	err := r.db.QueryRow(query, status, tenderId).Scan(
		&tender.Id, &tender.Name, &tender.Description, &tender.Status, &tender.ServiceType, &tender.OrganizationId, &tender.CreatedAt, &tender.Version)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, myErrors.ErrTenderNotFound
		}
		return nil, myErrors.ErrBadRequest
	}

	return &tender, nil
}
func (r *TenderRepoPostgres) EditTender(tenderId uuid.UUID, editedData *models.TenderEditRequest) (*models.TendersResponse, error) {
	query := `UPDATE tender SET `
	var args []interface{}
	argCounter := 1

	if editedData.Name != "" {
		query += `name = $` + fmt.Sprint(argCounter) + `, `
		args = append(args, editedData.Name)
		argCounter++
	}
	if editedData.Description != "" {
		query += `description = $` + fmt.Sprint(argCounter) + `, `
		args = append(args, editedData.Description)
		argCounter++
	}
	if editedData.ServiceType != "" {
		query += `service_type = $` + fmt.Sprint(argCounter) + `, `
		args = append(args, editedData.ServiceType)
		argCounter++
	}

	query += `version = version + 1, updated_at = CURRENT_TIMESTAMP WHERE id = $` + fmt.Sprint(argCounter)
	args = append(args, tenderId)

	query += ` RETURNING id, name, description, status, service_type, organization_id, created_at, version`

	var tender models.TendersResponse
	err := r.db.QueryRow(query, args...).Scan(
		&tender.Id, &tender.Name, &tender.Description, &tender.Status, &tender.ServiceType, &tender.OrganizationId, &tender.CreatedAt, &tender.Version)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, myErrors.ErrTenderNotFound
		}
		return nil, myErrors.ErrBadRequest
	}

	return &tender, nil
}
