#!/usr/bin/env bash

set -xe

report_dir=/tmp/karavel-e2e-report
mkdir -p $report_dir

tar -C $report_dir -xzvf $(sonobuoy retrieve)
npx xunit-viewer -r "$report_dir/plugins/karavel-e2e/results/global/results.xml" -o report.html
