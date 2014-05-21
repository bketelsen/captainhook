# Upstart

### Create the captainhook dirs

```
mkdir -p /opt/captainhook/bin
```

### Add the binary

```
cp captainhook /opt/captainhook/bin/captainhook
chmod +x /opt/captainhook/bin/captainhook
```

### Add the upstart config

```
cp captainhook.conf /etc/init/captainhook.conf
```

### Start the service

```
service captainhook start
```

### Check the logs

```
tail /var/log/upstart/captainhook.log
```
