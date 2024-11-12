# Changelog

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
