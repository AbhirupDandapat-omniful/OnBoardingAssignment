CREATE TABLE skus (
  id            UUID        PRIMARY KEY,
  tenant_id     UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  seller_id     UUID        NOT NULL REFERENCES sellers(id) ON DELETE CASCADE,
  code          TEXT        NOT NULL UNIQUE,
  name          TEXT        NOT NULL,
  description   TEXT        NULL,
  category_id   UUID        NULL REFERENCES categories(id),
  weight        NUMERIC     NULL,
  weight_unit   TEXT        NULL,
  length        NUMERIC     NULL,
  width         NUMERIC     NULL,
  height        NUMERIC     NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
