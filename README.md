
# BlogFlex

BlogFlex is an interactive blogging platform designed for bloggers and readers alike. Built with Go, HTMX, Templ, Bootstrap, Hasura (GraphQL + PostgreSQL) and Docker BlogFlex offers a comprehensive set of features for managing and enjoying content connected to multiple external services for enhanced user experience.

## Features

- User registration and authentication
- Session management with Gorilla Sessions
- JWT-based authentication middleware
- Post creation, listing, and detail views
- Commenting system
- Real-time updates with HTMX
- Responsive design with Bootstrap
- TinyMCE advanced WYSIWYG HTML editor
- Database management with Hasura using GraphQL with PostgreSQL
- Containerization with Docker

## Technologies Used

- **Backend**: Go, GraphQL, PostgreSQL, Hasura
- **Frontend**: HTMX, Templ, Bootstrap, TinyMCE WYSIWYG HTML editor
- **Containerization**: Docker

## Getting Started

### Prerequisites

- Go (version 1.18 or later)
- HTMX 
- Templ
- Docker
- Hasura Account
- [Air](https://github.com/cosmtrek/air) for live reloading

### Installation

1. **Clone the repository**

   ```sh
   git clone https://github.com/yourusername/blogflex.git
   cd blogflex
   ```

2. **Register a Hasura Account**

 Register for an account at [Hasura](https://hasura.io/).


3. **Install Go dependencies**

   ```sh
   go mod tidy
   ```

4. **Run the application**

   ```sh
   air
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
├── .vscode/
│   └── launch.json
├── internal/
│   ├── database/
│   │   └── database.go
│   ├── handlers/
│   │   ├── blog_handlers.go
│   │   ├── comment_handlers.go
│   │   ├── post_handlers.go
│   │   └── user_handlers.go
│   ├── helpers/
│   │   ├── format_time.go
│   │   ├── graphql.go
│   │   ├── logged_in.go
│   │   └── respond_error.go
│   │   └── truncate_words.go
│   └── middleware/
│   │   └── auth.go
│   └── models/
│       └── structs.go
├── router/
│   └── router.go
├── views/
│   ├── blog_list.templ
│   │   └── blog_list_templ.go
│   ├── blog_page.templ
│   │   └── blog_page_templ.go
│   ├── create.templ
│   │   └── create_templ.go
│   ├── detail.templ
│   │   └── detail_templ.go
│   ├── edit.templ
│   │   └── edit_templ.go
│   ├── index.templ
│   │   └── index_templ.go
│   ├── navbar_components.templ
│   │   └── navbar_components_templ.go
│   ├── navbar.templ
│   │   └── navbar_templ.go
│   └── post_list.templ
│       └── post_list_templ.go
├── .air.toml
├── .env
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## Development

### Running the Project Locally

To run the project locally, you have 3 options:

1. **Launch Debugger**:
   - Open your project in Visual Studio Code.
   - Set breakpoints as needed.
   - Launch the debugger by pressing `F5` or by selecting `Run > Start Debugging` from the menu.

2. **Run Air**:
   - Ensure you have [Air](https://github.com/cosmtrek/air) installed for live reloading.
   - Start Air by running the following command in your terminal:

     ```sh
     air
     ```

3. **Run `go run main.go` Command**:
   - Open your terminal.
   - Navigate to the project directory.
   - Run the following command to start the application:

     ```sh
     go run main.go
     ```

**Note**: Before running your project, make sure to generate the Templ files to get the most updated UI. You can do this by running:
```sh
npm run generate:templ

### Testing

To test the API endpoints, use the GraphQL queries in Hasura. Follow these steps:

1. **Access Hasura Console**:
   - Log in to your Hasura account at [Hasura](https://hasura.io/).
   - Navigate to your project's Hasura Console.

2. **Navigate to the API Tab**:
   - In the Hasura Console, go to the "API" tab.

3. **Run GraphQL Queries**:
   - Use the GraphQL query editor to write and execute your queries.
   - You can test various API endpoints by constructing appropriate GraphQL queries and mutations.

4. **Inspect Responses**:
   - Check the responses returned by the server to ensure your API is functioning correctly.

For more advanced testing, you can also use tools like Postman or Insomnia to send GraphQL requests to your Hasura endpoints.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [Gorilla Mux](https://github.com/gorilla/mux)
- [HTMX](https://htmx.org/)
- [Bootstrap](https://getbootstrap.com/)
- [TinyMCE](https://www.tiny.cloud/)
- [Hasura](https://hasura.io/)
- [Docker](https://www.docker.com/)

---

Enjoy using BlogFlex! If you encounter any issues, please feel free to open an issue or submit a pull request.
