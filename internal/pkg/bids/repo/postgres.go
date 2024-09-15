package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/satori/uuid"
	"zadanie-6105/internal/models"
	"zadanie-6105/internal/myErrors"
)

type BidRepoPostgres struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *BidRepoPostgres {
	return &BidRepoPostgres{
		db: db,
	}
}

func (r *BidRepoPostgres) CreateBid(bidData *models.BidRequest) (*models.BidResponse, error) {
	var tenderExists bool
	queryTender := `SELECT 1 FROM TENDER WHERE id = $1`
	if err := r.db.QueryRow(queryTender, bidData.TenderId).Scan(&tenderExists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myErrors.ErrTenderNotFound
		}
		return nil, myErrors.ErrBadRequest
	}

	query := `
		INSERT INTO bid (name, description, tender_id, author_type, author_id, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, description, status, tender_id, author_type, author_id, version, created_at
	`

	var bidResponse models.BidResponse
	err := r.db.QueryRow(
		query,
		bidData.Name, bidData.Description, bidData.TenderId, bidData.AuthorType, bidData.AuthorId, models.StatusCreated).Scan(
		&bidResponse.Id, &bidResponse.Name, &bidResponse.Description, &bidResponse.Status, &bidResponse.TenderId,
		&bidResponse.AuthorType, &bidResponse.AuthorId, &bidResponse.Version, &bidResponse.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &bidResponse, nil
}
func (r *BidRepoPostgres) SelectUserBids(limit, offset int32, username string) ([]*models.BidResponse, error) {
	userId, err := r.GetUserIdByUsername(username)
	if err != nil {
		return nil, err
	}
	query := `
		SELECT id, name, description, status, tender_id, author_type, author_id, version, created_at
		FROM bid
		WHERE author_id = $1
		ORDER BY name ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []*models.BidResponse
	for rows.Next() {
		var bid models.BidResponse
		err = rows.Scan(&bid.Id, &bid.Name, &bid.Description, &bid.Status, &bid.TenderId,
			&bid.AuthorType, &bid.AuthorId, &bid.Version, &bid.CreatedAt)
		if err != nil {
			return nil, err
		}
		bids = append(bids, &bid)
	}

	return bids, nil
}
func (r *BidRepoPostgres) SelectTenderBids(limit, offset int32, tenderId uuid.UUID, username string) ([]*models.BidResponse, error) {
	queryUser := `
			SELECT 1
			FROM organization_responsible AS o_r
			JOIN organization AS o ON o_r.organization_id = o.id
			JOIN tender AS t ON t.organization_id = o.id
			JOIN employee AS e ON e.id = o_r.user_id
			WHERE t.id = $1 AND e.username = $2
	`
	var isOrganizer bool
	err := r.db.QueryRow(queryUser, tenderId, username).Scan(&isOrganizer)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, myErrors.ErrForbidden
		}
		return nil, err
	}

	var tenderExists bool
	queryTender := `SELECT 1 FROM TENDER WHERE id = $1`
	if err = r.db.QueryRow(queryTender, tenderId).Scan(&tenderExists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myErrors.ErrTenderNotFound
		}
		return nil, myErrors.ErrBadRequest
	}

	query := `
		SELECT id, name, description, status, tender_id, author_type, author_id, version, created_at
		FROM bid
		WHERE tender_id = $1
		AND status = $2
		ORDER BY name ASC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(query, tenderId, models.StatusPublished, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []*models.BidResponse
	for rows.Next() {
		var bid models.BidResponse
		err = rows.Scan(&bid.Id, &bid.Name, &bid.Description, &bid.Status, &bid.TenderId,
			&bid.AuthorType, &bid.AuthorId, &bid.Version, &bid.CreatedAt)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, myErrors.ErrBidNotFound
			}
			return nil, err
		}
		bids = append(bids, &bid)
	}

	return bids, nil
}
func (r *BidRepoPostgres) SelectBidStatus(bidId uuid.UUID) (string, error) {
	query := `SELECT status FROM bid WHERE id = $1`

	var status string
	err := r.db.QueryRow(query, bidId).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", myErrors.ErrBidNotFound
		}
		return "", err
	}

	return status, nil
}
func (r *BidRepoPostgres) CheckBidAuthor(bidId uuid.UUID, username string) (bool, error) {
	authorId, err := r.GetUserIdByUsername(username)
	if err != nil {
		return false, err
	}
	query := `SELECT 1 FROM bid WHERE id = $1 AND author_id = $2`

	var exists bool
	err = r.db.QueryRow(query, bidId, authorId).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, myErrors.ErrForbidden
		}
		return false, err
	}

	return exists, nil
}
func (r *BidRepoPostgres) UpdateBidStatus(bidId uuid.UUID, status string) (*models.BidResponse, error) {
	query := `
		UPDATE bid
		SET status = $1
		WHERE id = $2
		RETURNING id, name, description, status, tender_id, author_type, author_id, version, created_at
	`

	var bid models.BidResponse
	err := r.db.QueryRow(query, status, bidId).Scan(&bid.Id, &bid.Name, &bid.Description, &bid.Status, &bid.TenderId,
		&bid.AuthorType, &bid.AuthorId, &bid.Version, &bid.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &bid, nil
}
func (r *BidRepoPostgres) UpdateBid(bidId uuid.UUID, editedData *models.BidEditRequest) (*models.BidResponse, error) {
	query := `UPDATE bid SET updated_at = CURRENT_TIMESTAMP, version = version + 1`

	var args []interface{}
	argCounter := 1
	if editedData.Name != "" {
		query += `, name = $` + fmt.Sprint(argCounter)
		args = append(args, editedData.Name)
		argCounter++
	}
	if editedData.Description != "" {
		query += `, description = $` + fmt.Sprint(argCounter)
		args = append(args, editedData.Description)
		argCounter++
	}

	query += ` WHERE id = $` + fmt.Sprint(argCounter)
	args = append(args, bidId)

	query += ` RETURNING id, name, description, status, tender_id, author_type, author_id, version, created_at`

	var bid models.BidResponse
	err := r.db.QueryRow(query, args...).Scan(&bid.Id, &bid.Name, &bid.Description, &bid.Status, &bid.TenderId,
		&bid.AuthorType, &bid.AuthorId, &bid.Version, &bid.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &bid, nil
}
func (r *BidRepoPostgres) SubmitDecision(bidId uuid.UUID, decision string, username string) (*models.BidResponse, error) {
	queryUser := `
			SELECT 1
			FROM organization_responsible AS o_r
			JOIN organization AS o ON o_r.organization_id = o.id
			JOIN tender AS t ON t.organization_id = o.id
			JOIN bid AS b ON t.id = b.tender_id
			JOIN employee AS e ON e.id = o_r.user_id
			WHERE b.id = $1 AND e.username = $2
	`
	var isOrganizer bool
	err := r.db.QueryRow(queryUser, bidId, username).Scan(&isOrganizer)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, myErrors.ErrForbidden
		}
		return nil, err
	}
	query := `
		UPDATE bid
		SET decision = $1, status = $2
		WHERE id = $3
		RETURNING id, name, description, status, tender_id, author_type, author_id, version, created_at
	`

	var bid models.BidResponse
	err = r.db.QueryRow(query, decision, models.StatusClosed, bidId).Scan(&bid.Id, &bid.Name, &bid.Description, &bid.Status, &bid.TenderId,
		&bid.AuthorType, &bid.AuthorId, &bid.Version, &bid.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &bid, nil
}

func (r *BidRepoPostgres) GetUserIdByUsername(username string) (uuid.UUID, error) {
	query := `SELECT id FROM employee WHERE username = $1`

	var userId uuid.UUID
	err := r.db.QueryRow(query, username).Scan(&userId)
	if err != nil {
		return uuid.Nil, err
	}

	return userId, nil
}
