// ProductInfoEditable.jsx
import { FiFilter } from '../../utils/iconImports';
import React, {useState} from 'react';
import PropTypes from 'prop-types';
import styled from 'styled-components';
/* ---------------------------------------------
   DESIGN TOKENS (compact spacing & sizes)
---------------------------------------------- */
const colorAccent = '#007bc7';
const colorAccentHover = '#0067a8';
const colorTextBase = '#333';
const colorTextStrong = '#1C1C1C';
const colorModalBackdrop = 'rgba(0, 0, 0, 0.45)';
const colorModalBG = '#fff';
const breakpoints = {
    mobile: '768px',
};
// Reduced spacing:
const spacing = {
    xs: '2px',
    sm: '4px',
    md: '6px',
    lg: '8px',
    xl: '12px',
};
// Slightly smaller font sizes to save space:
const fontSize = {
    xs: '0.7rem',   // ~11.2px
    sm: '0.8rem',   // ~12.8px
    md: '0.9rem',   // ~14.4px
    lg: '1rem',     // ~16px
    xl: '1.125rem', // ~18px
};
/* ---------------------------------------------
   Category Link
---------------------------------------------- */
function CategoryLinkRow({category}) {
    const handleCategoryClick = (e) => {
        e.preventDefault();
    };
    return (
        <StyledCategoryLink
            href="#"
            onClick={handleCategoryClick}
            aria-label={`View ${category} category`}
        >
            {category}
        </StyledCategoryLink>
    );
}
CategoryLinkRow.propTypes = {
    category: PropTypes.string.isRequired,
};
/* ---------------------------------------------
   Product Configuration Button & Modal
---------------------------------------------- */
function ProductConfigurationButton({
                                        shippingOptions,
                                        selectedShipping,
                                        onShippingChange,
                                        hasVariants,
                                        variants,
                                        selectedVariantId,
                                        onVariantChange,
                                        quantityOptions,
                                        selectedQuantity,
                                        onQuantityChange,
                                        productOptions,
                                        selectedProductOptions,
                                        onProductOptionChange,
                                    }) {
    const [open, setOpen] = useState(false);
    const [tempShipping, setTempShipping] = useState(selectedShipping || '');
    const [tempVariant, setTempVariant] = useState(selectedVariantId || '');
    const [tempQuantity, setTempQuantity] = useState(selectedQuantity || 1);
    const [tempProductOpts, setTempProductOpts] = useState({
        ...selectedProductOptions,
    });
    const showModal = () => {
        setTempShipping(selectedShipping || '');
        setTempVariant(selectedVariantId || '');
        setTempQuantity(selectedQuantity || 1);
        setTempProductOpts({...selectedProductOptions});
        setOpen(true);
    };
    const closeModal = () => setOpen(false);
    const handleSave = () => {
        onShippingChange(tempShipping);
        onVariantChange(tempVariant);
        onQuantityChange(tempQuantity);
        onProductOptionChange(tempProductOpts);
        closeModal();
    };
    const handleToggleProductOpt = (optName) => {
        setTempProductOpts((prev) => ({
            ...prev,
            [optName]: !prev[optName],
        }));
    };
    return (
        <>
            <ConfigButton onClick={showModal} aria-label="Configure product variants/options">
                <FiFilter/>
                Variants &amp; Options
            </ConfigButton>
            {open && (
                <ModalBackdrop onClick={closeModal} aria-hidden={!open}>
                    <ModalContent
                        onClick={(e) => e.stopPropagation()}
                        role="dialog"
                        aria-modal="true"
                        aria-label="Product configuration"
                    >
                        <ModalTitle>Choose Your Preferences</ModalTitle>
                        {hasVariants && (
                            <>
                                {!!shippingOptions?.length && (
                                    <OptionBlock>
                                        <OptionLabel htmlFor="shippingSelect">Shipping</OptionLabel>
                                        <SmallSelect
                                            id="shippingSelect"
                                            value={tempShipping}
                                            onChange={(e) => setTempShipping(e.target.value)}
                                        >
                                            {shippingOptions.map((opt) => (
                                                <option key={opt.code} value={opt.code}>
                                                    {`${opt.label} (${opt.cost}€)`}
                                                </option>
                                            ))}
                                        </SmallSelect>
                                    </OptionBlock>
                                )}
                                <OptionBlock>
                                    <OptionLabel htmlFor="variantSelect">Variant</OptionLabel>
                                    <SmallSelect
                                        id="variantSelect"
                                        value={tempVariant}
                                        onChange={(e) => setTempVariant(e.target.value)}
                                    >
                                        {variants.map((v) => {
                                            const detailStr = v.details?.length
                                                ? ` (${v.details.join(', ')})`
                                                : '';
                                            return (
                                                <option key={v.id} value={v.id}>
                                                    {`${v.name}${detailStr}`}
                                                </option>
                                            );
                                        })}
                                    </SmallSelect>
                                </OptionBlock>
                                {!!quantityOptions?.length && (
                                    <OptionBlock>
                                        <OptionLabel htmlFor="quantitySelect">Quantity</OptionLabel>
                                        <SmallSelect
                                            id="quantitySelect"
                                            value={tempQuantity}
                                            onChange={(e) => setTempQuantity(e.target.value)}
                                        >
                                            {quantityOptions.map((qty) => (
                                                <option key={qty} value={qty}>
                                                    {qty}
                                                </option>
                                            ))}
                                        </SmallSelect>
                                    </OptionBlock>
                                )}
                                {!!productOptions?.length && (
                                    <OptionBlock>
                                        <OptionLabel>Additional Options</OptionLabel>
                                        {productOptions.map((opt) => {
                                            const isChecked = !!tempProductOpts[opt.name];
                                            const labelText = opt.extraCost
                                                ? `${opt.name} (+${opt.extraCost}€)`
                                                : opt.name;
                                            return (
                                                <OptionCheckboxRow key={opt.name}>
                                                    <input
                                                        type="checkbox"
                                                        checked={isChecked}
                                                        onChange={() => handleToggleProductOpt(opt.name)}
                                                        aria-label={`Toggle ${opt.name}`}
                                                    />
                                                    {labelText}
                                                </OptionCheckboxRow>
                                            );
                                        })}
                                    </OptionBlock>
                                )}
                            </>
                        )}
                        <ModalFooter>
                            <ConfigButton type="button" onClick={closeModal} aria-label="Cancel configuration">
                                Cancel
                            </ConfigButton>
                            <ConfigButton type="button" onClick={handleSave} aria-label="Save configuration">
                                Save
                            </ConfigButton>
                        </ModalFooter>
                    </ModalContent>
                </ModalBackdrop>
            )}
        </>
    );
}
ProductConfigurationButton.propTypes = {
    shippingOptions: PropTypes.array,
    selectedShipping: PropTypes.string,
    onShippingChange: PropTypes.func.isRequired,
    hasVariants: PropTypes.bool,
    variants: PropTypes.array,
    selectedVariantId: PropTypes.string,
    onVariantChange: PropTypes.func.isRequired,
    quantityOptions: PropTypes.array,
    selectedQuantity: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
    onQuantityChange: PropTypes.func.isRequired,
    productOptions: PropTypes.array,
    selectedProductOptions: PropTypes.object,
    onProductOptionChange: PropTypes.func.isRequired,
};
ProductConfigurationButton.defaultProps = {
    shippingOptions: [],
    selectedShipping: '',
    hasVariants: false,
    variants: [],
    selectedVariantId: '',
    quantityOptions: [1, 2, 3, 4],
    selectedQuantity: 1,
    productOptions: [],
    selectedProductOptions: {},
};
/* ---------------------------------------------
   Main ProductInfoEditable Component
---------------------------------------------- */
function ProductInfoEditable({
                         category,
                         shippingOptions,
                         selectedShipping,
                         onShippingChange,
                         hasVariants,
                         variants,
                         selectedVariantId,
                         onVariantChange,
                         quantityOptions,
                         selectedQuantity,
                         onQuantityChange,
                         productOptions,
                         selectedProductOptions,
                         onProductOptionChange,
                     }) {
    const showConfigButton =
        (shippingOptions && shippingOptions.length > 0) ||
        (hasVariants && variants.length > 0) ||
        (productOptions && productOptions.length > 0);
    return (
        <RowWrapper>
            {/* Category link */}
            <CategoryLinkRow category={category}/>
            {/* Conditionally render config button */}
            {showConfigButton && (
                <ProductConfigurationButton
                    shippingOptions={shippingOptions}
                    selectedShipping={selectedShipping}
                    onShippingChange={onShippingChange}
                    hasVariants={hasVariants}
                    variants={variants}
                    selectedVariantId={selectedVariantId}
                    onVariantChange={onVariantChange}
                    quantityOptions={quantityOptions}
                    selectedQuantity={selectedQuantity}
                    onQuantityChange={onQuantityChange}
                    productOptions={productOptions}
                    selectedProductOptions={selectedProductOptions}
                    onProductOptionChange={onProductOptionChange}
                />
            )}
        </RowWrapper>
    );
}
ProductInfoEditable.propTypes = {
    category: PropTypes.string.isRequired,
    shippingOptions: PropTypes.array,
    selectedShipping: PropTypes.string,
    onShippingChange: PropTypes.func.isRequired,
    hasVariants: PropTypes.bool,
    variants: PropTypes.array,
    selectedVariantId: PropTypes.string,
    onVariantChange: PropTypes.func.isRequired,
    quantityOptions: PropTypes.array,
    selectedQuantity: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
    onQuantityChange: PropTypes.func.isRequired,
    productOptions: PropTypes.array,
    selectedProductOptions: PropTypes.object,
    onProductOptionChange: PropTypes.func.isRequired,
};
ProductInfoEditable.defaultProps = {
    shippingOptions: [],
    selectedShipping: '',
    hasVariants: false,
    variants: [],
    selectedVariantId: '',
    quantityOptions: [1, 2, 3, 4],
    selectedQuantity: 1,
    productOptions: [],
    selectedProductOptions: {},
};
export default ProductInfoEditable;
/* ---------------------------------------------
   STYLED COMPONENTS (Reduced spacing, line-clamp for Title)
---------------------------------------------- */
const RowWrapper = styled.div`
    display: flex;
    align-items: center;
    gap: ${spacing.sm};
    @media (max-width: ${breakpoints.mobile}) {
        gap: ${spacing.xs};
    }
`;
const StyledCategoryLink = styled.a`
    width: 50%;
    font-size: ${fontSize.md};
    font-weight: 400;
    color: ${colorAccent};
    text-decoration: none;
    &:hover {
        text-decoration: underline;
        color: ${colorAccentHover};
    }
    &:focus {
        outline: 2px solid ${colorAccent};
        outline-offset: 2px;
    }
    @media (max-width: ${breakpoints.mobile}) {
        font-size: ${fontSize.sm};
    }
`;
const ConfigButton = styled.button`
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: ${spacing.xs};
    padding: ${spacing.xs} ${spacing.sm};
    border: none;
    border-radius: 4px;
    background-color: ${colorAccent};
    color: #fff;
    font-weight: 600;
    font-size: ${fontSize.md};
    cursor: pointer;
    transition: background-color 0.2s ease, transform 0.1s ease;
    &:hover {
        background-color: ${colorAccentHover};
    }
    &:active {
        transform: scale(0.97);
        background-color: ${colorAccentHover};
    }
    &:focus {
        outline: 2px solid #fff;
        outline-offset: 2px;
    }
    svg {
        font-size: 1rem;
    }
    @media (max-width: ${breakpoints.mobile}) {
        padding: ${spacing.xs} ${spacing.xs};
        font-size: ${fontSize.sm};
        svg {
            font-size: 0.9rem;
        }
    }
`;
const ModalBackdrop = styled.div`
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: ${colorModalBackdrop};
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 9999;
`;
const ModalContent = styled.div`
    background: ${colorModalBG};
    padding: ${spacing.sm} ${spacing.md};
    border-radius: 6px;
    max-width: 420px;
    width: 90%;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
    display: flex;
    flex-direction: column;
    gap: ${spacing.sm};
    @media (max-width: ${breakpoints.mobile}) {
        padding: ${spacing.xs} ${spacing.sm};
        gap: ${spacing.xs};
    }
`;
const ModalTitle = styled.h2`
    margin: 0;
    font-size: ${fontSize.lg};
    color: ${colorTextBase};
    font-weight: 600;
    /* Force the title to two lines max */
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    line-height: 1.2;
    overflow: hidden;
    text-overflow: ellipsis;
    max-height: calc(2 * 1.2em); /* ensures 2 lines remain visible */
    @media (max-width: ${breakpoints.mobile}) {
        font-size: ${fontSize.md};
    }
`;
const OptionBlock = styled.div`
    display: flex;
    flex-direction: column;
    gap: ${spacing.xs};
    @media (max-width: ${breakpoints.mobile}) {
        gap: 2px;
    }
`;
const OptionLabel = styled.label`
    font-size: ${fontSize.sm};
    font-weight: 600;
    color: ${colorTextStrong};
    @media (max-width: ${breakpoints.mobile}) {
        font-size: ${fontSize.xs};
    }
`;
const SmallSelect = styled.select`
    font-size: ${fontSize.sm};
    padding: 3px 5px;
    background: #fff;
    color: ${colorTextBase};
    cursor: pointer;
    &:focus {
        outline: 2px solid ${colorAccent};
        outline-offset: 2px;
    }
    @media (max-width: ${breakpoints.mobile}) {
        font-size: ${fontSize.xs};
        padding: 2px 4px;
    }
`;
const OptionCheckboxRow = styled.label`
    display: flex;
    align-items: center;
    font-size: ${fontSize.sm};
    gap: ${spacing.xs};
    cursor: pointer;
    color: ${colorTextBase};
    input {
        cursor: pointer;
        &:focus {
            outline: 2px solid ${colorAccent};
            outline-offset: 2px;
        }
    }
    @media (max-width: ${breakpoints.mobile}) {
        font-size: ${fontSize.xs};
        gap: 2px;
    }
`;
const ModalFooter = styled.div`
    display: flex;
    justify-content: flex-end;
    gap: ${spacing.sm};
    @media (max-width: ${breakpoints.mobile}) {
        gap: ${spacing.xs};
    }
`;
