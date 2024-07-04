package domain

import "context"

type (
	PalmTreeLocationRepository interface {
		GetPalmTreesByUuid(ctx context.Context, id string) ([]PalmTree, error)
		PlantPalmTree(ctx context.Context, id string, param *PalmTree) error
	}

	PalmTree struct {
		Id     int64  `json:"id"`
		Uuid   string `json:"uuid"`
		X      int    `json:"x" validate:"gt=0"`
		Y      int    `json:"y" validate:"gt=0"`
		Height int    `json:"height" validate:"gte=1,lte=30"`
	}
)
