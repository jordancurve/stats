/*
Given two groups of numbers and a desired number of iterations N, permtest
runs N trials to estimate the probability that a random re-arrangement of the
numbers into two groups of the same size as the original would produce at least
as large a difference in means the difference between the means of the
original two groups.

In short, suppose you see some effect, some difference between two sets of
numbers.  permtest gives you the probability that an difference at least as big
could be achieved simply by randomly dividing those numbers into two sets.  The
lower the output from permtest, the more significant your effect is.

If the total number of possible arrangements is less than or equal to the
desired number of iterations N, permtest calculates the probability exactly
by trying all possible arrangements.
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

func main() {
	flag.Usage = func() {
		fmt.Printf("usage: permtest [-iter N] numbersFile1 numbersFile2\n")
	}
	nIter := flag.Int("iter", 0, "number of iterations")
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
	}
	rand.Seed(time.Now().UTC().UnixNano())
	a := readNums(flag.Arg(0))
	b := readNums(flag.Arg(1))
	sort.Float64s(a)
	sort.Float64s(b)
	fmt.Printf("%.3g\n", permTest(a, b, *nIter))
}

func readNums(path string) []float64 {
	fh, err := os.Open(path)
	if err != nil {
		panic(err)
	}
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

func sum(nums []float64) float64 {
	var sum float64
	for _, n := range nums {
		sum += n
	}
	return sum
}

// returns the median of a sorted list of numbers
func median(nums []float64) float64 {
	x, y := medianIndices(len(nums))
	return 0.5*nums[x] + 0.5*nums[y]
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// medianIndices returns the two indices at which the median of a sorted array with n elements would lie. The median will be calculate (array[x] + array[y])/2.
func medianIndices(n int) (int, int) {
	if n%2 == 0 {
		return (n - 1) / 2, n / 2
	}
	return (n - 1) / 2, (n - 1) / 2
}

func mean(nums []float64) float64 {
	var sum float64
	for _, f := range nums {
		sum += f
	}
	return sum / float64(len(nums))
}

// Given two groups of numbers and a desired number of iterations N, permtest
// runs N trials to estimate the probability that a random re-arrangement of the
// numbers into two groups of the same size as the original would produce at least
// as large a difference in means as the difference between the means of the
// original two groups.

func permTest(a, b []float64, nIter int) float64 {
	/* needed for median
	if !sort.Float64sAreSorted(a) {
		panic("permTest() called with unsorted first arg")
	}
	if !sort.Float64sAreSorted(b) {
		panic("permTest() called with unsorted second arg")
	}
	*/
	origDiff := math.Abs(mean(a) - mean(b))
	both := append(a, b...)
	// sort.Float64s(both) // for median
	numBigger := 0 // number of times difference in medians was larger than orig
	numTrials := 0
	sumBoth := sum(both)
	countBigger := func(c []int) {
		numTrials++
		sum1 := sumBoth
		sum2 := float64(0)
		for _, i := range c {
			sum1 -= both[i]
			sum2 += both[i]
		}
		mean1 := sum1 / float64(len(both)-len(c))
		mean2 := sum2 / float64(len(c))
		//fmt.Printf("%g / %g %g (%g)\n", mean1, mean2, math.Abs(mean1-mean2), origDiff)
		if math.Abs(mean1-mean2) >= origDiff {
			numBigger++
		}
		/* for median:
		sort.Ints(c)
		group1 := make([]float64, 0, len(c))
		group2 := make([]float64, 0, len(both)-len(c))
		for i := 0; i < len(both); i++ {
			if len(group1) < len(c) && i == c[len(group1)] {
				group1 = append(group1, both[i])
				continue
			} else {
				group2 = append(group2, both[i])
			}
		}
		// fmt.Printf("%v (%g) / (%g) = %g\n\n", c, median(group1), median(group2), math.Abs(median(group1)-median(group2)))
		if math.Abs(median(group1)-median(group2)) >= origDiff {
			numBigger++
		}
		*/
	}
	lenSmallerGroup := min(len(a), len(b))
	nComb := binom(len(a), len(b))
	if nIter == 0 {
		nIter = 3e7 / len(both)
		fmt.Printf("using %d iterations\n", nIter)
	}
	if nComb <= float64(nIter) {
		fmt.Fprintf(os.Stderr, "doing exact test\n")
		allCombs(len(both), lenSmallerGroup, countBigger)
	} else {
		for i := 0; i < nIter; i++ {
			randComb(len(both), lenSmallerGroup, countBigger)
		}
	}

	return float64(numBigger) / float64(numTrials)
}

// https://rosettacode.org/wiki/Combinations#Go
// Given non-negative integers m and n, generate all size m combinations of the integers from 0 to n-1 in sorted order (each combination is sorted and the entire table is sorted).

func allCombs(n, m int, emit func([]int)) {
	s := make([]int, m)
	last := m - 1
	var rc func(int, int)
	rc = func(i, next int) {
		for j := next; j < n; j++ {
			s[i] = j
			if i == last {
				emit(s)
			} else {
				rc(i+1, j+1)
			}
		}
		return
	}
	rc(0, 0)
}

// randComb(n, m, e) creates a random combination of m items from [0, n) and calls e() with it.
//
// http://stackoverflow.com/questions/2394246/algorithm-to-select-a-single-random-combination-of-values
// http://math.stackexchange.com/questions/178690/whats-the-proof-of-correctness-for-robert-floyds-algorithm-for-selecting-a-sin
func randComb(n, m int, emit func([]int)) {
	S := map[int]bool{}
	c := make([]int, m)
	for j := n - m; j < n; j++ {
		t := rand.Intn(j + 1)
		var x int
		if !S[t] {
			x = t
		} else {
			x = j
		}
		c[len(S)] = x
		S[x] = true
	}
	emit(c)
}

//  binom(n, m) = n+m choose m = (n+m)!/(n!m!)
func binom(n, m int) float64 {
	var big, small int
	if n > m {
		big, small = n, m
	} else {
		big, small = m, n
	}
	res := float64(1)
	s := small
	for i := big + small; i > big; i-- {
		res *= float64(i) / float64(s)
		s--
	}
	return res
}
