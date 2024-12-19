#!/bin/sh

# Create FTP user with specified username and password
if ! id -u "$FTP_USER" >/dev/null 2>&1; then
    adduser -D -h /home/"$FTP_USER" "$FTP_USER"
    echo "$FTP_USER:$FTP_PASS" | chpasswd
fi

# Add FTP user to vsftpd user list if not already present
if ! grep -q "$FTP_USER" /etc/vsftpd.userlist; then
    echo "$FTP_USER" >> /etc/vsftpd.userlist
fi

# Generate vsftpd configuration with dynamic passive ports
cat <<EOL > /etc/vsftpd/vsftpd.conf
listen=YES
anonymous_enable=NO
local_enable=YES
write_enable=YES
dirmessage_enable=YES
use_localtime=YES
xferlog_enable=YES
connect_from_port_20=YES
chroot_local_user=YES
allow_writeable_chroot=YES
seccomp_sandbox=NO

# Passive mode settings
pasv_enable=YES
pasv_min_port=$PASV_MIN_PORT
pasv_max_port=$PASV_MAX_PORT
pasv_address=$PASV_ADDRESS
pasv_promiscuous=YES

# Logging
xferlog_file=/var/log/vsftpd/xferlog
xferlog_std_format=YES
EOL

# Start vsftpd in the foreground
exec /usr/sbin/vsftpd /etc/vsftpd/vsftpd.conf
