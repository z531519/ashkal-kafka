

module.exports = {
  testEnvironment: "node",
  roots: ['<rootDir>/tests'],
  testMatch: ['**/tests/integration/**/*.spec.ts'],
  transform: {
    '^.+\\.ts?$': 'ts-jest'
  },
  preset: 'ts-jest',
  coveragePathIgnorePatterns: [
    '^.+\\.mocks\.ts?$'
  ],
  setupFilesAfterEnv: ['./jest.setup.js'],
  
};
