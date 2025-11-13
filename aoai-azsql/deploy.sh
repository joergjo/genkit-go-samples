#!/bin/bash
set -e

# Generate a URL-safe password (base64 encoding produces +, /, = which need to be replaced)
# Use alphanumeric characters plus hyphen and underscore which are URL-safe
password=$(openssl rand -base64 30 | tr -d '+/=' | tr -d '\n' | head -c 40)

# Configuration variables
resource_group="${AZ_RESOURCE_GROUP:-genkit-azsql-demo}"
location="${AZ_LOCATION:-northeurope}"
server_name="${AZ_SERVER_NAME:-genkit-demo-$RANDOM}"
database_name="${AZ_DATABASE_NAME:-genkit-demo}"
admin_user="${AZ_ADMIN_USER:-sqladmin}"
admin_password="${AZ_ADMIN_PASSWORD:-$password}"

echo "Creating Azure SQL Database resources..."

# Create resource group
echo "Creating resource group: $resource_group"
az group create \
    --name "$resource_group" \
    --location "$location"

# Create SQL server
echo "Creating SQL server: $server_name"
az sql server create \
    --name "$server_name" \
    --resource-group "$resource_group" \
    --location "$location" \
    --admin-user "$admin_user" \
    --admin-password "$admin_password"

# Configure firewall - allow Azure services
echo "Configuring firewall rules..."
az sql server firewall-rule create \
    --resource-group "$resource_group" \
    --server "$server_name" \
    --name "AllowAzureServices" \
    --start-ip-address 0.0.0.0 \
    --end-ip-address 0.0.0.0

# Get your current public IP and add firewall rule
my_ip=$(curl -s https://api.ipify.org)
az sql server firewall-rule create \
    --resource-group "$resource_group" \
    --server "$server_name" \
    --name "AllowMyIP" \
    --start-ip-address "$my_ip" \
    --end-ip-address "$my_ip"

# Create SQL database
echo "Creating SQL database: $database_name"
az sql db create \
    --resource-group "$resource_group" \
    --server "$server_name" \
    --name "$database_name" \
    --edition GeneralPurpose \
    --compute-model Serverless \
    --family Gen5 \
    --capacity 2

# Wait a moment for database to be ready
sleep 10

# Apply SQL script using sqlcmd
echo "Applying vector.sql script..."
if command -v sqlcmd &> /dev/null; then
    sqlcmd -S "$server_name.database.windows.net" \
        -d "$database_name" \
        -U "$admin_user" \
        -P "$admin_password" \
        -i vector.sql
    echo "Database schema applied successfully!"
else
    echo "Warning: sqlcmd not found. Please install SQL Server command-line tools."
    echo "You can manually apply vector.sql using:"
    echo "sqlcmd -S $server_name.database.windows.net -d $database_name -U $admin_user -P $admin_password -i vector.sql"
fi

echo ""
echo "Deployment complete!"
echo "Server: $server_name.database.windows.net"
echo "Database: $database_name"
echo "Admin User: $admin_user"
echo "Admin Password: $admin_password"
echo ""
echo "Go connection string:"
echo "sqlserver://$admin_user:$admin_password@$server_name.database.windows.net?database=$database_name"
# echo "Server=tcp:$server_name.database.windows.net,1433;Database=$database_name;User ID=$admin_user;Password=$admin_password;Encrypt=true;Connection Timeout=30;"
