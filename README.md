# TicTacToe

TicTacToe is a small microservice that acts as the backend for a tic tac toe game by using the Firebase realtime database.

## Contents
- [TicTacToe](#tictactoe)
  - [Contents](#contents)
  - [Requirements](#requirements)
  - [Getting started](#getting-started)
  - [API Documentation](#api-documentation)
  - [Implementation Status](#implementation-status)

## Requirements

- Go v1.15.x installed
- [golangci-lint](https://github.com/golangci/golangci-lint) installed on your local machine
  - used for linting and detecting common errors

## Getting started

1. Download the Go dependencies
  
    ```console
    $ go mod download
    ```
2. Clone (i.e make a copy) `.env.example` to `.env` and modify it as needed  
3. Build and execute the microservice

    ```console
    $ make buildrun
    ```

## API Documentation

The microservice exposes a [Swagger UI](https://swagger.io/tools/swagger-ui/) interface on `http://{host}:{port}/swagger/index.html`.

## Implementation Status
- [x] Allow only one move per player, i.e. one after the other
- [x] Detect if a ‘POST /move’ call was made by a player or the Firebase function
- [x] Serve correct status codes and descriptive response bodies for invalid input
- [x] Sanity check input and game state
- [ ] Bonus points for detecting if a game will result in a draw / deadlock early on
  - *Detection does not happen **early on**.*
- [x] Provide unit tests