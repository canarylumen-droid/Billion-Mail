-- Multi-Provider Relay Pool
-- 15 providers × 2 accounts = 30 slots, ~193k emails/month free combined

CREATE TABLE IF NOT EXISTS bm_relay_providers (
    id                  BIGSERIAL PRIMARY KEY,
    name                VARCHAR(100) NOT NULL,
    provider_type       VARCHAR(30) NOT NULL,
    slot                INT NOT NULL DEFAULT 1,        -- 1 or 2 (two accounts per provider)
    api_key             VARCHAR(500) NOT NULL DEFAULT '',
    api_key_2           VARCHAR(500) NOT NULL DEFAULT '', -- secondary credential (e.g. domain for Mailgun)
    daily_limit         INT NOT NULL DEFAULT 100,
    monthly_limit       INT NOT NULL DEFAULT 1000,
    daily_sent          INT NOT NULL DEFAULT 0,
    monthly_sent        INT NOT NULL DEFAULT 0,
    last_daily_reset    BIGINT NOT NULL DEFAULT 0,
    last_monthly_reset  BIGINT NOT NULL DEFAULT 0,
    priority            INT NOT NULL DEFAULT 0,
    status              VARCHAR(20) NOT NULL DEFAULT 'unconfigured', -- unconfigured, active, exhausted, error
    is_active           BOOLEAN NOT NULL DEFAULT FALSE,
    last_error          TEXT NOT NULL DEFAULT '',
    last_used_at        BIGINT NOT NULL DEFAULT 0,
    created_at          BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
    updated_at          BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
    UNIQUE(provider_type, slot)
);

CREATE INDEX IF NOT EXISTS idx_relay_providers_active ON bm_relay_providers(is_active, status);
CREATE INDEX IF NOT EXISTS idx_relay_providers_type ON bm_relay_providers(provider_type);

-- Seed 30 default rows (15 providers × 2 slots) with correct free tier limits
-- Users just fill in the api_key field to activate each slot

INSERT INTO bm_relay_providers (name, provider_type, slot, daily_limit, monthly_limit, priority) VALUES
  ('Resend (Account 1)',         'resend',       1, 100,  3000,  90),
  ('Resend (Account 2)',         'resend',       2, 100,  3000,  90),
  ('SendGrid (Account 1)',       'sendgrid',     1, 100,  3000,  85),
  ('SendGrid (Account 2)',       'sendgrid',     2, 100,  3000,  85),
  ('Brevo (Account 1)',          'brevo',        1, 300,  9000,  80),
  ('Brevo (Account 2)',          'brevo',        2, 300,  9000,  80),
  ('Mailjet (Account 1)',        'mailjet',      1, 200,  6000,  75),
  ('Mailjet (Account 2)',        'mailjet',      2, 200,  6000,  75),
  ('Mailersend (Account 1)',     'mailersend',   1, 100,  3000,  70),
  ('Mailersend (Account 2)',     'mailersend',   2, 100,  3000,  70),
  ('Sendpulse (Account 1)',      'sendpulse',    1, 500,  15000, 65),
  ('Sendpulse (Account 2)',      'sendpulse',    2, 500,  15000, 65),
  ('ZeptoMail (Account 1)',      'zeptomail',    1, 333,  10000, 60),
  ('ZeptoMail (Account 2)',      'zeptomail',    2, 333,  10000, 60),
  ('SMTP2GO (Account 1)',        'smtp2go',      1, 33,   1000,  55),
  ('SMTP2GO (Account 2)',        'smtp2go',      2, 33,   1000,  55),
  ('Elastic Email (Account 1)',  'elasticemail', 1, 100,  3000,  50),
  ('Elastic Email (Account 2)',  'elasticemail', 2, 100,  3000,  50),
  ('Mailgun (Account 1)',        'mailgun',      1, 33,   1000,  45),
  ('Mailgun (Account 2)',        'mailgun',      2, 33,   1000,  45),
  ('Postmark (Account 1)',       'postmark',     1, 4,    100,   40),
  ('Postmark (Account 2)',       'postmark',     2, 4,    100,   40),
  ('SparkPost (Account 1)',      'sparkpost',    1, 17,   500,   35),
  ('SparkPost (Account 2)',      'sparkpost',    2, 17,   500,   35),
  ('SocketLabs (Account 1)',     'socketlabs',   1, 1333, 40000, 30),
  ('SocketLabs (Account 2)',     'socketlabs',   2, 1333, 40000, 30),
  ('Netcore (Account 1)',        'netcore',      1, 100,  3000,  25),
  ('Netcore (Account 2)',        'netcore',      2, 100,  3000,  25),
  ('Mailtrap (Account 1)',       'mailtrap',     1, 33,   1000,  20),
  ('Mailtrap (Account 2)',       'mailtrap',     2, 33,   1000,  20)
ON CONFLICT (provider_type, slot) DO NOTHING;
