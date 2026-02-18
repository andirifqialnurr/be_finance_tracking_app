#!/usr/bin/env pwsh
# Test Runner Script for Finance Tracking API

param(
    [string]$TestType = "all"
)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Finance Tracking API - Test Suite" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

function Run-Tests {
    param(
        [string]$Package,
        [string]$Name
    )
    
    Write-Host "Running $Name..." -ForegroundColor Yellow
    go test -v $Package -count=1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ $Name passed" -ForegroundColor Green
    } else {
        Write-Host "✗ $Name failed" -ForegroundColor Red
    }
    Write-Host ""
}

function Run-TestsWithCoverage {
    param(
        [string]$Package,
        [string]$Name
    )
    
    Write-Host "Running $Name with coverage..." -ForegroundColor Yellow
    go test -v -cover -coverprofile=coverage.out $Package
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ $Name passed" -ForegroundColor Green
        go tool cover -func=coverage.out | Select-String "total:"
    } else {
        Write-Host "✗ $Name failed" -ForegroundColor Red
    }
    Write-Host ""
}

# Install dependencies
Write-Host "Installing test dependencies..." -ForegroundColor Cyan
go get github.com/stretchr/testify/assert
go get gorm.io/driver/postgres
Write-Host ""

switch ($TestType) {
    "unit" {
        Write-Host "Running Unit Tests Only" -ForegroundColor Cyan
        Write-Host ""
        Run-Tests "./test/repositories/..." "Repository Tests"
        Run-Tests "./test/services/..." "Service Tests"
    }
    
    "integration" {
        Write-Host "Running Integration Tests Only" -ForegroundColor Cyan
        Write-Host ""
        Run-Tests "./test/integration/..." "Integration Tests"
    }
    
    "handlers" {
        Write-Host "Running Handler/Endpoint Tests Only" -ForegroundColor Cyan
        Write-Host ""
        Run-Tests "./test/handlers/..." "Handler Tests"
    }
    
    "architecture" {
        Write-Host "Running Architecture Tests Only" -ForegroundColor Cyan
        Write-Host ""
        Run-Tests "./test/architecture/..." "Architecture Tests"
    }
    
    "coverage" {
        Write-Host "Running All Tests with Coverage" -ForegroundColor Cyan
        Write-Host ""
        Run-TestsWithCoverage "./test/..." "All Tests"
        
        Write-Host "Generating HTML coverage report..." -ForegroundColor Cyan
        go tool cover -html=coverage.out -o coverage.html
        Write-Host "Coverage report saved to coverage.html" -ForegroundColor Green
    }
    
    default {
        Write-Host "Running All Tests" -ForegroundColor Cyan
        Write-Host ""
        
        Run-Tests "./test/repositories/..." "Repository Tests"
        Run-Tests "./test/services/..." "Service Tests"
        Run-Tests "./test/handlers/..." "Handler Tests"
        Run-Tests "./test/integration/..." "Integration Tests"
        Run-Tests "./test/architecture/..." "Architecture Tests"
        
        Write-Host "========================================" -ForegroundColor Cyan
        Write-Host "Test Summary" -ForegroundColor Cyan
        Write-Host "========================================" -ForegroundColor Cyan
        Write-Host "All test suites completed!" -ForegroundColor Green
    }
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Test run completed" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
