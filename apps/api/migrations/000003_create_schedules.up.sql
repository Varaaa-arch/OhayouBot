CREATE TABLE schedules (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contact_id UUID REFERENCES contacts(id) ON DELETE CASCADE,
    message    TEXT NOT NULL,
    time       TEXT NOT NULL,
    enabled    BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT now()
);
