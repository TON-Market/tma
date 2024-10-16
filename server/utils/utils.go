package utils

import (
	"fmt"
	"github.com/tonkeeper/tongo/tlb"
	"math"
	"strconv"
)

func FloatToGrams(v float64) tlb.Grams {
	grams := uint64(math.Round(v * 1e9))
	return tlb.Grams(grams)
}

func GramsToString(g tlb.Grams) string {
	return strconv.FormatUint(uint64(g), 10)
}

//func GramsToStringInFloat(g tlb.Grams) string {
//	gramsInFloat := float64(g) / 1e9
//	rounded := math.Round(gramsInFloat*10) / 10
//	return strconv.FormatFloat(rounded, 'f', 1, 64)
//}

func GramsToStringInFloat(g tlb.Grams) string {
	gramsInFloat := float64(g) / 1e9
	rounded := math.Round(gramsInFloat*1e4) / 1e4
	return strconv.FormatFloat(rounded, 'f', 4, 64)
}

func FloatToString(f float64) string {
	return fmt.Sprintf("%.0f", f)
}
