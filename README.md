# Cronjob App for Upload to Cloud Provider
Little app that upload files in any readable local directories to folders in Cloud Provider. Mainly used conjunction
with [cron-backup](https://github.com/mdanialr/cron-backup).
App [cron-backup](https://github.com/mdanialr/cron-backup) for archiving local files and database, then this app
for uploading them to Cloud Provider.

# Features
* Upload any readable local directories as many as possible. (_make sure your cloud provider's capacity is sufficient_).
* Automatically create folders (**that doesn't exist yet**) in cloud provider based on config file.
* Option to delete files in cloud provider that exceed the maximum number of days.

# How to Use
1. Download the latest binary file from Releases.
2. Make directory where the download and extracted binary file will reside. We will use bin directory as example.
```bash
mkdir bin
```
3. Make sure the binary file is executable.
```bash
chmod u+x bin/cron-upload
```
4. Create configuration file.
```bash
touch app-config.yml
```
5. Fill in config file as needed. You can check app-config.yml.example in this repo for reference.
6. Prepare required files. See below.
7. Execute the binary file from directory where this config file exist, otherwise you will get error config file is not found.
8. Execute with `-refresh` params first to get refresh token then, exchange it with access token with `-init` params.
```bash
./bin/cron-upload -refresh -drive
./bin/cron-upload -init -drive
./bin/cron-upload -drive
```
8. Check logs file for any error. Maybe required fields are empty, etc.
9. (optional but recommended) Create a cronjob to run this app.
> Example
```bash
@daily cd /full/path/to/cron-upload && ./bin/cron-upload -drive
```

# Prepare Required Files (Google Drive)
1. Create OAuth client with 'Desktop Client'.
2. Download credential.json file.
3. Write the path where credential.json file reside to app-config in `provider.cred` segment.
4. Make sure credential.json file **readable** & **accessible** by this app.
5. You're good to go.

# Arguments
* `-refresh`: if used with provider argument (e.g. `-drive`), renew or init refresh token.
* `-init`: if used with provider argument (e.g. `-drive`), retrieve token for authentication against Google Drive provider.
* `-drive`: do the upload job using Google Drive provider.

# Under the Hood
1. `./bin/cron-upload -refresh -drive`. this will exchange credential.json file for authorization code and create json file that defined in app-config file that contain refresh token before exchange it for access token.
2. `./bin/cron-upload -init -drive`. this will exchange refresh token for access token and create ...token.json file that contain access token.
3. `./bin/cron-upload -drive`. this will do the upload job sequentially and automatically renew access token in ...token.json file if expired.