CREATE TABLE webhooks (
  id            UUID          PRIMARY KEY,
  tenant_id     UUID          NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  callback_url  TEXT          NOT NULL,
  events        TEXT[]        NOT NULL,
  headers       JSONB         NULL,
  is_active     BOOLEAN       NOT NULL DEFAULT TRUE,
  created_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);
