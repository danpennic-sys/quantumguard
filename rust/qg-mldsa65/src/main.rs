// SPDX-License-Identifier: MIT OR Apache-2.0
// Copyright (c) 2026 danpennic-sys
//! CLI: verify ML-DSA-65 over QuantumGuard payload_hash message bytes.
//!
//! Usage:
//!   qg-mldsa65-verify --pk <hex> --sig <hex> --hash <hex>
//!   qg-mldsa65-verify --pk-file pk.bin --sig-file sig.bin --bundle bundle.json

use qg_mldsa65::{message_from_bundle_json, verify_mldsa65, MESSAGE_LEN};
use std::env;
use std::fs;
use std::process;

fn usage() -> ! {
    eprintln!(
        "usage:\n  qg-mldsa65-verify --pk <hex> --sig <hex> --hash <hex>\n  qg-mldsa65-verify --pk-file <path> --sig-file <path> --bundle <bundle.json>"
    );
    process::exit(2);
}

fn main() {
    let args: Vec<String> = env::args().skip(1).collect();
    if args.is_empty() {
        usage();
    }

    let mut pk: Option<Vec<u8>> = None;
    let mut sig: Option<Vec<u8>> = None;
    let mut msg: Option<Vec<u8>> = None;

    let mut i = 0;
    while i < args.len() {
        match args[i].as_str() {
            "--pk" => {
                i += 1;
                if i >= args.len() {
                    usage();
                }
                pk = Some(hex::decode(&args[i]).unwrap_or_else(|e| {
                    eprintln!("pk hex: {e}");
                    process::exit(2);
                }));
            }
            "--sig" => {
                i += 1;
                if i >= args.len() {
                    usage();
                }
                sig = Some(hex::decode(&args[i]).unwrap_or_else(|e| {
                    eprintln!("sig hex: {e}");
                    process::exit(2);
                }));
            }
            "--hash" => {
                i += 1;
                if i >= args.len() {
                    usage();
                }
                msg = Some(hex::decode(&args[i]).unwrap_or_else(|e| {
                    eprintln!("hash hex: {e}");
                    process::exit(2);
                }));
            }
            "--pk-file" => {
                i += 1;
                if i >= args.len() {
                    usage();
                }
                pk = Some(fs::read(&args[i]).unwrap_or_else(|e| {
                    eprintln!("pk-file: {e}");
                    process::exit(2);
                }));
            }
            "--sig-file" => {
                i += 1;
                if i >= args.len() {
                    usage();
                }
                sig = Some(fs::read(&args[i]).unwrap_or_else(|e| {
                    eprintln!("sig-file: {e}");
                    process::exit(2);
                }));
            }
            "--bundle" => {
                i += 1;
                if i >= args.len() {
                    usage();
                }
                let data = fs::read(&args[i]).unwrap_or_else(|e| {
                    eprintln!("bundle: {e}");
                    process::exit(2);
                });
                let m = message_from_bundle_json(&data).unwrap_or_else(|e| {
                    eprintln!("bundle message: {e}");
                    process::exit(2);
                });
                msg = Some(m.to_vec());
            }
            _ => usage(),
        }
        i += 1;
    }

    let (Some(pk), Some(sig), Some(msg)) = (pk, sig, msg) else {
        usage();
    };

    if msg.len() != MESSAGE_LEN {
        eprintln!("hash must be {MESSAGE_LEN} bytes");
        process::exit(2);
    }

    let result = verify_mldsa65(&pk, &msg, &sig, "ml-dsa-65");
    println!("{}", serde_json::to_string_pretty(&result).unwrap());

    match result.verdict.as_str() {
        "PASS" => process::exit(0),
        "FAIL" => process::exit(1),
        _ => process::exit(3),
    }
}
