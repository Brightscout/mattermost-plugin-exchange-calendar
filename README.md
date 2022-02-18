# Mattermost Plugin Exchange MS Calendar

## Table of Contents

1. [Overview](#overview)
2. [Features](#features)
3. [Setup](#setup)
4. [Building the plugin](#building-the-plugin)
5. [Installation](#installation)
6. [Configuration](#configuration)
7. [Development](#development)

## Overview

This plugin supports a two-way integration between Mattermost and Microsoft Outlook Calendar. You can follow the instructions to [install](#installation) and [configure](#configuration) the plugin.

## Features

- Daily summary of calendar events.
- Automatic user status synchronization into Mattermost.
- Accept or decline calendar event invites from Mattermost.

## Setup
Make sure you have the following components installed:

  - Go - v1.16 - [Getting Started](https://golang.org/doc/install)
    > **Note:** If you have installed Go to a custom location, make sure the `$GOROOT` variable is set properly. Refer [Installing to a custom location](https://golang.org/doc/install#install).

## Building the plugin
Run the below command in the plugin repo to prepare a compiled, distributable plugin zip:

```bash
$ make dist
```
**Note**: On successful build, a `.tar.gz` file in `/dist` folder will be created that can be uploaded to Mattermost.

## Installation

### Uploading from a Github release
1. Go to the [releases page of this GitHub repository](https://github.com/Brightscout/mattermost-plugin-exchange-mscalendar/releases) and download the latest release for your Mattermost server.
2. Upload this file in the Mattermost **System Console > Plugins > Management** page to install the plugin. To learn more about how to upload a plugin, [see the documentation](https://docs.mattermost.com/administration/plugins.html#plugin-uploads).

### Uploading a local build
- Uploading an existing build present in the `dist` folder
Upload the zip file of the build present in the `dist` folder in the Mattermost **System Console > Plugins > Management** page to install the plugin. To learn more about how to upload a plugin, [see the documentation](https://docs.mattermost.com/administration/plugins.html#plugin-uploads).

## Configuration
- Go to the Microsoft Calendar plugin configuration page on Mattermost as **System Console > Plugins > Microsoft Calendar**

    <img src="https://user-images.githubusercontent.com/72438220/154666704-7f8c0162-4295-4c07-a528-8cf62b598afd.png" />

- On the Microsoft Calendar plugin configuration page you need to add data for the following fields
	- **Exchange Server Base URL**: Base URL of the Exchange server.
    
    <img src="https://user-images.githubusercontent.com/72438220/154667268-16b5cfbd-9250-4117-80a1-d6e460d8e898.png" />

	- **Exchange Server Authentication Key**: Authentication key set in EWS-server for authenticating API requests.
	You can click on the `Regenerate` button to generate a new key and make sure to add this key on the `EWS` server configuration to authenticate all the API calls made to the EWS server by this Mattermost plugin.
	<img src="https://user-images.githubusercontent.com/72438220/154667750-62deda36-3ecd-48b4-80b5-b36774fce3fc.png" />

## Development

To avoid having to manually install your plugin, build and deploy your plugin using one of the following options.

### Deploying with Local Mode

If your Mattermost server is running locally, you can enable [local mode](https://docs.mattermost.com/administration/mmctl-cli-tool.html#local-mode) to streamline deploying your plugin. Edit your server configuration as follows:

```json
{
    "ServiceSettings": {
        ...
        "EnableLocalMode": true,
        "LocalModeSocketLocation": "/var/tmp/mattermost_local.socket"
    }
}
```

and then deploy your plugin:
```
make deploy
```

You may also customize the Unix socket path:
```
export MM_LOCALSOCKETPATH=/var/tmp/alternate_local.socket
make deploy
```

If developing a plugin with a web app, watch for changes and deploy those automatically:
```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=j44acwd8obn78cdcx7koid4jkr
make watch
```

### Deploying with credentials

Alternatively, you can authenticate with the server's API with credentials:
```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_USERNAME=admin
export MM_ADMIN_PASSWORD=password
make deploy
```

or with a [personal access token](https://docs.mattermost.com/developer/personal-access-tokens.html):
```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=j44acwd8obn78cdcx7koid4jkr
make deploy
```
