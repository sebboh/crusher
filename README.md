[![Build Status](https://travis-ci.org/devonestes/crusher.svg?branch=master)](https://travis-ci.org/devonestes/crusher)
# Crusher - a command line tool for creating and managing database views
Crusher is a command line application that allows a user to have direct control over their own database views & materialized views on a production database. This should allow them to greatly increase their productivity by saving useful queries and chunks of code as views or materialized views, upon which they can build further queries, and also to share amongst other users of that database. It will also remove the dependency of the engineering to implement queries or modifications that users would like to see on the production database, removing a significant blocker from the user's workflow.

## Usage
### Commands

##### Creating a view
```
crusher create new_view.sql
```
This will create a new view on the production database named `new_view`.

##### Creating a materialized view
```
crusher create -m new_view.sql
```
This will create a new materalized view on the production database named `new_view`.

##### Updating a view
```
crusher update new_view.sql
```
This updates the view or materalized view named `new_view` by replacing the existing view with the query in the given file.

##### Refreshing a materalized view
```
crusher refresh new_view
```
This refreshes the materalized view named `new_view`

##### Help
```
crusher help
```
This shows information about how to use the program

### SQL File Requirements

For the sake of avoiding accidental issues with manipulating data on the production database, the SQL files that are passed to this program need to meet somewhat strict requriements.

* Every file must contain a single `select` statement. This means the file must begin with the word `select` - no spaces or comments in front of it.
* You cannot use a semi-colon ANYWHERE in the query.
* The name of the file must not be on the list of blacklisted table names.
* The words `create`, `delete`, `refresh`, `update`, `insert`, and `drop` cannot appear anywhere in the file. Words which contain those words, like `created_at`, are just fine.


### Configuration

There are two environment variables that need to be set before using.

1) `DB_URL`, which is the full URL for the postgres database that you wish to modify.

2) `BLACKLISTED_NAMES`, which is a list of table and/or view names that you do NOT want your users to be able to change (most likely by accident). The format for this should be `,table_1,table_2,table_3,...,table_10,`, with a comma before the first name, a comma between each table name, and a comma after the last name.

## Todos:

- Implement Slack integration to send notifications to a channel when a view is created, updated or refreshed

- Implement auth so the Slack integration notes who made the change

- Connect to GitHub so we can accept a link to a GitHub file in place of a file path on the local machine

- Deploy on a server and connect to GitHub webhooks so we can auto-deploy new or updated views whenever a new file is pushed to a given repo
