import UserMetrics from './UserMetrics.client';

export const metadata = {
  title: 'User Metrics - Admin',
  description: 'View user activity metrics and analytics',
};

export default function UserMetricsPage({ params: { locale, userId } }) {
  return <UserMetrics locale={locale} userId={userId} />;
}