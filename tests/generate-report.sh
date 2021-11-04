#!/usr/bin/env bash

set -xe

report_dir=/tmp/karavel-e2e-report
rm -rf $report_dir report.html
mkdir -p $report_dir

sonobuoy retrieve $report_dir --extract
npx xunit-viewer -r "$report_dir/plugins/karavel-e2e/results/global/report.xml" -o report.html

if [ -x "$(command -v xdg-open)" ]; then
  open_cmd="xdg-open"
elif [ -x "$(command -v xdg-open)" ]; then
  open_cmd="open"
fi

if ! [ -z "$open_cmd" ]
then
  $open_cmd report.html
fi
