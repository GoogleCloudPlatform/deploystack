# Changelog

## [1.9.0](https://github.com/GoogleCloudPlatform/deploystack/compare/config/v1.8.0...config/v1.9.0) (2023-03-07)


### Features

* enabling better manipulation of product info ([2f08319](https://github.com/GoogleCloudPlatform/deploystack/commit/2f083198a2c89429a9ba7adc4911f36faef934ab))
* moving config marshaling to the config package ([895da3e](https://github.com/GoogleCloudPlatform/deploystack/commit/895da3e0a0b42c6b27dbde647762490715f1e8de))


### Bug Fixes

* getting rid of unnecessary warnings ([772b53a](https://github.com/GoogleCloudPlatform/deploystack/commit/772b53a4125a4665d0467e2e76ddfd75cd885323))
* improving algorithm for finding main.tf ([cce3479](https://github.com/GoogleCloudPlatform/deploystack/commit/cce34793e09b805d0f79e125f1ec5a7de9c05af9))
* made path finding stuff work out of the box ([607e740](https://github.com/GoogleCloudPlatform/deploystack/commit/607e74096debdb5805f983b6add54fbec1cfef52))
* making terraform finder not panic if there is no tf files. ([92ed7d4](https://github.com/GoogleCloudPlatform/deploystack/commit/92ed7d498810b929eb1bd05b43d0f4541efad230))
* removed package dependencies on switching wd ([5240198](https://github.com/GoogleCloudPlatform/deploystack/commit/5240198f6ca467e8991e68686631395a2ceaac0c))

## [1.8.0](https://github.com/GoogleCloudPlatform/deploystack/compare/config/v1.5.1...config/v1.8.0) (2023-03-03)


### Features

* added authorsettings to replace hardsettings ([129a427](https://github.com/GoogleCloudPlatform/deploystack/commit/129a4272e8013693f630159f3bf5751275f1b60e))
* added authorsettings to replace hardsettings ([553979b](https://github.com/GoogleCloudPlatform/deploystack/commit/553979bf7520655c8467ce4ea4c6817f06447a4f))
* added multi stack detection to config operations. ([29b4ab5](https://github.com/GoogleCloudPlatform/deploystack/commit/29b4ab54ca9791ef00854b63fe5dddf5aeca90e4))
* added multi stack detection to config operations. ([a81267f](https://github.com/GoogleCloudPlatform/deploystack/commit/a81267ffa4abe07afae1f3942f6bc96e636b44e1))
* added new type of settings under the covers ([5eff3f6](https://github.com/GoogleCloudPlatform/deploystack/commit/5eff3f67686e71fb7ae09443341c8fc937024802))
* added new type of settings under the covers ([a21db57](https://github.com/GoogleCloudPlatform/deploystack/commit/a21db574eddf1ba49e727faf9171abbd14d47147))
* adding package config ([48c944b](https://github.com/GoogleCloudPlatform/deploystack/commit/48c944b97c8263e3032d92c90e9cc02f0cf7efe2))
* created new package to remove threat of circular dependency ([8b0738e](https://github.com/GoogleCloudPlatform/deploystack/commit/8b0738e28a839d2f9a21cb4c880ddd382d9017e2))


### Bug Fixes

* added more to the new setting features to work with tui ([ea6e1e5](https://github.com/GoogleCloudPlatform/deploystack/commit/ea6e1e534b1194a4c33d54c06dc4a5d529ab0a35))
* added more to the new setting features to work with tui ([0ca0104](https://github.com/GoogleCloudPlatform/deploystack/commit/0ca010434226093cfaa85f85b1f463b945c569e2))
* adding aliasof field to submodule ([cfb22ca](https://github.com/GoogleCloudPlatform/deploystack/commit/cfb22ca959e729c7b5e71d0d80d8db551ac0f6d9))
* adding test for init ([b2cce60](https://github.com/GoogleCloudPlatform/deploystack/commit/b2cce60d41cae75e779206738dba80e77c46fcfa))
* adding test for init ([a034b4d](https://github.com/GoogleCloudPlatform/deploystack/commit/a034b4d20bc6b715d3eca37a3aa8cb126e7d9a2e))
* made settings compatible with calling packages ([8eae051](https://github.com/GoogleCloudPlatform/deploystack/commit/8eae051360fc4b3620919f00e0d6d92ceeea07d2))
* made settings compatible with calling packages ([9131b4c](https://github.com/GoogleCloudPlatform/deploystack/commit/9131b4c8db61f0b35312060bbba6dd7c494054a6))
* tweaked find terraform behavior to be smarter ([96cdb5a](https://github.com/GoogleCloudPlatform/deploystack/commit/96cdb5a953fe2251d3b958675532a3e61ddc183b))
* tweaked find terraform behavior to be smarter ([73389e8](https://github.com/GoogleCloudPlatform/deploystack/commit/73389e8bdaa5962c1e0549329783eed0181cf496))

## [1.5.1](https://github.com/GoogleCloudPlatform/deploystack/compare/config/v1.5.0...config/v1.5.1) (2023-03-01)


### Bug Fixes

* tweaked find terraform behavior to be smarter ([1086322](https://github.com/GoogleCloudPlatform/deploystack/commit/1086322feaa93d71ee95d01e03f4ecb69771f1b7))

## [1.5.0](https://github.com/GoogleCloudPlatform/deploystack/compare/config/v1.4.0...config/v1.5.0) (2023-02-28)


### Features

* added multi stack detection to config operations. ([a81267f](https://github.com/GoogleCloudPlatform/deploystack/commit/a81267ffa4abe07afae1f3942f6bc96e636b44e1))

## [1.4.0](https://github.com/GoogleCloudPlatform/deploystack/compare/config/v1.3.0...config/v1.4.0) (2023-02-25)


### Features

* added authorsettings to replace hardsettings ([553979b](https://github.com/GoogleCloudPlatform/deploystack/commit/553979bf7520655c8467ce4ea4c6817f06447a4f))


### Bug Fixes

* made settings compatible with calling packages ([9131b4c](https://github.com/GoogleCloudPlatform/deploystack/commit/9131b4c8db61f0b35312060bbba6dd7c494054a6))

## [1.3.0](https://github.com/GoogleCloudPlatform/deploystack/compare/config/v1.2.1...config/v1.3.0) (2023-02-24)


### Features

* added new type of settings under the covers ([a21db57](https://github.com/GoogleCloudPlatform/deploystack/commit/a21db574eddf1ba49e727faf9171abbd14d47147))


### Bug Fixes

* added more to the new setting features to work with tui ([0ca0104](https://github.com/GoogleCloudPlatform/deploystack/commit/0ca010434226093cfaa85f85b1f463b945c569e2))

## [1.2.1](https://github.com/GoogleCloudPlatform/deploystack/compare/config/v1.2.0...config/v1.2.1) (2023-02-22)


### Bug Fixes

* adding test for init ([a034b4d](https://github.com/GoogleCloudPlatform/deploystack/commit/a034b4d20bc6b715d3eca37a3aa8cb126e7d9a2e))

## [1.2.0](https://github.com/GoogleCloudPlatform/deploystack/compare/config-v1.1.1...config/v1.2.0) (2023-02-21)


### Features

* adding package config ([48c944b](https://github.com/GoogleCloudPlatform/deploystack/commit/48c944b97c8263e3032d92c90e9cc02f0cf7efe2))
* created new package to remove threat of circular dependency ([8b0738e](https://github.com/GoogleCloudPlatform/deploystack/commit/8b0738e28a839d2f9a21cb4c880ddd382d9017e2))


### Bug Fixes

* adding aliasof field to submodule ([cfb22ca](https://github.com/GoogleCloudPlatform/deploystack/commit/cfb22ca959e729c7b5e71d0d80d8db551ac0f6d9))
* creating single source of truth for tf and repo data. ([91bd664](https://github.com/GoogleCloudPlatform/deploystack/commit/91bd664aa07b6b055eea63172278b4817890c605))
* forcing DVAs to be unique. ([dc9dfde](https://github.com/GoogleCloudPlatform/deploystack/commit/dc9dfde56f3cb6e5818621f6f3786b7ac4639387))
