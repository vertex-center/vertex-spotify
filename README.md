# Spotify Service `BETA`

A Spotify microservice for your Home Lab.

## Features

- [x] **Player state**

  The *Spotify Service* post the current Spotify player state on your [Redis Pub/Sub](https://redis.io/docs/manual/pubsub/) instance. Services like [Discord Service](https://github.com/quentinguidee/discord-service) can then reuse these events to display different things.

  <img width="231" alt="image" src="https://user-images.githubusercontent.com/12123721/219262662-e6dfaa9d-dfd6-4c7c-8e00-38e4d3c7a9ff.png">

- [ ] **Statistics**

  This service can aggregate your Spotify statistics. Other services can then reuse them easily.
  
  *Soon available on [cloud.sh](https://github.com/quentinguidee/cloud-sh-client)*

- [ ] *More coming soon...*

## Setup

*This service is currently a work-in-progress.*

- Install Redis. Run an instance with the `redis-server` command.
- Run `./spotifyservice`
