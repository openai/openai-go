# Changelog

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
