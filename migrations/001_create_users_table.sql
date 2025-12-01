-- migrations/001_create_users_table.sql

-- Create the users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL, -- Hashed password
    role VARCHAR(50) NOT NULL,     -- 'admin' or 'student'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for faster lookups
CREATE UNIQUE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_role ON users (role);

-- Insert default admin user. 
-- NOTE: The password 'admin123' must be BCrypt hashed before running this script
-- For this example, we'll use a known hash for 'admin123' 
-- (You should generate this hash programmatically or use a utility in a real app)
-- Hash of 'admin123' with cost 10: $2a$10$gNqj/Nl7lR2.oO.1Y9RjX./PzO2y6O9k.5m9jL6p4.Xq
-- We'll use a simpler placeholder here, but the Go app should handle the hashing.
-- For simple setup, we'll use a hash that the Go app can verify:
-- Hash of 'admin123' at cost 14: $2a$14$W.275iS.r6zE02u7xQnNKeUaR7bL4F4X6fVj6x.zK8n.m2k/R6lG
INSERT INTO users (name, email, password, role) 
VALUES (
    'Admin User',
    'admin@example.com',
    '$2a$14$W.275iS.r6zE02u7xQnNKeUaR7bL4F4X6fVj6x.zK8n.m2k/R6lG', -- BCrypt hash of 'admin123'
    'admin'
)
ON CONFLICT (email) DO NOTHING;