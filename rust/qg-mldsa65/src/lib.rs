// Copyright (c) 2026 danpennic-sys. All rights reserved.
// Proprietary. This crate is not granted under MIT or Apache-2.0.
// BinderSeal remains under the MIT License and is not part of this crate.
//! QuantumGuard ML-DSA-65 (FIPS 204) real verify path.
//! Contract: docs/QUANTUMGUARD-V2-PQC-CONTRACT.md
//! `verify_mldsa65` message = 32-byte payload_hash, alg = ml-dsa-65, empty context.
//! `verify_external` and `verify_prehash` are the FIPS 204 verify paths used by the
//! NIST ACVP suite (arbitrary message, context, and HashML-DSA). stub is always false.

use fips204::ml_dsa_65::{self, PK_LEN, SIG_LEN};
use fips204::traits::{SerDes, Verifier};
use serde::{Deserialize, Serialize};

pub const PUBLIC_KEY_LEN: usize = PK_LEN;
pub const SIGNATURE_LEN: usize = SIG_LEN;
pub const MESSAGE_LEN: usize = 32;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PqcResult {
    pub verdict: String,
    pub code: String,
    pub reason: String,
    pub alg: String,
    pub stub: bool,
}

/// Verify ML-DSA-65 over QuantumGuard message bytes (raw payload_hash).
pub fn verify_mldsa65(public_key: &[u8], message: &[u8], signature: &[u8], alg: &str) -> PqcResult {
    if alg != "ml-dsa-65" {
        return PqcResult {
            verdict: "INDETERMINATE".into(),
            code: "E_PQC_UNSUPPORTED_ALG".into(),
            reason: format!("unsupported alg: {alg:?}"),
            alg: alg.into(),
            stub: false,
        };
    }
    if message.len() != MESSAGE_LEN {
        return PqcResult {
            verdict: "INDETERMINATE".into(),
            code: "E_PQC_CANONICAL".into(),
            reason: format!(
                "payload_hash must be {MESSAGE_LEN} bytes, got {}",
                message.len()
            ),
            alg: "ml-dsa-65".into(),
            stub: false,
        };
    }
    // QuantumGuard evidence: pure ML-DSA over the 32-byte payload_hash, empty context.
    verify_external(public_key, message, signature, &[])
}

/// FIPS 204 `ML-DSA.Verify` for ML-DSA-65. Message and context are the NIST inputs,
/// not the QuantumGuard 32-byte payload hash. Context longer than 255 bytes is a
/// failed verify, matching the standard.
pub fn verify_external(
    public_key: &[u8],
    message: &[u8],
    signature: &[u8],
    context: &[u8],
) -> PqcResult {
    match decode_key_sig(public_key, signature) {
        Ok((pk, sig)) => verify_decoded(&pk, message, &sig, context),
        Err(result) => result,
    }
}

/// FIPS 204 `HashML-DSA.Verify` for ML-DSA-65. `digest` is PH(message), `hash_oid`
/// is the DER OID of PH. Used by the NIST pre-hash groups.
pub fn verify_prehash(
    public_key: &[u8],
    digest: &[u8],
    signature: &[u8],
    context: &[u8],
    hash_oid: &[u8],
) -> PqcResult {
    let (pk, sig) = match decode_key_sig(public_key, signature) {
        Ok(v) => v,
        Err(result) => return result,
    };
    if pk.hash_verify(digest, &sig, context, hash_oid) {
        pass()
    } else {
        fail()
    }
}

fn decode_key_sig(
    public_key: &[u8],
    signature: &[u8],
) -> Result<(ml_dsa_65::PublicKey, [u8; SIG_LEN]), PqcResult> {
    if public_key.len() != PUBLIC_KEY_LEN {
        return Err(malformed(format!(
            "public key must be {PUBLIC_KEY_LEN} bytes, got {}",
            public_key.len()
        )));
    }
    if signature.len() != SIGNATURE_LEN {
        return Err(malformed(format!(
            "signature must be {SIGNATURE_LEN} bytes, got {}",
            signature.len()
        )));
    }
    let pk_arr: [u8; PK_LEN] = match public_key.try_into() {
        Ok(a) => a,
        Err(_) => return Err(malformed("public key decode failed")),
    };
    let pk = match ml_dsa_65::PublicKey::try_from_bytes(pk_arr) {
        Ok(pk) => pk,
        Err(_) => return Err(malformed("public key decode failed")),
    };
    let sig_arr: [u8; SIG_LEN] = match signature.try_into() {
        Ok(a) => a,
        Err(_) => return Err(malformed("signature decode failed")),
    };
    Ok((pk, sig_arr))
}

fn verify_decoded(
    pk: &ml_dsa_65::PublicKey,
    message: &[u8],
    signature: &[u8; SIG_LEN],
    context: &[u8],
) -> PqcResult {
    if pk.verify(message, signature, context) {
        pass()
    } else {
        fail()
    }
}

fn pass() -> PqcResult {
    PqcResult {
        verdict: "PASS".into(),
        code: "OK".into(),
        reason: "valid".into(),
        alg: "ml-dsa-65".into(),
        stub: false,
    }
}

fn fail() -> PqcResult {
    PqcResult {
        verdict: "FAIL".into(),
        code: "E_PQC_SIG_FAIL".into(),
        reason: "ML-DSA-65 signature verification failed".into(),
        alg: "ml-dsa-65".into(),
        stub: false,
    }
}

fn malformed(reason: impl Into<String>) -> PqcResult {
    PqcResult {
        verdict: "INDETERMINATE".into(),
        code: "E_PQC_MALFORMED".into(),
        reason: reason.into(),
        alg: "ml-dsa-65".into(),
        stub: false,
    }
}

/// Message = hex-decoded payload_hash from v1 evidence JSON.
pub fn message_from_bundle_json(data: &[u8]) -> Result<[u8; MESSAGE_LEN], String> {
    #[derive(Deserialize)]
    struct Bundle {
        payload_hash: String,
    }
    let b: Bundle = serde_json::from_slice(data).map_err(|e| e.to_string())?;
    let bytes = hex::decode(b.payload_hash.trim()).map_err(|e| e.to_string())?;
    if bytes.len() != MESSAGE_LEN {
        return Err(format!(
            "payload_hash must be {MESSAGE_LEN} bytes, got {}",
            bytes.len()
        ));
    }
    let mut out = [0u8; MESSAGE_LEN];
    out.copy_from_slice(&bytes);
    Ok(out)
}

#[cfg(test)]
mod tests {
    use super::*;
    use fips204::traits::Signer;

    #[test]
    fn roundtrip_sign_verify() {
        let (pk, sk) = ml_dsa_65::try_keygen().expect("keygen");
        let msg = [0x11u8; 32];
        let sig = sk.try_sign(&msg, &[]).expect("sign");
        let r = verify_mldsa65(&pk.into_bytes(), &msg, &sig, "ml-dsa-65");
        assert_eq!(r.verdict, "PASS");
        assert_eq!(r.code, "OK");
        assert!(!r.stub);
    }

    #[test]
    fn reject_bad_sig() {
        let (pk, sk) = ml_dsa_65::try_keygen().expect("keygen");
        let msg = [0x22u8; 32];
        let mut sig = sk.try_sign(&msg, &[]).expect("sign");
        sig[0] ^= 0xff;
        let r = verify_mldsa65(&pk.into_bytes(), &msg, &sig, "ml-dsa-65");
        assert_eq!(r.verdict, "FAIL");
        assert_eq!(r.code, "E_PQC_SIG_FAIL");
    }

    #[test]
    fn unsupported_alg() {
        let r = verify_mldsa65(&[0u8; PK_LEN], &[0u8; 32], &[0u8; SIG_LEN], "dilithium3");
        assert_eq!(r.code, "E_PQC_UNSUPPORTED_ALG");
    }

    #[test]
    fn bad_message_len() {
        let r = verify_mldsa65(&[0u8; PK_LEN], &[0u8; 31], &[0u8; SIG_LEN], "ml-dsa-65");
        assert_eq!(r.code, "E_PQC_CANONICAL");
    }

    #[test]
    fn external_accepts_non_canonical_message() {
        let (pk, sk) = ml_dsa_65::try_keygen().expect("keygen");
        let msg = [0x5au8; 100];
        let ctx = b"ctx";
        let sig = sk.try_sign(&msg, ctx).expect("sign");
        let pk_bytes = pk.into_bytes();
        let got = verify_external(&pk_bytes, &msg, &sig, ctx);
        assert_eq!(got.verdict, "PASS");
        assert!(!got.stub);
        let gated = verify_mldsa65(&pk_bytes, &msg, &sig, "ml-dsa-65");
        assert_eq!(gated.code, "E_PQC_CANONICAL");
    }
}
