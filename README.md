
# Concurjob

CLI App to browse, filter, and apply for jobs and internships straight from your cmdline

Made faster using goroutines and concurrency


## CLI Reference

#### Build and Run Concurjob

```go
    go build -o concurjob main.go
    ./concurjob
```

| Options | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `-version` | `bool` | Current concurjob version |
| `-help` | `bool` | Information about cli options |
| `-intern` | `bool` | Filter for intern positions |
| `-fulltime` | `bool` | Filter for fulltime positions |
| `-flag` | `string` | Comma separated keywords to filter positions |
| `-limit` | `uint` | Limit job results using an unsigned int |
| `-company` | `string` | Filter jobs by their company name |

## TODOs

[] Additional filtering based on job description
[] Adding more sources of scraping - Both FT and Intern
[] Expanding beyond SWE and tech positions
[] Distributing and organizing main.go
[] Utilizing hashing O(1) lookup instead of slice iteration for -flags
