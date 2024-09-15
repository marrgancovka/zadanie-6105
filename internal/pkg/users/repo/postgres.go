package repo

import "database/sql"

type UserRepoPostgres struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *UserRepoPostgres {
	return &UserRepoPostgres{
		db: db,
	}
}

func (r *UserRepoPostgres) UserIsExists(username string) (bool, error) {
	query := `SELECT 1 FROM employee WHERE username = $1 LIMIT 1`

	var exists bool
	err := r.db.QueryRow(query, username).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
