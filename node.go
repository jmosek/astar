package main

import "math"

type Node struct {
	Obstacle   bool
	Visited    bool
	GlobalGoal float32
	LocalGoal  float32
	X, Y       int
	Neighbours []*Node
	Parent     *Node
}

// func (n *Node)
func (n1 *Node) DistanceTo(n2 *Node) float32 {
	return float32(math.Sqrt(float64((n1.X - n2.X) * (n1.X - n2.X) * (n1.Y - n2.Y) * (n1.Y - n2.Y))))
}
