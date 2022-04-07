# GHMGR : GitHub Manager 1.1.0
Support GitHub Enterprise/Organize management GitHub APIs tools 


## Installation

Download bin file and set env path (depending on OS)


## Login

### Login GitHub Personal Token

on Windows with 

Command Prompt

```cmd
set GHMGR_TOKEN=php_xxxxxxxxxxxxxxxxxx
set GHMGR_OWNER=xxxxx
```
PowerShell
```ps
$env:GHMGR_TOKEN="php_xxxxxxxxxxxxxxxxxx"
$env:GHMGR_OWNER="corp-ais"
```


on Unix or MacOS

```bash
export GHMGR_TOKEN=php_xxxxxxxxxxxxxxxxxx
export GHMGR_OWNER=xxxxx
```

or login with app (don't on support CMD,PS)

```bash
ghmgr login --token php_xxxxxxxxxxxxxxxxxx --owner xxxxx
```

or on pipeline you can use option --token recommand for secret

example

```bash
ghmgr list member -t teamname --token php_xxxxxxxxxxxxxxxxxx
```


## Load

### Init Cache

for the improve performance and query must init before use in the first time

```bash
ghmgr load cache 
```



## List

### List Team in Organization

Lists all teams in an organization that are visible to the authenticated user.
https://docs.github.com/en/rest/reference/teams#list-teams

```bash
 ghmgr list team
```

### List Membership of Teams

**option**

```-u,--username``` username

```bash
ghmgr list team --username [username]
```

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


### List member in team Email

**option**

```-t,--team``` team name (team in GitHub)

```-m show,--email show``` email show

```bash
ghmgr list member --team [teamname] --email show
```

### List member in team Exclude another team

**option**

```-t,--team``` team name (team in GitHub)

```-e ,--exclude``` exclude team member will be invisible

```bash
ghmgr list member --team [teamname] --exclude [teamname]
```



### List member in CSV file

for the files workspace in : ```report/input/teamname.csv```

**option**

```-f,--file``` Filename.CSV


```bash
ghmgr list member --file teamname.csv
```

### List member of Organization

**option**

```-o,--org``` Org


```bash
ghmgr list member --org
```

### List member of Organization Show Email , Teams

**option**

```-o,--org``` Org

```-m show,--email show``` Email

```-t show,--team show``` Teams


```bash
ghmgr list member --org --email show --team show
```

### List Dormant users of the organization from CSV file

**option**

for the file workspace in : ```report/input/dormant-users-report-xxxxxx.csv```

```-f,--file``` filename.csv


```bash
ghmgr list member dormant --file dormant-users-report-xxxxxx.csv
```

---

## Invite Member

### Invite member single command

**option**

```-t,--team``` team name (team in GitHub)

```-e,--email``` email

```bash
ghmgr invite member --team [teamname] --email name@domain.com
```



### Invite Team Member via CSV Template

**option**

```-f,--file``` file name

```bash
ghmgr invite member -f filename.csv
```

---

## Add

### Add or Update Membership of Team with Username

**option**

```-t,--team``` team name (team in GitHub)

```-u,--username``` GitHub username

```-r,--role``` Role of Team (maintainer, member)

```bash
ghmgr add member --team [teamname] --email [name@domain.com] -r maintainer
```

### Add or Update Membership of Team with Email

must load cache and add verify email account before use this cmd.

**option**

```-t,--team``` team name (team in GitHub)

```-m,--email``` email (verify)

```-r,--role``` Role of Team (maintainer, member)

```bash
ghmgr add member --team [teamname] --email [name@domain.com] -r maintainer
```


---

## Export

### Export Team Member to CSV

**option**

```-t,--team``` team name (team in GitHub)


```bash
ghmgr export member --team [teamname]
```



---


## Remove

### Remove Team Membership For User

To remove a membership between a user and a team (without Organization)


**option**

```-t,--team``` teamname (slug team name lowercase)

```-u,--username``` GitHub username

```bash
ghmgr remove member --team [teamname] --username [username]
```

### Remove Membership of Organization 

Removing a user from this list will remove them from all teams and they will no longer have any access to the organization's repositories


**option**

```-o,--org``` organization

```-u,--username``` GitHub username

```bash
ghmgr remove member --org --username [username]
```




### Remove dormant users of the organization from CSV file

for the file workspace in : ```report/input/dormant-users-report-xxxxxx.csv```

```-f,--file``` filename.csv

```-b,--backup``` the backup file for review

```bash
ghmgr remove member dormant --file dormant-users-report-xxxx.csv --backup
```

the backup file output into report/output/dormant-users-report-xxxxxx.csv-review-xxxx.csv

---

