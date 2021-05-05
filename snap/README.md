# EdgeX Camera Device Service Snap
[![snap store badge](https://raw.githubusercontent.com/snapcore/snap-store-badges/master/EN/%5BEN%5D-snap-store-black-uneditable.png)](https://snapcraft.io/edgex-device-camera)

This folder contains snap packaging for the EdgeX Camera Device Service Snap

The snap currently supports both `amd64` and `arm64` platforms.


## Snap configuration

Device services implement a service dependency check on startup which ensures that all of the runtime dependencies of a particular service are met before the service transitions to active state.

Snapd doesn't support orchestration between services in different snaps. It is therefore possible on a reboot for a device service to come up faster than all of the required services running in the main edgexfoundry snap. If this happens, it's possible that the device service repeatedly fails startup, and if it exceeds the systemd default limits, then it might be left in a failed state. This situation might be more likely on constrained hardware (e.g. RPi).

This snap therefore implements a basic retry loop with a maximum duration and sleep interval. If the dependent services are not available, the service sleeps for the defined interval (default: 1s) and then tries again up to a maximum duration (default: 60s). These values can be overridden with the following commands:
    
To change the maximum duration, use the following command:

```bash
$ sudo snap set edgex-device-camera startup-duration=60
```

To change the interval between retries, use the following command:

```bash
$ sudo snap set edgex-device-camera startup-interval=1
```

The service can then be started as follows. The "--enable" option
ensures that as well as starting the service now, it will be automatically started on boot:

```bash
$ sudo snap start --enable edgex-device-camera.device-camera-go
```

### Using a content interface to set device configuration

The `device-config` content interface allows another snap to seed this device
snap with both a configuration file and one or more device profiles. 


To use, create a new snap with a directory containing the configuration and device profile files. Your snapcraft.yaml file then needs to define a slot with read access to the directory you are sharing.

```
slots:
  device-config:
    interface: content  
    content: device-config
    write: 
      - $SNAP/config
```

where `$SNAP/config` is configuration directory your snap is providing to the device snap.

Then connect the plug in the device snap to the slot in your snap,
which will replace the configuration in the device snap. Do this with:

```bash
$ sudo snap connect edgex-device-camera:device-config your-snap:device-config
```

This needs to be done before the device service is started for the first time. Once you have set the configuration the device service can be started and it will then be configurated using the settings you provided:

```bash
$ sudo snap start edgex-device-camera.device-camera-go
```

**Note** - content interfaces from snaps installed from the Snap Store that have the same publisher connect automatically. For more information on snap content interfaces please refer to the snapcraft.io [Content Interface](https://snapcraft.io/docs/content-interface) documentation.
