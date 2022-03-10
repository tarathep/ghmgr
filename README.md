# GHMGR : GitHub Manager
Support GitHub Enterprise/Organize management GitHub APIs tools 


## List

### List member in team

**option**

```-t,--team``` team name (team in GitHub)

```bash
ghmgr list member --team [teamname]
```

### List member in team status pending

**option**

```-t,--team``` team name (team in GitHub)

```-p,--pending``` invited status pending


```bash
ghmgr list member --team [teamname] --pending
```

### List member in team role

**option**

```-t,--team``` team name (team in GitHub)

```-r,--role``` role team

```bash
ghmgr list member --team [teamname] --role
```

### List member in CSV file

**option**

```-f,--file``` Filename.CSV


```bash
ghmgr list member -file filename.csv
```


## Invite Member

### Invite member single command

**option**

```-t,--team``` team name (team in GitHub)

```-e,--email``` email

```-r,--role``` role of team (maintainer/ member)

```bash
ghmgr invite member --team [teamname] --email mail@domain.com --role member
```

### Invite Team Member via CSV Template

**option**

```-f,--file``` file name

```bash
ghmgr invite member -f filename.csv
```

## Export

### Export Team Member to CSV Template

**option**

```-t,--team``` team name (team in GitHub)


```bash
ghmgr export member --team [teamname]
```



