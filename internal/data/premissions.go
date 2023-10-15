package data

import (
	"context"
	"database/sql"
	"time"
)

type Premissions []string

func (p Premissions) Include(code string) bool {
	for i := range p {
		if code == p[i] {
			return true
		}
	}

	return false
}

type PremissionModel struct {
	DB *sql.DB
}

func (m PremissionModel) GetAllForUser(userID int64) (Premissions, error) {
	query := `
		select premissions.code
		from premissions
		inner join users_premissions on users_premissions.premission_id = premission.id
		inner join users on users_premissions.user_id = users.id
		where users.id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var premissions Premissions

	for rows.Next() {
		var premission string
		err := rows.Scan(&premission)
		if err != nil {
			return nil, err
		}

		premissions = append(premissions, premission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return premissions, nil
}
