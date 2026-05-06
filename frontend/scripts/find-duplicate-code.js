const fs = require('fs');
const path = require('path');
const crypto = require('crypto');

// Common patterns that indicate duplicate code
const duplicatePatterns = [
    // API call patterns
    /try\s*{\s*const\s+response\s*=\s*await\s+\w+\.\w+\(/g,
    // Error handling patterns
    /catch\s*\(error\)\s*{\s*if\s*\(process\.env\.NODE_ENV/g,
    // Loading states
    /if\s*\(loading\)\s*return\s*<div[^>]*>Loading/g,
    // Form validation
    /if\s*\(!\w+\s*\|\|\s*\w+\.length\s*===\s*0\)/g,
    // Price formatting
    /toFixed\(2\).*€/g,
    // Date formatting
    /new\s+Date\([^)]+\)\.toLocal/g
];

const codeSnippets = new Map();
const filesByPattern = new Map();

function extractCodeBlocks(content, filePath) {
    // Extract function bodies
    const functionPattern = /(?:function\s+\w+|const\s+\w+\s*=\s*(?:async\s*)?(?:\([^)]*\)|[^=])\s*=>)\s*{([^}]+(?:{[^}]*}[^}]*)*)}/g;
    const matches = [...content.matchAll(functionPattern)];
    
    matches.forEach(match => {
        const body = match[1].trim();
        // Skip very short functions
        if (body.length < 50) return;
        
        // Create a normalized version for comparison
        const normalized = body
            .replace(/\s+/g, ' ')
            .replace(/['"`]/g, '"')
            .replace(/\w+:/g, 'key:') // Normalize object keys
            .replace(/\d+/g, 'NUM'); // Normalize numbers
            
        const hash = crypto.createHash('md5').update(normalized).digest('hex');
        
        if (!codeSnippets.has(hash)) {
            codeSnippets.set(hash, {
                snippet: body.substring(0, 200),
                files: [],
                count: 0
            });
        }
        
        const data = codeSnippets.get(hash);
        data.files.push(filePath);
        data.count++;
    });
    
    // Check for specific patterns
    duplicatePatterns.forEach((pattern, index) => {
        const matches = content.match(pattern);
        if (matches && matches.length > 0) {
            const key = `pattern_${index}`;
            if (!filesByPattern.has(key)) {
                filesByPattern.set(key, {
                    pattern: pattern.source,
                    files: [],
                    count: 0
                });
            }
            const data = filesByPattern.get(key);
            data.files.push(filePath);
            data.count += matches.length;
        }
    });
}

function walkDir(dir) {
    const files = fs.readdirSync(dir);
    files.forEach(file => {
        const filePath = path.join(dir, file);
        const stat = fs.statSync(filePath);
        
        if (stat.isDirectory() && !file.includes('node_modules') && !file.startsWith('.')) {
            walkDir(filePath);
        } else if ((file.endsWith('.jsx') || file.endsWith('.js')) && !file.includes('.test.')) {
            try {
                const content = fs.readFileSync(filePath, 'utf8');
                extractCodeBlocks(content, filePath);
            } catch (err) {
                // Skip files that can't be read
            }
        }
    });
}

console.log('🔍 Analyzing code for duplicates...\n');

const srcDir = path.join(__dirname, '..', 'src');
walkDir(srcDir);

// Find actual duplicates (appearing in 2+ files)
const duplicates = Array.from(codeSnippets.entries())
    .filter(([hash, data]) => data.files.length > 1)
    .sort((a, b) => b[1].count - a[1].count);

console.log('📊 Duplicate Code Blocks Found:\n');

duplicates.slice(0, 10).forEach(([hash, data], index) => {
    console.log(`${index + 1}. Found in ${data.files.length} files (${data.count} occurrences):`);
    console.log(`   Preview: ${data.snippet.substring(0, 100)}...`);
    console.log(`   Files:`);
    data.files.slice(0, 3).forEach(file => {
        console.log(`   - ${path.relative(path.join(__dirname, '..'), file)}`);
    });
    if (data.files.length > 3) {
        console.log(`   ... and ${data.files.length - 3} more`);
    }
    console.log('');
});

console.log('\n📊 Common Patterns:\n');

const commonPatterns = Array.from(filesByPattern.entries())
    .filter(([key, data]) => data.files.length > 3)
    .sort((a, b) => b[1].count - a[1].count);

commonPatterns.forEach(([key, data], index) => {
    console.log(`${index + 1}. Pattern: ${data.pattern.substring(0, 50)}...`);
    console.log(`   Found in ${data.files.length} files (${data.count} total occurrences)`);
    console.log('');
});

console.log('\n💡 Recommendations:');
console.log('1. Create shared utility functions for common operations');
console.log('2. Extract common API error handling into a wrapper');
console.log('3. Create reusable loading and error components');
console.log('4. Standardize form validation logic');
console.log('5. Create shared formatting utilities (price, date, etc.)');

console.log('\n✨ Analysis complete!');