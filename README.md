# Mattermost Plugin Exchange MS Calendar

## Table of Contents

1. [License](#license)
2. [Overview](#overview)
3. [Features](#features)
4. [Setup](#setup)
5. [Building the plugin](#building-the-plugin)
6. [Installation](#installation)
7. [Configuration](#configuration)
8. [Development](#development)

## License

See the [LICENSE](./LICENSE) file for license rights and limitations.

## Overview

This plugin supports a two-way integration between Mattermost and Microsoft Outlook Calendar. For a stable production release, please download the latest version from the Plugin Marketplace and you can follow the instructions to [install](#installation) and [configure](#configuration) the plugin.

**Note:** This plugin supports only on-prem Microsoft Exchange server and not the online Microsoft server.

## Features

- Daily summary of calendar events.
- Automatic user status synchronization into Mattermost.
- Create calendar events from Mattermost.
- Get notifications for new calendar events on Mattermost.
- Accept or decline calendar event invites from Mattermost.

## Setup
Make sure you have the following components installed:

  - Go - v1.16 - [Getting Started](https://golang.org/doc/install)

**Note:** If you have installed Go to a custom location, make sure the `$GOROOT` variable is set properly. Refer [Installing to a custom location](https://golang.org/doc/install#install).

## Building the plugin
Run the below command in the plugin repo to prepare a compiled, distributable plugin zip:

```bash
$ make dist
```
**Note**: On successful build, a `.tar.gz` file in `/dist` folder will be created that can be uploaded to Mattermost.

## Installation

### Using a Github release
1. Go to the [releases page of this GitHub repository](https://github.com/Brightscout/mattermost-plugin-exchange-mscalendar/releases) and download the latest release for your Mattermost server.
2. Upload this file in the Mattermost **System Console > Plugins > Management** page to install the plugin. To learn more about how to upload a plugin, [see the documentation](https://docs.mattermost.com/administration/plugins.html#plugin-uploads).

### Using a local build
- Uploading an existing build present in the `dist` folder
Upload the zip file of the build present in the `dist` folder in the Mattermost **System Console > Plugins > Management** page to install the plugin. To learn more about how to upload a plugin, [see the documentation](https://docs.mattermost.com/administration/plugins.html#plugin-uploads).

## Configuration
- Go to the Microsoft Calendar plugin configuration page on Mattermost as **System Console > Plugins > Microsoft Calendar**.

    ![image](https://user-images.githubusercontent.com/72438220/154666704-7f8c0162-4295-4c07-a528-8cf62b598afd.png)

- On the Microsoft Calendar plugin configuration page, you need to add data for the following fields:
	- **Exchange EWS Proxy Server URL**: Base URL of the Exchange EWS Proxy Server.
    ![image](https://user-images.githubusercontent.com/72438220/155143980-2a20fe84-6c38-4205-89ba-c36244d50bdb.png)

	- **Exchange EWS Proxy Server Authentication Key**: Authentication key used by mattermost-plugin-exchange-ews-proxy for authenticating API requests.
	You can click on the `Regenerate` button to generate a new key and make sure to add this key on the mattermost-plugin-exchange-ews-proxy `.env` file to authenticate all the API calls made to the EWS server by this Mattermost plugin.
	![image](https://user-images.githubusercontent.com/72438220/155144336-2f98f3b4-553c-4827-9e1c-747775004fa3.png)

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
