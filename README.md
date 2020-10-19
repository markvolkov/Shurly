---
# Systems Engineering Project
> Submission for CloudFlare-Hiring 2020
> By Mark Volkov

# Output from running 3 profiles
> ./main -url="https://linktree-style-worker.markvolkov.workers.dev" -profile=10
```
{
   "totalRequestCount":10,
   "failedRequestCount":0,
   "successfulRequestCount":10,
   "originalRequest":{
      "address":"https://linktree-style-worker.markvolkov.workers.dev",
      "addressStripped":"linktree-style-worker.markvolkov.workers.dev",
      "requestedRoute":"/"
   },
   "minResponseSize":0,
   "maxResponseSize":255,
   "meanResponseTime":100,
   "medianResponseTime":0,
   "fastestResponseTime":0,
   "slowestResponseTime":1000
}
2020/10/19 15:26:51 Percentage of successful requests (10 / 10) = %100
```
> ./main -url="https://linktree-style-worker.markvolkov.workers.dev/links" -profile=10
```
{
   "totalRequestCount":10,
   "failedRequestCount":0,
   "successfulRequestCount":10,
   "originalRequest":{
      "address":"https://linktree-style-worker.markvolkov.workers.dev",
      "addressStripped":"linktree-style-worker.markvolkov.workers.dev",
      "requestedRoute":"/links"
   },
   "minResponseSize":0,
   "maxResponseSize":225,
   "meanResponseTime":100,
   "medianResponseTime":0,
   "fastestResponseTime":0,
   "slowestResponseTime":1000
}
2020/10/19 15:27:12 Percentage of successful requests (10 / 10) = %100
```
> ./main -url="https://www.netflix.com" -profile=10
```
{
  "totalRequestCount":10,
  "failedRequestCount":0,
  "successfulRequestCount":10,
  "originalRequest":{
     "address":"https://www.netflix.com",
     "addressStripped":"netflix.com",
     "requestedRoute":"/"
  },
  "minResponseSize":0,
  "maxResponseSize":251,
  "meanResponseTime":400,
  "medianResponseTime":500,
  "fastestResponseTime":0,
  "slowestResponseTime":1000
}
2020/10/19 15:27:36 Percentage of successful requests (10 / 10) = %100
```

# Simple but interesting findings...
I found that the url requests made to the links associated with my worker site were much more consistent than most other sites in terms of the max response size, the median response time, and the mean response time, I also found that most successful requests were completed in at most 1000ms or 1 second.
