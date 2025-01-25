# Yggdrasil Mock Server

This is a mock implementation of the Yggdrasil API.

## Features
- **User Authentication**: Token-based user authentication and token lifecycle management.
- **Character Management**: Create, manage, and retrieve Minecraft characters.
- **Texture Handling**: Upload, retrieve, and delete textures for Minecraft characters.
- **Session Server**: Simulates Minecraft session APIs for joining and checking server sessions.
- **Database Support**: Configurable database support for SQLite, MySQL, and PostgreSQL.

## Prerequisites
- **Go**: Version 1.20 or higher.
- **Database**:
   - SQLite 3 (default).
   - MySQL or PostgreSQL (optional).
- **Docker**: For containerized deployment.

## Installation

### Local Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/TeamLeak/yggdrasil-mock.git
   cd yggdrasil-mock
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the server:
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`.

### Docker Deployment

1. Clone the repository:
   ```bash
   git clone https://github.com/TeamLeak/yggdrasil-mock.git
   cd yggdrasil-mock
   ```

2. Build and run the Docker containers:
   ```bash
   docker-compose up --build
   ```

The server will start on `http://localhost:8080`.

### Configuration

The server uses a `config.yaml` file for configuration. If the file is not found, a default configuration will be generated. You can customize the following parameters:

```yaml
app:
  port: "8080"                # Server port
  app_url: "http://localhost" # Application base URL

database:
  type: "sqlite"              # Database type: sqlite, mysql, or postgres
  host: "localhost"           # Host for MySQL/PostgreSQL
  port: "3306"                # Port for MySQL/PostgreSQL
  user: "user"                # Database user
  password: "password"        # Database password
  name: "yggdrasil_db"        # Database name
  sqlite_file: "yggdrasil.db" # SQLite file path
```

### Generate Config File Manually

To generate a default `config.yaml` file, use:
```bash
go run main.go --generate-config
```

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

## Docker Volumes

- **Configuration Volume**: The `app_config` volume stores the `config.yaml` file.
- **Database Volume**: The `db_data` volume persists database data.

## OpenAPI Specification

The OpenAPI documentation is available in the [`openapi.yaml`](https://github.com/TeamLeak/yggdrasil-mock/blob/main/openapi.yaml) file.

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/TeamLeak/yggdrasil-mock/blob/main/LICENSE) file for details.
