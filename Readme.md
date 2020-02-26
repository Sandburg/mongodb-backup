# mongodb-backup

[![Docker Pulls](https://img.shields.io/docker/pulls/sandburg2011/mongodb-backup)](https://hub.docker.com/r/sandburg2011/mongodb-backup)

mongodb-backup is a simple MongoDB backup automation tool built with Go.
It can be run locally, as a cron job or within a k8s cluster.

#### Features

* dump mongodb to backup archive
* upload to gcloud storage

#### Install

mongodb-backup is available on Docker Hub at [sandburg2011/mongodb-backup](https://hub.docker.com/r/sandburg2011/mongodb-backup).

Supported tags:

* `sandburg2011/mongodb-backup:latest` latest stable release

#### Configuration

mongodb-backup is configurable via environment variables.
supported variables are:

 ENV_NAME| Description | Default
---------|-------------|----------
`GCS_BUCKET` | the bucketname to upload to | -
`GCS_KEY_FILE_PATH` | path to the GCS credential file | -
`MONGODB_HOST` | host of the database to backup | `localhost`
`MONGODB_PORT` | port of the database to backup | `27017`
`MONGODB_USER` | username credential for database | -
`MONGODB_PASSWORD` | password credential for database | -
`MAX_BACKUP_ITEMS` | maximum backup archives, deletes older | `10`
`BACKUP_DIR` | directory where the archive is beeing saved to | `/tmp`

#### Run

    docker run -d --env-file ./.env sandburg2011/mongodb-backup
    docker exec -it CONTAINER_ID ./backup

#### Restore

To restore the data you can download the archive file from the bucket and use mongorestore. In this example a file named backup.gz gets restored.

    mongorestore --gzip --archive=backup.gz --drop

If you want to restore just a specific database, use the nsInclude argument.

    mongorestore --gzip --archive=backup.gz --nsInclude="DATABASE_NAME.*" --drop

## Examples

##### Kubernetes Cron

Example for kubernetes Cron to backup the databse every hour and keep backups for one day.

```apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: backup-db
  namespace: graphql
spec:
  schedule: "0 */1 * * *"
  jobTemplate:
    spec:
      template:
        spec:
          volumes:
          - name: google-credentials
            secret:
              secretName: backup-service-credentials
          containers:
          - name: backup-db
            image: 'sandburg2011/mongodb-backup'
            args:
            - /backup
            volumeMounts:
            - name: google-credentials
              mountPath: /var/secrets/certs
            env:
            - name: GCS_BUCKET
              value: "backup_db"
            - name: GCS_KEY_FILE_PATH
              value: "/var/secrets/certs/credentials.json"
            - name: MONGODB_HOST
              value: "mongodb-client"
            - name: MONGODB_PORT
              value: "27017"
            - name: MAX_BACKUP_ITEMS
              value: "24"
          restartPolicy: OnFailure
```