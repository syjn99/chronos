bazel build --config=minimal --config=osx_arm64_docker //cmd/beacon-chain //cmd/validator //cmd/prysmctl //tools/enr-calculator
zip -j dist/chronos_osx_arm64_minimal.zip bazel-bin/cmd/beacon-chain/beacon-chain_/beacon-chain bazel-bin/cmd/validator/validator_/validator bazel-bin/tools/enr-calculator/enr-calculator_/enr-calculator bazel-bin/cmd/prysmctl/prysmctl_/prysmctl
rm -rf bazel-bin
