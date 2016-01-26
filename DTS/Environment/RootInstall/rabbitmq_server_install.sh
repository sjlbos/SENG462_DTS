# This script installs a RabbitMQ broker service in a unix/linux environment

echo "Adding RabbitMQ signing key..."
apt-key adv --keyserver pgp.mit.edu --recv-keys 0x056E8E56

echo "Adding RabbitMQ package repository..."
echo "deb http://www.rabbitmq.com/debian/ testing main" | tee /etc/apt/sources.list.d/rabbitmq.list

echo "Updating package list..."
apt-get update

echo "Installing package \"rabbitmq-server\""
apt-get -y install rabbitmq-server