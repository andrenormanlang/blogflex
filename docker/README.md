# Using Neon DB with BlogFlex

This guide explains how to configure BlogFlex to use Neon DB instead of a local PostgreSQL database.

## Setup Instructions

1. **Create a Neon DB Account and Database**
   - Sign up at [Neon](https://neon.tech/)
   - Create a new project and PostgreSQL database
   - Note down your connection details

2. **Update Your .env File**
   Add the following environment variables to your `.env` file:

   ```
   # Neon DB configuration
   NEON_DB_HOST=your-db-host.neon.tech
   NEON_DB_USER=your-db-user
   NEON_DB_PASSWORD=your-db-password
   NEON_DB_NAME=your-db-name
   NEON_DB_PORT=5432
   ```

3. **Update Hasura to Use Neon DB**
   - Go to your Hasura console
   - Navigate to Data → Manage → Database → Connect Database
   - Choose "Connect existing database"
   - Enter your Neon DB connection details
   - Make sure to check "Use SSL" option
   - Click "Connect Database"

4. **Run the Application with Docker**
   ```bash
   cd docker
   docker-compose up -d
   ```

## Troubleshooting

- If you encounter connection issues, make sure your Neon DB allows connections from your IP address
- Check that SSL is enabled for the connection
- Verify that your environment variables are correctly set in the .env file

## Benefits of Using Neon DB

- Serverless PostgreSQL database
- Automatic scaling
- Built-in branching for development and testing
- No need to manage local database instances
- Data persistence across deployments 