---
# Systems Engineering Project
> By Mark Volkov
>
> Profile the performance of your webpages from a client point of view!

# Build the script
> go build main.go

# Usage
> ./main -url="your-website.com" --> This will give a high level overview of your web page.
> 
> ./main -url="your-website.com -profile=amountOfRequests --> This will make amountofRequests requests to your web page and give you profiling information.

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
   "minResponseSize":51,
   "maxResponseSize":250,
   "meanResponseTime":125,
   "medianResponseTime":160,
   "fastestResponseTime":81,
   "slowestResponseTime":187
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
   "minResponseSize":35,
   "maxResponseSize":252,
   "meanResponseTime":59,
   "medianResponseTime":57,
   "fastestResponseTime":45,
   "slowestResponseTime":91
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
   "minResponseSize":12,
   "maxResponseSize":218,
   "meanResponseTime":492,
   "medianResponseTime":417,
   "fastestResponseTime":315,
   "slowestResponseTime":874
}
2020/10/19 15:27:36 Percentage of successful requests (10 / 10) = %100
```
