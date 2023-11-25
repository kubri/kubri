# Changelog

## [0.3.0](https://github.com/abemedia/appcast/compare/v0.2.0...v0.3.0) (2023-11-24)


### Features

* allow overriding URL in blob targets ([#155](https://github.com/abemedia/appcast/issues/155)) ([3156196](https://github.com/abemedia/appcast/commit/315619652b9c3840a178e7da437a3ecb76cd8207))
* **apt:** custom compression, support lz4 ([#212](https://github.com/abemedia/appcast/issues/212)) ([8f54a05](https://github.com/abemedia/appcast/commit/8f54a0522e9bf6e298e0d07ad328e25270de4469))
* custom URLs in source & target, fix windows paths ([#202](https://github.com/abemedia/appcast/issues/202)) ([e9c1c78](https://github.com/abemedia/appcast/commit/e9c1c78bd38b731fd07a56a3a950a83b506e1c24))
* pgp sign repo metadata, simplify crypto packages ([#207](https://github.com/abemedia/appcast/issues/207)) ([bc3ef36](https://github.com/abemedia/appcast/commit/bc3ef366e666bb34834e022f97374a364089d357))
* validate version constraint ([#194](https://github.com/abemedia/appcast/issues/194)) ([1e77b34](https://github.com/abemedia/appcast/commit/1e77b34164a9744757249a08db793602b9d63ecc))
* yum integration ([#196](https://github.com/abemedia/appcast/issues/196)) ([540b88c](https://github.com/abemedia/appcast/commit/540b88ca52c79d29cd1d4878abef96ca0f053747))


### Bug Fixes

* **apt:** deb encoding bugs ([#213](https://github.com/abemedia/appcast/issues/213)) ([6d564c4](https://github.com/abemedia/appcast/commit/6d564c40aa184eeb354107377d81e44164a79d14))
* skip empty version constraint ([#198](https://github.com/abemedia/appcast/issues/198)) ([cc9a300](https://github.com/abemedia/appcast/commit/cc9a3006bc9ff057a3a73f32764510e6d25348a7))


### Performance Improvements

* **apt:** reduce allocs on encoding metadata ([#211](https://github.com/abemedia/appcast/issues/211)) ([0b43413](https://github.com/abemedia/appcast/commit/0b4341385e80578f85841a567b19262f214159a1))

## [0.2.0](https://github.com/abemedia/appcast/compare/v0.1.0...v0.2.0) (2023-06-08)


### Features

* add targets, reuse data from target, remove flags ([#36](https://github.com/abemedia/appcast/issues/36)) ([8fc1d64](https://github.com/abemedia/appcast/commit/8fc1d646415f4fb82a74872f6af8bfff0667781d))
* build apt repository ([#40](https://github.com/abemedia/appcast/issues/40)) ([2ed31c4](https://github.com/abemedia/appcast/commit/2ed31c4a9d690296ccf62535405d779a2e937d29))
* github target ([#45](https://github.com/abemedia/appcast/issues/45)) ([cad51f0](https://github.com/abemedia/appcast/commit/cad51f090a595e64c4748a68582f48d98ea65484))
* publish appinstaller ([#52](https://github.com/abemedia/appcast/issues/52)) ([d0be246](https://github.com/abemedia/appcast/commit/d0be2462cd54118634ca3789a4ab7425736173cc))
* publish concurrently ([#105](https://github.com/abemedia/appcast/issues/105)) ([7dcf359](https://github.com/abemedia/appcast/commit/7dcf359e63697fab37ddf81ddda5210f618c35e4))
* skip integrations not in config ([#88](https://github.com/abemedia/appcast/issues/88)) ([be63bc5](https://github.com/abemedia/appcast/commit/be63bc5f379bda44896c9be3271f93147a8cee54))
* source & target as object, allow setting github target branch ([#48](https://github.com/abemedia/appcast/issues/48)) ([5f378ae](https://github.com/abemedia/appcast/commit/5f378aefff81d112efbc6324fa0cc3e0459d3959))
* **sparkle:** format description cdata ([#80](https://github.com/abemedia/appcast/issues/80)) ([4cfd9a7](https://github.com/abemedia/appcast/commit/4cfd9a773ad9c7cbd41c735864c1fce809f0611e))
* support tilde, caret & glob version constraints ([#50](https://github.com/abemedia/appcast/issues/50)) ([6a29402](https://github.com/abemedia/appcast/commit/6a29402d48ebc8234d68ba84bbb29ff3f7651fe6))
* upload packages, better version skipping ([#106](https://github.com/abemedia/appcast/issues/106)) ([4095ed7](https://github.com/abemedia/appcast/commit/4095ed734f37d3c5ae8ee2bcafaf82f298408c64))
* use contexts ([#29](https://github.com/abemedia/appcast/issues/29)) ([b857de0](https://github.com/abemedia/appcast/commit/b857de0fd6d89610a5967c8f03b357b60e26e1a7))
* version constraints include prerelease ([#49](https://github.com/abemedia/appcast/issues/49)) ([c4400d4](https://github.com/abemedia/appcast/commit/c4400d46a952d19683640e4838b63c05aa6c4cc6))


### Bug Fixes

* **apt:** handle unknown and empty control values ([#86](https://github.com/abemedia/appcast/issues/86)) ([d364c8b](https://github.com/abemedia/appcast/commit/d364c8bfc7cb68a337153457fd499b1e88bfdeee))
* **sparkle:** dsa sign ([#84](https://github.com/abemedia/appcast/issues/84)) ([0edf934](https://github.com/abemedia/appcast/commit/0edf934139bc7d122e58e2f80d4f7cbf330e2c61))


### Performance Improvements

* improve deb encoding ([#42](https://github.com/abemedia/appcast/issues/42)) ([50eaf57](https://github.com/abemedia/appcast/commit/50eaf57082d1a3bcc9542af2aae2dc9bd4991480))
