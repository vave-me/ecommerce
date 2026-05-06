const fs = require('fs');
const path = require('path');

// Files identified as needing cleanup
const filesToFix = [
    'src/app/[locale]/design/shortsSingle/page.jsx',
    'src/app/[locale]/user/page.jsx',
    'src/components/Header/AddDropdown.jsx',
    'src/components/Header/CategoryBar.jsx',
    'src/components/Header/SuggestionsList.jsx',
    'src/components/Header/TopicButtonsRow.jsx',
    'src/hooks/useFilterState.jsx',
    'src/hooks/useHeaderScroll.jsx',
    'src/hooks/useScrollNavigation.jsx',
    'src/utils/performanceOptimizer.js'
];

function fixFile(filePath) {
    try {
        let content = fs.readFileSync(filePath, 'utf8');
        let modified = false;
        
        // Fix setTimeout without cleanup
        content = content.replace(
            /(useEffect\s*\(\s*(?:\(\)|async\s*\(\))\s*=>\s*{[^}]*?)(setTimeout\s*\([^)]+\)\s*;)([^}]*?)(}\s*(?:,\s*\[[^\]]*\])?\s*\))/g,
            (match, before, setTimeoutCall, after, end) => {
                if (!after.includes('return') && !before.includes('const timeoutId')) {
                    modified = true;
                    // Extract the setTimeout call
                    const timeoutVar = `const timeoutId = ${setTimeoutCall}`;
                    const cleanup = `\n    return () => clearTimeout(timeoutId);\n  `;
                    return `${before}${timeoutVar}${after}${cleanup}${end}`;
                }
                return match;
            }
        );
        
        // Fix addEventListener without cleanup
        content = content.replace(
            /(useEffect\s*\(\s*(?:\(\)|async\s*\(\))\s*=>\s*{[^}]*?)(\w+\.addEventListener\s*\(\s*['"](\w+)['"]\s*,\s*(\w+)[^)]*\)\s*;)([^}]*?)(}\s*(?:,\s*\[[^\]]*\])?\s*\))/g,
            (match, before, addListenerCall, eventName, handlerName, after, end) => {
                if (!after.includes('return') && !after.includes('removeEventListener')) {
                    modified = true;
                    const target = addListenerCall.split('.')[0];
                    const cleanup = `\n    return () => ${target}.removeEventListener('${eventName}', ${handlerName});\n  `;
                    return `${before}${addListenerCall}${after}${cleanup}${end}`;
                }
                return match;
            }
        );
        
        // Fix setInterval without cleanup
        content = content.replace(
            /(useEffect\s*\(\s*(?:\(\)|async\s*\(\))\s*=>\s*{[^}]*?)(setInterval\s*\([^)]+\)\s*;)([^}]*?)(}\s*(?:,\s*\[[^\]]*\])?\s*\))/g,
            (match, before, setIntervalCall, after, end) => {
                if (!after.includes('return') && !before.includes('const intervalId')) {
                    modified = true;
                    const intervalVar = `const intervalId = ${setIntervalCall}`;
                    const cleanup = `\n    return () => clearInterval(intervalId);\n  `;
                    return `${before}${intervalVar}${after}${cleanup}${end}`;
                }
                return match;
            }
        );
        
        if (modified) {
            fs.writeFileSync(filePath, content, 'utf8');
            console.log(`✅ Fixed ${filePath}`);
            return true;
        } else {
            console.log(`ℹ️  No automatic fixes applied to ${filePath} - may need manual review`);
            return false;
        }
        
    } catch (error) {
        console.error(`❌ Error processing ${filePath}:`, error.message);
        return false;
    }
}

console.log('🔧 Fixing useEffect cleanup issues...\n');

let totalFixed = 0;
let needsManualReview = [];

filesToFix.forEach(relativePath => {
    const fullPath = path.join(__dirname, '..', relativePath);
    if (fs.existsSync(fullPath)) {
        if (fixFile(fullPath)) {
            totalFixed++;
        } else {
            needsManualReview.push(relativePath);
        }
    } else {
        console.log(`⚠️  File not found: ${relativePath}`);
    }
});

console.log(`\n✨ Summary:`);
console.log(`   Files automatically fixed: ${totalFixed}`);
console.log(`   Files needing manual review: ${needsManualReview.length}`);

if (needsManualReview.length > 0) {
    console.log('\n📋 Files needing manual review:');
    needsManualReview.forEach(file => {
        console.log(`   - ${file}`);
    });
}