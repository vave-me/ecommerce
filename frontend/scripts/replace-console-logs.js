#!/usr/bin/env node

/**
 * Script to replace console.* statements with logger utility
 * Run: node scripts/replace-console-logs.js
 */

const fs = require('fs');
const path = require('path');
const glob = require('glob');

// Directories to process
const DIRECTORIES = [
  'src/**/*.{js,jsx,ts,tsx}',
  '!src/utils/logger.js', // Exclude the logger utility itself
  '!src/**/*.test.{js,jsx,ts,tsx}', // Exclude test files
  '!src/**/__tests__/**', // Exclude test directories
];

// Patterns to replace
const REPLACEMENTS = [
  {
    pattern: /console\.log\(/g,
    replacement: 'logger.log(',
    import: 'logger',
  },
  {
    pattern: /console\.error\(/g,
    replacement: 'logger.error(',
    import: 'logger',
  },
  {
    pattern: /console\.warn\(/g,
    replacement: 'logger.warn(',
    import: 'logger',
  },
  {
    pattern: /console\.debug\(/g,
    replacement: 'logger.debug(',
    import: 'logger',
  },
  {
    pattern: /console\.info\(/g,
    replacement: 'logger.info(',
    import: 'logger',
  },
  {
    pattern: /console\.time\(/g,
    replacement: 'logger.time(',
    import: 'logger',
  },
  {
    pattern: /console\.timeEnd\(/g,
    replacement: 'logger.timeEnd(',
    import: 'logger',
  },
  {
    pattern: /console\.table\(/g,
    replacement: 'logger.table(',
    import: 'logger',
  },
  {
    pattern: /console\.group\(/g,
    replacement: 'logger.group(',
    import: 'logger',
  },
  {
    pattern: /console\.groupEnd\(/g,
    replacement: 'logger.groupEnd(',
    import: 'logger',
  },
];

// Special replacements for API files
const API_REPLACEMENTS = [
  {
    pattern: /console\.log\((.*API.*)\)/g,
    replacement: 'apiLogger.request($1)',
    import: '{ apiLogger }',
  },
  {
    pattern: /console\.error\((.*API.*)\)/g,
    replacement: 'apiLogger.error($1)',
    import: '{ apiLogger }',
  },
];

let filesProcessed = 0;
let replacementsMade = 0;

/**
 * Add import statement if not present
 */
function addImportIfNeeded(content, importName, filePath) {
  const loggerPath = path.relative(path.dirname(filePath), 'src/utils/logger').replace(/\\/g, '/');
  const importStatement = `import ${importName} from '${loggerPath.startsWith('.') ? loggerPath : './' + loggerPath}';`;
  
  // Check if import already exists
  const importRegex = new RegExp(`import.*from.*logger.*`);
  if (importRegex.test(content)) {
    return content;
  }
  
  // Add import after other imports or at the beginning
  const importMatch = content.match(/^(import[\s\S]*?;)/m);
  if (importMatch) {
    const lastImportIndex = content.lastIndexOf(importMatch[0]) + importMatch[0].length;
    return content.slice(0, lastImportIndex) + '\n' + importStatement + content.slice(lastImportIndex);
  } else {
    // If no imports, add at the beginning after 'use client' if present
    const useClientMatch = content.match(/^["']use client["'];?\s*\n/);
    if (useClientMatch) {
      return useClientMatch[0] + importStatement + '\n' + content.slice(useClientMatch[0].length);
    }
    return importStatement + '\n\n' + content;
  }
}

/**
 * Process a single file
 */
function processFile(filePath) {
  try {
    let content = fs.readFileSync(filePath, 'utf8');
    const originalContent = content;
    let needsImport = false;
    let importType = 'logger';
    
    // Check if it's an API file
    const isApiFile = filePath.includes('/api/');
    const replacements = isApiFile ? [...REPLACEMENTS, ...API_REPLACEMENTS] : REPLACEMENTS;
    
    // Apply replacements
    replacements.forEach(({ pattern, replacement, import: importName }) => {
      if (pattern.test(content)) {
        content = content.replace(pattern, replacement);
        needsImport = true;
        if (importName && importName !== 'logger') {
          importType = importName;
        }
      }
    });
    
    // If changes were made
    if (content !== originalContent) {
      // Add import if needed
      if (needsImport) {
        content = addImportIfNeeded(content, importType, filePath);
      }
      
      // Write back to file
      fs.writeFileSync(filePath, content, 'utf8');
      filesProcessed++;
      
      // Count replacements
      const originalMatches = originalContent.match(/console\.(log|error|warn|debug|info|time|timeEnd|table|group|groupEnd)\(/g) || [];
      replacementsMade += originalMatches.length;
      
      console.log(`✅ Processed: ${filePath} (${originalMatches.length} replacements)`);
    }
  } catch (error) {
    console.error(`❌ Error processing ${filePath}:`, error.message);
  }
}

/**
 * Main function
 */
async function main() {
  console.log('🔍 Starting console.log replacement...\n');
  
  // Find all files
  const files = [];
  for (const pattern of DIRECTORIES) {
    if (pattern.startsWith('!')) {
      // Ignore pattern
      continue;
    }
    const matches = glob.sync(pattern, { 
      ignore: DIRECTORIES.filter(p => p.startsWith('!')).map(p => p.slice(1)) 
    });
    files.push(...matches);
  }
  
  console.log(`Found ${files.length} files to process\n`);
  
  // Process each file
  files.forEach(processFile);
  
  console.log('\n✨ Replacement complete!');
  console.log(`📊 Files processed: ${filesProcessed}`);
  console.log(`🔄 Total replacements: ${replacementsMade}`);
  
  // Verify no console statements remain
  console.log('\n🔍 Verifying no console statements remain...');
  const remainingConsoles = [];
  
  files.forEach(filePath => {
    const content = fs.readFileSync(filePath, 'utf8');
    const matches = content.match(/console\.(log|error|warn|debug|info)\(/g);
    if (matches) {
      remainingConsoles.push({ file: filePath, count: matches.length });
    }
  });
  
  if (remainingConsoles.length > 0) {
    console.log('\n⚠️  Some console statements remain:');
    remainingConsoles.forEach(({ file, count }) => {
      console.log(`  - ${file}: ${count} statements`);
    });
  } else {
    console.log('✅ All console statements have been replaced!');
  }
}

// Run the script
main().catch(console.error);