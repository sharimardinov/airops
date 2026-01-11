#!/bin/bash
set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}=== AirOps Test Suite ===${NC}\n"

if ! command -v go &> /dev/null; then
  echo -e "${RED}Go is not installed${NC}"
  exit 1
fi

echo -e "${YELLOW}Go version:${NC}"
go version
echo ""

PASSED=0
FAILED=0

# coverage state
COVERAGE_INITIALIZED=0
rm -f coverage.out coverage.html coverage.tmp.out

append_coverage() {
  local file="$1"
  # если файл пустой/не создан — выходим
  [ -s "$file" ] || return 0

  if [ $COVERAGE_INITIALIZED -eq 0 ]; then
    # первый файл пишем целиком (с mode: set)
    cat "$file" > coverage.out
    COVERAGE_INITIALIZED=1
  else
    # у следующих выкидываем первую строку "mode: ..."
    tail -n +2 "$file" >> coverage.out
  fi
}

run_tests() {
  local package=$1
  local name=$2

  echo -e "${YELLOW}Running ${name}...${NC}"

  # Сохраняем полный вывод, но не ломаем код возврата
  set +e
  OUTPUT=$(go test "$package" -v -coverprofile=coverage.tmp.out 2>&1)
  STATUS=$?
  set -e

  # печатаем всё кроме "no test files" (чтоб не спамило)
  echo "$OUTPUT" | grep -v "no test files" || true

  if [ $STATUS -ne 0 ]; then
    echo -e "${RED}✗ ${name} failed${NC}\n"
    return 1
  fi

  append_coverage "coverage.tmp.out"
  rm -f coverage.tmp.out

  echo -e "${GREEN}✓ ${name} passed${NC}\n"
  return 0
}

# 1. gofmt
echo -e "${YELLOW}Checking code formatting...${NC}"
UNFORMATTED=$(gofmt -l . || true)
if [ -n "$UNFORMATTED" ]; then
  echo -e "${RED}Some files are not formatted:${NC}"
  echo "$UNFORMATTED"
  echo -e "${YELLOW}Run: gofmt -w .${NC}\n"
  ((FAILED++))
else
  echo -e "${GREEN}✓ All files are formatted${NC}\n"
  ((PASSED++))
fi

# 2. go vet
echo -e "${YELLOW}Running go vet...${NC}"
if go vet ./...; then
  echo -e "${GREEN}✓ go vet passed${NC}\n"
  ((PASSED++))
else
  echo -e "${RED}✗ go vet failed${NC}\n"
  ((FAILED++))
fi

# 3..6 tests (пакеты под твой проект — см. ниже)
if run_tests "./internal/domain/..." "Domain tests"; then ((PASSED++)); else ((FAILED++)); fi
if run_tests "./internal/app/usecase/..." "Use case tests"; then ((PASSED++)); else ((FAILED++)); fi

# ВАЖНО: эти пути у тебя не существуют — поэтому тесты падали.
# Подставь реальные директории:
# пример: ./internal/postgres/... или ./internal/repository/... или ./internal/transport/http/...
if run_tests "./internal/infrastructure/postgres/..." "Postgres tests"; then ((PASSED++)); else ((FAILED++)); fi
if run_tests "./internal/transport/http/..." "HTTP tests"; then ((PASSED++)); else ((FAILED++)); fi

# 7. integration (опционально)
if [ "${RUN_INTEGRATION_TESTS:-0}" = "1" ]; then
  echo -e "${YELLOW}Running integration tests...${NC}"

  export TEST_DATABASE_URL="${TEST_DATABASE_URL:-postgres://airops:secret@localhost:5432/airops_test?sslmode=disable}"

  echo -e "${YELLOW}Testing database connection...${NC}"
  if command -v psql >/dev/null 2>&1 && psql "$TEST_DATABASE_URL" -c "SELECT 1" >/dev/null 2>&1; then
    echo -e "${GREEN}✓ Database connection successful${NC}\n"
    if run_tests "./tests/integration/..." "Integration tests"; then ((PASSED++)); else ((FAILED++)); fi
  else
    echo -e "${RED}✗ Cannot connect to database (or psql not installed)${NC}\n"
    ((FAILED++))
  fi
else
  echo -e "${YELLOW}Skipping integration tests (set RUN_INTEGRATION_TESTS=1 to run)${NC}\n"
fi

# 8. Build
echo -e "${YELLOW}Testing build...${NC}"
if go build -o /tmp/airops ./cmd/api/main.go; then
  echo -e "${GREEN}✓ Build successful${NC}\n"
  rm -f /tmp/airops
  ((PASSED++))
else
  echo -e "${RED}✗ Build failed${NC}\n"
  ((FAILED++))
fi

# coverage report
if [ -s coverage.out ]; then
  echo -e "${YELLOW}Generating coverage report...${NC}"
  COVERAGE=$(go tool cover -func=coverage.out | awk '/total:/ {print $3}')
  echo -e "${GREEN}Total coverage: ${COVERAGE}${NC}\n"
  go tool cover -html=coverage.out -o coverage.html
  echo -e "${GREEN}HTML coverage report: coverage.html${NC}\n"
else
  echo -e "${YELLOW}No coverage data collected (likely no tests or no packages with coverprofile).${NC}\n"
fi

echo -e "${GREEN}=== Test Summary ===${NC}"
echo -e "Passed: ${GREEN}${PASSED}${NC}"
echo -e "Failed: ${RED}${FAILED}${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
  echo -e "${GREEN}All tests passed! 🎉${NC}"
  exit 0
else
  echo -e "${RED}Some tests failed${NC}"
  exit 1
fi