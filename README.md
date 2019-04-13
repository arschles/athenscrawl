# crathens

crathens makes a request to the `/catalog` endpoint of any Athens server (that you define in the `GOPROXY` env var) and, for each module in the returned list that has `github.com` at the beginning of the path:

- Requests the list of tags for the corresponding GH repository
- Sends a `/{module}/@v/{tag}.info` request to Athens for each tag in the response body

The crathens daemon rate limits both steps to prevent flooding the GH API (and getting rate limited) and to prevent overloading the Athens server at `GOPROXY`.
