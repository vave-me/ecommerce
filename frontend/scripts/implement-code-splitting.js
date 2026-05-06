const fs = require('fs');
const path = require('path');

// Heavy components that should be lazy loaded
const componentsToSplit = [
    {
        path: 'src/components/NewsletterEditor/TemplateEditor.jsx',
        importPath: '../NewsletterEditor/TemplateEditor',
        componentName: 'TemplateEditor'
    },
    {
        path: 'src/components/ListingManagement/OffersModalFormEditable.jsx',
        importPath: '../ListingManagement/OffersModalFormEditable',
        componentName: 'OffersModalFormEditable'
    },
    {
        path: 'src/components/AI/PromptInput.jsx',
        importPath: '../AI/PromptInput',
        componentName: 'PromptInput'
    },
    {
        path: 'src/components/classified/DetailedProductView.jsx',
        importPath: '../classified/DetailedProductView',
        componentName: 'DetailedProductView'
    },
    {
        path: 'src/components/services/DetailedServiceView.jsx',
        importPath: '../services/DetailedServiceView',
        componentName: 'DetailedServiceView'
    }
];

// Create lazy component wrappers
function createLazyWrapper(componentName, importPath) {
    return `import React, { lazy, Suspense } from 'react';

const ${componentName}Lazy = lazy(() => import('${importPath}'));

export default function ${componentName}(props) {
    return (
        <Suspense fallback={
            <div style={{ 
                padding: '2rem', 
                textAlign: 'center',
                minHeight: '200px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center'
            }}>
                <div className="loading-spinner">Loading...</div>
            </div>
        }>
            <${componentName}Lazy {...props} />
        </Suspense>
    );
}
`;
}

// Find files that import these components
function findImportingFiles(componentName, componentPath) {
    const srcDir = path.join(__dirname, '..', 'src');
    const importingFiles = [];
    
    function searchDir(dir) {
        const files = fs.readdirSync(dir);
        files.forEach(file => {
            const filePath = path.join(dir, file);
            const stat = fs.statSync(filePath);
            
            if (stat.isDirectory() && !file.includes('node_modules') && !file.startsWith('.')) {
                searchDir(filePath);
            } else if ((file.endsWith('.jsx') || file.endsWith('.js')) && !filePath.includes(componentPath)) {
                try {
                    const content = fs.readFileSync(filePath, 'utf8');
                    if (content.includes(componentName) && 
                        (content.includes(`from '`) || content.includes(`from "`))) {
                        importingFiles.push(filePath);
                    }
                } catch (err) {
                    // Skip files that can't be read
                }
            }
        });
    }
    
    searchDir(srcDir);
    return importingFiles;
}

console.log('🚀 Implementing code splitting for heavy components...\n');

componentsToSplit.forEach(({ path: componentPath, importPath, componentName }) => {
    const fullPath = path.join(__dirname, '..', componentPath);
    const wrapperPath = fullPath.replace('.jsx', '.lazy.jsx');
    
    console.log(`📦 Processing ${componentName}...`);
    
    // Check if component exists
    if (!fs.existsSync(fullPath)) {
        console.log(`  ⚠️  Component not found: ${componentPath}`);
        return;
    }
    
    // Create lazy wrapper
    const wrapperContent = createLazyWrapper(componentName, `./${path.basename(componentPath)}`);
    fs.writeFileSync(wrapperPath, wrapperContent, 'utf8');
    console.log(`  ✅ Created lazy wrapper: ${path.basename(wrapperPath)}`);
    
    // Find files that import this component
    const importingFiles = findImportingFiles(componentName, componentPath);
    console.log(`  📍 Found ${importingFiles.length} files importing ${componentName}`);
    
    if (importingFiles.length > 0) {
        console.log(`  📝 Files to update manually:`);
        importingFiles.slice(0, 5).forEach(file => {
            const relativePath = path.relative(path.join(__dirname, '..'), file);
            console.log(`     - ${relativePath}`);
        });
        if (importingFiles.length > 5) {
            console.log(`     ... and ${importingFiles.length - 5} more`);
        }
    }
    
    console.log('');
});

console.log('📝 Next steps:');
console.log('1. Update imports in the files listed above to use the .lazy versions');
console.log('2. Replace direct imports with the lazy wrappers');
console.log('3. Test that components load properly with suspense boundaries');
console.log('\nExample import change:');
console.log("  // Before:");
console.log("  import DetailedProductView from '../classified/DetailedProductView';");
console.log("  // After:");
console.log("  import DetailedProductView from '../classified/DetailedProductView.lazy';");

console.log('\n✨ Code splitting setup complete!');