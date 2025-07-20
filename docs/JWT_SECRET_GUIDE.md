# JWT Secret Generation Guide

Complete guide for generating secure JWT secrets for your Go application.

## 🔐 Security Requirements

**JWT secrets should be:**
- **Random**: Cryptographically secure random
- **Long**: At least 256 bits (32 bytes)
- **Unique**: Different for each environment
- **Secret**: Never commit to version control

## 🛠️ Generation Methods

### Method 1: Using Project Scripts (Recommended)

#### Generate New Secret
```bash
make generate-secret
```

#### Auto-Update .env File
```bash
make update-secret
```

#### Manual Generation
```bash
go run scripts/generate-jwt-secret.go
```

### Method 2: Using OpenSSL
```bash
openssl rand -base64 32
```

### Method 3: Using Python
```bash
python3 -c "import secrets; print(secrets.token_urlsafe(32))"
```

### Method 4: Using Node.js
```bash
node -e "console.log(require('crypto').randomBytes(32).toString('base64'))"
```

### Method 5: Online Generator (Development Only)
- Visit: https://www.allkeysgenerator.com/Random/Security-Encryption-Key-Generator.aspx
- Select: 256-bit
- **⚠️ Only use for development, never for production**

## 📁 Environment Configuration

### Development (.env)
```env
JWT_SECRET=your-generated-secret-here
```

### Production
Use environment variables or secret management:
```bash
export JWT_SECRET="your-production-secret"
```

### Docker
```yaml
version: '3.8'
services:
  app:
    environment:
      - JWT_SECRET=${JWT_SECRET}
```

## 🔄 Secret Rotation

### When to Rotate
- **Regular Schedule**: Every 90 days
- **Security Incident**: Immediately
- **Staff Changes**: When team members leave
- **Suspected Compromise**: Immediately

### How to Rotate
1. Generate new secret
2. Update environment variables
3. Restart application
4. Monitor for authentication issues

### Script for Rotation
```bash
# Generate and update
make update-secret

# Restart application
make run
```

## 🛡️ Security Best Practices

### Storage
- ✅ Use environment variables
- ✅ Use secret management services (AWS Secrets Manager, HashiCorp Vault)
- ❌ Never hardcode in source code
- ❌ Never commit to version control

### Access Control
- Limit who can access secrets
- Use different secrets per environment
- Log secret access (without values)

### Backup
- Store encrypted backups securely
- Have rotation procedures documented
- Test recovery procedures

## 🔍 Validation

### Check Secret Strength
```go
// scripts/validate-secret.go
package main

import (
    "encoding/base64"
    "fmt"
    "os"
)

func main() {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        fmt.Println("❌ JWT_SECRET not found")
        return
    }
    
    decoded, err := base64.URLEncoding.DecodeString(secret)
    if err != nil {
        fmt.Println("❌ Invalid base64 encoding")
        return
    }
    
    length := len(decoded)
    fmt.Printf("Secret length: %d bytes (%d bits)\n", length, length*8)
    
    if length >= 32 {
        fmt.Println("✅ Secret is sufficiently long")
    } else {
        fmt.Println("⚠️ Secret is too short, should be at least 32 bytes")
    }
}
```

### Test JWT Generation
```bash
# Start server and test login
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}'
```

## 📊 Different Environments

### Development
```env
JWT_SECRET=dev-secret-32-bytes-long-base64==
ENV=development
```

### Staging
```env
JWT_SECRET=staging-secret-different-from-dev=
ENV=staging
```

### Production
```env
JWT_SECRET=prod-secret-high-entropy-secure==
ENV=production
```

## 🚨 Emergency Procedures

### Compromised Secret
1. **Immediate Actions**:
   ```bash
   # Generate new secret
   make update-secret
   
   # Deploy immediately
   git add .env
   git commit -m "SECURITY: Rotate JWT secret"
   git push
   
   # Restart all services
   systemctl restart your-app
   ```

2. **User Impact**: All users will need to re-login

3. **Monitoring**: Watch for authentication errors

### Lost Secret
1. Check backups
2. Generate new secret
3. Update all environments
4. Notify team about re-login requirement

## 📝 Documentation Template

```markdown
## JWT Secret Management

**Current Secret**: Set in environment variable `JWT_SECRET`
**Last Rotated**: [DATE]
**Next Rotation**: [DATE + 90 days]
**Responsible**: [TEAM/PERSON]

### Emergency Contacts
- Security Team: security@company.com
- DevOps Team: devops@company.com

### Rotation Schedule
- Development: Monthly
- Staging: Quarterly  
- Production: Quarterly
```

## 🎯 Quick Commands Summary

```bash
# Generate new secret (view only)
make generate-secret

# Update .env file automatically
make update-secret

# Validate current secret
go run scripts/validate-secret.go

# Test authentication
curl -X POST localhost:3000/api/auth/login -H "Content-Type: application/json" -d '{"email":"test@example.com","password":"password"}'
```

Remember: **Security is only as strong as your weakest secret!** 🔒
