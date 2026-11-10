//! Every vector in the NIST ACVP ML-DSA FIPS 204 projections.
//!
//! Files are the gzipped `internalProjection.json` from
//! usnistgov/ACVP-Server @ 975de31eb83d87039ec88934fdc47d8c312b892d,
//! vendored by fips204 0.5.0 (`66129c49273f3daa1338f9634acdf84ee8d1faf8`).
//! ML-DSA-65 pure and HashML-DSA verifies go through this crate.
//! ML-DSA-44, ML-DSA-87, keyGen, and the internal ACVP hooks go through the
//! same fips204 0.5.0 this crate links. A skipped case fails the test.

#![allow(deprecated)]

use std::fs::File;
use std::io::Read;

use fips204::ml_dsa_44;
use fips204::ml_dsa_65;
use fips204::ml_dsa_87;
use fips204::pre_hash;
use fips204::traits::{KeyGen, SerDes, Signer, Verifier};
use flate2::read::GzDecoder;
use hex::decode;
use rand_core::{CryptoRng, RngCore};
use serde_json::Value;
use sha2::{Digest, Sha224, Sha256, Sha384, Sha512, Sha512_224, Sha512_256};
use sha3::digest::{ExtendableOutput, Update, XofReader};
use sha3::{Sha3_224, Sha3_256, Sha3_384, Sha3_512, Shake128, Shake256};

struct TestRng {
    data: Vec<Vec<u8>>,
}

impl RngCore for TestRng {
    fn next_u32(&mut self) -> u32 {
        unimplemented!()
    }
    fn next_u64(&mut self) -> u64 {
        unimplemented!()
    }
    fn fill_bytes(&mut self, out: &mut [u8]) {
        let x = self.data.pop().expect("test rng problem");
        out.copy_from_slice(&x);
    }
    fn try_fill_bytes(&mut self, out: &mut [u8]) -> Result<(), rand_core::Error> {
        self.fill_bytes(out);
        Ok(())
    }
}
impl CryptoRng for TestRng {}
impl TestRng {
    fn new() -> Self {
        Self { data: Vec::new() }
    }
    fn push(&mut self, new_data: &[u8]) {
        self.data.push(new_data.to_vec());
    }
}

struct Tally {
    total: usize,
    ran: usize,
}

impl Tally {
    fn new(doc: &Value) -> Self {
        let total = groups(doc).iter().map(|g| cases(g).len()).sum();
        Self { total, ran: 0 }
    }
    fn finish(&self, which: &str) {
        println!("NIST {which}: ran {} of {}", self.ran, self.total);
        assert_eq!(self.ran, self.total, "NIST {which}: some cases did not run");
        assert!(self.ran > 0, "NIST {which}: no cases ran");
    }
}

#[derive(Clone, Copy)]
enum Kind {
    Pure,
    PreHash,
    Internal,
    ExternalMu,
}

fn kind_of(group: &Value) -> Kind {
    match (
        group["signatureInterface"].as_str(),
        group["preHash"].as_str(),
        group["externalMu"].as_bool(),
    ) {
        (Some("external"), Some("pure"), Some(false)) => Kind::Pure,
        (Some("external"), Some("preHash"), Some(false)) => Kind::PreHash,
        (Some("internal"), Some("none"), Some(false)) => Kind::Internal,
        (Some("internal"), Some("none"), Some(true)) => Kind::ExternalMu,
        other => panic!("tgId {}: unknown group kind {other:?}", group["tgId"]),
    }
}

fn load(name: &str) -> Value {
    let path = format!("tests/data/{name}/internalProjection.json.gz");
    let mut text = String::new();
    GzDecoder::new(File::open(&path).unwrap_or_else(|e| panic!("open {path}: {e}")))
        .read_to_string(&mut text)
        .unwrap_or_else(|e| panic!("gunzip {path}: {e}"));
    serde_json::from_str(&text).unwrap()
}

fn groups(doc: &Value) -> &Vec<Value> {
    doc["testGroups"].as_array().expect("testGroups")
}
fn cases(group: &Value) -> &Vec<Value> {
    group["tests"].as_array().expect("tests")
}

fn check_group(group: &Value, known: &[&str]) {
    for key in group.as_object().expect("group").keys() {
        assert!(
            known.contains(&key.as_str()),
            "tgId {}: unknown group field {key}",
            group["tgId"]
        );
    }
    assert_eq!(group["testType"], "AFT", "tgId {}: testType", group["tgId"]);
}

fn hex_field(v: &Value, key: &str, loc: &str) -> Vec<u8> {
    decode(
        v[key]
            .as_str()
            .unwrap_or_else(|| panic!("{loc} missing {key}")),
    )
    .unwrap()
}

fn get_ph(hashname: &str, message: &[u8]) -> (&'static [u8], Vec<u8>) {
    match hashname {
        "SHA2-224" => (&pre_hash::SHA2_224[..], Sha224::digest(message).to_vec()),
        "SHA2-256" => (&pre_hash::SHA2_256[..], Sha256::digest(message).to_vec()),
        "SHA2-384" => (&pre_hash::SHA2_384[..], Sha384::digest(message).to_vec()),
        "SHA2-512" => (&pre_hash::SHA2_512[..], Sha512::digest(message).to_vec()),
        "SHA2-512/224" => (
            &pre_hash::SHA2_512_224[..],
            Sha512_224::digest(message).to_vec(),
        ),
        "SHA2-512/256" => (
            &pre_hash::SHA2_512_256[..],
            Sha512_256::digest(message).to_vec(),
        ),
        "SHA3-224" => (&pre_hash::SHA3_224[..], Sha3_224::digest(message).to_vec()),
        "SHA3-256" => (&pre_hash::SHA3_256[..], Sha3_256::digest(message).to_vec()),
        "SHA3-384" => (&pre_hash::SHA3_384[..], Sha3_384::digest(message).to_vec()),
        "SHA3-512" => (&pre_hash::SHA3_512[..], Sha3_512::digest(message).to_vec()),
        "SHAKE-128" => {
            let mut ret = [0u8; 32];
            let mut hasher = Shake128::default();
            hasher.update(message);
            hasher.finalize_xof().read(&mut ret);
            (&pre_hash::SHAKE_128[..], ret.to_vec())
        }
        "SHAKE-256" => {
            let mut ret = [0u8; 64];
            let mut hasher = Shake256::default();
            hasher.update(message);
            hasher.finalize_xof().read(&mut ret);
            (&pre_hash::SHAKE_256[..], ret.to_vec())
        }
        other => panic!("unknown digest algorithm {other}"),
    }
}

fn crate_bool(result: qg_mldsa65::PqcResult, loc: &str) -> bool {
    match result.verdict.as_str() {
        "PASS" => {
            assert_eq!(result.code, "OK", "{loc}");
            assert!(!result.stub, "{loc}");
            true
        }
        "FAIL" => {
            assert_eq!(result.code, "E_PQC_SIG_FAIL", "{loc}");
            assert!(!result.stub, "{loc}");
            false
        }
        other => panic!("{loc}: crate verdict {other} ({})", result.reason),
    }
}

fn crate_pure(pk: &[u8], message: &[u8], sig: &[u8], ctx: &[u8], loc: &str) -> bool {
    crate_bool(qg_mldsa65::verify_external(pk, message, sig, ctx), loc)
}

fn crate_prehash(
    pk: &[u8],
    message: &[u8],
    sig: &[u8],
    ctx: &[u8],
    hash_alg: &str,
    loc: &str,
) -> bool {
    let (oid, digest) = get_ph(hash_alg, message);
    crate_bool(qg_mldsa65::verify_prehash(pk, &digest, sig, ctx, oid), loc)
}

macro_rules! sigver_fips {
    ($pc:ident, $test:ident, $loc:ident, $pk_bytes:ident, $signature:ident, $kind:ident) => {{
        let pk = $pc::PublicKey::try_from_bytes($pk_bytes.as_slice().try_into().unwrap()).unwrap();
        match $signature.as_slice().try_into() {
            Err(_) => false,
            Ok(sig) => match $kind {
                Kind::Pure => pk.verify(
                    &hex_field($test, "message", &$loc),
                    sig,
                    &hex_field($test, "context", &$loc),
                ),
                Kind::PreHash => {
                    let message = hex_field($test, "message", &$loc);
                    let (oid, digest) = get_ph($test["hashAlg"].as_str().unwrap(), &message);
                    pk.hash_verify(&digest, sig, &hex_field($test, "context", &$loc), oid)
                }
                Kind::Internal => pk.verify_internal(&hex_field($test, "message", &$loc), sig),
                Kind::ExternalMu => {
                    let mu: [u8; 64] = hex_field($test, "mu", &$loc).try_into().unwrap();
                    pk.verify_mu(&mu, sig)
                }
            },
        }
    }};
}

macro_rules! siggen_one {
    ($pc:ident, $via_crate:expr, $test:ident, $loc:ident, $seed_format:ident, $sk_bytes:ident, $pk_bytes:ident, $sig_exp:ident, $kind:ident, $rnd:ident, $tg_deterministic:ident) => {{
        let sk = if $seed_format {
            let xi: [u8; 32] = hex_field($test, "seed", &$loc).try_into().unwrap();
            let (_pk, sk) = $pc::KG::keygen_from_seed(&xi);
            assert_eq!(
                $sk_bytes,
                sk.clone().into_bytes(),
                "{} seed expands to sk",
                $loc
            );
            sk
        } else {
            let arr: [u8; $pc::SK_LEN] = $sk_bytes.as_slice().try_into().unwrap();
            $pc::PrivateKey::try_from_bytes(arr).unwrap()
        };
        let pk2 = sk.get_public_key();
        assert_eq!($pk_bytes, pk2.into_bytes(), "{}", $loc);
        let pk1 = $pc::PublicKey::try_from_bytes($pk_bytes.as_slice().try_into().unwrap()).unwrap();
        let sig: [u8; $pc::SIG_LEN] = $sig_exp.as_slice().try_into().unwrap();
        let (sig_act, verified) = match $kind {
            Kind::Pure => {
                let message = hex_field($test, "message", &$loc);
                let context = hex_field($test, "context", &$loc);
                let sig_act = sk.try_sign_with_rng(&mut $rnd, &message, &context).unwrap();
                let verified = if $via_crate {
                    assert!(
                        crate_pure(&$pk_bytes, &message, &$sig_exp, &context, &$loc),
                        "{} crate",
                        $loc
                    );
                    true
                } else {
                    pk1.verify(&message, &sig, &context)
                };
                (sig_act, verified)
            }
            Kind::PreHash => {
                let message = hex_field($test, "message", &$loc);
                let context = hex_field($test, "context", &$loc);
                let (oid, digest) = get_ph($test["hashAlg"].as_str().unwrap(), &message);
                let sig_act = sk
                    .try_hash_sign_with_rng(&mut $rnd, &digest, &context, oid)
                    .unwrap();
                let verified = if $via_crate {
                    assert!(
                        crate_prehash(
                            &$pk_bytes,
                            &message,
                            &$sig_exp,
                            &context,
                            $test["hashAlg"].as_str().unwrap(),
                            &$loc
                        ),
                        "{} crate",
                        $loc
                    );
                    true
                } else {
                    pk1.hash_verify(&digest, &sig, &context, oid)
                };
                (sig_act, verified)
            }
            Kind::Internal => {
                let message = hex_field($test, "message", &$loc);
                (
                    sk.sign_internal(&mut $rnd, &message, !$tg_deterministic)
                        .unwrap(),
                    pk1.verify_internal(&message, &sig),
                )
            }
            Kind::ExternalMu => {
                let mu: [u8; 64] = hex_field($test, "mu", &$loc).try_into().unwrap();
                (
                    sk.sign_mu(&mut $rnd, &mu, !$tg_deterministic).unwrap(),
                    pk1.verify_mu(&mu, &sig),
                )
            }
        };
        assert_eq!($sig_exp, sig_act, "{}", $loc);
        assert!(verified, "{} verify", $loc);
    }};
}

#[test]
fn nist_keygen_all() {
    let v = load("ML-DSA-keyGen-FIPS204");
    let mut tally = Tally::new(&v);
    for group in groups(&v) {
        check_group(group, &["parameterSet", "testType", "tests", "tgId"]);
        let set = group["parameterSet"].as_str().unwrap();
        for test in cases(group) {
            let loc = format!("{set} tgId={} tcId={}", group["tgId"], test["tcId"]);
            let seed: [u8; 32] = hex_field(test, "seed", &loc).try_into().unwrap();
            let pk_exp = hex_field(test, "pk", &loc);
            let sk_exp = hex_field(test, "sk", &loc);
            match set {
                "ML-DSA-44" => {
                    let (pk, sk) = ml_dsa_44::KG::keygen_from_seed(&seed);
                    assert_eq!(pk_exp, pk.into_bytes(), "{loc}");
                    assert_eq!(sk_exp, sk.into_bytes(), "{loc}");
                }
                "ML-DSA-65" => {
                    let (pk, sk) = ml_dsa_65::KG::keygen_from_seed(&seed);
                    assert_eq!(pk_exp, pk.into_bytes(), "{loc}");
                    assert_eq!(sk_exp, sk.into_bytes(), "{loc}");
                }
                "ML-DSA-87" => {
                    let (pk, sk) = ml_dsa_87::KG::keygen_from_seed(&seed);
                    assert_eq!(pk_exp, pk.into_bytes(), "{loc}");
                    assert_eq!(sk_exp, sk.into_bytes(), "{loc}");
                }
                other => panic!("{loc}: unknown parameter set {other}"),
            }
            tally.ran += 1;
        }
    }
    tally.finish("keyGen");
    assert_eq!(tally.ran, 75);
}

#[test]
fn nist_siggen_all() {
    let (n, via) = siggen("ML-DSA-sigGen-FIPS204", "sigGen");
    assert_eq!(n, 360);
    assert_eq!(via, 60);
}

#[test]
fn nist_siggen_tr1_all() {
    let (n, via) = siggen("ML-DSA-sigGen-FIPS204-tr1", "sigGen-tr1");
    assert_eq!(n, 720);
    assert_eq!(via, 120);
}

#[test]
fn nist_sigver_all() {
    let v = load("ML-DSA-sigVer-FIPS204");
    let mut tally = Tally::new(&v);
    let mut via_crate = 0usize;
    for group in groups(&v) {
        check_group(
            group,
            &[
                "externalMu",
                "parameterSet",
                "preHash",
                "signatureInterface",
                "testType",
                "tests",
                "tgId",
            ],
        );
        let set = group["parameterSet"].as_str().unwrap();
        let kind = kind_of(group);
        for test in cases(group) {
            let loc = format!("{set} tgId={} tcId={}", group["tgId"], test["tcId"]);
            let signature = hex_field(test, "signature", &loc);
            let pk_bytes = hex_field(test, "pk", &loc);
            let expect = test["testPassed"].as_bool().unwrap();
            let got = match set {
                "ML-DSA-44" => sigver_fips!(ml_dsa_44, test, loc, pk_bytes, signature, kind),
                "ML-DSA-65" => {
                    let got = sigver_65(test, kind, &pk_bytes, &signature, &loc);
                    via_crate += 1;
                    got
                }
                "ML-DSA-87" => sigver_fips!(ml_dsa_87, test, loc, pk_bytes, signature, kind),
                other => panic!("{loc}: unknown parameter set {other}"),
            };
            assert_eq!(got, expect, "{loc} ({})", test["reason"]);
            tally.ran += 1;
        }
    }
    tally.finish("sigVer");
    assert_eq!(tally.ran, 180);
    assert_eq!(via_crate, 60);
}

fn sigver_65(test: &Value, kind: Kind, pk_bytes: &[u8], signature: &[u8], loc: &str) -> bool {
    match kind {
        Kind::Pure => {
            let message = hex_field(test, "message", loc);
            let context = hex_field(test, "context", loc);
            crate_pure(pk_bytes, &message, signature, &context, loc)
        }
        Kind::PreHash => {
            let message = hex_field(test, "message", loc);
            let context = hex_field(test, "context", loc);
            crate_prehash(
                pk_bytes,
                &message,
                signature,
                &context,
                test["hashAlg"].as_str().unwrap(),
                loc,
            )
        }
        Kind::Internal => {
            let message = hex_field(test, "message", loc);
            let pk = ml_dsa_65::PublicKey::try_from_bytes(pk_bytes.try_into().unwrap()).unwrap();
            let sig: [u8; ml_dsa_65::SIG_LEN] = match signature.try_into() {
                Ok(s) => s,
                Err(_) => return false,
            };
            pk.verify_internal(&message, &sig)
        }
        Kind::ExternalMu => {
            let mu: [u8; 64] = hex_field(test, "mu", loc).try_into().unwrap();
            let pk = ml_dsa_65::PublicKey::try_from_bytes(pk_bytes.try_into().unwrap()).unwrap();
            let sig: [u8; ml_dsa_65::SIG_LEN] = match signature.try_into() {
                Ok(s) => s,
                Err(_) => return false,
            };
            pk.verify_mu(&mu, &sig)
        }
    }
}

fn siggen(name: &str, which: &str) -> (usize, usize) {
    let v = load(name);
    let mut tally = Tally::new(&v);
    let mut via_crate = 0usize;
    for group in groups(&v) {
        check_group(
            group,
            &[
                "cornerCase",
                "deterministic",
                "externalMu",
                "keyFormat",
                "parameterSet",
                "preHash",
                "signatureInterface",
                "testType",
                "tests",
                "tgId",
            ],
        );
        let set = group["parameterSet"].as_str().unwrap();
        let kind = kind_of(group);
        let tg_deterministic = group["deterministic"].as_bool().unwrap();
        let seed_format = match group.get("keyFormat").map(Value::as_str) {
            None | Some(Some("expanded")) => false,
            Some(Some("seed")) => true,
            other => panic!("tgId {}: unknown keyFormat {other:?}", group["tgId"]),
        };
        for test in cases(group) {
            let loc = format!("{set} tgId={} tcId={}", group["tgId"], test["tcId"]);
            let sk_bytes = hex_field(test, "sk", &loc);
            let pk_bytes = hex_field(test, "pk", &loc);
            let sig_exp = hex_field(test, "signature", &loc);
            let rnd_field = test["rnd"].as_str();
            assert_eq!(rnd_field.is_none(), tg_deterministic, "{loc} rnd");
            let rnd_seed: [u8; 32] = match rnd_field {
                None => [0u8; 32],
                Some(s) => decode(s).unwrap().try_into().unwrap(),
            };
            let mut rnd = TestRng::new();
            rnd.push(&rnd_seed);
            let used = match set {
                "ML-DSA-44" => {
                    siggen_one!(
                        ml_dsa_44,
                        false,
                        test,
                        loc,
                        seed_format,
                        sk_bytes,
                        pk_bytes,
                        sig_exp,
                        kind,
                        rnd,
                        tg_deterministic
                    );
                    false
                }
                "ML-DSA-65" => {
                    siggen_one!(
                        ml_dsa_65,
                        true,
                        test,
                        loc,
                        seed_format,
                        sk_bytes,
                        pk_bytes,
                        sig_exp,
                        kind,
                        rnd,
                        tg_deterministic
                    );
                    matches!(kind, Kind::Pure | Kind::PreHash)
                }
                "ML-DSA-87" => {
                    siggen_one!(
                        ml_dsa_87,
                        false,
                        test,
                        loc,
                        seed_format,
                        sk_bytes,
                        pk_bytes,
                        sig_exp,
                        kind,
                        rnd,
                        tg_deterministic
                    );
                    false
                }
                other => panic!("{loc}: unknown parameter set {other}"),
            };
            if used {
                via_crate += 1;
            }
            tally.ran += 1;
        }
    }
    tally.finish(which);
    println!("NIST {which}: crate verifies {via_crate}");
    (tally.ran, via_crate)
}
