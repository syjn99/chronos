bazel build --remote_cache=http://192.168.2.200:9090 --sandbox_debug --verbose_failures --config=release --config=windows_amd64 //cmd/beacon-chain //cmd/validator //cmd/prysmctl //tools/enr-calculator //tools/bootnode
if [ $? -ne 0 ]; then
    echo "Bazel build failed."
    exit 1
fi
zip -j dist/chronos_windows.zip bazel-bin/cmd/beacon-chain/beacon-chain_/beacon-chain.exe bazel-bin/cmd/validator/validator_/validator.exe bazel-bin/tools/enr-calculator/enr-calculator_/enr-calculator.exe bazel-bin/cmd/prysmctl/prysmctl_/prysmctl.exe bazel-bin/cmd/bootnode/bootnode_/bootnode.exe
mkdir -p dist/bin/win
cp bazel-bin/cmd/beacon-chain/beacon-chain_/beacon-chain.exe dist/bin/win/chronos.exe
cp bazel-bin/cmd/validator/validator_/validator.exe dist/bin/win/validator.exe
rm -rf bazel-bin
