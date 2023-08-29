# docklogkeeper

[DockerHub image](https://hub.docker.com/r/nightlord189/docklogkeeper)

Simple Docker log viewer. Stores logs locally in SQLite
to allow access not to only realtime logs but also to past.

![Example](https://github.com/nightlord189/docklogkeeper/blob/master/site/screenshot1.png)

## Features
- simple search in logs history
- see logs of any Docker container (alive or stopped/removed)
- simple cookie-based authorization
- log rotation (by default, old logs will be deleted after 1 week)

## Usecase
When DockLogKeeper could be right choice for you?
+ You have multiple Docker containers on same host
+ Your containers don't generate millions of logs daily
+ You want to set up log solution by one click and see logs instantly without any scripts/settings/running big enterprise log solution

## How to use?
1. Run
```docker run --name docklogkeeper --env PASSWORD=YOUR_PASSWORD -d -v /var/run/docker.sock:/var/run/docker.sock -v docklogkeeper:/logs -p 3010:3010 nightlord189/docklogkeeper:latest```
2. Go to http://localhost:3010 and authorize with username **admin** and password from step 1.
3. See and search logs in convenient interface.
4. Select container from running on your host and you will see it's logs. 
5. Use "Update" button to get new logs instantly and "Next" to get older logs when you are scrolling to end.

## Configuration
You can use following environment variables:
+ PASSWORD - admin password. Default is keeper43.
+ LOG_LEVEL - level of logs of DockLogKeeper itself (trace, debug, info, warn, error, fatal, panic, disabled). Default is debug.
+ AUTH_SECRET - secret to encode auth cookies
+ HTTP_PORT - http port (default is 3010)
+ UPDATE_FREQUENCY - how often DockLogKeeper will check new logs from daemon (in seconds). Default value is 5.
+ LOG_RETENTION - when old logs will be deleted. Default value is 604800 (1 week).

## FAQ
1. Is it free?
Yes. It's open source and free project. You should run it on your own server.
2. Does DockLogKeeper send analytics or other data to 3rd parties? No, but in future releases anonymous analytics could be added. DockLogKeeper never will send any of your containers or logs data.
3. Is it well tested? Now DockLogKeeper is in alpha stage. But yes, testing is in future plans.

## Enhancement
+ Get realtime log updates by websocket
+ Search logs by regexp
+ Highlight <contains> value in search results
+ Automatic past logs loading on scrolling.
+ Make one-click-app for [CapRover](https://caprover.com/)