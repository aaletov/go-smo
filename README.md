# go-smo

## Notes
### Set manager
We have apropriate information about the system only 
at the moment of newest request, because no further 
reqs could be generated before it. Suppose at 1'st 
iteration we have 3 requests generated with genTime 
x1 < x2 < x3. Then, at the next iteration reqs would 
have genTimes x1 + a, x2 + b, x3 + c. If a < x2 - x1 
then pushing x2 in buffer at 1'st iteration would be 
incorrect, because request x1 + a comes before x2 
(but we know it only on second iteration). SetManager 
should consider this condition and operate only with the
newest request

Therefore, SetManager should contain queue with generation
time priority. On Collect() manager collects new requests
from sources, on ToBuffer() manager pushes all requests for
which there are older requests from other sources. Suppose
we have buffer with state:

```
[
  {"SourceNumber": 1, "GenTime": 1},
  {"SourceNumber": 2, "GenTime": 2}, 
  {"SourceNumber": 3, "GenTime": 3}, 
]
```
Then, we should push 1'st request into buffer, because no
requests could be generated before it, but we can't push
2'nd requests, because source 1 could generate next
request with GenTime == 1.5.