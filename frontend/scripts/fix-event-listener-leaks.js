const fs = require('fs');
const path = require('path');

// Find files that need event listener cleanup
const filesToAnalyze = [
    './src/utils/mobilePerformance.js',
    './src/utils/hydrationSafeUtils.js',
    './src/utils/hookOptimizations.js',
    './src/hooks/useMobileDetection.jsx',
    './src/hooks/useCleanup.jsx',
    './src/context/ThemeContext.jsx',
    './src/components/classified/ClassifiedCard.jsx',
    './src/components/Shared/MobileForm.jsx',
    './src/components/PerformanceMonitor/PerformanceMonitor.client.jsx',
    './src/components/PWA/PWAInitializer.jsx',
    './src/components/Header/MobileNavMenu.jsx',
    './src/components/Header/Header.jsx',
    './src/components/Header/AddOptionsSheetWithComposer.jsx',
    './src/components/Header/AddOptionsSheet.jsx',
    './src/app/[locale]/cart/CartPageWithOffline.jsx',
    './src/components/Filters/HorizontalFilters.jsx',
    './src/app/[locale]/admin/orders/OrdersManagement.client.jsx',
    './src/components/Modal/UnifiedModal.jsx'
];

let totalFixed = 0;
let totalAlreadyFixed = 0;
let issues = [];

function analyzeFile(filePath) {
    try {
        const content = fs.readFileSync(filePath, 'utf8');
        
        // Find addEventListener calls
        const addListenerPattern = /(\w+)\.addEventListener\s*\(\s*['"](\w+)['"]\s*,\s*([^)]+)\)/g;
        const matches = [...content.matchAll(addListenerPattern)];
        
        if (matches.length === 0) {
            return { hasListeners: false };
        }
        
        // Check if there's corresponding removeEventListener
        let unhandledListeners = [];
        
        matches.forEach(match => {
            const [fullMatch, target, event, handler] = match;
            const handlerName = handler.trim().replace(/[,\s].*/g, '');
            
            // Look for corresponding removeEventListener
            const removePattern = new RegExp(`${target}\\.removeEventListener\\s*\\(\\s*['"]${event}['"]\\s*,\\s*${handlerName.replace(/[.*+?^${}()|[\]\\]/g, '\\\\$&')}`, 'g');
            
            if (!removePattern.test(content)) {
                unhandledListeners.push({
                    target,
                    event,
                    handler: handlerName,
                    line: content.substring(0, match.index).split('\n').length
                });
            }
        });
        
        // Check if it's in a useEffect with cleanup
        const useEffectPattern = /useEffect\s*\(\s*\(\)\s*=>\s*{([^}]+(?:{[^}]*}[^}]*)*)}\s*,/g;
        const effectMatches = [...content.matchAll(useEffectPattern)];
        
        let properlyHandled = [];
        let needsFixing = [];
        
        unhandledListeners.forEach(listener => {
            let isInEffectWithCleanup = false;
            
            effectMatches.forEach(effectMatch => {
                const effectBody = effectMatch[1];
                if (effectBody.includes(`addEventListener`) && 
                    effectBody.includes(`${listener.event}`) &&
                    effectBody.includes('return')) {
                    isInEffectWithCleanup = true;
                }
            });
            
            if (isInEffectWithCleanup) {
                properlyHandled.push(listener);
            } else {
                needsFixing.push(listener);
            }
        });
        
        return {
            hasListeners: true,
            total: matches.length,
            unhandled: unhandledListeners.length,
            needsFixing: needsFixing,
            properlyHandled: properlyHandled
        };
        
    } catch (error) {
        console.error(`❌ Error analyzing ${filePath}:`, error.message);
        return { error: true };
    }
}

console.log('🔍 Analyzing event listener usage...\n');

filesToAnalyze.forEach(filePath => {
    const fullPath = path.join(__dirname, '..', filePath);
    if (fs.existsSync(fullPath)) {
        const result = analyzeFile(fullPath);
        
        if (result.error) {
            console.log(`❌ Error analyzing ${filePath}`);
        } else if (!result.hasListeners) {
            console.log(`ℹ️  No event listeners found in ${filePath}`);
        } else {
            const emoji = result.needsFixing.length === 0 ? '✅' : '⚠️';
            console.log(`${emoji} ${filePath}:`);
            console.log(`   Total listeners: ${result.total}`);
            console.log(`   Unhandled: ${result.unhandled}`);
            console.log(`   Needs fixing: ${result.needsFixing.length}`);
            console.log(`   Properly handled: ${result.properlyHandled.length}`);
            
            if (result.needsFixing.length > 0) {
                totalFixed += result.needsFixing.length;
                result.needsFixing.forEach(listener => {
                    console.log(`   ⚠️  Line ${listener.line}: ${listener.target}.addEventListener('${listener.event}', ${listener.handler})`);
                    issues.push({
                        file: filePath,
                        line: listener.line,
                        listener
                    });
                });
            } else {
                totalAlreadyFixed += result.properlyHandled.length;
            }
            console.log('');
        }
    } else {
        console.log(`⚠️  File not found: ${filePath}`);
    }
});

console.log('\n📊 Summary:');
console.log(`   Total event listeners needing cleanup: ${totalFixed}`);
console.log(`   Total already properly handled: ${totalAlreadyFixed}`);

if (issues.length > 0) {
    console.log('\n📋 Issues to fix:');
    issues.forEach(issue => {
        console.log(`   ${issue.file}:${issue.line} - ${issue.listener.target}.addEventListener('${issue.listener.event}')`);
    });
}

console.log('\n✨ Analysis complete!');