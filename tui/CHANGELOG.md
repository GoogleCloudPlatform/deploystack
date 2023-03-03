# Changelog

## [1.6.0](https://github.com/GoogleCloudPlatform/deploystack/compare/tui/v1.5.1...tui/v1.6.0) (2023-03-03)


### Features

* added a warning about slow lookups in a new project ([855ac53](https://github.com/GoogleCloudPlatform/deploystack/commit/855ac537cf7470ba4df3788e229790efdb3e7cbd))
* added ability for more than one stack to exist in the same repo ([959b797](https://github.com/GoogleCloudPlatform/deploystack/commit/959b797478114da0d8d08e336605645f3bd02e56))
* added authorsettings to replace hardsettings ([553979b](https://github.com/GoogleCloudPlatform/deploystack/commit/553979bf7520655c8467ce4ea4c6817f06447a4f))
* added caching to speed up these requests ([a0aaad8](https://github.com/GoogleCloudPlatform/deploystack/commit/a0aaad877b643181757d84e37621109f9e2f9d5a))
* added the progress bar. ([4e6a76c](https://github.com/GoogleCloudPlatform/deploystack/commit/4e6a76cead93920c619b6798b7caa3b743a58605))
* build a "don't put into settings flag" into the base configuration of pages ([d616a13](https://github.com/GoogleCloudPlatform/deploystack/commit/d616a138f29a10f56f7f14efb0636b34358f3621))
* created new package to remove threat of circular dependency ([8b0738e](https://github.com/GoogleCloudPlatform/deploystack/commit/8b0738e28a839d2f9a21cb4c880ddd382d9017e2))
* forgot to add prependProject setting. Working now ([a6b2e25](https://github.com/GoogleCloudPlatform/deploystack/commit/a6b2e25c2980b459d715e82ba4ec5c1cbbdac735))
* implemented a back function ([c65badc](https://github.com/GoogleCloudPlatform/deploystack/commit/c65badc1f5a2041e725ba43015ff8af8438948f8))
* updated tui to work with new configuration settings ([6a5a17a](https://github.com/GoogleCloudPlatform/deploystack/commit/6a5a17ac6fc097dcd563a5a5a0638972769b9e5d))


### Bug Fixes

* actually fixing the problem with project numbers ([5175fea](https://github.com/GoogleCloudPlatform/deploystack/commit/5175fea1dedd89b5d5da9c4a8174e4c9c7f92ecb))
* added an example for multi stack repos ([52ca81f](https://github.com/GoogleCloudPlatform/deploystack/commit/52ca81f34895b6a9a5b0ecc6e3941d9155ffaa1a))
* added billing account capture and made sure that was working ([d7c0dc3](https://github.com/GoogleCloudPlatform/deploystack/commit/d7c0dc3c51699da801e7eaccbcca42b453170723))
* adding a clean preprocessor for the last screen ([c3941d7](https://github.com/GoogleCloudPlatform/deploystack/commit/c3941d792be3663847ce42e06396f4da7806a8da))
* better default picker style ([4f962e5](https://github.com/GoogleCloudPlatform/deploystack/commit/4f962e53f04a886c4414bd8070c103a116987a7c))
* checking in todays changes ([18e3775](https://github.com/GoogleCloudPlatform/deploystack/commit/18e37755d37312bd26735c19dc163bd498dd2b5c))
* cloud shell always thinks it's in dark mode ([2fc66fd](https://github.com/GoogleCloudPlatform/deploystack/commit/2fc66fd214214226c5027b6e3ffa130aebe20152))
* converted spinner to use dsStyles to render color on Cloud Shell ([a7456b7](https://github.com/GoogleCloudPlatform/deploystack/commit/a7456b743a79585a481bd66242737ec31f7748d8))
* corrected issue where domain contact info was not being refreshed properly ([546af35](https://github.com/GoogleCloudPlatform/deploystack/commit/546af355c1a795017accddd23271e3a429db3480))
* correcting bug where Project forms looked bad. ([1a2e336](https://github.com/GoogleCloudPlatform/deploystack/commit/1a2e336cd7889d6b17d9d7050b7e5715496a4566))
* correcting dependency in test ([1431a27](https://github.com/GoogleCloudPlatform/deploystack/commit/1431a27bd75b21c7647cb6132f8355bce1c9872a))
* debugging might be screwing things up here. ([177cc54](https://github.com/GoogleCloudPlatform/deploystack/commit/177cc54fa769e927e76814feb28f1fca3d70b6c3))
* didn't like short form ([d4219c3](https://github.com/GoogleCloudPlatform/deploystack/commit/d4219c38c537dd4623f7c78a51fdd28855fe2dd8))
* fixed an issue where when the default value isn't set, the list isn't seleced ([6d69e80](https://github.com/GoogleCloudPlatform/deploystack/commit/6d69e80ccaef0d09c95859ddfe319e58c408cc2a))
* fixes dealing with changes in color rendering ([26bd638](https://github.com/GoogleCloudPlatform/deploystack/commit/26bd63893c6cfc23929d4517c1457a4fa06698a3))
* fixing breaking test ([9784da5](https://github.com/GoogleCloudPlatform/deploystack/commit/9784da55dbb6ec641b04fcdd0bbb74c942029a48))
* fixing error where ctr-c on the exit page causes a panic ([3974a8b](https://github.com/GoogleCloudPlatform/deploystack/commit/3974a8bbf92387104e5f644d4054cf3414b5560a))
* getting all the version stuff working properly ([140556a](https://github.com/GoogleCloudPlatform/deploystack/commit/140556afab66aad5233032c402072641f561d223))
* getting billing to work properly ([9894e2b](https://github.com/GoogleCloudPlatform/deploystack/commit/9894e2bcea439b197755ab9a324dbbbe7ed71e09))
* getting rid of old calls. ([a391e39](https://github.com/GoogleCloudPlatform/deploystack/commit/a391e39aa119211e49570f4f1e6a67b804a87325))
* got picker default values working properly ([ef623d7](https://github.com/GoogleCloudPlatform/deploystack/commit/ef623d748107ee0cb13e30082926bbeeb5015edb))
* initial import of tui code ([e53bb85](https://github.com/GoogleCloudPlatform/deploystack/commit/e53bb85709de1bc8c025ce0948b6cec348cdaab6))
* just some tweaks to improving linting ([5c716d7](https://github.com/GoogleCloudPlatform/deploystack/commit/5c716d7d05f1e7235f9319edacd5578669c081ef))
* list items limited to 50 chars and defaults now display correctly ([46b2e3b](https://github.com/GoogleCloudPlatform/deploystack/commit/46b2e3b8288c432e6ebff4dfefd756e9040dd38f))
* made currentProject a queue value instead of a global one ([c3bfb7b](https://github.com/GoogleCloudPlatform/deploystack/commit/c3bfb7bcbd46d781e84612dc4327204f93ba8897))
* made tui compatible with changes to config ([0ba76e4](https://github.com/GoogleCloudPlatform/deploystack/commit/0ba76e4dcf5270b9ac1cc71ceb96fcc395d5fa3b))
* make default options display normally if there are less than 10 of them. ([0cc9aa6](https://github.com/GoogleCloudPlatform/deploystack/commit/0cc9aa60a1fd71e18c7ca31ed96a801780abfc51))
* making sure tests work properly ([2d799e3](https://github.com/GoogleCloudPlatform/deploystack/commit/2d799e390ae6dffb66605c913907cb95f0b4a683))
* making sure that gcloud config project set is called ([4295f73](https://github.com/GoogleCloudPlatform/deploystack/commit/4295f739633f75e9360b76f97d662fc363cf73d1))
* making sure that project creation captures project number ([fb0d8d7](https://github.com/GoogleCloudPlatform/deploystack/commit/fb0d8d7cbb46a54c191e642846e39369ed128f77))
* making the progress bar work. ([11b0924](https://github.com/GoogleCloudPlatform/deploystack/commit/11b09240e043b7ce997346296bfc6e4379714ac4))
* making things more organized. ([2a69730](https://github.com/GoogleCloudPlatform/deploystack/commit/2a697304310a5ee0780b621d22854f74459bd109))
* massively changed the way color was rendered to make it simpler ([a5c652a](https://github.com/GoogleCloudPlatform/deploystack/commit/a5c652a9cdfc9a6820681ac734a7c1a51af8d300))
* more bug fixes ([5255ec0](https://github.com/GoogleCloudPlatform/deploystack/commit/5255ec06461588decdc0f4413740701419feab25))
* moved from parsing description files to having.a datastructure for products ([0758eb7](https://github.com/GoogleCloudPlatform/deploystack/commit/0758eb70801974f67863ccc8aa4eadf6be2e1ec8))
* moved project number retrieval to the actual flow. ([a546144](https://github.com/GoogleCloudPlatform/deploystack/commit/a546144dc53b5c0a02a9918d14daa96dd22288fc))
* needed to account for nils in casting operation ([4b0ea90](https://github.com/GoogleCloudPlatform/deploystack/commit/4b0ea90634dedde02b8ed34920639fdf341a944f))
* process is now exit(1) when user stops process ([ab4bbcf](https://github.com/GoogleCloudPlatform/deploystack/commit/ab4bbcf81370953a083b3224f01035401c9403e3))
* pruning directory files a little better. ([828d6ce](https://github.com/GoogleCloudPlatform/deploystack/commit/828d6ce19b80bd9f6463848008dea07b7a578eae))
* reduced the size of lists to 10 ([0182c5b](https://github.com/GoogleCloudPlatform/deploystack/commit/0182c5b992b6d359f80ab8e0340b798f17f3aa6a))
* removed a lot of unnecessary code ([315a948](https://github.com/GoogleCloudPlatform/deploystack/commit/315a948ca317de7f6634361fda93ea5146c36be8))
* removing generated file ([52e73fe](https://github.com/GoogleCloudPlatform/deploystack/commit/52e73fecb76a45e092ccb2fe47466e2c0366a549))
* resolves issue with color not showing in CloudShell ([df2a788](https://github.com/GoogleCloudPlatform/deploystack/commit/df2a788e2a4b7e9accb481e5b5c8c69fa8906b39))
* resolving issue where new accounts weren't getting billing attached ([fba7fad](https://github.com/GoogleCloudPlatform/deploystack/commit/fba7fadadeef1a7af510d4dfd58594a11d1e111e))
* restoreing full screen mode. ([3259419](https://github.com/GoogleCloudPlatform/deploystack/commit/3259419229d9abee0991a907350969def9ef5ede))
* testing color rendering for Cloud Shell ([9535e12](https://github.com/GoogleCloudPlatform/deploystack/commit/9535e12b8a0aa1c828b152fa751b9265697b84e1))
* testing if the alt screen is the problem on cloudshell ([8ce98f9](https://github.com/GoogleCloudPlatform/deploystack/commit/8ce98f999cfe31b61165180e3e389df2185c2ee4))
* testing the behavior of the queue more ([99d76fb](https://github.com/GoogleCloudPlatform/deploystack/commit/99d76fbf16c7bdf1b6d3d265e85ac718baa0f9eb))
* trying to correct a problem where Deploystack stops the terminal from working ([9b1d95e](https://github.com/GoogleCloudPlatform/deploystack/commit/9b1d95e1c57c3a9141280401518c74c14ed90811))
* trying to get all packages working properly ([5768e28](https://github.com/GoogleCloudPlatform/deploystack/commit/5768e2828529e17974feb8f1fe43feb706329f70))
* trying to update all the things to ensure proper dependencies ([e76fdae](https://github.com/GoogleCloudPlatform/deploystack/commit/e76fdae2123f6997268c7bc8a552e8f552e24a49))
* tweaking whitespace ([4b12003](https://github.com/GoogleCloudPlatform/deploystack/commit/4b12003b890966b8b4599325b2813202dc3d2004))
* update dependencies ([cb39106](https://github.com/GoogleCloudPlatform/deploystack/commit/cb391060f1b703dd38dec3233479a22e0dbc1130))
* updated colors to try and make it look better on Cloud Shell ([5d7f220](https://github.com/GoogleCloudPlatform/deploystack/commit/5d7f2206b3445d9e94430dec61df0b2e7db90d8a))
* updated tests ([00b392e](https://github.com/GoogleCloudPlatform/deploystack/commit/00b392e927603c46dedfb3418578256f3d6eb4a1))
* when you are just asking for the billing account, you don't need to attach it. ([0f9f547](https://github.com/GoogleCloudPlatform/deploystack/commit/0f9f5478f79d883ec7f857937e31a5bbcb83d3ea))
* working on a demo ui based on mock data ([cacabcc](https://github.com/GoogleCloudPlatform/deploystack/commit/cacabcce45d729eee47dee3ead038a485dc36e4e))

## [1.5.1](https://github.com/GoogleCloudPlatform/deploystack/compare/tui/v1.5.0...tui/v1.5.1) (2023-03-02)


### Bug Fixes

* cloud shell always thinks it's in dark mode ([5f36947](https://github.com/GoogleCloudPlatform/deploystack/commit/5f369479376aa285ac9061151c1f22b8058ff5fa))

## [1.5.0](https://github.com/GoogleCloudPlatform/deploystack/compare/tui/v1.4.0...tui/v1.5.0) (2023-02-28)


### Features

* added ability for more than one stack to exist in the same repo ([959b797](https://github.com/GoogleCloudPlatform/deploystack/commit/959b797478114da0d8d08e336605645f3bd02e56))


### Bug Fixes

* added an example for multi stack repos ([52ca81f](https://github.com/GoogleCloudPlatform/deploystack/commit/52ca81f34895b6a9a5b0ecc6e3941d9155ffaa1a))
* debugging might be screwing things up here. ([177cc54](https://github.com/GoogleCloudPlatform/deploystack/commit/177cc54fa769e927e76814feb28f1fca3d70b6c3))
* trying to correct a problem where Deploystack stops the terminal from working ([9b1d95e](https://github.com/GoogleCloudPlatform/deploystack/commit/9b1d95e1c57c3a9141280401518c74c14ed90811))

## [1.4.0](https://github.com/GoogleCloudPlatform/deploystack/compare/tui/v1.3.0...tui/v1.4.0) (2023-02-25)


### Features

* added authorsettings to replace hardsettings ([553979b](https://github.com/GoogleCloudPlatform/deploystack/commit/553979bf7520655c8467ce4ea4c6817f06447a4f))


### Bug Fixes

* made tui compatible with changes to config ([0ba76e4](https://github.com/GoogleCloudPlatform/deploystack/commit/0ba76e4dcf5270b9ac1cc71ceb96fcc395d5fa3b))

## [1.3.0](https://github.com/GoogleCloudPlatform/deploystack/compare/tui/v1.2.0...tui/v1.3.0) (2023-02-24)


### Features

* updated tui to work with new configuration settings ([6a5a17a](https://github.com/GoogleCloudPlatform/deploystack/commit/6a5a17ac6fc097dcd563a5a5a0638972769b9e5d))


### Bug Fixes

* list items limited to 50 chars and defaults now display correctly ([46b2e3b](https://github.com/GoogleCloudPlatform/deploystack/commit/46b2e3b8288c432e6ebff4dfefd756e9040dd38f))
* make default options display normally if there are less than 10 of them. ([0cc9aa6](https://github.com/GoogleCloudPlatform/deploystack/commit/0cc9aa60a1fd71e18c7ca31ed96a801780abfc51))
* reduced the size of lists to 10 ([0182c5b](https://github.com/GoogleCloudPlatform/deploystack/commit/0182c5b992b6d359f80ab8e0340b798f17f3aa6a))

## [1.2.0](https://github.com/GoogleCloudPlatform/deploystack/compare/tui-v1.1.1...tui/v1.2.0) (2023-02-21)


### Features

* added a warning about slow lookups in a new project ([855ac53](https://github.com/GoogleCloudPlatform/deploystack/commit/855ac537cf7470ba4df3788e229790efdb3e7cbd))
* added caching to speed up these requests ([a0aaad8](https://github.com/GoogleCloudPlatform/deploystack/commit/a0aaad877b643181757d84e37621109f9e2f9d5a))
* added the progress bar. ([4e6a76c](https://github.com/GoogleCloudPlatform/deploystack/commit/4e6a76cead93920c619b6798b7caa3b743a58605))
* build a "don't put into settings flag" into the base configuration of pages ([d616a13](https://github.com/GoogleCloudPlatform/deploystack/commit/d616a138f29a10f56f7f14efb0636b34358f3621))
* created new package to remove threat of circular dependency ([8b0738e](https://github.com/GoogleCloudPlatform/deploystack/commit/8b0738e28a839d2f9a21cb4c880ddd382d9017e2))
* forgot to add prependProject setting. Working now ([a6b2e25](https://github.com/GoogleCloudPlatform/deploystack/commit/a6b2e25c2980b459d715e82ba4ec5c1cbbdac735))
* implemented a back function ([c65badc](https://github.com/GoogleCloudPlatform/deploystack/commit/c65badc1f5a2041e725ba43015ff8af8438948f8))


### Bug Fixes

* actually fixing the problem with project numbers ([5175fea](https://github.com/GoogleCloudPlatform/deploystack/commit/5175fea1dedd89b5d5da9c4a8174e4c9c7f92ecb))
* added billing account capture and made sure that was working ([d7c0dc3](https://github.com/GoogleCloudPlatform/deploystack/commit/d7c0dc3c51699da801e7eaccbcca42b453170723))
* adding a clean preprocessor for the last screen ([c3941d7](https://github.com/GoogleCloudPlatform/deploystack/commit/c3941d792be3663847ce42e06396f4da7806a8da))
* better default picker style ([4f962e5](https://github.com/GoogleCloudPlatform/deploystack/commit/4f962e53f04a886c4414bd8070c103a116987a7c))
* checking in todays changes ([18e3775](https://github.com/GoogleCloudPlatform/deploystack/commit/18e37755d37312bd26735c19dc163bd498dd2b5c))
* converted spinner to use dsStyles to render color on Cloud Shell ([a7456b7](https://github.com/GoogleCloudPlatform/deploystack/commit/a7456b743a79585a481bd66242737ec31f7748d8))
* corrected issue where domain contact info was not being refreshed properly ([546af35](https://github.com/GoogleCloudPlatform/deploystack/commit/546af355c1a795017accddd23271e3a429db3480))
* correcting bug where Project forms looked bad. ([1a2e336](https://github.com/GoogleCloudPlatform/deploystack/commit/1a2e336cd7889d6b17d9d7050b7e5715496a4566))
* correcting dependency in test ([1431a27](https://github.com/GoogleCloudPlatform/deploystack/commit/1431a27bd75b21c7647cb6132f8355bce1c9872a))
* didn't like short form ([d4219c3](https://github.com/GoogleCloudPlatform/deploystack/commit/d4219c38c537dd4623f7c78a51fdd28855fe2dd8))
* fixed an issue where when the default value isn't set, the list isn't seleced ([6d69e80](https://github.com/GoogleCloudPlatform/deploystack/commit/6d69e80ccaef0d09c95859ddfe319e58c408cc2a))
* fixes dealing with changes in color rendering ([26bd638](https://github.com/GoogleCloudPlatform/deploystack/commit/26bd63893c6cfc23929d4517c1457a4fa06698a3))
* fixing breaking test ([9784da5](https://github.com/GoogleCloudPlatform/deploystack/commit/9784da55dbb6ec641b04fcdd0bbb74c942029a48))
* fixing error where ctr-c on the exit page causes a panic ([3974a8b](https://github.com/GoogleCloudPlatform/deploystack/commit/3974a8bbf92387104e5f644d4054cf3414b5560a))
* getting all the version stuff working properly ([140556a](https://github.com/GoogleCloudPlatform/deploystack/commit/140556afab66aad5233032c402072641f561d223))
* getting billing to work properly ([9894e2b](https://github.com/GoogleCloudPlatform/deploystack/commit/9894e2bcea439b197755ab9a324dbbbe7ed71e09))
* getting rid of old calls. ([a391e39](https://github.com/GoogleCloudPlatform/deploystack/commit/a391e39aa119211e49570f4f1e6a67b804a87325))
* got picker default values working properly ([ef623d7](https://github.com/GoogleCloudPlatform/deploystack/commit/ef623d748107ee0cb13e30082926bbeeb5015edb))
* initial import of tui code ([e53bb85](https://github.com/GoogleCloudPlatform/deploystack/commit/e53bb85709de1bc8c025ce0948b6cec348cdaab6))
* just some tweaks to improving linting ([5c716d7](https://github.com/GoogleCloudPlatform/deploystack/commit/5c716d7d05f1e7235f9319edacd5578669c081ef))
* made currentProject a queue value instead of a global one ([c3bfb7b](https://github.com/GoogleCloudPlatform/deploystack/commit/c3bfb7bcbd46d781e84612dc4327204f93ba8897))
* making sure tests work properly ([2d799e3](https://github.com/GoogleCloudPlatform/deploystack/commit/2d799e390ae6dffb66605c913907cb95f0b4a683))
* making sure that gcloud config project set is called ([4295f73](https://github.com/GoogleCloudPlatform/deploystack/commit/4295f739633f75e9360b76f97d662fc363cf73d1))
* making sure that project creation captures project number ([fb0d8d7](https://github.com/GoogleCloudPlatform/deploystack/commit/fb0d8d7cbb46a54c191e642846e39369ed128f77))
* making the progress bar work. ([11b0924](https://github.com/GoogleCloudPlatform/deploystack/commit/11b09240e043b7ce997346296bfc6e4379714ac4))
* making things more organized. ([2a69730](https://github.com/GoogleCloudPlatform/deploystack/commit/2a697304310a5ee0780b621d22854f74459bd109))
* massively changed the way color was rendered to make it simpler ([a5c652a](https://github.com/GoogleCloudPlatform/deploystack/commit/a5c652a9cdfc9a6820681ac734a7c1a51af8d300))
* more bug fixes ([5255ec0](https://github.com/GoogleCloudPlatform/deploystack/commit/5255ec06461588decdc0f4413740701419feab25))
* moved from parsing description files to having.a datastructure for products ([0758eb7](https://github.com/GoogleCloudPlatform/deploystack/commit/0758eb70801974f67863ccc8aa4eadf6be2e1ec8))
* moved project number retrieval to the actual flow. ([a546144](https://github.com/GoogleCloudPlatform/deploystack/commit/a546144dc53b5c0a02a9918d14daa96dd22288fc))
* needed to account for nils in casting operation ([4b0ea90](https://github.com/GoogleCloudPlatform/deploystack/commit/4b0ea90634dedde02b8ed34920639fdf341a944f))
* process is now exit(1) when user stops process ([ab4bbcf](https://github.com/GoogleCloudPlatform/deploystack/commit/ab4bbcf81370953a083b3224f01035401c9403e3))
* pruning directory files a little better. ([828d6ce](https://github.com/GoogleCloudPlatform/deploystack/commit/828d6ce19b80bd9f6463848008dea07b7a578eae))
* removed a lot of unnecessary code ([315a948](https://github.com/GoogleCloudPlatform/deploystack/commit/315a948ca317de7f6634361fda93ea5146c36be8))
* removing generated file ([52e73fe](https://github.com/GoogleCloudPlatform/deploystack/commit/52e73fecb76a45e092ccb2fe47466e2c0366a549))
* resolves issue with color not showing in CloudShell ([df2a788](https://github.com/GoogleCloudPlatform/deploystack/commit/df2a788e2a4b7e9accb481e5b5c8c69fa8906b39))
* resolving issue where new accounts weren't getting billing attached ([fba7fad](https://github.com/GoogleCloudPlatform/deploystack/commit/fba7fadadeef1a7af510d4dfd58594a11d1e111e))
* restoreing full screen mode. ([3259419](https://github.com/GoogleCloudPlatform/deploystack/commit/3259419229d9abee0991a907350969def9ef5ede))
* testing color rendering for Cloud Shell ([9535e12](https://github.com/GoogleCloudPlatform/deploystack/commit/9535e12b8a0aa1c828b152fa751b9265697b84e1))
* testing if the alt screen is the problem on cloudshell ([8ce98f9](https://github.com/GoogleCloudPlatform/deploystack/commit/8ce98f999cfe31b61165180e3e389df2185c2ee4))
* testing the behavior of the queue more ([99d76fb](https://github.com/GoogleCloudPlatform/deploystack/commit/99d76fbf16c7bdf1b6d3d265e85ac718baa0f9eb))
* trying to get all packages working properly ([5768e28](https://github.com/GoogleCloudPlatform/deploystack/commit/5768e2828529e17974feb8f1fe43feb706329f70))
* trying to update all the things to ensure proper dependencies ([e76fdae](https://github.com/GoogleCloudPlatform/deploystack/commit/e76fdae2123f6997268c7bc8a552e8f552e24a49))
* tweaking whitespace ([4b12003](https://github.com/GoogleCloudPlatform/deploystack/commit/4b12003b890966b8b4599325b2813202dc3d2004))
* update dependencies ([cb39106](https://github.com/GoogleCloudPlatform/deploystack/commit/cb391060f1b703dd38dec3233479a22e0dbc1130))
* updated colors to try and make it look better on Cloud Shell ([5d7f220](https://github.com/GoogleCloudPlatform/deploystack/commit/5d7f2206b3445d9e94430dec61df0b2e7db90d8a))
* updated tests ([00b392e](https://github.com/GoogleCloudPlatform/deploystack/commit/00b392e927603c46dedfb3418578256f3d6eb4a1))
* when you are just asking for the billing account, you don't need to attach it. ([0f9f547](https://github.com/GoogleCloudPlatform/deploystack/commit/0f9f5478f79d883ec7f857937e31a5bbcb83d3ea))
* working on a demo ui based on mock data ([cacabcc](https://github.com/GoogleCloudPlatform/deploystack/commit/cacabcce45d729eee47dee3ead038a485dc36e4e))
