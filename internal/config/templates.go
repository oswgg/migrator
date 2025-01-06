package config

const MigratorRCFile = `config_folder_path=config/
migrations_folder_path=user_migrations/`

const ConfigFile = `development:
  username: root
  password: root
  database: database
  migrations_table: user_migrations
  host: localhost
  port: 3306
  dialect: mysql
production:
  username: $USERNAME
  password: $PASSWORD
  database: $DATABASE
  migrations_table: user_migrations
  host: $HOST
  port: $PORT
  dialect: $DIALECT
test:
  username: $TEST_USERNAME
  password: $TEST_PASSWORD
  database: $TEST_DATABASE
  migrations_table: user_migrations
  host: $TEST_HOST
  port: $TEST_PORT
  dialect: $TEST_DIALECT
`
