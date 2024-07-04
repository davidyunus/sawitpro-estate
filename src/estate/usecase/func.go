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

func calculateTotalDistance(startX, startY int, palmTrees []domain.PalmTree, safetyLimit int) int {
	var totalDistance int
	currentHeight := 0
	x, y := startX, startY

	for _, tree := range palmTrees {
		// Horizontal distance
		totalDistance += (abs(tree.X-x) + abs(tree.Y-y)) * 10

		// Vertical distance
		newHeight := tree.Height + safetyLimit
		totalDistance += (newHeight - currentHeight)
		currentHeight = newHeight

		// Move to the tree's position
		x, y = tree.X, tree.Y
	}

	// Add distance for landing
	totalDistance += currentHeight

	return totalDistance
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

var generateUUID = func() string {
	return uuid.NewString()
}
