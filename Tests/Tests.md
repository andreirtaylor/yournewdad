# Testing Pics

- It has come fruition that I cant figure out how to easily look at the game while testing
boards for the tests. This file shows the layout of the board for the tests
- The file **gamestrings.go** contains all of the games listed in this file
- to use the strings for testing include the following at the top of your test

```go
data, err := NewMoveRequest(gameString#)

if err != nil {
  t.Logf("error: %v", err)
}
```

## gameString1

![](assets/Tests-ff3f1.png)

## gameString2

![](assets/Tests-d6761.png)

## gameString3

![](assets/Tests-316b9.png)

## gameString4

The snake in the lower left should always prefer to go up instead of down or to the right side.

![](assets/Tests-bf7c5.png)

## gameString5

![](assets/Tests-5820d.png)

## gameString6

![](assets/Tests-f3578.png)

## gameString7

![](assets/Tests-ab782.png)

## gameString8

![](assets/Tests-8554d.png)

## gameString9

![](assets/Tests-ba58c.png)

You are the snake with head at [13, 2]

## gameString10

![](assets/Tests-0312c.png)

You are the snake with head at [7, 18]

## gameString11

![](assets/Tests-68ff2.png)

## gameString12

![](assets/Tests-4b72f.png)

You are the red snake

## gameString13

![](assets/Tests-561f1.png)

You are the red snake

## gameString14

![](assets/Tests-0605d.png)

## gameString15

![](assets/Tests-353b5.png)

## gameString16

![](assets/Tests-a70ec.png)

## gameString17

![](assets/Tests-c03f3.png)

## gameString18

![](assets/Tests-0196a.png)
