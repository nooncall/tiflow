## Preparations

### Run integration tests locally

1. The following executables must be copied or generated or linked into these locations, `sync_diff_inspector` can be
   downloaded
   from [tidb-community-toolkit](https://download.pingcap.org/tidb-community-toolkit-v4.0.2-linux-amd64.tar.gz)
   , `tidb-server` related binaries can be downloaded
   from [tidb-community-server](https://download.pingcap.org/tidb-community-server-v4.0.2-linux-amd64.tar.gz):

   * `bin/tidb-server` # version >= 4.0.0-rc.1
   * `bin/tikv-server` # version >= 4.0.0-rc.1
   * `bin/pd-server`   # version >= 4.0.0-rc.1
   * `bin/pd-ctl`      # version >= 4.0.0-rc.1
   * `bin/tiflash`     # needs tiflash binary and some necessary so files
   * `bin/sync_diff_inspector`
   * [bin/go-ycsb](https://github.com/pingcap/go-ycsb)
   * [bin/etcdctl](https://github.com/etcd-io/etcd/tree/master/etcdctl)
   * [bin/jq](https://stedolan.github.io/jq/)
   * [bin/minio](https://github.com/minio/minio)

   > If you are running tests on MacOS, tidb related binaries can be downloaded from tiup mirrors, such as https://tiup-mirrors.pingcap.com/tidb-v4.0.2-darwin-amd64.tar.gz. And `sync_diff_inspector` can be compiled by yourself from source [tidb-tools](https://github.com/pingcap/tidb-tools)

   > All Tiflash required files can be found in [tidb-community-server](https://download.pingcap.org/tidb-community-server-v4.0.2-linux-amd64.tar.gz) packages. You should put `flash_cluster_manager`, `libtiflash_proxy.so` and `tiflash` into `bin` directory in TiCDC code base.

2. The following programs must be installed:

   * [mysql](https://dev.mysql.com/doc/mysql-installation-excerpt/5.7/en/) (the MySQL cli client,
     currently [mysql client 8.0 is not supported](https://github.com/pingcap/tidb/issues/14021))

3. The user used to execute the tests must have permission to create the folder /tmp/tidb_cdc_test. All test artifacts
   will be written into this folder.

### Run integration tests in docker

The following programs must be installed:

* [docker](https://docs.docker.com/get-docker/)
* [docker-compose](https://docs.docker.com/compose/install/)

We recommend that you provide docker with at least 6+ cores and 8G+ memory. Of course, the more resources, the better.

## Running

### Unit Test

1. Unit test does not need any dependencies, just running `make unit_test` in root dir of source code, or cd into
   directory of a test case and run single case via `GO111MODULE=on go test -check.f TestXXX`.

### Integration Test

#### Run integration tests locally

1. Run `make integration_test_build` to generate TiCDC related binaries for integration test

2. Run `make integration_test` to execute the integration tests. This command will

   1. Check that all required executables exist.
   2. Execute `tests/integration_tests/run.sh`

   > If want to run one integration test case only, just pass the CASE parameter, such as `make integration_test CASE=simple`.

   > There exists some environment variables that you can set by yourself, variable details can be found in [test_prepare](_utils/test_prepare).

   > `MySQL sink` will be used by default, if you want to test `Kafka sink`, please run with `make integration_test_kafka CASE=simple`.

3. After executing the tests, run `make coverage` to get a coverage report at `/tmp/tidb_cdc_test/all_cov.html`.

#### Run integration tests in docker

> **Warning:**
> These scripts and files may not work under the arm architecture,
> and we have not tested against it.
> Also, we currently use the PingCAP intranet address in our download scripts,
> so if you do not have access to the PingCAP intranet you will not be able to use these scripts.
> We will try to resolve these issues as soon as possible.

1. If you want to run kafka tests,
   run `CASE="clustered_index" docker-compose -f ./deployments/ticdc/docker-compose/docker-compose-kafka-integration.yml up --build`

2. If you want to run MySQL tests,
   run `CASE="clustered_index" docker-compose -f ./deployments/ticdc/docker-compose/docker-compose-mysql-integration.yml up --build`

3. Use the command `docker-compose -f ./deployments/ticdc/docker-compose/docker-compose-kafka-integration.yml down -v`
   to clean up the corresponding environment.

Some useful tips:

1. The log files for the test are mounted in the `./deployments/ticdc/docker-compose/logs` directory.

2. You can specify multiple tests to run in CASE, for example: `CASE="clustered_index kafka_messages"`. You can even
   use `CASE="*"` to indicate that you are running all tests。

3. You can specify in the [integration-test.Dockerfile](../../deployments/ticdc/docker/integration-test.Dockerfile)
   the version of other dependencies that you want to download, such as tidb, tikv, pd, etc.
   > For example, you can change `RUN ./download-integration-test-binaries.sh master` to `RUN ./download-integration-test-binaries.sh release-5.2`
   > to use the release-5.2 dependency.
   > Then rebuild the image with the [--no-cache](https://docs.docker.com/compose/reference/build/) flag.

## Writing new tests

New integration tests can be written as shell scripts in `tests/integration_tests/TEST_NAME/run.sh`. The script should
exit with a nonzero error code on failure.
