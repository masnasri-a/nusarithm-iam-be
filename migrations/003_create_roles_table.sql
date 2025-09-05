-- Migration: Create roles table
-- Created: 2025-09-05

CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    domain_id UUID NOT NULL REFERENCES domains(domain_id) ON DELETE CASCADE,
    role_name VARCHAR(255) NOT NULL,
    role_claims JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(domain_id, role_name)
);

-- Create index on domain_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_roles_domain_id ON roles(domain_id);

-- Create index on role_name for faster searches
CREATE INDEX IF NOT EXISTS idx_roles_role_name ON roles(role_name);

-- Create index on role_claims for JSON queries
CREATE INDEX IF NOT EXISTS idx_roles_claims ON roles USING GIN (role_claims);
