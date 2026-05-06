const fs = require('fs');
const path = require('path');
const { promisify } = require('util');
const readdir = promisify(fs.readdir);
const stat = promisify(fs.stat);

// Patterns that likely need cleanup
const cleanupPatterns = [
    /setTimeout/g,
    /setInterval/g,
    /addEventListener/g,
    /\.subscribe/g,
    /\.on\(/g,
    /observer\.observe/g,
    /new\s+WebSocket/g,
    /new\s+EventSource/g,
    /requestAnimationFrame/g
];

async function* walkDir(dir) {
    const files = await readdir(dir);
    for (const file of files) {
        const filePath = path.join(dir, file);
        const stats = await stat(filePath);
        if (stats.isDirectory()) {
            if (!file.includes('node_modules') && !file.startsWith('.')) {
                yield* walkDir(filePath);
            }
        } else if (file.endsWith('.jsx') || file.endsWith('.js') || file.endsWith('.tsx') || file.endsWith('.ts')) {
            yield filePath;
        }
    }
}

function analyzeUseEffect(content, filePath) {
    const issues = [];
    
    // Match useEffect calls
    const useEffectPattern = /useEffect\s*\(\s*(?:\(\)|async\s*\(\))\s*=>\s*{([^}]+(?:{[^}]*}[^}]*)*)}\s*(?:,\s*\[[^\]]*\])?\s*\)/g;
    const matches = [...content.matchAll(useEffectPattern)];
    
    matches.forEach((match, index) => {
        const effectBody = match[1];
        const hasReturn = /return\s+(?:\(\)|{|function)/.test(effectBody);
        
        // Check if effect contains patterns that need cleanup
        let needsCleanup = false;
        let foundPatterns = [];
        
        cleanupPatterns.forEach(pattern => {
            if (pattern.test(effectBody)) {
                needsCleanup = true;
                foundPatterns.push(pattern.source.replace(/\\/g, ''));
            }
        });
        
        if (needsCleanup && !hasReturn) {
            const line = content.substring(0, match.index).split('\n').length;
            issues.push({
                line,
                patterns: foundPatterns,
                snippet: effectBody.substring(0, 100).replace(/\n/g, ' ').trim() + '...'
            });
        }
    });
    
    return issues;
}

async function main() {
    console.log('🔍 Scanning for useEffect hooks without cleanup...\n');
    
    const srcDir = path.join(__dirname, '..', 'src');
    let totalIssues = 0;
    const allIssues = [];
    
    for await (const filePath of walkDir(srcDir)) {
        try {
            const content = fs.readFileSync(filePath, 'utf8');
            if (!content.includes('useEffect')) continue;
            
            const issues = analyzeUseEffect(content, filePath);
            if (issues.length > 0) {
                const relativePath = path.relative(path.join(__dirname, '..'), filePath);
                console.log(`⚠️  ${relativePath}:`);
                issues.forEach(issue => {
                    console.log(`   Line ${issue.line}: Needs cleanup for: ${issue.patterns.join(', ')}`);
                    console.log(`   Preview: ${issue.snippet}`);
                    allIssues.push({
                        file: relativePath,
                        line: issue.line,
                        patterns: issue.patterns
                    });
                });
                console.log('');
                totalIssues += issues.length;
            }
        } catch (error) {
            // Skip files with errors
        }
    }
    
    console.log(`\n📊 Summary:`);
    console.log(`   Total useEffect hooks needing cleanup: ${totalIssues}`);
    
    if (allIssues.length > 0) {
        console.log('\n📋 Priority files to fix:');
        // Group by file
        const byFile = {};
        allIssues.forEach(issue => {
            if (!byFile[issue.file]) byFile[issue.file] = 0;
            byFile[issue.file]++;
        });
        
        Object.entries(byFile)
            .sort((a, b) => b[1] - a[1])
            .slice(0, 10)
            .forEach(([file, count]) => {
                console.log(`   ${file} (${count} issues)`);
            });
    }
    
    console.log('\n✨ Scan complete!');
}

main().catch(console.error);