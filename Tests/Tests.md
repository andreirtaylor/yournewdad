# Testing Pics

- It has come fruition that I cant figure out how to easily look at the game
boards for the tests. This file shows the layout of the board for the tests

Not all tests are here

- used in
  - TestEfficientSpace
  - TestSetMinSnakePointInArea
![](assets/Tests-ff3f1.png)

- used in
  - TestExtraEfficientUseOfSpace
![](assets/Tests-316b9.png)


used in TestMetaDataOnlyOneSnake
![](assets/Tests-d6761.png)

data, err := NewMoveRequest(`{"you":"82557bbc-5ff2-4e51-8133-f6875d4f8d71","width":10,"turn":233,"snakes":[{"taunt":"battlesnake-go!","name":"7eef72e9-72fc-4c27-a387-898384639f46 (10x10)","id":"82557bbc-5ff2-4e51-8133-f6875d4f8d71","health_points":100,"coords":[[1,3],[0,3],[0,4],[0,5],[0,6],[0,7],[1,7],[2,7],[3,7],[3,8],[3,9],[4,9],[4,8],[4,7],[4,6],[4,5],[5,5],[5,4],[6,4],[7,4],[7,3],[6,3],[5,3],[4,3],[4,4],[3,4],[3,3],[3,2],[4,2],[5,2],[5,1],[4,1],[3,1],[2,1],[2,0],[3,0],[4,0],[5,0],[6,0],[7,0],[8,0],[9,0],[9,1],[9,2],[9,2]]}],"height":10,"game_id":"7eef72e9-72fc-4c27-a387-898384639f46","food":[[6,2],[7,5],[2,3]],"dead_snakes":[]}`)
if err != nil {
  t.Logf("error: %v", err)
}

# BUG

The snake in the lower left should always prefer to go up instead of down or to the right side.



![](assets/Tests-bf7c5.png)

{"you":"0086b0de-4d23-4189-9305-a7308240edb4","width":20,"turn":176,"snakes":[{"taunt":"Dad 2.0 Ready","name":"Your New Dad","id":"3a1cf2b6-ab7f-4870-b672-cb60d7ab4e67","health_points":90,"coords":[[1,15],[2,15],[3,15],[4,15],[4,14],[4,13],[4,12],[4,11],[4,10],[4,9],[3,9],[3,10],[3,11],[2,11],[1,11],[1,10],[1,9],[0,9],[0,8],[1,8],[2,8],[2,7],[3,7],[3,6],[3,5],[4,5],[5,5],[6,5],[7,5],[8,5],[9,5],[9,6],[9,7],[10,7]]},{"taunt":"Dad 2.0 Ready","name":"Your New Dad","id":"0086b0de-4d23-4189-9305-a7308240edb4","health_points":94,"coords":[[14,16],[13,16],[12,16],[11,16],[10,16],[10,15],[10,14],[11,14],[11,13],[10,13],[9,13],[8,13],[8,14],[8,15],[8,16],[8,17],[8,18],[7,18],[6,18],[6,19],[7,19],[8,19],[9,19],[10,19],[10,18],[11,18],[12,18],[13,18],[14,18],[15,18],[16,18],[17,18],[17,17],[17,16],[17,15],[16,15],[15,15],[15,14],[14,14],[13,14]]}],"height":20,"game_id":"1dd0a217-baef-46ca-bd17-d48a38d54436","food":[[0,6],[16,16],[15,19],[18,7],[8,4],[10,2],[0,15],[3,1],[13,2],[14,11]],"dead_snakes":[]}
