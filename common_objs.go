package cluster_algs

import "math"

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Cluster struct {
	Points   []Point `json:"points"`
	Centroid Point   `json:"centroid"`
}

func calculateCentroid(points []Point) Point {
	var sumX, sumY float64
	for _, p := range points {
		sumX += p.X
		sumY += p.Y
	}
	return Point{X: sumX / float64(len(points)), Y: sumY / float64(len(points))}
}
func GetEuclideanDistance(p1, p2 Point) float64 {
	return math.Sqrt(math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2))
}

func GetHaversineDistance(p1, p2 Point) float64 {
	// Радиус Земли в метрах
	const R = 6371000
	// Преобразование градусов в радианы
	fi1 := p1.X * math.Pi / 180
	fi2 := p2.X * math.Pi / 180
	deltaFi := (p2.X - p1.X) * math.Pi / 180
	deltaLambda := (p2.Y - p1.Y) * math.Pi / 180
	// Применение формулы Haversine
	a := (math.Sin(deltaFi/2)*math.Sin(deltaFi/2) +
		math.Cos(fi1)*math.Cos(fi2)*math.Pow(math.Sin(deltaLambda/2), 2.0))
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	// Расстояние в метрах
	return R * c
}
