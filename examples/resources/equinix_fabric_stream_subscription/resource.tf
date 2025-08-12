resource "equinix_fabric_stream_subscription" "SPLUNK" {
  type        = "STREAM_SUBSCRIPTION"
  name        = "<name>"
  description = "<description>"
  stream_id   = "<stream_id>"
  enabled     = true
  event_selector = {
    include = ["equinix.fabric.connection.*"]
  }
  metric_selector = {
    include = ["equinix.fabric.connection.*"]
  }
  sink = {
    type = "SPLUNK_HEC"
    uri  = "<splunk_uri>"
    settings = {
      event_index  = "<splunk_event_index>"
      metric_index = "<splunk_metric_index>"
      source       = "<splunk_source>"
    }
    credential = {
      type         = "ACCESS_TOKEN"
      access_token = "<splunk_access_token>"
    }
  }
}

resource "equinix_fabric_stream_subscription" "SLACK" {
  type        = "STREAM_SUBSCRIPTION"
  name        = "<name>"
  description = "<description>"
  stream_id   = "<stream_id>"
  enabled     = true
  sink = {
    type = "SLACK"
    uri  = "<slack_uri>"
  }
}

resource "equinix_fabric_stream_subscription" "PAGER_DUTY" {
  type        = "STREAM_SUBSCRIPTION"
  name        = "<name>"
  description = "<description>"
  stream_id   = "<stream_id>"
  enabled     = true
  sink = {
    type = "PAGERDUTY"
    host = "<pager_duty_host"
    settings = {
      change_uri       = "<pager_duty_change_uri>"
      alert_uri        = "<pager_duty_alert_uri>"
    }
    credential = {
      type            = "INTEGRATION_KEY"
      integration_key = "<pager_duty_integration_key>"
    }
  }
}

resource "equinix_fabric_stream_subscription" "DATADOG" {
  type        = "STREAM_SUBSCRIPTION"
  name        = "<name>"
  description = "<description>"
  stream_id   = "<stream_id>"
  enabled     = true
  sink = {
    type = "DATADOG"
    host = "<datadog_host>"
    settings = {
      source          = "Equinix"
      application_key = "<datadog_application_key>"
      event_uri       = "<datadog_event_uri>"
      metric_uri      = "<datadog_metric_uri>"
    }
    credential = {
      type    = "API_KEY"
      api_key = "<datadog_api_key>"
    }
  }
}

resource "equinix_fabric_stream_subscription" "MSTEAMS" {
  type        = "STREAM_SUBSCRIPTION"
  name        = "<name>"
  description = "<description>"
  stream_id   = "<stream_id>"
  enabled     = true
  sink = {
    type = "TEAMS"
    uri  = "<msteams_uri>"
  }
}

resource "equinix_fabric_stream_subscription" "MSTEAMS" {
  type        = "STREAM_SUBSCRIPTION"
  name        = "<name>"
  description = "<description>"
  stream_id   = "<stream_id>"
  enabled     = true
  sink = {
    type = "TEAMS"
    uri  = "<msteams_uri>"
  }
}

resource "equinix_fabric_stream_subscription" "servicenow" {
  type        = "STREAM_SUBSCRIPTION"
  name        = "<name>"
  description = "<description>"
  stream_id   = "<stream_id>"
  enabled     = true
  sink = {
    type = "SERVICENOW"
    host = "<servicenow_host>"
    settings = {
      source = "Equinix"
    }
    credential = {
      type = "USERNAME_PASSWORD"
      username = "<servicenow_username>"
      password = "<servicenow_password>"
    }
  }
}

resource "equinix_fabric_stream_subscription" "webhook" {
  type        = "STREAM_SUBSCRIPTION"
  name        = "<name>"
  description = "<description>"
  stream_id   = "<stream_id>"
  enabled     = true
  sink = {
    type = "WEBHOOK"
    settings = {
      format     = "<webhook_format>"
      event_uri  = "<webhook_event_uri>"
      metric_uri = "<webhook_metric_uri>"
    }
  }
}

resource "equinix_fabric_stream_subscription" "grafana" {
  type        = "STREAM_SUBSCRIPTION"
  name        = "<name>"
  description = "<description>"
  stream_id   = "<stream_id>"
  enabled     = true
  sink = {
    type = "WEBHOOK"
    settings = {
      format     = "OPENTELEMETRY"
      event_uri  = "<grafana_event_uri>"
      metric_uri = "<grafana_metric_uri>"
    }
  }
}
