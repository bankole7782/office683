# office683

a group of tools for any office

## Database

The project uses [flaarum](https://github.com/saenuma/flaarum)

## Production Setup

1. Install flaarum on a Ubuntu 20.04 server with the instructions at [sae.ng](https://sae.ng/flaarumtuts/pinstall) and make sure the server has a static external address

1. Point a domain or subdomain to the server's IP. This is necessary for https config.

1. ssh into the server from the step above and get the installation script with `wget https://sae.ng/install_office683.sh`

1. Make the downloaded script executable by running `sudo chmod +x install_office683.sh`

1. Execute the downloaded script by running `sudo ./install_office683.sh`

1. Edit the config at `/var/lib/office683/install.zconf`. All fields must be set.

1. Now start office683 services: `sudo snap restart o6sites` and `sudo snap restart o6ssl`

1. Make sure port `443` and port `80` are open in your server's firewall.

1. You can now view the project from the domain


## User Registration

Go to `/register`. You would need the admin_pass.

`admin_pass` is generated for every installation. You can get it from `/var/lib/office683/install.zconf` on the production server.
