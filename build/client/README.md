# Installation

Make sure you are on correct platform branch and the branch is up-to-date

```shell
#For Windows
git checkout windows
git pull

#For Linux
git checkout linux
git pull
```

# Configuration

- Ask the owner which IP and port the server is running on. You'll have to enter them each time you want to connect
- Ask for your client.pem and client.key. Create a folder at the root of the project called "cert" and put these two
  files in it.
- Ask for the key that used for encrypt messages. Put this file into the cert folder
