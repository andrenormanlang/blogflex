# BlogFlex

BlogFlex is an interactive blogging platform designed for bloggers and readers alike. Built with Go, HTMX, Templ, Tailwind CSS, MariaDB, Docker, and Gorm, BlogFlex offers a comprehensive set of features for managing and enjoying content connected to multiple external services for enhanced user experience.

## Features

- User registration and authentication
- Post creation, listing, and detail views
- Commenting system
- Real-time updates with HTMX
- Responsive design with Tailwind CSS
- Database management with Gorm and MariaDB
- Containerization with Docker

## Technologies Used

- **Backend**: Go, Gorm, MariaDB
- **Frontend**: HTMX, Templ, Tailwind CSS
- **Containerization**: Docker

## Getting Started

### Prerequisites

- Go (version 1.18 or later)
- Docker
- MariaDB
- Node.js and npm (for Tailwind CSS)

### Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/yourusername/blogflex.git
   cd blogflex

2. **Set up MariaDB**

   Ensure MariaDB is running and create a database named`blogflex\`. Update the database configuration in`internal/database/database.go\` with your MariaDB credentials.

 ``` sql
   CREATE DATABASE blogflex;
 ```

3. **Install Go dependencies**

 ``` bash
   go mod tidy
 ```

4. **Install Node.js dependencies**

 ```sh
   npm install
 ```

5. **Run Tailwind CSS in watch mode**

 ```bash
   npm run watch:css
 ```

6. **Run the application**

 ```bash
   CompileDaemon -command="go run main.go"
 ```

### API Endpoints

#### Users

- **Create a user**

  ```http
  POST /users
  ```

  ```json
  {
    "name": "Test User",
    "email": "testuser@example.com",
    "password": "password"
  }
  ```

- **List users**

  ```http
  GET /users
  ```

- **Get user details**

  ```http
  GET /users/{id}
  ```

#### Posts

- **Create a post**

  ```http
  POST /posts
  ```

```json
  {
    "title": "Test Post",
    "content": "This is the content of the test post.",
    "user_id": 1
  }
```

- **List posts**

```http
  GET /posts
```

- **Get post details**

```http
  GET /posts/{id}
```

#### Comments

- **Create a comment**

```http
  POST /posts/{postID}/comments
```

```json
  {
    "content": "This is a comment.",
    "post_id": 1,
    "user_id": 1
  }
```

- **List comments for a post**

```http
  GET /posts/{postID}/comments
```

## Directory Structure

```plaintext
blogflex/
├── internal/
│   ├── database/
│   │   └── database.go
│   ├── handlers/
│   │   ├── handlers.go
│   │   ├── post_handlers.go
│   │   └── user_handlers.go
│   └── models/
│       ├── comment.go
│       ├── post.go
│       └── user.go
├── router/
│   └── router.go
├── node_modules/
├── static/
│   ├── css/
│   │   ├── tailwind.css
│   │   └── tailwind.generated.css
│   ├── img/
│   └── js/
│       └── app.js
├── views/
│   ├── create.templ
│   ├── create_templ.go
│   ├── detail.templ
│   ├── detail_templ.go
│   ├── index.templ
│   ├── index_templ.go
│   └── post_list.templ
│       └── post_list_templ.go
├── blogflex.exe
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
├── package.json
└── README.md
```

## Development

### Running the Project Locally

To run the project locally, ensure MariaDB is running and accessible. Run the following commands:

```bash
docker-compose up
npm run watch:css
CompileDaemon -command="go run main.go"
```

### Generating Templ Files

To generate Templ files, run:

```bash
templ generate
```

### Testing

Use tools like Postman or Insomnia to test the API endpoints.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [Gorilla Mux](https://github.com/gorilla/mux)
- [HTMX](https://htmx.org/)
- [Tailwind CSS](https://tailwindcss.com/)
- [Gorm](https://gorm.io/)
- [Docker](https://www.docker.com/)

---

Enjoy using BlogFlex! If you encounter any issues, please feel free to open an issue or submit a pull request.
