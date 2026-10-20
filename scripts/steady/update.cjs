'use strict';

const { execFileSync } = require('child_process');
const fs = require('fs');
const os = require('os');
const path = require('path');

const [revision] = process.argv.slice(2);
if (!/^[a-f0-9]{40}$/u.test(revision ?? '')) {
  throw new Error('Usage: node scripts/steady/update.cjs <full-commit-sha>');
}

const manifestPath = path.join(__dirname, 'manifest.json');
const manifest = JSON.parse(fs.readFileSync(manifestPath, 'utf-8'));
const temporary = fs.mkdtempSync(path.join(os.tmpdir(), 'steady-update-'));
try {
  execFileSync('git', ['init', '--quiet', temporary], { stdio: 'inherit' });
  execFileSync(
    'git',
    ['-C', temporary, 'fetch', '--quiet', '--no-tags', '--', manifest.steady.repository, revision],
    { stdio: 'inherit' },
  );
  execFileSync(
    'git',
    [
      '-C',
      temporary,
      '-c',
      'advice.detachedHead=false',
      '-c',
      'core.autocrlf=false',
      '-c',
      'core.eol=lf',
      'checkout',
      '--quiet',
      '--detach',
      revision,
    ],
    { stdio: 'inherit' },
  );
  const digest = execFileSync(process.execPath, [path.join(__dirname, 'source-sha256.cjs'), temporary], {
    encoding: 'utf-8',
  }).trim();
  manifest.steady.revision = revision;
  manifest.steady.sha256 = digest;
  fs.writeFileSync(manifestPath, `${JSON.stringify(manifest, null, 2)}\n`);
  console.log(
    `Updated Steady to ${revision}. Run ./scripts/steady/install and node scripts/steady/test.cjs to verify it.`,
  );
} finally {
  if (fs.rmSync) {
    fs.rmSync(temporary, { recursive: true, force: true });
  } else {
    fs.rmdirSync(temporary, { recursive: true });
  }
}
