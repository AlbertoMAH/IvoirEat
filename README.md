# Go SaaS Application

This is a SaaS application for expense management, built with a Go backend and a Next.js frontend.

## Local Development Setup

This project is configured to run in a containerized environment using Docker and Docker Compose. This ensures a consistent and reproducible development setup.

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/) (usually included with Docker Desktop)
- [Visual Studio Code](https://code.visualstudio.com/)
- [Remote - Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) for VS Code

---

### Option 1: Using VS Code Dev Containers (Recommended)

This is the easiest and most integrated way to get started.

1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd <repository-folder>
    ```

2.  **Open in VS Code:** Open the cloned repository folder in Visual Studio Code.

3.  **Reopen in Container:** A notification will appear in the bottom-right corner asking if you want to "Reopen in Container". Click it.

    > If you don't see the notification, open the Command Palette (`Ctrl+Shift+P` or `Cmd+Shift+P`) and type `Remote-Containers: Reopen in Container`.

4.  **Done!** VS Code will now build and start the Docker containers and connect to the development environment. The backend, frontend, and database will be running.
    -   Frontend is available at `http://localhost:3000`
    -   Backend API is available at `http://localhost:8080`
    -   PostgreSQL database is available at `localhost:5432`

---

### Option 2: Using Docker Compose in the Terminal

If you prefer not to use the VS Code integration, you can manage the services directly from your terminal.

1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd <repository-folder>
    ```

2.  **Build and start the services:**
    ```bash
    docker-compose up --build -d
    ```
    The `-d` flag runs the containers in detached mode.

3.  **View logs:**
    ```bash
    # View logs for all services
    docker-compose logs -f

    # View logs for a specific service (e.g., backend)
    docker-compose logs -f backend
    ```

4.  **Stop the services:**
    ```bash
    docker-compose down
    ```

The application will be accessible at the same ports as mentioned above.
