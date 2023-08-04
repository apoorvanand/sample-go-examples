import influxdb

client = influxdb.InfluxDBClient('localhost', 8086, 'root', 'root', 'jmeter')

# The query to fetch response time, count, and errors.
query = """
SELECT mean("responseTime") AS responseTime, count(*) AS count, sum("errors") AS errors
FROM measurements
WHERE name = 'response_time'
GROUP BY time(1m)
ORDER BY time DESC
"""

# Fetch the results of the query.
results = client.query(query)

# Print the results of the query.
for result in results:
  print(result['time'], result['responseTime'], result['count'], result['errors']
