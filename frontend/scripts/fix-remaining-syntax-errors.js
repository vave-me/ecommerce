const fs = require('fs');
const path = require('path');

// Files with known syntax errors
const filesToFix = [
    './src/api/client/entityApi.jsx',
    './src/api/client/productsApi.jsx',
    './src/components/Feed/FeedItem.client.jsx',
    './src/utils/logger.js',
    './src/app/[locale]/services/[category]/[slug]/page.jsx'
];

// Pattern to find broken syntax like "}:`, error..." or "}: ${url}"
const brokenPatterns = [
    /\}:\s*`,\s*err(or)?\.{3}/g,
    /\}:\s*\$\{[^}]+\}`/g,
    /\}:\s*`[^`]*`\s*,/g
];

function fixFile(filePath) {
    try {
        let content = fs.readFileSync(filePath, 'utf8');
        let fixed = false;
        
        // Fix each broken pattern
        brokenPatterns.forEach(pattern => {
            if (pattern.test(content)) {
                content = content.replace(pattern, (match) => {
                    console.log(`  Found broken pattern in ${filePath}: ${match.substring(0, 30)}...`);
                    // Just remove the broken part
                    return '';
                });
                fixed = true;
            }
        });
        
        // Fix specific FeedItem.client.jsx issue
        if (filePath.includes('FeedItem.client.jsx')) {
            // Look for the broken catch block
            const catchPattern = /}\s*catch\s*\(error\)\s*{[^}]*}\s*}\s*{\s*return/g;
            if (catchPattern.test(content)) {
                content = content.replace(catchPattern, '} catch (error) {\n        // Error logged for debugging\n        if (process.env.NODE_ENV === \'development\') {\n            console.error(\'Error:\', error);\n        }\n        return');
                fixed = true;
                console.log(`  Fixed catch block structure in ${filePath}`);
            }
        }
        
        // Fix logger.js specific issues
        if (filePath.includes('logger.js')) {
            // Fix the broken console.log in apiLogger
            const apiLoggerPattern = /}\s*:\s*\$\{url\}`\s*,\s*data\s*\?\s*{\s*data\s*}\s*:\s*''\);/g;
            if (apiLoggerPattern.test(content)) {
                content = content.replace(apiLoggerPattern, 'console.log(`📡 API ${method}: ${url}`, data ? { data } : \'\');');
                fixed = true;
                console.log(`  Fixed apiLogger console.log in ${filePath}`);
            }
        }
        
        if (fixed) {
            fs.writeFileSync(filePath, content, 'utf8');
            console.log(`✅ Fixed ${filePath}`);
        } else {
            console.log(`ℹ️  No issues found in ${filePath}`);
        }
        
    } catch (error) {
        console.error(`❌ Error processing ${filePath}:`, error.message);
    }
}

console.log('🔧 Fixing remaining syntax errors...\n');

filesToFix.forEach(filePath => {
    const fullPath = path.join(__dirname, '..', filePath);
    if (fs.existsSync(fullPath)) {
        fixFile(fullPath);
    } else {
        console.log(`⚠️  File not found: ${filePath}`);
    }
});

console.log('\n✨ Done!');