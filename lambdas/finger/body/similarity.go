package body

import (
	"math"
	"reflect"
	"strings"
)

const (
	UP_RATE   = 0.3 // Gentle slope for values above basis
	DOWN_RATE = 3.3 // Aggressive slope for values below basis
)

type Number interface {
	uint64 | int64 | float64
}

func Avg(nums []float64) float64 {
	if len(nums) == 0 {
		return 0.0
	}

	sum := 0.0
	for _, num := range nums {
		sum += num
	}

	return sum / float64(len(nums))
}

func CompareValues(crit, comp reflect.Value) float64 {
	switch crit.Kind() {
	case reflect.String:
		return stringSimilarity(crit.String(), comp.String())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return numberSimilarity(crit.Uint(), comp.Uint())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return numberSimilarity(crit.Int(), comp.Int())

	case reflect.Float32, reflect.Float64:
		return numberSimilarity(crit.Float(), comp.Float())

	default:
		if reflect.DeepEqual(crit.Interface(), comp.Interface()) {
			return 1.0
		}
		return 0.0
	}
}

func numberSimilarity[T Number](crit, comp T) float64 {
	fCrit := float64(crit)
	fComp := float64(comp)

	if fCrit == fComp {
		return 1.0
	} else if fCrit == 0.0 || fComp == 0.0 {
		return 0.0
	}

	ratio := fComp / fCrit

	if ratio >= 1 {
		logDiff := math.Log10(ratio)
		return math.Pow(0.5, logDiff*UP_RATE)
	} else {
		logDiff := math.Log10(1 / ratio)
		return math.Pow(0.5, logDiff*DOWN_RATE)
	}
}

func stringSimilarity(crit, comp string) float64 {
	similarityPercentages := []float64{}
	critWords := strings.Fields(crit)
	compWords := strings.Fields(comp)

	for _, critWord := range critWords {
		bestSim := 0.0
		for _, compWord := range compWords {
			bestSim = max(bestSim, similarityPercentage(critWord, compWord))
			if bestSim == 1.0 {
				break
			}
		}
		similarityPercentages = append(similarityPercentages, bestSim)
	}

	return Avg(similarityPercentages)
}

func similarityPercentage(a, b string) float64 {
	distance := lavenshteinDistance(a, b)
	maxLen := max(len([]rune(a)), len([]rune(b)))
	if maxLen == 0 {
		return 100.0
	}
	return (1.0 - (float64(distance) / float64(maxLen)))
}

func lavenshteinDistance(a, b string) int {
	aRunes := []rune(a)
	bRunes := []rune(b)
	aLen := len(aRunes)
	bLen := len(bRunes)

	matrix := make([][]int, aLen+1)
	for i := range matrix {
		matrix[i] = make([]int, bLen+1)
	}

	for i := 0; i <= aLen; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= bLen; j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= aLen; i++ {
		for j := 1; j < bLen; j++ {
			cost := 0
			if aRunes[i-1] != bRunes[j-1] {
				cost = 1
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1,
				matrix[i][j-1]+1,
				matrix[i-1][j-1]+cost,
			)
		}
	}

	return matrix[aLen][bLen]
}
