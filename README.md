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
```

Save this to `$HOME/office683_data/install.zconf` if you are running office683 outside snapcraft

Save this to `$HOME/snap/office683/common/install.zconf` if you are running office683 inside snapcraft

The server is bound to `http://127.0.0.1:8387`

## User Registration

Go to `/register`. You would need the admin_pass.

## Production Config

Setup a HTTPS proxy to `http://127.0.0.1:8387` and disable external access to port '8387'
