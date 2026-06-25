-- ESP (Email Service Provider) Schema
-- BillionMail SMTP Relay, Domain Management, IP Pools, Delivery Analytics

-- SMTP relay credentials per user/account
CREATE TABLE IF NOT EXISTS bm_smtp_credentials (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    username    VARCHAR(100) NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL,
    description TEXT,
    status      SMALLINT NOT NULL DEFAULT 1, -- 1=active, 0=disabled
    daily_limit INT NOT NULL DEFAULT 1000,
    sent_today  INT NOT NULL DEFAULT 0,
    plan        VARCHAR(20) NOT NULL DEFAULT 'free', -- free, starter, pro
    ip_pool_id  BIGINT,
    created_at  BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
    updated_at  BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT
);

-- Sending domains (verified domains users send from)
CREATE TABLE IF NOT EXISTS bm_sending_domains (
    id               BIGSERIAL PRIMARY KEY,
    domain           VARCHAR(255) NOT NULL UNIQUE,
    status           VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, verified, failed
    dkim_selector    VARCHAR(50) NOT NULL DEFAULT 'bm1',
    dkim_private_key TEXT,
    dkim_public_key  TEXT,
    spf_record       TEXT,
    dmarc_record     TEXT,
    spf_verified     BOOLEAN NOT NULL DEFAULT FALSE,
    dkim_verified    BOOLEAN NOT NULL DEFAULT FALSE,
    dmarc_verified   BOOLEAN NOT NULL DEFAULT FALSE,
    last_checked_at  BIGINT,
    created_at       BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
    updated_at       BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT
);

-- IP pools (shared vs dedicated)
CREATE TABLE IF NOT EXISTS bm_ip_pools (
    id           BIGSERIAL PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    pool_type    VARCHAR(20) NOT NULL DEFAULT 'shared', -- shared, dedicated
    ips          JSONB NOT NULL DEFAULT '[]',
    status       VARCHAR(20) NOT NULL DEFAULT 'warming', -- warming, active, suspended
    warm_day     INT NOT NULL DEFAULT 0,
    daily_limit  INT NOT NULL DEFAULT 500,
    sent_today   INT NOT NULL DEFAULT 0,
    description  TEXT,
    created_at   BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
    updated_at   BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT
);

-- IP warming schedule (per IP per day limits)
CREATE TABLE IF NOT EXISTS bm_ip_warming (
    id          BIGSERIAL PRIMARY KEY,
    ip_pool_id  BIGINT NOT NULL REFERENCES bm_ip_pools(id) ON DELETE CASCADE,
    ip_address  VARCHAR(45) NOT NULL,
    day_number  INT NOT NULL DEFAULT 1,
    daily_limit INT NOT NULL DEFAULT 100,
    sent_count  INT NOT NULL DEFAULT 0,
    date_str    VARCHAR(10) NOT NULL, -- YYYY-MM-DD
    created_at  BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT
);

-- Delivery statistics (per domain per day)
CREATE TABLE IF NOT EXISTS bm_delivery_stats (
    id              BIGSERIAL PRIMARY KEY,
    credential_id   BIGINT,
    domain          VARCHAR(255),
    date_str        VARCHAR(10) NOT NULL, -- YYYY-MM-DD
    sent            INT NOT NULL DEFAULT 0,
    delivered       INT NOT NULL DEFAULT 0,
    bounced_soft    INT NOT NULL DEFAULT 0,
    bounced_hard    INT NOT NULL DEFAULT 0,
    complained      INT NOT NULL DEFAULT 0,
    opened          INT NOT NULL DEFAULT 0,
    clicked         INT NOT NULL DEFAULT 0,
    unsubscribed    INT NOT NULL DEFAULT 0,
    created_at      BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
    UNIQUE(credential_id, domain, date_str)
);

-- Suppression list (bounces, complaints, unsubscribes)
CREATE TABLE IF NOT EXISTS bm_suppression_list (
    id           BIGSERIAL PRIMARY KEY,
    email        VARCHAR(255) NOT NULL,
    reason       VARCHAR(20) NOT NULL, -- hard_bounce, soft_bounce, complaint, unsubscribe, manual
    domain       VARCHAR(255),
    raw_message  TEXT,
    created_at   BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
    UNIQUE(email, reason)
);

-- Haraka MTA server registry
CREATE TABLE IF NOT EXISTS bm_mta_servers (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    host        VARCHAR(255) NOT NULL,
    ssh_port    INT NOT NULL DEFAULT 22,
    ssh_user    VARCHAR(100) NOT NULL DEFAULT 'root',
    ssh_key     TEXT,
    status      VARCHAR(20) NOT NULL DEFAULT 'active', -- active, maintenance, offline
    region      VARCHAR(50),
    ip_count    INT NOT NULL DEFAULT 1,
    created_at  BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
    updated_at  BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT
);

CREATE INDEX IF NOT EXISTS idx_smtp_creds_username ON bm_smtp_credentials(username);
CREATE INDEX IF NOT EXISTS idx_sending_domains_domain ON bm_sending_domains(domain);
CREATE INDEX IF NOT EXISTS idx_delivery_stats_date ON bm_delivery_stats(date_str);
CREATE INDEX IF NOT EXISTS idx_suppression_email ON bm_suppression_list(email);
CREATE INDEX IF NOT EXISTS idx_suppression_domain ON bm_suppression_list(domain);
