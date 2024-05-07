# Use the official MySQL image as the base image
FROM mysql:8.0

# Set environment variables for MySQL

# Copy a custom MySQL configuration file to the container
COPY mysql/conf.d/cfg.cnf /etc/mysql/conf.d/
# Expose port 3306 for MySQL
EXPOSE 3306

# 启动 MySQL 服务
CMD ["mysqld"]