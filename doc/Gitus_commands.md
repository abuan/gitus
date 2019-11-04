# Gitus commands 

## Project commands

Create a new project:

```
gitus project create [<name>] 
			--description Description 
			--contributors Person1, Person2.. 
			
```

Modify a project:

```
gitus project modify [<name>]
			--description Description 
			--contributors Person1, Person2.. 
```

Display a project

```
gitus project display [<name>]
```

Display all the projects
```
gitus project display_all
```


Delete a project :

```
gitus project delete [<name>]
```


## Userstory commands

Create a new userstory

```
gitus userstory create [<name>]
			--description Description 
			--effort (0,1,3,5,8,13)
```

Modify a userstory

```
gitus userstory modify [<name>]
			--description Description 
			--effort (0,1,3,5,8,13)
```
		
Display a userstory

```
gitus userstory display [<name>]
```
Display all the userstories

```
gitus userstory display_all
```



Delete a userstory

```
gitus userstory delete [<name>]
```

## Tasks commands

Create a new task
```
gitus task create [<description>]
			--people Person1, Person2.. 
			--state : (TODO, IN PROGRESS, TO VERIFY)
```

Modify a task

```
gitus task modify [<description>]
			--people Person1, Person2.. 
			--state : (TODO, IN PROGRESS, TO VERIFY)	
```

Display a task

```
gitus task display [<description>]
```
Display all the tasks

```
gitus task display_all
```

Delete a task

```
gitus task delete [<description>]
```

# TO ADD


## Associate tasks, projects and userstories

Interaction between the 3 objects (tasks, projects and userstories), for instance : add a task linked to a userstory.

## Future features ?

Authentification when you add something into a gitus project to know the history of the project