package main

import (
	"log"
	"math"
	"slices"
	"time"
)

const MAX_TRIES = 10000

func SolveAStar(startNode *Node, endNode *Node) bool {

	startTime := time.Now()

	for x := 0; x < MAP_WIDTH; x++ {
		for y := 0; y < MAP_HEIGHT; y++ {
			node := game.Grid[x][y]

			for _, n := range node.Neighbours {
				n.GlobalGoal = math.MaxFloat32
				n.LocalGoal = math.MaxFloat32
				n.Visited = false
				n.Parent = nil
			}
		}
	}

	currentNode := startNode
	var notTestedNodes []*Node
	startNode.LocalGoal = 0
	tries := 0
	startNode.GlobalGoal = startNode.DistanceTo(endNode)

	notTestedNodes = append(notTestedNodes, startNode)

	for len(notTestedNodes) > 0 && currentNode != endNode {
		slices.SortFunc(notTestedNodes, func(a, b *Node) int {
			if a.GlobalGoal > b.GlobalGoal {
				return 1
			} else if a.GlobalGoal < b.GlobalGoal {
				return -1
			} else {
				return 0
			}
		})

		if len(notTestedNodes) > 0 && notTestedNodes[0].Visited {
			notTestedNodes = notTestedNodes[1:]
		}

		if len(notTestedNodes) == 0 {
			break
		}

		currentNode = notTestedNodes[0]
		currentNode.Visited = true

		for _, nb := range currentNode.Neighbours {

			if !nb.Visited && !nb.Obstacle {
				notTestedNodes = append(notTestedNodes, nb)
			}

			possibleLowerGoal := currentNode.LocalGoal + currentNode.DistanceTo(nb)

			if possibleLowerGoal < nb.LocalGoal {
				nb.Parent = currentNode
				nb.LocalGoal = possibleLowerGoal

				nb.GlobalGoal = nb.LocalGoal + nb.DistanceTo(endNode)
			}
		}

		tries++

		if tries > MAX_TRIES {
			return false
		}
	}

	log.Printf("astar took %d us", time.Since(startTime).Microseconds())

	return true
}
