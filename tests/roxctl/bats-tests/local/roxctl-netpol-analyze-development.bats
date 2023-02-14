#!/usr/bin/env bats

load "../helpers.bash"
out_dir=""
templated_fragment='"{{ printf "%s" ._thing.image }}"'

setup_file() {
    #command -v yq >/dev/null || skip "Tests in this file require yq"
    #echo "Using yq version: '$(yq4.16 --version)'" >&3
    # as of Aug 2022, we run yq version 4.16.2
    # remove binaries from the previous runs
    [[ -n "$NO_BATS_ROXCTL_REBUILD" ]] || rm -f "${tmp_roxctl}"/roxctl*
    echo "Testing roxctl version: '$(roxctl-development version)'" >&3
}

setup() {
  out_dir="$(mktemp -d -u)"
  ofile="$(mktemp)"
}

teardown() {
  rm -rf "$out_dir"
  rm -f "$ofile"
}

@test "roxctl-development analyze netpol generates connlist output" {
  assert_file_exist "${test_data}/np-guard/netpols-analysis-example/ns.yaml"
  assert_file_exist "${test_data}/np-guard/netpols-analysis-example/netpols.yaml"
  assert_file_exist "${test_data}/np-guard/netpols-analysis-example/kubernetes-manifests.yaml"
  echo "Writing connlist to ${ofile}" >&3
  run roxctl-development analyze netpol "${test_data}/np-guard/netpols-analysis-example"
  assert_success

  
  #echo "$output" >&3
  echo "$output" > "$ofile"
  assert_file_exist "$ofile"
  assert_output --partial 'default/checkoutservice[Deployment] => default/cartservice[Deployment] : TCP 7070'

}
