bazel build --config=windows_amd64 //cmd/beacon-chain //cmd/validator //cmd/prysmctl //tools/enr-calculator
zip -j dist/chronos_windows.zip bazel-bin/cmd/beacon-chain/beacon-chain_/beacon-chain.exe bazel-bin/cmd/validator/validator_/validator.exe bazel-bin/tools/enr-calculator/enr-calculator_/enr-calculator.exe bazel-bin/cmd/prysmctl/prysmctl_/prysmctl.exe
rm -rf bazel-bin
