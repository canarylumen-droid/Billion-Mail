-- Cold Outreach Platform Schema Migration
-- Run this against the BillionMail PostgreSQL database

-- ============================================================
-- Sequences table
-- ============================================================
CREATE TABLE IF NOT EXISTS bm_sequences (
    id           SERIAL PRIMARY KEY,
    name         VARCHAR(200)  NOT NULL,
    description  TEXT          NOT NULL DEFAULT '',
    status       VARCHAR(20)   NOT NULL DEFAULT 'draft',  -- draft|active|paused
    steps        JSONB         NOT NULL DEFAULT '[]',
    schedule     JSONB         NOT NULL DEFAULT '{}',
    calendly_url VARCHAR(500)  NOT NULL DEFAULT '',
    create_time  INTEGER       NOT NULL DEFAULT 0,
    update_time  INTEGER       NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_bm_sequences_status ON bm_sequences(status);

-- ============================================================
-- Leads table (cold outreach contacts, separate from newsletter contacts)
-- ============================================================
CREATE TABLE IF NOT EXISTS bm_leads (
    id             SERIAL PRIMARY KEY,
    email          VARCHAR(255) NOT NULL,
    first_name     VARCHAR(100) NOT NULL DEFAULT '',
    last_name      VARCHAR(100) NOT NULL DEFAULT '',
    company        VARCHAR(200) NOT NULL DEFAULT '',
    title          VARCHAR(200) NOT NULL DEFAULT '',
    phone          VARCHAR(50)  NOT NULL DEFAULT '',
    linkedin       VARCHAR(500) NOT NULL DEFAULT '',
    website        VARCHAR(500) NOT NULL DEFAULT '',
    attributes     JSONB        NOT NULL DEFAULT '{}',
    status         VARCHAR(30)  NOT NULL DEFAULT 'not_started', -- not_started|in_sequence|replied|interested|bounced|unsubscribed
    group_id       INTEGER      NOT NULL DEFAULT 0,
    sequence_id    INTEGER      NOT NULL DEFAULT 0,
    sequence_step  VARCHAR(50)  NOT NULL DEFAULT '',
    last_contacted INTEGER      NOT NULL DEFAULT 0,
    next_send      INTEGER      NOT NULL DEFAULT 0,
    create_time    INTEGER      NOT NULL DEFAULT 0,
    update_time    INTEGER      NOT NULL DEFAULT 0,
    UNIQUE(email)
);

CREATE INDEX IF NOT EXISTS idx_bm_leads_status      ON bm_leads(status);
CREATE INDEX IF NOT EXISTS idx_bm_leads_group_id    ON bm_leads(group_id);
CREATE INDEX IF NOT EXISTS idx_bm_leads_sequence_id ON bm_leads(sequence_id);
CREATE INDEX IF NOT EXISTS idx_bm_leads_email       ON bm_leads(email);

-- ============================================================
-- Sequence leads enrollment (many-to-many, tracks progress)
-- ============================================================
CREATE TABLE IF NOT EXISTS bm_sequence_leads (
    id            SERIAL PRIMARY KEY,
    sequence_id   INTEGER      NOT NULL,
    contact_id    INTEGER      NOT NULL,
    email         VARCHAR(255) NOT NULL,
    status        VARCHAR(30)  NOT NULL DEFAULT 'not_started',
    current_step  INTEGER      NOT NULL DEFAULT 0,
    next_send_time INTEGER     NOT NULL DEFAULT 0,
    stopped_at    INTEGER      NOT NULL DEFAULT 0,
    stop_reason   VARCHAR(50)  NOT NULL DEFAULT '',  -- replied|bounced|unsubscribed|manual
    create_time   INTEGER      NOT NULL DEFAULT 0,
    update_time   INTEGER      NOT NULL DEFAULT 0,
    UNIQUE(sequence_id, contact_id)
);

CREATE INDEX IF NOT EXISTS idx_bm_seq_leads_seq     ON bm_sequence_leads(sequence_id);
CREATE INDEX IF NOT EXISTS idx_bm_seq_leads_status  ON bm_sequence_leads(status);
CREATE INDEX IF NOT EXISTS idx_bm_seq_leads_next    ON bm_sequence_leads(next_send_time);

-- ============================================================
-- Sequence send logs (for analytics)
-- ============================================================
CREATE TABLE IF NOT EXISTS bm_sequence_logs (
    id           SERIAL PRIMARY KEY,
    sequence_id  INTEGER      NOT NULL,
    lead_id      INTEGER      NOT NULL,
    step_type    VARCHAR(30)  NOT NULL DEFAULT '',  -- initial|followup1|followup2
    message_id   VARCHAR(200) NOT NULL DEFAULT '',
    status       VARCHAR(30)  NOT NULL DEFAULT 'sent',  -- sent|delivered|bounced
    opened       SMALLINT     NOT NULL DEFAULT 0,
    clicked      SMALLINT     NOT NULL DEFAULT 0,
    replied      SMALLINT     NOT NULL DEFAULT 0,
    from_email   VARCHAR(255) NOT NULL DEFAULT '',
    sent_at      INTEGER      NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_bm_seq_logs_seq      ON bm_sequence_logs(sequence_id);
CREATE INDEX IF NOT EXISTS idx_bm_seq_logs_lead     ON bm_sequence_logs(lead_id);
CREATE INDEX IF NOT EXISTS idx_bm_seq_logs_sent_at  ON bm_sequence_logs(sent_at);

-- ============================================================
-- CRM inbox (inbound replies from leads)
-- ============================================================
CREATE TABLE IF NOT EXISTS bm_crm_inbox (
    id           SERIAL PRIMARY KEY,
    from_email   VARCHAR(255) NOT NULL,
    from_name    VARCHAR(200) NOT NULL DEFAULT '',
    to_email     VARCHAR(255) NOT NULL DEFAULT '',
    subject      VARCHAR(500) NOT NULL DEFAULT '',
    body_text    TEXT         NOT NULL DEFAULT '',
    body_html    TEXT         NOT NULL DEFAULT '',
    received_at  INTEGER      NOT NULL DEFAULT 0,
    status       VARCHAR(30)  NOT NULL DEFAULT 'unread',  -- unread|read|replied|archived
    lead_id      INTEGER      NOT NULL DEFAULT 0,
    sequence_id  INTEGER      NOT NULL DEFAULT 0,
    thread_id    VARCHAR(200) NOT NULL DEFAULT '',
    has_calendly SMALLINT     NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_bm_crm_inbox_status ON bm_crm_inbox(status);
CREATE INDEX IF NOT EXISTS idx_bm_crm_inbox_lead   ON bm_crm_inbox(lead_id);
CREATE INDEX IF NOT EXISTS idx_bm_crm_inbox_recv   ON bm_crm_inbox(received_at);

-- ============================================================
-- Mailbox daily send counters (for per-mailbox daily limits)
-- ============================================================
CREATE TABLE IF NOT EXISTS bm_mailbox_daily_counts (
    id          SERIAL PRIMARY KEY,
    mailbox     VARCHAR(255) NOT NULL,
    date_key    VARCHAR(10)  NOT NULL,  -- YYYY-MM-DD
    sent_count  INTEGER      NOT NULL DEFAULT 0,
    UNIQUE(mailbox, date_key)
);
