# BillionMail Relay Provider Pool — Complete Reference

15 providers × 2 accounts = **30 API slots** giving up to **~193,200 emails/month free** combined.

---

## Combined Free Capacity (2 accounts each)

| Provider | Daily/account | Monthly/account | × 2 accounts | Monthly Total |
|---|---|---|---|---|
| Resend | 100 | 3,000 | ×2 | 6,000 |
| SendGrid | 100 | 3,000 | ×2 | 6,000 |
| Brevo | 300 | 9,000 | ×2 | 18,000 |
| Mailjet | 200 | 6,000 | ×2 | 12,000 |
| Mailersend | 100 | 3,000 | ×2 | 6,000 |
| Sendpulse | 500 | 15,000 | ×2 | 30,000 |
| ZeptoMail | 333 | 10,000 | ×2 | 20,000 |
| SMTP2GO | 33 | 1,000 | ×2 | 2,000 |
| Elastic Email | 100 | 3,000 | ×2 | 6,000 |
| Mailgun | 33 | 1,000 | ×2 | 2,000 |
| Postmark | 4 | 100 | ×2 | 200 |
| SparkPost | 17 | 500 | ×2 | 1,000 |
| SocketLabs | 1,333 | 40,000 | ×2 | 80,000 |
| Netcore | 100 | 3,000 | ×2 | 6,000 |
| Mailtrap | 33 | 1,000 | ×2 | 2,000 |
| **TOTAL** | **3,287/day** | **~99k/month** | — | **~197,200/month** |

> SocketLabs free trial gives 40,000 for the first month. After trial it drops to ~1,000/month paid.
> Numbers are approximate — providers adjust free tiers over time.

---

## Provider Details & Environment Variables

### 1. Resend
- **Free tier:** 3,000 emails/month, 100/day
- **Sign up:** https://resend.com
- **Get API key:** resend.com/api-keys → Create API Key
- **API key format:** `re_xxxxxxxxxxxx`
- **Env vars:**
  ```
  RELAY_RESEND_1_KEY=re_your_api_key_here
  RELAY_RESEND_2_KEY=re_your_second_api_key_here
  ```

### 2. SendGrid
- **Free tier:** 3,000 emails/month (100/day forever free)
- **Sign up:** https://signup.sendgrid.com
- **Get API key:** Settings → API Keys → Create API Key (Full Access or Mail Send only)
- **API key format:** `SG.xxxxxxxxxxxx`
- **Env vars:**
  ```
  RELAY_SENDGRID_1_KEY=SG.your_api_key_here
  RELAY_SENDGRID_2_KEY=SG.your_second_api_key_here
  ```

### 3. Brevo (formerly Sendinblue)
- **Free tier:** 9,000 emails/month (300/day)
- **Sign up:** https://app.brevo.com/account/register
- **Get API key:** Account → SMTP & API → API Keys → Generate a new API key
- **Env vars:**
  ```
  RELAY_BREVO_1_KEY=xkeysib-your_api_key_here
  RELAY_BREVO_2_KEY=xkeysib-your_second_api_key_here
  ```

### 4. Mailjet
- **Free tier:** 6,000 emails/month (200/day)
- **Sign up:** https://app.mailjet.com/signup
- **Get API key:** Account Settings → REST API → API Key Management
- **Note:** Requires both Public Key (API Key) and Secret Key
- **Env vars:**
  ```
  RELAY_MAILJET_1_KEY=your_public_key
  RELAY_MAILJET_1_KEY2=your_secret_key
  RELAY_MAILJET_2_KEY=your_second_public_key
  RELAY_MAILJET_2_KEY2=your_second_secret_key
  ```

### 5. Mailersend
- **Free tier:** 3,000 emails/month (100/day)
- **Sign up:** https://app.mailersend.com/register
- **Get API key:** Settings → API Tokens → Generate new token
- **Env vars:**
  ```
  RELAY_MAILERSEND_1_KEY=mlsn.your_token_here
  RELAY_MAILERSEND_2_KEY=mlsn.your_second_token_here
  ```

### 6. Sendpulse
- **Free tier:** 15,000 emails/month (500/day)
- **Sign up:** https://sendpulse.com/en/signup
- **Get credentials:** Settings → API → Copy Client ID and Client Secret
- **Note:** Uses OAuth2 — needs Client ID + Client Secret (not a single API key)
- **Env vars:**
  ```
  RELAY_SENDPULSE_1_KEY=your_client_id
  RELAY_SENDPULSE_1_KEY2=your_client_secret
  RELAY_SENDPULSE_2_KEY=your_second_client_id
  RELAY_SENDPULSE_2_KEY2=your_second_client_secret
  ```

### 7. ZeptoMail (by Zoho)
- **Free tier:** 10,000 emails/month (333/day)
- **Sign up:** https://www.zoho.com/zeptomail/signup.html
- **Get API key:** ZeptoMail Dashboard → Mail Agents → your agent → Send Mail Token
- **Env vars:**
  ```
  RELAY_ZEPTOMAIL_1_KEY=Zoho-enczapikey your_token_here
  RELAY_ZEPTOMAIL_2_KEY=Zoho-enczapikey your_second_token_here
  ```

### 8. SMTP2GO
- **Free tier:** 1,000 emails/month (33/day)
- **Sign up:** https://www.smtp2go.com/pricing/
- **Get API key:** Settings → API Keys → Add API Key
- **Env vars:**
  ```
  RELAY_SMTP2GO_1_KEY=api-your_key_here
  RELAY_SMTP2GO_2_KEY=api-your_second_key_here
  ```

### 9. Elastic Email
- **Free tier:** 3,000 emails/month (100/day)
- **Sign up:** https://app.elasticemail.com/register
- **Get API key:** Settings → API → Create API Key (select "Full Access")
- **Env vars:**
  ```
  RELAY_ELASTICEMAIL_1_KEY=your_api_key_here
  RELAY_ELASTICEMAIL_2_KEY=your_second_api_key_here
  ```

### 10. Mailgun
- **Free tier:** 1,000 emails/month (Flex trial, 3 months)
- **Sign up:** https://signup.mailgun.com
- **Get API key:** Account → Security → API Keys → Private API key
- **Note:** Requires API key + sending domain (e.g. `mg.yourdomain.com`)
- **Env vars:**
  ```
  RELAY_MAILGUN_1_KEY=your-private-api-key
  RELAY_MAILGUN_1_KEY2=mg.yourdomain.com
  RELAY_MAILGUN_2_KEY=your-second-private-api-key
  RELAY_MAILGUN_2_KEY2=mg2.yourdomain.com
  ```

### 11. Postmark
- **Free tier:** 100 emails/month (developer test)
- **Sign up:** https://account.postmarkapp.com/sign_up
- **Get API key:** Server → API Tokens → Server API Token
- **Note:** Low free limit — useful as fallback, not primary
- **Env vars:**
  ```
  RELAY_POSTMARK_1_KEY=your-server-api-token
  RELAY_POSTMARK_2_KEY=your-second-server-api-token
  ```

### 12. SparkPost
- **Free tier:** 500 emails/month
- **Sign up:** https://app.sparkpost.com/join
- **Get API key:** Account → API Keys → Create API Key (Transmissions: Read/Write)
- **Env vars:**
  ```
  RELAY_SPARKPOST_1_KEY=your_api_key_here
  RELAY_SPARKPOST_2_KEY=your_second_api_key_here
  ```

### 13. SocketLabs
- **Free tier:** 40,000 emails first month (trial), ~1,000/month after
- **Sign up:** https://www.socketlabs.com/signup/
- **Get credentials:** Dashboard → SMTP Credentials → API Key + Server ID
- **Note:** Needs API Key + Server ID
- **Env vars:**
  ```
  RELAY_SOCKETLABS_1_KEY=your_api_key
  RELAY_SOCKETLABS_1_KEY2=your_server_id
  RELAY_SOCKETLABS_2_KEY=your_second_api_key
  RELAY_SOCKETLABS_2_KEY2=your_second_server_id
  ```

### 14. Netcore (formerly Pepipost)
- **Free tier:** 3,000 emails/month (100/day)
- **Sign up:** https://netcorecloud.com/email-api/
- **Get API key:** Dashboard → Settings → API Key
- **Env vars:**
  ```
  RELAY_NETCORE_1_KEY=your_api_key_here
  RELAY_NETCORE_2_KEY=your_second_api_key_here
  ```

### 15. Mailtrap
- **Free tier:** 1,000 emails/month (Email Sending product, not sandbox)
- **Sign up:** https://mailtrap.io/register/signup
- **Get API key:** Email Sending → API Keys → Generate
- **Note:** Use the "Email Sending" product, NOT the sandbox/inbox testing product
- **Env vars:**
  ```
  RELAY_MAILTRAP_1_KEY=your_api_token_here
  RELAY_MAILTRAP_2_KEY=your_second_api_token_here
  ```

---

## Quick Start: Maximum Free Capacity

To get the most free sends immediately:

1. **Sign up for Brevo (2 accounts)** — 18k/month, 300/day each. Fastest approval, no domain verification required on free tier.
2. **Sign up for Mailjet (2 accounts)** — 12k/month, instant activation.
3. **Sign up for Sendpulse (2 accounts)** — 30k/month, highest free tier.
4. **Sign up for SocketLabs (2 accounts)** — 80k in first month from trials.

These 4 providers alone give you **~140k emails in month 1**, all free.

---

## Adding API Keys in BillionMail

1. Go to **Settings → Relay Providers**
2. Find the provider group (e.g. "Brevo")
3. Enter your API key in "Account 1" slot
4. Click **Save Key** — slot activates automatically
5. Repeat for "Account 2" with your second account's key
6. Click **Test Pool Auto-Route** to verify routing works

BillionMail auto-routes through the highest-priority available slot. When Account 1 hits its daily limit, it switches to Account 2, then to the next provider — transparent to your users.
