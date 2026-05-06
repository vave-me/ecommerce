#!/usr/bin/env node

/**
 * Script to fix remaining syntax errors from console removal
 */

const fs = require('fs');
const path = require('path');

// Files with known syntax errors
const filesToFix = [
    {
        file: 'src/api/client/mediaApi.jsx',
        pattern: /\/\/ Error:.*?\n\s*\.\.\./g,
        replacement: '// Error logged for debugging'
    },
    {
        file: 'src/api/client/userApi.jsx',
        pattern: /\/\/ Error:.*?\n\s*\.\.\./g,
        replacement: '// Error logged for debugging'
    },
    {
        file: 'src/components/wishlist/WishlistItemCard.jsx',
        pattern: /\);\s*\n\s*\);/g,
        replacement: ');'
    },
    {
        file: 'src/context/AuthContext.jsx',
        pattern: /\n\s*\.toISOString\(\)[^}]*}/g,
        replacement: '\n                // Token info logged for debugging\n            }'
    }
];

function fixFile(fileInfo) {
    const filePath = path.join(process.cwd(), fileInfo.file);
    
    try {
        let content = fs.readFileSync(filePath, 'utf8');
        const originalContent = content;
        
        // Apply the pattern replacement
        content = content.replace(fileInfo.pattern, fileInfo.replacement);
        
        // Additional cleanup for broken syntax
        // Fix stray ellipsis
        content = content.replace(/\s*\.\.\.\s*\n/g, '\n');
        
        // Fix broken object notation
        content = content.replace(/,\s*\n\s*\.\.\./g, '');
        
        // Fix broken function calls
        content = content.replace(/\);\s*\);/g, ');');
        
        if (content !== originalContent) {
            fs.writeFileSync(filePath, content, 'utf8');
            console.log(`✅ Fixed: ${fileInfo.file}`);
            return true;
        } else {
            console.log(`ℹ️  No changes needed: ${fileInfo.file}`);
            return false;
        }
    } catch (error) {
        console.error(`❌ Error fixing ${fileInfo.file}:`, error.message);
        return false;
    }
}

// Additional function to find and fix common patterns
function findAndFixCommonPatterns() {
    const commonPatterns = [
        {
            pattern: /\/\/ Error:.*?\n\s*\.\.\.\s*\n/g,
            replacement: '// Error logged for debugging\n'
        },
        {
            pattern: /\/\/ Error:.*?\{[\s\S]*?\.\.\.\s*\}/g,
            replacement: '// Error details logged for debugging'
        },
        {
            pattern: /\);\s*\n\s*\);/g,
            replacement: ');'
        }
    ];
    
    const directories = ['src/api', 'src/components', 'src/context', 'src/features'];
    let totalFixed = 0;
    
    directories.forEach(dir => {
        const dirPath = path.join(process.cwd(), dir);
        if (!fs.existsSync(dirPath)) return;
        
        const files = findJsFiles(dirPath);
        
        files.forEach(file => {
            try {
                let content = fs.readFileSync(file, 'utf8');
                const originalContent = content;
                
                commonPatterns.forEach(({ pattern, replacement }) => {
                    content = content.replace(pattern, replacement);
                });
                
                if (content !== originalContent) {
                    fs.writeFileSync(file, content, 'utf8');
                    console.log(`✅ Fixed common patterns in: ${path.relative(process.cwd(), file)}`);
                    totalFixed++;
                }
            } catch (error) {
                console.error(`❌ Error processing ${file}:`, error.message);
            }
        });
    });
    
    return totalFixed;
}

function findJsFiles(dir, fileList = []) {
    const files = fs.readdirSync(dir);
    
    files.forEach(file => {
        const filePath = path.join(dir, file);
        const stat = fs.statSync(filePath);
        
        if (stat.isDirectory()) {
            if (!file.startsWith('.') && file !== 'node_modules') {
                findJsFiles(filePath, fileList);
            }
        } else if (file.endsWith('.js') || file.endsWith('.jsx')) {
            fileList.push(filePath);
        }
    });
    
    return fileList;
}

console.log('🔧 Fixing syntax errors from console removal...\n');

// Fix specific known files
let fixedCount = 0;
filesToFix.forEach(fileInfo => {
    if (fixFile(fileInfo)) {
        fixedCount++;
    }
});

console.log('\n🔍 Looking for common error patterns...\n');

// Find and fix common patterns
const additionalFixed = findAndFixCommonPatterns();

console.log('\n✨ Syntax error fix complete!');
console.log(`📊 Files fixed: ${fixedCount + additionalFixed}`);