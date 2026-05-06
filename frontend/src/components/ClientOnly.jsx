"use client";
import { useState, useEffect, memo } from 'react';
/**
 * ClientOnly component that only renders its children on the client side.
 * Use this for content that depends on browser-specific APIs or that needs
 * to avoid hydration mismatches between server and client.
 * 
 * OPTIMIZED: Memoized to prevent unnecessary re-renders
 */
const ClientOnly = memo(function ClientOnly({ children }) {
  const [isMounted, setIsMounted] = useState(false);
  useEffect(() => {
    setIsMounted(true);
  }, []);
  if (!isMounted) {
    return null;
  }
  return <>{children}</>;
});
export default ClientOnly;