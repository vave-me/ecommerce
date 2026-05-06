// src/components/ListingManagement/Button.jsx
import { FaSpinner } from '../../utils/iconImports';
import React, { memo } from 'react';
import styled, { css } from 'styled-components';
import PropTypes from 'prop-types';
 // For loading indicator
const Button = memo(({
                    children,
                    onClick,
                    primary,
                    type,
                    disabled,
                    fullWidth,
                    loading,
                    icon,
                    ...rest
                }) => {
    return (
        <StyledButton
            onClick={onClick}
            primary={primary}
            type={type}
            disabled={disabled || loading}
            fullWidth={fullWidth}
            aria-disabled={disabled || loading}
            loading={loading}
            {...rest}
        >
            {loading ? (
                <SpinnerIcon />
            ) : (
                <>
                    {icon && <IconWrapper>{icon}</IconWrapper>}
                    <ButtonText>{children}</ButtonText>
                </>
            )}
        </StyledButton>
    );
});
Button.displayName = 'Button';
Button.propTypes = {
    children: PropTypes.node.isRequired,
    onClick: PropTypes.func,
    primary: PropTypes.bool,
    type: PropTypes.oneOf(['button', 'submit', 'reset']),
    disabled: PropTypes.bool,
    fullWidth: PropTypes.bool,
    loading: PropTypes.bool,
    icon: PropTypes.element, // Optional: To include an icon
};
Button.defaultProps = {
    onClick: () => {},
    primary: false,
    type: 'button',
    disabled: false,
    fullWidth: false,
    loading: false,
    icon: null,
};
export default Button;
/**
 * Styled Components
 */
const StyledButton = styled.button`
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 12px 24px;
  background-color: ${({ primary }) => (primary ? '#1aa89e' : '#f5f7fa')};
  color: ${({ primary }) => (primary ? '#ffffff' : '#2c3e50')};
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.3s ease, transform 0.2s ease,
    box-shadow 0.3s ease;
  width: ${({ fullWidth }) => (fullWidth ? '100%' : 'auto')};
  height: 48px;
  position: relative;
  overflow: hidden;
  &:hover {
    background-color: ${({ primary }) =>
    primary ? '#178f7a' : '#e1e8ed'};
    transform: translateY(-2px);
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  }
  &:active {
    transform: translateY(0);
    box-shadow: none;
  }
  &:focus {
    outline: none;
    box-shadow: 0 0 0 3px rgba(26, 168, 158, 0.5);
  }
  &:disabled,
  &[aria-disabled='true'] {
    background-color: #bdc3c7;
    color: #7f8c8d;
    cursor: not-allowed;
    box-shadow: none;
  }
  /* Loading State */
  ${({ loading }) =>
    loading &&
    css`
      cursor: not-allowed;
      pointer-events: none;
    `}
  /* Responsive Styling */
  @media (max-width: 768px) {
    padding: 10px 20px;
    font-size: 14px;
    height: 44px;
  }
  @media (max-width: 480px) {
    padding: 8px 16px;
    font-size: 13px;
    height: 40px;
  }
`;
const SpinnerIcon = styled(FaSpinner)`
    animation: spin 1s linear infinite;
    font-size: 20px;
    color: inherit;
    @keyframes spin {
        0% {
            transform: rotate(0deg);
        }
        100% {
            transform: rotate(360deg);
        }
    }
`;
const IconWrapper = styled.span`
    display: inline-flex;
    margin-right: 8px;
    font-size: 18px;
    @media (max-width: 768px) {
        font-size: 16px;
        margin-right: 6px;
    }
    @media (max-width: 480px) {
        font-size: 14px;
        margin-right: 4px;
    }
`;
const ButtonText = styled.span`
    display: inline-block;
    line-height: 1;
`;
