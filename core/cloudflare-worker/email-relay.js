/**
 * BillionMail Cloudflare Email Relay Worker
 *
 * Deploy this to Cloudflare Workers (free tier: 100k requests/day).
 * It receives email payloads from BillionMail via HTTPS POST and
 * delivers them through Cloudflare's trusted IP infrastructure.
 *
 * NO port 25 needed on your server. Cloudflare handles delivery.
 *
 * Setup:
 *   1. Go to dash.cloudflare.com → Workers & Pages → Create Worker
 *   2. Paste this script
 *   3. Add Environment Variable: BILLIONMAIL_SECRET = any random string
 *   4. In Email Routing (your domain) → enable Email Routing
 *   5. Copy your Worker URL → paste into BillionMail Delivery Backend page
 */

export default {
  async fetch(request, env) {
    // Only accept POST requests
    if (request.method !== 'POST') {
      return new Response('Method not allowed', { status: 405 })
    }

    // Verify shared secret from BillionMail
    const authHeader = request.headers.get('X-BillionMail-Secret')
    if (!env.BILLIONMAIL_SECRET || authHeader !== env.BILLIONMAIL_SECRET) {
      return new Response('Unauthorized', { status: 401 })
    }

    let payload
    try {
      payload = await request.json()
    } catch {
      return new Response('Invalid JSON body', { status: 400 })
    }

    const { to, from, from_name, subject, html, text, reply_to, message_id } = payload

    if (!to || !from || !subject) {
      return new Response(JSON.stringify({ ok: false, error: 'Missing required fields: to, from, subject' }), {
        status: 400,
        headers: { 'Content-Type': 'application/json' }
      })
    }

    // Build RFC 2822 email message
    const boundary = `bm_${Date.now()}_${Math.random().toString(36).slice(2)}`
    const msgId = message_id || `<${Date.now()}.${Math.random().toString(36).slice(2)}@billionmail.worker>`

    const headers = [
      `From: ${from_name ? `"${from_name}" <${from}>` : from}`,
      `To: ${Array.isArray(to) ? to.join(', ') : to}`,
      `Subject: ${subject}`,
      `Message-ID: ${msgId}`,
      `Date: ${new Date().toUTCString()}`,
      `MIME-Version: 1.0`,
      `X-Mailer: BillionMail-Cloudflare-Worker/1.0`,
    ]

    if (reply_to) headers.push(`Reply-To: ${reply_to}`)

    let body
    if (html && text) {
      headers.push(`Content-Type: multipart/alternative; boundary="${boundary}"`)
      body = [
        `--${boundary}`,
        `Content-Type: text/plain; charset=UTF-8`,
        `Content-Transfer-Encoding: quoted-printable`,
        '',
        text,
        `--${boundary}`,
        `Content-Type: text/html; charset=UTF-8`,
        `Content-Transfer-Encoding: quoted-printable`,
        '',
        html,
        `--${boundary}--`,
      ].join('\r\n')
    } else if (html) {
      headers.push(`Content-Type: text/html; charset=UTF-8`)
      headers.push(`Content-Transfer-Encoding: quoted-printable`)
      body = html
    } else {
      headers.push(`Content-Type: text/plain; charset=UTF-8`)
      body = text || ''
    }

    const rawMessage = [...headers, '', body].join('\r\n')

    // Use Cloudflare Email Routing send_email binding
    // Requires: add email binding in Worker settings → "Send Email"
    try {
      if (env.EMAIL) {
        // Cloudflare send_email binding (preferred)
        const toAddresses = Array.isArray(to) ? to : [to]
        for (const recipient of toAddresses) {
          const message = new EmailMessage(from, recipient, rawMessage)
          await env.EMAIL.send(message)
        }
        return new Response(JSON.stringify({ ok: true, message_id: msgId, method: 'cf_email_binding' }), {
          headers: { 'Content-Type': 'application/json' }
        })
      } else {
        // Fallback: Cloudflare Email Routing API
        // Workers can also use fetch() to call the Cloudflare API
        return new Response(JSON.stringify({
          ok: false,
          error: 'EMAIL binding not configured. Add a "Send Email" binding to this Worker in the Cloudflare dashboard.',
          setup_url: 'https://developers.cloudflare.com/email-routing/email-workers/send-email-workers/'
        }), {
          status: 503,
          headers: { 'Content-Type': 'application/json' }
        })
      }
    } catch (err) {
      return new Response(JSON.stringify({ ok: false, error: String(err) }), {
        status: 500,
        headers: { 'Content-Type': 'application/json' }
      })
    }
  }
}
