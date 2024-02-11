/** @type {import('eslint').Linter.Config} */
module.exports = {
  extends: [
    'plugin:@typescript-eslint/recommended',
    'plugin:@typescript-eslint/stylistic',
    'plugin:@docusaurus/recommended',
    'airbnb',
    'airbnb/hooks',
    'airbnb-typescript',
    'plugin:react/jsx-runtime',
    'plugin:deprecation/recommended',
    'prettier',
  ],
  parser: '@typescript-eslint/parser',
  parserOptions: {
    project: './tsconfig.json',
    tsconfigRootDir: __dirname,
  },
  root: true,
  rules: {
    'react/jsx-props-no-spreading': 'off',
    'react/no-array-index-key': 'off',
  },
  overrides: [
    {
      files: ['docusaurus.config.ts'],
      rules: {
        'import/no-extraneous-dependencies': [2, { devDependencies: true }],
      },
    },
    {
      files: ['src/plugins/**/*'],
      rules: {
        'import/no-extraneous-dependencies': 'off',
      },
    },
  ],
}
