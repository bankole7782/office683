// lgcp is the launcher program for Google Cloud platform
package main

import (
	"fmt"
	"context"
	"google.golang.org/api/option"
	compute "google.golang.org/api/compute/v1"
	"os"
	"path/filepath"
	"github.com/gookit/color"
	"github.com/bankole7782/office683/office683_shared"
	"time"
	"strings"
	"strconv"
	"github.com/saenuma/zazabul"
	// "google.golang.org/api/option"
)


func waitForOperationRegion(project, region string, service *compute.Service, op *compute.Operation) error {
	ctx := context.Background()
	for {
		result, err := service.RegionOperations.Get(project, region, op.Name).Context(ctx).Do()
		if err != nil {
			return fmt.Errorf("Failed retriving operation status: %s", err)
		}

		if result.Status == "DONE" {
			if result.Error != nil {
				var errors []string
				for _, e := range result.Error.Errors {
					errors = append(errors, e.Message)
				}
				return fmt.Errorf("Operation failed with error(s): %s", strings.Join(errors, ", "))
			}
			break
		}
		time.Sleep(time.Second)
	}
	return nil
}


func waitForOperationZone(project, zone string, service *compute.Service, op *compute.Operation) error {
	ctx := context.Background()
	for {
		result, err := service.ZoneOperations.Get(project, zone, op.Name).Context(ctx).Do()
		if err != nil {
			return fmt.Errorf("Failed retriving operation status: %s", err)
		}

		if result.Status == "DONE" {
			if result.Error != nil {
				var errors []string
				for _, e := range result.Error.Errors {
					errors = append(errors, e.Message)
				}
				return fmt.Errorf("Operation failed with error(s): %s", strings.Join(errors, ", "))
			}
			break
		}
		time.Sleep(time.Second)
	}
	return nil
}


func main() {
	if len(os.Args) < 2 {
		color.Red.Println("Expecting a command. Run with help subcommand to view help.")
		os.Exit(1)
	}

	switch os.Args[1] {
  case "--help", "help", "h":
    fmt.Println(`lgcp creates and configures a office683 server on Google Cloud.

Supported Commands:

    init     Creates a config file for basic deployment (non-autoscaling). Edit to your own requirements.
             Some of the values can be gotten from Google Cloud's documentation.

    lh       Launches a configured instance based on the config created above. It expects a launch file (created from 'init' above)
             and a service account credentials file (gotten from Google Cloud).

      `)

  case "init":
    configFileName := "o683_" + time.Now().Format("20060102T150405") + ".zconf"

    userPath, err := office683_shared.GetUserPath()
    if err != nil {
    	panic(err)
    }

		writePath := filepath.Join(userPath, configFileName)

    var	tmpl = `// project is the Google Cloud Project name
// It can be created either from the Google Cloud Console or from the gcloud command
project:

// region is the Google Cloud Region name
// Specify the region you want to launch your office683 server in.
region:


// zone is the Google Cloud Zone which must be derived from the region above.
// for instance a region could be 'us-central1' and the zone could be 'us-central1-a'
zone:

// disk_size is the size of the root disk of the server. The data created is also stored in the root disk.
// It is measured in Gigabytes (GB) and a number is expected.
// 10 is the minimum.
disk_size: 20

// machine_type is the type of machine configuration to use to launch your office683 server.
// You must get this value from the Google Cloud Compute documentation if not it would fail.
// It is not necessary it must be an e2 instance.
machine_type: e2-highcpu-2

// name of the company that created this office tools information
company_name: Test1

// logo of the company on your website.
company_logo: https://sae.ng/static/logo.png

// admin_pass is the password used by all admins
admin_pass:

// admin_email is for contacting the admin to get access
admin_email: admin@admin.com

// flaarum_keystr is the key used in connecting to the flaarum server.
// you must set this after running this program.
// you can get it by sshing into your server and running 'flaarum.prod r'
flaarum_keystr: not-yet-set


// domain must be set after launching your server
domain:

`

    conf, err := zazabul.ParseConfig(tmpl)
    if err != nil {
    	panic(err)
    }

		conf.Update(map[string]string {
			"admin_pass": office683_shared.UntestedRandomString(50),
		})
    err = conf.Write(writePath)
    if err != nil {
    	panic(err)
    }

    fmt.Printf("Edit the file at '%s' before launching.\n", writePath)

  case "lh":
  	if len(os.Args) != 4 {
  		color.Red.Println("The lh command expects a launch file and a service account credentials file.")
  		os.Exit(1)
  	}

		userPath, err := office683_shared.GetUserPath()
    if err != nil {
    	panic(err)
    }

		inputPath := filepath.Join(userPath, os.Args[2])
  	conf, err := zazabul.LoadConfigFile(inputPath)
  	if err != nil {
  		panic(err)
  	}

  	for _, item := range conf.Items {
  		if item.Value == "" {
				if item.Name != "flaarum_keystr" && item.Name != "domain" {
					color.Red.Println("Every field in the launch file is compulsory.")
					os.Exit(1)
				}
  		}
  	}

		credentialsFilePath := filepath.Join(userPath, os.Args[3])

		instanceName := fmt.Sprintf("o683-%s", strings.ToLower(office683_shared.UntestedRandomString(4)))
		diskName := fmt.Sprintf("%s-disk", instanceName)
  	dataDiskName := fmt.Sprintf("%s-ddisk", instanceName)

  	diskSizeInt, err := strconv.ParseInt(conf.Get("disk_size"), 10, 64)
  	if err != nil {
  		color.Red.Println("The 'disk_size' variable must be a number greater or equal to 10")
  		os.Exit(1)
  	}

		rawInstallZconf, err := os.ReadFile(inputPath)
		if err != nil {
			panic(err)
		}
		var startupScript = `
#! /bin/bash
sudo apt update
sudo apt install nano

sudo snap install flaarum
sudo snap start flaarum.store
sudo flaarum.prod mpr

DATA_BTRFS=/var/snap/flaarum/common/data_btrfs
if  [ ! -d "$DATA_BTRFS" ]; then
	sudo mkfs.btrfs /dev/sdb
	sudo mkdir -p $DATA_BTRFS
fi
sudo mount -o discard,defaults /dev/sdb $DATA_BTRFS
sudo chmod a+w $DATA_BTRFS

sudo snap restart flaarum.store
sudo snap start flaarum.tindexer
sudo snap stop --disable flaarum.statsr

gcloud compute firewall-rules create o683https --direction ingress \
 --source-ranges 0.0.0.0/0 --rules tcp:443 --action allow

gcloud compute firewall-rules create o683http --direction ingress \
 --source-ranges 0.0.0.0/0 --rules tcp:80 --action allow

sudo snap install office683 --edge

`
		startupScript += "cat <<EOT > /var/snap/office683/common/install.zconf\n\n"
		startupScript += string(rawInstallZconf)
		startupScript += "EOT"

  	ctx := context.Background()

		computeService, err := compute.NewService(ctx, option.WithCredentialsFile(credentialsFilePath),
			option.WithScopes(compute.ComputeScope))
		if err != nil {
			panic(err)
		}

		op, err := computeService.Addresses.Insert(conf.Get("project"), conf.Get("region"), &compute.Address{
			AddressType: "EXTERNAL",
			Description: "IP address for a office683 server (" + instanceName + ").",
			Name: instanceName + "-ip",
		}).Context(ctx).Do()

		err = waitForOperationRegion(conf.Get("project"), conf.Get("region"), computeService, op)
		if err != nil {
			panic(err)
		}

		computeAddr, err := computeService.Addresses.Get(conf.Get("project"), conf.Get("region"), instanceName + "-ip").Context(ctx).Do()
		if err != nil {
			panic(err)
		}

		fmt.Println("Office683 server address: ", computeAddr.Address)

		prefix := "https://www.googleapis.com/compute/v1/projects/" + conf.Get("project")

		image, err := computeService.Images.GetFromFamily("ubuntu-os-cloud", "ubuntu-minimal-2004-lts").Context(ctx).Do()
		if err != nil {
			panic(err)
		}
		imageURL := image.SelfLink

		op, err = computeService.Disks.Insert(conf.Get("project"), conf.Get("zone"), &compute.Disk{
			Description: "Data disk for a office683 server (" + instanceName + ").",
			SizeGb: diskSizeInt,
			Type: prefix + "/zones/" + conf.Get("zone") + "/diskTypes/pd-ssd",
			Name: dataDiskName,
		}).Context(ctx).Do()
		err = waitForOperationZone(conf.Get("project"), conf.Get("zone"), computeService, op)
		if err != nil {
			panic(err)
		}

		instance := &compute.Instance{
			Name: instanceName,
			Description: "office683 instance",
			MachineType: prefix + "/zones/" + conf.Get("zone") + "/machineTypes/" + conf.Get("machine_type"),
			Disks: []*compute.AttachedDisk{
				{
					AutoDelete: true,
					Boot:       true,
					Type:       "PERSISTENT",

					InitializeParams: &compute.AttachedDiskInitializeParams{
						DiskName:    diskName,
						SourceImage: imageURL,
						DiskType: prefix + "/zones/" + conf.Get("zone") + "/diskTypes/pd-ssd",
						DiskSizeGb: 10,
					},
				},
				{
					AutoDelete: false,
					Boot:       false,
					Type:       "PERSISTENT",

					InitializeParams: &compute.AttachedDiskInitializeParams{
						DiskName:    dataDiskName,
						DiskType: prefix + "/zones/" + conf.Get("zone") + "/diskTypes/pd-ssd",
						DiskSizeGb: diskSizeInt,
					},
				},
			},
			NetworkInterfaces: []*compute.NetworkInterface{
				{
					AccessConfigs: []*compute.AccessConfig{
						{
							Type: "ONE_TO_ONE_NAT",
							Name: "External NAT",
							NatIP: computeAddr.Address,
						},
					},
				},
			},
			ServiceAccounts: []*compute.ServiceAccount{
				{
					Email: "default",
					Scopes: []string{
						compute.DevstorageFullControlScope,
						compute.ComputeScope,
					},
				},
			},
			Metadata: &compute.Metadata {
				Items: []*compute.MetadataItems {
					{
						Key: "startup-script",
						Value: &startupScript,
					},
				},
			},
			Tags: &compute.Tags {
				Items: []string{"https-server",},
			},
		}

		op, err = computeService.Instances.Insert(conf.Get("project"), conf.Get("zone"), instance).Do()
		if err != nil {
			panic(err)
		}
		err = waitForOperationZone(conf.Get("project"), conf.Get("zone"), computeService, op)
		if err != nil {
			panic(err)
		}

		fmt.Println("o683 Server Name: " + instanceName)


	default:
		color.Red.Println("Unexpected command. Run the lgcp with --help to find out the supported commands.")
		os.Exit(1)
  }

}
