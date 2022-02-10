# Telescope
Telescope is state alerter for Cosmos SDK based chains. 

## Install
You can download latest prebuild binary on [release page](https://github.com/gadost/telescope/releases/latest) or install via go (required >= 1.17.5):
```
go install github.com/gadost/telescope@latest
```

## Configure

```
# fill telescope.toml
telescope init

# fill <chainname>.toml
telescope config generate --name <chainname>

# start telescope service
telescope start
```

example telescope.toml
```toml
[settings]
#  check for new chain releases at github.com or not , default false
github_release_monitor = true  
[telegram]
# telegram as channel for alerts
enabled = true
# create new bot https://t.me/BotFather
token = "1234567:SecRetTokernByBoTfaTher"
# telegram chat_id . You can add bot to channel/group or send alert to DM. 
# Collect chat_id:
# send any message to channel where this bot added , then 
# curl https://api.telegram.org/bot<TOKEN>/getUpdates and find for
# "chat":"id" : "<CHAT_ID>"
chat_id = "-10000000"
[discord]
enabled = true
token = "SecRetTokernByBoTfaTher"
channel_id = 1234567890
[twilio]
# will be implemented in future updates
[mail]
# will be implemented in future updates
[sms]
# will be implemented in future updates
```

example chain config  <chain_name>.toml ( replace <chain_name> with chain name)
```toml
[info]
# enable alerts to telegram if enabled for current chain
telegram = true
# alert when voting power of your validator was changed by specified amount
voting_power_changes = 10
# how many blocks your validator can skip in a row ( means alert every X missed block in a row. if =1  -  every missed block)
blocks_missed_in_a_row = 10
# alert when peers count goes below  
peers_count = 10
# GitHub repository for new release alerts
github = "https://github.com/User/Repo"
[[node]]
# validator or sentry 
role = "validator"
# Node RPC 
rpc = "http://1.2.3.4:26657"
# Enable alerts for specified node
monitoring_enabled = true
[[node]]
role = "sentry"
rpc = "http://2.3.4.5:26657"
monitoring_enabled = true
# Monitoring missed block through this node 
network_monitoring_enabled = true
[[node]]
role = "sentry"
rpc = "http://3.4.5.6:26657"
monitoring_enabled = true
[[node]]
role = "sentry"
rpc = "http://7.8.9.1:26657"
monitoring_enabled = false
```

## Start monitoring

```
telescope start
```
default configdir `$HOME/.telescope/conf.d`
