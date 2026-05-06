module.exports = {
  extends: [
    "stylelint-config-standard",
    "stylelint-config-recommended",
  ],
  plugins: [],
  rules: {
    /* Disallow transition: all for performance */
    "declaration-property-value-disallowed-list": {
      "transition": [/\ball\b/i]
    },
    /* Prefer design-token shadows instead of rgba box-shadow values */
    "property-value-pattern": {
      "box-shadow": [/rgba\(/]
    }
  }
}; 