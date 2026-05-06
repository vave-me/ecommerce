/**
 * Next-intl Test Setup Utility
 * Provides real next-intl configuration for testing
 */

import { NextIntlClientProvider } from 'next-intl';
import React from 'react';

// Test messages for different locales
export const testMessages = {
  en: {
    HomePage: {
      title: 'Welcome',
      description: 'Test description',
      help: 'Help',
      contact: 'Contact',
      about: 'About',
      terms: 'Terms',
      privacy: 'Privacy',
      allRightsReserved: 'All rights reserved'
    },
    Seo: {
      title: 'sfx market – Live Marketplace',
      description: 'sfx markt is the live  marketplace that lets you buy, sell and connect with your community in real time.',
      keywords: 'sfx markt, marketplace, buy and sell locally, real‑time chat, SafePay'
    },
    Navigation: {
      home: 'Home',
      about: 'About',
      contact: 'Contact',
      login: 'Login',
      logout: 'Logout'
    },
    Common: {
      loading: 'Loading...',
      error: 'An error occurred',
      retry: 'Retry',
      cancel: 'Cancel',
      save: 'Save',
      delete: 'Delete',
      edit: 'Edit'
    }
  },
  pl: {
    HomePage: {
      title: 'Witamy',
      description: 'Opis testowy',
      help: 'Pomoc',
      contact: 'Kontakt',
      about: 'O nas',
      terms: 'Warunki',
      privacy: 'Prywatność',
      allRightsReserved: 'Wszystkie prawa zastrzeżone'
    },
    Seo: {
      title: 'sfx markt – Marketplace na Żywo',
      description: 'sfx markt to społecznościowy marketplace na żywo, który pozwala kupować, sprzedawać i łączyć się ze społecznością w czasie rzeczywistym.',
      keywords: 'sfx markt, marketplace, kupuj i sprzedawaj lokalnie, czat w czasie rzeczywistym, SafePay'
    },
    Navigation: {
      home: 'Strona główna',
      about: 'O nas',
      contact: 'Kontakt',
      login: 'Zaloguj',
      logout: 'Wyloguj'
    },
    Common: {
      loading: 'Ładowanie...',
      error: 'Wystąpił błąd',
      retry: 'Ponów',
      cancel: 'Anuluj',
      save: 'Zapisz',
      delete: 'Usuń',
      edit: 'Edytuj'
    }
  },
  de: {
    HomePage: {
      title: 'Willkommen',
      description: 'Test Beschreibung',
      help: 'Hilfe',
      contact: 'Kontakt',
      about: 'Über uns',
      terms: 'Bedingungen',
      privacy: 'Datenschutz',
      allRightsReserved: 'Alle Rechte vorbehalten'
    },
    Seo: {
      title: 'sfx markt – Live  Marktplatz',
      description: 'sfx markt ist der live  Marktplatz, der es Ihnen ermöglicht, in Echtzeit zu kaufen, zu verkaufen und sich mit Ihrer Gemeinschaft zu verbinden.',
      keywords: 'sfx markt, Marktplatz, lokal kaufen und verkaufen, Echtzeit-Chat, SafePay'
    },
    Navigation: {
      home: 'Startseite',
      about: 'Über uns',
      contact: 'Kontakt',
      login: 'Anmelden',
      logout: 'Abmelden'
    },
    Common: {
      loading: 'Laden...',
      error: 'Ein Fehler ist aufgetreten',
      retry: 'Wiederholen',
      cancel: 'Abbrechen',
      save: 'Speichern',
      delete: 'Löschen',
      edit: 'Bearbeiten'
    }
  }
};

/**
 * Test provider component that wraps components with NextIntlClientProvider
 */
export function NextIntlTestProvider({ 
  children, 
  locale = 'en', 
  messages = testMessages[locale] || testMessages.en,
  ...props 
}) {
  return (
    <NextIntlClientProvider 
      locale={locale} 
      messages={messages}
      timeZone="UTC"
      {...props}
    >
      {children}
    </NextIntlClientProvider>
  );
}

/**
 * Custom render function that includes NextIntlClientProvider
 */
export function renderWithNextIntl(
  ui,
  {
    locale = 'en',
    messages = testMessages[locale] || testMessages.en,
    ...renderOptions
  } = {}
) {
  const { render } = require('@testing-library/react');
  
  function Wrapper({ children }) {
    return (
      <NextIntlTestProvider locale={locale} messages={messages}>
        {children}
      </NextIntlTestProvider>
    );
  }

  return render(ui, { wrapper: Wrapper, ...renderOptions });
}

/**
 * Setup function for tests that need next-intl configuration
 */
export function setupNextIntlTest(locale = 'en') {
  const messages = testMessages[locale] || testMessages.en;
  
  return {
    locale,
    messages,
    provider: NextIntlTestProvider,
    render: (ui, options = {}) => renderWithNextIntl(ui, { locale, messages, ...options })
  };
} 