## GotToDo - A Todo App with Go, Templ, Tailwind CSS, and HTMX

GotToDo is a simple Todo application built with modern web technologies:

- **Backend**: Go with Echo web framework
- **Frontend**: Uses HTMX for interactivity with minimal JavaScript
- **Styling**: Tailwind CSS for a clean, responsive design
- **Templates**: Templ for type-safe HTML templating
- **Authentication**: GitHub OAuth integration
- **Database**: Supabase for data storage

## Features

- User authentication with GitHub OAuth or email/password
- Create, read, update, and delete todo items
- Mark todos as complete or incomplete
- Clean, responsive UI with Tailwind CSS
- Interactive UI with HTMX for minimal JavaScript
- Type-safe templating with Templ

## Project Structure

```
.
├── cmd/
│   └── server/           # Main application entry point
├── docs/                 # Documentation
├── internal/
│   ├── handlers/         # HTTP handlers
│   ├── models/           # Data models
│   ├── repositories/     # Data access layer
│   └── services/         # Business logic
├── migrations/           # Database migrations
├── pkg/
│   ├── auth/             # Authentication utilities
│   ├── config/           # Configuration management
│   └── database/         # Database utilities and client
├── ui/
│   └── templates/        # Templ templates for all UI components
│       ├── layout.templ  # Layout templates
│       ├── todo.templ    # Todo-related templates
│       ├── pages.templ   # Page templates
│       └── ajax.templ    # AJAX response templates
```

## Setup Instructions

### Prerequisites

- Go 1.21+
- Supabase account (for production use)

### Configuration

GotToDo uses a single JSON configuration file for all application settings. You can specify the configuration file path using the `-config` command-line flag when starting the application.

If no configuration file exists at the specified location (or the default location if no path is specified), a default configuration file will be created automatically.

Example configuration file:

```json
{
  "repository": {
    "type": "memory"
  },
  "server": {
    "port": "8080"
  },
  "database": {
    "supabase_url": "your_supabase_url",
    "supabase_anon_key": "your_supabase_anon_key",
    "supabase_db_url": "your_supabase_db_url"
  },
  "auth": {
    "github_client_id": "your_github_client_id",
    "github_client_secret": "your_github_client_secret",
    "github_redirect_url": "http://localhost:8080/auth/github/callback"
  }
}
```

Available repository types:
- `memory`: Stores todos in memory (data will be lost when the application restarts)
- `supabase`: Stores todos in a Supabase PostgreSQL database

### Running the Application

1. Install dependencies:
   ```
   go mod download
   ```

2. Generate Templ templates:
   ```
   cd ui && go install github.com/a-h/templ/cmd/templ@latest && templ generate
   ```

3. Run the application with the default configuration:
   ```
   go run cmd/server/main.go
   ```

   Or specify a custom configuration file path:
   ```
   go run cmd/server/main.go -config /path/to/config.json
   ```

4. Access the application at `http://localhost:8080` (or the port specified in your configuration)

## Why Templ?

Templ is a type-safe HTML templating language for Go that:

- Provides compile-time type checking for templates
- Integrates seamlessly with Go code
- Allows for component-based design
- Prevents common templating errors
- Makes refactoring safer

## Why HTMX?

HTMX allows us to build interactive web applications with minimal JavaScript by:

- Using HTML attributes to trigger AJAX requests
- Updating the DOM with server responses
- Providing smooth transitions and animations
- Reducing the need for client-side JavaScript

## Development

### Working with Templ Templates

Templ templates are defined in the `ui/templates` directory. After modifying templates, regenerate the Go code:

```
cd ui && templ generate
```

### Adding New Templates

1. Create a new `.templ` file in the appropriate directory
2. Define your templates using the Templ syntax
3. Regenerate the Go code
4. Use your templates in handlers

### Repository Implementation

The application supports multiple repository implementations:

- `MemoryTodoRepository`: In-memory storage for development
- `SupabaseTodoRepository`: Supabase PostgreSQL storage for production

## License

MIT
