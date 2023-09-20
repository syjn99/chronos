bazel build --config=minimal --config=linux_amd64 //cmd/beacon-chain //cmd/validator //cmd/prysmctl //tools/enr-calculator //tools/bootnode
if [ $? -ne 0 ]; then
    echo "Bazel build failed."
    exit 1
fi
zip -j dist/chronos_linux_amd64_minimal.zip bazel-bin/cmd/beacon-chain/beacon-chain_/beacon-chain bazel-bin/cmd/validator/validator_/validator bazel-bin/tools/enr-calculator/enr-calculator_/enr-calculator bazel-bin/cmd/prysmctl/prysmctl_/prysmctl bazel-bin/tools/bootnode/bootnode_/bootnode
rm -rf bazel-bin
