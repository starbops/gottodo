# Supabase Setup Guide for GotToDo

This document explains how to set up Supabase correctly for the GotToDo application, particularly focusing on UUID handling, which is critical for the application to work properly.

## Initial Setup

1. Create a Supabase account at [supabase.com](https://supabase.com) if you don't have one already.
2. Create a new project in Supabase.
3. Note your project URL and anon key, which you'll need for your `.env` file.

## Database Migration

Run the migration script from `migrations/001_create_todos_table.sql` in the Supabase SQL Editor:

```sql
-- Create extension for UUID support
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create todos table
CREATE TABLE IF NOT EXISTS todos (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create index on user_id for better query performance
CREATE INDEX IF NOT EXISTS idx_todos_user_id ON todos(user_id);

-- Add RLS (Row Level Security) policies
ALTER TABLE todos ENABLE ROW LEVEL SECURITY;

-- Create policy to ensure users can only see and manipulate their own todos
CREATE POLICY todo_user_policy ON todos
    USING (user_id = auth.uid())
    WITH CHECK (user_id = auth.uid());
```

## UUID Handling in Supabase

The most common issue when setting up the application with Supabase is related to UUID handling:

### Understanding the Error

If you see an error like:

```
ERROR: 42883: operator does not exist: character varying = uuid
HINT: No operator matches the given name and argument types. You might need to add explicit type casts.
```

This occurs because:

1. `auth.uid()` in Supabase returns a UUID value
2. But your `user_id` column might be defined as a different type (e.g., VARCHAR)

### How We Fixed It

1. Changed the `user_id` column type from VARCHAR to UUID:
   ```sql
   user_id UUID NOT NULL
   ```

2. Updated the application to handle UUIDs correctly:
   - Parse user IDs and todo IDs as UUIDs before sending to the database
   - Convert UUIDs back to strings when returning data from the database
   - Added validation to ensure all user IDs are valid UUIDs before using them

3. Added explicit UUID validation throughout the application:
   ```go
   // IsValidUUID checks if a string is a valid UUID
   func IsValidUUID(id string) bool {
       _, err := uuid.Parse(id)
       return err == nil
   }
   ```

## Authentication Setup

For GitHub OAuth to work correctly with Supabase:

1. Enable GitHub OAuth in Supabase:
   - Go to Authentication > Providers > GitHub
   - Toggle it on and enter your GitHub OAuth credentials

2. Make sure Supabase JWT settings match your application:
   - Go to Authentication > Settings
   - Ensure the JWT fields are properly configured

## Row Level Security (RLS)

The application leverages Supabase's Row Level Security to ensure data isolation between users:

1. RLS is enabled on the todos table:
   ```sql
   ALTER TABLE todos ENABLE ROW LEVEL SECURITY;
   ```

2. The policy ensures users can only see and modify their own todos:
   ```sql
   CREATE POLICY todo_user_policy ON todos
       USING (user_id = auth.uid())
       WITH CHECK (user_id = auth.uid());
   ```

3. The user_id column must be a UUID type to match auth.uid():
   ```sql
   user_id UUID NOT NULL
   ```

## Troubleshooting

If you still encounter issues:

1. Check the types in your database:
   ```sql
   SELECT column_name, data_type 
   FROM information_schema.columns 
   WHERE table_name = 'todos';
   ```

2. Check the type of auth.uid():
   ```sql
   SELECT pg_typeof(auth.uid());
   ```

3. If needed, you can use explicit type casting in your policies:
   ```sql
   CREATE POLICY todo_user_policy ON todos
       USING (user_id::text = auth.uid()::text)
       WITH CHECK (user_id::text = auth.uid()::text);
   ```

## Summary

The key points for making the application work with Supabase are:

1. Use UUID types consistently across your database and application
2. Ensure your application correctly handles UUID conversion
3. Set up Row Level Security to match the authentication system
4. Use the correct authentication flow for GitHub OAuth

If you follow these guidelines, you should have a functioning GotToDo application with proper data isolation and authentication. 