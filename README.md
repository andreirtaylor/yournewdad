## battlesnake-go

This snake runs on google cloud

to run locally with the docker container you have to specify that the **-host** is your local ip

i.e.

goapp serve -host=192.168.0.33

## TODO

- Follow your tail when you are large.
- Do a basic minmax to see where I can get to vs other snakes. only count the places I can get to first in calculations of moves vs spaces
make moves to trap other snakes

- Make Snake a Pacifist Avoid other snakes entirely
- make sure you dont move onto a spot that is right beside the tail
- fix bug wehre you dont return a move on another snake being able to kill you or 

![](assets/README-58f0d.png)
ie. this shouldnt be a death
