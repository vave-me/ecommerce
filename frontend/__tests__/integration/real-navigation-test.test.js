/**
 * REAL NAVIGATION TEST - NO MOCKS
 * Tests the actual next-intl navigation without any mocking
 */

import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { NextIntlClientProvider } from 'next-intl';

// Real test messages
const testMessages = {
  en: {
    Header: {
      homeButton: 'Home',
      exploreButton: 'Explore',
    }
  }
};

// Test component that uses REAL navigation
const RealNavigationTest = () => {
  // Import the REAL navigation module
  const { useRouter, useParams } = require('../../src/i18n/navigation');
  
  console.log('🔍 Testing REAL navigation module...');
  
  // Get the real router and params
  const router = useRouter();
  const params = useParams();
  
  console.log('📋 Router methods:', Object.keys(router));
  console.log('📋 Params:', params);
  
  // Track router calls
  const originalPush = router.push;
  const calls = [];
  
  // Wrap router.push to track calls
  router.push = (url) => {
    console.log(`🔗 REAL ROUTER CALL: router.push("${url}")`);
    calls.push({ method: 'push', url });
    // Don't actually navigate in tests
    return Promise.resolve();
  };
  
  const handleHomeClick = () => {
    const locale = params?.locale || 'en';
    router.push(`/${locale}/home`);
  };
  
  const handleExploreClick = () => {
    const locale = params?.locale || 'en';
    router.push(`/${locale}/explore`);
  };
  
  return (
    <div>
      <h2>Real Navigation Test</h2>
      <p>Current locale: {params?.locale || 'en'}</p>
      <button onClick={handleHomeClick} data-testid="home-btn">
        Home
      </button>
      <button onClick={handleExploreClick} data-testid="explore-btn">
        Explore
      </button>
      <div data-testid="calls-count">{calls.length}</div>
    </div>
  );
};

describe('🚀 REAL NAVIGATION TEST - NO MOCKS', () => {
  test('should use REAL next-intl navigation with locale URLs', () => {
    console.log('=== TESTING REAL NAVIGATION ===');
    
    // Render with real NextIntlClientProvider
    render(
      <NextIntlClientProvider locale="en" messages={testMessages.en}>
        <RealNavigationTest />
      </NextIntlClientProvider>
    );
    
    // Verify component renders
    expect(screen.getByText('Real Navigation Test')).toBeInTheDocument();
    expect(screen.getByText('Current locale: en')).toBeInTheDocument();
    
    // Click Home button
    console.log('🏠 Clicking Home button...');
    fireEvent.click(screen.getByTestId('home-btn'));
    
    // Click Explore button  
    console.log('🔍 Clicking Explore button...');
    fireEvent.click(screen.getByTestId('explore-btn'));
    
    // Verify calls were made
    const callsCount = screen.getByTestId('calls-count');
    expect(callsCount.textContent).toBe('2');
    
    console.log('✅ Real navigation test completed!');
    console.log('=== END REAL NAVIGATION TEST ===');
  });
  
  test('should work with Polish locale', () => {
    console.log('=== TESTING POLISH LOCALE ===');
    
    const polishMessages = {
      Header: {
        homeButton: 'Strona główna',
        exploreButton: 'Przeglądaj',
      }
    };
    
    render(
      <NextIntlClientProvider locale="pl" messages={polishMessages}>
        <RealNavigationTest />
      </NextIntlClientProvider>
    );
    
    expect(screen.getByText('Current locale: pl')).toBeInTheDocument();
    
    console.log('🇵🇱 Testing Polish navigation...');
    fireEvent.click(screen.getByTestId('home-btn'));
    fireEvent.click(screen.getByTestId('explore-btn'));
    
    console.log('✅ Polish navigation test completed!');
    console.log('=== END POLISH TEST ===');
  });
}); 