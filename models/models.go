package models

import (
	"context"
	"database/sql"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

// データベース接続用のモデルを返す
func NewModels(db *sql.DB) DBModel {
	return DBModel{
		DB: db,
	}
}

type User struct {
	ID        int       `json:"id"`
	LineID	  string    `json:"line_id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type WeightRecord struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	WeightNum int    	`json:"weight_num"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *DBModel) GetOneUser(line_id string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u User

	query := `
		select
			id, line_id, created_at, updated_at
		from
			users
		where line_id = $1`

	row := m.DB.QueryRowContext(ctx, query, line_id)

	err := row.Scan(
		&u.ID,
		&u.LineID,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return u, err
	}
	return u, nil
}

// ユーザの追加
func (m *DBModel) AddUser(line_id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into users (line_id, created_at, updated_at)
		values ($1, $2, $3)`

	_, err := m.DB.ExecContext(ctx, stmt,
		line_id,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}
	return nil
}

// ユーザ情報の更新
func (m *DBModel) UpdateUser(u User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into users (first_name, last_name, email, password, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6)`

	_, err := m.DB.ExecContext(ctx, stmt,
		u.LineID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}
	return nil
}

// ユーザの追加
func (m *DBModel) AddWeightRecord(id int, weight float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into weight_histories (user_id, weight_num, created_at, updated_at)
		values ($1, $2, $3, $4)`

	_, err := m.DB.ExecContext(ctx, stmt,
		id,
		weight,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}
	return nil
}

func (m *DBModel) GetMinWeight(id int) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var maxData float64

	query := `
		select
			min(w.weight_num)
		from
			users u
			left join weight_histories w on (u.id = w.user_id)
		where u.id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&maxData,
	)
	if err != nil {
		return 0, err
	}

	return maxData, nil
}

//　最新の記録(体重)を取得する
func (m *DBModel) GetLatestWeight(id int) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var latestData float64

	query := `
		select
			w.weight_num
		from
			users u
			left join weight_histories w on (u.id = w.user_id)
		where u.id = $1
		order by w.created_at desc
		limit 1
		`

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&latestData ,
	)
	if err != nil {
		return 0, err
	}

	return latestData , nil
}
