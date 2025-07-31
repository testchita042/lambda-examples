#!/bin/bash

# Load environment variables from customenv.conf
source ~/.config/environment.d/customenv.conf

# Deploy the lambda using chitac CLI
chitac deploy \
  --name mcp-hour \
  --handler Handler \
  --host mcp-hour.chitacloud.com \
  --method "GET,POST" \
  --path "*" \
  --response-content-type "text/event-stream" \
  --ttl 86400 \
  --public-invoke
