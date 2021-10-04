terraform {
  required_providers {
    fleet = {
      version = "~> 0.0.1"
      source  = "nchaulet.fr/tf/fleet"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = "2.15.0"
    }
  }
}


provider "fleet" {
  kibana_host = "http://localhost:5601"
  username    = "elastic"
  password    = "changeme"
}

resource "fleet_agent_policy" "test_policy" {
  config_json = <<EOL
  {
    "name": "myconfigtest1234568101234567101",
    "namespace": "production",
    "package_policies": [{
      "name": "log-2",
      "description": "",
      "namespace": "default",
      "policy_id": "958e4dd0-22bc-11ec-85e8-c3b96ba33e4a",
      "enabled": true,
      "output_id": "",
      "inputs": [
        {
          "type": "logfile",
          "policy_template": "logs",
          "enabled": true,
          "streams": [
            {
              "enabled": true,
              "data_stream": {
                "type": "logs",
                "dataset": "log.log"
              },
              "vars": {
                "paths": {
                  "type": "text",
                  "value": [
                    "/test.log"
                  ]
                },
                "data_stream.dataset": {
                  "value": "generic",
                  "type": "text"
                },
                "custom": {
                  "value": "",
                  "type": "yaml"
                }
              }
            }
          ]
        }
      ],
      "package": {
        "name": "log",
        "title": "Custom logs111 19",
        "version": "0.5.0"
      }
    }]
  }
  EOL
}

output "test_policy" {
  value     = fleet_agent_policy.test_policy
  sensitive = true
}

resource "docker_container" "foo" {
  image    = "docker.elastic.co/beats/elastic-agent:7.15.0-SNAPSHOT"
  name     = "foo"
  must_run = true
  host {
    host = "elasticsearch"
    ip   = "192.168.65.2"
  }
  host {
    host = "fleet-server"
    ip   = "192.168.65.2"
  }
  env = ["FLEET_ENROLL=1", "FLEET_URL=http://fleet-server:8220", "FLEET_INSECURE=1", "FLEET_ENROLLMENT_TOKEN=${fleet_agent_policy.test_policy.enrollment_token}"]
}

