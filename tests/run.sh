#!/usr/bin/env bash

results_dir="${RESULTS_DIR:-/tmp/results}"
mkdir -p ${results_dir}
go tool test2json -p github.com/karavel/e2e/tests/pkg -t e2e.test -test.v | tee ${results_dir}/report.log
echo ${results_dir}/report.log > ${results_dir}/done
