import { renderHook, act } from '@testing-library/react';
import { useCategoryKeyboardNav } from '@/hooks/useCategoryKeyboardNav.jsx';

describe('useCategoryKeyboardNav hook', () => {
  // Mock data setup
  const mockCategories = [
    { id: 'cat1', name: 'Category 1', parentId: null, subcategories: ['cat2', 'cat3'] },
    { id: 'cat2', name: 'Subcategory 1', parentId: 'cat1', subcategories: [] },
    { id: 'cat3', name: 'Subcategory 2', parentId: 'cat1', subcategories: ['cat4'] },
    { id: 'cat4', name: 'Nested Subcategory', parentId: 'cat3', subcategories: [] }
  ];
  
  const mockExpandedMap = {
    cat1: true,
    cat3: true
  };
  
  // Mock functions
  const mockToggleExpand = jest.fn();
  const mockOnSelect = jest.fn();
  const mockSetFocusedId = jest.fn();
  
  beforeEach(() => {
    jest.clearAllMocks();
  });
  
  test('should provide handleKeyDown function', () => {
    const { result } = renderHook(() => 
      useCategoryKeyboardNav(
        mockCategories, 
        mockExpandedMap, 
        mockToggleExpand, 
        mockOnSelect, 
        mockSetFocusedId
      )
    );
    
    expect(typeof result.current.handleKeyDown).toBe('function');
  });
  
  test('should handle ArrowDown to navigate to next category', () => {
    const { result } = renderHook(() => 
      useCategoryKeyboardNav(
        mockCategories, 
        mockExpandedMap, 
        mockToggleExpand, 
        mockOnSelect, 
        mockSetFocusedId
      )
    );
    
    const mockEvent = {
      key: 'ArrowDown',
      preventDefault: jest.fn(),
      stopPropagation: jest.fn()
    };
    
    act(() => {
      result.current.handleKeyDown(mockEvent, mockCategories[0]);
    });
    
    expect(mockEvent.preventDefault).toHaveBeenCalled();
    expect(mockEvent.stopPropagation).toHaveBeenCalled();
    expect(mockSetFocusedId).toHaveBeenCalledWith('cat2');
  });
  
  test('should handle ArrowUp to navigate to previous category', () => {
    const { result } = renderHook(() => 
      useCategoryKeyboardNav(
        mockCategories, 
        mockExpandedMap, 
        mockToggleExpand, 
        mockOnSelect, 
        mockSetFocusedId
      )
    );
    
    const mockEvent = {
      key: 'ArrowUp',
      preventDefault: jest.fn(),
      stopPropagation: jest.fn()
    };
    
    // Start from cat2 (index 1) and go up to cat1
    act(() => {
      result.current.handleKeyDown(mockEvent, mockCategories[1]);
    });
    
    expect(mockSetFocusedId).toHaveBeenCalledWith('cat1');
  });
  
  test('should handle ArrowRight to expand collapsed category', () => {
    // Create a new expandedMap where cat1 is collapsed
    const collapsedMap = { cat1: false };
    
    const { result } = renderHook(() => 
      useCategoryKeyboardNav(
        mockCategories, 
        collapsedMap, 
        mockToggleExpand, 
        mockOnSelect, 
        mockSetFocusedId
      )
    );
    
    const mockEvent = {
      key: 'ArrowRight',
      preventDefault: jest.fn(),
      stopPropagation: jest.fn()
    };
    
    act(() => {
      result.current.handleKeyDown(mockEvent, mockCategories[0]);
    });
    
    expect(mockToggleExpand).toHaveBeenCalledWith('cat1');
    expect(mockSetFocusedId).not.toHaveBeenCalled();
  });
  
  test('should handle ArrowRight to move focus to first child when category is expanded', () => {
    const { result } = renderHook(() => 
      useCategoryKeyboardNav(
        mockCategories, 
        mockExpandedMap,  // cat1 is already expanded
        mockToggleExpand, 
        mockOnSelect, 
        mockSetFocusedId
      )
    );
    
    const mockEvent = {
      key: 'ArrowRight',
      preventDefault: jest.fn(),
      stopPropagation: jest.fn()
    };
    
    act(() => {
      result.current.handleKeyDown(mockEvent, mockCategories[0]);
    });
    
    expect(mockToggleExpand).not.toHaveBeenCalled();
    expect(mockSetFocusedId).toHaveBeenCalledWith('cat2');
  });
  
  test('should handle ArrowLeft to collapse expanded category', () => {
    const { result } = renderHook(() => 
      useCategoryKeyboardNav(
        mockCategories, 
        mockExpandedMap,  // cat1 is expanded
        mockToggleExpand, 
        mockOnSelect, 
        mockSetFocusedId
      )
    );
    
    const mockEvent = {
      key: 'ArrowLeft',
      preventDefault: jest.fn(),
      stopPropagation: jest.fn()
    };
    
    act(() => {
      result.current.handleKeyDown(mockEvent, mockCategories[0]);
    });
    
    expect(mockToggleExpand).toHaveBeenCalledWith('cat1');
    expect(mockSetFocusedId).not.toHaveBeenCalled();
  });
  
  test('should handle ArrowLeft to move focus to parent when category is collapsed', () => {
    // Create a new expandedMap where cat3 is collapsed
    const modifiedMap = { cat1: true, cat3: false };
    
    const { result } = renderHook(() => 
      useCategoryKeyboardNav(
        mockCategories, 
        modifiedMap,
        mockToggleExpand, 
        mockOnSelect, 
        mockSetFocusedId
      )
    );
    
    const mockEvent = {
      key: 'ArrowLeft',
      preventDefault: jest.fn(),
      stopPropagation: jest.fn()
    };
    
    // Use cat4 which has cat3 as parent
    act(() => {
      result.current.handleKeyDown(mockEvent, mockCategories[3]);
    });
    
    expect(mockToggleExpand).not.toHaveBeenCalled();
    expect(mockSetFocusedId).toHaveBeenCalledWith('cat3');
  });
  
  test('should handle Enter/Space to select category', () => {
    const { result } = renderHook(() => 
      useCategoryKeyboardNav(
        mockCategories, 
        mockExpandedMap,
        mockToggleExpand, 
        mockOnSelect, 
        mockSetFocusedId
      )
    );
    
    // Test Enter key
    const enterEvent = {
      key: 'Enter',
      preventDefault: jest.fn(),
      stopPropagation: jest.fn()
    };
    
    act(() => {
      result.current.handleKeyDown(enterEvent, mockCategories[1]);
    });
    
    expect(mockOnSelect).toHaveBeenCalledWith(mockCategories[1]);
    
    mockOnSelect.mockClear();
    
    // Test Space key
    const spaceEvent = {
      key: ' ',
      preventDefault: jest.fn(),
      stopPropagation: jest.fn()
    };
    
    act(() => {
      result.current.handleKeyDown(spaceEvent, mockCategories[1]);
    });
    
    expect(mockOnSelect).toHaveBeenCalledWith(mockCategories[1]);
  });
  
  test('should do nothing for unhandled keys', () => {
    const { result } = renderHook(() => 
      useCategoryKeyboardNav(
        mockCategories, 
        mockExpandedMap,
        mockToggleExpand, 
        mockOnSelect, 
        mockSetFocusedId
      )
    );
    
    const mockEvent = {
      key: 'Tab',
      preventDefault: jest.fn(),
      stopPropagation: jest.fn()
    };
    
    act(() => {
      result.current.handleKeyDown(mockEvent, mockCategories[0]);
    });
    
    expect(mockEvent.preventDefault).not.toHaveBeenCalled();
    expect(mockEvent.stopPropagation).not.toHaveBeenCalled();
    expect(mockSetFocusedId).not.toHaveBeenCalled();
    expect(mockToggleExpand).not.toHaveBeenCalled();
    expect(mockOnSelect).not.toHaveBeenCalled();
  });
  
  test('should handle boundary conditions for navigation', () => {
    const { result } = renderHook(() => 
      useCategoryKeyboardNav(
        mockCategories, 
        mockExpandedMap,
        mockToggleExpand, 
        mockOnSelect, 
        mockSetFocusedId
      )
    );
    
    // Try to navigate up from first category
    const upEvent = {
      key: 'ArrowUp',
      preventDefault: jest.fn(),
      stopPropagation: jest.fn()
    };
    
    act(() => {
      result.current.handleKeyDown(upEvent, mockCategories[0]);
    });
    
    // Should still be on the first category
    expect(mockSetFocusedId).toHaveBeenCalledWith('cat1');
    
    mockSetFocusedId.mockClear();
    
    // Try to navigate down from last category
    const downEvent = {
      key: 'ArrowDown',
      preventDefault: jest.fn(),
      stopPropagation: jest.fn()
    };
    
    act(() => {
      result.current.handleKeyDown(downEvent, mockCategories[3]);
    });
    
    // Should still be on the last category
    expect(mockSetFocusedId).toHaveBeenCalledWith('cat4');
  });
}); 