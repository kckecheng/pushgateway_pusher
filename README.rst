About
======

A simple tool which supports extracting metrics from multiple tools (VDBench, iPerf and more) and push the metrics to a Prometheus Pushgateway for smooth observation.

Usage
------

::

  ./pushgateway_pusher --help
  <some tool> <tool options ...> | ./pushgateway_pusher -a vdbenchfs -f <fields definition yaml> -g <Prometheus pushgatway such as http://localhost:9021> -j <job name>

Notes
------

- The pattern to extract fields should follow https://pkg.go.dev/regexp/syntax
