#!/bin/bash

# Database Clear Script
# Menghapus semua data dari database

echo "🗑️  Cinema Booking API - Database Clear"
echo "======================================="

echo "⚠️  WARNING: This will DELETE ALL DATA from the database!"
echo "Are you sure you want to continue? (y/N)"
read -r response

if [[ ! "$response" =~ ^[Yy]$ ]]; then
    echo "❌ Operation cancelled."
    exit 0
fi

# Check if .env file exists
if [ ! -f .env ]; then
    echo "❌ Error: .env file not found!"
    exit 1
fi

echo "🔗 Extracting database connection info..."
MONGOSTRING=$(grep MONGOSTRING .env | cut -d '=' -f2)
DB_NAME="cinema_booking"

if [ -z "$MONGOSTRING" ]; then
    echo "❌ Error: MONGOSTRING not found in .env"
    exit 1
fi

# Create clear script
cat > scripts/clear_data.js << 'EOF'
// MongoDB clear script
const collections = ['pembayarans', 'tikets', 'jadwals', 'films', 'users'];

print("🗑️  Clearing collections in dependency order...");

collections.forEach(collection => {
    const result = db[collection].deleteMany({});
    print(`✅ Cleared ${collection}: ${result.deletedCount} documents deleted`);
});

print("\n✅ All collections cleared successfully!");
EOF

echo "🚀 Clearing database..."
mongosh "${MONGOSTRING}${DB_NAME}" --quiet scripts/clear_data.js

# Cleanup
rm -f scripts/clear_data.js

echo ""
echo "✅ Database clear completed!"
echo ""
echo "💡 To repopulate with sample data, run:"
echo "   ./scripts/run-seeder.sh"
