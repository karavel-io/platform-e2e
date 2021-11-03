#!/usr/bin/env bash

results_dir="${RESULTS_DIR:-/tmp/results}"
mkdir -p ${results_dir}
cd ${results_dir}
e2e.test -test.v | tee report.log
cat report.log | go-junit-report > report.xml
tar czf results.tar.gz *
echo ${results_dir}/results.tar.gz > ${results_dir}/done
