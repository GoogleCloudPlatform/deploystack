# Changelog

## [1.11.0](https://github.com/GoogleCloudPlatform/deploystack/compare/v1.10.0...v1.11.0) (2023-03-10)


### Features

* implemented New, with better ways of creating repos ([cc2b35a](https://github.com/GoogleCloudPlatform/deploystack/commit/cc2b35a602301f06c7bcf40d237e26b0d8b5ec83))

## [1.10.0](https://github.com/GoogleCloudPlatform/deploystack/compare/v1.9.0...v1.10.0) (2023-03-08)


### Features

* added suggestion to the dsexec ([b2a9a53](https://github.com/GoogleCloudPlatform/deploystack/commit/b2a9a53ba5021ba6fff76797d6fd9fe243ac9c58))


### Bug Fixes

* number of issues around mod updates ([a694472](https://github.com/GoogleCloudPlatform/deploystack/commit/a694472f68184030fad2f47bac32d12188921420))

## [1.9.0](https://github.com/GoogleCloudPlatform/deploystack/compare/v1.8.1...v1.9.0) (2023-03-07)


### Features

* adding ability to search terraform information. ([2f391f2](https://github.com/GoogleCloudPlatform/deploystack/commit/2f391f2d1cc84567c2a81a06064f56073ba4951a))
* adding features for trying repos without configs ([fb1a58c](https://github.com/GoogleCloudPlatform/deploystack/commit/fb1a58c1fcd7e0a9692a8a36ca1f0e1933563f67))
* adding suggested config capabilities to core ([c02bdfa](https://github.com/GoogleCloudPlatform/deploystack/commit/c02bdfa3b93b97add2c5f39d87eab045b20aef78))
* adding the ability to sort Blocks ([ef6b8a2](https://github.com/GoogleCloudPlatform/deploystack/commit/ef6b8a24bbaa7f5f12517b54386345747f18f4f5))
* enabling better manipulation of product info ([2f08319](https://github.com/GoogleCloudPlatform/deploystack/commit/2f083198a2c89429a9ba7adc4911f36faef934ab))
* got docker version shell script working perfectly ([ce9a125](https://github.com/GoogleCloudPlatform/deploystack/commit/ce9a125e3eda9b19eba43fc5ae9e80584a9f0eb4))
* moving config marshaling to the config package ([895da3e](https://github.com/GoogleCloudPlatform/deploystack/commit/895da3e0a0b42c6b27dbde647762490715f1e8de))


### Bug Fixes

* finding terraform now relies on main.tf files ([f5bf557](https://github.com/GoogleCloudPlatform/deploystack/commit/f5bf55777a92cd45eece41479862bd84748e13f7))
* getting rid of unnecessary warnings ([772b53a](https://github.com/GoogleCloudPlatform/deploystack/commit/772b53a4125a4665d0467e2e76ddfd75cd885323))
* improving algorithm for finding main.tf ([cce3479](https://github.com/GoogleCloudPlatform/deploystack/commit/cce34793e09b805d0f79e125f1ec5a7de9c05af9))
* made path finding stuff work out of the box ([607e740](https://github.com/GoogleCloudPlatform/deploystack/commit/607e74096debdb5805f983b6add54fbec1cfef52))
* making terraform finder not panic if there is no tf files. ([92ed7d4](https://github.com/GoogleCloudPlatform/deploystack/commit/92ed7d498810b929eb1bd05b43d0f4541efad230))
* refactoring broke some uses of code in places. ([ea3bd24](https://github.com/GoogleCloudPlatform/deploystack/commit/ea3bd24331b5d60ef0aedd1246516bceab8fcff2))
* removed package dependencies on switching wd ([5240198](https://github.com/GoogleCloudPlatform/deploystack/commit/5240198f6ca467e8991e68686631395a2ceaac0c))
* removing terminal clear. it's user hostile and hurts debugging ([f413228](https://github.com/GoogleCloudPlatform/deploystack/commit/f413228530b001beb0deb28a8b71d7f1c822a725))
* removing warnings for unrequired folders ([2d980c9](https://github.com/GoogleCloudPlatform/deploystack/commit/2d980c9b3b2531cb25408b43d9992475f6b7a308))

## [1.8.1](https://github.com/GoogleCloudPlatform/deploystack/compare/v1.8.0...v1.8.1) (2023-03-04)


### Bug Fixes

* polishing the deploystack shell script ([6ea6dd1](https://github.com/GoogleCloudPlatform/deploystack/commit/6ea6dd16161379c934a4c947fcc9259d38efb26e))
* testing caught a bug in billing caching ([880486e](https://github.com/GoogleCloudPlatform/deploystack/commit/880486e8671df47d5bfa55729c4c015a7691789f))
* trying to get releases correct ([b12504f](https://github.com/GoogleCloudPlatform/deploystack/commit/b12504f8ca86d41a77dea575295ad4a593a16b78))
* unpinned release please version. ([8a2f650](https://github.com/GoogleCloudPlatform/deploystack/commit/8a2f650580b845330e064dad95fadf5d9bfbf1fe))

## [1.8.0](https://github.com/GoogleCloudPlatform/deploystack/compare/v1.5.4...v1.8.0) (2023-03-03)


### ⚠ BREAKING CHANGES

* changing github to be simpler more general purpose

### Features

* add Dockerfile to create a Docker image ([d5a87db](https://github.com/GoogleCloudPlatform/deploystack/commit/d5a87db838b67a33065361ad06a227aa8dfaff3f))
* added a warning about slow lookups in a new project ([855ac53](https://github.com/GoogleCloudPlatform/deploystack/commit/855ac537cf7470ba4df3788e229790efdb3e7cbd))
* added ability for more than one stack to exist in the same repo ([4d7e928](https://github.com/GoogleCloudPlatform/deploystack/commit/4d7e9282d65f3d25cba7e13ee3e2c2184637628b))
* added ability for more than one stack to exist in the same repo ([959b797](https://github.com/GoogleCloudPlatform/deploystack/commit/959b797478114da0d8d08e336605645f3bd02e56))
* added authorsettings to replace hardsettings ([129a427](https://github.com/GoogleCloudPlatform/deploystack/commit/129a4272e8013693f630159f3bf5751275f1b60e))
* added authorsettings to replace hardsettings ([553979b](https://github.com/GoogleCloudPlatform/deploystack/commit/553979bf7520655c8467ce4ea4c6817f06447a4f))
* added caching to speed up these requests ([a0aaad8](https://github.com/GoogleCloudPlatform/deploystack/commit/a0aaad877b643181757d84e37621109f9e2f9d5a))
* added multi stack detection to config operations. ([29b4ab5](https://github.com/GoogleCloudPlatform/deploystack/commit/29b4ab54ca9791ef00854b63fe5dddf5aeca90e4))
* added multi stack detection to config operations. ([a81267f](https://github.com/GoogleCloudPlatform/deploystack/commit/a81267ffa4abe07afae1f3942f6bc96e636b44e1))
* added new type of settings under the covers ([5eff3f6](https://github.com/GoogleCloudPlatform/deploystack/commit/5eff3f67686e71fb7ae09443341c8fc937024802))
* added new type of settings under the covers ([a21db57](https://github.com/GoogleCloudPlatform/deploystack/commit/a21db574eddf1ba49e727faf9171abbd14d47147))
* added the progress bar. ([4e6a76c](https://github.com/GoogleCloudPlatform/deploystack/commit/4e6a76cead93920c619b6798b7caa3b743a58605))
* adding name to config and now searching out the value using git ([af07341](https://github.com/GoogleCloudPlatform/deploystack/commit/af073411ec2945e1151cb0913baa845c6c90c5c1))
* adding package config ([48c944b](https://github.com/GoogleCloudPlatform/deploystack/commit/48c944b97c8263e3032d92c90e9cc02f0cf7efe2))
* adding project retrieval to the mix ([12d6dcd](https://github.com/GoogleCloudPlatform/deploystack/commit/12d6dcd92e8f9b3b6a9bc8a9b987a40cb5998966))
* build a "don't put into settings flag" into the base configuration of pages ([d616a13](https://github.com/GoogleCloudPlatform/deploystack/commit/d616a138f29a10f56f7f14efb0636b34358f3621))
* changing github to be simpler more general purpose ([e0699fc](https://github.com/GoogleCloudPlatform/deploystack/commit/e0699fc6b7b0c74335bfe1c1cb1e5bd73e027c33))
* created new package to remove threat of circular dependency ([8b0738e](https://github.com/GoogleCloudPlatform/deploystack/commit/8b0738e28a839d2f9a21cb4c880ddd382d9017e2))
* forgot to add prependProject setting. Working now ([a6b2e25](https://github.com/GoogleCloudPlatform/deploystack/commit/a6b2e25c2980b459d715e82ba4ec5c1cbbdac735))
* implemented a back function ([c65badc](https://github.com/GoogleCloudPlatform/deploystack/commit/c65badc1f5a2041e725ba43015ff8af8438948f8))
* moving functionality for cloning github repos from GoogleCloud to main package ([a78eb4b](https://github.com/GoogleCloudPlatform/deploystack/commit/a78eb4b7568e2b33f7d5c2cbaad2b840c5dbde3e))
* moving functionality for cloning github repos from GoogleCloud to main package ([5da9d01](https://github.com/GoogleCloudPlatform/deploystack/commit/5da9d0130df6358a282a34043b1551823a87a409))
* moving functionality for cloning github repos from GoogleCloud to main package ([4de81c6](https://github.com/GoogleCloudPlatform/deploystack/commit/4de81c6eca04ac8841b926f8aaf821c2aea9a346))
* moving functionality for cloning github repos from GoogleCloud to main package ([7fc6d20](https://github.com/GoogleCloudPlatform/deploystack/commit/7fc6d20ef1f8c5a80bd1b63e936c4b51e347b5a1))
* moving functionality for cloning github repos from GoogleCloud to main package ([127e02d](https://github.com/GoogleCloudPlatform/deploystack/commit/127e02dd28c56aa79c7fc0c52c6fe7af481f54f4))
* moving functionality for cloning github repos from GoogleCloud to main package ([be01c77](https://github.com/GoogleCloudPlatform/deploystack/commit/be01c7761e187ec760bf84309ada1ec163f41f87))
* renamed package to better reflect purpose ([70850fc](https://github.com/GoogleCloudPlatform/deploystack/commit/70850fce76e7402261d9ae31177ffeaddec1722b))
* renamed package to make more clear its purpose ([d3d476e](https://github.com/GoogleCloudPlatform/deploystack/commit/d3d476ef2ae3da26910b3f44d6ce3edd5f07feab))
* updated tui to work with new configuration settings ([702e6c0](https://github.com/GoogleCloudPlatform/deploystack/commit/702e6c0798fd44bb4225773e3e59bf14fb78c523))
* updated tui to work with new configuration settings ([6a5a17a](https://github.com/GoogleCloudPlatform/deploystack/commit/6a5a17ac6fc097dcd563a5a5a0638972769b9e5d))


### Bug Fixes

* a number of the names of functions were pretty confusing for reuse. ([02c99cf](https://github.com/GoogleCloudPlatform/deploystack/commit/02c99cf898178ded80cbd394c0eb77bd9d75cd4c))
* absorbing some of the features we removed from github ([6257a47](https://github.com/GoogleCloudPlatform/deploystack/commit/6257a47e239d5a2a33c10fc2a3fd62761ef9a2c0))
* absorbing some of the features we removed from github ([3e9758e](https://github.com/GoogleCloudPlatform/deploystack/commit/3e9758e3fe14d51817103ef7314b20da43f73a61))
* actually fixing the problem with project numbers ([5175fea](https://github.com/GoogleCloudPlatform/deploystack/commit/5175fea1dedd89b5d5da9c4a8174e4c9c7f92ecb))
* add test for contact writing ([80f6679](https://github.com/GoogleCloudPlatform/deploystack/commit/80f66791fe86c160c056f241c99eab7335c1e670))
* add test for contact writing ([5927477](https://github.com/GoogleCloudPlatform/deploystack/commit/59274779d0c9b5dd6da2f00bb07c5ce2641948a0))
* added a test for iam ([b876d8a](https://github.com/GoogleCloudPlatform/deploystack/commit/b876d8a0fb64828efbd5db8667c7e961a44131ef))
* added an example for multi stack repos ([0ddf558](https://github.com/GoogleCloudPlatform/deploystack/commit/0ddf5588e30f1f6d87232d9144731b7813564ee0))
* added an example for multi stack repos ([52ca81f](https://github.com/GoogleCloudPlatform/deploystack/commit/52ca81f34895b6a9a5b0ecc6e3941d9155ffaa1a))
* added billing account capture and made sure that was working ([d7c0dc3](https://github.com/GoogleCloudPlatform/deploystack/commit/d7c0dc3c51699da801e7eaccbcca42b453170723))
* added Dockerfile for running DeployStack locally ([b1185b1](https://github.com/GoogleCloudPlatform/deploystack/commit/b1185b178ad7271b9acfa02052fd7c1661fb3342))
* added Dockerfile for running DeployStack locally ([2156a44](https://github.com/GoogleCloudPlatform/deploystack/commit/2156a44edd9a0c161751094708ea1a6c4386aa24))
* added more to the new setting features to work with tui ([ea6e1e5](https://github.com/GoogleCloudPlatform/deploystack/commit/ea6e1e534b1194a4c33d54c06dc4a5d529ab0a35))
* added more to the new setting features to work with tui ([0ca0104](https://github.com/GoogleCloudPlatform/deploystack/commit/0ca010434226093cfaa85f85b1f463b945c569e2))
* added test for checking for contacts ([e8679e7](https://github.com/GoogleCloudPlatform/deploystack/commit/e8679e738dca6a21f81ad82f2cee715e6b246f02))
* added test for checking for contacts ([5d9c940](https://github.com/GoogleCloudPlatform/deploystack/commit/5d9c940b277fba756bdb5f579f41c99a3e835a2a))
* added test for cloudbuild stuff ([2395694](https://github.com/GoogleCloudPlatform/deploystack/commit/2395694575c31f7c105db6b7d8c0766c93bc6693))
* adding a cache for multiple calls to the same reference to speed things up ([3906b1f](https://github.com/GoogleCloudPlatform/deploystack/commit/3906b1f2b061e197613df546f3aea4445666742b))
* adding a clean preprocessor for the last screen ([c3941d7](https://github.com/GoogleCloudPlatform/deploystack/commit/c3941d792be3663847ce42e06396f4da7806a8da))
* adding ability to check project existence ([5e4e6f8](https://github.com/GoogleCloudPlatform/deploystack/commit/5e4e6f89b183644940f857c1088f4291989d85e6))
* adding aliasof field to submodule ([cfb22ca](https://github.com/GoogleCloudPlatform/deploystack/commit/cfb22ca959e729c7b5e71d0d80d8db551ac0f6d9))
* adding test for Caching the domain contact ([26d0311](https://github.com/GoogleCloudPlatform/deploystack/commit/26d0311cec66a2f7e49e9d9ba9ae9d05201bdec8))
* adding test for Caching the domain contact ([8475386](https://github.com/GoogleCloudPlatform/deploystack/commit/847538682f3e029118ccbe7be306ae376bf6fb50))
* adding test for init ([b2cce60](https://github.com/GoogleCloudPlatform/deploystack/commit/b2cce60d41cae75e779206738dba80e77c46fcfa))
* adding test for init ([a034b4d](https://github.com/GoogleCloudPlatform/deploystack/commit/a034b4d20bc6b715d3eca37a3aa8cb126e7d9a2e))
* adding test for new package to CI/CD ([641cac8](https://github.com/GoogleCloudPlatform/deploystack/commit/641cac8680b3bcb0a241dea8174a7c8eaeb5f67b))
* adding test for secretmanager code. ([cbe4d04](https://github.com/GoogleCloudPlatform/deploystack/commit/cbe4d04abecc60147dec7287d23fa5f0bf5a0458))
* adjusting demo ([7f2ce4f](https://github.com/GoogleCloudPlatform/deploystack/commit/7f2ce4f1bde32525a52ddea776e6fa34c46a3e01))
* all requests for were failing because of bad projects ([35fdeb8](https://github.com/GoogleCloudPlatform/deploystack/commit/35fdeb865af96c9cf84bfecf999a4661e82a9dc4))
* allowing settings to be removed from stacks ([e1aa5bb](https://github.com/GoogleCloudPlatform/deploystack/commit/e1aa5bbe33452a13b10abf52b5416654b2c2ffcc))
* altering the verison to trigger a proper release ([e6126c0](https://github.com/GoogleCloudPlatform/deploystack/commit/e6126c0ace689218dc229bab5b60fcdf46b1966e))
* altering the verison to trigger a proper release ([e941f56](https://github.com/GoogleCloudPlatform/deploystack/commit/e941f563758b1f92658e42de32a6f51017cf4981))
* API changes caused test failure ([bcbda59](https://github.com/GoogleCloudPlatform/deploystack/commit/bcbda591755ac279b388ac81a4d1eccf84549fee))
* better default picker style ([4f962e5](https://github.com/GoogleCloudPlatform/deploystack/commit/4f962e53f04a886c4414bd8070c103a116987a7c))
* broken test fix. ([4df3414](https://github.com/GoogleCloudPlatform/deploystack/commit/4df34148b84e9958b19df01951ee66ff58da675d))
* broken test. ([fa29b23](https://github.com/GoogleCloudPlatform/deploystack/commit/fa29b23574589e40c8968d431427498a0ff9507e))
* broken tests ([f9de00d](https://github.com/GoogleCloudPlatform/deploystack/commit/f9de00d74ca2d86c23ce634288d99e91baf931b3))
* changing github to be simpler more general purpose ([cbad628](https://github.com/GoogleCloudPlatform/deploystack/commit/cbad62802e0a53bd1b12d505d99987f89283ecc5))
* changing to cloudshell to make sure this passes the test. ([5050a39](https://github.com/GoogleCloudPlatform/deploystack/commit/5050a3983b65a7c1e5f5d5b4cbbfb18000e3d29c))
* checking in todays changes ([18e3775](https://github.com/GoogleCloudPlatform/deploystack/commit/18e37755d37312bd26735c19dc163bd498dd2b5c))
* cloud shell always thinks it's in dark mode ([05a75e2](https://github.com/GoogleCloudPlatform/deploystack/commit/05a75e231771e76a654ad4bac771438dcafc8910))
* cloud shell always thinks it's in dark mode ([2fc66fd](https://github.com/GoogleCloudPlatform/deploystack/commit/2fc66fd214214226c5027b6e3ffa130aebe20152))
* commented out tests that were commented out before ([3755a5f](https://github.com/GoogleCloudPlatform/deploystack/commit/3755a5f5ca3c4f182e545112b4cc1ae619ded1b1))
* commenting out previously commented test ([a023d8a](https://github.com/GoogleCloudPlatform/deploystack/commit/a023d8a2115a7f6d83d52f989db175ba3433f6d2))
* commenting out tests that were commented out before ([d17cb02](https://github.com/GoogleCloudPlatform/deploystack/commit/d17cb02b8959e564b8d32fa9d78cff2dd7f535e2))
* converted spinner to use dsStyles to render color on Cloud Shell ([a7456b7](https://github.com/GoogleCloudPlatform/deploystack/commit/a7456b743a79585a481bd66242737ec31f7748d8))
* corrected error in test for project get ([8525be1](https://github.com/GoogleCloudPlatform/deploystack/commit/8525be1627891cbe535842e9d998186c98467e8e))
* corrected issue where domain contact info was not being refreshed properly ([546af35](https://github.com/GoogleCloudPlatform/deploystack/commit/546af355c1a795017accddd23271e3a429db3480))
* correcting billing check quality ([a830b43](https://github.com/GoogleCloudPlatform/deploystack/commit/a830b43c443dd6fea60b5ddd41ec3e096c7b9cd4))
* correcting broken test ([9775dff](https://github.com/GoogleCloudPlatform/deploystack/commit/9775dffb738836808574b2d640f9e536bec97462))
* correcting bug where Project forms looked bad. ([1a2e336](https://github.com/GoogleCloudPlatform/deploystack/commit/1a2e336cd7889d6b17d9d7050b7e5715496a4566))
* correcting dependency in test ([1431a27](https://github.com/GoogleCloudPlatform/deploystack/commit/1431a27bd75b21c7647cb6132f8355bce1c9872a))
* correcting failed test. ([bd2130b](https://github.com/GoogleCloudPlatform/deploystack/commit/bd2130b39466834a523a112a1624b7f6ff810a2f))
* correcting intermittent failures in test ([f1d8df0](https://github.com/GoogleCloudPlatform/deploystack/commit/f1d8df023a176c6e604c1baf0ef3687faabe1c1b))
* debugging might be screwing things up here. ([4abbc3e](https://github.com/GoogleCloudPlatform/deploystack/commit/4abbc3ed07847095f4970428d31ee6ac996e4ba8))
* debugging might be screwing things up here. ([177cc54](https://github.com/GoogleCloudPlatform/deploystack/commit/177cc54fa769e927e76814feb28f1fca3d70b6c3))
* dependecies out of date contributing errors ([472aca3](https://github.com/GoogleCloudPlatform/deploystack/commit/472aca384d29e48937f73a0faaadc371ad72f210))
* deploystack not breaking ([34faca4](https://github.com/GoogleCloudPlatform/deploystack/commit/34faca428f22648ef3ca94233ed7a6857be3fea4))
* deploystack not breaking ([220d98d](https://github.com/GoogleCloudPlatform/deploystack/commit/220d98d13c87f16c117a9c9ff161e37601ff4dbd))
* deploystack not breaking ([cac0e74](https://github.com/GoogleCloudPlatform/deploystack/commit/cac0e748b22c0317345a04bd84968208b4cd78ae))
* deploystack not breaking ([f936316](https://github.com/GoogleCloudPlatform/deploystack/commit/f9363163ce698714bcf707d9dacb8f5d34e3fad8))
* didn't like short form ([d4219c3](https://github.com/GoogleCloudPlatform/deploystack/commit/d4219c38c537dd4623f7c78a51fdd28855fe2dd8))
* don't throw errors if stuff isn't there log it. ([417addb](https://github.com/GoogleCloudPlatform/deploystack/commit/417addb923133f7d204eb5a535571260cf691a7f))
* don't throw errors if stuff isn't there log it. ([d73b15b](https://github.com/GoogleCloudPlatform/deploystack/commit/d73b15b1388379b7f14958455f03076db1a953e4))
* enabled multiple projects in a stack ([19b5603](https://github.com/GoogleCloudPlatform/deploystack/commit/19b5603edce2d822b1686db23b8f0c18b6ca491a))
* ensuring new components are tested ([92562b9](https://github.com/GoogleCloudPlatform/deploystack/commit/92562b92e498bcb5e3643471ef0034df21d29fcb))
* export terraform file information ([fedb9d9](https://github.com/GoogleCloudPlatform/deploystack/commit/fedb9d95b20d209b5759a36f7568dea8c81ae170))
* export terraform file information ([ac51436](https://github.com/GoogleCloudPlatform/deploystack/commit/ac514363d61ddc07f688ab9f69c0977b06981115))
* finalizing test design ([347aa25](https://github.com/GoogleCloudPlatform/deploystack/commit/347aa25bfd0da89b0124149f6ab73a24f43d2045))
* finding untested bits and fixing them. ([2f9896c](https://github.com/GoogleCloudPlatform/deploystack/commit/2f9896ca0da509bc8b0e2611dc706618c41b26dd))
* fix isn't taking ([0ac19c2](https://github.com/GoogleCloudPlatform/deploystack/commit/0ac19c288e8c93c62f89224db9ce00bdaa388814))
* fixed an issue where when the default value isn't set, the list isn't seleced ([6d69e80](https://github.com/GoogleCloudPlatform/deploystack/commit/6d69e80ccaef0d09c95859ddfe319e58c408cc2a))
* fixed example to use new ui ([99dc463](https://github.com/GoogleCloudPlatform/deploystack/commit/99dc463c34b7484ba4d84191674f569c430edca3))
* fixed issue where cloning with ssh caused errors ([5848091](https://github.com/GoogleCloudPlatform/deploystack/commit/5848091cbb0d1244f36feae761da5f6e1ed40ce7))
* fixes dealing with changes in color rendering ([26bd638](https://github.com/GoogleCloudPlatform/deploystack/commit/26bd63893c6cfc23929d4517c1457a4fa06698a3))
* fixing breaking test ([9784da5](https://github.com/GoogleCloudPlatform/deploystack/commit/9784da55dbb6ec641b04fcdd0bbb74c942029a48))
* fixing bug that caused project creation to fail ([303e61a](https://github.com/GoogleCloudPlatform/deploystack/commit/303e61a4f987ab07b1e0eaba667d3b520204e29b))
* fixing error where ctr-c on the exit page causes a panic ([3974a8b](https://github.com/GoogleCloudPlatform/deploystack/commit/3974a8bbf92387104e5f644d4054cf3414b5560a))
* fixing issues with testing billing code locally ([cbeb261](https://github.com/GoogleCloudPlatform/deploystack/commit/cbeb2618008bc2b048044f6d3554b905b30c9c09))
* fixing the demo app ([91f6211](https://github.com/GoogleCloudPlatform/deploystack/commit/91f6211ebde5b7c596cc447d937fc26ca5f1f5b2))
* getting all the version stuff working properly ([140556a](https://github.com/GoogleCloudPlatform/deploystack/commit/140556afab66aad5233032c402072641f561d223))
* getting billing to work properly ([9894e2b](https://github.com/GoogleCloudPlatform/deploystack/commit/9894e2bcea439b197755ab9a324dbbbe7ed71e09))
* getting docker version to build with configurable dependencies ([2cbc411](https://github.com/GoogleCloudPlatform/deploystack/commit/2cbc41119f8883e04e2991f94f07683fec3f3109))
* getting docker version to build with configurable dependencies ([1b70b39](https://github.com/GoogleCloudPlatform/deploystack/commit/1b70b397fcd600632d70c75c357c999dc2072b5c))
* getting repo version of executable working ([d938230](https://github.com/GoogleCloudPlatform/deploystack/commit/d938230f8b5b040afa505548d4b900c14c0813e2))
* getting repo version of executable working ([6fcec3b](https://github.com/GoogleCloudPlatform/deploystack/commit/6fcec3b4b5f76b4db05d2bf9e947a4dccf64db1d))
* getting rid of old calls. ([a391e39](https://github.com/GoogleCloudPlatform/deploystack/commit/a391e39aa119211e49570f4f1e6a67b804a87325))
* getting tests to pass ([c425b87](https://github.com/GoogleCloudPlatform/deploystack/commit/c425b875c35ea4a41546facefeb501daaf05fb70))
* got domain availability test to work ([7103da2](https://github.com/GoogleCloudPlatform/deploystack/commit/7103da2c83a1af3cb84ef4981a66ff93ef4506e6))
* got picker default values working properly ([ef623d7](https://github.com/GoogleCloudPlatform/deploystack/commit/ef623d748107ee0cb13e30082926bbeeb5015edb))
* initial implementation of test framework ([55be6be](https://github.com/GoogleCloudPlatform/deploystack/commit/55be6be840dd3ee220433cd896df79ae81a01db6))
* initial import of tui code ([e53bb85](https://github.com/GoogleCloudPlatform/deploystack/commit/e53bb85709de1bc8c025ce0948b6cec348cdaab6))
* initial move to move gcloud client stuff to own package ([22dd503](https://github.com/GoogleCloudPlatform/deploystack/commit/22dd50354c32014ec6b845cd30ebff8c33a64aa2))
* introducing tests for newly added features. ([272b41e](https://github.com/GoogleCloudPlatform/deploystack/commit/272b41e2cd3e80a115da693d032b30c4d3440295))
* issues with default control for machineType ([b7df60a](https://github.com/GoogleCloudPlatform/deploystack/commit/b7df60ab9230dd4329ac8c62777f77fc4f66b1a2))
* it sure is hard to call an ssh clone from within Cloud Build ([1eac57f](https://github.com/GoogleCloudPlatform/deploystack/commit/1eac57fcbe4dc77475230170ea03ef9268e80f24))
* just some tweaks to improving linting ([5c716d7](https://github.com/GoogleCloudPlatform/deploystack/commit/5c716d7d05f1e7235f9319edacd5578669c081ef))
* list items limited to 50 chars and defaults now display correctly ([5a068c4](https://github.com/GoogleCloudPlatform/deploystack/commit/5a068c4b3b7983aa6468a6c180378787cda72c15))
* list items limited to 50 chars and defaults now display correctly ([46b2e3b](https://github.com/GoogleCloudPlatform/deploystack/commit/46b2e3b8288c432e6ebff4dfefd756e9040dd38f))
* live project get is flakey. I should figure that out eventually. ([ce7cc1e](https://github.com/GoogleCloudPlatform/deploystack/commit/ce7cc1e91515e6cf59f4a472e661c2e01e09f50f))
* made currentProject a queue value instead of a global one ([c3bfb7b](https://github.com/GoogleCloudPlatform/deploystack/commit/c3bfb7bcbd46d781e84612dc4327204f93ba8897))
* made settings compatible with calling packages ([8eae051](https://github.com/GoogleCloudPlatform/deploystack/commit/8eae051360fc4b3620919f00e0d6d92ceeea07d2))
* made settings compatible with calling packages ([9131b4c](https://github.com/GoogleCloudPlatform/deploystack/commit/9131b4c8db61f0b35312060bbba6dd7c494054a6))
* made tui compatible with changes to config ([39a0104](https://github.com/GoogleCloudPlatform/deploystack/commit/39a010481404d8ef79a31718bc5c15e6d5f47ad1))
* made tui compatible with changes to config ([0ba76e4](https://github.com/GoogleCloudPlatform/deploystack/commit/0ba76e4dcf5270b9ac1cc71ceb96fcc395d5fa3b))
* make default options display normally if there are less than 10 of them. ([9326feb](https://github.com/GoogleCloudPlatform/deploystack/commit/9326febdc30ca5dcf8c6a92fdb0e4aa9db1be9cb))
* make default options display normally if there are less than 10 of them. ([0cc9aa6](https://github.com/GoogleCloudPlatform/deploystack/commit/0cc9aa60a1fd71e18c7ca31ed96a801780abfc51))
* make defaultUserAgent include the stack name ([57b52b3](https://github.com/GoogleCloudPlatform/deploystack/commit/57b52b3ffd9e676b960dd63cc55e1a8d7dc62617))
* make sure nil pointers don't break us. ([ee6c5d4](https://github.com/GoogleCloudPlatform/deploystack/commit/ee6c5d4274a37d63b2a17fcfedcec82c6fa9d437))
* make sure nil pointers don't break us. ([4a8835d](https://github.com/GoogleCloudPlatform/deploystack/commit/4a8835dee7baf5ecff2f67cf9fde62eed1efea95))
* make sure that this new feature doesn't block users without permission ([090ff20](https://github.com/GoogleCloudPlatform/deploystack/commit/090ff20ef281a422e40587af0d950f2cdd4d16a4))
* make test work locally ([d50f83b](https://github.com/GoogleCloudPlatform/deploystack/commit/d50f83bfaf420089314bd17ba61ab454c590331e))
* making all of the tests pass for the lass refactoring ([8a9db91](https://github.com/GoogleCloudPlatform/deploystack/commit/8a9db911185989c03c88d5a1971f018d952668d7))
* making storage tests work ([c82a5af](https://github.com/GoogleCloudPlatform/deploystack/commit/c82a5afb2bcc92ff80d0da825d8974bf50602afa))
* making sure contact doesn't get added back. ([97de87f](https://github.com/GoogleCloudPlatform/deploystack/commit/97de87f2aa69020ae8c173d64aff372bf6246b71))
* making sure tests work properly ([2d799e3](https://github.com/GoogleCloudPlatform/deploystack/commit/2d799e390ae6dffb66605c913907cb95f0b4a683))
* making sure that gcloud config project set is called ([4295f73](https://github.com/GoogleCloudPlatform/deploystack/commit/4295f739633f75e9360b76f97d662fc363cf73d1))
* making sure that project creation captures project number ([fb0d8d7](https://github.com/GoogleCloudPlatform/deploystack/commit/fb0d8d7cbb46a54c191e642846e39369ed128f77))
* making the idea of arguments work. ([a101a2c](https://github.com/GoogleCloudPlatform/deploystack/commit/a101a2ccadd804a7a6b88ee488586ecfba75a35d))
* making the progress bar work. ([11b0924](https://github.com/GoogleCloudPlatform/deploystack/commit/11b09240e043b7ce997346296bfc6e4379714ac4))
* making things more organized. ([2a69730](https://github.com/GoogleCloudPlatform/deploystack/commit/2a697304310a5ee0780b621d22854f74459bd109))
* massively changed the way color was rendered to make it simpler ([a5c652a](https://github.com/GoogleCloudPlatform/deploystack/commit/a5c652a9cdfc9a6820681ac734a7c1a51af8d300))
* merge resolution added a breaking change. ([34ca69f](https://github.com/GoogleCloudPlatform/deploystack/commit/34ca69f4e3558db2376fbb42fea90089dd516e9a))
* more bug fixes ([5255ec0](https://github.com/GoogleCloudPlatform/deploystack/commit/5255ec06461588decdc0f4413740701419feab25))
* more pruning ([f5066a9](https://github.com/GoogleCloudPlatform/deploystack/commit/f5066a9f3cb398e81882c4c4f3fee241d84c8875))
* more test data ([c671a18](https://github.com/GoogleCloudPlatform/deploystack/commit/c671a18f0df82401e1306fc9d7291a2a164e427c))
* more test data ([3caa800](https://github.com/GoogleCloudPlatform/deploystack/commit/3caa8006aefebb6db97798c03a5a97fd4e9b3278))
* more test data ([c92651c](https://github.com/GoogleCloudPlatform/deploystack/commit/c92651cc64225427d85800ff30548b788bb0813b))
* more test data ([d570e8f](https://github.com/GoogleCloudPlatform/deploystack/commit/d570e8fe749de218d0b6b19ab04c5f4af3afb0cb))
* more test data ([b480d13](https://github.com/GoogleCloudPlatform/deploystack/commit/b480d13769b7565c6dffe1ed25dbe12a3df17d48))
* more test data ([3346001](https://github.com/GoogleCloudPlatform/deploystack/commit/3346001a5c3cc03a1ea2878c55816b82b46472b1))
* moved from parsing description files to having.a datastructure for products ([0758eb7](https://github.com/GoogleCloudPlatform/deploystack/commit/0758eb70801974f67863ccc8aa4eadf6be2e1ec8))
* moved project number retrieval to the actual flow. ([a546144](https://github.com/GoogleCloudPlatform/deploystack/commit/a546144dc53b5c0a02a9918d14daa96dd22288fc))
* moving constants to the main deploystack project ([cc80588](https://github.com/GoogleCloudPlatform/deploystack/commit/cc805880581987370ec7c37bc550e88f6fcf66bc))
* need to find a permanent fix for this ([d522665](https://github.com/GoogleCloudPlatform/deploystack/commit/d522665f63f81608e79ed8e4f7a98d23f6da7f8e))
* needed to account for nils in casting operation ([4b0ea90](https://github.com/GoogleCloudPlatform/deploystack/commit/4b0ea90634dedde02b8ed34920639fdf341a944f))
* okay deploystack is changing a bit ([e301341](https://github.com/GoogleCloudPlatform/deploystack/commit/e3013418988afcd99620cd18db6307bbd4454245))
* okay deploystack is changing a bit ([1f735fb](https://github.com/GoogleCloudPlatform/deploystack/commit/1f735fb730695ecc6c62cf7ad6f10f619ff88bb9))
* okay deploystack is changing a bit ([629a26c](https://github.com/GoogleCloudPlatform/deploystack/commit/629a26c0ea6017da3ce502436057129d02c83a52))
* okay deploystack is changing a bit ([5d3f4e4](https://github.com/GoogleCloudPlatform/deploystack/commit/5d3f4e4cb46dfd60330f8c8cfddbb557003c491d))
* Ported over all of the gcloud functionality ([c380557](https://github.com/GoogleCloudPlatform/deploystack/commit/c38055786c572a6d90b6c8b24952f5a909f507c2))
* preventing error where canceling out of dsexec causes typing to stop showing ([42da885](https://github.com/GoogleCloudPlatform/deploystack/commit/42da885e44c2e5fdb2fe762ca3812c3c5cb0d9a2))
* preventing error where canceling out of dsexec causes typing to stop showing ([6cf41bc](https://github.com/GoogleCloudPlatform/deploystack/commit/6cf41bc01b494c5f2ca403c2b5a5a1109aa466b1))
* process is now exit(1) when user stops process ([ab4bbcf](https://github.com/GoogleCloudPlatform/deploystack/commit/ab4bbcf81370953a083b3224f01035401c9403e3))
* pruning directory files a little better. ([828d6ce](https://github.com/GoogleCloudPlatform/deploystack/commit/828d6ce19b80bd9f6463848008dea07b7a578eae))
* pull down repos ([3e6b3f5](https://github.com/GoogleCloudPlatform/deploystack/commit/3e6b3f55bb118caccbe5937d1723dc13f813ecf5))
* pull down repos ([d16c043](https://github.com/GoogleCloudPlatform/deploystack/commit/d16c043ddc321ec96f2ef9c282c0ae56e76d8cc7))
* pull down repos ([0eac901](https://github.com/GoogleCloudPlatform/deploystack/commit/0eac901a42d54fbba18941c1480d3273cf72b52d))
* pull down repos ([37e3d99](https://github.com/GoogleCloudPlatform/deploystack/commit/37e3d99d28f09ca0e8da2aaa4f4df95f27f33148))
* pull down repos ([344d5ef](https://github.com/GoogleCloudPlatform/deploystack/commit/344d5ef166b5bd66903b70bb04351d82503462ae))
* pull down repos ([47b66f2](https://github.com/GoogleCloudPlatform/deploystack/commit/47b66f27abf93ffa740693c1213f53c0f6f28ead))
* pulled over the rest of the functions to complete the interface for ui ([40d6741](https://github.com/GoogleCloudPlatform/deploystack/commit/40d674124d0681ded6c61d39b1906f758e43ed5a))
* reduced the size of lists to 10 ([ffe407e](https://github.com/GoogleCloudPlatform/deploystack/commit/ffe407e588ed20750ee9f3e103ce49d8c13e72d1))
* reduced the size of lists to 10 ([0182c5b](https://github.com/GoogleCloudPlatform/deploystack/commit/0182c5b992b6d359f80ab8e0340b798f17f3aa6a))
* refactoring labeledvalues control for comprehension and better testing ([38f2a1a](https://github.com/GoogleCloudPlatform/deploystack/commit/38f2a1ae46cab460050ab97080b7a08a1affeed7))
* removed a lot of unnecessary code ([315a948](https://github.com/GoogleCloudPlatform/deploystack/commit/315a948ca317de7f6634361fda93ea5146c36be8))
* removing a bunch of unused code ([786a7f7](https://github.com/GoogleCloudPlatform/deploystack/commit/786a7f7b4eec242ee0f495865ecb99fecfa4f6dc))
* removing breaking from root package ([e914f1b](https://github.com/GoogleCloudPlatform/deploystack/commit/e914f1bb83449ca6a6d03262a4e0a5017ed6809f))
* removing breaking from root package ([e35b25c](https://github.com/GoogleCloudPlatform/deploystack/commit/e35b25cf8e56ed6f664c206621df399ca8fd3c21))
* removing breaking from root package ([5d09eb8](https://github.com/GoogleCloudPlatform/deploystack/commit/5d09eb838671f78331699984d8365c132740ccc7))
* removing breaking from root package ([c72617e](https://github.com/GoogleCloudPlatform/deploystack/commit/c72617ed01d779077e85569daeda37276bdfb705))
* removing debugging content ([083dff0](https://github.com/GoogleCloudPlatform/deploystack/commit/083dff05c7ee051e9e4997c301a7c87c4d76d9a4))
* removing generated file ([52e73fe](https://github.com/GoogleCloudPlatform/deploystack/commit/52e73fecb76a45e092ccb2fe47466e2c0366a549))
* removing stack_name from the terraform variables file. ([819240b](https://github.com/GoogleCloudPlatform/deploystack/commit/819240b9da1e440cde8bf84236bcbcfc62b8f094))
* removing the need for defaultValue ([9c13e8c](https://github.com/GoogleCloudPlatform/deploystack/commit/9c13e8c55bd3f8b92b18740ac0b6acf3d512b3cb))
* renaming bits ([d52c9a2](https://github.com/GoogleCloudPlatform/deploystack/commit/d52c9a27f7287657b01c555e89252aaee10cef86))
* renaming bits ([d223e1f](https://github.com/GoogleCloudPlatform/deploystack/commit/d223e1f1d965cd69e75786bf664b55ecd972e31a))
* renaming bits ([d13e076](https://github.com/GoogleCloudPlatform/deploystack/commit/d13e076a851431fa528a2b7aef66ee76e77e6fe3))
* renaming bits ([019e900](https://github.com/GoogleCloudPlatform/deploystack/commit/019e9006d903f598fbe344b66d7aaf3b3e3ad84c))
* renaming bits ([9665a82](https://github.com/GoogleCloudPlatform/deploystack/commit/9665a8224826b746f8a9c06b3cf86eb637c69913))
* renaming bits ([1359491](https://github.com/GoogleCloudPlatform/deploystack/commit/1359491bbf2d0eae8c03fe9db056e23afcde5e9a))
* renaming services for clarity ([ad0351b](https://github.com/GoogleCloudPlatform/deploystack/commit/ad0351bee366b24e4dd7885b7f195ad276e8885b))
* renaming things to cut down on redundancy, long names and redundancy ([16455ef](https://github.com/GoogleCloudPlatform/deploystack/commit/16455ef6df7cc3b0d86fd914b3c2877bb431c542))
* resolves issue with color not showing in CloudShell ([df2a788](https://github.com/GoogleCloudPlatform/deploystack/commit/df2a788e2a4b7e9accb481e5b5c8c69fa8906b39))
* resolving issue where new accounts weren't getting billing attached ([fba7fad](https://github.com/GoogleCloudPlatform/deploystack/commit/fba7fadadeef1a7af510d4dfd58594a11d1e111e))
* restoreing full screen mode. ([3259419](https://github.com/GoogleCloudPlatform/deploystack/commit/3259419229d9abee0991a907350969def9ef5ede))
* reverting fix ([704d530](https://github.com/GoogleCloudPlatform/deploystack/commit/704d530ab4f19493bfa7d5b2a0699e86de198fce))
* reverting last change ([d16b18e](https://github.com/GoogleCloudPlatform/deploystack/commit/d16b18ec3e797adbf1c49b4605253ff026986ad7))
* reverting release-please config ([0eeddcf](https://github.com/GoogleCloudPlatform/deploystack/commit/0eeddcf2af22b23d4b7570a38014948478e16322))
* reverting release-please config ([9adfd7f](https://github.com/GoogleCloudPlatform/deploystack/commit/9adfd7ff7d352a92bd06065a875bc61af9c5c722))
* reverting reverting fix. ([8b95f62](https://github.com/GoogleCloudPlatform/deploystack/commit/8b95f6200f4afb92d0c2b10fca4ef01b34cc98d6))
* rolling back bad change. ([d1380c8](https://github.com/GoogleCloudPlatform/deploystack/commit/d1380c8e54280a9dd83d422876b40a4b71951ebb))
* setting up to be able to handle multiple stacks per repo ([a0cf51d](https://github.com/GoogleCloudPlatform/deploystack/commit/a0cf51d0ab1f3e1f84d76e0e0e7c0dfb949c3d30))
* setting up to be able to handle multiple stacks per repo ([6497186](https://github.com/GoogleCloudPlatform/deploystack/commit/6497186f83b61f1a9ce7c1fa90cf127f3c2f465c))
* ssh volume is no bueno. ([e730b7e](https://github.com/GoogleCloudPlatform/deploystack/commit/e730b7eaebcf67baf586adb9479249790418d23e))
* start using go workspaces to help ([f0e1f3f](https://github.com/GoogleCloudPlatform/deploystack/commit/f0e1f3f7352286e90bceaf3689c00bc7a1683b55))
* still trying to get build to pull down git via SSH ([d4c3117](https://github.com/GoogleCloudPlatform/deploystack/commit/d4c3117fe0e1e6ba41e88816192e9df6379cca0c))
* still trying to get please-release to work ([d0967bd](https://github.com/GoogleCloudPlatform/deploystack/commit/d0967bd529d66e68bd2d85e37b536967a230361e))
* stopped ui from showing non billing enabled projects incorrectly ([000f5ab](https://github.com/GoogleCloudPlatform/deploystack/commit/000f5abe716f79fe33ad6359ace7ed25a092d29c))
* Taskfile erred when there wasn't a script header. ([7e9a6b9](https://github.com/GoogleCloudPlatform/deploystack/commit/7e9a6b9ed404f65ae64b0f837bed1e39da71a71b))
* Taskfile erred when there wasn't a script header. ([5a11f31](https://github.com/GoogleCloudPlatform/deploystack/commit/5a11f3141a423eccb162fb92079285cab033bf7a))
* temporarily commenting out flaky test. ([5e8e2f3](https://github.com/GoogleCloudPlatform/deploystack/commit/5e8e2f3fa79e68e5b174f1003ce21272c8591abf))
* testing color rendering for Cloud Shell ([9535e12](https://github.com/GoogleCloudPlatform/deploystack/commit/9535e12b8a0aa1c828b152fa751b9265697b84e1))
* testing if the alt screen is the problem on cloudshell ([8ce98f9](https://github.com/GoogleCloudPlatform/deploystack/commit/8ce98f999cfe31b61165180e3e389df2185c2ee4))
* testing more of the code ([a74c38e](https://github.com/GoogleCloudPlatform/deploystack/commit/a74c38e419f4cf3803cda369c0c6ab751a197668))
* testing the behavior of the queue more ([99d76fb](https://github.com/GoogleCloudPlatform/deploystack/commit/99d76fbf16c7bdf1b6d3d265e85ac718baa0f9eb))
* tests broken because environment isn't right. ([e88931f](https://github.com/GoogleCloudPlatform/deploystack/commit/e88931f5dae41993069682fe3d19005f3b373e49))
* tests didn't get renamed ([bfa2066](https://github.com/GoogleCloudPlatform/deploystack/commit/bfa20663d6ab659e5fe1217c0dcc3dc0d3950d79))
* there are two expensive tests in this build, so upping the time and adding higher level machine ([bf62b97](https://github.com/GoogleCloudPlatform/deploystack/commit/bf62b97d1f00aa549a9ab6d76f0a9a92e9e23380))
* this is harder than it should be ([569fb52](https://github.com/GoogleCloudPlatform/deploystack/commit/569fb52ac6ca82b814c0513a4fd55b34b633f28c))
* this might fix ssh but break other tests. ([05c1c46](https://github.com/GoogleCloudPlatform/deploystack/commit/05c1c460977f9d61078a36b1c190524ebcb510c0))
* tightening things up to work with new package structure ([0577236](https://github.com/GoogleCloudPlatform/deploystack/commit/0577236e7c58b561477b0d6b142561d6abd6927d))
* trying to correct a problem where Deploystack stops the terminal from working ([ae49a99](https://github.com/GoogleCloudPlatform/deploystack/commit/ae49a99a211825ee4813a5440ad0c0b32ec40473))
* trying to correct a problem where Deploystack stops the terminal from working ([9b1d95e](https://github.com/GoogleCloudPlatform/deploystack/commit/9b1d95e1c57c3a9141280401518c74c14ed90811))
* trying to get all packages working properly ([5768e28](https://github.com/GoogleCloudPlatform/deploystack/commit/5768e2828529e17974feb8f1fe43feb706329f70))
* trying to get failing test working. ([6b6a544](https://github.com/GoogleCloudPlatform/deploystack/commit/6b6a5441cf58b6c9a4d5b85cbbb8f6556f3ebdaa))
* trying to get new verison working. ([0c1c084](https://github.com/GoogleCloudPlatform/deploystack/commit/0c1c084c686ecac53ec13bdf696e8617460c2582))
* trying to get new verison working. ([e6dbbd4](https://github.com/GoogleCloudPlatform/deploystack/commit/e6dbbd42bf14bf130c32d7950e570d586cfdd04e))
* trying to get new verison working. ([fac0c91](https://github.com/GoogleCloudPlatform/deploystack/commit/fac0c9159cf4fb044168223b50a0fc0b13767f4a))
* trying to get new verison working. ([8924c2b](https://github.com/GoogleCloudPlatform/deploystack/commit/8924c2b70128c7b12f71b71d35233c964ebc26b2))
* trying to get release-please to properly tag subpackages ([604e2d1](https://github.com/GoogleCloudPlatform/deploystack/commit/604e2d1e019071648d69fcbcf7359220119b9ea7))
* trying to update all the things to ensure proper dependencies ([e76fdae](https://github.com/GoogleCloudPlatform/deploystack/commit/e76fdae2123f6997268c7bc8a552e8f552e24a49))
* tweaked find terraform behavior to be smarter ([96cdb5a](https://github.com/GoogleCloudPlatform/deploystack/commit/96cdb5a953fe2251d3b958675532a3e61ddc183b))
* tweaked find terraform behavior to be smarter ([73389e8](https://github.com/GoogleCloudPlatform/deploystack/commit/73389e8bdaa5962c1e0549329783eed0181cf496))
* tweaking things for debugging ([9d62cf6](https://github.com/GoogleCloudPlatform/deploystack/commit/9d62cf642c5ab3718bad60a81a34253e869fca14))
* tweaking whitespace ([4b12003](https://github.com/GoogleCloudPlatform/deploystack/commit/4b12003b890966b8b4599325b2813202dc3d2004))
* typo in type - oh! ([7544d7d](https://github.com/GoogleCloudPlatform/deploystack/commit/7544d7dbc4997a26f7fa6b94e862643d8c9806b0))
* update dependencies ([d623a9c](https://github.com/GoogleCloudPlatform/deploystack/commit/d623a9c74183f4dc6dc1c7b7cf3a07083142d3d3))
* update dependencies ([cb39106](https://github.com/GoogleCloudPlatform/deploystack/commit/cb391060f1b703dd38dec3233479a22e0dbc1130))
* updated a lot of the capabilities of the test rig and added documentation ([b33e329](https://github.com/GoogleCloudPlatform/deploystack/commit/b33e3295e953afbceca164610feb4e54f55ea4a3))
* updated colors to try and make it look better on Cloud Shell ([5d7f220](https://github.com/GoogleCloudPlatform/deploystack/commit/5d7f2206b3445d9e94430dec61df0b2e7db90d8a))
* updated tests ([00b392e](https://github.com/GoogleCloudPlatform/deploystack/commit/00b392e927603c46dedfb3418578256f3d6eb4a1))
* updating errors to work with licensed files. ([49b0c66](https://github.com/GoogleCloudPlatform/deploystack/commit/49b0c66298af5dd9ad141bd115c59acdffcb216f))
* updating release-please config ([ce1505f](https://github.com/GoogleCloudPlatform/deploystack/commit/ce1505f6979f8c04a754f8084bc1b48d658f656f))
* updating tests ([a2c1617](https://github.com/GoogleCloudPlatform/deploystack/commit/a2c1617eab219bd2fe1b90ee5a69d1f79874d9bc))
* updating workspace settings ([3512aa3](https://github.com/GoogleCloudPlatform/deploystack/commit/3512aa34dd85c9fcc68cb030c8f797ad6e4edb3c))
* upgrade dependencies ([83b3ef0](https://github.com/GoogleCloudPlatform/deploystack/commit/83b3ef0ce7f07a5ed48ca5ac8d8bde88f083ffe3))
* upgrade dependencies ([9d80e3f](https://github.com/GoogleCloudPlatform/deploystack/commit/9d80e3fbb93328db775bc9651bc6b8bc3e045cac))
* when you are just asking for the billing account, you don't need to attach it. ([0f9f547](https://github.com/GoogleCloudPlatform/deploystack/commit/0f9f5478f79d883ec7f857937e31a5bbcb83d3ea))
* work! ([0b87451](https://github.com/GoogleCloudPlatform/deploystack/commit/0b8745196900632829822149c91f456e7d5aef58))
* working on a demo ui based on mock data ([cacabcc](https://github.com/GoogleCloudPlatform/deploystack/commit/cacabcce45d729eee47dee3ead038a485dc36e4e))
* working on better file searching ([f1d0439](https://github.com/GoogleCloudPlatform/deploystack/commit/f1d0439ee2518399aeed5e7a3d0f0c95f424d25a))

## [1.5.4](https://github.com/GoogleCloudPlatform/deploystack/compare/v1.5.3...v2.0.0) (2023-03-02)

### Features

* moving functionality for cloning github repos from GoogleCloud to main package ([a9ac5f5](https://github.com/GoogleCloudPlatform/deploystack/commit/a9ac5f58678038932e9361b9ac49d230e8695e82))


### Bug Fixes

* absorbing some of the features we removed from github ([32e4ddb](https://github.com/GoogleCloudPlatform/deploystack/commit/32e4ddb8f6ab05a7db1cd1a80d43ed63b8890ffe))
* added Dockerfile for running DeployStack locally ([c2d45cd](https://github.com/GoogleCloudPlatform/deploystack/commit/c2d45cda92e10988fa6553be39481bca90ba17e0))
* changing github to be simpler more general purpose ([3a2c058](https://github.com/GoogleCloudPlatform/deploystack/commit/3a2c058249fd32ec6f7d0b4451f57dbaa5249566))
* more test data ([a0eeac6](https://github.com/GoogleCloudPlatform/deploystack/commit/a0eeac68b50ae3de7213b166f421b27d50920363))
* pull down repos ([0b495a3](https://github.com/GoogleCloudPlatform/deploystack/commit/0b495a3980bde4f19e655fec8c02299cf76e70ce))
* renaming bits ([5d8b247](https://github.com/GoogleCloudPlatform/deploystack/commit/5d8b2472f1b16891c6dfb9c03bb6be9e61cf5e56))
* Taskfile erred when there wasn't a script header. ([ebf0ab1](https://github.com/GoogleCloudPlatform/deploystack/commit/ebf0ab1a80591aed8e4443030a5bb30e49cee28d))

## [1.5.3](https://github.com/GoogleCloudPlatform/deploystack/compare/v1.5.2...v1.5.3) (2023-03-02)


### Bug Fixes

* cloud shell always thinks it's in dark mode ([5f36947](https://github.com/GoogleCloudPlatform/deploystack/commit/5f369479376aa285ac9061151c1f22b8058ff5fa))

## [1.5.2](https://github.com/GoogleCloudPlatform/deploystack/compare/v1.5.1...v1.5.2) (2023-03-01)


### Bug Fixes

* tweaked find terraform behavior to be smarter ([1086322](https://github.com/GoogleCloudPlatform/deploystack/commit/1086322feaa93d71ee95d01e03f4ecb69771f1b7))

## [1.5.1](https://github.com/GoogleCloudPlatform/deploystack/compare/v1.5.0...v1.5.1) (2023-03-01)


### Bug Fixes

* don't throw errors if stuff isn't there log it. ([cd7da57](https://github.com/GoogleCloudPlatform/deploystack/commit/cd7da57db43625b95a49fca4603712ca94b4d170))
* export terraform file information ([289c988](https://github.com/GoogleCloudPlatform/deploystack/commit/289c988653b96d5b1085d232241c1578443539db))
* make sure nil pointers don't break us. ([ef3d435](https://github.com/GoogleCloudPlatform/deploystack/commit/ef3d435507e21fab2ad3ece6e4e11d0ecfd69aa3))
* setting up to be able to handle multiple stacks per repo ([5250188](https://github.com/GoogleCloudPlatform/deploystack/commit/525018825f7a8992b99a9f5f86c2045c367c56fd))

## [1.5.0](https://github.com/GoogleCloudPlatform/deploystack/compare/v1.4.1...v1.5.0) (2023-02-28)


### Features

* added ability for more than one stack to exist in the same repo ([959b797](https://github.com/GoogleCloudPlatform/deploystack/commit/959b797478114da0d8d08e336605645f3bd02e56))
* added multi stack detection to config operations. ([a81267f](https://github.com/GoogleCloudPlatform/deploystack/commit/a81267ffa4abe07afae1f3942f6bc96e636b44e1))


### Bug Fixes

* added an example for multi stack repos ([52ca81f](https://github.com/GoogleCloudPlatform/deploystack/commit/52ca81f34895b6a9a5b0ecc6e3941d9155ffaa1a))
* debugging might be screwing things up here. ([177cc54](https://github.com/GoogleCloudPlatform/deploystack/commit/177cc54fa769e927e76814feb28f1fca3d70b6c3))
* trying to correct a problem where Deploystack stops the terminal from working ([9b1d95e](https://github.com/GoogleCloudPlatform/deploystack/commit/9b1d95e1c57c3a9141280401518c74c14ed90811))

## [1.4.1](https://github.com/GoogleCloudPlatform/deploystack/compare/v1.4.0...v1.4.1) (2023-02-25)


### Bug Fixes

* upgrade dependencies ([9d80e3f](https://github.com/GoogleCloudPlatform/deploystack/commit/9d80e3fbb93328db775bc9651bc6b8bc3e045cac))

## [1.4.0](https://github.com/GoogleCloudPlatform/deploystack/compare/v1.3.0...v1.4.0) (2023-02-25)


### Features

* added authorsettings to replace hardsettings ([553979b](https://github.com/GoogleCloudPlatform/deploystack/commit/553979bf7520655c8467ce4ea4c6817f06447a4f))


### Bug Fixes

* made settings compatible with calling packages ([9131b4c](https://github.com/GoogleCloudPlatform/deploystack/commit/9131b4c8db61f0b35312060bbba6dd7c494054a6))
* made tui compatible with changes to config ([0ba76e4](https://github.com/GoogleCloudPlatform/deploystack/commit/0ba76e4dcf5270b9ac1cc71ceb96fcc395d5fa3b))

## [1.3.0](https://github.com/GoogleCloudPlatform/deploystack/compare/v1.2.3...v1.3.0) (2023-02-24)


### Features

* added new type of settings under the covers ([a21db57](https://github.com/GoogleCloudPlatform/deploystack/commit/a21db574eddf1ba49e727faf9171abbd14d47147))
* updated tui to work with new configuration settings ([6a5a17a](https://github.com/GoogleCloudPlatform/deploystack/commit/6a5a17ac6fc097dcd563a5a5a0638972769b9e5d))


### Bug Fixes

* added more to the new setting features to work with tui ([0ca0104](https://github.com/GoogleCloudPlatform/deploystack/commit/0ca010434226093cfaa85f85b1f463b945c569e2))
* list items limited to 50 chars and defaults now display correctly ([46b2e3b](https://github.com/GoogleCloudPlatform/deploystack/commit/46b2e3b8288c432e6ebff4dfefd756e9040dd38f))
* make default options display normally if there are less than 10 of them. ([0cc9aa6](https://github.com/GoogleCloudPlatform/deploystack/commit/0cc9aa60a1fd71e18c7ca31ed96a801780abfc51))
* reduced the size of lists to 10 ([0182c5b](https://github.com/GoogleCloudPlatform/deploystack/commit/0182c5b992b6d359f80ab8e0340b798f17f3aa6a))

## [1.2.3](https://github.com/GoogleCloudPlatform/deploystack/compare/v1.2.2...v1.2.3) (2023-02-22)


### Bug Fixes

* add test for contact writing ([5927477](https://github.com/GoogleCloudPlatform/deploystack/commit/59274779d0c9b5dd6da2f00bb07c5ce2641948a0))
* added test for checking for contacts ([5d9c940](https://github.com/GoogleCloudPlatform/deploystack/commit/5d9c940b277fba756bdb5f579f41c99a3e835a2a))
* adding test for Caching the domain contact ([8475386](https://github.com/GoogleCloudPlatform/deploystack/commit/847538682f3e029118ccbe7be306ae376bf6fb50))
* adding test for init ([a034b4d](https://github.com/GoogleCloudPlatform/deploystack/commit/a034b4d20bc6b715d3eca37a3aa8cb126e7d9a2e))

## [1.2.2](https://github.com/GoogleCloudPlatform/deploystack/compare/v1.2.1...v1.2.2) (2023-02-21)


### Bug Fixes

* reverting release-please config ([9adfd7f](https://github.com/GoogleCloudPlatform/deploystack/commit/9adfd7ff7d352a92bd06065a875bc61af9c5c722))

## [1.2.1](https://github.com/GoogleCloudPlatform/deploystack/compare/v1.2.0...v1.2.1) (2023-02-21)


### Bug Fixes

* update dependencies ([cb39106](https://github.com/GoogleCloudPlatform/deploystack/commit/cb391060f1b703dd38dec3233479a22e0dbc1130))

## [1.2.0](https://github.com/GoogleCloudPlatform/deploystack/compare/v1.1.4...v1.2.0) (2023-02-21)


### Features

* adding package config ([48c944b](https://github.com/GoogleCloudPlatform/deploystack/commit/48c944b97c8263e3032d92c90e9cc02f0cf7efe2))
* created new package to remove threat of circular dependency ([8b0738e](https://github.com/GoogleCloudPlatform/deploystack/commit/8b0738e28a839d2f9a21cb4c880ddd382d9017e2))


### Bug Fixes

* adding test for new package to CI/CD ([641cac8](https://github.com/GoogleCloudPlatform/deploystack/commit/641cac8680b3bcb0a241dea8174a7c8eaeb5f67b))

## [1.1.4](https://github.com/GoogleCloudPlatform/deploystack/compare/v1.1.3...v1.1.4) (2023-02-21)


### Bug Fixes

* trying to update all the things to ensure proper dependencies ([e76fdae](https://github.com/GoogleCloudPlatform/deploystack/commit/e76fdae2123f6997268c7bc8a552e8f552e24a49))

## [1.1.3](https://github.com/GoogleCloudPlatform/deploystack/compare/v1.1.2...v1.1.3) (2023-02-21)


### Bug Fixes

* still trying to get please-release to work ([d0967bd](https://github.com/GoogleCloudPlatform/deploystack/commit/d0967bd529d66e68bd2d85e37b536967a230361e))

## [1.1.2](https://github.com/GoogleCloudPlatform/deploystack/compare/v1.1.1...v1.1.2) (2023-02-21)


### Bug Fixes

* trying to get release-please to properly tag subpackages ([604e2d1](https://github.com/GoogleCloudPlatform/deploystack/commit/604e2d1e019071648d69fcbcf7359220119b9ea7))

## [1.1.1](https://github.com/GoogleCloudPlatform/deploystack/compare/v1.1.0...v1.1.1) (2023-02-21)


### Bug Fixes

* updating release-please config ([ce1505f](https://github.com/GoogleCloudPlatform/deploystack/commit/ce1505f6979f8c04a754f8084bc1b48d658f656f))

## [1.1.0](https://github.com/GoogleCloudPlatform/deploystack/compare/v1.0.1...v1.1.0) (2023-02-21)


### Features

* added a warning about slow lookups in a new project ([855ac53](https://github.com/GoogleCloudPlatform/deploystack/commit/855ac537cf7470ba4df3788e229790efdb3e7cbd))
* added caching to speed up these requests ([a0aaad8](https://github.com/GoogleCloudPlatform/deploystack/commit/a0aaad877b643181757d84e37621109f9e2f9d5a))
* adding project retrieval to the mix ([12d6dcd](https://github.com/GoogleCloudPlatform/deploystack/commit/12d6dcd92e8f9b3b6a9bc8a9b987a40cb5998966))
* build a "don't put into settings flag" into the base configuration of pages ([d616a13](https://github.com/GoogleCloudPlatform/deploystack/commit/d616a138f29a10f56f7f14efb0636b34358f3621))
* forgot to add prependProject setting. Working now ([a6b2e25](https://github.com/GoogleCloudPlatform/deploystack/commit/a6b2e25c2980b459d715e82ba4ec5c1cbbdac735))
* implemented a back function ([c65badc](https://github.com/GoogleCloudPlatform/deploystack/commit/c65badc1f5a2041e725ba43015ff8af8438948f8))
* renamed package to better reflect purpose ([70850fc](https://github.com/GoogleCloudPlatform/deploystack/commit/70850fce76e7402261d9ae31177ffeaddec1722b))
* renamed package to make more clear its purpose ([d3d476e](https://github.com/GoogleCloudPlatform/deploystack/commit/d3d476ef2ae3da26910b3f44d6ce3edd5f07feab))


### Bug Fixes

* actually fixing the problem with project numbers ([5175fea](https://github.com/GoogleCloudPlatform/deploystack/commit/5175fea1dedd89b5d5da9c4a8174e4c9c7f92ecb))
* commented out tests that were commented out before ([3755a5f](https://github.com/GoogleCloudPlatform/deploystack/commit/3755a5f5ca3c4f182e545112b4cc1ae619ded1b1))
* commenting out previously commented test ([a023d8a](https://github.com/GoogleCloudPlatform/deploystack/commit/a023d8a2115a7f6d83d52f989db175ba3433f6d2))
* commenting out tests that were commented out before ([d17cb02](https://github.com/GoogleCloudPlatform/deploystack/commit/d17cb02b8959e564b8d32fa9d78cff2dd7f535e2))
* converted spinner to use dsStyles to render color on Cloud Shell ([a7456b7](https://github.com/GoogleCloudPlatform/deploystack/commit/a7456b743a79585a481bd66242737ec31f7748d8))
* corrected issue where domain contact info was not being refreshed properly ([546af35](https://github.com/GoogleCloudPlatform/deploystack/commit/546af355c1a795017accddd23271e3a429db3480))
* correcting billing check quality ([a830b43](https://github.com/GoogleCloudPlatform/deploystack/commit/a830b43c443dd6fea60b5ddd41ec3e096c7b9cd4))
* fixes dealing with changes in color rendering ([26bd638](https://github.com/GoogleCloudPlatform/deploystack/commit/26bd63893c6cfc23929d4517c1457a4fa06698a3))
* fixing error where ctr-c on the exit page causes a panic ([3974a8b](https://github.com/GoogleCloudPlatform/deploystack/commit/3974a8bbf92387104e5f644d4054cf3414b5560a))
* getting rid of old calls. ([a391e39](https://github.com/GoogleCloudPlatform/deploystack/commit/a391e39aa119211e49570f4f1e6a67b804a87325))
* got domain availability test to work ([7103da2](https://github.com/GoogleCloudPlatform/deploystack/commit/7103da2c83a1af3cb84ef4981a66ff93ef4506e6))
* made currentProject a queue value instead of a global one ([c3bfb7b](https://github.com/GoogleCloudPlatform/deploystack/commit/c3bfb7bcbd46d781e84612dc4327204f93ba8897))
* making sure that gcloud config project set is called ([4295f73](https://github.com/GoogleCloudPlatform/deploystack/commit/4295f739633f75e9360b76f97d662fc363cf73d1))
* making sure that project creation captures project number ([fb0d8d7](https://github.com/GoogleCloudPlatform/deploystack/commit/fb0d8d7cbb46a54c191e642846e39369ed128f77))
* massively changed the way color was rendered to make it simpler ([a5c652a](https://github.com/GoogleCloudPlatform/deploystack/commit/a5c652a9cdfc9a6820681ac734a7c1a51af8d300))
* moved project number retrieval to the actual flow. ([a546144](https://github.com/GoogleCloudPlatform/deploystack/commit/a546144dc53b5c0a02a9918d14daa96dd22288fc))
* needed to account for nils in casting operation ([4b0ea90](https://github.com/GoogleCloudPlatform/deploystack/commit/4b0ea90634dedde02b8ed34920639fdf341a944f))
* process is now exit(1) when user stops process ([ab4bbcf](https://github.com/GoogleCloudPlatform/deploystack/commit/ab4bbcf81370953a083b3224f01035401c9403e3))
* resolves issue with color not showing in CloudShell ([df2a788](https://github.com/GoogleCloudPlatform/deploystack/commit/df2a788e2a4b7e9accb481e5b5c8c69fa8906b39))
* resolving issue where new accounts weren't getting billing attached ([fba7fad](https://github.com/GoogleCloudPlatform/deploystack/commit/fba7fadadeef1a7af510d4dfd58594a11d1e111e))
* restoreing full screen mode. ([3259419](https://github.com/GoogleCloudPlatform/deploystack/commit/3259419229d9abee0991a907350969def9ef5ede))
* testing color rendering for Cloud Shell ([9535e12](https://github.com/GoogleCloudPlatform/deploystack/commit/9535e12b8a0aa1c828b152fa751b9265697b84e1))
* testing if the alt screen is the problem on cloudshell ([8ce98f9](https://github.com/GoogleCloudPlatform/deploystack/commit/8ce98f999cfe31b61165180e3e389df2185c2ee4))
* tests didn't get renamed ([bfa2066](https://github.com/GoogleCloudPlatform/deploystack/commit/bfa20663d6ab659e5fe1217c0dcc3dc0d3950d79))
* updated colors to try and make it look better on Cloud Shell ([5d7f220](https://github.com/GoogleCloudPlatform/deploystack/commit/5d7f2206b3445d9e94430dec61df0b2e7db90d8a))
* when you are just asking for the billing account, you don't need to attach it. ([0f9f547](https://github.com/GoogleCloudPlatform/deploystack/commit/0f9f5478f79d883ec7f857937e31a5bbcb83d3ea))

## [1.0.1](https://github.com/GoogleCloudPlatform/deploystack/compare/v1.0.0...v1.0.1) (2023-02-16)


### Bug Fixes

* all requests for were failing because of bad projects ([35fdeb8](https://github.com/GoogleCloudPlatform/deploystack/commit/35fdeb865af96c9cf84bfecf999a4661e82a9dc4))
* reverting fix ([704d530](https://github.com/GoogleCloudPlatform/deploystack/commit/704d530ab4f19493bfa7d5b2a0699e86de198fce))
* reverting reverting fix. ([8b95f62](https://github.com/GoogleCloudPlatform/deploystack/commit/8b95f6200f4afb92d0c2b10fca4ef01b34cc98d6))

## 1.0.0 (2023-02-15)


### Features

* add Dockerfile to create a Docker image ([d5a87db](https://github.com/GoogleCloudPlatform/deploystack/commit/d5a87db838b67a33065361ad06a227aa8dfaff3f))
* add dvatracker and supporting libs ([0c1f91a](https://github.com/GoogleCloudPlatform/deploystack/commit/0c1f91a4144a11d068ab0091868ae19b33c8e2f9))
* added a deploystack documentation creator. ([16fd558](https://github.com/GoogleCloudPlatform/deploystack/commit/16fd558e9a9c0217dc38d3e7724fec8bad60179f))
* added ability to pass create labels for custom setting options. ([d0e60e3](https://github.com/GoogleCloudPlatform/deploystack/commit/d0e60e3ee38dc4c462d64f4721b0f3581158b84a))
* added tool that converts json to yaml ([d5e523b](https://github.com/GoogleCloudPlatform/deploystack/commit/d5e523bfcb280f51ad37da278d506ca4c25afe40))
* adding name to config and now searching out the value using git ([af07341](https://github.com/GoogleCloudPlatform/deploystack/commit/af073411ec2945e1151cb0913baa845c6c90c5c1))
* adding test creator to mix. ([c73e53c](https://github.com/GoogleCloudPlatform/deploystack/commit/c73e53cb268cca0dbb57f531f6d2259efb20b0eb))
* adding the ability to display documentation links. ([b1cbff3](https://github.com/GoogleCloudPlatform/deploystack/commit/b1cbff3f87eb167bbb00b9e64925b101f2b2460e))
* adding yaml parsing for config ([066b0e5](https://github.com/GoogleCloudPlatform/deploystack/commit/066b0e5815a7ba01b86e58ccdd549f77d5e11bd4))
* initial import of code that will speed up documentation. ([68fa1b2](https://github.com/GoogleCloudPlatform/deploystack/commit/68fa1b29e3a7d4230beac8a58ddab4019d5887aa))


### Bug Fixes

* a number of the names of functions were pretty confusing for reuse. ([02c99cf](https://github.com/GoogleCloudPlatform/deploystack/commit/02c99cf898178ded80cbd394c0eb77bd9d75cd4c))
* added a test for iam ([b876d8a](https://github.com/GoogleCloudPlatform/deploystack/commit/b876d8a0fb64828efbd5db8667c7e961a44131ef))
* added test for cloudbuild stuff ([2395694](https://github.com/GoogleCloudPlatform/deploystack/commit/2395694575c31f7c105db6b7d8c0766c93bc6693))
* Adding a user input to avoid timeouts on  the cloud shell prompt. ([4524fab](https://github.com/GoogleCloudPlatform/deploystack/commit/4524fab1f4c7b471fc4d8fb5bbfa51a961e93990))
* adding ability to check project existence ([5e4e6f8](https://github.com/GoogleCloudPlatform/deploystack/commit/5e4e6f89b183644940f857c1088f4291989d85e6))
* adding additional api features to deploystack. ([4eadc5c](https://github.com/GoogleCloudPlatform/deploystack/commit/4eadc5c6bb28a48b4a22580946e9a3f032300aca))
* adding aliasof field to submodule ([cfb22ca](https://github.com/GoogleCloudPlatform/deploystack/commit/cfb22ca959e729c7b5e71d0d80d8db551ac0f6d9))
* adding debugging help for project creation ([8b2c4f7](https://github.com/GoogleCloudPlatform/deploystack/commit/8b2c4f753a1b1394e07010d51405a2e7dba18d1c))
* adding more capabilities around gcloud management to core. ([5d6bd7c](https://github.com/GoogleCloudPlatform/deploystack/commit/5d6bd7c32b967ea583935b8cbb62067f9f2a5665))
* adding secret manager capabilities to package ([6741cd8](https://github.com/GoogleCloudPlatform/deploystack/commit/6741cd84efea641c2bcfc49185f02ace3faa803d))
* adding test for secretmanager code. ([cbe4d04](https://github.com/GoogleCloudPlatform/deploystack/commit/cbe4d04abecc60147dec7287d23fa5f0bf5a0458))
* adding the ability to work with triggers and scheduled jobs ([8332cbf](https://github.com/GoogleCloudPlatform/deploystack/commit/8332cbf7fb01c626ebe8b6c8a945384f729fba3e))
* allowing settings to be removed from stacks ([e1aa5bb](https://github.com/GoogleCloudPlatform/deploystack/commit/e1aa5bbe33452a13b10abf52b5416654b2c2ffcc))
* API changes caused test failure ([bcbda59](https://github.com/GoogleCloudPlatform/deploystack/commit/bcbda591755ac279b388ac81a4d1eccf84549fee))
* automating subpackage testing. ([c7eb566](https://github.com/GoogleCloudPlatform/deploystack/commit/c7eb566244899e3cbbc60e4884c43c77182dd162))
* bad behavior in sub lib ([a64f41a](https://github.com/GoogleCloudPlatform/deploystack/commit/a64f41ae7e896964aea66153c67e11bdd6e1589a))
* broken test fix. ([4df3414](https://github.com/GoogleCloudPlatform/deploystack/commit/4df34148b84e9958b19df01951ee66ff58da675d))
* broken test. ([fa29b23](https://github.com/GoogleCloudPlatform/deploystack/commit/fa29b23574589e40c8968d431427498a0ff9507e))
* broken tests ([f9de00d](https://github.com/GoogleCloudPlatform/deploystack/commit/f9de00d74ca2d86c23ce634288d99e91baf931b3))
* changing to cloudshell to make sure this passes the test. ([5050a39](https://github.com/GoogleCloudPlatform/deploystack/commit/5050a3983b65a7c1e5f5d5b4cbbfb18000e3d29c))
* consolidating deploystack - github code in one place ([093db5f](https://github.com/GoogleCloudPlatform/deploystack/commit/093db5f39f014ee6d3262048cbefa77b8626a3a5))
* consolidating tf processing in gcloudtf ([e07dad5](https://github.com/GoogleCloudPlatform/deploystack/commit/e07dad59174b3ffdf1e2cff6d8e04b6b58f47d46))
* correcting broken test ([9775dff](https://github.com/GoogleCloudPlatform/deploystack/commit/9775dffb738836808574b2d640f9e536bec97462))
* correcting dependency in test ([1431a27](https://github.com/GoogleCloudPlatform/deploystack/commit/1431a27bd75b21c7647cb6132f8355bce1c9872a))
* correcting failed test. ([bd2130b](https://github.com/GoogleCloudPlatform/deploystack/commit/bd2130b39466834a523a112a1624b7f6ff810a2f))
* correcting intermittent failures in test ([f1d8df0](https://github.com/GoogleCloudPlatform/deploystack/commit/f1d8df023a176c6e604c1baf0ef3687faabe1c1b))
* creating single source of truth for tf and repo data. ([91bd664](https://github.com/GoogleCloudPlatform/deploystack/commit/91bd664aa07b6b055eea63172278b4817890c605))
* dependecies out of date contributing errors ([472aca3](https://github.com/GoogleCloudPlatform/deploystack/commit/472aca384d29e48937f73a0faaadc371ad72f210))
* enabled multiple projects in a stack ([19b5603](https://github.com/GoogleCloudPlatform/deploystack/commit/19b5603edce2d822b1686db23b8f0c18b6ca491a))
* ensuring new components are tested ([92562b9](https://github.com/GoogleCloudPlatform/deploystack/commit/92562b92e498bcb5e3643471ef0034df21d29fcb))
* extending project creation ([45a3cb7](https://github.com/GoogleCloudPlatform/deploystack/commit/45a3cb7c0dfa426fe3af59f4af83708aed7307c1))
* finalizing test design ([347aa25](https://github.com/GoogleCloudPlatform/deploystack/commit/347aa25bfd0da89b0124149f6ab73a24f43d2045))
* finally got it working. ([1a074b5](https://github.com/GoogleCloudPlatform/deploystack/commit/1a074b52a29042dcbf9a19fed14316a1b48c087b))
* fix isn't taking ([0ac19c2](https://github.com/GoogleCloudPlatform/deploystack/commit/0ac19c288e8c93c62f89224db9ce00bdaa388814))
* fixed issue where cloning with ssh caused errors ([5848091](https://github.com/GoogleCloudPlatform/deploystack/commit/5848091cbb0d1244f36feae761da5f6e1ed40ce7))
* fixing broken cloud domain code ([8e72968](https://github.com/GoogleCloudPlatform/deploystack/commit/8e7296884adde46275b4da50bcaab882b3b6e669))
* fixing bug that caused project creation to fail ([303e61a](https://github.com/GoogleCloudPlatform/deploystack/commit/303e61a4f987ab07b1e0eaba667d3b520204e29b))
* fixing the demo app ([91f6211](https://github.com/GoogleCloudPlatform/deploystack/commit/91f6211ebde5b7c596cc447d937fc26ca5f1f5b2))
* FLAAAAAAAAKE ([3f6d148](https://github.com/GoogleCloudPlatform/deploystack/commit/3f6d148fc3b88d33e1d4c2129d1a2b0d366451d5))
* flaky test is gonna flake ([2a96f9a](https://github.com/GoogleCloudPlatform/deploystack/commit/2a96f9a4da7aa52ce083aba4f6666852276b1fff))
* flaky test now unflaked ([ed7215b](https://github.com/GoogleCloudPlatform/deploystack/commit/ed7215b220bc52206baa0471ef2c6bda5afbc4a2))
* forcing DVAs to be unique. ([dc9dfde](https://github.com/GoogleCloudPlatform/deploystack/commit/dc9dfde56f3cb6e5818621f6f3786b7ac4639387))
* getting all the version stuff working properly ([140556a](https://github.com/GoogleCloudPlatform/deploystack/commit/140556afab66aad5233032c402072641f561d223))
* getting tests to pass ([c425b87](https://github.com/GoogleCloudPlatform/deploystack/commit/c425b875c35ea4a41546facefeb501daaf05fb70))
* getting this folder selecting perfect. ([1a5f9b2](https://github.com/GoogleCloudPlatform/deploystack/commit/1a5f9b20fc98d1981bfa2d430b2327c04d8be0f1))
* getting this logic correct ([0c738d6](https://github.com/GoogleCloudPlatform/deploystack/commit/0c738d697d441af631b40041bbcaccf446b9887a))
* got documenation for devsite working. ([7422214](https://github.com/GoogleCloudPlatform/deploystack/commit/7422214e03c6b7f88123176b1f1b59fb1239003b))
* improved the creation of tests ([cfc4edc](https://github.com/GoogleCloudPlatform/deploystack/commit/cfc4edc0f3f968f8e6c84c42a493843175325515))
* improving project creation code ([af578db](https://github.com/GoogleCloudPlatform/deploystack/commit/af578db9c0ec0049e8f23dfe064c67b61a512864))
* initial implementation of test framework ([55be6be](https://github.com/GoogleCloudPlatform/deploystack/commit/55be6be840dd3ee220433cd896df79ae81a01db6))
* initial import of tui code ([e53bb85](https://github.com/GoogleCloudPlatform/deploystack/commit/e53bb85709de1bc8c025ce0948b6cec348cdaab6))
* initial move to move gcloud client stuff to own package ([22dd503](https://github.com/GoogleCloudPlatform/deploystack/commit/22dd50354c32014ec6b845cd30ebff8c33a64aa2))
* introducing tests for newly added features. ([272b41e](https://github.com/GoogleCloudPlatform/deploystack/commit/272b41e2cd3e80a115da693d032b30c4d3440295))
* issues with default control for machineType ([b7df60a](https://github.com/GoogleCloudPlatform/deploystack/commit/b7df60ab9230dd4329ac8c62777f77fc4f66b1a2))
* it sure is hard to call an ssh clone from within Cloud Build ([1eac57f](https://github.com/GoogleCloudPlatform/deploystack/commit/1eac57fcbe4dc77475230170ea03ef9268e80f24))
* left a bug in the flaky test fix. ([284919e](https://github.com/GoogleCloudPlatform/deploystack/commit/284919e43b65d9d3de858cb8f9a2f3181ebc8d67))
* make defaultUserAgent include the stack name ([57b52b3](https://github.com/GoogleCloudPlatform/deploystack/commit/57b52b3ffd9e676b960dd63cc55e1a8d7dc62617))
* make sure that this new feature doesn't block users without permission ([090ff20](https://github.com/GoogleCloudPlatform/deploystack/commit/090ff20ef281a422e40587af0d950f2cdd4d16a4))
* make test work locally ([d50f83b](https://github.com/GoogleCloudPlatform/deploystack/commit/d50f83bfaf420089314bd17ba61ab454c590331e))
* making all of the tests pass for the lass refactoring ([8a9db91](https://github.com/GoogleCloudPlatform/deploystack/commit/8a9db911185989c03c88d5a1971f018d952668d7))
* making Create Project to work again ([8715d97](https://github.com/GoogleCloudPlatform/deploystack/commit/8715d9737b803eb836e057ac4afefcd356c32315))
* making doccreator work with gcloudtf ([27bfb5e](https://github.com/GoogleCloudPlatform/deploystack/commit/27bfb5ed24acc7f9dc3138428f6d1a9f14610a73))
* making storage tests work ([c82a5af](https://github.com/GoogleCloudPlatform/deploystack/commit/c82a5afb2bcc92ff80d0da825d8974bf50602afa))
* making sure contact doesn't get added back. ([97de87f](https://github.com/GoogleCloudPlatform/deploystack/commit/97de87f2aa69020ae8c173d64aff372bf6246b71))
* making sure doc creation works with new stucture ([9935466](https://github.com/GoogleCloudPlatform/deploystack/commit/993546690117d6c9b021ae9b158b10b04a8a7723))
* making sure tests work properly ([2d799e3](https://github.com/GoogleCloudPlatform/deploystack/commit/2d799e390ae6dffb66605c913907cb95f0b4a683))
* making the idea of arguments work. ([a101a2c](https://github.com/GoogleCloudPlatform/deploystack/commit/a101a2ccadd804a7a6b88ee488586ecfba75a35d))
* making things more organized. ([2a69730](https://github.com/GoogleCloudPlatform/deploystack/commit/2a697304310a5ee0780b621d22854f74459bd109))
* merge resolution added a breaking change. ([34ca69f](https://github.com/GoogleCloudPlatform/deploystack/commit/34ca69f4e3558db2376fbb42fea90089dd516e9a))
* moving code to dsgithub package ([b9c9d04](https://github.com/GoogleCloudPlatform/deploystack/commit/b9c9d0405f98a0c0449c31566914ef717f3e6e1b))
* moving constants to the main deploystack project ([cc80588](https://github.com/GoogleCloudPlatform/deploystack/commit/cc805880581987370ec7c37bc550e88f6fcf66bc))
* need to find a permanent fix for this ([d522665](https://github.com/GoogleCloudPlatform/deploystack/commit/d522665f63f81608e79ed8e4f7a98d23f6da7f8e))
* needed to add gitkeep files to keep directories existing. ([fa7e760](https://github.com/GoogleCloudPlatform/deploystack/commit/fa7e7602e74df8a061516c1459412f4d5b445dc7))
* omitted line. ([b6042c8](https://github.com/GoogleCloudPlatform/deploystack/commit/b6042c834cb2f437dad16f470b8c3e254ce4cd76))
* Ported over all of the gcloud functionality ([c380557](https://github.com/GoogleCloudPlatform/deploystack/commit/c38055786c572a6d90b6c8b24952f5a909f507c2))
* problem with project name vs id in selector. ([cff9e70](https://github.com/GoogleCloudPlatform/deploystack/commit/cff9e70185d66f0461285cf17c824abb2a447886)), closes [#14](https://github.com/GoogleCloudPlatform/deploystack/issues/14) [#15](https://github.com/GoogleCloudPlatform/deploystack/issues/15)
* pulled over the rest of the functions to complete the interface for ui ([40d6741](https://github.com/GoogleCloudPlatform/deploystack/commit/40d674124d0681ded6c61d39b1906f758e43ed5a))
* refactoring ([37d01b2](https://github.com/GoogleCloudPlatform/deploystack/commit/37d01b2d871b65470daa0243ee60fb46f7ba09a8))
* refactoring labeledvalues control for comprehension and better testing ([38f2a1a](https://github.com/GoogleCloudPlatform/deploystack/commit/38f2a1ae46cab460050ab97080b7a08a1affeed7))
* remove local overrides ([6cbfb4b](https://github.com/GoogleCloudPlatform/deploystack/commit/6cbfb4b2daa947fd59c9717ff3402f203dd61859))
* remove project_name from tvars file. ([cb2cff9](https://github.com/GoogleCloudPlatform/deploystack/commit/cb2cff97019ed0049db46484e7bc2ee47e72c875))
* removing a bunch of unused code ([786a7f7](https://github.com/GoogleCloudPlatform/deploystack/commit/786a7f7b4eec242ee0f495865ecb99fecfa4f6dc))
* removing debugging content ([083dff0](https://github.com/GoogleCloudPlatform/deploystack/commit/083dff05c7ee051e9e4997c301a7c87c4d76d9a4))
* removing dev configs ([c30e5e6](https://github.com/GoogleCloudPlatform/deploystack/commit/c30e5e6a5b4ab2816a6f734bc673b229c5669b8e))
* removing generated file ([52e73fe](https://github.com/GoogleCloudPlatform/deploystack/commit/52e73fecb76a45e092ccb2fe47466e2c0366a549))
* removing stack_name from the terraform variables file. ([819240b](https://github.com/GoogleCloudPlatform/deploystack/commit/819240b9da1e440cde8bf84236bcbcfc62b8f094))
* removing the need for defaultValue ([9c13e8c](https://github.com/GoogleCloudPlatform/deploystack/commit/9c13e8c55bd3f8b92b18740ac0b6acf3d512b3cb))
* renaming and exposing some abilities. ([e67e52a](https://github.com/GoogleCloudPlatform/deploystack/commit/e67e52a1f9ffc27ed7de44b38339abd755e8954e))
* renaming services for clarity ([ad0351b](https://github.com/GoogleCloudPlatform/deploystack/commit/ad0351bee366b24e4dd7885b7f195ad276e8885b))
* renaming things to cut down on redundancy, long names and redundancy ([16455ef](https://github.com/GoogleCloudPlatform/deploystack/commit/16455ef6df7cc3b0d86fd914b3c2877bb431c542))
* rolling back bad change. ([d1380c8](https://github.com/GoogleCloudPlatform/deploystack/commit/d1380c8e54280a9dd83d422876b40a4b71951ebb))
* setting up testing for package to refactor a bit ([9ad840c](https://github.com/GoogleCloudPlatform/deploystack/commit/9ad840c8ee64043f8536f05d243e490c1e008b33))
* ssh volume is no bueno. ([e730b7e](https://github.com/GoogleCloudPlatform/deploystack/commit/e730b7eaebcf67baf586adb9479249790418d23e))
* start using go workspaces to help ([f0e1f3f](https://github.com/GoogleCloudPlatform/deploystack/commit/f0e1f3f7352286e90bceaf3689c00bc7a1683b55))
* still trying to get build to pull down git via SSH ([d4c3117](https://github.com/GoogleCloudPlatform/deploystack/commit/d4c3117fe0e1e6ba41e88816192e9df6379cca0c))
* stopped ui from showing non billing enabled projects incorrectly ([000f5ab](https://github.com/GoogleCloudPlatform/deploystack/commit/000f5abe716f79fe33ad6359ace7ed25a092d29c))
* stupid typos ([0c39a5c](https://github.com/GoogleCloudPlatform/deploystack/commit/0c39a5c37e8f31da5cda1cb1d4e236162dcb75ea))
* temporarily commenting out flaky test. ([5e8e2f3](https://github.com/GoogleCloudPlatform/deploystack/commit/5e8e2f3fa79e68e5b174f1003ce21272c8591abf))
* tests broken because environment isn't right. ([e88931f](https://github.com/GoogleCloudPlatform/deploystack/commit/e88931f5dae41993069682fe3d19005f3b373e49))
* there are two expensive tests in this build, so upping the time and adding higher level machine ([bf62b97](https://github.com/GoogleCloudPlatform/deploystack/commit/bf62b97d1f00aa549a9ab6d76f0a9a92e9e23380))
* this is harder than it should be ([569fb52](https://github.com/GoogleCloudPlatform/deploystack/commit/569fb52ac6ca82b814c0513a4fd55b34b633f28c))
* this might fix ssh but break other tests. ([05c1c46](https://github.com/GoogleCloudPlatform/deploystack/commit/05c1c460977f9d61078a36b1c190524ebcb510c0))
* tiniest of changes to setup ([705eb0b](https://github.com/GoogleCloudPlatform/deploystack/commit/705eb0baa50dcc53171153d0722e446f09ce5b3b))
* trying to get all packages working properly ([5768e28](https://github.com/GoogleCloudPlatform/deploystack/commit/5768e2828529e17974feb8f1fe43feb706329f70))
* trying to get failing test working. ([6b6a544](https://github.com/GoogleCloudPlatform/deploystack/commit/6b6a5441cf58b6c9a4d5b85cbbb8f6556f3ebdaa))
* trying to get stuff working. ([6de1b79](https://github.com/GoogleCloudPlatform/deploystack/commit/6de1b7971af563a7ef7b7b39c86c43219630a2f7))
* trying to update package issues ([3fb7a7b](https://github.com/GoogleCloudPlatform/deploystack/commit/3fb7a7b45de2e4bdfa83b07c3259750eb40a1af9))
* trying to update sub package stuff ([3ececa5](https://github.com/GoogleCloudPlatform/deploystack/commit/3ececa508b05427a65e56eddf6ccd0c79213b695))
* tweaking output to be lowercase ([7634376](https://github.com/GoogleCloudPlatform/deploystack/commit/76343761c067a753023ad229f61ac0fca035b6ca))
* typo in type - oh! ([7544d7d](https://github.com/GoogleCloudPlatform/deploystack/commit/7544d7dbc4997a26f7fa6b94e862643d8c9806b0))
* update tests so that it will detect changes to project selector ([6e1caf4](https://github.com/GoogleCloudPlatform/deploystack/commit/6e1caf418a459607c3a775557ace8ab6ce566362))
* updated a lot of the capabilities of the test rig and added documentation ([b33e329](https://github.com/GoogleCloudPlatform/deploystack/commit/b33e3295e953afbceca164610feb4e54f55ea4a3))
* updated dependencies ([341565f](https://github.com/GoogleCloudPlatform/deploystack/commit/341565f89aaabf849ba2b31df3c1fbe1c9bb2492))
* updating abilities of gcloudtf package ([20a336a](https://github.com/GoogleCloudPlatform/deploystack/commit/20a336a36cccda9f7f0f5c2cb34ba4912b9e970e))
* updating dependencies. ([6dc6ef9](https://github.com/GoogleCloudPlatform/deploystack/commit/6dc6ef9a10f7a543a3f26057dd66c2e3d654d673))
* updating some leaky code ([a38e440](https://github.com/GoogleCloudPlatform/deploystack/commit/a38e4402e684568ceb937624de4998ae9cf8d1c5))
* updating test creator to use gcloudtf issues. ([71c9a6b](https://github.com/GoogleCloudPlatform/deploystack/commit/71c9a6b009640e9924e0631e34925a3fbb60ec85))
* updating test creator to use gcloudtf issues. ([0002e13](https://github.com/GoogleCloudPlatform/deploystack/commit/0002e13907097b0df01c6ab40a2434741bd2ad76))
* updating tests ([a2c1617](https://github.com/GoogleCloudPlatform/deploystack/commit/a2c1617eab219bd2fe1b90ee5a69d1f79874d9bc))
* updating the dependencies ([2d9ee56](https://github.com/GoogleCloudPlatform/deploystack/commit/2d9ee56ebcd2e50cb4395b066282a0b005b658bd))
* updating workspace settings ([3512aa3](https://github.com/GoogleCloudPlatform/deploystack/commit/3512aa34dd85c9fcc68cb030c8f797ad6e4edb3c))
* work! ([0b87451](https://github.com/GoogleCloudPlatform/deploystack/commit/0b8745196900632829822149c91f456e7d5aef58))
* yet another typo ([ce930fa](https://github.com/GoogleCloudPlatform/deploystack/commit/ce930fa5355ef8a03960440af29f738cbfc77270))
