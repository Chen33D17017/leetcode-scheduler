# Leetcode Scheduler
This repository provides a web application which can schedule your review on problem in [Leetcode](https://leetcode.com/problemset/all/).
Inspired by [Anki](https://github.com/ankitects/anki), the review stradegy in this application is according to the difficulty of problem.
The harder problem will appear more frequently in your schedule so that your practice can be more effeciently.
## Demo
The demo website build on AWS: [leetcode-scheduler](http://ec2-13-230-102-44.ap-northeast-1.compute.amazonaws.com/login)(Not assure)

## Requirement
- Go
Please install the library below
    * [package html](https://godoc.org/golang.org/x/net/html)
    * [gorilla/mux](https://github.com/gorilla/mux)
    * [Go MySQL Driver](https://github.com/go-sql-driver/mysql)
- MySQL
  - The DDL of the database is in the other/*.sql
  - To import all the problem in leetcode, please see [leetcode-parser](https://github.com/Chen33D17017/Leetcode-parser)

If you want to deploy the website by yourself, add the config file after clone this repository

config.json

``` json
{
  "type": "mysql",
  "endpoint": "mysql ipaddress",
  "id": "mysql user name",
  "password": "password",
  "database": "leetcode_scheduler"
}
```