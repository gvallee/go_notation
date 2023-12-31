//
// Copyright (c) 2020-2023, NVIDIA CORPORATION. All rights reserved.
//
// See LICENSE.txt for license information
//

package notation

import (
	"fmt"
	"strconv"
	"strings"
)

func addRange(str string, start int, end int) string {
	if str == "" {
		return fmt.Sprintf("%d-%d", start, end)
	}
	return fmt.Sprintf("%s,%d-%d", str, start, end)
}

func addSingleton(str string, n int) string {

	if str == "" {
		return fmt.Sprintf("%d", n)
	}

	return fmt.Sprintf("%s,%d", str, n)
}

func CompressIntArray(array []int) string {
	compressedRep := ""
	for i := 0; i < len(array); i++ {
		start := i
		for i+1 < len(array) && array[i]+1 == array[i+1] {
			i++
		}
		if i != start {
			// We found a range
			compressedRep = addRange(compressedRep, array[start], array[i])
		} else {
			// We found a singleton
			compressedRep = addSingleton(compressedRep, array[i])
		}
	}
	return compressedRep
}

func GetNumberOfEltsFromCompressedNotation(str string) (int, error) {
	num := 0
	// ", " and "," are two possible delimiters but cannot be mixed
	t1 := strings.Split(str, ", ")
	if len(t1) == 1 {
		t1 = strings.Split(str, ",")
	}
	for _, t := range t1 {
		t2 := strings.Split(t, "-")
		if len(t2) == 2 {
			val1, err := strconv.Atoi(t2[0])
			if err != nil {
				return 0, err
			}
			val2, err := strconv.Atoi(t2[1])
			if err != nil {
				return 0, err
			}
			num += val2 - val1 + 1
		} else {
			num++
		}
	}
	return num, nil
}

func ConvertStringRangesToIntSlice(str string) ([]int, error) {
	var callIDs []int

	tokens := strings.Split(str, ", ")
	tokensNoSpace := strings.Split(str, ",")
	if len(tokens) == 1 && len(tokensNoSpace) > 1 {
		tokens = tokensNoSpace
	}

	for _, t := range tokens {
		tokens2 := strings.Split(t, "-")
		if len(tokens2) == 2 {
			val1, err := strconv.Atoi(tokens2[0])
			if err != nil {
				return callIDs, err
			}
			val2, err := strconv.Atoi(tokens2[1])
			if err != nil {
				return callIDs, err
			}
			for i := val1; i <= val2; i++ {
				callIDs = append(callIDs, i)
			}
		} else {
			for _, t2 := range tokens2 {
				n, err := strconv.Atoi(t2)
				if err != nil {
					return callIDs, fmt.Errorf("unable to parse compressed string: %s", str)
				}
				callIDs = append(callIDs, n)
			}
		}
	}

	return callIDs, nil
}

// ConvertStringRangesToStringSlice converts the notation of ranges into a slice of string.
// For example, 001-004,009 becomes the slice []int{"001","002","003","004","009"}
func ConvertStringRangesToStringSlice(str string) ([]string, error) {
	var ids []string

	tokens := strings.Split(str, ", ")
	tokensNoSpace := strings.Split(str, ",")
	if len(tokens) == 1 && len(tokensNoSpace) > 1 {
		tokens = tokensNoSpace
	}

	for _, t := range tokens {
		tokens2 := strings.Split(t, "-")
		if len(tokens2) == 2 {
			val1, err := strconv.Atoi(tokens2[0])
			if err != nil {
				return ids, err
			}
			val2, err := strconv.Atoi(tokens2[1])
			if err != nil {
				return ids, err
			}

			for num := val1; num <= val2; num++ {
				switch len(tokens2[0]) {
				case 1:
					ids = append(ids, fmt.Sprintf("%d", num))
				case 2:
					ids = append(ids, fmt.Sprintf("%02d", num))
				case 3:
					ids = append(ids, fmt.Sprintf("%03d", num))
				default:
					return ids, fmt.Errorf("node number of composed of %d digits, we support a max of 3", len(tokens2[0]))
				}
			}
		} else {
			ids = append(ids, tokens2...)
		}
	}

	return ids, nil
}

func IntSliceToString(s []int) string {
	if len(s) == 0 {
		return ""
	}

	str := strconv.Itoa(s[0])
	for i := 1; i < len(s); i++ {
		str += "," + strconv.Itoa(s[i])
	}

	return str
}
