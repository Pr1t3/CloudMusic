
# CloudMusic

CloudMusic is an online music player that allows users to upload and listen to their own music tracks. Built on a microservices architecture, it is part of the "Cloud" platform, which also includes CloudStorage.

## Features

- **Online Music Player**: Stream music directly from the web.
- **User Contributions**: Upload and share your own music.
- **Integration with Cloud Platform**: Seamless integration with other Cloud services like CloudStorage.

## Project Structure

The repository is organized into the following main directories and files:

- **ApiGateway**: Contains the API gateway configuration and code.
- **Services**: Includes various microservices that power the application.
  - **Api Gateway**: Routes requests between clients and internal services for efficient interaction.
  - **Authentication Service**: Manages user access and ensures the security of their musical data.
  - **Music Playback Service**: Handles requests for playing music tracks.
  - **Track Management Service**: Manages information about music tracks, genres, and artists.
  - **Playlist Management Service**: Organizes and manages user playlists for easy music access.
  - **File Storage Service**: Allows users to store their files on the server.
  - **Search Service**: Implements search functionality for music tracks, users, and genres.
  - **Client Service**: Facilitates user interaction with the web application's functionality.

## Technologies Used

- **Go**: The primary language used for the backend services.
- **HTML/CSS**: Used for the frontend interface

## Getting Started

To get started with CloudMusic, follow these steps:

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/Pr1t3/CloudMusic.git
   cd CloudMusic
   ```

2. **Start each service**:
   ```bash
   cd Service_name
   go run cmd/main.go
   ```
