# Question and Answer site

_19 June 12023_

## Given

- 10m registered users
- 100m views/month
- Upvotes and downvotes don't need instantaneous updates
- A user can edit their content
- **Top concern: serving content with high availability**

## Inferred

- 50KB web pages (just the HTML)
	- This is taken by averaging the size of one unpopular SO page and one popular SO page
- A question has one user
- A question has many answers
- An answer has one user
- An answer has many votes
- A vote has one user
- A user has many questions
- A user has many answers

## Assumed

- 10% of registered users are generating most of the question and answer content, though more than 50% of all registered users are providing votes
- Most questions have more than one answer

## Calculations

### Page views

#### With an even distribution of visits across days of the month and times of day

50KB web pages with 100m views/month -- 3.33m views/day -- 139k views/hour -- 40 views/second -- 50KB • 40 ≈ **2MB served per second**

#### With greater focus on the weekdays...
- 8 weekend days per month, where *site traffic is 50% of weekday traffic*
- 3.33m views/day for an even distribution --> weekend at 50% is 1.66m views/day
- Which means there's 8 * 1.66m views to distribute to the other 22 days (weekdays)
- 3.33 + 0.6 -> 3.93m views/weekday -- **this comes to something like 45 views/second and 2.25MB to serve per second**

#### ...and with emphasis on 8 hours of those weekdays
- We're looking at 163,750 views per hour if we're evenly distributing the 3.93m views
- If site traffic at 2am EST is 50% of 2pm EST traffic, then nighttime (off-hour) traffic is 81,875 views an hour. Those 81,875 * 16 views need to be distributed: 1.31m views spread out over 8 hours; this doubles daytime traffic to 327,500 views per daytime hour.
- That's about **91 views/second, and 4.55MB to serve per second**

## Data model

### Question

| Field | Type | Size |
|---|---|---|
| question_id | UUID | 16 bytes |
| user_id | UUID | 16 bytes |
| created | Datetime with timezone | 12 bytes |
| updated | Datetime with timezone | 12 bytes |
| text | string (utf-8) | 8 kilobytes |

_One row: <9KB_

### Answer

| Field | Type | Size |
|---|---|---|
| answer_id | UUID | 16 bytes |
| question_id | UUID | 16 bytes |
| user_id | UUID | 16 bytes |
| created | Datetime with timezone | 12 bytes |
| updated | Datetime with timezone | 12 bytes |
| text | string (utf-8) | 8 kilobytes |
| accepted | boolean | 1 byte |

_One row: <9KB_

### Votes

| Field | Type | Size |
|---|---|---|
| answer_id | UUID | 16 bytes |
| user_id | UUID | 16 bytes |
| up_or_down | enum | 4 bytes |

_One row: 36B_

### User

| Field | Type | Size |
|---|---|---|
| user_id  | UUID | 16 bytes |
| email | string (utf-8) | 64 bytes |
| display_name | string (utf-8) | 32 bytes |
| created | Datetime with timezone | 12 bytes |
| updated | Datetime with timezone | 12 bytes |

_One row: 136B_

Given the lack of a need for user profile pages or other indexes for individual users, we will omit any sort of associative table that quickly lets us retrieve all of a user's questions, answers, or votes.

### QuestionAnswer

Associative table

### QuestionVote

Associative table

### AnswerVote

Associative table

## System

Given the number of requests we're expecting per second, as well as the amount of content we would be serving per second, and our primary design directive being given as serving content with high availability, we don't actually need all that much to meet our goals.

### Just beyond the minimum design

```ascii

                         |---> [ server 1 ]------|
                         |          |            |--> [ replica 1 ]
[ client ] --> [ lb ] ---|     [ primary db ] -->|
                         |          |            |--> [ replica 2 ]
                         |---> [ server 2 ]------|

```

The main features of this design are a _load balancer_ for routing requests between _two instances of the application server_ for increased availability. There is _one primary database_ to which both servers write; there are _two read-only replicas_, from which the servers may read and the primary database may write.

#### Writing content

As an example: a client sends an *HTTP POST request* to add an answer for a question. The body looks like this:

```
user_id
question_id
text
```

The HTTP request first encounters the _load balancer_ (`lb`) which routes the request to *one of two servers*. After performing the necessary validation, a new row in the `answers` table is created in the _primary database_. The server responds with a `200 OK` status code. Asynchronously, the new row is added to _two read-only replicas_.

#### Reading content

A client issuing an _HTTP GET request_ for a question page kicks off a process for generating that page: the content (question, answers, votes, along with the respective user information) is fetched from _either read replica_. The application server then injects that data into a template and responds with an HTML page.

This process becomes repetitive and a strain on resources if requests spike for a single page; in this first proposed design, a spike would look like _an increase of one order of magnitude_ of requests for the same page per second (50KB page served by the application server after reading ~30KB of data from the read-only replicas). So, we can expect a slow-down if we're attempting to process 500 read requests per second.

### Improving performance

Use a content delivery network (CDN) to serve pages so that we may improve response time for our geographically-distributed user base and ease the load on our server when experiencing spikes.

When serving the webpage that contains the answer and its question, we can retrieve vote counts through a separate, asynchronous call. This lets us keep the HTML document on the CDN longer without having to bust the cache, as well as allows us to provide a more response experience for our primary concern: serving questions and answers.

Another vector to consider for improving read performance is how we handle tallying votes in the first place. Our most-simple approach can run a query asynchronously to get a count of each type of vote for each of the answers. To improve read speed, we could update our data model for `Answers` to add two more fields of the `int` type: `up_votes` and `down_votes`. A new answer would be initialized with a `0` value for both kinds of votes. As a user votes in either direction, we already record that action in the `Votes` table. A change to make to improve performance would be to periodically update the `up_votes` and `down_votes` values for each `answer` row after a vote is issued (this may be immediately after, or in a batch process to deal with instances where a great many of users are voting on a question at the same time).

## References

- Various knowledge cards in search results for Postgres type sizes
- [Fermi estimates on Postgres performance](https://www.citusdata.com/blog/2017/09/29/what-performance-can-you-expect-from-postgres)
- [How Many Requests Can a Real-World Node.js Server-Side App Handle?](https://plainenglish.io/blog/how-many-requests-can-handle-a-real-world-nodejs-server-side-application)
