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

-- Downgrade
-- DROP TABLE IF EXISTS todos;
-- DROP EXTENSION IF EXISTS "uuid-ossp"; 