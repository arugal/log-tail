#!/usr/bin/env bash
# compile for version
make
if [ $? -ne 0 ]; then
    echo "make error"
    exit 1
fi

log_tail_version=`./bin/log-tail --version`
echo "build version: log_tail_version"

# cross_compiles
make -f ./Makefile.cross-compiles

rm -rf ./packages
mkdir ./packages

os_all='linux windows darwin freebsd'
arch_all='386 amd64 arm arm64 mips64 mips64le mips mipsle'

for os in $os_all; do
    for arch in $arch_all; do
        log_tail_dir_name="log_tail_${log_tail_version}_${os}_${arch}"
        log_tail_path="./packages/log_tail_${log_tail_version}_${os}_${arch}"

        if [ "x${os}" = x"windows" ]; then
            if [ ! -f "./log_tail_${os}_${arch}.exe" ]; then
                continue
            fi
            mkdir ${log_tail_path}
            mv ./log_tail_${os}_${arch}.exe ${log_tail_path}/log-tail.exe
        else
            if [ ! -f "./log_tail_${os}_${arch}" ]; then
                continue
            fi
            mkdir ${log_tail_path}
            mv ./log_tail_${os}_${arch} ${log_tail_path}/log-tail
        fi
        cp ./LICENSE ${log_tail_path}
        cp -rf ./conf/* ${log_tail_path}

        # packages
        cd ./packages
        if [ "x${os}" = x"windows" ]; then
            zip -rq ${log_tail_dir_name}.zip ${log_tail_dir_name}
        else
            tar -zcf ${log_tail_dir_name}.tar.gz ${log_tail_dir_name}
        fi
        cd ..
        rm -rf ${log_tail_path}
    done
done
