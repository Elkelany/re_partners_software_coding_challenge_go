package orderpacks

import (
	"fmt"
	"slices"
)

type UseCaseCalculateOrderPacksRequest struct {
	PackSizes  []uint64 // Slice of available pack sizes.
	OrderItems uint64   // Number of items to order.
}

// OK validates the UseCaseCalculateOrderPacksRequest.
// Returns an error if OrderItems or any PackSizes are zero.
func (req UseCaseCalculateOrderPacksRequest) OK() error {
	if req.OrderItems == 0 {
		return fmt.Errorf("order items can't be zero")
	}

	for _, packSize := range req.PackSizes {
		if packSize == 0 {
			return fmt.Errorf("pack size can't be zero")
		}
	}

	return nil
}

// UseCaseCalculateOrderPacks implements the logic for calculating optimal order packs.
type UseCaseCalculateOrderPacks struct {
}

// Run calculates the optimal order packs to fulfill the order.
// It explores different combinations of packs and returns the least amount of items and as few packs as possible to fulfil the order.
func (u UseCaseCalculateOrderPacks) Run(request UseCaseCalculateOrderPacksRequest) (map[uint64]uint64, error) {
	// Validate the request
	if err := request.OK(); err != nil {
		return nil, err
	}

	packSizes := request.PackSizes

	// Sort pack sizes in ascending order.
	slices.Sort(packSizes)

	// Reverse the slice to have pack sizes in descending order.
	packSizes = reverse(packSizes)

	items := request.OrderItems

	// Calculate the initial order pack.
	initialOrderPacks := calculateInitialOrderPacks(items, packSizes)

	// If luckily the initial order pack total is equal to the items ordered, then we return it.
	if calculateTotalItems(initialOrderPacks) == items {
		return initialOrderPacks, nil
	}

	// Calculate variations of order packs.
	orderPacksVariations := calculateOrderPacksVariations(items, packSizes, initialOrderPacks)

	// Find the optimal order packs.
	optimalOrderPacks := findOptimalOrderPacks(orderPacksVariations)

	optimalOrderPacks = optimizePacksCount(optimalOrderPacks, packSizes)

	return optimalOrderPacks, nil
}

// optimizePacksCount optimizes the number of packs by attempting to consolidate smaller packs into larger ones.
func optimizePacksCount(orderPacks map[uint64]uint64, packSizes []uint64) map[uint64]uint64 {
	// Reverse the slice to have pack sizes in ascending order.
	packSizes = reverse(packSizes)

	for i, packSize := range packSizes {
		if i < len(packSizes)-1 &&
			orderPacks[packSize] > 1 &&
			(packSizes[i+1]/packSize) >= 2 &&
			orderPacks[packSize]*packSize >= packSizes[i+1] {
			orderPacks[packSizes[i+1]]++
			orderPacks[packSize] = orderPacks[packSize] - (packSizes[i+1] / packSize)
		}
	}

	return orderPacks
}

// reverse reverses a slice of uint64.
func reverse(slice []uint64) []uint64 {
	reversed := make([]uint64, len(slice))
	for i, e := range slice {
		reversed[len(slice)-1-i] = e
	}

	return reversed
}

// calculateInitialOrderPacks calculates an initial order packs.
func calculateInitialOrderPacks(items uint64, packSizes []uint64) map[uint64]uint64 {
	packCounts := make(map[uint64]uint64)
	for _, packSize := range packSizes {
		packCounts[packSize] = 0
	}

	for i, packSize := range packSizes {
		// If we have a pack size and a next pack size, if the remaining items is less than the current pack size
		// and less or equal than the next pack size, we can skip the current pack size
		if i < len(packSizes)-1 && items < packSizes[i] && items <= packSizes[i+1] {
			continue
		}

		// Handle remaining items in the last iteration.
		if i == len(packSizes)-1 && items < packSize {
			packCounts[packSize]++

			continue
		}

		q := items / packSize // Calculate the number of whole packs.
		r := items % packSize // Calculate the remaining items.

		packCounts[packSize] = q

		// Add one more pack of this pack size to fulfill the remaining items.
		if r > 0 {
			packCounts[packSize]++
		}

		break // Stop after finding the first suitable pack size
	}

	return packCounts
}

// calculateOrderPacksVariations calculates variations of order packs by adjusting pack counts.
func calculateOrderPacksVariations(
	items uint64,
	packSizes []uint64,
	initialOrderPacks map[uint64]uint64,
) []map[uint64]uint64 {
	originalItems := items

	variations := make([]map[uint64]uint64, 0)
	variations = append(variations, initialOrderPacks)

	currentPacks := map[uint64]uint64{}

	for loop := true; loop; {
		for i, packSize := range packSizes {
			if i == 0 && len(currentPacks) > 0 {
				// If no packs of this size, remove the current pack size and move to the next.
				if currentPacks[packSize] == 0 && len(packSizes) > 1 {
					packSizes = packSizes[1:]
					break
				}

				// If it's the smallest pack size, stop looping.
				if i == len(packSizes)-1 {
					loop = false
					continue
				}

				// Reduce count of largest pack size and add items back to the remaining items.
				currentPacks[packSize]--
				items = items + packSize

				continue
			}

			q := items / packSize // Calculate the number of whole packs.
			r := items % packSize // Calculate the remaining items.

			currentPacks[packSize] += q // Store the number of packs for this size.

			if r == 0 {
				variation := copyMap(currentPacks)
				variations = append(variations, variation) // Add variation.

				// If luckily the variation total is equal to the items ordered, then we stop looping.
				if calculateTotalItems(variation) == originalItems {
					loop = false
				}

				break
			}

			items = r // Update remaining items.

			// Handle the case where there are remaining items and it's the last pack size.
			// In this scenario, we add one more of the largest pack size to fulfill the remaining items.
			if r > 0 && i == len(packSizes)-1 {
				variation := copyMap(currentPacks) // Create a copy to test adding a pack.
				variation[packSize]++              // Add an extra pack.

				// If the new variation total is greater than or equal the last one, we skip and recalculate the packs.
				if calculateTotalItems(variation) >= calculateTotalItems(variations[len(variations)-1]) {
					items = r // Update remaining items.
					continue
				}

				variations = append(variations, variation) // Add variation.

				if calculateTotalItems(variation) == originalItems {
					loop = false
				}
				break
			}
		}
	}

	return variations
}

// total calculates the total number of items in an order pack.
func calculateTotalItems(orderPacks map[uint64]uint64) uint64 {
	total := uint64(0)
	for packSize, count := range orderPacks {
		total += packSize * count
	}

	return total
}

// copyMap creates a copy of a map[uint64]uint64.
func copyMap(src map[uint64]uint64) map[uint64]uint64 {
	dst := make(map[uint64]uint64)
	for key, value := range src {
		dst[key] = value
	}

	return dst
}

// findOptimalOrderPacks finds the optimal order pack with the smallest total items.
func findOptimalOrderPacks(orderPacksSlice []map[uint64]uint64) map[uint64]uint64 {
	optimalOrderPack := make(map[uint64]uint64)

	if len(orderPacksSlice) == 0 {
		return optimalOrderPack
	}

	optimalOrderPack = orderPacksSlice[0]

	for _, orderPacks := range orderPacksSlice {
		if calculateTotalItems(orderPacks) < calculateTotalItems(optimalOrderPack) {
			optimalOrderPack = orderPacks
		}
	}

	return optimalOrderPack
}
