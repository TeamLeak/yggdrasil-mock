# Yggdrasil Mock Server

This is a mock implementation of the Yggdrasil API.

## Features
- User authentication and token management.
- Character creation and management.
- Texture uploading and retrieval.
- Session server for joining and checking Minecraft game sessions.
- SQLite database for data storage.

## Prerequisites
- Go 1.20 or higher.
- SQLite 3.
- [Modernc SQLite driver](https://pkg.go.dev/modernc.org/sqlite).

## Installation

1. Clone the repository:

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the server:
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`.

## API Endpoints

| Method | Endpoint                                      | Description                                      |
|--------|----------------------------------------------|--------------------------------------------------|
| GET    | `/`                                          | Check if the server is running.                 |
| GET    | `/status`                                    | Get server statistics.                          |
| POST   | `/authserver/authenticate`                   | Authenticate a user and return a token.         |
| POST   | `/authserver/refresh`                        | Refresh an access token.                        |
| POST   | `/authserver/validate`                       | Validate an access token.                       |
| POST   | `/authserver/invalidate`                     | Invalidate an access token.                     |
| POST   | `/authserver/signout`                        | Sign out a user and revoke all their tokens.    |
| POST   | `/sessionserver/session/minecraft/join`      | Join a Minecraft server session.                |
| GET    | `/sessionserver/session/minecraft/hasJoined` | Check if a user has joined a server session.    |
| POST   | `/api/profiles/minecraft`                    | Query profiles by their usernames.              |
| GET    | `/sessionserver/session/minecraft/profile/:uuid` | Get the profile of a character by UUID.         |
| GET    | `/textures/:hash`                            | Retrieve texture data by hash.                  |
| DELETE | `/api/user/profile/:uuid/:textureType`       | Delete a texture for a character.               |
| PUT    | `/api/user/profile/:uuid/:textureType`       | Upload a texture for a character.               |

## OpenAPI Specification

The OpenAPI documentation is available in the `openapi.yaml` file.

## License
This project is licensed under the MIT License.