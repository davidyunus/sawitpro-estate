package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/davidyunus/sawitpro-estate/src/common"
	"github.com/davidyunus/sawitpro-estate/src/domain"
	mock_domain "github.com/davidyunus/sawitpro-estate/src/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewEstateUsecase(t *testing.T) {
	assert.NotNil(t, NewEstateUsecase(nil, nil))
}

func TestCreateEstate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	estateRepoMock := mock_domain.NewMockEstateRepository(ctrl)
	palmTreeLocationRepoMock := mock_domain.NewMockPalmTreeLocationRepository(ctrl)

	uc := &estateUsecase{
		estateRepo:           estateRepoMock,
		palmTreeLocationRepo: palmTreeLocationRepoMock,
	}

	type args struct {
		ctx   context.Context
		param *domain.Estate
	}
	tests := []struct {
		name       string
		args       args
		wantResult *domain.CreateEstateResponse
		wantErr    bool
		mock       func() func()
	}{
		{
			name: "success",
			args: args{
				ctx: ctx,
				param: &domain.Estate{
					Uuid:   common.UtUuid,
					Length: 5,
					Width:  5,
				},
			},
			wantResult: &domain.CreateEstateResponse{
				Id: common.UtUuid,
			},
			wantErr: false,
			mock: func() func() {
				tempGenerateUUID := generateUUID
				generateUUID = func() string {
					return common.UtUuid
				}

				estateRepoMock.EXPECT().CreateEstate(gomock.Any(), &domain.Estate{
					Uuid:   common.UtUuid,
					Length: 5,
					Width:  5,
				}).Return(nil)
				return func() {
					generateUUID = tempGenerateUUID
				}
			},
		},
		{
			name: "error exceed size estate",
			args: args{
				ctx: ctx,
				param: &domain.Estate{
					Uuid:   common.UtUuid,
					Length: 500,
					Width:  500,
				},
			},
			wantResult: nil,
			wantErr:    true,
			mock: func() func() {
				return func() {}
			},
		},
		{
			name: "error create estate",
			args: args{
				ctx: ctx,
				param: &domain.Estate{
					Uuid:   common.UtUuid,
					Length: 5,
					Width:  5,
				},
			},
			wantResult: nil,
			wantErr:    true,
			mock: func() func() {
				tempGenerateUUID := generateUUID
				generateUUID = func() string {
					return common.UtUuid
				}

				estateRepoMock.EXPECT().CreateEstate(gomock.Any(), &domain.Estate{
					Uuid:   common.UtUuid,
					Length: 5,
					Width:  5,
				}).Return(errors.New(common.UtSomeError))
				return func() {
					generateUUID = tempGenerateUUID
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := uc.CreateEstate(test.args.ctx, test.args.param)
			assert.Equal(t, test.wantErr, err != nil)
			assert.Equal(t, test.wantResult, got)
		})
	}
}

func TestPlantPalmTree(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	estateRepoMock := mock_domain.NewMockEstateRepository(ctrl)
	palmTreeLocationRepoMock := mock_domain.NewMockPalmTreeLocationRepository(ctrl)

	uc := &estateUsecase{
		estateRepo:           estateRepoMock,
		palmTreeLocationRepo: palmTreeLocationRepoMock,
	}

	type args struct {
		ctx   context.Context
		id    string
		param *domain.PalmTree
	}
	tests := []struct {
		name       string
		args       args
		wantResult *domain.PlantPalmTreeResponse
		wantErr    bool
		mock       func()
	}{
		{
			name: "success",
			args: args{
				ctx: ctx,
				id:  common.UtUuid,
				param: &domain.PalmTree{
					Uuid: common.UtUuid,
					X:    3,
					Y:    1,
				},
			},
			wantResult: &domain.PlantPalmTreeResponse{
				Id: common.UtUuid,
			},
			wantErr: false,
			mock: func() {
				estateRepoMock.EXPECT().GetEstateByUuid(gomock.Any(), common.UtUuid).Return(&domain.Estate{
					Uuid:   common.UtUuid,
					Length: 6,
					Width:  3,
				}, nil)

				palmTreeLocationRepoMock.EXPECT().GetPalmTreesByUuid(gomock.Any(), common.UtUuid).Return([]domain.PalmTree{
					{
						Uuid: common.UtUuid,
						X:    2,
						Y:    1,
					},
				}, nil)

				palmTreeLocationRepoMock.EXPECT().PlantPalmTree(gomock.Any(), common.UtUuid, &domain.PalmTree{
					Uuid: common.UtUuid,
					X:    3,
					Y:    1,
				}).Return(nil)
			},
		},
		{
			name: "error get estate",
			args: args{
				ctx: ctx,
				id:  common.UtUuid,
				param: &domain.PalmTree{
					Uuid: common.UtUuid,
					X:    3,
					Y:    1,
				},
			},
			wantResult: nil,
			wantErr:    true,
			mock: func() {
				estateRepoMock.EXPECT().GetEstateByUuid(gomock.Any(), common.UtUuid).Return(nil, errors.New(common.UtSomeError))
			},
		},
		{
			name: "error estate nil",
			args: args{
				ctx: ctx,
				id:  common.UtUuid,
				param: &domain.PalmTree{
					Uuid: common.UtUuid,
					X:    3,
					Y:    1,
				},
			},
			wantResult: nil,
			wantErr:    true,
			mock: func() {
				estateRepoMock.EXPECT().GetEstateByUuid(gomock.Any(), common.UtUuid).Return(nil, nil)
			},
		},
		{
			name: "error get palm trees",
			args: args{
				ctx: ctx,
				id:  common.UtUuid,
				param: &domain.PalmTree{
					Uuid: common.UtUuid,
					X:    3,
					Y:    1,
				},
			},
			wantResult: nil,
			wantErr:    true,
			mock: func() {
				estateRepoMock.EXPECT().GetEstateByUuid(gomock.Any(), common.UtUuid).Return(&domain.Estate{
					Uuid:   common.UtUuid,
					Length: 6,
					Width:  3,
				}, nil)

				palmTreeLocationRepoMock.EXPECT().GetPalmTreesByUuid(gomock.Any(), common.UtUuid).Return(nil, errors.New(common.UtSomeError))
			},
		},
		{
			name: "error location already filled",
			args: args{
				ctx: ctx,
				id:  common.UtUuid,
				param: &domain.PalmTree{
					Uuid: common.UtUuid,
					X:    3,
					Y:    1,
				},
			},
			wantResult: nil,
			wantErr:    true,
			mock: func() {
				estateRepoMock.EXPECT().GetEstateByUuid(gomock.Any(), common.UtUuid).Return(&domain.Estate{
					Uuid:   common.UtUuid,
					Length: 6,
					Width:  3,
				}, nil)

				palmTreeLocationRepoMock.EXPECT().GetPalmTreesByUuid(gomock.Any(), common.UtUuid).Return([]domain.PalmTree{
					{
						Uuid: common.UtUuid,
						X:    3,
						Y:    1,
					},
				}, nil)
			},
		},
		{
			name: "error plan palm tree",
			args: args{
				ctx: ctx,
				id:  common.UtUuid,
				param: &domain.PalmTree{
					Uuid: common.UtUuid,
					X:    3,
					Y:    1,
				},
			},
			wantResult: nil,
			wantErr:    true,
			mock: func() {
				estateRepoMock.EXPECT().GetEstateByUuid(gomock.Any(), common.UtUuid).Return(&domain.Estate{
					Uuid:   common.UtUuid,
					Length: 6,
					Width:  3,
				}, nil)

				palmTreeLocationRepoMock.EXPECT().GetPalmTreesByUuid(gomock.Any(), common.UtUuid).Return([]domain.PalmTree{
					{
						Uuid: common.UtUuid,
						X:    2,
						Y:    1,
					},
				}, nil)

				palmTreeLocationRepoMock.EXPECT().PlantPalmTree(gomock.Any(), common.UtUuid, &domain.PalmTree{
					Uuid: common.UtUuid,
					X:    3,
					Y:    1,
				}).Return(errors.New(common.UtSomeError))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := uc.PlantPalmTree(test.args.ctx, test.args.id, test.args.param)
			assert.Equal(t, test.wantErr, err != nil)
			assert.Equal(t, test.wantResult, got)
		})
	}
}

func TestGetTreeStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	estateRepoMock := mock_domain.NewMockEstateRepository(ctrl)
	palmTreeLocationRepoMock := mock_domain.NewMockPalmTreeLocationRepository(ctrl)

	uc := &estateUsecase{
		estateRepo:           estateRepoMock,
		palmTreeLocationRepo: palmTreeLocationRepoMock,
	}

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		wantResult *domain.GetTreeStatsResponse
		wantErr    bool
		mock       func()
	}{
		{
			name: "success",
			args: args{
				ctx: ctx,
				id:  common.UtUuid,
			},
			wantResult: &domain.GetTreeStatsResponse{
				Count:  3,
				Max:    10,
				Min:    5,
				Median: 7,
			},
			wantErr: false,
			mock: func() {
				estateRepoMock.EXPECT().GetEstateByUuid(gomock.Any(), common.UtUuid).Return(&domain.Estate{
					Uuid:   common.UtUuid,
					Length: 6,
					Width:  3,
				}, nil)

				palmTreeLocationRepoMock.EXPECT().GetPalmTreesByUuid(gomock.Any(), common.UtUuid).Return([]domain.PalmTree{
					{
						Uuid:   common.UtUuid,
						X:      2,
						Y:      1,
						Height: 10,
					},
					{
						Uuid:   common.UtUuid,
						X:      6,
						Y:      2,
						Height: 5,
					},
					{
						Uuid:   common.UtUuid,
						X:      4,
						Y:      2,
						Height: 7,
					},
				}, nil)
			},
		},
		{
			name: "success even tree count",
			args: args{
				ctx: ctx,
				id:  common.UtUuid,
			},
			wantResult: &domain.GetTreeStatsResponse{
				Count:  4,
				Max:    10,
				Min:    5,
				Median: 7,
			},
			wantErr: false,
			mock: func() {
				estateRepoMock.EXPECT().GetEstateByUuid(gomock.Any(), common.UtUuid).Return(&domain.Estate{
					Uuid:   common.UtUuid,
					Length: 6,
					Width:  3,
				}, nil)

				palmTreeLocationRepoMock.EXPECT().GetPalmTreesByUuid(gomock.Any(), common.UtUuid).Return([]domain.PalmTree{
					{
						Uuid:   common.UtUuid,
						X:      2,
						Y:      1,
						Height: 10,
					},
					{
						Uuid:   common.UtUuid,
						X:      6,
						Y:      2,
						Height: 5,
					},
					{
						Uuid:   common.UtUuid,
						X:      4,
						Y:      2,
						Height: 7,
					},
					{
						Uuid:   common.UtUuid,
						X:      4,
						Y:      2,
						Height: 7,
					},
				}, nil)
			},
		},
		{
			name: "error get estate",
			args: args{
				ctx: ctx,
				id:  common.UtUuid,
			},
			wantResult: nil,
			wantErr:    true,
			mock: func() {
				estateRepoMock.EXPECT().GetEstateByUuid(gomock.Any(), common.UtUuid).Return(nil, errors.New(common.UtSomeError))
			},
		},
		{
			name: "error estate not found",
			args: args{
				ctx: ctx,
				id:  common.UtUuid,
			},
			wantResult: nil,
			wantErr:    true,
			mock: func() {
				estateRepoMock.EXPECT().GetEstateByUuid(gomock.Any(), common.UtUuid).Return(nil, nil)
			},
		},
		{
			name: "error get palm trees",
			args: args{
				ctx: ctx,
				id:  common.UtUuid,
			},
			wantResult: nil,
			wantErr:    true,
			mock: func() {
				estateRepoMock.EXPECT().GetEstateByUuid(gomock.Any(), common.UtUuid).Return(&domain.Estate{
					Uuid:   common.UtUuid,
					Length: 6,
					Width:  3,
				}, nil)

				palmTreeLocationRepoMock.EXPECT().GetPalmTreesByUuid(gomock.Any(), common.UtUuid).Return(nil, errors.New(common.UtSomeError))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := uc.GetTreeStats(test.args.ctx, test.args.id)
			assert.Equal(t, test.wantErr, err != nil)
			assert.Equal(t, test.wantResult, got)
		})
	}
}

func TestGetDroneFlyingDistance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	estateRepoMock := mock_domain.NewMockEstateRepository(ctrl)
	palmTreeLocationRepoMock := mock_domain.NewMockPalmTreeLocationRepository(ctrl)

	uc := &estateUsecase{
		estateRepo:           estateRepoMock,
		palmTreeLocationRepo: palmTreeLocationRepoMock,
	}

	type args struct {
		ctx         context.Context
		id          string
		maxDistance int
	}
	tests := []struct {
		name       string
		args       args
		wantResult *domain.GetDroneFlyingDistanceResponse
		wantErr    bool
		mock       func()
	}{
		{
			name: "success",
			args: args{
				ctx: ctx,
				id:  common.UtUuid,
			},
			wantResult: &domain.GetDroneFlyingDistanceResponse{
				Distance: 242,
			},
			wantErr: false,
			mock: func() {
				estateRepoMock.EXPECT().GetEstateByUuid(gomock.Any(), common.UtUuid).Return(&domain.Estate{
					Uuid:   common.UtUuid,
					Length: 6,
					Width:  3,
				}, nil)

				palmTreeLocationRepoMock.EXPECT().GetPalmTreesByUuid(gomock.Any(), common.UtUuid).Return([]domain.PalmTree{
					{
						Uuid:   common.UtUuid,
						X:      3,
						Y:      1,
						Height: 10,
					},
					{
						Uuid:   common.UtUuid,
						X:      6,
						Y:      2,
						Height: 5,
					},
					{
						Uuid:   common.UtUuid,
						X:      4,
						Y:      2,
						Height: 7,
					},
					{
						Uuid:   common.UtUuid,
						X:      3,
						Y:      2,
						Height: 15,
					},
					{
						Uuid:   common.UtUuid,
						X:      5,
						Y:      3,
						Height: 30,
					},
				}, nil)
			},
		},
		{
			name: "success with distance",
			args: args{
				ctx:         ctx,
				id:          common.UtUuid,
				maxDistance: 100,
			},
			wantResult: &domain.GetDroneFlyingDistanceResponse{
				Distance: 100,
				Rest: &domain.Rest{
					X: 4,
					Y: 2,
				},
			},
			wantErr: false,
			mock: func() {
				estateRepoMock.EXPECT().GetEstateByUuid(gomock.Any(), common.UtUuid).Return(&domain.Estate{
					Uuid:   common.UtUuid,
					Length: 6,
					Width:  3,
				}, nil)

				palmTreeLocationRepoMock.EXPECT().GetPalmTreesByUuid(gomock.Any(), common.UtUuid).Return([]domain.PalmTree{
					{
						Uuid:   common.UtUuid,
						X:      3,
						Y:      1,
						Height: 10,
					},
					{
						Uuid:   common.UtUuid,
						X:      6,
						Y:      2,
						Height: 5,
					},
					{
						Uuid:   common.UtUuid,
						X:      4,
						Y:      2,
						Height: 7,
					},
					{
						Uuid:   common.UtUuid,
						X:      3,
						Y:      2,
						Height: 15,
					},
					{
						Uuid:   common.UtUuid,
						X:      5,
						Y:      3,
						Height: 30,
					},
				}, nil)
			},
		},
		{
			name: "error get estate",
			args: args{
				ctx:         ctx,
				id:          common.UtUuid,
				maxDistance: 100,
			},
			wantResult: nil,
			wantErr:    true,
			mock: func() {
				estateRepoMock.EXPECT().GetEstateByUuid(gomock.Any(), common.UtUuid).Return(nil, errors.New(common.UtSomeError))
			},
		},
		{
			name: "error get estate nil",
			args: args{
				ctx:         ctx,
				id:          common.UtUuid,
				maxDistance: 100,
			},
			wantResult: nil,
			wantErr:    true,
			mock: func() {
				estateRepoMock.EXPECT().GetEstateByUuid(gomock.Any(), common.UtUuid).Return(nil, nil)
			},
		},
		{
			name: "error get palm trees",
			args: args{
				ctx: ctx,
				id:  common.UtUuid,
			},
			wantResult: nil,
			wantErr:    true,
			mock: func() {
				estateRepoMock.EXPECT().GetEstateByUuid(gomock.Any(), common.UtUuid).Return(&domain.Estate{
					Uuid:   common.UtUuid,
					Length: 6,
					Width:  3,
				}, nil)

				palmTreeLocationRepoMock.EXPECT().GetPalmTreesByUuid(gomock.Any(), common.UtUuid).Return(nil, errors.New(common.UtSomeError))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := uc.GetDroneFlyingDistance(test.args.ctx, test.args.id, test.args.maxDistance)
			assert.Equal(t, test.wantErr, err != nil)
			assert.Equal(t, test.wantResult, got)
		})
	}
}
