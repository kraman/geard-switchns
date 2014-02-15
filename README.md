geard-switchns
==============

Utility command which switches into a specified docker container's namespace and execute a command.
It allows two use-cases:

* Admin commands

Usage:

    geard-switchns <docker container name> <command> [args]...
        
Allows a user with CAP_SYS_ADMIN capability to switch into a specified docker container and execute a command.
Typical use for this would be to run admin commands within a container.

* User SSH

Add geard-switchns as a command to the ```.authorized_keys``` file. Eg:

    command="/usr/sbin/geard-switchns" ssh-rsa AAAA...== user@host

When the user SSH's into the host machine, SSH runs ```geard-switchns```. The utility then looks up a docker container
with the same name as the username and starts a bash shell within the container.

License
=======

Apache Software License (ASL) 2.0.