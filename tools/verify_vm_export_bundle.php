<?php
// SPDX-License-Identifier: MIT OR Apache-2.0
// Copyright (c) 2026 danpennic-sys
declare(strict_types=1);

/**
 * VM-EXPORT-BUNDLE-001 conformance verifier (PHP reference).
 * Verify-only. Fail-closed. Emits JSON to stdout.
 *
 * Mirrors tools/verify_vm_export_bundle.py check order and codes.
 * Exit: 0 = PASS case, 1 = FAIL, 2 = harness error.
 */

function jcs_like_canonical(mixed $v): string {
    return json_encode($v, JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES);
}

function sort_keys_recursive(mixed $v): mixed {
    if (!is_array($v)) {
        return $v;
    }
    $isList = array_is_list($v);
    if (!$isList) {
        ksort($v);
    }
    foreach ($v as $k => $child) {
        $v[$k] = sort_keys_recursive($child);
    }
    return $v;
}

function compute_state_hash(mixed $state): string {
    $sorted = sort_keys_recursive($state);
    $canonical = jcs_like_canonical($sorted);
    return hash('sha256', $canonical);
}

function compute_replay_merkle_root(array $events): string {
    if (count($events) === 0) {
        return '';
    }
    $leaves = [];
    foreach ($events as $e) {
        $leaves[] = hex2bin($e['hash']);
    }
    while (count($leaves) > 1) {
        $next = [];
        for ($i = 0; $i < count($leaves); $i += 2) {
            if ($i + 1 === count($leaves)) {
                $next[] = $leaves[$i];
            } else {
                $next[] = hash('sha256', "\x01" . $leaves[$i] . $leaves[$i + 1], true);
            }
        }
        $leaves = $next;
    }
    return bin2hex($leaves[0]);
}

function verify_case(array $input): array {
    $hasState = !empty($input['state_hash']);
    $hasMerkle = !empty($input['replay_merkle_root']);

    if ($hasState && $hasMerkle) {
        return ['verdict' => 'FAIL', 'code' => 'BOTH_HASHES_CLAIMED'];
    }

    $files = $input['files'] ?? [];
    $seen = [];
    foreach ($files as $f) {
        $name = $f['name'] ?? '';
        if (str_contains($name, '/') || str_contains($name, '\\') ||
            str_starts_with($name, '..') || in_array($name, ['.', '..', ''], true)) {
            return [
                'verdict' => 'FAIL',
                'code' => 'ZIP_SLIP',
                'detail' => 'invalid entry name "' . $name . '" (subdirectories not allowed)',
            ];
        }
        if (isset($seen[$name])) {
            return [
                'verdict' => 'FAIL',
                'code' => 'DUPLICATE_PATH',
                'detail' => 'path appears more than once: ' . $name,
            ];
        }
        $seen[$name] = true;
    }

    if (isset($input['extracted_files']) && is_array($input['extracted_files'])) {
        $allowed = ['bundle.json' => true, 'signature.ed25519' => true, 'public_key.pem' => true];
        foreach ($files as $f) {
            $allowed[$f['name']] = true;
        }
        foreach ($input['extracted_files'] as $name) {
            if (!isset($allowed[$name])) {
                return [
                    'verdict' => 'FAIL',
                    'code' => 'EXTRA_FILE',
                    'detail' => 'unauthorized file present in bundle: ' . $name,
                ];
            }
        }
    }

    if (!empty($input['require_adr_anchor'])) {
        return [
            'verdict' => 'FAIL',
            'code' => 'NOT_IMPLEMENTED',
            'detail' => 'adr_anchor: NOT IMPLEMENTED in v1.0 — see VM-EXPORT-BUNDLE-001 §10.2',
        ];
    }

    if ($hasState) {
        if (!isset($input['state'])) {
            return ['verdict' => 'FAIL', 'code' => 'STATE_HASH_MISMATCH', 'detail' => 'missing state'];
        }
        $computed = compute_state_hash($input['state']);
        $claimed = $input['state_hash'];
        if ($computed !== $claimed) {
            $outerCa = $input['created_at'] ?? null;
            $stateCa = $input['state']['created_at'] ?? null;
            if ($outerCa && $stateCa && $outerCa !== $stateCa) {
                return [
                    'verdict' => 'FAIL',
                    'code' => 'CREATED_AT_TAMPER',
                    'detail' => 'created_at is inside the signed region; any modification invalidates root_hash/signature',
                ];
            }
            return ['verdict' => 'FAIL', 'code' => 'STATE_HASH_MISMATCH'];
        }
        return ['verdict' => 'PASS', 'code' => 'OK', 'computed_state_hash' => $computed];
    }

    if ($hasMerkle) {
        $events = $input['events'] ?? [];
        $computed = compute_replay_merkle_root($events);
        $claimed = $input['replay_merkle_root'];
        if ($computed !== $claimed) {
            return ['verdict' => 'FAIL', 'code' => 'REPLAY_MERKLE_MISMATCH'];
        }
        return ['verdict' => 'PASS', 'code' => 'OK', 'computed_replay_merkle_root' => $computed];
    }

    return [
        'verdict' => 'FAIL',
        'code' => 'NOT_IMPLEMENTED',
        'detail' => 'no hash claimed and no --no-* suppression',
    ];
}

function main(array $argv): int {
    if (count($argv) < 2) {
        fwrite(STDERR, "usage: php verify_vm_export_bundle.php <input.json> [--expected expected.json]\n");
        return 2;
    }
    $inputPath = $argv[1];
    $expectedPath = null;
    for ($i = 2; $i < count($argv); $i++) {
        if ($argv[$i] === '--expected' && isset($argv[$i + 1])) {
            $expectedPath = $argv[++$i];
        }
    }

    $raw = @file_get_contents($inputPath);
    if ($raw === false) {
        fwrite(STDERR, "read error: $inputPath\n");
        return 2;
    }
    $data = json_decode($raw, true);
    if (!is_array($data)) {
        fwrite(STDERR, "malformed JSON\n");
        return 2;
    }

    $result = verify_case($data);
    echo json_encode($result, JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE), "\n";

    if ($expectedPath !== null) {
        $exp = json_decode((string) file_get_contents($expectedPath), true);
        if (($result['verdict'] ?? null) !== ($exp['verdict'] ?? null) ||
            ($result['code'] ?? null) !== ($exp['code'] ?? null)) {
            fwrite(STDERR, "SELF-CHECK FAILED\n");
            return 2;
        }
    }

    return ($result['verdict'] ?? '') === 'PASS' ? 0 : 1;
}

exit(main($argv));
