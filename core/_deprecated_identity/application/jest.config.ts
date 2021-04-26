/*
 * For a detailed explanation regarding each configuration property and type check, visit:
 * https://jestjs.io/docs/en/configuration.html
 */

export default {
  coverageDirectory: "coverage",
  // An array of file extensions your modules use
  moduleFileExtensions: ["js", "json", "jsx", "ts", "tsx", "node"],
  // A list of paths to directories that Jest should use to search for files in
  roots: ["<rootDir>"],
  // The test environment that will be used for testing
  testEnvironment: "node",
  // The regexp pattern or array of patterns that Jest uses to detect test files
  testRegex: "(/__tests__/.*|(\\.|/)(test|spec))\\.tsx?$",
  // A map from regular expressions to paths to transformers
  transform: {
    "^.+\\.tsx?$": "ts-jest",
  },
};
