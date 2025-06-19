CREATE TABLE inventory_transactions (
  id               UUID        PRIMARY KEY,
  tenant_id        UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  hub_id           UUID        NOT NULL REFERENCES hubs(id) ON DELETE CASCADE,
  sku_id           UUID        NOT NULL REFERENCES skus(id) ON DELETE CASCADE,
  delta            BIGINT      NOT NULL,
  transaction_type TEXT        NOT NULL,
  reference_id     TEXT        NULL,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
