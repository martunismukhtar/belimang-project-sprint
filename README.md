# Go! Belimang

## Requirements

- Install Docker
- Install VS Code Extension: "Dev Containers" (for development inside a container)
- Install Go v1.25 (for vanilla/local development)

> **Note:**  
> Go Fiber & PostgreSQL is included for local development using  devcontainer
> PostgreSQL env credential :
> - POSTGRES_USER= abc
> - POSTGRES_PASSWORD= password
> - POSTGRES_DB= belimang
> - ports= 5432
---

## Getting Started

### For VS Code Developers 👋

1. **Clone this repository**
    ```bash
    git clone <this-repo>
    cd belimang
    code .
    ```

2. **Open in a Dev Container**
    
    Open the project in a [Dev Container](https://containers.dev/) for a ready-to-use development environment.
    - Install the "Dev Containers" extension in VS Code.
    - Press `CTRL + SHIFT + P`, then select `Rebuild and Reopen in Container` (or `Reopen in Container`).

3. **Start Developing**
    ```bash
    # Run development using hot reload inside dev container shell
    air
    # Open http://localhost:3000 for web view
    ```
---

### For Vanilla Developers 🍦

1. **Clone this repository**
    ```bash
    git clone <this-repo>
    cd belimang
    ```

2. **Install Go v1.25**

    https://go.dev/dl/

3. **Start Developing**
    ```bash
    # Build and run the belimang app using Docker Compose
    docker compose up
    # Open http://localhost:3000 for web view

    # After making changes, restart the app:
    docker compose down
    docker compose up
    ```

---

## Contributing

1. Create a branch for each requirement/feature.
2. Make a Pull Request (PR) into the main branch.

> **Remember:**  
> Commit `go.sum` and `go.mod` after adding any new packages during development.

---