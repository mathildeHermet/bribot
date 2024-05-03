# bribot

## Description

This bot aims to sends on discord channel configured via `--webhook-url` the info relatives to a CTF.

## Usage

`bribot` only has `reminder` command for now.

```bash
➜  bribot git:(main) ✗ ./bribot --help         
Manage discord bot for CTF annoncements.

Usage:
  bribot [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  reminder    Remind CTF url, timelines adn credentials.

Flags:
  -d, --dry-run   Do not perform the operation, only show what would be sent.
  -h, --help      help for bribot
  -v, --verbose   Make the operation more talkative.
      --version   version for bribot

Use "bribot [command] --help" for more information about a command.
```

`bribot reminder` flags to configure

```bash
➜  bribot git:(fix_local_time) ✗ bribot reminder --help                       
Enroll in a specific CTF challenge.

Usage:
  bribot reminder [flags]

Examples:
bribot reminder --log-level debug --webhook-url https://discord.com/api/webhooks/id/pwd --ctf-name punkctf--ctf-url https://ctf.example.com --ctf-date-start 2024-01-01 00:00:00 UTC --ctf-date-end 2024-02-01 00:00:00 UTC --ctf-date-tz UTC --team-name teamName --team-password teamPassword

Flags:
      --ctf-date-end string     Specify the end date. Format should be 'YYYY-MM-DD HH:MM:SS'.
      --ctf-date-start string   Specify the start date. Format should be 'YYYY-MM-DD HH:MM:SS'. 
      --ctf-date-tz string      Specify the timezone of date in input. EX: 'Europe/London'. (default "Europe/Paris")
      --ctf-name string         Specify the CTF name
      --ctf-url string          Specify the CTF url
  -h, --help                    help for reminder
      --log-level string        Specify the log level (debug, info, warn, error) (default "info")
      --team-name string        Specify the team name for this CTF.
      --team-password string    Specify the team password for this CTF.
      --webhook-url string      Specify the webhook url

Global Flags:
  -d, --dry-run   Do not perform the operation, only show what would be sent.
  -v, --verbose   Make the operation more talkative.
```
