# balancer-api

Golang REST API for storing balance sheet records.

## REST API Endpoints

GET `/records`

- returns a list of stored records

### Example Response

```
[
	{
		"id": 1,
		"name": "test 1",
		"balance": 23.01,
		"type": "ASSET"
	},
	{
		"id": 2,
		"name": "test 2",
		"balance": -234234.4,
		"type": "LIABILITY"
	},
	...
]
```

POST `/records`

- creates a new record

### Example Request Body

```
{
	"name": "Test 3",
	"type": "Asset",
	"balance": 100.50
}
```

### Example Response

```
{
	"recordId": 1
}
```

PUT `/records/{id}`

- updates a record by `id`
- returns 200 on successful update

### Example Request Body

```
{
	"name": "updated",
	"type": "asset",
	"balance": 20
}
```

DELETE `/records/{id}`

- deletes a record by `id`

### Example Response

```
{
	"recordId": 1
}
```

GET `/records/net`

- returns the sum of all records

### Example Response

```
{
	"value": 1000.4
}
```

GET `/records/sum`

- Query Parameters:
  - `type`: "asset" or "liability"
- returns the sum of all records matching the query parameter `type`

```
{
	"value": 100.62
}
```

## Testing

BDD tests are written for all HTTP handlers using Ginkgo and Gomega.

To run tests, either run `go test ./...` from the project root or, if you have ginkgo installed, run `ginkgo test ./...`

View a coverage report with `go test ./... -coverprofile <filename>` or with ginkgo, `ginkgo test -r --cover`
