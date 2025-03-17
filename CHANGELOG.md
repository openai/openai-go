# Changelog

## 0.1.0-alpha.65 (2025-03-17)

Full Changelog: [v0.1.0-alpha.64...v0.1.0-alpha.65](https://github.com/openai/openai-go/compare/v0.1.0-alpha.64...v0.1.0-alpha.65)

### Bug Fixes

* **ci:** add workflow back in ([f695f00](https://github.com/openai/openai-go/commit/f695f0013d805623548b42c346a167ee15dd9405))
* revert bad update ([e1ebde6](https://github.com/openai/openai-go/commit/e1ebde6c189d93efd9bae193a8686242446b86f3))

## 0.1.0-alpha.64 (2025-03-15)

Full Changelog: [v0.1.0-alpha.63...v0.1.0-alpha.64](https://github.com/openai/openai-go/compare/v0.1.0-alpha.63...v0.1.0-alpha.64)

### ⚠ BREAKING CHANGES

* **client:** improve naming of some variants ([#89](https://github.com/openai/openai-go/issues/89))

### Features

* add azure, examples, and message constructors ([fb2df0f](https://github.com/openai/openai-go/commit/fb2df0fe22002f1826bfaa1cb008c45db375885c))
* add SKIP_BREW env var to ./scripts/bootstrap ([#255](https://github.com/openai/openai-go/issues/255)) ([151c5e7](https://github.com/openai/openai-go/commit/151c5e7106b467174703089ca21845780f121c03))
* add support for error property in stream ([#29](https://github.com/openai/openai-go/issues/29)) ([0c7d6e5](https://github.com/openai/openai-go/commit/0c7d6e5fc62fd0ed686d5fd254376ff35eb9903a))
* **api:** add chatgpt-4o-latest model ([#24](https://github.com/openai/openai-go/issues/24)) ([110d1f0](https://github.com/openai/openai-go/commit/110d1f025a3e9259ae87b8ed7c4824de32d489fb))
* **api:** add file search result details to run steps ([#32](https://github.com/openai/openai-go/issues/32)) ([c1862bd](https://github.com/openai/openai-go/commit/c1862bda1d70ce21d4af85583f975a3aa9cc1e1f))
* **api:** add gpt-4.5-preview ([#242](https://github.com/openai/openai-go/issues/242)) ([961bf16](https://github.com/openai/openai-go/commit/961bf16109133c94298241f1240d2c0203c8ead7))
* **api:** add gpt-4o-2024-11-20 model ([#131](https://github.com/openai/openai-go/issues/131)) ([8fe1011](https://github.com/openai/openai-go/commit/8fe1011b33c5d05a14fe249e520894eec7c74334))
* **api:** add gpt-4o-audio-preview model for chat completions ([#88](https://github.com/openai/openai-go/issues/88)) ([f4a76d0](https://github.com/openai/openai-go/commit/f4a76d0aff1127efe4faa8052334a39f6303cada))
* **api:** add new, expressive voices for Realtime and Audio in Chat Completions ([#101](https://github.com/openai/openai-go/issues/101)) ([f946acc](https://github.com/openai/openai-go/commit/f946acc71a92f885bed87f0d4e724fb40cae0f14))
* **api:** add o1 models ([#49](https://github.com/openai/openai-go/issues/49)) ([698a0c9](https://github.com/openai/openai-go/commit/698a0c9c09ffc8b070a796bf4b10b4ba39a17815))
* **api:** add o3-mini ([#195](https://github.com/openai/openai-go/issues/195)) ([1dc8887](https://github.com/openai/openai-go/commit/1dc888754d531b7d18768b17a44c5415fa6bf3ea))
* **api:** add omni-moderation model ([#63](https://github.com/openai/openai-go/issues/63)) ([7402f24](https://github.com/openai/openai-go/commit/7402f24732a207c67c98248d58431520081bc324))
* **api:** add support for predicted outputs ([#110](https://github.com/openai/openai-go/issues/110)) ([73c798a](https://github.com/openai/openai-go/commit/73c798a65bd0aa7d241d5a2fa21eb8c880f8e769))
* **api:** add support for storing chat completions ([#228](https://github.com/openai/openai-go/issues/228)) ([3da23d8](https://github.com/openai/openai-go/commit/3da23d8f861f8c73c91a7d1443a35ea20475f91e))
* **api:** new o1 and GPT-4o models + preference fine-tuning ([#142](https://github.com/openai/openai-go/issues/142)) ([a9e2f35](https://github.com/openai/openai-go/commit/a9e2f35573f21296a452955d6ded82385812f649))
* **api:** support storing chat completions, enabling evals and model distillation in the dashboard ([#72](https://github.com/openai/openai-go/issues/72)) ([b0eae50](https://github.com/openai/openai-go/commit/b0eae50c466a7f6628da792ecbef97a398489a69))
* **api:** update enum values, comments, and examples ([#181](https://github.com/openai/openai-go/issues/181)) ([29e5479](https://github.com/openai/openai-go/commit/29e547924355b3d0cd64bfed807418b7cacc179e))
* **api:** updates ([#138](https://github.com/openai/openai-go/issues/138)) ([525573e](https://github.com/openai/openai-go/commit/525573e0d8bb0d69b183e418ded2b987f7da7c73))
* **api:** updates ([#259](https://github.com/openai/openai-go/issues/259)) ([aa5cb47](https://github.com/openai/openai-go/commit/aa5cb47722c42f17808a15cf15dd58a7a3cfa17f))
* **api:** updates ([#5](https://github.com/openai/openai-go/issues/5)) ([f92a25c](https://github.com/openai/openai-go/commit/f92a25c87da861702e1792b6799b6601861d4d01))
* **assistants:** add polling helpers and examples ([#84](https://github.com/openai/openai-go/issues/84)) ([eaa9194](https://github.com/openai/openai-go/commit/eaa91946e78c278f0eacffd6f19c4691d9e23eeb))
* **client:** accept RFC6838 JSON content types ([#256](https://github.com/openai/openai-go/issues/256)) ([9a8f472](https://github.com/openai/openai-go/commit/9a8f472bc60a208f446d5a56a60a3b824e486786))
* **client:** allow custom baseurls without trailing slash ([#254](https://github.com/openai/openai-go/issues/254)) ([32b7eb4](https://github.com/openai/openai-go/commit/32b7eb453c7d8a7e2982877baff56c46e3e47c50))
* **client:** improve default client options support ([0c621c0](https://github.com/openai/openai-go/commit/0c621c0572cedf1b2d7976020d037d688c9956b7))
* **client:** improve default client options support ([#266](https://github.com/openai/openai-go/issues/266)) ([e68b1cd](https://github.com/openai/openai-go/commit/e68b1cd08f7ed7f008b9b006fc3f9ccb2497fa27))
* **client:** improve naming of some variants ([#89](https://github.com/openai/openai-go/issues/89)) ([6bb0f75](https://github.com/openai/openai-go/commit/6bb0f75b9b00fcdeffa9c19aa5de0cc48e12c168))
* **client:** send `X-Stainless-Timeout` header ([#204](https://github.com/openai/openai-go/issues/204)) ([72405f0](https://github.com/openai/openai-go/commit/72405f00fbb168f85b9b8b40aee8a0dbadce67f3))
* **client:** send retry count header ([#60](https://github.com/openai/openai-go/issues/60)) ([01ed6ab](https://github.com/openai/openai-go/commit/01ed6ab70bae02e882cd4ffd4e0455f0ab03511e))
* **examples/structure-outputs:** created an example for using structured outputs ([d4303e8](https://github.com/openai/openai-go/commit/d4303e8f8c4cfd7f260681ebed600749b6262028))
* extract out `ImageModel`, `AudioModel`, `SpeechModel` ([#3](https://github.com/openai/openai-go/issues/3)) ([4b90869](https://github.com/openai/openai-go/commit/4b90869569ec034cb78995d9c56720cad5920577))
* make enums not nominal ([#4](https://github.com/openai/openai-go/issues/4)) ([a559359](https://github.com/openai/openai-go/commit/a55935938f21c2a4b1d682538bd3481861a126a3))
* move pagination package from internal to packages ([#81](https://github.com/openai/openai-go/issues/81)) ([c7476f7](https://github.com/openai/openai-go/commit/c7476f78f69b23044436236081f0faf405c98b1c))
* **pagination:** avoid fetching when has_more: false ([#218](https://github.com/openai/openai-go/issues/218)) ([978707d](https://github.com/openai/openai-go/commit/978707dd3b64230c7835e83548167c8617444462))
* publish ([c329601](https://github.com/openai/openai-go/commit/c329601324226e28ff18d6ccecfdde41cedd3b5a))
* simplify content union ([#18](https://github.com/openai/openai-go/issues/18)) ([b228103](https://github.com/openai/openai-go/commit/b2281035a6d6abeb02a912a3543c98ac70ff456f))
* **stream-accumulators:** added streaming accumulator helpers and example ([ecfdb64](https://github.com/openai/openai-go/commit/ecfdb64214f4a8d75bff6a29d068c55357ed7815))
* support assistants stream ([7c00c63](https://github.com/openai/openai-go/commit/7c00c6340ffed1557b5763bba1790fb0f55c8530))
* support deprecated markers ([#178](https://github.com/openai/openai-go/issues/178)) ([2c21e34](https://github.com/openai/openai-go/commit/2c21e3404639defedcad990acc1f530bbcb2781f))
* **vector store:** improve chunking strategy type names ([#40](https://github.com/openai/openai-go/issues/40)) ([5d0740f](https://github.com/openai/openai-go/commit/5d0740f77b32c89fc2bebfea98fe8b939d78a155))


### Bug Fixes

* **api/types:** correct audio duration & role types ([#209](https://github.com/openai/openai-go/issues/209)) ([f480273](https://github.com/openai/openai-go/commit/f480273953cb9eef9b367867dfc9a3b4af3ea858))
* **api:** add missing file rank enum + more metadata ([#248](https://github.com/openai/openai-go/issues/248)) ([a47089e](https://github.com/openai/openai-go/commit/a47089e5d9ce08ba0fc6461e3d4bb8941ca79291))
* **api:** add missing reasoning effort + model enums ([#215](https://github.com/openai/openai-go/issues/215)) ([5b53a1d](https://github.com/openai/openai-go/commit/5b53a1d531a05599d8d517b3c1e98c8c8e40ccbd))
* **api:** escape key values when encoding maps ([#116](https://github.com/openai/openai-go/issues/116)) ([a29c08e](https://github.com/openai/openai-go/commit/a29c08e27eb6fb403cfca60247bf1d23d82ec26e))
* **audio:** correct response_format translations type ([#62](https://github.com/openai/openai-go/issues/62)) ([c46777b](https://github.com/openai/openai-go/commit/c46777ba01b6d0f3be506521ed5c7420f1fd0f20))
* **beta:** pass beta header by default ([#75](https://github.com/openai/openai-go/issues/75)) ([e0a5caa](https://github.com/openai/openai-go/commit/e0a5caa65f645c3531bc2a9795b0571949305dbd))
* **client:** don't truncate manually specified filenames ([#230](https://github.com/openai/openai-go/issues/230)) ([86febfc](https://github.com/openai/openai-go/commit/86febfce3ae78704595e4450cb4e1373e7932a78))
* **client:** no panic on missing BaseURL ([#121](https://github.com/openai/openai-go/issues/121)) ([9e252ee](https://github.com/openai/openai-go/commit/9e252ee90db5549d910d46dd998fd499dac60b22))
* correct required fields for flattened unions ([#120](https://github.com/openai/openai-go/issues/120)) ([9d6e6f2](https://github.com/openai/openai-go/commit/9d6e6f2ace15ac7c843b0585265b1cea1e46e78d))
* deserialization of struct unions that implement json.Unmarshaler ([#11](https://github.com/openai/openai-go/issues/11)) ([7c0847a](https://github.com/openai/openai-go/commit/7c0847aa2ae15b4442ab0625d8a780ed684c275e))
* do not call path.Base on ContentType ([#225](https://github.com/openai/openai-go/issues/225)) ([e1c1a55](https://github.com/openai/openai-go/commit/e1c1a55e7de733fa70194a3475bcf943c0568ba7))
* **examples/fine-tuning:** used an old constant name ([#34](https://github.com/openai/openai-go/issues/34)) ([fb78f1d](https://github.com/openai/openai-go/commit/fb78f1df84dc5a5187157eeeb9a21b3853c328f6))
* **example:** use correct model ([#86](https://github.com/openai/openai-go/issues/86)) ([0f3578d](https://github.com/openai/openai-go/commit/0f3578d2c3e4f40bfeea511e385d5c74d2b06e4e))
* fix apijson.Port for embedded structs ([#174](https://github.com/openai/openai-go/issues/174)) ([13ecf6c](https://github.com/openai/openai-go/commit/13ecf6c0b82085e54e56807cb7dcd981768c50ca))
* fix apijson.Port for embedded structs ([#177](https://github.com/openai/openai-go/issues/177)) ([1820b25](https://github.com/openai/openai-go/commit/1820b2592eee7aa55212c3b0bcfc3077c5d7f3a7))
* fix early cancel when RequestTimeout is provided for streaming requests ([#221](https://github.com/openai/openai-go/issues/221)) ([08320be](https://github.com/openai/openai-go/commit/08320be359825fbceb4642aa6a60142ceec8ce2e))
* fix unicode encoding for json ([#193](https://github.com/openai/openai-go/issues/193)) ([4a905a7](https://github.com/openai/openai-go/commit/4a905a7a596d7e1a7f0de9bec7f49c1026d36621))
* flush stream response when done event is sent ([#172](https://github.com/openai/openai-go/issues/172)) ([cf1a6a5](https://github.com/openai/openai-go/commit/cf1a6a5b384fbd3de948494d50a1b19cfba79fdd))
* **requestconfig:** copy over more fields when cloning ([#44](https://github.com/openai/openai-go/issues/44)) ([3c7aa48](https://github.com/openai/openai-go/commit/3c7aa48609798b665b61a222067af0f8bbdc0f40))
* **responses:** correct computer use enum value ([#261](https://github.com/openai/openai-go/issues/261)) ([73d820b](https://github.com/openai/openai-go/commit/73d820bd96debf858adee674ec330f7b9d477342))
* **responses:** correct reasoning output type ([#262](https://github.com/openai/openai-go/issues/262)) ([a636da3](https://github.com/openai/openai-go/commit/a636da379ea7f41ed73e320fa85f06219a71e8ac))
* **stream:** ensure .Close() doesn't panic ([#194](https://github.com/openai/openai-go/issues/194)) ([71821a8](https://github.com/openai/openai-go/commit/71821a8938562f79d859d8939b9d20b2ef2ad3ca))
* **stream:** ensure .Close() doesn't panic ([#201](https://github.com/openai/openai-go/issues/201)) ([a75c812](https://github.com/openai/openai-go/commit/a75c812270296b3c308b4f044e2b81071bdb8f8b))
* **streaming:** correctly accumulate tool calls and roles ([#55](https://github.com/openai/openai-go/issues/55)) ([321ff9e](https://github.com/openai/openai-go/commit/321ff9e778b3c70addac90abbbc5e06e63581de2))
* **types:** correct metadata type + other fixes ([1dc8887](https://github.com/openai/openai-go/commit/1dc888754d531b7d18768b17a44c5415fa6bf3ea))
* update stream error handling ([#213](https://github.com/openai/openai-go/issues/213)) ([b2a763d](https://github.com/openai/openai-go/commit/b2a763dcde75cc6de5ee2b7d5e18fdcda1aaa90c))


### Chores

* add back custom code that was reverted ([4b46d02](https://github.com/openai/openai-go/commit/4b46d02aa506b9b521c88abcd5ea78ffe5083088))
* **api:** bump spec version ([#154](https://github.com/openai/openai-go/issues/154)) ([4fc775f](https://github.com/openai/openai-go/commit/4fc775fa9bd22a0a0df370e2ef20f6485a11f3e3))
* **api:** delete deprecated method ([#208](https://github.com/openai/openai-go/issues/208)) ([dcb01cc](https://github.com/openai/openai-go/commit/dcb01cc7ba479e486b6bc57c73b46b634f316dca))
* bump Go to v1.21 ([#12](https://github.com/openai/openai-go/issues/12)) ([e4c3228](https://github.com/openai/openai-go/commit/e4c322840f7ea441a78dd723e46e27b048d16c1e))
* bump license year ([#151](https://github.com/openai/openai-go/issues/151)) ([5e724f9](https://github.com/openai/openai-go/commit/5e724f9816b2d9697970c56c62adab3a62f2b6cd))
* bump openapi url ([#136](https://github.com/openai/openai-go/issues/136)) ([b9bf99d](https://github.com/openai/openai-go/commit/b9bf99dbf047cb3a32efafa7bcf6979ac7ba15ab))
* **ci:** bump prism mock server version ([#10](https://github.com/openai/openai-go/issues/10)) ([00f9455](https://github.com/openai/openai-go/commit/00f9455692c52fb37544d3f657090b216667d8ec))
* **ci:** codeowners file ([#9](https://github.com/openai/openai-go/issues/9)) ([be41ac2](https://github.com/openai/openai-go/commit/be41ac2ce87efacf17748cb9dd2d3b1b4a43180e))
* **docs:** add docstring explaining streaming pattern ([#205](https://github.com/openai/openai-go/issues/205)) ([bfabf9d](https://github.com/openai/openai-go/commit/bfabf9d23191c2050dde7844da06bab321b5a361))
* **docs:** fix maxium typo ([#69](https://github.com/openai/openai-go/issues/69)) ([29dfb56](https://github.com/openai/openai-go/commit/29dfb56cb755fce02e0d923bcd742f9d780151b7))
* **docs:** remove some duplicative api.md entries ([#65](https://github.com/openai/openai-go/issues/65)) ([532c9a0](https://github.com/openai/openai-go/commit/532c9a05d7637d891aa76a9ef2ec3254c0f20942))
* **examples:** minor formatting changes ([#14](https://github.com/openai/openai-go/issues/14)) ([85aaaa5](https://github.com/openai/openai-go/commit/85aaaa5af7242ca6a2d12a61a127d4b1ae08f7d9))
* fix GetNextPage docstring ([#78](https://github.com/openai/openai-go/issues/78)) ([a736116](https://github.com/openai/openai-go/commit/a736116b2544976afe6a4f95ceb066fc959ece67))
* **internal:** fix devcontainers setup ([#236](https://github.com/openai/openai-go/issues/236)) ([25b0137](https://github.com/openai/openai-go/commit/25b0137b0d161554ed9ce0cab2b71c60aa7a17e0))
* **internal:** remove CI condition ([#271](https://github.com/openai/openai-go/issues/271)) ([c61ef3a](https://github.com/openai/openai-go/commit/c61ef3a83a5ac7cf810331c1f5c8b4bffd385b19))
* **internal:** remove extra empty newlines ([#268](https://github.com/openai/openai-go/issues/268)) ([df22608](https://github.com/openai/openai-go/commit/df2260865fb588c84707a067128ff03561cccef6))
* **internal:** rename `streaming.go` ([#176](https://github.com/openai/openai-go/issues/176)) ([fb192c6](https://github.com/openai/openai-go/commit/fb192c675b008273962f88ca599dd84f5ddf7823))
* **internal:** spec update ([#130](https://github.com/openai/openai-go/issues/130)) ([7fd444a](https://github.com/openai/openai-go/commit/7fd444a5724bdeef6cdc4b9e9db4f445eaa020e2))
* **internal:** spec update ([#145](https://github.com/openai/openai-go/issues/145)) ([b8ba547](https://github.com/openai/openai-go/commit/b8ba547746107d6c2081bc66e746b82114ac9b44))
* **internal:** spec update ([#146](https://github.com/openai/openai-go/issues/146)) ([d4bcfc0](https://github.com/openai/openai-go/commit/d4bcfc0fd6906d50c2e522e18d3d7cc154b0cc71))
* **internal:** spec update ([#158](https://github.com/openai/openai-go/issues/158)) ([fd7fe8c](https://github.com/openai/openai-go/commit/fd7fe8c49e12cf2bee9b96b1af0ad509ee745fa0))
* **internal:** streaming refactors ([#165](https://github.com/openai/openai-go/issues/165)) ([2fbb02c](https://github.com/openai/openai-go/commit/2fbb02c1c18b92e966266b3340ee15deaf408d34))
* **internal:** update spec link ([#53](https://github.com/openai/openai-go/issues/53)) ([4915187](https://github.com/openai/openai-go/commit/4915187abeec2e3a1e6290004fc4c18a7f1129ca))
* **internal:** update spec version ([#95](https://github.com/openai/openai-go/issues/95)) ([ba37fcc](https://github.com/openai/openai-go/commit/ba37fcc41ef1c5d63a218b9bf9139029264bc274))
* **internal:** updates ([#2](https://github.com/openai/openai-go/issues/2)) ([103d454](https://github.com/openai/openai-go/commit/103d454b1d1fab5032c6af54e62ffe65f0cacb35))
* **internal:** updates ([#6](https://github.com/openai/openai-go/issues/6)) ([316e623](https://github.com/openai/openai-go/commit/316e6231c27728f4031f822287389c67e914739a))
* move ChatModel type to shared ([#250](https://github.com/openai/openai-go/issues/250)) ([304ec6b](https://github.com/openai/openai-go/commit/304ec6b670cfea141177e098b622d517c2aa3a9d))
* refactor client tests ([#187](https://github.com/openai/openai-go/issues/187)) ([b956e3a](https://github.com/openai/openai-go/commit/b956e3a350d27965407154319a30eef09131fc3c))
* **tests:** limit array example length ([#128](https://github.com/openai/openai-go/issues/128)) ([45fa490](https://github.com/openai/openai-go/commit/45fa4909359bbcda7acc6a72a1cf2f1e709c7cf1))
* **types:** define FilePurpose enum ([#22](https://github.com/openai/openai-go/issues/22)) ([1daff0f](https://github.com/openai/openai-go/commit/1daff0f73ce40ae61e2d4b217d22811cbb373c27))
* **types:** improve type name for embedding models ([#57](https://github.com/openai/openai-go/issues/57)) ([57736f9](https://github.com/openai/openai-go/commit/57736f9c857865e5a26c62b8cfba3c064a784e7b))
* **types:** rename vector store chunking strategy ([#169](https://github.com/openai/openai-go/issues/169)) ([6076a58](https://github.com/openai/openai-go/commit/6076a584b57291ae915dd0596c2b6a4331a4f080))


### Documentation

* add missing docs for some enums ([#114](https://github.com/openai/openai-go/issues/114)) ([3d9fbd8](https://github.com/openai/openai-go/commit/3d9fbd85098e8ff5ee89872846c8151028179299))
* document raw responses ([#197](https://github.com/openai/openai-go/issues/197)) ([9d06d7a](https://github.com/openai/openai-go/commit/9d06d7a1a78431ba01e6018fe6c2a9a4b75c3db6))
* **examples:** fix typo ([#207](https://github.com/openai/openai-go/issues/207)) ([bf8afc3](https://github.com/openai/openai-go/commit/bf8afc3a81676b7675c314aa1c1e7e4941e9430f))
* improve and reference contributing documentation ([#73](https://github.com/openai/openai-go/issues/73)) ([cd4dcc1](https://github.com/openai/openai-go/commit/cd4dcc17fd18353730c74e9993700423d9df2e0d))
* **readme:** add an alpha warning ([#27](https://github.com/openai/openai-go/issues/27)) ([42ecbf8](https://github.com/openai/openai-go/commit/42ecbf8bd79d2469896e872b4c97f30cd9b203d8))
* **readme:** added some examples to readme ([#39](https://github.com/openai/openai-go/issues/39)) ([a714dde](https://github.com/openai/openai-go/commit/a714dde3784397e902208b9b47cf45a83307aa74))
* **readme:** fix example snippet ([#118](https://github.com/openai/openai-go/issues/118)) ([2af88cb](https://github.com/openai/openai-go/commit/2af88cbdf2fa07f59d7bf929f3df4aac35a579b0))
* **readme:** fix misplaced period ([#156](https://github.com/openai/openai-go/issues/156)) ([42bbc45](https://github.com/openai/openai-go/commit/42bbc454ba50b88844cac58abfd254f889b56790))
* **readme:** fix typo ([#148](https://github.com/openai/openai-go/issues/148)) ([07d3e40](https://github.com/openai/openai-go/commit/07d3e40f2e5f605b10818b97b580ebd46589b9b2))
* **readme:** smaller readme snippets with links to examples ([#46](https://github.com/openai/openai-go/issues/46)) ([082e6ae](https://github.com/openai/openai-go/commit/082e6aedc9b3a80402511c062f3331c90be6d4d0))
* update CONTRIBUTING.md ([#51](https://github.com/openai/openai-go/issues/51)) ([871b758](https://github.com/openai/openai-go/commit/871b7580fe8da75cc4c2525612a094006cc8bcf8))
* update URLs from stainlessapi.com to stainless.com ([#243](https://github.com/openai/openai-go/issues/243)) ([fcb72c5](https://github.com/openai/openai-go/commit/fcb72c51429536374d2017e4303c7d90ded7a7b8))


### Refactors

* sort fields for squashed union structs ([#111](https://github.com/openai/openai-go/issues/111)) ([e927fb0](https://github.com/openai/openai-go/commit/e927fb06a76a9ad5e6061a24cf78ee2f43b9cb26))

## 0.1.0-alpha.63 (2025-03-12)

Full Changelog: [v0.1.0-alpha.62...v0.1.0-alpha.63](https://github.com/openai/openai-go/compare/v0.1.0-alpha.62...v0.1.0-alpha.63)

### Features

* add SKIP_BREW env var to ./scripts/bootstrap ([#255](https://github.com/openai/openai-go/issues/255)) ([175ff8a](https://github.com/openai/openai-go/commit/175ff8a9fd945152693a873b20beaa1c4a4b0fc7))
* **api:** updates ([#259](https://github.com/openai/openai-go/issues/259)) ([b95257b](https://github.com/openai/openai-go/commit/b95257b2950d77f5a634aa08a82925300c318031))
* **client:** accept RFC6838 JSON content types ([#256](https://github.com/openai/openai-go/issues/256)) ([5941c83](https://github.com/openai/openai-go/commit/5941c838786018b77fdd8a72eaef5e798b95f81a))
* **client:** allow custom baseurls without trailing slash ([#254](https://github.com/openai/openai-go/issues/254)) ([de1216a](https://github.com/openai/openai-go/commit/de1216a51301e6aeb987bdbff707d2f85b374593))


### Bug Fixes

* **responses:** correct computer use enum value ([#261](https://github.com/openai/openai-go/issues/261)) ([e5a07c6](https://github.com/openai/openai-go/commit/e5a07c6dbd762d567de2049bb218330b480ae35e))
* **responses:** correct reasoning output type ([#262](https://github.com/openai/openai-go/issues/262)) ([ecad35d](https://github.com/openai/openai-go/commit/ecad35dce6297a7905af3f4760e46c482ac9f89f))


### Chores

* move ChatModel type to shared ([#250](https://github.com/openai/openai-go/issues/250)) ([34fbacc](https://github.com/openai/openai-go/commit/34fbacc8308683c2dca3368ad0d42cca3aee8deb))


### Refactors

* tidy up dependencies ([#257](https://github.com/openai/openai-go/issues/257)) ([d367e14](https://github.com/openai/openai-go/commit/d367e141c63b2585feb57670a9d7e762dfc7ca75))

## 0.1.0-alpha.62 (2025-03-05)

Full Changelog: [v0.1.0-alpha.61...v0.1.0-alpha.62](https://github.com/openai/openai-go/compare/v0.1.0-alpha.61...v0.1.0-alpha.62)

### Bug Fixes

* **api:** add missing file rank enum + more metadata ([#248](https://github.com/openai/openai-go/issues/248)) ([78e98d1](https://github.com/openai/openai-go/commit/78e98d18b319bc0de2c00543d75771166a42db73))

## 0.1.0-alpha.61 (2025-02-27)

Full Changelog: [v0.1.0-alpha.60...v0.1.0-alpha.61](https://github.com/openai/openai-go/compare/v0.1.0-alpha.60...v0.1.0-alpha.61)

### Documentation

* update URLs from stainlessapi.com to stainless.com ([#243](https://github.com/openai/openai-go/issues/243)) ([98019cf](https://github.com/openai/openai-go/commit/98019cf8c6b51e0b00ce8e58aca3b438f225f48f))

## 0.1.0-alpha.60 (2025-02-27)

Full Changelog: [v0.1.0-alpha.59...v0.1.0-alpha.60](https://github.com/openai/openai-go/compare/v0.1.0-alpha.59...v0.1.0-alpha.60)

### Features

* **api:** add gpt-4.5-preview ([#242](https://github.com/openai/openai-go/issues/242)) ([0a7488c](https://github.com/openai/openai-go/commit/0a7488ce04a50e513d1dfe805540ddf84914bd50))


### Chores

* **internal:** fix devcontainers setup ([#236](https://github.com/openai/openai-go/issues/236)) ([b27a9db](https://github.com/openai/openai-go/commit/b27a9db77dc86b57312831f6c0e22a1bcb4967ed))

## 0.1.0-alpha.59 (2025-02-15)

Full Changelog: [v0.1.0-alpha.58...v0.1.0-alpha.59](https://github.com/openai/openai-go/compare/v0.1.0-alpha.58...v0.1.0-alpha.59)

### Bug Fixes

* **client:** don't truncate manually specified filenames ([#230](https://github.com/openai/openai-go/issues/230)) ([853b748](https://github.com/openai/openai-go/commit/853b7483c07dc6f3b820f28bf3c5f097c3d440ad))

## 0.1.0-alpha.58 (2025-02-13)

Full Changelog: [v0.1.0-alpha.57...v0.1.0-alpha.58](https://github.com/openai/openai-go/compare/v0.1.0-alpha.57...v0.1.0-alpha.58)

### Features

* **api:** add support for storing chat completions ([#228](https://github.com/openai/openai-go/issues/228)) ([e3cb85e](https://github.com/openai/openai-go/commit/e3cb85ea8020c774557a8b2283f41f146b7efd94))

## 0.1.0-alpha.57 (2025-02-10)

Full Changelog: [v0.1.0-alpha.56...v0.1.0-alpha.57](https://github.com/openai/openai-go/compare/v0.1.0-alpha.56...v0.1.0-alpha.57)

### Bug Fixes

* do not call path.Base on ContentType ([#225](https://github.com/openai/openai-go/issues/225)) ([7dda9a8](https://github.com/openai/openai-go/commit/7dda9a8792c5d0725614d8867d42db45cb408a91))

## 0.1.0-alpha.56 (2025-02-07)

Full Changelog: [v0.1.0-alpha.55...v0.1.0-alpha.56](https://github.com/openai/openai-go/compare/v0.1.0-alpha.55...v0.1.0-alpha.56)

### Features

* **pagination:** avoid fetching when has_more: false ([#218](https://github.com/openai/openai-go/issues/218)) ([22dfd12](https://github.com/openai/openai-go/commit/22dfd12bd06caaa6750b54c934b970288f43c67f))


### Bug Fixes

* fix early cancel when RequestTimeout is provided for streaming requests ([#221](https://github.com/openai/openai-go/issues/221)) ([4843296](https://github.com/openai/openai-go/commit/48432968fbb130a60d2a678df4e19bcc41d4a5e1))

## 0.1.0-alpha.55 (2025-02-05)

Full Changelog: [v0.1.0-alpha.54...v0.1.0-alpha.55](https://github.com/openai/openai-go/compare/v0.1.0-alpha.54...v0.1.0-alpha.55)

### Bug Fixes

* **api:** add missing reasoning effort + model enums ([#215](https://github.com/openai/openai-go/issues/215)) ([a2345e6](https://github.com/openai/openai-go/commit/a2345e67316c4571354ab8587a94aa477fadf81b))
* update stream error handling ([#213](https://github.com/openai/openai-go/issues/213)) ([2f82244](https://github.com/openai/openai-go/commit/2f82244b2c6c8584c4e8a91db8863a4cd0b41a2b))

## 0.1.0-alpha.54 (2025-02-05)

Full Changelog: [v0.1.0-alpha.53...v0.1.0-alpha.54](https://github.com/openai/openai-go/compare/v0.1.0-alpha.53...v0.1.0-alpha.54)

### Bug Fixes

* **streaming:** correctly decode assistant events ([38ded46](https://github.com/openai/openai-go/commit/38ded4694480071a768eb4d2790ba4552e001506))


### Chores

* add UnionUnmarshaler for responses that are interfaces ([#211](https://github.com/openai/openai-go/issues/211)) ([185d848](https://github.com/openai/openai-go/commit/185d848cd3e9efb48f7468acc06a8b78f3a7b785))

## 0.1.0-alpha.53 (2025-02-05)

Full Changelog: [v0.1.0-alpha.52...v0.1.0-alpha.53](https://github.com/openai/openai-go/compare/v0.1.0-alpha.52...v0.1.0-alpha.53)

### Bug Fixes

* **api/types:** correct audio duration & role types ([#209](https://github.com/openai/openai-go/issues/209)) ([bb8cc1a](https://github.com/openai/openai-go/commit/bb8cc1a938ba142068261170d7c82a445c2f0c6c))


### Chores

* **api:** delete deprecated method ([#208](https://github.com/openai/openai-go/issues/208)) ([0a927ba](https://github.com/openai/openai-go/commit/0a927ba16dfc5cde2a368e8a1040f5ba3cda7708))
* **docs:** add docstring explaining streaming pattern ([#205](https://github.com/openai/openai-go/issues/205)) ([0bdb37f](https://github.com/openai/openai-go/commit/0bdb37f7efd9b338c32a1a83b90cecdb74f8ecce))


### Documentation

* **examples:** fix typo ([#207](https://github.com/openai/openai-go/issues/207)) ([05796de](https://github.com/openai/openai-go/commit/05796de39d3c1c88251fc42674bec1a53730c3d2))

## 0.1.0-alpha.52 (2025-02-03)

Full Changelog: [v0.1.0-alpha.51...v0.1.0-alpha.52](https://github.com/openai/openai-go/compare/v0.1.0-alpha.51...v0.1.0-alpha.52)

### Features

* **client:** send `X-Stainless-Timeout` header ([#204](https://github.com/openai/openai-go/issues/204)) ([4ccf9c9](https://github.com/openai/openai-go/commit/4ccf9c9be331773407ec062c9da563d5b20831fe))


### Bug Fixes

* **stream:** ensure .Close() doesn't panic ([#194](https://github.com/openai/openai-go/issues/194)) ([02983a3](https://github.com/openai/openai-go/commit/02983a322b264af40105b7c742e6fe24cbb396d3))
* **stream:** ensure .Close() doesn't panic ([#201](https://github.com/openai/openai-go/issues/201)) ([2df52a9](https://github.com/openai/openai-go/commit/2df52a9c7ec4839f104c0f30edfda183693fce8c))


### Documentation

* document raw responses ([#197](https://github.com/openai/openai-go/issues/197)) ([8400879](https://github.com/openai/openai-go/commit/8400879b1c226a62dc95feedc97a9c2718c8210f))

## 0.1.0-alpha.51 (2025-01-31)

Full Changelog: [v0.1.0-alpha.50...v0.1.0-alpha.51](https://github.com/openai/openai-go/compare/v0.1.0-alpha.50...v0.1.0-alpha.51)

### Features

* **api:** add o3-mini ([#195](https://github.com/openai/openai-go/issues/195)) ([c5689d0](https://github.com/openai/openai-go/commit/c5689d01773a7ac0c10e95c1c8badadde251924e))


### Bug Fixes

* fix unicode encoding for json ([#193](https://github.com/openai/openai-go/issues/193)) ([3bd3c60](https://github.com/openai/openai-go/commit/3bd3c60561b1213dc4a5fe8b27c30ddef8234726))
* **types:** correct metadata type + other fixes ([c5689d0](https://github.com/openai/openai-go/commit/c5689d01773a7ac0c10e95c1c8badadde251924e))

## 0.1.0-alpha.50 (2025-01-29)

Full Changelog: [v0.1.0-alpha.49...v0.1.0-alpha.50](https://github.com/openai/openai-go/compare/v0.1.0-alpha.49...v0.1.0-alpha.50)

### Chores

* refactor client tests ([#187](https://github.com/openai/openai-go/issues/187)) ([2752c07](https://github.com/openai/openai-go/commit/2752c076dcc87e7c08d6785ba0a00ec9eba5a0c1))

## 0.1.0-alpha.49 (2025-01-22)

Full Changelog: [v0.1.0-alpha.48...v0.1.0-alpha.49](https://github.com/openai/openai-go/compare/v0.1.0-alpha.48...v0.1.0-alpha.49)

### Features

* **api:** update enum values, comments, and examples ([#181](https://github.com/openai/openai-go/issues/181)) ([a074981](https://github.com/openai/openai-go/commit/a07498136304d74bd706684341fb9dcce6e8075c))
* Minor text change: Update readme to say beta instead of alpha ([2b766ab](https://github.com/openai/openai-go/commit/2b766ab7054cb649d40db0e7ac50c370e070043f))
* support deprecated markers ([#178](https://github.com/openai/openai-go/issues/178)) ([3d6f52f](https://github.com/openai/openai-go/commit/3d6f52f0f5b30f1f064cfbba8f61d28c1094bb0a))

## 0.1.0-alpha.48 (2025-01-21)

Full Changelog: [v0.1.0-alpha.47...v0.1.0-alpha.48](https://github.com/openai/openai-go/compare/v0.1.0-alpha.47...v0.1.0-alpha.48)

### Bug Fixes

* fix apijson.Port for embedded structs ([#174](https://github.com/openai/openai-go/issues/174)) ([b9bc4bf](https://github.com/openai/openai-go/commit/b9bc4bf94438100057c9c95a199e82ec6a48e12e))
* fix apijson.Port for embedded structs ([#177](https://github.com/openai/openai-go/issues/177)) ([a85df33](https://github.com/openai/openai-go/commit/a85df33a9fe89dffb5ce00ee297d173bb40018ed))


### Chores

* **internal:** rename `streaming.go` ([#176](https://github.com/openai/openai-go/issues/176)) ([8c54a3b](https://github.com/openai/openai-go/commit/8c54a3bfe8ed07346c2eade2cabfff1f6d97a7d8))

## 0.1.0-alpha.47 (2025-01-20)

Full Changelog: [v0.1.0-alpha.46...v0.1.0-alpha.47](https://github.com/openai/openai-go/compare/v0.1.0-alpha.46...v0.1.0-alpha.47)

### Bug Fixes

* flush stream response when done event is sent ([#172](https://github.com/openai/openai-go/issues/172)) ([fa793de](https://github.com/openai/openai-go/commit/fa793de4e849c5e3ba23fbe4d6fd1533f08d9fe6))

## 0.1.0-alpha.46 (2025-01-17)

Full Changelog: [v0.1.0-alpha.45...v0.1.0-alpha.46](https://github.com/openai/openai-go/compare/v0.1.0-alpha.45...v0.1.0-alpha.46)

### Chores

* **internal:** streaming refactors ([#165](https://github.com/openai/openai-go/issues/165)) ([168a030](https://github.com/openai/openai-go/commit/168a0305dc423cb09cefbfe0ad7c9c240a6a41b8))
* **types:** rename vector store chunking strategy ([#169](https://github.com/openai/openai-go/issues/169)) ([2dc79a6](https://github.com/openai/openai-go/commit/2dc79a62b58e78275a68ee19d85b9be7873a3491))

## 0.1.0-alpha.45 (2025-01-09)

Full Changelog: [v0.1.0-alpha.44...v0.1.0-alpha.45](https://github.com/openai/openai-go/compare/v0.1.0-alpha.44...v0.1.0-alpha.45)

### Chores

* **internal:** spec update ([#158](https://github.com/openai/openai-go/issues/158)) ([0ca04bd](https://github.com/openai/openai-go/commit/0ca04bdb233cf8fc5ea9ff497363a466bb95759a))

## 0.1.0-alpha.44 (2025-01-08)

Full Changelog: [v0.1.0-alpha.43...v0.1.0-alpha.44](https://github.com/openai/openai-go/compare/v0.1.0-alpha.43...v0.1.0-alpha.44)

### Documentation

* **readme:** fix misplaced period ([#156](https://github.com/openai/openai-go/issues/156)) ([438fe84](https://github.com/openai/openai-go/commit/438fe84825a5d73a649d72c64ba52bf2f59e8847))

## 0.1.0-alpha.43 (2025-01-03)

Full Changelog: [v0.1.0-alpha.42...v0.1.0-alpha.43](https://github.com/openai/openai-go/compare/v0.1.0-alpha.42...v0.1.0-alpha.43)

### Chores

* **api:** bump spec version ([#154](https://github.com/openai/openai-go/issues/154)) ([e55cffe](https://github.com/openai/openai-go/commit/e55cffe979c1e573ca2fe713365f31c62bb00b0e))

## 0.1.0-alpha.42 (2025-01-02)

Full Changelog: [v0.1.0-alpha.41...v0.1.0-alpha.42](https://github.com/openai/openai-go/compare/v0.1.0-alpha.41...v0.1.0-alpha.42)

### Chores

* bump license year ([#151](https://github.com/openai/openai-go/issues/151)) ([c3e375e](https://github.com/openai/openai-go/commit/c3e375eca1eaece3b73b6918e5c07474717decc5))

## 0.1.0-alpha.41 (2024-12-21)

Full Changelog: [v0.1.0-alpha.40...v0.1.0-alpha.41](https://github.com/openai/openai-go/compare/v0.1.0-alpha.40...v0.1.0-alpha.41)

### Chores

* **internal:** spec update ([#146](https://github.com/openai/openai-go/issues/146)) ([a9a67ba](https://github.com/openai/openai-go/commit/a9a67ba2f487d23adc7238c1151cab76a7fd9f01))


### Documentation

* **readme:** fix typo ([#148](https://github.com/openai/openai-go/issues/148)) ([1bb2322](https://github.com/openai/openai-go/commit/1bb23222e794c8439a56f765dee0c8d4e0db6eac))

## 0.1.0-alpha.40 (2024-12-17)

Full Changelog: [v0.1.0-alpha.39...v0.1.0-alpha.40](https://github.com/openai/openai-go/compare/v0.1.0-alpha.39...v0.1.0-alpha.40)

### Features

* **api:** new o1 and GPT-4o models + preference fine-tuning ([#142](https://github.com/openai/openai-go/issues/142)) ([9207561](https://github.com/openai/openai-go/commit/920756111c0f7725eb5746b22c1393e327dc2990))


### Chores

* **internal:** spec update ([#145](https://github.com/openai/openai-go/issues/145)) ([bc0dba4](https://github.com/openai/openai-go/commit/bc0dba40f69294f460b3d448a6a093e487509dca))

## 0.1.0-alpha.39 (2024-12-05)

Full Changelog: [v0.1.0-alpha.38...v0.1.0-alpha.39](https://github.com/openai/openai-go/compare/v0.1.0-alpha.38...v0.1.0-alpha.39)

### Features

* **api:** updates ([#138](https://github.com/openai/openai-go/issues/138)) ([67badff](https://github.com/openai/openai-go/commit/67badffd9f6b24979e52a71548437516b8d97e6f))


### Chores

* bump openapi url ([#136](https://github.com/openai/openai-go/issues/136)) ([179a1cd](https://github.com/openai/openai-go/commit/179a1cd7104d13c2014fa9fe3b1551e4b8829d71))

## 0.1.0-alpha.38 (2024-11-20)

Full Changelog: [v0.1.0-alpha.37...v0.1.0-alpha.38](https://github.com/openai/openai-go/compare/v0.1.0-alpha.37...v0.1.0-alpha.38)

### Features

* **api:** add gpt-4o-2024-11-20 model ([#131](https://github.com/openai/openai-go/issues/131)) ([8dabbd3](https://github.com/openai/openai-go/commit/8dabbd3d3f338986470049c1a3842e47a7194c1e))


### Chores

* **internal:** spec update ([#130](https://github.com/openai/openai-go/issues/130)) ([23476d8](https://github.com/openai/openai-go/commit/23476d8131b009e38f7f22024d925fd7fe76e4db))
* **tests:** limit array example length ([#128](https://github.com/openai/openai-go/issues/128)) ([5560e6b](https://github.com/openai/openai-go/commit/5560e6ba5f77e774d6be40e71628c1b62f8b3005))

## 0.1.0-alpha.37 (2024-11-12)

Full Changelog: [v0.1.0-alpha.36...v0.1.0-alpha.37](https://github.com/openai/openai-go/compare/v0.1.0-alpha.36...v0.1.0-alpha.37)

### Bug Fixes

* **client:** no panic on missing BaseURL ([#121](https://github.com/openai/openai-go/issues/121)) ([1a8b841](https://github.com/openai/openai-go/commit/1a8b8415caabf9e4f33fc2b095b31dd926376b47))

## 0.1.0-alpha.36 (2024-11-11)

Full Changelog: [v0.1.0-alpha.35...v0.1.0-alpha.36](https://github.com/openai/openai-go/compare/v0.1.0-alpha.35...v0.1.0-alpha.36)

### Bug Fixes

* correct required fields for flattened unions ([#120](https://github.com/openai/openai-go/issues/120)) ([8fe865b](https://github.com/openai/openai-go/commit/8fe865b5cb230bdac498ae295b16a7388cfdde6c))


### Documentation

* **readme:** fix example snippet ([#118](https://github.com/openai/openai-go/issues/118)) ([7f1803b](https://github.com/openai/openai-go/commit/7f1803b44183e8796a40eab6dee440cf16813e3c))

## 0.1.0-alpha.35 (2024-11-10)

Full Changelog: [v0.1.0-alpha.34...v0.1.0-alpha.35](https://github.com/openai/openai-go/compare/v0.1.0-alpha.34...v0.1.0-alpha.35)

### Bug Fixes

* **api:** escape key values when encoding maps ([#116](https://github.com/openai/openai-go/issues/116)) ([a2bcd73](https://github.com/openai/openai-go/commit/a2bcd7394f725ea0e653d7d4c145f3f48b36c1c3))

## 0.1.0-alpha.34 (2024-11-07)

Full Changelog: [v0.1.0-alpha.33...v0.1.0-alpha.34](https://github.com/openai/openai-go/compare/v0.1.0-alpha.33...v0.1.0-alpha.34)

### Documentation

* add missing docs for some enums ([#114](https://github.com/openai/openai-go/issues/114)) ([f01913f](https://github.com/openai/openai-go/commit/f01913f1432a64304de9232bd36624c9506e02ab))

## 0.1.0-alpha.33 (2024-11-05)

Full Changelog: [v0.1.0-alpha.32...v0.1.0-alpha.33](https://github.com/openai/openai-go/compare/v0.1.0-alpha.32...v0.1.0-alpha.33)

### Features

* **api:** add support for predicted outputs ([#110](https://github.com/openai/openai-go/issues/110)) ([ab88fa9](https://github.com/openai/openai-go/commit/ab88fa960917bedd15d2ffdf50d0d7169afd661a))


### Refactors

* sort fields for squashed union structs ([#111](https://github.com/openai/openai-go/issues/111)) ([f7e4ac8](https://github.com/openai/openai-go/commit/f7e4ac83cd345e58c824a489a2883bd4ef7717f6))

## 0.1.0-alpha.32 (2024-10-30)

Full Changelog: [v0.1.0-alpha.31...v0.1.0-alpha.32](https://github.com/openai/openai-go/compare/v0.1.0-alpha.31...v0.1.0-alpha.32)

### Features

* **api:** add new, expressive voices for Realtime and Audio in Chat Completions ([#101](https://github.com/openai/openai-go/issues/101)) ([f946acc](https://github.com/openai/openai-go/commit/f946acc71a92f885bed87f0d4e724fb40cae0f14))

## 0.1.0-alpha.31 (2024-10-23)

Full Changelog: [v0.1.0-alpha.30...v0.1.0-alpha.31](https://github.com/openai/openai-go/compare/v0.1.0-alpha.30...v0.1.0-alpha.31)

### Chores

* **internal:** update spec version ([#95](https://github.com/openai/openai-go/issues/95)) ([0cb6f6a](https://github.com/openai/openai-go/commit/0cb6f6abd428a5bd314902708ab12bc12a1b978f))

## 0.1.0-alpha.30 (2024-10-22)

Full Changelog: [v0.1.0-alpha.29...v0.1.0-alpha.30](https://github.com/openai/openai-go/compare/v0.1.0-alpha.29...v0.1.0-alpha.30)

### ⚠ BREAKING CHANGES

* **client:** improve naming of some variants ([#89](https://github.com/openai/openai-go/issues/89))

### Features

* **client:** improve naming of some variants ([#89](https://github.com/openai/openai-go/issues/89)) ([12ac070](https://github.com/openai/openai-go/commit/12ac070611061e98ae2aaaeefa8eb661ff7f995f))

## 0.1.0-alpha.29 (2024-10-17)

Full Changelog: [v0.1.0-alpha.28...v0.1.0-alpha.29](https://github.com/openai/openai-go/compare/v0.1.0-alpha.28...v0.1.0-alpha.29)

### Features

* **api:** add gpt-4o-audio-preview model for chat completions ([#88](https://github.com/openai/openai-go/issues/88)) ([03da9c9](https://github.com/openai/openai-go/commit/03da9c984e6b0fdb0d8da7a8e4fde29fc45b784d))
* **assistants:** add polling helpers and examples ([#84](https://github.com/openai/openai-go/issues/84)) ([eab25dd](https://github.com/openai/openai-go/commit/eab25dde95f7dd50b712326714ce55e93432b4dc))


### Bug Fixes

* **example:** use correct model ([#86](https://github.com/openai/openai-go/issues/86)) ([6dad9b2](https://github.com/openai/openai-go/commit/6dad9b256b86b069ec9445ae86bfb2f2c3764b66))

## 0.1.0-alpha.28 (2024-10-16)

Full Changelog: [v0.1.0-alpha.27...v0.1.0-alpha.28](https://github.com/openai/openai-go/compare/v0.1.0-alpha.27...v0.1.0-alpha.28)

### Features

* move pagination package from internal to packages ([#81](https://github.com/openai/openai-go/issues/81)) ([8875bdc](https://github.com/openai/openai-go/commit/8875bdc847467b322bd9b6c54c027d97a79c5f16))

## 0.1.0-alpha.27 (2024-10-14)

Full Changelog: [v0.1.0-alpha.26...v0.1.0-alpha.27](https://github.com/openai/openai-go/compare/v0.1.0-alpha.26...v0.1.0-alpha.27)

### Chores

* fix GetNextPage docstring ([#78](https://github.com/openai/openai-go/issues/78)) ([490f8f0](https://github.com/openai/openai-go/commit/490f8f0ae34cc6769a7555cf77fef0192963ad06))

## 0.1.0-alpha.26 (2024-10-08)

Full Changelog: [v0.1.0-alpha.25...v0.1.0-alpha.26](https://github.com/openai/openai-go/compare/v0.1.0-alpha.25...v0.1.0-alpha.26)

### Bug Fixes

* **beta:** pass beta header by default ([#75](https://github.com/openai/openai-go/issues/75)) ([cb66b47](https://github.com/openai/openai-go/commit/cb66b474fb86646501314456fce6acc4b31a2026))

## 0.1.0-alpha.25 (2024-10-02)

Full Changelog: [v0.1.0-alpha.24...v0.1.0-alpha.25](https://github.com/openai/openai-go/compare/v0.1.0-alpha.24...v0.1.0-alpha.25)

### Documentation

* improve and reference contributing documentation ([#73](https://github.com/openai/openai-go/issues/73)) ([03a8261](https://github.com/openai/openai-go/commit/03a8261970011b2be7e101ec095a0eb93b361a04))

## 0.1.0-alpha.24 (2024-10-01)

Full Changelog: [v0.1.0-alpha.23...v0.1.0-alpha.24](https://github.com/openai/openai-go/compare/v0.1.0-alpha.23...v0.1.0-alpha.24)

### Features

* **api:** support storing chat completions, enabling evals and model distillation in the dashboard ([#72](https://github.com/openai/openai-go/issues/72)) ([1e50f54](https://github.com/openai/openai-go/commit/1e50f549ef135d7494c9260c4638c6054fe06c74))


### Chores

* **docs:** fix maxium typo ([#69](https://github.com/openai/openai-go/issues/69)) ([3a5c6a6](https://github.com/openai/openai-go/commit/3a5c6a657ac8d821e95e07b442f00140b5332c93))

## 0.1.0-alpha.23 (2024-09-29)

Full Changelog: [v0.1.0-alpha.22...v0.1.0-alpha.23](https://github.com/openai/openai-go/compare/v0.1.0-alpha.22...v0.1.0-alpha.23)

### Chores

* **docs:** remove some duplicative api.md entries ([#65](https://github.com/openai/openai-go/issues/65)) ([13a1ca2](https://github.com/openai/openai-go/commit/13a1ca2eb6320c797d6e278bfe258e1e7f27e031))

## 0.1.0-alpha.22 (2024-09-26)

Full Changelog: [v0.1.0-alpha.21...v0.1.0-alpha.22](https://github.com/openai/openai-go/compare/v0.1.0-alpha.21...v0.1.0-alpha.22)

### Features

* **api:** add omni-moderation model ([#63](https://github.com/openai/openai-go/issues/63)) ([9ca9ebb](https://github.com/openai/openai-go/commit/9ca9ebb1f40c056642d987445ea0cc8d60a1d15f))

## 0.1.0-alpha.21 (2024-09-25)

Full Changelog: [v0.1.0-alpha.20...v0.1.0-alpha.21](https://github.com/openai/openai-go/compare/v0.1.0-alpha.20...v0.1.0-alpha.21)

### Features

* **client:** send retry count header ([#60](https://github.com/openai/openai-go/issues/60)) ([8797500](https://github.com/openai/openai-go/commit/87975004c4917be6b59c34b1252b6a393412a754))


### Bug Fixes

* **audio:** correct response_format translations type ([#62](https://github.com/openai/openai-go/issues/62)) ([4b8df65](https://github.com/openai/openai-go/commit/4b8df6595d2d416c3589d6a270ebdf247bbe18af))

## 0.1.0-alpha.20 (2024-09-20)

Full Changelog: [v0.1.0-alpha.19...v0.1.0-alpha.20](https://github.com/openai/openai-go/compare/v0.1.0-alpha.19...v0.1.0-alpha.20)

### Chores

* **types:** improve type name for embedding models ([#57](https://github.com/openai/openai-go/issues/57)) ([05fe24e](https://github.com/openai/openai-go/commit/05fe24eec2797848bbe866ad3c4bfa8da4a61b77))

## 0.1.0-alpha.19 (2024-09-18)

Full Changelog: [v0.1.0-alpha.18...v0.1.0-alpha.19](https://github.com/openai/openai-go/compare/v0.1.0-alpha.18...v0.1.0-alpha.19)

### Features

* fix(streaming): correctly accumulate tool calls and roles ([#55](https://github.com/openai/openai-go/issues/55)) ([89651e4](https://github.com/openai/openai-go/commit/89651e4ebb80179b2fcc92d3c573679683a39201))

## 0.1.0-alpha.18 (2024-09-16)

Full Changelog: [v0.1.0-alpha.17...v0.1.0-alpha.18](https://github.com/openai/openai-go/compare/v0.1.0-alpha.17...v0.1.0-alpha.18)

### Chores

* **internal:** update spec link ([#53](https://github.com/openai/openai-go/issues/53)) ([0fefed1](https://github.com/openai/openai-go/commit/0fefed1b392ea99ce2fa68e22a9ee53f60476037))


### Documentation

* update CONTRIBUTING.md ([#51](https://github.com/openai/openai-go/issues/51)) ([fe2d656](https://github.com/openai/openai-go/commit/fe2d656eaee480a87c4bf00eb5937fb167018ec9))

## 0.1.0-alpha.17 (2024-09-12)

Full Changelog: [v0.1.0-alpha.16...v0.1.0-alpha.17](https://github.com/openai/openai-go/compare/v0.1.0-alpha.16...v0.1.0-alpha.17)

### Features

* **api:** add o1 models ([#49](https://github.com/openai/openai-go/issues/49)) ([37d160c](https://github.com/openai/openai-go/commit/37d160cef58d3aca3f8dfc8c50b0eb8b516c1bcb))


### Documentation

* **readme:** smaller readme snippets with links to examples ([#46](https://github.com/openai/openai-go/issues/46)) ([dcea342](https://github.com/openai/openai-go/commit/dcea34213655ce8f9d84979d2f3d9dfa1f7459a3))

## 0.1.0-alpha.16 (2024-09-10)

Full Changelog: [v0.1.0-alpha.15...v0.1.0-alpha.16](https://github.com/openai/openai-go/compare/v0.1.0-alpha.15...v0.1.0-alpha.16)

### Bug Fixes

* **requestconfig:** copy over more fields when cloning ([#44](https://github.com/openai/openai-go/issues/44)) ([6e02130](https://github.com/openai/openai-go/commit/6e02130c086c21e7f0895d18d6ed98fefb56f4d0))

## 0.1.0-alpha.15 (2024-09-05)

Full Changelog: [v0.1.0-alpha.14...v0.1.0-alpha.15](https://github.com/openai/openai-go/compare/v0.1.0-alpha.14...v0.1.0-alpha.15)

### Features

* **vector store:** improve chunking strategy type names ([#40](https://github.com/openai/openai-go/issues/40)) ([4932cca](https://github.com/openai/openai-go/commit/4932ccac47b4b7976366244aab5810fa44292350))

## 0.1.0-alpha.14 (2024-09-03)

Full Changelog: [v0.1.0-alpha.13...v0.1.0-alpha.14](https://github.com/openai/openai-go/compare/v0.1.0-alpha.13...v0.1.0-alpha.14)

### Features

* **examples/structure-outputs:** created an example for using structured outputs ([7d1e71e](https://github.com/openai/openai-go/commit/7d1e71e72b8c55d5b7228b72d967e4cae8165280))
* **stream-accumulators:** added streaming accumulator helpers and example ([29e80e7](https://github.com/openai/openai-go/commit/29e80e7dfb4571e93e616981ddc950e3058b6203))


### Bug Fixes

* **examples/fine-tuning:** used an old constant name ([#34](https://github.com/openai/openai-go/issues/34)) ([5d9ec26](https://github.com/openai/openai-go/commit/5d9ec26407b15c7effceb999bba3dfbeefc0adf2))


### Documentation

* **readme:** added some examples to readme ([#39](https://github.com/openai/openai-go/issues/39)) ([2dbfa62](https://github.com/openai/openai-go/commit/2dbfa62ffc89ead88e0fed586684a6b757836752))

## 0.1.0-alpha.13 (2024-08-29)

Full Changelog: [v0.1.0-alpha.12...v0.1.0-alpha.13](https://github.com/openai/openai-go/compare/v0.1.0-alpha.12...v0.1.0-alpha.13)

### Features

* **api:** add file search result details to run steps ([#32](https://github.com/openai/openai-go/issues/32)) ([f6a1f12](https://github.com/openai/openai-go/commit/f6a1f12acbaf158af8009debcc2019d1b9e19104))

## 0.1.0-alpha.12 (2024-08-23)

Full Changelog: [v0.1.0-alpha.11...v0.1.0-alpha.12](https://github.com/openai/openai-go/compare/v0.1.0-alpha.11...v0.1.0-alpha.12)

### Features

* support assistants stream ([0647f03](https://github.com/openai/openai-go/commit/0647f03c55fe8ec654f6a8fd98d77384d9df6b9d))

## 0.1.0-alpha.11 (2024-08-22)

Full Changelog: [v0.1.0-alpha.10...v0.1.0-alpha.11](https://github.com/openai/openai-go/compare/v0.1.0-alpha.10...v0.1.0-alpha.11)

### Features

* add support for error property in stream ([#29](https://github.com/openai/openai-go/issues/29)) ([73f9342](https://github.com/openai/openai-go/commit/73f93429e1319387f1a95208166b3e871ce4e03a))

## 0.1.0-alpha.10 (2024-08-21)

Full Changelog: [v0.1.0-alpha.9...v0.1.0-alpha.10](https://github.com/openai/openai-go/compare/v0.1.0-alpha.9...v0.1.0-alpha.10)

### Documentation

* **readme:** add an alpha warning ([#27](https://github.com/openai/openai-go/issues/27)) ([3f1cc3b](https://github.com/openai/openai-go/commit/3f1cc3bbf19daa48e83aacb6906b9776726d7154))

## 0.1.0-alpha.9 (2024-08-16)

Full Changelog: [v0.1.0-alpha.8...v0.1.0-alpha.9](https://github.com/openai/openai-go/compare/v0.1.0-alpha.8...v0.1.0-alpha.9)

### Features

* **api:** add chatgpt-4o-latest model ([#24](https://github.com/openai/openai-go/issues/24)) ([112c7f3](https://github.com/openai/openai-go/commit/112c7f31917596b6c029a1f00643647375e8c8c8))

## 0.1.0-alpha.8 (2024-08-15)

Full Changelog: [v0.1.0-alpha.7...v0.1.0-alpha.8](https://github.com/openai/openai-go/compare/v0.1.0-alpha.7...v0.1.0-alpha.8)

### Chores

* **types:** define FilePurpose enum ([#22](https://github.com/openai/openai-go/issues/22)) ([2a7c699](https://github.com/openai/openai-go/commit/2a7c699e4fb21f848aa5d260da9d2a5c471866d1))

## 0.1.0-alpha.7 (2024-08-12)

Full Changelog: [v0.1.0-alpha.6...v0.1.0-alpha.7](https://github.com/openai/openai-go/compare/v0.1.0-alpha.6...v0.1.0-alpha.7)

### Chores

* add back custom code that was reverted ([2557bd8](https://github.com/openai/openai-go/commit/2557bd8b5f1748adf67d9208ceaeea3250d93b14))

## 0.1.0-alpha.6 (2024-08-12)

Full Changelog: [v0.1.0-alpha.5...v0.1.0-alpha.6](https://github.com/openai/openai-go/compare/v0.1.0-alpha.5...v0.1.0-alpha.6)

### Features

* Adding in Azure support ([3225b7c](https://github.com/openai/openai-go/commit/3225b7c6028c0c5ab9420416b6bb8b31a5383218))

## 0.1.0-alpha.5 (2024-08-12)

Full Changelog: [v0.1.0-alpha.4...v0.1.0-alpha.5](https://github.com/openai/openai-go/compare/v0.1.0-alpha.4...v0.1.0-alpha.5)

### Features

* simplify content union ([#18](https://github.com/openai/openai-go/issues/18)) ([51877bf](https://github.com/openai/openai-go/commit/51877bf8f16e348a531aa54f0f49e9d71390a485))

## 0.1.0-alpha.4 (2024-08-12)

Full Changelog: [v0.1.0-alpha.3...v0.1.0-alpha.4](https://github.com/openai/openai-go/compare/v0.1.0-alpha.3...v0.1.0-alpha.4)

### Chores

* **examples:** minor formatting changes ([#14](https://github.com/openai/openai-go/issues/14)) ([8d4490b](https://github.com/openai/openai-go/commit/8d4490b78dcc0edee3264448e3fa3f3781d04258))

## 0.1.0-alpha.3 (2024-08-12)

Full Changelog: [v0.1.0-alpha.2...v0.1.0-alpha.3](https://github.com/openai/openai-go/compare/v0.1.0-alpha.2...v0.1.0-alpha.3)

### Chores

* bump Go to v1.21 ([#12](https://github.com/openai/openai-go/issues/12)) ([db5efda](https://github.com/openai/openai-go/commit/db5efdaad3848b8f130f279e6760d9d525e02bda))

## 0.1.0-alpha.2 (2024-08-10)

Full Changelog: [v0.1.0-alpha.1...v0.1.0-alpha.2](https://github.com/openai/openai-go/compare/v0.1.0-alpha.1...v0.1.0-alpha.2)

### Bug Fixes

* deserialization of struct unions that implement json.Unmarshaler ([#11](https://github.com/openai/openai-go/issues/11)) ([7c0847a](https://github.com/openai/openai-go/commit/7c0847aa2ae15b4442ab0625d8a780ed684c275e))


### Chores

* **ci:** bump prism mock server version ([#10](https://github.com/openai/openai-go/issues/10)) ([00f9455](https://github.com/openai/openai-go/commit/00f9455692c52fb37544d3f657090b216667d8ec))
* **ci:** codeowners file ([#9](https://github.com/openai/openai-go/issues/9)) ([be41ac2](https://github.com/openai/openai-go/commit/be41ac2ce87efacf17748cb9dd2d3b1b4a43180e))
* **internal:** updates ([#6](https://github.com/openai/openai-go/issues/6)) ([316e623](https://github.com/openai/openai-go/commit/316e6231c27728f4031f822287389c67e914739a))

## 0.1.0-alpha.1 (2024-08-06)

Full Changelog: [v0.0.1-alpha.0...v0.1.0-alpha.1](https://github.com/openai/openai-go/compare/v0.0.1-alpha.0...v0.1.0-alpha.1)

### Features

* add azure, examples, and message constructors ([fb2df0f](https://github.com/openai/openai-go/commit/fb2df0fe22002f1826bfaa1cb008c45db375885c))
* **api:** updates ([#5](https://github.com/openai/openai-go/issues/5)) ([9f525e8](https://github.com/openai/openai-go/commit/9f525e85d8fe13cce2a18a1a48179bc5a6d1f094))
* extract out `ImageModel`, `AudioModel`, `SpeechModel` ([#3](https://github.com/openai/openai-go/issues/3)) ([f085893](https://github.com/openai/openai-go/commit/f085893d109a9e841d1df13df4c71cae06018758))
* make enums not nominal ([#4](https://github.com/openai/openai-go/issues/4)) ([9f77005](https://github.com/openai/openai-go/commit/9f77005474b8a38cbfc09f22ec3b81d1de62d3c3))
* publish ([c329601](https://github.com/openai/openai-go/commit/c329601324226e28ff18d6ccecfdde41cedd3b5a))


### Chores

* **internal:** updates ([#2](https://github.com/openai/openai-go/issues/2)) ([5976d8d](https://github.com/openai/openai-go/commit/5976d8d8b9a94cd78e4d86f704137f4b43224a08))
