# Slack Emoji Uploader

[![Build Status](https://github.com/Zoraone/slack-emoji-uploader/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/Zoraone/slack-emoji-uploader/actions/workflows/ci.yml)

Simple app to upload emojis into a slack space. Made as a personal project to transfer emojis from one workspace to another.

## Configuration
Requires a user to enter a slack workspace, a slack user token, and a JSON file containing the emojis to upload.

### Slack User Token
Currently, the Slack user token needs to be retrieved manually. The token can be retrieved by logging into a Slack session, opening the browser's developer tools, then running `console.log(window.boot_data.api_token)` in the console. This token should be from the slack workspace which the emojis will be uploaded **to**.

### Slack Emoji JSON File
The JSON file which contains the emojis uses the format which is returned in the response from Slack when calling the [emoji.list](https://api.slack.com/methods/emoji.list) method in the Slack API.

## Uploading Emojis
When an emoji is retrieved, the app will attempt to compare the emojis in the JSON file to the existing emojis that are in the target workspace. Any duplicated will be left off of the list.

Emojis can be uploaded individually once they have been retrieved an displayed on the emoji list.

## Implemented using:
* [Astilectron](https://github.com/asticode/astilectron) for Electron bindings.
* [Materialize CSS](https://materializecss.com/) for UI styling.
* [VueJS](https://vuejs.org/) for frontend JS framework.