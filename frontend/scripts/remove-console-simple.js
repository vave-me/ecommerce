#!/usr/bin/env node

/**
 * Simple script to remove console statements
 * Uses built-in Node.js modules only
 */

const fs = require('fs');
const path = require('path');

let filesProcessed = 0;
let statementsRemoved = 0;

/**
 * Recursively find all JS/JSX files
 */
function findFiles(dir, fileList = []) {
  const files = fs.readdirSync(dir);
  
  files.forEach(file => {
    const filePath = path.join(dir, file);
    const stat = fs.statSync(filePath);
    
    if (stat.isDirectory()) {
      // Skip node_modules and other directories
      if (!file.startsWith('.') && file !== 'node_modules' && file !== 'build' && file !== 'dist') {
        findFiles(filePath, fileList);
      }
    } else if (file.endsWith('.js') || file.endsWith('.jsx') || file.endsWith('.ts') || file.endsWith('.tsx')) {
      // Skip test files and logger utility
      if (!file.includes('.test.') && !file.includes('.spec.') && filePath !== 'src/utils/logger.js') {
        fileList.push(filePath);
      }
    }
  });
  
  return fileList;
}

/**
 * Remove console statements from a file
 */
function processFile(filePath) {
  try {
    let content = fs.readFileSync(filePath, 'utf8');
    const originalContent = content;
    
    // Remove console.log statements (including multi-line)
    content = content.replace(/console\s*\.\s*log\s*\([^)]*\)\s*;?/g, '');
    
    // Remove console.error but keep the error tracking logic
    content = content.replace(/console\s*\.\s*error\s*\(([^)]*)\)\s*;?/g, (match, args) => {
      // If it's in a catch block, keep minimal error handling
      if (content.includes('catch')) {
        return `// Error: ${args.slice(0, 50)}...`;
      }
      return '';
    });
    
    // Remove other console methods
    content = content.replace(/console\s*\.\s*(warn|debug|info|time|timeEnd|table|group|groupEnd)\s*\([^)]*\)\s*;?/g, '');
    
    // Clean up empty lines
    content = content.replace(/\n\s*\n\s*\n/g, '\n\n');
    
    // If changes were made
    if (content !== originalContent) {
      fs.writeFileSync(filePath, content, 'utf8');
      filesProcessed++;
      
      // Count removals
      const matches = originalContent.match(/console\s*\.\s*(log|error|warn|debug|info)\s*\(/g) || [];
      statementsRemoved += matches.length;
      
      console.log(`✅ Processed: ${filePath} (${matches.length} statements removed)`);
    }
  } catch (error) {
    console.error(`❌ Error processing ${filePath}:`, error.message);
  }
}

/**
 * Main function
 */
function main() {
  console.log('🔍 Starting console statement removal...\n');
  
  // Find all files in src directory
  const srcPath = path.join(process.cwd(), 'src');
  const files = findFiles(srcPath);
  
  console.log(`Found ${files.length} files to process\n`);
  
  // Process each file
  files.forEach(processFile);
  
  console.log('\n✨ Removal complete!');
  console.log(`📊 Files processed: ${filesProcessed}`);
  console.log(`🔄 Total statements removed: ${statementsRemoved}`);
}

// Run the script
main();