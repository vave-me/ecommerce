import { redirect } from 'next/navigation';
export default function HomePage() {
  // This is a server component, so we can redirect on the server
  redirect('/');
  // This return is just for TypeScript - it will never be reached
  return null;
} 