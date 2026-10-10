-- 001_traditional_schema.sql
-- Baseline 3-tier relational schema for traditional task management

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Projects Table
CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    lead VARCHAR(100) NOT NULL DEFAULT '',
    target_date DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Tasks Table
CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    priority VARCHAR(50) NOT NULL DEFAULT 'Medium' CHECK (priority IN ('Low', 'Medium', 'High', 'Critical')),
    status VARCHAR(50) NOT NULL DEFAULT 'Pending' CHECK (status IN ('Pending', 'InProgress', 'Completed')),
    due_date DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for frequent query paths
CREATE INDEX IF NOT EXISTS idx_tasks_project_id ON tasks(project_id);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);

-- Helper trigger for automatic updated_at timestamps
CREATE OR REPLACE FUNCTION update_timestamp_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TRIGGER IF EXISTS trg_projects_updated_at ON projects;
CREATE TRIGGER trg_projects_updated_at
    BEFORE UPDATE ON projects
    FOR EACH ROW
    EXECUTE FUNCTION update_timestamp_column();

DROP TRIGGER IF EXISTS trg_tasks_updated_at ON tasks;
CREATE TRIGGER trg_tasks_updated_at
    BEFORE UPDATE ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION update_timestamp_column();

-- Seed Sample Data for immediate testing
INSERT INTO projects (id, name, description, lead, target_date)
VALUES 
    ('11111111-1111-1111-1111-111111111111', 'Cloud Migration 2026', 'Migrate legacy microservices to Kubernetes', 'Alex Rivera', CURRENT_DATE + INTERVAL '30 days'),
    ('22222222-2222-2222-2222-222222222222', 'Observability Stack', 'Deploy Prometheus and Grafana dashboards', 'Devon Vance', CURRENT_DATE + INTERVAL '14 days')
ON CONFLICT (id) DO NOTHING;

INSERT INTO tasks (id, project_id, title, description, priority, status, due_date)
VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111', 'Containerize Auth Service', 'Write multi-stage Dockerfile and test locally', 'High', 'Completed', CURRENT_DATE - INTERVAL '2 days'),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '11111111-1111-1111-1111-111111111111', 'Provision Managed PostgreSQL', 'Set up RDS/Cloud SQL instance with replica', 'Critical', 'InProgress', CURRENT_DATE + INTERVAL '5 days'),
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', '22222222-2222-2222-2222-222222222222', 'Configure AlertManager Webhooks', 'Route severity=critical alerts to PagerDuty', 'Medium', 'Pending', CURRENT_DATE + INTERVAL '10 days')
ON CONFLICT (id) DO NOTHING;
