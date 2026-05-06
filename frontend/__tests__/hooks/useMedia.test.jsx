import { renderHook, act } from '@testing-library/react';
import { useMedia, processMediaResponse } from '@/hooks/useMedia.jsx';
import { getMediaByItem } from '@/api/mediaApi.jsx';

jest.mock('@/api/mediaApi.jsx', () => ({
  getMediaByItem: jest.fn()
}));

describe('processMediaResponse function', () => {
  test('should return empty array for invalid response', () => {
    expect(processMediaResponse(null)).toEqual([]);
    expect(processMediaResponse({})).toEqual([]);
    expect(processMediaResponse({ media: {} })).toEqual([]);
    expect(processMediaResponse({ media: { mediaOrder: null } })).toEqual([]);
  });
  
  test('should correctly process image media items', () => {
    const response = {
      media: {
        mediaOrder: [
          {
            mediaItemId: 'img1',
            url: 'https://example.com/image.jpg',
            altText: 'Example image',
            displayOrder: 0
          }
        ]
      }
    };
    
    const result = processMediaResponse(response);
    
    expect(result).toEqual([
      {
        id: 'img1',
        type: 'image',
        src: 'https://example.com/image.jpg',
        alt: 'Example image',
        displayOrder: 0,
        poster: undefined
      }
    ]);
  });
  
  test('should correctly identify video media types by extension', () => {
    const response = {
      media: {
        mediaOrder: [
          {
            mediaItemId: 'vid1',
            url: 'https://example.com/video.mp4',
            altText: 'Example video',
            displayOrder: 0
          },
          {
            mediaItemId: 'vid2',
            url: 'https://example.com/video.mov',
            altText: 'Another video',
            displayOrder: 1
          },
          {
            mediaItemId: 'vid3',
            url: 'https://example.com/video.webm',
            altText: 'WEBM video',
            displayOrder: 2
          }
        ]
      }
    };
    
    const result = processMediaResponse(response);
    
    expect(result.every(item => item.type === 'video')).toBe(true);
  });
  
  test('should identify video media by path pattern', () => {
    const response = {
      media: {
        mediaOrder: [
          {
            mediaItemId: 'vid1',
            url: 'https://example.com/video/12345',
            altText: 'Video by path pattern',
            displayOrder: 0
          }
        ]
      }
    };
    
    const result = processMediaResponse(response);
    
    expect(result[0].type).toBe('video');
  });
  
  test('should use item type property if provided', () => {
    const response = {
      media: {
        mediaOrder: [
          {
            mediaItemId: 'vid1',
            url: 'https://example.com/content.bin',
            type: 'video',
            altText: 'Video with explicit type',
            displayOrder: 0
          }
        ]
      }
    };
    
    const result = processMediaResponse(response);
    
    expect(result[0].type).toBe('video');
  });
  
  test('should handle missing properties gracefully', () => {
    const response = {
      media: {
        mediaOrder: [
          {
            url: 'https://example.com/image.jpg'
            // Missing mediaItemId, altText, displayOrder
          }
        ]
      }
    };
    
    const result = processMediaResponse(response);
    
    expect(result).toEqual([
      {
        id: 'media-0',
        type: 'image',
        src: 'https://example.com/image.jpg',
        alt: 'Media 1',
        displayOrder: 0,
        poster: undefined
      }
    ]);
  });
});

describe('useMedia hook', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });
  
  test('should return initial empty state when itemId is not provided', () => {
    const { result } = renderHook(() => useMedia(null));
    
    expect(result.current.media).toEqual([]);
    expect(result.current.loading).toBe(false);
    expect(result.current.error).toBe(null);
    expect(result.current.thumbnail).toBe(null);
    expect(result.current.mediaId).toBe(null);
    expect(result.current.isEmpty).toBe(true);
  });
  
  test('should fetch and process media when itemId is provided', async () => {
    const mockMediaResponse = {
      media: {
        id: 'media123',
        mediaOrder: [
          {
            mediaItemId: 'img1',
            url: 'https://example.com/image.jpg',
            altText: 'Test image',
            displayOrder: 0
          }
        ]
      }
    };
    
    getMediaByItem.mockResolvedValue(mockMediaResponse);
    
    let hook;
    
    await act(async () => {
      hook = renderHook(() => useMedia('item123'));
    });
    
    expect(getMediaByItem).toHaveBeenCalledWith('item123');
    
    expect(hook.result.current.media).toEqual([
      {
        id: 'img1',
        type: 'image',
        src: 'https://example.com/image.jpg',
        alt: 'Test image',
        displayOrder: 0,
        poster: undefined
      }
    ]);
    
    expect(hook.result.current.loading).toBe(false);
    expect(hook.result.current.error).toBe(null);
    expect(hook.result.current.thumbnail).toBe('https://example.com/image.jpg');
    expect(hook.result.current.mediaId).toBe('media123');
    expect(hook.result.current.isEmpty).toBe(false);
  });
  
  test('should handle API errors gracefully', async () => {
    const testError = new Error('API error');
    getMediaByItem.mockRejectedValue(testError);
    
    let hook;
    
    await act(async () => {
      hook = renderHook(() => useMedia('item123'));
    });
    
    expect(hook.result.current.media).toEqual([]);
    expect(hook.result.current.loading).toBe(false);
    expect(hook.result.current.error).toBe(testError);
    expect(hook.result.current.isEmpty).toBe(true);
  });
  
  test('should set first video as thumbnail if no images available', async () => {
    const mockMediaResponse = {
      media: {
        id: 'media123',
        mediaOrder: [
          {
            mediaItemId: 'vid1',
            url: 'https://example.com/video.mp4',
            altText: 'Test video',
            displayOrder: 0
          }
        ]
      }
    };
    
    getMediaByItem.mockResolvedValue(mockMediaResponse);
    
    let hook;
    
    await act(async () => {
      hook = renderHook(() => useMedia('item123'));
    });
    
    expect(hook.result.current.thumbnail).toBe('https://example.com/video.mp4');
  });
}); 