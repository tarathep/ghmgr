# GHMGR : GitHub Manager
Support GitHub Enterprise/Organize management GitHub APIs tools 


## List

### List member in team

**option**

```-t,--team``` team name (team in GitHub)


```bash
ghmgr list member -t [teamname]
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

### Invite member via CSV Template

**option**

```-f,--file``` file name

```bash
ghmgr invite member -f filename.csv
```

