import ActivityMonitoring from './ActivityMonitoring.client';
import '../admin-theme.css';

export const metadata = {
  title: 'Activity Monitoring - Admin',
  description: 'Monitor user activity and platform interactions',
};

export default function ActivityMonitoringPage() {
  return <ActivityMonitoring />;
}