import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { ErrorBoundary } from '../../src/components/ErrorBoundary';
import { errorHandler } from '../../src/utils/globalErrorHandler';

// Mock error handler
jest.mock('../../src/utils/globalErrorHandler', () => ({
    errorHandler: {
        handleError: jest.fn()
    }
}));

// Component that throws an error
const ThrowError = ({ shouldThrow = false }) => {
    if (shouldThrow) {
        throw new Error('Test error');
    }
    return <div>No error</div>;
};

// Component that throws in useEffect
const ThrowErrorInEffect = () => {
    React.useEffect(() => {
        throw new Error('Effect error');
    }, []);
    return <div>Component with effect</div>;
};

describe('ErrorBoundary', () => {
    // Suppress console errors during tests
    const originalError = console.error;
    beforeAll(() => {
        console.error = jest.fn();
    });

    afterAll(() => {
        console.error = originalError;
    });

    beforeEach(() => {
        jest.clearAllMocks();
    });

    describe('error catching', () => {
        it('should render children when no error occurs', () => {
            render(
                <ErrorBoundary>
                    <div>Test content</div>
                </ErrorBoundary>
            );

            expect(screen.getByText('Test content')).toBeInTheDocument();
        });

        it('should catch errors and display fallback UI', () => {
            render(
                <ErrorBoundary>
                    <ThrowError shouldThrow={true} />
                </ErrorBoundary>
            );

            expect(screen.getByText(/Something went wrong/i)).toBeInTheDocument();
            expect(screen.queryByText('No error')).not.toBeInTheDocument();
        });

        it('should use custom fallback when provided', () => {
            const CustomFallback = () => <div>Custom error UI</div>;

            render(
                <ErrorBoundary fallback={<CustomFallback />}>
                    <ThrowError shouldThrow={true} />
                </ErrorBoundary>
            );

            expect(screen.getByText('Custom error UI')).toBeInTheDocument();
        });

        it('should use function fallback with error details', () => {
            const fallbackFn = jest.fn((error, errorInfo, reset) => (
                <div>
                    <p>Error: {error.message}</p>
                    <button onClick={reset}>Reset</button>
                </div>
            ));

            render(
                <ErrorBoundary fallback={fallbackFn}>
                    <ThrowError shouldThrow={true} />
                </ErrorBoundary>
            );

            expect(screen.getByText('Error: Test error')).toBeInTheDocument();
            expect(fallbackFn).toHaveBeenCalledWith(
                expect.any(Error),
                expect.objectContaining({ componentStack: expect.any(String) }),
                expect.any(Function)
            );
        });
    });

    describe('error handling', () => {
        it('should call errorHandler.handleError', () => {
            render(
                <ErrorBoundary name="TestBoundary">
                    <ThrowError shouldThrow={true} />
                </ErrorBoundary>
            );

            expect(errorHandler.handleError).toHaveBeenCalledWith(
                expect.any(Error),
                expect.objectContaining({
                    context: 'React Error Boundary',
                    metadata: expect.objectContaining({
                        componentStack: expect.any(String),
                        errorBoundary: 'TestBoundary'
                    })
                })
            );
        });

        it('should generate unique error ID', () => {
            render(
                <ErrorBoundary>
                    <ThrowError shouldThrow={true} />
                </ErrorBoundary>
            );

            const errorId = screen.getByText(/Error ID:/i).textContent;
            expect(errorId).toMatch(/Error ID: ERR_\d+_[a-z0-9]+/);
        });

        it('should show reload button when showReload is true', () => {
            const mockReload = jest.fn();
            Object.defineProperty(window, 'location', {
                value: { reload: mockReload },
                writable: true
            });

            render(
                <ErrorBoundary showReload={true}>
                    <ThrowError shouldThrow={true} />
                </ErrorBoundary>
            );

            const reloadButton = screen.getByText('Reload Page');
            expect(reloadButton).toBeInTheDocument();

            userEvent.click(reloadButton);
            expect(mockReload).toHaveBeenCalled();
        });

        it('should not show reload button when showReload is false', () => {
            render(
                <ErrorBoundary showReload={false}>
                    <ThrowError shouldThrow={true} />
                </ErrorBoundary>
            );

            expect(screen.queryByText('Reload Page')).not.toBeInTheDocument();
        });
    });

    describe('reset functionality', () => {
        it('should reset error state when reset is called', async () => {
            let shouldThrow = true;
            const { rerender } = render(
                <ErrorBoundary>
                    <ThrowError shouldThrow={shouldThrow} />
                </ErrorBoundary>
            );

            expect(screen.getByText(/Something went wrong/i)).toBeInTheDocument();

            // Fix the error condition
            shouldThrow = false;

            // Click reset
            const resetButton = screen.getByText('Try Again');
            await userEvent.click(resetButton);

            // Re-render with fixed component
            rerender(
                <ErrorBoundary>
                    <ThrowError shouldThrow={shouldThrow} />
                </ErrorBoundary>
            );

            expect(screen.getByText('No error')).toBeInTheDocument();
            expect(screen.queryByText(/Something went wrong/i)).not.toBeInTheDocument();
        });

        it('should call onReset callback when provided', async () => {
            const onReset = jest.fn();

            render(
                <ErrorBoundary onReset={onReset}>
                    <ThrowError shouldThrow={true} />
                </ErrorBoundary>
            );

            const resetButton = screen.getByText('Try Again');
            await userEvent.click(resetButton);

            expect(onReset).toHaveBeenCalled();
        });
    });

    describe('error details', () => {
        it('should show error details in development', () => {
            const originalEnv = process.env.NODE_ENV;
            process.env.NODE_ENV = 'development';

            render(
                <ErrorBoundary>
                    <ThrowError shouldThrow={true} />
                </ErrorBoundary>
            );

            const details = screen.getByText(/Show Details/i);
            expect(details).toBeInTheDocument();

            userEvent.click(details);
            expect(screen.getByText(/Test error/)).toBeInTheDocument();

            process.env.NODE_ENV = originalEnv;
        });

        it('should not show error details in production', () => {
            const originalEnv = process.env.NODE_ENV;
            process.env.NODE_ENV = 'production';

            render(
                <ErrorBoundary>
                    <ThrowError shouldThrow={true} />
                </ErrorBoundary>
            );

            expect(screen.queryByText(/Show Details/i)).not.toBeInTheDocument();

            process.env.NODE_ENV = originalEnv;
        });
    });

    describe('error boundary nesting', () => {
        it('should catch errors in nested components', () => {
            render(
                <ErrorBoundary name="Outer">
                    <div>
                        <ErrorBoundary name="Inner">
                            <ThrowError shouldThrow={true} />
                        </ErrorBoundary>
                    </div>
                </ErrorBoundary>
            );

            // Inner boundary should catch the error
            expect(errorHandler.handleError).toHaveBeenCalledWith(
                expect.any(Error),
                expect.objectContaining({
                    metadata: expect.objectContaining({
                        errorBoundary: 'Inner'
                    })
                })
            );
        });

        it('should propagate to parent boundary if child boundary throws', () => {
            const BrokenBoundary = ({ children }) => {
                throw new Error('Boundary error');
            };

            render(
                <ErrorBoundary name="Parent">
                    <BrokenBoundary>
                        <div>Content</div>
                    </BrokenBoundary>
                </ErrorBoundary>
            );

            expect(screen.getByText(/Something went wrong/i)).toBeInTheDocument();
        });
    });

    describe('withErrorBoundary HOC', () => {
        it('should wrap component with error boundary', () => {
            const { withErrorBoundary } = require('../../src/components/ErrorBoundary');
            
            const TestComponent = () => <div>Test Component</div>;
            const WrappedComponent = withErrorBoundary(TestComponent, {
                name: 'TestHOC'
            });

            render(<WrappedComponent />);
            expect(screen.getByText('Test Component')).toBeInTheDocument();
        });

        it('should pass props to wrapped component', () => {
            const { withErrorBoundary } = require('../../src/components/ErrorBoundary');
            
            const TestComponent = ({ message }) => <div>{message}</div>;
            const WrappedComponent = withErrorBoundary(TestComponent);

            render(<WrappedComponent message="Hello World" />);
            expect(screen.getByText('Hello World')).toBeInTheDocument();
        });
    });
});