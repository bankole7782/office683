# office683

a group of tools for any office

## Database

The project uses [flaarum](https://github.com/saenuma/flaarum)

## Development Environment Config

```
// name of the company that created this office tools information
company_name: Test1

// logo of the company on your website.
company_logo: https://sae.ng/static/logo.png

// admin_pass is the password used by all admins
admin_pass: MClhRqmAUYMJndr0m3R3sAslKHrNimFIRkuTaicx02lh00Js1Z

// admin_email is for contacting the admin to get access
admin_email: admin@admin.com

// flaarum_keystr is the key used in connecting to the flaarum server.
// you must set this after running this program.
// you can get it by sshing into your server and running 'flaarum.prod r'
flaarum_keystr: not-yet-set


// domain must be set after launching your server
domain:
```

Save this to `$HOME/office683_data/install.zconf` if you are running office683 outside snapcraft

Save this to `/var/snap/office683/common/install.zconf` if you are running office683 inside snapcraft

The server is bound to `http://127.0.0.1:8387`


## Production Setup

1. Install office683 `sudo snap install office683`

1. Generate a config using `office683 init` and fill

1. Create a service account in Google Cloud and store the downloaded json in your office683 folder ($HOME/snap/office683/common/)

1. Launch your server using `office683 lh {conf} {json}` where conf is generated in step 2 and {json} is gotten from step 3

1. Point a domain to the server's IP. This is necessary for https config.

1. SSH into your server and run `flaarum.prod r` to get your flaarum key

1. Edit the config at `/var/snap/office683/common/install.zconf` and add your flaarum key, domain in their expected slots

1. Now restart office683 services: `sudo snap restart office683.sites` and `sudo snap restart office683.ssl`

1. You can now view the project from the domain


## User Registration

Go to `/register`. You would need the admin_pass.

`admin_pass` is generated for every installation. You can get it from `/var/snap/office683/common/install.zconf` on the production server.
