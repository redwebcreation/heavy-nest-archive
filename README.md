# Nest

Nest is a container orchestration tool for a single server.

## Installation
```bash
curl -L https://github.com/wormable/nest/releases/latest/download/nest -o nest
chmod +x nest
sudo mv nest /usr/local/bin
```

## Getting Started
 
Nest has a single configuration file in `/etc/nest/nest.json`, you can create one with the following command:
```bash
nest init
```

This command will create a `/etc/nest/nest.json` file with the following content:
```json
{
  
}
```