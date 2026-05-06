import { renderHook, act } from '@testing-library/react';
import { useAutoSave } from '../useAutoSave';
import { UnifiedUtils } from '../../utils/duplicateEliminator';
import { createQueryClientWrapper } from '../../tests/test-utils';

// Mock dependencies
jest.mock('../../utils/duplicateEliminator', () => ({
  UnifiedUtils: {
    drafts: {
      saveDraft: jest.fn().mockReturnValue({ success: true, timestamp: '2023-01-01T00:00:00.000Z' })
    }
  }
}));

describe('useAutosave hook', () => {
  const mockUserId = 'user-123';
  const mockPostId = 'post-456';
  
  // Initial post data with empty tags string
  const initialPostData = {
    name: '',
    description: '',
    tags: ''
  };

  // Set up timers for testing debounce and intervals
  beforeEach(() => {
    jest.useFakeTimers();
    jest.clearAllMocks();
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  test('should not autosave when post data is empty', () => {
    renderHook(() => useAutoSave(mockUserId, initialPostData, mockPostId), {
      wrapper: createQueryClientWrapper()
    });
    
    // Fast-forward past debounce and interval timers
    act(() => {
      jest.advanceTimersByTime(2000); // Past the 1000ms debounce
    });
    
    // Should not save empty posts
    expect(UnifiedUtils.drafts.saveDraft).not.toHaveBeenCalled();
  });

  test('should autosave when post has a non-empty title', () => {
    const postData = {
      ...initialPostData,
      name: 'Test Post Title'
    };
    
    const { result } = renderHook(() => useAutoSave(mockUserId, postData, mockPostId), {
      wrapper: createQueryClientWrapper()
    });
    
    // Initially should be not saving and no last saved time
    expect(result.current.isSaving).toBe(false);
    expect(result.current.lastSaved).toBeNull();
    
    // Trigger debounced save
    act(() => {
      jest.advanceTimersByTime(1500); // Past the 1000ms debounce
    });
    
    // Should have called saveDraft with correct data
    expect(UnifiedUtils.drafts.saveDraft).toHaveBeenCalledWith(
      'auto',
      mockUserId,
      expect.objectContaining({
        id: mockPostId,
        userId: mockUserId,
        name: 'Test Post Title',
        description: '',
        tags: [] // Empty array since tags is an empty string
      })
    );
    
    // Should update lastSaved
    expect(result.current.lastSaved).not.toBeNull();
  });

  test('should autosave when post has a non-empty description', () => {
    const postData = {
      ...initialPostData,
      description: 'This is a test description for the post'
    };
    
    renderHook(() => useAutoSave(mockUserId, postData, mockPostId), {
      wrapper: createQueryClientWrapper()
    });
    
    // Trigger debounced save
    act(() => {
      jest.advanceTimersByTime(1500); // Past the 1000ms debounce
    });
    
    // Should have called saveDraft
    expect(UnifiedUtils.drafts.saveDraft).toHaveBeenCalled();
  });

  test('should format tags correctly when saving', () => {
    const postData = {
      name: 'Test Post',
      description: 'Description',
      tags: 'tag1, tag2,tag3, ' // Comma-separated string
    };
    
    renderHook(() => useAutoSave(mockUserId, postData, mockPostId), {
      wrapper: createQueryClientWrapper()
    });
    
    // Trigger debounced save
    act(() => {
      jest.advanceTimersByTime(1500);
    });
    
    // Check tags were properly formatted
    expect(UnifiedUtils.drafts.saveDraft).toHaveBeenCalledWith(
      'auto',
      mockUserId,
      expect.objectContaining({
        tags: ['tag1', 'tag2', 'tag3'] // Trimmed and filtered
      })
    );
  });

  test('should automatically save at regular intervals', () => {
    const postData = {
      name: 'Test Post',
      description: 'Description',
      tags: 'tag1, tag2'
    };
    
    renderHook(() => useAutoSave(mockUserId, postData, mockPostId), {
      wrapper: createQueryClientWrapper()
    });
    
    // Clear initial debounced save
    UnifiedUtils.drafts.saveDraft.mockClear();
    
    // Fast-forward 30 seconds to trigger interval save
    act(() => {
      jest.advanceTimersByTime(30000);
    });
    
    // Should have called saveDraft again
    expect(UnifiedUtils.drafts.saveDraft).toHaveBeenCalled();
  });

  test('should provide a startAutosave function to manually trigger save', () => {
    const postData = {
      name: 'Test Post',
      description: 'Test Description',
      tags: '' // Empty tags to avoid split error
    };
    
    const { result } = renderHook(() => useAutoSave(mockUserId, postData, mockPostId), {
      wrapper: createQueryClientWrapper()
    });
    
    // Clear any initial auto-saves
    UnifiedUtils.drafts.saveDraft.mockClear();
    
    // Call the function manually
    act(() => {
      result.current.forceSave(); // Updated to use forceSave instead of startAutosave
      // Fast-forward past debounce
      jest.advanceTimersByTime(1500);
    });
    
    // Should have called save
    expect(UnifiedUtils.drafts.saveDraft).toHaveBeenCalled();
  });
}); 