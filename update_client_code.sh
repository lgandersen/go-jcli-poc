#!/bin/sh
cd /host/jocker && sudo mix openapi.spec.json --spec Jocker.Engine.API.Spec
cd /host/jcli && oapi-codegen ../jocker/openapi.json > client/jocker_engine_gen.go
