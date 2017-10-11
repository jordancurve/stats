package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

func main() {
	myName := filepath.Base(os.Args[0])
	f := readAllNums(os.Args[1:])
	var res float64
	switch myName {
	case "median":
		sort.Float64s(f)
		res = median(f)
	case "mode":
		res = mode(f)
	case "mean":
		res = mean(f)
	default:
		panic("unknown action: " + myName)
	}
	fmt.Printf("%g\n", res)
}

func readAllNums(files []string) []float64 {
	if len(files) == 0 {
		return readNums(os.Stdin)
	}
	f := []float64{}
	for _, file := range files {
		fh, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		f = append(f, readNums(fh)...)
	}
	return f
}

func readNums(fh *os.File) []float64 {
	scanner := bufio.NewScanner(fh)
	scanner.Split(bufio.ScanWords)
	nums := []float64{}
	for scanner.Scan() {
		num, err := strconv.ParseFloat(scanner.Text(), 64)
		if err != nil {
			panic(err)
		}
		nums = append(nums, num)
	}
	return nums
}

// returns the median of a sorted list of numbers
func median(nums []float64) float64 {
	x, y := medianIndices(len(nums))
	return 0.5*nums[x] + 0.5*nums[y]
}

func mean(nums []float64) float64 {
	sum := float64(0)
	for _, x := range nums {
		sum += x
	}
	return sum / float64(len(nums))
}

// medianIndices returns the two indices at which the median of a sorted array with n elements would lie. The median will be calculate (array[x] + array[y])/2.
func medianIndices(n int) (int, int) {
	if n%2 == 0 {
		return (n - 1) / 2, n / 2
	}
	return (n - 1) / 2, (n - 1) / 2
}

func mode(nums []float64) float64 {
	var mode float64
	count := map[float64]int{}
	maxCount := -1
	for _, f := range nums {
		count[f]++
		if count[f] > maxCount {
			mode = f
			maxCount = count[f]
		}
	}
	return mode
}
