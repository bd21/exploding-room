# Exploding-room

A disposable chat room.  My first Go project.  Check it out at [explodingroom.io](http://explodingroom.io)

![exploding room](https://i.imgur.com/zjYRrFY.png)

It's written in GO and uses the Gorilla Websockets library.  The frontend uses Bootstrap and Websockets to communicate with the backend.  


It's currently hosted on an EC2 micro instance with a gig of ram.  Connections are very cheap so for scaling this is more I/O bound than compute bound.


## Features
* Asynchronous chat
* Support for creating and joining rooms

## Wishlist
* Destroy rooms after x minutes
* Persist messages and load old messages when joining a room
* Automatic reconnecting when connection is lost
* Support for switching to short polling when WS are not avaliable
