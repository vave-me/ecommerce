#!/usr/bin/env node

/**
 * This script ensures the ads.txt file is properly copied to the output directory
 * for Next.js application deployment.
 */

const fs = require('fs');
const path = require('path');

// Configuration
const PUBLIC_DIR = path.join(process.cwd(), 'public');
const BUILD_DIR_OUTPUT = path.join(process.cwd(), '.next/static');
const SRC_ADS_TXT = path.join(PUBLIC_DIR, 'ads.txt');
const DEST_ADS_TXT = path.join(BUILD_DIR_OUTPUT, 'ads.txt');

// Ensure the build output directory exists
function ensureDirectoryExists(directory) {
  if (!fs.existsSync(directory)) {
    console.log(`Creating directory: ${directory}`);
    fs.mkdirSync(directory, { recursive: true });
  }
}

// Copy the ads.txt file
function copyAdsTxt() {
  try {
    // Check if source file exists
    if (!fs.existsSync(SRC_ADS_TXT)) {
      console.error('❌ Source ads.txt file not found at:', SRC_ADS_TXT);
      return false;
    }

    // Ensure destination directory exists
    ensureDirectoryExists(BUILD_DIR_OUTPUT);

    // Copy the file
    fs.copyFileSync(SRC_ADS_TXT, DEST_ADS_TXT);
    console.log(`✅ ads.txt copied to: ${DEST_ADS_TXT}`);
    return true;
  } catch (error) {
    console.error('❌ Error copying ads.txt file:', error.message);
    return false;
  }
}

// Main execution
console.log('📝 Starting ads.txt copy process...');
const success = copyAdsTxt();
console.log('📝 ads.txt copy process complete.');

// Also copy to root of .next directory
try {
  const ROOT_DEST = path.join(process.cwd(), '.next', 'ads.txt');
  fs.copyFileSync(SRC_ADS_TXT, ROOT_DEST);
  console.log(`✅ ads.txt also copied to: ${ROOT_DEST}`);
} catch (error) {
  console.error('❌ Error copying ads.txt to root of .next:', error.message);
}

// Also check if we need to place in standalone directory for container deployments
const STANDALONE_DIR = path.join(process.cwd(), '.next/standalone');
if (fs.existsSync(STANDALONE_DIR)) {
  try {
    const STANDALONE_DEST = path.join(STANDALONE_DIR, 'ads.txt');
    fs.copyFileSync(SRC_ADS_TXT, STANDALONE_DEST);
    console.log(`✅ ads.txt also copied to standalone directory: ${STANDALONE_DEST}`);
  } catch (error) {
    console.error('❌ Error copying ads.txt to standalone directory:', error.message);
  }
}

// Exit with appropriate status code
process.exit(success ? 0 : 1); 