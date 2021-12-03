# Nest

Nest is a container orchestration tool for a single server.

## Installation

```bash
curl -L https://github.com/wormable/nest/releases/latest/download/nest -o nest
chmod +x nest
sudo mv nest /usr/local/bin
```

## Getting Started

Nest has a single configuration file in `~/.config/nest/nest.json`, you can create one with the following command:

```bash
nest init
```

This command will create a `~/.config/nest/nest.json` file with the following content:

```json
{
  "$schema": "https://raw.githubusercontent.com/redwebcreation/nest/main/globals/schema.json",
  "applications": {},
  "log_policies": [
    {
      "name": "default",
      "rules": [
        {
          "level": "info"
        }
      ]
    }
  ]
}
```