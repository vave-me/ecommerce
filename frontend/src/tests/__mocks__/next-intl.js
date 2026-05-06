// Mock for next-intl package
export const useTranslations = () => {
  return (key, params = {}) => {
    // Replace placeholders in format {placeholder}
    if (typeof key === 'string' && params) {
      let result = key;
      Object.entries(params).forEach(([param, value]) => {
        result = result.replace(`{${param}}`, value);
      });
      return result;
    }
    return key;
  };
};
export const useFormatter = () => ({
  dateTime: () => '01/01/2023',
  number: (num) => num.toString(),
  list: (list) => list.join(', '),
  relative: () => 'just now'
});
export default {
  useTranslations,
  useFormatter,
  NextIntlClientProvider: ({ children }) => children
}; 