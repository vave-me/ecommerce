// File: src/components/Spinner.jsx
import React, {memo} from 'react';
import styles from './Spinner.module.css';
const Spinner = memo(() => (
    <div className={styles.spinnerContainer}>
        {/* First bouncing circle */}
        <div className={styles.bounce}/>
        {/* Second bouncing circle */}
        <div className={`${styles.bounce} ${styles.bounce2}`}/>
    </div>
));
Spinner.displayName = 'Spinner';
export default Spinner;
