"use client";
import React, { memo } from 'react';
/**
 * Utility for composing providers in an efficient way
 * Avoids deep nesting in JSX and makes provider composition more flexible
 * 
 * @param {Object} props - Component props
 * @param {Array} props.providers - Array of provider components
 * @param {React.ReactNode} props.children - Child elements
 * @returns {React.ReactElement} - Composed providers with children
 */
export const ProvidersComposer = ({ providers = [], children }) => {
  return providers.reduceRight(
    (accumulator, Provider) => <Provider>{accumulator}</Provider>,
    children
  );
};
export default memo(ProvidersComposer); 