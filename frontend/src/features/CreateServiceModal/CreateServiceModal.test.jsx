import React from 'react';
import { screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import '@testing-library/jest-dom';
import CreateServiceModal from './index';
import { useAuth } from '../../context/AuthContext';
import { addService, updateService } from '../../api/client/servicesApi';
import { createMedia } from '../../api/client/mediaApi';
import { fetchMainCategories } from '../../api/categories';
import { renderWithProviders } from '../../tests/test-utils';
import { renderModalWithLoading, waitForStep } from '../../tests/modal-test-utils';

// Mock the required modules
jest.mock('../../context/AuthContext');
jest.mock('../../api/client/servicesApi');
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
        name: 'Test Service',
        description: 'This is a test service',
        basePrice: '100',
        categoryId: 'cat1',
        categorySlug: 'category-1'
      });
    };

    return (
      <div data-testid="basic-info-step">
        <h2>Basic Information</h2>
        <label htmlFor="name">Service Name</label>
        <input id="name" data-testid="name-input" />
        <label htmlFor="description">Description</label>
        <input id="description" data-testid="description-input" />
        <label htmlFor="basePrice">Base Price</label>
        <input id="basePrice" data-testid="price-input" />
        <label htmlFor="category">Category</label>
        <select id="category" data-testid="category-select">
          <option value="cat1">Category 1</option>
          <option value="cat2">Category 2</option>
        </select>
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
  const mockOptionalInfoStep = ({ onComplete, onBack, isLoading }) => {
    const handleComplete = () => {
      onComplete({
        shippingCost: '10'
      });
    };

    return (
      <div data-testid="optional-settings-step">
        <h2>Optional Settings</h2>
        <label htmlFor="shippingCost">Shipping Cost</label>
        <input id="shippingCost" data-testid="shipping-cost-input" />
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
        lat: 40.7128,
        lng: -74.0060
      });
    };

    return (
      <div data-testid="finalize-step">
        <h2>Finalize</h2>
        {isSuccess ? <div>Success!</div> : (
          <>
            <label htmlFor="location">Location</label>
            <input id="location" data-testid="location-input" />
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

describe('CreateServiceModal Component', () => {
  // Common props
  const mockOnClose = jest.fn();
  
  // Mock auth context
  const mockUser = { 
    userId: 'user123',
    authenticated: true 
  };
  
  // Mock API responses
  const mockServiceId = 'service123';
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
    addService.mockResolvedValue({ id: mockServiceId });
    updateService.mockResolvedValue({ id: mockServiceId });
    createMedia.mockResolvedValue({ id: mockMediaId });
    fetchMainCategories.mockResolvedValue({ categories: mockCategories });
  });

  test('renders the modal with step 1 initially', async () => {
    await renderModalWithLoading(<CreateServiceModal onClose={mockOnClose} />);
    
    expect(screen.getByTestId('modal-container')).toBeInTheDocument();
    expect(screen.getByTestId('basic-info-step')).toBeInTheDocument();
  });

  // Helper functions for multistep tests
  const getToMediaStep = async () => {
    await renderModalWithLoading(<CreateServiceModal onClose={mockOnClose} />);
  
    // Submit the basic info form
    const continueBtn = screen.getByText('Continue');
    fireEvent.click(continueBtn);
    
    // Should now be on media upload step
    await waitForStep('media-upload-step');
    
    return { continueBtn: screen.getByText('Continue') };
  };
  
  test('Basic Info Step (Step 1) > can navigate to step 2', async () => {
    await renderModalWithLoading(<CreateServiceModal onClose={mockOnClose} />);
    
    // Submit the basic info form
    const continueBtn = screen.getByText('Continue');
    fireEvent.click(continueBtn);
    
    // Should now be on media upload step
    await waitForStep('media-upload-step');
  });

  test('Media Upload Step (Step 2) > can navigate back to step 1', async () => {
    await getToMediaStep();
    
    // Go back to step 1
    const backBtn = screen.getByText('Back');
    fireEvent.click(backBtn);
    
    // Should now be back on basic info step
    await waitForStep('basic-info-step');
  });

  test('renders modal with correct title in create mode', async () => {
    await renderModalWithLoading(<CreateServiceModal onClose={mockOnClose} />);
    expect(screen.getByText('Create Service')).toBeInTheDocument();
  });

  test('renders modal with correct title in edit mode', async () => {
    await renderModalWithLoading(
      <CreateServiceModal onClose={mockOnClose} editMode={true} initialServiceData={{ id: 'test' }} />
    );
    expect(screen.getByText('Edit Service')).toBeInTheDocument();
  });
  
  test('clicking close button calls onClose', async () => {
    await renderModalWithLoading(<CreateServiceModal onClose={mockOnClose} />);
    const closeButton = screen.getByRole('button', { name: /close/i });
    await userEvent.click(closeButton);
    expect(mockOnClose).toHaveBeenCalled();
  });
  
  test('cannot navigate to incomplete steps', async () => {
    await renderModalWithLoading(<CreateServiceModal onClose={mockOnClose} />);
    
    // Try to click on step 3
    const step3Button = screen.getByText(/Optional Settings/i);
    await userEvent.click(step3Button);
    
    // Should still be on step 1
    expect(screen.getByTestId('basic-info-step')).toBeInTheDocument();
  });

  test('successfully submits basic info and advances to step 2', async () => {
    await renderModalWithLoading(<CreateServiceModal onClose={mockOnClose} />);
    
    // Continue button
    const continueButton = screen.getByRole('button', { name: /Continue/i });
    await userEvent.click(continueButton);
    
    // Verify API call
    await waitFor(() => {
      expect(addService).toHaveBeenCalledWith(expect.objectContaining({
        name: 'Test Service',
        description: 'This is a test service',
        basePrice: 100,
        userId: mockUser.userId
      }));
    });
    
    // Verify navigation to next step
    await waitForStep('media-upload-step');
  });
}); 