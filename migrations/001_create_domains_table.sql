-- Migration: Create domains table
-- Created: 2025-09-05

CREATE TABLE IF NOT EXISTS domains (
    domain_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    domain VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on domain for faster lookups
CREATE INDEX IF NOT EXISTS idx_domains_domain ON domains(domain);

-- Create index on name for faster searches
CREATE INDEX IF NOT EXISTS idx_domains_name ON domains(name);
