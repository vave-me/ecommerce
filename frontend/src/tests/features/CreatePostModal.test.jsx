import React from 'react';
import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import CreatePostModal from '../../features/CreatePostModal';
import * as postsApi from '../../api/client/postsApi';
import * as mediaApi from '../../api/client/mediaApi';
import { useAuth } from '../../context/AuthContext';

// Mock the AuthContext
jest.mock('../../context/AuthContext', () => ({
  useAuth: jest.fn()
}));

// Mock next-intl translation hook
jest.mock('next-intl', () => ({
  useTranslations: () => (key) => key
}));

// Mock the child components
jest.mock('../../features/CreatePostModal/components/steps/BasicInfoStep/BasicInfoStep', () => ({
  __esModule: true,
  BasicInfoStep: jest.fn(({ onSubmit, initialData, errors }) => (
    <div data-testid="basic-info-step">
      <button 
        data-testid="submit-basic-info" 
        onClick={() => onSubmit({
          name: 'Test Post',
          description: 'This is a detailed description for the test post that should be long enough to pass validation.',
          tags: 'tag1, tag2, tag3'
        })}
      >
        Submit Basic Info
      </button>
      {errors?.name && <div role="alert">{errors.name}</div>}
      {errors?.description && <div role="alert">{errors.description}</div>}
      <pre>{JSON.stringify(initialData, null, 2)}</pre>
    </div>
  ))
}));

jest.mock('../../features/CreatePostModal/components/steps/MediaUploadStep/MediaUploadStep', () => ({
  __esModule: true,
  MediaUploadStep: jest.fn(({ onComplete, initialImages, initialVideoUrl, mediaId }) => (
    <div data-testid="media-upload-step">
      <div>Media ID: {mediaId}</div>
      <button 
        data-testid="submit-media-upload" 
        onClick={() => onComplete({
          images: ['post1.jpg', 'post2.jpg'],
          videoUrl: 'https://example.com/video.mp4'
        })}
      >
        Submit Media
      </button>
      <pre>{JSON.stringify({ initialImages, initialVideoUrl }, null, 2)}</pre>
    </div>
  ))
}));

jest.mock('../../features/CreatePostModal/components/steps/OptionalSettingsStep/OptionalSettingsStep', () => ({
  __esModule: true,
  OptionalSettingsStep: jest.fn(({ onPublish, isLoading }) => (
    <div data-testid="optional-settings-step">
      <button 
        data-testid="submit-publish" 
        onClick={onPublish}
        disabled={isLoading}
      >
        Publish Post
      </button>
      {isLoading && <div>Loading...</div>}
    </div>
  ))
}));

jest.mock('../../features/CreatePostModal/components/steps/SuccessStep/SuccessStep', () => ({
  __esModule: true,
  SuccessStep: jest.fn(({ onViewDashboard, onClose }) => (
    <div data-testid="success-step">
      <button data-testid="view-dashboard" onClick={onViewDashboard}>
        View Dashboard
      </button>
      <button data-testid="close-modal" onClick={onClose}>
        Close Modal
      </button>
    </div>
  ))
}));

// Mock the FormActions component that's used within steps
jest.mock('../../common/components/FormActions', () => ({
  FormActions: jest.fn(({ onPrimaryAction, onCancel, isPrimaryDisabled }) => (
    <div className="form-actions">
      <button 
        onClick={onPrimaryAction} 
        disabled={isPrimaryDisabled}
        data-testid="primary-action"
      >
        Primary Action
      </button>
      <button onClick={onCancel} data-testid="cancel-action">
        Cancel
      </button>
    </div>
  ))
}));

// Mock the ErrorAlert component
jest.mock('../../common/components/ErrorAlert', () => ({
  ErrorAlert: ({ message }) => <div role="alert">{message}</div>
}));

// Mock the API functions
jest.mock('../../api/client/postsApi', () => ({
  addPost: jest.fn(),
  updatePost: jest.fn()
}));

jest.mock('../../api/client/mediaApi', () => ({
  createMedia: jest.fn(),
  addImage: jest.fn(),
  addVideo: jest.fn()
}));

// Mock the hooks
jest.mock('../../hooks/useMobileDetection', () => ({
  useMobileDetection: jest.fn(() => false)
}));

jest.mock('../../hooks/useFocusTrap', () => ({
  useFocusTrap: jest.fn(() => ({ current: null }))
}));

jest.mock('../../hooks/useAutoSave', () => ({
  useAutosave: jest.fn(() => ({ lastSaved: null, isSaving: false }))
}));

describe('CreatePostModal Component', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    
    // Mock successful user authentication
    useAuth.mockReturnValue({
      user: { userId: 'user123' },
      isUserLoggedIn: true
    });
    
    // Mock successful API responses
    postsApi.addPost.mockResolvedValue({ id: 'post123' });
    postsApi.updatePost.mockResolvedValue({ id: 'post123' });
    mediaApi.createMedia.mockResolvedValue({ id: 'media123' });
  });

  it('renders in create mode', () => {
    render(<CreatePostModal onClose={() => {}} />);
    expect(screen.getByTestId('basic-info-step')).toBeInTheDocument();
  });

  it('handles edit mode with initial data correctly', () => {
    const initialData = {
      id: 'post123',
      name: 'Existing Post',
      description: 'Existing Description',
      tags: ['tag1', 'tag2'],
      mediaId: 'media123',
      images: ['image1.jpg']
    };

    render(
      <CreatePostModal 
        onClose={() => {}} 
        editMode={true}
        initialData={initialData}
      />
    );

    // Verify BasicInfoStep receives the initial data
    const basicInfoStepElement = screen.getByTestId('basic-info-step');
    expect(basicInfoStepElement).toContainHTML('Existing Post');
    expect(basicInfoStepElement).toContainHTML('Existing Description');
    expect(basicInfoStepElement).toContainHTML('tag1, tag2');
  });

  it('navigates through all steps and creates a post successfully', async () => {
    render(<CreatePostModal onClose={() => {}} />);

    // Step 1: Submit Basic Info
    expect(screen.getByTestId('basic-info-step')).toBeInTheDocument();
    fireEvent.click(screen.getByTestId('submit-basic-info'));

    // Step 2: Submit Media Upload
    await waitFor(() => {
      expect(screen.getByTestId('media-upload-step')).toBeInTheDocument();
    });
    fireEvent.click(screen.getByTestId('submit-media-upload'));

    // Step 3: Publish
    await waitFor(() => {
      expect(screen.getByTestId('optional-settings-step')).toBeInTheDocument();
    });
    fireEvent.click(screen.getByTestId('submit-publish'));

    // Verify the post is created successfully and APIs are called
    await waitFor(() => {
      expect(postsApi.addPost).toHaveBeenCalled();
      expect(mediaApi.createMedia).toHaveBeenCalled();
      expect(postsApi.updatePost).toHaveBeenCalled();
      expect(screen.getByTestId('success-step')).toBeInTheDocument();
    });
  });

  it('handles API errors during post creation', async () => {
    // Mock API error
    postsApi.addPost.mockRejectedValueOnce(new Error('Failed to create post'));

    render(<CreatePostModal onClose={() => {}} />);

    // Submit basic info
    fireEvent.click(screen.getByTestId('submit-basic-info'));

    // Check for error message
    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument();
    });
  });

  it('handles API errors during publishing', async () => {
    // Make the first call succeed but the publishing step fail
    postsApi.addPost.mockResolvedValueOnce({ id: 'post123' });
    postsApi.updatePost.mockRejectedValueOnce(new Error('Failed to publish post'));

    render(<CreatePostModal onClose={() => {}} />);

    // Complete all steps
    fireEvent.click(screen.getByTestId('submit-basic-info'));
    
    await waitFor(() => {
      expect(screen.getByTestId('media-upload-step')).toBeInTheDocument();
    });
    fireEvent.click(screen.getByTestId('submit-media-upload'));
    
    await waitFor(() => {
      expect(screen.getByTestId('optional-settings-step')).toBeInTheDocument();
    });
    fireEvent.click(screen.getByTestId('submit-publish'));

    // Verify error handling
    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument();
    });
  });

  it('requires user login to submit a post', async () => {
    // Mock user not logged in
    useAuth.mockReturnValue({
      user: null,
      isUserLoggedIn: false
    });

    render(<CreatePostModal onClose={() => {}} />);
    
    // Try to submit basic info
    fireEvent.click(screen.getByTestId('submit-basic-info'));
    
    // Expect login error message
    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument();
    });
    
    // Verify API calls are not made
    expect(postsApi.addPost).not.toHaveBeenCalled();
  });

  it('validates form inputs before submission', async () => {
    // Mock the BasicInfoStep submission to return invalid data
    require('../../features/CreatePostModal/components/steps/BasicInfoStep/BasicInfoStep').BasicInfoStep.mockImplementationOnce(
      ({ onSubmit }) => (
        <div data-testid="basic-info-step">
          <button 
            data-testid="submit-invalid-form" 
            onClick={() => onSubmit({
              name: '', // Invalid: empty name
              description: 'Too short', // Invalid: too short
              tags: 'tag1, tag2'
            })}
          >
            Submit Invalid Form
          </button>
        </div>
      )
    );

    render(<CreatePostModal onClose={() => {}} />);

    // Submit invalid form
    fireEvent.click(screen.getByTestId('submit-invalid-form'));

    // Verify API calls are not made and we stay on the same step
    await waitFor(() => {
      expect(postsApi.addPost).not.toHaveBeenCalled();
      expect(screen.getByTestId('basic-info-step')).toBeInTheDocument();
    });
  });

  it('creates and updates media container correctly', async () => {
    render(<CreatePostModal onClose={() => {}} />);

    // Step 1: Submit Basic Info
    fireEvent.click(screen.getByTestId('submit-basic-info'));

    // Verify media container is created
    await waitFor(() => {
      expect(mediaApi.createMedia).toHaveBeenCalledWith({
        itemId: 'post123',
        itemType: 'post',
        userId: 'user123'
      });
      expect(screen.getByTestId('media-upload-step')).toBeInTheDocument();
      expect(screen.getByText('Media ID: media123')).toBeInTheDocument();
    });
  });

  it('handles media creation errors', async () => {
    // First make post creation succeed but media creation fail
    postsApi.addPost.mockResolvedValueOnce({ id: 'post123' });
    mediaApi.createMedia.mockRejectedValueOnce(new Error('Failed to create media'));

    render(<CreatePostModal onClose={() => {}} />);

    // Submit basic info
    fireEvent.click(screen.getByTestId('submit-basic-info'));

    // Check for error message
    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument();
    });
  });

  it('supports closing the modal', () => {
    const onCloseMock = jest.fn();
    render(<CreatePostModal onClose={onCloseMock} />);
    
    // Click the close button in the modal header (X button)
    fireEvent.click(screen.getByLabelText('Close Modal'));
    
    expect(onCloseMock).toHaveBeenCalled();
  });
}); 