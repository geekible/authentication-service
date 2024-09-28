## Authentication Server
Configuration
You will need to create a config.yaml file in the root of the application folder as below

auth_service:     
&nbsp;&nbsp;&nbsp;&nbsp;client_id: ""     
&nbsp;&nbsp;&nbsp;&nbsp;client_secret: ""      
service:     
&nbsp;&nbsp;&nbsp;&nbsp;port: 0000     
&nbsp;&nbsp;&nbsp;&nbsp;environment: dev     
&nbsp;&nbsp;&nbsp;&nbsp;log_file_path: "."     
&nbsp;&nbsp;&nbsp;&nbsp;log_file_name: "authentiction_service.log"     
database:     
&nbsp;&nbsp;&nbsp;&nbsp;host: db     
&nbsp;&nbsp;&nbsp;&nbsp;username: postgres     
&nbsp;&nbsp;&nbsp;&nbsp;password:      
&nbsp;&nbsp;&nbsp;&nbsp;dbname: auth_db     
&nbsp;&nbsp;&nbsp;&nbsp;port: 5432     
