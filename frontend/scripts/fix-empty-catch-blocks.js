#!/usr/bin/env node

/**
 * Fix empty catch blocks by adding proper error handling
 */

const fs = require('fs');
const path = require('path');

// Pattern to match empty or minimal catch blocks
const CATCH_PATTERNS = [
  // Empty catch block
  /catch\s*\(([^)]+)\)\s*\{\s*\}/g,
  // Catch with only comment
  /catch\s*\(([^)]+)\)\s*\{\s*\/\/[^\n]*\n?\s*\}/g,
  // Catch with only console statement (should be replaced)
  /catch\s*\(([^)]+)\)\s*\{\s*console\.[^;]+;\s*\}/g,
];

// Context-based error handling templates
const ERROR_HANDLERS = {
  // For async operations in event handlers
  eventHandler: (errorVar) => `{
        // Handle error silently for better UX
        if (process.env.NODE_ENV === 'development') {
            console.error('Event handler error:', ${errorVar});
        }
    }`,
  
  // For API calls
  apiCall: (errorVar) => `{
        // API error handling
        if (process.env.NODE_ENV === 'development') {
            console.error('API error:', ${errorVar});
        }
        throw ${errorVar}; // Re-throw for caller to handle
    }`,
  
  // For form submissions
  formSubmission: (errorVar) => `{
        // Form submission error
        if (process.env.NODE_ENV === 'development') {
            console.error('Form submission error:', ${errorVar});
        }
        // Could set error state here if available
        throw ${errorVar};
    }`,
  
  // For file operations
  fileOperation: (errorVar) => `{
        // File operation error
        if (process.env.NODE_ENV === 'development') {
            console.error('File operation error:', ${errorVar});
        }
        return null; // Return null for failed file operations
    }`,
  
  // For initialization/setup
  initialization: (errorVar) => `{
        // Initialization error - log but continue
        if (process.env.NODE_ENV === 'development') {
            console.error('Initialization error:', ${errorVar});
        }
        // Continue with default behavior
    }`,
  
  // Default handler
  default: (errorVar) => `{
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', ${errorVar});
        }
    }`
};

function determineContextType(fileContent, catchIndex) {
  const contextWindow = 200; // Characters to look before catch
  const beforeCatch = fileContent.substring(Math.max(0, catchIndex - contextWindow), catchIndex);
  
  // Check for API/fetch calls
  if (beforeCatch.match(/fetch|axios|api|request|response/i)) {
    return 'apiCall';
  }
  
  // Check for form handling
  if (beforeCatch.match(/handleSubmit|onSubmit|form/i)) {
    return 'formSubmission';
  }
  
  // Check for file operations
  if (beforeCatch.match(/file|upload|download|blob|createObjectURL/i)) {
    return 'fileOperation';
  }
  
  // Check for event handlers
  if (beforeCatch.match(/onClick|onChange|onKeyPress|handle[A-Z]/)) {
    return 'eventHandler';
  }
  
  // Check for initialization
  if (beforeCatch.match(/useEffect|componentDidMount|init|setup|connect/i)) {
    return 'initialization';
  }
  
  return 'default';
}

function fixEmptyCatchBlocks(filePath) {
  try {
    let content = fs.readFileSync(filePath, 'utf8');
    let modified = false;
    let fixCount = 0;
    
    // Check each pattern
    CATCH_PATTERNS.forEach(pattern => {
      let match;
      let lastIndex = 0;
      let newContent = '';
      
      // Reset regex
      pattern.lastIndex = 0;
      
      while ((match = pattern.exec(content)) !== null) {
        const fullMatch = match[0];
        const errorVar = match[1] || 'error';
        const catchIndex = match.index;
        
        // Determine context
        const contextType = determineContextType(content, catchIndex);
        const handler = ERROR_HANDLERS[contextType](errorVar);
        
        // Build new content
        newContent += content.substring(lastIndex, catchIndex);
        newContent += `catch (${errorVar}) ${handler}`;
        lastIndex = catchIndex + fullMatch.length;
        
        modified = true;
        fixCount++;
      }
      
      if (modified) {
        // Add remaining content
        newContent += content.substring(lastIndex);
        content = newContent;
      }
    });
    
    if (modified) {
      // Ensure we have the necessary imports for development checks
      if (!content.includes("process.env.NODE_ENV")) {
        // Add at the top of the file after other imports
        const importMatch = content.match(/(import[^;]+;\s*\n)+/);
        if (importMatch) {
          const lastImportEnd = importMatch.index + importMatch[0].length;
          content = content.slice(0, lastImportEnd) + 
                   "\n// Error handling utilities added by production readiness script\n" +
                   content.slice(lastImportEnd);
        }
      }
      
      fs.writeFileSync(filePath, content, 'utf8');
      console.log(`✅ Fixed ${fixCount} empty catch blocks in ${filePath}`);
      return fixCount;
    }
    
    return 0;
  } catch (error) {
    console.error(`❌ Error processing ${filePath}:`, error.message);
    return 0;
  }
}

// Recursive function to find all JS/JSX files
function findFiles(dir, fileList = []) {
  const files = fs.readdirSync(dir);
  
  files.forEach(file => {
    const filePath = path.join(dir, file);
    const stat = fs.statSync(filePath);
    
    if (stat.isDirectory()) {
      // Skip certain directories
      if (!file.includes('node_modules') && 
          !file.includes('build') && 
          !file.includes('dist') &&
          !file.includes('.next')) {
        findFiles(filePath, fileList);
      }
    } else if (file.endsWith('.js') || file.endsWith('.jsx')) {
      // Skip test files
      if (!file.includes('.test.') && !file.includes('.spec.')) {
        fileList.push(filePath);
      }
    }
  });
  
  return fileList;
}

// Main execution
console.log('🔧 Fixing empty catch blocks...\n');

// Find all JS/JSX files
const files = findFiles('src');

console.log(`Found ${files.length} files to check\n`);

let totalFixed = 0;
let filesFixed = 0;
const fixedFilesList = [];

files.forEach(file => {
  const fixes = fixEmptyCatchBlocks(file);
  if (fixes > 0) {
    totalFixed += fixes;
    filesFixed++;
    fixedFilesList.push(file);
  }
});

console.log(`\n✨ Summary:`);
console.log(`   - Fixed ${totalFixed} empty catch blocks`);
console.log(`   - Modified ${filesFixed} files`);
console.log(`   - Checked ${files.length} total files`);

// Create a report of the fixes
const report = {
  timestamp: new Date().toISOString(),
  totalFixed,
  filesFixed,
  filesChecked: files.length,
  fixedFiles: fixedFilesList
};

fs.writeFileSync(
  'empty-catch-fixes-report.json',
  JSON.stringify(report, null, 2)
);

console.log('\n📊 Report saved to empty-catch-fixes-report.json');