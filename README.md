# Mattermost Plugin Exchange Calendar

## Table of Contents

- [Mattermost Plugin Exchange Calendar](#mattermost-plugin-exchange-calendar)
  - [Table of Contents](#table-of-contents)
  - [License](#license)
  - [Overview](#overview)
  - [Features](#features)
  - [Setup](#setup)
  - [Building the plugin](#building-the-plugin)
  - [Installation](#installation)
    - [Using a Github release](#using-a-github-release)
    - [Using a local build](#using-a-local-build)
  - [Configuration](#configuration)
  - [Development](#development)
    - [Deploying with Local Mode](#deploying-with-local-mode)
    - [Deploying with credentials](#deploying-with-credentials)

## License

See the [LICENSE](./LICENSE) file for license rights and limitations.

## Overview

This plugin supports a two-way integration between Mattermost and a Microsoft Exchange Server Calendar. For a stable production release, please download the latest version from the Plugin Marketplace and you can follow the instructions to [install](#installation) and [configure](#configuration) the plugin.

**Note:** This plugin only supports the integration with on-premise Microsoft Exchange Server 2016+. For Azure/Office365 support please see [Mattermost Microsoft Calendar Plugin](https://github.com/mattermost/mattermost-plugin-mscalendar).

## Features

- Daily summary of calendar events.
- Automatic user status synchronization into Mattermost.
- Create calendar events from Mattermost.
- Get notifications for new calendar events on Mattermost.
- Accept or decline calendar event invites from Mattermost.

## Setup

Make sure you have the following components installed:

- Go - v1.16 - [Getting Started](https://golang.org/doc/install)
    > **Note:** If you have installed Go to a custom location, make sure the `$GOROOT` variable is set properly. Refer [Installing to a custom location](https://golang.org/doc/install#install).

- Make

## Building the plugin

Run the following command in the plugin repo to prepare a compiled, distributable plugin zip:

```bash
make dist
```

**Note**: After a successful build, a `.tar.gz` file in `/dist` folder will be created which can be uploaded to Mattermost.

## Installation

### Using a Github release

1. Go to the [releases page of this GitHub repository](https://github.com/Brightscout/mattermost-plugin-exchange-calendar/releases) and download the latest release for your Mattermost server.
2. Upload this file in the Mattermost **System Console > Plugins > Management** page to install the plugin. To learn more about how to upload a plugin, [see the documentation](https://docs.mattermost.com/administration/plugins.html#plugin-uploads).

### Using a local build

Upload the zip file created during the build and found in the `dist` folder using the Mattermost **System Console > Plugins > Management** page to install the plugin. To learn more about how to upload a plugin, [see the documentation](https://docs.mattermost.com/administration/plugins.html#plugin-uploads).

## Configuration

- Go to the Microsoft Calendar plugin configuration page on Mattermost as **System Console > Plugins > Microsoft Exchange Calendar**.

    ![image](https://user-images.githubusercontent.com/72438220/154666704-7f8c0162-4295-4c07-a528-8cf62b598afd.png)

- On the Microsoft Exchange Calendar plugin configuration page, you need to enter data for the following fields:
  - **EWS Proxy Server URL**: The base URL of the EWS Proxy Server.
    ![image](https://user-images.githubusercontent.com/72438220/155143980-2a20fe84-6c38-4205-89ba-c36244d50bdb.png)

  - **EWS Proxy Server Authentication Key**: The authentication key used by the [mattermost-plugin-exchange-ews-proxy](https://github.com/Brightscout/mattermost-plugin-exchange-ews-proxy) for authenticating API requests from this plugin.
 You can click on the `Regenerate` button to generate a new key. Ensure that the key is set in the mattermost-plugin-exchange-ews-proxy `.env` file so that the proxy can authenticate all the API calls made by this Mattermost plugin.
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

```bash
make deploy
```

You may also customize the Unix socket path:

```bash
export MM_LOCALSOCKETPATH=/var/tmp/alternate_local.socket
make deploy
```

If developing a plugin with a web app, watch for changes and deploy those automatically:

```bash
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=j44acwd8obn78cdcx7koid4jkr
make watch
```

### Deploying with credentials

Alternatively, you can authenticate with the server's API with credentials:

```bash
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_USERNAME=admin
export MM_ADMIN_PASSWORD=password
make deploy
```

or with a [personal access token](https://docs.mattermost.com/developer/personal-access-tokens.html):

```bashs
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=j44acwd8obn78cdcx7koid4jkr
make deploy
```
