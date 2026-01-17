# Changelog

## [0.7.1](https://github.com/kubri/kubri/compare/v0.7.0...v0.7.1) (2026-01-17)


### Performance Improvements

* use IndexByte instead of Index for single byte searches ([#824](https://github.com/kubri/kubri/issues/824)) ([1cd2280](https://github.com/kubri/kubri/commit/1cd228070ffacd4b6736f3043cbf0655b249ee56))

## [0.7.0](https://github.com/kubri/kubri/compare/v0.6.0...v0.7.0) (2025-04-10)


### Features

* **apk:** support riscv64 & loongarch64 repositories ([#643](https://github.com/kubri/kubri/issues/643)) ([fa9599b](https://github.com/kubri/kubri/commit/fa9599b88f21579e7b3c80a500fa9f50f471497c))
* **apt:** use better compression for indices ([#640](https://github.com/kubri/kubri/issues/640)) ([da08ba2](https://github.com/kubri/kubri/commit/da08ba244334b67d6989c692e325591cd839e19a))
* **arch:** arch linux repository builder ([#638](https://github.com/kubri/kubri/issues/638)) ([a2d261a](https://github.com/kubri/kubri/commit/a2d261a0a7c33fa6010e05f7fc92f23a3fd185e1))
* **keys:** export keys, allow ed25519 PEM import ([#622](https://github.com/kubri/kubri/issues/622)) ([fcc610f](https://github.com/kubri/kubri/commit/fcc610f8f8a3cdccf9d2bf3ab9211cb62d37b90c))
* **rsa:** increase key size to 4096 ([#651](https://github.com/kubri/kubri/issues/651)) ([5dfc2b2](https://github.com/kubri/kubri/commit/5dfc2b2d2cd4104851715f3eb5fcbfc572b41238))
* **sparkle:** support ed25519 signatures for windows ([#642](https://github.com/kubri/kubri/issues/642)) ([ff19ae1](https://github.com/kubri/kubri/commit/ff19ae1bf168a9366d11a16e6af36f6b0c10619d))


### Performance Improvements

* **deb,desc:** improve float encode performance ([#657](https://github.com/kubri/kubri/issues/657)) ([c2fe830](https://github.com/kubri/kubri/commit/c2fe83060bbab66208813433dcaf4e9a5beb5ff2))

## [0.6.0](https://github.com/kubri/kubri/compare/v0.5.2...v0.6.0) (2025-03-17)


### ⚠ BREAKING CHANGES

* **s3:** require full URL for endpoint option, remove disable-ssl option ([#514](https://github.com/kubri/kubri/issues/514))

### Features

* **s3:** require full URL for endpoint option, remove disable-ssl option ([#514](https://github.com/kubri/kubri/issues/514)) ([9fbe359](https://github.com/kubri/kubri/commit/9fbe359e7d9c1931021ba1dd51cbe519146be24e))
* **sparkle:** support windows-arm64 ([#610](https://github.com/kubri/kubri/issues/610)) ([81bebc9](https://github.com/kubri/kubri/commit/81bebc9caa1aa37bc369b373379b3c0e11d76274))

## [0.5.2](https://github.com/kubri/kubri/compare/v0.5.1...v0.5.2) (2024-07-07)


### Bug Fixes

* **source/gitlab:** upload asset uses wrong path ([#442](https://github.com/kubri/kubri/issues/442)) ([1b0ccec](https://github.com/kubri/kubri/commit/1b0ccec0ae423aa50744bc471267bccde70b1afd))

## [0.5.1](https://github.com/kubri/kubri/compare/v0.5.0...v0.5.1) (2024-05-19)


### Bug Fixes

* **rpm:** panic when package has no group ([#411](https://github.com/kubri/kubri/issues/411)) ([60ef176](https://github.com/kubri/kubri/commit/60ef176a29406812744179c4000fb106af9ca29c))

## [0.5.0](https://github.com/kubri/kubri/compare/v0.4.0...v0.5.0) (2024-02-12)


### Features

* **apt:** upload gpg public key ([#336](https://github.com/kubri/kubri/issues/336)) ([6b30d3f](https://github.com/kubri/kubri/commit/6b30d3f53cd6756d9917f803461645e157f1aa55))
* rename to kubri ([#315](https://github.com/kubri/kubri/issues/315)) ([adac6a4](https://github.com/kubri/kubri/commit/adac6a40c20307baa12a7ec33737540c2cb73094))
* **target/github:** allow non-existent branch ([#312](https://github.com/kubri/kubri/issues/312)) ([7ecd098](https://github.com/kubri/kubri/commit/7ecd0987f9bb91a4bedce64e6e23eea5752e7a7e))


### Bug Fixes

* **apk:** key name missing .rsa.pub suffix ([#337](https://github.com/kubri/kubri/issues/337)) ([f33e1ff](https://github.com/kubri/kubri/commit/f33e1ffce29e293433af7fe948cdff4ab923e150))

## [0.4.0](https://github.com/kubri/kubri/compare/v0.3.0...v0.4.0) (2024-02-04)


### Features

* **apk:** apk repository builder ([#237](https://github.com/kubri/kubri/issues/237)) ([007d18d](https://github.com/kubri/kubri/commit/007d18d7f543d310cf7fe86b394d25e757f31473))
* **apk:** implement apk pipe ([#281](https://github.com/kubri/kubri/issues/281)) ([d8e5b38](https://github.com/kubri/kubri/commit/d8e5b38e1e32ebe916b872206de511f1085b60f0))
* **appinstaller:** group on-launch config under single parent ([#242](https://github.com/kubri/kubri/issues/242)) ([dc08d28](https://github.com/kubri/kubri/commit/dc08d28412446f8cb2bd64e7dcbca0203e9dc742))
* better error messages on failed integration ([#250](https://github.com/kubri/kubri/issues/250)) ([e08fde5](https://github.com/kubri/kubri/commit/e08fde5b180201b3b8127488695d5fb548557f93))
* better version constraints ([#249](https://github.com/kubri/kubri/issues/249)) ([2755b2c](https://github.com/kubri/kubri/commit/2755b2cfce1e47ddcccdfadbb03c72217cd0b5ba))
* **cli:** manage RSA keys ([#244](https://github.com/kubri/kubri/issues/244)) ([cb5de41](https://github.com/kubri/kubri/commit/cb5de41d7a2c623cca7b51c6166457275f5eb14e))
* generate jsonschema from config ([#287](https://github.com/kubri/kubri/issues/287)) ([fe5e907](https://github.com/kubri/kubri/commit/fe5e9070743e664e160acfa88c4abdfd4b3e9160))
* **source/gitlab:** allow setting custom URL (for self-hosted) ([#240](https://github.com/kubri/kubri/issues/240)) ([3ed31b4](https://github.com/kubri/kubri/commit/3ed31b4b66a12511b20737b258e232fefdb6239e))
* use raw base64 ed25519 keys ([#262](https://github.com/kubri/kubri/issues/262)) ([b18608a](https://github.com/kubri/kubri/commit/b18608aa6e8cbf3b6b6ea8fc445feed714fde7ea))
* validate config ([#290](https://github.com/kubri/kubri/issues/290)) ([fd12557](https://github.com/kubri/kubri/commit/fd125570f07107072adb659d3b2c8938eb3294c9))


### Bug Fixes

* **apk:** incorrect version constraint, handle target errors ([#305](https://github.com/kubri/kubri/issues/305)) ([8c5b11e](https://github.com/kubri/kubri/commit/8c5b11e1bf00c85ba01c07ebf94bad55bf5e715d))
* **appinstaller:** wrong extension for appxbundle ([#282](https://github.com/kubri/kubri/issues/282)) ([591f90d](https://github.com/kubri/kubri/commit/591f90d9c32abf9f14262e93a217e0981e2e1ad2))
* dirname validation ([#297](https://github.com/kubri/kubri/issues/297)) ([416b2df](https://github.com/kubri/kubri/commit/416b2df3650057c42b36ef9f1be391e580ad6db6))
* **pgp:** invalid signature on fedora 39 ([#227](https://github.com/kubri/kubri/issues/227)) ([9c2049f](https://github.com/kubri/kubri/commit/9c2049f03bd83478eff84d0611080484e54d3c40))
* **sparkle:** invalid ed signature ([#248](https://github.com/kubri/kubri/issues/248)) ([04e0ae6](https://github.com/kubri/kubri/commit/04e0ae64560d2b2aefa8e15fc07adc9e487a2221))

## [0.3.0](https://github.com/kubri/kubri/compare/v0.2.0...v0.3.0) (2023-11-24)


### Features

* allow overriding URL in blob targets ([#155](https://github.com/kubri/kubri/issues/155)) ([3156196](https://github.com/kubri/kubri/commit/315619652b9c3840a178e7da437a3ecb76cd8207))
* **apt:** custom compression, support lz4 ([#212](https://github.com/kubri/kubri/issues/212)) ([8f54a05](https://github.com/kubri/kubri/commit/8f54a0522e9bf6e298e0d07ad328e25270de4469))
* custom URLs in source & target, fix windows paths ([#202](https://github.com/kubri/kubri/issues/202)) ([e9c1c78](https://github.com/kubri/kubri/commit/e9c1c78bd38b731fd07a56a3a950a83b506e1c24))
* pgp sign repo metadata, simplify crypto packages ([#207](https://github.com/kubri/kubri/issues/207)) ([bc3ef36](https://github.com/kubri/kubri/commit/bc3ef366e666bb34834e022f97374a364089d357))
* validate version constraint ([#194](https://github.com/kubri/kubri/issues/194)) ([1e77b34](https://github.com/kubri/kubri/commit/1e77b34164a9744757249a08db793602b9d63ecc))
* yum integration ([#196](https://github.com/kubri/kubri/issues/196)) ([540b88c](https://github.com/kubri/kubri/commit/540b88ca52c79d29cd1d4878abef96ca0f053747))


### Bug Fixes

* **apt:** deb encoding bugs ([#213](https://github.com/kubri/kubri/issues/213)) ([6d564c4](https://github.com/kubri/kubri/commit/6d564c40aa184eeb354107377d81e44164a79d14))
* skip empty version constraint ([#198](https://github.com/kubri/kubri/issues/198)) ([cc9a300](https://github.com/kubri/kubri/commit/cc9a3006bc9ff057a3a73f32764510e6d25348a7))


### Performance Improvements

* **apt:** reduce allocs on encoding metadata ([#211](https://github.com/kubri/kubri/issues/211)) ([0b43413](https://github.com/kubri/kubri/commit/0b4341385e80578f85841a567b19262f214159a1))

## [0.2.0](https://github.com/kubri/kubri/compare/v0.1.0...v0.2.0) (2023-06-08)


### Features

* add targets, reuse data from target, remove flags ([#36](https://github.com/kubri/kubri/issues/36)) ([8fc1d64](https://github.com/kubri/kubri/commit/8fc1d646415f4fb82a74872f6af8bfff0667781d))
* build apt repository ([#40](https://github.com/kubri/kubri/issues/40)) ([2ed31c4](https://github.com/kubri/kubri/commit/2ed31c4a9d690296ccf62535405d779a2e937d29))
* github target ([#45](https://github.com/kubri/kubri/issues/45)) ([cad51f0](https://github.com/kubri/kubri/commit/cad51f090a595e64c4748a68582f48d98ea65484))
* publish appinstaller ([#52](https://github.com/kubri/kubri/issues/52)) ([d0be246](https://github.com/kubri/kubri/commit/d0be2462cd54118634ca3789a4ab7425736173cc))
* publish concurrently ([#105](https://github.com/kubri/kubri/issues/105)) ([7dcf359](https://github.com/kubri/kubri/commit/7dcf359e63697fab37ddf81ddda5210f618c35e4))
* skip integrations not in config ([#88](https://github.com/kubri/kubri/issues/88)) ([be63bc5](https://github.com/kubri/kubri/commit/be63bc5f379bda44896c9be3271f93147a8cee54))
* source & target as object, allow setting github target branch ([#48](https://github.com/kubri/kubri/issues/48)) ([5f378ae](https://github.com/kubri/kubri/commit/5f378aefff81d112efbc6324fa0cc3e0459d3959))
* **sparkle:** format description cdata ([#80](https://github.com/kubri/kubri/issues/80)) ([4cfd9a7](https://github.com/kubri/kubri/commit/4cfd9a773ad9c7cbd41c735864c1fce809f0611e))
* support tilde, caret & glob version constraints ([#50](https://github.com/kubri/kubri/issues/50)) ([6a29402](https://github.com/kubri/kubri/commit/6a29402d48ebc8234d68ba84bbb29ff3f7651fe6))
* upload packages, better version skipping ([#106](https://github.com/kubri/kubri/issues/106)) ([4095ed7](https://github.com/kubri/kubri/commit/4095ed734f37d3c5ae8ee2bcafaf82f298408c64))
* use contexts ([#29](https://github.com/kubri/kubri/issues/29)) ([b857de0](https://github.com/kubri/kubri/commit/b857de0fd6d89610a5967c8f03b357b60e26e1a7))
* version constraints include prerelease ([#49](https://github.com/kubri/kubri/issues/49)) ([c4400d4](https://github.com/kubri/kubri/commit/c4400d46a952d19683640e4838b63c05aa6c4cc6))


### Bug Fixes

* **apt:** handle unknown and empty control values ([#86](https://github.com/kubri/kubri/issues/86)) ([d364c8b](https://github.com/kubri/kubri/commit/d364c8bfc7cb68a337153457fd499b1e88bfdeee))
* **sparkle:** dsa sign ([#84](https://github.com/kubri/kubri/issues/84)) ([0edf934](https://github.com/kubri/kubri/commit/0edf934139bc7d122e58e2f80d4f7cbf330e2c61))


### Performance Improvements

* improve deb encoding ([#42](https://github.com/kubri/kubri/issues/42)) ([50eaf57](https://github.com/kubri/kubri/commit/50eaf57082d1a3bcc9542af2aae2dc9bd4991480))
