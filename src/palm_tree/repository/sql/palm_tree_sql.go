package sql

import (
	"context"
	"database/sql"
	"time"

	"github.com/davidyunus/sawitpro-estate/src/domain"
	"github.com/davidyunus/sawitpro-estate/src/helper"
	"github.com/labstack/gommon/log"
)

type palmTreeLocationRepositorySql struct {
	conn    *sql.DB
	manager *helper.Manager
}

func NewPalmTreeRepositorySql(conn *sql.DB, manager *helper.Manager) domain.PalmTreeLocationRepository {
	return &palmTreeLocationRepositorySql{
		conn:    conn,
		manager: manager,
	}
}

func (p *palmTreeLocationRepositorySql) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.PalmTree, error) {
	rows, err := p.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	result := []domain.PalmTree{}
	for rows.Next() {
		palmTree := domain.PalmTree{}

		err = rows.Scan(
			&palmTree.Id,
			&palmTree.Uuid,
			&palmTree.X,
			&palmTree.Y,
			&palmTree.Height,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, palmTree)
	}

	return result, nil
}

func (p *palmTreeLocationRepositorySql) GetPalmTreesByUuid(ctx context.Context, id string) ([]domain.PalmTree, error) {
	result, err := p.fetch(ctx, QueryGetByUuid, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *palmTreeLocationRepositorySql) PlantPalmTree(ctx context.Context, id string, param *domain.PalmTree) error {
	var dbConn interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	} = p.conn

	tx, _ := ctx.Value(p.manager.GetKey()).(*sql.Tx)
	if tx != nil {
		dbConn = tx
	}

	_, err := dbConn.ExecContext(ctx, QueryPlantPalmTree,
		id,
		param.X,
		param.Y,
		param.Height,
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}
