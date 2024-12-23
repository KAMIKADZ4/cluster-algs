package cluster_algs

import (
	"math"
	"math/rand"
	"time"
)

type KMeansOptions struct {
	K             int
	MaxIterations int
	GetDistance   func(Point, Point) float64
}

func initializeCentroids(points []Point, options KMeansOptions) []Point {
	random := rand.New(rand.NewSource(int64(time.Now().UnixNano())))
	centroids := make([]Point, 0, options.K)
	centroids = append(centroids, points[random.Intn(len(points))])
	for len(centroids) < options.K {
		// Вычисляем минимальные расстояния от каждой точки до ближайшего центроида
		distances := make([]float64, len(points))
		totalDistance := 0.0
		for i, p := range points {
			minDist := math.MaxFloat64
			for _, c := range centroids {
				d := options.GetDistance(p, c)
				if d < minDist {
					minDist = d
				}
			}
			distances[i] = minDist
			totalDistance += minDist
		}
		// Выбираем следующую точку как центроид с вероятностью, пропорциональной квадрату расстояния
		r := rand.Float64() * totalDistance
		cumulative := 0.0
		for i, d := range distances {
			cumulative += d
			if cumulative >= r {
				centroids = append(centroids, points[i])
				break
			}
		}
	}
	return centroids
}

func KMeans(points []Point, options KMeansOptions) []Cluster {
	if options.MaxIterations == 0 {
		options.MaxIterations = 100
	}
	if options.GetDistance == nil {
		options.GetDistance = getEuclideanDistance
	}
	clusters := make([]Cluster, options.K)
	centroids := initializeCentroids(points, options)
	for iter := 0; iter < options.MaxIterations; iter++ {
		// Создаем пустые кластеры
		for i := range clusters {
			clusters[i].Centroid = centroids[i]
			clusters[i].Points = []Point{}
		}
		// Распределяем точки по кластерам
		for _, p := range points {
			minDist := math.MaxFloat64
			nearestCluster := 0
			for i, c := range centroids {
				dist := options.GetDistance(p, c)
				if dist < minDist {
					minDist = dist
					nearestCluster = i
				}
			}
			clusters[nearestCluster].Points = append(clusters[nearestCluster].Points, p)
		}
		// Обновляем центроиды кластеров
		newCentroids := make([]Point, options.K)
		for i, cluster := range clusters {
			if len(cluster.Points) > 0 {
				newCentroids[i] = calculateCentroid(cluster.Points)
			} else {
				// Если кластер пустой, оставляем старый центроид
				newCentroids[i] = cluster.Centroid
			}
		}
		// Проверяем, изменились ли центроиды
		converged := true
		for i := range centroids {
			if options.GetDistance(centroids[i], newCentroids[i]) > 1e-6 {
				converged = false
				break
			}
		}
		centroids = newCentroids
		if converged {
			break
		}
	}
	return clusters
}
