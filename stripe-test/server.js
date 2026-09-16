require('dotenv').config();
const express = require('express');
const stripe = require('stripe')(process.env.STRIPE_SECRET_KEY);
const fs = require('fs');
const path = require('path');

const app = express();
const PORT = process.env.PORT || 4242;
const PAYMENTS_FILE = path.join(__dirname, 'payments.json');

// Ensure payments.json exists
if (!fs.existsSync(PAYMENTS_FILE)) {
  fs.writeFileSync(PAYMENTS_FILE, JSON.stringify([], null, 2));
}

function loadPayments() {
  return JSON.parse(fs.readFileSync(PAYMENTS_FILE, 'utf8'));
}

function savePayment(record) {
  const payments = loadPayments();
  // Idempotency: skip if session already recorded
  if (payments.some(p => p.sessionId === record.sessionId)) {
    console.log(`Idempotent skip: ${record.sessionId}`);
    return;
  }
  payments.push(record);
  fs.writeFileSync(PAYMENTS_FILE, JSON.stringify(payments, null, 2));
  console.log(`Recorded payment: ${record.sessionId}`);
}

// Stripe requires raw body for signature verification
app.post('/webhook', express.raw({ type: 'application/json' }), (req, res) => {
  const sig = req.headers['stripe-signature'];
  let event;

  try {
    event = stripe.webhooks.constructEvent(
      req.body,
      sig,
      process.env.STRIPE_WEBHOOK_SECRET
    );
  } catch (err) {
    console.error('Webhook signature verification failed:', err.message);
    return res.status(400).send(`Webhook Error: ${err.message}`);
  }

  if (event.type === 'checkout.session.completed') {
    const session = event.data.object;

    // Only record successful paid sessions
    if (session.payment_status === 'paid') {
      savePayment({
        sessionId: session.id,
        paymentIntentId: session.payment_intent,
        amount: session.amount_total,
        currency: session.currency,
        customerEmail: session.customer_details?.email || null,
        createdAt: new Date().toISOString(),
        livemode: session.livemode
      });
    }
  }

  res.json({ received: true });
});

// JSON body for other routes
app.use(express.json());
app.use(express.static('public'));

app.post('/create-checkout-session', async (req, res) => {
  try {
    const session = await stripe.checkout.sessions.create({
      mode: 'payment',
      payment_method_types: ['card'],
      line_items: [
        {
          price_data: {
            currency: 'usd',
            product_data: {
              name: 'QuantumGuard Test Payment',
              description: 'Isolated $10 test payment'
            },
            unit_amount: 1000 // $10.00
          },
          quantity: 1
        }
      ],
      success_url: `${req.protocol}://${req.get('host')}/success.html?session_id={CHECKOUT_SESSION_ID}`,
      cancel_url: `${req.protocol}://${req.get('host')}/cancel.html`
    });

    res.redirect(303, session.url);
  } catch (err) {
    console.error(err);
    res.status(500).json({ error: err.message });
  }
});

app.get('/', (req, res) => {
  res.send(`
    <h1>QuantumGuard Stripe Test Slice</h1>
    <form action="/create-checkout-session" method="POST">
      <button type="submit">Pay $10 (Test Mode)</button>
    </form>
    <p>Webhook is the only confirmation boundary. Success page is not payment proof.</p>
  `);
});

app.listen(PORT, () => {
  console.log(`QuantumGuard Stripe test slice running on http://localhost:${PORT}`);
  console.log('Use Stripe CLI for local webhooks:');
  console.log('  stripe listen --forward-to localhost:4242/webhook');
});
