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
6. Execute the binary file from directory where this config file exist, otherwise you will get error config file is not found.
7. Execute _-init_ params first to get token.json for authentication for Google Drive provider.
```bash
./bin/cron-upload -init -drive
./bin/cron-upload -drive
```
8. (optional but recommended) Create a cronjob to run this app.
> Example
```bash
@daily cd /full/path/to/cron-upload && ./bin/cron-upload -drive
```

# Arguments
* `-init`: if used with provider argument (e.g. `-drive`), retrieve token for authentication against Google Drive provider.
* `-drive`: do the upload job using Google Drive provider.

# Notes
The most important thing for using Google Drive provider is their Oauth2 credentials, so you will need to fill in auth.json
file that consists of four keys then add its path to the config file:
* "refresh": refresh key. You will get this when exchanging authorization code. usually has `1//` as its first characters.
* "client_id": Client ID for this credential.
* "client_secret": Client secret for this credential.
* "token_uri": You will get this in json response when exchanging authorization code for token. usually `https://oauth2.googleapis.com/token`.

How to get these values you can head over to [this](https://stackoverflow.com/questions/19766912/how-do-i-authorise-an-app-web-or-installed-without-user-intervention) link.

As long as you could prepare these auth.json file with their values then you are good to go and use this app.
