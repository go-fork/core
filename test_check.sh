#!/bin/bash
# test_check_with_subtests.sh
# List all test cases and sub-tests, then run coverage for each

echo "Extracting all test cases and sub-tests..."
# Run tests with -v to capture sub-test names
go test -v . | grep '=== RUN' | awk '{print $3}' | sort -u > test_list.txt

echo "Running coverage for each test and sub-test..."
while read -r test; do
  echo "Running coverage for $test"
  go test -cover -run "^${test}$" .
done < test_list.txt

# Clean up
rm test_list.txt