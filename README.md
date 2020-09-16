# Exploding-room

A chat room where the chat room disappears every 24 hours.  Check it out at [explodingroom.io](http://explodingroom.io)

![exploding room](https://i.imgur.com/zjYRrFY.png)

It's written in GO and uses the Gorilla Websockets library.  The frontend uses Bootstrap and Websockets to communicate with the backend.  


It's hosted on an EC2 micro instance with a gig of ram.  Connections are very cheap so for scaling this is more I/O bound than compute bound.


It's all held together by duct tape - I wanted to learn GO and stand this up as fast as possible, not showcase excellent GO.


## Features
* Asynchronous chat
* Support for creating and joining rooms

## Wishlist
* Persist messages and load old messages when joining a room
* Automatic reconnecting when connection is lost
* Support for switching to short polling when WS is not avaliable
