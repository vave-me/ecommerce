// TextEditor/components/MenuBar/ButtonGroup.jsx
import React from 'react';
import PropTypes from 'prop-types';
import styles from '../../TextEditor.module.css';
export function ButtonGroup({ children }) {
    return (
        <div className={styles.buttonGroup}>
            {children}
        </div>
    );
}
ButtonGroup.propTypes = {
    children: PropTypes.node.isRequired
};