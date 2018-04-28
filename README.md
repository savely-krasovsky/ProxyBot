# ProxyBot
Ready to deploy SOCKS5 server and Telegram Bot with embedded database to store users credentials. Great for deployment and using in small communities and chats.

# Config
To start using bot you need `_config.yml` file:
```yaml
# Required params
token: telegram_bot_token
addr: fdqn_or_ip

# Optionally you can change port, but 1080 is default
port: 1080

# In private mode, bot will require adding user manually, else it will just register anyone automatically
private: true

# Also you can change limit of maximum users, 100 is default
limit: 100

# Connections limit per user
connsperuser: 10

# Setting admin id will give access to admin commands
adminid: 123456

# Only for development with restricted api.telegram.org
proxy:
    addr: fdqn_or_ip
    port: 1080
    username: user
    password: pass
    
# Verbose logs
verbose: false
```

# Commands
* `/start` - automatically registers new user (if bot is not in private mode) or just shows you actual credentials
* `/redeem {invitation code}` - redeems invitation code and registers user
* `/update {username} {password}` - updates user's creds, without args generates default pair
* `/remove` - totally removes user from database

Admins only:
* `/make_invitation` - makes invitation that can be redeemed only once
