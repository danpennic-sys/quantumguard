# QuantumGuard Stripe Test Slice

This directory is **not** the Operator Pack.

| | This harness | Operator Pack (the seat) |
|---|---|---|
| Amount | $10 | $490 once, or $1,490/year |
| Mode | Stripe test mode only | Live Payment Link, created in the Dashboard |
| Proves | Webhook signature + `payments.json` | A paid seat |
| Connects to `qg` | No | No. Fulfillment is a human, after payment |

List price and fulfillment rules: `../MARKETING.md` and `../docs/OPERATOR-PACK-CHECKOUT.md`. Do not commit `sk_live_`, `sk_test_`, or `whsec_` values.

Isolated Stripe Checkout test flow for QuantumGuard.

## Design rules

- Test mode only
- $10 one-time payment
- `/create-checkout-session` → Stripe-hosted Checkout
- `/webhook` with Stripe signature verification
- Idempotent recording into `payments.json`
- Success / cancel pages exist but **are not** treated as payment proof
- Webhook (`checkout.session.completed` + `payment_status === 'paid'`) is the sole confirmation boundary
- No card data is stored
- No connection to any QuantumGuard verifier

## Quick start

```bash
cp .env.example .env
# fill in sk_test_... and later the whsec_...

npm install
npm start
```

In a second terminal:

```bash
stripe listen --forward-to localhost:4242/webhook
```

Copy the printed `whsec_...` into `.env` and restart the server if needed.

Open http://localhost:4242 → Pay $10 → use test card `4242 4242 4242 4242`.

After the webhook fires you should see exactly one record in `payments.json`.

## Test path

```
QuantumGuard
     ↓
Create Checkout Session
     ↓
Stripe Test Checkout
     ↓
Test card
     ↓
Stripe webhook
     ↓
Signature verified
     ↓
checkout.session.completed
     ↓
payments.json
```
