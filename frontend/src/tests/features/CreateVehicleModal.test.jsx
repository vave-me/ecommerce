import React from 'react';
import { render, screen, waitFor, act, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import CreateVehicleModal from '../../features/CreateVehicleModal';
import * as vehiclesApi from '../../api/client/vehiclesApi';
import * as mediaApi from '../../api/client/mediaApi';
import * as categoriesApi from '../../api/categories';

// Mock the AuthContext
jest.mock('../../context/AuthContext', () => ({
  useAuth: () => ({
    user: { userId: 'user123' },
    isUserLoggedIn: true
  })
}));

// Mock the child components
jest.mock('../../features/CreateVehicleModal/components/steps/BasicInfoStep/BasicInfoStep', () => ({
  __esModule: true,
  BasicInfoStep: jest.fn(({ onSubmit, initialData }) => (
    <div data-testid="basic-info-step">
      <button 
        data-testid="submit-basic-info" 
        onClick={() => onSubmit({
          name: 'Test Vehicle',
          description: 'Test Description',
          basePrice: '25000',
          categoryId: 'cat123',
          categorySlug: 'automotive',
          condition: 'new',
          brand: 'Toyota',
          model: 'Camry',
          performance: '200',
          fuelType: 'gasoline',
          transmissionType: 'automatic',
          numberOfOwners: '1',
          accidentFree: true,
          year: '2023',
          mileage: '10000',
          vin: '1HGCM82633A123456',
          tags: ['sedan', 'family'],
          negotiable: false,
          userType: 'private',
          hasVariants: false,
          attributes: [],
        })}
      >
        Submit Basic Info
      </button>
      <pre>{JSON.stringify(initialData, null, 2)}</pre>
    </div>
  ))
}));

jest.mock('../../features/CreateVehicleModal/components/steps/MediaUploadStep/MediaUploadStep', () => ({
  __esModule: true,
  MediaUploadStep: jest.fn(({ onComplete, initialData, mediaId }) => (
    <div data-testid="media-upload-step">
      <button 
        data-testid="submit-media-upload" 
        onClick={() => onComplete({
          thumbnail: 'vehicle1.jpg',
          images: ['vehicle1.jpg']
        })}
      >
        Submit Media
      </button>
      <div>Media ID: {mediaId}</div>
      <pre>{JSON.stringify(initialData, null, 2)}</pre>
    </div>
  ))
}));

jest.mock('../../features/CreateJobModal/components/steps/OptionalSettingsStep/OptionalSettingsStep', () => ({
  __esModule: true,
  OptionalInfoStep: jest.fn(({ onComplete, initialData }) => (
    <div data-testid="optional-info-step">
      <button 
        data-testid="submit-optional-info" 
        onClick={() => onComplete({
          weight: '1500',
          height: '150',
          width: '180',
          depth: '480',
          manageStocks: true,
          stock: '1',
          hasVariants: false,
          attributes: [],
          options: []
        })}
      >
        Submit Optional Info
      </button>
      <pre>{JSON.stringify(initialData, null, 2)}</pre>
    </div>
  ))
}));

jest.mock('../../features/CreateVehicleModal/components/steps/FinalizeStep/FinalizeStep', () => ({
  __esModule: true,
  FinalizeStep: jest.fn(({ dealData, initialLocation, onFinalize, isSuccess, onClose }) => (
    <div data-testid="finalize-step">
      <button 
        data-testid="submit-finalize" 
        onClick={() => onFinalize({ latitude: 40.7128, longitude: -74.0060 })}
      >
        Finalize Vehicle
      </button>
      {isSuccess && (
        <div data-testid="success-message">Vehicle successfully created!</div>
      )}
      <button data-testid="close-modal" onClick={onClose}>
        Close Modal
      </button>
      <pre>{JSON.stringify({ dealData, initialLocation }, null, 2)}</pre>
    </div>
  ))
}));

// Mock API functions
jest.mock('../../api/client/vehiclesApi', () => ({
  addVehicle: jest.fn().mockResolvedValue({ id: 'vehicle123' }),
  updateVehicle: jest.fn().mockResolvedValue({ id: 'vehicle123' }),
  finalizeVehicle: jest.fn().mockResolvedValue({ success: true })
}));

jest.mock('../../api/client/mediaApi', () => ({
  createMedia: jest.fn()
}));

jest.mock('../../api/categories', () => ({
  getCategories: jest.fn()
}));

describe('CreateVehicleModal Component', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    
    // Mock successful API responses
    vehiclesApi.addVehicle.mockResolvedValue({ id: 'vehicle123' });
    vehiclesApi.updateVehicle.mockResolvedValue({ id: 'vehicle123' });
    vehiclesApi.finalizeVehicle.mockResolvedValue({ success: true });
    mediaApi.createMedia.mockResolvedValue({ id: 'media123' });
    categoriesApi.getCategories.mockResolvedValue({
      categories: [
        { id: 'cat123', name: 'Automotive', slug: 'automotive' }
      ]
    });
  });

  it('renders in create mode', () => {
    render(<CreateVehicleModal onClose={() => {}} />);
    expect(screen.getByTestId('basic-info-step')).toBeInTheDocument();
  });

  it('fetches categories on mount', async () => {
    render(<CreateVehicleModal onClose={() => {}} />);
    await waitFor(() => {
      expect(categoriesApi.getCategories).toHaveBeenCalled();
    });
  });

  it('navigates through all steps and creates a vehicle', async () => {
    render(<CreateVehicleModal onClose={() => {}} />);

    // Step 1: Submit Basic Info
    expect(screen.getByTestId('basic-info-step')).toBeInTheDocument();
    fireEvent.click(screen.getByTestId('submit-basic-info'));

    // Step 2: Submit Media Upload
    await waitFor(() => {
      expect(screen.getByTestId('media-upload-step')).toBeInTheDocument();
    });
    fireEvent.click(screen.getByTestId('submit-media-upload'));

    // Step 3: Submit Optional Info
    await waitFor(() => {
      expect(screen.getByTestId('optional-info-step')).toBeInTheDocument();
    });
    fireEvent.click(screen.getByTestId('submit-optional-info'));

    // Step 4: Finalize
    await waitFor(() => {
      expect(screen.getByTestId('finalize-step')).toBeInTheDocument();
    });
    fireEvent.click(screen.getByTestId('submit-finalize'));

    // Verify API calls - update the expectation to check for updateVehicle instead of finalizeVehicle
    await waitFor(() => {
      expect(vehiclesApi.addVehicle).toHaveBeenCalled();
      expect(mediaApi.createMedia).toHaveBeenCalled();
      expect(vehiclesApi.updateVehicle).toHaveBeenCalled();
      expect(screen.getByTestId('success-message')).toBeInTheDocument();
    });
  });

  it('handles edit mode correctly', async () => {
    const initialVehicleData = {
      name: 'Existing Vehicle',
      description: 'Existing Description',
      basePrice: '30000',
      categoryId: 'cat123',
      categorySlug: 'automotive',
      brand: 'Honda',
      model: 'Accord',
      condition: 'excellent',
      performance: '180',
      fuelType: 'hybrid',
      transmissionType: 'automatic',
      numberOfOwners: '2',
      accidentFree: true,
      year: '2022',
      mileage: '20000',
      vin: '2HGES16575H123456',
      tags: ['sedan', 'hybrid'],
      negotiable: true,
      userType: 'dealer',
      hasVariants: false,
      attributes: []
    };

    render(
      <CreateVehicleModal 
        onClose={() => {}} 
        editMode={true}
        vehicleId="vehicle123"
        initialVehicleData={initialVehicleData}
      />
    );

    // Verify that BasicInfoStep contains the initial data
    await waitFor(() => {
      expect(screen.getByTestId('basic-info-step')).toBeInTheDocument();
    });
    
    // Verify that the initialData is passed correctly
    const basicInfoStepElement = screen.getByTestId('basic-info-step');
    expect(basicInfoStepElement).toContainHTML('Existing Vehicle');
    expect(basicInfoStepElement).toContainHTML('Existing Description');
    expect(basicInfoStepElement).toContainHTML('30000');
  });

  it('handles errors during vehicle update', async () => {
    vehiclesApi.updateVehicle.mockRejectedValue(new Error('Failed to save vehicle details.'));

    render(
      <CreateVehicleModal 
        onClose={() => {}} 
        editMode={true}
        vehicleId="vehicle123"
        initialVehicleData={{
          name: 'Test Vehicle',
          description: 'Test Description',
          basePrice: '25000'
        }}
      />
    );

    // Complete all steps
    fireEvent.click(screen.getByTestId('submit-basic-info'));
    
    // Wait for API call to complete and step transition
    await waitFor(() => {
      expect(screen.getByTestId('media-upload-step')).toBeInTheDocument();
    });
    fireEvent.click(screen.getByTestId('submit-media-upload'));
    
    await waitFor(() => {
      expect(screen.getByTestId('optional-info-step')).toBeInTheDocument();
    });
    fireEvent.click(screen.getByTestId('submit-optional-info'));

    // Verify error handling
    await waitFor(() => {
      expect(screen.getByText('Failed to save vehicle details.')).toBeInTheDocument();
    });
  });

  it('handles errors during finalization', async () => {
    // Mock errors for updateVehicle (used in the finalizing step)
    vehiclesApi.updateVehicle
      .mockResolvedValueOnce({ id: 'vehicle123' }) // For step 3
      .mockRejectedValueOnce(new Error('Failed to update vehicle')); // For step 4

    render(<CreateVehicleModal onClose={() => {}} />);

    // Step 1: Submit Basic Info
    fireEvent.click(screen.getByTestId('submit-basic-info'));
    
    // Step 2: Submit Media Upload
    await waitFor(() => {
      expect(screen.getByTestId('media-upload-step')).toBeInTheDocument();
    });
    fireEvent.click(screen.getByTestId('submit-media-upload'));
    
    // Step 3: Submit Optional Info - this needs to succeed to reach the finalize step
    await waitFor(() => {
      expect(screen.getByTestId('optional-info-step')).toBeInTheDocument();
    });
    fireEvent.click(screen.getByTestId('submit-optional-info'));
    
    // Step 4: Attempt to finalize (should fail)
    await waitFor(() => {
      expect(screen.getByTestId('finalize-step')).toBeInTheDocument();
    });
    fireEvent.click(screen.getByTestId('submit-finalize'));
    
    // Check for error message - using a more flexible approach
    await waitFor(() => {
      const errorElement = screen.getByRole('alert');
      expect(errorElement).toBeInTheDocument();
    });
  });
}); 