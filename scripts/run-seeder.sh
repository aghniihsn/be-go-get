#!/bin/bash

# Database Seeder Script
# Menjalankan seeder untuk mengisi database dengan data sample

echo "🌱 Cinema Booking API - Database Seeder"
echo "========================================="

# Check if .env file exists
if [ ! -f .env ]; then
    echo "❌ Error: .env file not found!"
    echo "Please create .env file with your database configuration."
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed!"
    echo "Please install Go to run the seeder."
    exit 1
fi

echo "📝 Checking Go modules..."
go mod tidy

echo "🚀 Running database seeder..."
cd scripts/seeder
go run main.go

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Database seeding completed successfully!"
    echo ""
    echo "📋 Seeded Data Summary:"
    echo "   👤 Users: 5 (1 admin + 4 regular users)"
    echo "   🎬 Films: 6 popular movies"
    echo "   📅 Jadwals: 18 schedules (3 per film)"
    echo "   🎫 Tikets: 8 sample tickets"
    echo "   💳 Pembayarans: 6 payment records"
    echo ""
    echo "🔐 Test Credentials:"
    echo "   Admin: admin@cinema.com / admin123"
    echo "   User: john@example.com / password123"
    echo ""
    echo "🚀 You can now start testing the API!"
else
    echo "❌ Database seeding failed!"
    exit 1
fi
