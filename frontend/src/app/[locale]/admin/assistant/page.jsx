import { Suspense } from 'react';
import AssistantManagement from './AssistantManagement.client';

export const metadata = {
  title: 'Assistant Management | Admin',
  description: 'Manage AI assistants and chatbots'
};

export default function AssistantPage() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <AssistantManagement />
    </Suspense>
  );
}