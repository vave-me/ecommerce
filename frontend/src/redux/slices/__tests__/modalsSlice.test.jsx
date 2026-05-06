import modalsReducer, {
  openCommentsFullModal,
  closeCommentsFullModal,
  openMessageModal,
  closeMessageModal,
  openAddProductModal,
  closeAddProductModal,
  openAddPostModal,
  closeAddPostModal,
  openAddVehicleModal,
  closeAddVehicleModal,
  openAddDealModal,
  closeAddDealModal,
  openAddPropertyModal,
  closeAddPropertyModal,
  openAddServiceModal,
  closeAddServiceModal,
  openAddJobModal,
  closeAddJobModal,
  openProductModal,
  closeProductModal,
  openAddVideoModal,
  closeAddVideoModal
} from '../modalsSlice';

describe('modalsSlice', () => {
  const initialState = {
    commentsFullModalOpen: false,
    commentsFullModalItemId: null,
    commentsFullItemType: null,
    commentsFullCategoryId: null,
    messageModalOpen: false,
    messageModalItemId: null,
    messageRecipientId: null,
    addProductModalOpen: false,
    addVehicleModalOpen: false,
    addDealModalOpen: false,
    addPropertyModalOpen: false,
    addJobModalOpen: false,
    addServiceModalOpen: false,
    addPostModalOpen: false,
    isProductModalOpen: false,
    selectedProduct: null,
    isVideoModalOpen: false,
    openModalsCount: 0,
  };

  beforeEach(() => {
    // Mock console.log to avoid cluttering test output
    jest.spyOn(console, 'log').mockImplementation(() => {});
  });

  test('should return the initial state', () => {
    expect(modalsReducer(undefined, { type: undefined })).toEqual(initialState);
  });

  // Comments modal tests
  test('should handle openCommentsFullModal', () => {
    const payload = {
      itemId: 'item123',
      itemType: 'product',
      categoryId: 'cat456'
    };
    
    expect(
      modalsReducer(initialState, openCommentsFullModal(payload))
    ).toEqual({
      ...initialState,
      commentsFullModalOpen: true,
      commentsFullModalItemId: 'item123',
      commentsFullItemType: 'product',
      commentsFullCategoryId: 'cat456',
      openModalsCount: 1
    });
  });

  test('should handle closeCommentsFullModal', () => {
    const stateWithOpenCommentsModal = {
      ...initialState,
      commentsFullModalOpen: true,
      commentsFullModalItemId: 'item123',
      commentsFullItemType: 'product',
      commentsFullCategoryId: 'cat456',
      openModalsCount: 1
    };
    
    expect(
      modalsReducer(stateWithOpenCommentsModal, closeCommentsFullModal())
    ).toEqual({
      ...initialState,
      openModalsCount: 0
    });
  });

  // Message modal tests
  test('should handle openMessageModal', () => {
    const payload = {
      itemId: 'item123',
      recipientId: 'user456'
    };
    
    expect(
      modalsReducer(initialState, openMessageModal(payload))
    ).toEqual({
      ...initialState,
      messageModalOpen: true,
      messageModalItemId: 'item123',
      messageRecipientId: 'user456',
      openModalsCount: 1
    });
  });

  test('should handle closeMessageModal', () => {
    const stateWithOpenMessageModal = {
      ...initialState,
      messageModalOpen: true,
      messageModalItemId: 'item123',
      messageRecipientId: 'user456',
      openModalsCount: 1
    };
    
    expect(
      modalsReducer(stateWithOpenMessageModal, closeMessageModal())
    ).toEqual({
      ...initialState,
      openModalsCount: 0
    });
  });

  // Product modal tests
  test('should handle openAddProductModal', () => {
    expect(
      modalsReducer(initialState, openAddProductModal())
    ).toEqual({
      ...initialState,
      addProductModalOpen: true,
      openModalsCount: 1
    });
  });

  test('should handle closeAddProductModal', () => {
    const stateWithOpenAddProductModal = {
      ...initialState,
      addProductModalOpen: true,
      openModalsCount: 1
    };
    
    expect(
      modalsReducer(stateWithOpenAddProductModal, closeAddProductModal())
    ).toEqual({
      ...initialState,
      openModalsCount: 0
    });
  });

  // Post modal tests
  test('should handle openAddPostModal', () => {
    expect(
      modalsReducer(initialState, openAddPostModal())
    ).toEqual({
      ...initialState,
      addPostModalOpen: true,
      openModalsCount: 1
    });
  });

  test('should handle closeAddPostModal', () => {
    const stateWithOpenAddPostModal = {
      ...initialState,
      addPostModalOpen: true,
      openModalsCount: 1
    };
    
    expect(
      modalsReducer(stateWithOpenAddPostModal, closeAddPostModal())
    ).toEqual({
      ...initialState,
      openModalsCount: 0
    });
  });

  // Vehicle modal tests
  test('should handle openAddVehicleModal', () => {
    expect(
      modalsReducer(initialState, openAddVehicleModal())
    ).toEqual({
      ...initialState,
      addVehicleModalOpen: true,
      openModalsCount: 1
    });
  });

  test('should handle closeAddVehicleModal', () => {
    const stateWithOpenAddVehicleModal = {
      ...initialState,
      addVehicleModalOpen: true,
      openModalsCount: 1
    };
    
    expect(
      modalsReducer(stateWithOpenAddVehicleModal, closeAddVehicleModal())
    ).toEqual({
      ...initialState,
      openModalsCount: 0
    });
  });

  // Product detail modal tests
  test('should handle openProductModal', () => {
    const product = { id: 'prod123', name: 'Test Product' };
    
    expect(
      modalsReducer(initialState, openProductModal({ product }))
    ).toEqual({
      ...initialState,
      isProductModalOpen: true,
      selectedProduct: product,
      openModalsCount: 1
    });
  });

  test('should handle closeProductModal', () => {
    const stateWithOpenProductModal = {
      ...initialState,
      isProductModalOpen: true,
      selectedProduct: { id: 'prod123', name: 'Test Product' },
      openModalsCount: 1
    };
    
    expect(
      modalsReducer(stateWithOpenProductModal, closeProductModal())
    ).toEqual({
      ...initialState,
      openModalsCount: 0
    });
  });

  // Video modal tests
  test('should handle openAddVideoModal', () => {
    expect(
      modalsReducer(initialState, openAddVideoModal())
    ).toEqual({
      ...initialState,
      isVideoModalOpen: true,
      openModalsCount: 1
    });
  });

  test('should handle closeAddVideoModal', () => {
    const stateWithOpenVideoModal = {
      ...initialState,
      isVideoModalOpen: true,
      openModalsCount: 1
    };
    
    expect(
      modalsReducer(stateWithOpenVideoModal, closeAddVideoModal())
    ).toEqual({
      ...initialState,
      openModalsCount: 0
    });
  });

  // Deal modal tests
  test('should handle openAddDealModal', () => {
    expect(
      modalsReducer(initialState, openAddDealModal())
    ).toEqual({
      ...initialState,
      addDealModalOpen: true,
      openModalsCount: 1
    });
  });

  test('should handle closeAddDealModal', () => {
    const stateWithOpenDealModal = {
      ...initialState,
      addDealModalOpen: true,
      openModalsCount: 1
    };
    
    expect(
      modalsReducer(stateWithOpenDealModal, closeAddDealModal())
    ).toEqual({
      ...initialState,
      openModalsCount: 0
    });
  });

  // Property modal tests
  test('should handle openAddPropertyModal', () => {
    expect(
      modalsReducer(initialState, openAddPropertyModal())
    ).toEqual({
      ...initialState,
      addPropertyModalOpen: true,
      openModalsCount: 1
    });
  });

  test('should handle closeAddPropertyModal', () => {
    const stateWithOpenPropertyModal = {
      ...initialState,
      addPropertyModalOpen: true,
      openModalsCount: 1
    };
    
    expect(
      modalsReducer(stateWithOpenPropertyModal, closeAddPropertyModal())
    ).toEqual({
      ...initialState,
      openModalsCount: 0
    });
  });

  // Service modal tests
  test('should handle openAddServiceModal', () => {
    expect(
      modalsReducer(initialState, openAddServiceModal())
    ).toEqual({
      ...initialState,
      addServiceModalOpen: true,
      openModalsCount: 1
    });
  });

  test('should handle closeAddServiceModal', () => {
    const stateWithOpenServiceModal = {
      ...initialState,
      addServiceModalOpen: true,
      openModalsCount: 1
    };
    
    expect(
      modalsReducer(stateWithOpenServiceModal, closeAddServiceModal())
    ).toEqual({
      ...initialState,
      openModalsCount: 0
    });
  });

  // Job modal tests
  test('should handle openAddJobModal', () => {
    expect(
      modalsReducer(initialState, openAddJobModal())
    ).toEqual({
      ...initialState,
      addJobModalOpen: true,
      openModalsCount: 1
    });
  });

  test('should handle closeAddJobModal', () => {
    const stateWithOpenJobModal = {
      ...initialState,
      addJobModalOpen: true,
      openModalsCount: 1
    };
    
    expect(
      modalsReducer(stateWithOpenJobModal, closeAddJobModal())
    ).toEqual({
      ...initialState,
      openModalsCount: 0
    });
  });

  // Test multiple open modals
  test('should correctly track multiple open modals', () => {
    let state = initialState;
    
    // Open first modal
    state = modalsReducer(state, openAddProductModal());
    expect(state.openModalsCount).toBe(1);
    
    // Open second modal
    state = modalsReducer(state, openAddPostModal());
    expect(state.openModalsCount).toBe(2);
    
    // Close first modal
    state = modalsReducer(state, closeAddProductModal());
    expect(state.openModalsCount).toBe(1);
    
    // Close second modal
    state = modalsReducer(state, closeAddPostModal());
    expect(state.openModalsCount).toBe(0);
  });
}); 