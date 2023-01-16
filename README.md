# fluent-bit redis-metrics output plugin

Original project was https://github.com/majst01/fluent-bit-go-redis-output.

This is a fork that handles a specific log structure with the following fields:
 - 'a': "application" (scope of the counters)
 - 'm': "metric" (name of the metric)
 - 'v': value (if the metric is an accumulative counter)
 - 'd': value (if the metric is a discrete value that replaces previous value - eg: a gauge)

for each application and each day, a redis hash is created, which contains one key for each metric.
