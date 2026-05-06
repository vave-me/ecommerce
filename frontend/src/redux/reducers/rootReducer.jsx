// src/redux/reducers/rootReducer.js
import {combineReducers} from 'redux';
import listingFiltersReducer from '../slices/listingFiltersSlice';
import modalsReducer from '../slices/modalsSlice';
import appModeReducer from '../slices/appModeSlice';
import uiPreferencesReducer from '../slices/uiPreferencesSlice';

const rootReducer = combineReducers({
    // 'listingFilters' is the name you use in useSelector:
    listingFilters: listingFiltersReducer,
    modals: modalsReducer,
    appMode: appModeReducer,
    uiPreferences: uiPreferencesReducer,
});
export default rootReducer;
