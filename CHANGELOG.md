### v2.0.0-beta1

#### Features

 - [`14fb519`](https://github.com/deis/logger/commit/14fb519650261a4ec6f46229df43190c1d621135) .travis.yml: have this job notify its sister job in Jenkins
 - [`05134a0`](https://github.com/deis/logger/commit/05134a0200e80f5f089a488803dca4133cfa8222) logger: Add healthz endpoint
 - [`707ede6`](https://github.com/deis/logger/commit/707ede62deb6f449bccbb4e1590d8aa9b9c9b1fc) makefile: - Adding a makefile that provides some convience targets for the development workflow.

#### Fixes

 - [`a70d48e`](https://github.com/deis/logger/commit/a70d48eedc902f0fd1581f88ca9a153bfeb8b973) logger: Use only unprivileged ports
 - [`f52405e`](https://github.com/deis/logger/commit/f52405e2ac7e8c31115f47ff20d1f49fe3167c72) logger: Drop bad messages instead of panicking
 - [`e2a38a6`](https://github.com/deis/logger/commit/e2a38a642c24b63737813b98fdcfdaacb676bdbc) server.go: Use labels when determining container name
 - [`b01e4c5`](https://github.com/deis/logger/commit/b01e4c553a89161aca6f4b9ddca9f75a3596e27b) deploy: When building images from master tag those builds as canary
 - [`41abb90`](https://github.com/deis/logger/commit/41abb90a963b5dedfc52581dc1353ecaff3e4767) tests: - Fix tests from refactoring to v2
 - [`46a7695`](https://github.com/deis/logger/commit/46a7695b3072bc45b270d859fd458240999a10e4) env vars: Adding and removing unnecessary env vars

#### Documentation

 - [`277f489`](https://github.com/deis/logger/commit/277f489cfdf8e95b4dbefcc2db15bf5af0db27b6) readme: Update readme with installation instructions

#### Maintenance

 - [`a16379c`](https://github.com/deis/logger/commit/a16379cfb1818207ddb605868fa68de39b4b7dfe) manifests: Put all manfiests in deis namespace
 - [`5197ce7`](https://github.com/deis/logger/commit/5197ce7fb12ac95984ba416dec211aa674e21b3c) logger: Remove drain capability
 - [`ba4427c`](https://github.com/deis/logger/commit/ba4427cd24d5d8260b7aa9deebca9e231611e9fa) logger: Upgrade to alpine 3.3
 - [`cb1f240`](https://github.com/deis/logger/commit/cb1f2406309969cd47eac8614afe7f5c642a1b21) ci: Always run docker build during CI
 - [`4141415`](https://github.com/deis/logger/commit/41414153d18101ba419129945f8f4bb725a5833a) deploy.sh: Produce git-tagged images on PR merges
 - [`f55f851`](https://github.com/deis/logger/commit/f55f85142681854c7f8ea65dcf318f5f2ce6ab7f) release.sh: Add a release script for cutting stable releases
 - [`3d02cbb`](https://github.com/deis/logger/commit/3d02cbb47e1832b5896513e0986067a28f0d7038) makefile: Update makefile to use deis instead of deisci.
 - [`64c831b`](https://github.com/deis/logger/commit/64c831bd079a79e216c6d274d305299cdc6151b0) (all): add boilerplate files
