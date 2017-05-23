# YourNewDad

This snake was my entry into the 2017 battlesnake competition. It defeated the most bounty snakes and was a contender in the semi finals of the competition.

## Running

There are no dependencies just build with

`go build .`

Then run with

`./yournewdad`

The server runs by default on port 9000 although this is configurable with environment variables

`PORT=8080 ./yournewdad`

will run the server on port 8080 instead of 9000

Alternatively, if you just want to play against this snake it will perpetually live at http://yournewdad.herokuapp.com/

## Testing

There are a lot of tests included with this snake most of them test behaviour. To run the tests run

`go test` 

while in the directory

To add new tests when running the server add the environment variable `YND_LOG`

`YND_LOG=true ./yournewdad`

This will print out the incomming json blob on every request. You then add this to the `gamestrings.go` file and you can reference the blob in your test (see `filters_test.go`) for examples

```go
...
data, err := NewMoveRequest(gameString17)

if err != nil {
        t.Errorf("%v", err)
}
...

```

`NewMoveRequest` is a function that parses the json blob into the internal data structure used in this program.

## Improvements

The snake typically dies by running head on into other snakes or by starving if you can find a way to detect head on colisions that would be a major improvement into the usability of the snake.

It would also be really nice to have more metrics about the win loss rate of the snake overallI want to implement some sort of mechanism that tests if you have won or lost the game then stores this information so that you can view the win loss and kill ratio of the snake.


