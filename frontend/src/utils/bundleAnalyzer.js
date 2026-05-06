/**
 * BUNDLE SIZE ANALYZER
 * 
 * FIX 61: Bundle optimization and dead code elimination
 * Analyzes bundle composition and identifies optimization opportunities
 */
import fs from 'fs';
import path from 'path';
export class BundleAnalyzer {
  constructor(options = {}) {
    this.srcDir = options.srcDir || 'src';
    this.buildDir = options.buildDir || '.next';
    this.excludePatterns = options.excludePatterns || [
      'node_modules',
      '.git',
      '.next',
      'coverage',
      '__tests__',
      '*.test.*',
      '*.spec.*'
    ];
  }
  /**
   * Analyze import statements to identify optimization opportunities
   */
  async analyzeImports() {
    const files = this.getJavaScriptFiles();
    const importAnalysis = {
      reactImports: {},
      lucideImports: {},
      libraryImports: {},
      internalImports: {},
      unusedImports: [],
      duplicateImports: {},
      heavyImports: []
    };
    for (const file of files) {
      try {
        const content = fs.readFileSync(file, 'utf8');
        const imports = this.extractImports(content);
        imports.forEach(importStatement => {
          this.categorizeImport(importStatement, file, importAnalysis);
        });
      } catch (error) {
        // File operation error
        if (process.env.NODE_ENV === 'development') {
            console.error('File operation error:', error);
        }
        return null; // Return null for failed file operations
    }
    }
    return this.generateImportReport(importAnalysis);
  }
  /**
   * Get all JavaScript/TypeScript files
   */
  getJavaScriptFiles() {
    const files = [];
    const scanDirectory = (dir) => {
      try {
        const items = fs.readdirSync(dir);
        for (const item of items) {
          const fullPath = path.join(dir, item);
          const stat = fs.statSync(fullPath);
          if (stat.isDirectory() && !this.shouldExclude(item)) {
            scanDirectory(fullPath);
          } else if (stat.isFile() && this.isJavaScriptFile(item)) {
            files.push(fullPath);
          }
        }
      } catch (error) {
        // File operation error
        if (process.env.NODE_ENV === 'development') {
            console.error('File operation error:', error);
        }
        return null; // Return null for failed file operations
    }
    };
    scanDirectory(this.srcDir);
    return files;
  }
  /**
   * Check if directory/file should be excluded
   */
  shouldExclude(item) {
    return this.excludePatterns.some(pattern => 
      item.includes(pattern) || item.startsWith('.')
    );
  }
  /**
   * Check if file is a JavaScript/TypeScript file
   */
  isJavaScriptFile(filename) {
    return /\.(js|jsx|ts|tsx)$/.test(filename);
  }
  /**
   * Extract import statements from file content
   */
  extractImports(content) {
    const importRegex = /import\s+(?:(?:[\w*\s{},]*)\s+from\s+)?['""]([^'"]*)['""];?/g;
    const imports = [];
    let match;
    while ((match = importRegex.exec(content)) !== null) {
      imports.push({
        full: match[0],
        source: match[1],
        specifiers: this.extractImportSpecifiers(match[0])
      });
    }
    return imports;
  }
  /**
   * Extract import specifiers (what's being imported)
   */
  extractImportSpecifiers(importStatement) {
    // Extract named imports
    const namedMatch = importStatement.match(/import\s*\{([^}]+)\}/);
    if (namedMatch) {
      return namedMatch[1]
        .split(',')
        .map(spec => spec.trim())
        .filter(spec => spec.length > 0);
    }
    // Extract default imports
    const defaultMatch = importStatement.match(/import\s+(\w+)/);
    if (defaultMatch) {
      return [defaultMatch[1]];
    }
    return [];
  }
  /**
   * Categorize imports for analysis
   */
  categorizeImport(importStatement, file, analysis) {
    const { source, specifiers } = importStatement;
    // React imports
    if (source === 'react') {
      specifiers.forEach(spec => {
        if (!analysis.reactImports[spec]) {
          analysis.reactImports[spec] = [];
        }
        analysis.reactImports[spec].push(file);
      });
    }
    // Lucide React imports
    else if (source === 'lucide-react') {
      specifiers.forEach(spec => {
        if (!analysis.lucideImports[spec]) {
          analysis.lucideImports[spec] = [];
        }
        analysis.lucideImports[spec].push(file);
      });
    }
    // Library imports (external packages)
    else if (!source.startsWith('.') && !source.startsWith('/')) {
      if (!analysis.libraryImports[source]) {
        analysis.libraryImports[source] = [];
      }
      analysis.libraryImports[source].push(file);
    }
    // Internal imports
    else {
      if (!analysis.internalImports[source]) {
        analysis.internalImports[source] = [];
      }
      analysis.internalImports[source].push(file);
    }
  }
  /**
   * Generate comprehensive import analysis report
   */
  generateImportReport(analysis) {
    const report = {
      summary: {
        totalReactImports: Object.keys(analysis.reactImports).length,
        totalLucideImports: Object.keys(analysis.lucideImports).length,
        totalLibraryImports: Object.keys(analysis.libraryImports).length,
        totalInternalImports: Object.keys(analysis.internalImports).length
      },
      optimizationOpportunities: {
        // React hooks that could be consolidated
        multipleReactHooks: this.findMultipleReactHookUsage(analysis.reactImports),
        // Unused icon imports
        singleUseLucideIcons: this.findSingleUseLucideIcons(analysis.lucideImports),
        // Heavy libraries with few usages
        heavyLibraries: this.findHeavyLibraries(analysis.libraryImports)
      },
      recommendations: []
    };
    // Generate recommendations
    this.generateRecommendations(report);
    return report;
  }
  /**
   * Find files with multiple React hook imports that could use useReducer
   */
  findMultipleReactHookUsage(reactImports) {
    const hookFiles = {};
    const stateHooks = ['useState', 'useReducer'];
    const effectHooks = ['useEffect', 'useLayoutEffect'];
    const memoHooks = ['useCallback', 'useMemo'];
    Object.entries(reactImports).forEach(([hook, files]) => {
      if ([...stateHooks, ...effectHooks, ...memoHooks].includes(hook)) {
        files.forEach(file => {
          if (!hookFiles[file]) {
            hookFiles[file] = { state: 0, effect: 0, memo: 0, total: 0 };
          }
          if (stateHooks.includes(hook)) hookFiles[file].state++;
          if (effectHooks.includes(hook)) hookFiles[file].effect++;
          if (memoHooks.includes(hook)) hookFiles[file].memo++;
          hookFiles[file].total++;
        });
      }
    });
    // Find files with excessive hook usage
    return Object.entries(hookFiles)
      .filter(([file, counts]) => counts.total > 5 || counts.state > 3)
      .map(([file, counts]) => ({ file, ...counts }));
  }
  /**
   * Find Lucide icons used only once
   */
  findSingleUseLucideIcons(lucideImports) {
    return Object.entries(lucideImports)
      .filter(([icon, files]) => files.length === 1)
      .map(([icon, files]) => ({ icon, file: files[0] }));
  }
  /**
   * Find heavy libraries with minimal usage
   */
  findHeavyLibraries(libraryImports) {
    const heavyLibraries = [
      'lodash', 'moment', 'antd', 'material-ui', 
      'bootstrap', 'jquery', 'three', 'd3'
    ];
    return Object.entries(libraryImports)
      .filter(([lib, files]) => {
        return heavyLibraries.some(heavy => lib.includes(heavy)) && files.length < 3;
      })
      .map(([lib, files]) => ({ library: lib, usage: files.length, files }));
  }
  /**
   * Generate optimization recommendations
   */
  generateRecommendations(report) {
    const { optimizationOpportunities } = report;
    // React hooks optimization
    if (optimizationOpportunities.multipleReactHooks.length > 0) {
      report.recommendations.push({
        type: 'react-hooks',
        priority: 'high',
        title: 'Consolidate React Hooks',
        description: `Found ${optimizationOpportunities.multipleReactHooks.length} files with excessive hook usage. Consider using useReducer for state management.`,
        files: optimizationOpportunities.multipleReactHooks.slice(0, 5),
        estimatedSavings: '10-20KB bundle size, improved re-render performance'
      });
    }
    // Lucide icons optimization
    if (optimizationOpportunities.singleUseLucideIcons.length > 10) {
      report.recommendations.push({
        type: 'lucide-icons',
        priority: 'medium',
        title: 'Optimize Icon Imports',
        description: `Found ${optimizationOpportunities.singleUseLucideIcons.length} icons used only once. Consider creating an icon bundle.`,
        estimatedSavings: '5-15KB bundle size reduction'
      });
    }
    // Heavy libraries
    if (optimizationOpportunities.heavyLibraries.length > 0) {
      report.recommendations.push({
        type: 'heavy-libraries',
        priority: 'high',
        title: 'Replace Heavy Libraries',
        description: 'Found heavy libraries with minimal usage. Consider lighter alternatives.',
        libraries: optimizationOpportunities.heavyLibraries,
        estimatedSavings: '50-200KB bundle size reduction'
      });
    }
  }
  /**
   * Generate comprehensive optimization report
   */
  async generateOptimizationReport() {
    const importAnalysis = await this.analyzeImports();
    const report = {
      timestamp: new Date().toISOString(),
      importAnalysis,
      summary: {
        totalOptimizationOpportunities: importAnalysis.recommendations.length,
        estimatedBundleSavings: this.calculateEstimatedSavings(importAnalysis),
        priorityActions: this.getPriorityActions(importAnalysis)
      }
    };
    return report;
  }
  /**
   * Calculate estimated bundle savings
   */
  calculateEstimatedSavings(importAnalysis) {
    let totalSavings = 0;
    importAnalysis.recommendations.forEach(rec => {
      switch (rec.type) {
        case 'heavy-libraries':
          totalSavings += 100; // KB
          break;
        case 'react-hooks':
          totalSavings += 15; // KB
          break;
        case 'lucide-icons':
          totalSavings += 10; // KB
          break;
      }
    });
    return `${totalSavings}KB estimated bundle size reduction`;
  }
  /**
   * Get priority actions for optimization
   */
  getPriorityActions(importAnalysis) {
    const actions = [];
    // High priority: Heavy libraries
    const heavyLibRecs = importAnalysis.recommendations.filter(r => r.type === 'heavy-libraries');
    if (heavyLibRecs.length > 0) {
      actions.push({
        priority: 1,
        action: 'Replace heavy libraries with lighter alternatives',
        impact: 'High bundle size reduction'
      });
    }
    // Low priority: React hooks consolidation
    const hookRecs = importAnalysis.recommendations.filter(r => r.type === 'react-hooks');
    if (hookRecs.length > 0) {
      actions.push({
        priority: 2,
        action: 'Consolidate React hooks with useReducer',
        impact: 'Better state management and re-render optimization'
      });
    }
    return actions;
  }
  /**
   * Print formatted optimization report
   */
  printReport(report) {
    // Summary
    // Priority Actions
    report.summary.priorityActions.forEach(action => {
      });
    // Import Analysis
    // Recommendations
    report.importAnalysis.recommendations.forEach((rec, i) => {
      });
  }
}
export default BundleAnalyzer; 