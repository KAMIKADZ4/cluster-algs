package cluster_algs

import (
	"math"
	"sort"
)

type AffinityPropagationOptions struct {
	MaxIterations int
	Damping       float64
	GetDistance   func(Point, Point) float64
}

// Находит медиану матрицы сходства
func findMedian(S [][]float64) float64 {
	var values []float64
	for i := 0; i < len(S); i++ {
		values = append(values, S[i]...)
	}
	sort.Float64s(values)
	if len(values)%2 == 0 {
		return (values[len(values)/2-1] + values[len(values)/2]) / 2
	} else {
		return values[len(values)/2]
	}
}

func AffinityPropagation(points []Point, options AffinityPropagationOptions) []Cluster {
	if options.GetDistance == nil {
		options.GetDistance = GetEuclideanDistance
	}

	n := len(points)
	if n == 0 {
		return []Cluster{}
	}

	// Матрица сходства (similarity): S(i, k) = -distance(i, k)
	S := make([][]float64, n)
	for i := 0; i < n; i++ {
		S[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			S[i][j] = -options.GetDistance(points[i], points[j])
		}
	}

	// Устанавливаем предпочтения (preferences) для значений на диагонале
	// Обычно используют медиану или минимальное значение матрицы сходства
	preferences := findMedian(S)
	for i := 0; i < n; i++ {
		S[i][i] = preferences
	}

	// Матрицы ответственности (responsibility) и доступности (availability)
	R := make([][]float64, n)
	A := make([][]float64, n)
	for i := range R {
		R[i] = make([]float64, n)
		A[i] = make([]float64, n)
	}

	// Итерации алгоритма
	for iter := 0; iter < options.MaxIterations; iter++ {
		// Обновление ответственности
		for i := 0; i < n; i++ {
			max1, max2 := math.Inf(-1), math.Inf(-1)
			max1Ind := -1
			for k := 0; k < n; k++ {
				val := A[i][k] + S[i][k]
				if val > max1 {
					max2 = max1
					max1 = val
					max1Ind = k
				} else if val > max2 {
					max2 = val
				}
			}
			for k := 0; k < n; k++ {
				if k == max1Ind {
					R[i][k] = options.Damping*R[i][k] + (S[i][k]-max2)*(1-options.Damping)
				} else {
					R[i][k] = options.Damping*R[i][k] + (S[i][k]-max1)*(1-options.Damping)
				}
			}
		}
		// подсчет суммы в столбце исключая диагональных значений
		column_sums := make([]float64, n)
		for k := 0; k < n; k++ {
			sum := 0.0
			for i := 0; i < n; i++ {
				if i != k {
					sum += max(0.0, R[i][k])
				}
			}
			column_sums[k] = sum
		}
		// Обновление доступности
		for i := 0; i < n; i++ {
			for k := 0; k < n; k++ {
				if i == k {
					// Заполнение значении на диагонали
					A[i][k] = (1-options.Damping)*column_sums[k] + options.Damping*A[i][k]
				} else {
					A[i][k] = (1-options.Damping)*math.Min(0, R[k][k]+column_sums[k]-R[i][k]) + options.Damping*A[i][k]
				}
			}
		}
	}

	centroidCluster := map[int]*Cluster{}
	for i := 0; i < n; i++ {
		maxVal := math.Inf(-1)
		clusterIndex := -1
		for k := 0; k < n; k++ {
			if val := R[i][k] + A[i][k]; val > maxVal {
				maxVal = val
				clusterIndex = k
			}
		}
		if clusterIndex != -1 {
			if _, ok := centroidCluster[clusterIndex]; !ok {
				centroidCluster[clusterIndex] = &Cluster{Points: []Point{}}
			}
			centroidCluster[clusterIndex].Points = append(centroidCluster[clusterIndex].Points, points[i])
		}
	}
	clusters := make([]Cluster, 0, len(centroidCluster))
	for _, val := range centroidCluster {
		val.Centroid = calculateCentroid(val.Points)
		clusters = append(clusters, *val) // Добавляем ключ в срез
	}
	return clusters
}
