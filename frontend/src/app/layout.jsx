"use client";
import './global.css';
// This root layout intentionally doesn't render any HTML structure
// The actual HTML/body will be rendered by the [locale]/layout.jsx
export default function RootLayout({ children }) {
  // Simply pass children to localized layouts
  return children;
}
