package usecase

import (
	"sort"

	"github.com/davidyunus/sawitpro-estate/src/domain"
	"github.com/google/uuid"
)

func calculateMedian(arr []int) float64 {
	sort.Ints(arr)
	n := len(arr)

	if n%2 == 1 {
		return float64(arr[n/2])
	} else {
		middle1 := arr[(n/2)-1]
		middle2 := arr[n/2]
		return float64(middle1+middle2) / 2.0
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func generateCoordinates(length, width int) []domain.PalmTree {
	var coordinates []domain.PalmTree
	for y := 1; y <= width; y++ {
		if y%2 != 0 {
			for x := 1; x <= length; x++ {
				coordinates = append(coordinates, domain.PalmTree{X: x, Y: y, Height: 0})
			}
		} else {
			for x := length; x >= 1; x-- {
				coordinates = append(coordinates, domain.PalmTree{X: x, Y: y, Height: 0})
			}
		}
	}
	return coordinates
}

func mergePalmTrees(coordinates []domain.PalmTree, palmTrees []domain.PalmTree) []domain.PalmTree {
	for i := range coordinates {
		for _, palmTree := range palmTrees {
			if coordinates[i].X == palmTree.X && coordinates[i].Y == palmTree.Y {
				coordinates[i].Height = palmTree.Height
				break
			}
		}
	}
	return coordinates
}

func removeTrailingZeroHeightCoordinates(coordinates []domain.PalmTree) []domain.PalmTree {
	for i := len(coordinates) - 1; i >= 0; i-- {
		if coordinates[i].Height != 0 {
			break
		}
		coordinates = coordinates[:i]
	}
	return coordinates
}

var generateUUID = func() string {
	return uuid.NewString()
}
