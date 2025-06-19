CREATE TABLE hubs (
  id            UUID        PRIMARY KEY,
  tenant_id     UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  seller_id     UUID        NOT NULL REFERENCES sellers(id) ON DELETE CASCADE,
  name          TEXT        NOT NULL,
  location      TEXT        NOT NULL,
  address       TEXT        NULL,
  contact_email TEXT        NULL,
  contact_phone TEXT        NULL,
  timezone      TEXT        NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
