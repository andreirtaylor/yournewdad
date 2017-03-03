## battlesnake-go

This snake runs on google cloud

to run locally with the docker container you have to specify that the **-host** is your local ip

i.e.

goapp serve -host=192.168.0.33

## TODO

- Make Snake a Pacifist Avoid other snakes entirely
- fix bug wehre you dont return a move on another snake being able to kill you or 
- sort the snake array so that I dont go for food that is on the boundary of another snake
- watch out for head on colisions
    - test if the second thing on the snake is closer to me that the first thing
    - i.e the distance from me to their head is 3 and their second thing is 4 or greater.
- Clean up source code

![](assets/README-58f0d.png)
ie. this shouldnt be a death
