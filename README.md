Crusher - a command line tool for creating and managing database views

Basic pseudocode for v1

# Startup the program
- Parse arguments to command line execution
- Import config settings from included .env file. These variables will have the blacklisted table names & database config settings
- If args[1] == 'create' or 'update', check if there is a second argument which is a path to a .sql file

# Check for accidental messing with production tables & data
- The file name will become the view name. If the file is there, ensure the file name is not one of the blacklisted table names (the names of the "real" tables & views in our production environment)

# Check for accidental SQL injection
- If the name is good, parse the .sql file to ensure that it does not include any of the banned words ('create', 'delete', 'refresh', 'update', 'insert', 'drop')
- If none of the banned words are there, make sure that the first word is 'select'
- If the first word is 'select', make sure the last character is ';' and that there is only one ';' in the file

# Build query and execute
- Build the query according to the arguments given
- Connect to the database
- Execute query after successful connection to database


Todos post MVP

- Implement Slack integration to send notifications to #analysis-team when a view is created, updated or refreshed

- Implement auth so the Slack integration notes who made the change

- Connect to GitHub so we can accept a link to a GitHub file in place of a file path on the local machine

- Deploy on a server and connect to GitHub webhooks so we can auto-deploy new or updated views whenever a new file is pushed to SQLovers
