package usecase

import (
	"context"
	"fmt"
	"log"
	"sort"

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

func (e *estateUsecase) GetDroneFlyingDistance(ctx context.Context, id string) (*domain.GetDroneFlyingDistanceResponse, error) {
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

	sort.Slice(trees, func(i, j int) bool {
		if trees[i].Y == trees[j].Y {
			return trees[i].X < trees[j].X
		}
		return trees[i].Y < trees[j].Y
	})

	trees = []domain.PalmTree{
		{X: 3, Y: 1, Height: 10},
		{X: 6, Y: 2, Height: 5},
		{X: 4, Y: 2, Height: 7},
		{X: 3, Y: 2, Height: 15},
		{X: 5, Y: 3, Height: 30},
	}
	log.Println(`trees`, trees)
	startX, startY := 1, 1
	// totalDistance := 0.0
	// latestHeight := 0
	safetyLimit := 1

	totalDistance := calculateTotalDistance(startX, startY, trees, safetyLimit)

	fmt.Printf("Total flying distance: %v meters\n", totalDistance)

	return &domain.GetDroneFlyingDistanceResponse{
		Distance: int(totalDistance),
	}, nil
}
