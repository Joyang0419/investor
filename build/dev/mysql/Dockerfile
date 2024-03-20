# Use the official MySQL image as the base image
FROM mysql:latest

# Set environment variables for MySQL
ENV MYSQL_ROOT_PASSWORD=root
ENV MYSQL_DATABASE=investor

# Copy a custom MySQL configuration file to the container
COPY mysql/conf.d/cfg.cnf /etc/mysql/conf.d/
# Expose port 3306 for MySQL
EXPOSE 3306