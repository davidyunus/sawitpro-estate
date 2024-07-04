package sql

import (
	"context"
	"database/sql"

	"github.com/davidyunus/sawitpro-estate/src/domain"
	"github.com/davidyunus/sawitpro-estate/src/helper"
	"github.com/labstack/gommon/log"
)

type estateRepositorySql struct {
	conn    *sql.DB
	manager *helper.Manager
}

func NewEstateRepositorySql(conn *sql.DB, manager *helper.Manager) domain.EstateRepository {
	return &estateRepositorySql{
		conn:    conn,
		manager: manager,
	}
}

func (m *estateRepositorySql) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.Estate, error) {
	rows, err := m.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	result := []domain.Estate{}
	for rows.Next() {
		estate := domain.Estate{}

		err = rows.Scan(
			&estate.Uuid,
			&estate.Length,
			&estate.Width,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, estate)
	}

	return result, nil
}

func (e *estateRepositorySql) CreateEstate(ctx context.Context, param *domain.Estate) error {
	var dbConn interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	} = e.conn

	tx, _ := ctx.Value(e.manager.GetKey()).(*sql.Tx)
	if tx != nil {
		dbConn = tx
	}

	_, err := dbConn.ExecContext(ctx, QueryCreateEstate,
		param.Uuid,
		param.Length,
		param.Width,
		helper.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (e *estateRepositorySql) GetEstateByUuid(ctx context.Context, id string) (*domain.Estate, error) {
	result, err := e.fetch(ctx, QueryGetByUuid, id)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}
	return &result[0], nil
}
