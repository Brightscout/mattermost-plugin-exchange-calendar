# Mattermost Microsoft Exchange Calendar Plugin

## Table of Contents

- [License](#license)
- [Overview](#overview)
- [Features](#features)
- [Prerequisite](#prerequisite)
- [Installation](#installation)
- [Configuration](#configuration)
- [Development](#development)
  - [Setup](#setup)
  - [Building the plugin](#building-the-plugin)
  - [Deploying with Local Mode](#deploying-with-local-mode)
  - [Deploying with credentials](#deploying-with-credentials)

## License

See the [LICENSE](./LICENSE) file for license rights and limitations.

## Overview

This plugin supports a two-way integration between Mattermost and a Microsoft Exchange Calendar. For a stable production release, please download the latest version from the Plugin Marketplace and follow the instructions to [install](#installation) and [configure](#configuration) the plugin.

**Note:** This plugin only supports the integration with on-premise Microsoft Exchange Server 2016+. For Azure/Office365 support please see [Mattermost Microsoft Calendar Plugin](https://github.com/mattermost/mattermost-plugin-mscalendar).

## Features

- Daily summary of calendar events.
- View your calendar events for the week.
- Automatically update your status when you are in a meeting.
- Get notifications for new calendar events on Mattermost.
- Accept or decline calendar event invites from Mattermost.

## Prerequisite

This plugin communicates with Microsoft Exchange Server through a companion service, the [Mattermost Plugin Exchange EWS Proxy](https://github.com/Brightscout/mattermost-plugin-exchange-ews-proxy). Please ensure this service is installed and available to the Mattermost servers running this plugin.

## Installation

1. Go to the [releases page of this GitHub repository](https://github.com/Brightscout/mattermost-plugin-exchange-calendar/releases) and download the latest release for your Mattermost server.
2. Upload this file in the Mattermost **System Console > Plugins > Management** page to install the plugin. To learn more about how to upload a plugin, [see the documentation](https://docs.mattermost.com/administration/plugins.html#plugin-uploads).

## Configuration

- Go to the Microsoft Exchange Calendar plugin configuration page on Mattermost as **System Console > Plugins > Microsoft Exchange Calendar**.

- On the Microsoft Exchange Calendar plugin configuration page, you need to configure the following:
  - **EWS Proxy Server URL**: Enter the base URL of the EWS Proxy Server.

  - **EWS Proxy Server Authentication Key**: The authentication key used by the [mattermost-plugin-exchange-ews-proxy](https://github.com/Brightscout/mattermost-plugin-exchange-ews-proxy) to authenticate API requests from this plugin.
  
 You may use the `Regenerate` button to generate a new key. Ensure that the key is configured in the mattermost-plugin-exchange-ews-proxy's `.env` file so that the proxy can authenticate API requests made by this plugin.

  - **Auto-Connect Users**: When set to 'true', all the users on Mattermost are automatically connected to Exchange via EWS proxy server.

## Development

### Setup

Make sure you have the following components installed:  

- Go - v1.16 - [Getting Started](https://golang.org/doc/install)
    > **Note:** If you have installed Go to a custom location, make sure the `$GOROOT` variable is set properly. Refer [Installing to a custom location](https://golang.org/doc/install#install).

- Make

### Building the plugin

Run the following command in the plugin repo to prepare a compiled, distributable plugin zip:

```bash
make dist
```

After a successful build, a `.tar.gz` file in `/dist` folder will be created which can be uploaded to Mattermost. To avoid having to manually install your plugin, deploy your plugin using one of the following options.

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

```bash
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=j44acwd8obn78cdcx7koid4jkr
make deploy
```
