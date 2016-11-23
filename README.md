# Togglr

Togglr is a tool for managing your Toggl account and help in you work routine.

# Install

```sh
go get -u github.com/pedronasser/togglr
```

# Quick start

Now that you have installed `togglr` in your system, login to your Toggl account. 

```sh
togglr login
```

Now that you are logged, check your projects.

```sh
togglr projects
```

Checking your account summary

```sh
togglr summary
```

Is a good practice to define an alias for each project, so you don't need to know their ID everytime.

```sh
#                     PROJECT ID     ALIAS
togglr projects alias 98172389127 mycoolproject
```

## Configurations

Some important configurations:

```sh
togglr config name "Your Name"
```

```sh
#                HOUR RATE
togglr config rate 100
```

```sh
#                 PREFERRED CURRENCY
togglr config currency USD
```

## Control your timers

### Starting your timer

```sh
#            PROJECT ID or ALIAS        DESCRIPTION
togglr start myultraawesomeproject "Doing some awesome stuff"
```

### Stopping timer

```sh
togglr stop
```

## Sending an invoice 

### Generating PDFs

Example

```sh
togglr invoice --detailed --client "Client Company" --projects "mycoolproject,othercoolproject"
```

# What's next

- Sending invoices to email
- Update timer description
- Better timer/project management
- Rate per project
- Better summary 