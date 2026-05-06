import { createSlice } from '@reduxjs/toolkit';

const initialState = {
  isVaveAdVisible: true,
  lastDismissed: null,
};

const adsSlice = createSlice({
  name: 'ads',
  initialState,
  reducers: {
    dismissVaveAd: (state) => {
      state.isVaveAdVisible = false;
      state.lastDismissed = new Date().toISOString();
    },
    showVaveAd: (state) => {
      state.isVaveAdVisible = true;
    },
  },
});

export const { dismissVaveAd, showVaveAd } = adsSlice.actions;
export default adsSlice.reducer;