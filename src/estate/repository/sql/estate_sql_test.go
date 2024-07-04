package sql

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/davidyunus/sawitpro-estate/src/common"
	"github.com/davidyunus/sawitpro-estate/src/domain"
	"github.com/davidyunus/sawitpro-estate/src/helper"
	"github.com/stretchr/testify/assert"
)

func init() {
	err := helper.InitTime()
	if err != nil {
		log.Fatal(err)
	}
}

func TestNewEstateRepositorySql(t *testing.T) {
	assert.NotNil(t, NewEstateRepositorySql(nil, nil))
}

func TestCreateEstate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	now := helper.Now()
	tempNow := helper.Now
	helper.Now = func() time.Time {
		return now
	}
	defer func() {
		helper.Now = tempNow
	}()
	repo := estateRepositorySql{
		conn:    db,
		manager: helper.NewManager(db, common.TransactionContextKey),
	}

	type args struct {
		ctx   context.Context
		param *domain.Estate
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mock    func()
	}{
		{
			name: "success",
			args: args{
				ctx: ctx,
				param: &domain.Estate{
					Uuid:   common.UtUuid,
					Length: 6,
					Width:  3,
				},
			},
			wantErr: false,
			mock: func() {
				mock.ExpectExec("INSERT").
					WithArgs("uuid", 6, 3, now).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "error",
			args: args{
				ctx: ctx,
				param: &domain.Estate{
					Uuid:   common.UtUuid,
					Length: 6,
					Width:  3,
				},
			},
			wantErr: true,
			mock: func() {
				mock.ExpectExec("INSERT").
					WithArgs("uuid", 6, 3, now).
					WillReturnError(errors.New(common.UtSomeError))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			err := repo.CreateEstate(test.args.ctx, test.args.param)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}

func TestGetEstateByUuid(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	repo := &estateRepositorySql{
		conn:    db,
		manager: helper.NewManager(db, common.TransactionContextKey),
	}

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		wantResult *domain.Estate
		wantErr    bool
		mock       func()
	}{
		{
			name: "success",
			args: args{
				ctx: ctx,
				id:  common.UtUuid,
			},
			wantResult: &domain.Estate{
				Uuid:   common.UtUuid,
				Length: 6,
				Width:  3,
			},
			wantErr: false,
			mock: func() {
				rows := sqlmock.NewRows([]string{"uuid", "length", "widty"}).
					AddRow(common.UtUuid, 6, 3)

				mock.ExpectQuery("SELECT").WithArgs(common.UtUuid).WillReturnRows(rows)
			},
		},
		{
			name: "error",
			args: args{
				ctx: ctx,
				id:  common.UtUuid,
			},
			wantResult: nil,
			wantErr:    true,
			mock: func() {
				mock.ExpectQuery("SELECT").WithArgs(common.UtUuid).WillReturnError(errors.New(common.UtSomeError))
			},
		},
		{
			name: "success no rows",
			args: args{
				ctx: ctx,
				id:  common.UtUuid,
			},
			wantResult: nil,
			wantErr:    false,
			mock: func() {
				rows := sqlmock.NewRows([]string{"uuid", "length", "widty"})

				mock.ExpectQuery("SELECT").WithArgs(common.UtUuid).WillReturnRows(rows)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := repo.GetEstateByUuid(test.args.ctx, test.args.id)
			assert.Equal(t, test.wantErr, err != nil)
			assert.Equal(t, test.wantResult, got)
		})
	}
}
