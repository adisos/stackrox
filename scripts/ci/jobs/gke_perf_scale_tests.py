#!/usr/bin/env -S python3 -u

"""
Run the perf scale test in a GKE cluster
"""
import os
from runners import ClusterTestRunner
from clusters import GKECluster
from pre_tests import NullPreTest
from ci_tests import NullTest
from post_tests import NullPostTest

os.environ["ORCHESTRATOR_FLAVOR"] = "k8s"

ClusterTestRunner(
    cluster=GKECluster("perf-scale-test"),
    pre_test=NullPreTest(),
    test=NullTest(),
    post_test=NullPostTest(),
).run()
