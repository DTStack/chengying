const appConfig = require('./mock/appConfig');
const VERSION = JSON.stringify(require("./package.json").version); // app version.

module.exports = {
  globals: {
    'APPCONFIG': appConfig,
    APP: {
      VERSION: VERSION
    }
  },
  roots: [
    "<rootDir>/src"
  ],
  transform: {
    "^.+\\.[t|j]sx?$": "babel-jest",
    ".+\\.(css|styl|less|sass|scss)$": "jest-transform-css"
  },
  // testRegex: "(/__tests__/.*|(\\.|/)+(test|spec))\\.tsx?$",
  testMatch: [
    '**/__tests__/**/(*.)+(spec|test).[jt]s?(x)',
    '**/test/**/(*.)+(spec|test).[jt]s?(x)'
  ],
  transformIgnorePatterns: [`node_modules`],// Ignore modules without dt-common dir
  // testPathIgnorePatterns: ['/node_modules/'],
  snapshotSerializers: ["enzyme-to-json/serializer"],
  setupFilesAfterEnv: ["<rootDir>/src/setupEnzyme.ts"],
  moduleFileExtensions: [
    "ts",
    "tsx",
    "js",
    "jsx",
    "json",
    "node"
  ],
  moduleNameMapper: {
    '\\.(jpg|jpeg|png|gif|eot|otf|webp|svg|ttf|woff|woff2|mp4|webm|wav|mp3|m4a|aac|oga)$': '<rootDir>/mock/fileMock.ts',
    '\\.(css|scss|less)$': '<rootDir>/mock/styleMock.ts',
    '^@/(.*)$': '<rootDir>/src/$1'
  },
  testResultsProcessor: "jest-sonar-reporter"
}
