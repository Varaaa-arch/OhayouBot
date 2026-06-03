CREATE TABLE whatsapp_sessions (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_data TEXT NOT NULL,
    created_at   TIMESTAMPTZ DEFAULT now()
);
