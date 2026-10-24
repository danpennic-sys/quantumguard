// SPDX-License-Identifier: MIT OR Apache-2.0
// Copyright (c) 2026 danpennic-sys
//! QuantumGuard ML-DSA-65 (FIPS 204) real verify path.
//! Contract: docs/QUANTUMGUARD-V2-PQC-CONTRACT.md
//! Message = 32-byte payload_hash. Alg = ml-dsa-65. Empty context.
//! stub is always false — real verify or form reject.

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
pub fn verify_mldsa65(
    public_key: &[u8],
    message: &[u8],
    signature: &[u8],
    alg: &str,
) -> PqcResult {
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
            reason: format!("payload_hash must be {MESSAGE_LEN} bytes, got {}", message.len()),
            alg: "ml-dsa-65".into(),
            stub: false,
        };
    }
    if public_key.len() != PUBLIC_KEY_LEN {
        return PqcResult {
            verdict: "INDETERMINATE".into(),
            code: "E_PQC_MALFORMED".into(),
            reason: format!("public key must be {PUBLIC_KEY_LEN} bytes, got {}", public_key.len()),
            alg: "ml-dsa-65".into(),
            stub: false,
        };
    }
    if signature.len() != SIGNATURE_LEN {
        return PqcResult {
            verdict: "INDETERMINATE".into(),
            code: "E_PQC_MALFORMED".into(),
            reason: format!("signature must be {SIGNATURE_LEN} bytes, got {}", signature.len()),
            alg: "ml-dsa-65".into(),
            stub: false,
        };
    }

    let pk_arr: [u8; PK_LEN] = match public_key.try_into() {
        Ok(a) => a,
        Err(_) => {
            return PqcResult {
                verdict: "INDETERMINATE".into(),
                code: "E_PQC_MALFORMED".into(),
                reason: "public key decode failed".into(),
                alg: "ml-dsa-65".into(),
                stub: false,
            };
        }
    };
    let pk = match ml_dsa_65::PublicKey::try_from_bytes(&pk_arr) {
        Ok(pk) => pk,
        Err(_) => {
            return PqcResult {
                verdict: "INDETERMINATE".into(),
                code: "E_PQC_MALFORMED".into(),
                reason: "public key decode failed".into(),
                alg: "ml-dsa-65".into(),
                stub: false,
            };
        }
    };
    let sig_arr: [u8; SIG_LEN] = match signature.try_into() {
        Ok(a) => a,
        Err(_) => {
            return PqcResult {
                verdict: "INDETERMINATE".into(),
                code: "E_PQC_MALFORMED".into(),
                reason: "signature decode failed".into(),
                alg: "ml-dsa-65".into(),
                stub: false,
            };
        }
    };

    // Empty context: Pure ML-DSA (CONTRACT-PQC-1).
    if pk.verify(message, &sig_arr, &[]) {
        PqcResult {
            verdict: "PASS".into(),
            code: "OK".into(),
            reason: "valid".into(),
            alg: "ml-dsa-65".into(),
            stub: false,
        }
    } else {
        PqcResult {
            verdict: "FAIL".into(),
            code: "E_PQC_SIG_FAIL".into(),
            reason: "ML-DSA-65 signature verification failed".into(),
            alg: "ml-dsa-65".into(),
            stub: false,
        }
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
        return Err(format!("payload_hash must be {MESSAGE_LEN} bytes, got {}", bytes.len()));
    }
    let mut out = [0u8; MESSAGE_LEN];
    out.copy_from_slice(&bytes);
    Ok(out)
}

#[cfg(test)]
mod tests {
    use super::*;
    use fips204::traits::{KeyGen, Signer};

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
}
