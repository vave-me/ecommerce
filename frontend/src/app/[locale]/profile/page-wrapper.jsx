import ProfilePage from "./page";
import ErrorBoundary from "../../../components/ErrorBoundary";

export default function ProfilePageWrapper() {
    return (
        <ErrorBoundary 
            name="ProfilePage" 
            fallback={
                <div style={{ padding: '2rem', textAlign: 'center' }}>
                    <h2>Profile Error</h2>
                    <p>Unable to load profile. Please try again.</p>
                    <button onClick={() => window.location.reload()}>Refresh Page</button>
                </div>
            }
        >
            <ProfilePage />
        </ErrorBoundary>
    );
}