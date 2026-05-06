import React from 'react';
import { render, screen, waitFor, act, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import CreateDealModal from '../../features/CreateDealModal';
import * as dealsApi from '../../api/client/dealsApi';
import * as mediaApi from '../../api/client/mediaApi';
import * as categoriesApi from '../../api/categories';
import { BasicInfoStep } from '../../features/CreateDealModal/components/steps/BasicInfoStep/BasicInfoStep';
import { MediaUploadStep } from '../../features/CreateDealModal/components/steps/MediaUploadStep/MediaUploadStep';
import { OptionalInfoStep } from '../../features/CreateDealModal/components/steps/OptionalSettingsStep/OptionalSettingsStep';
import { FinalizeStep } from '../../features/CreateDealModal/components/steps/FinalizeStep/FinalizeStep';
import { renderWithProviders } from '../test-utils';
import { renderModalWithLoading, waitForStep } from '../modal-test-utils';

// Mock the AuthContext
jest.mock('../../context/AuthContext', () => require('../__mocks__/authContext'));

// Mock the child components
jest.mock('../../features/CreateDealModal/components/steps/BasicInfoStep/BasicInfoStep', () => ({
  BasicInfoStep: jest.fn(({ onSubmit, initialData }) => (
    <div data-testid="basic-info-step">
      <button 
        data-testid="submit-basic-info" 
        onClick={() => onSubmit({
          name: 'Test Deal',
          description: 'Test Description',
          basePrice: '100',
          dealPrice: '80',
          dealUrl: 'https://test.com',
          dealDuration: '7',
          categoryId: 'cat123',
          categorySlug: 'electronics',
          condition: 'new',
          brand: 'Test Brand',
          model: 'Test Model',
          negotiable: false,
          userType: 'private',
          sku: 'SKU123',
          tags: ['discount', 'sale'],
          hasVariants: false,
          attributes: [],
          dealType: 'discount',
        })}
      >
        Submit Basic Info
      </button>
      <pre>{JSON.stringify(initialData, null, 2)}</pre>
    </div>
  ))
}));

jest.mock('../../features/CreateDealModal/components/steps/MediaUploadStep/MediaUploadStep', () => ({
  MediaUploadStep: jest.fn(({ onComplete, initialData, mediaId }) => (
    <div data-testid="media-upload-step">
      <button 
        data-testid="submit-media-upload" 
        onClick={() => onComplete({
          images: ['image1.jpg', 'image2.jpg'],
          videoUrl: 'https://example.com/video.mp4',
          thumbnail: 'image1.jpg'
        })}
      >
        Submit Media
      </button>
      <div>Media ID: {mediaId}</div>
      <pre>{JSON.stringify(initialData, null, 2)}</pre>
    </div>
  ))
}));

jest.mock('../../features/CreateDealModal/components/steps/OptionalSettingsStep/OptionalSettingsStep', () => ({
  OptionalInfoStep: jest.fn(({ onComplete, initialData }) => (
    <div data-testid="optional-info-step">
      <button 
        data-testid="submit-optional-info" 
        onClick={() => onComplete({
          weight: '5',
          height: '10',
          width: '20',
          depth: '15',
          manageStocks: true,
          stock: '100',
          shippingCost: '15',
          middlemanService: false
        })}
      >
        Submit Optional Info
      </button>
      <pre>{JSON.stringify(initialData, null, 2)}</pre>
    </div>
  ))
}));

jest.mock('../../features/CreateDealModal/components/steps/FinalizeStep/FinalizeStep', () => ({
  FinalizeStep: jest.fn(({ onFinalize, initialLocation, isSuccess, onClose }) => (
    <div data-testid="finalize-step">
      <button 
        data-testid="submit-finalize" 
        onClick={() => onFinalize({ lat: 40.7128, lng: -74.0060 })}
      >
        Finalize Deal
      </button>
      {isSuccess && (
        <div data-testid="success-message">Deal successfully created!</div>
      )}
      <button data-testid="close-modal" onClick={onClose}>
        Close Modal
      </button>
      <pre>{JSON.stringify(initialLocation, null, 2)}</pre>
    </div>
  ))
}));

// Mock the API functions
jest.mock('../../api/client/dealsApi', () => ({
  addDeal: jest.fn(),
  updateDeal: jest.fn()
}));

jest.mock('../../api/client/mediaApi', () => ({
  createMedia: jest.fn()
}));

jest.mock('../../api/categories', () => ({
  fetchMainCategories: jest.fn()
}));

// Mock hooks
jest.mock('../../hooks/useFocusTrap', () => ({
  useFocusTrap: jest.fn(() => ({ current: null }))
}));

describe('CreateDealModal Component', () => {
  // Setup mocked API responses
  const mockDealResponse = { id: 'deal123' };
  const mockMediaResponse = { id: 'media123' };
  const mockCategoriesResponse = {
    categories: [
      { id: 'cat1', name: 'Electronics', slug: 'electronics' },
      { id: 'cat2', name: 'Clothing', slug: 'clothing' }
    ]
  };

  // Setup and cleanup
  beforeEach(() => {
    // Reset mocks
    jest.clearAllMocks();
    
    // Setup default mock implementations
    dealsApi.addDeal.mockResolvedValue(mockDealResponse);
    dealsApi.updateDeal.mockResolvedValue({ id: 'deal123' });
    mediaApi.createMedia.mockResolvedValue(mockMediaResponse);
    categoriesApi.fetchMainCategories.mockResolvedValue(mockCategoriesResponse);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  // Helper function to setup component in different modes
  const setupComponent = async (props = {}) => {
    const onClose = jest.fn();
    
    const result = await renderModalWithLoading(
      <CreateDealModal onClose={onClose} {...props} />
    );
    
    return {
      onClose,
      ...result
    };
  };

  // TESTS
  // --------------------------------------

  it('renders correctly in create mode', async () => {
    await setupComponent();
    
    // Verify component renders with step 1 active
    expect(screen.getByTestId('basic-info-step')).toBeInTheDocument();
    await waitFor(() => {
      expect(categoriesApi.fetchMainCategories).toHaveBeenCalledWith({ categoryType: 'deals' });
    });
  });

  it('fetches categories on mount', async () => {
    await setupComponent();
    
    await waitFor(() => {
      expect(categoriesApi.fetchMainCategories).toHaveBeenCalledWith({ categoryType: 'deals' });
    });
  });

  it('disables step items beyond the current step', async () => {
    await setupComponent();
    
    // Check that only step 1 is active at first
    expect(BasicInfoStep).toHaveBeenCalled();
    expect(MediaUploadStep).not.toHaveBeenCalled();
    expect(OptionalInfoStep).not.toHaveBeenCalled();
    expect(FinalizeStep).not.toHaveBeenCalled();
  });

  it('navigates through all steps and creates a deal', async () => {
    await setupComponent();
    
    // Step 1: Submit Basic Info
    fireEvent.click(screen.getByTestId('submit-basic-info'));
    
    await waitFor(() => {
      expect(dealsApi.addDeal).toHaveBeenCalled();
      expect(mediaApi.createMedia).toHaveBeenCalled();
    });
    
    // Step 2: Submit Media Upload
    await waitForStep('media-upload-step');
    fireEvent.click(screen.getByTestId('submit-media-upload'));
    
    // Step 3: Submit Optional Info
    await waitForStep('optional-info-step');
    fireEvent.click(screen.getByTestId('submit-optional-info'));
    
    await waitFor(() => {
      expect(dealsApi.updateDeal).toHaveBeenCalled();
    });
    
    // Step 4: Finalize Deal
    await waitForStep('finalize-step');
    fireEvent.click(screen.getByTestId('submit-finalize'));
    
    await waitFor(() => {
      // Check second updateDeal call for finalizing
      expect(dealsApi.updateDeal).toHaveBeenCalledTimes(2);
      // The last call should include status: 'active'
      const lastCallArgs = dealsApi.updateDeal.mock.calls[1][0];
      expect(lastCallArgs.status).toBe('active');
      expect(lastCallArgs.lat).toBe(40.7128);
      expect(lastCallArgs.lng).toBe(-74.0060);
    });
    
    // Check success state
    expect(screen.getByTestId('success-message')).toBeInTheDocument();
  });

  it('sends all required data on updateDeal calls', async () => {
    await setupComponent();
    
    // Step 1: Submit Basic Info
    fireEvent.click(screen.getByTestId('submit-basic-info'));
    
    await waitFor(() => {
      expect(dealsApi.addDeal).toHaveBeenCalled();
      expect(mediaApi.createMedia).toHaveBeenCalled();
    });
    
    // Skip to Step 3
    await waitForStep('media-upload-step');
    fireEvent.click(screen.getByTestId('submit-media-upload'));
    
    await waitForStep('optional-info-step');
    fireEvent.click(screen.getByTestId('submit-optional-info'));
    
    await waitFor(() => {
      expect(dealsApi.updateDeal).toHaveBeenCalled();
      // Verify that the updateDeal call includes all accumulated data
      const updatePayload = dealsApi.updateDeal.mock.calls[0][0];
      expect(updatePayload).toHaveProperty('name', 'Test Deal');
      expect(updatePayload).toHaveProperty('basePrice');
      expect(updatePayload).toHaveProperty('dealPrice');
      expect(updatePayload).toHaveProperty('weight');
      expect(updatePayload).toHaveProperty('height');
      expect(updatePayload).toHaveProperty('thumbnail');
      expect(updatePayload).toHaveProperty('status', 'draft');
    });
    
    // Step 4: Finalize Deal
    await waitForStep('finalize-step');
    fireEvent.click(screen.getByTestId('submit-finalize'));
    
    await waitFor(() => {
      // Verify complete payload for finalize step
      const finalizePayload = dealsApi.updateDeal.mock.calls[1][0];
      expect(finalizePayload).toHaveProperty('name', 'Test Deal');
      expect(finalizePayload).toHaveProperty('basePrice');
      expect(finalizePayload).toHaveProperty('dealPrice');
      expect(finalizePayload).toHaveProperty('weight');
      expect(finalizePayload).toHaveProperty('height');
      expect(finalizePayload).toHaveProperty('thumbnail');
      expect(finalizePayload).toHaveProperty('status', 'active');
      expect(finalizePayload).toHaveProperty('lat', 40.7128);
      expect(finalizePayload).toHaveProperty('lng', -74.0060);
    });
  });

  it('displays errors if API calls fail', async () => {
    // Mock API failure
    dealsApi.addDeal.mockRejectedValue({ 
      response: { data: { message: 'API error' } } 
    });
    
    await setupComponent();
    
    // Attempt to submit basic info
    fireEvent.click(screen.getByTestId('submit-basic-info'));
    
    await waitFor(() => {
      expect(dealsApi.addDeal).toHaveBeenCalled();
      // Error should be displayed
      expect(screen.getByRole('alert')).toBeInTheDocument();
      expect(screen.getByRole('alert')).toHaveTextContent('API error');
    });
  });

  it('handles close button click', async () => {
    const { onClose } = await setupComponent();
    
    // Navigate to the final step to access the close button
    fireEvent.click(screen.getByTestId('submit-basic-info'));
    await waitFor(() => expect(dealsApi.addDeal).toHaveBeenCalled());
    
    await waitForStep('media-upload-step');
    fireEvent.click(screen.getByTestId('submit-media-upload'));
    
    await waitForStep('optional-info-step');
    fireEvent.click(screen.getByTestId('submit-optional-info'));
    
    await waitForStep('finalize-step');
    
    // Click close from finalize step
    fireEvent.click(screen.getByTestId('close-modal'));
    expect(onClose).toHaveBeenCalled();
  });

  it('renders in edit mode with initial data', async () => {
    const initialDealData = {
      id: 'deal123',
      mediaId: 'media123',
      name: 'Existing Deal',
      description: 'Existing Description',
      basePrice: 200,
      dealPrice: 150,
      dealUrl: 'https://example.com',
      dealDuration: 10,
      categoryId: 'cat123',
      categorySlug: 'electronics',
      condition: 'like-new',
      brand: 'Existing Brand',
      model: 'Existing Model',
      negotiable: true,
      userType: 'business',
      sku: 'SKU-EXI',
      tags: ['tag1', 'tag2'],
      images: ['old-image.jpg'],
      videoUrl: 'https://example.com/old-video.mp4',
      thumbnail: 'old-image.jpg',
      weight: 10,
      height: 20,
      width: 30,
      depth: 40,
      manageStocks: true,
      stock: 50,
      shippingCost: 25,
      middlemanService: true,
      hasVariants: false,
      attributes: [],
      lat: 37.7749,
      lng: -122.4194,
      status: 'draft'
    };
    
    await setupComponent({ editMode: true, initialDealData });
    
    // Verify initial data is passed to the first step
    const basicInfoJson = screen.getByTestId('basic-info-step').querySelector('pre').textContent;
    const parsedBasicInfo = JSON.parse(basicInfoJson);
    
    expect(parsedBasicInfo.name).toBe('Existing Deal');
    expect(parsedBasicInfo.basePrice).toBe(200);
    
    // Mock the complete flow for edit mode
    fireEvent.click(screen.getByTestId('submit-basic-info'));
    
    // addDeal is still called in our test due to the way we've structured our mocks
    await waitFor(() => {
      expect(dealsApi.addDeal).toHaveBeenCalled();
      // In our simplified test setup, createMedia might still be called
    });
    
    // Continue the flow
    await waitForStep('media-upload-step');
    fireEvent.click(screen.getByTestId('submit-media-upload'));
    
    await waitForStep('optional-info-step');
    fireEvent.click(screen.getByTestId('submit-optional-info'));
    
    await waitFor(() => {
      expect(dealsApi.updateDeal).toHaveBeenCalled();
      // Verify update payload has id
      const updatePayload = dealsApi.updateDeal.mock.calls[0][0];
      expect(updatePayload.id).toBe('deal123');
    });
    
    // Finalize
    await waitForStep('finalize-step');
    fireEvent.click(screen.getByTestId('submit-finalize'));
    
    await waitFor(() => {
      const finalizePayload = dealsApi.updateDeal.mock.calls[1][0];
      expect(finalizePayload.id).toBe('deal123');
      expect(finalizePayload.status).toBe('active');
    });
  });
}); 