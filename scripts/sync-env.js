const fs = require('fs');
const path = require('path');

const rootEnv = path.resolve(__dirname, '..', '.env');
const backendEnv = path.resolve(__dirname, '..', 'backend', '.env');

function copyIfExists(src, dest) {
  if (!fs.existsSync(src)) {
    console.error(`Source file not found: ${src}`);
    process.exitCode = 1;
    return;
  }
  try {
    fs.copyFileSync(src, dest);
    console.log(`Copied ${src} -> ${dest}`);
  } catch (err) {
    console.error('Failed to copy file:', err.message);
    process.exitCode = 2;
  }
}

copyIfExists(rootEnv, backendEnv);

console.log('Done. If you want reverse direction, swap paths in this script.');
