# Changelog

## [0.0.10](https://github.com/pascaliske/misp-operator/compare/v0.0.9...v0.0.10) (2026-08-07)


### Features

* **chart:** add note about values schema reference ([21436e6](https://github.com/pascaliske/misp-operator/commit/21436e60823b3a616a46fc81419083c9aeea0766))
* **chart:** provide values schema for intellisense ([8347b03](https://github.com/pascaliske/misp-operator/commit/8347b03c3b1c429e175a450454c0e15887c0dd27))
* **chart:** support network policies for metrics and webhooks ([1b79467](https://github.com/pascaliske/misp-operator/commit/1b79467d16d7712b94a77e5955d13301e2ad14f7))


### Bug Fixes

* **chart:** prefix webhook configurations to prevent collisions ([cf2f4b0](https://github.com/pascaliske/misp-operator/commit/cf2f4b044a9f76a95eef76d7656f5a9fbd3c05f7))

## [0.0.9](https://github.com/pascaliske/misp-operator/compare/v0.0.8...v0.0.9) (2026-08-07)


### Features

* automatically update version number in README and Chart.yaml ([66200f9](https://github.com/pascaliske/misp-operator/commit/66200f932f49094581e0770c6c59eed35fd4d552))
* **chart:** generate values documentarion using helm-docs ([aebddfb](https://github.com/pascaliske/misp-operator/commit/aebddfb373e2dbf6ab8544e656cd319f08faacb5))
* **chart:** support extra volumes and extra volume mounts ([6ccfc0f](https://github.com/pascaliske/misp-operator/commit/6ccfc0fa3545a8308015d7bf619e29b78e39ceca))
* refine misp instance watcher to exclude status only events ([c7b2bcc](https://github.com/pascaliske/misp-operator/commit/c7b2bcc0ad5baa63ae137a0d97eb048932e05450))


### Bug Fixes

* **chart:** setup webhook dependencies ([082a101](https://github.com/pascaliske/misp-operator/commit/082a101f329c05521a041641bb02a6192792fcec))

## [0.0.8](https://github.com/pascaliske/misp-operator/compare/v0.0.7...v0.0.8) (2026-08-01)


### Bug Fixes

* ensure at least one of passwordSecretRef or enableEmptyPassword is set for cache ([9f30924](https://github.com/pascaliske/misp-operator/commit/9f3092458bc3ec07a64adb3b65ee309caa4689dc))
* explicitly mark fields as required ([2df6383](https://github.com/pascaliske/misp-operator/commit/2df63839c74fd095769213f47212e8edcbfa5f0e))
* update spec and status documentation comments ([f5a4a5e](https://github.com/pascaliske/misp-operator/commit/f5a4a5e7421b9a3338b411d44962bf88e14f4da9))

## [0.0.7](https://github.com/pascaliske/misp-operator/compare/v0.0.6...v0.0.7) (2026-07-31)


### Features

* provide basic crds documentation via https://doc.crds.dev ([b001a35](https://github.com/pascaliske/misp-operator/commit/b001a3589aa0120cacb53a0b9bb7d0a1c00c5a49))
* switch update strategy of modules to rolling updates ([e174f36](https://github.com/pascaliske/misp-operator/commit/e174f36729de812981cb90aa7d4b0ce874f2bce8))


### Bug Fixes

* **deps:** update kubernetes monorepo to v0.36.3 ([e280262](https://github.com/pascaliske/misp-operator/commit/e28026248f19191cf7b8a089b98789d4822ed980))
* **deps:** update kubernetes monorepo to v0.36.3 ([d6bc3bc](https://github.com/pascaliske/misp-operator/commit/d6bc3bc7f46b2527044df2780f384140b62c552a))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.32.0 ([2e2d980](https://github.com/pascaliske/misp-operator/commit/2e2d9805cc9b479f5e9eb667aca8065c65cd7d89))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.32.0 ([c3cd749](https://github.com/pascaliske/misp-operator/commit/c3cd7498d654c0b5d667b38a6595651c860bf88b))
* **deps:** update module github.com/onsi/gomega to v1.42.1 ([49915ee](https://github.com/pascaliske/misp-operator/commit/49915ee11a2c2f0cef4302b6c0ae8b10ecf6519d))
* **deps:** update module github.com/onsi/gomega to v1.42.1 ([de5c237](https://github.com/pascaliske/misp-operator/commit/de5c2376624a230f588e31e0181d14f60e3c8465))

## [0.0.6](https://github.com/pascaliske/misp-operator/compare/v0.0.5...v0.0.6) (2026-07-20)


### Bug Fixes

* **ci:** add missing id-token permission for signing the chart ([8876adc](https://github.com/pascaliske/misp-operator/commit/8876adcc039b3b0e7b52f714693754c56b20e555))

## [0.0.5](https://github.com/pascaliske/misp-operator/compare/v0.0.4...v0.0.5) (2026-07-20)


### Bug Fixes

* **ci:** add missing cosign installer action to chart release job ([ac1150f](https://github.com/pascaliske/misp-operator/commit/ac1150fad686eb69608c9a4bad507f03636d90a1))

## [0.0.4](https://github.com/pascaliske/misp-operator/compare/v0.0.3...v0.0.4) (2026-07-20)


### Bug Fixes

* **ci:** add version tag as annotation to cosign signature ([5628a34](https://github.com/pascaliske/misp-operator/commit/5628a349c062121bf131dc5dbdfd49991c8b83ed))
* **ci:** ensure helm files are generated before packaging chart ([352ff27](https://github.com/pascaliske/misp-operator/commit/352ff27c9b4bcb15f76638c225a4f508f5550b42))
* **ci:** upload installer to draft release before finalization ([d3d5b30](https://github.com/pascaliske/misp-operator/commit/d3d5b30e7f2456b66659e45fcb0aaa3403835cef))

## [0.0.3](https://github.com/pascaliske/misp-operator/compare/v0.0.2...v0.0.3) (2026-07-20)


### Features

* enable provenance and sbom features ([21319c3](https://github.com/pascaliske/misp-operator/commit/21319c3647d330352684c695fb8f1fefeb5b7650))
* sign chart using sigstore/cosign as well ([f33c6ae](https://github.com/pascaliske/misp-operator/commit/f33c6ae92ea5d66a64e8728bb3b99ab7018e7238))
* sign image using sigstore/cosign ([4f3b865](https://github.com/pascaliske/misp-operator/commit/4f3b8658b701a5ca257e8a6c2172014e371b11d9))

## [0.0.2](https://github.com/pascaliske/misp-operator/compare/v0.0.1...v0.0.2) (2026-07-19)


### Features

* trigger initial release ([dd634f9](https://github.com/pascaliske/misp-operator/commit/dd634f9927b2673b366e77ac0e3111f36e953389))
