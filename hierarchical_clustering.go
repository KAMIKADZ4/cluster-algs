package cluster_algs

import (
	"math"
)

type HierarchicalClusteringOptions struct {
	GetDistance func(Point, Point) float64
	K           int
}

func HierarchicalClustering(points []Point, options HierarchicalClusteringOptions) []Cluster {
	if options.GetDistance == nil {
		options.GetDistance = GetEuclideanDistance
	}
	// Инициализируем каждый элемент как отдельный кластер
	clusters := make([]Cluster, len(points))
	for i, p := range points {
		clusters[i] = Cluster{Points: []Point{p}, Centroid: p}
	}
	// Выполняем объединение кластеров, пока не останется один кластер
	for len(clusters) > options.K {
		// Найти два ближайших кластера
		minDist := math.MaxFloat64
		var mergeA, mergeB int
		for i := 0; i < len(clusters); i++ {
			for j := i + 1; j < len(clusters); j++ {
				dist := options.GetDistance(clusters[i].Centroid, clusters[j].Centroid)
				if dist < minDist {
					minDist = dist
					mergeA, mergeB = i, j
				}
			}
		}
		// Объединить два ближайших кластера
		newCluster := Cluster{
			Points: append(clusters[mergeA].Points, clusters[mergeB].Points...),
		}
		newCluster.Centroid = calculateCentroid(
			[]Point{clusters[mergeA].Centroid, clusters[mergeB].Centroid},
		)
		// Удалить старые кластеры и добавить новый
		if mergeA > mergeB {
			mergeA, mergeB = mergeB, mergeA
		}
		clusters = append(clusters[:mergeA], clusters[mergeA+1:]...) // Удаляем mergeA
		clusters = append(clusters[:mergeB-1], clusters[mergeB:]...) // Удаляем mergeB
		clusters = append(clusters, newCluster)
	}
	return clusters
}
