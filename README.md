# GoTToDo

A simple todo application built with the GoTTH stack:

- **Go**: Backend language
- **Templ**: HTML templating
- **Tailwind CSS**: Styling
- **HTMX**: Frontend interactivity

## Features

- GitHub OAuth authentication
- Email/password authentication as fallback
- Create, read, update, and delete todo items
- Mark todo items as complete/incomplete
- Data isolation between users with Row Level Security in Supabase
- Modern, responsive UI

## Technical Architecture

### Authentication

The application uses GitHub OAuth for authentication, providing a secure and seamless user experience. The authentication flow is as follows:

1. User clicks "Login with GitHub" button
2. User is redirected to GitHub to authorize the application
3. GitHub redirects back to the app with an authorization code
4. The app exchanges the code for an access token
5. The app retrieves the user information from GitHub
6. The app creates or retrieves the user in the database
7. The app creates a session for the user

### Data Storage

All todo entries and user data are stored in Supabase, utilizing PostgreSQL with the following features:

- Row Level Security (RLS) for data isolation between users
- SQL queries for CRUD operations
- Prepared statements to prevent SQL injection
- Transaction support for data integrity
- Connection pooling for performance

## Prerequisites

- Go 1.21 or higher
- Supabase account (for database and authentication)
- GitHub OAuth App credentials

## Setup

1. Clone the repository:

```bash
git clone https://github.com/starbops/gottodo.git
cd gottodo
```

2. Copy the example environment file and update it with your configuration:

```bash
cp .env.example .env
```

3. Update the `.env` file with your Supabase credentials and GitHub OAuth app credentials:

```
# Supabase configuration
SUPABASE_URL=https://your-project-id.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_DB_URL=postgres://postgres:postgres@localhost:5432/gottodo

# GitHub OAuth configuration
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
GITHUB_REDIRECT_URL=http://localhost:8080/auth/github/callback
```

4. Create a GitHub OAuth app:
   - Go to GitHub Settings > Developer settings > OAuth Apps > New OAuth App
   - Set the callback URL to http://localhost:8080/auth/github/callback
   - Copy the Client ID and Client Secret to your .env file

5. Set up your Supabase database:
   - Create a new Supabase project
   - Navigate to the SQL Editor
   - Run the migration script from `migrations/001_create_todos_table.sql`
   
   **Important**: The migration script creates a `todos` table with UUIDs for IDs. Make sure 
   your auth configurations are set up correctly:
   
   - In the Supabase Dashboard, go to Authentication > Settings > General
   - Ensure "Service Role JWT Templates" has the correct claim for user IDs
   - Make sure user IDs are stored as UUIDs in your Supabase Auth settings

   The application expects user IDs to be valid UUIDs to match with Supabase's `auth.uid()` function.

6. Install dependencies:

```bash
go mod tidy
```

7. Run the application:

```bash
go run cmd/server/main.go
```

8. Open your browser and navigate to `http://localhost:8080`

## Troubleshooting Supabase Setup

If you encounter errors like:

```
ERROR: 42883: operator does not exist: character varying = uuid
HINT: No operator matches the given name and argument types. You might need to add explicit type casts.
```

This indicates a type mismatch in the SQL script. Make sure:

1. Your `user_id` column in the todos table is defined as a UUID type
2. The RLS policy is comparing two values of the same type
3. Supabase's auth.uid() function returns a UUID that matches the user_id column type

If you need to convert data types in existing policies:
```sql
-- Example of explicit type casting if needed
CREATE POLICY todo_user_policy ON todos
    USING (user_id::text = auth.uid()::text)
    WITH CHECK (user_id::text = auth.uid()::text);
```

## Project Structure

```
gottodo/
├── cmd/                    # Application entry points
│   └── server/             # Web server
├── internal/               # Private application code
│   ├── handlers/           # HTTP handlers
│   ├── middleware/         # HTTP middleware
│   ├── models/             # Data models
│   ├── repositories/       # Data access layer
│   └── services/           # Business logic
├── pkg/                    # Public libraries
│   ├── auth/               # Authentication
│   ├── database/           # Database connection
│   └── config/             # Configuration
├── ui/                     # User interface
│   ├── components/         # UI components
│   ├── layouts/            # Page layouts
│   ├── pages/              # Page templates
│   └── static/             # Static assets
└── migrations/             # Database migrations
```

## Testing

The application includes comprehensive unit tests for critical components. To run the tests:

```bash
make test
```

For test coverage:

```bash
make test-coverage
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
