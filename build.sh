#!/bin/bash
set -e

echo "Starting build process..."
echo "Current directory: $(pwd)"
echo "Listing contents:"
ls -la

echo "Changing to frontend directory..."
cd frontend

echo "Frontend directory contents:"
ls -la

echo "Installing dependencies..."
npm ci

echo "Building project..."
npm run build

echo "Build completed successfully!"
ls -la dist/
