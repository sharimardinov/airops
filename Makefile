SHELL := /bin/bash

BASE_URL ?= http://localhost:8080

K6 := k6
K6_DIR := test/k6

# Даты для stats/routes (обязательные)
FROM ?= 2025-12-12
TO   ?= 2026-01-11

# load_mix
RPS ?= 300
RPS_STATS ?= 20
DURATION ?= 60s
PRE_VUS ?= 200
MAX_VUS ?= 2000
PRE_VUS_STATS ?= 50
MAX_VUS_STATS ?= 500

# ceiling_one
START_RPS ?= 200
R1 ?= 1000
R2 ?= 3000
R3 ?= 5000
CEIL_PRE_VUS ?= 500
CEIL_MAX_VUS ?= 5000

.PHONY: help k6-smoke k6-load k6-ceiling k6-stats k6-all

help:
	@echo "Targets:"
	@echo "  make k6-smoke     BASE_URL=$(BASE_URL)"
	@echo "  make k6-load      RPS=$(RPS) RPS_STATS=$(RPS_STATS) DURATION=$(DURATION)"
	@echo "  make k6-ceiling   START_RPS=$(START_RPS) R1=$(R1) R2=$(R2) R3=$(R3)"
	@echo "  make k6-stats     FROM=$(FROM) TO=$(TO) RPS=.. DURATION=.."
	@echo ""
	@echo "Examples:"
	@echo "  make k6-smoke"
	@echo "  make k6-load RPS=500 RPS_STATS=30 DURATION=90s"
	@echo "  make k6-ceiling R1=2000 R2=4000 R3=6000"
	@echo "  make k6-stats FROM=2025-12-01 TO=2026-01-01 RPS=50 DURATION=60s"

k6-smoke:
	$(K6) run \
	  -e BASE_URL="$(BASE_URL)" \
	  -e DURATION=10s \
	  "$(K6_DIR)/smoke.js"

k6-load:
	$(K6) run \
	  -e BASE_URL="$(BASE_URL)" \
	  -e RPS="$(RPS)" \
	  -e RPS_STATS="$(RPS_STATS)" \
	  -e DURATION="$(DURATION)" \
	  -e PRE_VUS="$(PRE_VUS)" \
	  -e MAX_VUS="$(MAX_VUS)" \
	  -e PRE_VUS_STATS="$(PRE_VUS_STATS)" \
	  -e MAX_VUS_STATS="$(MAX_VUS_STATS)" \
	  "$(K6_DIR)/load_mix.js"

k6-ceiling:
	$(K6) run \
	  -e BASE_URL="$(BASE_URL)" \
	  -e START_RPS="$(START_RPS)" \
	  -e R1="$(R1)" \
	  -e R2="$(R2)" \
	  -e R3="$(R3)" \
	  -e PRE_VUS="$(CEIL_PRE_VUS)" \
	  -e MAX_VUS="$(CEIL_MAX_VUS)" \
	  "$(K6_DIR)/ceiling_one.js"

k6-stats:
	$(K6) run \
	  -e BASE_URL="$(BASE_URL)" \
	  -e FROM="$(FROM)" \
	  -e TO="$(TO)" \
	  -e RPS="$(RPS_STATS)" \
	  -e DURATION="$(DURATION)" \
	  -e PRE_VUS="$(PRE_VUS_STATS)" \
	  -e MAX_VUS="$(MAX_VUS_STATS)" \
	  "$(K6_DIR)/stats_routes_only.js"

k6-all: k6-smoke k6-load k6-stats