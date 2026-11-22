# Operator Pack checkout

The seat. Not a verifier change. Do this in the Stripe Dashboard. Do not add a secret to this repo.

`v1.0.0` is already the frozen drop: https://github.com/danpennic-sys/quantumguard/releases/tag/v1.0.0
Those binaries are unsigned. `SHA256SUMS.txt` on that release is the integrity record. This runbook does not cut a new tag.

## Status

**LIVE** as of 2026-10-28.

- One-time ($490): https://buy.stripe.com/5kQaEYc9AewBaJleoN0ZW01
- Yearly ($1,495): https://buy.stripe.com/cNibJ24H8gEJcRteoN0ZW02

2026-11-18 re-read: the yearly link shows $1,495.00 due today, unit_amount 149500, interval year. The earlier 149000 figure did not match the charge. Do not open a replacement price at 149000.

2026-11-22 re-read, no charge sent. Both links are still livemode and the amounts are unchanged.

- One-time link: mode `payment`, price `price_1UO8XSKCM4PgZoVF1frIZxdg`, 49000 cents USD, one-time. The product name on that price is `Operator Pages`, not the runbook name below.
- Yearly link: mode `subscription`, price `price_1UVgKaKCM4PgZoVFD5lGs1ft`, 149500 cents USD, interval `year`. The product name on that price is `Quantomguard yearly`, spelled that way.

The public landing checkout is a third live link, `https://buy.stripe.com/dRmaEYb5w2NT2cP2G50ZW03`. It is mode `subscription` and includes both prices above. It is not a substitute for either single link.

## What the buyer is paying for

- Operator Pages QuantumGuard, Operator Pack, single site
- $490 once or $1,495 per year
- Not post-quantum. Not QuantumGuard, Inc.

`stripe-test/` stays a $10 test-mode harness. It is not this SKU.

## Steps (both links live)

1. Stripe Dashboard → Payment Links → New. One link per price.
2. Product name: `Operator Pages QuantumGuard — Operator Pack (single site)`.
3. Price: 49000 cents USD, one-time. Second link: 149500 cents USD, yearly. The live yearly link already charges 149500. Do not recreate it at 149000.
4. Collect the buyer email. Do not collect a GitHub username on the card form unless you are ready to read it the same day.
5. Run the link in test mode first. Card `4242 4242 4242 4242`. Confirm Stripe shows paid. Success-page HTML is not proof. The Dashboard payment is proof.
6. Switch that link to live mode. Put the live URL in the README purchase line in the same hour it works. **Done for both.**
7. Fulfill by hand:
   - Add the buyer as a collaborator on the private `quantumguard-operator` repo, or
   - Send the `v1.0.0` release assets plus `SHA256SUMS.txt` and `CONTRACT-v1.md`.
8. Do not grant access from `stripe-test` or from a success page. A webhook you have not deployed must not mint a seat.

## What not to do

- Do not move or delete tag `v1.0.0`.
- Do not put `sk_live_`, `sk_test_`, or `whsec_` in git.
- Do not describe the Rust `qg-mldsa65` crate as included in the pack.
