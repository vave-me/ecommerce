import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { NavBarProvider, NavBarContext } from '../../context/NavBarContext.jsx';

// Test component that uses the NavBar context
const TestComponent = () => {
  const { 
    showNavbars, 
    toggleNavbars, 
    isMobile
  } = React.useContext(NavBarContext);

  return (
    <div>
      <div data-testid="navbar-status">{showNavbars ? 'Open' : 'Closed'}</div>
      <div data-testid="mobile-status">{isMobile ? 'Mobile' : 'Desktop'}</div>
      
      <button data-testid="toggle-navbar" onClick={toggleNavbars}>Toggle NavBar</button>
    </div>
  );
};

describe('NavBarContext', () => {
  it('initializes with closed state on desktop', () => {
    // Mock window.innerWidth to simulate desktop
    global.innerWidth = 1024;
    global.dispatchEvent(new Event('resize'));

    render(
      <NavBarProvider>
        <TestComponent />
      </NavBarProvider>
    );

    // Initial states should be closed on desktop
    expect(screen.getByTestId('navbar-status')).toHaveTextContent('Closed');
    expect(screen.getByTestId('mobile-status')).toHaveTextContent('Desktop');
  });

  it('initializes with open state on mobile', () => {
    // Mock window.innerWidth to simulate mobile
    global.innerWidth = 480;
    global.dispatchEvent(new Event('resize'));

    render(
      <NavBarProvider>
        <TestComponent />
      </NavBarProvider>
    );

    // Initial states should be open on mobile
    expect(screen.getByTestId('navbar-status')).toHaveTextContent('Open');
    expect(screen.getByTestId('mobile-status')).toHaveTextContent('Mobile');
  });

  it('toggles navbar when toggleNavbars is called', () => {
    // Mock window.innerWidth to simulate desktop (closed by default)
    global.innerWidth = 1024;
    global.dispatchEvent(new Event('resize'));

    render(
      <NavBarProvider>
        <TestComponent />
      </NavBarProvider>
    );

    // Initially closed on desktop
    expect(screen.getByTestId('navbar-status')).toHaveTextContent('Closed');

    // Toggle once - should open
    fireEvent.click(screen.getByTestId('toggle-navbar'));
    expect(screen.getByTestId('navbar-status')).toHaveTextContent('Open');

    // Toggle again - should close
    fireEvent.click(screen.getByTestId('toggle-navbar'));
    expect(screen.getByTestId('navbar-status')).toHaveTextContent('Closed');
  });
}); 