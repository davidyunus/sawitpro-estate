package domain

import "context"

type (
	EstateUsecase interface {
		CreateEstate(ctx context.Context, param *Estate) (*CreateEstateResponse, error)
		PlantPalmTree(ctx context.Context, id string, param *PalmTree) (*PlantPalmTreeResponse, error)
		GetTreeStats(ctx context.Context, id string) (*GetTreeStatsResponse, error)
		GetDroneFlyingDistance(ctx context.Context, id string, maxDistance int) (*GetDroneFlyingDistanceResponse, error)
	}

	EstateRepository interface {
		CreateEstate(ctx context.Context, param *Estate) error
		GetEstateByUuid(ctx context.Context, id string) (*Estate, error)
	}

	Estate struct {
		Uuid   string `json:"uuid"`
		Length int    `json:"length" validate:"gt=0"`
		Width  int    `json:"width" validate:"gt=0"`
	}

	CreateEstateResponse struct {
		Id string `json:"id"`
	}

	PlantPalmTreeResponse struct {
		Id string `json:"id"`
	}

	GetTreeStatsResponse struct {
		Count  int `json:"count"`
		Max    int `json:"max"`
		Min    int `json:"min"`
		Median int `json:"median"`
	}

	GetDroneFlyingDistanceResponse struct {
		Distance int   `json:"distance"`
		Rest     *Rest `json:"rest,omitempty"`
	}

	Rest struct {
		X int `json:"x"`
		Y int `json:"y"`
	}
)
