terraform {
  required_providers {
    fleet = {
      version = "~> 0.0.1"
      source  = "nchaulet.fr/tf/fleet"
    }
  }
}


provider "fleet" {
  kibana_host = "http://localhost:5601"
  username    = "elastic"
  password    = "changeme"
}

resource "fleet_agent_policy" "test_policy" {
  name = "test123"
}

output "test_policy" {
  value = fleet_agent_policy.test_policy
}
