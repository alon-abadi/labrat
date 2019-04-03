package main

import (
	"fmt"
	"hash/crc32"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	const numTests int = 1

	distribution1 := append([]float32{}, 0.5, 0.25, 0.25)
	distribution2 := append([]float32{}, 0.7, 0.3)
	distribution3 := append([]float32{}, 0.25, 0.25, 0.25, 0.25)

	var groups [][]float32
	groups = append(groups, distribution1, distribution2, distribution3)

	for i := 0; i < numTests; i++ {
		runTest(getRandomDistribution(groups), determineEvaluationMethod())
	}
}

func getRandomDistribution(groups [][]float32) []float32 {
	rand.Seed(time.Now().UnixNano())
	return groups[rand.Intn(len(groups))]
}

func runTest(group []float32, method string) {
	fmt.Println(fmt.Sprintf("Method: %s, Distribution: %+v", method, group))

	const numIterations int = 100
	const minDriverID int = 2000000
	const maxDriverID int = 99999999

	buckets := map[string][]int{}
	totalAllBuckets := 0

	for i := 0; i < numIterations; i++ {
		rand.Seed(time.Now().UnixNano())
		driverID := rand.Intn(maxDriverID-minDriverID) + minDriverID
		if evaluate(method, driverID, group, buckets) {
			totalAllBuckets = totalAllBuckets + 1
		}
	}

	for k, v := range buckets {
		actualPercentage := float32(len(v)) / float32(totalAllBuckets)
		i, _ := strconv.Atoi(k)
		requiredPercentage := group[i]
		totalInBucket := len(v)
		fmt.Printf("[%s]: total => %d, required: %f, actual: %f, deviation: %.2f%%\n", k, totalInBucket, requiredPercentage, actualPercentage, math.Abs(float64(requiredPercentage-actualPercentage))*100)
	}

	fmt.Println()
}

func evaluate(method string, driverID int, distribution []float32, buckets map[string][]int) bool {
	var prob float32
	groupsSum := float32(0)
	rand.Seed(time.Now().UnixNano())

	if method == "hash" {
		h := crc32Num(strconv.Itoa(driverID), "salt")
		mod := (h % 100)
		prob = float32(mod) / 100
	} else {
		prob = rand.Float32()
	}

	for k, v := range distribution {
		if prob > groupsSum && prob <= groupsSum+v {
			buckets[groupName(k)] = append(buckets[groupName(k)], driverID)
			return true
		}
		groupsSum = groupsSum + v
	}
	return false
}

func groupName(k int) string {
	return fmt.Sprintf("%d", k)
}

func crc32Num(driverID string, salt string) uint {
	return uint(crc32.ChecksumIEEE([]byte(salt+driverID))) % 1000
}

func determineEvaluationMethod() string {
	if len(os.Args) > 1 && os.Args[1] == "hash" {
		return "hash"
	}
	return "weighted_random"
}
