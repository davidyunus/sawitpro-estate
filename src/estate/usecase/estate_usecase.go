package usecase

import (
	"context"

	"github.com/davidyunus/sawitpro-estate/src/domain"
)

type estateUsecase struct {
	estateRepo           domain.EstateRepository
	palmTreeLocationRepo domain.PalmTreeLocationRepository
}

func NewEstateUsecase(estateRepo domain.EstateRepository, palmTreeLocationRepo domain.PalmTreeLocationRepository) domain.EstateUsecase {
	return &estateUsecase{
		estateRepo:           estateRepo,
		palmTreeLocationRepo: palmTreeLocationRepo,
	}
}

func (e *estateUsecase) CreateEstate(ctx context.Context, param *domain.Estate) (*domain.CreateEstateResponse, error) {
	plots := 100
	maxEstateSize := 50000
	estateSize := param.Width * param.Length * plots

	if estateSize > maxEstateSize {
		return nil, domain.ErrMaxSizeEstate
	}

	id := generateUUID()
	err := e.estateRepo.CreateEstate(ctx, &domain.Estate{
		Uuid:   id,
		Length: param.Length,
		Width:  param.Width,
	})
	if err != nil {
		return nil, err
	}

	return &domain.CreateEstateResponse{
		Id: id,
	}, nil
}

func (e *estateUsecase) PlantPalmTree(ctx context.Context, id string, param *domain.PalmTree) (*domain.PlantPalmTreeResponse, error) {
	estate, err := e.estateRepo.GetEstateByUuid(ctx, id)
	if err != nil {
		return nil, err
	}
	if estate == nil {
		return nil, domain.ErrEstateNotFound
	}

	trees, err := e.palmTreeLocationRepo.GetPalmTreesByUuid(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, tree := range trees {
		if tree.X == param.X && tree.Y == param.Y {
			return nil, domain.ErrLocationFilled
		}
	}

	err = e.palmTreeLocationRepo.PlantPalmTree(ctx, id, param)
	if err != nil {
		return nil, err
	}

	return &domain.PlantPalmTreeResponse{
		Id: estate.Uuid,
	}, nil
}

func (e *estateUsecase) GetTreeStats(ctx context.Context, id string) (*domain.GetTreeStatsResponse, error) {
	estate, err := e.estateRepo.GetEstateByUuid(ctx, id)
	if err != nil {
		return nil, err
	}
	if estate == nil {
		return nil, domain.ErrEstateNotFound
	}

	trees, err := e.palmTreeLocationRepo.GetPalmTreesByUuid(ctx, id)
	if err != nil {
		return nil, err
	}

	treeStatsResp := &domain.GetTreeStatsResponse{}
	treesHeight := []int{}
	for _, tree := range trees {
		treeStatsResp.Count++
		if treeStatsResp.Max < tree.Height {
			treeStatsResp.Max = tree.Height
		}
		if treeStatsResp.Min > tree.Height || treeStatsResp.Min == 0 {
			treeStatsResp.Min = tree.Height
		}
		treesHeight = append(treesHeight, tree.Height)
	}
	treeStatsResp.Median = int(calculateMedian(treesHeight))

	return treeStatsResp, nil
}

func (e *estateUsecase) GetDroneFlyingDistance(ctx context.Context, id string, maxDistance int) (*domain.GetDroneFlyingDistanceResponse, error) {
	estate, err := e.estateRepo.GetEstateByUuid(ctx, id)
	if err != nil {
		return nil, err
	}
	if estate == nil {
		return nil, domain.ErrEstateNotFound
	}

	palmTrees, err := e.palmTreeLocationRepo.GetPalmTreesByUuid(ctx, id)
	if err != nil {
		return nil, err
	}

	palmTreeArr := generateCoordinates(estate.Length, estate.Width)
	mergedCoordinates := mergePalmTrees(palmTreeArr, palmTrees)
	finalCoordinates := removeTrailingZeroHeightCoordinates(mergedCoordinates)

	horizontalDistance := 0
	verticalDistance := 0
	lastHeight := 0
	for i, coordinate := range finalCoordinates {
		if coordinate.Height != 0 {
			if lastHeight == 0 {
				verticalDistance++
				lastHeight = coordinate.Height
				verticalDistance += lastHeight
			} else if coordinate.Height != 0 {
				add := abs(lastHeight - coordinate.Height)
				verticalDistance += add
				lastHeight = coordinate.Height
			}
			if len(finalCoordinates)-1 == i {
				verticalDistance++
				verticalDistance += coordinate.Height
			}
		}
		horizontalDistance += 10
		if horizontalDistance+verticalDistance >= maxDistance && maxDistance != 0 {
			return &domain.GetDroneFlyingDistanceResponse{
				Distance: maxDistance,
				Rest: &domain.Rest{
					X: coordinate.X,
					Y: coordinate.Y,
				},
			}, nil
		}
	}

	totalDistance := horizontalDistance + verticalDistance
	return &domain.GetDroneFlyingDistanceResponse{
		Distance: int(totalDistance),
	}, nil
}
