## Overview

Portkey is a tool lets us securely connect to services on remote servers, even if those services are not directly accessible from your location. It uses SSH to create a secure pathway (or tunnel) between your local machine and the remote service, making it feel like the service is running locally on your own computer

![portkey](/img/portkey.jpeg)


### Adding services to the JSON

```json
    "service_name": [
    "IP",
    "connection_type", 
    "port_number"
]
```
Add the IP of server, if we are trying to access the service remotely, configure a proxyjump. If everything is hosted on a single machine, you can just use the public IP in this section. Connection type refers to type of socket - unix or TCP