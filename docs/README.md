# appjet

## High Level Architecture Diagram
<img width="783" alt="image" src="https://github.com/user-attachments/assets/41e3d49b-ef43-416f-9f07-0100ad91f9fa" />

## How to Run

### Server Side

- Dependency: Docker installed
- Steps:
      - 1) Docker already installed
      - 2) Open terminal: run docker compose up in the same folder as you placed the docker-compose.yaml file present in this projects directory: https://github.com/guilhermemalhado1/appjet/tree/master/server/etc/infrastructure
      - 3) Double click in the appjet-server executable

### Client Side

- Dependency: None
- Configuration: Fill the Config.json file
                 <img width="565" alt="image" src="https://github.com/user-attachments/assets/2d9882f4-59b8-4373-bff2-d5befd1bfa8c" />

  Note:. The config.json file, must be place in the same directory as your appjet-cli executable.
- Steps:

  1) Open terminal on the directory of the appjet-cli executable
  2) Run one of the following commands: 
<img width="544" alt="image" src="https://github.com/user-attachments/assets/54405162-77f2-4315-a8e5-bbeebc1b2981" />

## Metrics And Monitoring

### Third Party tools used

#### Grafana & InfluxDB

There is the optional feature of colecting logs and metrics and displaying that on grafana.
In order to do that you must only go on the server URL where appjet-decision-manager is running/3001 and login with grafana credentials. (They are set on your .env file within infrastructure folder and you can change the credentials as you wish).

You can also create your own dashboards and place them on the following directory: https://github.com/guilhermemalhado1/appjet/tree/master/server/etc/infrastructure/dashboards 
Then you can simply upload them into your grafana instance in the desired server.


