#!/bin/bash

# Database Verification Script
# Verifikasi hasil seeding database

echo "🔍 Cinema Booking API - Database Verification"
echo "============================================="

# Check if .env file exists
if [ ! -f .env ]; then
    echo "❌ Error: .env file not found!"
    exit 1
fi

echo "📊 Connecting to database and checking collections..."

# Create temporary verification script
cat > scripts/verify_data.js << 'EOF'
// MongoDB verification script
const collections = ['users', 'films', 'jadwals', 'tikets', 'pembayarans'];

collections.forEach(collection => {
    const count = db[collection].countDocuments();
    print(`📋 ${collection}: ${count} documents`);
});

print("\n👤 Sample Users:");
db.users.find({}, {username: 1, email: 1, role: 1}).forEach(user => {
    print(`   - ${user.username} (${user.email}) [${user.role}]`);
});

print("\n🎬 Sample Films:");
db.films.find({}, {title: 1, genre: 1, duration: 1}).limit(3).forEach(film => {
    print(`   - ${film.title} (${film.genre}, ${film.duration}min)`);
});

print("\n📅 Sample Jadwals:");
db.jadwals.aggregate([
    {
        $lookup: {
            from: "films",
            localField: "film_id",
            foreignField: "_id",
            as: "film"
        }
    },
    { $limit: 3 },
    {
        $project: {
            "film.title": 1,
            "tanggal": 1,
            "waktu": 1,
            "ruangan": 1,
            "harga": 1
        }
    }
]).forEach(jadwal => {
    const filmTitle = jadwal.film[0] ? jadwal.film[0].title : "Unknown Film";
    print(`   - ${filmTitle} | ${jadwal.tanggal} ${jadwal.waktu} | ${jadwal.ruangan} | Rp ${jadwal.harga}`);
});

print("\n🎫 Sample Tikets:");
db.tikets.aggregate([
    {
        $lookup: {
            from: "users",
            localField: "user_id", 
            foreignField: "_id",
            as: "user"
        }
    },
    { $limit: 3 },
    {
        $project: {
            "user.username": 1,
            "kursi": 1,
            "status": 1,
            "tanggal_pembelian": 1
        }
    }
]).forEach(tiket => {
    const username = tiket.user[0] ? tiket.user[0].username : "Unknown User";
    print(`   - ${username} | Kursi ${tiket.kursi} | ${tiket.status}`);
});

print("\n💳 Sample Pembayarans:");
db.pembayarans.find({}, {jumlah: 1, metode_pembayaran: 1, status: 1}).limit(3).forEach(payment => {
    print(`   - Rp ${payment.jumlah} via ${payment.metode_pembayaran} [${payment.status}]`);
});

print("\n✅ Database verification completed!");
EOF

echo "🔗 Extracting database connection info..."
MONGOSTRING=$(grep MONGOSTRING .env | cut -d '=' -f2)
DB_NAME="cinema_booking"

if [ -z "$MONGOSTRING" ]; then
    echo "❌ Error: MONGOSTRING not found in .env"
    exit 1
fi

# Run verification
echo "🚀 Running verification script..."
mongosh "${MONGOSTRING}${DB_NAME}" --quiet scripts/verify_data.js

# Cleanup
rm -f scripts/verify_data.js

echo ""
echo "✅ Database verification completed!"
echo ""
echo "🔐 Test these credentials in your API:"
echo "   Admin: admin@cinema.com / admin123"
echo "   User: john@example.com / password123"
