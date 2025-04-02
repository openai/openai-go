# Changelog

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
