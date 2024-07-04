package http

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/davidyunus/sawitpro-estate/src/common"
	"github.com/davidyunus/sawitpro-estate/src/domain"
	mock_domain "github.com/davidyunus/sawitpro-estate/src/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewEstateHandler(t *testing.T) {
	NewEstateHandler(echo.New(), nil)
}

func TestCreateEstate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	estateMock := mock_domain.NewMockEstateUsecase(ctrl)
	handler := &estateHandler{
		estateUsecase: estateMock,
	}

	tests := []struct {
		name       string
		args       string
		wantResult string
		mock       func()
	}{
		{
			name: "success",
			args: `{"length":6,"width":3}`,
			wantResult: `{"code":201,"message":"Success create estate","data":{"id":"uuid"},"errors":null}
`,
			mock: func() {
				estateMock.EXPECT().CreateEstate(gomock.Any(), &domain.Estate{
					Length: 6,
					Width:  3,
				}).Return(&domain.CreateEstateResponse{
					Id: common.UtUuid,
				}, nil)
			},
		},
		{
			name: "error create estate",
			args: `{"length":6,"width":3}`,
			wantResult: `{"code":400,"message":"some error","data":null,"errors":"some error"}
`,
			mock: func() {
				estateMock.EXPECT().CreateEstate(gomock.Any(), &domain.Estate{
					Length: 6,
					Width:  3,
				}).Return(nil, errors.New(common.UtSomeError))
			},
		},
		{
			name: "error json decode",
			args: `{"length":"aaa","width":3}`,
			wantResult: `{"code":400,"message":"invalid input","data":null,"errors":"invalid input"}
`,
			mock: func() {},
		},
		{
			name: "error invalid input",
			args: `{"length":0,"width":3}`,
			wantResult: `{"code":400,"message":"invalid input","data":null,"errors":"invalid input"}
`,
			mock: func() {},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/estate", strings.NewReader(test.args))
			req.Header.Set(common.UtContentType, common.ContentTypeJson)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			test.mock()

			if assert.NoError(t, handler.CreateEstate(c)) {
				assert.Equal(t, test.wantResult, rec.Body.String())
			}
		})
	}
}

func TestPlantPalmTree(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	estateMock := mock_domain.NewMockEstateUsecase(ctrl)
	handler := &estateHandler{
		estateUsecase: estateMock,
	}

	tests := []struct {
		name       string
		args       string
		wantResult string
		mock       func()
	}{
		{
			name: "success",
			args: `{"x":3,"y":1,"height":10}`,
			wantResult: `{"code":201,"message":"Success plant palm tree","data":{"id":"uuid"},"errors":null}
`,
			mock: func() {
				estateMock.EXPECT().PlantPalmTree(gomock.Any(), common.UtUuid, &domain.PalmTree{
					X:      3,
					Y:      1,
					Height: 10,
				}).Return(&domain.PlantPalmTreeResponse{
					Id: common.UtUuid,
				}, nil)
			},
		},
		{
			name: "error plant palm tree",
			args: `{"x":3,"y":1,"height":10}`,
			wantResult: `{"code":404,"message":"estate not found","data":null,"errors":"estate not found"}
`,
			mock: func() {
				estateMock.EXPECT().PlantPalmTree(gomock.Any(), common.UtUuid, &domain.PalmTree{
					X:      3,
					Y:      1,
					Height: 10,
				}).Return(nil, domain.ErrEstateNotFound)
			},
		},
		{
			name: "error json decoder",
			args: `{"x":"aaa","y":1,"height":10}`,
			wantResult: `{"code":400,"message":"invalid input","data":null,"errors":"invalid input"}
`,
			mock: func() {},
		},
		{
			name: "error json decoder",
			args: `{"x":0,"y":1,"height":10}
`,
			wantResult: `{"code":400,"message":"invalid input","data":null,"errors":"invalid input"}
`,
			mock: func() {},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/estate/uuid/tree", strings.NewReader(test.args))
			req.Header.Set(common.UtContentType, common.ContentTypeJson)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues("uuid")

			test.mock()

			if assert.NoError(t, handler.PlantPalmTree(c)) {
				assert.Equal(t, test.wantResult, rec.Body.String())
			}
		})
	}
}

func TestGetTreeStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	estateMock := mock_domain.NewMockEstateUsecase(ctrl)
	handler := &estateHandler{
		estateUsecase: estateMock,
	}

	tests := []struct {
		name       string
		args       string
		wantResult string
		mock       func()
	}{
		{
			name: "success",
			args: common.UtUuid,
			wantResult: `{"code":200,"message":"Success get tree stats","data":{"count":3,"max":30,"min":5,"median":15},"errors":null}
`,
			mock: func() {
				estateMock.EXPECT().GetTreeStats(gomock.Any(), common.UtUuid).Return(&domain.GetTreeStatsResponse{
					Count:  3,
					Max:    30,
					Min:    5,
					Median: 15,
				}, nil)
			},
		},
		{
			name: "error get tree stats",
			args: common.UtUuid,
			wantResult: `{"code":404,"message":"estate not found","data":null,"errors":"estate not found"}
`,
			mock: func() {
				estateMock.EXPECT().GetTreeStats(gomock.Any(), common.UtUuid).Return(nil, domain.ErrEstateNotFound)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/estate/%s/stats", test.args), nil)
			req.Header.Set(common.UtContentType, common.ContentTypeJson)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(common.UtUuid)

			test.mock()

			if assert.NoError(t, handler.GetTreeStats(c)) {
				assert.Equal(t, test.wantResult, rec.Body.String())
			}
		})
	}
}

func TestGetDroneFlyingDistance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	estateMock := mock_domain.NewMockEstateUsecase(ctrl)
	handler := &estateHandler{
		estateUsecase: estateMock,
	}
	type args struct {
		id          string
		maxDistance string
	}

	tests := []struct {
		name       string
		args       args
		wantResult string
		mock       func()
	}{
		{
			name: "success",
			args: args{
				id: common.UtUuid,
			},
			wantResult: `{"code":200,"message":"Success get drone flying distance","data":{"distance":100},"errors":null}
`,
			mock: func() {
				estateMock.EXPECT().GetDroneFlyingDistance(gomock.Any(), common.UtUuid, 0).Return(&domain.GetDroneFlyingDistanceResponse{
					Distance: 100,
				}, nil)
			},
		},
		{
			name: "error get drone flying distance",
			args: args{
				id: common.UtUuid,
			},
			wantResult: `{"code":404,"message":"estate not found","data":null,"errors":"estate not found"}
`,
			mock: func() {
				estateMock.EXPECT().GetDroneFlyingDistance(gomock.Any(), common.UtUuid, 0).Return(nil, domain.ErrEstateNotFound)
			},
		},
		{
			name: "error max distance param",
			args: args{
				id:          common.UtUuid,
				maxDistance: "aaa",
			},
			wantResult: `{"code":400,"message":"strconv.Atoi: parsing \"aaa\": invalid syntax","data":null,"errors":"strconv.Atoi: parsing \"aaa\": invalid syntax"}
`,
			mock: func() {
				// estateMock.EXPECT().GetDroneFlyingDistance(gomock.Any(), common.UtUuid, 0).Return(nil, domain.ErrEstateNotFound)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/estate/%s/drone-plan?max-distance=%v", test.args.id, test.args.maxDistance), nil)
			req.Header.Set(common.UtContentType, common.ContentTypeJson)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(common.UtUuid)

			test.mock()

			if assert.NoError(t, handler.GetDroneFlyingDistance(c)) {
				assert.Equal(t, test.wantResult, rec.Body.String())
			}
		})
	}
}
