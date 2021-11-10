# go-notion

A Go wrapper around the Notion API.

| ⚠ This package is new and under active development.

## How to Use

Install the package

```bash
go get github.com/bmorrisondev/go-notion
```

Create a client with your integration token (more info on that [here](https://developers.notion.com/docs/authorization))

```go

client := NotionClient{
  Token: "MY_INTEGRATION_TOKEN",
}

```

Example: Query a database

```go
	shareAt := "Share At"
	asc := "ascending"
	filter := QueryFilter{
		Filter: &Filter{
			Property: &shareAt,
			Date: &DateFilter{
				Before: &now,
			},
		},
		Sorts: &[]Sort{
			{
				Property:  &shareAt,
				Direction: &asc,
			},
		},
	}

	results, err := client.QueryDatabase("DATABASE_ID", &filter)
```


## How to Contribute

- Fork the repo
- Make your changes
- Submit a PR

## Contact Me

For more info, you can reach me on the Learn Build Teach Discord at [https://discord.gg/vM2bagU](https://discord.gg/vM2bagU).