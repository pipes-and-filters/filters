# Filters
Filters Package for Golang

The Filters package is designed to allow the chaining of cli
commands using stdin and stdout.

[![GoDoc](https://godoc.org/github.com/pipes-and-filters/filters?status.svg)](https://godoc.org/github.com/pipes-and-filters/filters)
[![Build Status](https://travis-ci.org/pipes-and-filters/filters.svg?branch=master)](https://travis-ci.org/pipes-and-filters/filters)
[![Go Report Card](https://goreportcard.com/badge/github.com/pipes-and-filters/filters)](https://goreportcard.com/report/github.com/pipes-and-filters/filters)

The command allows setting the commands explicitly or loading from a yaml file.

##Filter
A Filter is any cli command, but primarily focused on the idea of unix cli commands being used as byte filters used to mangle and transform text and piping commands together.

##Chain
A Chain is a sequential list of Filters with each piping its output into the next.

##Exec
From either a chain or filter an Exec can be generated.

Input for an Exec is provided as a io.Reader and output is assigned as an io.Writer.

On an Exec from a  chain, input is assigned to the first process and output assigned to the last.  Everything in between is automatically piped to the next Filter.

##Chains
Chains allow the application to load multiple chains in one yaml file and retrive them by title.

##Yaml file definitions
Filter YAML Example
```yml
Name: 'cat'
Domain: 'bash'
Version: '1.0'
Command: 'cat'
Arguments:
```

Chain YAML Example
```yml
Chain:
- Name: 'cat'
  Domain: 'bash'
  Version: '1.0'
  Command: 'cat'
  Arguments:
- Name: 'grep'
  Domain: 'bash'
  Version: '1.0'
  Command: 'grep'
  Arguments:
          - 'wrong'
- Name: 'xargs'
  Domain: 'bash'
  Version: '1.0'
  Command: 'xargs'
  Arguments:
          - '-n'
          - '3'
```
Chains YAML Example
```yml
FirstChain:
        Chain:
        - Name: 'ls'
          Domain: 'bash'
          Version: '1.0'
          Command: 'ls'
          Arguments:
          VCS:
                  Type: 'git'
                  Location: 'github.com'
        - Name: 'grep'
          Domain: 'bash'
          Version: '1.0'
          Command: 'grep'
          Arguments:
                  - 'filters'
          VCS:
                  Type: 'git'
                  Location: 'github.com'
        - Name: 'xargs'
          Domain: 'bash'
          Version: '1.0'
          Command: 'xargs'
          Arguments:
                  - '-n'
                  - '4'
          VCS:
                  Type: 'git'
                  Location: 'github.com'
SecondChain:
        Chain:
        - Name: 'ls'
          Domain: 'bash'
          Version: '1.0'
          Command: 'ls'
          Arguments:
          VCS:
                  Type: 'git'
                  Location: 'github.com'
        - Name: 'grep'
          Domain: 'bash'
          Version: '1.0'
          Command: 'grep'
          Arguments:
                  - 'filters'
          VCS:
                  Type: 'git'
                  Location: 'github.com'
        - Name: 'xargs'
          Domain: 'bash'
          Version: '1.0'
          Command: 'xargs'
          Arguments:
                  - '-n'
                  - '4'
          VCS:
                  Type: 'git'
                  Location: 'github.com'
```

See the [godocs](https://godoc.org/github.com/pipes-and-filters/filters) for more information and examples.
