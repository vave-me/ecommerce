const fs = require('fs');
const path = require('path');

// Card components that should use React.memo
const componentsToMemoize = [
    'src/components/classified/ClassifiedCard.jsx',
    'src/components/services/ServiceCard.jsx',
    'src/components/PostCard/PostCard.jsx',
    'src/components/wishlist/WishlistItemCard.jsx',
    'src/components/Listings/ListingCard.jsx',
    'src/components/Feed/AIResponseCard.jsx',
    'src/components/shared/CardTitle.jsx',
    'src/components/shared/CardImageContainer.jsx',
    'src/components/AI/ProductCard.jsx'
];

function addMemoToComponent(filePath) {
    try {
        let content = fs.readFileSync(filePath, 'utf8');
        
        // Skip if already memoized
        if (content.includes('React.memo') || content.includes('memo(')) {
            return { status: 'already-memoized' };
        }
        
        // Check if React import includes memo
        const hasReactImport = content.includes('import React');
        const hasMemoImport = content.includes('memo');
        
        // Add memo to React import if needed
        if (hasReactImport && !hasMemoImport) {
            content = content.replace(
                /import\s+React(\s*,\s*{[^}]*})?/,
                (match, namedImports) => {
                    if (namedImports) {
                        // Add memo to existing named imports
                        return match.replace('}', ', memo }');
                    } else {
                        // Add named import
                        return 'import React, { memo }';
                    }
                }
            );
        } else if (!hasReactImport) {
            // Add import at the top
            content = `import { memo } from 'react';\n${content}`;
        }
        
        // Find the component export and wrap with memo
        const patterns = [
            // export default function ComponentName
            /export\s+default\s+function\s+(\w+)/,
            // export default ComponentName
            /export\s+default\s+(\w+)(?![\w.])/,
            // const ComponentName = () => {}; export default ComponentName;
            /const\s+(\w+)\s*=\s*(?:\([^)]*\)|[^=])\s*=>\s*[\s\S]*?export\s+default\s+\1/
        ];
        
        let modified = false;
        
        // Try pattern 1: export default function
        content = content.replace(
            /export\s+default\s+function\s+(\w+)/,
            (match, componentName) => {
                modified = true;
                return `function ${componentName}`;
            }
        );
        
        if (modified) {
            // Add memo export at the end
            const componentMatch = content.match(/function\s+(\w+)/);
            if (componentMatch) {
                const componentName = componentMatch[1];
                content += `\n\nexport default memo(${componentName});`;
            }
        } else {
            // Try pattern 2: const Component = 
            const constPattern = /const\s+(\w+)\s*=\s*(?:\([^)]*\)|[^=])\s*=>/;
            const constMatch = content.match(constPattern);
            
            if (constMatch) {
                const componentName = constMatch[1];
                // Replace export default ComponentName with memo wrapped version
                content = content.replace(
                    new RegExp(`export\\s+default\\s+${componentName}(?![\\w.])`),
                    `export default memo(${componentName})`
                );
                modified = true;
            }
        }
        
        if (modified) {
            fs.writeFileSync(filePath, content, 'utf8');
            return { status: 'success' };
        } else {
            return { status: 'no-changes' };
        }
        
    } catch (error) {
        return { status: 'error', error: error.message };
    }
}

console.log('🎯 Adding React.memo to card components...\n');

let successCount = 0;
let alreadyMemoized = 0;
let errors = 0;

componentsToMemoize.forEach(relativePath => {
    const fullPath = path.join(__dirname, '..', relativePath);
    
    if (!fs.existsSync(fullPath)) {
        console.log(`⚠️  File not found: ${relativePath}`);
        errors++;
        return;
    }
    
    const result = addMemoToComponent(fullPath);
    
    switch (result.status) {
        case 'success':
            console.log(`✅ Added memo to ${path.basename(relativePath)}`);
            successCount++;
            break;
        case 'already-memoized':
            console.log(`ℹ️  Already memoized: ${path.basename(relativePath)}`);
            alreadyMemoized++;
            break;
        case 'no-changes':
            console.log(`⚠️  Could not modify: ${path.basename(relativePath)} - may need manual review`);
            errors++;
            break;
        case 'error':
            console.log(`❌ Error processing ${path.basename(relativePath)}: ${result.error}`);
            errors++;
            break;
    }
});

console.log('\n📊 Summary:');
console.log(`   Successfully memoized: ${successCount}`);
console.log(`   Already memoized: ${alreadyMemoized}`);
console.log(`   Errors or manual review needed: ${errors}`);

if (successCount > 0) {
    console.log('\n✨ Components are now wrapped with React.memo for better performance!');
    console.log('📝 Note: Make sure these components receive stable props to benefit from memoization.');
}