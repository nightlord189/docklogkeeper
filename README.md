# docklogkeeper

[DockerHub image](https://hub.docker.com/r/nightlord189/docklogkeeper)

Simple Docker log viewer. DockLogKeeper offers a streamlined approach to viewing Docker logs. It persistently stores logs in a local SQLite database, ensuring both real-time and historical logs are readily accessible.

![Example](https://github.com/nightlord189/docklogkeeper/blob/master/site/screenshot1.png)

## Features
- Persistent storage of your Docker logs.
- Simple search in your logs' historical data.
- Access logs from any Docker container, regardless of its status (active or terminated).
- Simple cookie-based authorization.
- Log rotation: by default, older logs are pruned after a span of 1 week.

## Usecase
When DockLogKeeper could be right choice for you?
+ You manage multiple Docker containers on the same host.
+ Your containers don't generate millions of logs each day.
+ You want to set up log solution by one click and see logs instantly without any scripts/settings/running big enterprise log solution

## How to use?
1. Run command below:
    ```
    docker run --name docklogkeeper --env PASSWORD=YOUR_PASSWORD -d -v /var/run/docker.sock:/var/run/docker.sock -v docklogkeeper:/logs -p 3010:3010 nightlord189/docklogkeeper:1.0.2
    ```
2. Navigate to http://localhost:3010. Authenticate using the username **admin** and the password specified in step 1.
3. Explore and search your logs via an user-friendly interface.
4. Choose a container from those running on your host to view its logs.
5. Use "Update" button to get new logs instantly and "Next" to get older logs when you are scrolling to end.

### Docker-Compose
```
version: "3"
services:
  docklogkeeper:
    container_name: docklogkeeper
    image: nightlord189/docklogkeeper:1.0.2
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - docklogkeeper:/logs
    ports:
      - 3010:3010

volumes:
  docklogkeeper:
```

### [CapRover](https://caprover.com):
To deploy DockLogKeeper as one-click app add following public repo to your CapRover:  
```https://caprover-one-click-apps.nightlord189.vercel.app```

Source code of this repo is [here](https://github.com/nightlord189/caprover-one-click-apps).

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
2. Does DockLogKeeper send analytics or other data to 3rd parties? Not. However, future releases may incorporate anonymous analytics. Rest assured, DockLogKeeper will never transmit your container data or logs.
3. Is it well tested? Currently, DockLogKeeper is in its alpha stage. Comprehensive testing is on our roadmap.

## Enhancement
+ Get realtime log updates via websocket.
+ Search logs by regexp.
+ Automatic past logs loading on scrolling.