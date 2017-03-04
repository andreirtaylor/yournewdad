package main

func headOn(data *MoveRequest, sid int) bool {
	theirHead := data.Snakes[sid].Head()
	myHead := data.Snakes[data.MyIndex].Head()

	myHeadToTheirHead := myHead.Dist(theirHead)
	if myHeadToTheirHead.X+myHeadToTheirHead.Y != 4 || len(data.Snakes[sid].Coords) < data.MyLength {
		return false
	}
	// returns the first piece of a snakes body
	myFirstBody := &(data.Snakes[data.MyIndex].Coords[1])
	theirFirstBody := &(data.Snakes[sid].Coords[1])
	myHeadToTheirBody := myHead.Dist(theirFirstBody)
	theirHeadToMyBody := theirHead.Dist(myFirstBody)

	if (totalDist(myHeadToTheirBody) > totalDist(myHeadToTheirHead)) && (totalDist(theirHeadToMyBody) > totalDist(myHeadToTheirHead)) {
		return true
	}
	return false
}

func totalDist(p *Point) int {
	return p.X + p.Y
}
