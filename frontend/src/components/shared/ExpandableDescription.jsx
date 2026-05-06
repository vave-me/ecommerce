"use client";
import React, { useState, useRef, useEffect, memo } from 'react';
import { ChevronDown, ChevronUp } from '@/icons';
import styles from './ExpandableDescription.module.css';
/**
 * ExpandableDescription - Professional, social media-inspired expandable text component
 * Optimized for readability, mobile experience, and modern design standards
 * 
 * Design Philosophy:
 * - Typography inspired by Twitter/LinkedIn/Facebook standards
 * - Compact, professional appearance with excellent readability
 * - Smooth, fast interactions with subtle animations
 * - Fully accessible with proper ARIA attributes
 * - Mobile-first responsive design
 * 
 * @param {string} text - The description text to display
 * @param {number} maxLines - Maximum lines to show when collapsed (default: 3)
 * @param {string} className - Additional CSS classes
 * @param {string} moreText - Text for expand button (default: "Show more")
 * @param {string} lessText - Text for collapse button (default: "Show less")
 * @param {boolean} compact - Use compact mode with tighter spacing (default: false)
 * @param {string} variant - Typography variant: 'body' | 'small' | 'large' (default: 'body')
 * @param {React.ReactNode} additionalContent - Additional content to show when expanded (e.g., tags)
 */
const ExpandableDescription = memo(({ 
    text, 
    maxLines = 3, 
    className = "", 
    moreText = "Show more",
    lessText = "Show less",
    compact = false,
    variant = "body",
    additionalContent = null
}) => {
    const [isExpanded, setIsExpanded] = useState(false);
    const [needsExpansion, setNeedsExpansion] = useState(false);
    const [isCalculating, setIsCalculating] = useState(true);
    const textRef = useRef(null);
    const measureRef = useRef(null);
    useEffect(() => {
        if (!text || text.trim() === '') {
            // If we have additional content but no text, still show the expansion capability
            if (additionalContent) {
                setNeedsExpansion(true);
                setIsCalculating(false);
            } else {
                setIsCalculating(false);
            }
            return;
        }
        const calculateTruncation = () => {
            if (!textRef.current) return;
            const element = textRef.current;
            const computedStyle = window.getComputedStyle(element);
            const lineHeight = parseFloat(computedStyle.lineHeight);
            const fontSize = parseFloat(computedStyle.fontSize);
            // Calculate actual line height if it's 'normal'
            const actualLineHeight = isNaN(lineHeight) ? fontSize * 1.4 : lineHeight;
            const maxHeight = Math.ceil(actualLineHeight * maxLines);
            // Create a temporary element to measure full height
            const tempElement = element.cloneNode(true);
            tempElement.style.cssText = `
                position: absolute;
                visibility: hidden;
                height: auto;
                max-height: none;
                -webkit-line-clamp: none;
                display: block;
                width: ${element.offsetWidth}px;
                font-family: ${computedStyle.fontFamily};
                font-size: ${computedStyle.fontSize};
                line-height: ${computedStyle.lineHeight};
                font-weight: ${computedStyle.fontWeight};
                letter-spacing: ${computedStyle.letterSpacing};
                word-spacing: ${computedStyle.wordSpacing};
                padding: 0;
                margin: 0;
                border: 0;
            `;
            document.body.appendChild(tempElement);
            const fullHeight = tempElement.scrollHeight;
            document.body.removeChild(tempElement);
            const shouldTruncate = fullHeight > maxHeight + 2; // 2px tolerance
            // Show expansion if text needs truncation OR if we have additional content
            setNeedsExpansion(shouldTruncate || additionalContent);
            setIsCalculating(false);
        };
        // Use ResizeObserver for better performance
        const resizeObserver = new ResizeObserver(() => {
            calculateTruncation();
        });
        if (textRef.current) {
            resizeObserver.observe(textRef.current);
            calculateTruncation();
        }
        return () => {
            resizeObserver.disconnect();
        };
    }, [text, maxLines, additionalContent]);
    const toggleExpansion = () => {
        setIsExpanded(!isExpanded);
    };
    const handleKeyDown = (event) => {
        if (event.key === 'Enter' || event.key === ' ') {
            event.preventDefault();
            toggleExpansion();
        }
    };
    if ((!text || text.trim() === '') && !additionalContent) {
        return null;
    }
    const containerClasses = [
        styles.container,
        compact && styles.compact,
        styles[variant],
        className
    ].filter(Boolean).join(' ');
    const textClasses = [
        styles.text,
        isExpanded ? styles.expanded : styles.collapsed,
        isCalculating && styles.calculating
    ].filter(Boolean).join(' ');
    return (
        <div className={containerClasses}>
            {text && text.trim() !== '' && (
                <div
                    ref={textRef}
                    className={textClasses}
                    style={{
                        WebkitLineClamp: isExpanded ? 'none' : maxLines,
                    }}
                    data-testid="expandable-text"
                >
                    {text}
                </div>
            )}
            {/* Additional content shown only when expanded */}
            {isExpanded && additionalContent && (
                <div className={styles.additionalContent}>
                    {additionalContent}
                </div>
            )}
            {needsExpansion && !isCalculating && (
                <button
                    className={styles.toggleButton}
                    onClick={toggleExpansion}
                    onKeyDown={handleKeyDown}
                    aria-expanded={isExpanded}
                    aria-label={isExpanded ? `Collapse description` : `Expand description`}
                    data-testid="expand-toggle"
                >
                    <span className={styles.buttonText}>
                        {isExpanded ? lessText : moreText}
                    </span>
                    <span className={styles.buttonIcon} aria-hidden="true">
                        {isExpanded ? 
                            <ChevronUp size={12} strokeWidth={2.5} /> : 
                            <ChevronDown size={12} strokeWidth={2.5} />
                        }
                    </span>
                </button>
            )}
        </div>
    );
});
ExpandableDescription.displayName = 'ExpandableDescription';
export default ExpandableDescription; 