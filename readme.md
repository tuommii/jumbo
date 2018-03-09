# Jumbo
Add games played with your friends and check your stats. WIP.

## Screenshots

![Screenshot](/static/images/add.png "Screenshot 1")
![Screenshot](/static/images/search.png "Screenshot 2")

## API
* /api/create/player
* /api/create/game
* /api/create/match

* /api/delete/player
* /api/delete/game
* /api/delete/match

Example for adding player

`curl -d "playerName=Jack Bauer" -X POST http://lol:lol@localhost:3000/api/create/player`

## TODO

* [ ] GUI for managing games, players and matches
* [ ] PostgreSQL
* [ ] UI/Javascript/CSS improvements
* [ ] All flash messages