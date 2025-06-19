CREATE TABLE inventory (
  hub_id            UUID        NOT NULL REFERENCES hubs(id) ON DELETE CASCADE,
  sku_id            UUID        NOT NULL REFERENCES skus(id) ON DELETE CASCADE,
  quantity_on_hand  BIGINT      NOT NULL DEFAULT 0,
  quantity_reserved BIGINT      NOT NULL DEFAULT 0,
  min_threshold     BIGINT      NOT NULL DEFAULT 0,
  max_threshold     BIGINT      NOT NULL DEFAULT 0,
  updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (hub_id, sku_id)
);
