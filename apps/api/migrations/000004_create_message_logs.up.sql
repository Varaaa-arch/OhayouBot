CREATE TABLE message_logs (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contact_id UUID REFERENCES contacts(id) ON DELETE CASCADE,
    message    TEXT NOT NULL,
    status     TEXT NOT NULL,
    sent_at    TIMESTAMPTZ DEFAULT now()
);
