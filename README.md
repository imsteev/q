# q
### q's a query notepad - remember useful queries!

todo
- [ ] auto-complete by query name
- [ ] ability to edit query in default shell editor

installation
```zsh

# choose where you want to store data. the file itself does not need to exist
export QUERYLIST_FILE_PATH=~/Desktop/.querylist.json

# clone this repo and install the command
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

> q -d revenue-from-tickets
Deleted query: "revenue-from-saas"

> q
Queries
1. revenue-from-saas
```
