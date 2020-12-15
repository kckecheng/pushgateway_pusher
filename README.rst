About
======

This tool extracts metrics from benchmark tools (VDBench, iPerf and more) and push the collected metrics to a Prometheus Pushgateway for prometheus integration.

Usage
------

::

  ./pushgateway_pusher --help
  <benchmark tool> <options ...> | ./pushgateway_pusher -a <application name> -f <fields definition yaml> -g <Prometheus pushgatway such as http://localhost:9091> -j <job name>

Notes
------

- The pattern to extract fields should follow `go regexp syntax <https://pkg.go.dev/regexp/syntax>`_
- This tool depends on line buffering. If a program (such as iperf3) does not use line buffering, this tool won't work. The workaround on Linux is changing the buffer options with **stdbuf** as below (I have no idea on how to achieve the same effect on Windows:():

  ::

    stdbuf -oL -eL iperf3  -c 192.168.100.100 -t 3600 -i 10 -f M |\
      ./pushgateway_pusher -a iperf3 -g http://192.168.100.200:9091 -j iperf_job1 -f iperf3_tcp.yaml
- Since regular expression patterns to extract metrics may be complicated, turn on debug output as below for troubleshooting:

  ::

    # exprot DEBUG=false
    export DEBUG=true
    ...
