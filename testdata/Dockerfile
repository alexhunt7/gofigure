FROM alpine:3.8
RUN apk add --no-cache openssh && \
    rm -f /etc/ssh/*_key
RUN addgroup testuser && \
    adduser -D --gecos "" -G testuser -s /bin/sh testuser && \
    echo 'testuser:password' | chpasswd && \
    mkdir /home/testuser/.ssh
COPY ./ssh_host_rsa_key /etc/ssh/ssh_host_rsa_key
COPY ./ssh_host_rsa_key.pub /home/testuser/.ssh/authorized_keys
RUN chmod u+rwX,go-rwx -R /home/testuser/.ssh /etc/ssh/ssh_host_rsa_key && \
    chown -R testuser:testuser /home/testuser/.ssh
EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]
