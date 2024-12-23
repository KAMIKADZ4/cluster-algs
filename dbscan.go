package cluster_algs

type DbscanOptions struct {
	Eps         float64
	MinPts      int
	GetDistance func(Point, Point) float64
}

// Функция для расширения кластера
func expandCluster(pointIndex int, points []Point, visited []bool, options DbscanOptions) Cluster {
	clusterPoints := []Point{points[pointIndex]}
	nearIndexes := getNearPointIndexes(points[pointIndex], points, options)

	for len(nearIndexes) > 0 {
		index := nearIndexes[0]
		nearIndexes = nearIndexes[1:]
		// for i := 0; i < len(nearIndexes); i++ {
		// 	index := nearIndexes[i]

		if !visited[index] {
			visited[index] = true

			newNeighborIndexes := getNearPointIndexes(points[index], points, options)
			nearIndexes = append(nearIndexes, newNeighborIndexes...)
			clusterPoints = append(clusterPoints, points[index])
		}
	}

	centroid := calculateCentroid(clusterPoints)
	return Cluster{Points: clusterPoints, Centroid: centroid}
}

// Функция для поиска соседей точки
func getNearPointIndexes(targetPoint Point, points []Point, options DbscanOptions) []int {
	neighbors := []int{}
	for i, p := range points {
		if options.GetDistance(targetPoint, p) <= options.Eps {
			neighbors = append(neighbors, i)
		}
	}
	return neighbors
}

func Dbscan(points []Point, options DbscanOptions) ([]Cluster, []Point) {
	if options.GetDistance == nil {
		options.GetDistance = GetEuclideanDistance
	}

	clusters := []Cluster{}
	visited := make([]bool, len(points))
	noise := []Point{}

	for i := range points {
		if visited[i] {
			continue
		}
		visited[i] = true

		cluster := expandCluster(i, points, visited, options)
		if len(cluster.Points) >= options.MinPts {
			clusters = append(clusters, cluster)
		} else {
			noise = append(noise, cluster.Points...)
		}
	}

	return clusters, noise
}
