/**
 * Router Debug Test
 * Debug test to understand why mock router is not being called
 */

import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { jest } from '@jest/globals';

// Test the router mock directly
describe('Router Debug Test', () => {
  test('should verify router mock is working', () => {
    // Import the navigation module
    const navigation = require('../../src/i18n/navigation');
    
    console.log('=== ROUTER DEBUG ===');
    console.log('Navigation module:', Object.keys(navigation));
    
    // Test the useRouter hook directly
    const router = navigation.useRouter();
    console.log('Router from useRouter:', router);
    console.log('Router methods:', Object.keys(router));
    
    // Test calling router.push
    router.push('/test-url');
    console.log('Router push calls:', router.push.mock.calls);
    
    // Test setting a custom router
    const customRouter = {
      push: jest.fn(),
      replace: jest.fn(),
      back: jest.fn(),
      forward: jest.fn(),
      refresh: jest.fn(),
      prefetch: jest.fn(),
      pathname: '/custom'
    };
    
    navigation.__setMockRouter(customRouter);
    
    const newRouter = navigation.useRouter();
    console.log('New router after setting custom:', newRouter);
    console.log('Is same as custom router:', newRouter === customRouter);
    
    newRouter.push('/custom-test');
    console.log('Custom router push calls:', customRouter.push.mock.calls);
    
    console.log('=== END ROUTER DEBUG ===');
  });

  test('should test simple button with router', () => {
    const navigation = require('../../src/i18n/navigation');
    
    // Create a simple test component
    const TestComponent = () => {
      const router = navigation.useRouter();
      
      return (
        <button onClick={() => router.push('/test-route')}>
          Test Button
        </button>
      );
    };
    
    // Set up custom router
    const mockRouter = {
      push: jest.fn(),
      replace: jest.fn(),
      back: jest.fn(),
      forward: jest.fn(),
      refresh: jest.fn(),
      prefetch: jest.fn(),
      pathname: '/'
    };
    
    navigation.__setMockRouter(mockRouter);
    
    render(<TestComponent />);
    
    const button = screen.getByText('Test Button');
    fireEvent.click(button);
    
    console.log('=== SIMPLE BUTTON TEST ===');
    console.log('Mock router push calls:', mockRouter.push.mock.calls);
    console.log('Expected: [["/test-route"]]');
    console.log('=== END SIMPLE BUTTON TEST ===');
    
    expect(mockRouter.push).toHaveBeenCalledWith('/test-route');
  });
}); 