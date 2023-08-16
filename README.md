# chess-server
Server for hosting chess games


## Installation

```bash
go get .
./install_mysql.sh
```


Create a config file
```bash
cat << EOF
DBUSER: root
DBPASS: [PASSWORD MADE ON INSTALL] 
EOF
```

Add `server` database to mysql

```bash
mysql -u root -p
```
Enter password
```sql
CREATE DATABASE server;
```


Note: if running on WSL, you might need to add systemctl. 
Add the following to `/etc/wsl.conf`, then restart 

On WSL:
```conf
[boot]
systemctl=true
```

In powershell:
```powershell
wsl.exe --shutdown
```
