// TitlePriceIconsRow.jsx
import { FaCheckCircle } from '../../utils/iconImports';
import React from 'react';
import PropTypes from 'prop-types';
import styled from 'styled-components';
const colorTextStrong = '#1C1C1C';
const colorCondition = '#4B9EA0';
const spacing = {
    xs: '4px',
    sm: '6px',
    md: '8px',
    lg: '12px',
    xl: '16px',
};
const breakpoints = {
    mobile: '480px',
    tablet: '768px',
};
// A small set of typical font sizes
const fontSize = {
    xs: '0.75rem',    // ~12px
    sm: '0.875rem',   // ~14px
    md: '1rem',       // ~16px
    lg: '1.125rem',   // ~18px
};
/* --------------------------------------------
   COMPONENT
--------------------------------------------- */
function TitleRowEditable({
                      condition,
                      productName,
                  }) {
    return (
        <Row>
            <TitleColumn>
                <ProductTitle>
                    <MinimalBadgeCondition>
                        <FaCheckCircle className="badge-icon"/>
                        {condition}
                    </MinimalBadgeCondition>
                    {productName}
                </ProductTitle>
            </TitleColumn>
        </Row>
    );
}
TitleRowEditable.propTypes = {
    id: PropTypes.oneOfType([PropTypes.string, PropTypes.number]).isRequired,
    condition: PropTypes.string.isRequired,
    productName: PropTypes.string.isRequired,
    priceCutPercent: PropTypes.number,
    productPrice: PropTypes.string.isRequired,
    canUseMiddleman: PropTypes.bool,
    handleAddToCart: PropTypes.func,
    handleOpenOffersForm: PropTypes.func,
};
TitleRowEditable.defaultProps = {
    priceCutPercent: 0,
    canUseMiddleman: false,
};
export default TitleRowEditable;
const Row = styled.div`
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: ${spacing.md};
    padding: ${spacing.md} 0;
    @media (max-width: ${breakpoints.tablet}) {
        gap: ${spacing.sm};
        padding: ${spacing.sm} 0;
    }
    @media (max-width: ${breakpoints.mobile}) {
        gap: ${spacing.xs};
        padding: ${spacing.xs} 0;
    }
`;
const MinimalBadgeCondition = styled.span`
    display: inline-flex;
    align-items: center;
    gap: ${spacing.xs};
    font-size: ${fontSize.xs};
    font-weight: 600;
    text-transform: uppercase;
    color: ${colorCondition};
    margin: ${spacing.xs};
    .badge-icon {
        font-size: 0.8rem;
    }
    @media (max-width: ${breakpoints.tablet}) {
        font-size: 0.7rem;
        .badge-icon {
            font-size: 0.7rem;
        }
    }
    @media (max-width: ${breakpoints.mobile}) {
        font-size: 0.65rem;
        .badge-icon {
            font-size: 0.65rem;
        }
    }
`;
const TitleColumn = styled.div`
    flex: 1;
    min-width: 60%;
    @media (max-width: ${breakpoints.mobile}) {
        min-width: auto;
    }
`;
const ProductTitle = styled.h2`
    font-size: ${fontSize.md};
    font-weight: 600;
    color: ${colorTextStrong};
    line-height: 1.3;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    word-break: break-word;
    @media (max-width: ${breakpoints.tablet}) {
        font-size: ${fontSize.sm};
    }
    @media (max-width: ${breakpoints.mobile}) {
        font-size: 0.8rem;
    }
`;