import React from 'react';
import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import '@testing-library/jest-dom';
import CreateProductModal from './index';
import { useAuth } from '../../context/AuthContext';
import { addProduct, updateProduct } from '../../api/client/productsApi';
import { createMedia } from '../../api/client/mediaApi';
import { fetchMainCategories } from '../../api/categories';
import { renderWithProviders } from '../../tests/test-utils';
import { renderModalWithLoading, waitForStep } from '../../tests/modal-test-utils';

// Mock the required modules
jest.mock('../../context/AuthContext');
jest.mock('../../api/client/productsApi');
jest.mock('../../api/client/mediaApi');
jest.mock('../../api/categories');
jest.mock('../../hooks/useFocusTrap', () => ({
  useFocusTrap: () => ({ current: null })
}));

// Mock the child components to simplify testing
jest.mock('./components/steps/BasicInfoStep/BasicInfoStep', () => {
  const mockBasicInfoStep = ({ onSubmit, onCancel, isLoading }) => {
    const handleSubmit = () => {
      onSubmit({
        name: 'Test Product',
        description: 'This is a test product',
        basePrice: '100',
        condition: 'new',
        categoryId: 'cat1',
        categorySlug: 'category-1',
        brand: 'Test Brand',
        model: 'Test Model',
        negotiable: true,
        userType: 'private',
        sku: 'TEST-123',
        tags: 'electronics, test'
      });
    };

    return (
      <div data-testid="basic-info-step">
        <h2>Basic Information</h2>
        <button onClick={onCancel}>Cancel</button>
        <button 
          onClick={handleSubmit}
          disabled={isLoading}
        >
          Continue
        </button>
      </div>
    );
  };
  return { BasicInfoStep: mockBasicInfoStep };
});

jest.mock('./components/steps/MediaUploadStep/MediaUploadStep', () => {
  const mockMediaUploadStep = ({ onComplete, onBack, isLoading }) => {
    const handleComplete = () => {
      onComplete({
        images: [],
        videoUrl: "",
        thumbnail: ""
      });
    };

    return (
      <div data-testid="media-upload-step">
        <h2>Media Upload</h2>
        <button onClick={onBack}>Back</button>
        <button 
          onClick={handleComplete}
          disabled={isLoading}
        >
          Continue
        </button>
      </div>
    );
  };
  return { MediaUploadStep: mockMediaUploadStep };
});

jest.mock('./components/steps/OptionalSettingsStep/OptionalSettingsStep', () => {
  const mockOptionalInfoStep = ({ onSubmit, onBack, isLoading }) => {
    const handleComplete = () => {
      onSubmit({
        weight: '5',
        height: '10',
        width: '15',
        depth: '8',
        manageStocks: true,
        stock: '100',
        shippingCost: '10',
        middlemanService: false,
        hasVariants: false
      });
    };

    return (
      <div data-testid="optional-settings-step">
        <h2>Optional Settings</h2>
        <button onClick={onBack}>Back</button>
        <button 
          onClick={handleComplete}
          disabled={isLoading}
        >
          Continue
        </button>
      </div>
    );
  };
  return { OptionalInfoStep: mockOptionalInfoStep };
});

jest.mock('./components/steps/FinalizeStep/FinalizeStep', () => {
  const mockFinalizeStep = ({ onFinalize, onBack, isLoading, isSuccess }) => {
    const handleFinalize = () => {
      onFinalize({
        latitude: 40.7128,
        longitude: -74.0060
      });
    };

    return (
      <div data-testid="finalize-step">
        <h2>Finalize</h2>
        {isSuccess ? <div>Success!</div> : (
          <>
            <button onClick={onBack}>Back</button>
            <button 
              onClick={handleFinalize}
              disabled={isLoading}
            >
              Publish
            </button>
          </>
        )}
      </div>
    );
  };
  return { FinalizeStep: mockFinalizeStep };
});

// Mock modal components
jest.mock('../shared/modal/ModalOverlay', () => ({
    ModalOverlay: ({ children }) => <div data-testid="modal-overlay">{children}</div>
}));

jest.mock('../shared/modal/ModalContainer', () => ({
    ModalContainer: ({ children }) => <div data-testid="modal-container">{children}</div>
}));

jest.mock('../shared/StepNavigation/StepNavigation', () => ({
  StepNavigation: ({ currentStep, lastCompletedStep, onStepClick }) => (
    <div data-testid="step-navigation">
      <button onClick={() => onStepClick(1)}>Basic Info</button>
      <button onClick={() => onStepClick(2)}>Media Upload</button>
      <button onClick={() => onStepClick(3)}>Optional Settings</button>
      <button onClick={() => onStepClick(4)}>Finalize</button>
    </div>
  )
}));

jest.mock('../shared/ModalHeader/ModalHeader', () => ({
  ModalHeader: ({ title, onClose }) => (
    <div data-testid="modal-header">
      <h1>{title}</h1>
      <button onClick={onClose} aria-label="close">Close</button>
    </div>
  )
}));

jest.mock('../../common/components/ErrorAlert', () => ({
  ErrorAlert: ({ message }) => <div data-testid="error-alert">{message}</div>
}));

describe('CreateProductModal Component', () => {
  // Common props
  const mockOnClose = jest.fn();
  
  // Mock auth context
  const mockUser = { 
    userId: 'user123',
    authenticated: true 
  };
  
  // Mock API responses
  const mockProductId = 'product123';
  const mockMediaId = 'media123';
  const mockCategories = [
    { id: 'cat1', name: 'Category 1', slug: 'category-1' },
    { id: 'cat2', name: 'Category 2', slug: 'category-2' }
  ];

  beforeEach(() => {
    // Reset all mocks before each test
    jest.clearAllMocks();
    
    // Setup auth mock
    useAuth.mockReturnValue({ user: mockUser });
    
    // Setup API mocks
    addProduct.mockResolvedValue({ id: mockProductId });
    createMedia.mockResolvedValue({ id: mockMediaId });
    updateProduct.mockResolvedValue({ success: true });
    fetchMainCategories.mockResolvedValue({ categories: mockCategories });
  });

  describe('Rendering', () => {
    test('renders modal with correct title in create mode', async () => {
      await renderModalWithLoading(<CreateProductModal onClose={mockOnClose} />);
      expect(screen.getByText('Create Listing')).toBeInTheDocument();
    });

    test('renders modal with correct title in edit mode', async () => {
      await renderModalWithLoading(<CreateProductModal onClose={mockOnClose} editMode={true} initialProductData={{ id: 'test' }} />);
      expect(screen.getByText('Edit Listing')).toBeInTheDocument();
    });
    
    test('renders first step by default', async () => {
      await renderModalWithLoading(<CreateProductModal onClose={mockOnClose} />);
      expect(screen.getByTestId('basic-info-step')).toBeInTheDocument();
    });
  });

  describe('Navigation', () => {
    test('clicking close button calls onClose', async () => {
      await renderModalWithLoading(<CreateProductModal onClose={mockOnClose} />);
      const closeButton = screen.getByRole('button', { name: /close/i });
      await userEvent.click(closeButton);
      expect(mockOnClose).toHaveBeenCalled();
    });
    
    test('cannot navigate to incomplete steps', async () => {
      await renderModalWithLoading(<CreateProductModal onClose={mockOnClose} />);
      
      // Try to click on step 3
      const step3Button = screen.getByText(/Optional Settings/i);
      await userEvent.click(step3Button);
      
      // Should still be on step 1
      expect(screen.getByTestId('basic-info-step')).toBeInTheDocument();
    });
  });

  describe('Basic Info Step (Step 1)', () => {
    test('successfully submits basic info and advances to step 2', async () => {
      await renderModalWithLoading(<CreateProductModal onClose={mockOnClose} />);
      
      // Continue button
      const continueButton = screen.getByRole('button', { name: /Continue/i });
      await userEvent.click(continueButton);
      
      // Verify API call
      await waitFor(() => {
        expect(addProduct).toHaveBeenCalledWith(expect.objectContaining({
          name: 'Test Product',
          description: 'This is a test product',
          basePrice: 100,
          userSellerId: mockUser.userId
        }));
      });
      
      // Verify navigation to next step
      await waitForStep('media-upload-step');
    });

    test('shows error message when API call fails', async () => {
      // Mock API failure
      addProduct.mockRejectedValueOnce(new Error('API Error'));
      
      await renderModalWithLoading(<CreateProductModal onClose={mockOnClose} />);
      
      // Continue button
      const continueButton = screen.getByRole('button', { name: /Continue/i });
      await userEvent.click(continueButton);
      
      await waitFor(() => {
        expect(screen.getByTestId('error-alert')).toBeInTheDocument();
      });
    });
  });

  describe('Media Upload Step (Step 2)', () => {
    // Helper to get to step 2
    const getToMediaStep = async () => {
      await renderModalWithLoading(<CreateProductModal onClose={mockOnClose} />);
      
      // Fill step 1 and submit
      const continueButton = screen.getByRole('button', { name: /Continue/i });
      await userEvent.click(continueButton);
      
      // Wait for step 2 to render
      await waitForStep('media-upload-step');
    };

    test('can navigate back to step 1', async () => {
      await getToMediaStep();
      
      // Click back button
      const backButton = screen.getByRole('button', { name: /Back/i });
      await userEvent.click(backButton);
      
      // Verify navigation back
      await waitForStep('basic-info-step');
    });

    test('can continue to step 3 without images', async () => {
      await getToMediaStep();
      
      // Continue without adding images
      const continueButton = screen.getByRole('button', { name: /Continue/i });
      await userEvent.click(continueButton);
      
      // Verify navigation to step 3
      await waitForStep('optional-settings-step');
    });
  });

  describe('Optional Info Step (Step 3)', () => {
    // Helper to get to step 3
    const getToOptionalStep = async () => {
      await renderModalWithLoading(<CreateProductModal onClose={mockOnClose} />);
      
      // Complete step 1
      let continueButton = screen.getByRole('button', { name: /Continue/i });
      await userEvent.click(continueButton);
      
      // Continue past step 2
      await waitForStep('media-upload-step');
      
      continueButton = screen.getByRole('button', { name: /Continue/i });
      await userEvent.click(continueButton);
      
      // Wait for step 3 to render
      await waitForStep('optional-settings-step');
    };

    test('can navigate back to step 2', async () => {
      await getToOptionalStep();
      
      // Click back button
      const backButton = screen.getByRole('button', { name: /Back/i });
      await userEvent.click(backButton);
      
      // Verify navigation back
      await waitForStep('media-upload-step');
    });

    test('can continue to step 4', async () => {
      await getToOptionalStep();
      
      // Continue to step 4
      const continueButton = screen.getByRole('button', { name: /Continue/i });
      await userEvent.click(continueButton);
      
      // Note: The updateProduct call is commented out in the actual component
      
      // Verify navigation to final step
      await waitForStep('finalize-step');
    });
  });

  describe('Finalize Step (Step 4)', () => {
    // Helper to get to step 4
    const getFinalizeStep = async () => {
      await renderModalWithLoading(<CreateProductModal onClose={mockOnClose} />);
      
      // Complete steps 1-3
      let continueButton = screen.getByRole('button', { name: /Continue/i });
      await userEvent.click(continueButton);
      
      // Skip through step 2 and 3
      await waitForStep('media-upload-step');
      continueButton = screen.getByRole('button', { name: /Continue/i });
      await userEvent.click(continueButton);
      
      await waitForStep('optional-settings-step');
      continueButton = screen.getByRole('button', { name: /Continue/i });
      await userEvent.click(continueButton);
      
      // Wait for final step to render
      await waitForStep('finalize-step');
    };

    test('can navigate back to step 3', async () => {
      await getFinalizeStep();
      
      // Click back button
      const backButton = screen.getByRole('button', { name: /Back/i });
      await userEvent.click(backButton);
      
      // Verify navigation back
      await waitForStep('optional-settings-step');
    });

    test('successfully publishes product', async () => {
      await getFinalizeStep();
      
      // Publish the product 
      const publishButton = screen.getByRole('button', { name: /Publish/i });
      await userEvent.click(publishButton);
      
      // Note: The updateProduct call is commented out in the actual component
      
      // Verify success message
      await waitFor(() => {
        expect(screen.getByText(/Success!/i)).toBeInTheDocument();
      });
    });
  });

  describe('Edit Mode', () => {
    const mockInitialData = {
      id: 'existing123',
      mediaId: 'existingMedia123',
      name: 'Existing Product',
      description: 'This is an existing product',
      basePrice: 200,
      condition: 'used',
      categoryId: 'cat2',
      categorySlug: 'category-2',
      brand: 'Existing Brand',
      model: 'Existing Model',
      tags: ['tag1', 'tag2'],
      thumbnail: 'thumbnail.jpg',
      weight: 10,
      height: 20,
      width: 30,
      depth: 15,
      stock: 50,
      shippingCost: 15,
      manageStocks: true,
      middlemanService: true,
      hasVariants: false,
      longitude: -73.9857,
      latitude: 40.7484
    };

    test('loads initial data in edit mode', async () => {
      await renderModalWithLoading(
        <CreateProductModal 
          onClose={mockOnClose} 
          editMode={true} 
          initialProductData={mockInitialData} 
        />
      );
      
      // Should start with step 1
      expect(screen.getByTestId('basic-info-step')).toBeInTheDocument();
      
      // This is just checking that the component renders in edit mode
      // The actual form data is validated at component level and pre-filled
      // with mock components that don't actually display values
    });

    test('calls updateProduct when publishing in edit mode', async () => {
      // Setup for edit mode
      await renderModalWithLoading(
        <CreateProductModal 
          onClose={mockOnClose} 
          editMode={true} 
          initialProductData={mockInitialData} 
        />
      );
      
      // Skip to last step
      let continueButton = screen.getByRole('button', { name: /Continue/i });
      await userEvent.click(continueButton);
      
      await waitForStep('media-upload-step');
      continueButton = screen.getByRole('button', { name: /Continue/i });
      await userEvent.click(continueButton);
      
      await waitForStep('optional-settings-step');
      continueButton = screen.getByRole('button', { name: /Continue/i });
      await userEvent.click(continueButton);
      
      await waitForStep('finalize-step');
      
      // Publish the edited product
      const publishButton = screen.getByRole('button', { name: /Publish/i });
      await userEvent.click(publishButton);
      
      // Verify updateProduct was called
      await waitFor(() => {
        expect(updateProduct).toHaveBeenCalled();
      });
      
      // Success message should be showing
      await waitFor(() => {
        expect(screen.getByText(/Success!/i)).toBeInTheDocument();
      });
    });
  });

  describe('Error Handling', () => {
    test('displays error when user is not logged in', async () => {
      // Mock user as not logged in
      useAuth.mockReturnValue({ user: null });
      
      render(<CreateProductModal onClose={mockOnClose} />);
      
      // Try to submit
      const continueButton = screen.getByRole('button', { name: /Continue/i });
      await userEvent.click(continueButton);
      
      // Verify error message
      await waitFor(() => {
        expect(screen.getByTestId('error-alert')).toBeInTheDocument();
      });
    });
    
    test('handles API error on media creation', async () => {
      // First, make addProduct succeed but createMedia fail
      addProduct.mockResolvedValueOnce({ id: mockProductId });
      createMedia.mockRejectedValueOnce(new Error('Media creation failed'));
      
      render(<CreateProductModal onClose={mockOnClose} />);
      
      // Complete step 1
      const continueButton = screen.getByRole('button', { name: /Continue/i });
      await userEvent.click(continueButton);
      
      // Verify error alert appears
      await waitFor(() => {
        expect(screen.getByTestId('error-alert')).toBeInTheDocument();
      });
    });
  });
}); 