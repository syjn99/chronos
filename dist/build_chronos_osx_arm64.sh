bazel build --remote_cache=http://192.168.2.200:9090 --config=release --config=osx_arm64 //cmd/beacon-chain //cmd/validator //cmd/prysmctl //tools/enr-calculator
if [ $? -ne 0 ]; then
    echo "Bazel build failed."
    exit 1
fi
zip -j dist/chronos_osx_arm64.zip bazel-bin/cmd/beacon-chain/beacon-chain_/beacon-chain bazel-bin/cmd/validator/validator_/validator bazel-bin/tools/enr-calculator/enr-calculator_/enr-calculator bazel-bin/cmd/prysmctl/prysmctl_/prysmctl bazel-bin/tools/bootnode/bootnode_/bootnode
mkdir -p dist/bin/mac/arm64
cp bazel-bin/cmd/beacon-chain/beacon-chain_/beacon-chain dist/bin/mac/arm64/chronos
cp bazel-bin/cmd/validator/validator_/validator dist/bin/mac/arm64/validator
rm -rf bazel-bin
