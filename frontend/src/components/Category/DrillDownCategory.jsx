"use client";
import React, { useState, memo } from "react";
import { useTranslations } from "next-intl"; //  Import hook
import DrillDownCategoryPanel from "./DrillDownCategoryPanel"; // Assuming this handles its own translations
/**
 * The main orchestrator for a multi-level "drill-down" UI.
 * Manages the category path stack and renders the active panel.
 * Uses next-intl for translations.
 * 
 * OPTIMIZED: React.memo for better category navigation performance
 */
const DrillDownCategory = memo(function DrillDownCategory({
    rootCategory, // The starting point category object
    // onFinalSelect = () => {}, // Prop not used in this component, likely used by DrillDownCategoryPanel
}) {
    const t = useTranslations('DrillDownCategory'); //  Instantiate hook with namespace
    // A stack of categories => path to current
    // e.g. [rootCategory, subCategory, subSubCategory, ... ]
    const [path, setPath] = useState(rootCategory ? [rootCategory] : []);
    const activeCategory = path.length > 0 ? path[path.length - 1] : null;
    const parentCategory = path.length > 1 ? path[path.length - 2] : null;
    // Called by DrillDownCategoryPanel when a subcategory is chosen
    const handleSelectCategory = (cat) => {
        // If the chosen subcategory has more children, we push it.
        // The decision to call onFinalSelect might happen inside DrillDownCategoryPanel
        // based on whether the selected 'cat' has children itself.
        setPath((prev) => [...prev, cat]);
    };
    // Called by DrillDownCategoryPanel when the back button is clicked
    const handleBack = () => {
        if (path.length > 1) {
            setPath((prev) => prev.slice(0, prev.length - 1));
        }
    };
    // Basic styling, adjust as needed
    const containerStyle = {
        position: "relative",
        width: "100%", // Example: make it full width
        maxWidth: "360px", // Example: constrain max width
        border: "1px solid #e0e0e0", // Example border
        borderRadius: "8px", // Example rounding
        overflow: "hidden", // Important for panel transitions if any
        backgroundColor: "#ffffff", // Example background
    };
    return (
        <div style={containerStyle}>
            {path.length === 0 ? (
                <div style={{ padding: '20px', textAlign: 'center', color: '#666' }}>
                    {/*   Use translation */}
                    {t('noRootCategory')}
                </div>
            ) : (
                // Render the active panel. It receives the current category and its parent.
                <DrillDownCategoryPanel
                    key={activeCategory?.id} // Add key for potential transition effects
                    category={activeCategory}
                    parent={parentCategory}
                    onBack={handleBack} // Pass handler for back action
                    onSelectCategory={handleSelectCategory} // Pass handler for selecting a subcategory
                    // onFinalSelect={onFinalSelect} // Pass final select if needed by panel
                />
            )}
        </div>
    );
}, (prevProps, nextProps) => {
    // Only re-render if rootCategory actually changed
    return (
        prevProps.rootCategory?.id === nextProps.rootCategory?.id &&
        prevProps.rootCategory?.name === nextProps.rootCategory?.name
    );
});
export default DrillDownCategory;
// Basic PropTypes (adjust as needed)
// DrillDownCategory.propTypes = {
//    rootCategory: PropTypes.object, // Define shape if possible
//    onFinalSelect: PropTypes.func,
// };