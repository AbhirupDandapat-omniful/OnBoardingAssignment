CREATE TABLE tenants (
  id          UUID        PRIMARY KEY,
  name        TEXT        NOT NULL,
  metadata    JSONB       NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
