#!/usr/bin/env bash

results_dir="${RESULTS_DIR:-/tmp/results}"
mkdir -p ${results_dir}
go tool test2json -p github.com/karavel/e2e/tests/pkg -t e2e.test -test.v | tee report.log
cat report.log | go-test-report --title 'Karavel E2E Tests' --output report.html
mv report.log report.html ${results_dir}
cd ${results_dir}
tar czf results.tar.gz *
echo ${results_dir}/results.tar.gz > ${results_dir}/done
