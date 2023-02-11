# CLI for Upload local files to Cloud Provider
CLI app to upload files in any readable local directories to Cloud Provider.

Mainly used in conjunction with [cron-backup](https://github.com/mdanialr/cron-backup).
The [cron-backup](https://github.com/mdanialr/cron-backup) is for archiving local files and database dump, then this app is for uploading them to Cloud Provider.

# Features
* Upload any readable local directories as many as possible. (_make sure your cloud provider's capacity is sufficient_).
* Automatically create folders (**that doesn't exist yet**) in cloud provider based on config file.
* Option to set the maximum number of days (retain) delete files in cloud provider that exceed the maximum number of days (retain).
* Option to set the number of worker for upload and or delete job.
* Option to set the chunk size (in byte) of the file when uploading.

# How to Use
1. Download the latest binary file from Releases.
2. Extract the downloaded binary file and make sure it's executable.
    ```bash
    tar -xzf cron-upload....tar.gz
    chmod u+x cron-upload
    ```
3. Create configuration file from the template.
    ```bash
    cp app.yml.example app.yml
    ```
4. Edit the app config file as needed. You can check the template for explanation of each field. 
5. Try to execute and check if there is any error in the app config file.
   ```bash
   ./cron-upload -test
   ```
6. Check the logs file for any error. Maybe failed to upload or delete files, etc.
7. Create a cronjob to run this app. (*optional but recommended*) 
   
    **Example**:
    ```bash
    @daily cd /full/path/to/cron-upload && ./cron-upload -log file
    ```

# Supported Cloud Provider
Currently, support Google Drive & AWS S3 Bucket as the cloud provider. See below how to setup the prerequisite for each provider. 
## Google Drive
1. Create Google Service Account and download the credential in json format. You can follow this awesome [tutorial](https://www.labnol.org/google-api-service-account-220404),
   but following until the [#4](https://www.labnol.org/google-api-service-account-220404#4-share-a-drive-folder) step will be sufficient. Use the shared folder's name as `root` in app config file.
2. Put the full file path where the downloaded credential is to the app config, like so.
    ```yml
    provider:
      name: drive
      cred: /full/path/to/credential.json
    ```
## AWS S3 Bucket
1. Create an administrative user with `AmazonS3FullAccess` permission. see [docs](https://docs.aws.amazon.com/AmazonS3/latest/userguide/setting-up-s3.html#create-an-admin).
2. Create new access key for that admin user and select `Other` option in **Access key best practices & alternatives** step. see [docs](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html#Using_CreateAccessKey).
3. Save the `Access key` & `Secret access key` to json file like so.
    ```json
    {
        "key": "KEY",
        "secret": "SECRET"
    }
    ```
4. Login to console using the admin user from step #1.
5. Create new S3 bucket. see [docs](https://docs.aws.amazon.com/AmazonS3/latest/userguide/creating-bucket.html). 
6. Create new folder and use that folder as `root` in app config file.
7. Put the full file path where the saved json file in step #4 is to the app config, like so.
    ```yml
    provider:
      name: s3
      cred: /full/path/to/credential.json
      region: eu-central-1 # set to where your bucket's region. see step #5
      bucket: chum-bucket # your bucket name. should be created beforehand. see step #5
    ```

# Arguments
* `-path`: set where to find the config file. Default is set to current directory where the binary file is run.
* `-log`: set where to write the log. Default is set to stdout.
  You can change it to `-log file` to write the log to file in the directory that you set in config file.
* `-test`: run all sort of tests such as, validations for the config file, try to create folder, upload & delete files to cloud provider and also check if there is any error.

# License
This project is licensed under the **MIT License** - see the [LICENSE](LICENSE "LICENSE") file for details.
