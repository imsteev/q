# q
query notepad - remember your useful queries!

installation
```zsh
# clone this repo and cd into it
go install .
```

commands
```
> q <query-name>          # list all
> q <query-name> <query>  # add query
> q -d <query-name>       # delete query
```

example
```
> q
Queries

> q revenue-from-tickets 'select ...'
[revenue-from-tickets]	select ...

> q revenue-from-saas 'select ...'
[revenue-from-saas]	select ...

> q
Queries
1. revenue-from-tickets
2. revenue-from-saas

# delete revenue-from-tickets
~/go/src/q [main] q -d revenue-from-tickets
Deleted query: "revenue-from-saas"

> q
Queries
1. revenue-from-saas
```
