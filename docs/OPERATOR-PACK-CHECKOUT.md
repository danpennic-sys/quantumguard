# Operator Pack checkout

The seat. Not a verifier change. Do this in the Stripe Dashboard. Do not add a secret to this repo.

`v1.0.0` is already the frozen drop: https://github.com/danpennic-sys/quantumguard/releases/tag/v1.0.0
Those binaries are unsigned. `SHA256SUMS.txt` on that release is the integrity record. This runbook does not cut a new tag.

## What the buyer is paying for

- Operator Pages QuantumGuard, Operator Pack, single site
- $490 once, or $1,490 per year (separate link)
- Not post-quantum. Not QuantumGuard, Inc.

`stripe-test/` stays a $10 test-mode harness. It is not this SKU.

## Steps

1. Stripe Dashboard → Payment Links → New. One link per price.
2. Product name: `Operator Pages QuantumGuard — Operator Pack (single site)`.
3. Price: 49000 cents USD, one-time. Second link: 149000 cents USD, yearly.
4. Collect the buyer email. Do not collect a GitHub username on the card form unless you are ready to read it the same day.
5. Run the link in test mode first. Card `4242 4242 4242 4242`. Confirm Stripe shows paid. Success-page HTML is not proof. The Dashboard payment is proof.
6. Switch that link to live mode. Put the live URL in the README purchase line, replacing the mailto-only path, in a follow-up commit the hour it works. Until that commit, this file's status is NOT LIVE.
7. Fulfill by hand:
   - Add the buyer as a collaborator on the private `quantumguard-operator` repo, or
   - Send the `v1.0.0` release assets plus `SHA256SUMS.txt` and `CONTRACT-v1.md`.
8. Do not grant access from `stripe-test` or from a success page. A webhook you have not deployed must not mint a seat.

## What not to do

- Do not move or delete tag `v1.0.0`.
- Do not put `sk_live_`, `sk_test_`, or `whsec_` in git.
- Do not describe the Rust `qg-mldsa65` crate as included in the pack.
