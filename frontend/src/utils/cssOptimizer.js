/**
 * CSS OPTIMIZER UTILITY - PHASE 7
 * Automatically replaces heavy CSS effects with optimized alternatives
 * 
 * FIX 95: CSS Performance Optimizer
 */
/**
 * Replace heavy gradients with optimized solid colors
 */
export const optimizeGradients = (cssString) => {
  const gradientReplacements = {
    'linear-gradient(135deg, #e6f3ff 0%, #cce7ff 100%)': '#e6f3ff',
    'linear-gradient(45deg, #0066cc, #0052a3)': '#0066cc',
    'linear-gradient(135deg, #fed7d7 0%, #feb2b2 100%)': '#fed7d7',
    'linear-gradient(135deg, #c6f6d5 0%, #9ae6b4 100%)': '#c6f6d5',
    'linear-gradient(135deg, #fef3c7, #fed7aa)': '#fef3c7',
    'linear-gradient(90deg, #10b981, #f59e0b, #dc2626)': '#10b981',
    'linear-gradient(135deg, #dc2626, #b91c1c)': '#dc2626',
    'linear-gradient(135deg, #b91c1c, #991b1b)': '#b91c1c',
    'linear-gradient(135deg, #6b7280, #4b5563)': '#6b7280',
    'linear-gradient(135deg, #4b5563, #374151)': '#4b5563',
    'linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%)': '#8b5cf6',
    'linear-gradient(135deg, #22c55e 0%, #16a34a 100%)': '#22c55e',
    'linear-gradient(135deg, #ef4444 0%, #dc2626 100%)': '#ef4444',
    'linear-gradient(135deg, #f1f5f9 0%, #e2e8f0 100%)': '#f1f5f9',
    'linear-gradient(135deg, #2980b9 0%, #6366f1 100%)': '#2980b9',
    'linear-gradient(145deg, #f8fafc 0%, #f1f5f9 100%)': '#f8fafc'
  };
  let optimizedCSS = cssString;
  Object.entries(gradientReplacements).forEach(([gradient, solidColor]) => {
    const regex = new RegExp(gradient.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'gi');
    optimizedCSS = optimizedCSS.replace(regex, solidColor);
  });
  return optimizedCSS;
};
/**
 * Remove backdrop-filter effects for better performance
 */
export const removeBackdropFilters = (cssString) => {
  const backdropFilters = [
    'backdrop-filter: blur(2px);',
    'backdrop-filter: blur(3px);',
    'backdrop-filter: blur(4px);',
    'backdrop-filter: blur(5px);',
    'backdrop-filter: blur(10px);',
    'backdrop-filter: blur(12px);',
    'backdrop-filter: blur(16px);',
    'backdrop-filter: blur(20px);'
  ];
  let optimizedCSS = cssString;
  backdropFilters.forEach(filter => {
    const regex = new RegExp(filter.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'gi');
    optimizedCSS = optimizedCSS.replace(regex, '/* backdrop-filter removed for performance */');
  });
  return optimizedCSS;
};
/**
 * Replace heavy box-shadows with optimized versions
 */
export const optimizeBoxShadows = (cssString) => {
  const shadowReplacements = {
    'box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25), 0 10px 10px -5px rgba(0, 0, 0, 0.04);': 'box-shadow: 0 20px 25px rgba(0, 0, 0, 0.15);',
    'box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);': 'box-shadow: 0 10px 15px rgba(0, 0, 0, 0.1);',
    'box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);': 'box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);'
  };
  let optimizedCSS = cssString;
  Object.entries(shadowReplacements).forEach(([heavyShadow, lightShadow]) => {
    const regex = new RegExp(heavyShadow.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'gi');
    optimizedCSS = optimizedCSS.replace(regex, lightShadow);
  });
  return optimizedCSS;
};
/**
 * Replace transform scale animations with opacity-based ones
 */
export const optimizeTransforms = (cssString) => {
  const transformReplacements = {
    'transform: scale(1.05);': 'opacity: 0.9;',
    'transform: scale(1.1);': 'opacity: 0.85;',
    'transform: scale(1.02);': 'opacity: 0.95;',
    'transform: scale(0.95);': 'opacity: 0.9; transform: translateY(2px);',
    'transform: scale(0.98);': 'opacity: 0.95; transform: translateY(1px);'
  };
  let optimizedCSS = cssString;
  Object.entries(transformReplacements).forEach(([heavyTransform, lightTransform]) => {
    const regex = new RegExp(heavyTransform.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'gi');
    optimizedCSS = optimizedCSS.replace(regex, lightTransform);
  });
  return optimizedCSS;
};
/**
 * Optimize transition properties to be more specific
 */
export const optimizeTransitions = (cssString) => {
  const transitionReplacements = {
    'transition: all 0.2s ease;': 'transition: background-color 0.2s ease, border-color 0.2s ease, opacity 0.2s ease;',
    'transition: all 0.3s ease;': 'transition: background-color 0.3s ease, border-color 0.3s ease, opacity 0.3s ease;',
    'transition: all 0.15s ease;': 'transition: background-color 0.15s ease, border-color 0.15s ease, opacity 0.15s ease;'
  };
  let optimizedCSS = cssString;
  Object.entries(transitionReplacements).forEach(([genericTransition, specificTransition]) => {
    const regex = new RegExp(genericTransition.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'gi');
    optimizedCSS = optimizedCSS.replace(regex, specificTransition);
  });
  return optimizedCSS;
};
/**
 * Add will-change property only where needed
 */
export const optimizeWillChange = (cssString) => {
  // Remove unnecessary will-change properties
  let optimizedCSS = cssString.replace(/will-change:\s*auto;?/gi, '');
  // Add will-change for elements that actually animate
  const animatedSelectors = [
    '.optimized-button:hover',
    '.optimized-card:hover',
    '.gpu-accelerated'
  ];
  animatedSelectors.forEach(selector => {
    const regex = new RegExp(`(${selector.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')}\\s*{[^}]*)(})`, 'gi');
    optimizedCSS = optimizedCSS.replace(regex, (match, beforeClosing, closing) => {
      if (!beforeClosing.includes('will-change')) {
        return beforeClosing + 'will-change: transform;' + closing;
      }
      return match;
    });
  });
  return optimizedCSS;
};
/**
 * Complete CSS optimization pipeline
 */
export const optimizeCSS = (cssString) => {
  let optimized = cssString;
  // Apply all optimizations
  optimized = optimizeGradients(optimized);
  optimized = removeBackdropFilters(optimized);
  optimized = optimizeBoxShadows(optimized);
  optimized = optimizeTransforms(optimized);
  optimized = optimizeTransitions(optimized);
  optimized = optimizeWillChange(optimized);
  return optimized;
};
/**
 * CSS Performance Analysis
 */
export const analyzeCSS = (cssString) => {
  const analysis = {
    gradients: (cssString.match(/linear-gradient|radial-gradient/gi) || []).length,
    backdropFilters: (cssString.match(/backdrop-filter/gi) || []).length,
    heavyShadows: (cssString.match(/box-shadow:[^;]*rgba\([^)]*\)[^;]*rgba/gi) || []).length,
    transforms: (cssString.match(/transform:\s*scale/gi) || []).length,
    genericTransitions: (cssString.match(/transition:\s*all/gi) || []).length,
    totalSize: cssString.length
  };
  analysis.performanceScore = Math.max(0, 100 - (
    analysis.gradients * 2 +
    analysis.backdropFilters * 5 +
    analysis.heavyShadows * 3 +
    analysis.transforms * 1 +
    analysis.genericTransitions * 1
  ));
  return analysis;
};
/**
 * Generate optimization report
 */
export const generateOptimizationReport = (originalCSS, optimizedCSS) => {
  const originalAnalysis = analyzeCSS(originalCSS);
  const optimizedAnalysis = analyzeCSS(optimizedCSS);
  const report = {
    original: originalAnalysis,
    optimized: optimizedAnalysis,
    improvements: {
      gradientsReduced: originalAnalysis.gradients - optimizedAnalysis.gradients,
      backdropFiltersRemoved: originalAnalysis.backdropFilters - optimizedAnalysis.backdropFilters,
      heavyShadowsReduced: originalAnalysis.heavyShadows - optimizedAnalysis.heavyShadows,
      transformsOptimized: originalAnalysis.transforms - optimizedAnalysis.transforms,
      transitionsImproved: originalAnalysis.genericTransitions - optimizedAnalysis.genericTransitions,
      sizeReduced: originalAnalysis.totalSize - optimizedAnalysis.totalSize,
      performanceGain: optimizedAnalysis.performanceScore - originalAnalysis.performanceScore
    }
  };
  return report;
};
export default {
  optimizeCSS,
  analyzeCSS,
  generateOptimizationReport,
  optimizeGradients,
  removeBackdropFilters,
  optimizeBoxShadows,
  optimizeTransforms,
  optimizeTransitions,
  optimizeWillChange
}; 