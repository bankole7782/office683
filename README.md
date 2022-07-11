# office683

A group of tools for any office.

The developers of this program have no access to your data. Your data would be stored by you on your own servers.

## Currently Included Programs

1. **Cabinet**: A store of files. You can upload files created from your local programs here

1. **Documents**: A documents editor. It uses markdown. Its not like Microsoft Word.

1. **Events**: A Record of Events. Prepares workers for future events and keeps memories of past events.


## Database

This project uses [flaarum](https://github.com/saenuma/flaarum)

## Production Setup

1. Install flaarum on a Ubuntu server with the instructions at [sae.ng](https://sae.ng/flaarumtuts/pinstall) and make sure the server has a static external address

1. Point a domain or subdomain to the server's IP. This is necessary for https config.

1. ssh into the server from the step above and get the installation script with `wget https://sae.ng/install_office683.sh`

1. Make the downloaded script executable by running `sudo chmod +x install_office683.sh`

1. Execute the downloaded script by running `sudo ./install_office683.sh`

1. Edit the config at `/var/lib/office683/install.zconf`. All fields must be set.

1. Make sure port `443` and port `80` are open in your server's firewall.

1. Now start office683 services: `sudo systemctl start o6sites` and `sudo systemctl start o6ssl`

1. You can now view the project from the domain


## User Registration

Go to `/register`. You would need the admin_pass.

`admin_pass` is generated for every installation. You can get it from `/var/lib/office683/install.zconf` on the production server.
