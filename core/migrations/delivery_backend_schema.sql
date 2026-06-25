-- Delivery Backend Configuration
-- Stores which delivery method BillionMail uses to actually send emails

CREATE TABLE IF NOT EXISTS bm_delivery_backends (
    id           BIGSERIAL PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    backend_type VARCHAR(30) NOT NULL, -- cloudflare_worker, haraka, oracle_vm, smtp_relay
    config       JSONB NOT NULL DEFAULT '{}',
    status       VARCHAR(20) NOT NULL DEFAULT 'active', -- active, disabled, testing
    priority     INT NOT NULL DEFAULT 0,
    daily_sent   INT NOT NULL DEFAULT 0,
    daily_limit  INT NOT NULL DEFAULT 100000,
    is_default   BOOLEAN NOT NULL DEFAULT FALSE,
    last_tested  BIGINT,
    last_error   TEXT,
    created_at   BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
    updated_at   BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT
);

CREATE INDEX IF NOT EXISTS idx_delivery_backends_type ON bm_delivery_backends(backend_type);
CREATE INDEX IF NOT EXISTS idx_delivery_backends_default ON bm_delivery_backends(is_default);
