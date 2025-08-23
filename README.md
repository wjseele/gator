# Welcome to Gator
A CLI RSS reader that just about barely works.
But it works. And that's what matters.

# Installation

## Dependencies
Gator depends on PostgreSQL and goose to bring the database up. You also need a ```.gatorconfig.json``` file in your home directory (```~/```).

Install PostgreSQL with your favorite package manager. If you're using Linux, I assume you know how. If you're on MacOS, use brew. The application
was built on version 17, which you can install with ```brew install postgresql@17``` and then start with ```brew services start postgresql@17```.
You can check if your PostgreSQL server is running with ```psql "postgres://<USERNAME>:@localhost:5432"```. Initialize your database with
```CREATE DATABASE gator;```. 

You then need something called goose to bring the database up. You can install that with ```brew install goose```. After it's installed, go to the
```sql/schema``` folder in your gator download and run ```goose "postgres://<USERNAME>:@localhost:5432/gator" up```. This will ensure the right tables
are present.

Now that the database is up and running, open your ```.gatorconfig.json``` and tell it where to find the database. Add the following:
```{"db_url":"postgres://<USERNAME>:@localhost:5432/gator?sslmode=disable"}```. The ```?sslmode=disable``` is so it doesn't complain about you connecting
to localhost in an insecure way.

That's it for dependencies. 
## Installing Gator
Well, now that everything is set up, this is the easy part. Go back to the root folder of gator and run ```go install gator```.
You can now run gator from your command line. 

# Usage

## Commands
Gator takes one command, outputs whatever it needs to output and closes down again. The only exception to this is the agg(regation) command.
Arguments are separated by spaces ("Here is a blog" are 4 arguments). Here's what's available:
|Command| # of arguments | Type of argument | Explanation |
|-|-|-|-|
|addfeed|2|Name, url|Adds a feed to the database.|
|agg|1|Time interval|Time between querying servers. Be nice. 1m or 30s is fine|
|browse|0 or 1|Int|Number of posts you'd like to see.|
|feeds|0||List the feeds in the database.|
|follow|1|url|Follow a feed already in the database.|
|following|0||List the feeds that the current user is following.|
|login|1|username|Login as a registered user.|
|register|1|username|Register a new user.|
|reset|0||Wipe all data.|
|unfollow|1|url|Unfollow a feed.|
|users|0||List all registered users.|

## Typical Usage
Register yourself (you'll automatically be logged in). If you weren't the last user, login with your username. Add a feed.
If it already exist, follow it instead. Run agg to get posts. It'll run forever, so either open another terminal session to use gator's other functions,
or kill it with ```ctrl-c``` after some time. It'll cycle between feeds, so the more feeds you have, the longer it'll have to run for a full update.
Use browse to peruse the posts. It defaults to the latest 2, run browse with a different number if you want to see more or less.

# Potential Future Developments
If I ever get back to this (who knows), I'll probably have gator run REPL, agg concurrently and only have it update feeds that haven't been updated in
at least an hour. It should probably also have a man page. And SQLite embedded so it doesn't require a dependency dance. 
