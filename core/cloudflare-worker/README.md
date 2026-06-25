# BillionMail Cloudflare Email Relay Worker

Zero-cost email delivery using Cloudflare's globally trusted IP infrastructure.
No VPS. No port 25. No external ESP dependency.

## How It Works

```
BillionMail (Railway/any host)
    │
    │  HTTPS POST (no port 25 needed)
    ▼
Cloudflare Worker (your account, free tier)
    │
    │  send_email() — Cloudflare's trusted IPs
    ▼
Gmail / Outlook / Yahoo  ✅ Inbox
```

## Free Tier Limits

| Resource | Free Limit |
|---|---|
| Worker requests | 100,000/day |
| CPU time per request | 10ms |
| Email sends | Tied to Cloudflare Email Routing limits |

## Setup (5 minutes)

### Step 1: Enable Email Routing on your domain
1. Go to [dash.cloudflare.com](https://dash.cloudflare.com)
2. Select your domain → **Email** → **Email Routing**
3. Enable it (adds MX records automatically)

### Step 2: Deploy the Worker
1. Go to **Workers & Pages** → **Create Worker**
2. Paste the contents of `email-relay.js`
3. Click **Deploy**

### Step 3: Add Email Binding
1. In your Worker → **Settings** → **Bindings**
2. Add **Send Email** binding → name it `EMAIL`
3. Set allowed sender: your verified domain

### Step 4: Add Secret Variable
1. In your Worker → **Settings** → **Environment Variables**
2. Add: `BILLIONMAIL_SECRET` = (any random string, e.g. `openssl rand -hex 32`)

### Step 5: Connect to BillionMail
1. Copy your Worker URL: `https://your-worker.your-subdomain.workers.dev`
2. Go to BillionMail → **Delivery Backend** → **Cloudflare Workers**
3. Paste the URL + your secret → **Test Connection** → **Save**

## Oracle Cloud Free VM Alternative (Higher Volume)

For higher volume (unlimited sends), use Oracle Cloud Always Free:
- 2 AMD VMs free forever (unlike AWS 12-month trial)
- Port 25 available on request
- Install Haraka: `npm install -g Haraka && haraka -i /etc/haraka`

See the MTA Servers tab in BillionMail for full setup guide.
