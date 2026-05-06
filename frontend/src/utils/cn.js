import { clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';
/**
 * Utility function to merge Tailwind CSS classes efficiently
 * Combines clsx for conditional classes and tailwind-merge for deduplication
 * 
 * @param {...any} inputs - Class names to merge
 * @returns {string} - Merged class string
 */
export function cn(...inputs) {
  return twMerge(clsx(inputs));
}
/**
 * Alternative implementation without external dependencies
 * Use this if you prefer not to add clsx and tailwind-merge
 */
export function cnSimple(...classes) {
  return classes
    .filter(Boolean)
    .join(' ')
    .replace(/\s+/g, ' ')
    .trim();
} 