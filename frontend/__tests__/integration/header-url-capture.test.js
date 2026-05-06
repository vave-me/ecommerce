/**
 * Header URL Capture Test
 * Simple test to capture what URLs the Header component actually calls
 */

import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { jest } from '@jest/globals';
import Header from '../../src/components/Header/Header';
import { renderWithRealProviders } from '../utils/test-setup';

describe('Header URL Capture Test', () => {
  test('should capture actual URLs called by Header navigation buttons', async () => {
    const { mockRouter } = renderWithRealProviders(<Header />, {
      locale: 'en',
      isMobile: false,
      showNavbars: true,
      isClient: true
    });

    console.log('=== TESTING HEADER NAVIGATION URLs ===');

    // Test Home button
    try {
      const homeButton = screen.getByText('Home');
      console.log('✅ Found Home button with text "Home"');
      
      // Clear any previous calls
      mockRouter.push.mockClear();
      
      fireEvent.click(homeButton);
      
      // Give it a moment for the click to process
      await new Promise(resolve => setTimeout(resolve, 100));
      
      console.log('🔗 Home button router calls:', mockRouter.push.mock.calls);
      console.log('🔗 Total calls made:', mockRouter.push.mock.calls.length);
      
      if (mockRouter.push.mock.calls.length > 0) {
        console.log('🎯 ACTUAL URL CALLED:', mockRouter.push.mock.calls[0][0]);
      } else {
        console.log('❌ No router calls made for Home button');
      }
    } catch (error) {
      console.log('❌ Home button not found with text "Home"');
      console.log('Available buttons:', screen.getAllByRole('button').map(btn => btn.textContent || btn.getAttribute('aria-label')));
    }

    // Test Explore button
    try {
      const exploreButton = screen.getByText('Explore');
      console.log('✅ Found Explore button with text "Explore"');
      
      // Clear previous calls
      mockRouter.push.mockClear();
      
      fireEvent.click(exploreButton);
      
      // Give it a moment for the click to process
      await new Promise(resolve => setTimeout(resolve, 100));
      
      console.log('🔗 Explore button router calls:', mockRouter.push.mock.calls);
      
      if (mockRouter.push.mock.calls.length > 0) {
        console.log('🎯 ACTUAL URL CALLED:', mockRouter.push.mock.calls[0][0]);
      } else {
        console.log('❌ No router calls made for Explore button');
      }
    } catch (error) {
      console.log('❌ Explore button not found with text "Explore"');
    }

    console.log('=== END URL CAPTURE TEST ===');
  });

  test('should test with Polish locale', async () => {
    const { mockRouter } = renderWithRealProviders(<Header />, {
      locale: 'pl',
      isMobile: false,
      showNavbars: true,
      isClient: true
    });

    console.log('=== TESTING POLISH LOCALE URLs ===');

    // Test Home button in Polish
    try {
      const homeButton = screen.getByText('Strona główna');
      console.log('✅ Found Polish Home button with text "Strona główna"');
      
      // Clear any previous calls
      mockRouter.push.mockClear();
      
      fireEvent.click(homeButton);
      
      // Give it a moment for the click to process
      await new Promise(resolve => setTimeout(resolve, 100));
      
      console.log('🔗 Polish Home button router calls:', mockRouter.push.mock.calls);
      
      if (mockRouter.push.mock.calls.length > 0) {
        console.log('🎯 ACTUAL POLISH URL CALLED:', mockRouter.push.mock.calls[0][0]);
      } else {
        console.log('❌ No router calls made for Polish Home button');
      }
    } catch (error) {
      console.log('❌ Polish Home button not found');
      console.log('Available buttons:', screen.getAllByRole('button').map(btn => btn.textContent || btn.getAttribute('aria-label')));
    }

    console.log('=== END POLISH LOCALE TEST ===');
  });
}); 