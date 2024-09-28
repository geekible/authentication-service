## Authentication Server
#Configuration
You will need to create a config.yaml file in the root of the application folder as below

auth_service:
  client_id: ""
  client_secret: ""
service:
  port: 0000
  environment: dev
  log_file_path: "."
  log_file_name: "authentiction_service.log"
database:
  host: db
  username: postgres
  password: !
  dbname: auth_db
  port: 5432
