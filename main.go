package main

import (
	"bufio"
	"fmt"
	"github.com/juliangruber/go-intersect"
	"os"
	"strconv"
	"strings"
	"sync"
)

// TODO: modify so that x0 and x1 are 2D (for real)
func getSNNDistance(x0, x1 [][]float64) float64 {
	inter := intersect.Sorted(x0, x1)
	return 1.0 - (float64(len(inter)) / float64(len(x0)))
}

// in go:
// [[1, 2, 3, 4, 5], [4, 5, 7, 8, 9]]
//
// in file:
// 1, 2, 3, 4, 5
// 4, 5, 7, 8, 9
func writeDistanceMatrix(matrix [][]float64, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for row := range matrix {
		for col := range matrix[row] {
			// Unless this is the last one, put a ", " after
			var err error
			if col == len(matrix[row]) - 1 {
				_, err = fmt.Fprint(w, matrix[row][col])
			} else {
				_, err = fmt.Fprint(w, fmt.Sprintf("%f, ", matrix[row][col]))
			}
			if err != nil {
				return err
			}
		}
		fmt.Fprint(w, "\n")
	}

	return w.Flush()
}


// in file:
// 1, 2, 3, 4, 5
// 4, 5, 7, 8, 9
//
// 1, 2, 3, 4, 5
// 4, 5, 7, 8, 9
//
// in go:
// [[[1, 2, 3, 4, 5], [4, 5, 7, 8, 9]], [[1, 2, 3, 4, 5], [4, 5, 7, 8, 9]]]
//
// the last row is required to have double newline !
func readNeighbors(filename string) ([][][]float64, error) {
	var ret3D [][][]float64
	var ret2D [][]float64

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var innerRet []float64

		l := scanner.Text()

		if l == "" {
			ret3D = append(ret3D, ret2D)
			ret2D = [][]float64{}
			continue
		}

		for _, strVal := range strings.Split(l, ",") {
			// strip out whitespace from string value to that ParseFloat is happy :)
			strStrippedVal := strings.ReplaceAll(strVal, " ", "")

			// convert string to a float64
			floatVal, err := strconv.ParseFloat(strStrippedVal, 64) // 64 is the number of bits
			if err != nil {
				return nil, err
			}

			innerRet = append(innerRet, floatVal)
		}

		ret2D = append(ret2D, innerRet)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ret3D, nil
}


// def get_snn_distance(x0, x1):
//    return 1 - (len(x0.intersection(x1)) / len(x0))


//snn_distance_matrix = np.zeros((len(neighbors), len(neighbors)))
//for i in range(len(neighbors)):
//	for j in range(i):
//		dist = get_snn_distance(neighbors[i], neighbors[j])
//		snn_distance_matrix[i][j] = dist
//		snn_distance_matrix[j][i] = dist

//func ()

func main() {
	neighbors, err := readNeighbors("text.csv")
	if err != nil {
		panic(err)
	}
	fmt.Println(neighbors)


	// fill snnDistanceMatrix with zeroes, make it (# neighbors x # neighbors)
	var snnDistanceMatrix [][]float64
	for _ = range neighbors {
		snnDistanceMatrix = append(snnDistanceMatrix, make([]float64, len(neighbors)))
	}


	//
	// TODO: do the code

	var wg sync.WaitGroup

	for i := range neighbors {
		for j := range neighbors[i] {
			wg.Add(1)
			go func() {
				dist := getSNNDistance(neighbors[i], neighbors[j])

				snnDistanceMatrix[i][j] = dist
				snnDistanceMatrix[j][i] = dist

				wg.Done()
			}()
		}
	}

	wg.Wait()



	err = writeDistanceMatrix(snnDistanceMatrix, "text-copy.csv")
	if err != nil {
		panic(err)
	}
}
