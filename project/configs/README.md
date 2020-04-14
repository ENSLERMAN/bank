#### Please create file `config.toml` in this folder and configure your server

````
bind_addr = ":8080"
log_level = "debug"

[store]
database_url = "host=yourhost dbname=yourdbname sslmode=disable password=password user=user"
````