# Changelog

## [3.55.0](https://github.com/openai/openai-go/compare/v3.54.0...v3.55.0) (2026-09-02)


### Features

* **api:** add compute_units to Responses and Chat Completions usage ([#866](https://github.com/openai/openai-go/issues/866)) ([4d06294](https://github.com/openai/openai-go/commit/4d062949c62507e56514af8c7beb186dc09ac075))
* **api:** update usage APIs and documentation ([#876](https://github.com/openai/openai-go/issues/876)) ([d8a7e6c](https://github.com/openai/openai-go/commit/d8a7e6c70d1bf67e744ea8e4fb617b622372ca0e))


### Bug Fixes

* isolate provider dependencies from root consumers ([#874](https://github.com/openai/openai-go/issues/874)) ([8de3444](https://github.com/openai/openai-go/commit/8de344439ca10236fc958a468bfbd5807737a15f))

## [3.54.0](https://github.com/openai/openai-go/compare/v3.53.0...v3.54.0) (2026-08-27)


### Features

* **api:** make function call output call IDs optional ([#864](https://github.com/openai/openai-go/issues/864)) ([7095c46](https://github.com/openai/openai-go/commit/7095c4615c3033d78c460e7bac646ba8062e0e16))


### Bug Fixes

* **auth:** harden X.509 workload identity contract ([#863](https://github.com/openai/openai-go/issues/863)) ([adefcb3](https://github.com/openai/openai-go/commit/adefcb377ad07600d995d964c80f15d13fefd468))

## [3.53.0](https://github.com/openai/openai-go/compare/v3.52.0...v3.53.0) (2026-08-25)


### Features

* **api:** Add obfuscation field to ChatCompletionChunk ([#813](https://github.com/openai/openai-go/issues/813)) ([4b08348](https://github.com/openai/openai-go/commit/4b08348ed0d61dc0d1f2a7b75f9808e7a51ab9db))
* **api:** add project residency and usage quantity units ([#853](https://github.com/openai/openai-go/issues/853)) ([2fd626c](https://github.com/openai/openai-go/commit/2fd626c905da2fee3d9ce0d1afb30f6d92fe6099))
* **api:** add Realtime call creation and recursive filters ([#860](https://github.com/openai/openai-go/issues/860)) ([c396b57](https://github.com/openai/openai-go/commit/c396b579f36ab8e154b9f08390a70e6fcfe15217))
* **auth:** add caller-attested X.509 transport capability ([#855](https://github.com/openai/openai-go/issues/855)) ([64b2ef2](https://github.com/openai/openai-go/commit/64b2ef2c2b1c0799b382b4f9d9e3bb79cedbcd4a))
* **auth:** complete X.509 token lifecycle and setup guidance ([#858](https://github.com/openai/openai-go/issues/858)) ([3e3a4e0](https://github.com/openai/openai-go/commit/3e3a4e0a73ff392bb8c4c15523b0a22e12f07b0e))
* **auth:** implement pinned X.509 token exchange ([#856](https://github.com/openai/openai-go/issues/856)) ([a84eb49](https://github.com/openai/openai-go/commit/a84eb49d65391ab4618b5c653c73332a591981f5))
* **option:** add attested X.509 workload identity authentication ([#857](https://github.com/openai/openai-go/issues/857)) ([b509c12](https://github.com/openai/openai-go/commit/b509c12cc12a22886f304bf62374c3667020bed2))
* support named data-residency endpoints ([#805](https://github.com/openai/openai-go/issues/805)) ([6e72880](https://github.com/openai/openai-go/commit/6e7288096828717805cfccbea3b15f054326ff8e))


### Bug Fixes

* **apijson:** preserve unknown union array variants ([#768](https://github.com/openai/openai-go/issues/768)) ([0ab7ee6](https://github.com/openai/openai-go/commit/0ab7ee6680fe3887e9eea545706600a8a890549b))
* **apijson:** reject non-object root payloads ([#769](https://github.com/openai/openai-go/issues/769)) ([2959fbd](https://github.com/openai/openai-go/commit/2959fbd636ee8dacbbb65416ec2f18bea55e136e))
* **azure:** isolate provider authentication defaults ([#835](https://github.com/openai/openai-go/issues/835)) ([4e5a07a](https://github.com/openai/openai-go/commit/4e5a07a6a205761650d73593fe75cd79ea7a25ec))
* **azure:** require secure credential transport ([#821](https://github.com/openai/openai-go/issues/821)) ([1dcf764](https://github.com/openai/openai-go/commit/1dcf7643c7661de54b00d189aa421b26db36bc8f))
* **azure:** validate deployment routing bodies ([#826](https://github.com/openai/openai-go/issues/826)) ([54189a0](https://github.com/openai/openai-go/commit/54189a0b76d51f440bfb842547617beea3dec21d))
* bound retry and polling delays ([#828](https://github.com/openai/openai-go/issues/828)) ([6d39b8b](https://github.com/openai/openai-go/commit/6d39b8b554a21202dac16dd5bbdc158ecbb09b06))
* don't drop tool_calls when content is empty but present ([#757](https://github.com/openai/openai-go/issues/757)) ([4a73240](https://github.com/openai/openai-go/commit/4a73240556735716521d97c24df0512d2028d715)), closes [#756](https://github.com/openai/openai-go/issues/756)
* enforce configured request origins ([#823](https://github.com/openai/openai-go/issues/823)) ([a71b2d1](https://github.com/openai/openai-go/commit/a71b2d17ed0c5b5e9f1ba1817f9b86a1ac7f2b5f))
* harden multipart file metadata encoding ([#825](https://github.com/openai/openai-go/issues/825)) ([a4632e5](https://github.com/openai/openai-go/commit/a4632e570c087bc5b12e3412982fc6d1cc36ac69))
* harden workload identity transport and request lifecycles ([#861](https://github.com/openai/openai-go/issues/861)) ([4ce8802](https://github.com/openai/openai-go/commit/4ce8802217a8bdf0ee9959eb292e3004e888fce7))
* harden X.509 refresh races and empty-body replay ([#859](https://github.com/openai/openai-go/issues/859)) ([f95adf1](https://github.com/openai/openai-go/commit/f95adf17f6894eaf8d4e3a630cfba703615a4dd6))
* ignore empty tool call deltas in accumulator ([#688](https://github.com/openai/openai-go/issues/688)) ([c2c8590](https://github.com/openai/openai-go/commit/c2c85905588a6f9cea14964cdf70105f9e5ab691))
* make chat stream accumulation linear ([#819](https://github.com/openai/openai-go/issues/819)) ([0b7382a](https://github.com/openai/openai-go/commit/0b7382a4c1bc2b917b5939eb88005dbd654a472e))
* make debug logging metadata-only ([#815](https://github.com/openai/openai-go/issues/815)) ([2e0743f](https://github.com/openai/openai-go/commit/2e0743fd569b477247379ff51d80028de3dc450f))
* normalize registered stream decoder content types ([#790](https://github.com/openai/openai-go/issues/790)) ([a689437](https://github.com/openai/openai-go/commit/a689437a709ca07555386f273d2e048170f0dcc3))
* preserve request options across pagination ([#820](https://github.com/openai/openai-go/issues/820)) ([ee97c12](https://github.com/openai/openai-go/commit/ee97c12e413e3a1299196e63fee01a85bc7310ac))
* **responses:** decode shell output delta stream events ([#803](https://github.com/openai/openai-go/issues/803)) ([bb12b7a](https://github.com/openai/openai-go/commit/bb12b7a7d292ea6f4e2c8a662a2d0c3e1aeef8e9))
* **sdk:** lock and verify the Steady mock server ([#811](https://github.com/openai/openai-go/issues/811)) ([37c0f0a](https://github.com/openai/openai-go/commit/37c0f0ab07fbe5f5b2775c12ca3500acd0caad8c))
* **ssestream:** finalize stream resources promptly ([#839](https://github.com/openai/openai-go/issues/839)) ([71a7aef](https://github.com/openai/openai-go/commit/71a7aefbe265cbff6088e617fef741e8dfbb1061))
* track stream accumulator events for every choice ([#806](https://github.com/openai/openai-go/issues/806)) ([8fa5f45](https://github.com/openai/openai-go/commit/8fa5f4504641b33280c3a9cefda82e4cbb668c35))
* validate stream accumulator indices ([#822](https://github.com/openai/openai-go/issues/822)) ([02959ba](https://github.com/openai/openai-go/commit/02959ba995fbe8c98dad306541ca3ea9ece14284))
* **webhooks:** honor request-scoped verification options ([#814](https://github.com/openai/openai-go/issues/814)) ([da53088](https://github.com/openai/openai-go/commit/da53088d2320ecf44d8eb75d0ccbdeb5e8c12210))
* **webhooks:** validate prefixed signing secrets ([#824](https://github.com/openai/openai-go/issues/824)) ([c88e19a](https://github.com/openai/openai-go/commit/c88e19a1c266ef8fbc140969c862c1eb635142d8))


### Chores

* **api:** Add file header comments to source files ([#809](https://github.com/openai/openai-go/issues/809)) ([557cc1c](https://github.com/openai/openai-go/commit/557cc1ceed473541eadab147c9dbc468df90caab))
* **api:** Clarify image background transparency docs ([#831](https://github.com/openai/openai-go/issues/831)) ([6dfcd9b](https://github.com/openai/openai-go/commit/6dfcd9b01bc201830df3aff3820a1664ea05e21b))
* **api:** Document allowed image generation model values ([#816](https://github.com/openai/openai-go/issues/816)) ([c70532e](https://github.com/openai/openai-go/commit/c70532ec68987507eaacb5257270ae24aea5608d))
* **api:** update generation metadata ([#842](https://github.com/openai/openai-go/issues/842)) ([ceb67ae](https://github.com/openai/openai-go/commit/ceb67ae80c8bcfe26759ee6fb764d6e265e4304c))
* set a 2,000-line custom-code budget ([#840](https://github.com/openai/openai-go/issues/840)) ([31ce143](https://github.com/openai/openai-go/commit/31ce1430b9848d7c3716b081a72faf2900154186))
* **shared:** restore generated alias placement ([#830](https://github.com/openai/openai-go/issues/830)) ([7d0e09e](https://github.com/openai/openai-go/commit/7d0e09e4bec1d20a888fc9a877ff214e608792cd))


### Documentation

* add Markdown-only Go SDK review skill ([#795](https://github.com/openai/openai-go/issues/795)) ([3e82f23](https://github.com/openai/openai-go/commit/3e82f23220a804fbc5b55bd517cb7cd94c6eb8b5))
* add recurring improvement skill ([#818](https://github.com/openai/openai-go/issues/818)) ([cb08fd8](https://github.com/openai/openai-go/commit/cb08fd888dd45446cf037b1ac2caf7e7f6f765bf))
* clean up changelog release notes ([#696](https://github.com/openai/openai-go/issues/696)) ([5ba3879](https://github.com/openai/openai-go/commit/5ba3879a9277f3d833d143ab9b559a9dbb2f590b))
* improve webhook request handling examples ([#817](https://github.com/openai/openai-go/issues/817)) ([c152406](https://github.com/openai/openai-go/commit/c15240667e236baf12eed21e81922c3b655f8c8e))
* migrate structured outputs example to Responses API ([#804](https://github.com/openai/openai-go/issues/804)) ([dd30934](https://github.com/openai/openai-go/commit/dd309349035acf67d81ce9b2b11309e29a6a99f0)), closes [#482](https://github.com/openai/openai-go/issues/482)
* strengthen Go SDK security contribution guidance ([#797](https://github.com/openai/openai-go/issues/797)) ([b4ca1a7](https://github.com/openai/openai-go/commit/b4ca1a75215e59c13542eed462af9ad0549e58c0))


### Refactors

* **chat:** isolate message conversion ([#836](https://github.com/openai/openai-go/issues/836)) ([3175b15](https://github.com/openai/openai-go/commit/3175b151b474b7a14be16688f2a470252e1e0d27))
* **chat:** reuse generated assistant constructor ([#834](https://github.com/openai/openai-go/issues/834)) ([7a8f196](https://github.com/openai/openai-go/commit/7a8f19666cf536dc9aa7a17b75fb3be32cb41a2d))
* **responses:** isolate output text aggregation ([#838](https://github.com/openai/openai-go/issues/838)) ([d38d045](https://github.com/openai/openai-go/commit/d38d04517a2c47a70cb58c740270fb7e0ce3043d))
* **vector-stores:** isolate batch orchestration ([#833](https://github.com/openai/openai-go/issues/833)) ([a2fa6c1](https://github.com/openai/openai-go/commit/a2fa6c18bfddac98262f0c4db7c2ff073b54c882))
* **vector-stores:** isolate file orchestration ([#832](https://github.com/openai/openai-go/issues/832)) ([803983f](https://github.com/openai/openai-go/commit/803983f0ab5606e50f036d98ebec48bcf959c1f9))
* **webhooks:** isolate verification helpers ([#837](https://github.com/openai/openai-go/issues/837)) ([401eaf3](https://github.com/openai/openai-go/commit/401eaf328ba7a9ad48a0c8a32820b58b54e0a3ef))


### Build System

* **deps:** bump actions/checkout from 6.0.2 to 7.0.1 ([#848](https://github.com/openai/openai-go/issues/848)) ([fbddc21](https://github.com/openai/openai-go/commit/fbddc21d13c33e55222dbc5a34df4996136b6a0e))
* **deps:** bump actions/github-script from 7.1.0 to 9.0.0 ([#849](https://github.com/openai/openai-go/issues/849)) ([9d7c09e](https://github.com/openai/openai-go/commit/9d7c09e7cb761486cfba37bd11e18f17860d9ae6))
* **deps:** bump github.com/aws/aws-sdk-go-v2/config in / ([#850](https://github.com/openai/openai-go/issues/850)) ([a5fedc0](https://github.com/openai/openai-go/commit/a5fedc01b16f88a082aad74c49326f63d374ad48))
* **deps:** bump github.com/Azure/azure-sdk-for-go/sdk/azcore in / ([#845](https://github.com/openai/openai-go/issues/845)) ([f0079ac](https://github.com/openai/openai-go/commit/f0079ac64a790f49b9c8bcec5cc062e19dc3e6b0))
* **deps:** bump github.com/ebitengine/oto/v3 in /examples ([#851](https://github.com/openai/openai-go/issues/851)) ([65e8a8f](https://github.com/openai/openai-go/commit/65e8a8f024c019b6bf156b65753651458ab9805c))
* **deps:** bump the codeql group with 2 updates ([#847](https://github.com/openai/openai-go/issues/847)) ([29b2a11](https://github.com/openai/openai-go/commit/29b2a1123ade1d9d587501c58911341fd212bb0d))

## [3.52.0](https://github.com/openai/openai-go/compare/v3.51.0...v3.52.0) (2026-08-17)


### Features

* add Bedrock Runtime endpoint support (SDK-290) ([#793](https://github.com/openai/openai-go/issues/793)) ([8e8444e](https://github.com/openai/openai-go/commit/8e8444ecf7918b3370356fa0d8478788e2c8196d))

## [3.51.0](https://github.com/openai/openai-go/compare/v3.50.0...v3.51.0) (2026-08-14)


### Features

* **api:** Add model constants for Daybreak and gpt-5.6-cyber ([df00e97](https://github.com/openai/openai-go/commit/df00e9760b4584dc0986b7383ecb843a9cf67789))
* **api:** add WebSocket stream IDs ([#776](https://github.com/openai/openai-go/issues/776)) ([0b6316e](https://github.com/openai/openai-go/commit/0b6316e74982b8403aac3854405ed90153dfb2e0))
* **api:** add workload identity access token issued event ([#766](https://github.com/openai/openai-go/issues/766)) ([c61af30](https://github.com/openai/openai-go/commit/c61af303735128ab487ec1486cdf4983047bcbc0))
* **api:** deprecate Sora video APIs ([#775](https://github.com/openai/openai-go/issues/775)) ([602c5e4](https://github.com/openai/openai-go/commit/602c5e48c7c51827cd656b4f613f3680bd858a7c))
* **api:** Ultrafast tier, structured MCP and websocket errors, separate websocket events ([#789](https://github.com/openai/openai-go/issues/789)) ([f8b63fd](https://github.com/openai/openai-go/commit/f8b63fd68a81b586a8227e1032d173b039bd888e))


### Bug Fixes

* **api:** Add model constants for Daybreak and gpt-5.6-cyber ([#763](https://github.com/openai/openai-go/issues/763)) ([df00e97](https://github.com/openai/openai-go/commit/df00e9760b4584dc0986b7383ecb843a9cf67789))
* **api:** clarify audio upload metadata requirements ([#764](https://github.com/openai/openai-go/issues/764)) ([e75e3aa](https://github.com/openai/openai-go/commit/e75e3aaf598f6257445984cfedaf3dd40d790a29))
* handle ignored errors in generated and runtime Go code ([#778](https://github.com/openai/openai-go/issues/778)) ([ba5f249](https://github.com/openai/openai-go/commit/ba5f249454f32f1ba16a184c2290e64b32cd1bee))
* recover approved Go correctness analyzers after stacked merges ([#782](https://github.com/openai/openai-go/issues/782)) ([a98cf53](https://github.com/openai/openai-go/commit/a98cf53f7a98a977ab90750f3365c65b276df38d))


### Chores

* **api:** No public API change; update file header comments ([#755](https://github.com/openai/openai-go/issues/755)) ([d0a42cc](https://github.com/openai/openai-go/commit/d0a42cc753ffb56b7f268ecc2b9975ca58203215))
* enforce Go lint suppression hygiene ([#771](https://github.com/openai/openai-go/issues/771)) ([7a3bb22](https://github.com/openai/openai-go/commit/7a3bb22e4b29b489efaecc9e539c44be8fa52d2d))
* enforce gofmt across the repository ([#770](https://github.com/openai/openai-go/issues/770)) ([41b2f1f](https://github.com/openai/openai-go/commit/41b2f1f4d7b30399d143e6cb456b6c94ce7f51e2)), closes [#771](https://github.com/openai/openai-go/issues/771)
* remove Stainless attribution and infrastructure ([#765](https://github.com/openai/openai-go/issues/765)) ([9b3d395](https://github.com/openai/openai-go/commit/9b3d395105c3272d0d703f8e1e3e37776f7b42a5))


### Documentation

* **api:** describe response stream event unions ([#783](https://github.com/openai/openai-go/issues/783)) ([4cc4021](https://github.com/openai/openai-go/commit/4cc4021ca46415326b78377b07f221e11f979b4d))
* define Go code quality policy ([#767](https://github.com/openai/openai-go/issues/767)) ([4119aa3](https://github.com/openai/openai-go/commit/4119aa3b19feb7c20757d78d4fd6ed0c09233a9b))
* define Go ignored-error handling policy ([#777](https://github.com/openai/openai-go/issues/777)) ([1941f81](https://github.com/openai/openai-go/commit/1941f81fa2b3353feec2d39c5dd7fd2b85e58780))


### Build System

* **deps:** bump github.com/aws/aws-sdk-go-v2/config in / ([#759](https://github.com/openai/openai-go/issues/759)) ([e8dfbe0](https://github.com/openai/openai-go/commit/e8dfbe097d282f0c0aee74f6aed00a6b4decb7c0))
* **deps:** bump the codeql group with 2 updates ([#760](https://github.com/openai/openai-go/issues/760)) ([9bbfebf](https://github.com/openai/openai-go/commit/9bbfebf71e721d0575ccbd56d0aaeacc4a4cef10))
* **deps:** enforce an eight-day dependency cooldown ([#762](https://github.com/openai/openai-go/issues/762)) ([744e64f](https://github.com/openai/openai-go/commit/744e64ff17ba988371534849e22555255d90bb56))

## [3.50.0](https://github.com/openai/openai-go/compare/v3.49.0...v3.50.0) (2026-08-03)


### Features

* **api:** promote Go SDK changes ([#748](https://github.com/openai/openai-go/issues/748)) ([9e4b4d5](https://github.com/openai/openai-go/commit/9e4b4d5a5ccadc94ea8b8a54148e7a0e52797e94))


### Build System

* **deps:** bump github.com/aws/aws-sdk-go-v2/config in / ([#745](https://github.com/openai/openai-go/issues/745)) ([af89b42](https://github.com/openai/openai-go/commit/af89b4276381d98193419f46429f79d74b490ec9))
* **deps:** bump the codeql group with 2 updates ([#746](https://github.com/openai/openai-go/issues/746)) ([41c8a70](https://github.com/openai/openai-go/commit/41c8a70fef4c2f7067a7eb6bc095c94d1782a2fe))

## 3.49.0 (2026-07-31)

Full Changelog: [v3.48.0...v3.49.0](https://github.com/openai/openai-go/compare/v3.48.0...v3.49.0)

### Features

* **api:** content provenance checks ([5a46c55](https://github.com/openai/openai-go/commit/5a46c55e1714343afcf8b99c69c80fb67801e2a8))


### Documentation

* demonstrate mutual TLS with custom HTTP client ([#741](https://github.com/openai/openai-go/issues/741)) ([e862c49](https://github.com/openai/openai-go/commit/e862c4964615058d518d53a15277c1a2d5945f02))

## 3.48.0 (2026-07-30)

Full Changelog: [v3.47.0...v3.48.0](https://github.com/openai/openai-go/compare/v3.47.0...v3.48.0)

### Features

* **api:** fast tier ([cffac28](https://github.com/openai/openai-go/commit/cffac28fd20783bb48a5e69e102e6f327f2adb83))

## 3.47.0 (2026-07-28)

Full Changelog: [v3.46.0...v3.47.0](https://github.com/openai/openai-go/compare/v3.46.0...v3.47.0)

### Features

* **api:** transcription model updates ([1707706](https://github.com/openai/openai-go/commit/1707706e0ac9d8c91d781e218ee3e6c9ac4fc41d))


### Bug Fixes

* windows platform header normalization ([#674](https://github.com/openai/openai-go/issues/674)) ([706d3a8](https://github.com/openai/openai-go/commit/706d3a8c5c76a211a8cf3db6322eb47ef3713d6f))


### Build System

* **deps:** bump actions/download-artifact from 7.0.0 to 8.0.1 ([#739](https://github.com/openai/openai-go/issues/739)) ([2c2a363](https://github.com/openai/openai-go/commit/2c2a3637f2c50a01a7b2558c41c68913ca693e81))
* **deps:** bump actions/setup-go from 6.4.0 to 7.0.0 ([#737](https://github.com/openai/openai-go/issues/737)) ([1fcc7ce](https://github.com/openai/openai-go/commit/1fcc7cef36599d8a0d8f3858fa6c4fe6d27bf241))
* **deps:** bump actions/setup-node from 6.3.0 to 7.0.0 ([#734](https://github.com/openai/openai-go/issues/734)) ([15562ec](https://github.com/openai/openai-go/commit/15562ec9fc84b9e26b6e3cad42f7284405b5b36e))
* **deps:** bump actions/upload-artifact from 6.0.0 to 7.0.1 ([#738](https://github.com/openai/openai-go/issues/738)) ([16fcc6a](https://github.com/openai/openai-go/commit/16fcc6a835f0791764841f5a6fc3f0d00b6c1e4d))
* **deps:** bump github.com/aws/aws-sdk-go-v2/config in / ([#733](https://github.com/openai/openai-go/issues/733)) ([279abc1](https://github.com/openai/openai-go/commit/279abc1d002c6a87a07470aeae1c3da1e35cfb41))
* **deps:** bump github/codeql-action from 4.37.2 to 4.37.3 ([#736](https://github.com/openai/openai-go/issues/736)) ([6032a46](https://github.com/openai/openai-go/commit/6032a4681950b0bc426e5528b5ca4d2a2590b815))

## 3.46.0 (2026-07-23)

Full Changelog: [v3.45.0...v3.46.0](https://github.com/openai/openai-go/compare/v3.45.0...v3.46.0)

### Features

* **api:** add support for `spend_limit` admin apis ([8de8d9b](https://github.com/openai/openai-go/commit/8de8d9b41f2ad94562b42643a06c85a6c7e3a8f4))


### Bug Fixes

* **api:** safely format project spend limit paths ([ac981e4](https://github.com/openai/openai-go/commit/ac981e4c1feb83e1be1e7e26055819ccc3eda820))

## 3.45.0 (2026-07-22)

Full Changelog: [v3.44.0...v3.45.0](https://github.com/openai/openai-go/compare/v3.44.0...v3.45.0)

### Go version requirement

openai-go v3.45.0 requires Go 1.25 or newer. This release upgrades
`golang.org/x/net` to v0.57.0 to address security findings and aligns the SDK
with supported Go toolchains. v3.44.0 is the final release that builds with Go
1.22–1.24; older SDK releases do not receive guaranteed security backports. No
exported SDK API or `/v3` import path changes are included.

### Features

* **api:** accept `None` for prompt_cache_key/safety_identifier ([711460e](https://github.com/openai/openai-go/commit/711460ea11963cd250be5e420269b25488b02ef9))
* **stlc:** configurable CI runner and private-production-repo support in workflow templates ([6669f50](https://github.com/openai/openai-go/commit/6669f5011fc02c2160024d26ab6f9240560a778d))

## 3.44.0 (2026-07-17)

Full Changelog: [v3.43.0...v3.44.0](https://github.com/openai/openai-go/compare/v3.43.0...v3.44.0)

### Features

* **api:** /organization/projects/{project_id}/service_accounts/{service_account_id}/api_keys" endpoint ([28d2a6a](https://github.com/openai/openai-go/commit/28d2a6ab978f78d10b29087e70539ed230340274))
* **api:** manual updates ([a531d4a](https://github.com/openai/openai-go/commit/a531d4a994c75177a7ed3edaf0c909db188147ce))
* **api:** manual updates ([8073e1e](https://github.com/openai/openai-go/commit/8073e1e1a20907ae89433008157c238bccd4eca9))


### Bug Fixes

* **api:** format service account API key path ([999e165](https://github.com/openai/openai-go/commit/999e16536e7dc5ea4040d13b94c413cfc9df98ec))
* format service account API key path ([9fb19bb](https://github.com/openai/openai-go/commit/9fb19bbb13c067a12c14153f8ba5e57b233cd927))

## 3.43.0 (2026-07-14)

Full Changelog: [v3.42.0...v3.43.0](https://github.com/openai/openai-go/compare/v3.42.0...v3.43.0)

### Features

* **api:** add owner_project_access to APIKeyListParams ([df2e35d](https://github.com/openai/openai-go/commit/df2e35dc39e43e8da55f1006084d27a1041005fc))


### Bug Fixes

* **example:** set the 'include usage' field to true under stream options to get accurate usage information ([#125](https://github.com/openai/openai-go/issues/125)) ([c38b6c3](https://github.com/openai/openai-go/commit/c38b6c36c911d47b369518b7222363befa9c1e25))
* **ssestream:** skip SSE blocks without data ([#621](https://github.com/openai/openai-go/issues/621)) ([114224d](https://github.com/openai/openai-go/commit/114224dd71cf6695a01de8353352145832662e84))


### Documentation

* fix Azure API version link typo ([#630](https://github.com/openai/openai-go/issues/630)) ([4d7f649](https://github.com/openai/openai-go/commit/4d7f6493d248af66e42bfbb2bfb6c3bd4ca04b9c))
* fix function call comment wording in example ([#639](https://github.com/openai/openai-go/issues/639)) ([050ab8a](https://github.com/openai/openai-go/commit/050ab8af70b562ff4c81a16476bf393534849d92))
* links in README.md ([#308](https://github.com/openai/openai-go/issues/308)) ([35501ce](https://github.com/openai/openai-go/commit/35501ce5ec0449128bd84d57b24841ef541c3122))

## 3.42.0 (2026-07-09)

Full Changelog: [v3.41.1...v3.42.0](https://github.com/openai/openai-go/compare/v3.41.1...v3.42.0)

### Features

* **api:** gpt-5.6-sol updates ([8c062a6](https://github.com/openai/openai-go/commit/8c062a6e8fa8ae91daf96f97e301447d04ecd371))


### Bug Fixes

* **api:** restore programmatic response types ([2384620](https://github.com/openai/openai-go/commit/2384620bab9bede6a6b57cc2689bf86084d1e8f3))
* **api:** restore response union metadata ([7b2e975](https://github.com/openai/openai-go/commit/7b2e9755e6532a33d14085411a9f57d7709ada59))

## 3.41.1 (2026-07-06)

Full Changelog: [v3.41.0...v3.41.1](https://github.com/openai/openai-go/compare/v3.41.0...v3.41.1)

### Documentation

* fix Azure API versioning link ([#649](https://github.com/openai/openai-go/issues/649)) ([a8bd01a](https://github.com/openai/openai-go/commit/a8bd01aa80ba28bb3b5944cdac24a1aadc18fdbe))

## 3.41.0 (2026-06-17)

Full Changelog: [v3.40.0...v3.41.0](https://github.com/openai/openai-go/compare/v3.40.0...v3.41.0)

### Features

* **api:** update OpenAPI spec or Stainless config ([868891f](https://github.com/openai/openai-go/commit/868891f7fd02f5f92fa86a8b7a0fcd7a8c950bb1))

## 3.40.0 (2026-06-16)

Full Changelog: [v3.39.0...v3.40.0](https://github.com/openai/openai-go/compare/v3.39.0...v3.40.0)

### Features

* **api:** admin spend_alerts ([26992e5](https://github.com/openai/openai-go/commit/26992e5e1d4553d51cc63dddb0c0817a6940448c))
* **api:** manual updates ([28f7e09](https://github.com/openai/openai-go/commit/28f7e092c03b7458c69c0dd467fbcfcfb7a31864))
* **api:** update OpenAPI spec or Stainless config ([57c307e](https://github.com/openai/openai-go/commit/57c307e3893744bfc186d23cbb5c0d1f2e9d798a))

## 3.39.0 (2026-06-03)

Full Changelog: [v3.38.0...v3.39.0](https://github.com/openai/openai-go/compare/v3.38.0...v3.39.0)

### Features

* **api:** responses.moderation and chat_completions.moderation ([7a2dac0](https://github.com/openai/openai-go/commit/7a2dac0ddf6b92dfcaa45b190dbe7f51368e199c))

## 3.38.0 (2026-06-01)

Full Changelog: [v3.37.0...v3.38.0](https://github.com/openai/openai-go/compare/v3.37.0...v3.38.0)

### Features

* **api:** manual updates ([d7dac81](https://github.com/openai/openai-go/commit/d7dac8192c292d882cb765d3379e10427930979e))
* **api:** workload identity in audit logs, additional_tools item in responses, fix ActionSearch.query to be optional. ([4c3981c](https://github.com/openai/openai-go/commit/4c3981cfabc98ec1d67c2baf7c6a891ef6640f96))

## 3.37.0 (2026-05-21)

Full Changelog: [v3.36.0...v3.37.0](https://github.com/openai/openai-go/compare/v3.36.0...v3.37.0)

### Features

* **api:** api update ([7f7416e](https://github.com/openai/openai-go/commit/7f7416ea4f6953a2861189dee6391515c3b995a9))
* **api:** manual updates ([d646562](https://github.com/openai/openai-go/commit/d6465620413df87d971e7e37ae74bef4c70076b1))
* **api:** update OpenAPI spec or Stainless config ([b34b78a](https://github.com/openai/openai-go/commit/b34b78a83433003a6168fffd175cc963ad719495))
* **client:** optimize json encoder for internal types ([93adc6e](https://github.com/openai/openai-go/commit/93adc6e6247e8ce830152c3df0980a3154aa098a))


### Bug Fixes

* **go:** format generated admin paths ([1dd8f5e](https://github.com/openai/openai-go/commit/1dd8f5ec0adeeefef6a56068b5532ba5e3b3290e))
* **go:** format generated project permission paths ([b751c37](https://github.com/openai/openai-go/commit/b751c37ce2d6348545d75451dfc253dd7dda0f4f))


### Chores

* **api:** docs updates ([08bc80e](https://github.com/openai/openai-go/commit/08bc80ea58a19ba0725942c1f3afbcfb043851a0))

## 3.36.0 (2026-05-13)

Full Changelog: [v3.35.0...v3.36.0](https://github.com/openai/openai-go/compare/v3.35.0...v3.36.0)

### Features

* **api:** add service_tier parameter to response compact method ([bacd2c0](https://github.com/openai/openai-go/commit/bacd2c0bcf980e8d424d67446fb4d9c4ea897d24))


### Bug Fixes

* **go:** avoid panic when http.DefaultTransport is wrapped ([95a0250](https://github.com/openai/openai-go/commit/95a0250a9c770674f8deacb3a3fc1175e6808967))

## 3.35.0 (2026-05-07)

Full Changelog: [v3.34.0...v3.35.0](https://github.com/openai/openai-go/compare/v3.34.0...v3.35.0)

### Features

* **api:** add quantity field to admin organization usage costs results ([899b6de](https://github.com/openai/openai-go/commit/899b6dec70bd180a98497150b892285261882625))
* **api:** launch realtime translate + update image 2 ([1fbf3a6](https://github.com/openai/openai-go/commit/1fbf3a68e418f84993db1c291255c467e47dbe37))
* **api:** manual updates ([cd9047f](https://github.com/openai/openai-go/commit/cd9047f7416cedba79c7479d96acce3f59ac2c86))
* **api:** realtime 2 ([5b51410](https://github.com/openai/openai-go/commit/5b514101d40d0b9f36194bab0c7e49a7f6098c03))


### Bug Fixes

* **api:** fix imagegen `size` enum regression ([995778e](https://github.com/openai/openai-go/commit/995778eac4b36433a7d116f9a09999f34e78b91f))
* **types:** add format tags for URI and timestamp fields across models ([a6c2198](https://github.com/openai/openai-go/commit/a6c219806da98fa140b837a1063d3ee7f5f51e8d))


### Chores

* **internal:** regenerate SDK with no functional changes ([71c80ef](https://github.com/openai/openai-go/commit/71c80ef9e4407d6714f9d0013bc9579dfab1a223))
* redact api-key headers in debug logs ([ff51a02](https://github.com/openai/openai-go/commit/ff51a02d264f06535bac6000865b8ad27e4e8a72))


### Documentation

* **api:** document web_search_call.results output data option ([6c2d8ef](https://github.com/openai/openai-go/commit/6c2d8eff506aae6b4fea17fdca94748bd4b2cd3e))

## 3.34.0 (2026-05-01)

Full Changelog: [v3.33.0...v3.34.0](https://github.com/openai/openai-go/compare/v3.33.0...v3.34.0)

### Features

* **api:** add group_type, external_key_id, user fields to admin organization models ([460f400](https://github.com/openai/openai-go/commit/460f400cc195d603631a13d6ce867d7bf9ca7fb3))
* **api:** add support for Admin API Keys per endpoint ([51a0561](https://github.com/openai/openai-go/commit/51a0561845cb7f567e28074cddbe8bf1936f553a))
* **api:** admin API updates ([a0ad21e](https://github.com/openai/openai-go/commit/a0ad21e5584562bc7c9c951a6872479bd1aa7773))
* **api:** manual updates ([b60b3bc](https://github.com/openai/openai-go/commit/b60b3bcbc0bb9f11cd74d1b1cfba116ab1893465))
* **api:** manual updates ([f4e0c1a](https://github.com/openai/openai-go/commit/f4e0c1a9f98d838d494b71b4d32e6f7fc5ffab2d))


### Bug Fixes

* **auth:** prefer admin keys for admin requests ([d128349](https://github.com/openai/openai-go/commit/d1283491698c89bf53acc945db02c5094741021d))
* **types:** add format annotations to ChatKit timestamp and URL fields ([0b7d7d0](https://github.com/openai/openai-go/commit/0b7d7d036b4588877bb3f2fc1d46a56f23150468))
* **types:** add format annotations to timestamp and URL fields across resources ([0ff11a7](https://github.com/openai/openai-go/commit/0ff11a7f0c9eff426e1e408b0b70bc101e01fca8))
* **types:** change timestamp fields to float64 in Response ([d714e3c](https://github.com/openai/openai-go/commit/d714e3c84592cd2630b45e8f634ede88b19a4bb5))

## 3.33.0 (2026-04-28)

Full Changelog: [v3.32.0...v3.33.0](https://github.com/openai/openai-go/compare/v3.32.0...v3.33.0)

### Features

* **go:** add default http client with timeout ([25dd39a](https://github.com/openai/openai-go/commit/25dd39aec0ca97f51be1fd32accd5b10ff441365))
* support setting headers via env ([b4eaa78](https://github.com/openai/openai-go/commit/b4eaa78969c6c908ab81ec4e6c15fa9803e7de17))


### Bug Fixes

* **types:** change PromptCacheRetention enum in-memory to in_memory ([15e7432](https://github.com/openai/openai-go/commit/15e7432bd0a3fd4891be46f5d124ebc44301b611))


### Chores

* **internal:** more robust bootstrap script ([eb0cbcf](https://github.com/openai/openai-go/commit/eb0cbcf152867ad154c827dcf2c3c91adb773af9))
* **tests:** bump steady to v0.22.1 ([ef28f59](https://github.com/openai/openai-go/commit/ef28f59c84a30ffd0d2342bc8d09bea5803c89a6))


### Documentation

* **api:** add rate limit and vector store guidance to file upload ([6c8c7fa](https://github.com/openai/openai-go/commit/6c8c7faa70bd73a673744107c3487ec2a7362dca))
* **api:** update rate limit documentation in file upload method ([54fde2c](https://github.com/openai/openai-go/commit/54fde2c59e0f3b05741a93ec38cbe96495977848))

## 3.32.0 (2026-04-16)

Full Changelog: [v3.31.0...v3.32.0](https://github.com/openai/openai-go/compare/v3.31.0...v3.32.0)

### Features

* **api:** Add detail to InputFileContent ([b8e782f](https://github.com/openai/openai-go/commit/b8e782f8cf858ca4967553d221b7846471bc9ac7))
* **api:** add OAuthErrorCode type ([d54fff1](https://github.com/openai/openai-go/commit/d54fff1c958456d0441133015c5a34be30bfcec0))
* **api:** add prompt_cache_retention parameter to response compact ([e2da10d](https://github.com/openai/openai-go/commit/e2da10ddc8a28d1ea2bd47c907ba9149af700b9d))


### Bug Fixes

* fix for union type names ([0bc8ad5](https://github.com/openai/openai-go/commit/0bc8ad5e3372b432deee5286f392929ffc262c90))


### Documentation

* improve examples ([c3c36ad](https://github.com/openai/openai-go/commit/c3c36addfcb65a7b76f9d7d8ddabcf4c0d60511e))

## 3.31.0 (2026-04-08)

Full Changelog: [v3.30.0...v3.31.0](https://github.com/openai/openai-go/compare/v3.30.0...v3.31.0)

### Features

* **api:** add phase field to conversation messages ([563d747](https://github.com/openai/openai-go/commit/563d747fff0370d16e5cffff5a11e8e02e17bb8b))
* **api:** add web_search_call.results to ResponseIncludable ([6ef5b20](https://github.com/openai/openai-go/commit/6ef5b20a234ebdc85bd580539989a3933f06b5df))
* **client:** add support for short-lived tokens ([#799](https://github.com/openai/openai-go/issues/799)) ([038871e](https://github.com/openai/openai-go/commit/038871ef135200279e45bd8920b0d142c58e51c0))
* **internal:** support comma format in multipart form encoding ([905b473](https://github.com/openai/openai-go/commit/905b473d7bfb14ffd1db7c7b48efc1fdb84d1750))


### Bug Fixes

* fix issue with unmarshaling in some cases ([a017bc6](https://github.com/openai/openai-go/commit/a017bc6b810fcb01558736cd742f180dce342e49))
* prevent duplicate ? in query params ([cbb2c75](https://github.com/openai/openai-go/commit/cbb2c754223558d6fd19f34642a0aa89adba42e7))
* **types:** remove web_search_call.results from ResponseIncludable ([180597f](https://github.com/openai/openai-go/commit/180597fddba09d631a2fdc786116b3ad4f2e5418))


### Chores

* **ci:** support opting out of skipping builds on metadata-only commits ([12e3a5a](https://github.com/openai/openai-go/commit/12e3a5afdbacd2fad57640ba5486226898d6df71))
* remove unnecessary error check for url parsing ([cfe9c41](https://github.com/openai/openai-go/commit/cfe9c411d11b6f097d6e0563767bf231abb9975a))
* **tests:** bump steady to v0.20.1 ([8e4ef11](https://github.com/openai/openai-go/commit/8e4ef117bdec78698d72a00b0a14d12cf0482561))
* **tests:** bump steady to v0.20.2 ([ad31f81](https://github.com/openai/openai-go/commit/ad31f81e22e4f2cc931e11625755f92978a18be9))
* update docs for api:"required" ([2d16ebc](https://github.com/openai/openai-go/commit/2d16ebc57c3eca70c5bcb421ee8592549722e7f4))


### Documentation

* **api:** add multi-file ingestion guidance to vectorstorefile/filebatch params ([dbba33f](https://github.com/openai/openai-go/commit/dbba33f751ffa7e7a36678c4305a3e7630a4494f))

## 3.30.0 (2026-03-25)

Full Changelog: [v3.29.0...v3.30.0](https://github.com/openai/openai-go/compare/v3.29.0...v3.30.0)

### Features

* **api:** add keys field to computer action types ([7abb3c4](https://github.com/openai/openai-go/commit/7abb3c4b4d1fb3d2757dd70c5a201de4190f298a))


### Bug Fixes

* **api:** align SDK response types with expanded item schemas ([324c584](https://github.com/openai/openai-go/commit/324c584ec31315ec18c169918e524741e74d395a))
* **types:** generate shared enum types that are not referenced by other schemas ([e2f2920](https://github.com/openai/openai-go/commit/e2f2920f510348ea8b7e26e6b2460f802d34998c))
* **types:** make Type required, correct enum in ResponseInputMessageItem ([8c2c7a5](https://github.com/openai/openai-go/commit/8c2c7a59329e95da117ffdde7b36f43e1fbdbacb))


### Chores

* **ci:** skip lint on metadata-only changes ([40cf884](https://github.com/openai/openai-go/commit/40cf884295e1631dc637f6e3e1316e8665e5e694))
* **client:** fix multipart serialisation of Default() fields ([06c04ae](https://github.com/openai/openai-go/commit/06c04aeca47e7348de94573eefbfaf2e110c4573))
* **internal:** support default value struct tag ([f1afb2a](https://github.com/openai/openai-go/commit/f1afb2ab20307c50da7ec421ac5cd3e5a9c33fe3))
* **internal:** update gitignore ([bde0a0d](https://github.com/openai/openai-go/commit/bde0a0d3e4a8b63e2ce28fae207f9546abac0722))
* **tests:** bump steady to v0.19.4 ([a0561c3](https://github.com/openai/openai-go/commit/a0561c3014b5d5fbeee9d7de58dad3f9be9e2d9f))
* **tests:** bump steady to v0.19.5 ([56c5c36](https://github.com/openai/openai-go/commit/56c5c3631cb1dec68b8cb68c1e03d62e69fe3a70))
* **tests:** bump steady to v0.19.6 ([44ebe37](https://github.com/openai/openai-go/commit/44ebe3744e70460f2cacd94ccf21957e30d697a1))
* **tests:** bump steady to v0.19.7 ([1b04072](https://github.com/openai/openai-go/commit/1b04072ef94542b2a0007e4088ef09d4f5387dcf))


### Refactors

* **tests:** switch from prism to steady ([2185358](https://github.com/openai/openai-go/commit/21853589203e202637f9ebc7f7be5d4c73f2a471))

## 3.29.0 (2026-03-17)

Full Changelog: [v3.28.0...v3.29.0](https://github.com/openai/openai-go/compare/v3.28.0...v3.29.0)

### Features

* **api:** 5.4 nano and mini model slugs ([c4ebb48](https://github.com/openai/openai-go/commit/c4ebb48b3314d41ae1ca8006ccf7d5b9facad61e))
* **api:** add in and nin values to comparison filter type ([d060d44](https://github.com/openai/openai-go/commit/d060d4473d6ca87e0427eea709f87d5edd0e7c78))


### Chores

* **internal:** tweak CI branches ([6379626](https://github.com/openai/openai-go/commit/63796269b563770b6cbb24144edade0e57051e50))

## 3.28.0 (2026-03-14)

Full Changelog: [v3.27.0...v3.28.0](https://github.com/openai/openai-go/compare/v3.27.0...v3.28.0)

### Features

* **api:** add /v1/videos endpoint option to batch ([7b2d67e](https://github.com/openai/openai-go/commit/7b2d67e3d65737572d89536d16ed81a3ce39688f))
* **api:** add defer_loading field to function tools ([6d4b683](https://github.com/openai/openai-go/commit/6d4b6833e5b0b29a9b1d0c99062a231290e8b93f))
* **api:** custom voices ([d00b782](https://github.com/openai/openai-go/commit/d00b782c32db4c953b8e39edc5a77504693c70f3))

### ⚠ BREAKING CHANGES

* **api:** The `voice` parameter type has changed from a `string` to a `string | {id: string}`. This is a breaking change for Go.

## 3.27.0 (2026-03-13)

Full Changelog: [v3.26.0...v3.27.0](https://github.com/openai/openai-go/compare/v3.26.0...v3.27.0)

### Features

* **api:** add video character/edit/extend, remove Azure/webhook/polling/accumulator ([20da043](https://github.com/openai/openai-go/commit/20da043643286aa444450f4bf4ab6f68c5401455))
* **api:** add video edits/extensions/character, remove Azure/webhook/polling helpers ([fa9413f](https://github.com/openai/openai-go/commit/fa9413f9b7c00ec36b124c4b3714daf2b2cd978d))
* **api:** api update ([c88c6c9](https://github.com/openai/openai-go/commit/c88c6c9732f7b7a8cd739b8006a56d564ab59183))
* **api:** sora api improvements: character api, video extensions/edits, higher resolution exports. ([93f6779](https://github.com/openai/openai-go/commit/93f6779c0c2a1beb7f68b570c61aae8fda52fd1f))


### Chores

* **internal:** codegen related update ([d531232](https://github.com/openai/openai-go/commit/d5312325255bcac5934636d320653c338bb9622e))
* **internal:** codegen related update ([1748c11](https://github.com/openai/openai-go/commit/1748c114b2d642b38fdacfef7823cbd4130c29c5))
* **internal:** codegen related update ([531ece1](https://github.com/openai/openai-go/commit/531ece1703028189ba6081a168f3a979481dcf0c))
* **internal:** codegen related update ([ce9f435](https://github.com/openai/openai-go/commit/ce9f4357c22b533fff0d3385058dbadc033928e8))

## 3.26.0 (2026-03-05)

Full Changelog: [v3.25.0...v3.26.0](https://github.com/openai/openai-go/compare/v3.25.0...v3.26.0)

### Features

* **api:** The GA computer tool now uses the ComputerTool type. The `computer_use_preview` tool has moved to ComputerUsePreviewTool ([347418b](https://github.com/openai/openai-go/commit/347418be8d4fa33881d9ac30f6c7132f2f545f2b))

## 3.25.0 (2026-03-05)

Full Changelog: [v3.24.0...v3.25.0](https://github.com/openai/openai-go/compare/v3.24.0...v3.25.0)

### Features

* **api:** gpt-5.4, tool search tool, and new computer tool ([101826d](https://github.com/openai/openai-go/commit/101826dd757a0213aecb4eaa6332866657b9aa83))
* **api:** remove Phase from input/output messages, PromptCacheKey from responses ([961b8ca](https://github.com/openai/openai-go/commit/961b8ca27923beca8aa08d4a8e3382c2da9d61db))


### Bug Fixes

* **api:** internal schema fixes ([fe5f7cd](https://github.com/openai/openai-go/commit/fe5f7cdb34d11dd18caa503716cae1512b245053))
* **api:** manual updates ([70b02c8](https://github.com/openai/openai-go/commit/70b02c8f63c98a17813dc6cb7f7707fb2bba81c5))
* **api:** readd phase ([548aff8](https://github.com/openai/openai-go/commit/548aff8ad8b96518f5549ec3bc98da71e9b7f540))


### Chores

* **internal:** codegen related update ([ab733b9](https://github.com/openai/openai-go/commit/ab733b91db39e99e292696530340333c065e04b9))
* **internal:** codegen related update ([23d1831](https://github.com/openai/openai-go/commit/23d1831cb5ca6f61ca8575737cec17e2f347818b))
* **internal:** reduce warnings ([2963312](https://github.com/openai/openai-go/commit/2963312c075fa9a30abad32b1e90813229b22129))

## 3.24.0 (2026-02-24)

Full Changelog: [v3.23.0...v3.24.0](https://github.com/openai/openai-go/compare/v3.23.0...v3.24.0)

### Features

* **api:** add phase ([72366d8](https://github.com/openai/openai-go/commit/72366d895c78b5188a590ee7f9b572b567447b32))


### Bug Fixes

* **api:** fix phase enum ([5712ebf](https://github.com/openai/openai-go/commit/5712ebf1f30e63d148a72c451f1df48620b14a2d))
* **api:** phase docs ([b67dd66](https://github.com/openai/openai-go/commit/b67dd6680110c013c1660c37dde5467e2cd50030))


### Chores

* **internal:** move custom custom `json` tags to `api` ([0735303](https://github.com/openai/openai-go/commit/0735303849ece03c57adbb0f899d7f3a0b60bc50))
* **internal:** refactor sse event parsing ([45dc6bb](https://github.com/openai/openai-go/commit/45dc6bb18120330de4470739a07b13f08d7f7666))

## 3.23.0 (2026-02-24)

Full Changelog: [v3.22.1...v3.23.0](https://github.com/openai/openai-go/compare/v3.22.1...v3.23.0)

### Features

* **api:** add gpt-realtime-1.5 and gpt-audio-1.5 models to realtime session ([9076e2f](https://github.com/openai/openai-go/commit/9076e2f2fab882d5a8ba9992096e5997902c5589))

## 3.22.1 (2026-02-23)

Full Changelog: [v3.22.0...v3.22.1](https://github.com/openai/openai-go/compare/v3.22.0...v3.22.1)

### Bug Fixes

* allow canceling a request while it is waiting to retry ([54672cf](https://github.com/openai/openai-go/commit/54672cf6b2c21a0e2ac0d2a7d7bed8680eee3e44))


### Chores

* update mock server docs ([3ac84dd](https://github.com/openai/openai-go/commit/3ac84dd90e21d9374c0141d86c07d21e0914c6b0))


### Documentation

* **api:** add batch size limit to vector store file batch parameters ([f751c40](https://github.com/openai/openai-go/commit/f751c40c522d6fba5c4eb244fd25f28c2317ca33))
* **api:** clarify safety_identifier max length in chat completions and responses ([8257f9b](https://github.com/openai/openai-go/commit/8257f9b0e4d63038f2b088a42399e3a80e9c9bb6))
* **api:** enhance method docstrings across audio/chat/completion/skill/upload/video APIs ([38b8f63](https://github.com/openai/openai-go/commit/38b8f63a16f9a50bb176561f7842baf976e88316))

## 3.22.0 (2026-02-13)

Full Changelog: [v3.21.0...v3.22.0](https://github.com/openai/openai-go/compare/v3.21.0...v3.22.0)

### Features

* **api:** container network_policy and skills ([8e5ea23](https://github.com/openai/openai-go/commit/8e5ea2344580eefa186040cae3583dacf459e0b9))


### Bug Fixes

* **encoder:** correctly serialize NullStruct ([a6cb49e](https://github.com/openai/openai-go/commit/a6cb49ef3743b1c0b58d24f8dbc5e16e1a5f5852))


### Documentation

* update comment ([bc316d7](https://github.com/openai/openai-go/commit/bc316d7b10fa928289e0560123b0de16099edfb7))

## 3.21.0 (2026-02-10)

Full Changelog: [v3.20.0...v3.21.0](https://github.com/openai/openai-go/compare/v3.20.0...v3.21.0)

### Features

* **api:** support for images in batch api ([e23aeb1](https://github.com/openai/openai-go/commit/e23aeb1b13bfd089cc73d3097c9635b687446f82))

## 3.20.0 (2026-02-10)

Full Changelog: [v3.19.0...v3.20.0](https://github.com/openai/openai-go/compare/v3.19.0...v3.20.0)

### Features

* **api:** skills and hosted shell ([9e191de](https://github.com/openai/openai-go/commit/9e191de75f67a6a693c8b25ac9ab1b9288673993))

## 3.19.0 (2026-02-09)

Full Changelog: [v3.18.0...v3.19.0](https://github.com/openai/openai-go/compare/v3.18.0...v3.19.0)

### Features

* **api:** responses context_management ([199f230](https://github.com/openai/openai-go/commit/199f23025ab098f2ac0ac9a99dee37235613c287))

## 3.18.0 (2026-02-05)

Full Changelog: [v3.17.0...v3.18.0](https://github.com/openai/openai-go/compare/v3.17.0...v3.18.0)

### Features

* **api:** add shell_call_output status field ([67a75d7](https://github.com/openai/openai-go/commit/67a75d755e815f6d6fdf4ac48314472a94c8613f))
* **api:** image generation actions for responses; ResponseFunctionCallArgumentsDoneEvent.name ([2c57016](https://github.com/openai/openai-go/commit/2c57016b7c7f45072c59f193e567a55ecbda21fd))


### Bug Fixes

* **client:** undo change to web search Find action ([e340256](https://github.com/openai/openai-go/commit/e340256509214ee386de32b993f5ec4ebba43d38))
* **client:** update type for `find_in_page` action ([4b5d499](https://github.com/openai/openai-go/commit/4b5d4993e82ada68276bb5560bb2cd8b457aa3da))


### Chores

* **client:** improve example values ([c86a65c](https://github.com/openai/openai-go/commit/c86a65cefd55eb18568f4b7d2660c82dc90af4ad))


### Documentation

* split `api.md` by standalone resources ([aeed37b](https://github.com/openai/openai-go/commit/aeed37b814d37ad3d59111b7665d48bf220cbf9e))

## 3.17.0 (2026-01-27)

Full Changelog: [v3.16.0...v3.17.0](https://github.com/openai/openai-go/compare/v3.16.0...v3.17.0)

### Features

* **api:** api update ([a456c60](https://github.com/openai/openai-go/commit/a456c60498b33b7da048cc64bdff76c49a904117))
* **api:** api updates ([21fd4a9](https://github.com/openai/openai-go/commit/21fd4a9534d5ef531c6a9bc497b90c14b68ebda3))
* **client:** add a convenient param.SetJSON helper ([1b35ece](https://github.com/openai/openai-go/commit/1b35ece947934982609557e6acacfd3526289de1))
* **client:** add a StreamError type to access raw events from sse streams ([fd14b30](https://github.com/openai/openai-go/commit/fd14b30e3cc2d14438a100be12627692e8ea045a))


### Bug Fixes

* **api:** mark assistants as deprecated ([9a8c9af](https://github.com/openai/openai-go/commit/9a8c9af8cf988069e543487c09a0897806408d67))
* **client:** retain streaming when user sets request body ([3a7a22e](https://github.com/openai/openai-go/commit/3a7a22ec90c5ff44203321bbff94f9541a80601f))
* **docs:** add missing pointer prefix to api.md return types ([dd641d9](https://github.com/openai/openai-go/commit/dd641d924ead979b4369b952f99387880879938d))


### Chores

* **internal:** codegen related update ([74d3989](https://github.com/openai/openai-go/commit/74d3989139a94407020f7bf43c8351c4dffe412c))
* **internal:** update `actions/checkout` version ([2db54a5](https://github.com/openai/openai-go/commit/2db54a5d05e3240c415cb91320ede5704331828a))

## 3.16.0 (2026-01-09)

Full Changelog: [v3.15.0...v3.16.0](https://github.com/openai/openai-go/compare/v3.15.0...v3.16.0)

### Features

* **api:** add new Response completed_at prop ([bff6331](https://github.com/openai/openai-go/commit/bff6331c1b428de935966f59f2465e77df08f075))


### Bug Fixes

* **client:** use the correct order of params for vector store file and batch polling ([ef32641](https://github.com/openai/openai-go/commit/ef32641b81da84c7d0524c372ee8b45cda71fe2c))


### Chores

* **internal:** codegen related update ([1e529a4](https://github.com/openai/openai-go/commit/1e529a4b48c55a6bb611f5dcaf0ad1bfbd6f729d))
* **internal:** use different example values for some enums ([a2836ee](https://github.com/openai/openai-go/commit/a2836eef250be42d9e0f135c36823219129cad1e))


### Documentation

* update URL version ([dc00e14](https://github.com/openai/openai-go/commit/dc00e14cae07daaefabeb5371daa12f90bb82dc8))

## 3.15.0 (2025-12-19)

Full Changelog: [v3.14.0...v3.15.0](https://github.com/openai/openai-go/compare/v3.14.0...v3.15.0)

### Bug Fixes

* rebuild ([8205ae7](https://github.com/openai/openai-go/commit/8205ae7c00de1bd4a543381ba61b34f9b5676eda))


### Chores

* add float64 to valid types for RegisterFieldValidator ([e67d89d](https://github.com/openai/openai-go/commit/e67d89d39bc14af7309df702592ae152d1dfd60b))

## 3.14.0 (2025-12-16)

Full Changelog: [v3.13.0...v3.14.0](https://github.com/openai/openai-go/compare/v3.13.0...v3.14.0)

### Features

* **api:** gpt-image-1.5 ([6102f02](https://github.com/openai/openai-go/commit/6102f029e7ccbffe1dcf4d53b38f7da49dfbdbaa))

## 3.13.0 (2025-12-15)

Full Changelog: [v3.12.0...v3.13.0](https://github.com/openai/openai-go/compare/v3.12.0...v3.13.0)

### Features

* **api:** api update ([20b5112](https://github.com/openai/openai-go/commit/20b51126dc55b5fa357ae848593873d46514d820))
* **api:** fix grader input list, add dated slugs for sora-2 ([e8f0b76](https://github.com/openai/openai-go/commit/e8f0b76c55abdcca2920372f74e08621d8a530b9))


### Bug Fixes

* **azure:** correct Azure OpenAI API URL construction and auth ([3ba3736](https://github.com/openai/openai-go/commit/3ba3736c4b1a6138c05df5ccb64944a3dca6ea74))

## 3.12.0 (2025-12-11)

Full Changelog: [v3.11.0...v3.12.0](https://github.com/openai/openai-go/compare/v3.11.0...v3.12.0)

### Features

* **api:** gpt 5.2 ([56b5d41](https://github.com/openai/openai-go/commit/56b5d410cb7ee90d2c7ddd4fb8650bf1dd818855))
* **encoder:** support bracket encoding form-data object members ([a2cbacf](https://github.com/openai/openai-go/commit/a2cbacff1f0189e81ec4091a33f1ec350bbabd09))

## 3.11.0 (2025-12-10)

Full Changelog: [v3.10.0...v3.11.0](https://github.com/openai/openai-go/compare/v3.10.0...v3.11.0)

### Features

* **api:** make model required for the responses/compact endpoint ([05f8f4d](https://github.com/openai/openai-go/commit/05f8f4de214bfdb5ad8946e7c5252b14a41e3122))


### Bug Fixes

* **mcp:** correct code tool API endpoint ([03d08f9](https://github.com/openai/openai-go/commit/03d08f934e5ee8294af09eb7278c1c7c4233e1f7))
* rename param to avoid collision ([17d276d](https://github.com/openai/openai-go/commit/17d276d797f505a0254b112c3f7926490d1d789e))


### Chores

* elide duplicate aliases ([2bf988e](https://github.com/openai/openai-go/commit/2bf988eee937ced3ec3f53bd7773b6eab07bbdbf))
* **internal:** codegen related update ([9b1a3e9](https://github.com/openai/openai-go/commit/9b1a3e99d31006ad9db6dd8e9dc2902a9b57cf02))

## 3.10.0 (2025-12-04)

Full Changelog: [v3.9.0...v3.10.0](https://github.com/openai/openai-go/compare/v3.9.0...v3.10.0)

### Features

* **api:** gpt-5.1-codex-max and responses/compact ([1e1ca2a](https://github.com/openai/openai-go/commit/1e1ca2a6369c79a79bb54df3ee40b2d5604a21c5))

## 3.9.0 (2025-12-01)

Full Changelog: [v3.8.1...v3.9.0](https://github.com/openai/openai-go/compare/v3.8.1...v3.9.0)

### Features

* **api:** gpt 5.1 ([470f91f](https://github.com/openai/openai-go/commit/470f91faac304e518019be9f7b12e6270af63bbd))


### Bug Fixes

* **api:** align types of input items / output items for typescript ([5b89d3b](https://github.com/openai/openai-go/commit/5b89d3ba03968ee9f5b49e7e065495c3c5c77710))
* **client:** correctly specify Accept header with */* instead of empty ([fbadb4e](https://github.com/openai/openai-go/commit/fbadb4e8b1a81c99a7b3936da483ee9542de2c23))


### Chores

* bump gjson version ([305831f](https://github.com/openai/openai-go/commit/305831feb6c39d1f9f6e85c2e9f94f6c7f0dcd45))
* fix empty interfaces ([2aaa980](https://github.com/openai/openai-go/commit/2aaa980c2f0cac814065e4e5e294b151500c2e3f))

## 3.8.1 (2025-11-04)

Full Changelog: [v3.8.0...v3.8.1](https://github.com/openai/openai-go/compare/v3.8.0...v3.8.1)

### Bug Fixes

* **api:** fix nullability of logprobs ([b5aeb99](https://github.com/openai/openai-go/commit/b5aeb999e5088db4f9d1232a202a568a4a283019))

## 3.8.0 (2025-11-03)

Full Changelog: [v3.7.0...v3.8.0](https://github.com/openai/openai-go/compare/v3.7.0...v3.8.0)

### Features

* **api:** Realtime API token_limits, Hybrid searching ranking options ([9495f4a](https://github.com/openai/openai-go/commit/9495f4aa72bd5784fe3291637349ae1b706d8f8c))


### Chores

* **internal:** grammar fix (it's -&gt; its) ([879772d](https://github.com/openai/openai-go/commit/879772dc881a3d4fecf1425afc94a4bc141e9fb8))

## 3.7.0 (2025-10-28)

Full Changelog: [v3.6.1...v3.7.0](https://github.com/openai/openai-go/compare/v3.6.1...v3.7.0)

### Features

* **api:** remove InputAudio from ResponseInputContent ([cf50c53](https://github.com/openai/openai-go/commit/cf50c53f779784e1ee73b7d815456afaa3e1c447))
* **azure:** allow passing custom scopes ([#541](https://github.com/openai/openai-go/issues/541)) ([dffa08e](https://github.com/openai/openai-go/commit/dffa08ece6c860ae1f87a01a5b8c26f18ce7ab2b))


### Bug Fixes

* **api:** docs updates ([94d54c1](https://github.com/openai/openai-go/commit/94d54c1e19d0d58875f56058042e06410b23ac49))

## 3.6.1 (2025-10-20)

Full Changelog: [v3.6.0...v3.6.1](https://github.com/openai/openai-go/compare/v3.6.0...v3.6.1)

### Bug Fixes

* **api:** fix discriminator propertyName for ResponseFormatJsonSchema ([57b0505](https://github.com/openai/openai-go/commit/57b0505361029563b5fd56fb6085b58e813936cc))

## 3.6.0 (2025-10-20)

Full Changelog: [v3.5.0...v3.6.0](https://github.com/openai/openai-go/compare/v3.5.0...v3.6.0)

### Features

* **api:** Add responses.input_tokens.count ([a43f2ce](https://github.com/openai/openai-go/commit/a43f2cef132d4cbd4a4a3dedf600f2da0a1ea2f5))


### Bug Fixes

* **api:** internal openapi updates ([7ad9b02](https://github.com/openai/openai-go/commit/7ad9b02d1e86cb3235c779e4e0f6e2ee226662d3))

## 3.5.0 (2025-10-17)

Full Changelog: [v3.4.0...v3.5.0](https://github.com/openai/openai-go/compare/v3.4.0...v3.5.0)

### Features

* **api:** api update ([1aa78dd](https://github.com/openai/openai-go/commit/1aa78dda7aae7b72ce021250b5357ead8db36f46))

## 3.4.0 (2025-10-16)

Full Changelog: [v3.3.0...v3.4.0](https://github.com/openai/openai-go/compare/v3.3.0...v3.4.0)

### Features

* **api:** Add support for gpt-4o-transcribe-diarize on audio/transcriptions endpoint ([ee32400](https://github.com/openai/openai-go/commit/ee32400f70d6d16c583978c574806648bdeecd91))


### Chores

* **api:** internal updates ([74c8031](https://github.com/openai/openai-go/commit/74c8031304013f5d7c24bd9db93d73da80efba9f))
* **client:** undo more naming changes ([db441bc](https://github.com/openai/openai-go/commit/db441bcb7fb830743d9489589a3a48ca79d2f80a))
* **client:** undo some naming changes ([a5aa3d6](https://github.com/openai/openai-go/commit/a5aa3d6e0d6773f838f826bbd68f96b70fef0653))

## 3.3.0 (2025-10-10)

Full Changelog: [v3.2.0...v3.3.0](https://github.com/openai/openai-go/compare/v3.2.0...v3.3.0)

### Features

* **api:** comparison filter in/not in ([d6daca0](https://github.com/openai/openai-go/commit/d6daca0eedd998f49d8bfde0c3caba74d762c0d6))

## 3.2.0 (2025-10-06)

Full Changelog: [v3.1.0...v3.2.0](https://github.com/openai/openai-go/compare/v3.1.0...v3.2.0)

### Features

* **api:** dev day 2025 launches ([d40a768](https://github.com/openai/openai-go/commit/d40a7689c769fd8b581fa753c5b748805c8d7bd1))

## 3.1.0 (2025-10-02)

Full Changelog: [v3.0.1...v3.1.0](https://github.com/openai/openai-go/compare/v3.0.1...v3.1.0)

### Features

* **api:** add support for realtime calls ([565ca67](https://github.com/openai/openai-go/commit/565ca678729182ae35c634ff7791383273b29993))

## 3.0.1 (2025-10-01)

Full Changelog: [v3.0.0...v3.0.1](https://github.com/openai/openai-go/compare/v3.0.0...v3.0.1)

### Bug Fixes

* **api:** add status, approval_request_id to MCP tool call ([a7f95e4](https://github.com/openai/openai-go/commit/a7f95e4ef4335a8eb3fc0e51e0b70b11b144e5aa))

## 3.0.0 (2025-09-30)

Full Changelog: [v2.7.1...v3.0.0](https://github.com/openai/openai-go/compare/v2.7.1...v3.0.0)

### ⚠ BREAKING CHANGES

* **api:** `ResponseFunctionToolCallOutputItem.output` and `ResponseCustomToolCallOutput.output` now return `string | Array<ResponseInputText | ResponseInputImage | ResponseInputFile>` instead of `string` only. This may break existing callsites that assume `output` is always a string.

### Features

* **api:** Support images and files for function call outputs in responses, BatchUsage ([21901ef](https://github.com/openai/openai-go/commit/21901ef84eac7028b92939c3e54c4ec7f2c8663f))

## 2.7.1 (2025-09-29)

Full Changelog: [v2.7.0...v2.7.1](https://github.com/openai/openai-go/compare/v2.7.0...v2.7.1)

### Bug Fixes

* bugfix for setting JSON keys with special characters ([f9ae028](https://github.com/openai/openai-go/commit/f9ae0283fe34fef6a8a7909655423b45795e41fc))

## 2.7.0 (2025-09-23)

Full Changelog: [v2.6.1...v2.7.0](https://github.com/openai/openai-go/compare/v2.6.1...v2.7.0)

### Features

* **api:** gpt-5-codex ([b0eac3e](https://github.com/openai/openai-go/commit/b0eac3ed2bcf2b7f0a5d6d68c9e13e7b7e409f0f))

## 2.6.1 (2025-09-22)

Full Changelog: [v2.6.0...v2.6.1](https://github.com/openai/openai-go/compare/v2.6.0...v2.6.1)

### Bug Fixes

* **api:** fix mcp tool name ([6de601a](https://github.com/openai/openai-go/commit/6de601aa71c7325938f839cdc0f3b7c808a5d7f8))
* use slices.Concat instead of sometimes modifying r.Options ([7312ee7](https://github.com/openai/openai-go/commit/7312ee73efec4bd523e18b9524072d6dcf8bab09))


### Chores

* **api:** openapi updates for conversations ([4a7d204](https://github.com/openai/openai-go/commit/4a7d204a4e1140babdeab43bdac59dfc8dae95b5))
* bump minimum go version to 1.22 ([8396ab5](https://github.com/openai/openai-go/commit/8396ab5d918bf068b6d6f342c825ba32d1d982b0))
* do not install brew dependencies in ./scripts/bootstrap by default ([d519b81](https://github.com/openai/openai-go/commit/d519b8100047bae1dbd1458112097c1c21880977))
* update more docs for 1.22 ([1b0514d](https://github.com/openai/openai-go/commit/1b0514df9508a652a11cb8efa70ac30eaa088dbe))

## 2.6.0 (2025-09-19)

Full Changelog: [v2.5.0...v2.6.0](https://github.com/openai/openai-go/compare/v2.5.0...v2.6.0)

### Features

* **api:** add reasoning_text ([6ebf50d](https://github.com/openai/openai-go/commit/6ebf50d756f06d951cdccff432615835bbf3165f))

## 2.5.0 (2025-09-17)

Full Changelog: [v2.4.3...v2.5.0](https://github.com/openai/openai-go/compare/v2.4.3...v2.5.0)

### Features

* **api:** type updates for conversations, reasoning_effort and results for evals ([3e68a60](https://github.com/openai/openai-go/commit/3e68a60d764645c5bfc9003f61525401268ef3a1))

## 2.4.3 (2025-09-15)

Full Changelog: [v2.4.2...v2.4.3](https://github.com/openai/openai-go/compare/v2.4.2...v2.4.3)

### Chores

* **api:** docs and spec refactoring ([e67af66](https://github.com/openai/openai-go/commit/e67af66b35df49267ccc7e3af73220d8f51339e9))

## 2.4.2 (2025-09-12)

Full Changelog: [v2.4.1...v2.4.2](https://github.com/openai/openai-go/compare/v2.4.1...v2.4.2)

### Chores

* **api:** Minor docs and type updates for realtime ([d92ea48](https://github.com/openai/openai-go/commit/d92ea4850f3720ba7a372f7bc9f8ecff07392ba0))

## 2.4.1 (2025-09-10)

Full Changelog: [v2.4.0...v2.4.1](https://github.com/openai/openai-go/compare/v2.4.0...v2.4.1)

### Chores

* **api:** fix realtime GA types ([012b83e](https://github.com/openai/openai-go/commit/012b83e3fa37a69d39eeaf6b227c37f5d3e42134))

## 2.4.0 (2025-09-08)

Full Changelog: [v2.3.1...v2.4.0](https://github.com/openai/openai-go/compare/v2.3.1...v2.4.0)

### Features

* **api:** ship the RealtimeGA API shape ([2b6c6db](https://github.com/openai/openai-go/commit/2b6c6db63e4871f3fa12a29c568365ac09290b9d))

## 2.3.1 (2025-09-05)

Full Changelog: [v2.3.0...v2.3.1](https://github.com/openai/openai-go/compare/v2.3.0...v2.3.1)

### Bug Fixes

* **internal:** unmarshal correctly when there are multiple discriminators ([98596b2](https://github.com/openai/openai-go/commit/98596b2183dcf3a13297b0dc07b0efc015dff83f))

## 2.3.0 (2025-09-03)

Full Changelog: [v2.2.2...v2.3.0](https://github.com/openai/openai-go/compare/v2.2.2...v2.3.0)

### Features

* **api:** Add gpt-realtime models ([3cf6a34](https://github.com/openai/openai-go/commit/3cf6a3484108786df49cd8e44356fc5fcaf58d8a))

## 2.2.2 (2025-09-02)

Full Changelog: [v2.2.1...v2.2.2](https://github.com/openai/openai-go/compare/v2.2.1...v2.2.2)

### Bug Fixes

* update url to refresh pkg.go.dev ([edf94ce](https://github.com/openai/openai-go/commit/edf94ce95a9f3fae87722a338c213dcf57ac1bf2))
* use release please annotations on more places ([2ff82f9](https://github.com/openai/openai-go/commit/2ff82f98ae636ff942cbdd8b909854f01279af90))

## 2.2.1 (2025-09-02)

Full Changelog: [v2.2.0...v2.2.1](https://github.com/openai/openai-go/compare/v2.2.0...v2.2.1)

### Chores

* **api:** manual updates for ResponseInputAudio ([8c0ebe5](https://github.com/openai/openai-go/commit/8c0ebe566fb03be01cd772a80eb2581b46b78f5c))

## 2.2.0 (2025-09-02)

Full Changelog: [v2.1.1...v2.2.0](https://github.com/openai/openai-go/compare/v2.1.1...v2.2.0)

### Features

* **api:** Add connectors support for MCP tool ([35888bc](https://github.com/openai/openai-go/commit/35888bcd26c7633e8ea68f9213cc3977b8ac49eb))
* **api:** add web search filters ([6f2c71d](https://github.com/openai/openai-go/commit/6f2c71d4e28971fc73e7e291d40f3b875d9cc42a))
* **api:** adding support for /v1/conversations to the API ([5b7c31b](https://github.com/openai/openai-go/commit/5b7c31bde9c1086d3fb71c88dfdf74228845b22e))
* **api:** realtime API updates ([130fc8e](https://github.com/openai/openai-go/commit/130fc8ea5ba39e6c1457ed6d26ef827d931a6242))
* **client:** add support for verifying signatures on incoming webhooks ([f7c8dbb](https://github.com/openai/openai-go/commit/f7c8dbb6b5bd5bab72b7d146dc255d543b0b5a71))


### Bug Fixes

* **azure:** compatibility with edit image endpoint ([#477](https://github.com/openai/openai-go/issues/477)) ([d156eec](https://github.com/openai/openai-go/commit/d156eeca37bc86a5d8e1c973063a8425744810f1))
* close body before retrying ([8dfed35](https://github.com/openai/openai-go/commit/8dfed35f11a00970ad804ab985cf393c2332ea8f))


### Chores

* **internal/ci:** setup breaking change detection ([0af0cd0](https://github.com/openai/openai-go/commit/0af0cd01302d3859a4e43554ed8e665007f69aad))
* **internal:** version bump ([3265795](https://github.com/openai/openai-go/commit/3265795fffa44fb40d65a800b300807d9f1e7b2b))

## 2.1.1 (2025-08-20)

Full Changelog: [v2.1.0...v2.1.1](https://github.com/openai/openai-go/compare/v2.1.0...v2.1.1)

### Chores

* **api:** accurately represent shape for verbosity on Chat Completions ([f81197b](https://github.com/openai/openai-go/commit/f81197b4b02f3aa022bc363d6db6949d0d105d92))

## 2.1.0 (2025-08-18)

Full Changelog: [v2.0.2...v2.1.0](https://github.com/openai/openai-go/compare/v2.0.2...v2.1.0)

### Features

* **api:** add new text parameters, expiration options ([323154c](https://github.com/openai/openai-go/commit/323154ccec2facf80d9ada76ed3c35553cb8896d))


### Documentation

* give https its missing "h" in Azure OpenAI REST API link ([#480](https://github.com/openai/openai-go/issues/480)) ([8a401c9](https://github.com/openai/openai-go/commit/8a401c9eecbe4936de487447be09757859001009))

## 2.0.2 (2025-08-09)

Full Changelog: [v2.0.1...v2.0.2](https://github.com/openai/openai-go/compare/v2.0.1...v2.0.2)

### Chores

* **internal:** update comment in script ([4be24de](https://github.com/openai/openai-go/commit/4be24dee6ab7b116ed34e50d56c99c1a36c0ef9d))
* update @stainless-api/prism-cli to v5.15.0 ([eca22af](https://github.com/openai/openai-go/commit/eca22af6f1d1f2ac36fbee365616210c12267bb1))

## 2.0.1 (2025-08-08)

Full Changelog: [v2.0.0...v2.0.1](https://github.com/openai/openai-go/compare/v2.0.0...v2.0.1)

### Bug Fixes

* **client:** fix verbosity parameter location in Responses ([6e2e903](https://github.com/openai/openai-go/commit/6e2e903e7c63a3e2a5aef5c81bdae55d220c0292))

## 2.0.0 (2025-08-07)

Full Changelog: [v1.12.0...v2.0.0](https://github.com/openai/openai-go/compare/v1.12.0...v2.0.0)

### Breaking changes

With the launch of `custom` tools in Chat Completions, `function` tools have been renamed to clarify the difference between the two.

`ChatCompletionToolParam` has become a union and is now named `ChatCompletionToolUnionParam`.

Older versions of the SDK used function tools: to migrate 


```diff
- openai.ChatCompletionToolParam{
-  Function: openai.FunctionDefinitionParam{
+ openai.ChatCompletionFunctionTool(
+  openai.FunctionDefinitionParam{
    Name:        "get_weather",
    Description: openai.String("Get weather at the given location"),
    Parameters: openai.FunctionParameters{ … },
+  },
+ )
- },
```

### Features

* **api:** adds GPT-5 and new API features: platform.openai.com/docs/guides/gpt-5 ([af46c88](https://github.com/openai/openai-go/commit/af46c885ea2414ba2b960f5d3accce89699a6250))
* **api:** manual updates ([219f209](https://github.com/openai/openai-go/commit/219f2092a6d7f1952d119b5b4ec32512956825ff))
* **client:** remove HTML escaping in JSON ([aea5ebc](https://github.com/openai/openai-go/commit/aea5ebccacb4fd854197dbf2547821860a62debc))
* **client:** rename union helpers ([645e881](https://github.com/openai/openai-go/commit/645e881dee5799d81fb4fd40d6494a296710d0ce))
* **client:** support optional json html escaping ([1d2336b](https://github.com/openai/openai-go/commit/1d2336b0d92f810fad3bf5faf5bf9e74975adf61))


### Bug Fixes

* **client:** revert path param changes ([9239f06](https://github.com/openai/openai-go/commit/9239f06bf0cb537d80980cee140a90d07b6d14f2))


### Chores

* change readme warning and minimum version ([1d0e22f](https://github.com/openai/openai-go/commit/1d0e22f85593a70f006f285f4461a05243b0fd74))
* document breaking changes ([afaa2b8](https://github.com/openai/openai-go/commit/afaa2b8482e8d10ea508716ad9b241517c9affa1))
* migrate examples ([9c57dd7](https://github.com/openai/openai-go/commit/9c57dd72515aab1c6d05d604870c5d0cf7fc1652))

## 1.12.0 (2025-07-30)

Full Changelog: [v1.11.1...v1.12.0](https://github.com/openai/openai-go/compare/v1.11.1...v1.12.0)

### Features

* **api:** manual updates ([16312ea](https://github.com/openai/openai-go/commit/16312ea2fea76c7cd2db4f38dfa10e0839f52d3e))


### Chores

* **client:** refactor streaming slightly to better future proof it ([0b9cb85](https://github.com/openai/openai-go/commit/0b9cb85a6bf0f2386e5db13aed34fbfad645efbe))

## 1.11.1 (2025-07-22)

Full Changelog: [v1.11.0...v1.11.1](https://github.com/openai/openai-go/compare/v1.11.0...v1.11.1)

### Bug Fixes

* **client:** process custom base url ahead of time ([cc1c23e](https://github.com/openai/openai-go/commit/cc1c23e3b1f4645004cb07b75816e3df445e73df))


### Chores

* **api:** event shapes more accurate ([2acd10d](https://github.com/openai/openai-go/commit/2acd10df4df52d1954d9ee3a98e5a4e56531533b))

## 1.11.0 (2025-07-16)

Full Changelog: [v1.10.3...v1.11.0](https://github.com/openai/openai-go/compare/v1.10.3...v1.11.0)

### Features

* **api:** manual updates ([97ed7fd](https://github.com/openai/openai-go/commit/97ed7fd1d432ad0144ec76bcebb61c9aaa1148de))

## 1.10.3 (2025-07-15)

Full Changelog: [v1.10.2...v1.10.3](https://github.com/openai/openai-go/compare/v1.10.2...v1.10.3)

## 1.10.2 (2025-07-15)

Full Changelog: [v1.10.1...v1.10.2](https://github.com/openai/openai-go/compare/v1.10.1...v1.10.2)

### Chores

* **api:** update realtime specs, build config ([3d2afda](https://github.com/openai/openai-go/commit/3d2afda006bd1f9e7ebde27b2873efa67e5e480d))

## 1.10.1 (2025-07-11)

Full Changelog: [v1.10.0...v1.10.1](https://github.com/openai/openai-go/compare/v1.10.0...v1.10.1)

### Chores

* **api:** specification cleanup ([5dbf6d2](https://github.com/openai/openai-go/commit/5dbf6d2cebe770d980db7888d705d1642ccd9cbc))
* lint tests in subpackages ([02f440d](https://github.com/openai/openai-go/commit/02f440dc6d899d7816b9fec9c47c09b393a7dd6c))

## 1.10.0 (2025-07-10)

Full Changelog: [v1.9.0...v1.10.0](https://github.com/openai/openai-go/compare/v1.9.0...v1.10.0)

### Features

* **api:** add file_url, fix event ID ([cb33971](https://github.com/openai/openai-go/commit/cb339714b65249844a87009192b2cf1508329673))

## 1.9.0 (2025-07-10)

Full Changelog: [v1.8.3...v1.9.0](https://github.com/openai/openai-go/compare/v1.8.3...v1.9.0)

### Features

* **client:** expand max streaming buffer size ([44390c8](https://github.com/openai/openai-go/commit/44390c81fdf33144f088b3ee8fef02269634dbe9))

## 1.8.3 (2025-07-08)

Full Changelog: [v1.8.2...v1.8.3](https://github.com/openai/openai-go/compare/v1.8.2...v1.8.3)

### Chores

* **ci:** only run for pushes and fork pull requests ([d6aab99](https://github.com/openai/openai-go/commit/d6aab99dadf267201add9812ba34ab2d5c70e0f4))
* **internal:** fix lint script for tests ([9c0a745](https://github.com/openai/openai-go/commit/9c0a74553c57ea5c29fb55f5ca2e122ca96031a4))
* lint tests ([2bd38d2](https://github.com/openai/openai-go/commit/2bd38d248cf2097254d1821a44c87827805732d1))

## 1.8.2 (2025-06-27)

Full Changelog: [v1.8.1...v1.8.2](https://github.com/openai/openai-go/compare/v1.8.1...v1.8.2)

### Bug Fixes

* don't try to deserialize as json when ResponseBodyInto is []byte ([74ad0f8](https://github.com/openai/openai-go/commit/74ad0f8fab0f956234503a9ba26fbd395944dcf8))
* **pagination:** check if page data is empty in GetNextPage ([c9becdc](https://github.com/openai/openai-go/commit/c9becdc9908f2a1961160837c6ab8cd9064e7854))

## 1.8.1 (2025-06-26)

Full Changelog: [v1.8.0...v1.8.1](https://github.com/openai/openai-go/compare/v1.8.0...v1.8.1)

### Chores

* **api:** remove unsupported property ([e22316a](https://github.com/openai/openai-go/commit/e22316adcd8f2c5aa672b12453cbd287de0e1878))
* **docs:** update README to include links to docs on Webhooks ([7bb8f85](https://github.com/openai/openai-go/commit/7bb8f8549fdd98997b1d145cbae98ff0146b4e43))

## 1.8.0 (2025-06-26)

Full Changelog: [v1.7.0...v1.8.0](https://github.com/openai/openai-go/compare/v1.7.0...v1.8.0)

### Features

* **api:** webhook and deep research support ([f6a7e7d](https://github.com/openai/openai-go/commit/f6a7e7dcd8801facc4f8d981f1ca43786c10de1e))


### Chores

* **internal:** add tests for breaking change detection ([339522d](https://github.com/openai/openai-go/commit/339522d38cd31b0753a8df37b8924f7e7dfb0b1d))

## 1.7.0 (2025-06-23)

Full Changelog: [v1.6.0...v1.7.0](https://github.com/openai/openai-go/compare/v1.6.0...v1.7.0)

### Features

* **api:** make model and inputs not required to create response ([19f0b76](https://github.com/openai/openai-go/commit/19f0b76378d35b3d81c60c85bf2e64d6bf85b9c2))
* **api:** update api shapes for usage and code interpreter ([d24d42c](https://github.com/openai/openai-go/commit/d24d42cba60e565627e8ffb1cac63a5085ddb6da))
* **client:** add escape hatch for null slice & maps ([9c633d6](https://github.com/openai/openai-go/commit/9c633d6f1dbcc0b153f42f831ee7e13d6fe62296))


### Chores

* fix documentation of null map ([8f3a134](https://github.com/openai/openai-go/commit/8f3a134e500b1b7791ab855adaef2d7b10d2d1c3))

## 1.6.0 (2025-06-17)

Full Changelog: [v1.5.0...v1.6.0](https://github.com/openai/openai-go/compare/v1.5.0...v1.6.0)

### Features

* **api:** add reusable prompt IDs ([280c698](https://github.com/openai/openai-go/commit/280c698015eba5f6bd47e2fce038eb401f6ef0f2))
* **api:** manual updates ([740f840](https://github.com/openai/openai-go/commit/740f84006ac283a25f5ad96aaf845a3c8a51c6ac))
* **client:** add debug log helper ([5715c49](https://github.com/openai/openai-go/commit/5715c491c483f8dab4ea2a900c400384f6810024))


### Chores

* **ci:** enable for pull requests ([9ed793a](https://github.com/openai/openai-go/commit/9ed793a51010423db464a7b7bd263d2fd275967f))

## 1.5.0 (2025-06-10)

Full Changelog: [v1.4.0...v1.5.0](https://github.com/openai/openai-go/compare/v1.4.0...v1.5.0)

### Features

* **api:** Add o3-pro model IDs ([3bbd0b8](https://github.com/openai/openai-go/commit/3bbd0b8f09030a6c571900d444742c4fc2a3c211))

## 1.4.0 (2025-06-09)

Full Changelog: [v1.3.0...v1.4.0](https://github.com/openai/openai-go/compare/v1.3.0...v1.4.0)

### Features

* **client:** allow overriding unions ([27c6299](https://github.com/openai/openai-go/commit/27c6299cb4ac275c6542b5691d81b795e65eeff6))


### Bug Fixes

* **client:** cast to raw message when converting to params ([a3282b0](https://github.com/openai/openai-go/commit/a3282b01a8d9a2c0cd04f24b298bf2ffcd160ebd))

## 1.3.0 (2025-06-03)

Full Changelog: [v1.2.1...v1.3.0](https://github.com/openai/openai-go/compare/v1.2.1...v1.3.0)

### Features

* **api:** add new realtime and audio models, realtime session options ([8b8f62b](https://github.com/openai/openai-go/commit/8b8f62b8e185f3fe4aaa99e892df5d35638931a1))

## 1.2.1 (2025-06-02)

Full Changelog: [v1.2.0...v1.2.1](https://github.com/openai/openai-go/compare/v1.2.0...v1.2.1)

### Bug Fixes

* **api:** Fix evals and code interpreter interfaces ([7e244c7](https://github.com/openai/openai-go/commit/7e244c73caad6b4768cced9a798452f03b1165c8))
* fix error ([a200fca](https://github.com/openai/openai-go/commit/a200fca92c3fa413cf724f424077d1537fa2ca3e))


### Chores

* make go mod tidy continue on error ([48f41c2](https://github.com/openai/openai-go/commit/48f41c2993bf6181018da859ae759951261f9ee2))

## 1.2.0 (2025-05-29)

Full Changelog: [v1.1.0...v1.2.0](https://github.com/openai/openai-go/compare/v1.1.0...v1.2.0)

### Features

* **api:** Config update for pakrym-stream-param ([84d59d5](https://github.com/openai/openai-go/commit/84d59d5cbc7521ddcc04435317903fd4ec3d17f6))


### Bug Fixes

* **client:** return binary content from `get /containers/{container_id}/files/{file_id}/content` ([f8c8de1](https://github.com/openai/openai-go/commit/f8c8de18b720b224267d54da53d7d919ed0fdff3))


### Chores

* deprecate Assistants API ([027470e](https://github.com/openai/openai-go/commit/027470e066ea6bbca1aeeb4fb9a8a3430babb84c))
* **internal:** fix release workflows ([fd46533](https://github.com/openai/openai-go/commit/fd4653316312755ccab7435fca9fb0a2d8bf8fbb))

## 1.1.0 (2025-05-22)

Full Changelog: [v1.0.0...v1.1.0](https://github.com/openai/openai-go/compare/v1.0.0...v1.1.0)

### Features

* **api:** add container endpoint ([2bd777d](https://github.com/openai/openai-go/commit/2bd777d6813b5dfcd3a2d339047a944c478dcd64))
* **api:** new API tools ([e7e2123](https://github.com/openai/openai-go/commit/e7e2123de7cafef515e07adde6edd45a7035b610))
* **api:** new streaming helpers for background responses ([422a0db](https://github.com/openai/openai-go/commit/422a0db3c674135e23dd200f5d8d785bd0be33e6))


### Chores

* **docs:** grammar improvements ([f4b23dd](https://github.com/openai/openai-go/commit/f4b23dd31facfc8839310854521b48060ef76be2))
* improve devcontainer setup ([dfdaeec](https://github.com/openai/openai-go/commit/dfdaeec2d6dd5cd679514d60c49b68c5df9e1b1e))

## 1.0.0 (2025-05-19)

Full Changelog: [v0.1.0-beta.11...v1.0.0](https://github.com/openai/openai-go/compare/v0.1.0-beta.11...v1.0.0)

### ⚠ BREAKING CHANGES

* **client:** rename file array param variant
* **api:** improve naming and remove assistants
* **accumulator:** update casing ([#401](https://github.com/openai/openai-go/issues/401))

### Features

* **api:** improve naming and remove assistants ([4c623b8](https://github.com/openai/openai-go/commit/4c623b88a9025db1961cc57985eb7374342f43e7))


### Bug Fixes

* **accumulator:** update casing ([#401](https://github.com/openai/openai-go/issues/401)) ([d59453c](https://github.com/openai/openai-go/commit/d59453c95b89fdd0b51305778dec0a39ce3a9d2a))
* **client:** correctly set stream key for multipart ([0ec68f0](https://github.com/openai/openai-go/commit/0ec68f0d779e7726931b1115eca9ae81eab59ba8))
* **client:** don't panic on marshal with extra null field ([9c15332](https://github.com/openai/openai-go/commit/9c153320272d212beaa516d4c70d54ae8053a958))
* **client:** increase max stream buffer size ([9456455](https://github.com/openai/openai-go/commit/945645559c5d68d9e28cf445d9c3b83e5fc6bd35))
* **client:** rename file array param variant ([4cfcf86](https://github.com/openai/openai-go/commit/4cfcf869280e7531fbbc8c00db0dd9271d07c423))
* **client:** use scanner for streaming ([aa58806](https://github.com/openai/openai-go/commit/aa58806bffc3aed68425c480414ddbb4dac3fa78))


### Chores

* **docs:** typo fix ([#400](https://github.com/openai/openai-go/issues/400)) ([bececf2](https://github.com/openai/openai-go/commit/bececf24cd0324b7c991b7d7f1d3eff6bf71f996))
* **examples:** migrate enum ([#447](https://github.com/openai/openai-go/issues/447)) ([814dd8b](https://github.com/openai/openai-go/commit/814dd8b6cfe4eeb535dc8ecd161a409ea2eb6698))
* **examples:** migrate to latest version ([#444](https://github.com/openai/openai-go/issues/444)) ([1c8754f](https://github.com/openai/openai-go/commit/1c8754ff905ed023f6381c8493910d63039407de))
* **examples:** remove beta assisstants examples ([#445](https://github.com/openai/openai-go/issues/445)) ([5891583](https://github.com/openai/openai-go/commit/589158372be9c0517b5508f9ccd872fdb1fe480b))
* **example:** update fine-tuning ([#450](https://github.com/openai/openai-go/issues/450)) ([421e3c5](https://github.com/openai/openai-go/commit/421e3c5065ace2d5ddd3d13a036477fff9123e5f))

## 0.1.0-beta.11 (2025-05-16)

Full Changelog: [v0.1.0-beta.10...v0.1.0-beta.11](https://github.com/openai/openai-go/compare/v0.1.0-beta.10...v0.1.0-beta.11)

### ⚠ BREAKING CHANGES

* **client:** clearer array variant names
* **client:** rename resp package
* **client:** improve core function names
* **client:** improve union variant names
* **client:** improve param subunions & deduplicate types

### Features

* **api:** add image sizes, reasoning encryption ([0852fb3](https://github.com/openai/openai-go/commit/0852fb3101dc940761f9e4f32875bfcf3669eada))
* **api:** add o3 and o4-mini model IDs ([3fabca6](https://github.com/openai/openai-go/commit/3fabca6b5c610edfb7bcd0cab5334a06444df0b0))
* **api:** Add reinforcement fine-tuning api support ([831a124](https://github.com/openai/openai-go/commit/831a12451cfce907b5ae4d294b9c2ac95f40d97a))
* **api:** adding gpt-4.1 family of model IDs ([1ef19d4](https://github.com/openai/openai-go/commit/1ef19d4cc94992dc435d7d5f28b30c9b1d255cd4))
* **api:** adding new image model support ([bf17880](https://github.com/openai/openai-go/commit/bf17880e182549c5c0fc34ec05df3184f223bc00))
* **api:** manual updates ([11f5716](https://github.com/openai/openai-go/commit/11f5716afa86aa100f80f3fa127e1d49203e5e21))
* **api:** responses x eval api ([183aaf7](https://github.com/openai/openai-go/commit/183aaf700f1d7ffad4ac847627d9ace65379c459))
* **api:** Updating Assistants and Evals API schemas ([47ca619](https://github.com/openai/openai-go/commit/47ca619fa1b439cf3a68c98e48e9bf1942f0568b))
* **client:** add dynamic streaming buffer to handle large lines ([8e6aad6](https://github.com/openai/openai-go/commit/8e6aad6d54fc73f1fcc174e1f06c9b3cf00c2689))
* **client:** add helper method to generate constant structs ([ff82809](https://github.com/openai/openai-go/commit/ff828094b561fc11184fed83f04424b6f68f7781))
* **client:** add support for endpoint-specific base URLs in python ([072dce4](https://github.com/openai/openai-go/commit/072dce46486d373fa0f0de5415f5270b01c2d972))
* **client:** add support for reading base URL from environment variable ([0d37268](https://github.com/openai/openai-go/commit/0d372687d673990290bad583f1906a2b121960b2))
* **client:** clearer array variant names ([a5d8b5d](https://github.com/openai/openai-go/commit/a5d8b5d6b161e3083184586840b2cbe0606d8de1))
* **client:** experimental support for unmarshalling into param structs ([5234875](https://github.com/openai/openai-go/commit/523487582e15a47e2f409f183568551258f4b8fe))
* **client:** improve param subunions & deduplicate types ([8a78f37](https://github.com/openai/openai-go/commit/8a78f37c25abf10498d16d210de3078f491ff23e))
* **client:** rename resp package ([4433516](https://github.com/openai/openai-go/commit/443351625ee290937a25425719b099ce785bd21b))
* **client:** support more time formats ([ec171b2](https://github.com/openai/openai-go/commit/ec171b2405c46f9cf04560760da001f7133d2fec))
* fix lint ([9c50a1e](https://github.com/openai/openai-go/commit/9c50a1eb9f93b578cb78085616f6bfab69f21dbc))


### Bug Fixes

* **client:** clean up reader resources ([710b92e](https://github.com/openai/openai-go/commit/710b92eaa7e94c03aeeca7479668677b32acb154))
* **client:** correctly update body in WithJSONSet ([f2d7118](https://github.com/openai/openai-go/commit/f2d7118295dd3073aa449426801d02e6f60bdaa3))
* **client:** improve core function names ([9f312a9](https://github.com/openai/openai-go/commit/9f312a9b14f5424d44d5834f1b82f3d3fcd57db2))
* **client:** improve union variant names ([a2c3de9](https://github.com/openai/openai-go/commit/a2c3de9e6c9f6e406b953f6de2eb78d1e72ec1b5))
* **client:** include path for type names in example code ([69561c5](https://github.com/openai/openai-go/commit/69561c549e18bd16a3641d62769479b125a4e955))
* **client:** resolve issue with optional multipart files ([910d173](https://github.com/openai/openai-go/commit/910d1730e97a03898e5dee7c889844a2ccec3e56))
* **client:** time format encoding fix ([ca17553](https://github.com/openai/openai-go/commit/ca175533ac8a17d36be1f531bbaa89c770da3f58))
* **client:** unmarshal responses properly ([fc9fec3](https://github.com/openai/openai-go/commit/fc9fec3c466ba9f633c3f7a4eebb5ebd3b85e8ac))
* handle empty bodies in WithJSONSet ([8372464](https://github.com/openai/openai-go/commit/83724640c6c00dcef1547dcabace309f17d14afc))
* **pagination:** handle errors when applying options ([eebf84b](https://github.com/openai/openai-go/commit/eebf84bf19f0eb6d9fa21e64bb83b0258e8cb42c))


### Chores

* **ci:** add timeout thresholds for CI jobs ([26b0dd7](https://github.com/openai/openai-go/commit/26b0dd760c142ca3aa287e8441bbe44cc8b3be0b))
* **ci:** only use depot for staging repos ([7682154](https://github.com/openai/openai-go/commit/7682154fdbcbe2a2ffdb2df590647a1712d52275))
* **ci:** run on more branches and use depot runners ([d7badbc](https://github.com/openai/openai-go/commit/d7badbc0d17bcf3cffec332f65cb68e531cb3176))
* **docs:** document pre-request options ([4befa5a](https://github.com/openai/openai-go/commit/4befa5a48ca61372715f36c45e72eb159d95bf2d))
* **docs:** update respjson package name ([9a00229](https://github.com/openai/openai-go/commit/9a002299a91e1145f053c51b1a4de10298fd2f43))
* **readme:** improve formatting ([a847e8d](https://github.com/openai/openai-go/commit/a847e8df45f725f9652fcea53ce57d3b9046efc7))
* **utils:** add internal resp to param utility ([239c4e2](https://github.com/openai/openai-go/commit/239c4e2cb32c7af71ab14668ccc2f52ea59653f9))


### Documentation

* update documentation links to be more uniform ([f5f0bb0](https://github.com/openai/openai-go/commit/f5f0bb05ee705d84119806f8e703bf2e0becb1fa))

## 0.1.0-beta.10 (2025-04-14)

Full Changelog: [v0.1.0-beta.9...v0.1.0-beta.10](https://github.com/openai/openai-go/compare/v0.1.0-beta.9...v0.1.0-beta.10)

### Chores

* **internal:** expand CI branch coverage ([#369](https://github.com/openai/openai-go/issues/369)) ([258dda8](https://github.com/openai/openai-go/commit/258dda8007a69b9c2720b225ee6d27474d676a93))
* **internal:** reduce CI branch coverage ([a2f7c03](https://github.com/openai/openai-go/commit/a2f7c03eb984d98f29f908df103ea1743f2e3d9a))

## 0.1.0-beta.9 (2025-04-09)

Full Changelog: [v0.1.0-beta.8...v0.1.0-beta.9](https://github.com/openai/openai-go/compare/v0.1.0-beta.8...v0.1.0-beta.9)

### Chores

* workaround build errors ([#366](https://github.com/openai/openai-go/issues/366)) ([adeb003](https://github.com/openai/openai-go/commit/adeb003cab8efbfbf4424e03e96a0f5e728551cb))

## 0.1.0-beta.8 (2025-04-09)

Full Changelog: [v0.1.0-beta.7...v0.1.0-beta.8](https://github.com/openai/openai-go/compare/v0.1.0-beta.7...v0.1.0-beta.8)

### Features

* **api:** Add evalapi to sdk ([#360](https://github.com/openai/openai-go/issues/360)) ([88977d1](https://github.com/openai/openai-go/commit/88977d1868dbbe0060c56ba5dac8eb19773e4938))
* **api:** manual updates ([#363](https://github.com/openai/openai-go/issues/363)) ([5d068e0](https://github.com/openai/openai-go/commit/5d068e0053172db7f5b75038aa215eee074eeeed))
* **client:** add escape hatch to omit required param fields ([#354](https://github.com/openai/openai-go/issues/354)) ([9690d6b](https://github.com/openai/openai-go/commit/9690d6b49f8b00329afc038ec15116750853e620))
* **client:** support custom http clients ([#357](https://github.com/openai/openai-go/issues/357)) ([b5a624f](https://github.com/openai/openai-go/commit/b5a624f658cad774094427b36b05e446b41e8c52))


### Chores

* **docs:** readme improvements ([#356](https://github.com/openai/openai-go/issues/356)) ([b2f8539](https://github.com/openai/openai-go/commit/b2f8539d6316e3443aa733be2c95926696119c13))
* **internal:** fix examples ([#361](https://github.com/openai/openai-go/issues/361)) ([de398b4](https://github.com/openai/openai-go/commit/de398b453d398299eb80c15f8fdb2bcbef5eeed6))
* **internal:** skip broken test ([#362](https://github.com/openai/openai-go/issues/362)) ([cccead9](https://github.com/openai/openai-go/commit/cccead9ba916142ac8fbe6e8926d706511e32ae3))
* **tests:** improve enum examples ([#359](https://github.com/openai/openai-go/issues/359)) ([e0b9739](https://github.com/openai/openai-go/commit/e0b9739920114d6e991d3947b67fdf62cfaa09c7))

## 0.1.0-beta.7 (2025-04-07)

Full Changelog: [v0.1.0-beta.6...v0.1.0-beta.7](https://github.com/openai/openai-go/compare/v0.1.0-beta.6...v0.1.0-beta.7)

### Features

* **client:** make response union's AsAny method type safe ([#352](https://github.com/openai/openai-go/issues/352)) ([1252f56](https://github.com/openai/openai-go/commit/1252f56c917e57d6d2b031501b2ff5f89f87cf87))


### Chores

* **docs:** doc improvements ([#350](https://github.com/openai/openai-go/issues/350)) ([80debc8](https://github.com/openai/openai-go/commit/80debc824eaacb4b07c8f3e8b1d0488d860d5be5))

## 0.1.0-beta.6 (2025-04-04)

Full Changelog: [v0.1.0-beta.5...v0.1.0-beta.6](https://github.com/openai/openai-go/compare/v0.1.0-beta.5...v0.1.0-beta.6)

### Features

* **api:** manual updates ([4e39609](https://github.com/openai/openai-go/commit/4e39609d499b88039f1c90cc4b56e26f28fd58ea))
* **client:** support unions in query and forms ([#347](https://github.com/openai/openai-go/issues/347)) ([cf8af37](https://github.com/openai/openai-go/commit/cf8af373ab7c019c75e886855009ffaca320d0e3))

## 0.1.0-beta.5 (2025-04-03)

Full Changelog: [v0.1.0-beta.4...v0.1.0-beta.5](https://github.com/openai/openai-go/compare/v0.1.0-beta.4...v0.1.0-beta.5)

### Features

* **api:** manual updates ([563cc50](https://github.com/openai/openai-go/commit/563cc505f2ab17749bb77e937342a6614243b975))
* **client:** omitzero on required id parameter ([#339](https://github.com/openai/openai-go/issues/339)) ([c0b4842](https://github.com/openai/openai-go/commit/c0b484266ccd9faee66873916d8c0c92ea9f1014))


### Bug Fixes

* **client:** return error on bad custom url instead of panic ([#341](https://github.com/openai/openai-go/issues/341)) ([a06c5e6](https://github.com/openai/openai-go/commit/a06c5e632242e53d3fdcc8964931acb533a30b7e))
* **client:** support multipart encoding array formats ([#342](https://github.com/openai/openai-go/issues/342)) ([5993b28](https://github.com/openai/openai-go/commit/5993b28309d02c2d748b54d98934ef401dcd193a))
* **client:** unmarshal stream events into fresh memory ([#340](https://github.com/openai/openai-go/issues/340)) ([52c3e08](https://github.com/openai/openai-go/commit/52c3e08f51d471d728e5acd16b3c304b51be2d03))

## 0.1.0-beta.4 (2025-04-02)

Full Changelog: [v0.1.0-beta.3...v0.1.0-beta.4](https://github.com/openai/openai-go/compare/v0.1.0-beta.3...v0.1.0-beta.4)

### Features

* **api:** manual updates ([bc4fe73](https://github.com/openai/openai-go/commit/bc4fe73eec9c4d39229e4beae8eaafb55b1d3364))
* **api:** manual updates ([aa7ff10](https://github.com/openai/openai-go/commit/aa7ff10b0616a6b2ece45cb10e9c83f25e35aded))


### Chores

* **docs:** update file uploads in README ([#333](https://github.com/openai/openai-go/issues/333)) ([471c452](https://github.com/openai/openai-go/commit/471c4525c94e83cf4b78cb6c9b2f65a8a27bf3ce))
* **internal:** codegen related update ([#335](https://github.com/openai/openai-go/issues/335)) ([48422dc](https://github.com/openai/openai-go/commit/48422dcca333ab808ccb02506c033f1c69d2aa19))
* Remove deprecated/unused remote spec feature ([c5077a1](https://github.com/openai/openai-go/commit/c5077a154a6db79b73cf4978bdc08212c6da6423))

## 0.1.0-beta.3 (2025-03-28)

Full Changelog: [v0.1.0-beta.2...v0.1.0-beta.3](https://github.com/openai/openai-go/compare/v0.1.0-beta.2...v0.1.0-beta.3)

### ⚠ BREAKING CHANGES

* **client:** add enums ([#327](https://github.com/openai/openai-go/issues/327))

### Features

* **api:** add `get /chat/completions` endpoint ([e8ed116](https://github.com/openai/openai-go/commit/e8ed1168576c885cb26fbf819b9c8d24975749bd))
* **api:** add `get /responses/{response_id}/input_items` endpoint ([8870c26](https://github.com/openai/openai-go/commit/8870c26f010a596adcf37ac10dba096bdd4394e3))


### Bug Fixes

* **client:** add enums ([#327](https://github.com/openai/openai-go/issues/327)) ([b0e3afb](https://github.com/openai/openai-go/commit/b0e3afbd6f18fd9fc2a5ea9174bd7ec0ac0614db))


### Chores

* add hash of OpenAPI spec/config inputs to .stats.yml ([104b786](https://github.com/openai/openai-go/commit/104b7861bb025514999b143f7d1de45d2dab659f))
* add request options to client tests ([#321](https://github.com/openai/openai-go/issues/321)) ([f5239ce](https://github.com/openai/openai-go/commit/f5239ceecf36835341eac5121ed1770020c4806a))
* **api:** updates to supported Voice IDs ([#325](https://github.com/openai/openai-go/issues/325)) ([477727a](https://github.com/openai/openai-go/commit/477727a44b0fb72493c4749cc60171e0d30f98ec))
* **docs:** improve security documentation ([#319](https://github.com/openai/openai-go/issues/319)) ([0271053](https://github.com/openai/openai-go/commit/027105363ab30ac3e189234908169faf94e0ca49))
* fix typos ([#324](https://github.com/openai/openai-go/issues/324)) ([dba15f7](https://github.com/openai/openai-go/commit/dba15f74d63814ce16f778e1017a209a42f46179))

## 0.1.0-beta.2 (2025-03-22)

Full Changelog: [v0.1.0-beta.1...v0.1.0-beta.2](https://github.com/openai/openai-go/compare/v0.1.0-beta.1...v0.1.0-beta.2)

### Bug Fixes

* **client:** elide fields in ToAssistantParam ([#309](https://github.com/openai/openai-go/issues/309)) ([1fcd837](https://github.com/openai/openai-go/commit/1fcd83753ea806745d278a5b94797bbee0f018ed))

## 0.1.0-beta.1 (2025-03-22)

Full Changelog: [v0.1.0-alpha.67...v0.1.0-beta.1](https://github.com/openai/openai-go/compare/v0.1.0-alpha.67...v0.1.0-beta.1)

### Chores

* **docs:** clarify breaking changes ([#306](https://github.com/openai/openai-go/issues/306)) ([db4bd1f](https://github.com/openai/openai-go/commit/db4bd1f5304aa523a6b62da6e2571487d4248518))

## 0.1.0-alpha.67 (2025-03-21)

Full Changelog: [v0.1.0-alpha.66...v0.1.0-alpha.67](https://github.com/openai/openai-go/compare/v0.1.0-alpha.66...v0.1.0-alpha.67)

### ⚠ BREAKING CHANGES

* **api:** migrate to v2

### Features

* **api:** migrate to v2 ([9377508](https://github.com/openai/openai-go/commit/9377508e45ae485d11c3199d6d3d91d345f1b76e))
* **api:** new models for TTS, STT, + new audio features for Realtime ([#298](https://github.com/openai/openai-go/issues/298)) ([48fa064](https://github.com/openai/openai-go/commit/48fa064202a6e4a3e850d435b29f6fe9a1fe53f4))


### Chores

* **internal:** bugfix ([0d8c1f4](https://github.com/openai/openai-go/commit/0d8c1f4e801785728b6ad3342146fe38874d6c04))


### Documentation

* add migration guide ([#302](https://github.com/openai/openai-go/issues/302)) ([19e32fa](https://github.com/openai/openai-go/commit/19e32fa595e65048bb129e813c697991117abca2))
